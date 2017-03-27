---
title: dict
linktitle: dict
description: Creates a dictionary `(map[string, interface{})` that expects parameters added in a value:object fashion.
godocref:
workson: []
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-26
categories: [functions]
tags: [dictionary]
ns:
signature: []
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
needsexamples: true
---

`dict` creates a dictionary `(map[string, interface{})` that expects parameters added in a value:object fashion.

Invalid combinations---e.g., keys that are not strings or an uneven number of parameters---will result in an exception being thrown. `dict` is especially useful for passing maps to partials being added to a template.

For example, the following snippet passes a map with the keys "important, content" into "foo.html"

{{% code file="dict-example.html" %}}
```html
{{$important := .Site.Params.SomethingImportant }}
{{range .Site.Params.Bar}}
    {{partial "foo" (dict "content" . "important" $important)}}
{{end}}
```
{{% /code %}}

These keys can then be called in `foo.html` as follows:

```golang
Important {{.important}}
{{.content}}
```

`dict` also allows you to create a map on the fly to pass into your [partial templates][partials]

{{% code file="dict-create-map.html" %}}
```golang
{{partial "foo" (dict "important" "Smiles" "content" "You should do more")}}
```
{{% /code %}}

[partials]: /templates/partials/