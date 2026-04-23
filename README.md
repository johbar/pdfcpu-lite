# pdfcpu-lite a Go PDF processor lib


[pdfcpu-lite](https://github.com/johbar/pdfcpu-lite) is a PDF processing library written in [Go](https://go.dev/). It is compatible with all PDF versions with basic support and ongoing improvement for PDF 2.0 (ISO-32000-2).

## Motivation for this fork

This opinionated fork aims at faster builds and smaller binaries for library users of [pdfcpu.io](https://pdfcpu.io).

pdfcpu-lite

* removes all test data and bundled certificates
* drops support for certificates and signing
* omits the CLI/CMD packages
* drops handling and creation of config.yaml files
* avoids some dependencies entirely

If you are interested in the CLI or removed functionality, stick with the upstream repository and Go module [pdfcpu](https://github.com/pdfcpu/pdfcpu)

## Documentation

See upstream docs of [pdfcpu.io](https://pdfcpu.io)

## Contributors

Thanks 💚 to all contributors:

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->

||||||||
| :---: | :---: | :---: | :---: | :---: |  :---: | :---: |
| [<img src="https://avatars1.githubusercontent.com/u/11322155?v=4" width="100px"/><br/><sub><b>Horst Rutter</b></sub>](https://github.com/hhrutter) | [<img src="https://avatars0.githubusercontent.com/u/5140211?v=4" width="100px"/><br/><sub><b>haldyr</b></sub>](https://github.com/haldyr) | [<img src="https://avatars3.githubusercontent.com/u/20608155?v=4" width="100px"/><br/><sub><b>Vyacheslav</b></sub>](https://github.com/SimePel) | [<img src="https://avatars1.githubusercontent.com/u/617459?v=4" width="100px"/><br/><sub><b>Erik Unger</b></sub>](https://github.com/ungerik) | [<img src="https://avatars1.githubusercontent.com/u/13079058?v=4" width="100px"/><br/><sub><b>Richard Wilkes</b></sub>](https://github.com/richardwilkes) | [<img src="https://avatars1.githubusercontent.com/u/16303386?s=400&v=4" width="100px"/><br/><sub><b>minenok-tutu</b></sub>](https://github.com/minenok-tutu) | [<img src="https://avatars0.githubusercontent.com/u/1965445?s=400&v=4" width="100px"/><br/><sub><b>Mateusz Burniak</b></sub>](https://github.com/matbur) |
| [<img src="https://avatars2.githubusercontent.com/u/1175110?s=400&v=4" width="100px"/><br/><sub><b>Dmitry Harnitski</b></sub>](https://github.com/dharnitski) | [<img src="https://avatars0.githubusercontent.com/u/1074083?s=400&v=4" width="100px"/><br/><sub><b>ryarnyah</b></sub>](https://github.com/ryarnyah) | [<img src="https://avatars0.githubusercontent.com/u/13267?s=400&v=4" width="100px"/><br/><sub><b>Sam Giffney</b></sub>](https://github.com/s01ipsist) | [<img src="https://avatars3.githubusercontent.com/u/32948066?s=400&v=4" width="100px"/><br /><sub><b>Carlos Eduardo Witte</b></sub>](https://github.com/cewitte) | [<img src="https://avatars1.githubusercontent.com/u/2374948?s=400&u=a36e5f8da8dc1c102bc4d283f25e4c61cae7f985&v=4" width="100px"/><br/><sub><b>minusworld</b></sub>](https://github.com/minusworld) | [<img src="https://avatars0.githubusercontent.com/u/18538487?s=400&u=b9e628dfc60f672a887be2ed04a791195829943e&v=4" width="100px"/><br/><sub><b>Witold Konior</b></sub>](https://github.com/jozuenoon) | [<img src="https://avatars0.githubusercontent.com/u/630151?s=400&v=4" width="100px"/><br/><sub><b>joonas.fi</b></sub>](https://github.com/joonas-fi) |
| [<img src="https://avatars3.githubusercontent.com/u/10349817?s=400&u=93bacb23bd2909d5b6c5b644a8d4cdd947422ee1&v=4" width="100px"/><br/><sub><b>Henrik Reinstädtler</b></sub>](https://github.com/henrixapp) | [<img src="https://avatars1.githubusercontent.com/u/72016286?s=400&v=4" width="100px"/><br/><sub><b>VMorozov-wh</b></sub>](https://github.com/VMorozov-wh) | [<img src="https://avatars0.githubusercontent.com/u/31929422?s=400&v=4" width="100px"/><br/><sub><b>Benoit KUGLER</b></sub>](https://github.com/benoitkugler) | [<img src="https://avatars.githubusercontent.com/u/704919?s=400&v=4" width="100px"/><br/><sub><b>Adam Greenhall</b></sub>](https://github.com/adamgreenhall) | [<img src="https://avatars.githubusercontent.com/u/5201812?s=400&u=8a0a9fca4560be71d4923299ddebf877854eea54&v=4" width="100px"/><br/><sub><b>moritamori</b></sub>](https://github.com/moritamori) | [<img src="https://avatars.githubusercontent.com/u/41904529?s=400&u=044396494285ad806e86d1936c390b3071ce57c0&v=4" width="100px"/><br/><sub><b>JanBaryla</b></sub>](https://github.com/JanBaryla) | [<img src="https://avatars.githubusercontent.com/u/43145244?s=400&u=89a689f1a854ce0f57ae2a0333c82bfdc5723bb9&v=4" width="100px"/><br/><sub><b>TheDiscordian</b></sub>](https://github.com/TheDiscordian) |
| [<img src="https://avatars.githubusercontent.com/u/15472552?v=4" width="100px"/><br/><sub><b>Rafael Garcia Argente</b></sub>](https://github.com/rgargente) | [<img src="https://avatars.githubusercontent.com/u/710057?v=4" width="100px"/><br/><sub><b>truyet</b></sub>](https://github.com/truyet) | [<img src="https://avatars.githubusercontent.com/u/5031217?v=4" width="100px"/><br/><sub><b>Christian Nicola</b></sub>](https://github.com/christiannicola) | [<img src="https://avatars.githubusercontent.com/u/3233970?v=4" width="100px"/><br/><sub><b>Benjamin Krill</b></sub>](https://github.com/kben) | [<img src="https://avatars.githubusercontent.com/u/26521615?v=4" width="100px"/><br/><sub><b>Peter Wyatt</b></sub>](https://github.com/petervwyatt) | [<img src="https://avatars.githubusercontent.com/u/3142701?v=4" width="100px"/><br/><sub><b>Kroum Tzanev</b></sub>](https://github.com/kpym) | [<img src="https://avatars.githubusercontent.com/u/992878?v=4" width="100px"/><br/><sub><b>Stefan Huber</b></sub>](https://github.com/signalwerk) |
| [<img src="https://avatars.githubusercontent.com/u/59667587?v=4" width="100px"/><br/><sub><b>Juan Iscar</b></sub>](https://github.com/juaismar) | [<img src="https://avatars.githubusercontent.com/u/20135478?v=4" width="100px"/><br/><sub><b>Eng Zer Jun</b></sub>](https://github.com/Juneezee) | [<img src="https://avatars.githubusercontent.com/u/28459131?v=4" width="100px"/><br/><sub><b>Dmitry Ivanov</b></sub>](https://github.com/hant0508)|[<img src="https://avatars.githubusercontent.com/u/16866547?v=4" width="100px"/><br/><sub><b>Rene Kaufmann</b></sub>](https://github.com/HeavyHorst)|[<img src="https://avatars.githubusercontent.com/u/26827864?v=4" width="100px"/><br/><sub><b>Christian Heusel</b></sub>](https://github.com/christian-heusel) | [<img src="https://avatars.githubusercontent.com/u/305673?v=4" width="100px"/><br/><sub><b>Chris</b></sub>](https://github.com/freshteapot) | [<img src="https://avatars.githubusercontent.com/u/2892794?v=4" width="100px"/><br/><sub><b>Lukasz Czaplinski</b></sub>](https://github.com/scoiatael) |
[<img src="https://avatars.githubusercontent.com/u/49206635?v=4" width="100px"/><br/><sub><b>Joel Silva Schutz</b></sub>](https://github.com/joelschutz) | [<img src="https://avatars.githubusercontent.com/u/28219076?v=4" width="100px"/><br/><sub><b>semvis123</b></sub>](https://github.com/semvis123) | [<img src="https://avatars.githubusercontent.com/u/8717479?v=4"  width="100px"/><br/><sub><b>guangwu</b></sub>](https://github.com/testwill) | [<img src="https://avatars.githubusercontent.com/u/4014912?v=4"  width="100px"/><br/><sub><b>Yoshiki Nakagawa</b></sub>](https://github.com/yyoshiki41) | [<img src="https://avatars.githubusercontent.com/u/432860?v=4"  width="100px"/><br/><sub><b>Steve van Loben Sels</b></sub>](https://github.com/stevevls) | [<img src="https://avatars.githubusercontent.com/u/6083533?v=4" width="100px"/><br/><sub><b>Yaofu</b></sub>](https://github.com/mygityf) | [<img src="https://avatars.githubusercontent.com/u/15724278?v=4" width="100px"/><br/><sub><b>vsenko</b></sub>](https://github.com/vsenko) |
[<img src="https://avatars.githubusercontent.com/u/16507?v=4" width="100px"/><br/><sub><b>Alexis Hildebrandt</b></sub>](https://github.com/afh) | [<img src="https://avatars.githubusercontent.com/u/1395040?v=4" width="100px"/><br/><sub><b>Sivukhin Nikita</b></sub>](https://github.com/sivukhin)  | [<img src="https://avatars.githubusercontent.com/u/247730?v=4"  width="100px"/><br/><sub><b>Joachim Bauch</b></sub>](https://github.com/fancycode) | [<img src="https://avatars.githubusercontent.com/u/127291996?v=4"  width="100px"/><br/><sub><b>kalimit</b></sub>](https://github.com/kalimit) | [<img src="https://avatars.githubusercontent.com/u/5080535?v=4"  width="100px"/><br/><sub><b>Andreas Erhard</b></sub>](https://github.com/xelan) | [<img src="https://avatars.githubusercontent.com/u/32378535?v=4"  width="100px"/><br/><sub><b>Matsumoto Toshi</b></sub>](https://github.com/toshi1127) | [<img src="https://avatars.githubusercontent.com/u/440634?v=4"  width="100px"/><br/><sub><b>Carl Wilson</b></sub>](https://github.com/carlwilson) |
[<img src="https://avatars.githubusercontent.com/u/9918222?v=4" width="100px"/><br/><sub><b>LNAhri</b></sub>](https://github.com/LNAhri) | [<img src="https://avatars.githubusercontent.com/u/142796877?v=4" width="100px"/><br/><sub><b>vishal</b></sub>](https://github.com/vishal-at) | [<img src="https://avatars.githubusercontent.com/u/18169566?v=4" width="100px"/><br/><sub><b>Andreas Deininger</b></sub>](https://github.com/deining) | [<img src="https://avatars.githubusercontent.com/u/5825735?v=4" width="100px"/><br/><sub><b>Robert Raines</b></sub>](https://github.com/solintllc-robert) | [<img src="https://avatars.githubusercontent.com/u/316176?v=4" width="100px"/><br/><sub><b>Frank Anderson</b></sub>](https://github.com/frob) |  [<img src="https://avatars.githubusercontent.com/u/20972350?v=4" width="100px"/><br/><sub><b>Sven Lilienthal</b></sub>](https://github.com/SveLil) |  [<img src="https://avatars.githubusercontent.com/u/1900106?v=4" width="100px"/><br/><sub><b>Florian Kinder</b></sub>](https://github.com/fank) |
[<img src="https://avatars.githubusercontent.com/u/113192632?v=4" width="100px"/><br/><sub><b>mdmcconnell</b></sub>](https://github.com/mdmcconnell) | [<img src="https://avatars.githubusercontent.com/u/195061?v=4" width="100px"/><br/><sub><b>Bradley Erickson</b></sub>](https://github.com/13rac1) | [<img src="https://avatars.githubusercontent.com/u/10998835?v=4" width="100px"/><br/><sub><b>doronbehar</b></sub>](https://github.com/doronbehar) |[<img src="https://avatars.githubusercontent.com/u/48005710?v=4" width="100px"/><br/><sub><b>joeyave</b></sub>](https://github.com/joeyave) ||||

<!-- ALL-CONTRIBUTORS-LIST:END - Do not remove or modify this section -->

---
## Disclaimer

Use of pdfcpu assumes compliance with all applicable copyrights for any processed PDF content, including embedded resources such as fonts and images.

Gopher artwork by [Renee French](https://www.instagram.com/reneefrenchphoto)

---

## License

Apache License 2.0

---

<img referrerpolicy="no-referrer-when-downgrade" src="https://static.scarf.sh/a.png?x-pxid=592f747a-ac99-42e3-802a-80dbf2b94519" width="1" height="1" />
