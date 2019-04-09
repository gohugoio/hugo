
---
date: 2019-04-08
title: "Hugo 0.55.0: The early Easter Egg Edition!"
description: "Faster, virtualized Output Formats, revised shortcodes, new return keyword, and much more …"
categories: ["Releases"]
---

Hugo `0.55` is **the early Easter Egg Edition** with lots of great improvements and fixes. The original motivation for this release was to prepare for [Issue #5074](https://github.com/gohugoio/hugo/issues/5074), but the structural changes needed for that paved the way for lots of others. Please study the list of changes below, and especially the **Notes** section, but some headlines include:

## Virtualized Output Formats

[Custom Output Formats](https://gohugo.io/templates/output-formats) has been a really useful feature, but it has had some annoying and not so obvious restrictions that are now lifted. Now all `Page` collections are aware of the output format being rendered. This means, to give some examples, that:

* In a `RSS` template, listing pages with their content will use output format specific shortcode templates even if the pages themselves are not configured to output to that output format.
* Using `.Render` when looping over a `Page` collection will now work as expected.
* Every Output Format can be paginated.

We have now also added a new `Permalinkable` configuration attribute, which is enabled by default for `HTML` and `AMP`.

## Shortcodes Revised

Shortcodes using the `{{%/* */%}}` as the outer-most delimiter will now be fully rendered when sent to the content renderer (e.g. Blackfriday for Markdown), meaning they can be part of the generated table of contents, footnotes, etc.

If you want the old behavior, you can put the following line in the start of your shortcode template:

```
{{ $_hugo_config := `{ "version": 1 }` }}
```

But using the `{{</* */>}}` delimiter will, in most cases, be a better alternative, possibly in combination with the `markdownify` template func.

See [#5763](https://github.com/gohugoio/hugo/issues/5763).

## New Return Keyword for Partials

Hugo's `partial` and `partialCached` are great for template macros. With the new `return` keyword you can write partials as proper functions that can return any type:

```go-html-template
{{ $v := add . 42 }}
{{ return $v }}
```

See [#5783](https://github.com/gohugoio/hugo/issues/5783).

## .Page on Taxonomy nodes

The taxonomy nodes now have a `.Page` accessor which makes it much simpler to get a proper `.Title` etc. This means that older and clumsy constructs can be simplified. Some examples:

```go-html-template
<ul>
    {{ range .Data.Terms.Alphabetical }}
            <li><a href="{{ .Page.Permalink }}">{{ .Page.Title }}</a> {{ .Count }}</li>
    {{ end }}
</ul>
```

```go-html-template
<ul>
    {{ range .Site.Taxonomies.tags }}
            <li><a href="{{ .Page.Permalink }}">{{ .Page.Title }}</a> {{ .Count }}</li>
    {{ end }}
</ul>
```

See [#5719](https://github.com/gohugoio/hugo/issues/5719).

## And it's Faster!

This version is also the **fastest to date**. A site building benchmark shows more than 10% decrease in both build time and memory consumption, but that depends on the site. It’s quite a challenge to consistently add significant new functionality and simultaneously improve performance. Also, note that we are now more honest about build times reported (`Total in 1234 ms`). We now do all initialization in the `Build` step, so you may get a higher time reported if you, as an example, have `--enableGitInfo` enabled, which now is included in the reported time.

![Benchmark](https://pbs.twimg.com/media/D3kGYiMXsAUjYxS.png)

## Thanks!

This release represents **59 contributions by 10 contributors** to the main Hugo code base.
[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@mcdee](https://github.com/mcdee), [@quasilyte](https://github.com/quasilyte), and [@danielcompton](https://github.com/danielcompton) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **36 contributions by 21 contributors**. A special thanks to [@bep](https://github.com/bep), [@peaceiris](https://github.com/peaceiris), [@budparr](https://github.com/budparr), and [@tinymachine](https://github.com/tinymachine) for their work on the documentation site.

As this release has required a significant effort with manual testing, a special thanks go to [@onedrawingperday](https://github.com/onedrawingperday) (the 300 theme demo sites have been invaluable to check for API-breakage!), [@adiabatic](https://github.com/adiabatic), and [@divinerites](https://github.com/divinerites).

Hugo now has:

* 34077+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 439+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 306+ [themes](http://themes.gohugo.io/)


## Notes
* `{{ %` as the outer-most shortcode delimiter means "pass the inner content to the content renderer" (e.g. Blackfriday);  the old behavior can be had, see  [#5763](https://github.com/gohugoio/hugo/issues/5763)
* `preserveTaxonomyNames`configuration option is removed. Use `.Page.Title`.
* We no longer limit the number of pages passed to the `RSS` Output Format. We have moved that limit to the internal `RSS` template, and you can do so yourself using the `Config.Services.RSS.Limit` in your custom template.
* We no longer add XML headers to Output Formats that output XML (`<?xml version="1.0" encoding="utf-8" standalone="yes" ?>`). This header is moved to the templates. If you have custom RSS or sitemap templates you may want to add the XML declaration to these. Since they, by default, is handled by Go's HTML template package, you must do something like this to make sure it's preserved: `{{ printf "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?>" | safeHTML }}`
* More honest build times reported (`Total in 1234 ms`). We now do all initialization in the `Build` step, so you may get a higher time reported if you, as an example, have `--enableGitInfo` enabled, which now is included in the reported time.
* The taxonomy nodes now have a `.Page` accessor which makes it much simpler to get a proper `.Title` etc. see [#5719](https://github.com/gohugoio/hugo/issues/5719).
* The template keywords `with` and `if` now work properly for zero and interface types, see [#5739](https://github.com/gohugoio/hugo/issues/5739)
* Taxonomy terms lists (`Page` of `Kind` `taxonomyTerm`) without any date(s) set (e.g. from front matter) will now pick the latest dates from the child pages. This is in line with how other branch nodes work in Hugo.
* A new configuration option, `Permalinkable`, is added to Output Format and enabled that by default for `HTML` and `AMP` types. See  [#5706](https://github.com/gohugoio/hugo/issues/5706) 
* `--stepAnalysis` is removed. If you want to really understand the latency in your project in Hugo, try the new `--trace` flag and pass that file to the many tools that read [Go Trace files](https://golang.org/pkg/runtime/trace/). There are also some newly hidden flags in `--profile-cpu`, `--profile-men`, `--profile-mutex`, hidden because they are considered to be only of interest to developers.
* Chroma is updated with many fixes and new lexers, see [#5797](https://github.com/gohugoio/hugo/issues/5797)
* We now support `Page`-relative aliases, e.g. aliases that do not start with a `/`, see [#5757](https://github.com/gohugoio/hugo/issues/5757)
* We now support context-relative (language) URLs in front matter, meaning that in most cases `url: /en/mypage` can be replaced with the more portable `url: mypage`. See [#5704](https://github.com/gohugoio/hugo/issues/5704)

## Enhancements

### Templates

* Allow the partial template func to return any type [a55640de](https://github.com/gohugoio/hugo/commit/a55640de8e3944d3b9f64b15155148a0e35cb31e) [@bep](https://github.com/bep) [#5783](https://github.com/gohugoio/hugo/issues/5783)

### Output

* Add missing JSON tag [b6a60f71](https://github.com/gohugoio/hugo/commit/b6a60f718e376066456da37e7bb997a7697edc31) [@bep](https://github.com/bep) 

### Core

* Log warning on relative front matter url with lang [f34e6172](https://github.com/gohugoio/hugo/commit/f34e6172cf2a4d1d1aef22304ecbc7c8e2d142ff) [@bep](https://github.com/bep) [#5818](https://github.com/gohugoio/hugo/issues/5818)
* Consider summary in front matter for .Summary [3a62d547](https://github.com/gohugoio/hugo/commit/3a62d54745e2cbfda6772390830042908d725c71) [@mcdee](https://github.com/mcdee) [#5800](https://github.com/gohugoio/hugo/issues/5800)
* Buffer the render pages chan [95029551](https://github.com/gohugoio/hugo/commit/950295516da882dcc51d83f70835dde230a0b4d6) [@bep](https://github.com/bep) 
* Re-work "fast render" logic in the new flow [d0d661df](https://github.com/gohugoio/hugo/commit/d0d661dffd19d5ed6efbd4dd2c572bad008bd859) [@bep](https://github.com/bep) [#5811](https://github.com/gohugoio/hugo/issues/5811)[#5784](https://github.com/gohugoio/hugo/issues/5784)
* Allow relative URLs in front matter [5185fb06](https://github.com/gohugoio/hugo/commit/5185fb065b0f8a4142c29ee3e3cd917e917280a4) [@bep](https://github.com/bep) [#5704](https://github.com/gohugoio/hugo/issues/5704)
* Allow page-relative aliases [92baa14f](https://github.com/gohugoio/hugo/commit/92baa14fd3f45c0917c5988235cd1a0f8692f171) [@bep](https://github.com/bep) [#5757](https://github.com/gohugoio/hugo/issues/5757)
* Add a simple test for jsonify of Site [8bfd3a54](https://github.com/gohugoio/hugo/commit/8bfd3a54a4142c397cab69bfa9699e5b5db9b40b) [@bep](https://github.com/bep) [#5780](https://github.com/gohugoio/hugo/issues/5780)
* Do not fall back to site title if not set in content file [9bc6187b](https://github.com/gohugoio/hugo/commit/9bc6187b8337c4a370bd3f21130a764d9ef6f7b3) [@bep](https://github.com/bep) [#5784](https://github.com/gohugoio/hugo/issues/5784)
* Add a test for home page with no title [bceda1b2](https://github.com/gohugoio/hugo/commit/bceda1b288f0ad6282916826b596cb1fe19983bb) [@bep](https://github.com/bep) [#5784](https://github.com/gohugoio/hugo/issues/5784)
* Add String() to fileInfo [a7ee9b0b](https://github.com/gohugoio/hugo/commit/a7ee9b0bb98f519e485655af578fb35d755e5c44) [@bep](https://github.com/bep) [#5784](https://github.com/gohugoio/hugo/issues/5784)
* Remove unused slice [3011f36c](https://github.com/gohugoio/hugo/commit/3011f36c27ecde309325e6c75ca377f4f87fa97a) [@bep](https://github.com/bep) 
* Adjust site benchmark [34c49d78](https://github.com/gohugoio/hugo/commit/34c49d788c102a370006e476d6f6143a51b2a03d) [@bep](https://github.com/bep) 
* Adjust test for Go 1.12 [b4148cd1](https://github.com/gohugoio/hugo/commit/b4148cd1d9ea889b81070d3e84a37bd5d23e5746) [@bep](https://github.com/bep) 

### Other

* Misc paginator adjustments [612a06f0](https://github.com/gohugoio/hugo/commit/612a06f0671125be6b42ec2982a18080005994c8) [@bep](https://github.com/bep) [#5825](https://github.com/gohugoio/hugo/issues/5825)
* Update to Go 1.12.2 and Go 1.11.7 [3db4a1cf](https://github.com/gohugoio/hugo/commit/3db4a1cf7ab12343ce5705ac56aa7ca6ea1677b6) [@bep](https://github.com/bep) [#5819](https://github.com/gohugoio/hugo/issues/5819)
* Adjust rlimit logic [708d4cee](https://github.com/gohugoio/hugo/commit/708d4ceebd491c6a89f271311eb8d94d6b5d58bc) [@bep](https://github.com/bep) [#5821](https://github.com/gohugoio/hugo/issues/5821)
* Add information about summary front matter variable [ed65bda3](https://github.com/gohugoio/hugo/commit/ed65bda3b43f6149e41ddb049cbb295a82473bc9) [@mcdee](https://github.com/mcdee) 
* Regenerate JSON wrapper [ebab291c](https://github.com/gohugoio/hugo/commit/ebab291c0e321d23b098684bacaf830a3979e310) [@bep](https://github.com/bep) 
* Add missing GitInfo to Page [75467cd7](https://github.com/gohugoio/hugo/commit/75467cd7852852305549a6c71ac503bb4a57e716) [@bep](https://github.com/bep) 
* Add support for sha384 [d1553b4b](https://github.com/gohugoio/hugo/commit/d1553b4b0f83e4a4305d2b4ab9ba6e305637f134) [@bep](https://github.com/bep) [#5815](https://github.com/gohugoio/hugo/issues/5815)
* Add HUGO_NUMWORKERMULTIPLIER [87b16abd](https://github.com/gohugoio/hugo/commit/87b16abd93ff60acd245776d5b0d914fd580c259) [@bep](https://github.com/bep) [#5814](https://github.com/gohugoio/hugo/issues/5814)
* Use YAML for the benchmark compare [8559f5c2](https://github.com/gohugoio/hugo/commit/8559f5c29f20b7b5188f93f8b1d9e510e3dee4f5) [@bep](https://github.com/bep) 
* Update to imaging v1.6.0 [032e6802](https://github.com/gohugoio/hugo/commit/032e6802d1f34cc41f6d1275fdd2deab8bbe5480) [@bep](https://github.com/bep) [#5812](https://github.com/gohugoio/hugo/issues/5812)
* Adjust the howSimilar logic vs strings [4494a01b](https://github.com/gohugoio/hugo/commit/4494a01b794ab785c64c8e93c61ccbfa845bc478) [@bep](https://github.com/bep) 
* Implement compare.ProbablyEqer for the core slices [e91e222c](https://github.com/gohugoio/hugo/commit/e91e222cd21213961d1e6206e1523bee2c21fa0c) [@bep](https://github.com/bep) [#5808](https://github.com/gohugoio/hugo/issues/5808)
* Regenerate docshelper data [bfdc4496](https://github.com/gohugoio/hugo/commit/bfdc44964af82807fa91407132d47b6bf52704c3) [@bep](https://github.com/bep) [#5799](https://github.com/gohugoio/hugo/issues/5799)
* Update Chroma [cc8515f1](https://github.com/gohugoio/hugo/commit/cc8515f18767298da4c6d712d1fd747c7950150b) [@bep](https://github.com/bep) [#5780](https://github.com/gohugoio/hugo/issues/5780)
* Regenerate CLI docs [bb533ca5](https://github.com/gohugoio/hugo/commit/bb533ca5e1c778c95ed7014eab99c8cc1bd4c85e) [@bep](https://github.com/bep) [#5779](https://github.com/gohugoio/hugo/issues/5779)
* Update Afero [10bb614a](https://github.com/gohugoio/hugo/commit/10bb614a70db22c01c9a52054ede35bc0a01aa24) [@bep](https://github.com/bep) [#5673](https://github.com/gohugoio/hugo/issues/5673)
* Avoid nilpointer on no File on Page [4dae52af](https://github.com/gohugoio/hugo/commit/4dae52af680e6ff2c8cdeb4ce1f219330b27001c) [@bep](https://github.com/bep) [#5781](https://github.com/gohugoio/hugo/issues/5781)
* Improve the "feature not available" error [794d4052](https://github.com/gohugoio/hugo/commit/794d4052b87c98943588b35e1cfecc06e6a0c7f2) [@bep](https://github.com/bep) 
* Re-introduce .Page.Page [91ef9655](https://github.com/gohugoio/hugo/commit/91ef9655aaf2adea3a044bf9a464908084917a98) [@bep](https://github.com/bep) [#5784](https://github.com/gohugoio/hugo/issues/5784)
* Apply staticcheck recommendations [b5f39d23](https://github.com/gohugoio/hugo/commit/b5f39d23b86f9cb83c51da9fe4abb4c19c01c3b7) [@bep](https://github.com/bep) 
* Run gofmt -s [d30e8454](https://github.com/gohugoio/hugo/commit/d30e845485b416e1c48fade14694b12a9fe59b6b) [@bep](https://github.com/bep) 
* Make Page an interface [597e418c](https://github.com/gohugoio/hugo/commit/597e418cb02883418f2cebb41400e8e61413f651) [@bep](https://github.com/bep) [#5074](https://github.com/gohugoio/hugo/issues/5074)[#5763](https://github.com/gohugoio/hugo/issues/5763)[#5758](https://github.com/gohugoio/hugo/issues/5758)[#5090](https://github.com/gohugoio/hugo/issues/5090)[#5204](https://github.com/gohugoio/hugo/issues/5204)[#4695](https://github.com/gohugoio/hugo/issues/4695)[#5607](https://github.com/gohugoio/hugo/issues/5607)[#5707](https://github.com/gohugoio/hugo/issues/5707)[#5719](https://github.com/gohugoio/hugo/issues/5719)[#3113](https://github.com/gohugoio/hugo/issues/3113)[#5706](https://github.com/gohugoio/hugo/issues/5706)[#5767](https://github.com/gohugoio/hugo/issues/5767)[#5723](https://github.com/gohugoio/hugo/issues/5723)[#5769](https://github.com/gohugoio/hugo/issues/5769)[#5770](https://github.com/gohugoio/hugo/issues/5770)[#5771](https://github.com/gohugoio/hugo/issues/5771)[#5759](https://github.com/gohugoio/hugo/issues/5759)[#5776](https://github.com/gohugoio/hugo/issues/5776)[#5777](https://github.com/gohugoio/hugo/issues/5777)[#5778](https://github.com/gohugoio/hugo/issues/5778)
* List future and expired dates in CSV format [44f5c1c1](https://github.com/gohugoio/hugo/commit/44f5c1c14cb1f42cc5f01739c289e9cfc83602af) [@danielcompton](https://github.com/danielcompton) [#5610](https://github.com/gohugoio/hugo/issues/5610)
* Update to Go 1.12.1 and Go 1.11.6 [984a73af](https://github.com/gohugoio/hugo/commit/984a73af9e5b5145297723f26faa38f29ca2918d) [@bep](https://github.com/bep) [#5755](https://github.com/gohugoio/hugo/issues/5755)
* Update Viper [79d517d8](https://github.com/gohugoio/hugo/commit/79d517d86c02e879bc4a43ab86b817c589b61485) [@bep](https://github.com/bep) 
* Update to Go 1.12 [b9e75afd](https://github.com/gohugoio/hugo/commit/b9e75afd6c007a6af8b71caeebc4a5a24c270861) [@bep](https://github.com/bep) [#5716](https://github.com/gohugoio/hugo/issues/5716)
* Remove Gitter dev chat link [dfc72d61](https://github.com/gohugoio/hugo/commit/dfc72d61a522f5cb926271d9391a8670f064d198) [@bep](https://github.com/bep) 
* Update Travis config to work for forked builds [bdf47e8d](https://github.com/gohugoio/hugo/commit/bdf47e8da80f87b7689badf48a6b8672c048d7e4) [@grahamjamesaddis](https://github.com/grahamjamesaddis) 
* Add skipHTML option to blackfriday config [75904332](https://github.com/gohugoio/hugo/commit/75904332f3bedcfe656856821d4c9560a177cc51) [@arrtchiu](https://github.com/arrtchiu) 
* Update stretchr/testify to 1.3.0. [60c0eb4e](https://github.com/gohugoio/hugo/commit/60c0eb4e892baedd533424b47baf7039c0005f87) [@QuLogic](https://github.com/QuLogic) 
* Rewrite relative action URLS [c154c2f7](https://github.com/gohugoio/hugo/commit/c154c2f7b2a6703dbde7f6bd2a1817a39c6fd2ea) [@larson004](https://github.com/larson004) [#5701](https://github.com/gohugoio/hugo/issues/5701)
* Support Docker args TAGS, WORKDIR, CGO; speed up repetitive builds [075b17ee](https://github.com/gohugoio/hugo/commit/075b17ee1d621e0ebbcecf1063f8f68a00ac221a) [@tonymet](https://github.com/tonymet) 
* Support nested keys/fields with missing values with the `where` function [908692fa](https://github.com/gohugoio/hugo/commit/908692fae5c5840a0db8c7dd389b59dd3b8026b9) [@tryzniak](https://github.com/tryzniak) [#5637](https://github.com/gohugoio/hugo/issues/5637)[#5416](https://github.com/gohugoio/hugo/issues/5416)
* Update debouncer version [7e4b18c5](https://github.com/gohugoio/hugo/commit/7e4b18c5ae409435760ebd86ff9ee3061db34a5d) [@bep](https://github.com/bep) 

## Fixes

### Templates

* Fix mutex unlock [e54213f5](https://github.com/gohugoio/hugo/commit/e54213f5257267ed232b2465337c39ddc8c73388) [@bep](https://github.com/bep) 
* Fix template truth logic [02eaddc2](https://github.com/gohugoio/hugo/commit/02eaddc2fbe92c26e67d9f82dd9aabecbbf2106c) [@bep](https://github.com/bep) [#5738](https://github.com/gohugoio/hugo/issues/5738)
* Fix strings.HasPrefix args order [72010429](https://github.com/gohugoio/hugo/commit/7201042946dde78d5ea4fea9cb006fb4dded55c1) [@quasilyte](https://github.com/quasilyte) 

### Core

* Fix default date assignment for sections [1d9dde82](https://github.com/gohugoio/hugo/commit/1d9dde82a0577d93eea8ed0a7ec0b4ae3068eb19) [@bep](https://github.com/bep) [#5784](https://github.com/gohugoio/hugo/issues/5784)
* Fix the GOMAXPROCS env get [415ca967](https://github.com/gohugoio/hugo/commit/415ca9673d3bd3c06ab94f3d83897c892fce5f27) [@bep](https://github.com/bep) [#5813](https://github.com/gohugoio/hugo/issues/5813)
* Fix benchmark for YAML front matter [e2dc432f](https://github.com/gohugoio/hugo/commit/e2dc432fe287a280aeba94bafdcce85b7a8646c6) [@bep](https://github.com/bep) 
* Fix alias path for AMP and similar [f9d6feca](https://github.com/gohugoio/hugo/commit/f9d6feca0802cd83c4d843244ce389cf7c792cec) [@bep](https://github.com/bep) [#5760](https://github.com/gohugoio/hugo/issues/5760)

### Other

* Fix image publish ordering issue [439ab033](https://github.com/gohugoio/hugo/commit/439ab0339d9ac6972caabaa55fa41887ace839cb) [@bep](https://github.com/bep) [#5730](https://github.com/gohugoio/hugo/issues/5730)
* Fix doLiveReload logic [4a2a8aff](https://github.com/gohugoio/hugo/commit/4a2a8afff2021c8e967254c76c159147da7e78fa) [@bep](https://github.com/bep) [#5754](https://github.com/gohugoio/hugo/issues/5754)
* Fix args order in strings.TrimPrefix [483cf19d](https://github.com/gohugoio/hugo/commit/483cf19d5de05e8a83fd1be6934baa169c7fd7c8) [@quasilyte](https://github.com/quasilyte) 





