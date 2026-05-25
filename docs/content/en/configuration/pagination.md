---
title: Configure pagination
linkTitle: Pagination
description: Configure pagination.
categories: []
keywords: []
---

This is the default configuration:

{{< code-toggle config=pagination />}}

disableAliases
: (`bool`) Whether to disable alias generation for the first pager. Default is `false`.

pagerSize
: (`int`) The number of pages per pager. Default is `10`.

path
: (`string`) The segment of each pager URL indicating that the target page is a pager. Default is `page`.

With multilingual projects you can define the pagination behavior for each language:

{{< code-toggle file=hugo >}}
[languages.en]
contentDir = 'content/en'
direction = 'ltr'
label = 'English'
locale = 'en-US'
weight = 1
[languages.en.pagination]
disableAliases = true
pagerSize = 10
path = 'page'
[languages.de]
contentDir = 'content/de'
direction = 'ltr'
label = 'Deutsch'
locale = 'de-DE'
weight = 2
[languages.de.pagination]
disableAliases = true
pagerSize = 20
path = 'blatt'
{{< /code-toggle >}}
