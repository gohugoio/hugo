---
title: Hugo's Lookup Order
linktitle: Template Lookup Order
description: Hugo searches for the layout to use for a given page in a well defined order, starting from the most specific.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-07-05
categories: [templates,fundamentals]
keywords: [templates]
menu:
  docs:
    parent: "templates"
    weight: 15
  quicklinks:
weight: 15
sections_weight: 15
draft: false
aliases: [/templates/lookup/]
toc: true
---

## Hugo Layouts Lookup Rules

Hugo takes the parameters listed below into consideration when choosing a layout for a given page. They are listed in a priority order. This should feel natural, but look at the table below for concrete examples of the different parameter variations.


Kind
: The page `Kind` (the home page is one). See the example tables below per kind. This also determines if it is a **single page** (i.e. a regular content page. We then look for a template in `_default/single.html` for HTML) or a **list page** (section listings, home page, taxonomy lists, taxonomy terms. We then look for a template in `_default/list.html` for HTML).

Layout
: Can be set in page front matter.

Output Format
: See [Custom Output Formats](/templates/output-formats). An output format has both a `name` (e.g. `rss`, `amp`, `html`) and a `suffix` (e.g. `xml`, `html`). We prefer matches with both (e.g. `index.amp.html`, but look for less specific templates.

Language
: We will consider a language code in the template name. If the site language is `fr`, `index.fr.amp.html` will win over `index.amp.html`, but `index.amp.html` will be chosen before `index.fr.html`.

Type
: Is value of `type` if set in front matter, else it is the name of the root section (e.g. "blog"). It will always have a value, so if not set, the value is "page". 

Section
: Is relevant for `section`, `taxonomy` and `taxonomyTerm` types.

{{% note %}}
**Tip:** The examples below looks long and complex. That is the flexibility talking. Most Hugo sites contain just a handful of templates:

```bash
├── _default
│   ├── baseof.html
│   ├── list.html
│   └── single.html
└── index.html
```
{{% /note %}}


## Hugo Layouts Lookup Rules With Theme

In Hugo, layouts can live in either the project's or the themes' layout folders, and the most specific layout will be chosen. Hugo will interleave the lookups listed below, finding the most specific one either in the project or themes.

## Examples: Layout Lookup for Regular Pages

{{< datatable-filtered "output" "layouts" "Kind == page" "Example" "OutputFormat" "Suffix" "Template Lookup Order" >}}

## Examples: Layout Lookup for Home Page

{{< datatable-filtered "output" "layouts" "Kind == home" "Example" "OutputFormat" "Suffix" "Template Lookup Order" >}}

## Examples: Layout Lookup for Section Pages

{{< datatable-filtered "output" "layouts" "Kind == section" "Example" "OutputFormat" "Suffix" "Template Lookup Order" >}}

## Examples: Layout Lookup for Taxonomy List Pages

{{< datatable-filtered "output" "layouts" "Kind == taxonomy" "Example" "OutputFormat" "Suffix" "Template Lookup Order" >}}

## Examples: Layout Lookup for Taxonomy Terms Pages

{{< datatable-filtered "output" "layouts" "Kind == taxonomyTerm" "Example" "OutputFormat" "Suffix" "Template Lookup Order" >}}




