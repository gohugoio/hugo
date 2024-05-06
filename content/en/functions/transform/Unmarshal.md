---
title: transform.Unmarshal
description: Parses serialized data and returns a map or an array. Supports CSV, JSON, TOML, YAML, and XML.
categories: []
keywords: []
action:
  aliases: [unmarshal]
  related:
    - functions/transform/Remarshal
    - functions/resources/Get
    - functions/resources/GetRemote
    - functions/encoding/Jsonify
  returnType: any
  signatures: ['transform.Unmarshal [OPTIONS] INPUT']
toc: true
aliases: [/functions/transform.unmarshal]
---

The input can be a string or a [resource].

## Unmarshal a string

```go-html-template
{{ $string := `
title: Les Misérables
author: Victor Hugo
`}}

{{ $book := unmarshal $string }}
{{ $book.title }} → Les Misérables
{{ $book.author }} → Victor Hugo
```

## Unmarshal a resource

Use the `transform.Unmarshal` function with global, page, and remote resources.

### Global resource

A global resource is a file within the assets directory, or within any directory mounted to the assets directory.

```text
assets/
└── data/
    └── books.json
```

```go-html-template
{{ $data := dict }}
{{ $path := "data/books.json" }}
{{ with resources.Get $path }}
  {{ with . | transform.Unmarshal }}
    {{ $data = . }}
  {{ end }}
{{ else }}
  {{ errorf "Unable to get global resource %q" $path }}
{{ end }}

{{ range where $data "author" "Victor Hugo" }}
  {{ .title }} → Les Misérables
{{ end }}
```

### Page resource

A page resource is a file within a [page bundle].

```text
content/
├── post/
│   └── book-reviews/
│       ├── books.json
│       └── index.md
└── _index.md
```

```go-html-template
{{ $data := dict }}
{{ $path := "books.json" }}
{{ with .Resources.Get $path }}
  {{ with . | transform.Unmarshal }}
    {{ $data = . }}
  {{ end }}
{{ else }}
  {{ errorf "Unable to get page resource %q" $path }}
{{ end }}

{{ range where $data "author" "Victor Hugo" }}
  {{ .title }} → Les Misérables
{{ end }}
```

### Remote resource

A remote resource is a file on a remote server, accessible via HTTP or HTTPS.

```go-html-template
{{ $data := dict }}
{{ $url := "https://example.org/books.json" }}
{{ with resources.GetRemote $url }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else }}
    {{ $data = . | transform.Unmarshal }}
  {{ end }}
{{ else }}
  {{ errorf "Unable to get remote resource %q" $url }}
{{ end }}

{{ range where $data "author" "Victor Hugo" }}
  {{ .title }} → Les Misérables
{{ end }}
```

{{% note %}}
When retrieving remote data, a misconfigured server may send a response header with an incorrect [Content-Type]. For example, the server may set the Content-Type header to `application/octet-stream` instead of `application/json`.

In these cases, pass the resource `Content` through the `transform.Unmarshal` function instead of passing the resource itself. For example, in the above, do this instead:

`{{ $data = .Content | transform.Unmarshal }}`

[Content-Type]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
{{% /note %}}

## Options

When unmarshaling a CSV file, provide an optional map of options.

delimiter
: (`string`) The delimiter used, default is `,`.

comment
: (`string`) The comment character used in the CSV. If set, lines beginning with the comment character without preceding whitespace are ignored.

lazyQuotes {{< new-in 0.122.0 >}}
: (`bool`) If true, a quote may appear in an unquoted field and a non-doubled quote may appear in a quoted field. Default is `false`.

```go-html-template
{{ $csv := "a;b;c" | transform.Unmarshal (dict "delimiter" ";") }}
```

## Working with XML

When unmarshaling an XML file, do not include the root node when accessing data. For example, after unmarshaling the RSS feed below, access the feed title with `$data.channel.title`.

```xml
<?xml version="1.0" encoding="utf-8" standalone="yes"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>Books on Example Site</title>
    <link>https://example.org/books/</link>
    <description>Recent content in Books on Example Site</description>
    <language>en-US</language>
    <atom:link href="https://example.org/books/index.xml" rel="self" type="application/rss+xml" />
    <item>
      <title>The Hunchback of Notre Dame</title>
      <description>Written by Victor Hugo</description>
      <link>https://example.org/books/the-hunchback-of-notre-dame/</link>
      <pubDate>Mon, 09 Oct 2023 09:27:12 -0700</pubDate>
      <guid>https://example.org/books/the-hunchback-of-notre-dame/</guid>
    </item>
    <item>
      <title>Les Misérables</title>
      <description>Written by Victor Hugo</description>
      <link>https://example.org/books/les-miserables/</link>
      <pubDate>Mon, 09 Oct 2023 09:27:11 -0700</pubDate>
      <guid>https://example.org/books/les-miserables/</guid>
    </item>
  </channel>
</rss>
```

Get the remote data:

```go-html-template
{{ $data := dict }}
{{ $url := "https://example.org/books/index.xml" }}
{{ with resources.GetRemote $url }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else }}
    {{ $data = . | transform.Unmarshal }}
  {{ end }}
{{ else }}
  {{ errorf "Unable to get remote resource %q" $url }}
{{ end }}
```

Inspect the data structure:

```go-html-template
<pre>{{ debug.Dump $data }}</pre>
```

List the book titles:

```go-html-template
{{ with $data.channel.item }}
  <ul>
    {{ range . }}
      <li>{{ .title }}</li>
    {{ end }}
  </ul>
{{ end }}
```

Hugo renders this to:

```html
<ul>
  <li>The Hunchback of Notre Dame</li>
  <li>Les Misérables</li>
</ul>
```

### XML attributes and namespaces

Let's add a `lang` attribute to the `title` nodes of our RSS feed, and a namespaced node for the ISBN number:

```xml
<?xml version="1.0" encoding="utf-8" standalone="yes"?>
<rss version="2.0"
  xmlns:atom="http://www.w3.org/2005/Atom"
  xmlns:isbn="http://schemas.isbn.org/ns/1999/basic.dtd"
>
  <channel>
    <title>Books on Example Site</title>
    <link>https://example.org/books/</link>
    <description>Recent content in Books on Example Site</description>
    <language>en-US</language>
    <atom:link href="https://example.org/books/index.xml" rel="self" type="application/rss+xml" />
    <item>
      <title lang="fr">The Hunchback of Notre Dame</title>
      <description>Written by Victor Hugo</description>
      <isbn:number>9780140443530</isbn:number>
      <link>https://example.org/books/the-hunchback-of-notre-dame/</link>
      <pubDate>Mon, 09 Oct 2023 09:27:12 -0700</pubDate>
      <guid>https://example.org/books/the-hunchback-of-notre-dame/</guid>
    </item>
    <item>
      <title lang="en">Les Misérables</title>
      <description>Written by Victor Hugo</description>
      <isbn:number>9780451419439</isbn:number>
      <link>https://example.org/books/les-miserables/</link>
      <pubDate>Mon, 09 Oct 2023 09:27:11 -0700</pubDate>
      <guid>https://example.org/books/les-miserables/</guid>
    </item>
  </channel>
</rss>
```

After retrieving the remote data, inspect the data structure:

```go-html-template
<pre>{{ debug.Dump $data }}</pre>
```

Each item node looks like this:

```json
{
  "description": "Written by Victor Hugo",
  "guid": "https://example.org/books/the-hunchback-of-notre-dame/",
  "link": "https://example.org/books/the-hunchback-of-notre-dame/",
  "number": "9780140443530",
  "pubDate": "Mon, 09 Oct 2023 09:27:12 -0700",
  "title": {
    "#text": "The Hunchback of Notre Dame",
    "-lang": "fr"
  }
}
```

The title keys do not begin with an underscore or a letter---they are not valid [identifiers]. Use the [`index`] function to access the values:

```go-html-template
{{ with $data.channel.item }}
  <ul>
    {{ range . }}
      {{ $title := index .title "#text" }}
      {{ $lang := index .title "-lang" }}
      {{ $ISBN := .number }}
      <li>{{ $title }} ({{ $lang }}) {{ $ISBN }}</li>
    {{ end }}
  </ul>
{{ end }}
```

Hugo renders this to:

```html
<ul>
  <li>The Hunchback of Notre Dame (fr) 9780140443530</li>
  <li>Les Misérables (en) 9780451419439</li>
</ul>
```

[`index`]: /functions/collections/indexfunction/
[identifiers]: https://go.dev/ref/spec#Identifiers
[resource]: /getting-started/glossary/#resource
[page bundle]: /content-management/page-bundles/
