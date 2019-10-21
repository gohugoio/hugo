---
date: 2017-05-22T17:53:58-04:00
categories: ["Releases"]
description: "Hugo 0.21 brings full support for shortcodes per Output Format, the last vital piece of that puzzle"
link: ""
title: "Hugo 0.21"
draft: false
author: bep
aliases: [/0-21/]
---

Hugo `0.21` brings full support for shortcodes per [Output Format](https://gohugo.io/extras/output-formats/) ([#3220](https://github.com/gohugoio/hugo/issues/3220)), the last vital piece of that puzzle. This is especially useful for `Google AMP` with its many custom media tags.

This release represents **126 contributions by 29 contributors** to the main Hugo code base. Since last main release Hugo has **gained 850 stars and 7 additional themes**.

Hugo now has:

* 17156&#43; [stars](https://github.com/gohugoio/hugo/stargazers)
* 457&#43; [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 163&#43; [themes](http://themes.gohugo.io/)

[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@moorereason](https://github.com/moorereason), [@bogem](https://github.com/bogem), and [@munnerz](https://github.com/munnerz) for their ongoing contributions. And as always a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the documentation and the themes site in pristine condition.

## Other Highlights

On a more technical side, [@moorereason](https://github.com/moorereason) and [@bep](https://github.com/bep) have introduced namespaces for Hugo&#39;s many template funcs ([#3042](https://github.com/gohugoio/hugo/issues/3042) ). There are so many now, and adding more into that big pile would be a sure path toward losing control.  Now they are nicely categorised into namespaces with its own tests and examples, with an API that the documentation site can use to make sure it is correct and up-to-date.

## Notes

* The deprecated `.Extension`, `.Now` and `.TargetPath` will now `ERROR` [544f0a63](https://github.com/gohugoio/hugo/commit/544f0a6394b0e085d355e8217fc5bb3d96c12a98) [@bep](https://github.com/bep) 
* The config settings and flags `disable404`, `disableRSS`, `disableSitemap`, `disableRobotsTXT` are now deprecated. Use `disableKinds`. [5794a265](https://github.com/gohugoio/hugo/commit/5794a265b41ffdeebfd8485eecf65cf4088d49d6) [@bep](https://github.com/bep) [#3345](https://github.com/gohugoio/hugo/issues/3345) 

## Enhancements

### Templates

* Log a WARNING on wrong usage of `IsSet` [38661c17](https://github.com/gohugoio/hugo/commit/38661c17bb8c31c9f31ee18f8eba5e3bfddd5574) [@moorereason](https://github.com/moorereason) [#3092](https://github.com/gohugoio/hugo/issues/3092) 
* Add support for ellipsed paginator navigator, making paginators with lots of pages more compact  [b6ea492b](https://github.com/gohugoio/hugo/commit/b6ea492b7a6325d04d44eeb00a990a3a0e29e0c0) [@bep](https://github.com/bep) [#3466](https://github.com/gohugoio/hugo/issues/3466) 
* Add support for interfaces to `intersect` [f1c29b01](https://github.com/gohugoio/hugo/commit/f1c29b017bbd88e701cd5151dd186e868672ef89) [@moorereason](https://github.com/moorereason) [#1952](https://github.com/gohugoio/hugo/issues/1952) 
* Add `NumFmt` function [93b3b138](https://github.com/gohugoio/hugo/commit/93b3b1386714999d716e03b131f77234248f1724) [@moorereason](https://github.com/moorereason) [#1444](https://github.com/gohugoio/hugo/issues/1444) 
* Add template function namespaces [#3418](https://github.com/gohugoio/hugo/issues/3418)  [#3042](https://github.com/gohugoio/hugo/issues/3042)  [@moorereason](https://github.com/moorereason)  [@bep](https://github.com/bep) 
* Add translation links to the default sitemap template [90d3fbf1](https://github.com/gohugoio/hugo/commit/90d3fbf1da93a279cfe994a226ae82cf5441deab) [@rayjolt](https://github.com/rayjolt) [#2569](https://github.com/gohugoio/hugo/issues/2569) 
* Allow text partials in HTML templates and the other way around [1cf29200](https://github.com/gohugoio/hugo/commit/1cf29200b4bb0a9c006155ec76759b7f4b1ad925) [@bep](https://github.com/bep) [#3273](https://github.com/gohugoio/hugo/issues/3273) 

### Output

* Refactor site rendering with an &#34;output format context&#34;. In this release, this is used for shortcode handling only, but this paves the way for future niceness [1e4d082c](https://github.com/gohugoio/hugo/commit/1e4d082cf5b92fedbc60b1b4f0e9d1ee6ec45e33) [@bep](https://github.com/bep) [#3397](https://github.com/gohugoio/hugo/issues/3397)  [2bcbf104](https://github.com/gohugoio/hugo/commit/2bcbf104006e0ec03be4fd500f2519301d460f8c) [@bep](https://github.com/bep) [#3220](https://github.com/gohugoio/hugo/issues/3220) 


### Core

* Handle `shortcode` per `Output Format` [af72db80](https://github.com/gohugoio/hugo/commit/af72db806f2c1c0bf1dfe5832275c41eeba89906) [@bep](https://github.com/bep) [#3220](https://github.com/gohugoio/hugo/issues/3220) 
* Improve shortcode error message [58d9cbd3](https://github.com/gohugoio/hugo/commit/58d9cbd31bcf7c296a39860fd7e566d10faaff28) [@bep](https://github.com/bep) 
* Avoid `index.md` in `/index/index.html` [fea4fd86](https://github.com/gohugoio/hugo/commit/fea4fd86a324bf9679df23f8289887d91b42e919) [@bep](https://github.com/bep) [#3396](https://github.com/gohugoio/hugo/issues/3396) 
* Make missing `GitInfo` a `WARNING` [5ad2f176](https://github.com/gohugoio/hugo/commit/5ad2f17693a9860be76ef8089c8728d2b59d6b04) [@bep](https://github.com/bep) [#3376](https://github.com/gohugoio/hugo/issues/3376) 
* Prevent decoding `pageParam` in common cases [e98f885b](https://github.com/gohugoio/hugo/commit/e98f885b8af27f5473a89d31d0b1f02e61e8a5ec) [@bogem](https://github.com/bogem) 
* Ignore non-source files on partial rebuild [b5b6e81c](https://github.com/gohugoio/hugo/commit/b5b6e81c0269abf9b0f4bc6a127744a25344e5c6) [@xofyarg](https://github.com/xofyarg) [#3325](https://github.com/gohugoio/hugo/issues/3325) 
* Log `WARNING` only on unknown `/data` files [ab692e73](https://github.com/gohugoio/hugo/commit/ab692e73dea3ddfe979c88ee236cc394e47e82f1) [@bep](https://github.com/bep) [#3361](https://github.com/gohugoio/hugo/issues/3361) 
* Avoid processing the same notify event twice [3b677594](https://github.com/gohugoio/hugo/commit/3b67759495c9268c30e6ba2d8c7e3b75d52d2960) [@bep](https://github.com/bep) 
* Only show `rssURI` deprecation `WARNING` if it is actually set [cfd3af8e](https://github.com/gohugoio/hugo/commit/cfd3af8e691119461effa4385251b9d3818e2291) [@bep](https://github.com/bep) [#3319](https://github.com/gohugoio/hugo/issues/3319) 

### Docs

* Add documentation on slug translation [635b3bb4](https://github.com/gohugoio/hugo/commit/635b3bb4eb873978c7d52e6c0cb85da0c4d25299) [@xavib](https://github.com/xavib) 
* Replace `cdn.mathjax.org` with `cdnjs.cloudflare.com` [4b637ac0](https://github.com/gohugoio/hugo/commit/4b637ac041d17b22187f5ccd0f65461f0065aaa9) [@takuti](https://github.com/takuti) 
* Add notes about some output format behaviour [162d3a58](https://github.com/gohugoio/hugo/commit/162d3a586d36cabf6376a76b096fd8b6414487ae) [@jpatters](https://github.com/jpatters) 
* Add `txtpen` as alternative commenting service [7cdc244a](https://github.com/gohugoio/hugo/commit/7cdc244a72de4c08edc0008e37aec83d945dccdf) [@rickyhan](https://github.com/rickyhan) 

### Other

* Embed `Page` in `WeightedPage` [ebf677a5](https://github.com/gohugoio/hugo/commit/ebf677a58360126d8b9a1e98d086aa4279f53181) [@bep](https://github.com/bep) [#3435](https://github.com/gohugoio/hugo/issues/3435) 
* Improve the detection of untranslated strings [a40d1f6e](https://github.com/gohugoio/hugo/commit/a40d1f6ed2aedddc99725658993258cd557640ed) [@bogem](https://github.com/bogem) [#2607](https://github.com/gohugoio/hugo/issues/2607) 
* Make first letter of the Hugo commands flags&#39; usage lowercase [f0f69d03](https://github.com/gohugoio/hugo/commit/f0f69d03c551acb8ac2eeedaad579cf0b596f9ef) [@bogem](https://github.com/bogem) 
* Import `Octopress` image tag in `Jekyll importer` [5f3ad1c3](https://github.com/gohugoio/hugo/commit/5f3ad1c31985450fab8d6772e9cbfcb57cf5cc53) [@buynov](https://github.com/buynov) 

## Fixes

### Templates

*  Do not lower case template names [6d2ea0f7](https://github.com/gohugoio/hugo/commit/6d2ea0f7d7e8a54b8edfc36e52ff74266c30dc27) [@bep](https://github.com/bep) [#3333](https://github.com/gohugoio/hugo/issues/3333) 

### Output

* Fix output format mixup in example [10287263](https://github.com/gohugoio/hugo/commit/10287263f529181d3169668b044cb84e2e3b049a) [@bep](https://github.com/bep) [#3481](https://github.com/gohugoio/hugo/issues/3481) 
* Fix base theme vs project base template logic [077005e5](https://github.com/gohugoio/hugo/commit/077005e514b1ed50d84ceb90c7c72f184cb04521) [@bep](https://github.com/bep) [#3323](https://github.com/gohugoio/hugo/issues/3323) 

### Core
* Render `404` in default language only [154e18dd](https://github.com/gohugoio/hugo/commit/154e18ddb9ad205055d5bd4827c87f3f0daf499f) [@mitchchn](https://github.com/mitchchn) [#3075](https://github.com/gohugoio/hugo/issues/3075) 
* Fix `RSSLink` vs `RSS` `Output Format` [e682fcc6](https://github.com/gohugoio/hugo/commit/e682fcc62233b47cf5bdcaf598ac0657ef089471) [@bep](https://github.com/bep) [#3450](https://github.com/gohugoio/hugo/issues/3450) 
* Add default config for `ignoreFiles`, making that option work when running in server mode [42f4ce15](https://github.com/gohugoio/hugo/commit/42f4ce15a9d68053da36f9efcf7a7d975cc59559) [@chaseadamsio](https://github.com/chaseadamsio) 
* Fix output formats override when no outputs definition given [6e2f2dd8](https://github.com/gohugoio/hugo/commit/6e2f2dd8d3ca61c92a2ee8824fbf05cadef08425) [@bep](https://github.com/bep) [#3447](https://github.com/gohugoio/hugo/issues/3447) 
* Fix handling of zero-length files [0e87b18b](https://github.com/gohugoio/hugo/commit/0e87b18b66d2c8ba9e2abc429630cb03f5b093d6) [@bep](https://github.com/bep) [#3355](https://github.com/gohugoio/hugo/issues/3355) 
* Must recreate `Paginator` on live-reload [45c74526](https://github.com/gohugoio/hugo/commit/45c74526686f6a2afa02bcee767d837d6b9dd028) [@bep](https://github.com/bep) [#3315](https://github.com/gohugoio/hugo/issues/3315) 

### Docs

* Fix incorrect path in `templates/list` [27e88154](https://github.com/gohugoio/hugo/commit/27e88154af2dd9af6d0523d6e67b612e6336f91c) [@MunifTanjim](https://github.com/MunifTanjim) 
* Fixed incorrect specification of directory structure [a28fbca6](https://github.com/gohugoio/hugo/commit/a28fbca6dcfa80b6541f5ef6c8c12cd1804ae9ed) [@TejasQ](https://github.com/TejasQ) 
* Fix `bash` command in `tutorials/github-pages-blog` [c9976155](https://github.com/gohugoio/hugo/commit/c99761555c014e4d041438d5d7e53a6cbaee4492) [@hansott](https://github.com/hansott) 
* Fix `.Data.Pages` range in example [b5e32eb6](https://github.com/gohugoio/hugo/commit/b5e32eb60993b4656918af2c959ae217a68c461e) [@hxlnt](https://github.com/hxlnt) 

### Other

* Fix data race in live-reload close, avoiding some rare panics [355736ec](https://github.com/gohugoio/hugo/commit/355736ec357c81dfb2eb6851ee019d407090c5ec) [@bep](https://github.com/bep) [#2625](https://github.com/gohugoio/hugo/issues/2625) 
* Skip `.git` directories in file scan [94b5be67](https://github.com/gohugoio/hugo/commit/94b5be67fc73b87d114d94a7bb1a33ab997f30f1) [@bogem](https://github.com/bogem) [#3468](https://github.com/gohugoio/hugo/issues/3468) 