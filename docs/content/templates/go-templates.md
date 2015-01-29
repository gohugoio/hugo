---
aliases:
- /layout/go-templates/
- /layouts/go-templates/
date: 2013-07-01
menu:
  main:
    parent: layout
next: /templates/functions
prev: /templates/overview
title: Go Template Primer
weight: 15
---

Hugo uses the excellent [Go][] [html/template][gohtmltemplate] library for
its template engine. It is an extremely lightweight engine that provides a very
small amount of logic. In our experience it is just the right amount of
logic to be able to create a good static website. If you have used other
template systems from different languages or frameworks, you will find a lot of
similarities in Go templates.

This document is a brief primer on using Go templates. The [Go docs][gohtmltemplate]
go into more depth and cover features that aren't mentioned here.

## Introduction to Go Templates

Go templates provide an extremely simple template language. It adheres to the
belief that only the most basic of logic belongs in the template or view layer.
One consequence of this simplicity is that Go templates parse very quickly.

A unique characteristic of Go templates is they are content aware. Variables and
content will be sanitized depending on the context of where they are used. More
details can be found in the [Go docs][gohtmltemplate].

## Basic Syntax

Go lang templates are HTML files with the addition of variables and
functions.

**Go variables and functions are accessible within {{ }}**

Accessing a predefined variable "foo":

    {{ foo }}

**Parameters are separated using spaces**

Calling the `add` function with input of 1, 2:

    {{ add 1 2 }}

**Methods and fields are accessed via dot notation**

Accessing the Page Parameter "bar"

    {{ .Params.bar }}

**Parentheses can be used to group items together**

    {{ if or (isset .Params "alt") (isset .Params "caption") }} Caption {{ end }}


## Variables

Each Go template has a struct (object) made available to it. In Hugo, each
template is passed either a page or a node struct depending on which type of
page you are rendering. More details are available on the
[variables](/layout/variables/) page.

A variable is accessed by referencing the variable name.

    <title>{{ .Title }}</title>

Variables can also be defined and referenced.

    {{ $address := "123 Main St."}}
    {{ $address }}


## Functions

Go template ships with a few functions which provide basic functionality. The Go
template system also provides a mechanism for applications to extend the
available functions with their own. [Hugo template
functions](/layout/functions/) provide some additional functionality we believe
are useful for building websites. Functions are called by using their name
followed by the required parameters separated by spaces. Template
functions cannot be added without recompiling Hugo.

**Example 1: Adding numbers**

    {{ add 1 2 }}

**Example 2: Comparing numbers**

    {{ lt 1 2 }}

(There are more boolean operators, detailed in the
[template documentation](http://golang.org/pkg/text/template/#hdr-Functions).)

## Includes

When including another template, you will pass to it the data it will be
able to access. To pass along the current context, please remember to
include a trailing dot. The templates location will always be starting at
the /layout/ directory within Hugo.

**Example:**

    {{ template "partials/header.html" . }}

And, starting with Hugo v0.12, you may also use the `partial` call
for [partial templates](/templates/partials/):

    {{ partial "header.html" . }}


## Logic

Go templates provide the most basic iteration and conditional logic.

### Iteration

Just like in Go, the Go templates make heavy use of `range` to iterate over
a map, array or slice. The following are different examples of how to use
range.

**Example 1: Using Context**

    {{ range array }}
        {{ . }}
    {{ end }}

**Example 2: Declaring value variable name**

    {{range $element := array}}
        {{ $element }}
    {{ end }}

**Example 2: Declaring key and value variable name**

    {{range $index, $element := array}}
        {{ $index }}
        {{ $element }}
    {{ end }}

### Conditionals

`if`, `else`, `with`, `or` & `and` provide the framework for handling conditional
logic in Go Templates. Like `range`, each statement is closed with `end`.

Go Templates treat the following values as false:

* false
* 0
* any array, slice, map, or string of length zero

**Example 1: `if`**

    {{ if isset .Params "title" }}<h4>{{ index .Params "title" }}</h4>{{ end }}

**Example 2: `if` … `else`**

    {{ if isset .Params "alt" }}
        {{ index .Params "alt" }}
    {{else}}
        {{ index .Params "caption" }}
    {{ end }}

**Example 3: `and` & `or`**

    {{ if and (or (isset .Params "title") (isset .Params "caption")) (isset .Params "attr")}}

**Example 4: `with`**

An alternative way of writing "`if`" and then referencing the same value
is to use "`with`" instead. `with` rebinds the context `.` within its scope,
and skips the block if the variable is absent.

The first example above could be simplified as:

    {{ with .Params.title }}<h4>{{ . }}</h4>{{ end }}

**Example 5: `if` … `else if`**

    {{ if isset .Params "alt" }}
        {{ index .Params "alt" }}
    {{ else if isset .Params "caption" }}
        {{ index .Params "caption" }}
    {{ end }}

## Pipes

One of the most powerful components of Go templates is the ability to
stack actions one after another. This is done by using pipes. Borrowed
from Unix pipes, the concept is simple, each pipeline's output becomes the
input of the following pipe.

Because of the very simple syntax of Go templates, the pipe is essential
to being able to chain together function calls. One limitation of the
pipes is that they only can work with a single value and that value
becomes the last parameter of the next pipeline.

A few simple examples should help convey how to use the pipe.

**Example 1:**

    {{ if eq 1 1 }} Same {{ end }}

is the same as

    {{ eq 1 1 | if }} Same {{ end }}

It does look odd to place the `if` at the end, but it does provide a good
illustration of how to use the pipes.

**Example 2:**

    {{ index .Params "disqus_url" | html }}

Access the page parameter called "disqus_url" and escape the HTML.

**Example 3:**

    {{ if or (or (isset .Params "title") (isset .Params "caption")) (isset .Params "attr")}}
    Stuff Here
    {{ end }}

Could be rewritten as

    {{  isset .Params "caption" | or isset .Params "title" | or isset .Params "attr" | if }}
    Stuff Here
    {{ end }}

### Internet Explorer conditional comments using Pipes

By default, Go Templates remove HTML comments from output. This has the unfortunate side effect of removing Internet Explorer conditional comments. As a workaround, use something like this:

    {{ "<!--[if lt IE 9]>" | safeHtml }}
      <script src="html5shiv.js"></script>
    {{ "<![endif]-->" | safeHtml }}

Alternatively, use the backtick (`` ` ``) to quote the IE conditional comments, avoiding the tedious task of escaping every double quotes (`"`) inside, as demonstrated in the [examples](http://golang.org/pkg/text/template/#hdr-Examples) in the Go text/template documentation, e.g.:

```
{{ `<!--[if lt IE 7]><html class="no-js lt-ie9 lt-ie8 lt-ie7"><![endif]-->` | safeHtml }}
```

## Context (a.k.a. the dot)

The most easily overlooked concept to understand about Go templates is that `{{ . }}`
always refers to the current context. In the top level of your template, this
will be the data set made available to it. Inside of a iteration, however, it will have
the value of the current item. When inside of a loop, the context has changed:
`{{ . }}` will no longer refer to the data available to the entire page. If you need
to
access this from within the loop, you will likely want to do one of the following:

1. Set it to a variable instead of depending on the context.  For example:

        {{ $title := .Site.Title }}
        {{ range .Params.tags }}
          <li>
            <a href="{{ $baseurl }}/tags/{{ . | urlize }}">{{ . }}</a>
            - {{ $title }}
          </li>
        {{ end }}

    Notice how once we have entered the loop the value of `{{ . }}` has changed. We
    have defined a variable outside of the loop so we have access to it from within
    the loop.

2. Use `$.` to access the global context from anywhere.
   Here is an equivalent example:

        {{ range .Params.tags }}
          <li>
            <a href="{{ $baseurl }}/tags/{{ . | urlize }}">{{ . }}</a>
            - {{ $.Site.Title }}
          </li>
        {{ end }}

    This is because `$`, a special variable, is set to the starting value
    of `.` the dot by default,
    a [documented feature](http://golang.org/pkg/text/template/#hdr-Variables)
    of Go text/template.  Very handy, eh?

    > However, this little magic would cease to work if someone were to
    > mischievously redefine `$`, e.g. `{{ $ := .Site }}`.
    > *(No, don't do it!)*
    > You may, of course, recover from this mischief by using `{{ $ := . }}`
    > in a global context to reset `$` to its default value.

# Hugo Parameters

Hugo provides the option of passing values to the template language
through the site configuration (for sitewide values), or through the meta
data of each specific piece of content. You can define any values of any
type (supported by your front matter/config format) and use them however
you want to inside of your templates.


## Using Content (page) Parameters

In each piece of content, you can provide variables to be used by the
templates. This happens in the [front matter](/content/front-matter/).

An example of this is used in this documentation site. Most of the pages
benefit from having the table of contents provided. Sometimes the TOC just
doesn't make a lot of sense. We've defined a variable in our front matter
of some pages to turn off the TOC from being displayed.

Here is the example front matter:

```
---
title: "Permalinks"
date: "2013-11-18"
aliases:
  - "/doc/permalinks/"
groups: ["extras"]
groups_weight: 30
notoc: true
---
```

Here is the corresponding code inside of the template:

      {{ if not .Params.notoc }}
        <div id="toc" class="well col-md-4 col-sm-6">
        {{ .TableOfContents }}
        </div>
      {{ end }}



## Using Site (config) Parameters
In your top-level configuration file (e.g., `config.yaml`) you can define site
parameters, which are values which will be available to you in partials.

For instance, you might declare:

```yaml
params:
  CopyrightHTML: "Copyright &#xA9; 2013 John Doe. All Rights Reserved."
  TwitterUser: "spf13"
  SidebarRecentLimit: 5
```

Within a footer layout, you might then declare a `<footer>` which is only
provided if the `CopyrightHTML` parameter is provided, and if it is given,
you would declare it to be HTML-safe, so that the HTML entity is not escaped
again.  This would let you easily update just your top-level config file each
January 1st, instead of hunting through your templates.

```
{{if .Site.Params.CopyrightHTML}}<footer>
<div class="text-center">{{.Site.Params.CopyrightHTML | safeHtml}}</div>
</footer>{{end}}
```

An alternative way of writing the "`if`" and then referencing the same value
is to use "`with`" instead. With rebinds the context `.` within its scope,
and skips the block if the variable is absent:

```
{{with .Site.Params.TwitterUser}}<span class="twitter">
<a href="https://twitter.com/{{.}}" rel="author">
<img src="/images/twitter.png" width="48" height="48" title="Twitter: {{.}}"
 alt="Twitter"></a>
</span>{{end}}
```

Finally, if you want to pull "magic constants" out of your layouts, you can do
so, such as in this example:

```
<nav class="recent">
  <h1>Recent Posts</h1>
  <ul>{{range first .Site.Params.SidebarRecentLimit .Site.Recent}}
    <li><a href="{{.RelPermalink}}">{{.Title}}</a></li>
  {{end}}</ul>
</nav>
```


[go]: http://golang.org/
[gohtmltemplate]: http://golang.org/pkg/html/template/

# Template example: Show only upcoming events

Go allows you to do more than what's shown here.  Using Hugo's
[`where`](/templates/functions/#toc_4) function and Go built-ins, we can list
only the items from `content/events/` whose date (set in the front matter) is in
the future:

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
