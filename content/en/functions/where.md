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

`where` filters an array to only the elements containing a matching
value for a given field.

It works in a similar manner to the [`where` keyword in
SQL][wherekeyword].

```go-html-template
{{ range where .Pages "Section" "foo" }}
  {{ .Content }}
{{ end }}
```

It can be used by dot-chaining the second argument to refer to a nested element of a value.

```
+++
series: golang
+++
```

```go-html-template
{{ range where .Site.Pages "Params.series" "golang" }}
   {{ .Content }}
{{ end }}
```

It can also be used with the logical operators `!=`, `>=`, `in`, etc. Without an operator, `where` compares a given field with a matching value equivalent to `=`.

```go-html-template
{{ range where .Pages "Section" "!=" "foo" }}
   {{ .Content }}
{{ end }}
```

The following logical operators are available with `where`:

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

## Use `where` with `Booleans`
When using booleans you should not put quotation marks.
```go-html-template
{{range where .Pages ".Draft" true}}
        <p>{{.Title}}</p>
{{end}}
```
  

## Use `where` with `intersect`

```go-html-template
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

Using `first` and `where` together can be very
powerful. Below snippet gets a list of posts only from [**main
sections**](#mainsections), sorts it using the [default
ordering](/templates/lists/) for lists (i.e., `weight => date`), and
then ranges through only the first 5 posts in that list:

{{< code file="first-and-where-together.html" >}}
{{ range first 5 (where site.RegularPages "Type" "in" site.Params.mainSections) }}
   {{ .Content }}
{{ end }}
{{< /code >}}

## Nest `where` Clauses

You can also nest `where` clauses to drill down on lists of content by more than one parameter. The following first grabs all pages in the "blog" section and then ranges through the result of the first `where` clause and finds all pages that are *not* featured:

```go-html-template
{{ range where (where .Pages "Section" "blog" ) ".Params.featured" "!=" true }}
```

## Unset Fields

Filtering only works for set fields. To check whether a field is set or exists, you can use the operand `nil`.

This can be useful to filter a small amount of pages from a large pool. Instead of set field on all pages, you can set field on required pages only.

Only the following operators are available for `nil`

* `=`, `==`, `eq`: True if the given field is not set.
* `!=`, `<>`, `ne`: True if the given field is set.

```go-html-template
{{ range where .Pages ".Params.specialpost" "!=" nil }}
   {{ .Content }}
{{ end }}
```

## Portable `where` filters -- `site.Params.mainSections` {#mainsections}

**This is especially important for themes.**

To list the most relevant pages on the front page or similar, you
should use the `site.Params.mainSections` list instead of comparing
section names to hard-coded values like `"posts"` or `"post"`.

```go-html-template
{{ $pages := where site.RegularPages "Type" "in" site.Params.mainSections }}
```

If the user has not set this config parameter in their site config, it
will default to the _section with the most pages_.

The user can override the default in `config.toml`:

```toml
[params]
  mainSections = ["blog", "docs"]
```

[intersect]: /functions/intersect/
[wherekeyword]: https://www.techonthenet.com/sql/where.php
