---
title: Go Template Primer
linktitle:
description:
godocref: https://golang.org/pkg/html/template/
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
weight:
tags: []
categories: [templates]
draft: false
slug:
aliases: [/templates/go-templates/]
toc: true
notesforauthors:
---

Hugo uses the excellent [Go html/template][] library, an extremely lightweight engine that provides just the right amount of logic to be able to create a good static website. If you have used other template systems from different languages or frameworks, you will find a lot of similarities in Go templates.

This document is a brief primer on using Go templates. The [Go docs][gohtmltemplate] go into more depth and cover features that aren't mentioned here.

## Introduction to Go Templates

Go templates provide an extremely simple template language. It adheres to the belief that only the most basic of logic belongs in the template or view layer. One consequence of this simplicity is that Go templates parse very quickly.

A unique characteristic of Go templates is they are content aware. Variables and content will be sanitized depending on the context of where they are used. More details can be found in the [Go docs][gohtmltemplate].

## Basic Syntax

Golang templates are HTML files with the addition of [variables][variablesparams] and [functions][hugofunctions]. Golang template variables and functions are accessible within `{{ }}`.

### Accessing a Predefined Variable

```golang
{{ foo }}
```

Parameters for functions are separated using spaces. The following example calls the `add` function with inputs of `1` and `2`:

```golang
{{ add 1 2 }}
```d

**Methods and fields are accessed via dot notation**

Accessing the Page Parameter "bar"

```golang
{{ .Params.bar }}
```

**Parentheses can be used to group items together**

```golang
{{ if or (isset .Params "alt") (isset .Params "caption") }} Caption {{ end }}
```

## Variables

Each Go template has a struct (object) made available to it. In Hugo, each
template is passed page struct. More details are available in the [variables and params section][variablesparams].

A variable is accessed by referencing the variable name.

```golang
<title>{{ .Title }}</title>
```

Variables can also be defined and referenced.

```golang
{{ $address := "123 Main St."}}
{{ $address }}
```

## Functions

Go template ships with a few functions which provide basic functionality. The Go template system also provides a mechanism for applications to extend the available functions with their own. [Hugo template functions][hugofunctions] provide some additional functionality we believe are useful for building websites. Functions are called by using their name followed by the required parameters separated by spaces. Template functions cannot be added without recompiling Hugo.

### Example 1: Adding Numbers

```golang
{{ add 1 2 }}
```

### Example 2: Comparing Numbers

```golang
{{ lt 1 2 }}
```

{{% note "Additional Boolean Operators" %}}
There are more boolean operators than those listed in the Hugo docs in the [Golang template documentation](http://golang.org/pkg/text/template/#hdr-Functions).
{{% /note %}}

## Includes

When including another template, you will pass to it the data it will be
able to access. To pass along the current context, please remember to
include a trailing dot. The templates location will always be starting at
the /layout/ directory within Hugo.

### Template and Partial Examples

```golang
{{ template "partials/header.html" . }}
```

And, starting with Hugo v0.12, you may also use the `partial` call
for [partial templates][]:

```golang
{{ partial "header.html" . }}
```

## Logic

Go templates provide the most basic iteration and conditional logic.

### Iteration

Just like in Go, the Go templates make heavy use of `range` to iterate over
a map, array or slice. The following are different examples of how to use
range.

**Example 1: Using Context**

```golang
{{ range array }}
    {{ . }}
{{ end }}
```

**Example 2: Declaring value variable name**

```golang
{{range $element := array}}
    {{ $element }}
{{ end }}
```

**Example 2: Declaring key and value variable name**

```golang
{{range $index, $element := array}}
   {{ $index }}
   {{ $element }}
{{ end }}
```

### Conditionals

`if`, `else`, `with`, `or` & `and` provide the framework for handling conditional logic in Go Templates. Like `range`, each statement is closed with an `{{end}}`.

Go Templates treat the following values as false:

* false
* 0
* any array, slice, map, or string of length zero

**Example 1: `if`**

```golang
{{ if isset .Params "title" }}<h4>{{ index .Params "title" }}</h4>{{ end }}
```

**Example 2: `if` … `else`**

```golang
{{ if isset .Params "alt" }}
    {{ index .Params "alt" }}
{{else}}
    {{ index .Params "caption" }}
{{ end }}
```

**Example 3: `and` & `or`**

```golang
{{ if and (or (isset .Params "title") (isset .Params "caption")) (isset .Params "attr")}}
```

**Example 4: `with`**

An alternative way of writing "`if`" and then referencing the same value
is to use "`with`" instead. `with` rebinds the context `.` within its scope,
and skips the block if the variable is absent.

The first example above could be simplified as:

    {{ with .Params.title }}<h4>{{ . }}</h4>{{ end }}

**Example 5: `if` … `else if`**

```golang
{{ if isset .Params "alt" }}
    {{ index .Params "alt" }}
{{ else if isset .Params "caption" }}
    {{ index .Params "caption" }}
{{ end }}
```

## Pipes

One of the most powerful components of Go templates is the ability to stack actions one after another. This is done by using pipes. Borrowed from Unix pipes, the concept is simple, each pipeline's output becomes the input of the following pipe.

Because of the very simple syntax of Go templates, the pipe is essential to being able to chain together function calls. One limitation of the pipes is that they only can work with a single value and that value becomes the last parameter of the next pipeline.

A few simple examples should help convey how to use the pipe.

**Example 1:**

```golang
{{ shuffle (seq 1 5) }}
```

is the same as

```golang
{{ (seq 1 5) | shuffle }}
```

**Example 2:**

```golang
{{ index .Params "disqus_url" | html }}
```

Access the page parameter called "disqus_url" and escape the HTML.

The `index` function is a built in to [Go][] built-in. [You can read more about `index` in the Godocs][]. The Godocs have the following to say about`index`:

> ...returns the result of indexing its first argument by the following arguments. Thus "index x 1 2 3" is, in Go syntax, `x[1][2][3]`. Each indexed item must be a map, slice, or array.

**Example 3:**

    {{ if or (or (isset .Params "title") (isset .Params "caption")) (isset .Params "attr") }}
    Stuff Here
    {{ end }}

Could be rewritten as

    {{ if isset .Params "caption" | or isset .Params "title" | or isset .Params "attr" }}
    Stuff Here
    {{ end }}

### Internet Explorer Conditional Comments

By default, Go Templates remove HTML comments from output. This has the unfortunate side effect of removing Internet Explorer conditional comments. As a workaround, use something like this:

    {{ "<!--[if lt IE 9]>" | safeHTML }}
      <script src="html5shiv.js"></script>
    {{ "<![endif]-->" | safeHTML }}

Alternatively, use the backtick (`` ` ``) to quote the IE conditional comments, avoiding the tedious task of escaping every double quotes (`"`) inside, as demonstrated in the [examples](http://golang.org/pkg/text/template/#hdr-Examples) in the Go text/template documentation, e.g.:

```
{{ `<!--[if lt IE 7]><html class="no-js lt-ie9 lt-ie8 lt-ie7"><![endif]-->` | safeHTML }}
```

## Context (aka "the dot")

The most easily overlooked concept to understand about Go templates is that `{{ . }}` always refers to the current context. In the top level of your template, this will be the data set made available to it. Inside of a iteration, however, it will have the value of the current item. When inside of a loop, the context has changed: `{{ . }}` will no longer refer to the data available to the entire page. If you need to access this from within the loop, you will likely want to do one of the following:

### Define Variable Independent of Context

variable instead of depending on the context.  For example:

{{% input "range-through-tags-w-variable.html" %}}
```html
{{ $title := .Site.Title }}
{{ $base := .Site.BaseURL }}
<ul class="tags">
{{ range .Params.tags }}
    <li>
        <a href="{{ $base }}tags/{{ . | urlize }}">{{ . }}</a>
        - {{ $title }}
    </li>
{{ end }}
</ul>
```
{{% /input %}}

{{% note %}}
Notice how once we have entered the loop, the value of `{{ . }}` has changed. We have defined a variable outside of the loop (`{{$title}}`) so we have access to it from within the loop.
{{% /note %}}

### Use `$.` to Access the Global Context

from anywhere. Here is an equivalent example:

{{% input "range-through-tags-w-global.html" %}}
```html
{{ $base := .Site.BaseURL }}
<ul class="tags">
{{ range .Params.tags }}
  <li>
    <a href="{{$base}}tags/{{ . | urlize }}">{{ . }}</a>
            - {{ $.Site.Title }}
  </li>
{{ end }}
</ul>
```
{{% /input %}}

This is because `$`, a special variable, is set to the starting value of `.` ("the dot") by default. This is a [documented feature of Go text/template][].

{{% warning "Don't Redefine the Dot" %}}
The built-in magic of `$` would cease to work if someone were to mischievously redefine the special character; e.g. `{{ $ := .Site }}`. *Don't do it.* You may, of course, recover from this mischief by using `{{ $ := . }}` in a global context to reset `$` to its default value.
{{% /warning %}}

## Whitespace

Go 1.6 includes the ability to trim the whitespace from either side of a Go tag by including a hyphen (`-`) and space immediately beside the corresponding `{{` or `}}` delimiter.

For instance, the following Go template:

```html
<div>
  {{ .Title }}
</div>
```

will include the newlines and horizontal tab in its HTML output:

```html
<div>
  Hello, World!
</div>
```

whereas using

```html
<div>
  {{- .Title -}}
</div>
```

in that case will output simply

```html
<div>Hello, World!</div>
```

Go considers the following characters as whitespace:

* <kbd>space</kbd>
* horizontal <kbd>tab</kbd>
* carriage <kbd>return</kbd>
* newline

## Hugo Parameters

Hugo provides the option of passing values to the template language through the site configuration (for sitewide values), or through the meta data of each specific piece of content. You can define any values of any type (supported by your front matter/config format) and use them however you want to inside of your templates.

## Using Content (page) Parameters

In each piece of content, you can provide variables to be used by the templates. This happens in the [front matter][].

An example of this is used in this documentation site. Most of the pages benefit from having the table of contents provided. Sometimes the TOC just doesn't make a lot of sense. We've defined a variable in our front matter of some pages to turn off the TOC from being displayed.

Here is the example front matter:

```yaml
---
title: "Permalinks"
lastmod: 2015-11-30
date: "2013-11-18"
aliases:
  - "/doc/permalinks/"
groups: ["extras"]
groups_weight: 30
notoc: true
---
```

Here is the corresponding code inside of the template:

```html
{{ if not .Params.notoc }}
    <div id="toc" class="well col-md-4 col-sm-6">
    {{ .TableOfContents }}
    </div>
{{ end }}
```

## Using Site (config) Parameters

In your top-level configuration file (e.g., `config.yaml`), you can define site-level parameters that are available to you as variables throughout your templates.

For instance, you might declare:

{{% input "config.yaml" %}}
```yaml
params:
  CopyrightHTML: "Copyright &#xA9; 2013 John Doe. All Rights Reserved."
  TwitterUser: "spf13"
  SidebarRecentLimit: 5
```
{{% /input %}}

Within a footer layout, you might then declare a `<footer>` which is only provided if the `CopyrightHTML` parameter is provided, and if it is given, you would declare it to be HTML-safe, so that the HTML entity is not escaped again.  This would let you easily update just your top-level config file each January 1st, instead of hunting through your templates.

```html
{{if .Site.Params.CopyrightHTML}}<footer>
<div class="text-center">{{.Site.Params.CopyrightHTML | safeHTML}}</div>
</footer>{{end}}
```

An alternative way of writing the "`if`" and then referencing the same value is to use "`with`" instead. With rebinds the context `.` within its scope, and skips the block if the variable is absent:

```html
{{with .Site.Params.TwitterUser}}<span class="twitter">
<a href="https://twitter.com/{{.}}" rel="author">
<img src="/images/twitter.png" width="48" height="48" title="Twitter: {{.}}"
 alt="Twitter"></a>
</span>{{end}}
```

Finally, if you want to pull "magic constants" out of your layouts, you can do so, such as in this example:

```html
<nav class="recent">
  <h1>Recent Posts</h1>
  <ul>{{range first .Site.Params.SidebarRecentLimit .Site.Pages}}
    <li><a href="{{.RelPermalink}}">{{.Title}}</a></li>
  {{end}}</ul>
</nav>
```

## Template example: Show only upcoming events

Go allows you to do more than what's shown here.  Using Hugo's [`where` function][] and Go built-ins, we can list only the items from `content/events/` whose date (set in the front matter) is in the future:

{{% input "show-upcoming-dates.html" %}}
```golang
<h4>Upcoming Events</h4>
<ul class="upcoming-events">
{{ range where .Data.Pages.ByDate "Section" "events" }}
  {{ if ge .Date.Unix .Now.Unix }}
    <li><span class="event-type">{{ .Type | title }} —</span>
      {{ .Title }}
      on <span class="event-date">
      {{ .Date.Format "2 January at 3:04pm" }}</span>
      at {{ .Params.place }}
    </li>
  {{ end }}
{{ end }}
```
{{% /input %}}

[`where` function]: /functions/where/
[documented feature of Go text/template]: http://golang.org/pkg/text/template/#hdr-Variables
[front matter]: /content-management/front-matter/
[Go html/template]: http://golang.org/pkg/html/template/ "Godocs references for Golang's html templating"
[gohtmltemplate]: http://golang.org/pkg/html/template/ "Godocs references for Golang's html templating"
[hugofunctions]: /functions/ "Link to section for Hugo's templating functions"
[partial templates]: /templates/partials-templates/ "Link to the partial templates page inside of the templating section of the Hugo docs"
[variablesparams]: /variables-and-params/ "Link to the list page for the Variables and Params section of the site."
[You can read more about `index` in the Godocs]: http://golang.org/pkg/text/template/ "Godocs page for index function"