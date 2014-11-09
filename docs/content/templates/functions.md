---
aliases:
- /layout/functions/
date: 2013-07-01
linktitle: Functions
menu:
  main:
    parent: layout
next: /templates/variables
prev: /templates/go-templates
title: Hugo Template Functions
weight: 20
---

Hugo uses the excellent Go html/template library for its template engine.
It is an extremely lightweight engine that provides a very small amount of
logic. In our experience it is just the right amount of logic to be able
to create a good static website.

Go templates are lightweight but extensible. Hugo has added the following
functions to the basic template logic.

(Go itself supplies built-in functions, including comparison operators
and other basic tools; these are listed in the
[Go template documentation](http://golang.org/pkg/text/template/#hdr-Functions).)

## General

### isset
Return true if the parameter is set.
Takes either a slice, array or channel and an index or a map and a key as input.

e.g. {{ if isset .Params "project_url" }} {{ index .Params "project_url" }}{{ end }}

### echoParam
If parameter is set, then echo it.

e.g. {{echoParam .Params "project_url" }}

### eq
Return true if the parameters are equal.

e.g.
    {{ if eq .Section "blog" }}current{{ end}}"

### first
Slices an array to only the first X elements.

Works on [lists](/templates/list/), [taxonomies](/taxonomies/displaying/), [terms](/templates/terms/), [groups](/templates/list/)

e.g.
    {{ range first 10 .Data.Pages }}
        {{ .Render "summary"}}
    {{ end }}

### where
Filters an array to only elements containing a matching value for a given field.

Works on [lists](/templates/list/), [taxonomies](/taxonomies/displaying/), [terms](/templates/terms/), [groups](/templates/list/)

e.g.

    {{ range where .Data.Pages "Section" "post" }}
       {{ .Content}}
    {{ end }}

*where and first can be stacked*

e.g.

    {{ range first 5 (where .Data.Pages "Section" "post") }}
       {{ .Content}}
    {{ end }}

### in
Checks if an element is in an array (or slice) and returns a boolean.  The elements supported are strings, integers and floats (only float64 will match as expected).  In addition, it can also check if a substring exists in a string.

e.g.
    {{ if in .Params.tags "Git" }}Follow me on GitHub!{{ end }}
or
    {{ if in "this string contains a substring" "substring" }}Substring found!{{ end }}

### intersect
Given two arrays (or slices), this function will return the common elements in the arrays.  The elements supported are strings, integers and floats (only float64).

A useful example of this functionality is a 'similar posts' block.  Create a list of links to posts where any of the tags in the current post match any tags in other posts.

e.g.
    <ul>
    {{ $page_link := .Permalink }}
    {{ $tags := .Params.tags }}
    {{ range .Site.Recent }}
        {{ $page := . }}
        {{ $has_common_tags := intersect $tags .Params.tags | len | lt 0 }}
        {{ if and $has_common_tags (ne $page_link $page.Permalink) }}
            <li><a href="{{ $page.Permalink }}">{{ $page.Title }}</a></li>
        {{ end }}
    {{ end }}
    </ul>


## Math

### add
Adds two integers.

e.g. {{add 1 2}} → 3

### sub
Subtracts two integers.

e.g. {{sub 3 2}} → 1

### div
Divides two integers.

e.g. {{div 6 3}} → 2

### mul
Multiplies two integers.

e.g. {{mul 2 3}} → 6

### mod
Modulus of two integers.

e.g. {{mod 15 3}} → 0

### modBool
Boolean of modulus of two integers.
true if modulus is 0.

e.g. {{modBool 15 3}} → true

## Strings

### urlize
Takes a string and sanitizes it for usage in URLs, converts spaces to "-".

e.g. &lt;a href="/tags/{{ . | urlize }}"&gt;{{ . }}&lt;/a&gt;

### safeHtml
Declares the provided string as "safe" so Go templates will not filter it.

e.g. {{ .Params.CopyrightHTML | safeHtml }}

### lower
Convert all characters in string to lowercase.

e.g. {{lower "BatMan"}} → "batman"

### upper
Convert all characters in string to uppercase.

e.g. {{upper "BatMan"}} → "BATMAN"

### title
Convert all characters in string to titlecase.

e.g. {{title "BatMan"}} → "Batman"

### highlight
Take a string of code and a language, uses Pygments to return the syntax
highlighted code in HTML. Used in the [highlight shortcode](/extras/highlighting).
