
---
date: 2020-11-27
title: "0.79.0"
description: "0.79.0"
categories: ["Releases"]
---

	Hugo `0.79.0` is a small, but useful release. You can now set custom `.Params` in your [menu](https://gohugo.io/content-management/menus/) configuration, and you can now also override deeply nested snake_cased configuration variables with [OS environment variables](https://gohugo.io/getting-started/configuration/#configure-with-environment-variables). Other than that we have refreshed all the core upstream dependencies. A special thanks to [@alecthomas](https://github.com/alecthomas) (some new [Chroma lexers](https://github.com/alecthomas/chroma/releases/tag/v0.8.2) and fixes) and [@evanw](https://github.com/evanw) ([ESBuild](https://github.com/evanw/esbuild)).

This release represents **33 contributions by 8 contributors** to the main Hugo code base. [@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout [@AdamKorcz](https://github.com/AdamKorcz), and [@davidejones](https://github.com/davidejones) for their ongoing contributions. And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour), [@coliff](https://github.com/coliff) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **13 contributions by 11 contributors**. A special thanks to [@Valac01](https://github.com/Valac01), [@bep](https://github.com/bep), [@mhansen](https://github.com/mhansen), and [@chanjarster](https://github.com/chanjarster) for their work on the documentation site.

Hugo now has:

* 48392+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 437+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 361+ [themes](http://themes.gohugo.io/)

## Enhancements

### Templates

* Add more layout lookup tests [34061706](https://github.com/gohugoio/hugo/commit/34061706e6a9631d92ae3d01e0458eee7bc251cc) [@moorereason](https://github.com/moorereason) [#7964](https://github.com/gohugoio/hugo/issues/7964)

### Other

* bump gopkg.in/yaml.v2 from 2.3.0 to 2.4.0 [17e0bbe8](https://github.com/gohugoio/hugo/commit/17e0bbe821b508cea936bcfd5c1c181bdb8ad70d) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Allow setting the delimiter used for setting config via OS env, e.g. HUGO_ [7e223b3b](https://github.com/gohugoio/hugo/commit/7e223b3baaef68d6e6f99e28f162362c81deffba) [@bep](https://github.com/bep) [#7829](https://github.com/gohugoio/hugo/issues/7829)
* Update to github.com/evanw/esbuild 0.8.11 to 0.8.14 [8a6e7060](https://github.com/gohugoio/hugo/commit/8a6e70605350255920100c5c085bb9ea6576d972) [@bep](https://github.com/bep) [#7986](https://github.com/gohugoio/hugo/issues/7986)
* bump github.com/google/go-cmp from 0.5.2 to 0.5.3 [6f7633df](https://github.com/gohugoio/hugo/commit/6f7633df7d2c06e32eac628f9c7809dfee75eeed) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Remove unneeded meta tag from blog example [a546059a](https://github.com/gohugoio/hugo/commit/a546059a9c0b4541f6c9e292f2bb065c1b6115d9) [@coliff](https://github.com/coliff) 
* bump github.com/getkin/kin-openapi from 0.30.0 to 0.31.0 [b5d906e3](https://github.com/gohugoio/hugo/commit/b5d906e31e716328e2c0fbbdbfe6fc5b2ff98886) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Regen docshelper [fd70bdaf](https://github.com/gohugoio/hugo/commit/fd70bdafe7dc5d18c9a2f22c49acc3a8de376e8e) [@bep](https://github.com/bep) 
* Add menu params [8f5c9a74](https://github.com/gohugoio/hugo/commit/8f5c9a747fcebb02bb99f5de272046411eb15370) [@davidejones](https://github.com/davidejones) [#7951](https://github.com/gohugoio/hugo/issues/7951)
* Preserve url set in frontmatter without sanitizing [e4fcb672](https://github.com/gohugoio/hugo/commit/e4fcb672ed8bae21fd9780292b54fea3040dd877) [@satotake](https://github.com/satotake) [#6007](https://github.com/gohugoio/hugo/issues/6007)
* Add file deleted by accident [18c13adc](https://github.com/gohugoio/hugo/commit/18c13adcd46bdff963311fdba9eaa9b5a299106e) [@bep](https://github.com/bep) [#7972](https://github.com/gohugoio/hugo/issues/7972)
* Regenerate docshelper" [20a35374](https://github.com/gohugoio/hugo/commit/20a35374a3c90adb32a90a5f671afb15165210be) [@bep](https://github.com/bep) [#7972](https://github.com/gohugoio/hugo/issues/7972)
* Regenerate docshelper [caf16c20](https://github.com/gohugoio/hugo/commit/caf16c20853947138883f6460682e19733895f52) [@bep](https://github.com/bep) 
* Update to Chroma v0.8.2 [b298c06e](https://github.com/gohugoio/hugo/commit/b298c06e0551e3eba10b39ae5668b7a6a36a08a7) [@bep](https://github.com/bep) [#7970](https://github.com/gohugoio/hugo/issues/7970)
* bump github.com/evanw/esbuild from 0.8.8 to 0.8.11 [55e290af](https://github.com/gohugoio/hugo/commit/55e290af41ad1c92af13679d4a84d64985d41456) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/getkin/kin-openapi from 0.26.0 to 0.30.0 [506a190a](https://github.com/gohugoio/hugo/commit/506a190a82cc5564012a1228b4179637b64e58eb) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/evanw/esbuild from 0.8.6 to 0.8.8 [fc81de64](https://github.com/gohugoio/hugo/commit/fc81de643934e84bb1e1392f6200559ee0ada9b6) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Let ESBuild handle all imports from node_modules [78f227b6](https://github.com/gohugoio/hugo/commit/78f227b664d86c30fbb25f7a953b7ef8f2dacf13) [@bep](https://github.com/bep) [#7948](https://github.com/gohugoio/hugo/issues/7948)
* bump github.com/evanw/esbuild from 0.8.5 to 0.8.6 [5e03f644](https://github.com/gohugoio/hugo/commit/5e03f644a4507f51bdbcdb42b65ce4e99095374f) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/evanw/esbuild from 0.8.4 to 0.8.5 [a92ef20f](https://github.com/gohugoio/hugo/commit/a92ef20ff6e43ba05844539b60782e8190712cdc) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/getkin/kin-openapi from 0.22.1 to 0.26.0 [0d54a844](https://github.com/gohugoio/hugo/commit/0d54a844061e808dd5b4ff4874b2e4bd9df4d556) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Update GH docs to say "main" as default branch [943f3c93](https://github.com/gohugoio/hugo/commit/943f3c932f5f67ab52bf8e0636e57751dc9b1891) [@maco](https://github.com/maco) 
* Updated year in header [4f20bf29](https://github.com/gohugoio/hugo/commit/4f20bf29eb246a2e65508175fdd5f25b44e98370) [@AdamKorcz](https://github.com/AdamKorcz) 
* Added first fuzzer [4c613d5d](https://github.com/gohugoio/hugo/commit/4c613d5d5d60b80a262e968ae8a4525eba8619a2) [@AdamKorcz](https://github.com/AdamKorcz) 
* bump github.com/frankban/quicktest from 1.11.1 to 1.11.2 [82a182e5](https://github.com/gohugoio/hugo/commit/82a182e52c4165b4f51d0cc8ef0f21df5d628c69) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump golang.org/x/text from 0.3.3 to 0.3.4 [dfc662b2](https://github.com/gohugoio/hugo/commit/dfc662b2086430dde96c18ccb6b92bba4f1be428) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/evanw/esbuild from 0.8.3 to 0.8.4 [2f0917cc](https://github.com/gohugoio/hugo/commit/2f0917cc014557e201a9348664736d608a7fa131) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Disable NPM test on Travis on Windows [3437174c](https://github.com/gohugoio/hugo/commit/3437174c3a7b96925b82b351ac87530b4fa796a5) [@bep](https://github.com/bep) 
* Install nodejs on Windows [f66302ca](https://github.com/gohugoio/hugo/commit/f66302ca0579171ffd1730eb8f33dd05af3d9a00) [@bep](https://github.com/bep) 
* Remove external source map option [944150ba](https://github.com/gohugoio/hugo/commit/944150bafbbb5c3e807ba3688174e70764dbdc64) [@bep](https://github.com/bep) [#7932](https://github.com/gohugoio/hugo/issues/7932)

## Fixes

### Other

* Fix memory usage in writeStats [d162bbd7](https://github.com/gohugoio/hugo/commit/d162bbd7990b6a523bdadcd10bf60fcb43ecf270) [@bep](https://github.com/bep) [#7945](https://github.com/gohugoio/hugo/issues/7945)
* Fix server rebuild issue with partials referenced from render hooks [e442cf30](https://github.com/gohugoio/hugo/commit/e442cf30a215e33b49ce588a9098147282bd883f) [@bep](https://github.com/bep) [#7990](https://github.com/gohugoio/hugo/issues/7990)
* Misc fixes [bf2837a3](https://github.com/gohugoio/hugo/commit/bf2837a314eaf70135791984a423b0b09f58741d) [@bep](https://github.com/bep) [#7924](https://github.com/gohugoio/hugo/issues/7924)[#7923](https://github.com/gohugoio/hugo/issues/7923)





