---
title: data.GetCSV
linkTitle: getCSV
description: Returns an array of arrays from a local or remote CSV file, or an error if the file does not exist.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [getCSV]
  returnType: '[]string'
  signatures: [data.GetCSV SEPARATOR PATHPART...]
relatedFunctions:
  - data.GetCSV
  - data.GetJSON
toc: true
---

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

## Global resource alternative

Consider using `resources.Get` with [`transform.Unmarshal`] when accessing a global resource.

```text
my-project/
└── assets/
    └── data/
        └── pets.csv
```

```go-html-template
{{ $data := "" }}
{{ $p := "data/pets.csv" }}
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
        └── my-pets/
            ├── index.md
            └── pets.csv
```

```go-html-template
{{ $data := "" }}
{{ $p := "pets.csv" }}
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

[`transform.Unmarshal`]: /functions/transform/unmarshal
