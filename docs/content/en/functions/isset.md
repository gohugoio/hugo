---
title: isset
linktitle: isset
description: Returns true if the parameter is set.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: []
signature: ["isset COLLECTION INDEX", "isset COLLECTION KEY"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

Takes either a slice, array, or channel and an index or a map and a key as input.

```
{{ if isset .Params "project_url" }} {{ index .Params "project_url" }}{{ end }}
```

{{% warning %}}
All site-level configuration keys are stored as lower case. Therefore, a `myParam` key-value set in your [site configuration file](/getting-started/configuration/) needs to be accessed with `{{if isset .Site.Params "myparam"}}` and *not* with `{{if isset .Site.Params "myParam"}}`. Note that you can still access the same config key with `.Site.Params.myParam` *or* `.Site.Params.myparam`, for example, when using [`with`](/functions/with).
This restriction also applies when accessing page-level front matter keys from within [shortcodes](/content-management/shortcodes/).
{{% /warning %}}

