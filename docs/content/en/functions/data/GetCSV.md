---
title: data.GetCSV
description: Returns an array of arrays from a local or remote CSV file, or an error if the file does not exist.
categories: []
keywords: []
action:
  aliases: [getCSV]
  related:
    - functions/data/GetJSON
    - functions/resources/Get
    - functions/resources/GetRemote
    - methods/page/Resources
  returnType: '[][]string'
  signatures: ['data.GetCSV SEPARATOR INPUT... [OPTIONS]']
toc: true
expiryDate: 2025-02-19 # deprecated 2024-02-19
---

{{% deprecated-in 0.123.0 %}}
Instead, use [`transform.Unmarshal`] with a [global], [page], or [remote] resource.

See the [remote data example].

[`transform.Unmarshal`]: /functions/transform/unmarshal/
[global]: /getting-started/glossary/#global-resource
[page]: /getting-started/glossary/#page-resource
[remote data example]: /functions/resources/getremote/#remote-data
[remote]: /getting-started/glossary/#remote-resource
{{% /deprecated-in %}}

Given the following directory structure:

```text
my-project/
└── other-files/
    └── pets.csv
```

Access the data with either of the following:

```go-html-template
{{ $data := getCSV "," "other-files/pets.csv" }}
{{ $data := getCSV "," "other-files/" "pets.csv" }}
```

{{% note %}}
When working with local data, the file path is relative to the working directory.

You must not place CSV files in the project's data directory.
{{% /note %}}

Access remote data with either of the following:

```go-html-template
{{ $data := getCSV "," "https://example.org/pets.csv" }}
{{ $data := getCSV "," "https://example.org/" "pets.csv" }}
```

The resulting data structure is an array of arrays:

```json
[
  ["name","type","breed","age"],
  ["Spot","dog","Collie","3"],
  ["Felix","cat","Malicious","7"]
]
```

## Options

Add headers to the request by providing an options map:

```go-html-template
{{ $opts := dict "Authorization" "Bearer abcd" }}
{{ $data := getCSV "," "https://example.org/pets.csv" $opts }}
```

Add multiple headers using a slice:

```go-html-template
{{ $opts := dict "X-List" (slice "a" "b" "c") }}
{{ $data := getCSV "," "https://example.org/pets.csv" $opts }}
```

## Global resource alternative

Consider using the [`resources.Get`] function with [`transform.Unmarshal`] when accessing a global resource.

```text
my-project/
└── assets/
    └── data/
        └── pets.csv
```

```go-html-template
{{ $data := dict }}
{{ $p := "data/pets.csv" }}
{{ with resources.Get $p }}
  {{ $opts := dict "delimiter" "," }}
  {{ $data = . | transform.Unmarshal $opts }}
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
        └── my-pets/
            ├── index.md
            └── pets.csv
```

```go-html-template
{{ $data := dict }}
{{ $p := "pets.csv" }}
{{ with .Resources.Get $p }}
  {{ $opts := dict "delimiter" "," }}
  {{ $data = . | transform.Unmarshal $opts }}
{{ else }}
  {{ errorf "Unable to get resource %q" $p }}
{{ end }}
```

## Remote resource alternative

Consider using the [`resources.GetRemote`] function with [`transform.Unmarshal`] when accessing a remote resource to improve error handling and cache control.

```go-html-template
{{ $data := dict }}
{{ $u := "https://example.org/pets.csv" }}
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

[`Resources.Get`]: /methods/page/resources/
[`resources.GetRemote`]: /functions/resources/getremote/
[`resources.Get`]: /functions/resources/get/
[`transform.Unmarshal`]: /functions/transform/unmarshal/
