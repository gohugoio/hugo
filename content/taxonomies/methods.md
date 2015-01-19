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

**.Get(term)** Returns the WeightedPages for a term. <br>
**.Count(term)** The number of pieces of content assigned to this term.<br>
**.Alphabetical** Returns an OrderedTaxonomy (slice) ordered by Term. <br>
**.ByCount** Returns an OrderedTaxonomy (slice) ordered by number of entries. <br>

## OrderedTaxonomy

Since Maps are unordered, an OrderedTaxonomy is a special structure that has a defined order.

    []struct {
        Name          string
        WeightedPages WeightedPages
    }

Each element of the slice has:

**.Term**  The Term used.<br>
**.WeightedPages**  A slice of Weighted Pages.<br>
**.Count** The number of pieces of content assigned to this term.<br>
**.Pages**  All Pages assigned to this term. All [list methods](/templates/list/) are available to this.<br>

## WeightedPages

WeightedPages is simply a slice of WeightedPage.

    type WeightedPages []WeightedPage

**.Count(term)** The number of pieces of content assigned to this term.<br>
**.Pages** Returns a slice of pages, which then can be ordered using any of the [list methods](/templates/list/). <br>







