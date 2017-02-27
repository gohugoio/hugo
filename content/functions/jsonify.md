---
title: jsonify
linktitle: jsonify
description: Encodes a given object to JSON.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: []
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`jsonify` encodes a given object to JSON.

```
{{ dict "title" .Title "content" .Plain | jsonify }}
```

