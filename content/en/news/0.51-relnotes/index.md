
---
date: 2018-11-07
title: "Hugo 0.51: The 30K Stars Edition!"
description: "Bug fixes, new template functions and more error improvements."
categories: ["Releases"]
---

Hugo reached [30 000 stars on GitHub](https://github.com/gohugoio/hugo/stargazers) this week, which is a good occasion to do a follow-up release of the great Hugo `0.50`. This is mostly a bug fix release, but it also adds some useful new functionality, two examples are the new template funcs [complement](https://gohugo.io/functions/complement/) and [symdiff](https://gohugo.io/functions/symdiff/). This release also continues the work on improving Hugo's error messages. And with `.Position` now available on shortcodes, you can also improve your own error messages inside your custom shortcodes:


```bash
{{ with .Get "name" }}
{{ else }}
{{ errorf "missing value for param 'name': %s" .Position }}
{{ end }}
```

When the above fails, you will see an `ERROR` log similar to the below:

```bash
ERROR 2018/11/07 10:05:55 missing value for param name: "/sites/hugoDocs/content/en/variables/shortcodes.md:32:1"
```
 
This release represents **31 contributions by 5 contributors** to the main Hugo code base.
[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@krisbudhram](https://github.com/krisbudhram), [@LorenzCK](https://github.com/LorenzCK), and [@coliff](https://github.com/coliff) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **6 contributions by 5 contributors**. A special thanks to [@ikemo3](https://github.com/ikemo3), [@maiki](https://github.com/maiki), [@morya](https://github.com/morya), and [@regisphilibert](https://github.com/regisphilibert) for their work on the documentation site.


Hugo now has:

* 30095+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 441+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 276+ [themes](http://themes.gohugo.io/)


## Notes

* Remove deprecated useModTimeAsFallback [0bc4b024](https://github.com/gohugoio/hugo/commit/0bc4b0246dd6b7d71f8676a52644077a4f70ec8f) [@bep](https://github.com/bep) 
* Bump to ERROR for the deprecated Pages.Sort [faeb55c1](https://github.com/gohugoio/hugo/commit/faeb55c1d827f0ea994551a103ff4f7448786d39) [@bep](https://github.com/bep) 
* Deprecate .Site.Ref and .Site.RelRef [6c6a6c87](https://github.com/gohugoio/hugo/commit/6c6a6c87ec2b5ac7342e268ab47861429230f7f4) [@bep](https://github.com/bep) [#5386](https://github.com/gohugoio/hugo/issues/5386)

## Enhancements

### Templates

* Properly handle pointer types in complement/symdiff [79a06aa4](https://github.com/gohugoio/hugo/commit/79a06aa4b64b526c242dfa41f2c7bc24e1352d5b) [@bep](https://github.com/bep) 
* Add collections.SymDiff [488776b6](https://github.com/gohugoio/hugo/commit/488776b6498d1377718133d42daa87ce1236215d) [@bep](https://github.com/bep) [#5410](https://github.com/gohugoio/hugo/issues/5410)
* Add collections.Complement [42d8dfc8](https://github.com/gohugoio/hugo/commit/42d8dfc8c88af03ea926a59bc2332acc70cca5f6) [@bep](https://github.com/bep) [#5400](https://github.com/gohugoio/hugo/issues/5400)

### Core

* Improve error message on duplicate menu items [3a44920e](https://github.com/gohugoio/hugo/commit/3a44920e79ef86003555d8a4860c29257b2914f0) [@bep](https://github.com/bep) 
* Add .Position to shortcode [33a7b36f](https://github.com/gohugoio/hugo/commit/33a7b36fd42ee31dd79115ec6639bed24247332f) [@bep](https://github.com/bep) [#5371](https://github.com/gohugoio/hugo/issues/5371)

### Other

* Document shortcode error handling [e456e34b](https://github.com/gohugoio/hugo/commit/e456e34bdbde058243eb0a5d3c0017748639e08e) [@bep](https://github.com/bep) 
* Document symdiff [5d14d04a](https://github.com/gohugoio/hugo/commit/5d14d04ac678ad24e4946ed2a581ab71b3834def) [@bep](https://github.com/bep) 
* Document complement [ddcb4028](https://github.com/gohugoio/hugo/commit/ddcb402859b50193bfd6d8b752b568d26d14f603) [@bep](https://github.com/bep) 
* Update minify [d212f609](https://github.com/gohugoio/hugo/commit/d212f60949b6afefbe5aa79394f98dbddf7be068) [@bep](https://github.com/bep) 
* Re-generate CLI docs [2998fa0c](https://github.com/gohugoio/hugo/commit/2998fa0cd5bad161b9c802d2409d8c9c81155011) [@bep](https://github.com/bep) 
* Add --minify to hugo server [5b1edd28](https://github.com/gohugoio/hugo/commit/5b1edd281a493bdb27af4dc3c8fae7e10dd54830) [@bep](https://github.com/bep) 
* Make WARN the new default log log level [4b7d3e57](https://github.com/gohugoio/hugo/commit/4b7d3e57a40214a1269eda59731aa22a8f4463dd) [@bep](https://github.com/bep) [#5203](https://github.com/gohugoio/hugo/issues/5203)
* Regenerate the docs helper [486bc46a](https://github.com/gohugoio/hugo/commit/486bc46a5217a9d70fe0d14ab9261d7b4eb026d6) [@bep](https://github.com/bep) 
* Skip watcher event files if matched in ignoreFiles [f8446188](https://github.com/gohugoio/hugo/commit/f8446188dbec8378f34f0fea39161a49fcc46083) [@krisbudhram](https://github.com/krisbudhram) 
* Update Chroma [d523aa4b](https://github.com/gohugoio/hugo/commit/d523aa4bb03e913f55c2f89544e6112e320c975a) [@bep](https://github.com/bep) [#5392](https://github.com/gohugoio/hugo/issues/5392)
* Add file (line/col) info to ref/relref errors [1d18eb05](https://github.com/gohugoio/hugo/commit/1d18eb0574a57c3e9f468659d076a666a3dd76f2) [@bep](https://github.com/bep) [#5371](https://github.com/gohugoio/hugo/issues/5371)
* Improve log color regexp [d3a98325](https://github.com/gohugoio/hugo/commit/d3a98325c31d7f02f0762e589a4986e55b2a0da2) [@bep](https://github.com/bep) 
* Correct minor typo (#5372) [e65268f2](https://github.com/gohugoio/hugo/commit/e65268f2c2dd5ac54681d3266564901d99ed3ea3) [@coliff](https://github.com/coliff) 

## Fixes

### Templates

* Fix the docshelper [61f210dd](https://github.com/gohugoio/hugo/commit/61f210dd7abe5de77c27dc6a6995a3ad5e77afa1) [@bep](https://github.com/bep) 
* Fix BOM issue in templates [3a786a24](https://github.com/gohugoio/hugo/commit/3a786a248d3eff6e732aa94e87d6e88196e5147a) [@bep](https://github.com/bep) [#4895](https://github.com/gohugoio/hugo/issues/4895)

### Output

* Fix ANSI character output regression on Windows [b8725f51](https://github.com/gohugoio/hugo/commit/b8725f5181f6a2709274a82c1c3fdfd8f2e3e28c) [@LorenzCK](https://github.com/LorenzCK) [#5377](https://github.com/gohugoio/hugo/issues/5377)

### Core

* Fix changing paginators in lazy render [b8b8436f](https://github.com/gohugoio/hugo/commit/b8b8436fcca17c152e94cae2a1acad32efc3946c) [@bep](https://github.com/bep) [#5406](https://github.com/gohugoio/hugo/issues/5406)
* Fix REF_NOT_FOUND logging to include page path [6180c85f](https://github.com/gohugoio/hugo/commit/6180c85fb8f95e01446b74c50cab3f0480305fe4) [@bep](https://github.com/bep) [#5371](https://github.com/gohugoio/hugo/issues/5371)
* Fix broken manual summary handling [b2a676f5](https://github.com/gohugoio/hugo/commit/b2a676f5f09a3eea360887b099b9d5fc25a88492) [@bep](https://github.com/bep) [#5381](https://github.com/gohugoio/hugo/issues/5381)
* Fix deadlock when content building times out [729593c8](https://github.com/gohugoio/hugo/commit/729593c842794eaf7127050953a5c2256d332051) [@bep](https://github.com/bep) [#5375](https://github.com/gohugoio/hugo/issues/5375)

### Other

* Fix spelling [47506d16](https://github.com/gohugoio/hugo/commit/47506d164467eb7ddbcada81b767d8df5f9c8786) [@qeesung](https://github.com/qeesung) 
* Fix shortcode directly following a shortcode delimiter [d16a7a33](https://github.com/gohugoio/hugo/commit/d16a7a33ff1f22b9fa357189a901a4f1de4e65e7) [@bep](https://github.com/bep) [#5402](https://github.com/gohugoio/hugo/issues/5402)
* Fix recently broken error template [2bd9d909](https://github.com/gohugoio/hugo/commit/2bd9d9099db267831731ed2d2200eb09305df9fc) [@bep](https://github.com/bep) 





