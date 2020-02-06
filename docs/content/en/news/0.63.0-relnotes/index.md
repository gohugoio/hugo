
---
date: 2020-01-23
title: "Improved base templates, and faster!"
description: "In Hugo 0.63 we have improved the base template lookup logic, and this simplification also made Hugo faster …"
categories: ["Releases"]
---

**Note:** There is already a [patch release](/news/0.63.1-relnotes/) with some fixes.

Hugo `0.63` is, in general, **considerably faster and more memory effective** (see the site-building benchmarks below comparing it to `v0.62`). Not that we are particularly concerned about Hugo's build speed. We leave that to others. But we would hate if it got slower, so we have a comprehensive benchmark suite. And when we needed to simplify the template handling code to solve a concurrency issue, it also became more effective. And as a bonus, we also finally got the [base template lookup order](https://gohugo.io/templates/base/#base-template-lookup-order) that you really, really wanted!

```bash
name                              old time/op    new time/op    delta
SiteNew/Bundle_with_image-16        13.2ms ± 2%    10.7ms ± 0%  -19.13%  (p=0.029 n=4+4)
SiteNew/Bundle_with_JSON_file-16    13.1ms ± 0%    10.8ms ± 0%  -17.50%  (p=0.029 n=4+4)
SiteNew/Tags_and_categories-16      47.7ms ± 1%    43.7ms ± 2%   -8.43%  (p=0.029 n=4+4)
SiteNew/Canonify_URLs-16            52.3ms ± 6%    49.5ms ± 7%     ~     (p=0.200 n=4+4)
SiteNew/Deep_content_tree-16        77.7ms ± 0%    71.6ms ± 1%   -7.84%  (p=0.029 n=4+4)
SiteNew/Many_HTML_templates-16      44.0ms ± 2%    37.5ms ± 1%  -14.79%  (p=0.029 n=4+4)
SiteNew/Page_collections-16         58.4ms ± 1%    52.5ms ± 1%  -10.09%  (p=0.029 n=4+4)

name                              old alloc/op   new alloc/op   delta
SiteNew/Bundle_with_image-16        3.81MB ± 0%    2.22MB ± 0%  -41.57%  (p=0.029 n=4+4)
SiteNew/Bundle_with_JSON_file-16    3.60MB ± 0%    2.01MB ± 0%  -44.09%  (p=0.029 n=4+4)
SiteNew/Tags_and_categories-16      19.3MB ± 1%    14.2MB ± 0%  -26.52%  (p=0.029 n=4+4)
SiteNew/Canonify_URLs-16            70.7MB ± 0%    69.0MB ± 0%   -2.30%  (p=0.029 n=4+4)
SiteNew/Deep_content_tree-16        37.0MB ± 0%    31.2MB ± 0%  -15.78%  (p=0.029 n=4+4)
SiteNew/Many_HTML_templates-16      17.5MB ± 0%    10.6MB ± 0%  -39.68%  (p=0.029 n=4+4)
SiteNew/Page_collections-16         25.8MB ± 0%    21.2MB ± 0%  -17.80%  (p=0.029 n=4+4)

name                              old allocs/op  new allocs/op  delta
SiteNew/Bundle_with_image-16         52.3k ± 0%     26.1k ± 0%  -50.08%  (p=0.029 n=4+4)
SiteNew/Bundle_with_JSON_file-16     52.3k ± 0%     26.1k ± 0%  -50.06%  (p=0.029 n=4+4)
SiteNew/Tags_and_categories-16        337k ± 1%      272k ± 0%  -19.20%  (p=0.029 n=4+4)
SiteNew/Canonify_URLs-16              422k ± 0%      395k ± 0%   -6.33%  (p=0.029 n=4+4)
SiteNew/Deep_content_tree-16          400k ± 0%      314k ± 0%  -21.41%  (p=0.029 n=4+4)
SiteNew/Many_HTML_templates-16        247k ± 0%      143k ± 0%  -41.84%  (p=0.029 n=4+4)
SiteNew/Page_collections-16           282k ± 0%      207k ± 0%  -26.31%  (p=0.029 n=4+4)
```

This release represents **35 contributions by 9 contributors** to the main Hugo code base. [@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@hcwong](https://github.com/hcwong), [@flother](https://github.com/flother), and [@RemcodM](https://github.com/RemcodM) for their ongoing contributions.

And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **13 contributions by 10 contributors**. A special thanks to [@bep](https://github.com/bep), [@jeffscottlevine](https://github.com/jeffscottlevine), [@davidsneighbour](https://github.com/davidsneighbour), and [@nicfits](https://github.com/nicfits) for their work on the documentation site.


Hugo now has:

* 41091+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 440+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 286+ [themes](http://themes.gohugo.io/)

## Notes

* `.Type` on-page now returns an empty string if not set in the front matter or if `.Section` returns empty. See [#6760](https://github.com/gohugoio/hugo/issues/6760).
* Baseof template (e.g. `baseof.html`) lookup order is completely revised/expanded. See [base template lookup order](https://gohugo.io/templates/base/#base-template-lookup-order). We now do template resolution at execution time following the same rules as the template it applies to (e.g. `_default/single.html`). This is an obvious improvement and your site should work as before.
* Shortcode params now supports params with raw string literals (strings surrounded with \`\`) that supports newlines.
* Note: In a base template (e.g. `baseof.html`), the first template block must be a `define`.

## Enhancements

### Templates

* Rework template management to get rid of concurrency issues [c6d650c8](https://github.com/gohugoio/hugo/commit/c6d650c8c8b22fdc7ddedc1e42a3ca698e1390d6) [@bep](https://github.com/bep) [#6716](https://github.com/gohugoio/hugo/issues/6716)[#6760](https://github.com/gohugoio/hugo/issues/6760)[#6768](https://github.com/gohugoio/hugo/issues/6768)[#6778](https://github.com/gohugoio/hugo/issues/6778)
* Put Go's internal template funcs in Hugo's map [1cf23541](https://github.com/gohugoio/hugo/commit/1cf235412f98b42aefe368e99a0e9e95bae6eef7) [@bep](https://github.com/bep) [#6717](https://github.com/gohugoio/hugo/issues/6717)

### Output

* Add base template lookup variant to docs.json [cafb1d53](https://github.com/gohugoio/hugo/commit/cafb1d53c0927e2aef8abff1bf9095c90c6f3067) [@bep](https://github.com/bep) 

### Core

* Disable a test assertion on ARM [836c2426](https://github.com/gohugoio/hugo/commit/836c24261f9f175254256fb326d92a3db47e1c75) [@bep](https://github.com/bep) [#6655](https://github.com/gohugoio/hugo/issues/6655)
* Some more benchmark adjustments [ddd75f21](https://github.com/gohugoio/hugo/commit/ddd75f212110a3d6643a07301e377415f3d163bd) [@bep](https://github.com/bep) 
* Adjust site benchmarks [4ed6ebef](https://github.com/gohugoio/hugo/commit/4ed6ebef4ca71572a19bb890cb4c026a688b2b5b) [@bep](https://github.com/bep) 
* Add a benchmark with lots of templates [ea05c0e8](https://github.com/gohugoio/hugo/commit/ea05c0e8456e8dec71ffd796148355b0d8b36eb0) [@bep](https://github.com/bep) 

### Other

* Regen docs helper [4f466db6](https://github.com/gohugoio/hugo/commit/4f466db666dded1b6c6d1e6926e170f22164433a) [@bep](https://github.com/bep) 
* Allow multiple arguments in ne/ge/gt/le/lt functions Treat op arg1 arg2 arg3 ... as (arg1 op arg2) && (arg1 op arg3) and so on for ne/ge/gt/le/lt. [0c251be6](https://github.com/gohugoio/hugo/commit/0c251be66bf3ad4abafbc47583e394ca4e6ffcf1) [@le0tan](https://github.com/le0tan) [#6619](https://github.com/gohugoio/hugo/issues/6619)
* Update go-org [8585b388](https://github.com/gohugoio/hugo/commit/8585b388d27abde1ab6b6c63ad6addf4066ec8dd) [@niklasfasching](https://github.com/niklasfasching) 
* Add support for newline characters in raw string shortcode [21ca2e9c](https://github.com/gohugoio/hugo/commit/21ca2e9ce4255bfad2bb0576aff087a240acf70a) [@hcwong](https://github.com/hcwong) 
* Update github.com/alecthomas/chroma [3efa1d81](https://github.com/gohugoio/hugo/commit/3efa1d81219a6e7b41c9676e9cab446741f69055) [@ghislainbourgeois](https://github.com/ghislainbourgeois) 
* Update minify to v2.7.2 [65ec8fe8](https://github.com/gohugoio/hugo/commit/65ec8fe827efef5a14c4e1bc440a6df97d2f20a2) [@bep](https://github.com/bep) [#6756](https://github.com/gohugoio/hugo/issues/6756)
* Update Goldmark to v1.1.21 [d3e8ab2e](https://github.com/gohugoio/hugo/commit/d3e8ab2e39dcc27853b163079f4a82364286fe82) [@flother](https://github.com/flother) [#6571](https://github.com/gohugoio/hugo/issues/6571)
* Allow raw string literals in shortcode params [da814556](https://github.com/gohugoio/hugo/commit/da814556567eab9ba0ac5fef5314c3ad5ee50ccd) [@hcwong](https://github.com/hcwong) 
* Update github.com/gohugoio/testmodBuilder [0c0bb372](https://github.com/gohugoio/hugo/commit/0c0bb372858b5de58c15ccd300144e0bc205ffad) [@bep](https://github.com/bep) 
* Update direct dependencies [94cfdf6b](https://github.com/gohugoio/hugo/commit/94cfdf6befd657e46c9458b23f17d851cd2f7037) [@bep](https://github.com/bep) [#6719](https://github.com/gohugoio/hugo/issues/6719)
* Update to new CSS config [45138017](https://github.com/gohugoio/hugo/commit/451380177868e48127a33362aa8d553b90516fb5) [@bep](https://github.com/bep) [#6719](https://github.com/gohugoio/hugo/issues/6719)
* Update to Minify v2.7.0 [56354a63](https://github.com/gohugoio/hugo/commit/56354a63bb73271224a9300a4742dc1a2f551202) [@bep](https://github.com/bep) 
* Add support for freebsd/arm64 [aead8108](https://github.com/gohugoio/hugo/commit/aead8108b80d77e23c68a47fd8d86464310130be) [@dmgk](https://github.com/dmgk) [#6719](https://github.com/gohugoio/hugo/issues/6719)
* Update releasenotes_writer.go [df6e9efd](https://github.com/gohugoio/hugo/commit/df6e9efd8f345707932231ea23dc8713afb5b026) [@bep](https://github.com/bep) 
* Adjust auto ID space handling [9b6e6146](https://github.com/gohugoio/hugo/commit/9b6e61464b09ffe3423fb8d7c72bddb7a9ed5b98) [@bep](https://github.com/bep) [#6710](https://github.com/gohugoio/hugo/issues/6710)
* Document the new autoHeadingIDType setting [d62ede8e](https://github.com/gohugoio/hugo/commit/d62ede8e9e5883e7ebb023e49b82f07b45edc1c7) [@bep](https://github.com/bep) [#6707](https://github.com/gohugoio/hugo/issues/6707)[#6616](https://github.com/gohugoio/hugo/issues/6616)
* Regenerate docshelper [81b7e48a](https://github.com/gohugoio/hugo/commit/81b7e48a55092203aeee8785799e6fed3928760e) [@bep](https://github.com/bep) [#6707](https://github.com/gohugoio/hugo/issues/6707)[#6616](https://github.com/gohugoio/hugo/issues/6616)
* Add an optional Blackfriday auto ID strategy [16e7c112](https://github.com/gohugoio/hugo/commit/16e7c1120346bd853cf6510ffac8e94824bf2c7f) [@bep](https://github.com/bep) [#6707](https://github.com/gohugoio/hugo/issues/6707)
* Make the autoID type config a string [8f071fc1](https://github.com/gohugoio/hugo/commit/8f071fc159ce9a0fc0ea14a73bde8f299bedd109) [@bep](https://github.com/bep) [#6707](https://github.com/gohugoio/hugo/issues/6707)
* markup/goldmark: Simplify code [5ee1f087](https://github.com/gohugoio/hugo/commit/5ee1f0876f3ec8b79d6305298185dc821ead2d28) [@bep](https://github.com/bep) 
* Make auto IDs GitHub compatible [a82d2700](https://github.com/gohugoio/hugo/commit/a82d2700fcc772aada15d65b8f76913ca23f7404) [@bep](https://github.com/bep) [#6616](https://github.com/gohugoio/hugo/issues/6616)
* Support files in content mounts [ff6253bc](https://github.com/gohugoio/hugo/commit/ff6253bc7cf745e9c0127ddc9006da3c2c00c738) [@bep](https://github.com/bep) [#6684](https://github.com/gohugoio/hugo/issues/6684)[#6696](https://github.com/gohugoio/hugo/issues/6696)
* Update alpine base image in Dockerfile to 3.11 [aa4ccb8a](https://github.com/gohugoio/hugo/commit/aa4ccb8a1e9b8aa17397acf34049a2aa16b0b6cb) [@RemcodM](https://github.com/RemcodM) 

## Fixes

### Templates

* Fix eq when > 2 args [2fefc016](https://github.com/gohugoio/hugo/commit/2fefc01606fddb119f368c89fb2dedd452ad6547) [@bep](https://github.com/bep) [#6786](https://github.com/gohugoio/hugo/issues/6786)

### Core

* Fix relative .Page.GetPage from bundle [196a9df5](https://github.com/gohugoio/hugo/commit/196a9df585c4744e3280f37c1c24e469fce14b8c) [@bep](https://github.com/bep) [#6705](https://github.com/gohugoio/hugo/issues/6705)
* Fix inline shortcode regression [5509954c](https://github.com/gohugoio/hugo/commit/5509954c7e8b0ce8d5ea903b0ab639ea14b69acb) [@bep](https://github.com/bep) [#6677](https://github.com/gohugoio/hugo/issues/6677)

### Other

* Fix 0.62.1 server rebuild slowdown regression [17af79a0](https://github.com/gohugoio/hugo/commit/17af79a03e249a731cf5634ffea23ca00774333d) [@bep](https://github.com/bep) [#6784](https://github.com/gohugoio/hugo/issues/6784)
* Fix blog not building [d61bee5e](https://github.com/gohugoio/hugo/commit/d61bee5e0916b5d2b388e66ef85c336312a21a06) [@colonelpopcorn](https://github.com/colonelpopcorn) [#6752](https://github.com/gohugoio/hugo/issues/6752)





