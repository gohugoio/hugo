---
title: getenv
description: Returns the value of an environment variable.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: []
signature: ["getenv VARIABLE"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

Takes a string containing the name of the variable as input. Returns
an empty string if the variable is not set, otherwise returns the
value of the variable. 

```
{{ getenv "HOME" }}
```

{{% note %}}
In Unix-like environments, the variable must also be exported in order to be seen by `hugo`.
{{% /note %}}
