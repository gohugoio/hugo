---
title: collections.KeyVals
description: Returns a KeyVals struct.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [keyVals]
    returnType: types.KeyValues
    signatures: [collections.KeyVals KEY VALUE...]
aliases: [/functions/keyvals]
---

The primary application for this function is the definition of the `namedSlices` value in the options map passed to the [`Related`] method on the `Pages` object.

[`Related`]: /methods/pages/related/

See [related content](/content-management/related-content/).

```go-html-template
{{ $kv := keyVals "foo" "a" "b" "c" }}
```

The resulting data structure is:

```json
{
  "Key": "foo",
  "Values": [
    "a",
    "b",
    "c"
  ]
}
```

To extract the key and values:

```go-html-template
{{ $kv.Key }} → foo
{{ $kv.Values }} → [a b c]
```
