---
title: Configure markup
linkTitle: Markup
description: Configure markup.
categories: []
keywords: []
aliases: [/getting-started/configuration-markup/]
---

## Default handler

In its default configuration, Hugo uses [Goldmark] to render Markdown to HTML.

{{< code-toggle file=hugo >}}
[markup]
defaultMarkdownHandler = 'goldmark'
{{< /code-toggle >}}

Files with ending with `.md`, `.mdown`, or `.markdown` are processed as Markdown, unless you've explicitly set a different format using the `markup` field in your front matter.

To use a different renderer for Markdown files, specify one of `asciidocext`, `org`, `pandoc`, or `rst` in your site configuration.

`defaultMarkdownHandler`|Renderer
:--|:--
`asciidocext`|[AsciiDoc]
`goldmark`|[Goldmark]
`org`|[Emacs Org Mode]
`pandoc`|[Pandoc]
`rst`|[reStructuredText]

To use AsciiDoc, Pandoc, or reStructuredText you must install the relevant renderer and update your [security policy].

> [!note]
> Unless you need a unique capability provided by one of the alternative Markdown handlers, we strongly recommend that you use the default setting. Goldmark is fast, well maintained, conforms to the [CommonMark] specification, and is compatible with [GitHub Flavored Markdown] (GFM).

## Goldmark

This is the default configuration for the Goldmark Markdown renderer:

{{< code-toggle config=markup.goldmark />}}

### Extensions

The extensions below, excluding Extras and Passthrough, are enabled by default.

Extension|Documentation|Enabled
:--|:--|:-:
`cjk`|[Goldmark Extensions: CJK]|:heavy_check_mark:
`definitionList`|[PHP Markdown Extra: Definition lists]|:heavy_check_mark:
`extras`|[Hugo Goldmark Extensions: Extras]||
`footnote`|[PHP Markdown Extra: Footnotes]|:heavy_check_mark:
`linkify`|[GitHub Flavored Markdown: Autolinks]|:heavy_check_mark:
`passthrough`|[Hugo Goldmark Extensions: Passthrough]||
`strikethrough`|[GitHub Flavored Markdown: Strikethrough]|:heavy_check_mark:
`table`|[GitHub Flavored Markdown: Tables]|:heavy_check_mark:
`taskList`|[GitHub Flavored Markdown: Task list items]|:heavy_check_mark:
`typographer`|[Goldmark Extensions: Typographer]|:heavy_check_mark:

#### Extras

{{< new-in 0.126.0 />}}

Enable [deleted text], [inserted text], [mark text], [subscript], and [superscript] elements in Markdown.

Element|Markdown|Rendered
:--|:--|:--
Deleted text|`~~foo~~`|`<del>foo</del>`
Inserted text|`++bar++`|`<ins>bar</ins>`
Mark text|`==baz==`|`<mark>baz</mark>`
Subscript|`H~2~O`|`H<sub>2</sub>O`
Superscript|`1^st^`|`1<sup>st</sup>`

To avoid a conflict when enabling the "subscript" feature of the Extras extension, if you want to render subscript and strikethrough text concurrently you must:

1. Disable the Strikethrough extension
1. Enable the "deleted text" feature of the Extras extension

For example:

{{< code-toggle file=hugo >}}
[markup.goldmark.extensions]
strikethrough = false

[markup.goldmark.extensions.extras.delete]
enable = true

[markup.goldmark.extensions.extras.subscript]
enable = true
{{< /code-toggle >}}

#### Passthrough

{{< new-in 0.122.0 />}}

Enable the Passthrough extension to include mathematical equations and expressions in Markdown using LaTeX markup. See [mathematics in Markdown] for details.

#### Typographer

The Typographer extension replaces certain character combinations with HTML entities as specified below:

Markdown|Replaced by|Description
:--|:--|:--
`...`|`&hellip;`|horizontal ellipsis
`'`|`&rsquo;`|apostrophe
`--`|`&ndash;`|en dash
`---`|`&mdash;`|em dash
`«`|`&laquo;`|left angle quote
`“`|`&ldquo;`|left double quote
`‘`|`&lsquo;`|left single quote
`»`|`&raquo;`|right angle quote
`”`|`&rdquo;`|right double quote
`’`|`&rsquo;`|right single quote

### Settings explained

Most of the Goldmark settings above are self-explanatory, but some require explanation.

duplicateResourceFiles
: {{< new-in 0.123.0 />}}
: (`bool`) Whether to duplicate shared page resources for each language on multilingual single-host sites. See [multilingual page resources] for details. Default is `false`.

  > [!note]
  > With multilingual single-host sites, setting this parameter to `false` will enable Hugo's [embedded link render hook] and [embedded image render hook]. This is the default configuration for multilingual single-host sites.

parser.wrapStandAloneImageWithinParagraph
: (`bool`) Whether to wrap image elements without adjacent content within a `p` element when rendered. This is the default Markdown behavior. Set to `false` when using an [image render hook] to render standalone images as `figure` elements. Default is `true`.

parser.autoDefinitionTermID
: {{< new-in 0.144.0 />}}
: (`bool`) Whether to automatically add `id` attributes to description list terms (i.e., `dt` elements). When `true`, the `id` attribute of each `dt` element is accessible through the [`Fragments.Identifiers`] method on a `Page` object.

parser.autoHeadingID
: (`bool`) Whether to automatically add `id` attributes to headings (i.e., `h1`, `h2`, `h3`, `h4`, `h5`, and `h6` elements).

parser.autoIDType
: (`string`) The strategy used to automatically generate `id` attributes, one of `github`, `github-ascii` or `blackfriday`.

  - `github` produces GitHub-compatible `id` attributes
  - `github-ascii` drops any non-ASCII characters after accent normalization
  - `blackfriday` produces `id` attributes compatible with the Blackfriday Markdown renderer

  This is also the strategy used by the [anchorize](/functions/urls/anchorize) template function. Default is `github`.

parser.attribute.block
: (`bool`) Whether to enable [Markdown attributes] for block elements. Default is `false`.

parser.attribute.title
: (`bool`) Whether to enable [Markdown attributes] for headings. Default is `true`.

renderHooks.image.enableDefault
: {{< new-in 0.123.0 />}}
: (`bool`) Whether to enable the [embedded image render hook]. Default is `false`.

  > [!note]
  > The embedded image render hook is automatically enabled for multilingual single-host sites if [duplication of shared page resources] is disabled. This is the default configuration for multilingual single-host sites.

renderHooks.link.enableDefault
: {{< new-in 0.123.0 />}}
: (`bool`) Whether to enable the [embedded link render hook]. Default is `false`.

  > [!note]
  > The embedded link render hook is automatically enabled for multilingual single-host sites if [duplication of shared page resources] is disabled. This is the default configuration for multilingual single-host sites.

renderer.hardWraps
: (`bool`) Whether to replace newline characters within a paragraph with `br` elements. Default is `false`.

renderer.unsafe
: (`bool`) Whether to render raw HTML mixed within Markdown. This is unsafe unless the content is under your control. Default is `false`.

## AsciiDoc

This is the default configuration for the AsciiDoc renderer:

{{< code-toggle config=markup.asciidocExt />}}

### Settings explained

attributes
: (`map`) A map of key-value pairs, each a document attribute. See Asciidoctor's [attributes].

backend
: (`string`) The backend output file format. Default is `html5`.

extensions
: (`string array`) An array of enabled extensions, one or more of `asciidoctor-html5s`, `asciidoctor-bibtex`, `asciidoctor-diagram`, `asciidoctor-interdoc-reftext`, `asciidoctor-katex`, `asciidoctor-latex`, `asciidoctor-mathematical`, or `asciidoctor-question`.

  > [!note]
  > To mitigate security risks, entries in the extension array may not contain forward slashes (`/`), backslashes (`\`), or periods. Due to this restriction, extensions must be in Ruby's `$LOAD_PATH`.

failureLevel
: (`string`) The minimum logging level that triggers a non-zero exit code (failure). Default is `fatal`.

noHeaderOrFooter
: (`bool`) Whether to output an embeddable document, which excludes the header, the footer, and everything outside the body of the document. Default is `true`.

preserveTOC
: (`bool`) Whether to preserve the table of contents (TOC) rendered by Asciidoctor. By default, to make the TOC compatible with existing themes, Hugo removes the TOC rendered by Asciidoctor. To render the TOC, use the [`TableOfContents`] method on a `Page` object in your templates. Default is `false`.

safeMode
: (`string`) The safe mode level, one of `unsafe`, `safe`, `server`, or `secure`. Default is `unsafe`.

sectionNumbers
: (`bool`) Whether to number each section title. Default is `false`.

trace
: (`bool`) Whether to include backtrace information on errors. Default is `false`.

verbose
: (`bool`) Whether to verbosely print processing information and configuration file checks to stderr. Default is `false`.

workingFolderCurrent
: (`bool`) Whether to set the working directory to be the same as that of the AsciiDoc file being processed, allowing [includes] to work with relative paths. Set to `true` to render diagrams with the [asciidoctor-diagram] extension. Default is `false`.

### Configuration example

{{< code-toggle file=hugo >}}
[markup.asciidocExt]
    extensions = ["asciidoctor-html5s", "asciidoctor-diagram"]
    workingFolderCurrent = true
    [markup.asciidocExt.attributes]
        my-base-url = "https://example.com/"
        my-attribute-name = "my value"
{{< /code-toggle >}}

### Syntax highlighting

Follow the steps below to enable syntax highlighting.

#### Step 1

Set the `source-highlighter` attribute in your site configuration. For example:

{{< code-toggle file=hugo >}}
[markup.asciidocExt.attributes]
source-highlighter = 'rouge'
{{< /code-toggle >}}

#### Step 2

Generate the highlighter CSS. For example:

```text
rougify style monokai.sublime > assets/css/syntax.css
```

#### Step 3

In your base template add a link to the CSS file:

```go-html-template {file="layouts/baseof.html"}
<head>
  ...
  {{ with resources.Get "css/syntax.css" }}
    <link rel="stylesheet" href="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous">
  {{ end }}
  ...
</head>
```

Then add the code to be highlighted to your markup:

```text
[#hello,ruby]
----
require 'sinatra'

get '/hi' do
  "Hello World!"
end
----
```

### Troubleshooting

Run `hugo --logLevel debug` to examine Hugo's call to the Asciidoctor executable:

```txt
INFO 2019/12/22 09:08:48 Rendering book-as-pdf.adoc with C:\Ruby26-x64\bin\asciidoctor.bat using asciidoc args [--no-header-footer -r asciidoctor-html5s -b html5s -r asciidoctor-diagram --base-dir D:\prototypes\hugo_asciidoc_ddd\docs -a outdir=D:\prototypes\hugo_asciidoc_ddd\build -] ...
```

## Highlight

This is the default configuration.

{{< code-toggle config=markup.highlight />}}

{{% include "/_common/syntax-highlighting-options.md" %}}

## Table of contents

This is the default configuration for the table of contents, applicable to Goldmark and Asciidoctor:

{{< code-toggle config=markup.tableOfContents />}}

startLevel
: (`int`) Heading levels less than this value will be excluded from the table of contents. For example, to exclude `h1` elements from the table of contents, set this value to `2`. Default is `2`.

endLevel
: (`int`) Heading levels greater than this value will be excluded from the table of contents. For example, to exclude `h4`, `h5`, and `h6` elements from the table of contents, set this value to `3`. Default is `3`.

ordered
: (`bool`) Whether to generates an ordered list instead of an unordered list. Default is `false`.

[`Fragments.Identifiers`]: /methods/page/fragments/#identifiers
[`TableOfContents`]: /methods/page/tableofcontents/
[asciidoctor-diagram]: https://asciidoctor.org/docs/asciidoctor-diagram/
[attributes]: https://asciidoctor.org/docs/asciidoc-syntax-quick-reference/#attributes-and-substitutions
[CommonMark]: https://spec.commonmark.org/current/
[deleted text]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/del
[duplication of shared page resources]: /configuration/markup/#duplicateresourcefiles
[duplication of shared page resources]: /configuration/markup/#duplicateresourcefiles
[embedded image render hook]: /render-hooks/images/#default
[embedded image render hook]: /render-hooks/images/#default
[embedded link render hook]: /render-hooks/links/#default
[embedded link render hook]: /render-hooks/links/#default
[GitHub Flavored Markdown]: https://github.github.com/gfm/
[GitHub Flavored Markdown: Autolinks]: https://github.github.com/gfm/#autolinks-extension-
[GitHub Flavored Markdown: Strikethrough]: https://github.github.com/gfm/#strikethrough-extension-
[GitHub Flavored Markdown: Tables]: https://github.github.com/gfm/#tables-extension-
[GitHub Flavored Markdown: Task list items]: https://github.github.com/gfm/#task-list-items-extension-
[Goldmark]: https://github.com/yuin/goldmark/
[Goldmark Extensions: CJK]: https://github.com/yuin/goldmark?tab=readme-ov-file#cjk-extension
[Goldmark Extensions: Typographer]: https://github.com/yuin/goldmark?tab=readme-ov-file#typographer-extension
[Hugo Goldmark Extensions: Extras]: https://github.com/gohugoio/hugo-goldmark-extensions?tab=readme-ov-file#extras-extension
[Hugo Goldmark Extensions: Passthrough]: https://github.com/gohugoio/hugo-goldmark-extensions?tab=readme-ov-file#passthrough-extension
[image render hook]: /render-hooks/images/
[includes]: https://docs.asciidoctor.org/asciidoc/latest/syntax-quick-reference/#includes
[inserted text]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/ins
[mark text]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/mark
[Markdown attributes]: /content-management/markdown-attributes/
[mathematics in Markdown]: content-management/mathematics/
[multilingual page resources]: /content-management/page-resources/#multilingual
[PHP Markdown Extra: Definition lists]: https://michelf.ca/projects/php-markdown/extra/#def-list
[PHP Markdown Extra: Footnotes]: https://michelf.ca/projects/php-markdown/extra/#footnotes
[security policy]: /configuration/security/
[subscript]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/sub
[superscript]: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/sup
[AsciiDoc]: https://asciidoc.org/
[Emacs Org Mode]: https://orgmode.org/
[Pandoc]: https://pandoc.org/
[reStructuredText]: https://docutils.sourceforge.io/rst.html
