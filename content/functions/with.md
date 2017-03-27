---
title: with
linktitle: with
description: Rebinds the context (`.`) within its scope and skips the block if the variable is absent.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-03-12
categories: [functions]
tags: [conditionals,fundamentals]
ns:
signature: ["with INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
---

An alternative way of writing the "`if`" and then referencing the same value is to use `with` instead. `with` rebinds the context (`.`) within its scope and skips the block if the variable is absent:

{{% code file="layouts/partials/twitter.html" %}}
```html
{{with .Site.Params.TwitterUser}}<span class="twitter">
<a href="https://twitter.com/{{.}}" rel="author">
<img src="/images/twitter.png" width="48" height="48" title="Twitter: {{.}}"
 alt="Twitter"></a>
</span>{{end}}
```
{{% /code %}}
