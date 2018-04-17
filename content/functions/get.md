---
title: .Get
description: Accesses positional and ordered parameters in shortcode declaration.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [shortcodes]
signature: [".Get INDEX", ".Get KEY"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
needsexample: true
---


`.Get` is specifically used when creating your own [shortcode template][sc], to access the [positional and named](/templates/shortcode-templates/#positional-vs-named-parameters) parameters passed to it. When used with a numeric INDEX, it queries positional parameters (starting with 0). With a string KEY, it queries named parameters.

When accessing a named parameter that does not exist, `.Get` returns an empty string instead of interrupting the build. The same goes with positional parameters in hugo version 0.40 and after. This allows you to chain `.Get` with `if`, `with`, `default` or `cond` to check for parameter existence. For example, you may now use:

```
{{ $quality := default "100" (.Get 1) }}
```

[sc]: /templates/shortcode-templates/




