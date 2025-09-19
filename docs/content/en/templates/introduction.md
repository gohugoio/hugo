---
title: Introduction to templating
linkTitle: Introduction
description: An introduction to Hugo's templating syntax.
categories: []
keywords: []
weight: 10
---

{{< newtemplatesystem >}}

{{% glossary-term template %}}

Templates use [variables], [functions], and [methods] to transform your content, resources, and data into a published page.

> [!note]
> Hugo uses Go's [text/template] and [html/template] packages.
>
> The text/template package implements data-driven templates for generating textual output, while the html/template package implements data-driven templates for generating HTML output safe against code injection.
>
> By default, Hugo uses the html/template package when rendering HTML files.

For example, this HTML template initializes the `$v1` and `$v2` variables, then displays them and their product within an HTML paragraph.

```go-html-template
{{ $v1 := 6 }}
{{ $v2 := 7 }}
<p>The product of {{ $v1 }} and {{ $v2 }} is {{ mul $v1 $v2 }}.</p>
```

While HTML templates are the most common, you can create templates for any [output format](g) including CSV, JSON, RSS, and plain text.

## Context

The most important concept to understand before creating a template is _context_, the data passed into each template. The data may be a simple value, or more commonly [objects](g) and associated [methods](g).

For example, a _page_ template receives a `Page` object, and the `Page` object provides methods to return values or perform actions.

### Current context

Within a template, the dot (`.`) represents the current context.

```go-html-template {file="layouts/page.html"}
<h2>{{ .Title }}</h2>
```

In the example above the dot represents the `Page` object, and we call its [`Title`] method to return the title as defined in [front matter].

The current context may change within a template. For example, at the top of a template the context might be a `Page` object, but we rebind the context to another value or object within [`range`] or [`with`] blocks.

```go-html-template {file="layouts/page.html"}
<h2>{{ .Title }}</h2>

{{ range slice "foo" "bar" }}
  <p>{{ . }}</p>
{{ end }}

{{ with "baz" }}
  <p>{{ . }}</p>
{{ end }}
```

In the example above, the context changes as we `range` through the [slice](g) of values. In the first iteration the context is "foo", and in the second iteration the context is "bar". Inside of the `with` block the context is "baz". Hugo renders the above to:

```html
<h2>My Page Title</h2>
<p>foo</p>
<p>bar</p>
<p>baz</p>
```

### Template context

Within a `range` or `with` block you can access the context passed into the template by prepending a dollar sign (`$`) to the dot:

```go-html-template {file="layouts/page.html"}
{{ with "foo" }}
  <p>{{ $.Title }} - {{ . }}</p>
{{ end }}
```

Hugo renders this to:

```html
<p>My Page Title - foo</p>
```

> [!note]
> Make sure that you thoroughly understand the concept of _context_ before you continue reading. The most common templating errors made by new users relate to context.

## Actions

In the examples above the paired opening and closing braces represent the beginning and end of a template action, a data evaluation or control structure within a template.

A template action may contain literal values ([boolean](g), [string](g), [integer](g), and [float](g)), variables, functions, and methods.

```go-html-template {file="layouts/page.html"}
{{ $convertToLower := true }}
{{ if $convertToLower }}
  <h2>{{ strings.ToLower .Title }}</h2>
{{ end }}
```

In the example above:

- `$convertToLower` is a variable
- `true` is a literal boolean value
- `strings.ToLower` is a function that converts all characters to lowercase
- `Title` is a method on a the `Page` object

Hugo renders the above to:

```html
  
  
    <h2>my page title</h2>
  
```

### Whitespace

Notice the blank lines and indentation in the previous example? Although irrelevant in production when you typically minify the output, you can remove the adjacent whitespace by using template action delimiters with hyphens:

```go-html-template {file="layouts/page.html"}
{{- $convertToLower := true -}}
{{- if $convertToLower -}}
  <h2>{{ strings.ToLower .Title }}</h2>
{{- end -}}
```

Hugo renders this to:

```html
<h2>my page title</h2>
```

Whitespace includes spaces, horizontal tabs, carriage returns, and newlines.

### Pipes

Within a template action you may [pipe](g) a value to a function or method. The piped value becomes the final argument to the function or method. For example, these are equivalent:

```go-html-template
{{ strings.ToLower "Hugo" }} → hugo
{{ "Hugo" | strings.ToLower }} → hugo
```

You can pipe the result of one function or method into another. For example, these are equivalent:

```go-html-template
{{ strings.TrimSuffix "o" (strings.ToLower "Hugo") }} → hug
{{ "Hugo" | strings.ToLower | strings.TrimSuffix "o" }} → hug
```

These are also equivalent:

```go-html-template
{{ mul 6 (add 2 5) }} → 42
{{ 5 | add 2 | mul 6 }} → 42
```

> [!note]
> Remember that the piped value becomes the final argument to the function or method to which you are piping.

### Line splitting

You can split a template action over two or more lines. For example, these are equivalent:

```go-html-template
{{ $v := or $arg1 $arg2 }}

{{ $v := or 
  $arg1
  $arg2
}}
```

You can also split [raw string literals](g) over two or more lines. For example, these are equivalent:

```go-html-template
{{ $msg := "This is line one.\nThis is line two." }}

{{ $msg := `This is line one.
This is line two.`
}}
```

## Variables

A variable is a user-defined [identifier](g) prepended with a dollar sign (`$`), representing a value of any data type, initialized or assigned within a template action. For example, `$foo` and `$bar` are variables.

Variables may contain [scalars](g), [slices](g), [maps](g), or [objects](g).

Use `:=` to initialize a variable, and use `=` to assign a value to a variable that has been previously initialized. For example:

```go-html-template
{{ $total := 3 }}
{{ range slice 7 11 21 }}
  {{ $total = add $total . }}
{{ end }}
{{ $total }} → 42
```

Variables initialized inside of an `if`, `range`, or `with` block are scoped to the block. Variables initialized outside of these blocks are scoped to the template.

With variables that represent a slice or map, use the [`index`] function to return the desired value.

```go-html-template
{{ $slice := slice "foo" "bar" "baz" }}
{{ index $slice 2 }} → baz

{{ $map := dict "a" "foo" "b" "bar" "c" "baz" }}
{{ index $map "c" }} → baz
```

> [!note]
> Slices and arrays are zero-based; element 0 is the first element.

With variables that represent a map or object, [chain](g) identifiers to return the desired value or to access the desired method.

```go-html-template
{{ $map := dict "a" "foo" "b" "bar" "c" "baz" }}
{{ $map.c }} → baz

{{ $homePage := .Site.Home }}
{{ $homePage.Title }} → My Homepage
```

> [!note]
> As seen above, object and method names are capitalized. Although not required, to avoid confusion we recommend beginning variable and map key names with a lowercase letter or underscore.

## Functions

Used within a template action, a function takes one or more arguments and returns a value. Unlike methods, functions are not associated with an object.

Go's text/template and html/template packages provide a small set of functions, operators, and statements for general use. See the [go-templates] section of the function documentation for details.

Hugo provides hundreds of custom [functions] categorized by namespace. For example, the `strings` namespace includes these and other functions:

Function|Alias
:--|:--
[`strings.ToLower`](/functions/strings/tolower)|`lower`
[`strings.ToUpper`](/functions/strings/toupper)|`upper`
[`strings.Replace`](/functions/strings/replace)|`replace`

As shown above, frequently used functions have an alias. Use aliases in your templates to reduce code length.

When calling a function, separate the arguments from the function, and from each other, with a space. For example:

```go-html-template
{{ $total := add 1 2 3 4 }}
```

## Methods

Used within a template action and associated with an object, a method takes zero or more arguments and either returns a value or performs an action.

The most commonly accessed objects are the [`Page`] and [`Site`] objects. This is a small sampling of the [methods] available to each object.

Object|Method|Description
:--|:--|:--
`Page`|[`Date`](methods/page/date/)|Returns the date of the given page.
`Page`|[`Params`](methods/page/params/)|Returns a map of custom parameters as defined in the front matter of the given page.
`Page`|[`Title`](methods/page/title/)|Returns the title of the given page.
`Site`|[`Data`](methods/site/data/)|Returns a data structure composed from the files in the `data` directory.
`Site`|[`Params`](methods/site/params/)|Returns a map of custom parameters as defined in the site configuration.
`Site`|[`Title`](methods/site/title/)|Returns the title as defined in the site configuration.

Chain the method to its object with a dot (`.`) as shown below, remembering that the leading dot represents the [current context].

```go-html-template {file="layouts/page.html"}
{{ .Site.Title }} → My Site Title
{{ .Page.Title }} → My Page Title
```

The context passed into most templates is a `Page` object, so this is equivalent to the previous example:

```go-html-template {file="layouts/page.html"}
{{ .Site.Title }} → My Site Title
{{ .Title }} → My Page Title
```

Some methods take an argument. Separate the argument from the method with a space. For example:

```go-html-template {file="layouts/page.html"}
{{ $page := .Page.GetPage "/books/les-miserables" }}
{{ $page.Title }} → Les Misérables
```

## Comments

> [!note]
> Do not attempt to use HTML comment delimiters to comment out template code.
>
> Hugo strips HTML comments when rendering a page, but first evaluates any template code within the HTML comment delimiters. Depending on the template code within the HTML comment delimiters, this could cause unexpected results or fail the build.

Template comments are similar to template actions. Paired opening and closing braces represent the beginning and end of a comment. For example:

```text
{{/* This is an inline comment. */}}
{{- /* This is an inline comment with adjacent whitespace removed. */ -}}
```

Code within a comment is not parsed, executed, or displayed. Comments may be inline, as shown above, or in block form:

```text
{{/*
This is a block comment.
*/}}

{{- /*
This is a block comment with
adjacent whitespace removed.
*/ -}}
```

You may not nest one comment inside of another.

To render an HTML comment, pass a string through the [`safeHTML`] template function. For example:

```go-html-template
{{ "<!-- I am an HTML comment. -->" | safeHTML }}
{{ printf "<!-- This is the %s site. -->" .Site.Title | safeHTML }}
```

## Include

Use the [`template`] function to include one or more of Hugo's [embedded templates]:

```go-html-template
{{ partial "google_analytics.html" . }}
{{ partial "opengraph" . }}
{{ partial "pagination.html" . }}
{{ partial "schema.html" . }}
{{ partial "twitter_cards.html" . }}
```

Use the [`partial`] or [`partialCached`] function to include one or more [partial templates]:

```go-html-template
{{ partial "breadcrumbs.html" . }}
{{ partialCached "css.html" . }}
```

Create your _partial_ templates in the `layouts/_partials` directory.

> [!note]
> In the examples above, note that we are passing the current context (the dot) to each of the templates.

## Examples

This limited set of contrived examples demonstrates some of concepts described above. Please see the [functions], [methods], and [templates] documentation for specific examples.

### Conditional blocks

See documentation for [`if`], [`else`], and [`end`].

```go-html-template
{{ $var := 42 }}
{{ if eq $var 6 }}
  {{ print "var is 6" }}
{{ else if eq $var 7 }}
  {{ print "var is 7" }}
{{ else if eq $var 42 }}
  {{ print "var is 42" }}
{{ else }}
  {{ print "var is something else" }}
{{ end }}
```

### Logical operators

See documentation for [`and`] and [`or`].

```go-html-template
{{ $v1 := true }}
{{ $v2 := false }}
{{ $v3 := false }}
{{ $result := false }}

{{ if and $v1 $v2 $v3 }}
  {{ $result = true }}
{{ end }}
{{ $result }} → false

{{ if or $v1 $v2 $v3 }}
  {{ $result = true }}
{{ end }}
{{ $result }} → true
```

### Loops

See documentation for [`range`], [`else`], and [`end`].

```go-html-template
{{ $s := slice "foo" "bar" "baz" }}
{{ range $s }}
  <p>{{ . }}</p>
{{ else }}
  <p>The collection is empty</p>
{{ end }}
```

To loop a specified number of times:

```go-html-template
{{ $s := slice }}
{{ range 3 }}
  {{ $s = $s | append . }}
{{ end }}
{{ $s }} → [0 1 2]
```

### Rebind context

See documentation for [`with`], [`else`], and [`end`].

```go-html-template
{{ $var := "foo" }}
{{ with $var }}
  {{ . }} → foo
{{ else }}
  {{ print "var is falsy" }}
{{ end }}
```

To test multiple conditions:

```go-html-template
{{ $v1 := 0 }}
{{ $v2 := 42 }}
{{ with $v1 }}
  {{ . }}
{{ else with $v2 }}
  {{ . }} → 42
{{ else }}
  {{ print "v1 and v2 are falsy" }}
{{ end }}
```

### Access site parameters

See documentation for the [`Params`](/methods/site/params/) method on a `Site` object.

With this site configuration:

{{< code-toggle file=hugo >}}
title = 'ABC Widgets'
baseURL = 'https://example.org'
[params]
  subtitle = 'The Best Widgets on Earth'
  copyright-year = '2023'
  [params.author]
    email = 'jsmith@example.org'
    name = 'John Smith'
  [params.layouts]
    rfc_1123 = 'Mon, 02 Jan 2006 15:04:05 MST'
    rfc_3339 = '2006-01-02T15:04:05-07:00'
{{< /code-toggle >}}

Access the custom site parameters by chaining the identifiers:

```go-html-template
{{ .Site.Params.subtitle }} → The Best Widgets on Earth
{{ .Site.Params.author.name }} → John Smith

{{ $layout := .Site.Params.layouts.rfc_1123 }}
{{ .Site.Lastmod.Format $layout }} → Tue, 17 Oct 2023 13:21:02 PDT
```

### Access page parameters

See documentation for the [`Params`](/methods/page/params/) method on a `Page` object.

By way of example, consider this front matter:

{{< code-toggle file=content/annual-conference.md fm=true >}}
title = 'Annual conference'
date = 2023-10-17T15:11:37-07:00
[params]
display_related = true
key-with-hyphens = 'must use index function'
[params.author]
  email = 'jsmith@example.org'
  name = 'John Smith'
{{< /code-toggle >}}

The `title` and `date` fields are standard [front matter fields], while the other fields are user-defined.

Access the custom fields by [chaining](g) the [identifiers](g) when needed:

```go-html-template
{{ .Params.display_related }} → true
{{ .Params.author.email }} → jsmith@example.org
{{ .Params.author.name }} → John Smith
```

In the template example above, each of the keys is a valid identifier. For example, none of the keys contains a hyphen. To access a key that is not a valid identifier, use the [`index`] function:

```go-html-template
{{ index .Params "key-with-hyphens" }} → must use index function
```

[`and`]: /functions/go-template/and
[`else`]: /functions/go-template/else/
[`end`]: /functions/go-template/end/
[`if`]: /functions/go-template/if/
[`index`]: /functions/collections/indexfunction/
[`or`]: /functions/go-template/or
[`Page`]: /methods/page/
[`partial`]: /functions/partials/include/
[`partialCached`]: /functions/partials/includecached/
[`range`]: /functions/go-template/range/
[`safeHTML`]: /functions/safe/html
[`Site`]: /methods/site/
[`template`]: /functions/go-template/template/
[`Title`]: /methods/page/title
[`with`]: /functions/go-template/with/
[current context]: #current-context
[embedded templates]: /templates/embedded/
[front matter]: /content-management/front-matter/
[front matter fields]: /content-management/front-matter/#fields
[functions]: /functions/
[go-templates]: /functions/go-template/
[html/template]: https://pkg.go.dev/html/template
[methods]: /methods/
[partial templates]: /templates/types/#partial
[templates]: /templates/
[text/template]: https://pkg.go.dev/text/template
[variables]: #variables
