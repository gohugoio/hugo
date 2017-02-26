---
title: Introduction to Lists in Hugo
linktitle: Hugo Lists
description: Lists have a specific meaning and usage in Hugo. If you want to know how to render your site homepage, section homepage, taxonomy list, or taxonomy terms list, start here.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [lists,sections,rss,taxonomies,terms]
weight: 22
draft: false
aliases: [/templates/lists-in-hugo/]
toc: true
needsreview: true
---

## What is a List Page Template?

A list page template is a template used to render multiple pieces of content in a single HTML page (with the exception of the homepage, which has a [dedicated template][homepage]).

Hugo uses the term *list* in its truest sense; i.e. a sequential arrangement of material, especially in alphabetical or numerical order. Hugo uses list templates on any output HTML page where content is traditionally listed:

* [Taxonomy terms pages][taxonomy terms]
* [Taxonomy list pages][taxonomy lists]
* [Section list pages][section templates]
* [RSS][]

The idea of a list page comes from the [hierarchical mental model of the web][mentalmodel] and is best demonstrated visually:

![Image demonstrating a hierarchical website sitemap.](/images/site-hierarchy.svg)




```yaml
---
aliases:
- /doc/using-index-md/
lastmod: 2017-02-22
date: 2017-02-22
linktitle: Using _index.md
menu:
  main:
    parent: content
prev: /content/example
next: /themes/overview
notoc: true
title: Using _index.md
weight: 70
---
```



## List Defaults

### Default Templates

Since section lists and taxonomy lists (N.B., *not* [taxonomy terms lists][]) are both *lists* with regards to their templates, both have the same terminating default of `_default/list.html`---or `themes/<THEME>/layouts/_default/list.html` in the case of a themed project---in their *lookup orders*. In addition, both [section lists][sectiontemplates] and [taxonomy lists][taxonomies] have their own default list templates in `_default`:

#### Default Section Templates

1. `layouts/section/<SECTIONNAME>.html`
2. `layouts/section/list.html`
3. `layouts/_default/section.html`
4. `layouts/_default/list.html`

### Understanding `.Data.Pages`


### Taxonomy Lists

A Taxonomy will be rendered at /`PLURAL`/`TERM`/ (e.g.&nbsp;http://spf13.com/topics/golang/) from:

* /layouts/taxonomy/`SINGULAR`.html (e.g.&nbsp;`/layouts/taxonomy/topic.html`)
* /layouts/\_default/taxonomy.html
* /layouts/\_default/list.html
* /themes/`THEME`/layouts/taxonomy/`SINGULAR`.html
* /themes/`THEME`/layouts/\_default/taxonomy.html
* /themes/`THEME`/layouts/\_default/list.html

Note that a taxonomy list page can also have a content file with front matter,  see [Source Organization](/overview/source-directory/).

### Taxonomy RSS

A Taxonomy’s RSS will be rendered at /`PLURAL`/`TERM`/index.xml (e.g.&nbsp;http://spf13.com/topics/golang/index.xml)

*Hugo ships with its own [RSS 2.0][] template. In most cases this will
be sufficient, and an RSS template will not need to be provided by the
user.*

Hugo provides the ability for you to define any RSS type you wish, and
can have different RSS files for each section and taxonomy.

* /layouts/taxonomy/`SINGULAR`.rss.xml
* /layouts/\_default/rss.xml
* /themes/`THEME`/layouts/taxonomy/`SINGULAR`.rss.xml
* /themes/`THEME`/layouts/\_default/rss.xml

## Variables

A list page is a `Page` and have all the [page variables](/templates/variables/)
and [site variables](/templates/variables/) available to use in the templates.

Taxonomy pages will additionally have:

**.Data.`Singular`** The taxonomy itself.<br>



{{% note %}}
If `where` or `first` receives invalid input or a field name that doesn’t exist, it will return an error and stop site generation. `where` and `first` also work on taxonomy list templates *and* taxonomy terms templates. (See [Taxonomy Templates](/templates/taxonomy-templates/).)
{{% /note %}}


[directorystructure]: /getting-started/directory-structure/
[homepage]: /templates/homepage-template/
[homepage]: /templates/homepage-template/
[limitkeyword]: https://www.techonthenet.com/sql/select_limit.php
[mentalmodel]: http://webstyleguide.com/wsg3/3-information-architecture/3-site-structure.html
[partials]: /templates/partials/
[RSS 2.0]: http://cyber.law.harvard.edu/rss/rss.html "RSS 2.0 Specification"
[RSS]: /templates/rss-templates/
[sections]: /content-management/sections/
[section templates]: /templates/section-templates
[taxonomy lists]: /templates/taxonomy-templates/#taxonomy-list-templates/
[taxonomy terms]: /templates/taxonomy-templates/#taxonomy-terms-templates/
