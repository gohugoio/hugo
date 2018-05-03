
---
date: 2018-04-23
title: "Hugo 0.40: The Revival of the Shortcodes"
description: "Hugo 0.40: Shortcodes with `.Content` (almost) always available, processed in order of appearance, several new template funcs …"
categories: ["Releases"]
---

Hugo `0.40` is **The Revival of the Shortcodes**. Shortcodes is one of the prime features in Hugo. Really useful, but it has had some known shortcomings. [@bep](https://github.com/bep) has been really busy the last week to address those.

In summary:

* `.Content` for a page retrieved in a query in a shortcode is now _almost_ always available. Note that shortcodes can include content that can include shortcodes that can include content... It is possible to bite your tail. See more below.
* Shortcodes are now processed and rendered in their order of appearance.
* Related to the above, we have now added a zero-based `.Ordinal` to the shortcode.


The first bullet above resolves some surprising behaviour when reading other pages' content from shortcodes. Before this release, that behaviour was undefined. Note that this has never been an issue from regular templates.

It will still not be possible to get **the current shortcode's  page's rendered content**. That would have impressed Einstein.

The new and well defined rules are:

* `.Page.Content` from a shortcode will be empty. The related `.Page.Truncated` `.Page.Summary`, `.Page.WordCount`, `.Page.ReadingTime`, `.Page.Plain` and `.Page.PlainWords` will also have empty values.
* For _other pages_ (retrieved via `.Page.Site.GetPage`, `.Site.Pages` etc.) the `.Content` is there to use as you please as long as you don't have infinite content recursion in your shortcode/content setup. See below.
* `.Page.TableOfContents` is good to go (but does not support shortcodes in headlines; this is unchanged)

If you get into a situation of infinite recursion, the `.Content` will be empty. Run `hugo -v` for more information.

This release represents **19 contributions by 3 contributors** to the main Hugo code base.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **13 contributions by 6 contributors**. A special thanks to [@kaushalmodi](https://github.com/kaushalmodi), [@anthonyfok](https://github.com/anthonyfok), [@bep](https://github.com/bep), and [@regisphilibert](https://github.com/regisphilibert) for their work on the documentation site.


Hugo now has:

* 25071+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 446+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 222+ [themes](http://themes.gohugo.io/)

## Notes

* We have added a `timeout` configuration setting. This is currently only used to time out the `.Content` creation, to bail out of recursive recursions. Run `hugo -v` to see potential warnings about this. The `timeout` is set default to `10000` (10 seconds).

## Enhancements

### Templates

* Add `path.Ext`, `path.Dir` and `path.Base` [47e7788b](https://github.com/gohugoio/hugo/commit/47e7788b3c30de6fb895522096baf2c13598c317) [@bep](https://github.com/bep) 
* Make `fileExist` use the same filesystem as `readFile` [51af1d2e](https://github.com/gohugoio/hugo/commit/51af1d2eadcad89e8c2906c05549352ef69ab016) [@bep](https://github.com/bep) [#4633](https://github.com/gohugoio/hugo/issues/4633)

### Other

* Add `.Page.BundleType` [402f6788](https://github.com/gohugoio/hugo/commit/402f6788ee955ad2aace84e8fba1625db7b356d9) [@bep](https://github.com/bep) [#4662](https://github.com/gohugoio/hugo/issues/4662)
* Add zero-based `Ordinal` to shortcode [3decf4a3](https://github.com/gohugoio/hugo/commit/3decf4a327157e98d3da3502b6d777de63437c39) [@bep](https://github.com/bep) [#3359](https://github.com/gohugoio/hugo/issues/3359)
* Process and render shortcodes in their order of appearance [85535084](https://github.com/gohugoio/hugo/commit/85535084dea4d3e3adf1ebd08ae57b39d76e1904) [@bep](https://github.com/bep) [#3359](https://github.com/gohugoio/hugo/issues/3359)
* Init the content and shortcodes early [19084eaf](https://github.com/gohugoio/hugo/commit/19084eaf74246feac61d618c55031369520dfa8e) [@bep](https://github.com/bep) [#4632](https://github.com/gohugoio/hugo/issues/4632)
* Prepare child page resources before the page itself [3238e14f](https://github.com/gohugoio/hugo/commit/3238e14fdfeedf189a5af122e20bff040ac059bd) [@bep](https://github.com/bep) [#4632](https://github.com/gohugoio/hugo/issues/4632)
* Make .Content (almost) always available in shortcodes [4d26ab33](https://github.com/gohugoio/hugo/commit/4d26ab33dcef704086f43828d1dfb4b8beae2593) [@bep](https://github.com/bep) [#4632](https://github.com/gohugoio/hugo/issues/4632)[#4653](https://github.com/gohugoio/hugo/issues/4653)[#4655](https://github.com/gohugoio/hugo/issues/4655)
* Add language merge support for Pages in resource.Resources [47c05c47](https://github.com/gohugoio/hugo/commit/47c05c47e0b663632a649ee5d256acc1a32fe9e4) [@bep](https://github.com/bep) [#4644](https://github.com/gohugoio/hugo/issues/4644)
* Improve .Content vs shortcodes [e590cc26](https://github.com/gohugoio/hugo/commit/e590cc26eb1363a4b84603f051b20bd43fd1f7bd) [@bep](https://github.com/bep) [#4632](https://github.com/gohugoio/hugo/issues/4632)
* Improve .Get docs [74520d2c](https://github.com/gohugoio/hugo/commit/74520d2cfd39bb4428182e26c57afa9df83ce7b5) [@paulcmal](https://github.com/paulcmal) 
* Update missing positional parameter test for .Get [e2b277bb](https://github.com/gohugoio/hugo/commit/e2b277bba5935c0686cb83f132eae021ef2dc5e1) [@paulcmal](https://github.com/paulcmal) 
* Improve error message in metadata parse [d681ea55](https://github.com/gohugoio/hugo/commit/d681ea55a0a59b7096dacd194ee0cb8fe15b0757) [@bep](https://github.com/bep) [#3696](https://github.com/gohugoio/hugo/issues/3696)
* Add some context to front matter parse error [159bed34](https://github.com/gohugoio/hugo/commit/159bed34c3a850d58d08a36ddc40372ed96af2db) [@bep](https://github.com/bep) [#4638](https://github.com/gohugoio/hugo/issues/4638)
* Updated GetCSV error message (#4636) [5cc944ff](https://github.com/gohugoio/hugo/commit/5cc944ffd77289ab0b8efd69d628fb11d1280993) [@CubeLuke](https://github.com/CubeLuke) 
* .Get doesn't crash on missing positional param fixes #4619 [236f0c84](https://github.com/gohugoio/hugo/commit/236f0c840b45e0c41fcbb2fb6ee556c0fb2d4859) [@paulcmal](https://github.com/paulcmal) 
* fix syntax signature [cd6a2612](https://github.com/gohugoio/hugo/commit/cd6a261242b63555ac2c3ca7a8462b874b490701) [@paulcmal](https://github.com/paulcmal) 







