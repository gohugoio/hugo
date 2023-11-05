---
title: encoding.Jsonify
description: Encodes a given object to JSON.
categories: []
keywords: []
action:
  aliases: [jsonify]
  returnType: template.HTML
  related:
    - functions/transform/Remarshal
    - functions/transform/Unmarshal
  signatures:
    - encoding.Jsonify INPUT
    - encoding.Jsonify OPTIONS INPUT
aliases: [/functions/jsonify]
---

To customize the printing of the JSON, pass an options map as the first
argument.  Supported options are "prefix" and "indent".  Each JSON element in
the output will begin on a new line beginning with *prefix* followed by one or
more copies of *indent* according to the indentation nesting.

```go-html-template
{{ dict "title" .Title "content" .Plain | jsonify }}
{{ dict "title" .Title "content" .Plain | jsonify (dict "indent" "  ") }}
{{ dict "title" .Title "content" .Plain | jsonify (dict "prefix" " " "indent" "  ") }}
```

## Options

indent
: (`string`) Indentation to use. Default is "".

prefix
: (`string`) Indentation prefix. Default is "".

noHTMLEscape (false)
: (`bool`) Disable escaping of problematic HTML characters inside JSON quoted strings. The default behavior is to escape `&`, `<`, and `>` to `\u0026`, `\u003c`, and `\u003e` to avoid certain safety problems that can arise when embedding JSON in HTML. Default is `false`.

See also the `.PlainWords`, `.Plain`, and `.RawContent` [page variables].

[page variables]: /variables/page/
