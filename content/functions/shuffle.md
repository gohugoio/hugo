---
title: shuffle
# linktitle:
description: Returns a random permutation of a given array or slice.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-04-30
keywords: [ordering]
categories: [functions]
menu:
  docs:
    parent: "functions"
signature: ["shuffle COLLECTION"]
workson: []
hugoversion:
relatedfuncs: [seq]
deprecated: false
draft: false
aliases: []
---

{{< code file="shuffle-input.html" >}}
<!-- Shuffled sequence = -->
<div>{{ shuffle (seq 1 5) }}</div>
<!-- Shuffled slice =  -->
<div>{{ shuffle (slice "foo" "bar" "buzz") }}</div>
{{< /code >}}

This example would return the following:

{{< output file="shuffle-output.html" >}}
<!-- Shuffled sequence =  -->
<div>2 5 3 1 4</div>
<!-- Shuffled slice =  -->
<div>buzz foo bar</div>
{{< /output >}}

This example also makes use of the [slice](/functions/slice/) and [seq](/functions/seq/) functions.
