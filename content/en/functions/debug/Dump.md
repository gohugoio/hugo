---
title: debug.Dump
description: Returns an object dump as a string.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: string
  signatures: [debug.Dump VALUE]
relatedFunctions: []
---

```go-html-template
{{ $data := "" }}
{{ $p := "data/books.json" }}
{{ with resources.Get $p }}
  {{ $opts := dict "delimiter" "," }}
  {{ $data = . | transform.Unmarshal $opts }}
{{ else }}
  {{ errorf "Unable to get resource %q" $p }}
{{ end }}
```

```go-html-template
<pre>{{ debug.Dump $data }}</pre>
```

```text
[]interface {}{
  map[string]interface {}{
    "author": "Victor Hugo",
    "rating": 5.0,
    "title": "Les Misérables",
  },
  map[string]interface {}{
    "author": "Victor Hugo",
    "rating": 4.0,
    "title": "The Hunchback of Notre Dame",
  },
}
```



{{% note %}}
Output from this function may change from one release to the next. Use for debugging only.
{{% /note %}}
