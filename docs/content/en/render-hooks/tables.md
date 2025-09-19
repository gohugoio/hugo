---
title: Table render hooks
linkTitle: Tables
description: Create table render hook templates to override the rendering of Markdown tables to HTML.
categories: []
keywords: []
---

{{< new-in 0.134.0 />}}

## Context

Table _render hook_ templates receive the following [context](g):

Attributes
: (`map`) The [Markdown attributes], available if you configure your site as follows:

  {{< code-toggle file=hugo >}}
  [markup.goldmark.parser.attribute]
  block = true
  {{< /code-toggle >}}

Ordinal
: (`int`) The zero-based ordinal of the table on the page.

Page
: (`page`) A reference to the current page.

PageInner
: (`page`) A reference to a page nested via the [`RenderShortcodes`] method. [See details](#pageinner-details).

Position
: (`string`) The position of the table within the page content.

THead
: (`slice`) A slice of table header rows, where each element is a slice of table cells.

TBody
: (`slice`) A slice of table body rows, where each element is a slice of table cells.

[Markdown attributes]: /content-management/markdown-attributes/
[`RenderShortcodes`]: /methods/page/rendershortcodes

## Table cells

Each table cell within the slice of slices returned by the `THead` and `TBody` methods has the following fields:

Alignment
: (`string`) The alignment of the text within the table cell, one of `left`, `center`, or `right`.

Text
: (`template.HTML`) The text within the table cell.

## Example

In its default configuration, Hugo renders Markdown tables according to the [GitHub Flavored Markdown specification]. To create a render hook that does the same thing:

[GitHub Flavored Markdown specification]: https://github.github.com/gfm/#tables-extension-

```go-html-template {file="layouts/_markup/render-table.html" copy=true}
<table
  {{- range $k, $v := .Attributes }}
    {{- if $v }}
      {{- printf " %s=%q" $k $v | safeHTMLAttr }}
    {{- end }}
  {{- end }}>
  <thead>
    {{- range .THead }}
      <tr>
        {{- range . }}
          <th
            {{- with .Alignment }}
              {{- printf " style=%q" (printf "text-align: %s" .) | safeHTMLAttr }}
            {{- end -}}
          >
            {{- .Text -}}
          </th>
        {{- end }}
      </tr>
    {{- end }}
  </thead>
  <tbody>
    {{- range .TBody }}
      <tr>
        {{- range . }}
          <td
            {{- with .Alignment }}
              {{- printf " style=%q" (printf "text-align: %s" .) | safeHTMLAttr }}
            {{- end -}}
          >
            {{- .Text -}}
          </td>
        {{- end }}
      </tr>
    {{- end }}
  </tbody>
</table>
```

{{% include "/_common/render-hooks/pageinner.md" %}}
