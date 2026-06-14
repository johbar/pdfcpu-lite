/*
Copyright 2018 The pdfcpu Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package filter

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"

	"github.com/johbar/pdfcpu-lite/pkg/log"
)

// Portions of this code are based on ideas of image/png: reader.go:readImagePass
// PNG is documented here: www.w3.org/TR/PNG-Filters.html

// PDF allows a prediction step prior to compression applying TIFF or PNG prediction.
// Predictor algorithm.
const (
	PredictorNo      = 1  // No prediction.
	PredictorTIFF    = 2  // Use TIFF prediction for all rows.
	PredictorNone    = 10 // Use PNGNone for all rows.
	PredictorSub     = 11 // Use PNGSub for all rows.
	PredictorUp      = 12 // Use PNGUp for all rows.
	PredictorAverage = 13 // Use PNGAverage for all rows.
	PredictorPaeth   = 14 // Use PNGPaeth for all rows.
	PredictorOptimum = 15 // Use the optimum PNG prediction for each row.
)

// For predictor > 2 PNG filters (see RFC 2083) get applied and the first byte of each pixelrow defines
// the prediction algorithm used for all pixels of this row.
const (
	PNGNone    = 0x00
	PNGSub     = 0x01
	PNGUp      = 0x02
	PNGAverage = 0x03
	PNGPaeth   = 0x04
)

// rowBufPool recycles the per-row scratch slices used by decodePostProcess.
// Each pool entry is a *[]byte so that the pointer itself can be stored and
// retrieved without escaping the slice header to the heap on every get/put.
var rowBufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 256) // reasonable starting capacity for typical rows
		return &b
	},
}

// getRowBuf retrieves a zeroed []byte of exactly length n from the pool,
// growing the backing array only when necessary.
func getRowBuf(n int) *[]byte {
	bp := rowBufPool.Get().(*[]byte)
	if cap(*bp) >= n {
		*bp = (*bp)[:n]
	} else {
		*bp = make([]byte, n)
	}
	clear(*bp) // zero so prior content never bleeds into a new stream
	return bp
}

// putRowBuf returns a row buffer to the pool. The caller must not use the
// slice after this call.
func putRowBuf(bp *[]byte) {
	rowBufPool.Put(bp)
}

// bufPool recycles bytes.Buffer objects used to collect decoded output.
// Buffers are Reset() before reuse so no prior data leaks out.
var bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

// pooledBuffer wraps a *bytes.Buffer that came from bufPool and returns it
// to the pool when Close is called.  It exposes io.ReadCloser so callers
// that only need io.Reader still work (io.ReadCloser satisfies io.Reader).
type pooledBuffer struct {
	*bytes.Buffer
}

// Close returns the underlying buffer to the pool. The caller must not read
// from this pooledBuffer after Close returns.
func (pb *pooledBuffer) Close() error {
	pb.Buffer.Reset()
	bufPool.Put(pb.Buffer)
	pb.Buffer = nil
	return nil
}

// validBPC lists the BitsPerComponent values permitted by the PDF spec.
// Package-level to avoid a per-call heap allocation in parameters().
var validBPC = []int{1, 2, 4, 8, 16}

// zlibReaderPool recycles zlib.Reader instances across DecodeLength calls.
//
// zlib.NewReader internally calls flate.NewReader, which allocates a
// decompressor carrying a ~32 KB sliding-window buffer plus Huffman decode
// tables — by far the largest per-stream allocation in this pipeline.
//
// Both zlib.Reader and its embedded flate.Reader implement the Resetter
// interface, so Reset(r, nil) re-reads the 2-byte zlib header from the new
// reader and redirects the entire decompressor stack in-place with no new
// heap allocation.
//
// Critical rule: do NOT call Close() before returning a reader to this pool.
// Close() tears down the internal flate.Reader; a closed reader cannot be
// Reset. Reinitialization via Reset is sufficient and replaces Close entirely.
var zlibReaderPool sync.Pool

// getZlibReader returns a zlib reader positioned at the start of r.
// It prefers a pooled instance (reset via zlib.Resetter) over allocating a
// new one. If Reset reports an invalid zlib header the pooled entry is
// discarded and zlib.NewReader is used as fallback — keeping the fast path
// allocation-free while still handling corrupt-header edge cases safely.
func getZlibReader(r io.Reader) (io.ReadCloser, error) {
	if v := zlibReaderPool.Get(); v != nil {
		zr := v.(io.ReadCloser)
		if err := zr.(zlib.Resetter).Reset(r, nil); err == nil {
			return zr, nil
		}
		// Reset failed (e.g. pooled instance left in an unrecoverable state).
		// Fall through to allocate a fresh reader; the discarded entry is GC'd.
	}
	return zlib.NewReader(r)
}

// putZlibReader returns zr to the pool without closing it.
// The caller must not use zr after this call.
func putZlibReader(zr io.ReadCloser) {
	zlibReaderPool.Put(zr)
}

// validPredictors is a package-level slice so that decodePostProcess does not
// allocate a fresh []int literal on every invocation.
var validPredictors = []int{
	PredictorTIFF,
	PredictorNone,
	PredictorSub,
	PredictorUp,
	PredictorAverage,
	PredictorPaeth,
	PredictorOptimum,
}

type flate struct {
	baseFilter
}

// Encode implements encoding for a Flate filter.
func (f flate) Encode(r io.Reader) (io.Reader, error) {
	if log.TraceEnabled() {
		log.Trace.Println("EncodeFlate begin")
	}

	// TODO Optional decode parameters may need predictor preprocessing.

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	defer w.Close()

	written, err := io.Copy(w, r)
	if err != nil {
		return nil, err
	}

	if log.TraceEnabled() {
		log.Trace.Printf("EncodeFlate end: %d bytes written\n", written)
	}

	return &b, nil
}

// Decode implements decoding for a Flate filter.
func (f flate) Decode(r io.Reader) (io.Reader, error) {
	return f.DecodeLength(r, -1)
}

func (f flate) DecodeLength(r io.Reader, maxLen int64) (io.Reader, error) {
	if log.TraceEnabled() {
		log.Trace.Println("DecodeFlate begin")
	}

	// Obtain a reusable zlib reader from the pool. Reset re-reads the
	// 2-byte zlib header from r and reinitializes the flate decompressor
	// in-place, avoiding the ~32 KB window-buffer allocation of NewReader.
	rc, err := getZlibReader(r)
	if err != nil {
		return nil, err
	}
	// Return rc to the pool instead of closing it.  Close would free the
	// internal flate decompressor; we want Reset on the next call to reuse it.
	defer putZlibReader(rc)

	// Optional decode parameters need postprocessing.
	return f.decodePostProcess(rc, maxLen)
}

// passThru copies rin (up to maxLen bytes if ≥ 0) into a pooled buffer and
// returns it as an io.ReadCloser. Callers that previously received a
// *bytes.Buffer continue to work because *bytes.Buffer implements io.Reader;
// callers that can call Close() will return the buffer to the pool.
func passThru(rin io.Reader, maxLen int64) (*pooledBuffer, error) {
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	var err error
	if maxLen < 0 {
		_, err = io.Copy(b, rin)
	} else {
		_, err = io.CopyN(b, rin, maxLen)
	}
	if err != nil && strings.Contains(err.Error(), "invalid checksum") {
		if log.CLIEnabled() {
			log.CLI.Println("skipped: truncated zlib stream")
		}
		err = nil
	}
	if err == io.ErrUnexpectedEOF {
		// Workaround for missing support for partial flush in compress/flate.
		// See also https://github.com/golang/go/issues/31514
		if log.ReadEnabled() {
			log.Read.Println("flateDecode: ignoring unexpected EOF")
		}
		err = nil
	}
	if err != nil {
		b.Reset()
		bufPool.Put(b)
		return nil, err
	}
	return &pooledBuffer{b}, nil
}

func intMemberOf(i int, list []int) bool {
	return slices.Contains(list, i)
}

// Each prediction value implies (a) certain row filter(s).
// func validateRowFilter(f, p int) error {

// 	switch p {

// 	case PredictorNone:
// 		if !intMemberOf(f, []int{PNGNone, PNGSub, PNGUp, PNGAverage, PNGPaeth}) {
// 			return errors.Errorf("pdfcpu: validateRowFilter: PredictorOptimum, unexpected row filter #%02x", f)
// 		}
// 		// if f != PNGNone {
// 		// 	return errors.Errorf("validateRowFilter: expected row filter #%02x, got: #%02x", PNGNone, f)
// 		// }

// 	case PredictorSub:
// 		if f != PNGSub {
// 			return errors.Errorf("pdfcpu: validateRowFilter: expected row filter #%02x, got: #%02x", PNGSub, f)
// 		}

// 	case PredictorUp:
// 		if f != PNGUp {
// 			return errors.Errorf("pdfcpu: validateRowFilter: expected row filter #%02x, got: #%02x", PNGUp, f)
// 		}

// 	case PredictorAverage:
// 		if f != PNGAverage {
// 			return errors.Errorf("pdfcpu: validateRowFilter: expected row filter #%02x, got: #%02x", PNGAverage, f)
// 		}

// 	case PredictorPaeth:
// 		if f != PNGPaeth {
// 			return errors.Errorf("pdfcpu: validateRowFilter: expected row filter #%02x, got: #%02x", PNGPaeth, f)
// 		}

// 	case PredictorOptimum:
// 		if !intMemberOf(f, []int{PNGNone, PNGSub, PNGUp, PNGAverage, PNGPaeth}) {
// 			return errors.Errorf("pdfcpu: validateRowFilter: PredictorOptimum, unexpected row filter #%02x", f)
// 		}

// 	default:
// 		return errors.Errorf("pdfcpu: validateRowFilter: unexpected predictor #%02x", p)

// 	}

// 	return nil
// }

func applyHorDiff(row []byte, colors int) ([]byte, error) {
	// This works for 8 bits per color only.
	for i := 1; i < len(row)/colors; i++ {
		for j := range colors {
			row[i*colors+j] += row[(i-1)*colors+j]
		}
	}
	return row, nil
}

func processRow(pr, cr []byte, p, colors, bytesPerPixel int) ([]byte, error) {
	//fmt.Printf("pr(%v) =\n%s\n", &pr, hex.Dump(pr))
	//fmt.Printf("cr(%v) =\n%s\n", &cr, hex.Dump(cr))

	if p == PredictorTIFF {
		return applyHorDiff(cr, colors)
	}

	// Apply the filter.
	cdat := cr[1:]
	pdat := pr[1:]

	// Get row filter from 1st byte
	f := int(cr[0])

	// The value of Predictor supplied by the decoding filter need not match the value
	// used when the data was encoded if they are both greater than or equal to 10.

	switch f {

	case PNGNone:
		// No operation.

	case PNGSub:
		for i := bytesPerPixel; i < len(cdat); i++ {
			cdat[i] += cdat[i-bytesPerPixel]
		}

	case PNGUp:
		for i, p := range pdat {
			cdat[i] += p
		}

	case PNGAverage:
		// The average of the two neighboring pixels (left and above).
		// Raw(x) - floor((Raw(x-bpp)+Prior(x))/2)
		for i := range bytesPerPixel {
			cdat[i] += pdat[i] / 2
		}
		for i := bytesPerPixel; i < len(cdat); i++ {
			cdat[i] += uint8((int(cdat[i-bytesPerPixel]) + int(pdat[i])) / 2)
		}

	case PNGPaeth:
		filterPaeth(cdat, pdat, bytesPerPixel)

	}

	return cdat, nil
}

func (f flate) parameters() (colors, bpc, columns int, err error) {
	// Colors, int
	// The number of interleaved colour components per sample.
	// Valid values are 1 to 4 (PDF 1.0) and 1 or greater (PDF 1.3). Default value: 1.
	// Used by PredictorTIFF only.
	colors, found := f.parms["Colors"]
	if !found {
		colors = 1
	} else if colors == 0 {
		return 0, 0, 0, fmt.Errorf("pdfcpu: filter FlateDecode: \"Colors\" must be > 0")
	}

	// BitsPerComponent, int
	// The number of bits used to represent each colour component in a sample.
	// Valid values are 1, 2, 4, 8, and (PDF 1.5) 16. Default value: 8.
	// Used by PredictorTIFF only.
	bpc, found = f.parms["BitsPerComponent"]
	if !found {
		bpc = 8
	} else if !intMemberOf(bpc, validBPC) {
		return 0, 0, 0, fmt.Errorf("pdfcpu: filter FlateDecode: Unexpected \"BitsPerComponent\": %d", bpc)
	}

	// Columns, int
	// The number of samples in each row. Default value: 1.
	columns, found = f.parms["Columns"]
	if !found {
		columns = 1
	}

	return colors, bpc, columns, nil
}

func checkBufLen(b bytes.Buffer, maxLen int64) bool {
	return maxLen < 0 || int64(b.Len()) < maxLen
}

func process(w io.Writer, pr, cr []byte, predictor, colors, bytesPerPixel int) error {
	d, err := processRow(pr, cr, predictor, colors, bytesPerPixel)
	if err != nil {
		return err
	}

	_, err = w.Write(d)

	return err
}

// decodePostProcess
func (f flate) decodePostProcess(r io.Reader, maxLen int64) (io.Reader, error) {
	predictor, found := f.parms["Predictor"]
	if !found || predictor == PredictorNo {
		return passThru(r, maxLen)
	}

	// Use the package-level slice to avoid a per-call []int allocation.
	if !intMemberOf(predictor, validPredictors) {
		return nil, fmt.Errorf("pdfcpu: filter FlateDecode: undefined \"Predictor\" %d", predictor)
	}

	colors, bpc, columns, err := f.parameters()
	if err != nil {
		return nil, err
	}

	bytesPerPixel := (bpc*colors + 7) / 8
	rowSize := (bpc*colors*columns + 7) / 8

	m := rowSize
	if predictor != PredictorTIFF {
		// PNG prediction uses a row filter byte prefixing the pixelbytes of a row.
		m++
	}

	// cr and pr are the bytes for the current and previous row.
	// Both are obtained from rowBufPool so their backing arrays are reused
	// across calls.  We keep the original *[]byte pointers in crPtr/prPtr so
	// that the end-of-loop swap (pr, cr = cr, pr) on the slice headers does
	// not confuse the pool's Put calls in the defers.
	crPtr := getRowBuf(m)
	prPtr := getRowBuf(m)
	defer putRowBuf(crPtr)
	defer putRowBuf(prPtr)

	cr := *crPtr
	pr := *prPtr

	// Output buffer obtained from the pool; pre-grown with a rough capacity
	// estimate (number of rows × rowSize) to minimise internal copies.
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	if maxLen > 0 {
		b.Grow(int(maxLen)) // caller already knows the expected byte count
	} else if rowSize > 0 {
		// A modest initial reservation; avoids the first few doublings.
		b.Grow(rowSize * 64)
	}

	for checkBufLen(*b, maxLen) {

		// Read decompressed bytes for one pixel row.
		n, err := io.ReadFull(r, cr)
		if err != nil {
			if err != io.EOF {
				b.Reset()
				bufPool.Put(b)
				return nil, err
			}
			// eof
			if n == 0 {
				break
			}
		}

		if n != m {
			b.Reset()
			bufPool.Put(b)
			return nil, fmt.Errorf("pdfcpu: filter FlateDecode: read error, expected %d bytes, got: %d", m, n)
		}

		if err := process(b, pr, cr, predictor, colors, bytesPerPixel); err != nil {
			b.Reset()
			bufPool.Put(b)
			return nil, err
		}

		if err == io.EOF {
			break
		}

		// Swap slice headers; the pool pointers (crPtr/prPtr) are unaffected.
		pr, cr = cr, pr
	}

	if maxLen < 0 && b.Len()%rowSize > 0 {
		log.Info.Printf("failed postprocessing: %d %d\n", b.Len(), rowSize)
		b.Reset()
		bufPool.Put(b)
		return nil, errors.New("pdfcpu: filter FlateDecode: postprocessing failed")
	}

	return &pooledBuffer{b}, nil
}
