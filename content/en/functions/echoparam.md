---
title: echoParam
description: Prints a parameter if it is set.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: collections
relatedFuncs: []
signature:
  - collections.EchoParam COLLECTION KEY
  - echoParam COLLECTION KEY
---

For example, consider this site configuration:

{{< code-toggle file=hugo copy=false >}}
[params.footer]
poweredBy = 'Hugo'
{{< /code-toggle >}}

To print the value of `poweredBy`:

```go-html-template
{{ echoParam site.Params.footer "poweredby" }} → Hugo
```

{{% note %}}
When using the `echoParam` function you must reference the key using lower case. See the previous example.

The `echoParam` function will be deprecated in a future release. Instead, use either of the constructs below.
{{% /note %}}

```go-html-template
{{ site.Params.footer.poweredBy }} → Hugo
{{ index site.Params.footer "poweredBy" }} → Hugo
```
