---
title: seq
# linktitle:
description: Creates a sequence of integers.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: []
signature: ["seq LAST", "seq FIRST LAST", "seq FIRST INCREMENT LAST"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
---

It's named and used in the model of [GNU's seq][].

```
3 → 1, 2, 3
1 2 4 → 1, 3
-3 → -1, -2, -3
1 4 → 1, 2, 3, 4
1 -2 → 1, 0, -1, -2
```

## Example: `seq` with `range` and `after`

You can use `seq` in combination with `range` and `after`. The following will return 19 elements:

```
{{ range after 1 (seq 20)}}
{{ end }}
```

However, when ranging with an index, the following may be less confusing in that `$indexStartingAt1` and `$num` will return `1,2,3 ... 20`:

```
{{ range $index, $num := (seq 20) }}
$indexStartingAt1 := (add $index 1)
{{ end }}
```


[GNU's seq]: http://www.gnu.org/software/coreutils/manual/html_node/seq-invocation.html#seq-invocation
