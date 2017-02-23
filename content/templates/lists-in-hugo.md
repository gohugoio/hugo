---
title: Lists in Hugo
linktitle: Lists in Hugo
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [lists,sections,rss,taxonomies,terms]
weight: 25
draft: false
aliases: [/templates/list/]
toc: true
needsreview: true
---

## What is a "List" Template

A list page template is a template used to render multiple pieces of content in a single HTML page (with the exception of the homepage, which has a [dedicated template][homepage]).

Hugo uses the term *list* in its truest sense: a sequential arrangement of material, especially in alphabetical or numerical order. Hugo uses list templates on any output HTML page where content is being listed (e.g., [taxonomies][], [sections][], and [RSS][]). The idea of a list page comes from the [hierarchical mental model of the web][mentalmodel] and is best demonstrated visually:

![Image demonstrating a hierarchical website sitemap.](/images/site-hierarchy.svg)

## List Defaults

### Default Templates

Since section lists and taxonomy lists (N.B., *not* [taxonomy terms lists][]) are both *lists* with regards to their templates, both of these templates have the same terminating default of `_default/list.html`---or `themes/mytheme/layouts/_default/list.html` in the case of a themed project---in their *lookup orders*. In addition, both [section lists][sections] and [taxonomy lists][taxonomies] have their own default list templates in `_default`:

#### Default Section Templates

1. `layouts/section/sectionname.html`

### Understanding `.Data.Pages`

{{% note "The Confusion over `.Data`" %}}
**Mention something here about the difference between .Data.Pages and .Site.Data maybe?**
{{% /note %}}


[homepage]: /templates/homepage-template/
[mentalmodel]: http://webstyleguide.com/wsg3/3-information-architecture/3-site-structure.html
[RSS]: /templates/rss-templates/
[sections]: /templates/section-templates
[taxonomies]: /templates/taxonomy-templates/#taxonomy-list-templates/
[taxonomy terms lists]: /templates/#taxonomy-terms-templates/
