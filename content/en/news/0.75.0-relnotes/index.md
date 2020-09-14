
---
date: 2020-09-14
title: "NPM Pack"
description: "Hugo 0.75 comes with a new \"hugo mod npm pack\" command, several improvements re. Hugo Modules and the Node tools, and more."
categories: ["Releases"]
---

Hugo `0.75.0` brings several improvements to Hugo Modules, a new CLI command to bridge the JavaScript dependencies into Hugo, a refresh of the versions of the most important upstream dependencies, and more. There are also some good bug fixes in this release. One notable one is covered by [this commit](https://github.com/gohugoio/hugo/commit/4055c121847847d8bd6b95a928185daee065091b) -- which covers a "stale content scenario in server" when you include content or page data via `GetPage` from a shortcode.

## NPM Pack

The new CLI command is called `hugo mod npm pack`.  We have marked it as experimental. It works great, go ahead and use it, but we need to test this out in real projects to get a feel of it; it is likely that it will change/improve in the upcoming versions of Hugo. The command creates a consolidated `package.json` from the project and all of its [theme components](https://gohugo.io/hugo-modules/theme-components/). On version conflicts, the version closest to the project is selected. We may revise that strategy in the future ([minimal version selection](https://about.sourcegraph.com/blog/the-pain-that-minimal-version-selection-solves/) maybe?), but this should give both control and the least amount of surprise for the site owner.

So, why did we do this? JavaScript is often a background actor in a Hugo project, and it doesn't make sense to publish it to a NPM registry. The JS dependencies are mostly build tools (PostCSS, TailwindCSS, Babel), `devDependencies`. This has been working fine as long as you kept the JS config files (including `package.json`) in the project, adding duplication/work when using ready-to-use theme components. These tools work best when you have everything below a single file tree, which is very much different to how [Hugo Modules](https://gohugo.io/hugo-modules/) work. An example of a module with TailwindCSS:

```bash
tailwind-module
├── assets
│   └── css
├── package.json
├── postcss.config.js
└── tailwind.config.js
```

If you included `tailwind-module` in a Hugo project and processed it with `PostCSS`, this is what happened in earlier versions:

* It used the directory where the `postcss.config.js` lives as a starting point to look  for any `require`'d dependency.
* TailwindCSS would, on the other hand, load its configuration file relative to where `PostCSS` was invoked (the project directory).

The above just doesn't work and here is the technical notes on how we have fixed this:

* The new `hugo mod npm pack` creates a consolidated `package.json` based on files named `package.hugo.json` it finds in the dependency tree (one is created for you the first time you run this command). The end result will always be `package.json`, which works seamlessly with `npm install` invoked automatically by Netlify and other CI vendors.
* The main project's `node_modules` folder is added to [NODE_PATH](https://medium.com/nafsadh/setting-up-node-path-for-using-global-packages-via-require-642eb711c725) when running `PostCSS`and `Babel`.
* We have introduced a new mount point `assets/_jsconfig` where we, by default, mount the JS configuration files that we're interested in. This is where Hugo will start looking for these files, and the files' filenames will also be available in the Node environment, so you can do:

```js
let tailwindConfig = process.env.HUGO_FILE_TAILWIND_CONFIG_JS || './tailwind.config.js';
const tailwind = require('tailwindcss')(tailwindConfig);
```

## Module Enhancements

* We have added a `noVendor` Glob pattern config to the module config [d4611c43](https://github.com/gohugoio/hugo/commit/d4611c4322dabfd8d2520232be578388029867db) [@bep](https://github.com/bep) [#7647](https://github.com/gohugoio/hugo/issues/7647). This allows you to only vendor a subset of your dependencies.
* We have added `ignoreImports` option to module imports config [20af9a07](https://github.com/gohugoio/hugo/commit/20af9a078189ce1e92a1d2047c90fba2a4e91827) [@bep](https://github.com/bep) [#7646](https://github.com/gohugoio/hugo/issues/7646), which allows you to import a module and load its config, but not follow its imports.
* We have deprecated `--ignoreVendor` in favour of a `--ignoreVendorPaths`, a patch matching Glob pattern [9a1e6d15](https://github.com/gohugoio/hugo/commit/9a1e6d15a31ec667b2ff9cf20e43b1daca61e004) [@bep](https://github.com/bep). A typical use for this would be when you have vendored your dependencies, but want to edit one of them.


## Statistics

This release represents **79 contributions by 19 contributors** to the main Hugo code base. [@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@dependabot[bot]](https://github.com/apps/dependabot), [@moorereason](https://github.com/moorereason), and [@jmooring](https://github.com/jmooring) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour), [@coliff](https://github.com/coliff) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **24 contributions by 15 contributors**. A special thanks to [@jmooring](https://github.com/jmooring), [@bep](https://github.com/bep), [@jornane](https://github.com/jornane), and [@inwardmovement](https://github.com/inwardmovement) for their work on the documentation site.


Hugo now has:

* 46596+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 438+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 352+ [themes](http://themes.gohugo.io/)

## Notes
* We now build with Go 1.15, which means that we no longer build release binaries for MacOS 32-bit.
* You may now get an error message about "error calling partial: partials that returns a value needs a non-zero argument.". This error situation was not caught earlier, and comes from a limitation in Go's templates: If you use the `return` keyword in a partial, the argument you pass to that partial (e.g. the ".") cannot be zero (and 0 and "" is considered a zero argument).

## Enhancements

### Templates

* Print layout name if it was specified when showing missing layout file error [9df60b62](https://github.com/gohugoio/hugo/commit/9df60b62f9c4e36a269f0c6e9a69bee9dc691031) [@richtera](https://github.com/richtera) [#7617](https://github.com/gohugoio/hugo/issues/7617)
* Add limit support to replaceRE [cdfd1c99](https://github.com/gohugoio/hugo/commit/cdfd1c99baa22d69e865294dfcd783811f96c880) [@moorereason](https://github.com/moorereason) [#7586](https://github.com/gohugoio/hugo/issues/7586)
* Extend merge to accept multiple parameters [047af7cf](https://github.com/gohugoio/hugo/commit/047af7cfe5e9aa740b85e0f9974a2d31a0ef4c08) [@moorereason](https://github.com/moorereason) [#7595](https://github.com/gohugoio/hugo/issues/7595)
* Add limit option to replace template function [f9ebaaed](https://github.com/gohugoio/hugo/commit/f9ebaaed1be1e4a26eef2aebd2c7554c979f29fa) [@moorereason](https://github.com/moorereason) [#7586](https://github.com/gohugoio/hugo/issues/7586)

### Output

* Respect mediatypes for deploy [12f6a1cd](https://github.com/gohugoio/hugo/commit/12f6a1cdc0aedf4319367af57bda3c94150d6a84) [@satotake](https://github.com/satotake) [#6861](https://github.com/gohugoio/hugo/issues/6861)

### Other

* Set PWD in environment when running the Node apps [377ad87a](https://github.com/gohugoio/hugo/commit/377ad87a51e0ef3619af4fe1be6aeee14c215c0a) [@bep](https://github.com/bep) 
* already -> already [292b0e26](https://github.com/gohugoio/hugo/commit/292b0e26ec9253398f7289dcf096691f63de2d96) [@dholbach](https://github.com/dholbach) 
* Regen docs helper [be2404c8](https://github.com/gohugoio/hugo/commit/be2404c8b17d3275cc82d9e659b9e41dddea7ded) [@bep](https://github.com/bep) 
* Regenerate CLI docs [c8da8eb1](https://github.com/gohugoio/hugo/commit/c8da8eb1f5551e6d141843daab41cb0ddbb0de4b) [@bep](https://github.com/bep) 
* Add "hugo mod npm pack" [85ba9bff](https://github.com/gohugoio/hugo/commit/85ba9bfffba9bfd0b095cb766f72700d4c211e31) [@bep](https://github.com/bep) [#7644](https://github.com/gohugoio/hugo/issues/7644)[#7656](https://github.com/gohugoio/hugo/issues/7656)[#7675](https://github.com/gohugoio/hugo/issues/7675)
* bump github.com/aws/aws-sdk-go from 1.34.21 to 1.34.22 [4fad43c8](https://github.com/gohugoio/hugo/commit/4fad43c8bd528f1805e78c50cd2e33822351c183) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Add support to linkable line anchors on Chroma [fb0f2cc7](https://github.com/gohugoio/hugo/commit/fb0f2cc718a54fd0774a0367e0a60718b5731de5) [@fjorgemota](https://github.com/fjorgemota) [#7622](https://github.com/gohugoio/hugo/issues/7622)
* Bump bundled Node.js from v8.12.0 to v12.18.3 [748fd4cb](https://github.com/gohugoio/hugo/commit/748fd4cb0d083de7c173d4b04b874358750fc900) [@anthonyfok](https://github.com/anthonyfok) [#7278](https://github.com/gohugoio/hugo/issues/7278)
* Change confinement from strict to classic" [b82f440c](https://github.com/gohugoio/hugo/commit/b82f440c59a5bf466c0f4c0431af6099216b0e37) [@anthonyfok](https://github.com/anthonyfok) 
* bump github.com/getkin/kin-openapi from 0.14.0 to 0.22.0 [c8143efa](https://github.com/gohugoio/hugo/commit/c8143efa5d21d20bcf3fa1d4f3fb292e460f90d8) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/aws/aws-sdk-go from 1.34.20 to 1.34.21 [c80132bb](https://github.com/gohugoio/hugo/commit/c80132bbe50f443a8be06dcbc51b855a5a5f8fa2) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/spf13/viper from 1.6.1 to 1.7.1 [75fa4c5c](https://github.com/gohugoio/hugo/commit/75fa4c5c950a43e33dfadfa138f61126b548ac40) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Run "go mod tidy" [fd7969e0](https://github.com/gohugoio/hugo/commit/fd7969e0b09e282d1cd83281bc0f5a62080afe5a) [@bep](https://github.com/bep) 
* Update to Goldmark v1.2.1 [b7fa3c4b](https://github.com/gohugoio/hugo/commit/b7fa3c4bba73f873bda71ba028ef46ce58aad908) [@bep](https://github.com/bep) 
* bump github.com/aws/aws-sdk-go from 1.27.1 to 1.34.20 [746ba803](https://github.com/gohugoio/hugo/commit/746ba803afee8f0f56ee0655cc55087f1822d39c) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/mitchellh/mapstructure from 1.1.2 to 1.3.3 [612b7d37](https://github.com/gohugoio/hugo/commit/612b7d376f1c50abe1fe6fe5188d576c1f5f1743) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Change confinement from strict to classic [6f4ff1a4](https://github.com/gohugoio/hugo/commit/6f4ff1a4617ec42861d255db718286ceaf4f6c8e) [@anthonyfok](https://github.com/anthonyfok) 
* bump github.com/spf13/cobra from 0.0.5 to 0.0.7 [ddeca459](https://github.com/gohugoio/hugo/commit/ddeca45933ab6e58c1b5187ad58dd261c9059009) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/sanity-io/litter from 1.2.0 to 1.3.0 [31f2091f](https://github.com/gohugoio/hugo/commit/31f2091f5803129b97c2a3f6245acc8b788235c7) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Add noVendor to module config [d4611c43](https://github.com/gohugoio/hugo/commit/d4611c4322dabfd8d2520232be578388029867db) [@bep](https://github.com/bep) [#7647](https://github.com/gohugoio/hugo/issues/7647)
* Add ignoreImports to module imports config [20af9a07](https://github.com/gohugoio/hugo/commit/20af9a078189ce1e92a1d2047c90fba2a4e91827) [@bep](https://github.com/bep) [#7646](https://github.com/gohugoio/hugo/issues/7646)
* Make ignoreVendor a glob pattern [9a1e6d15](https://github.com/gohugoio/hugo/commit/9a1e6d15a31ec667b2ff9cf20e43b1daca61e004) [@bep](https://github.com/bep) [#7642](https://github.com/gohugoio/hugo/issues/7642)
* bump github.com/gorilla/websocket from 1.4.1 to 1.4.2 [84adecf9](https://github.com/gohugoio/hugo/commit/84adecf97baa91ab18cb26812fa864b4451d3c5f) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/fsnotify/fsnotify from 1.4.7 to 1.4.9 [573558a0](https://github.com/gohugoio/hugo/commit/573558a078c6aaa671de0224c2d62b6d451d667c) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/kyokomi/emoji [8b10c22f](https://github.com/gohugoio/hugo/commit/8b10c22f822f0874890d2d6df68439450b83ef89) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/markbates/inflect from 1.0.0 to 1.0.4 [195bd124](https://github.com/gohugoio/hugo/commit/195bd1243b350e7a7814e0c893d17c3c408039c7) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/frankban/quicktest from 1.7.2 to 1.10.2 [6a544ece](https://github.com/gohugoio/hugo/commit/6a544ece24c37c98e2e4770fab350d76a0553f6a) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Encode & in livereload injected code [4b430d45](https://github.com/gohugoio/hugo/commit/4b430d456afee9c6da5e5ab46084a05469be1430) [@axllent](https://github.com/axllent) 
* bump github.com/niklasfasching/go-org from 1.3.1 to 1.3.2 [b9f10c75](https://github.com/gohugoio/hugo/commit/b9f10c75cb74c1976fbbf3d9e8dcdd4f3d46e790) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/bep/golibsass from 0.6.0 to 0.7.0 [537c598e](https://github.com/gohugoio/hugo/commit/537c598e9a4d8b8b47f5bffbcf59f72e9a1902c1) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump golang.org/x/text from 0.3.2 to 0.3.3 [67348676](https://github.com/gohugoio/hugo/commit/67348676f703f3ad3f778da1cdfa0fe001e5f925) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/evanw/esbuild from 0.6.5 to 0.6.32 [f9cc0ec7](https://github.com/gohugoio/hugo/commit/f9cc0ec76ee84451583a16a0abb9b09d298c7e00) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/nicksnyder/go-i18n from 1.10.0 to 1.10.1 [b5483eed](https://github.com/gohugoio/hugo/commit/b5483eed6e8c07809fc818192e0ce00d9496565c) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Revert "Update dependabot.yml" [90285f47](https://github.com/gohugoio/hugo/commit/90285f47504f8f2e30254745dd795d4ef007e205) [@bep](https://github.com/bep) 
* Update replaceRE func [f7c1b5fe](https://github.com/gohugoio/hugo/commit/f7c1b5fe1c22ba5f16e3fa442df6a8a70711f23f) [@moorereason](https://github.com/moorereason) 
* Update replace func [183e8626](https://github.com/gohugoio/hugo/commit/183e8626070a5f55c11648082e3060e35231d934) [@moorereason](https://github.com/moorereason) 
* Update merge function [f50ee6bb](https://github.com/gohugoio/hugo/commit/f50ee6bbe5ec0c0a1f7c21da6629faaed23bbe71) [@moorereason](https://github.com/moorereason) 
* Update dependabot.yml [c0655ba6](https://github.com/gohugoio/hugo/commit/c0655ba6ce5db54e8fec2c0e2bef9965b9fb90fc) [@bep](https://github.com/bep) 
* Create dependabot.yml [a2dda22c](https://github.com/gohugoio/hugo/commit/a2dda22c368adbffbba74c8c388cc10299801692) [@bep](https://github.com/bep) 
* Remove Pygments from requirements.txt [910d81a6](https://github.com/gohugoio/hugo/commit/910d81a6927c30ad1126c1bfaab1155b970f6442) [@bep](https://github.com/bep) 
* Regen CLI docs [8c490a73](https://github.com/gohugoio/hugo/commit/8c490a73b3735a0db46abba9bbe15de5ed2167e1) [@bep](https://github.com/bep) 
* Regen docs helper [e6cd9da4](https://github.com/gohugoio/hugo/commit/e6cd9da42d415552ae69e6b0afae823fd5e0003c) [@bep](https://github.com/bep) 
* markup/asciidocext: Revert trace=true [dcf25c0b](https://github.com/gohugoio/hugo/commit/dcf25c0b49eefef0572ec66337a5721bfde22233) [@bep](https://github.com/bep) 
* Update to Go 1.15.1 and 1.14.8 [e820b366](https://github.com/gohugoio/hugo/commit/e820b366b91729313c68be04b413e8894efc4421) [@bep](https://github.com/bep) [#7627](https://github.com/gohugoio/hugo/issues/7627)
* Add support for .TableOfContents [3ba7c925](https://github.com/gohugoio/hugo/commit/3ba7c92530a80f2f04fe57705ab05c247a6e8437) [@npiganeau](https://github.com/npiganeau) [#1687](https://github.com/gohugoio/hugo/issues/1687)
* Add a test case [19ef27b9](https://github.com/gohugoio/hugo/commit/19ef27b98edca53c4138b01c0f7c7bfbeb5ffcf1) [@bep](https://github.com/bep) [#7619](https://github.com/gohugoio/hugo/issues/7619)
* Add SourceMap flag with inline option [c6b661de](https://github.com/gohugoio/hugo/commit/c6b661de826f3ed8768a97a5178b4e020cb2ace1) [@richtera](https://github.com/richtera) [#7607](https://github.com/gohugoio/hugo/issues/7607)
* Remove logic that hides 'Building Sites' message after build completes [d39636a5](https://github.com/gohugoio/hugo/commit/d39636a5fc6bb82b3e0bd013858c7d116faa0c6b) [@jwarner112](https://github.com/jwarner112) [#7579](https://github.com/gohugoio/hugo/issues/7579)
* Improve stderr logging for PostCSS and simlilar [ec374204](https://github.com/gohugoio/hugo/commit/ec37420468157284651ef6b04b30420b249179e2) [@bep](https://github.com/bep) [#7584](https://github.com/gohugoio/hugo/issues/7584)
* Fail on  partials with return when given none or a zero argument [ae63c2b5](https://github.com/gohugoio/hugo/commit/ae63c2b5c94f68fbabd5dbd821630e747e8959a4) [@bep](https://github.com/bep) [#7572](https://github.com/gohugoio/hugo/issues/7572)[#7528](https://github.com/gohugoio/hugo/issues/7528)
* Update to Go 1.15 [e627449c](https://github.com/gohugoio/hugo/commit/e627449c0a2f1d2ffac29357c4f1832fc5462870) [@bep](https://github.com/bep) [#7554](https://github.com/gohugoio/hugo/issues/7554)
* Revert "Update stale.yml" [c2235c6a](https://github.com/gohugoio/hugo/commit/c2235c6a62d29e0a9e2e274eb340358a445b695d) [@bep](https://github.com/bep) 
* Update stale.yml [4f69ade7](https://github.com/gohugoio/hugo/commit/4f69ade7118302abff97169d17bfa9baac6a711c) [@bep](https://github.com/bep) 
* Remove trailing whitespace and tabs from RSS templates [5f425901](https://github.com/gohugoio/hugo/commit/5f42590144579c318a444ea2ce46d5c3fbbbfe6e) [@solarkennedy](https://github.com/solarkennedy) 
* Add uninstall target [21dbfa1f](https://github.com/gohugoio/hugo/commit/21dbfa1f111ca2f066e06af68f267932ce6cf04f) [@felicianotech](https://github.com/felicianotech) 
* Update Chroma to 0.8.0 [e5591e89](https://github.com/gohugoio/hugo/commit/e5591e89d3a71560b70c5f0ded33f2c9465ffe5a) [@jmooring](https://github.com/jmooring) [#7517](https://github.com/gohugoio/hugo/issues/7517)
* Update go-org to v1.3.1 [88929bc2](https://github.com/gohugoio/hugo/commit/88929bc23f5a830645c4e2cdac60aa43f480a478) [@niklasfasching](https://github.com/niklasfasching) 
* Collect transition attributes as classes [00e00da2](https://github.com/gohugoio/hugo/commit/00e00da233ab4d643de90bafca00f60ee0bbe785) [@bep](https://github.com/bep) [#7509](https://github.com/gohugoio/hugo/issues/7509)
* Add option for setting bundle format [0256959a](https://github.com/gohugoio/hugo/commit/0256959a358bb26b983c9d9496862b0fdf387621) [@bep](https://github.com/bep) [#7503](https://github.com/gohugoio/hugo/issues/7503)
* Simplify options handling [eded9ac2](https://github.com/gohugoio/hugo/commit/eded9ac2a05b9a7244c25c70ca8f761b69b33385) [@bep](https://github.com/bep) [#7499](https://github.com/gohugoio/hugo/issues/7499)
* make sure documentation intro text only appears once [8d725128](https://github.com/gohugoio/hugo/commit/8d72512825b4cee12dc1952004f48fd076a3517b) [@TheHippo](https://github.com/TheHippo) 
* Add es5 build target [e81aef0a](https://github.com/gohugoio/hugo/commit/e81aef0a954623e4a19062d1534bd8c2af97102a) [@bep](https://github.com/bep) 
* esbuild v0.6.5 [9f919147](https://github.com/gohugoio/hugo/commit/9f9191471ec501f1f957020726f939c9ef48e193) [@bep](https://github.com/bep) 
* Add .Defines to js.Build options [35011bcb](https://github.com/gohugoio/hugo/commit/35011bcb26b6fcfcbd77dc05aa8246ca45b2c2ba) [@bep](https://github.com/bep) [#7489](https://github.com/gohugoio/hugo/issues/7489)

## Fixes

### Other

* Fix AsciiDoc TOC with code [6a848cbc](https://github.com/gohugoio/hugo/commit/6a848cbc3a2487c8b015e715c2de44aef6051080) [@helfper](https://github.com/helfper) [#7649](https://github.com/gohugoio/hugo/issues/7649)
* markup/asciidocext: Fix broken test [4949bdc2](https://github.com/gohugoio/hugo/commit/4949bdc2ef98a1aebe5536c554d214f15c574a81) [@bep](https://github.com/bep) 
* Fix some change detection issues on server reloads [4055c121](https://github.com/gohugoio/hugo/commit/4055c121847847d8bd6b95a928185daee065091b) [@bep](https://github.com/bep) [#7623](https://github.com/gohugoio/hugo/issues/7623)[#7624](https://github.com/gohugoio/hugo/issues/7624)[#7625](https://github.com/gohugoio/hugo/issues/7625)
* Fixed misspelled words [ad01aea3](https://github.com/gohugoio/hugo/commit/ad01aea3f426206c2b70bbd97c5d29562dfe954d) [@aurkenb](https://github.com/aurkenb) 
* Fix a typo in CONTRIBUTING.md [f3cb0be3](https://github.com/gohugoio/hugo/commit/f3cb0be35adddfe43423a19116994b53817d97f7) [@capnfabs](https://github.com/capnfabs) 
* Revert "Fix ellipsis display logic in pagination template" [bffc4e12](https://github.com/gohugoio/hugo/commit/bffc4e12fe6d255e1fb8d28943993afc7e99e010) [@jmooring](https://github.com/jmooring) 
* Fix ellipsis display logic in pagination template [2fa851e6](https://github.com/gohugoio/hugo/commit/2fa851e6500752c0cea1da5cfdfc6d99e0a81a71) [@jmooring](https://github.com/jmooring) [#7523](https://github.com/gohugoio/hugo/issues/7523)
* Fix Asciidoctor args [45c665d3](https://github.com/gohugoio/hugo/commit/45c665d396ed368261f4a63ceee753c7f6dc5bf9) [@helfper](https://github.com/helfper) [#7493](https://github.com/gohugoio/hugo/issues/7493)
* Fix date format in internal schema template [a06c06a5](https://github.com/gohugoio/hugo/commit/a06c06a5c202de85ff47792b7468bfaeec2fea12) [@jmooring](https://github.com/jmooring) [#7495](https://github.com/gohugoio/hugo/issues/7495)
* Fix baseof block regression [c91dbe4c](https://github.com/gohugoio/hugo/commit/c91dbe4ce9c30623ba6e686fd17efae935aa0cc5) [@bep](https://github.com/bep) [#7478](https://github.com/gohugoio/hugo/issues/7478)





