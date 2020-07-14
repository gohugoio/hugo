
---
date: 2020-07-13
title: "Native JS Bundler, Open API Support, Inline Partials"
description: "Hugo 0.74.0 brings blazingly fast native JavaScript bundling, with minification, tree shaking, scope hoisting for ES6 modules, and transpilation of JSX and newer JS syntax down to ES6. And more."
categories: ["Releases"]
---

**Note:** The documentation site isn't updated with all of the main new things below. We will get to it soon. See https://github.com/gohugoio/hugoDocs/issues/1171

This release comes with native JavaScript bundling (and minifier), with [import](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/import) support (from `node_modules` etc.), tree shaking, scope hoisting for ES6 modules, transpilation of JSX and newer JS syntax down to ES6, JavaScript/JSX and TypeScript/TSX support. And it's _very fast_. [Benchmarks](https://github.com/evanw/esbuild#benchmarks) rates it **at least 100x faster** than the other JavaScript bundlers included. This new feature is backed by the very impressive [ESBuild](https://github.com/evanw/esbuild) project by [@evanw](https://github.com/evanw). Many thanks to [@remko](https://github.com/remko) for the integration work.

A very simple example building a TypeScript file:

```go-html-template
{{ $js := resources.Get "js/main.ts" | js.Build }}
```
This release also comes with Open API 3-support. This makes it much easier to create "Swagger styled" API-documentation. The below will unmarshal your YAML file into [this object graph](https://godoc.org/github.com/getkin/kin-openapi/openapi3#Swagger):

```go-html-template
{{ $api := resources.Get "api/openapi.yaml" | openapi3.Unmarshal }}
```

Hugo's Asciidoc integration has also gotten a face lift. A big shoutout to [@muenchhausen](https://github.com/muenchhausen) and [@bwklein](https://github.com/bwklein) for their work on this.

And finally, [partials](https://gohugo.io/templates/partials/#inline-partials) can now be defined inline -- and that is way more useful than it sounds.
 

This release represents **23 contributions by 9 contributors** to the main Hugo code base. [@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@niklasfasching](https://github.com/niklasfasching), [@bwklein](https://github.com/bwklein), and [@muenchhausen](https://github.com/muenchhausen) for their ongoing contributions.

And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour), [@coliff](https://github.com/coliff) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **8 contributions by 7 contributors**. A special thanks to [@OmarEmaraDev](https://github.com/OmarEmaraDev), [@regisphilibert](https://github.com/regisphilibert), [@coliff](https://github.com/coliff), and [@jessicahuynh](https://github.com/jessicahuynh) for their work on the documentation site.


Hugo now has:

* 45377+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 438+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 331+ [themes](http://themes.gohugo.io/)

## Enhancements

### Templates

* Add strings.Count [028b3567](https://github.com/gohugoio/hugo/commit/028b356787426dbc190ce9868fbc9a6400c2996e) [@bep](https://github.com/bep) [#7453](https://github.com/gohugoio/hugo/issues/7453)
* Add debug.Dump [defd7106](https://github.com/gohugoio/hugo/commit/defd7106bf79a502418ec373bdb82742b16f777f) [@bep](https://github.com/bep) [#3957](https://github.com/gohugoio/hugo/issues/3957)

### Output

* Add proper Media Type handling in js.Build [9df98ec4](https://github.com/gohugoio/hugo/commit/9df98ec49ca9fa326125ccfee626b6e46c6ab14b) [@bep](https://github.com/bep) [#732](https://github.com/gohugoio/hugo/issues/732)

### Core

* Add missing zero check on file [ccfaeb67](https://github.com/gohugoio/hugo/commit/ccfaeb678b312535928af3451324a54f2c7cb199) [@bep](https://github.com/bep) 

### Other

* Regenerate docs helper [25e3da33](https://github.com/gohugoio/hugo/commit/25e3da3343b5cd4bbcd11fa76b382fb089971840) [@bep](https://github.com/bep) 
* Add js.Build asset bundling [2fc33807](https://github.com/gohugoio/hugo/commit/2fc33807077cd25bf91f2298bf1a8ace126881a7) [@remko](https://github.com/remko) [#7321](https://github.com/gohugoio/hugo/issues/7321)
* Add openapi3.Unmarshal [12a65e76](https://github.com/gohugoio/hugo/commit/12a65e76df9470d9563b91a22969ddb41b7c19aa) [@bep](https://github.com/bep) [#7442](https://github.com/gohugoio/hugo/issues/7442)[#7443](https://github.com/gohugoio/hugo/issues/7443)
* Remove trailing hyphen from auto heading ID [58c0f5e6](https://github.com/gohugoio/hugo/commit/58c0f5e6171cbf8e3ed8d73ac95a7b85168c5b2f) [@jmooring](https://github.com/jmooring) [#6798](https://github.com/gohugoio/hugo/issues/6798)
* Ensure that non-trivial default flag values are passed through. [a1c3e3c1](https://github.com/gohugoio/hugo/commit/a1c3e3c1f32bcbc3b3aa6921bdee98a9f795a2da) [@vangent](https://github.com/vangent) 
* Update formats.md doc for new allowed extensions. [e9f87c4e](https://github.com/gohugoio/hugo/commit/e9f87c4e3feee937d05504763935805fec26213c) [@bwklein](https://github.com/bwklein) 
* Update config.go to add two Asciidoctor extensions [beb6c03b](https://github.com/gohugoio/hugo/commit/beb6c03bc8f476b753e5f3e3bc7a4a2e3f8ad355) [@bwklein](https://github.com/bwklein) 
* Add support for inline partials [4a3efea7](https://github.com/gohugoio/hugo/commit/4a3efea7efe59cd3de7d0eb352836ab395a2b6b3) [@bep](https://github.com/bep) [#7444](https://github.com/gohugoio/hugo/issues/7444)
* Add support for native Org dates in frontmatter [c66dc6c7](https://github.com/gohugoio/hugo/commit/c66dc6c74fa3bbe308ccaade8c76071b49908129) [@sometimesfood](https://github.com/sometimesfood) 
* Update go-org to v1.3.0 [127d5feb](https://github.com/gohugoio/hugo/commit/127d5feb32b466c4a0035e81f86684920dd88cfe) [@niklasfasching](https://github.com/niklasfasching) 
* Update go-org to v1.2.0 [2d42ba91](https://github.com/gohugoio/hugo/commit/2d42ba912ba945230aa0be23c3c8256cba40ce99) [@niklasfasching](https://github.com/niklasfasching) 
* Update bug_report.md [5b7b5dea](https://github.com/gohugoio/hugo/commit/5b7b5dea1fe3494995c6a9c3368087abf47cdc12) [@bep](https://github.com/bep) 
* Remove some unused code [057b1377](https://github.com/gohugoio/hugo/commit/057b1377c5f4d0d80ee299293db06384a475ad19) [@bep](https://github.com/bep) 
* Add an option to print memory usage at intervals [48dbb593](https://github.com/gohugoio/hugo/commit/48dbb593f7cc0dceb55d232ac198e82f3df1c964) [@bep](https://github.com/bep) 
* Rework external asciidoctor integration [f0266e2e](https://github.com/gohugoio/hugo/commit/f0266e2ef3487bc57dd05402002fc816e3b40195) [@muenchhausen](https://github.com/muenchhausen) 
* Enable the embedded template test when race detector is off [77aa385b](https://github.com/gohugoio/hugo/commit/77aa385b84dbc1805ff7e34dafeadb181905c689) [@bep](https://github.com/bep) [#5926](https://github.com/gohugoio/hugo/issues/5926)
* Merge branch 'release-0.73.0' [545a1c1c](https://github.com/gohugoio/hugo/commit/545a1c1cedc93d091050bae07c02fc2435ad2d20) [@bep](https://github.com/bep) 
* Updated installation instruction about Sass/SCSS support [0b579db8](https://github.com/gohugoio/hugo/commit/0b579db80fba1bde7dab07ea92d622dd6214dcfb) [@mateusz-szczyrzyca](https://github.com/mateusz-szczyrzyca) 

## Fixes

### Other

* Fix server reload when non-HTML shortcode changes [42e150fb](https://github.com/gohugoio/hugo/commit/42e150fbfac736bd49bc7e50cb8cdf9f81386f59) [@bep](https://github.com/bep) [#7448](https://github.com/gohugoio/hugo/issues/7448)


