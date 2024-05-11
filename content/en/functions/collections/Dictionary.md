---
title: collections.Dictionary
description: Returns a map composed of the given key-value pairs.
categories: []
keywords: []
action:
  aliases: [dict]
  related:
    - functions/collections/Slice
  returnType: mapany
  signatures: ['collections.Dictionary [VALUE...]']
aliases: [/functions/dict]
---

Specify the key-value pairs as individual arguments:

```go-html-template
{{ $m := dict "a" 1 "b" 2 }}
```

The above produces this data structure:

```json
{
  "a": 1,
  "b": 2
}
```

To create an empty map:

```go-html-template
{{ $m := dict }}
```


Note that the `key` can be either a `string` or a `string slice`. The latter is useful to create a deeply nested structure, e.g.:

```go-html-template
{{ $m := dict (slice "a" "b" "c") "value" }}
```

The above produces this data structure:

```json
{
  "a": {
    "b": {
      "c": "value"
    }
  }
}
```
