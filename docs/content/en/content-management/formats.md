---
title: Content formats
description: Create your content using Markdown, HTML, Emacs Org Mode, AsciiDoc, Pandoc, or reStructuredText.
categories: []
keywords: []
aliases: [/content/markdown-extras/,/content/supported-formats/,/doc/supported-formats/]
---

## Introduction

You may mix content formats throughout your site. For example:

```tree
content/
└── posts/
    ├── post-1.md
    ├── post-2.adoc
    ├── post-3.org
    ├── post-4.pandoc
    ├── post-5.rst
    └── post-6.html
```

Regardless of content format, all content must have [front matter][], preferably including both `title` and `date`.

Hugo selects the content renderer based on the `markup` identifier in front matter, falling back to the file extension. See the [classification](#classification) table below for a list of markup identifiers and recognized file extensions.

## Formats

### Markdown

Create your content in [Markdown][] preceded by front matter.

Markdown is Hugo's default content format. Hugo natively renders Markdown to HTML using [Goldmark][]. Goldmark is fast and conforms to the [CommonMark][] and [GitHub Flavored Markdown][] specifications. You can configure Goldmark in your [project configuration][configure goldmark].

Hugo provides custom Markdown features including:

[Attributes][]
: Apply HTML attributes such as `class` and `id` to Markdown images and block elements including blockquotes, fenced code blocks, headings, horizontal rules, lists, paragraphs, and tables.

[Extensions][]
: Leverage the embedded Markdown extensions to create tables, definition lists, footnotes, task lists, inserted text, mark text, subscripts, superscripts, and more.

[Mathematics][]
: Include mathematical equations and expressions in Markdown using LaTeX markup.

[Render hooks][]
: Override the conversion of Markdown to HTML when rendering fenced code blocks, headings, images, and links. For example, render every standalone image as an HTML `figure` element.

### HTML

Create your content in [HTML][] preceded by front matter. The content is typically what you would place within an HTML document's `body` or `main` element.

> [!NOTE]
> The HTML content format is denied by default. See [`security.allowContent`][].

### Emacs Org Mode

Create your content in the [Emacs Org Mode][] format preceded by front matter. You can use Org Mode keywords for front matter. See [details][].

### AsciiDoc

Create your content in the [AsciiDoc][] format preceded by front matter. Hugo renders AsciiDoc content to HTML using the Asciidoctor executable. You must install Asciidoctor and its dependencies (Ruby) to render the AsciiDoc content format.

You can configure the AsciiDoc renderer in your [project configuration][configure asciidoc].

In its default configuration, Hugo passes these CLI flags when calling the Asciidoctor executable:

```sh
--no-header-footer
```

The CLI flags passed to the Asciidoctor executable depend on configuration. You may inspect the flags when building your project:

```sh
hugo build --logLevel info
```

### Pandoc

Create your content in the [Pandoc][] format preceded by front matter. Hugo renders Pandoc content to HTML using the Pandoc executable. You must install Pandoc to render the Pandoc content format.

Hugo passes these CLI flags when calling the Pandoc executable:

```sh
--mathjax
```

### reStructuredText

Create your content in the [reStructuredText][] format preceded by front matter. Hugo renders reStructuredText content to HTML using [Docutils][], specifically rst2html. You must install Docutils and its dependencies (Python) to render the reStructuredText content format.

Hugo passes these CLI flags when calling the rst2html executable:

```sh
--leave-comments --initial-header-level=2
```

## Classification

{{% include "/_common/content-format-table.md" %}}

When converting content to HTML, Hugo uses:

- Native renderers for Markdown, HTML, and Emacs Org mode
- External renderers for AsciiDoc, Pandoc, and reStructuredText

Native renderers are faster than external renderers.

[AsciiDoc]: https://asciidoc.org/
[Attributes]: /content-management/markdown-attributes/
[CommonMark]: https://spec.commonmark.org/current/
[Docutils]: https://docutils.sourceforge.io/
[Emacs Org Mode]: https://orgmode.org/
[Extensions]: /configuration/markup/#extensions
[GitHub Flavored Markdown]: https://github.github.com/gfm/
[Goldmark]: https://github.com/yuin/goldmark
[HTML]: https://developer.mozilla.org/en-US/docs/Learn_web_development/Getting_started/Your_first_website/Creating_the_content
[Markdown]: https://daringfireball.net/projects/markdown/
[Mathematics]: /content-management/mathematics/
[Pandoc]: https://pandoc.org/MANUAL.html#pandocs-markdown
[Render hooks]: /render-hooks/introduction/
[`security.allowContent`]: /configuration/security/#allowcontent
[configure asciidoc]: /configuration/markup/#asciidoc
[configure goldmark]: /configuration/markup/#goldmark
[details]: /content-management/front-matter/#emacs-org-mode
[front matter]: /content-management/front-matter/
[reStructuredText]: https://docutils.sourceforge.io/rst.html
