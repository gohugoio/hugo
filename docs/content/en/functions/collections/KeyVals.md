---
title: collections.KeyVals
linkTitle: keyVals
description: Returns a KeyVals struct.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [keyVals]
  returnType: KeyValues
  signatures: [collections.KeyVals KEY VALUES...]
relatedFunctions: []
aliases: [/functions/keyvals]
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
