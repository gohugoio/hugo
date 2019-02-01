---
title: len
linktitle: len
description: Returns the length of a variable according to its type.
godocref: https://golang.org/pkg/builtin/#len
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-04-18
categories: [functions]
keywords: []
signature: ["len INPUT"]
workson: [lists,taxonomies,terms]
hugoversion:
relatedfuncs: []
deprecated: false
toc: false
aliases: []
---

`len` is a built-in function in Go that returns the length of a variable according to its type. From the Go documentation:

> Array: the number of elements in v.
>
> Pointer to array: the number of elements in *v (even if v is nil).
>
> Slice, or map: the number of elements in v; if v is nil, len(v) is zero.
>
> String: the number of bytes in v.
>
> Channel: the number of elements queued (unread) in the channel buffer; if v is nil, len(v) is zero.

`len` is also considered a [fundamental function for Hugo templating][].

## `len` Example 1: Longer Headings

You may want to append a class to a heading according to the length of the string therein. The following templating checks to see if the title's length is greater than 80 characters and, if so, adds a `long-title` class to the `<h1>`:

{{< code file="check-title-length.html" >}}
<header>
    <h1{{if gt (len .Title) 80}} class="long-title"{{end}}>{{.Title}}</h1>
</header>
{{< /code >}}

## `len` Example 2: Counting Pages with `where`

The following templating uses [`where`][] in conjunction with `len` to
figure out the total number of content pages in a `posts` [section][]:

{{< code file="how-many-posts.html" >}}
{{ $posts := (where .Site.RegularPages "Section" "==" "posts") }}
{{ $postCount := len $posts }}
{{< /code >}}

Note the use of `.RegularPages`, a [site variable][] that counts all regular content pages but not the `_index.md` pages used to add front matter and content to [list templates][].


[fundamental function for Hugo templating]: /templates/introduction/
[list templates]: /templates/lists/
[section]: /content-management/sections/
[site variable]: /variables/site/
[`where`]: /functions/where/
