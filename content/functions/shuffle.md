---
title: shuffle
linktitle:
description: Returns a random permutation of a given array or slice.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [ordering]
categories: [functions]
ns:
signature: ["shuffle COLLECTION"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
needsexamples: true
---

`shuffle` returns a random permutation of a given array or slice:

{{% code file="shuffle-input.html" %}}
```html
<div class="shuffle-sequence">{{ shuffle (seq 1 5) }}</div>
<div class="shuffle-slice">{{ shuffle (slice "foo" "bar" "buzz") }}</div>
```
{{% /code %}}

This example would return the following:

{{% output file="shuffle-output.html" %}}
```html
<div class="shuffle-seq">2 5 3 1 4</div>
<div class="shuffle-slice">buzz foo bar</div>
```
{{% /output %}}

This example also makes use of the [slice](/functions/slice/) and [seq](/functions/seq/) functions.
