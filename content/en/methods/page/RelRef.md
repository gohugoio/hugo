---
title: RelRef
description: Returns the relative URL of the page with the given path, language, and output format.
categories: []
keywords: []
action:
  related:
    - methods/page/Ref
    - functions/urls/Ref
    - functions/urls/RelRef
  returnType: string
  signatures: [PAGE.RelRef OPTIONS]
---

The map of option contains:

path
: (`string`) The path to the page, relative to the content directory. Required.

lang
: (`string`) The language (site) to search for the page. Default is the current language. Optional.

outputFormat
: (`string`) The output format to search for the page. Default is the current output format. Optional.

The examples below show the rendered output when visiting a page on the English language version of the site:

```go-html-template
{{ $opts := dict "path" "/books/book-1" }}
{{ .RelRef $opts }} → /en/books/book-1/

{{ $opts := dict "path" "/books/book-1" "lang" "de" }}
{{ .RelRef $opts }} → /de/books/book-1/

{{ $opts := dict "path" "/books/book-1" "lang" "de" "outputFormat" "json" }}
{{ .RelRef $opts }} → /de/books/book-1/index.json
```

By default, Hugo will throw an error and fail the build if it cannot resolve the path. You can change this to a warning in your site configuration, and specify a URL to return when the path cannot be resolved.

{{< code-toggle file=hugo >}}
refLinksErrorLevel = 'warning'
refLinksNotFoundURL = '/some/other/url'
{{< /code-toggle >}}
