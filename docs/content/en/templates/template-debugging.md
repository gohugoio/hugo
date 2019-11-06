---
title: Template Debugging
# linktitle: Template Debugging
description: You can use Go templates' `printf` function to debug your Hugo  templates. These snippets provide a quick and easy visualization of the variables available to you in different contexts.
godocref: https://golang.org/pkg/fmt/
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
keywords: [debugging,troubleshooting]
menu:
  docs:
    parent: "templates"
    weight: 180
weight: 180
sections_weight: 180
draft: false
aliases: []
toc: false
---

Here are some snippets you can add to your template to answer some common questions.

These snippets use the `printf` function available in all Go templates.  This function is an alias to the Go function, [fmt.Printf](https://golang.org/pkg/fmt/).

## What Variables are Available in this Context?

You can use the template syntax, `$.`, to get the top-level template context from anywhere in your template. This will print out all the values under, `.Site`.

```
{{ printf "%#v" $.Site }}
```

This will print out the value of `.Permalink`:


```
{{ printf "%#v" .Permalink }}
```


This will print out a list of all the variables scoped to the current context
(`.`, aka ["the dot"][tempintro]).


```
{{ printf "%#v" . }}
```


When developing a [homepage][], what does one of the pages you're looping through look like?

```
{{ range .Pages }}
    {{/* The context, ".", is now each one of the pages as it goes through the loop */}}
    {{ printf "%#v" . }}
{{ end }}
```

## Why Am I Showing No Defined Variables?

Check that you are passing variables in the `partial` function:

```
{{ partial "header" }}
```

This example will render the header partial, but the header partial will not have access to any contextual variables. You need to pass variables explicitly. For example, note the addition of ["the dot"][tempintro].

```
{{ partial "header" . }}
```

The dot (`.`) is considered fundamental to understanding Hugo templating. For more information, see [Introduction to Hugo Templating][tempintro].

[homepage]: /templates/homepage/
[tempintro]: /templates/introduction/
