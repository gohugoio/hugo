
---
date: 2017-10-16
title: "Hugo 0.30: Race Car Edition!"
description: "Fast Render Mode boosts live reloading!"
categories: ["Releases"]
images:
- images/blog/hugo-30-poster.png
---

	
Hugo `0.30` is the **Race Car Edition**. Hugo is already very very fast, but wants much more. So we added **Fast Render Mode**. It is hard to explain, so start the Hugo development server with `hugo server` and start editing. Live reloads just got so much faster! The "how and what" is discussed at length in [other places](https://github.com/gohugoio/hugo/pull/3959), but the short version is that we now re-render only the parts of the site that you are working on.

The second performance-related feature is a follow-up to the Template Metrics added in Hugo `0.29`. Now, if you add the flag `--templateMetricsHints`, we will calculate a score for how your partials can be cached (with the `partialCached` template func).

This release also more or less makes the really fast Chroma highlighter a complete alternative to Pygments. Most notable is the new table `linenos` support ([7c30e2cb](https://github.com/gohugoio/hugo/commit/7c30e2cbb08fdf0e61f80c7f1aa29909aeca4211) [@bep](https://github.com/bep) [#3915](https://github.com/gohugoio/hugo/issues/3915)), which makes copy-and-paste code blocks much easier.

This release represents **31 contributions by 10 contributors** to the main Hugo code base.
[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contribution, but also a big shoutout to [@moorereason](https://github.com/moorereason), [@digitalcraftsman](https://github.com/digitalcraftsman), and [@bmon](https://github.com/bmon) for their ongoing contributions.
And as always a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the documentation and the themes site in pristine condition.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **26 contributions by 15 contributors**. A special thanks to [@bep](https://github.com/bep), [@digitalcraftsman](https://github.com/digitalcraftsman), [@moorereason](https://github.com/moorereason), and [@kaushalmodi](https://github.com/kaushalmodi) for their work on the documentation site.

Hugo now has:

* 20195+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 454+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 180+ [themes](http://themes.gohugo.io/)

## Notes

* Running `hugo server` will now run with the new "Fast Render Mode" default on. To turn it off, run `hugo server --disableFastRender` or set `disableFastRender=true` in your site config.
* There have been several fixes and enhancements in the Chroma highlighter. One is that it now creates Pygments compatible CSS classes, which means that you may want to re-generate the stylesheet. See the [Syntax Highlighting Doc](https://gohugo.io/content-management/syntax-highlighting/).

## Enhancements

### Performance
* Only re-render the view(s) you're working on [60bd332c](https://github.com/gohugoio/hugo/commit/60bd332c1f68e49e6ac439047e7c660865189380) [@bep](https://github.com/bep) [#3962](https://github.com/gohugoio/hugo/issues/3962)
* Detect `partialCached` candidates [5800a20a](https://github.com/gohugoio/hugo/commit/5800a20a258378440e203a6c4a4343f5077755df) [@bep](https://github.com/bep) 
* Move metrics output to the end of the site build [b277cb33](https://github.com/gohugoio/hugo/commit/b277cb33e4dfa7440fca3b7888026944ce056154) [@moorereason](https://github.com/moorereason) 

### Templates

* Output `xmlns:xhtml` only if there are translations available [0859d9df](https://github.com/gohugoio/hugo/commit/0859d9dfe647db3b8a192da38ad7efb5480a29a1) [@jamieconnolly](https://github.com/jamieconnolly) 
* Add `errorf` template function [4fc67fe4](https://github.com/gohugoio/hugo/commit/4fc67fe44a3c65fc7faaed21d5fa5bb5f87edf2c) [@bmon](https://github.com/bmon) [#3817](https://github.com/gohugoio/hugo/issues/3817)
* Add `os.FileExists` template function [28188789](https://github.com/gohugoio/hugo/commit/2818878994e906c292cbe00cb2a83f1531a21f32) [@digitalcraftsman](https://github.com/digitalcraftsman) [#3839](https://github.com/gohugoio/hugo/issues/3839)
* Add `float` template function [57adc539](https://github.com/gohugoio/hugo/commit/57adc539fc98dcb6fba8070b9611b8bd545f6f7f) [@x3ro](https://github.com/x3ro) [#3307](https://github.com/gohugoio/hugo/issues/3307)
* Rework the partial test and benchmarks [e2e8bcbe](https://github.com/gohugoio/hugo/commit/e2e8bcbec34702a27047b91b6b007a15f1fc0797) [@bep](https://github.com/bep) 

### Other

* Change `SummaryLength` to be configurable (#3924) [8717a60c](https://github.com/gohugoio/hugo/commit/8717a60cc030f4310c1779c0cdd51db37ad636cd) [@bmon](https://github.com/bmon) [#3734](https://github.com/gohugoio/hugo/issues/3734)
* Replace `make` with `mage` in CircleCI build [fe71cb6f](https://github.com/gohugoio/hugo/commit/fe71cb6f5f83cdc8374cf1fc35a6d48102bd4b12) [@bep](https://github.com/bep) [#3969](https://github.com/gohugoio/hugo/issues/3969)
* Add table `linenos` support for Chroma highlighter [7c30e2cb](https://github.com/gohugoio/hugo/commit/7c30e2cbb08fdf0e61f80c7f1aa29909aeca4211) [@bep](https://github.com/bep) [#3915](https://github.com/gohugoio/hugo/issues/3915)
* Replace `make` with `mage` [8d2580f0](https://github.com/gohugoio/hugo/commit/8d2580f07c0253e12524a4b5c13165f876d00b21) [@bep](https://github.com/bep) [#3937](https://github.com/gohugoio/hugo/issues/3937)
* Create `magefile` from `Makefile` [384a6ac4](https://github.com/gohugoio/hugo/commit/384a6ac4bd2de16fcd6a1c952e7ca41b66023a12) [@natefinch](https://github.com/natefinch) 
* Clean up lint in various packages [47fdfd51](https://github.com/gohugoio/hugo/commit/47fdfd5196cd24a23b30afe1d88969ffb413ab59) [@moorereason](https://github.com/moorereason) 

## Fixes

* Make sure `Date` and `PublishDate` are always set to a value if one is available [6a30874f](https://github.com/gohugoio/hugo/commit/6a30874f19610a38e846e120aac03c68e12f9b7b) [@bep](https://github.com/bep) [#3854](https://github.com/gohugoio/hugo/issues/3854)
* Add correct config file name to verbose server log [15ec031d](https://github.com/gohugoio/hugo/commit/15ec031d9818d239bfbff525c00cd99cc3118a96) [@mdhender](https://github.com/mdhender) 
