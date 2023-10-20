---
title: compare.Default
linkTitle: default
description: Allows setting a default value that can be returned if a first value is not set.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [default]
  returnType: any
  signatures: [compare.Default DEFAULT INPUT]
relatedFunctions:
  - compare.Conditional
  - compare.Default
aliases: [/functions/default]
---

`default` checks whether a given value is set and returns a default value if it is not. *Set* in this context means different things depending on the data type:

* non-zero for numeric types and times
* non-zero length for strings, arrays, slices, and maps
* any boolean or struct value
* non-nil for any other types

`default` function examples reference the following content page:

{{< code file="content/posts/default-function-example.md" >}}
---
title: Sane Defaults
seo_title:
date: 2017-02-18
font:
oldparam: The default function helps make your templating DRYer.
newparam:
---
{{< /code >}}

`default` can be written in more than one way:

```go-html-template
{{ .Params.font | default "Roboto" }}
{{ default "Roboto" .Params.font }}
```

Both of the above `default` function calls return `Roboto`.

A `default` value, however, does not need to be hard coded like the previous example. The `default` value can be a variable or pulled directly from the front matter using dot notation:

```go-html-template
{{ $old := .Params.oldparam }}
<p>{{ .Params.newparam | default $old }}</p>
```

Which would return:

```html
<p>The default function helps make your templating DRYer.</p>
```

And then using dot notation

```go-html-template
<title>{{ .Params.seo_title | default .Title }}</title>
```

Which would return

```html
<title>Sane Defaults</title>
```

The following have equivalent return values but are far less terse. This demonstrates the utility of `default`:

Using `if`:

```go-html-template
<title>{{ if .Params.seo_title }}{{ .Params.seo_title }}{{ else }}{{ .Title }}{{ end }}</title>
=> Sane Defaults
```

Using `with`:

```go-html-template
<title>{{ with .Params.seo_title }}{{ . }}{{ else }}{{ .Title }}{{ end }}</title>
=> Sane Defaults
```
