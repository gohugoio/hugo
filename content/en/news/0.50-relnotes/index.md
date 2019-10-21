
---
date: 2018-10-29
title: "Hugo 0.50: Errors so Good, You’ll Want to Fail!"
description: "Errors with full filename and line and column number, shown in browser. And improved Fast Render Mode …"
categories: ["Releases"]
---

Hugo `0.50` brings **greatly improved error messages**, and we now also show them in the browser. Having error messages with filename, line- and column number greatly simplify troubleshooting. Many editors (like VS Code) even let you click and navigate directly to the problematic line. If your editor requires a different log format, you can set it via the `HUGO_FILE_LOG_FORMAT` OS environment variable:


```bash
env HUGO_FILE_LOG_FORMAT="\":file::line::col\"" hugo server
```

But this release isn't all about error handling. Getting line- and column number into "every" error also meant that we had to consolidate and simplify some code paths, which, as a nice side effect, made Hugo a little bit faster. Benchmarks show it running **about 5% faster and consume about 8% less memory**.

Also, we have now implemented **"render on demand"** in Hugo's Fast Render Mode (default when running `hugo server`). This means that you should now always see updated content when navigating around the site after a change.

This release represents **88 contributions by 14 contributors** to the main Hugo code base.
[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@moorereason](https://github.com/moorereason), [@anthonyfok](https://github.com/anthonyfok), and [@GregorioMartinez](https://github.com/GregorioMartinez) for their ongoing contributions.

And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **14 contributions by 9 contributors**. A special thanks to [@bep](https://github.com/bep), [@NotWoods](https://github.com/NotWoods), [@Nick-Rivera](https://github.com/Nick-Rivera), and [@tomanistor](https://github.com/tomanistor) for their work on the documentation site.

Hugo now has:

* 29842+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 441+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 275+ [themes](http://themes.gohugo.io/)

## Notes

* You should not get stale content in Fast Render Mode anymore.
* Errors will now show up in the browser by default, turn it off by running `hugo server --disableBrowserError`
* `jsonify` will now produce pretty/indented output

## Enhancements

### Templates

* Handle truncated identifiers in Go template errors [2d7709d1](https://github.com/gohugoio/hugo/commit/2d7709d15584e4c11138cd7fe92717a2a58e4585) [@bep](https://github.com/bep) [#5346](https://github.com/gohugoio/hugo/issues/5346)
* Update Jsonify to return pretty-print output [5a52cd5f](https://github.com/gohugoio/hugo/commit/5a52cd5f920bb3d067ab1682adece9f813c67ba1) [@SeanPrashad](https://github.com/SeanPrashad) [#5040](https://github.com/gohugoio/hugo/issues/5040)
* Improve the Execute panic error message [0fe4ff18](https://github.com/gohugoio/hugo/commit/0fe4ff18751156fa072e1f83077e49a8597e7dcd) [@bep](https://github.com/bep) [#5327](https://github.com/gohugoio/hugo/issues/5327)
* Use .Lastmod in embedded schema template [c21e5179](https://github.com/gohugoio/hugo/commit/c21e5179ce9a370c416c01fbe9472be1fb5c6650) [@akshaybabloo](https://github.com/akshaybabloo) [#5320](https://github.com/gohugoio/hugo/issues/5320)
* Cast IsSet key to int for indexed types [0d5110d0](https://github.com/gohugoio/hugo/commit/0d5110d03324380cb4a288d3fa08c1b86ba227da) [@moorereason](https://github.com/moorereason) [#3681](https://github.com/gohugoio/hugo/issues/3681)
* Add a delimiter parameter to lang.NumFmt [ce264b93](https://github.com/gohugoio/hugo/commit/ce264b936ce9f589bd889e18762881cff8bc9cd0) [@moorereason](https://github.com/moorereason) [#5260](https://github.com/gohugoio/hugo/issues/5260)

### Core

* Adjust error test to make it pass on Go tip [acc14b46](https://github.com/gohugoio/hugo/commit/acc14b4646d849e09e8da37552d4f4f777d0fce2) [@bep](https://github.com/bep) 
* Rename some page_* files [e3ed4a83](https://github.com/gohugoio/hugo/commit/e3ed4a83b8e92ce9bf070f7b41780798b006e848) [@bep](https://github.com/bep) 
* Get file context in "config parse failed" errors [ed7b3e26](https://github.com/gohugoio/hugo/commit/ed7b3e261909fe425ef64216f12806840c45b205) [@bep](https://github.com/bep) [#5325](https://github.com/gohugoio/hugo/issues/5325)
* Improve errors in /i18n handlling [2bf686ee](https://github.com/gohugoio/hugo/commit/2bf686ee217808186385bfcf6156f15bbdb33651) [@bep](https://github.com/bep) [#5324](https://github.com/gohugoio/hugo/issues/5324)
* Improve errors in /data handlling [9f74dc2a](https://github.com/gohugoio/hugo/commit/9f74dc2a52b6f568b5a060b7a4be47196804b01f) [@bep](https://github.com/bep) [#5324](https://github.com/gohugoio/hugo/issues/5324)
* Continue the file context/line number errors work [d1661b82](https://github.com/gohugoio/hugo/commit/d1661b823af25c50d3bbe5366ea40a3cdd52e237) [@bep](https://github.com/bep) [#5324](https://github.com/gohugoio/hugo/issues/5324)
* Remove the now superflous Source struct [7930d213](https://github.com/gohugoio/hugo/commit/7930d2132a3c36c1aaca20f16f56978c84656b0a) [@bep](https://github.com/bep) [#5324](https://github.com/gohugoio/hugo/issues/5324)
* Redo the summary delimiter logic [44da60d8](https://github.com/gohugoio/hugo/commit/44da60d869578423dea529db62ed613588a2a560) [@bep](https://github.com/bep) [#5324](https://github.com/gohugoio/hugo/issues/5324)
* Integrate new page parser [1e3e3400](https://github.com/gohugoio/hugo/commit/1e3e34002dae3d4a980141efcc86886e7de5bef8) [@bep](https://github.com/bep) [#5324](https://github.com/gohugoio/hugo/issues/5324)
* Use []byte in shortcode parsing [1b7ecfc2](https://github.com/gohugoio/hugo/commit/1b7ecfc2e176315b69914756c70b46306561e4d1) [@bep](https://github.com/bep) [#5324](https://github.com/gohugoio/hugo/issues/5324)
* Use stdlib context package [4b4af2c5](https://github.com/gohugoio/hugo/commit/4b4af2c52e658d516dd4bfaf59fef4f054dabec3) [@GregorioMartinez](https://github.com/GregorioMartinez) 
* Normalize permalink path segments [fae48d74](https://github.com/gohugoio/hugo/commit/fae48d7457de96969ec53349194dcbfa45adc269) [@moorereason](https://github.com/moorereason) [#5223](https://github.com/gohugoio/hugo/issues/5223)[#4926](https://github.com/gohugoio/hugo/issues/4926)
* Improve error message for bad taxonomy weights [d3b81ee5](https://github.com/gohugoio/hugo/commit/d3b81ee58e8fd3a0ab8265a2898d66cbcdf6a7c1) [@moorereason](https://github.com/moorereason) 
* Cast taxonomy weight parameters to int [1fd30d46](https://github.com/gohugoio/hugo/commit/1fd30d462ee7f67fde6f29d170af1d225258322b) [@moorereason](https://github.com/moorereason) [#4628](https://github.com/gohugoio/hugo/issues/4628)
* Allow nil to be unwrapped as *Page [498d6299](https://github.com/gohugoio/hugo/commit/498d6299581bead0f582431b8133d8b5f8760618) [@moorereason](https://github.com/moorereason) [#5043](https://github.com/gohugoio/hugo/issues/5043)
* Be a litle more specific in NextPage TODO [fb732d53](https://github.com/gohugoio/hugo/commit/fb732d5322381ee7f3a849258419cef7bbf1487b) [@bep](https://github.com/bep) 
* Introduce Page.NextPage and Page.PrevPage [ad705aac](https://github.com/gohugoio/hugo/commit/ad705aac0649fa3102f7639bc4db65d45e108ee2) [@felicianotech](https://github.com/felicianotech) [#1061](https://github.com/gohugoio/hugo/issues/1061)

### Other

* Update go.sum [7082a5d1](https://github.com/gohugoio/hugo/commit/7082a5d14382acfc300ae4f66d07a14100e2358c) [@bep](https://github.com/bep) 
* Update minify [aa281b51](https://github.com/gohugoio/hugo/commit/aa281b5135db2e84b9e21b5f38a6cb63cf3ab158) [@bep](https://github.com/bep) [#5261](https://github.com/gohugoio/hugo/issues/5261)
* Regenerate CLI docs [32501987](https://github.com/gohugoio/hugo/commit/325019872467ee152ea56329a796acf35dec6fb3) [@bep](https://github.com/bep) 
* Make sure the global logger also gets colored labels [9c88a8a5](https://github.com/gohugoio/hugo/commit/9c88a8a55adf7779039504fa77d74ec80d658c40) [@bep](https://github.com/bep) [#4414](https://github.com/gohugoio/hugo/issues/4414)
* Avoid using the global logger [95e72f5e](https://github.com/gohugoio/hugo/commit/95e72f5e8e4634fbbb2ea7ece2156487230ad1d4) [@bep](https://github.com/bep) [#4414](https://github.com/gohugoio/hugo/issues/4414)
* Add color to ERROR and WARN [1c7b7b4e](https://github.com/gohugoio/hugo/commit/1c7b7b4ef293aa133e5b55f3ebb2d37d8839532f) [@bep](https://github.com/bep) [#4414](https://github.com/gohugoio/hugo/issues/4414)
* Make the file error log format configurable [1ad117cb](https://github.com/gohugoio/hugo/commit/1ad117cbe2903aa9d029f90750acf633eb2a51a2) [@bep](https://github.com/bep) [#5352](https://github.com/gohugoio/hugo/issues/5352)
* Allow a mix of slice types in append/Scratch.Add [dac7092a](https://github.com/gohugoio/hugo/commit/dac7092a9cb22d59db28fb15af15f7b14ff47588) [@bep](https://github.com/bep) [#5361](https://github.com/gohugoio/hugo/issues/5361)
* Allow .Data.Integrity to be accessed on its own [b27ccf34](https://github.com/gohugoio/hugo/commit/b27ccf34bf4e5ee618a66fa11c68a9690e395034) [@bep](https://github.com/bep) [#5296](https://github.com/gohugoio/hugo/issues/5296)
* Update minify [83c873ff](https://github.com/gohugoio/hugo/commit/83c873ff37ddd379181540021232f026e7678486) [@bep](https://github.com/bep) [#5261](https://github.com/gohugoio/hugo/issues/5261)
* Update cast [a2440dc0](https://github.com/gohugoio/hugo/commit/a2440dc0e2d46ef774305cd5e4fea5ff2bdd5f11) [@bep](https://github.com/bep) [#5340](https://github.com/gohugoio/hugo/issues/5340)
* Truncate the error log on repeated config errors [1e9ac3dc](https://github.com/gohugoio/hugo/commit/1e9ac3dcc21e8f78d3f0a0ba4f35f6c142dfa6bc) [@bep](https://github.com/bep) 
* Regenerate CLI docs [40e99672](https://github.com/gohugoio/hugo/commit/40e99672b6f697a614485aace07ca84268f6c787) [@bep](https://github.com/bep) [#5354](https://github.com/gohugoio/hugo/issues/5354)
* Serialize image processing [3a3badfd](https://github.com/gohugoio/hugo/commit/3a3badfd1d4b1d4c9863ecaf029512d36136fa0f) [@bep](https://github.com/bep) [#5220](https://github.com/gohugoio/hugo/issues/5220)
* Only show Ansi escape codes if in a terminal [df021317](https://github.com/gohugoio/hugo/commit/df021317a964a482cd1cd579de5a12d50faf0d08) [@bep](https://github.com/bep) 
* Read disableFastRender from flag even if it's not changed [78a4c2e3](https://github.com/gohugoio/hugo/commit/78a4c2e32ef9ea8e92bb7bb3586e4c22b02eb494) [@bep](https://github.com/bep) [#5353](https://github.com/gohugoio/hugo/issues/5353)
* Use overflow-x: auto; for browser errors [d4ebfea1](https://github.com/gohugoio/hugo/commit/d4ebfea1fffdc35059f42a46387e0aaf0ea877d2) [@bep](https://github.com/bep) 
* Remove the ANSI color for the browser error version [93aa6261](https://github.com/gohugoio/hugo/commit/93aa6261b4fc8caa74afef97b6304ea35dfd7d0e) [@bep](https://github.com/bep) 
* Add some color to the relevant filenames in terminal log [deff9e15](https://github.com/gohugoio/hugo/commit/deff9e154bc0371af56741ddb22cb1f9e392838a) [@bep](https://github.com/bep) [#5344](https://github.com/gohugoio/hugo/issues/5344)
* Run gofmt -s [889aca05](https://github.com/gohugoio/hugo/commit/889aca054a267506a1c7cfaa3992d324764d6358) [@bep](https://github.com/bep) 
* Resolve error handling/parser related TODOs [6636cf1b](https://github.com/gohugoio/hugo/commit/6636cf1bea77d20ef2a72a45fae59ac402fb133b) [@bep](https://github.com/bep) [#5324](https://github.com/gohugoio/hugo/issues/5324)
* Improve handling of JSON errors [f669ef6b](https://github.com/gohugoio/hugo/commit/f669ef6bec25155d015b6ab231c53caef4fa5cdc) [@bep](https://github.com/bep) [#5324](https://github.com/gohugoio/hugo/issues/5324)
* Convert the rest to new page parser code paths [eb038cfa](https://github.com/gohugoio/hugo/commit/eb038cfa0a8ada29dfcba1204ec5c432da9ed7e0) [@bep](https://github.com/bep) [#5324](https://github.com/gohugoio/hugo/issues/5324)
* Consolidate the metadata decoders [129c27ee](https://github.com/gohugoio/hugo/commit/129c27ee6e9fed98dbfebeaa272fd52757b475b2) [@bep](https://github.com/bep) [#5324](https://github.com/gohugoio/hugo/issues/5324)
* Use []byte in page lexer [27f5a906](https://github.com/gohugoio/hugo/commit/27f5a906a2a34e3b8348c8baeea48355352b5bbb) [@bep](https://github.com/bep) [#5324](https://github.com/gohugoio/hugo/issues/5324)
* Add front matter etc. support [2fdc4a24](https://github.com/gohugoio/hugo/commit/2fdc4a24d5450a98cf38a4456e8e0e8e97a3343d) [@bep](https://github.com/bep) [#5324](https://github.com/gohugoio/hugo/issues/5324)
* File renames and splitting [f6863e1e](https://github.com/gohugoio/hugo/commit/f6863e1ef725f654a4c869ef4955f9add6908a46) [@bep](https://github.com/bep) [#5324](https://github.com/gohugoio/hugo/issues/5324)
* Move the shortcode parser to the new pageparser package [d6c16afd](https://github.com/gohugoio/hugo/commit/d6c16afde0ce62cfea73447f30d6ed2b8ef4b411) [@bep](https://github.com/bep) [#5324](https://github.com/gohugoio/hugo/issues/5324)
* Avoid panic in error handler on config errors [6f3716dc](https://github.com/gohugoio/hugo/commit/6f3716dc22e373097a38d053f5415feca602f330) [@bep](https://github.com/bep) 
* Prevent stale content in Fast Render Mode [4a366fcf](https://github.com/gohugoio/hugo/commit/4a366fcfee24b3a5a5045b16c3b87b76147adf5e) [@bep](https://github.com/bep) [#5281](https://github.com/gohugoio/hugo/issues/5281)
* Allow date and slug from filename for leaf bundles [1f42e47e](https://github.com/gohugoio/hugo/commit/1f42e47e475c0cd684426dd230de411d4c385a3c) [@Japanuspus](https://github.com/Japanuspus) [#4558](https://github.com/gohugoio/hugo/issues/4558)
* Show server error info in browser [35fbfb19](https://github.com/gohugoio/hugo/commit/35fbfb19a173b01bc881f2bbc5d104136633a7ec) [@bep](https://github.com/bep) [#5284](https://github.com/gohugoio/hugo/issues/5284)[#5290](https://github.com/gohugoio/hugo/issues/5290)[#5325](https://github.com/gohugoio/hugo/issues/5325)[#5324](https://github.com/gohugoio/hugo/issues/5324)
* Add .gitignore "hugo new site" [92979d92](https://github.com/gohugoio/hugo/commit/92979d92889ff8991acfadd42982c2e55d01b163) [@napei](https://github.com/napei) 
* Optimize integrity string generation [0a3340e9](https://github.com/gohugoio/hugo/commit/0a3340e95254597bc8a9feb250f2733b7d51edf8) [@moorereason](https://github.com/moorereason) 
* Add help text to "hugo new" [6b21ac3e](https://github.com/gohugoio/hugo/commit/6b21ac3e67cb101255e8c3d9dbf076391a9eed8d) [@elliotforbes](https://github.com/elliotforbes) 
* Set "extended" tag based on build_url on Launchpad [d1442053](https://github.com/gohugoio/hugo/commit/d14420539ac04d31dde2252eee66d7e4c7749780) [@anthonyfok](https://github.com/anthonyfok) 
* Call rst2html directly on *nix [3d4a9882](https://github.com/gohugoio/hugo/commit/3d4a9882bfc81215fb4f9eba8859324958747d4a) [@shreyanshk](https://github.com/shreyanshk) 
* Update URLs to stop 301 redirects [bdca9727](https://github.com/gohugoio/hugo/commit/bdca9727944e4cbb5a9372a8404e948ffea7c31c) [@benabbottnz](https://github.com/benabbottnz) 
* Merge branch 'release-0.49.2' [604ddb90](https://github.com/gohugoio/hugo/commit/604ddb90c5d6f1ca5583be1ec0ea8e48f014741a) [@bep](https://github.com/bep) 
* Release 0.49.2 [c397f2c0](https://github.com/gohugoio/hugo/commit/c397f2c08087cf4cda3abe2d146e30f58d6d3216) [@bep](https://github.com/bep) 
* Merge branch 'release-0.49.1' [3583dd6d](https://github.com/gohugoio/hugo/commit/3583dd6d713c243808b5e8724b32565ceaf66104) [@bep](https://github.com/bep) 
* Relase 0.49.1 [235acf22](https://github.com/gohugoio/hugo/commit/235acf22321475895442ce49ca5d16be273c1e1f) [@bep](https://github.com/bep) 
* Improve append in Scratch [23f48c30](https://github.com/gohugoio/hugo/commit/23f48c300cb5ffe0fe43c88464f38c68831a17ad) [@bep](https://github.com/bep) [#5275](https://github.com/gohugoio/hugo/issues/5275)
* Add GOPATH Hugo building tip [b5e17f7c](https://github.com/gohugoio/hugo/commit/b5e17f7c837ce796e1094c8033fa7084510402fb) [@bep](https://github.com/bep) 
* Consolidate MakeSegment vs MakePathSanitized [e421696d](https://github.com/gohugoio/hugo/commit/e421696d02bfb8764ae57238e211ce0e85e9782e) [@bep](https://github.com/bep) [#4926](https://github.com/gohugoio/hugo/issues/4926)
* Render Markdown in figure shortcode "caption" and "attr" params [68181703](https://github.com/gohugoio/hugo/commit/6818170308994b5f01dec7a559f92d6c7c5ca100) [@kaushalmodi](https://github.com/kaushalmodi) 
* Re-organize the figure shortcode for better readability [c5279064](https://github.com/gohugoio/hugo/commit/c5279064df9664d6b2ad277e2fba1e4bb3b0f4be) [@kaushalmodi](https://github.com/kaushalmodi) 
* Update README & CONTRIBUTING [152cffb1](https://github.com/gohugoio/hugo/commit/152cffb13a237651c2277dc6c2c9e4172d58b3df) [@GregorioMartinez](https://github.com/GregorioMartinez) 
* Add custom x-nodejs plugin to support ppc64el and s390x [91f49c07](https://github.com/gohugoio/hugo/commit/91f49c0700dde13e16f42c745584a0bef60c6fe2) [@anthonyfok](https://github.com/anthonyfok) 
* Fetch mage with GO111MODULE=off [a475bf12](https://github.com/gohugoio/hugo/commit/a475bf125cd76dacc1bf7ccbcc263a7b59efc510) [@anthonyfok](https://github.com/anthonyfok) 
* Use build-snaps instead of building go from source [fa873a6c](https://github.com/gohugoio/hugo/commit/fa873a6cb3f0fa81002fcd725ecd52fc4b9df48f) [@anthonyfok](https://github.com/anthonyfok) 
* Skip "mage -v test" due to build failure on Launchpad [52ac85fb](https://github.com/gohugoio/hugo/commit/52ac85fbc4d4066b5e13df454593597df0166262) [@anthonyfok](https://github.com/anthonyfok) 
* Move snapcraft.yaml to snap/snapcraft.yaml [27d42111](https://github.com/gohugoio/hugo/commit/27d4211187d4617f4b3afa970f91349567886748) [@anthonyfok](https://github.com/anthonyfok) 
* Update the temp docker script [48413d76](https://github.com/gohugoio/hugo/commit/48413d76f44ecfc9b90f9df63974080f6b285667) [@bep](https://github.com/bep) 

## Fixes

### Templates

* Fix baseof.html in error message [646a52a5](https://github.com/gohugoio/hugo/commit/646a52a5c5f52890f2d0270b68ba0f863047484a) [@bep](https://github.com/bep) [#5288](https://github.com/gohugoio/hugo/issues/5288)
* Fix handling of different interface types in Slice [10ac2ec4](https://github.com/gohugoio/hugo/commit/10ac2ec4466090957e1f6897906ddeb1e0b13673) [@bep](https://github.com/bep) [#5269](https://github.com/gohugoio/hugo/issues/5269)

### Core

* Fix test on Windows [083311d0](https://github.com/gohugoio/hugo/commit/083311d0336ced35909b3375950f7817ecf95ed0) [@bep](https://github.com/bep) 
* Fix FuzzyWordCount test error message [06d28a46](https://github.com/gohugoio/hugo/commit/06d28a464d086880f52dd850f91e668ab957b31f) [@GregorioMartinez](https://github.com/GregorioMartinez) 

### Other

* Fix archetype handling of directories in theme [78578632](https://github.com/gohugoio/hugo/commit/78578632f545283741a01f024a6ccedc0b695a30) [@bep](https://github.com/bep) [#5318](https://github.com/gohugoio/hugo/issues/5318)
* Simple doc fix in CONTRIBUTING.md [3a308912](https://github.com/gohugoio/hugo/commit/3a3089121b852332b5744d1f566959c8cf93cef4) [@krisbudhram](https://github.com/krisbudhram) 
* Fix type checking in Append [2159d77f](https://github.com/gohugoio/hugo/commit/2159d77f368eb1f78e51dd94133554f88052d85f) [@bep](https://github.com/bep) [#5303](https://github.com/gohugoio/hugo/issues/5303)
* Fix go plugin build failure by renaming go.mod [3033a9a3](https://github.com/gohugoio/hugo/commit/3033a9a37eb66c08e60f9fe977f29d22bd646857) [@anthonyfok](https://github.com/anthonyfok) 





