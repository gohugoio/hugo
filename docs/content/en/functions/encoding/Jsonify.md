---
title: encoding.Jsonify
linkTitle: jsonify
description: Encodes a given object to JSON.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [jsonify]
  returnType: template.HTML
  signatures:
    - encoding.Jsonify INPUT
    - encoding.Jsonify OPTIONS INPUT
relatedFunctions:
  - encoding.Jsonify
  - transform.Remarshal
  - transform.Unmarshal
aliases: [/functions/jsonify]
---

To customize the printing of the JSON, pass a map of options as the first
argument.  Supported options are "prefix" and "indent".  Each JSON element in
the output will begin on a new line beginning with *prefix* followed by one or
more copies of *indent* according to the indentation nesting.


```go-html-template
{{ dict "title" .Title "content" .Plain | jsonify }}
{{ dict "title" .Title "content" .Plain | jsonify (dict "indent" "  ") }}
{{ dict "title" .Title "content" .Plain | jsonify (dict "prefix" " " "indent" "  ") }}
```

## Options

indent ("")
: Indentation to use.

prefix ("")
: Indentation prefix.

noHTMLEscape (false)
: Disable escaping of problematic HTML characters inside JSON quoted strings. The default behavior is to escape &, <, and > to \u0026, \u003c, and \u003e to avoid certain safety problems that can arise when embedding JSON in HTML.

See also the `.PlainWords`, `.Plain`, and `.RawContent` [page variables][pagevars].

[pagevars]: /variables/page/
