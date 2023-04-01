---
title: delimit
description: Loops through any array, slice, or map and returns a string of all the values separated by a delimiter.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [iteration]
signature: ["delimit COLLECTION DELIMIT LAST"]
relatedfuncs: []
---

`delimit` called in your template takes the form of

```go-html-template
{{ delimit array/slice/map delimiter optionallastdelimiter }}
```

`delimit` loops through any array, slice, or map and returns a string of all the values separated by a delimiter, the second argument in the function call. There is an optional third parameter that lets you choose a different delimiter to go between the last two values in the loop.

To maintain a consistent output order, maps will be sorted by keys and only a slice of the values will be returned.

The examples of `delimit` that follow all use the same front matter:

{{< code-toggle file="content/about.md" copy=false fm=true >}}
title: About
tags: [ "tag1", "tag2", "tag3" ]
{{< /code-toggle >}}

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
