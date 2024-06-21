---
title: collections.Complement
description: Returns the elements of the last collection that are not in any of the others.
categories: []
keywords: []
action:
  aliases: [complement]
  related:
    - functions/collections/Intersect
    - functions/collections/SymDiff
    - functions/collections/Union
  returnType: any
  signatures: ['collections.Complement COLLECTION [COLLECTION...]']
aliases: [/functions/complement]
---

To find the elements within `$c3` that do not exist in `$c1` or `$c2`:

```go-html-template
{{ $c1 := slice 3 }}
{{ $c2 := slice 4 5 }}
{{ $c3 := slice 1 2 3 4 5 }}

{{ complement $c1 $c2 $c3 }} → [1 2]
```

{{% note %}}
Make your code simpler to understand by using a [chained pipeline]:

[chained pipeline]: https://pkg.go.dev/text/template#hdr-Pipelines
{{% /note %}}

```go-html-template
{{ $c3 | complement $c1 $c2 }} → [1 2]
```

You can also use the `complement` function with page collections. Let's say your site has five content types:

```text
content/
├── blog/
├── books/
├── faqs/
├── films/
└── songs/
```

To list everything except blog articles (`blog`) and frequently asked questions (`faqs`):

```go-html-template
{{ $blog := where site.RegularPages "Type" "blog" }}
{{ $faqs := where site.RegularPages "Type" "faqs" }}
{{ range site.RegularPages | complement $blog $faqs }}
  <a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
{{ end }}
```

{{% note %}}
Although the example above demonstrates the `complement` function, you could use the [`where`] function as well:

[`where`]: /functions/collections/where/
{{% /note %}}

```go-html-template
{{ range where site.RegularPages "Type" "not in" (slice "blog" "faqs") }}
  <a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
{{ end }}
```

In this example we use the `complement` function to remove [stop words] from a sentence:

```go-html-template
{{ $text := "The quick brown fox jumps over the lazy dog" }}
{{ $stopWords := slice "a" "an" "in" "over" "the" "under" }}
{{ $filtered := split $text " " | complement $stopWords }}

{{ delimit $filtered " " }} → The quick brown fox jumps lazy dog
```

[stop words]: https://en.wikipedia.org/wiki/Stop_word
