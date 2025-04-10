---
title: RelRef
description: Returns the relative URL of the page with the given path, language, and output format.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [PAGE.RelRef OPTIONS]
---

## Usage

The `RelRef` method accepts a single argument: an options map.

## Options

{{% include "_common/ref-and-relref-options.md" %}}

## Examples

The following examples show the rendered output for a page on the English version of the site:

```go-html-template
{{ $opts := dict "path" "/books/book-1" }}
{{ .RelRef $opts }} → /en/books/book-1/

{{ $opts := dict "path" "/books/book-1" "lang" "de" }}
{{ .RelRef $opts }} → /de/books/book-1/

{{ $opts := dict "path" "/books/book-1" "lang" "de" "outputFormat" "json" }}
{{ .RelRef $opts }} → /de/books/book-1/index.json
```

## Error handling

{{% include "_common/ref-and-relref-error-handling.md" %}}
