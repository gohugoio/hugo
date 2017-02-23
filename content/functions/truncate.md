---
title: truncate
linktitle: truncate
description:
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: []
toc: false
draft: false
aliases: []
---

Truncate a text to a max length without cutting words or leaving unclosed HTML tags. Since Go templates are HTML-aware, truncate will handle normal strings vs HTML strings intelligently. It's important to note that if you have a raw string that contains HTML tags that you want treated as HTML, you will need to convert the string to HTML using the safeHTML template function before sending the value to truncate; otherwise, the HTML tags will be escaped by truncate.

`{{ "<em>Keep my HTML</em>" | safeHTML | truncate 10 }}` → `<em>Keep my …</em>`