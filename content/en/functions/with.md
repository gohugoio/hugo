---
title: with
# linktitle: with
description: Rebinds the context (`.`) within its scope and skips the block if the variable is absent or empty.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-03-12
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [conditionals]
signature: ["with INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
---

An alternative way of writing an `if` statement and then referencing the same value is to use `with` instead. `with` rebinds the context (`.`) within its scope and skips the block if the variable is absent, unset or empty.

The set of *empty* values is defined by [the Go templates package](https://golang.org/pkg/text/template/). Empty values include `false`, the number zero, and the empty string.

If you want to render a block if an index or key is present in a slice, array, channel or map, regardless of whether the value is empty, you should use [`isset`](/functions/isset) instead.

The following example checks for a [user-defined site variable](/variables/site/) called `twitteruser`. If the key-value is not set, the following will render nothing:

{{< code file="layouts/partials/twitter.html" >}}
{{ with .Site.Params.twitteruser }}<span class="twitter">
<a href="https://twitter.com/{{.}}" rel="author">
<img src="/images/twitter.png" width="48" height="48" title="Twitter: {{.}}"
 alt="Twitter"></a>
</span>{{ end }}
{{< /code >}}
