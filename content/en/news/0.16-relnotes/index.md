---
date: 2016-06-06T13:54:06-04:00
description: "Hugo 0.16 is our best and biggest release ever. The Hugo community has outdone itself with continued performance improvements, beautiful themes and increased stability."
title: "Hugo 0.16"
categories: ["Releases"]
aliases: [/0-16/]
---

Hugo 0.16 is our best and biggest release ever. The Hugo community has
outdone itself with continued performance improvements,
[beautiful themes](http://themes.gohugo.io) for all types of sites from project
sites to documentation to blogs to portfolios, and increased stability.

This release represents **over 550 contributions by over 110 contributors** to
the main Hugo code base. Since last release Hugo has **gained 3500 stars, 90
contributors and 23 additional themes.**

This release celebrates 3 years since  {{< gh "@spf13" >}} wrote the first lines
of Hugo. During those 3 years Hugo has accomplished some major milestones
including...

* 10,000+ stars on GitHub
* 320+ contributors
* 90+ themes
* 1000s of happy websites
* Many subprojects like {{< gh "@spf13/cobra">}}, {{< gh "@spf13/viper">}} and
  {{< gh "@spf13/afero">}} which have experienced broad usage across the Go
  ecosystem.

{{< gh "@bep" >}} led the development of Hugo for the 3rd consecutive release
with nearly half of the contributions to 0.16 in addition to his considerable
contributions as lead maintainer. {{< gh "@anthonyfok" >}}, {{< gh
"@DigitalCraftsman" >}}, {{< gh "@MooreReason" >}} all made significant
contributions. A special thanks to {{< gh "@abourget" >}} for his considerable
work on multilingual support. Due to its broad impact we wanted to spend more
time testing it and it will be included in Hugo's next release.

### Highlights

**Partial Builds:** Prior to this release Hugo would always reread and rebuild
the entire site. This release introduces support for reactive site building
while watching (`hugo server`). Hugo will watch the filesystem for changes and
only re-read the changed files. Depending on the files change Hugo will
intelligently re-render only the needed portion of the site. Performance gains
depend on the operation performed and size of the site. In our testing build
times decreased anywhere from 10% to 99%.

**Template Improvements:** Template improvements continue to be a mainstay of each Hugo release. Hugo 0.16 adds support for the new `block` keyword introduced in Go 1.6 -- think base templates with default sections -- as well as many new template functions.

**Polish:** As Hugo matures releases will inevitably contain fewer huge new features. This release represents hundreds of small improvements across ever facet of Hugo which will make for a much better experience for all of our users. Worth mentioning here is the curious bug where live reloading didn't work in some editors on OS X, including the popular TextMate 2. This is now fixed. Oh, and now any error will exit with an error code, a big thing for automated deployments.

### New Features
* Support reading configuration variables from the OS environment {{<gh 2090 >}}
* Add emoji support {{<gh 1892>}}
* Add `themesDir` option to configuration {{<gh 1556>}}
* Add support for Go 1.6 `block` keyword in templates {{<gh 1832>}}
* Partial static sync {{<gh 1644>}}
* Source file based relative linking (à la GitHub) {{<gh 0x0f6b334b6715253b030c4e783b88e911b6e53e56>}}
* Add `ByLastmod` sort function to pages. {{<gh 0xeb627ca16de6fb5e8646279edd295a8bf0f72bf1 >}}
* New templates functions:
    * `readFile` {{<gh 1551 >}}
    * `countwords` and `countrunes` {{<gh 1440>}}
    * `default` {{<gh 1943>}}
    * `hasPrefix` {{<gh 1243>}}
    * `humanize` {{<gh 1818>}}
    * `jsonify` {{<gh 0x435e996c4fd48e9009ffa9f83a19fb55f0777dbd>}}
    * `md5` and `sha1` {{<gh 1932>}}
    * `replaceRE` {{<gh 1845>}}
    * `findRE` {{<gh 2048>}}
    * `shuffle` {{<gh 1942>}}
    * `slice` {{<gh 1902>}}
    * `plainify` {{<gh 1915>}}

### Enhancements

* Hugo now exits with error code on any error. This is a big thing for
  automated deployments. {{<gh 740 >}}
* Print error when `/index.html` is zero-length {{<gh 947>}}
* Enable dirname and filename bash autocompletion for more flags {{<gh
  0x666ddd237791b56fd048992dca9a27d1af50a10e>}}
* Improve error handling in commands {{<gh 1502>}}
* Add sanity checks for `hugo import jekyll` {{<gh 1625 >}}
* Add description to `Page.Params` {{<gh 1484>}}
* Add async version of Google Analytics internal template {{<gh 1711>}}
* Add autostart option to YouTube shortcode {{<gh 1784>}}
* Set Date and Lastmod for main home page {{<gh 1903>}}
* Allow URL with extension in front matter {{<gh 1923>}}
* Add list support in Scratch {{<gh
  0xeaba04e82bdfc5d4c29e970f11b4aab9cc0efeaa>}}
* Add file option to gist shortcode {{<gh 1955>}}
* Add config layout and content directory CLI options {{<gh 1698>}}
* Add boolean value comparison to `where` template function {{<gh
  0xf3c74c9db484c8961e70cb3458f9e41e7832fa12>}}
* Do not write to cache when `ignoreCache` is set  {{<gh 2067>}}
* Add option to disable rendering of 404 page  {{<gh 2037>}}
* Mercurial is no longer needed to build Hugo {{<gh 2062 >}}
* Do not create `robots.txt` by default {{<gh 2049>}}
* Disable syntax guessing for PygmentsCodeFences by default.  To enable syntax
  guessing again, add the following to your config file:
  `PygmentsCodeFencesGuessSyntax = true` {{<gh 2034>}}
* Make `ByCount` sort consistently {{<gh 1930>}}
* Add `Scratch` to shortcode {{<gh 2000>}}
* Add support for symbolic links for content, layout, static, theme  {{<gh 1855 >}}
* Add '+' as one of the valid characters in URLs specified in the front matter {{<gh 1290 >}}
* Make alias redirect output URLs relative when `RelativeURLs = true` {{<gh 2093 >}}
* Hugo injects meta generator tag on homepage if missing {{<gh 2182 >}}

### Fixes

* Fix file change watcher for TextMate 2 and friends on OS X {{<gh 1053 >}}
* Make dynamic reloading of config file reliable on all platform {{<gh 1684 >}}
* Hugo now works on Linux/arm64 {{<gh 1772 >}}
* `plainIDAnchors` now defaults to `true`  {{<gh 2057>}}
* Win32 and ARM builds fixed {{<gh 1716>}}
* Copy static dir files without theme's static dir {{<gh 1656>}}
* Make `noTimes` command flag work {{<gh 1657>}}
* Change most global CLI flags into local ones {{<gh 1624>}}
* Remove transformation of menu URLs {{<gh 1239>}}
* Do not fail on unknown Jekyll file {{<gh 1705>}}
* Use absolute path when editing with editor {{<gh 1589>}}
* Fix hugo server "Watching for changes" path display {{<gh 1721>}}
* Do not strip special characters out of URLs {{<gh 1292>}}
* Fix `RSSLink` when uglyURLs are enabled {{<gh 175>}}
* Get BaseURL from viper in server mode {{<gh 1821>}}
* Fix shortcode handling in RST {{<gh 1904>}}
* Use default sitemap configuration for homepage {{<gh 1304>}}
* Exit if specific port is unavailable in server mode {{<gh 1901>}}
* Fix regression in "section menus for lazy blogger" {{<gh 2065>}}

****

## **0.15.0**  November 25, 2015

The v0.15.0 Hugo release brings a lot of polish to Hugo. Exactly 6 months after
the 0.14 release, Hugo has seen massive growth and changes. Most notably, this
is Hugo's first release under the Apache 2.0 license. With this license change
we hope to expand the great community around Hugo and make it easier for our
many users to contribute.  This release represents over **377 contributions by
87 contributors** to the main Hugo repo and hundreds of improvements to the
libraries Hugo uses. Hugo also launched a [new theme
showcase](http://themes.gohugo.io) and participated in
[Hacktoberfest](https://hacktoberfest.digitalocean.com).

Hugo now has:

* 6700 (+2700) stars on GitHub
* 235 (+75) contributors
* 65 (+30) themes


**Template Improvements:** This release takes Hugo to a new level of speed and
usability. Considerable work has been done adding features and performance to
the template system which now has full support of Ace, Amber and Go Templates.

**Hugo Import:** Have a Jekyll site, but dreaming of porting it to Hugo? This
release introduces a new `hugo import jekyll`command that makes this easier
than ever.

**Performance Improvements:** Just when you thought Hugo couldn't get any faster,
Hugo continues to improve in speed while adding features. Notably Hugo 0.15
introduces the ability to render and serve directly from memory resulting in
30%+ lower render times.

Huge thanks to all who participated in this release. A special thanks to
{{< gh "@bep" >}} who led the development of Hugo this release again,
{{< gh "@anthonyfok" >}},
{{< gh "@eparis" >}},
{{< gh "@tatsushid" >}} and
{{< gh "@DigitalCraftsman" >}}.


### New features
* new `hugo import jekyll` command. {{< gh 1469 >}}
* The new `Param` convenience method on `Page` and `Node` can be used to get the most specific parameter value for a given key. {{< gh 1462 >}}
* Several new information elements have been added to `Page` and `Node`:
    * `RuneCount`: The number of [runes](http://blog.golang.org/strings) in the content, excluding any whitespace. This may be a good alternative to `.WordCount`  for Japanese and other CJK languages where a word-split by spaces makes no sense.  {{< gh 1266 >}}
    * `RawContent`: Raw Markdown as a string. One use case may be of embedding remarkjs.com slides.
    * `IsHome`: tells the truth about whether you're on the home page or not.

### Improvements
* `hugo server` now builds ~30%+ faster by rendering to memory instead of disk. To get the old behavior, start the server with `--renderToDisk=true`.
* Hugo now supports dynamic reloading of the config file when watching.
* We now use a custom-built `LazyFileReader` for reading file contents, which means we don't read media files in `/content` into memory anymore -- and file reading is now performed in parallel on multicore PCs. {{< gh 1181 >}}
* Hugo is now built with `Go 1.5` which, among many other improvements, have fixed the last known data race in Hugo. {{< gh 917 >}}
* Paginator now also supports page groups. {{< gh 1274 >}}
* Markdown improvements:
    * Hugo now supports GitHub-flavoured markdown code fences for highlighting for `md`-files (Blackfriday rendered markdown) and `mmark` files (MMark rendered markdown). {{< gh 362 1258 >}}
    * Several new Blackfriday options are added:
        * Option to disable Blackfriday's `Smartypants`.
        * Option for Blackfriday to open links in a new window/tab. {{< gh 1220 >}}
        * Option to disable Blackfriday's LaTeX style dashes {{< gh 1231 >}}
        * Definition lists extension support.
* `Scratch` now has built-in `map` support.
* We now fall back to `link title` for the default page sort. {{< gh 1299 >}}
* Some notable new configuration options:
    *  `IgnoreFiles` can be set with a list of Regular Expressions that matches files to be ignored during build. {{< gh 1189 >}}
    * `PreserveTaxonomyNames`, when set to `true`, will preserve what you type as the taxonomy name both in the folders created and the taxonomy `key`, but it will be normalized for the URL.  {{< gh 1180 >}}
* `hugo gen` can now generate man files, bash auto complete and markdown documentation
* Hugo will now make suggestions when a command is mistyped
* Shortcodes now have a boolean `.IsNamedParams` property. {{< gh 1597 >}}

### New Template Features
* All template engines:
    * The new `dict` function that could be used to pass maps into a template. {{< gh 1463 >}}
    * The new `pluralize` and `singularize` template funcs.
    * The new `base64Decode` and `base64Encode` template funcs.
    * The `sort` template func now accepts field/key chaining arguments and pointer values. {{< gh 1330 >}}
    * Several fixes for `slicestr` and `substr`, most importantly, they now have full `utf-8`-support. {{< gh 1190 1333 1347 >}}
    * The new `last` template function allows the user to select the last `N` items of a slice. {{< gh 1148 >}}
    * The new `after` func allows the user to select the items after the `Nth` item. {{< gh 1200 >}}
    * Add `time.Time` type support to the `where`, `ge`, `gt`, `le`, and `lt` template functions.
    * It is now possible to use constructs like `where Values ".Param.key" nil` to filter pages that doesn't have a particular parameter. {{< gh 1232 >}}
    * `getJSON`/`getCSV`: Add retry on invalid content. {{< gh 1166 >}}
    *   The new `readDir` func lists local files. {{< gh 1204 >}}
    * The new `safeJS` function allows the embedding of content into JavaScript contexts in Go templates.
    * Get the main site RSS link from any page by accessing the `.Site.RSSLink` property. {{< gh 1566 >}}
* Ace templates:
    * Base templates now also works in themes. {{< gh 1215 >}}.
    * And now also on Windows. {{< gh 1178 >}}
* Full support for Amber templates including all template functions.
* A built-in template for Google Analytics. {{< gh 1505 >}}
* Hugo is now shipped with new built-in shortcodes: {{< gh 1576 >}}
  * `youtube` for YouTube videos
  * `vimeo` for Vimeo videos
  * `gist` for GitHub gists
  * `tweet` for Twitter Tweets
  * `speakerdeck` for Speakerdeck slides


### Bugfixes
* Fix data races in page sorting and page reversal. These operations are now also cached. {{< gh 1293 >}}
* `page.HasMenuCurrent()` and `node.HasMenuCurrent()` now work correctly in multi-level nested menus.
* Support `Fish and Chips` style section titles. Previously, this would end up as  `Fish And Chips`. Now, the first character is made toupper, but the rest are preserved as-is. {{< gh 1176 >}}
* Hugo now removes superfluous p-tags around shortcodes. {{< gh 1148 >}}

### Notices
* `hugo server` will watch by default now.
* Some fields and methods were deprecated in `0.14`. These are now removed, so the error message isn't as friendly if you still use the old values. So please change:
    *   `getJson` to `getJSON`, `getCsv` to `getCSV`, `safeHtml` to
  `safeHTML`, `safeCss` to `safeCSS`, `safeUrl` to `safeURL`, `Url` to `URL`,
  `UrlPath` to `URLPath`, `BaseUrl` to `BaseURL`, `Recent` to `Pages`.

### Known Issues

Using the Hugo v0.15 32-bit Windows or ARM binary, running `hugo server` would crash or hang due to a [memory alignment issue](https://golang.org/pkg/sync/atomic/#pkg-note-BUG) in [Afero](https://github.com/spf13/afero).  The bug was discovered shortly after the v0.15.0 release and has since been [fixed](https://github.com/spf13/afero/pull/23) by {{< gh "@tpng" >}}.  If you encounter this bug, you may either compile Hugo v0.16-DEV from source, or use the following solution/workaround:

* **64-bit Windows users: Please use [hugo_0.15_windows_amd64.zip](https://github.com/gohugoio/hugo/releases/download/v0.15/hugo_0.15_windows_amd64.zip)** (amd64 == x86-64).  It is only the 32-bit hugo_0.15_windows_386.zip that crashes/hangs (see {{< gh 1621 >}} and {{< gh 1628 >}}).
* **32-bit Windows and ARM users: Please run `hugo server --renderToDisk` as a workaround** until Hugo v0.16 is released (see [“hugo server” returns runtime error on armhf](https://discourse.gohugo.io/t/hugo-server-returns-runtime-error-on-armhf/2293) and {{< gh 1716 >}}).
