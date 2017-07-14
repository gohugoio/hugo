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
menu:
  docs:
    parent: "functions"
tags: [dictionary]
ns:
signature: ["dict KEY VALUE [KEY VALUE]..."]
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

`dict` also allows you to create a map on the fly to pass into your [partial templates][partials]:

{{% code file="dict-create-map.html" %}}
```golang
{{partial "foo" (dict "important" "Smiles" "content" "You should do more")}}
```
{{% /code %}}

## Example: `dict` with Embedded SVGs

Let's suppose you want to pass values into embedded SVG icons on the fly in your template without having to change your CSS. The following is the XML for an "external link" SVG icon:

{{% code file="layouts/partials/svgs/external-links.svg" copy="false" %}}
```xml
<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" fill="#000000" width="40" height="40" viewBox="0 0 32 32" aria-label="External Link">
<path d="M25.152 16.576v5.696q0 2.144-1.504 3.648t-3.648 1.504h-14.848q-2.144 0-3.648-1.504t-1.504-3.648v-14.848q0-2.112 1.504-3.616t3.648-1.536h12.576q0.224 0 0.384 0.16t0.16 0.416v1.152q0 0.256-0.16 0.416t-0.384 0.16h-12.576q-1.184 0-2.016 0.832t-0.864 2.016v14.848q0 1.184 0.864 2.016t2.016 0.864h14.848q1.184 0 2.016-0.864t0.832-2.016v-5.696q0-0.256 0.16-0.416t0.416-0.16h1.152q0.256 0 0.416 0.16t0.16 0.416zM32 1.152v9.12q0 0.48-0.352 0.8t-0.8 0.352-0.8-0.352l-3.136-3.136-11.648 11.648q-0.16 0.192-0.416 0.192t-0.384-0.192l-2.048-2.048q-0.192-0.16-0.192-0.384t0.192-0.416l11.648-11.648-3.136-3.136q-0.352-0.352-0.352-0.8t0.352-0.8 0.8-0.352h9.12q0.48 0 0.8 0.352t0.352 0.8z"></path>
</svg>
```
{{% /code %}}

You can then call this [partial template][partials] using `{{ partial "svgs/external-links.svg" . }}`, but what if you want to pass specific values such as `fill` `height` and `width` into the icon? To allow for more flexibility, you can abstract these values into variables that you will later define when calling the partial template:

{{% code file="layouts/partials/svgs/external-links.svg" download="external-links.svg" %}}
```xml
<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" fill="{{ .fill }}" width="{{ .size }}" height="{{ .size }}" viewBox="0 0 32 32" aria-label="External Link">
<path d="M25.152 16.576v5.696q0 2.144-1.504 3.648t-3.648 1.504h-14.848q-2.144 0-3.648-1.504t-1.504-3.648v-14.848q0-2.112 1.504-3.616t3.648-1.536h12.576q0.224 0 0.384 0.16t0.16 0.416v1.152q0 0.256-0.16 0.416t-0.384 0.16h-12.576q-1.184 0-2.016 0.832t-0.864 2.016v14.848q0 1.184 0.864 2.016t2.016 0.864h14.848q1.184 0 2.016-0.864t0.832-2.016v-5.696q0-0.256 0.16-0.416t0.416-0.16h1.152q0.256 0 0.416 0.16t0.16 0.416zM32 1.152v9.12q0 0.48-0.352 0.8t-0.8 0.352-0.8-0.352l-3.136-3.136-11.648 11.648q-0.16 0.192-0.416 0.192t-0.384-0.192l-2.048-2.048q-0.192-0.16-0.192-0.384t0.192-0.416l11.648-11.648-3.136-3.136q-0.352-0.352-0.352-0.8t0.352-0.8 0.8-0.352h9.12q0.48 0 0.8 0.352t0.352 0.8z"></path>
</svg>
```
{{% /code %}}

You can then pass these values to customize the SVG by calling the partials with `dict`:

{{% code file="layouts/_default/list.html" %}}
```html
...
{{ partial "svg/link-ext.svg" (dict "fill" "#01589B" "size" "10") }}
...
```
{{% /code %}}

The above partial will generate the following code when Hugo builds your site:

```html
<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" fill="#01589B" width="10" height="10" viewBox="0 0 32 32" aria-label="External Link">
<path d="M25.152 16.576v5.696q0 2.144-1.504 3.648t-3.648 1.504h-14.848q-2.144 0-3.648-1.504t-1.504-3.648v-14.848q0-2.112 1.504-3.616t3.648-1.536h12.576q0.224 0 0.384 0.16t0.16 0.416v1.152q0 0.256-0.16 0.416t-0.384 0.16h-12.576q-1.184 0-2.016 0.832t-0.864 2.016v14.848q0 1.184 0.864 2.016t2.016 0.864h14.848q1.184 0 2.016-0.864t0.832-2.016v-5.696q0-0.256 0.16-0.416t0.416-0.16h1.152q0.256 0 0.416 0.16t0.16 0.416zM32 1.152v9.12q0 0.48-0.352 0.8t-0.8 0.352-0.8-0.352l-3.136-3.136-11.648 11.648q-0.16 0.192-0.416 0.192t-0.384-0.192l-2.048-2.048q-0.192-0.16-0.192-0.384t0.192-0.416l11.648-11.648-3.136-3.136q-0.352-0.352-0.352-0.8t0.352-0.8 0.8-0.352h9.12q0.48 0 0.8 0.352t0.352 0.8z"></path>
</svg>
```


[partials]: /templates/partials/
