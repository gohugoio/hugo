
---
date: 2017-09-11
title: "Hugo 0.27: Fast and Flexible Related Content!"
description: "Makes it easy to add \"See Also\" sections etc. to your site."
categories: ["Releases"]
images:
- images/blog/hugo-27-poster.png
---

	
Hugo `0.27`comes with fast and flexible **Related Content** ([3b4f17bb](https://github.com/gohugoio/hugo/commit/3b4f17bbc9ff789faa581ac278ad109d1ac5b816) [@bep](https://github.com/bep) [#98](https://github.com/gohugoio/hugo/issues/98)). To add this to your site, put something like this in your single page template:

```go-html-template
{{ $related := .Site.RegularPages.Related . | first 5 }}
{{ with $related }}
<h3>See Also</h3>
<ul>
	{{ range . }}
	<li><a href="{{ .RelPermalink }}">{{ .Title }}</a></li>
	{{ end }}
</ul>
{{ end }}
```

The above translates to _list the five regular pages mostly related to the current page_. See the [Related Content Documentation](https://gohugo.io/content-management/related/) for details and configuration options.

This release represents **37 contributions by 9 contributors** to the main Hugo code base.

[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@moorereason](https://github.com/moorereason), [@yihui](https://github.com/yihui), and [@oneleaftea](https://github.com/oneleaftea) for their ongoing contributions.

And as always a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the documentation and the themes site in pristine condition.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **44 contributions by 30 contributors**. A special thanks to [@bep](https://github.com/bep), [@sdomino](https://github.com/sdomino), [@gotgenes](https://github.com/gotgenes), and [@digitalcraftsman](https://github.com/digitalcraftsman) for their work on the documentation site.


Hugo now has:

* 19464+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 455+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 178+ [themes](http://themes.gohugo.io/)

## Notes

* We now only strip p tag in `markdownify` if there is only one paragraph. This allows blocks of paragraphs to be "markdownified" [33ae10b6](https://github.com/gohugoio/hugo/commit/33ae10b6ade67cd9618970121d7de5fd2ce7d781) [@bep](https://github.com/bep) [#3040](https://github.com/gohugoio/hugo/issues/3040)

## Enhancements

### Templates

* Add `time.Duration` and `time.ParseDuration` template funcs [f4bf2141](https://github.com/gohugoio/hugo/commit/f4bf214137ebd24a0d12f16d3a98d9038e6eabd3) [@bep](https://github.com/bep) [#3828](https://github.com/gohugoio/hugo/issues/3828)
* Add `cond` (ternary) template func [0462c96a](https://github.com/gohugoio/hugo/commit/0462c96a5a9da3e8adc78d96acd39575a8b46c40) [@bep](https://github.com/bep) [#3860](https://github.com/gohugoio/hugo/issues/3860)
* Prepare for template metrics [d000cf60](https://github.com/gohugoio/hugo/commit/d000cf605091c6999b72d6c632752289bc680223) [@bep](https://github.com/bep) 
* Add `strings.TrimLeft` and `TrimRight` [7674ad73](https://github.com/gohugoio/hugo/commit/7674ad73825c61eecc4003475fe0577f225fe579) [@moorereason](https://github.com/moorereason) 
* compare, hugolib, tpl: Add `Eqer` interface [08f48b91](https://github.com/gohugoio/hugo/commit/08f48b91d68d3002b887ddf737456ff0cc4e786d) [@bep](https://github.com/bep) [#3807](https://github.com/gohugoio/hugo/issues/3807)
* Only strip p tag in `markdownify` if only one paragraph [33ae10b6](https://github.com/gohugoio/hugo/commit/33ae10b6ade67cd9618970121d7de5fd2ce7d781) [@bep](https://github.com/bep) [#3040](https://github.com/gohugoio/hugo/issues/3040)
* Cleanup `strings.TrimPrefix` and `TrimSuffix` [29a2da05](https://github.com/gohugoio/hugo/commit/29a2da0593b081cdd61b93c6328af2c9ea4eb20f) [@moorereason](https://github.com/moorereason) 

### Output

* Improve the base template (aka `baseof.html`) identification [0019ce00](https://github.com/gohugoio/hugo/commit/0019ce002449d671a20a69406da37b10977f9493) [@bep](https://github.com/bep) 

### Core

* Implement "related content" [3b4f17bb](https://github.com/gohugoio/hugo/commit/3b4f17bbc9ff789faa581ac278ad109d1ac5b816) [@bep](https://github.com/bep) [#98](https://github.com/gohugoio/hugo/issues/98)
* Add `Page.Equals` [f0f49ed9](https://github.com/gohugoio/hugo/commit/f0f49ed9b0c9b4545a45c95d56340fcbf4aafbef) [@bep](https://github.com/bep) 
* Rewrite `replaceDivider` to reduce memory allocation [71ae9b45](https://github.com/gohugoio/hugo/commit/71ae9b4533083be185c5314c9c5b273cc3bd07bd) [@bep](https://github.com/bep) 


### Other

* Set up Hugo release flow on `CircleCI` [d2249c50](https://github.com/gohugoio/hugo/commit/d2249c50991ba7b00b092aca6e315ca1a4de75a1) [@bep](https://github.com/bep) [#3779](https://github.com/gohugoio/hugo/issues/3779)
* Maintain the scroll position if possible [7231d5a8](https://github.com/gohugoio/hugo/commit/7231d5a829f8d97336a2120afde1260db6ee6541) [@yihui](https://github.com/yihui) [#3824](https://github.com/gohugoio/hugo/issues/3824)
* Add an `iFrame` title to the `YouTube` shortcode [919bc921](https://github.com/gohugoio/hugo/commit/919bc9210a69c801c7304c0b529df93d1dca27aa) [@nraboy](https://github.com/nraboy) 
* Remove the theme submodule from /docs [ea2cc26b](https://github.com/gohugoio/hugo/commit/ea2cc26b390476f1c605405604f8c92afd09b6ee) [@bep](https://github.com/bep) [#3791](https://github.com/gohugoio/hugo/issues/3791)
* Add support for multiple config files via `--config a.toml,b.toml,c.toml` [0f9f73cc](https://github.com/gohugoio/hugo/commit/0f9f73cce5c3f1f05be20bcf1d23b2332623d7f9) [@jgielstra](https://github.com/jgielstra) 
* Render task list item inside `label` for correct accessibility [c8257f8b](https://github.com/gohugoio/hugo/commit/c8257f8b726478ca70dc8984cdcc17b31e4bdc0c) [@danieka](https://github.com/danieka) [#3303](https://github.com/gohugoio/hugo/issues/3303)
* Normalize `UniqueID` between Windows & Linux [0abdeeef](https://github.com/gohugoio/hugo/commit/0abdeeef6740a3cbba0db95374853d040f2022b8) [@Shywim](https://github.com/Shywim) 


## Fixes

### Output

* Fix taxonomy term base template lookup [f88fe312](https://github.com/gohugoio/hugo/commit/f88fe312cb35f7de1615c095edd2f898303dd23b) [@bep](https://github.com/bep) [#3856](https://github.com/gohugoio/hugo/issues/3856)
* Fix `published` front matter handling [202510fd](https://github.com/gohugoio/hugo/commit/202510fdc92d52a20baeaa7edb1091f6882bd95f) [@bep](https://github.com/bep) [#3867](https://github.com/gohugoio/hugo/issues/3867)








