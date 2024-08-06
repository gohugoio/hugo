---
title: collections.Partition
description: Partitions a slice of pages into a slice of slices of pages with fixed (maximal) length.
categories: []
keywords: []
action:
  related: []
  returnType: [][]page.Pages
  signatures: [collections.Partition N PAGES]
---

{{< new-in 0.131.0 >}}

A set of 5 pages can be partitioned as show in this example:

```go-html-template
{{ .Site.RegularPages | collections.Partition 5 }} → [Pages(5)]
{{ .Site.RegularPages | collections.Partition 4 }} → [Pages(4) Pages(1)]
{{ .Site.RegularPages | collections.Partition 3 }} → [Pages(3) Pages(2)]
{{ .Site.RegularPages | collections.Partition 2 }} → [Pages(2) Pages(2) Pages(1)]
{{ .Site.RegularPages | collections.Partition 1 }} → [Pages(1) Pages(1) Pages(1) Pages(1) Pages(1)]
```
