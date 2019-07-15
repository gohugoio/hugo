---
date: 2015-05-26T01:32:45Z
description: "The v0.14.0 Hugo release brings of the most demanded features to Hugo. The foundation of Hugo is stabilizing nicely and a lot of polish has been added. We’ve expanded support for additional content types with support for AsciiDoc, Restructured Text, HTML and Markdown."
title: "Hugo 0.14"
categories: ["Releases"]
---

The v0.14.0 Hugo release brings of the most demanded features to Hugo. The foundation of Hugo is stabilizing nicely and a lot of polish has been added. We’ve expanded support for additional content types with support for AsciiDoc, Restructured Text, HTML and Markdown. Some of these types depend on external libraries as there does not currently exist native support in Go. We’ve tried to make the experience as seamless as possible. Look for more improvements here in upcoming releases.

A lot of work has been done to improve the user experience, with extra polish to the Windows experience. Hugo errors are more helpful overall and Hugo now can detect if it’s being run in Windows Explorer and provide additional instructions to run it via the command prompt.

The Hugo community continues to grow. Hugo has over 4000 stars on github, 165 contributors, 35 themes and 1000s of happy users. It is now the 5th most popular static site generator (by Stars) and has the 3rd largest contributor community.

This release represents over **240 contributions by 36 contributors** to the main Hugo codebase.

Big shout out to [@bep](https://github.com/bep) who led the development of Hugo this release, [@anthonyfok](https://github.com/anthonyfok), [@eparis](https://github.com/eparis), [@SchumacherFM](https://github.com/SchumacherFM), [@RickCogley](https://github.com/RickCogley) & [@mdhender](https://github.com/mdhender) for their significant contributions and [@tatsushid](https://github.com/tatsushid) for his continuous improvements to the templates. Also a big thanks to all the theme creators. 11 new themes have been added since last release and the [hugoThemes repo now has previews of all of them](https://github.com/spf13/hugoThemes/blob/master/README.md#theme-list).

Hugo also depends on a lot of other great projects. A big thanks to all of our dependencies inclding:
[cobra](https://github.com/spf13/cobra), [viper](https://github.com/spf13/viper), [blackfriday](https://github.com/russross/blackfriday), [pflag](https://github.com/spf13/pflag), [HugoThemes](https://github.com/spf13/hugothemes), [BurntSushi/Toml](github.com/BurntSushi/toml), [goYaml](https://github.com/go-yaml/yaml/tree/v2), and the Go standard library.

## New features
- Support for all file types in content directory.
  - If dedicated file type handler isn’t found it will be copied to the destination.
- Add `AsciiDoc` support using external helpers.
- Add experimental support for [`Mmark`](https://github.com/miekg/mmark) markdown processor
- Bash autocomplete support via `genautocomplete` command
- Add section menu support for a [Section Menu for "the Lazy Blogger"](http://gohugo.io/extras/menus.md#section-menu-for-the-lazy-blogger")
- Add support for `Ace` base templates
- Adding `RelativeURLs = true` to site config will now make all the relative URLs relative to the content root.
- New template functions:
  - `getenv`
  - The string functions `substr` and `slicestr`
    *`seq`, a sequence generator very similar to its Gnu counterpart
  - `absURL` and `relURL`, both of which takes the `BaseURL` setting into account

## Improvements
- Highlighting with `Pygments` is now cached to disk -- expect a major speed boost if you use it!
- More Pygments highlighting options, including `line numbers`
- Show help information to Windows users who try to double click on `hugo.exe`.
- Add `bind` flag to `hugo server` to set the interface to which the server will bind
- Add support for `canonifyurls` in `srcset`
- Add shortcode support for HTML (content) files
- Allow the same `shortcode` to  be used with or without inline content
- Configurable RSS output filename

## Bugfixes
- Fix panic with paginator and zero pages in result set.
- Fix crossrefs on Windows.
- Fix `eq` and `ne` template functions when used with a raw number combined with the result of `add`, `sub` etc.
- Fix paginator with uglyurls
- Fix [#998](https://github.com/spf13/hugo/issues/988), supporting UTF8 characters in Permalinks.

## Notices
- To get variable and function names in line with the rest of the Go community, a set of variable and function names has been deprecated: These will still  work in 0.14, but will be removed in 0.15. What to do should be obvious by  the build log; `getJson` to `getJSON`, `getCsv` to `getCSV`, `safeHtml` to   `safeHTML`, `safeCss` to `safeCSS`, `safeUrl` to `safeURL`, `Url` to `URL`,  `UrlPath` to `URLPath`, `BaseUrl` to `BaseURL`, `Recent` to `Pages`,  `Indexes` to `Taxonomies`.
