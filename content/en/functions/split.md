---
title: split
# linktitle: split
description: Returns an array of strings by splitting STRING by DELIM.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings]
signature: ["split STRING DELIM"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

* `{{ split "tag1,tag2,tag3" "," }}` → ["tag1", "tag2", "tag3"]

{{% note %}}
`split` essentially does the opposite of [delimit]({{< ref "functions/delimit" >}}). While `split` creates a slice from a string, `delimit` creates a string from a slice.
{{% /note %}}
