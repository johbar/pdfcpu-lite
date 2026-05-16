/*
Copyright 2026 The pdfcpu-lite Authors.

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
package model

import (
	"fmt"
	"io"

	"github.com/johbar/pdfcpu-lite/pkg/pdfcpu/types"
)

// PageContentInto appends the content in PDF syntax for page dict d to buf.
// This enables (external) memory pooling.
func (xRefTable *XRefTable) PageContentInto(d types.Dict, pageNr int, buf io.Writer) error {
	o, _ := d.Find("Contents")
	if o == nil {
		return ErrNoContent
	}

	o, err := xRefTable.Dereference(o)
	if err != nil || o == nil {
		return err
	}

	n := 0
	switch o := o.(type) {

	case types.StreamDict:
		// no further processing.
		if err := xRefTable.decodeContentStream(&o, pageNr); err != nil {
			return err
		}
		n, _ = buf.Write(o.Content)

	case types.Array:
		// process array of content stream dicts.
		for _, o := range o {
			if o == nil {
				continue
			}
			o, _, err := xRefTable.DereferenceStreamDict(o)
			if err != nil {
				return fmt.Errorf("page %d content decode: %v", pageNr, err)
			}
			if o == nil {
				continue
			}
			if err := xRefTable.decodeContentStream(o, pageNr); err != nil {
				return err
			}
			m, _ := buf.Write(o.Content)
			n += m
		}

	default:
		return fmt.Errorf("pdfcpu: page content must be stream dict or array")
	}

	if n == 0 {
		return ErrNoContent
	}

	return nil
}
