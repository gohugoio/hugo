---
title: TableOfContents
description: Returns a table of contents for the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/Fragments
  returnType: template.HTML
  signatures: [PAGE.TableOfContents]
aliases: [/content-management/toc/]
---

The `TableOfContents` method on a `Page` object returns an ordered or unordered list of the Markdown [ATX] and [setext] headings within the page content.

[atx]: https://spec.commonmark.org/0.30/#atx-headings
[setext]: https://spec.commonmark.org/0.30/#setext-headings

This template code:

```go-html-template
{{ .TableOfContents }}
```

Produces this HTML:

```html
<nav id="TableOfContents">
  <ul>
    <li><a href="#section-1">Section 1</a>
      <ul>
        <li><a href="#section-11">Section 1.1</a></li>
        <li><a href="#section-12">Section 1.2</a></li>
      </ul>
    </li>
    <li><a href="#section-2">Section 2</a></li>
  </ul>
</nav>
```

By default, the `TableOfContents` method returns an unordered list of level 2 and level 3 headings. You can adjust this in your site configuration:

{{< code-toggle file=hugo >}}
[markup.tableOfContents]
endLevel = 3
ordered = false
startLevel = 2
{{< /code-toggle >}}
