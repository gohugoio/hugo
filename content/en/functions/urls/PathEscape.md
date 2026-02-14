---
title: urls.PathEscape 
description: Returns the given string, replacing all percent-encoded sequences with the corresponding unescaped characters.
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
{{ urls.PathEscape "A/b/c?d=é&f=g+h" }} → A%2Fb%2Fc%3Fd=%C3%A9&f=g+h
```

[`urls.PathUnescape`]: /functions/urls/PathUnescape/
