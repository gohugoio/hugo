---
title: Template Debugging
linktitle: Template Debugging
description: You can use Go templates' `printf` function to debug your Hugo  templates. These snippets provide a quick and easy visualization of the variables available to you in different contexts.
godocref: http://golang.org/pkg/fmt/
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [debugging,troubleshooting]
weight: 180
draft: false
aliases: [/templates/debugging/]
toc: false
---


Here are some snippets you can add to your template to answer some common questions.

These snippets use the `printf` function available in all Go templates.  This function is an alias to the Go function, [fmt.Printf](http://golang.org/pkg/fmt/).

## What Variables are Available in this Context?

You can use the template syntax, `$.`, to get the top-level template context from anywhere in your template. This will print out all the values under, `.Site`.

{{% code file="get-top-level-syntax.sh" %}}
```golang
{{ printf "%#v" $.Site }}
```
{{% /code %}}

This will print out the value of `.Permalink`:

{{% code file="get-permalink.sh" %}}
```golang
{{ printf "%#v" .Permalink }}
```
{{% /code %}}

This will print out a list of all the variables scoped to the current context
(aka [The dot, "`.`"][thedot]).

{{% code file="get-all-vars-current-context.sh" %}}
```golang
{{ printf "%#v" . }}
```
{{% /code %}}

When writing a [Homepage][hometemplate], what does one of the pages you're looping through look like?

```golang
{{ range .Data.Pages }}
    {{/* The context, ".", is now each one of the pages as it goes through the loop */}}
    {{ printf "%#v" . }}
{{ end }}
```

{{% note "`.Date.Pages` on the Homepage" %}}
Using `.Data.Pages` on the homepage is the equivalent of writing `.Site.Pages`.
{{% /note %}}

## Why Am I Showing No Defined Variables?

Check that you are passing variables in the `partial` function:

```
{{ partial "header" }}
```

This example will render the header partial, but the header partial will not have access to any contextual variables. You need to pass variables explicitly. For example note the addition of [the dot][thedot].

```
{{ partial "header" . }}
```

The dot (`.`) is considered fundamental to understand Hugo templating. For more information, see the [Go Template Primer][primer].

[hometemplate]: /templates/homepage-template/
[primer]: /templates/go-template-primer/
[thedot]: /functions/the-dot/