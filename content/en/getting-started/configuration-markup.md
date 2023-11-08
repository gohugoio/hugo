---
title: Configure markup
description: Configure rendering of markup to HTML.
categories: [getting started,fundamentals]
keywords: [configuration,highlighting]
menu:
  docs:
    parent: getting-started
    weight: 50
weight: 50
slug: configuration-markup
toc: true
---

## Default handler

By default, Hugo uses [Goldmark] to render markdown to HTML.

{{< code-toggle file=hugo >}}
[markup]
defaultMarkdownHandler = 'goldmark'
{{< /code-toggle >}}

Files with the `.md` or `.markdown` extension are processed as markdown, provided that you have not specified a different [content format] using the `markup` field in front matter.

To use a different renderer for markdown files, specify one of `asciidocext`, `org`, `pandoc`, or `rst` in your site configuration.

defaultMarkdownHandler|Description
:--|:--
`asciidocext`|[AsciiDoc]
`goldmark`|[Goldmark]
`org`|[Emacs Org Mode]
`pandoc`|[Pandoc]
`rst`|[reStructuredText]

To use Asciidoc, Pandoc, or reStructuredText you must install the relevant renderer and update your [security policy].

{{% note %}}
Unless you need a unique capability provided by one of the alternate markdown handlers, we strongly recommend that you use the default setting. Goldmark is fast, well maintained, conforms to the [CommonMark] specification, and is compatible with [GitHub Flavored Markdown] (GFM).

[commonmark]: https://spec.commonmark.org/0.30/
[github flavored markdown]: https://github.github.com/gfm/
{{% /note %}}

[asciidoc]: https://asciidoc.org/
[content format]: /content-management/formats/#list-of-content-formats
[emacs org mode]: https://orgmode.org/
[goldmark]: https://github.com/yuin/goldmark/
[pandoc]: https://pandoc.org/
[restructuredtext]: https://docutils.sourceforge.io/rst.html
[security policy]: /about/security-model/#security-policy

## Goldmark

This is the default configuration for the Goldmark markdown renderer:

{{< code-toggle config=markup.goldmark />}}

For details on the extensions, refer to [this section](https://github.com/yuin/goldmark/#built-in-extensions) of the Goldmark documentation

Some settings explained:

hardWraps
: By default, Goldmark ignores newlines within a paragraph. Set to `true` to render newlines as `<br>` elements.

unsafe
: By default, Goldmark does not render raw HTML and potentially dangerous links. If you have lots of inline HTML and/or JavaScript, you may need to turn this on.

typographer
: The typographer extension replaces certain character combinations with HTML entities as specified below:

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

attribute
: Enable custom attribute support for titles and blocks by adding attribute lists inside single curly brackets (`{.myclass class="class1 class2" }`) and placing it _after the Markdown element it decorates_, on the same line for titles and on a new line directly below for blocks.

Hugo supports adding attributes (e.g. CSS classes) to Markdown blocks, e.g. tables, lists, paragraphs etc.

A blockquote with a CSS class:

```md
> foo
> bar
{.myclass}
```

There are some current limitations: For tables you can currently only apply it to the full table, and for lists the `ul`/`ol`-nodes only, e.g.:

```md
* Fruit
  * Apple
  * Orange
  * Banana
  {.fruits}
* Dairy
  * Milk
  * Cheese
  {.dairies}
{.list}
```

Note that attributes in [code fences](/content-management/syntax-highlighting/#highlighting-in-code-fences) must come after the opening tag, with any other highlighting processing instruction, e.g.:

````txt
```go {.myclass linenos=table,hl_lines=[8,"15-17"],linenostart=199}
// ... code
```
````

autoHeadingIDType ("github")
: The strategy used for creating auto IDs (anchor names). Available types are `github`, `github-ascii` and `blackfriday`. `github` produces GitHub-compatible IDs, `github-ascii` will drop any non-Ascii characters after accent normalization, and `blackfriday` will make the IDs compatible with Blackfriday, the default Markdown engine before Hugo 0.60. Note that if Goldmark is your default Markdown engine, this is also the strategy used in the [anchorize](/functions/urls/anchorize) template func.

## Asciidoc

This is the default configuration for the AsciiDoc markdown renderer:

{{< code-toggle config=markup.asciidocExt />}}

attributes
: (`map`) Variables to be referenced in your AsciiDoc file. This is a list of variable name/value maps. See Asciidoctor’s [attributes].

[attributes]: https://asciidoctor.org/docs/asciidoc-syntax-quick-reference/#attributes-and-substitutions

backend:
: (`string`) Don’t change this unless you know what you are doing.

extensions
: (`[]string`) Possible extensions are `asciidoctor-html5s`, `asciidoctor-bibtex`, `asciidoctor-diagram`, `asciidoctor-interdoc-reftext`, `asciidoctor-katex`, `asciidoctor-latex`, `asciidoctor-mathematical`, and `asciidoctor-question`.

failureLevel
: (`string`) The minimum logging level that triggers a non-zero exit code (failure).

noHeaderOrFooter
: (`bool`) Output an embeddable document, which excludes the header, the footer, and everything outside the body of the document. Don’t change this unless you know what you are doing.

preserveTOC
: (`bool`) By default, Hugo removes the table of contents generated by Asciidoctor and provides it through the built-in variable `.TableOfContents` to enable further customization and better integration with the various Hugo themes. This option can be set to true to preserve Asciidoctor’s TOC in the generated page.

safeMode
: (`string`) Safe mode level `unsafe`, `safe`, `server`, or `secure`. Don’t change this unless you know what you are doing.

sectionNumbers
: (`bool`) Auto-number section titles.

trace
: (`bool`) Include backtrace information on errors.

verbose
: (`bool`) Verbosely print processing information and configuration file checks to stderr.

workingFolderCurrent
: (`bool`) Sets the working directory to be the same as that of the AsciiDoc file being processed, so that [include] will work with relative paths. This setting uses the asciidoctor cli parameter --base-dir and attribute outdir=. For rendering diagrams with [asciidoctor-diagram], `workingFolderCurrent` must be set to `true`.

[asciidoctor-diagram]: https://asciidoctor.org/docs/asciidoctor-diagram/
[include]: https://asciidoctor.org/docs/asciidoc-syntax-quick-reference/#include-files

Notice that for security concerns only extensions that do not have path separators (either `\`, `/` or `.`) are allowed. That means that extensions can only be invoked if they are in the Ruby's `$LOAD_PATH` (ie. most likely, the extension has been installed by the user). Any extension declared relative to the website's path will not be accepted.

Example of how to set extensions and attributes:

```yml
[markup.asciidocExt]
    extensions = ["asciidoctor-html5s", "asciidoctor-diagram"]
    workingFolderCurrent = true
    [markup.asciidocExt.attributes]
        my-base-url = "https://example.com/"
        my-attribute-name = "my value"
```

In a complex Asciidoctor environment it is sometimes helpful to debug the exact call to your external helper with all
parameters. Run Hugo with `-v`. You will get an output like

```txt
INFO 2019/12/22 09:08:48 Rendering book-as-pdf.adoc with C:\Ruby26-x64\bin\asciidoctor.bat using asciidoc args [--no-header-footer -r asciidoctor-html5s -b html5s -r asciidoctor-diagram --base-dir D:\prototypes\hugo_asciidoc_ddd\docs -a outdir=D:\prototypes\hugo_asciidoc_ddd\build -] ...
```

## Highlight

This is the default `highlight` configuration. Note that some of these settings can be set per code block, see [Syntax Highlighting](/content-management/syntax-highlighting/).

{{< code-toggle config=markup.highlight />}}

For `style`, see these galleries:

* [Short snippets](https://xyproto.github.io/splash/docs/all.html)
* [Long snippets](https://xyproto.github.io/splash/docs/longer/all.html)

For CSS, see [Generate Syntax Highlighter CSS](/content-management/syntax-highlighting/#generate-syntax-highlighter-css).

## Table of contents

{{< code-toggle config=markup.tableOfContents />}}

These settings only works for the Goldmark renderer:

startLevel
: The heading level, values starting at 1 (`h1`), to start render the table of contents.

endLevel
: The heading level, inclusive, to stop render the table of contents.

ordered
: If `true`, generates an ordered list instead of an unordered list.

## Render hooks

See [Markdown Render Hooks](/templates/render-hooks/).
