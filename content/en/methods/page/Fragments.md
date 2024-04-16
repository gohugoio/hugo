---
title: Fragments
description: Returns a data structure of the fragments in the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/TableOfContents
  returnType: tableofcontents.Fragments
  signatures: [PAGE.Fragments]
toc: true
---

{{< new-in 0.111.0 >}}

In a URL, whether absolute or relative, the [fragment] links to an `id` attribute of an HTML element on the page.

```text
/articles/article-1#section-2
------------------- ---------
       path         fragment
```

Hugo assigns an `id` attribute to each Markdown [ATX] and [setext] heading within the page content. You can override the `id` with a [Markdown attribute] as needed. This creates the relationship between an entry in the [table of contents] (TOC) and a heading on the page.

Use the `Fragments` method on a `Page` object to create a table of contents with the `Fragments.ToHTML` method, or by [walking] the `Fragments.Map` data structure.

## Methods

Headings
: (`map`) A nested map of all headings on the page. Each map contains the following keys: `ID`, `Level`, `Title` and `Headings`. To inspect the data structure:

```go-html-template
<pre>{{ debug.Dump .Fragments.Headings }}</pre>
```

HeadingsMap
: (`slice`) A slice of maps of all headings on the page, with first-level keys for each heading. Each map contains the following keys: `ID`, `Level`, `Title` and `Headings`. To inspect the data structure:

```go-html-template
<pre>{{ debug.Dump .Fragments.HeadingsMap }}</pre>
```

Identifiers
: (`slice`) A slice containing the `id` of each heading on the page. To inspect the data structure:

```go-html-template
<pre>{{ debug.Dump .Fragments.Identifiers }}</pre>
```

Identifiers.Contains ID
: (`bool`) Reports whether one or more headings on the page has the given `id` attribute, useful for validating fragments within a link [render hook].

```go-html-template
{{ .Fragments.Identifiers.Contains "section-2" }} → true
```

Identifiers.Count ID
: (`int`) The number of headings on a page with the given `id` attribute, useful for detecting duplicates.

```go-html-template
{{ .Fragments.Identifiers.Count "section-2" }} → 1
```

ToHTML
: (`template.HTML`) Returns a TOC as a nested list, either ordered or unordered, identical to the HTML returned by the [`TableOfContents`] method. This method take three arguments: the start level&nbsp;(`int`), the end level&nbsp;(`int`), and a boolean (`true` to return an ordered list, `false` to return an unordered list).

Use this method when you want to control the start level, end level, or list type independently from the table of contents settings in your site configuration.

```go-html-template
{{ $startLevel := 2 }}
{{ $endLevel := 3 }}
{{ $ordered := true }}
{{ .Fragments.ToHTML $startLevel $endLevel $ordered }}
```

Hugo renders this to:

```html
<nav id="TableOfContents">
  <ol>
    <li><a href="#section-1">Section 1</a>
      <ol>
        <li><a href="#section-11">Section 1.1</a></li>
        <li><a href="#section-12">Section 1.2</a></li>
      </ol>
    </li>
    <li><a href="#section-2">Section 2</a></li>
  </ol>
</nav>
```

{{% note %}}
It is safe to use the `Fragments` methods within a render hook, even for the current page.

When using the `Fragments` methods within a shortcode, call the shortcode using the `{{</* */>}}` notation. If you use the `{{%/* */%}}` notation, the rendered shortcode is included in the creation of the fragments map, resulting in a circular loop.
{{% /note %}}

[atx]: https://spec.commonmark.org/0.30/#atx-headings
[fragment]: /getting-started/glossary/#fragment
[markdown attribute]: /getting-started/glossary/#markdown-attribute
[setext]: https://spec.commonmark.org/0.30/#setext-headings
[table of contents]: /methods/page/tableofcontents/
[walking]: /getting-started/glossary/#walk
[`tableofcontents`]: /methods/page/tableofcontents/
[render hook]: /getting-started/glossary/#render-hook
