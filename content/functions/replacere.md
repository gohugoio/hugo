---
title: replacere
linktitle: replaceRE
description: Replaces all occurrences of a regular expression with the replacement pattern.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [regex]
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

Replaces all occurrences of a regular expression with the replacement pattern.

```golang
{{ replaceRE "^https?://([^/]+).*" "$1" "http://gohugo.io/docs" }}` → "gohugo.io"
{{ "http://gohugo.io/docs" | replaceRE "^https?://([^/]+).*" "$1" }}` → "gohugo.io"
```