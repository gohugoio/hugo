---
title: delimit
linktitle: delimit
description: "`delimit` loops through any array, slice, or map and returns a string of all the values separated by a delimiter."
godocref:
qref: loops through array, slice, or map and returns string of all values separated by a delimiter.
workson: []
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [iteration]
toc: false
signature: ""
workson: [lists,taxonomies,terms]
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
---

`delimit` loops through any array, slice, or map and returns a string of all the values separated by a delimiter, the second argument in the function call. There is an optional third parameter that lets you choose a different delimiter to go between the last two values in the loop.

To maintain a consistent output order, maps will be sorted by keys and only a slice of the values will be returned.

Examples of `delimit` use the following front matter:

{{% code file="delimit-example-front-matter.toml" nocopy="true" %}}
```toml
+++
title: I love Delimit
tags: [ "tag1", "tag2", "tag3" ]
+++
```
{{% /code %}}

`delimit` called in your template takes the form of

```
{{ delimit array/slice/map delimiter optionallastdelimiter}}
```

{{% code file="delimit-pages-tags.html" %}}
```html
<p>Tags: {{ delimit .Params.tags ", " }}</p>
```
{{% /code %}}

{{% output "delimit-pages-tags-output.html" %}}
```html
<p>Tags: tag1, tag2, tag3</p>
```
{{% /output %}}

Here is the same example but with the optional "last" delimiter:

{{% code file="delimit-page-tags-final-and.html" %}}
```golang
Tags: {{ delimit .Params.tags ", " ", and " }}
```
{{% /code %}}

{{% output "delimit-page-tags-final-and-output.html" %}}
```html
<p>Tags: tag1, tag2, and tag3</p>
```
{{% /output %}}


[lists]: /templates/lists-in-hugo/
[taxonomies]: /templates/taxonomy-templates/#taxonomy-list-templates
[terms]: /templates/taxonomy-templates/#terms-list-templates