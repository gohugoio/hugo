---
title: collections.D
description: Returns a slice of sequentially ordered random integers.
categories: []
keywords: [random]
params:
  functions_and_methods:
    returnType: '[]int'
    signatures: [collections.D SEED N HIGH]
---

{{< new-in 0.149.0 />}}

The `collections.D` function returns a slice of `N` sequentially ordered unique random integers in the half-open [interval](g) [0, `HIGH`) using the provided [`SEED`](g) value. This function implements J. S. Vitter's Method&nbsp;D[^1] for sequential random sampling, a fast and efficient algorithm for this task.

See [this article][] for a detailed explanation.

## Examples

```go-html-template
{{ collections.D 6 7 42 }} → [4, 9, 10, 20, 22, 24, 41]
```

The example above generates the _same_ random numbers each time it is called. To generate a _different_ set of 7 random numbers in the same range, change the seed value.

```go-html-template
{{ collections.D 2 7 42 }} → [3, 11, 19, 25, 32, 33, 38]
```

A common use case is the selection of random pages from a page collection. For example, to render a list of 5 random pages using the [day of the year][] as the seed value:

```go-html-template
<ul>
  {{ $p := site.RegularPages }}
  {{ range collections.D time.Now.YearDay 5 ($p | len) }}
    {{ with (index $p .) }}
      <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
    {{ end }}
  {{ end }}
</ul>
```

The construct above is significantly faster than using the [`collections.Shuffle`][] function.

## Seed value

Choosing an appropriate seed value depends on your objective.

Objective|Seed example
:--|:--
Consistent result|`42`
Different result on each call|`int time.Now.UnixNano`
Same result per day|`time.Now.YearDay`
Same result per page|`hash.FNV32a .Path`
Different result per page per day|`hash.FNV32a (print .Path time.Now.YearDay)`

> [!note]
> The slice created by this function is limited to 1 million elements.

[^1]: J. S. Vitter, "An efficient algorithm for sequential random sampling," _ACM Trans. Math. Soft._, vol. 13, pp. 58&ndash;67, Mar. 1987.

[`collections.Shuffle`]: /functions/collections/shuffle/
[day of the year]: /methods/time/yearday/
[this article]: https://getkerf.wordpress.com/2016/03/30/the-best-algorithm-no-one-knows-about/
