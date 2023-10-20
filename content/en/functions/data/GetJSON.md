---
title: data.GetJSON
linkTitle: getJSON
description: Returns a JSON object from a local or remote JSON file, or an error if the file does not exist.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [getJSON]
  returnType: any
  signatures: [data.GetJSON PATHPART...]
relatedFunctions:
  - data.GetCSV
  - data.GetJSON
toc: true
---

Given the following directory structure:

```text
my-project/
└── other-files/
    └── books.json
```

Access the data with either of the following:

```go-html-template
{{ $data := getCSV "," "other-files/books.json" }}
{{ $data := getCSV "," "other-files/" "books.json" }}
```

Access remote data with either of the following:

```go-html-template
{{ $data := getCSV "," "https://example.org/books.json" }}
{{ $data := getCSV "," "https://example.org/" "books.json" }}
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

## Global resource alternative

Consider using `resources.Get` with [`transform.Unmarshal`] when accessing a global resource.

```text
my-project/
└── assets/
    └── data/
        └── books.json
```

```go-html-template
{{ $data := "" }}
{{ $p := "data/books.json" }}
{{ with resources.Get $p }}
  {{ $opts := dict "delimiter" "," }}
  {{ $data = . | transform.Unmarshal $opts }}
{{ else }}
  {{ errorf "Unable to get resource %q" $p }}
{{ end }}
```

## Page resource alternative

Consider using `.Resources.Get` with [`transform.Unmarshal`] when accessing a page resource.

```text
my-project/
└── content/
    └── posts/
        └── reading-list/
            ├── books.json
            └── index.md
```

```go-html-template
{{ $data := "" }}
{{ $p := "books.json" }}
{{ with .Resources.Get $p }}
  {{ $opts := dict "delimiter" "," }}
  {{ $data = . | transform.Unmarshal $opts }}
{{ else }}
  {{ errorf "Unable to get resource %q" $p }}
{{ end }}
```

## Remote resource alternative

Consider using `resources.GetRemote` with [`transform.Unmarshal`] for improved error handling when accessing a remote resource.

```go-html-template
{{ $data := "" }}
{{ $u := "https://example.org/books.json" }}
{{ with resources.GetRemote $u }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else }}
    {{ $opts := dict "delimiter" "," }}
    {{ $data = . | transform.Unmarshal $opts }}
  {{ end }}
{{ else }}
  {{ errorf "Unable to get remote resource %q" $u }}
{{ end }}
```

[`transform.Unmarshal`]: /functions/transform/unmarshal
