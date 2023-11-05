---
title: Content formats
description: Both HTML and Markdown are supported content formats.
categories: [content management]
keywords: [markdown,asciidoc,pandoc,content format]
menu:
  docs:
    parent: content-management
    weight: 40
weight: 40
toc: true
aliases: [/content/markdown-extras/,/content/supported-formats/,/doc/supported-formats/]
---

You can put any file type into your `/content` directories, but Hugo uses the `markup` front matter value if set or the file extension (see `Markup identifiers` in the table below) to determine if the markup needs to be processed, e.g.:

* Markdown converted to HTML
* [Shortcodes](/content-management/shortcodes/) processed
* Layout applied

## List of content formats

The current list of content formats in Hugo:

| Name  | Markup identifiers | Comment |
| ------------- | ------------- |-------------|
| Goldmark  | `markdown`, `goldmark`  |Note that you can set the default handler of `md` and `markdown` to something else, see [Configure Markup](/getting-started/configuration-markup/).|
|Emacs Org-Mode|`org`|See [go-org](https://github.com/niklasfasching/go-org).|
|AsciiDoc|`asciidocext`, `adoc`, `ad`|Needs [Asciidoctor][ascii] installed.|
|RST|`rst`|Needs [RST](https://docutils.sourceforge.io/rst.html) installed.|
|Pandoc|`pandoc`, `pdc`|Needs [Pandoc](https://www.pandoc.org/) installed.|
|HTML|`html`, `htm`|To be treated as a content file, with layout, shortcodes etc., it must have front matter. If not, it will be copied as-is.|

The `markup identifier` is fetched from either the `markup` variable in front matter or from the file extension. For markup-related configuration, see [Configure Markup](/getting-started/configuration-markup/).

## External helpers

Some of the formats in the table above need external helpers installed on your PC. For example, for AsciiDoc files,
Hugo will try to call the `asciidoctor` command. This means that you will have to install the associated
tool on your machine to be able to use these formats.

Hugo passes reasonable default arguments to these external helpers by default:

- `asciidoctor`: `--no-header-footer -`
- `rst2html`: `--leave-comments --initial-header-level=2`
- `pandoc`: `--mathjax`

{{% note %}}
Because additional formats are external commands, generation performance will rely heavily on the performance of the external tool you are using. As this feature is still in its infancy, feedback is welcome.
{{% /note %}}

### Asciidoctor

The Asciidoctor community offers a wide set of tools for the AsciiDoc format that can be installed additionally to Hugo.
[See the Asciidoctor docs for installation instructions](https://asciidoctor.org/docs/install-toolchain/). Make sure that also all
optional extensions like `asciidoctor-diagram` or `asciidoctor-html5s` are installed if required.

{{% note %}}
External `asciidoctor` command requires Hugo rendering to _disk_ to a specific destination directory. It is required to run Hugo with the command option `--destination`.
{{% /note %}}

Some Asciidoctor parameters can be customized in Hugo. See&nbsp;[details].

[details]: /getting-started/configuration-markup/#asciidoc

## Learn markdown

Markdown syntax is simple enough to learn in a single sitting. The following are excellent resources to get you up and running:

* [Daring Fireball: Markdown, John Gruber (Creator of Markdown)][fireball]
* [Markdown Cheatsheet, Adam Pritchard][mdcheatsheet]
* [Markdown Tutorial (Interactive), Garen Torikian][mdtutorial]
* [The Markdown Guide, Matt Cone][mdguide]

[ascii]: https://asciidoctor.org/
[config]: /getting-started/configuration/
[developer tools]: /tools/
[fireball]: https://daringfireball.net/projects/markdown/
[gfmtasks]: https://guides.github.com/features/mastering-markdown/#syntax
[helperssource]: https://github.com/gohugoio/hugo/blob/77c60a3440806067109347d04eb5368b65ea0fe8/helpers/general.go#L65
[hl]: /content-management/syntax-highlighting/
[hlsc]: /content-management/shortcodes/#highlight
[hugocss]: /css/style.css
[ietf]: https://tools.ietf.org/html/
[mathjaxdocs]: https://docs.mathjax.org/en/latest/
[mdcheatsheet]: https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet
[mdguide]: https://www.markdownguide.org/
[mdtutorial]: https://www.markdowntutorial.com/
[org]: https://orgmode.org/
[pandoc]: https://www.pandoc.org/
[rest]: https://docutils.sourceforge.io/rst.html
[sc]: /content-management/shortcodes/
[sct]: /templates/shortcode-templates/
