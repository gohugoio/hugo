---
title: delimit
linktitle: delimit
description: loops through any array, slice, or map and returns a string of all the values separated by a delimiter.
qref: d
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [iteration]
toc: false
draft: false
aliases: []
---

`delimit` loops through any array, slice, or map and returns a string of all the values separated by the delimiter. There is an optional third parameter that lets you choose a different delimiter to go between the last two values.

Maps will be sorted by the keys, and only a slice of the values will be returned, keeping a consistent output order.

`delimit` works on [lists][], [taxonomies][], and [terms][].

Examples of `delimit` use the following front matter:

{{% input "delimit-example-front-matter.toml" nocopy %}}
```toml
+++
title: I love Delimit
tags: [ "tag1", "tag2", "tag3" ]
+++
```
{{% /input %}}

`delimit` called in your template takes the form of

```
{{ delimit array/slice/map delimiter optionallastdelimiter}}
```

{{% input "delimit-pages-tags.html" %}}
```html
<p>Tags: {{ delimit .Params.tags ", " }}</p>
```
{{% /input %}}

{{% output "delimit-pages-tags-output.html" %}}
```html
<p>Tags: tag1, tag2, tag3</p>
```
{{% /output %}}

Here is the same example but with the optional "last" delimiter:

{{% input "delimit-page-tags-final-and.html" %}}
```golang
Tags: {{ delimit .Params.tags ", " ", and " }}
```
{{% /input %}}

{{% output "delimit-page-tags-final-and-output.html" %}}
```html
<p>Tags: tag1, tag2, and tag3</p>
```
{{% /output %}}


[lists]: /templates/lists-in-hugo/
[taxonomies]: /templates/taxonomy-templates/#taxonomy-list-templates
[terms]: /templates/taxonomy-templates/#terms-list-templates