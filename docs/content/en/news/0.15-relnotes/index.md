---
date: 2015-12-19T09:45:24Z
description: "The v0.15.0 Hugo release brings a lot of polish to Hugo. Exactly 6 months after the 0.14 release, Hugo has seen massive growth and changes. Most notably, this is Hugo's first release under the Apache 2.0 license."
title: "Hugo 0.15"
categories: ["Releases"]
---

The v0.15.0 Hugo release brings a lot of polish to Hugo. Exactly 6 months after the 0.14 release, Hugo has seen massive growth and changes. Most notably, this is Hugo's first release under the Apache 2.0 license. With this license change we hope to expand the great community around Hugo and make it easier for our many users to contribute.  This release represents over **377 contributions by 87 contributors** to the main Hugo repo and hundreds of improvements to the libraries Hugo uses. Hugo also launched a [new theme showcase](http://themes.gohugo.io) and participated in [Hacktoberfest](https://hacktoberfest.digitalocean.com).

Hugo now has:
- 6700 (+2700) stars on GitHub
- 235 (+75) contributors
- 65 (+30) themes

**Template Improvements:** This release takes Hugo to a new level of speed and usability. Considerable work has been done adding features and performance to the template system which now has full support of Ace, Amber and Go Templates.

**Hugo Import:** Have a Jekyll site, but dreaming of porting it to Hugo? This release introduces a new `hugo import jekyll`command that makes this easier than ever.

**Performance Improvements:** Just when you thought Hugo couldn't get any faster, Hugo continues to improve in speed while adding features. Notably Hugo 0.15 introduces the ability to render and serve directly from memory resulting in 30%+ lower render times.

Huge thanks to all who participated in this release. A special thanks to [@bep](https://github.com/bep) who led the development of Hugo this release again, [@anthonyfok](https://github.com/anthonyfok), [@eparis](https://github.com/eparis), [@tatsushid](https://github.com/tatsushid) and [@DigitalCraftsman](https://github.com/digitalcraftsman/).

## New features
- new `hugo import jekyll` command. [#1469](https://github.com/spf13/hugo/pull/1469)
- The new `Param` convenience method on `Page` and `Node` can be used to get the most specific parameter value for a given key. [#1462](https://github.com/spf13/hugo/issues/1462)
- Several new information elements have been added to `Page` and `Node`:
  - `RuneCount`: The number of [runes](http://blog.golang.org/strings) in the content, excluding any whitespace. This may be a good alternative to `.WordCount`  for Japanese and other CJK languages where a word-split by spaces makes no sense.  [#1266](https://github.com/spf13/hugo/issues/1266)
  - `RawContent`: Raw Markdown as a string. One use case may be of embedding remarkjs.com slides.
  - `IsHome`: tells the truth about whether you're on the home page or not.

## Improvements
- `hugo server` now builds ~30%+ faster by rendering to memory instead of disk. To get the old behavior, start the server with `--renderToDisk=true`.
- Hugo now supports dynamic reloading of the config file when watching.
- We now use a custom-built `LazyFileReader` for reading file contents, which means we don't read media files in `/content` into memory anymore -- and file reading is now performed in parallel on multicore PCs. [#1181](https://github.com/spf13/hugo/issues/1181)
- Hugo is now built with `Go 1.5` which, among many other improvements, have fixed the last known data race in Hugo. [#917](https://github.com/spf13/hugo/issues/917)
- Paginator now also supports page groups. [#1274](https://github.com/spf13/hugo/issues/1274)
- Markdown improvements:
  - Hugo now supports GitHub-flavoured markdown code fences for highlighting for `md`-files (Blackfriday rendered markdown) and `mmark` files (MMark rendered markdown). [#362](https://github.com/spf13/hugo/issues/362) [#1258](https://github.com/spf13/hugo/issues/1258)
  - Several new Blackfriday options are added:
    - Option to disable Blackfriday's `Smartypants`.
    - Option for Blackfriday to open links in a new window/tab. [#1220](https://github.com/spf13/hugo/issues/1220)
    - Option to disable Blackfriday's LaTeX style dashes [#1231](https://github.com/spf13/hugo/issues/1231)
    - Definition lists extension support.
- `Scratch` now has built-in `map` support.
- We now fall back to `link title` for the default page sort. [#1299](https://github.com/spf13/hugo/issues/1299)
- Some notable new configuration options:
  -  `IgnoreFiles` can be set with a list of Regular Expressions that matches files to be ignored during build. [#1189](https://github.com/spf13/hugo/issues/1189)
  - `PreserveTaxonomyNames`, when set to `true`, will preserve what you type as the taxonomy name both in the folders created and the taxonomy `key`, but it will be normalized for the URL.  [#1180](https://github.com/spf13/hugo/issues/1180)
- `hugo gen` can now generate man files, bash auto complete and markdown documentation
- Hugo will now make suggestions when a command is mistyped
- Shortcodes now have a boolean `.IsNamedParams` property. [#1597](https://github.com/spf13/hugo/pull/1597)

## New Template Features
- All template engines:
  - The new `dict` function that could be used to pass maps into a template. [#1463](https://github.com/spf13/hugo/pull/1463)
  - The new `pluralize` and `singularize` template funcs.
  - The new `base64Decode` and `base64Encode` template funcs.
  - The `sort` template func now accepts field/key chaining arguments and pointer values. [#1330](https://github.com/spf13/hugo/issues/1330)
  - Several fixes for `slicestr` and `substr`, most importantly, they now have full `utf-8`-support. [#1190](https://github.com/spf13/hugo/issues/1190) [#1333](https://github.com/spf13/hugo/issues/1333) [#1347](https://github.com/spf13/hugo/issues/1347)
  - The new `last` template function allows the user to select the last `N` items of a slice. [#1148](https://github.com/spf13/hugo/issues/1148)
  - The new `after` func allows the user to select the items after the `Nth` item. [#1200](https://github.com/spf13/hugo/pull/1200)
  - Add `time.Time` type support to the `where`, `ge`, `gt`, `le`, and `lt` template functions.
  - It is now possible to use constructs like `where Values ".Param.key" nil` to filter pages that doesn't have a particular parameter. [#1232](https://github.com/spf13/hugo/issues/1232)
  - `getJSON`/`getCSV`: Add retry on invalid content. [#1166](https://github.com/spf13/hugo/issues/1166)
  -   The new `readDir` func lists local files. [#1204](https://github.com/spf13/hugo/pull/1204)
  - The new `safeJS` function allows the embedding of content into JavaScript contexts in Go templates.
  - Get the main site RSS link from any page by accessing the `.Site.RSSLink` property. [#1566](https://github.com/spf13/hugo/pull/1566)
- Ace templates:
  - Base templates now also works in themes. [#1215](https://github.com/spf13/hugo/issues/1215).
  - And now also on Windows. [#1178](https://github.com/spf13/hugo/issues/1178)
- Full support for Amber templates including all template functions.
- A built-in template for Google Analytics. [#1505](https://github.com/spf13/hugo/pull/1505)
- Hugo is now shipped with new built-in shortcodes: [#1576](https://github.com/spf13/hugo/issues/1576)
  - `youtube` for YouTube videos
  - `vimeo` for Vimeo videos
  - `gist` for GitHub gists
  - `tweet` for Twitter Tweets
  - `speakerdeck` for Speakerdeck slides

## Bugfixes
- Fix data races in page sorting and page reversal. These operations are now also cached. [#1293](https://github.com/spf13/hugo/issues/1293)
- `page.HasMenuCurrent()` and `node.HasMenuCurrent()` now work correctly in multi-level nested menus.
- Support `Fish and Chips` style section titles. Previously, this would end up as  `Fish And Chips`. Now, the first character is made toupper, but the rest are preserved as-is. [#1176](https://github.com/spf13/hugo/issues/1176)
- Hugo now removes superfluous p-tags around shortcodes. [#1148](https://github.com/spf13/hugo/issues/1148)

## Notices
- `hugo server` will watch by default now.
- Some fields and methods were deprecated in `0.14`. These are now removed, so the error message isn't as friendly if you still use the old values. So please change:
  -   `getJson` to `getJSON`, `getCsv` to `getCSV`, `safeHtml` to
    `safeHTML`, `safeCss` to `safeCSS`, `safeUrl` to `safeURL`, `Url` to `URL`,
    `UrlPath` to `URLPath`, `BaseUrl` to `BaseURL`, `Recent` to `Pages`.

## Known Issues

Using the Hugo v0.15 32-bit Windows or ARM binary, running `hugo server` would crash or hang due to a [memory alignment issue](https://golang.org/pkg/sync/atomic/#pkg-note-BUG) in [Afero](https://github.com/spf13/afero).  The bug was discovered shortly after the v0.15.0 release and has since been [fixed](https://github.com/spf13/afero/pull/23) by @tpng.  If you encounter this bug, you may either compile Hugo v0.16-DEV from source, or use the following solution/workaround:
- **64-bit Windows users: Please use [hugo_0.15_windows_amd64.zip](https://github.com/spf13/hugo/releases/download/v0.15/hugo_0.15_windows_amd64.zip)** (amd64 == x86-64).  It is only the 32-bit hugo_0.15_windows_386.zip that crashes/hangs (see #1621 and #1628).
- **32-bit Windows and ARM users: Please run `hugo server --renderToDisk` as a workaround** until Hugo v0.16 is released (see [“hugo server” returns runtime error on armhf](https://discuss.gohugo.io/t/hugo-server-returns-runtime-error-on-armhf/2293) and #1716).
