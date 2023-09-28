---
title: keyVals
description: Returns a KeyVals struct.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: collections
relatedFuncs: []
signature:
 - collections.KeyVals KEY VALUES...
 - keyVals KEY VALUES...

---

The primary application for this function is the definition of the `namedSlices` parameter in the options map passed to the `.Related` method on the `Page` object.

See [related content](/content-management/related).

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
