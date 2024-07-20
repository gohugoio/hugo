---
title: Ref
description: Returns the absolute URL of the page with the given path, language, and output format.
categories: []
keywords: []
action:
  related:
    - methods/shortcode/RelRef
    - functions/urls/RelRef
    - functions/urls/Ref
  returnType: string
  signatures: [SHORTCODE.Ref OPTIONS]
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
{{ .Ref $opts }} → https://example.org/en/books/book-1/

{{ $opts := dict "path" "/books/book-1" "lang" "de" }}
{{ .Ref $opts }} → https://example.org/de/books/book-1/

{{ $opts := dict "path" "/books/book-1" "lang" "de" "outputFormat" "json" }}
{{ .Ref $opts }} → https://example.org/de/books/book-1/index.json
```

By default, Hugo will throw an error and fail the build if it cannot resolve the path. You can change this to a warning in your site configuration, and specify a URL to return when the path cannot be resolved.

{{< code-toggle file=hugo >}}
refLinksErrorLevel = 'warning'
refLinksNotFoundURL = '/some/other/url'
{{< /code-toggle >}}
