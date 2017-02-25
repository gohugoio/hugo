---
title: dict
linktitle:
description:
godocref:
workson: []
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: []
categories: [functions]
toc:
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
needsexamples: true
---

Creates a dictionary `(map[string, interface{})`, expects parameters added in value:object fasion.
Invalid combinations like keys that are not strings or uneven number of parameters, will result in an exception thrown.
Useful for passing maps to partials when adding to a template.

e.g. Pass into "foo.html" a map with the keys "important, content"

{{% code file="dict-example.html" %}}
```html
{{$important := .Site.Params.SomethingImportant }}
{{range .Site.Params.Bar}}
    {{partial "foo" (dict "content" . "important" $important)}}
{{end}}
```
{{% /code %}}

And then in `foo.html`:

```golang
Important {{.important}}
{{.content}}
```

`dict` also allows you to create a map on the fly to pass into

{{% code file="dict-create-map.html" %}}
```golang
{{partial "foo" (dict "important" "Smiles" "content" "You should do more")}}
```
{{% /code %}}


