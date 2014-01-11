---
title: "Templates"
date: "2013-07-01"
aliases: ["/doc/templates/"]
linktitle: "Overview"
groups: ["layout"]
groups_weight: 10
---

Hugo uses the excellent golang html/template library for its template engine.
It is an extremely lightweight engine that provides a very small amount of
logic. In our experience that it is just the right amount of logic to be able
to create a good static website

If you are new to go's templates the [go template primer](/layout/go-templates)
is a great place to start.

## Template roles

There are 6 different kinds of templates that Hugo works with.

### [Homepage](/layout/homepage/)
The homepage of your site

### [RSS](/layout/rss/)
Used to render all rss documents

### [Index](/layout/indexes)
Page that list multiple pieces of content

### [Content](/layout/content)
Render a single piece of content

### [Views](/layout/views)
Different view of a single piece of content type

### [Chrome](/layout/chrome)
Support for the above templates
