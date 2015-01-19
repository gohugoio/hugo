---
date: 2014-05-26
linktitle: Structure & Methods
menu:
  main:
    parent: taxonomy
next: /extras/aliases
prev: /taxonomies/ordering
title: Using Taxonomies
weight: 75
---

Hugo makes a set of values and methods available on the various Taxonomy structures.

## Taxonomy Methods

A Taxonomy is a `map[string]WeightedPages`.

<dl>
<dt><code>.Get(term)</code></dt><dd>Returns the WeightedPages for a term.</dd>
<dt><code>.Count(term)</code></dt><dd>The number of pieces of content assigned to this term.</dd>
<dt><code>.Alphabetical</code></dt><dd>Returns an OrderedTaxonomy (slice) ordered by Term.</dd>
<dt><code>.ByCount</code></dt><dd>Returns an OrderedTaxonomy (slice) ordered by number of entries.</dd>
</dl>

## OrderedTaxonomy

Since Maps are unordered, an OrderedTaxonomy is a special structure that has a defined order.

    []struct {
        Name          string
        WeightedPages WeightedPages
    }

Each element of the slice has:

<dl>
<dt><code>.Term</code></dt><dd>The Term used.</dd>
<dt><code>.WeightedPages</code></dt><dd>A slice of Weighted Pages.</dd>
<dt><code>.Count</code></dt><dd>The number of pieces of content assigned to this term.</dd>
<dt><code>.Pages</code></dt><dd>All Pages assigned to this term. All <a href="/templates/list/">list methods</a> are available to this.</dd>
</dl>

## WeightedPages

WeightedPages is simply a slice of WeightedPage.

    type WeightedPages []WeightedPage

<dl>
<dt><code>.Count(term)</code></dt><dd>The number of pieces of content assigned to this term.</dd>
<dt><code>.Pages</code></dt><dd>Returns a slice of pages, which then can be ordered using any of the <a href="/templates/list/">list methods</a>.</dd>
</dl>
