---
title: Pages variables
description: Use these methods with a collection of Page objects.
categories: [variables]
keywords: [pages]
menu:
  docs:
    parent: variables
    weight: 60
weight: 60
aliases: [/variables/site-variables/]
toc: true
---

{{% include "variables/_common/consistent-terminology.md" %}}

## All methods

Use any of these methods with page collections in your templates.

{{< list-pages-in-section path=/methods/pages titlePrefix=. >}}

## Sort by

Use these methods to sort page collections.

{{< list-pages-in-section path=/methods/pages filter=methods_pages_sort filterType=include titlePrefix=. omitElementIDs=true >}}

## Group by

Use these methods to group page collections.

{{< list-pages-in-section path=/methods/pages filter=methods_pages_group filterType=include titlePrefix=. omitElementIDs=true >}}

## Navigation

Use these methods to create navigation links between pages.

{{< list-pages-in-section path=/methods/pages filter=methods_pages_navigation filterType=include titlePrefix=. omitElementIDs=true >}}
