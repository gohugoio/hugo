---
title: urls.PathEscape 
description: Returns the given string, applying percent-encoding to special characters and reserved delimiters so it can be safely used as a segment within a URL path.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [urls.PathEscape INPUT]
---

{{< new-in v0.153.0 />}}

The `urls.PathEscape` function does the inverse transformation of [`urls.PathUnescape`][].

```go-html-template
{{ urls.PathEscape "my café" }} → my%20caf%C3%A9
```

Use this function to escape a string so that it can be safely used as an individual segment within a URL path.

[`urls.PathUnescape`]: /functions/urls/PathUnescape/
