---
title: collections.Where 
description: Returns the given collection, removing elements that do not satisfy the comparison condition.
categories: []
keywords: []
action:
  aliases: [where]
  related: []
  returnType: any
  signatures: ['collections.Where COLLECTION KEY [OPERATOR] VALUE']
toc: true
aliases: [/functions/where]
---

The `where` function returns the given collection, removing elements that do not satisfy the comparison condition. The comparison condition is composed of the `KEY`, `OPERATOR`, and `VALUE` arguments:

```text
collections.Where COLLECTION KEY [OPERATOR] VALUE
                             --------------------
                             comparison condition
```

Hugo will test for equality if you do not provide an `OPERATOR` argument. For example:

```go-html-template
{{ $pages := where .Site.RegularPages "Section" "books" }}
{{ $books := where .Site.Data.books "genres" "suspense" }}
```

## Arguments

The where function takes three or four arguments. The `OPERATOR` argument is optional.

COLLECTION
: (`any`) A [page collection] or a [slice] of [maps].

[maps]: /getting-started/glossary/#map
[page collection]: /getting-started/glossary/#page-collection
[slice]: /getting-started/glossary/#slice

KEY
: (`string`) The key of the page or map value to compare with `VALUE`. With page collections, commonly used comparison keys are `Section`, `Type`, and `Params`. To compare with a member of the page `Params` map, [chain] the subkey as shown below:

```go-html-template
{{ $result := where .Site.RegularPages "Params.foo" "bar" }}
```

[chain]: /getting-started/glossary/#chain

OPERATOR
: (`string`) The logical comparison [operator](#operators).

VALUE
: (`any`) The value with which to compare. The values to compare must have comparable data types. For example:

Comparison|Result
:--|:--
`"123" "eq" "123"`|`true`
`"123" "eq" 123`|`false`
`false "eq" "false"`|`false`
`false "eq" false`|`true`

When one or both of the values to compare is a slice, use the `in`, `not in`, or `intersect` operators as described below.

## Operators

Use any of the following logical operators:

`=`, `==`, `eq`
: (`bool`) Reports whether the given field value is equal to `VALUE`.

`!=`, `<>`, `ne`
: (`bool`) Reports whether the given field value is not equal to `VALUE`.

`>=`, `ge`
: (`bool`) Reports whether the given field value is greater than or equal to `VALUE`.

`>`, `gt`
: `true` Reports whether the given field value is greater than `VALUE`.

`<=`, `le`
: (`bool`) Reports whether the given field value is less than or equal to `VALUE`.

`<`, `lt`
: (`bool`) Reports whether the given field value is less than `VALUE`.

`in`
: (`bool`) Reports whether the given field value is a member of `VALUE`. Compare string to slice, or string to string. See&nbsp;[details](/functions/collections/in).

`not in`
: (`bool`) Reports whether the given field value is not a member of `VALUE`. Compare string to slice, or string to string. See&nbsp;[details](/functions/collections/in).

`intersect`
: (`bool`) Reports whether the given field value (a slice) contains one or more elements in common with `VALUE`. See&nbsp;[details](/functions/collections/intersect).

`like` {{< new-in 0.116.0 >}}
: (`bool`) Reports whether the given field value matches the regular expression specified in `VALUE`. Use the `like` operator to compare `string` values. The `like` operator returns `false` when comparing other data types to the regular expression.

{{% note %}}
The examples below perform comparisons within a page collection, but the same comparisons are applicable to a slice of maps.
{{% /note %}}

## String comparison

Compare the value of the given field to a [`string`]:

[`string`]: /getting-started/glossary/#string

```go-html-template
{{ $pages := where .Site.RegularPages "Section" "eq" "books" }}
{{ $pages := where .Site.RegularPages "Section" "ne" "books" }}
```

## Numeric comparison

Compare the value of the given field to an [`int`] or [`float`]:

[`int`]: /getting-started/glossary/#int
[`float`]: /getting-started/glossary/#float

```go-html-template
{{ $books := where site.RegularPages "Section" "eq" "books" }}

{{ $pages := where $books "Params.price" "eq" 42 }}
{{ $pages := where $books "Params.price" "ne" 42.67 }}
{{ $pages := where $books "Params.price" "ge" 42 }}
{{ $pages := where $books "Params.price" "gt" 42.67 }}
{{ $pages := where $books "Params.price" "le" 42 }}
{{ $pages := where $books "Params.price" "lt" 42.67 }}
```

## Boolean comparison

Compare the value of the given field to a [`bool`]:

[`bool`]: /getting-started/glossary/#bool

```go-html-template
{{ $books := where site.RegularPages "Section" "eq" "books" }}

{{ $pages := where $books "Params.fiction" "eq" true }}
{{ $pages := where $books "Params.fiction" "eq" false }}
{{ $pages := where $books "Params.fiction" "ne" true }}
{{ $pages := where $books "Params.fiction" "ne" false }}
```

## Member comparison

Compare a [`scalar`] to a [`slice`].

[`scalar`]: /getting-started/glossary/#scalar
[`slice`]: /getting-started/glossary/#slice

For example, to return a collection of pages where the `color` page parameter is either "red" or "yellow":

```go-html-template
{{ $fruit := where site.RegularPages "Section" "eq" "fruit" }}

{{ $colors := slice "red" "yellow" }}
{{ $pages := where $fruit "Params.color" "in" $colors }}
```

To return a collection of pages where the "color" page parameter is neither "red" nor "yellow":

```go-html-template
{{ $fruit := where site.RegularPages "Section" "eq" "fruit" }}

{{ $colors := slice "red" "yellow" }}
{{ $pages := where $fruit "Params.color" "not in" $colors }}
```

## Intersection comparison

Compare a [`slice`] to a [`slice`], returning collection elements with common values. This is frequently used when comparing taxonomy terms.

For example, to return a collection of pages where any of the terms in the "genres" taxonomy are "suspense" or "romance":

```go-html-template
{{ $books := where site.RegularPages "Section" "eq" "books" }}

{{ $genres := slice "suspense" "romance" }}
{{ $pages := where $books "Params.genres" "intersect" $genres }}
```

## Regular expression comparison

{{< new-in 0.116.0 >}}

To return a collection of pages where the "author" page parameter begins with either "victor" or "Victor":

```go-html-template
{{ $pages := where .Site.RegularPages "Params.author" "like" `(?i)^victor` }}
```

{{% include "functions/_common/regular-expressions.md" %}}

{{% note %}}
Use the `like` operator to compare string values. Comparing other data types will result in an empty collection.
{{% /note %}}

## Date comparison

### Predefined dates

There are four predefined front matter dates: [`date`], [`publishDate`], [`lastmod`], and [`expiryDate`]. Regardless of the front matter data format (TOML, YAML, or JSON) these are [`time.Time`] values, allowing precise comparisons.

[`date`]: /methods/page/date/
[`publishdate`]: /methods/page/publishdate/
[`lastmod`]: /methods/page/lastmod/
[`expirydate`]: /methods/page/expirydate/
[`time.Time`]: https://pkg.go.dev/time#Time

For example, to return a collection of pages that were created before the current year:

```go-html-template
{{ $startOfYear := time.AsTime (printf "%d-01-01" now.Year) }}
{{ $pages := where .Site.RegularPages "Date" "lt" $startOfYear }}
```

### Custom dates

With custom front matter dates, the comparison depends on the front matter data format (TOML, YAML, or JSON). 

{{% note %}}
Using TOML for pages with custom front matter dates enables precise date comparisons.
{{% /note %}}

With TOML, date values are first-class citizens. TOML has a date data type while JSON and YAML do not. If you quote a TOML date, it is a string. If you do not quote a TOML date value, it is [`time.Time`] value, enabling precise comparisons.

In the TOML example below, note that the event date is not quoted.

{{< code file="content/events/2024-user-conference.md" >}}
+++
title = '2024 User Conference"
eventDate = 2024-04-01
+++
{{< /code >}}

To return a collection of future events:

```go-html-template
{{ $events := where .Site.RegularPages "Type" "events" }}
{{ $futureEvents := where $events "Params.eventDate" "gt" now }}
```

When working with YAML or JSON, or quoted TOML values, custom dates are strings; you cannot compare them with `time.Time` values. String comparisons may be possible if the custom date layout is consistent from one page to the next. To be safe, filter the pages by ranging through the collection:

```go-html-template
{{ $events := where .Site.RegularPages "Type" "events" }}
{{ $futureEvents := slice }}
{{ range $events }}
  {{ if gt (time.AsTime .Params.eventDate) now }}
    {{ $futureEvents = $futureEvents | append . }}
  {{ end }}
{{ end }}
```

## Nil comparison

To return a collection of pages where the "color" parameter is present in front matter, compare to `nil`:

```go-html-template
{{ $pages := where .Site.RegularPages "Params.color" "ne" nil }}
```

To return a collection of pages where the "color" parameter is not present in front matter, compare to `nil`:

```go-html-template
{{ $pages := where .Site.RegularPages "Params.color" "eq" nil }}
```

In both examples above, note that `nil` is not quoted.

## Nested comparison

These are equivalent:

```go-html-template
{{ $pages := where .Site.RegularPages "Type" "tutorials" }}
{{ $pages = where $pages "Params.level" "eq" "beginner" }}
```

```go-html-template
{{ $pages := where (where .Site.RegularPages "Type" "tutorials") "Params.level" "eq" "beginner" }}
```

## Portable section comparison

Useful for theme authors, avoid hardcoding section names by using the `where` function with the [`MainSections`] method on a `Site` object.

[`MainSections`]: /methods/site/mainsections/

```go-html-template
{{ $pages := where .Site.RegularPages "Section" "in" .Site.MainSections }}
```

With this construct, a theme author can instruct users to specify their main sections in the site configuration:

{{< code-toggle file=hugo >}}
[params]
mainSections = ['blog','galleries']
{{< /code-toggle >}}

If `params.mainSections` is not defined in the site configuration, the `MainSections` method returns a slice with one element---the top level section with the most pages.

## Boolean/undefined comparison

Consider this site content:

```text
content/
├── posts/
│   ├── _index.md
│   ├── post-1.md  <-- front matter: exclude = false
│   ├── post-2.md  <-- front matter: exclude = true
│   └── post-3.md  <-- front matter: exclude not defined
└── _index.md
```

The first two pages have an "exclude" field in front matter, but the last page does not. When testing for _equality_, the third page is _excluded_ from the result. When testing for _inequality_, the third page is _included_ in the result.

### Equality test

This template:

```go-html-template
<ul>
  {{ range where .Site.RegularPages "Params.exclude" "eq" false }}
    <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
  {{ end }}
</ul>
```

Is rendered to:

```html
<ul>
  <li><a href="/posts/post-1/">Post 1</a></li>
</ul>
```

This template:

```go-html-template
<ul>
  {{ range where .Site.RegularPages "Params.exclude" "eq" true }}
    <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
  {{ end }}
</ul>
```

Is rendered to:

```html
<ul>  
  <li><a href="/posts/post-2/">Post 2</a></li>
</ul>
```

### Inequality test

This template:

```go-html-template
<ul>
  {{ range where .Site.RegularPages "Params.exclude" "ne" false }}
    <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
  {{ end }}
</ul>
```

Is rendered to:

```html
<ul>
  <li><a href="/posts/post-2/">Post 2</a></li>
  <li><a href="/posts/post-3/">Post 3</a></li>
</ul>
```

This template:

```go-html-template
<ul>
  {{ range where .Site.RegularPages "Params.exclude" "ne" true }}
    <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
  {{ end }}
</ul>
```

Is rendered to:

```html
<ul>
  <li><a href="/posts/post-1/">Post 1</a></li>
  <li><a href="/posts/post-3/">Post 3</a></li>
</ul>
```

To exclude a page with an undefined field from a boolean _inequality_ test:

1. Create a collection using a boolean comparison
2. Create a collection using a nil comparison
3. Subtract the second collection from the first collection using the [`collections.Complement`] function.

[`collections.Complement`]: /functions/collections/complement/

This template:

```go-html-template
{{ $p1 := where .Site.RegularPages "Params.exclude" "ne" true }}
{{ $p2 := where .Site.RegularPages "Params.exclude" "eq" nil  }}
<ul>
  {{ range $p1 | complement $p2 }}
    <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
  {{ end }}
</ul>
```

Is rendered to:

```html
<ul>
  <li><a href="/posts/post-1/">Post 1</a></li>
</ul>
```

This template:

```go-html-template
{{ $p1 := where .Site.RegularPages "Params.exclude" "ne" false }}
{{ $p2 := where .Site.RegularPages "Params.exclude" "eq" nil  }}
<ul>
  {{ range $p1 | complement $p2 }}
    <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
  {{ end }}
</ul>
```

Is rendered to:

```html
<ul>
  <li><a href="/posts/post-1/">Post 2</a></li>
</ul>
```
