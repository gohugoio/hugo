---
title: apply
description: Given a map, array, or slice, `apply` returns a new slice with a function applied over it.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [advanced]
signature: ["apply COLLECTION FUNCTION [PARAM...]"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
---

{{< todo >}}
POTENTIAL NEW CONTENT: see apply/sequence discussion: https://discourse.gohugo.io/t/apply-printf-on-a-sequence/5722;
{{< /todo >}}

`apply` expects at least three parameters, depending on the function being applied.

1. The first parameter is the sequence to operate on.
2. The second parameter is the name of the function as a string, which must be the name of a valid [Hugo function][functions].
3. After that, the parameters to the applied function are provided, with the string `"."` standing in for each element of the sequence the function is to be applied against.

Here is an example of a content file with `names:` as a front matter field:

```
+++
names: [ "Derek Perkins", "Joe Bergevin", "Tanner Linsley" ]
+++
```

You can then use `apply` as follows:

```
{{ apply .Params.names "urlize" "." }}
```

Which will result in the following:

```
"derek-perkins", "joe-bergevin", "tanner-linsley"
```

This is *roughly* equivalent to using the following with [range][]:

```
{{ range .Params.names }}{{ . | urlize }}{{ end }}
```

However, it is not possible to provide the output of a range to the [`delimit` function][delimit], so you need to `apply` it.

If you have `post-tag-list.html` and `post-tag-link.html` as [partials][], you *could* use the following snippets, respectively:

{{< code file="layouts/partials/post-tag-list.html" copy="false" >}}
{{ with .Params.tags }}
<div class="tags-list">
  Tags:
  {{ $len := len . }}
  {{ if eq $len 1 }}
    {{ partial "post-tag-link" (index . 0) }}
  {{ else }}
    {{ $last := sub $len 1 }}
    {{ range first $last . }}
      {{ partial "post-tag-link" . }},
    {{ end }}
    {{ partial "post-tag-link" (index . $last) }}
  {{ end }}
</div>
{{ end }}
{{< /code >}}

{{< code file="layouts/partials/post-tag-link.html" copy="false" >}}
<a class="post-tag post-tag-{{ . | urlize }}" href="/tags/{{ . | urlize }}">{{ . }}</a>
{{< /code >}}

This works, but the complexity of `post-tag-list.html` is fairly high. The Hugo template needs to perform special behavior for the case where there’s only one tag, and it has to treat the last tag as special. Additionally, the tag list will be rendered something like `Tags: tag1 , tag2 , tag3` because of the way that the HTML is generated and then interpreted by a browser.

This first version of `layouts/partials/post-tag-list.html` separates all of the operations for ease of reading. The combined and DRYer version is shown next:

```
{{ with .Params.tags }}
    <div class="tags-list">
      Tags:
      {{ $sort := sort . }}
      {{ $links := apply $sort "partial" "post-tag-link" "." }}
      {{ $clean := apply $links "chomp" "." }}
      {{ delimit $clean ", " }}
    </div>
{{ end }}
```

Now in the completed version, you can sort the tags, convert the tags to links with `layouts/partials/post-tag-link.html`, [chomp][] off stray newlines, and join the tags together in a delimited list for presentation. Here is an even DRYer version of the preceding example:

{{< code file="layouts/partials/post-tag-list.html" download="post-tag-list.html" >}}
    {{ with .Params.tags }}
    <div class="tags-list">
      Tags:
      {{ delimit (apply (apply (sort .) "partial" "post-tag-link" ".") "chomp" ".") ", " }}
    </div>
    {{ end }}
{{< /code >}}

{{% note %}}
`apply` does not work when receiving the sequence as an argument through a pipeline.
{{% /note %}}

[chomp]: /functions/chomp/ "See documentation for the chomp function"
[delimit]: /functions/delimit/ "See documentation for the delimit function"
[functions]: /functions/ "See the full list of Hugo functions to see what can be passed as an argument to the apply function."
[partials]: /templates/partials/
[range]: /functions/range/ "Learn the importance of the range function, a fundamental keyword in both Hugo templates and the Go programming language."
