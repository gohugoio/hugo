
---
date: 2019-07-25
title: "Hugo 0.56.0: Hugo Modules and Deployment"
description: "This release adds powerful module support with dependency management for all component types, including content. And we now have built-in deployment support to GCS, S3, or Azure."
categories: ["Releases"]
---

	
**Hugo 0.56.0** is filled with improvements, but there are two main headliners: **Hugo Modules** and **Hugo Deploy**.  

**Hugo Deploy** is implemented by [@vangent](https://github.com/vangent) and brings built-in deployment support for GCS, S3, or Azure using the Hugo CLI. See the [Hugo Deploy Documentation](https://gohugo.io/hosting-and-deployment/hugo-deploy/) for more information.

**Hugo Modules** is very much a community effort on the design and specification side, but [@bep](https://github.com/bep) has driven the implementation. Some notes about what all of this is about:

* A new `module` configuration section where you can import almost anything. You can configure both your own file mounts and the file mounts of the modules you import. This is the new recommended way of configuring what you earlier put in `configDir`, `staticDir` etc. And it also allows you to mount folders in non-Hugo-projects, e.g. the `SCSS` folder in the Bootstrap GitHub project.
* A module consists of a set of mounts to the standard 7 component types in Hugo: `static`, `content`, `layouts`, `data`, `assets`, `i18n`, and `archetypes`. Yes, Theme Components can now include content, which should be very useful, especially in bigger multilingual projects.
* Modules not in your local file cache will be downloaded automatically and even "hot replaced" while the server is running.
* Hugo Modules supports and encourages semver versioned modules, and uses the minimal version selection algorithm to resolve versions.
* A new set of CLI commands are provided to manage all of this: `hugo mod init`,  `hugo mod get`,  `hugo mod graph`,  `hugo mod tidy`, and  `hugo mod vendor`.

**Hugo Modules is powered by Go Modules.**

This is all very much brand new and there are only a few example projects around:

* https://github.com/bep/docuapi is a theme that has been ported to **Hugo Modules** while testing this feature. It is a good example of a non-Hugo-project mounted into Hugo's folder structure. It even shows a JS Bundler implementation in regular Go templates.
* https://github.com/bep/my-modular-site is a very simple site used for testing.

See the [Hugo Modules Documentation](https://gohugo.io/hugo-modules/) for more information.

This release represents **104 contributions by 19 contributors** to the main Hugo code base.
[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@vangent](https://github.com/vangent), [@niklasfasching](https://github.com/niklasfasching), and [@coliff](https://github.com/coliff) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **48 contributions by 37 contributors**. A special thanks to [@davidsneighbour](https://github.com/davidsneighbour), [@bep](https://github.com/bep), [@BCNelson](https://github.com/BCNelson), and [@coliff](https://github.com/coliff) for their work on the documentation site.


Hugo now has:

* 36902+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 440+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 328+ [themes](http://themes.gohugo.io/)


## Notes

* We  have removed the "auto theme namespacing" of  params from theme configuration. This was an undocumented and hidden feature that wasn't useful in practice.
* We have revised and improved the symlinks support in Hugo: In earlier versions, symlinks were only fully supported for the content folders. With the introduction of the new very flexible file mounts, with content support even for what we have traditionally named "themes", we needed a more precise definition of symlink support in Hugo:

  * Symlinks are not supported outside of the main project ((the project you run `hugo` or `hugo server` from).
  * In the main project `static` mounts, only symlinks to files are supported.
  * In all other mounts in the main project, both file and directory symlinks are allowed. 

## Enhancements

### Templates

* Add Merge function [c624a779](https://github.com/gohugoio/hugo/commit/c624a77992c4f7c1bcb5f659e89332d8588986b7) [@bep](https://github.com/bep) [#5992](https://github.com/gohugoio/hugo/issues/5992)
* Regenerate templates [b2a3d464](https://github.com/gohugoio/hugo/commit/b2a3d4644bb5a505db662b2927af6f80856a3076) [@bep](https://github.com/bep) 
* Unwrap any interface value in sort and where [8d898ad6](https://github.com/gohugoio/hugo/commit/8d898ad6672e0ccb62c5a29b6fccab24d980f104) [@bep](https://github.com/bep) [#5989](https://github.com/gohugoio/hugo/issues/5989)
* Convert numeric values to float64 and compare them [fb007e9a](https://github.com/gohugoio/hugo/commit/fb007e9ae56f295abe9835485f98dcf3cc362420) [@tryzniak](https://github.com/tryzniak) [#5685](https://github.com/gohugoio/hugo/issues/5685)
* Provide more detailed errors in Where [f76e5011](https://github.com/gohugoio/hugo/commit/f76e50118b8b5dd8989d068db35222bfa0a242d8) [@moorereason](https://github.com/moorereason) 
* Return error on invalid input in in [7fbfedf0](https://github.com/gohugoio/hugo/commit/7fbfedf01367ff076c3c875b183789b769b99241) [@bep](https://github.com/bep) [#5875](https://github.com/gohugoio/hugo/issues/5875)
* Make Pages etc. work with the in func [06f56fc9](https://github.com/gohugoio/hugo/commit/06f56fc983d460506d39b3a6f638b1632af07073) [@bep](https://github.com/bep) [#5875](https://github.com/gohugoio/hugo/issues/5875)
* Make Pages etc. work in uniq [d7a67dcb](https://github.com/gohugoio/hugo/commit/d7a67dcb51829b12d492d3f2ee4f6e2a3834da63) [@bep](https://github.com/bep) [#5852](https://github.com/gohugoio/hugo/issues/5852)
* Handle late transformation of templates [2957795f](https://github.com/gohugoio/hugo/commit/2957795f5276cc9bc8d438da2d7d9b61defea225) [@bep](https://github.com/bep) [#5865](https://github.com/gohugoio/hugo/issues/5865)

### Output

* Remove comma after URL in new site output [de7b9475](https://github.com/gohugoio/hugo/commit/de7b9475c049e5df5b076d80799ebcbee3eb84c8) [@snorfalorpagus](https://github.com/snorfalorpagus) 

### Core

* Add a symdiff test [072aa7f1](https://github.com/gohugoio/hugo/commit/072aa7f110ddc8a8b9dbc6d4cad3e5ba6c2ac4d0) [@bep](https://github.com/bep) 
* Add testfile to .gitignore [7611078d](https://github.com/gohugoio/hugo/commit/7611078daef32306ab31fe360db9895cdd3626d3) [@bep](https://github.com/bep) 
* Add another site benchmark [dc1d4a92](https://github.com/gohugoio/hugo/commit/dc1d4a9200c54e631775f34725140fd5009aa561) [@bep](https://github.com/bep) 
* Prevent parallel server rebuilds [95ce2a40](https://github.com/gohugoio/hugo/commit/95ce2a40e734bb82b69f9a64270faf3ed69c92cc) [@bep](https://github.com/bep) [#5885](https://github.com/gohugoio/hugo/issues/5885)[#5968](https://github.com/gohugoio/hugo/issues/5968)
* Disable racy test [ad5703a9](https://github.com/gohugoio/hugo/commit/ad5703a91712cd245231ba8fdbc49632c794a165) [@bep](https://github.com/bep) [#5926](https://github.com/gohugoio/hugo/issues/5926)
* Avoid recloning of shortcode templates [69a56420](https://github.com/gohugoio/hugo/commit/69a56420aec5bf5abb846701d4a5ec67fe060d96) [@bep](https://github.com/bep) [#5890](https://github.com/gohugoio/hugo/issues/5890)
* No links for bundled pages [0775c98e](https://github.com/gohugoio/hugo/commit/0775c98e6c5b700e46adaaf190fc3f693a6ab002) [@bep](https://github.com/bep) [#5882](https://github.com/gohugoio/hugo/issues/5882)
* Add some OutputFormats.Get tests [7aeeb60d](https://github.com/gohugoio/hugo/commit/7aeeb60d7ee71690461df92ff41cb8b2f7f5aa61) [@bep](https://github.com/bep) [#5877](https://github.com/gohugoio/hugo/issues/5877)
* Add some integration tests for in/uniq using Pages [6c80acbd](https://github.com/gohugoio/hugo/commit/6c80acbd5e314dd92fc075551ffabafaae01dca7) [@bep](https://github.com/bep) [#5875](https://github.com/gohugoio/hugo/issues/5875)[#5852](https://github.com/gohugoio/hugo/issues/5852)
* Add more tests for Permalinkable [35f41834](https://github.com/gohugoio/hugo/commit/35f41834ea3a8799b9b7eda360cf8d30b1b727ba) [@bep](https://github.com/bep) [#5849](https://github.com/gohugoio/hugo/issues/5849)
* Add a test for parent's resources in shortcode [8d7607ae](https://github.com/gohugoio/hugo/commit/8d7607aed10b3fe7373126ff5fa7dae36c559d7f) [@bep](https://github.com/bep) [#5833](https://github.com/gohugoio/hugo/issues/5833)

### Other

* Add Hugo Modules docs [77bf2991](https://github.com/gohugoio/hugo/commit/77bf2991b1690bcde8c6570cd4c12f2041d93096) [@bep](https://github.com/bep) 
* Block symlink dir traversal for /static [e5f22997](https://github.com/gohugoio/hugo/commit/e5f229974166402f51e4ee0695ffb4d1e09fa174) [@bep](https://github.com/bep) 
* modules: Gofmt [87a07282](https://github.com/gohugoio/hugo/commit/87a07282a2f01779e098cde0aaee1bae34dc32e6) [@bep](https://github.com/bep) 
* Rename disabled => disable in config [882d678b](https://github.com/gohugoio/hugo/commit/882d678bbf2a149a90e2aed4341d7f6fc2cb394d) [@bep](https://github.com/bep) 
* Regenerate CLI docs [215d2ed8](https://github.com/gohugoio/hugo/commit/215d2ed883d5adbde1d119d33e2f2e88c8435f41) [@bep](https://github.com/bep) 
* Regenerate data helpers [23adc0c2](https://github.com/gohugoio/hugo/commit/23adc0c2d96d33b426836b09683d2440d6186728) [@bep](https://github.com/bep) 
* Add Hugo Modules [9f5a9207](https://github.com/gohugoio/hugo/commit/9f5a92078a3f388b52d597b5a59af5c933a112d2) [@bep](https://github.com/bep) [#5973](https://github.com/gohugoio/hugo/issues/5973)[#5996](https://github.com/gohugoio/hugo/issues/5996)[#6010](https://github.com/gohugoio/hugo/issues/6010)[#5911](https://github.com/gohugoio/hugo/issues/5911)[#5940](https://github.com/gohugoio/hugo/issues/5940)[#6074](https://github.com/gohugoio/hugo/issues/6074)[#6082](https://github.com/gohugoio/hugo/issues/6082)[#6092](https://github.com/gohugoio/hugo/issues/6092)
* Tidy [47953148](https://github.com/gohugoio/hugo/commit/47953148b6121441d0147c960a99829c53b5a5ba) [@bep](https://github.com/bep) 
* Update xerrors [ac101aba](https://github.com/gohugoio/hugo/commit/ac101aba4fde02dc8273a3f06a3b4267ca501a3d) [@bep](https://github.com/bep) 
* Ignore errors in go mod download [58a47ccd](https://github.com/gohugoio/hugo/commit/58a47ccde00e2425364eaa5c2123e0718a2ee3f2) [@bep](https://github.com/bep) 
* Update Chroma [95b1d301](https://github.com/gohugoio/hugo/commit/95b1d3013b4717f8b02093a99d1d0c4a6a1ca929) [@bep](https://github.com/bep) [#6088](https://github.com/gohugoio/hugo/issues/6088)
* Change postcss to check for local installation under node_modules/.bin [a5604e18](https://github.com/gohugoio/hugo/commit/a5604e18b0cb260b7748949b12c05814ced50815) [@ericselin](https://github.com/ericselin) [#5091](https://github.com/gohugoio/hugo/issues/5091)
* Add org to front matter formats [020086cb](https://github.com/gohugoio/hugo/commit/020086cb2b0396909d831abf66b8f1455e6f7e6a) [@niklasfasching](https://github.com/niklasfasching) 
* Update go-org [8524baee](https://github.com/gohugoio/hugo/commit/8524baee167fa6a8684569e7acff225c16c301c7) [@niklasfasching](https://github.com/niklasfasching) 
* Pagination - do not render href if no next item [88c8a15b](https://github.com/gohugoio/hugo/commit/88c8a15be18a0bd1bd9b9cb175f7a68f2b9cd355) [@markmandel](https://github.com/markmandel) 
* Include path to source page in non-relative ref/relref warning [59c4bc52](https://github.com/gohugoio/hugo/commit/59c4bc52ed6b146baa6cca97d054004297ea098a) [@justinvp](https://github.com/justinvp) 
* Update link to prevent redirect [ff10aa52](https://github.com/gohugoio/hugo/commit/ff10aa522dd4d2741e8f4e50a4a40a0854232f23) [@coliff](https://github.com/coliff) 
* Update URLs [9f258d2b](https://github.com/gohugoio/hugo/commit/9f258d2b8f98a137c5c8c4586a3db3e3c53a0040) [@coliff](https://github.com/coliff) 
* Introduce '#+KEY[]:' array notation [fad183c4](https://github.com/gohugoio/hugo/commit/fad183c4ae55069be9246e64ab1c8b2f43d08d06) [@niklasfasching](https://github.com/niklasfasching) 
* Replace goorgeous with go-org [b6867bf8](https://github.com/gohugoio/hugo/commit/b6867bf8068fcaaddf1cb7478f4d52a9c1be1411) [@niklasfasching](https://github.com/niklasfasching) 
* Default --target to the first deployment target [9df57154](https://github.com/gohugoio/hugo/commit/9df57154ee3e3185d024bfe376101b404d8b7cc4) [@vangent](https://github.com/vangent) 
* Add safety barrier between concatenated javascript resources [35abce27](https://github.com/gohugoio/hugo/commit/35abce27cabee43cc751db55a75b927f26275833) [@vincent99](https://github.com/vincent99) 
* Update CLI doc for "long" form [8914fe7e](https://github.com/gohugoio/hugo/commit/8914fe7ed7e7e55e07be32564159310c90e2dbd4) [@vangent](https://github.com/vangent) 
* Drop dashes in http header matcher attributes; other changes from code review [b0f536fb](https://github.com/gohugoio/hugo/commit/b0f536fb276f63df0b0b8d92aeda80affb4b6a34) [@vangent](https://github.com/vangent) 
* Add documentation for "hugo deploy" [1384d77a](https://github.com/gohugoio/hugo/commit/1384d77a04d7027d719993c6f54f892b5e7974db) [@vangent](https://github.com/vangent) 
* remove TODO comment about subfolders; handled by GoCDK blob URLs [5e83f425](https://github.com/gohugoio/hugo/commit/5e83f4256279149879a8e88cb02679dd00e8da2b) [@vangent](https://github.com/vangent) 
* Update gocloud.dev to v0.15.0. [b376b268](https://github.com/gohugoio/hugo/commit/b376b2685a2e21961052a0946ab8a6acc076d4da) [@vangent](https://github.com/vangent) 
* Return nil when not found in resources.Get [4c560020](https://github.com/gohugoio/hugo/commit/4c560020bc0c50f8004873be8adf83698b7c095a) [@bep](https://github.com/bep) [#5999](https://github.com/gohugoio/hugo/issues/5999)
* Update Viper [d44d3ea8](https://github.com/gohugoio/hugo/commit/d44d3ea89865baa33170244cac72a7ce26419b15) [@bep](https://github.com/bep) [#5954](https://github.com/gohugoio/hugo/issues/5954)
* Remove references to Google+ [d1cf53f5](https://github.com/gohugoio/hugo/commit/d1cf53f5f4d71b9461e0fe58193b5a8418b572e7) [@brunoamaral](https://github.com/brunoamaral) 
* Update gitmap to get CommitDate field [811ee996](https://github.com/gohugoio/hugo/commit/811ee996a625b5eb3464a34a1623260b11c0bfd3) [@bep](https://github.com/bep) 
* Create new 'hugo list all' command [5b4b8bb3](https://github.com/gohugoio/hugo/commit/5b4b8bb3c1ecb30e7a38ed44eb795f1d972cd320) [@rusnasonov](https://github.com/rusnasonov) [#5904](https://github.com/gohugoio/hugo/issues/5904)
* Medium -> Hugo exporting tool [2278b0eb](https://github.com/gohugoio/hugo/commit/2278b0eb02ccdd3c2d4358d39074767d33fecb71) [@Ahirice](https://github.com/Ahirice) 
* Switch base image for final build [41974303](https://github.com/gohugoio/hugo/commit/41974303f3e5e6d7deb87a791aba512dbf1b9c59) [@brianredbeard](https://github.com/brianredbeard) [#5970](https://github.com/gohugoio/hugo/issues/5970)[#5056](https://github.com/gohugoio/hugo/issues/5056)
* Merge branch 'release-0.55.6' [6b3f1a10](https://github.com/gohugoio/hugo/commit/6b3f1a10028c81b776a5797bbd02c86957f8f042) [@bep](https://github.com/bep) 
* Release 0.55.6 [9b48c5d6](https://github.com/gohugoio/hugo/commit/9b48c5d6bd56741dac714324a6ae59f6374bccdc) [@bep](https://github.com/bep) 
* Update to Go 1.12.5 and Go 1.11.10 [71b8d8b6](https://github.com/gohugoio/hugo/commit/71b8d8b6a4127acacda8ada08cd61d7bfb18e962) [@bep](https://github.com/bep) [#5944](https://github.com/gohugoio/hugo/issues/5944)
* Support configuration of upload order [527cf1ab](https://github.com/gohugoio/hugo/commit/527cf1ab03fe4746885e90a197bc25decad88f89) [@vangent](https://github.com/vangent) 
* Support invalidating a CloudFront CDN cache [f4956d9a](https://github.com/gohugoio/hugo/commit/f4956d9aae69b1cb5715114cf5242fd80a9cabc7) [@vangent](https://github.com/vangent) 
* Move the package below /langs [2838d58b](https://github.com/gohugoio/hugo/commit/2838d58b1daa0f6a337125c5a64d06215901c5d6) [@bep](https://github.com/bep) 
* compute MD5 by reading if List didn't provide one [f330e869](https://github.com/gohugoio/hugo/commit/f330e869e42dc782a48c045aea5d29a134e225cb) [@vangent](https://github.com/vangent) 
* Use proxy.golang.org [0091b1f8](https://github.com/gohugoio/hugo/commit/0091b1f893aba6a0060c392f58fcc0351ee0db66) [@bep](https://github.com/bep) 
* Add a "deploy" command [c7165589](https://github.com/gohugoio/hugo/commit/c7165589b3be5923f1275f0024818e6ae889d881) [@vangent](https://github.com/vangent) 
* Release 0.55.5 [e33ed29b](https://github.com/gohugoio/hugo/commit/e33ed29b754ac1982321e87e54a37c9bb59e53ca) [@bep](https://github.com/bep) 
* Regenerate docs helper [cee181c3](https://github.com/gohugoio/hugo/commit/cee181c3a67fe04b8e0c9f2807c5aa2871df474e) [@bep](https://github.com/bep) 
* Update blackfriday to v1.5.2 [1cbb501b](https://github.com/gohugoio/hugo/commit/1cbb501be8aa83b08865fbb6ad5aee254946712f) [@dbirks](https://github.com/dbirks) 
* Release 0.55.4 [1707f1a5](https://github.com/gohugoio/hugo/commit/1707f1a5f52b8036d675d6ce66fd844effcab9fe) [@bep](https://github.com/bep) 
* Avoid rebuilding the Translations map for every lookup [4756ec3c](https://github.com/gohugoio/hugo/commit/4756ec3cd8ef998f889619fe11be70cc900e2b75) [@bep](https://github.com/bep) [#5892](https://github.com/gohugoio/hugo/issues/5892)
* Init mem profile at the end [4c3c5120](https://github.com/gohugoio/hugo/commit/4c3c5120389cc95edc63b8f18a0eee786aa0c5e2) [@bep](https://github.com/bep) 
* Release 0.55.3 [c85b726f](https://github.com/gohugoio/hugo/commit/c85b726f8a3cca2f06f766e5982dc0023d0dd12c) [@bep](https://github.com/bep) 
* Regenerate docs helper [75b16e30](https://github.com/gohugoio/hugo/commit/75b16e30ec55e82a8024cc4d27880d9b79e0fa41) [@bep](https://github.com/bep) 
* Replace IsDraft with Draft in list command [3e421bd4](https://github.com/gohugoio/hugo/commit/3e421bd47cd35061df89c1c127ec8fa4ae368449) [@bep](https://github.com/bep) [#5873](https://github.com/gohugoio/hugo/issues/5873)
* Release 0.55.2 [fcd63a86](https://github.com/gohugoio/hugo/commit/fcd63a865e731077a0d508084460b6fe6ec82c33) [@bep](https://github.com/bep) 
* Release 0.55.1 [adb776b2](https://github.com/gohugoio/hugo/commit/adb776b22babc0d312ce7b83abbe6f5107c180d7) [@bep](https://github.com/bep) 
* Remove the space in `. RelPermalink` [7966c0b5](https://github.com/gohugoio/hugo/commit/7966c0b5b7b2297527f8be9040b793de5e4e3f48) [@yihui](https://github.com/yihui) 

## Fixes

### Templates

* Fix slice type handling in sort [e8a716b2](https://github.com/gohugoio/hugo/commit/e8a716b23a1ca78cf29460daacd4ba49bbc05ad1) [@bep](https://github.com/bep) [#6023](https://github.com/gohugoio/hugo/issues/6023)
* Fix internal templates usage of safeHTMLAttr [e22b3f54](https://github.com/gohugoio/hugo/commit/e22b3f54c3d8ce6567c21c63beab0b03cf7983ea) [@rhcarvalho](https://github.com/rhcarvalho) [#5236](https://github.com/gohugoio/hugo/issues/5236)[#5246](https://github.com/gohugoio/hugo/issues/5246)
* Fix nil compare in eq/ne for interface values [66b143a0](https://github.com/gohugoio/hugo/commit/66b143a01d1c192619839b732ce188923ab15d60) [@bep](https://github.com/bep) [#5905](https://github.com/gohugoio/hugo/issues/5905)
* Fix hugo package name and add godocs [4f93f8c6](https://github.com/gohugoio/hugo/commit/4f93f8c670b26258dc7e3a613c38dbc86d8eda76) [@moorereason](https://github.com/moorereason) 

### Output

* Fix permalink in sitemap etc. when multiple permalinkable output formats [6b76841b](https://github.com/gohugoio/hugo/commit/6b76841b052b97625b8995f326d758b89f5c2349) [@bep](https://github.com/bep) [#5910](https://github.com/gohugoio/hugo/issues/5910)
* Fix links for non-HTML output formats [c7dd66bf](https://github.com/gohugoio/hugo/commit/c7dd66bfe2e32430f9b1a3126c67014e40d8405e) [@bep](https://github.com/bep) [#5877](https://github.com/gohugoio/hugo/issues/5877)
* Fix menu URL when multiple permalinkable output formats [ea529c84](https://github.com/gohugoio/hugo/commit/ea529c847ebc0267c6d0426cc8f77d5c76c73fe4) [@bep](https://github.com/bep) [#5849](https://github.com/gohugoio/hugo/issues/5849)

### Core

* Fix broken test [fa28df10](https://github.com/gohugoio/hugo/commit/fa28df1058e0131364cea2e3ac7f80e934d024a1) [@bep](https://github.com/bep) 
* Fix bundle path when slug is set [3e6cb2cb](https://github.com/gohugoio/hugo/commit/3e6cb2cb77e16be5b6ddd4ae55d5fc6bfba2d226) [@bep](https://github.com/bep) [#4870](https://github.com/gohugoio/hugo/issues/4870)
* Fix PrevInSection/NextInSection for nested sections [bcbed4eb](https://github.com/gohugoio/hugo/commit/bcbed4ebdaf55b67abc521d69bba456c041a7e7d) [@bep](https://github.com/bep) [#5883](https://github.com/gohugoio/hugo/issues/5883)
* Fix shortcode version=1 logic [33c73811](https://github.com/gohugoio/hugo/commit/33c738116c26e2ac37f4bd48159e8e3197fd7b39) [@bep](https://github.com/bep) [#5831](https://github.com/gohugoio/hugo/issues/5831)
* Fix Pages reinitialization on rebuilds [9b17cbb6](https://github.com/gohugoio/hugo/commit/9b17cbb62a056ea7e26b1146cbf3ba42f5acf805) [@bep](https://github.com/bep) [#5833](https://github.com/gohugoio/hugo/issues/5833)
* Fix shortcode namespace issue [56550d1e](https://github.com/gohugoio/hugo/commit/56550d1e449f45ebee398ac8a9e3b9818b3ee60e) [@bep](https://github.com/bep) [#5863](https://github.com/gohugoio/hugo/issues/5863)
* Fix false WARNINGs in lang prefix check [7881b096](https://github.com/gohugoio/hugo/commit/7881b0965f8b83d03379e9ed102cd0c3bce297e2) [@bep](https://github.com/bep) [#5860](https://github.com/gohugoio/hugo/issues/5860)
* Fix bundle resource publishing when multiple output formats [49d0a826](https://github.com/gohugoio/hugo/commit/49d0a82641581aa7dd66b9d5e8c7d75e23260083) [@bep](https://github.com/bep) [#5858](https://github.com/gohugoio/hugo/issues/5858)
* Fix panic for unused taxonomy content files [b799b12f](https://github.com/gohugoio/hugo/commit/b799b12f4a693dfeae8a5a362f131081a727bb8f) [@bep](https://github.com/bep) [#5847](https://github.com/gohugoio/hugo/issues/5847)
* Fix dates for sections with dates in front matter [70148672](https://github.com/gohugoio/hugo/commit/701486728e21bc0c6c78c2a8edb988abdf6116c7) [@bep](https://github.com/bep) [#5854](https://github.com/gohugoio/hugo/issues/5854)
* Fix simple menu config [9e9a1f92](https://github.com/gohugoio/hugo/commit/9e9a1f92baf151f8d840d6b5b963945d1410ce25) [@bep](https://github.com/bep) 

### Other

* Fix test on Windows [e5b6e208](https://github.com/gohugoio/hugo/commit/e5b6e2085aba74767ace269cd5f8a746230b4fa4) [@bep](https://github.com/bep) 
* Fix livereload for @import case [2fc0abd2](https://github.com/gohugoio/hugo/commit/2fc0abd22a37d90b6f1032eef46191a7bddf41bd) [@bep](https://github.com/bep) [#6106](https://github.com/gohugoio/hugo/issues/6106)
* Fix typo s/Meny/Menu/ [90b0127f](https://github.com/gohugoio/hugo/commit/90b0127f63e9cd5bf3a8bd4282237db224a3c263) [@kaushalmodi](https://github.com/kaushalmodi) 
* Add tests; fix Windows [5dc6d0df](https://github.com/gohugoio/hugo/commit/5dc6d0df94076e116934c83b837e2dd416efa784) [@vangent](https://github.com/vangent) 
* Fix concurrent initialization order [009076e5](https://github.com/gohugoio/hugo/commit/009076e5ee88fc46c95a9afd34f82f9386aa282a) [@bep](https://github.com/bep) [#5901](https://github.com/gohugoio/hugo/issues/5901)
* Fix WeightedPages in union etc. [f2795d4d](https://github.com/gohugoio/hugo/commit/f2795d4d2cef30170af43327f3ff7114923833b1) [@bep](https://github.com/bep) [#5850](https://github.com/gohugoio/hugo/issues/5850)
* Fix [4d425a86](https://github.com/gohugoio/hugo/commit/4d425a86f5c03a5cca27d4e0f99d61acbb938d80) [@bep](https://github.com/bep) 
* Fix paginator refresh on server change [f7375c49](https://github.com/gohugoio/hugo/commit/f7375c497239115cd30ae42af6b4d298e4e7ad7d) [@bep](https://github.com/bep) [#5838](https://github.com/gohugoio/hugo/issues/5838)
* Fix .RSSLinke deprecation message [3b86b4a9](https://github.com/gohugoio/hugo/commit/3b86b4a9f5ce010c9714d813d5b8ecddda22c69f) [@bep](https://github.com/bep) [#4427](https://github.com/gohugoio/hugo/issues/4427)





