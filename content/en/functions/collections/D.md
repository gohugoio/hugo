---
title: collections.D
description: Returns a sorted slice of unique random integers based on a given seed, count, and maximum value.
categories: []
keywords: [random]
params:
  functions_and_methods:
    aliases: []
    returnType: '[]int'
    signatures: [collections.D SEED N HIGH]
---

{{< new-in 0.149.0 />}}

The `collections.D` function returns a sorted slice of unique random integers in the half-open [interval](g) `[0, HIGH)` using the provided [`SEED`](g) value. The number of elements in the resulting slice is `N` or `HIGH`, whichever is less.

- `N` and `H` must be integers in the closed interval `[0, 1000000]`
- `SEED` must be an integer in the closed interval `[0, 2^64 - 1]`

## Return values

Condition|Return value
:--|:--|:--
`N <= HIGH`|A sorted random sample of size `N` using J. S. Vitter's [Method D][] for sequential random sampling
`N > HIGH`|The full, sorted range `[0, HIGH)` of size `HIGH`
`N == 0`|An empty slice
`N < 0`|Error
`N > 10^6`|Error
`HIGH == 0`|An empty slice
`HIGH < 0`|Error
`HIGH > 10^6`|Error
`SEED < 0`|Error
{.no-wrap-first-col}

## Examples

```go-html-template
{{ collections.D 6 7 42 }} → [4, 9, 10, 20, 22, 24, 41]
```

The example above generates the _same_ random numbers each time it is called. To generate a _different_ set of 7 random numbers in the same range, change the seed value.

```go-html-template
{{ collections.D 2 7 42 }} → [3, 11, 19, 25, 32, 33, 38]
```

When `N` is greater than `HIGH`, this function returns the full, sorted range [0, `HIGH`) of size `HIGH`:

```go-html-template
{{ collections.D 6 42 7 }} → [0 1 2 3 4 5 6] 
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

[`collections.Shuffle`]: /functions/collections/shuffle/
[day of the year]: /methods/time/yearday/
[Method D]: https://getkerf.wordpress.com/2016/03/30/the-best-algorithm-no-one-knows-about/
