---
title: Page variables
description: Use these methods with a Page object.
categories: [variables]
keywords: [page]
menu:
  docs:
    parent: variables
    weight: 50
weight: 50
toc: true
---

{{% include "variables/_common/consistent-terminology.md" %}}

## All methods

Use any of these methods in your templates.

{{< list-pages-in-section path=/methods/page titlePrefix=. >}}

## Dates

Use these methods to access content dates.

{{< list-pages-in-section path=/methods/page filter=methods_page_dates filterType=include titlePrefix=. omitElementIDs=true >}}

## Multilingual

Use these methods with your multilingual projects.

{{< list-pages-in-section path=/methods/page filter=methods_page_multilinugual filterType=include titlePrefix=. omitElementIDs=true >}}

## Navigation

Use these methods to create navigation links between pages.

{{< list-pages-in-section path=/methods/page filter=methods_page_navigation filterType=include titlePrefix=. omitElementIDs=true >}}

## Page collections

Range through these collections when rendering lists on [section] pages, [taxonomy] pages, [term] pages, and the home page.

[section]: /getting-started/glossary/#section
[taxonomy]: /getting-started/glossary/#taxonomy
[term]: /getting-started/glossary/#term
[context]: /getting-started/glossary/#context

{{< list-pages-in-section path=/methods/page filter=methods_page_page_collections filterType=include titlePrefix=. omitElementIDs=true >}}

## Parameters

Use these methods to access page parameters.

{{< list-pages-in-section path=/methods/page filter=methods_page_parameters filterType=include titlePrefix=. omitElementIDs=true >}}

## Sections

Use these methods to access section pages, and their ancestors and descendants. See&nbsp;[details].

[details]: /content-management/sections/

{{< list-pages-in-section path=/methods/page filter=methods_page_sections filterType=include titlePrefix=. omitElementIDs=true >}}
