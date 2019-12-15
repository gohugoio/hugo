
---
date: 2019-11-27
title: "Now CommonMark Compliant!"
description: "Goldmark -- CommonMark compliant,  GitHub flavored, fast and flexible -- is the new default library for Markdown in Hugo."
categories: ["Releases"]
---

[Goldmark](https://github.com/yuin/goldmark/) by [@yuin](https://github.com/yuin) is now the new default library used for Markdown in Hugo. It's CommonMark compliant and  GitHub flavored, and both fast and flexible. Blackfriday, the old default, has served us well, but there have been formatting and portability issues that were hard to work around. The "CommonMark compliant" part is the main selling feature of Goldmark, but with that you also get attribute syntax on headers and code blocks (for code blocks you can turn on/off line numbers and highlight line ranges), strikethrough support and an improved and configurable implementation of `TableOfContents`. See [Markup Configuration](https://gohugo.io/getting-started/configuration-markup/) for an overview of extensions.

Please read the [Notes Section](#notes) and the updated documentation. We suggest you start with [List of content formats in Hugo](https://gohugo.io/content-management/formats/#list-of-content-formats). Goldmark is better, but the feature set is not fully comparable and it may be more stricter in some areas (there are 17 rules for how a [headline](https://spec.commonmark.org/0.29/#emphasis-and-strong-emphasis) should look like); if you have any problems you cannot work around, see [Configure Markup](https://gohugo.io/getting-started/configuration-markup/#configure-markup) for a way to change the default Markdown handler.

Also, if you have lots of inline HTML in your Markdown files, you may have to enable the `unsafe` mode:

{{< code-toggle file="config" >}}
markup:
  goldmark:
    renderer:
      unsafe: true
{{< /code-toggle >}}

This release represents **62 contributions by 10 contributors** to the main Hugo code base. [@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@anthonyfok](https://github.com/anthonyfok), [@max-arnold](https://github.com/max-arnold), and [@trimbo](https://github.com/trimbo) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) and [@davidsneighbour](https://github.com/davidsneighbour) for great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **8 contributions by 4 contributors**. A special thanks to [@bep](https://github.com/bep), [@jasdeepgill](https://github.com/jasdeepgill), [@luucamay](https://github.com/luucamay), and [@jkreft-usgs](https://github.com/jkreft-usgs) for their work on the documentation site.


Hugo now has:

* 39668+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 440+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 274+ [themes](http://themes.gohugo.io/)


## Notes

* Permalink config now supports Go date format strings. [#6489](https://github.com/gohugoio/hugo/pull/6489)
* We have removed the option to use Pygments as a highlighter. [#4491](https://github.com/gohugoio/hugo/pull/4491)
* Config option for code highlighting of code fences in Markdown is now default on. This is what most people wants.
* There are some differences in the feature set of Goldmark and Blackfriday. See the documentation for details.
* The highlight shortcode/template func and the code fence attributes now share the same API regarding line numbers and highlight ranges.
* The `Total in ...` for the `hugo` command now includes the configuration and modules loading, which should make it more honest/accurate.
* The image logic in the 3 SEO internal templates twitter_cards.html,  opengraph.html, and schema.html is consolidated: `images` page param first, then bundled image matching `*feature*`, `*cover*` or `*thumbnail*`, then finally `images` site param.
* Deprecate mmark [33d73330](https://github.com/gohugoio/hugo/commit/33d733300a4f0b765234706e51bb7e077fdc2471) [@bep](https://github.com/bep) [#6486](https://github.com/gohugoio/hugo/issues/6486)

## Enhancements

### Templates

* Featured and Site.Params image support for Schema [c91970c0](https://github.com/gohugoio/hugo/commit/c91970c08ddf8c22ca4f967c2cc864c483987ac7) [@max-arnold](https://github.com/max-arnold) 
* Add support for featured and global image to OpenGraph template [25a6b336](https://github.com/gohugoio/hugo/commit/25a6b33693992e8c6d9c35bc1e781ce3e2bca4be) [@max-arnold](https://github.com/max-arnold) 
* Allow dict to create nested structures [a2670bf4](https://github.com/gohugoio/hugo/commit/a2670bf460e10ed5de69f90abbe7c4e2b32068cf) [@bep](https://github.com/bep) [#6497](https://github.com/gohugoio/hugo/issues/6497)
* Add collections.Reverse [90d0cdf2](https://github.com/gohugoio/hugo/commit/90d0cdf236b54000bfe444ba3a00236faaa28790) [@bep](https://github.com/bep) [#6499](https://github.com/gohugoio/hugo/issues/6499)
* Make index work with slice as the last arg [95ef93be](https://github.com/gohugoio/hugo/commit/95ef93be667afb480184175a319584fd651abf03) [@bep](https://github.com/bep) [#6496](https://github.com/gohugoio/hugo/issues/6496)
* Add some index map test cases [9f46a72c](https://github.com/gohugoio/hugo/commit/9f46a72c7eec25a4b9dea387d5717173b8d9ec72) [@bep](https://github.com/bep) [#3974](https://github.com/gohugoio/hugo/issues/3974)

### Output

* Add some more output if modules download takes time [14a1de14](https://github.com/gohugoio/hugo/commit/14a1de14fb1ec93444ba5dd028fdad8959924545) [@bep](https://github.com/bep) [#6519](https://github.com/gohugoio/hugo/issues/6519)
* Add some more output if loading modules takes time [2dcc1318](https://github.com/gohugoio/hugo/commit/2dcc1318d1d9ed849d040115aa5ba6191a1c102a) [@bep](https://github.com/bep) [#6519](https://github.com/gohugoio/hugo/issues/6519)

### Core

* Disable test assertion on Windows [dd1e5fc0](https://github.com/gohugoio/hugo/commit/dd1e5fc0b43739941372c0c27b75977380acd582) [@bep](https://github.com/bep) 
* Adjust .Site.Permalinks deprecation level [03b369e6](https://github.com/gohugoio/hugo/commit/03b369e6726ed8a732c07db48f7209425c434bbe) [@bep](https://github.com/bep) 
* Remove .Site.Ref/RelRef [69fd1c60](https://github.com/gohugoio/hugo/commit/69fd1c60d8bcf6d1cea4bfea852f62df8891ee81) [@bep](https://github.com/bep) 
* Increase default timeout value to 30s [a8e9f838](https://github.com/gohugoio/hugo/commit/a8e9f8389a61471fa372c815b216511201b56823) [@bep](https://github.com/bep) [#6502](https://github.com/gohugoio/hugo/issues/6502)
* Add a benchmark [0cf85c07](https://github.com/gohugoio/hugo/commit/0cf85c071aba57de8c6567fba166ed8332d01bac) [@bep](https://github.com/bep) 

### Other

* Add some internal template image tests [dcde8af8](https://github.com/gohugoio/hugo/commit/dcde8af8c6ab39eb34b5e1d6030d1aa2fe6923ca) [@bep](https://github.com/bep) [#6542](https://github.com/gohugoio/hugo/issues/6542)
* Update Goldmark [b0c7749f](https://github.com/gohugoio/hugo/commit/b0c7749fa1efca04839b767e1d48d617f3556867) [@bep](https://github.com/bep) 
* Use HUGO_ENV if set [5c5231e0](https://github.com/gohugoio/hugo/commit/5c5231e09e20953dc262df7d3b351a35f1c4b058) [@bep](https://github.com/bep) [#6456](https://github.com/gohugoio/hugo/issues/6456)
* Make the image cache more robust [d6f7a9e2](https://github.com/gohugoio/hugo/commit/d6f7a9e28dfd5abff08b6aaf6fb3493c46bd1e39) [@bep](https://github.com/bep) [#6501](https://github.com/gohugoio/hugo/issues/6501)
* Update to Go 1.13.4 and Go 1.12.13 [031f948f](https://github.com/gohugoio/hugo/commit/031f948f87ac97ca49d0a487a392a8a0c6afb699) [@bep](https://github.com/bep) 
* Restore -v behaviour [71597bd1](https://github.com/gohugoio/hugo/commit/71597bd1adfb016a3ea1977068c37dce92d49458) [@bep](https://github.com/bep) 
* Update Goldmark [82219128](https://github.com/gohugoio/hugo/commit/8221912869cf863d64ae7b50e0085589dc18e4d2) [@bep](https://github.com/bep) 
* Improve grammar in README.md [e1175ae8](https://github.com/gohugoio/hugo/commit/e1175ae83a365e0b17ec5904194e68ff3833e15a) [@jasdeepgill](https://github.com/jasdeepgill) 
* Replace the temp for with a dependency [a2d77f4a](https://github.com/gohugoio/hugo/commit/a2d77f4a803ce27802ea653a4aab53b89c37b488) [@bep](https://github.com/bep) 
* Update Chroma [b546417a](https://github.com/gohugoio/hugo/commit/b546417a27f8c59c8c7ccaebfef6bca03f5c4ac4) [@bep](https://github.com/bep) 
* Update Goldmark [4175b046](https://github.com/gohugoio/hugo/commit/4175b0468680b076a5e5f90450157a98f841789b) [@bep](https://github.com/bep) 
* markup/tableofcontents: GoDoc etc. [55f951cb](https://github.com/gohugoio/hugo/commit/55f951cbba69c29daabca57eeff5661d132fa162) [@bep](https://github.com/bep) 
* Minor cleanups [20f351ee](https://github.com/gohugoio/hugo/commit/20f351ee4cd40b3b53e33805fc6226c837290ed7) [@moorereason](https://github.com/moorereason) 
* Add Goldmark as the new default markdown handler [bfb9613a](https://github.com/gohugoio/hugo/commit/bfb9613a14ab2d93a4474e5486d22e52a9d5e2b3) [@bep](https://github.com/bep) [#5963](https://github.com/gohugoio/hugo/issues/5963)[#1778](https://github.com/gohugoio/hugo/issues/1778)[#6355](https://github.com/gohugoio/hugo/issues/6355)
* Add parallel task executor helper [628efd6e](https://github.com/gohugoio/hugo/commit/628efd6e293d27984a3f5ba33522f8edd19d69d6) [@bep](https://github.com/bep) 
* Update homepage.md [14a985f8](https://github.com/gohugoio/hugo/commit/14a985f8abc527d4e8487fcd5fa742e1ab2a00ed) [@bep](https://github.com/bep) 
* Do not check for remote modules if main project is vendored [20ec9fa2](https://github.com/gohugoio/hugo/commit/20ec9fa2bbd69dc47dfc9f1db40c954e08520071) [@bep](https://github.com/bep) [#6506](https://github.com/gohugoio/hugo/issues/6506)
* Add hint when dir not empty [1a36ce9b](https://github.com/gohugoio/hugo/commit/1a36ce9b0903e02a5068aed5f807ed9d21f48ece) [@YaguraStation](https://github.com/YaguraStation) [#4825](https://github.com/gohugoio/hugo/issues/4825)
* Headless bundles should not be listed in .Pages [d1d1f240](https://github.com/gohugoio/hugo/commit/d1d1f240a25945b37eebe8a9a3f439f290832b33) [@bep](https://github.com/bep) [#6492](https://github.com/gohugoio/hugo/issues/6492)
* Support Go time format strings in permalinks [70a1aa34](https://github.com/gohugoio/hugo/commit/70a1aa345b95bcf325f19c6e7184bcd6f885e454) [@look](https://github.com/look) 
* Increase timeout to 30000 for mage -v check [cafecca4](https://github.com/gohugoio/hugo/commit/cafecca440e495ec915cc6290fe09d2a343e9c95) [@anthonyfok](https://github.com/anthonyfok) 
* Prepare for Goldmark [5f6b6ec6](https://github.com/gohugoio/hugo/commit/5f6b6ec68936ebbbf590894c02a1a3ecad30735f) [@bep](https://github.com/bep) [#5963](https://github.com/gohugoio/hugo/issues/5963)
* Update quicktest [366ee4d8](https://github.com/gohugoio/hugo/commit/366ee4d8da1c2b0c1751e9bf6d54638439735296) [@bep](https://github.com/bep) 
* Use pointer receiver for ContentSpec [9abd3967](https://github.com/gohugoio/hugo/commit/9abd396789d007193145db9246d5daf1640bbb8a) [@bep](https://github.com/bep) 
* Allow arm64 to fail [ad4c56b5](https://github.com/gohugoio/hugo/commit/ad4c56b5512226e74fb4ed6f10630d26d93e9eb6) [@bep](https://github.com/bep) 
* Add a JSON roundtrip test [3717db1f](https://github.com/gohugoio/hugo/commit/3717db1f90797f4e2a5d546472fb6b6df072d435) [@bep](https://github.com/bep) [#6472](https://github.com/gohugoio/hugo/issues/6472)
* Update .travis.yml for arm64 support, etc. [ae4fde08](https://github.com/gohugoio/hugo/commit/ae4fde0866b2a10f0a414e0d76c4ff09bed3776e) [@anthonyfok](https://github.com/anthonyfok) 
* Skip Test386 on non-AMD64 architectures [c6d69d0c](https://github.com/gohugoio/hugo/commit/c6d69d0c95c42915956c210dbac8b884682d4a3e) [@anthonyfok](https://github.com/anthonyfok) 
* Switch to mage builds, various optimizations [ed268232](https://github.com/gohugoio/hugo/commit/ed2682325aeb8fd1c8139077d14a5f6906757a4e) [@jakejarvis](https://github.com/jakejarvis) 
* Add exception for new test image [66fe68ff](https://github.com/gohugoio/hugo/commit/66fe68ffc98974936e157b18cf6bd9266ee081a4) [@anthonyfok](https://github.com/anthonyfok) [#6439](https://github.com/gohugoio/hugo/issues/6439)
* Adjust benchmark templates [c5e1e824](https://github.com/gohugoio/hugo/commit/c5e1e8241a3b9f922f4a5134064ab2847174a959) [@bep](https://github.com/bep) 
* Update quicktest [3e8b5a5c](https://github.com/gohugoio/hugo/commit/3e8b5a5c0157fdcf93588a42fbc90b3cd898f6b1) [@bep](https://github.com/bep) 
* Do not attempt to build if there is no config file [e6aa6edb](https://github.com/gohugoio/hugo/commit/e6aa6edb4c5f37feb1f2bb8c0f3f80933c7adf5f) [@ollien](https://github.com/ollien) [#5896](https://github.com/gohugoio/hugo/issues/5896)

## Fixes

### Output

* Fix mage check on darwin and add debugging output [8beaa4c2](https://github.com/gohugoio/hugo/commit/8beaa4c25efb593d0363271000a3667b96567976) [@trimbo](https://github.com/trimbo) 

### Core

* Fix cascade in server mode [01766439](https://github.com/gohugoio/hugo/commit/01766439246add22a6e6d0c12f932610be55cd8a) [@bep](https://github.com/bep) [#6538](https://github.com/gohugoio/hugo/issues/6538)
* Fix .Sections vs siblings [da535235](https://github.com/gohugoio/hugo/commit/da53523599b43261520a22d77019b390aaa072e7) [@bep](https://github.com/bep) [#6365](https://github.com/gohugoio/hugo/issues/6365)
* Fix recently broken timeout config [e3451371](https://github.com/gohugoio/hugo/commit/e3451371bdb68015f89c8c0f7d8ea0a19fff8df5) [@bep](https://github.com/bep) 
* Fix emoji handling inside shortcodes [812688fc](https://github.com/gohugoio/hugo/commit/812688fc2f3e220ac35cad9f0445a2548f0cc603) [@bep](https://github.com/bep) [#6504](https://github.com/gohugoio/hugo/issues/6504)
* Fix ref/relref anhcor handling [c26d00db](https://github.com/gohugoio/hugo/commit/c26d00db648a4b475d94c9ed8e21dafb6efa1776) [@bep](https://github.com/bep) [#6481](https://github.com/gohugoio/hugo/issues/6481)

### Other

* Fix language handling in ExecuteAsTemplate [96f09659](https://github.com/gohugoio/hugo/commit/96f09659ce8752c32a2a6429c9faf23be4faa091) [@bep](https://github.com/bep) [#6331](https://github.com/gohugoio/hugo/issues/6331)
* Fix potential data race [03e2d746](https://github.com/gohugoio/hugo/commit/03e2d7462dec17c2f623a13db709f9efc88182af) [@bep](https://github.com/bep) [#6478](https://github.com/gohugoio/hugo/issues/6478)
* Fix jekyll metadata import on individual posts [8a89b858](https://github.com/gohugoio/hugo/commit/8a89b8582f0f681dc28961adb05ab0bf66da9543) [@trimbo](https://github.com/trimbo) [#5576](https://github.com/gohugoio/hugo/issues/5576)
* Fix Params case handling in the index, sort and where func [a3fe5e5e](https://github.com/gohugoio/hugo/commit/a3fe5e5e35f311f22b6b4fc38abfcf64cd2c7d6f) [@bep](https://github.com/bep) 
* Fix GetPage Params case issue [cd07e6d5](https://github.com/gohugoio/hugo/commit/cd07e6d57b158a76f812e8c4c9567dbc84f57939) [@bep](https://github.com/bep) [#5946](https://github.com/gohugoio/hugo/issues/5946)
* Update to Chroma v0.6.9 for Java lexer fix [8483b53a](https://github.com/gohugoio/hugo/commit/8483b53aefc3c6b52f9917e6e5af9c4d2e98df66) [@anthonyfok](https://github.com/anthonyfok) [#6476](https://github.com/gohugoio/hugo/issues/6476)
* Update past go-cmp's checkptr fix [c3d433af](https://github.com/gohugoio/hugo/commit/c3d433af56071d42aeb3f85854bd30db64ed70b8) [@anthonyfok](https://github.com/anthonyfok) 
* Fix crash in multilingual content fs [33c474b9](https://github.com/gohugoio/hugo/commit/33c474b9b3bd470670740f30c5131071ce906b22) [@bep](https://github.com/bep) [#6463](https://github.com/gohugoio/hugo/issues/6463)
* Update to Chroma v0.6.8 to fix a crash [baa97508](https://github.com/gohugoio/hugo/commit/baa975082c6809c8a02a8109ec3062a2b7d48344) [@bep](https://github.com/bep) [#6450](https://github.com/gohugoio/hugo/issues/6450)





