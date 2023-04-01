---
title: shuffle
description: Returns a random permutation of a given array or slice.
keywords: [ordering]
categories: [functions]
menu:
  docs:
    parent: functions
signature: ["shuffle COLLECTION"]
relatedfuncs: [seq]
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
