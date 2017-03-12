---
title: Introduction to Hugo Lists
linktitle: Hugo Lists Introduction
description: Lists have a specific meaning and usage in Hugo when it comes to rendering your site homepage, section page, taxonomy list, or taxonomy terms list.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [lists,sections,rss,taxonomies,terms]
weight: 22
draft: false
aliases: [/templates/list/,/layout/indexes/]
toc: true
wip: true
---

## What is a List Page Template?

A list page template is a template used to render multiple pieces of content in a single HTML page. The exception to this rule is the homepage, which is still a list but has its own [dedicated template][homepage]).

Hugo uses the term *list* in its truest sense; i.e. a sequential arrangement of material, especially in alphabetical or numerical order. Hugo uses list templates on any output HTML page where content is traditionally listed:

* [Taxonomy terms pages][taxterms]
* [Taxonomy list pages][taxlists]
* [Section list pages][sectiontemps]
* [RSS][rss]

The idea of a list page comes from the [hierarchical mental model of the web][mentalmodel] and is best demonstrated visually:

![Image demonstrating a hierarchical website sitemap.](/images/site-hierarchy.svg)

## List Defaults

### Default Templates

Since section lists and taxonomy lists (N.B., *not* [taxonomy terms lists][taxterms]) are both *lists* with regards to their templates, both have the same terminating default of `_default/list.html`---or `themes/<THEME>/layouts/_default/list.html` in the case of a themed project---in their *lookup orders*. In addition, both [section lists][sectiontemps] and [taxonomy lists][taxlists] have their own default list templates in `_default`:

#### Default Section Templates

1. `layouts/section/<SECTIONNAME>.html`
2. `layouts/section/list.html`
3. `layouts/_default/section.html`
4. `layouts/_default/list.html`


### Taxonomy RSS

A Taxonomy’s RSS will be rendered at `/<PLURAL>/<TERM>/index.xml` (e.g.&nbsp;http://spf13.com/topics/golang/index.xml).

{{% note %}}
Most use cases will find that the [RSS 2.0][] template that ships with Hugo is sufficient for their needs.
{{% /note %}}

Hugo provides the ability for you to define any RSS type you wish. You can can have different RSS files for each section and taxonomy:

1. `/layouts/taxonomy/<SINGULAR>.rss.xml`
1. `/layouts/_default/rss.xml`
1. `/themes/<THEME>/layouts/taxonomy/<SINGULAR>.rss.xml`
1. `/themes/<THEME>/layouts/_default/rss.xml`

## List Page Variables

A list page is a `Page` and has all the [page variables][pagevars]
and [site variables][sitevars] available for use in templates.

Taxonomy list pages will additionally have:

**.Data.`Singular`** The taxonomy itself.<br> [See Taxonomy Variables][taxvars]

{{% note %}}
If `where` or `first` receives invalid input or a field name that doesn’t exist, it will return an error and stop site generation. `where` and `first` also work on taxonomy list templates *and* taxonomy terms templates. (See [Taxonomy Templates](/templates/taxonomy-templates/).)
{{% /note %}}

[directorystructure]: /getting-started/directory-structure/
[homepage]: /templates/homepage/
[homepage]: /templates/homepage/
[limitkeyword]: https://www.techonthenet.com/sql/select_limit.php
[mentalmodel]: http://webstyleguide.com/wsg3/3-information-architecture/3-site-structure.html
[pagevars]: /variables/pagevars/
[partials]: /templates/partials/
[RSS 2.0]: http://cyber.law.harvard.edu/rss/rss.html "RSS 2.0 Specification"
[rss]: /templates/rss/
[sections]: /content-management/sections/
[sectiontemps]: /templates/section-templates
[sitevars]: /variables/site/
[taxlists]: /templates/taxonomy-templates/#taxonomy-list-templates/
[taxvars]: /templates/taxonomy-variables/
[taxterms]: /templates/taxonomy-templates/#taxonomy-terms-templates/
