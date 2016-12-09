---
lastmod: 2015-12-23
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

.Get(term)
: Returns the WeightedPages for a term.

.Count(term)
: The number of pieces of content assigned to this term.

.Alphabetical
: Returns an OrderedTaxonomy (slice) ordered by Term.

.ByCount
: Returns an OrderedTaxonomy (slice) ordered by number of entries.

## OrderedTaxonomy

Since Maps are unordered, an OrderedTaxonomy is a special structure that has a defined order.

```go
[]struct {
	Name          string
	WeightedPages WeightedPages
}
```

Each element of the slice has:

.Term
: The Term used.

.WeightedPages
: A slice of Weighted Pages.

.Count
: The number of pieces of content assigned to this term.

.Pages
: All Pages assigned to this term. All [list methods](/templates/list/) are available to this.

## WeightedPages

WeightedPages is simply a slice of WeightedPage.

```go
type WeightedPages []WeightedPage
```

.Count(term)
: The number of pieces of content assigned to this term.

.Pages
: Returns a slice of pages, which then can be ordered using any of the [list methods](/templates/list/).
