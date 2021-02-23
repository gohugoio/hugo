---
title: "transform.Unmarshal"
description: "`transform.Unmarshal` (alias `unmarshal`) parses the input and converts it into a map or an array. Supported formats are JSON, TOML, YAML and CSV."
date: 2018-12-23
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: []
signature: ["RESOURCE or STRING | transform.Unmarshal [OPTIONS]"]
hugoversion: "0.53"
aliases: []
---

The function accepts either a `Resource` created in [Hugo Pipes](/hugo-pipes/) or via [Page Bundles](/content-management/page-bundles/), or simply a string. The two examples below will produce the same map:

```go-html-template
{{ $greetings := "hello = \"Hello Hugo\"" | transform.Unmarshal }}`
```

```go-html-template
{{ $greetings := "hello = \"Hello Hugo\"" | resources.FromString "data/greetings.toml" | transform.Unmarshal }}
```

In both the above examples, you get a map you can work with:

```go-html-template
{{ $greetings.hello }}
```

The above prints `Hello Hugo`.

## CSV Options

Unmarshal with CSV as input has some options you can set:

delimiter
: The delimiter used, default is `,`.

comment
: The comment character used in the CSV. If set, lines beginning with the comment character without preceding whitespace are ignored.:

Example:

```go-html-template
{{ $csv := "a;b;c" | transform.Unmarshal (dict "delimiter" ";") }}
```
