
---
date: 2020-02-20
title: "0.65.0: Hugo Reloaded!"
description: "Draft, expire, resource bundling, and fine grained publishing control for any page. And it's faster."
categories: ["Releases"]
---

**Hugo 0.65** generalizes how a page is packaged and published to be applicable to **any page**. This should solve some of the most common issues we see people ask and talk about on the [issue tracker](https://github.com/gohugoio/hugo/issues) and on the [forum](https://discourse.gohugo.io/).

## Release Highlights

### New in Hugo Core

Any [branch node](https://gohugo.io/content-management/page-bundles/#branch-bundles) can now bundle resources (images, data files etc.), even the taxonomy nodes (e.g. /categories).

List pages (sections and the home page) can now be added to taxonomies.

The front matter fields that control when and if to publish a piece of content (`draft`, `publishDate`, `expiryDate`) now also works for list pages, and is recursive.

We have added a new `_build` front matter keyword to provide fine-grained control over page publishing. The default values:

```yaml
_build:
  # Whether to add it to any of the page collections.
  # Note that the page can still be found with .Site.GetPage.
  list: true
  
  # Whether to render it.
  render: true
  
  # Whether to publish its resources. These will still be published on demand,
  # but enabling this can be useful if the originals (e.g. images) are
  # never used.
  publishResources: true
```

Note that all front matter keywords can be set in the [cascade](https://gohugo.io/content-management/front-matter#front-matter-cascade) on a branch node, which would be especially useful for `_build`.

We have also upgraded to the latest LibSass (v3.6.3). Nothing remarkable functional new here, but it makes Hugo ready for the upcoming [Dart Backport](https://github.com/sass/libsass/pull/2918).

And finally, we have added a `GetTerms` method on `Page`, making  listing the terms defined on this page in the given taxonomy much simpler:

```go-html-template
<ul>
    {{ range (.GetTerms "tags") }}
        <li><a href="{{ .Permalink }}">{{ .LinkTitle }}</a></li>
   {{ end }}
</ul>
```

### New in Hugo Modules

There are several improvements to the tooling used in [Hugo Modules](https://gohugo.io/hugo-modules/). One bug fix, but also some improvements to make it easier to manage:

* You can now recursively update your modules with `hugo mod get -u ./...`
* `hugo mod clean` will now only clean the cache for the current project and now also takes an optional module path pattern, e.g. `hugo mod clean --pattern "github.com/**"`
* A new command `hugo mod verify` is added to verify that the module cache matches the hashes in `go.sum`. Run with `hugo mod verify --clean` to delete any modules that fail this check.

See [hugo mod](https://gohugo.io/commands/hugo_mod/#see-also).

### Performance

The new features listed above required a structural simplification, and we do watch our weight when doing this. And the benchmarks show that Hugo should, in general, be slightly faster. This is especially true if you're using taxonomies, and the partial rebuilding on content changes should be considerably faster.

## Numbers

This release represents **34 contributions by 6 contributors** to the main Hugo code base.[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@satotake](https://github.com/satotake), [@QuLogic](https://github.com/QuLogic), and [@JaymoKang](https://github.com/JaymoKang) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **7 contributions by 4 contributors**. A special thanks to [@coliff](https://github.com/coliff), [@bep](https://github.com/bep), [@tibnew](https://github.com/tibnew), and [@nerg4l](https://github.com/nerg4l) for their work on the documentation site.

Hugo now has:

* 41724+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 439+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 299+ [themes](http://themes.gohugo.io/)

## Notes

* `.GetPage "members.md"` (the Page method) will now only do relative lookups, which is what most people would expect.
* There have been a slight change of how disableKinds for regular pages: They will not be rendered on its own, but will be added to the site collections.

## Enhancements

### Templates

* Adjust the RSS taxonomy logic [d73e3738](https://github.com/gohugoio/hugo/commit/d73e37387ca0012bd58bd3f36a0477854b41ab6e) [@bep](https://github.com/bep) [#6909](https://github.com/gohugoio/hugo/issues/6909)

### Output

* Handle disabled RSS even if it's defined in outputs [da54787c](https://github.com/gohugoio/hugo/commit/da54787cfa97789624e467a4451dfeb50f563e41) [@bep](https://github.com/bep) 

### Other

* Regenerate CLI docs [a5ebdf7d](https://github.com/gohugoio/hugo/commit/a5ebdf7d17e6c6a9dc686cf8f7cd8e0a1bab5f2d) [@bep](https://github.com/bep) 
* Improve "hugo mod clean" [dce210ab](https://github.com/gohugoio/hugo/commit/dce210ab56fc885818fc5d1a084a1c3ba84e7929) [@bep](https://github.com/bep) [#6907](https://github.com/gohugoio/hugo/issues/6907)
* Add "hugo mod verify" [0b96aba0](https://github.com/gohugoio/hugo/commit/0b96aba022d51cf9939605c029bb8dba806653a1) [@bep](https://github.com/bep) [#6907](https://github.com/gohugoio/hugo/issues/6907)
* Add Page.GetTerms [fa520a2d](https://github.com/gohugoio/hugo/commit/fa520a2d983b982394ad10088393fb303e48980a) [@bep](https://github.com/bep) [#6905](https://github.com/gohugoio/hugo/issues/6905)
* Add a list terms benchmark [7489a864](https://github.com/gohugoio/hugo/commit/7489a864591b6df03f435f40696c6ceeb4776ec9) [@bep](https://github.com/bep) [#6905](https://github.com/gohugoio/hugo/issues/6905)
* Use the tree for taxonomy.Pages() [b2dcd53e](https://github.com/gohugoio/hugo/commit/b2dcd53e3c0240c4afd21d1818fd180c2d1b9d34) [@bep](https://github.com/bep) 
* Add some cagegories to the site collections benchmarks [36983e61](https://github.com/gohugoio/hugo/commit/36983e6189a717f1d4d1da6652621d7f8fe186ad) [@bep](https://github.com/bep) 
* Do not try to get local themes in "hugo mod get" [20f2211f](https://github.com/gohugoio/hugo/commit/20f2211fce55e1811629245f9e5e4a2ac754d788) [@bep](https://github.com/bep) [#6893](https://github.com/gohugoio/hugo/issues/6893)
* Update goldmark-highlighting [a21a9373](https://github.com/gohugoio/hugo/commit/a21a9373e06091ab70d8a5f4da8ff43f7c609b4b) [@satotake](https://github.com/satotake) 
* Support "hugo mod get -u ./..." [775c7c24](https://github.com/gohugoio/hugo/commit/775c7c2474d8797c96c9ac529a3cd93c0c2d3514) [@bep](https://github.com/bep) [#6828](https://github.com/gohugoio/hugo/issues/6828)
* Introduce a tree map for all content [eada236f](https://github.com/gohugoio/hugo/commit/eada236f87d9669885da1ff647672bb3dc6b4954) [@bep](https://github.com/bep) [#6312](https://github.com/gohugoio/hugo/issues/6312)[#6087](https://github.com/gohugoio/hugo/issues/6087)[#6738](https://github.com/gohugoio/hugo/issues/6738)[#6412](https://github.com/gohugoio/hugo/issues/6412)[#6743](https://github.com/gohugoio/hugo/issues/6743)[#6875](https://github.com/gohugoio/hugo/issues/6875)[#6034](https://github.com/gohugoio/hugo/issues/6034)[#6902](https://github.com/gohugoio/hugo/issues/6902)[#6173](https://github.com/gohugoio/hugo/issues/6173)[#6590](https://github.com/gohugoio/hugo/issues/6590)
* Another benchmark rename [e5329f13](https://github.com/gohugoio/hugo/commit/e5329f13c02b87f0c30f8837759c810cd90ff8da) [@bep](https://github.com/bep) 
* Rename the Edit benchmarks [5b145ddc](https://github.com/gohugoio/hugo/commit/5b145ddc4c951a827e1ac00444dc4719e53e0885) [@bep](https://github.com/bep) 
* Refactor a benchmark to make it runnable as test [54bdcaac](https://github.com/gohugoio/hugo/commit/54bdcaacaedec178554e696f34647801bbe61362) [@bep](https://github.com/bep) 
* Add benchmark for content edits [1622510a](https://github.com/gohugoio/hugo/commit/1622510a5c651b59a79f64e9dc3cacd24832ec0b) [@bep](https://github.com/bep) 
* Add "go mod verify" to build scripts [56d0b658](https://github.com/gohugoio/hugo/commit/56d0b658879bbf476810d013176d6568553aa71e) [@bep](https://github.com/bep) 
* Add git to Dockerfile [75c3787f](https://github.com/gohugoio/hugo/commit/75c3787fc254d933fa11e5c39d978bfa1a21a371) [@JaymoKang](https://github.com/JaymoKang) 
* Update go.sum [9babb1f0](https://github.com/gohugoio/hugo/commit/9babb1f0c4fca048b0339f6ce3618f88d34e0457) [@bep](https://github.com/bep) 
* Rename doWithCommandeer to cfgInit/cfgSetAndInit [8a5124d6](https://github.com/gohugoio/hugo/commit/8a5124d6b38156cb6f765ac7492513ac7c0d90b2) [@MarkRosemaker](https://github.com/MarkRosemaker) 
* Update golibsass [898a0a96](https://github.com/gohugoio/hugo/commit/898a0a96afd472fad8fe70be71f6cb00a4267c4a) [@bep](https://github.com/bep) [#6885](https://github.com/gohugoio/hugo/issues/6885)
* Shuffle test files before insertion [3b721110](https://github.com/gohugoio/hugo/commit/3b721110d560c8831c282e6e7a5c510fe7a5129a) [@bep](https://github.com/bep) 
* Update to LibSass v3.6.3 [40ba7e6d](https://github.com/gohugoio/hugo/commit/40ba7e6d63c1a0734f257a642e46eb1572116a32) [@bep](https://github.com/bep) [#6862](https://github.com/gohugoio/hugo/issues/6862)
* Update Go version requirement [23ea4318](https://github.com/gohugoio/hugo/commit/23ea43180b84e35d99e88083a83e7ca1916b3b36) [@bep](https://github.com/bep) [#6853](https://github.com/gohugoio/hugo/issues/6853)

## Fixes

### Templates

* Fix RSS template for the terms listing [aa3e1830](https://github.com/gohugoio/hugo/commit/aa3e1830568cabaa8bf3277feeba6cb48746e40c) [@bep](https://github.com/bep) [#6909](https://github.com/gohugoio/hugo/issues/6909)

### Other

* Fix lazy publishing with publishResources=false [9bdedb25](https://github.com/gohugoio/hugo/commit/9bdedb251c7cd8f8af800c7d9914cf84292c5c50) [@bep](https://github.com/bep) [#6914](https://github.com/gohugoio/hugo/issues/6914)
* Fix goMinorVersion on non-final Go releases [c7975b48](https://github.com/gohugoio/hugo/commit/c7975b48b6532823868a6aa8c93eb76caa46c570) [@QuLogic](https://github.com/QuLogic) 
* Fix taxonomy [1b7acfe7](https://github.com/gohugoio/hugo/commit/1b7acfe7634a5d7bbc597ef4dddf4babce5666c5) [@bep](https://github.com/bep) 
* Fix RenderString for pages without content [19e12caf](https://github.com/gohugoio/hugo/commit/19e12caf8c90516e3b803ae8a40b907bd89dc96c) [@bep](https://github.com/bep) [#6882](https://github.com/gohugoio/hugo/issues/6882)
* Fix chroma highlight [3c568ad0](https://github.com/gohugoio/hugo/commit/3c568ad0139c79e5c0596ca40637512d71401afc) [@satotake](https://github.com/satotake) [#6877](https://github.com/gohugoio/hugo/issues/6877)[#6856](https://github.com/gohugoio/hugo/issues/6856)
* Fix mount with hole regression [b78576fd](https://github.com/gohugoio/hugo/commit/b78576fd38a76bbdaab5ad21228c8e5a559090b1) [@bep](https://github.com/bep) [#6854](https://github.com/gohugoio/hugo/issues/6854)
* Fix bundle resource ordering regression [18888e09](https://github.com/gohugoio/hugo/commit/18888e09bbb5325bdd63f2cd93116ff490dd37ab) [@bep](https://github.com/bep) [#6851](https://github.com/gohugoio/hugo/issues/6851)
* Fix note about CGO [7f0ebd4a](https://github.com/gohugoio/hugo/commit/7f0ebd4a3c9e016afddc2cf5e7dfe6a820aa099a) [@moorereason](https://github.com/moorereason) 





