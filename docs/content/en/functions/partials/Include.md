---
title: partials.Include
description: Executes the given template, optionally passing context. If the  partial template contains a return statement, returns the given value, else returns the rendered output.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [partial]
    returnType: any
    signatures: ['partials.Include NAME [CONTEXT]']
aliases: [/functions/partial]
---

Without a [`return`] statement, the `partial` function returns a string of type `template.HTML`. With a `return` statement, the `partial` function can return any data type.

In this example we have three _partial_ templates:

```text
layouts/
└── _partials/
    ├── average.html
    ├── breadcrumbs.html
    └── footer.html
```

The "average" partial returns the average of one or more numbers. We pass the numbers in context:

```go-html-template
{{ $numbers := slice 1 6 7 42 }}
{{ $average := partial "average.html" $numbers }}
```

The "breadcrumbs" partial renders [breadcrumb navigation], and needs to receive the current page in context:

```go-html-template
{{ partial "breadcrumbs.html" . }}
```

The "footer" partial renders the site footer. In this contrived example, the footer does not need access to the current page, so we can omit context:

```go-html-template
{{ partial "footer.html" }}
```

You can pass anything in context: a page, a page collection, a scalar value, a slice, or a map. In this example we pass the current page and three scalar values:

```go-html-template
{{ $ctx := dict 
  "page" .
  "name" "John Doe" 
  "major" "Finance"
  "gpa" 4.0
}}
{{ partial "render-student-info.html" $ctx }}
```

Then, within the _partial_ template:

```go-html-template
<p>{{ .name }} is majoring in {{ .major }}.</p>
<p>Their grade point average is {{ .gpa }}.</p>
<p>See <a href="{{ .page.RelPermalink }}">details.</a></p>
```

To return a value from a _partial_ template, it must contain only one `return` statement, placed at the end of the template:

```go-html-template
{{ $result := "" }}
{{ if math.ModBool . 2 }}
  {{ $result = "even" }}
{{ else }}
  {{ $result = "odd" }}
{{ end }}
{{ return $result }}
```

See&nbsp;[details][`return`].

[`return`]: /functions/go-template/return/
[breadcrumb navigation]: /content-management/sections/#ancestors-and-descendants
