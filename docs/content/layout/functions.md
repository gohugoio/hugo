---
title: "Template Functions"
date: "2013-07-01"
linktitle: "Template Functions"
groups: ["layout"]
groups_weight: 70
---

Hugo uses the excellent golang html/template library for its template engine.
It is an extremely lightweight engine that provides a very small amount of
logic. In our experience that it is just the right amount of logic to be able
to create a good static website.

Go templates are lightweight but extensible. Hugo has added the following
functions to the basic template logic.

Golang documentation for the built-in functions can be found [here](http://golang.org/pkg/text/template/)

## General

### isset
Return true if the parameter is set.
Takes either a slice, array or channel and an index or a map and a key as input.

eg. {{ if isset .Params "project_url" }} {{ index .Params "project_url" }}{{ end }}

### echoParam
If parameter is set, then echo it.

eg. {{echoParam .Params "project_url" }}

### first
Slices an array to only the first X elements.

eg.
    {{ range first 10 .Data.Pages }}
        {{ .Render "summary"}}
    {{ end }}


## Math

### add
Adds two integers.

eg {{add 1 2}} -> 3

### sub
Subtracts two integers.

eg {{sub 3 2}} -> 1

### div
Divides two integers.

eg {{div 6 3}} -> 2

### mul
Multiplies two integers.

eg {{mul 2 3}} -> 6

### mod
Modulus of two integers.

eg {{mod 15 3}} -> 0

### modBool
Boolean of modulus of two integers.
true if modulus is 0.

eg {{modBool 15 3}} -> true

## Strings

### urlize
Takes a string and sanitizes it for usage in urls, converts spaces to "-".

eg. &lt;a href="/tags/{{ . | urlize }}"&gt;{{ . }}&lt;/a&gt;

### safeHtml
Declares the provided string as "safe" so go templates will not filter it.

eg. {{ .Params.CopyrightHTML | safeHtml }}

### lower
Convert all characters in string to lowercase.

eg {{lower "BatMan"}} -> "batman"

### upper
Convert all characters in string to uppercase.

eg {{upper "BatMan"}} -> "BATMAN"

### title
Convert all characters in string to titlecase.

eg {{title "BatMan"}} -> "Batman"

### highlight
Take a string of code and a language, uses pygments to return the syntax
highlighted code in html. Used in the [highlight shortcode](/extras/highlight).
