---
title: delimit
description: Loops through any array, slice, or map and returns a string of all the values separated by a delimiter.
godocref:
workson: []
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [iteration]
toc: false
signature: ["delimit COLLECTION DELIMIT LAST"]
workson: [lists,taxonomies,terms]
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
---

`delimit` called in your template takes the form of

```
{{ delimit array/slice/map delimiter optionallastdelimiter}}
```

`delimit` loops through any array, slice, or map and returns a string of all the values separated by a delimiter, the second argument in the function call. There is an optional third parameter that lets you choose a different delimiter to go between the last two values in the loop.

To maintain a consistent output order, maps will be sorted by keys and only a slice of the values will be returned.

The examples of `delimit` that follow all use the same front matter:

{{< code file="delimit-example-front-matter.toml" nocopy="true" >}}
+++
title: I love Delimit
tags: [ "tag1", "tag2", "tag3" ]
+++
{{< /code >}}

{{< code file="delimit-page-tags-input.html" >}}
<p>Tags: {{ delimit .Params.tags ", " }}</p>
{{< /code >}}

{{< output file="delimit-page-tags-output.html" >}}
<p>Tags: tag1, tag2, tag3</p>
{{< /output >}}

Here is the same example but with the optional "last" delimiter:

{{< code file="delimit-page-tags-final-and-input.html" >}}
Tags: {{ delimit .Params.tags ", " ", and " }}
{{< /code >}}

{{< output file="delimit-page-tags-final-and-output.html" >}}
<p>Tags: tag1, tag2, and tag3</p>
{{< /output >}}


[lists]: /templates/lists/
[taxonomies]: /templates/taxonomy-templates/#taxonomy-list-templates
[terms]: /templates/taxonomy-templates/#terms-list-templates
