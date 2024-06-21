---
title: Content formats
description: Create your content using Markdown, HTML, Emacs Org Mode, AsciiDoc, Pandoc, or reStructuredText.
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

## Introduction

You may mix content formats throughout your site. For example:

```text
content/
└── posts/
    ├── post-1.md
    ├── post-2.adoc
    ├── post-3.org
    ├── post-4.pandoc
    ├── post-5.rst
    └── post-6.html
```

Regardless of content format, all content must have [front matter], preferably including both `title` and `date`.

Hugo selects the content renderer based on the `markup` identifier in front matter, falling back to the file extension. See the [classification](#classification) table below for a list of markup identifiers and recognized file extensions.

## Formats

### Markdown

Create your content in [Markdown] preceded by front matter.

Markdown is Hugo's default content format. Hugo natively renders Markdown to HTML using [Goldmark]. Goldmark is fast and conforms to the [CommonMark] and [GitHub Flavored Markdown] specifications. You can [configure Goldmark] in your site configuration.

Hugo provides custom Markdown features including:

[Attributes]
: Apply HTML attributes such as `class` and `id` to Markdown images and block elements including blockquotes, fenced code blocks, headings, horizontal rules, lists, paragraphs, and tables.

[Extensions]
: Leverage the embedded Markdown extensions to create tables, definition lists, footnotes, task lists, inserted text, mark text, subscripts, superscripts, and more.

[Mathematics]
: Include mathematical equations and expressions in Markdown using LaTeX or TeX typesetting syntax.

[Render hooks]
: Override the conversion of Markdown to HTML when rendering fenced code blocks, headings, images, and links. For example, render every standalone image as an HTML `figure` element.

### HTML

Create your content in [HTML] preceded by front matter. The content is typically what you would place within an HTML document's `body` or `main` element.

### Emacs Org Mode

Create your content in the [Emacs Org Mode] format preceded by front matter. You can use Org Mode keywords for front matter. See [details](/content-management/front-matter/#emacs-org-mode)).

### AsciiDoc

Create your content in the [AsciiDoc] format preceded by front matter. Hugo renders AsciiDoc content to HTML using the Asciidoctor executable. You must install Asciidoctor and its dependencies (Ruby) to use the AsciiDoc content format.

You can [configure the AsciiDoc renderer] in your site configuration.

In its default configuration, Hugo passes these CLI flags when calling the Asciidoctor executable:

```text
--no-header-footer
```

The CLI flags passed to the Asciidoctor executable depend on configuration. You may inspect the flags when building your site:

```text
hugo --logLevel info
```

### Pandoc

Create your content in the [Pandoc] format preceded by front matter. Hugo renders Pandoc content to HTML using the Pandoc executable. You must install Pandoc to use the Pandoc content format.

Hugo passes these CLI flags when calling the Pandoc executable:

```text
--mathjax
```

### reStructuredText

Create your content in the [reStructuredText] format preceded by front matter. Hugo renders reStructuredText content to HTML using [Docutils], specifically rst2html. You must install Docutils and its dependencies (Python) to use the reStructuredText content format.

Hugo passes these CLI flags when calling the rst2html executable:

```text
--leave-comments --initial-header-level=2
```

## Classification

Content format|Media type|Identifier|File extensions
:--|:--|:--|:--
Markdown|`text/markdown`|`markdown`|`markdown`,`md`, `mdown`
HTML|`text/html`|`html`|`htm`, `html`
Emacs Org Mode|`text/org`|`org`|`org`
AsciiDoc|`text/asciidoc`|`asciidoc`|`ad`, `adoc`, `asciidoc`
Pandoc|`text/pandoc`|`pandoc`|`pandoc`, `pdc`
reStructuredText|`text/rst`|`rst`|`rst`

When converting content to HTML, Hugo uses:

- Native renderers for Markdown, HTML, and Emacs Org mode
- External renderers for AsciiDoc, Pandoc, and reStructuredText

Native renderers are faster than external renderers.

[AsciiDoc]: https://asciidoc.org/
[Asciidoctor]: https://asciidoctor.org/
[Attributes]: /content-management/markdown-attributes/
[CommonMark]: https://spec.commonmark.org/current/
[Docutils]: https://docutils.sourceforge.io/
[Emacs Org Mode]: https://orgmode.org/
[Extensions]: /getting-started/configuration-markup/#goldmark-extensions
[GitHub Flavored Markdown]: https://github.github.com/gfm/
[Goldmark]: https://github.com/yuin/goldmark
[HTML]: https://developer.mozilla.org/en-US/docs/Learn/Getting_started_with_the_web/HTML_basics
[Markdown]: https://daringfireball.net/projects/markdown/
[Mathematics]: /content-management/mathematics/
[Pandoc]: https://pandoc.org/
[Render hooks]: https://gohugo.io/render-hooks/introduction/
[configure Goldmark]: /getting-started/configuration-markup/#goldmark
[configure the AsciiDoc renderer]: /getting-started/configuration-markup/#asciidoc
[front matter]: /content-management/front-matter/
[reStructuredText]: https://docutils.sourceforge.io/rst.html
