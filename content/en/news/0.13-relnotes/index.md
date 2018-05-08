---
date: 2015-02-22T04:10:06Z
description: "Hugo 0.13 is the largest Hugo release to date. The release introduced some long sought after features (pagination, sequencing, data loading, tons of template improvements) as well as major internal improvements. In addition to the code changes, the Hugo community has grown significantly and now has over 3000 stars on GitHub, 134 contributors, 24 themes and 1000s of happy users."
title: "Hugo 0.13"
categories: ["Releases"]
---

The v0.13.0 release is the largest Hugo release to date. The release introduced
some long sought after features (pagination, sequencing, data loading, tons of
template improvements) as well as major internal improvements. In addition to
the code changes, the Hugo community has grown significantly and now has over
3000 stars on GitHub, 134 contributors, 24 themes and 1000s of happy users.

This release represents **448 contributions by 65 contributors**

A special shout out to [@bep](https://github.com/bep) and
[@anthonyfok](https://github.com/anthonyfok) for their new role as Hugo
maintainers and their tremendous contributions this release.

### New major features
- Support for [data files](http://gohugo.io/extras/datafiles/) in [YAML](http://yaml.org/),
  [JSON](http://www.json.org/), or [TOML](https://github.com/toml-lang/toml)
  located in the `data` directory ([#885](https://github.com/spf13/hugo/issues/885))
- Support for [dynamic content](http://gohugo.io/extras/dynamiccontent/) by loading JSON & CSV
  from remote sources via GetJson and GetCsv in short codes or other layout
  files ([#748](https://github.com/spf13/hugo/issues/748))
- [Pagination support](http://gohugo.io/extras/pagination/) for home page, sections and
  taxonomies ([#750](https://github.com/spf13/hugo/issues/750))
- Universal sequencing support
  - A new, generic Next/Prev functionality is added to all lists of pages
    (sections, taxonomies, etc.)
  - Add in-section [Next/Prev](http://gohugo.io/templates/variables/) content pointers
- `Scratch` -- [a "scratchpad"](http://gohugo.io/extras/scratch) for your node- and page-scoped
  variables
- [Cross Reference](http://gohugo.io/extras/crossreferences/) support to easily link documents
  together with the ref and relref shortcodes.
- [Ace](http://ace.yoss.si/) template engine support ([#541](https://github.com/spf13/hugo/pull/541))
- A new [shortcode](http://gohugo.io/extras/shortcodes/) token of `{{</* */>}}` (raw HTML)
  alongside the existing `{{%/* */%}}` (Markdown)
- A top level `Hugo` variable (on Page & Node) is added with various build
  information
- Several new ways to order and group content:
  - `ByPublishDate`
  - `GroupByPublishDate(format, order)`
  - `GroupByParam(key, order)`
  - `GroupByParamDate(key, format, order)`
- Hugo has undergone a major refactoring, with a new handler system and a
  generic file system. This sounds and is technical, but will pave the way for
  new features and make Hugo even speedier

### Notable enhancements to existing features
- The [shortcode](http://gohugo.io/extras/shortcodes/) handling is rewritten for speed and
  better error messages.
- Several improvements to the [template functions](http://gohugo.io/templates/functions/):
  - `where` is now even more powerful and accepts SQL-like syntax with the
    operators `==`, `eq`; `!=`, `<>`, `ne`; `>=`, `ge`; `>`, `gt`; `<=`,
    `le`; `<`, `lt`; `in`, `not in`
  - `where` template function now also accepts dot chaining key argument
    (e.g. `"Params.foo.bar"`)
- New template functions:
  - `apply`
  - `chomp`
  - `delimit`
  - `sort`
  - `markdownify`
  - `in` and `intersect`
  - `trim`
  - `replace`
  - `dateFormat`
- Several [configurable improvements related to Markdown
  rendering](http://gohugo.io/overview/configuration/#configure-blackfriday-rendering:a66b35d20295cb764719ac8bd35837ec):
  - Configuration of footnote rendering
  - Optional support for smart angled quotes, e.g. `"Hugo"` → «Hugo»
  - Enable descriptive header IDs
- URLs in XML output is now correctly canonified ([#725](https://github.com/spf13/hugo/issues/725), [#728](https://github.com/spf13/hugo/issues/728), and part
  of [#789](https://github.com/spf13/hugo/issues/789))

### Other improvements
- Internal change to use byte buffer pool significantly lowering memory usage
  and providing measurable performance improvements overall
- Changes to docs:
  - A new [Troubleshooting](http://gohugo.io/troubleshooting/overview/) section is added
  - It's now searchable through Google Custom Search ([#753](https://github.com/spf13/hugo/issues/753))
  - Some new great tutorials:
    - [Automated deployments with
      Wercker](http://gohugo.io/tutorials/automated-deployments/)
    - [Creating a new theme](http://gohugo.io/tutorials/creating-a-new-theme/)
- [`hugo new`](http://gohugo.io/content/archetypes/) now copies the content in addition to the front matter
- Improved unit test coverage
- Fixed a lot of Windows-related path issues
- Improved error messages for template and rendering errors
- Enabled soft LiveReload of CSS and images ([#490](https://github.com/spf13/hugo/pull/490))
- Various fixes in RSS feed generation ([#789](https://github.com/spf13/hugo/issues/789))
- `HasMenuCurrent` and `IsMenuCurrent` is now supported on Nodes
- A bunch of [bug fixes](https://github.com/spf13/hugo/commits/master)
