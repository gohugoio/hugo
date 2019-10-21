
---
date: 2018-08-29
title: "This One Goes to 11!"
description: "With Go 1.11, Hugo finally gets support for variable overwrites in templates!"
categories: ["Releases"]
---

Hugo `0.48` is built with the brand new Go 1.11. On the technical side this means that Hugo now uses [Go Modules](https://github.com/golang/go/wiki/Modules) for the build. The big new functional thing in Go 1.11 for Hugo is added support for [variable overwrites](https://github.com/golang/go/issues/10608). This means that you can now do this and get the expected result:

```go-html-template
{{ $var := "Hugo Page" }}
{{ if .IsHome }}
	{{ $var = "Hugo Home" }}
{{ end }}
Var is {{ $var }}
```

The above may look obvious, but has not been possible until now. In Hugo we have had `.Scratch` as a workaround for this, but Go 1.11 will help clean up a lot of templates.

This release represents **23 contributions by 5 contributors** to the main Hugo code base. [@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@anthonyfok](https://github.com/anthonyfok), [@vsopvsop](https://github.com/vsopvsop), and [@moorereason](https://github.com/moorereason) for their ongoing contributions. And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **15 contributions by 12 contributors**. A special thanks to [@bep](https://github.com/bep), [@kaushalmodi](https://github.com/kaushalmodi), [@regisphilibert](https://github.com/regisphilibert), and [@anthonyfok](https://github.com/anthonyfok) for their work on the documentation site.


Hugo now has:

* 28275+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 441+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 252+ [themes](http://themes.gohugo.io/)

## Enhancements

* Add a test for template variable overwrite [0c8a4154](https://github.com/gohugoio/hugo/commit/0c8a4154838e32a33d34202fd4fa0187aa502190) [@bep](https://github.com/bep) 
* Include language code in REF_NOT_FOUND errors [94d0e79d](https://github.com/gohugoio/hugo/commit/94d0e79d33994b9a9d26a4d020500acdcc71e58c) [@bep](https://github.com/bep) [#5110](https://github.com/gohugoio/hugo/issues/5110)
* Improve minifier MIME type resolution [ebb56e8b](https://github.com/gohugoio/hugo/commit/ebb56e8bdbfaf4f955326017e40b2805850871e9) [@bep](https://github.com/bep) [#5093](https://github.com/gohugoio/hugo/issues/5093)
* Update to Go 1.11 [6b9934a2](https://github.com/gohugoio/hugo/commit/6b9934a26615ea614b1774770532cae9762a58d3) [@bep](https://github.com/bep) [#5115](https://github.com/gohugoio/hugo/issues/5115)
* Set GO111MODULE=on for mage install [c7f05779](https://github.com/gohugoio/hugo/commit/c7f057797ca7bfc781d5a2bbf181bb52360f160f) [@bep](https://github.com/bep) [#5115](https://github.com/gohugoio/hugo/issues/5115)
* Add instruction to install PostCSS when missing [08d14113](https://github.com/gohugoio/hugo/commit/08d14113b60ff70ffe922e8098e289b099a70e0f) [@anthonyfok](https://github.com/anthonyfok) [#5111](https://github.com/gohugoio/hugo/issues/5111)
* Update snapcraft build config to Go 1.11 [94d6d678](https://github.com/gohugoio/hugo/commit/94d6d6780fac78e9ed5ed58ecdb9821ad8f0f27c) [@bep](https://github.com/bep) [#5115](https://github.com/gohugoio/hugo/issues/5115)
* Use Go 1.11 modules with Mage [45c9c45d](https://github.com/gohugoio/hugo/commit/45c9c45d1d0d99443fa6bb524a1265fa9ba95e0e) [@bep](https://github.com/bep) [#5115](https://github.com/gohugoio/hugo/issues/5115)
* Add go.mod [fce32c07](https://github.com/gohugoio/hugo/commit/fce32c07fb80e9929bc2660cf1e681e93009d24b) [@bep](https://github.com/bep) [#5115](https://github.com/gohugoio/hugo/issues/5115)
* Update Travis to Go 1.11 and Go 1.10.4 [d32ff16f](https://github.com/gohugoio/hugo/commit/d32ff16fd61f55874e81d73759afa099b8bdcb57) [@bep](https://github.com/bep) [#5115](https://github.com/gohugoio/hugo/issues/5115)
* Skip installing postcss due to failure on build server [66f688f7](https://github.com/gohugoio/hugo/commit/66f688f7120560ca787c1a23e3e7fbc3aa617956) [@anthonyfok](https://github.com/anthonyfok) 

## Fixes

* Keep end tags [e6eda2a3](https://github.com/gohugoio/hugo/commit/e6eda2a370aa1184e0afaf12e95dbd6f8b63ace5) [@vsopvsop](https://github.com/vsopvsop) 
* Fix permissions when creating new folders [f4675fa0](https://github.com/gohugoio/hugo/commit/f4675fa0f0fae2358adfaea49e8da824ee094495) [@bep](https://github.com/bep) [#5128](https://github.com/gohugoio/hugo/issues/5128)
* Fix handling of taxonomy terms containing slashes [fff13253](https://github.com/gohugoio/hugo/commit/fff132537b4094221f4f099e2251f3cda613060f) [@moorereason](https://github.com/moorereason) [#4090](https://github.com/gohugoio/hugo/issues/4090)
* Fix build on armv7 [8999de19](https://github.com/gohugoio/hugo/commit/8999de193c18b7aa07b44e5b7d9f443a8572e117) [@caarlos0](https://github.com/caarlos0) [#5101](https://github.com/gohugoio/hugo/issues/5101)





