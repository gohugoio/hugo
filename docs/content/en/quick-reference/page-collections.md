---
title: Page collections
description: A quick reference guide to Hugo's page collections.
categories: [quick reference]
keywords: []
menu:
  docs:
    parent: quick-reference
    weight: 50
weight: 50
toc: true
---

## Page

Use these `Page` methods when rendering lists on [section] pages, [taxonomy] pages, [term] pages, and the home page.

[section]: /getting-started/glossary/#section
[taxonomy]: /getting-started/glossary/#taxonomy
[term]: /getting-started/glossary/#term

{{< list-pages-in-section path=/methods/page filter=methods_page_page_collections filterType=include omitElementIDs=true titlePrefix=PAGE. >}}

## Site

Use these `Site` methods when rendering lists on any page.

{{< list-pages-in-section path=/methods/site filter=methods_site_page_collections filterType=include omitElementIDs=true titlePrefix=SITE. >}}

## Filter

Use the [`where`] function to filter page collections.

[`where`]: /functions/collections/where/

## Sort

By default, Hugo sorts page collections by:

1. [Weight]
2. [Date] in descending order
3. [LinkTitle] falling back to [Title]
4. [Filename] if the page is backed by a file

[Date]: /methods/page/date/
[Weight]: /methods/page/weight/
[LinkTitle]: /methods/page/linktitle/
[Title]: /methods/page/title/
[Filename]: /methods/page/file/#filename

Use these methods to sort page collections.

{{< list-pages-in-section path=/methods/pages filter=methods_pages_sort filterType=include titlePrefix=. omitElementIDs=true titlePrefix=PAGES. >}}

## Group

Use these methods to group page collections.

{{< list-pages-in-section path=/methods/pages filter=methods_pages_group filterType=include titlePrefix=. omitElementIDs=true titlePrefix=PAGES. >}}
