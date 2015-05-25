---
aliases:
- /doc/debugging/
- /layout/debugging/
date: 2015-05-22
linktitle: Debugging
menu:
  main:
    parent: layout
prev: /templates/404
title: Template Debugging
weight: 110
---


# Template Debugging

Here are some snippets you can add to your template to answer some common questions.
These snippets use the `printf` function available in all Go templates.  This function is
an alias to the Go function, [fmt.Printf](http://golang.org/pkg/fmt/).

### What type of page is this?

Does Hugo consider this page to be a "Node" or a "Page"? (Put this snippet at
the top level of your template. Don't use it inside of a `range` loop.)

    {{ printf "%T" . }}


### What variables are available in this context?

You can use the template syntax, `$.`, to get the top-level template context
from anywhere in your template.  This will print out all the values under, `.Site`.

    {{ printf "%#v" $.Site }}

This will print out the value of `.Permalink`, which is available on both Nodes
and Pages.

    {{ printf "%#v" .Permalink }}

This will print out a list of all the variables scoped to the current context
(a.k.a. The dot, "`.`").

    {{ printf "%#v" . }}

When writing a [Homepage](/templates/homepage), what does one of the pages
you're looping through look like?

```
{{ range .Data.Pages }}
    {{/* The context, ".", is now a Page */}}
    {{ printf "%#v" . }}
{{ end }}
```
