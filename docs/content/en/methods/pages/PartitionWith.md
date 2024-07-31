---
title: PartitionWith
description: Partitions a slice of pages into a slice of slices of pages with fixed (maximal) length.
categories: []
keywords: []
action:
  related: []
  returnType: [][]page.Pages
  signatures: [PAGES.PartitionWith N]
---

{{< new-in 0.131.0 >}}

This template shows how to partition a set of 5 pages into chunks of some size:

```go-html-template
{{ $allPages := .Site.RegularPages }}

<p>Number pages: {{ len $allPages }}</p>

<ul>
  <li>Chunk size: 5 => Partitions: {{ $allPages.PartitionWith 5 }}</li>
  <li>Chunk size: 4 => Partitions: {{ $allPages.PartitionWith 4 }}</li>
  <li>Chunk size: 3 => Partitions: {{ $allPages.PartitionWith 3 }}</li>
  <li>Chunk size: 2 => Partitions: {{ $allPages.PartitionWith 2 }}</li>
  <li>Chunk size: 1 => Partitions: {{ $allPages.PartitionWith 1 }}</li>
</ul>
```

Which generates the following output:

```html
<p>Number pages: 5</p>

<ul>
  <li>Chunk size: 5 => Partitions: [Pages(5)]</li>
  <li>Chunk size: 4 => Partitions: [Pages(4) Pages(1)]</li>
  <li>Chunk size: 3 => Partitions: [Pages(3) Pages(2)]</li>
  <li>Chunk size: 2 => Partitions: [Pages(2) Pages(2) Pages(1)]</li>
  <li>Chunk size: 1 => Partitions: [Pages(1) Pages(1) Pages(1) Pages(1) Pages(1)]</li>
</ul>
```
