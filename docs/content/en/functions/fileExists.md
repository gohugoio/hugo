---
title: "fileExists"
linktitle: "fileExists"
date: 2017-08-31T22:38:22+02:00
description: Checks whether a file exists under the given path.
godocref:
publishdate: 2017-08-31T22:38:22+02:00
lastmod: 2017-08-31T22:38:22+02:00
categories: [functions]
menu:
  docs:
    parent: "functions"
signature: ["fileExists PATH"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`fileExists` allows you to check if a file exists under a given path, e.g. before inserting code into a template:

```
{{ if (fileExists "static/img/banner.jpg") -}}
<img src="{{ "img/banner.jpg" | absURL }}" />
{{- end }}
```

In the example above, a banner from the `static` folder should be shown if the given path points to an existing file.