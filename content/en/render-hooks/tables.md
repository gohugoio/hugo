---
title: Table render hooks
linkTitle: Tables
description: Create a table render hook to override the rendering of Markdown tables to HTML.
categories: [render hooks]
keywords: []
menu:
  docs:
    parent: render-hooks
    weight: 90
weight: 90
toc: true
---

{{< new-in 0.134.0 >}}

## Context

Table render hook templates receive the following [context]:

[context]: /getting-started/glossary/#context

###### Attributes

(`map`) The [Markdown attributes], available if you configure your site as follows:

[Markdown attributes]: /content-management/markdown-attributes/

{{< code-toggle file=hugo >}}
[markup.goldmark.parser.attribute]
block = true
{{< /code-toggle >}}

###### Ordinal

(`int`) The zero-based ordinal of the table on the page.

###### Page

(`page`) A reference to the current page.

###### PageInner

(`page`) A reference to a page nested via the [`RenderShortcodes`] method. [See details](#pageinner-details).

[`RenderShortcodes`]: /methods/page/rendershortcodes

###### Position

(`string`) The position of the table within the page content.

###### THead
(`slice`) A slice of table header rows, where each element is a slice of table cells.

###### TBody
(`slice`) A slice of table body rows, where each element is a slice of table cells.

## Table cells

Each table cell within the slice of slices returned by the `THead` and `TBody` methods has the following fields:

###### Alignment
(`string`) The alignment of the text within the table cell, one of `left`, `center`, or `right`.

###### Text
(`template.HTML`) The text within the table cell.

## Example

In its default configuration, Hugo renders Markdown tables according to the [GitHub Flavored Markdown specification]. To create a render hook that does the same thing:

[GitHub Flavored Markdown specification]: https://github.github.com/gfm/#tables-extension-

{{< code file=layouts/_default/_markup/render-table.html copy=true >}}
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
          <th {{ printf "style=%q" (printf "text-align: %s" .Alignment) | safeHTMLAttr }}>
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
          <td {{ printf "style=%q" (printf "text-align: %s" .Alignment) | safeHTMLAttr }}>
            {{- .Text -}}
          </td>
        {{- end }}
      </tr>
    {{- end }}
  </tbody>
</table>
{{< /code >}}

{{% include "/render-hooks/_common/pageinner.md" %}}
