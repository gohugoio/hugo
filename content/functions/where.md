---
title: where
linktitle: where
description: Filters an array to only the elements containing a matching value for a given field.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [filtering]
signature:
workson: [lists,taxonomies,terms,groups]
hugoversion:
relatedfuncs: [intersect,first]
deprecated: false
toc: true
needsexample: true
---

`where` filters an array to only the elements containing a matching value for a given field.

```golang
{{ range where .Data.Pages "Section" "post" }}
  {{ .Content }}
{{ end }}
```

It can be used by dot-chaining the second argument to refer to a nested element of a value.

```toml
+++
series: golang
+++
```

```golang
{{ range where .Site.Pages "Params.series" "golang" }}
   {{ .Content }}
{{ end }}
```

It can also be used with the logical operators `!=`, `>=`, `in`, etc. Without an operator, `where` compares a given field with a matching value equivalent to `=`.

```golang
{{ range where .Data.Pages "Section" "!=" "post" }}
   {{ .Content }}
{{ end }}
```

The following logical operators are vailable with `where`:

* `=`, `==`, `eq`: True if a given field value equals a matching value
* `!=`, `<>`, `ne`: True if a given field value doesn't equal a matching value
* `>=`, `ge`: True if a given field value is greater than or equal to a matching value
* `>`, `gt`: True if a given field value is greater than a matching value
* `<=`, `le`: True if a given field value is lesser than or equal to a matching value
* `<`, `lt`: True if a given field value is lesser than a matching value
* `in`: True if a given field value is included in a matching value. A matching value must be an array or a slice
* `not in`: True if a given field value isn't included in a matching value. A matching value must be an array or a slice
* `intersect`: True if a given field value that is a slice / array of strings or integers contains elements in common with the matching value. It follows the same rules as the intersect function.

## Using `where` with `intersect`

```golang
{{ range where .Site.Pages ".Params.tags" "intersect" .Params.tags }}
  {{ if ne .Permalink $.Permalink }}
    {{ .Render "summary" }}
  {{ end }}
{{ end }}
```

## Using `where` with `first`

```golang
{{ range first 5 (where .Data.Pages "Section" "post") }}
   {{ .Content }}
{{ end }}
```

## Nesting `where` Clauses

**Needs Example**

## Unset Fields

Filtering only works for set fields. To check whether a field is set or exists, you can use the operand `nil`.

This can be useful to filter a small amount of pages from a large pool. Instead of set field on all pages, you can set field on required pages only.

Only the following operators are available for `nil`

* `=`, `==`, `eq`: True if the given field is not set.
* `!=`, `<>`, `ne`: True if the given field is set.

```golang
{{ range where .Data.Pages ".Params.specialpost" "!=" nil }}
   {{ .Content }}
{{ end }}
```

