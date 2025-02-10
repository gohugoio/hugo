---
title: data.GetJSON
description: Returns a JSON object from a local or remote JSON file, or an error if the file does not exist.
categories: []
keywords: []
action:
  aliases: [getJSON]
  related: []
  returnType: any
  signatures: ['data.GetJSON INPUT... [OPTIONS]']
toc: true
expiryDate: 2026-02-19 # deprecated 2024-02-19 in v0.123.0
---

{{% deprecated-in 0.123.0 %}}
Instead, use [`transform.Unmarshal`] with a [global resource](g), [page resource](g), or [remote resource](g).

See the [remote data example].

[`transform.Unmarshal`]: /functions/transform/unmarshal/
[remote data example]: /functions/resources/getremote/#remote-data
{{% /deprecated-in %}}

Given the following directory structure:

```text
my-project/
└── other-files/
    └── books.json
```

Access the data with either of the following:

```go-html-template
{{ $data := getJSON "other-files/books.json" }}
{{ $data := getJSON "other-files/" "books.json" }}
```

{{% note %}}
When working with local data, the file path is relative to the working directory.
{{% /note %}}

Access remote data with either of the following:

```go-html-template
{{ $data := getJSON "https://example.org/books.json" }}
{{ $data := getJSON "https://example.org/" "books.json" }}
```

The resulting data structure is a JSON object:

```json
[
  {
    "author": "Victor Hugo",
    "rating": 5,
    "title": "Les Misérables"
  },
  {
    "author": "Victor Hugo",
    "rating": 4,
    "title": "The Hunchback of Notre Dame"
  }
]
```

## Options

Add headers to the request by providing an options map:

```go-html-template
{{ $opts := dict "Authorization" "Bearer abcd" }}
{{ $data := getJSON "https://example.org/books.json" $opts }}
```

Add multiple headers using a slice:

```go-html-template
{{ $opts := dict "X-List" (slice "a" "b" "c") }}
{{ $data := getJSON "https://example.org/books.json" $opts }}
```

## Global resource alternative

Consider using the [`resources.Get`] function with [`transform.Unmarshal`] when accessing a global resource.

```text
my-project/
└── assets/
    └── data/
        └── books.json
```

```go-html-template
{{ $data := dict }}
{{ $p := "data/books.json" }}
{{ with resources.Get $p }}
  {{ $data = . | transform.Unmarshal }}
{{ else }}
  {{ errorf "Unable to get resource %q" $p }}
{{ end }}
```

## Page resource alternative

Consider using the [`Resources.Get`] method with [`transform.Unmarshal`] when accessing a page resource.

```text
my-project/
└── content/
    └── posts/
        └── reading-list/
            ├── books.json
            └── index.md
```

```go-html-template
{{ $data := dict }}
{{ $p := "books.json" }}
{{ with .Resources.Get $p }}
  {{ $data = . | transform.Unmarshal }}
{{ else }}
  {{ errorf "Unable to get resource %q" $p }}
{{ end }}
```

## Remote resource alternative

Consider using the [`resources.GetRemote`] function with [`transform.Unmarshal`] when accessing a remote resource to improve error handling and cache control.

```go-html-template
{{ $data := dict }}
{{ $url := "https://example.org/books.json" }}
{{ with try (resources.GetRemote $url) }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else with .Value }}
    {{ $data = . | transform.Unmarshal }}
  {{ else }}
    {{ errorf "Unable to get remote resource %q" $url }}
  {{ end }}
{{ end }}
```

[`Resources.Get`]: /methods/page/resources/
[`resources.GetRemote`]: /functions/resources/getremote/
[`resources.Get`]: /functions/resources/get/
[`transform.Unmarshal`]: /functions/transform/unmarshal/
