
---
date: 2018-07-22
title: "Hugo 0.45: Revival of ref, relref and GetPage"
description: "Hugo 0.45 adds relative page lookups, language support in ref/relref and several Hugo Pipes improvements."
categories: ["Releases"]
---

	
Hugo `0.45` is the **revival of ref, relref and GetPage**. [@vassudanagunta](https://github.com/vassudanagunta) and [@bep](https://github.com/bep) have done some great work improving the API and implementation for the helper functions used to **get one page**. Before this release, the API was a little bit clumsy and the result potentially ambiguous in some situations.

Now you can simply do:

```go-html-template
{{ with .Site.GetPage "/blog/my-post.md" }}{{ .Title }}{{ end }}
```

Or to get a section page:


```go-html-template
{{ with .Site.GetPage "/blog" }}{{ .Title }}{{ end }}
```

We have also added a `.GetPage` method on `Page` and added support for page-relative linking. This means that the leading slash (`/`) now has a meaning. For `.Site.GetPage`, all lookups will start at the content root. But for lookups with a `Page` context, paths without a leading slash will be treated as relative to the page.

This means that the following example will find the page in the current section:

```go-html-template
{{</* ref "my-post.md" */>}}
```

You can also use the `..` to refer to a page one level up etc.:

```go-html-template
{{</* ref "../my-post.md" */>}}
```

We have now also added language support to `ref` and `relref`, so you can link to a page in another language:

```go-html-template
{{</* relref path="document.md" lang="ja" */>}}
```

To link to a given Output Format of a document, you can use this syntax:

```go-html-template
{{</* relref path="document.md" outputFormat="rss" */>}}
```

To make working with these reflinks on bigger sites easier to work with, we have also improved the error logging, and added two new configuration settings:

* refLinksErrorLevel: ERROR (default, will fail the build when a reflink cannot be resolved) or WARNING.
* refLinksNotFoundURL: Set this to an URL placeholder used when no reference could be resolved.

Visit the [Hugo Docs](https://gohugo.io/content-management/cross-references) for more information.

We have also done some important improvements and fixes in **Hugo Pipes** in this release: SCSS source maps on Windows now works, we now support project-local `PostCSS` installation, and we have added `IncludePaths` to `SCSS` options, making it possible to include, say, a path below `node_modules` in the SASS/SCSS build.

This release represents **31 contributions by 4 contributors** to the main Hugo code base.
[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@vassudanagunta](https://github.com/vassudanagunta), [@hairmare](https://github.com/hairmare), and [@garrmcnu](https://github.com/garrmcnu) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the themes site in pristine condition and to [@kaushalmodi](https://github.com/kaushalmodi) for his great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **10 contributions by 8 contributors**. A special thanks to [@kaushalmodi](https://github.com/kaushalmodi), [@Hanzei](https://github.com/Hanzei), [@KurtTrowbridge](https://github.com/KurtTrowbridge), and [@regisphilibert](https://github.com/regisphilibert) for their work on the documentation site.


Hugo now has:

* 27334+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 443+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 238+ [themes](http://themes.gohugo.io/)

## Notes
* `.Site.GetPage` with more than 2 arguments will not work anymore. This means that `{{ .Site.GetPage "page" "blog" "my-post.md" }}` will fail. `{{ .Site.GetPage "page" "blog/my-post.md" }}` will work, but we recommend you use the simpler `{{ .Site.GetPage "/blog/my-post.md" }}`
* Relative paths in `relref` or `ref` that finds its match not relative to the page itself will work, but we now print a warning saying that you should correct it to an absolute path. E.g. `{{</* ref "blog/my-post.md" */>}}` => `{{</* ref "/blog/my-post.md" */>}}`.

## Enhancements

* Print a WARNING about relative non-relative ref/relref matches [a451c49f](https://github.com/gohugoio/hugo/commit/a451c49fde1da6e2cc436a2b7d383ee772b1f893) [@bep](https://github.com/bep) [#4973](https://github.com/gohugoio/hugo/issues/4973)
* Allow untyped nil to be merged in lang.Merge [ff16c42e](https://github.com/gohugoio/hugo/commit/ff16c42ed0965e1c8acf6e6a6dcda3ea50c107f2) [@bep](https://github.com/bep) [#4977](https://github.com/gohugoio/hugo/issues/4977)
* Get rid of the utils package [062510cf](https://github.com/gohugoio/hugo/commit/062510cf1f7b79aed2efe88c5b9340d009bdec0e) [@bep](https://github.com/bep) 
* Update hugo_windows.go [4e1d0cd9](https://github.com/gohugoio/hugo/commit/4e1d0cd9f1d43d133d669a019a84117cadd41955) [@bep](https://github.com/bep) 
* Add IncludePaths config option [166483fe](https://github.com/gohugoio/hugo/commit/166483fe1227b0c59c6b4d88cfdfaf7d7b0d79c5) [@bep](https://github.com/bep) [#4921](https://github.com/gohugoio/hugo/issues/4921)
* Increase refLinker test coverage [8278384b](https://github.com/gohugoio/hugo/commit/8278384b9680cfdcecef9c668638ad483012857f) [@vassudanagunta](https://github.com/vassudanagunta) 
* Add test coverage for recent ref overhaul [2bac3715](https://github.com/gohugoio/hugo/commit/2bac3715448e90e197ada7cc73c87f696c19def6) [@vassudanagunta](https://github.com/vassudanagunta) [#4969](https://github.com/gohugoio/hugo/issues/4969)
* Update ref, relref, GetPage docs [1eb8b36b](https://github.com/gohugoio/hugo/commit/1eb8b36b3802e72bc2c16965461ef1899bb073b3) [@bep](https://github.com/bep) 
* Document refLinksErrorLevel and refLinksNotFoundURL [00c74ee7](https://github.com/gohugoio/hugo/commit/00c74ee7ffae71fd5f47d555160354a775e26151) [@bep](https://github.com/bep) [#4964](https://github.com/gohugoio/hugo/issues/4964)
* Add configurable ref/relref error handling and notFoundURL [e25aa655](https://github.com/gohugoio/hugo/commit/e25aa655f4227ac064be5fe770d517a80acd46b2) [@bep](https://github.com/bep) [#4964](https://github.com/gohugoio/hugo/issues/4964)
* Try node_modules/postcss-cli/bin/postcss first [ebe4d39f](https://github.com/gohugoio/hugo/commit/ebe4d39f175f73e4f130972cb3d74ef0af5d5761) [@bep](https://github.com/bep) [#4952](https://github.com/gohugoio/hugo/issues/4952)
* Add optional lang as argument to rel/relref [d741064b](https://github.com/gohugoio/hugo/commit/d741064bebe2f4663a7ba12556dccc3dffe08629) [@bep](https://github.com/bep) [#4956](https://github.com/gohugoio/hugo/issues/4956)
* Simplify .Site.GetPage etc. [3eb313fe](https://github.com/gohugoio/hugo/commit/3eb313fef495a39731dafa6bddbf77760090230d) [@bep](https://github.com/bep) [#4147](https://github.com/gohugoio/hugo/issues/4147)[#4727](https://github.com/gohugoio/hugo/issues/4727)[#4728](https://github.com/gohugoio/hugo/issues/4728)[#4728](https://github.com/gohugoio/hugo/issues/4728)[#4726](https://github.com/gohugoio/hugo/issues/4726)[#4652](https://github.com/gohugoio/hugo/issues/4652)
* Unify page lookups [b93417aa](https://github.com/gohugoio/hugo/commit/b93417aa1d3d38a9e56bad25937e0e638a113faf) [@vassudanagunta](https://github.com/vassudanagunta) [#4147](https://github.com/gohugoio/hugo/issues/4147)[#4727](https://github.com/gohugoio/hugo/issues/4727)[#4728](https://github.com/gohugoio/hugo/issues/4728)[#4728](https://github.com/gohugoio/hugo/issues/4728)[#4726](https://github.com/gohugoio/hugo/issues/4726)[#4652](https://github.com/gohugoio/hugo/issues/4652)
* Improve error message [4c240800](https://github.com/gohugoio/hugo/commit/4c240800a4275244c9e0847cd6707383180f1ac3) [@bep](https://github.com/bep) 
* Remove unused code [2f2bc7ff](https://github.com/gohugoio/hugo/commit/2f2bc7ff70b90fb11580cc092ef3883bf68d8ad7) [@bep](https://github.com/bep) 

## Fixes

* Avoid server panic on TOML mistake in i18n [75acff5f](https://github.com/gohugoio/hugo/commit/75acff5f20d0d41ffa1ae20402001c7a82f077cb) [@bep](https://github.com/bep) [#4942](https://github.com/gohugoio/hugo/issues/4942)
* Only set 'allThemes' if there are themes in the config file [38204c4a](https://github.com/gohugoio/hugo/commit/38204c4ab6fa2aa2ab8bd06ddb3e07b66e5f9646) [@garrmcnu](https://github.com/garrmcnu) [#4851](https://github.com/gohugoio/hugo/issues/4851)
* Fix potential server panic with drafts/future enabled [1ab4658c](https://github.com/gohugoio/hugo/commit/1ab4658c0d5ea2927f04bd748206e5b139a6326e) [@bep](https://github.com/bep) [#4965](https://github.com/gohugoio/hugo/issues/4965)
* Mark shortcode changes as content changes in server mode [12679b40](https://github.com/gohugoio/hugo/commit/12679b408362a93a3c6159588d6291a3b7ed5548) [@bep](https://github.com/bep) [#4965](https://github.com/gohugoio/hugo/issues/4965)
* Fix source maps on Windows [f01505c9](https://github.com/gohugoio/hugo/commit/f01505c910a325acc18742ac6b3637aa01975e37) [@bep](https://github.com/bep) [#4968](https://github.com/gohugoio/hugo/issues/4968)
* Fix typo-logic bug in GetPage [b56d9a12](https://github.com/gohugoio/hugo/commit/b56d9a1294e692d096bff442e0b1fec61a8c2b0f) [@vassudanagunta](https://github.com/vassudanagunta) 
* Enable test case fixed by commit 501543d4 [d6fde8fa](https://github.com/gohugoio/hugo/commit/d6fde8fa131f3852fa98a8ec5c360e736486cf54) [@vassudanagunta](https://github.com/vassudanagunta) 
* Fix theme config for Work Fs [5c9d5413](https://github.com/gohugoio/hugo/commit/5c9d5413a4e2cc8d44a8b2d7dff04e6523ba2a29) [@bep](https://github.com/bep) [#4951](https://github.com/gohugoio/hugo/issues/4951)
* Fix addkit link to account for i18n [fd1f4a78](https://github.com/gohugoio/hugo/commit/fd1f4a7860c4b989865b47c727239cf924a52fa4) [@hairmare](https://github.com/hairmare) 
