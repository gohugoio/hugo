---
title: Template debugging
description: You can use Go templates' `printf` function to debug your Hugo  templates. These snippets provide a quick and easy visualization of the variables available to you in different contexts.
categories: [templates]
keywords: [debugging,troubleshooting]
menu:
  docs:
    parent: templates
    weight: 240
weight: 240
---

Here are some snippets you can add to your template to answer some common questions.

These snippets use the `printf` function available in all Go templates.  This function is an alias to the Go function, [fmt.Printf](https://pkg.go.dev/fmt).

## What variables are available in this context?

You can use the template syntax, `$.`, to get the top-level template context from anywhere in your template. This will print out all the values under, `.Site`.

```go-html-template
{{ printf "%#v" $.Site }}
```

This will print out the value of `.Permalink`:

```go-html-template
{{ printf "%#v" .Permalink }}
```

This will print out a list of all the variables scoped to the current context
(`.`, aka ["the dot"][tempintro]).

```go-html-template
{{ printf "%#v" . }}
```

When developing a [homepage], what does one of the pages you're looping through look like?

```go-html-template
{{ range .Pages }}
    {{/* The context, ".", is now each one of the pages as it goes through the loop */}}
    {{ printf "%#v" . }}
{{ end }}
```

In some cases it might be more helpful to use the following snippet on the current context to get a pretty printed view of what you're able to work with:

```go-html-template
<pre>{{ . | jsonify (dict "indent" " ") }}</pre>
```

Note that Hugo will throw an error if you attempt to use this construct to display context that includes a page collection (e.g., the context passed to home, section, taxonomy, and term templates).

## Why am I showing no defined variables?

Check that you are passing variables in the `partial` function:

```go-html-template
{{ partial "header.html" }}
```

This example will render the header partial, but the header partial will not have access to any contextual variables. You need to pass variables explicitly. For example, note the addition of ["the dot"][tempintro].

```go-html-template
{{ partial "header.html" . }}
```

The dot (`.`) is considered fundamental to understanding Hugo templating. For more information, see [Introduction to Hugo Templating][tempintro].

[homepage]: /templates/homepage/
[tempintro]: /templates/introduction/
