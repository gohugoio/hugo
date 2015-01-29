---
aliases:
- /indexes/overview/
- /doc/indexes/
- /extras/indexes
date: 2013-07-01
linktitle: Overview
menu:
  main:
    identifier: taxonomy overview
    parent: taxonomy
next: /taxonomies/usage
prev: /templates/404
title: Taxonomy Overview
weight: 10
---

Hugo includes support for user-defined groupings of content called
taxonomies.[^1] Taxonomies give us a way to classify our content so we can
demonstrate relationships in a variety of logical ways.

[^1]: Taxonomies were called *indexes* in Hugo prior to v0.11.

The default taxonomies for Hugo are *tags* and *categories*. These
taxonomies are common to many website systems (e.g. WordPress, Drupal,
Jekyll). Unlike all of those systems, Hugo makes it trivial to customize
the taxonomies you will be using for your site however you wish. Another
good use for taxonomies is to group a set of posts into a series. Other
common uses would include *categories*, *tags*, *groups*, *series* and many
more.

When taxonomies are used (and templates are provided), Hugo will
automatically create pages listing all of the taxonomies, their terms
and all of the content attached to those terms.

## Definitions

**Taxonomy:** A categorization that can be used to classify content

**Term:** A key within that taxonomy 

**Value:** A piece of content assigned to that Term

## Example

For example, if I was writing about movies, I may want the following
taxonomies:

* Actors
* Directors
* Studios
* Genre
* Year
* Awards

I would then specify in each movie’s front-matter the specific terms for
each of those taxonomies. Hugo would then automatically create pages for
each Actor, Director, Studio, Genre, Year and Award listing all of the
Movies that matched that specific Actor, Director, etc.


### Taxonomy Organization

Let’s use an example to demonstrate the different labels in action.
From the perspective of the taxonomy, it could be visualized as:

    Actor                    <- Taxonomy
        Bruce Willis         <- Term
            The Six Sense    <- Content
            Unbreakable      <- Content
            Moonrise Kingdom <- Content
        Samuel L. Jackson    <- Term
            Unbreakable      <- Content
            The Avengers     <- Content
            xXx              <- Content

From the perspective of the content, it would appear differently, though
the data and labels used are the same:

    Unbreakable                 <- Content
        Actors                  <- Taxonomy
            Bruce Willis        <- Term
            Samuel L. Jackson   <- Term
        Director                <- Taxonomy
            M. Night Shyamalan  <- Term
        ...
    Moonrise Kingdom            <- Content
        Actors                  <- Taxonomy
            Bruce Willis        <- Term
            Bill Murray         <- Term
        Director                <- Taxonomy
            Wes Anderson        <- Term
        ...

