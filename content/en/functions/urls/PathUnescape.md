---
title: urls.PathUnescape 
description: Returns the given string, applying percent-encoding to special characters and reserved delimiters so it can be safely used as a segment within a URL path.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [urls.PathUnescape INPUT]
---

{{< new-in v0.153.0 />}}

The `urls.PathUnescape` function does the inverse transformation of [`urls.PathEscape`][].

```go-html-template
{{ urls.PathUnescape "A%2Fb%2Fc%3Fd=%C3%A9&f=g+h" }} → A/b/c?d=é&f=g+h
```

[`urls.PathEscape`]: /functions/urls/PathEscape/
