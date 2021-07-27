
---
date: 2021-02-19
title: "Hugo 0.81.0: The Smorgasbord Edition"
description: "Attribute lists (e.g. CSS classes) for Markdown blocks, newlines in template actions/blocks, native Apple M1 ARM64 binary, it's faster, and more …"
categories: ["Releases"]
toc: true
---

**Hugo 0.81.0** is the first release of this decade, it is the fastest to date, and it's packed with useful new features.

## Newlines in Template Actions and Commands

You can now have newlines within template actions and pipelines. This means that you can now do this and similar:

```go-html-template
{{ dict 
	"country" "Norway" 
	"population" "5 millions"
	"language" "Norwegian"
	"language_code" "nb"
	"weather" "freezing cold"
	"capitol" "Oslo"
	"largest_city" "Oslo"
	"currency"  "Norwegian krone"
	"dialing_code" "+47" 
}}
```

Note that the above construction will fail in Hugo versions < `0.81.0`.

## Attribute Lists after Markdown Blocks

Hugo already supports adding attribute lists (e.g CSS classes) after titles. We now also allow adding attribute lists after Markdown blocks, e.g. tables, lists, paragraphs etc.:

```
> foo
> bar
{.myclass}
```

See [Configure Goldmark](https://gohugo.io/getting-started/configuration-markup#goldmark).

## Performance

This release is the fastest Hugo to date, see details in the benchmarks below. This is [especially true](https://gohugo.io/news/hugo-macos-intel-vs-arm/) if you use the new ARM64 MacOS binary (only works on [Apple M1](https://en.wikipedia.org/wiki/Apple_M1) devices).

### Site Building and Rebuilding Benchmarks: v0.80.0 => v0.81.0

```
name                                      old time/op    new time/op    delta
SiteNew/Edit_Bundle_with_image-16            771µs ± 6%     817µs ± 7%     ~     (p=0.200 n=4+4)
SiteNew/Edit_Bundle_with_JSON_file-16        728µs ± 2%     737µs ± 1%     ~     (p=0.343 n=4+4)
SiteNew/Edit_Tags_and_categories-16         16.6ms ± 5%    16.3ms ± 3%     ~     (p=0.686 n=4+4)
SiteNew/Edit_Canonify_URLs-16               29.4ms ± 6%    26.9ms ± 4%   -8.37%  (p=0.029 n=4+4)
SiteNew/Edit_Deep_content_tree-16           33.8ms ± 3%    31.2ms ± 3%   -7.53%  (p=0.029 n=4+4)
SiteNew/Edit_Many_HTML_templates-16         12.1ms ± 2%    11.6ms ± 1%   -3.94%  (p=0.029 n=4+4)
SiteNew/Edit_Page_collections-16            20.6ms ± 1%    19.8ms ± 0%   -3.57%  (p=0.029 n=4+4)
SiteNew/Edit_List_terms-16                  3.91ms ± 1%    3.81ms ± 2%   -2.52%  (p=0.029 n=4+4)
SiteNew/Regular_Bundle_with_image-16        6.15ms ± 2%    5.53ms ± 2%  -10.11%  (p=0.029 n=4+4)
SiteNew/Regular_Bundle_with_JSON_file-16    6.26ms ± 4%    5.76ms ± 4%   -7.98%  (p=0.029 n=4+4)
SiteNew/Regular_Tags_and_categories-16      26.2ms ± 2%    25.5ms ± 1%   -2.42%  (p=0.029 n=4+4)
SiteNew/Regular_Canonify_URLs-16            34.7ms ± 8%    33.8ms ± 9%     ~     (p=0.486 n=4+4)
SiteNew/Regular_Deep_content_tree-16        43.8ms ± 1%    43.4ms ± 6%     ~     (p=0.343 n=4+4)
SiteNew/Regular_Many_HTML_templates-16      21.5ms ± 1%    19.7ms ± 2%   -8.54%  (p=0.029 n=4+4)
SiteNew/Regular_Page_collections-16         30.7ms ± 2%    28.2ms ± 1%   -8.23%  (p=0.029 n=4+4)
SiteNew/Regular_List_terms-16               9.70ms ± 1%    8.95ms ± 0%   -7.72%  (p=0.029 n=4+4)

name                                      old alloc/op   new alloc/op   delta
SiteNew/Edit_Bundle_with_image-16            437kB ± 0%     428kB ± 0%   -2.02%  (p=0.029 n=4+4)
SiteNew/Edit_Bundle_with_JSON_file-16        216kB ± 0%     207kB ± 0%   -4.20%  (p=0.029 n=4+4)
SiteNew/Edit_Tags_and_categories-16         10.4MB ± 0%     9.7MB ± 0%   -6.08%  (p=0.029 n=4+4)
SiteNew/Edit_Canonify_URLs-16               84.4MB ± 0%    85.2MB ± 0%   +0.87%  (p=0.029 n=4+4)
SiteNew/Edit_Deep_content_tree-16           26.6MB ± 0%    25.6MB ± 0%   -3.57%  (p=0.029 n=4+4)
SiteNew/Edit_Many_HTML_templates-16         6.03MB ± 0%    5.75MB ± 0%   -4.57%  (p=0.029 n=4+4)
SiteNew/Edit_Page_collections-16            14.8MB ± 0%    14.2MB ± 0%   -4.10%  (p=0.029 n=4+4)
SiteNew/Edit_List_terms-16                  1.83MB ± 0%    1.73MB ± 0%   -5.51%  (p=0.029 n=4+4)
SiteNew/Regular_Bundle_with_image-16        1.93MB ± 0%    1.90MB ± 0%   -1.43%  (p=0.029 n=4+4)
SiteNew/Regular_Bundle_with_JSON_file-16    1.72MB ± 0%    1.69MB ± 0%   -1.72%  (p=0.029 n=4+4)
SiteNew/Regular_Tags_and_categories-16      14.3MB ± 0%    13.6MB ± 0%   -4.80%  (p=0.029 n=4+4)
SiteNew/Regular_Canonify_URLs-16            89.7MB ± 0%    90.2MB ± 0%   +0.61%  (p=0.029 n=4+4)
SiteNew/Regular_Deep_content_tree-16        30.4MB ± 0%    29.2MB ± 0%   -3.95%  (p=0.029 n=4+4)
SiteNew/Regular_Many_HTML_templates-16      9.26MB ± 0%    8.94MB ± 0%   -3.47%  (p=0.029 n=4+4)
SiteNew/Regular_Page_collections-16         18.5MB ± 0%    17.7MB ± 0%   -4.25%  (p=0.029 n=4+4)
SiteNew/Regular_List_terms-16               4.00MB ± 0%    3.85MB ± 0%   -3.55%  (p=0.029 n=4+4)

name                                      old allocs/op  new allocs/op  delta
SiteNew/Edit_Bundle_with_image-16            3.99k ± 0%     4.07k ± 0%   +1.80%  (p=0.029 n=4+4)
SiteNew/Edit_Bundle_with_JSON_file-16        3.99k ± 0%     4.06k ± 0%   +1.81%  (p=0.029 n=4+4)
SiteNew/Edit_Tags_and_categories-16           241k ± 0%      245k ± 0%   +1.67%  (p=0.029 n=4+4)
SiteNew/Edit_Canonify_URLs-16                 364k ± 0%      321k ± 0%  -11.78%  (p=0.029 n=4+4)
SiteNew/Edit_Deep_content_tree-16             264k ± 0%      268k ± 0%   +1.53%  (p=0.029 n=4+4)
SiteNew/Edit_Many_HTML_templates-16          90.3k ± 0%     91.0k ± 0%   +0.83%  (p=0.029 n=4+4)
SiteNew/Edit_Page_collections-16              153k ± 0%      156k ± 0%   +1.46%  (p=0.029 n=4+4)
SiteNew/Edit_List_terms-16                   30.4k ± 0%     30.9k ± 0%   +1.54%  (p=0.029 n=4+4)
SiteNew/Regular_Bundle_with_image-16         23.3k ± 0%     23.1k ± 0%   -0.57%  (p=0.029 n=4+4)
SiteNew/Regular_Bundle_with_JSON_file-16     23.3k ± 0%     23.1k ± 0%   -0.59%  (p=0.029 n=4+4)
SiteNew/Regular_Tags_and_categories-16        284k ± 0%      288k ± 0%   +1.58%  (p=0.029 n=4+4)
SiteNew/Regular_Canonify_URLs-16              387k ± 0%      343k ± 0%  -11.41%  (p=0.029 n=4+4)
SiteNew/Regular_Deep_content_tree-16          307k ± 0%      309k ± 0%   +0.52%  (p=0.029 n=4+4)
SiteNew/Regular_Many_HTML_templates-16        129k ± 0%      129k ± 0%   +0.35%  (p=0.029 n=4+4)
SiteNew/Regular_Page_collections-16           199k ± 0%      200k ± 0%   +0.55%  (p=0.029 n=4+4)
SiteNew/Regular_List_terms-16                53.5k ± 0%     53.4k ± 0%   -0.18%  (p=0.029 n=4+4)
```



## Native Arm Binary for Apple M1

We finally provide native Hugo binary for [Apple M1](https://en.wikipedia.org/wiki/Apple_M1) devices. Download the file named `hugo_0.81.0_macOS-ARM64.tar.gz`.

## JavaScript Building

* Add `inject` config option allowing you to automatically replace a global variable with an import from another file relative to `assets`. [32b86076](https://github.com/gohugoio/hugo/commit/32b86076ee1c0833b538b84e1cc9e6d79babecf2) [@bep](https://github.com/bep) [#8164](https://github.com/gohugoio/hugo/issues/8164)
* Add `shims` config option to swap out a component with another. A common use case is to load dependencies like React from a CDN  (with _shims_) when in production, but running with the full bundled `node_modules` dependency during development [e19a046c](https://github.com/gohugoio/hugo/commit/e19a046c4be9b0654884259b9df94f41561e4fc3) [@bep](https://github.com/bep) [#8165](https://github.com/gohugoio/hugo/issues/8165)
* Add external source map support to js.Build and Babel [2c8b5d91](https://github.com/gohugoio/hugo/commit/2c8b5d9165011c4b24b494e661ae60dfc7bb7d1b) [@richtera](https://github.com/richtera) [#8132](https://github.com/gohugoio/hugo/issues/8132)
* Fix nilpointer in js.Build error handling [a1fe552f](https://github.com/gohugoio/hugo/commit/a1fe552fc9e622a15010a94281f604eb85bebd84) [@bep](https://github.com/bep) [#8162](https://github.com/gohugoio/hugo/issues/8162)

Also see [js.Build Options](https://gohugo.io/hugo-pipes/js#options).

## Hugo Modules

There are several [Hugo Modules](https://gohugo.io/hugo-modules/)-related improvements in this release:

* Allow absolute paths for any modules resolved via project replacement [3a5ee0d2](https://github.com/gohugoio/hugo/commit/3a5ee0d2d6e344b12efc7a97354ec3480c4c578b) [@bep](https://github.com/bep) [#8240](https://github.com/gohugoio/hugo/issues/8240)
* Add config option modules.vendorClosest [bdfbcf6f](https://github.com/gohugoio/hugo/commit/bdfbcf6f4b4ab53a617ab76f72e8aa28da6067de) [@bep](https://github.com/bep) [#8235](https://github.com/gohugoio/hugo/issues/8235)[#8242](https://github.com/gohugoio/hugo/issues/8242)
* Throw an error running hugo mod vendor on mountless module [4ffaeaf1](https://github.com/gohugoio/hugo/commit/4ffaeaf15536596c94dc73b393ca7894e3bd5e2c) [@bep](https://github.com/bep) 

## Minify - Keep Comments

Keep comments when running `hugo --minify` with a new setting:

{{< code-toggle file="config" >}}
[minify.tdewolff.html]
keepComments = true
{{< /code-toggle >}}

The default value for this setting is `false`.

## Statistics

This release represents **59 contributions by 14 contributors** to the main Hugo code base.[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@moorereason](https://github.com/moorereason) and [@benmezger](https://github.com/benmezger) for their ongoing contributions. And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour), [@coliff](https://github.com/coliff) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **20 contributions by 13 contributors**. A special thanks to [@bep](https://github.com/bep), [@gagarine](https://github.com/gagarine), [@fridde](https://github.com/fridde), and [@NicoHood](https://github.com/NicoHood) for their work on the documentation site.


Hugo now has:

* 50152+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 435+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 370+ [themes](http://themes.gohugo.io/)


## Notes

* We have updated to Beta 6 of the Dart Sass Protocol which is not backwards compatible, so if you use Dart Sass you need to also update [that binary](https://gohugo.io/hugo-pipes/scss-sass/#options).
* `hugo gen autocomplete` now default to `stdout`; you can change this by setting `--completionfile`. As an added bonus we now also support auto completion for zsh, fish and powershell.

## Changelog

* Make the build green again [fe77f743](https://github.com/gohugoio/hugo/commit/fe77f7434bc0d7a9b54af69014eb28dbea2b236b) [@bep](https://github.com/bep) 
* Regenerate internal templates [c6080655](https://github.com/gohugoio/hugo/commit/c60806550a21062936b0d02708c9a8c240cafa9d) [@bep](https://github.com/bep) 
* Update date logic of opengraph and schema internal templates [ffd9dac4](https://github.com/gohugoio/hugo/commit/ffd9dac4218b8f1709de04f7131ca661715fc481) [@djatwood](https://github.com/djatwood) 
* Synch Go templates fork with Go 1.16dev [cf3e077d](https://github.com/gohugoio/hugo/commit/cf3e077da304e6f4d7c22f8296e1382335d055c6) [@bep](https://github.com/bep) 
* Exclude pages without Permalink from sitemap [4867cd1d](https://github.com/gohugoio/hugo/commit/4867cd1dea34ee53fb73cede2bcff4792e470104) [@Jaza](https://github.com/Jaza) 
* Add default user-agent header for getJSON requests [35def0ae](https://github.com/gohugoio/hugo/commit/35def0ae4560bb86febd12663bf5602485ad4b20) [@peacecwz](https://github.com/peacecwz) 
* remove 1mb limit for readFile. [ee9c1367](https://github.com/gohugoio/hugo/commit/ee9c1367635eab446fcf9baa1ab8b4066882548e) [@avdva](https://github.com/avdva) 
* Do not return errors in substr for out-of-bounds cases [8a26ab0b](https://github.com/gohugoio/hugo/commit/8a26ab0bc5dd9fa34e1362681fc08b0e522cd4ea) [@moorereason](https://github.com/moorereason) [#8113](https://github.com/gohugoio/hugo/issues/8113)
* Add missing test scenario for strings.Substr [788e50ad](https://github.com/gohugoio/hugo/commit/788e50ad3a55609ed49ce0b7ee98965c181fe9cf) [@moorereason](https://github.com/moorereason) 
* Regen CLI docs [9e99950c](https://github.com/gohugoio/hugo/commit/9e99950c6ebf82c85ee52a8de85e284a506a2f90) [@bep](https://github.com/bep) 
* Regen docs helper [1b364b00](https://github.com/gohugoio/hugo/commit/1b364b003f68df3adb1644769fe69810d85e3897) [@bep](https://github.com/bep) 
* Run go mod tidy [88b93a09](https://github.com/gohugoio/hugo/commit/88b93a09dc79518d7fbd14681eeeea3411dab1dd) [@bep](https://github.com/bep) 
* Add arm64 to Darwinextended build and add vendorInfo [29fb456c](https://github.com/gohugoio/hugo/commit/29fb456c9e63ee1a2314bf4a7227a5146e7f9b31) [@bep](https://github.com/bep) [#8003](https://github.com/gohugoio/hugo/issues/8003)
* Update Travis, GitHub, CircleCI and Snap to Go 1.16 (only) [718fba7d](https://github.com/gohugoio/hugo/commit/718fba7d63424017cb3b9774c33e7acc69c68af6) [@bep](https://github.com/bep) 
* Pull in latest Go 1.16 template source [e77b2e3a](https://github.com/gohugoio/hugo/commit/e77b2e3aa0b24c5ca960905772335b27845705eb) [@bep](https://github.com/bep) 
* Add breaking tests for "map read and map write in templates" [b5485aea](https://github.com/gohugoio/hugo/commit/b5485aeae7e1f73f18835fbf0b8eedc305d450d0) [@bep](https://github.com/bep) [#7293](https://github.com/gohugoio/hugo/issues/7293)
* Pull in latest Go template source [ccb822eb](https://github.com/gohugoio/hugo/commit/ccb822eb5afad210432eb46ec3727e3536a87f58) [@bep](https://github.com/bep) 
* Expand template newline testcase to commands [21e9eb18](https://github.com/gohugoio/hugo/commit/21e9eb18acc2a2f8d8b97f096615b836e65091a2) [@bep](https://github.com/bep) 
* Add a test case for Go 1.16 template action newlines [ae57ba6a](https://github.com/gohugoio/hugo/commit/ae57ba6a9dee87347fa2d5e8c6865f390989622e) [@bep](https://github.com/bep) 
* Update github.com/tdewolff/minify/v2 v2.6.2 => v2.9.13 [66beac99](https://github.com/gohugoio/hugo/commit/66beac99c64b5e5fe7bec0bda437ba5858d49a36) [@bep](https://github.com/bep) [#8258](https://github.com/gohugoio/hugo/issues/8258)
* bump github.com/frankban/quicktest from 1.11.2 to 1.11.3 [968dd7a7](https://github.com/gohugoio/hugo/commit/968dd7a711063934af84bd1c017c58a1e66f51bb) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/getkin/kin-openapi from 0.32.0 to 0.39.0 [38f29e81](https://github.com/gohugoio/hugo/commit/38f29e817f2058ed56f96fb8e628315f3ab5d7f9) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/aws/aws-sdk-go from 1.36.33 to 1.37.11 [cd87813a](https://github.com/gohugoio/hugo/commit/cd87813aa0327ec7a7e6f023dadcea5a3e6a9fef) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/sanity-io/litter from 1.3.0 to 1.5.0 [4e815b06](https://github.com/gohugoio/hugo/commit/4e815b063e4af70f21b6796688025675253bec65) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/olekukonko/tablewriter from 0.0.4 to 0.0.5 [652a59d3](https://github.com/gohugoio/hugo/commit/652a59d38523e23e39376cba9c554abbe87b198d) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Update to esbuild v0.8.46 [84f0ec7f](https://github.com/gohugoio/hugo/commit/84f0ec7f80855dcc9b123418bcbf816b5efa2cdf) [@bep](https://github.com/bep) 
* Add config option modules.vendorClosest [bdfbcf6f](https://github.com/gohugoio/hugo/commit/bdfbcf6f4b4ab53a617ab76f72e8aa28da6067de) [@bep](https://github.com/bep) [#8235](https://github.com/gohugoio/hugo/issues/8235)[#8242](https://github.com/gohugoio/hugo/issues/8242)
* bump google.golang.org/api from 0.26.0 to 0.40.0 [a9b0fea6](https://github.com/gohugoio/hugo/commit/a9b0fea6a3aec658912a8db134824dee4a9b6369) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Change version string format and add VendorInfo to help with issue triaging [e8df0977](https://github.com/gohugoio/hugo/commit/e8df09774534abe6131eb455b4f5c614fb438983) [@anthonyfok](https://github.com/anthonyfok) 
* Allow absolute paths for any modules resolved via project replacement [3a5ee0d2](https://github.com/gohugoio/hugo/commit/3a5ee0d2d6e344b12efc7a97354ec3480c4c578b) [@bep](https://github.com/bep) [#8240](https://github.com/gohugoio/hugo/issues/8240)
* Throw an error running hugo mod vendor on mountless module [4ffaeaf1](https://github.com/gohugoio/hugo/commit/4ffaeaf15536596c94dc73b393ca7894e3bd5e2c) [@bep](https://github.com/bep) 
* Add PowerShell completion support [5f621df2](https://github.com/gohugoio/hugo/commit/5f621df2570236a08cd21e8dd1c60502ec3db328) [@anthonyfok](https://github.com/anthonyfok) [#8122](https://github.com/gohugoio/hugo/issues/8122)
* Refer to mage instead of make in comment regarding commitHash [7118f89c](https://github.com/gohugoio/hugo/commit/7118f89cf35246767e26dcb5e747469ffa61f473) [@anthonyfok](https://github.com/anthonyfok) 
* Add attributes support for blocks (tables etc.) [2681633d](https://github.com/gohugoio/hugo/commit/2681633db8d340d2dc59cf801419874d572fc704) [@bep](https://github.com/bep) [#7548](https://github.com/gohugoio/hugo/issues/7548)
* Update to Goldmark v1.3.2 [1b247282](https://github.com/gohugoio/hugo/commit/1b2472825664763c0b88807b0d193e73553423ec) [@bep](https://github.com/bep) [#8143](https://github.com/gohugoio/hugo/issues/8143)
* Update to Dart Sass Protocol beta6 [441b11be](https://github.com/gohugoio/hugo/commit/441b11beec3cf0371ff9a2898f220a0bf00faf8c) [@bep](https://github.com/bep) 
* Write to stdout by default [d36fd5b3](https://github.com/gohugoio/hugo/commit/d36fd5b3ee6989203de2a29b1de67521fd1c8ea5) [@benmezger](https://github.com/benmezger) 
* Remove powershell support [a7c515e1](https://github.com/gohugoio/hugo/commit/a7c515e1b56e8cab34ca2647b4116904df9c8250) [@benmezger](https://github.com/benmezger) 
* Add zsh, fish and powershell completion support [216b00f3](https://github.com/gohugoio/hugo/commit/216b00f358dbfa36b34ff515d7f4f88387156db8) [@benmezger](https://github.com/benmezger) [#4296](https://github.com/gohugoio/hugo/issues/4296)
* Enable NPM tests on Windows [14494379](https://github.com/gohugoio/hugo/commit/144943798c2a199ed256ae901a14d3c918055eba) [@bep](https://github.com/bep) [#8196](https://github.com/gohugoio/hugo/issues/8196)
* Update to esbuild v0.8.39 [440fdb0e](https://github.com/gohugoio/hugo/commit/440fdb0eb96b3230ddefee732b0c1afe52a37228) [@bep](https://github.com/bep) [#8189](https://github.com/gohugoio/hugo/issues/8189)
* Trim whitespace in elements written to hugo_stats.json [b2a48dce](https://github.com/gohugoio/hugo/commit/b2a48dce58abd3a661aa198af3277ef12f44cce0) [@pmatiash](https://github.com/pmatiash) [#7958](https://github.com/gohugoio/hugo/issues/7958)
* bump github.com/aws/aws-sdk-go from 1.35.0 to 1.36.33 [2f9dadae](https://github.com/gohugoio/hugo/commit/2f9dadae4072960bbaec3656347e20eec238288c) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Remove mention of a file size limit for readFile [ed3071b7](https://github.com/gohugoio/hugo/commit/ed3071b753c8dec83a2c054624e49b204553ecd3) [@avdva](https://github.com/avdva) 
* Add Inject config option [32b86076](https://github.com/gohugoio/hugo/commit/32b86076ee1c0833b538b84e1cc9e6d79babecf2) [@bep](https://github.com/bep) [#8164](https://github.com/gohugoio/hugo/issues/8164)
* Add Shims option [e19a046c](https://github.com/gohugoio/hugo/commit/e19a046c4be9b0654884259b9df94f41561e4fc3) [@bep](https://github.com/bep) [#8165](https://github.com/gohugoio/hugo/issues/8165)
* bump github.com/spf13/afero from 1.4.1 to 1.5.1 [07ad283f](https://github.com/gohugoio/hugo/commit/07ad283f686904e5835f621d73ed342ba2a48eb3) [@eclipseo](https://github.com/eclipseo) 
* Add external source map support to js.Build and Babel [2c8b5d91](https://github.com/gohugoio/hugo/commit/2c8b5d9165011c4b24b494e661ae60dfc7bb7d1b) [@richtera](https://github.com/richtera) [#8132](https://github.com/gohugoio/hugo/issues/8132)
* Run go mod tidy [4d2b6fc4](https://github.com/gohugoio/hugo/commit/4d2b6fc4c0e714f3f1ed345d6d75ed1662948791) [@bep](https://github.com/bep) 
* Update go-org to v1.4.0 [212e5e55](https://github.com/gohugoio/hugo/commit/212e5e554284bc9368e52a512ed09be5a0224d3e) [@niklasfasching](https://github.com/niklasfasching) 
* Adjust log level [4fdec67b](https://github.com/gohugoio/hugo/commit/4fdec67b1155ae1cdf051582d9ab387286b71a07) [@bep](https://github.com/bep) 
* Add temporary patch to fix template data race [9650e568](https://github.com/gohugoio/hugo/commit/9650e568418a316e71ad94d7e27caf544a4a2d0d) [@bep](https://github.com/bep) [#7293](https://github.com/gohugoio/hugo/issues/7293)
* Fix race condition in text template baseof [241b7483](https://github.com/gohugoio/hugo/commit/241b7483ea954653512d4895ad6bacf79ee26ddc) [@moorereason](https://github.com/moorereason) 
* Fix metrics hint tracking [0004a733](https://github.com/gohugoio/hugo/commit/0004a733c85cee991a8a170e93cd69c326cc8f2f) [@moorereason](https://github.com/moorereason) [#8125](https://github.com/gohugoio/hugo/issues/8125)
* Fix potential path issue on Windows [b60e9279](https://github.com/gohugoio/hugo/commit/b60e9279ab95030828eb4f822be96250284c4d8d) [@bep](https://github.com/bep) 
* Fix some humanize issues [bf55afd7](https://github.com/gohugoio/hugo/commit/bf55afd71f2fdb47272ebf1188c9cc87df47b233) [@susiwen8](https://github.com/susiwen8) [#7912](https://github.com/gohugoio/hugo/issues/7912)
* Fix handling of legacy attribute config [e6dd3128](https://github.com/gohugoio/hugo/commit/e6dd312812c7c711986af2d60f2999d116b82ea0) [@bep](https://github.com/bep) [#7548](https://github.com/gohugoio/hugo/issues/7548)
* Support translation files with suffix *.yml [92c6c404](https://github.com/gohugoio/hugo/commit/92c6c40419bdc13b8bb422a212d1d79240356651) [@bep](https://github.com/bep) [#8212](https://github.com/gohugoio/hugo/issues/8212)
* Fix nilpointer in js.Build error handling [a1fe552f](https://github.com/gohugoio/hugo/commit/a1fe552fc9e622a15010a94281f604eb85bebd84) [@bep](https://github.com/bep) [#8162](https://github.com/gohugoio/hugo/issues/8162)



