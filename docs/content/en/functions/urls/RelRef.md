---
title: urls.RelRef
description: Returns the relative URL of the page with the given path, language, and output format.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [relref]
    returnType: string
    signatures:
      - urls.RelRef PAGE PATH
      - urls.RelRef PAGE OPTIONS
aliases: [/functions/relref]
---

## Usage

The `relref` function takes two arguments:

1. The context for resolving relative paths (typically the current page).
1. Either the target page's path or an options map (see below).

## Options

{{% include "_common/ref-and-relref-options.md" %}}

## Examples

The following examples show the rendered output for a page on the English version of the site:

```go-html-template
{{ relref . "/books/book-1" }} → /en/books/book-1/

{{ $opts := dict "path" "/books/book-1" }}
{{ relref . $opts }} → /en/books/book-1/

{{ $opts := dict "path" "/books/book-1" "lang" "de" }}
{{ relref . $opts }} → /de/books/book-1/

{{ $opts := dict "path" "/books/book-1" "lang" "de" "outputFormat" "json" }}
{{ relref . $opts }} → /de/books/book-1/index.json
```

## Error handling

{{% include "_common/ref-and-relref-error-handling.md" %}}
