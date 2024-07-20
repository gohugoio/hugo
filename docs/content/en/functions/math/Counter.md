---
title: math.Counter
description: Increments and returns a global counter.
categories: []
keywords: []
action:
  aliases: []
  related: []
  returnType: uint64
  signatures: [math.Counter]
---

The counter is global for both monolingual and multilingual sites, and its initial value for each build is&nbsp;1.

```go-html-template
{{ warnf "single.html called %d times" math.Counter }}
```

```sh
WARN  single.html called 1 times
WARN  single.html called 2 times
WARN  single.html called 3 times
```

Use this function to:

- Create unique warnings as shown above; the [`warnf`] function suppresses duplicate messages
- Create unique target paths for the `resources.FromString` function where the target path is also the cache key

[`warnf`]: /functions/fmt/warnf/
[`resources.FromString`]: /functions/resources/fromstring/

{{% note %}}
Due to concurrency, the value returned in a given template for a given page will vary from one build to the next. You cannot use this function to assign a static id to each page.
{{% /note %}}
