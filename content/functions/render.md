---
title: render
linktitle: Render
description:
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [views]
categories: [functions]
toc:
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
---

Takes a view to render the content with.  The view is an alternate layout, and should be a file name that points to a template in one of the locations specified in the documentation for [Content Views](/templates/views).

This function is only available on a piece of content, and in list context.

This example could render a piece of content using the content view located at `/layouts/_default/summary.html`:

```golang
{{ range .Data.Pages }}
    {{ .Render "summary"}}
{{ end }}
```