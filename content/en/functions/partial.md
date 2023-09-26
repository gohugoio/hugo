---
title: partial
description: Executes the named partial template. If the partial contains a return statement, returns that value, else returns the rendered output.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: partials
relatedFuncs:
  - partials.Include
  - partials.IncludeCached
signature:
  - partials.Include LAYOUT [CONTEXT]
  - partial LAYOUT [CONTEXT]
---

In this example we have three partial templates:

```text
layouts/
└── partials/
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
{{ partial "breadcrumbs.html" }}
```

You can pass anything in context: a page, a page collection, a scalar value, a slice, or a map. For example:

```go-html-template
{{ $student := dict 
  "name" "John Doe" 
  "major" "Finance"
  "gpa" 4.0
}}
{{ partial "render-student-info.html" $student }}
```

Then, within the partial template:

```go-html-template
<p>{{ .name }} is majoring in {{ .major }}. Their grade point average is {{ .gpa }}.</p>
```


[breadcrumb navigation]: /content-management/sections/#ancestors-and-descendants
