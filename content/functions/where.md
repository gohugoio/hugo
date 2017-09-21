---
title: where
# linktitle: where
description: Filters an array to only the elements containing a matching value for a given field.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [filtering]
signature: ["where COLLECTION KEY [OPERATOR] MATCH"]
workson: [lists,taxonomies,terms,groups]
hugoversion:
relatedfuncs: [intersect,first,after,last]
deprecated: false
toc: true
needsexample: true
---

`where` filters an array to only the elements containing a matching value for a given field.

```
{{ range where .Data.Pages "Section" "post" }}
  {{ .Content }}
{{ end }}
```

It can be used by dot-chaining the second argument to refer to a nested element of a value.

```
+++
series: golang
+++
```

```
{{ range where .Site.Pages "Params.series" "golang" }}
   {{ .Content }}
{{ end }}
```

It can also be used with the logical operators `!=`, `>=`, `in`, etc. Without an operator, `where` compares a given field with a matching value equivalent to `=`.

```
{{ range where .Data.Pages "Section" "!=" "post" }}
   {{ .Content }}
{{ end }}
```

The following logical operators are vailable with `where`:

`=`, `==`, `eq`
: `true` if a given field value equals a matching value

`!=`, `<>`, `ne`
: `true` if a given field value doesn't equal a matching value

`>=`, `ge`
: `true` if a given field value is greater than or equal to a matching value

`>`, `gt`
: `true` if a given field value is greater than a matching value

`<=`, `le`
: `true` if a given field value is lesser than or equal to a matching value

`<`, `lt`
: `true` if a given field value is lesser than a matching value

`in`
: `true` if a given field value is included in a matching value; a matching value must be an array or a slice

`not in`
: `true` if a given field value isn't included in a matching value; a matching value must be an array or a slice

`intersect`
: `true` if a given field value that is a slice/array of strings or integers contains elements in common with the matching value; it follows the same rules as the [`intersect` function][intersect].

## Use `where` with `intersect`

```
{{ range where .Site.Pages ".Params.tags" "intersect" .Params.tags }}
  {{ if ne .Permalink $.Permalink }}
    {{ .Render "summary" }}
  {{ end }}
{{ end }}
```

You can also put the returned value of the `where` clauses into a variable:

{{< code file="where-intersect-variables.html" >}}
{{ $v1 := where .Site.Pages "Params.a" "v1" }}
{{ $v2 := where .Site.Pages "Params.b" "v2" }}
{{ $filtered := $v1 | intersect $v2 }}
{{ range $filtered }}
{{ end }}
{{< /code >}}

## Use `where` with `first`

The following grabs the first five content files in `post` using the [default ordering](/templates/lists/) for lists (i.e., `weight => date`):

{{< code file="where-with-first.html" >}}
{{ range first 5 (where .Data.Pages "Section" "post") }}
   {{ .Content }}
{{ end }}
{{< /code >}}

## Nest `where` Clauses

You can also nest `where` clauses to drill down on lists of content by more than one parameter. The following first grabs all pages in the "blog" section and then ranges through the result of the first `where` clause and finds all pages that are *not* featured:

```
{{ range where (where .Data.Pages "Section" "blog" ) ".Params.featured" "!=" "true" }}
```

## Unset Fields

Filtering only works for set fields. To check whether a field is set or exists, you can use the operand `nil`.

This can be useful to filter a small amount of pages from a large pool. Instead of set field on all pages, you can set field on required pages only.

Only the following operators are available for `nil`

* `=`, `==`, `eq`: True if the given field is not set.
* `!=`, `<>`, `ne`: True if the given field is set.

```
{{ range where .Data.Pages ".Params.specialpost" "!=" nil }}
   {{ .Content }}
{{ end }}
```

[intersect]: /functions/intersect/
