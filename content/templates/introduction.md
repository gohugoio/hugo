---
title: Introduction to Hugo Templating
linktitle: Introduction
description: Hugo uses Go's `html/template` and `text/template` libraries as the basis for the templating.
godocref: https://golang.org/pkg/html/template/
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-25
categories: [templates,fundamentals]
keywords: [go]
menu:
  docs:
    parent: "templates"
    weight: 10
weight: 10
sections_weight: 10
draft: false
aliases: [/templates/introduction/,/layouts/introduction/,/layout/introduction/, /templates/go-templates/]
toc: true
---

{{% note %}}
The following is only a primer on Go templates. For an in-depth look into Go templates, check the official [Go docs](http://golang.org/pkg/html/template/).
{{% /note %}}

Go templates provide an extremely simple template language that adheres to the belief that only the most basic of logic belongs in the template or view layer.

{{< youtube gnJbPO-GFIw >}}

## Basic Syntax

Golang templates are HTML files with the addition of [variables][variables] and [functions][functions]. Golang template variables and functions are accessible within `{{ }}`.

### Access a Predefined Variable

```
{{ foo }}
```

Parameters for functions are separated using spaces. The following example calls the `add` function with inputs of `1` and `2`:

```
{{ add 1 2 }}
```

#### Methods and Fields are Accessed via dot Notation

Accessing the Page Parameter `bar` defined in a piece of content's [front matter][].

```
{{ .Params.bar }}
```

#### Parentheses Can be Used to Group Items Together

```
{{ if or (isset .Params "alt") (isset .Params "caption") }} Caption {{ end }}
```

## Variables

Each Go template gets a data object. In Hugo, each template is passed a `Page`. See [variables][] for more information.

This is how you access a `Page` variable from a template:

```
<title>{{ .Title }}</title>
```

Values can also be stored in custom variables and referenced later:

```
{{ $address := "123 Main St."}}
{{ $address }}
```

{{% warning %}}
Variables defined inside `if` conditionals and similar are not visible on the outside. See [https://github.com/golang/go/issues/10608](https://github.com/golang/go/issues/10608).

Hugo has created a workaround for this issue in [Scratch](/functions/scratch).

{{% /warning %}}

## Functions

Go templates only ship with a few basic functions but also provide a mechanism for applications to extend the original set.

[Hugo template functions][functions] provide additional functionality specific to building websites. Functions are called by using their name followed by the required parameters separated by spaces. Template functions cannot be added without recompiling Hugo.

### Example 1: Adding Numbers

```
{{ add 1 2 }}
=> 3
```

### Example 2: Comparing Numbers

```
{{ lt 1 2 }}
=> true (i.e., since 1 is less than 2)
```

Note that both examples make use of Go template's [math functions][].

{{% note "Additional Boolean Operators" %}}
There are more boolean operators than those listed in the Hugo docs in the [Golang template documentation](http://golang.org/pkg/text/template/#hdr-Functions).
{{% /note %}}

## Includes

When including another template, you will pass to it the data it will be
able to access. To pass along the current context, please remember to
include a trailing dot. The templates location will always be starting at
the `/layouts/` directory within Hugo.

### Template and Partial Examples

```
{{ template "partials/header.html" . }}
```

Starting with Hugo v0.12, you may also use the `partial` call
for [partial templates][partials]:

```
{{ partial "header.html" . }}
```

## Logic

Go templates provide the most basic iteration and conditional logic.

### Iteration

Just like in Go, the Go templates make heavy use of `range` to iterate over
a map, array, or slice. The following are different examples of how to use
range.

#### Example 1: Using Context

```
{{ range array }}
    {{ . }}
{{ end }}
```

#### Example 2: Declaring Value => Variable name

```
{{range $element := array}}
    {{ $element }}
{{ end }}
```

#### Example 3: Declaring Key-Value Variable Name

```
{{range $index, $element := array}}
   {{ $index }}
   {{ $element }}
{{ end }}
```

### Conditionals

`if`, `else`, `with`, `or`, and `and` provide the framework for handling conditional logic in Go Templates. Like `range`, each statement is closed with an `{{end}}`.

Go Templates treat the following values as false:

* false
* 0
* any zero-length array, slice, map, or string

#### Example 1: `if`

```
{{ if isset .Params "title" }}<h4>{{ index .Params "title" }}</h4>{{ end }}
```

#### Example 2: `if` … `else`

```
{{ if isset .Params "alt" }}
    {{ index .Params "alt" }}
{{else}}
    {{ index .Params "caption" }}
{{ end }}
```

#### Example 3: `and` & `or`

```
{{ if and (or (isset .Params "title") (isset .Params "caption")) (isset .Params "attr")}}
```

#### Example 4: `with`

An alternative way of writing "`if`" and then referencing the same value
is to use "`with`" instead. `with` rebinds the context `.` within its scope
and skips the block if the variable is absent.

The first example above could be simplified as:

```
{{ with .Params.title }}<h4>{{ . }}</h4>{{ end }}
```

#### Example 5: `if` … `else if`

```
{{ if isset .Params "alt" }}
    {{ index .Params "alt" }}
{{ else if isset .Params "caption" }}
    {{ index .Params "caption" }}
{{ end }}
```

## Pipes

One of the most powerful components of Go templates is the ability to stack actions one after another. This is done by using pipes. Borrowed from Unix pipes, the concept is simple: each pipeline's output becomes the input of the following pipe.

Because of the very simple syntax of Go templates, the pipe is essential to being able to chain together function calls. One limitation of the pipes is that they can only work with a single value and that value becomes the last parameter of the next pipeline.

A few simple examples should help convey how to use the pipe.

### Example 1: `shuffle`

The following two examples are functionally the same:

```
{{ shuffle (seq 1 5) }}
```


```
{{ (seq 1 5) | shuffle }}
```

### Example 2: `index`

The following accesses the page parameter called "disqus_url" and escapes the HTML. This example also uses the [`index` function][index], which is built into Go templates:

```
{{ index .Params "disqus_url" | html }}
```

### Example 3: `or` with `isset`

```
{{ if or (or (isset .Params "title") (isset .Params "caption")) (isset .Params "attr") }}
Stuff Here
{{ end }}
```

Could be rewritten as

```
{{ if isset .Params "caption" | or isset .Params "title" | or isset .Params "attr" }}
Stuff Here
{{ end }}
```

### Example 4: Internet Explorer Conditional Comments

By default, Go Templates remove HTML comments from output. This has the unfortunate side effect of removing Internet Explorer conditional comments. As a workaround, use something like this:

```
{{ "<!--[if lt IE 9]>" | safeHTML }}
  <script src="html5shiv.js"></script>
{{ "<![endif]-->" | safeHTML }}
```

Alternatively, you can use the backtick (`` ` ``) to quote the IE conditional comments, avoiding the tedious task of escaping every double quotes (`"`) inside, as demonstrated in the [examples](http://golang.org/pkg/text/template/#hdr-Examples) in the Go text/template documentation:

```
{{ `<!--[if lt IE 7]><html class="no-js lt-ie9 lt-ie8 lt-ie7"><![endif]-->` | safeHTML }}
```

## Context (aka "the dot")

The most easily overlooked concept to understand about Go templates is that `{{ . }}` always refers to the current context. In the top level of your template, this will be the data set made available to it. Inside of an iteration, however, it will have the value of the current item in the loop; i.e., `{{ . }}` will no longer refer to the data available to the entire page. If you need to access page-level data (e.g., page params set in front matter) from within the loop, you will likely want to do one of the following:

### 1. Define a Variable Independent of Context

The following shows how to define a variable independent of the context.

{{< code file="tags-range-with-page-variable.html" >}}
{{ $title := .Site.Title }}
<ul>
{{ range .Params.tags }}
    <li>
        <a href="/tags/{{ . | urlize }}">{{ . }}</a>
        - {{ $title }}
    </li>
{{ end }}
</ul>
{{< /code >}}

{{% note %}}
Notice how once we have entered the loop (i.e. `range`), the value of `{{ . }}` has changed. We have defined a variable outside of the loop (`{{$title}}`) that we've assigned a value so that we have access to the value from within the loop as well.
{{% /note %}}

### 2. Use `$.` to Access the Global Context

`$` has special significance in your templates. `$` is set to the starting value of `.` ("the dot") by default. This is a [documented feature of Go text/template][dotdoc]. This means you have access to the global context from anywhere. Here is an equivalent example of the preceding code block but now using `$` to grab `.Site.Title` from the global context:

{{< code file="range-through-tags-w-global.html" >}}
<ul>
{{ range .Params.tags }}
  <li>
    <a href="/tags/{{ . | urlize }}">{{ . }}</a>
            - {{ $.Site.Title }}
  </li>
{{ end }}
</ul>
{{< /code >}}

{{% warning "Don't Redefine the Dot" %}}
The built-in magic of `$` would cease to work if someone were to mischievously redefine the special character; e.g. `{{ $ := .Site }}`. *Don't do it.* You may, of course, recover from this mischief by using `{{ $ := . }}` in a global context to reset `$` to its default value.
{{% /warning %}}

## Whitespace

Go 1.6 includes the ability to trim the whitespace from either side of a Go tag by including a hyphen (`-`) and space immediately beside the corresponding `{{` or `}}` delimiter.

For instance, the following Go template will include the newlines and horizontal tab in its HTML output:

```
<div>
  {{ .Title }}
</div>
```

Which will output:

```
<div>
  Hello, World!
</div>
```

Leveraging the `-` in the following example will remove the extra white space surrounding the `.Title` variable and remove the newline:

```
<div>
  {{- .Title -}}
</div>
```

Which then outputs:

```
<div>Hello, World!</div>
```

Go considers the following characters whitespace:

* <kbd>space</kbd>
* horizontal <kbd>tab</kbd>
* carriage <kbd>return</kbd>
* newline

## Hugo Parameters

Hugo provides the option of passing values to your template layer through your [site configuration][config] (i.e. for site-wide values) or through the metadata of each specific piece of content (i.e. the [front matter][]). You can define any values of any type and use them however you want in your templates, as long as the values are supported by the front matter format specified via `metaDataFormat` in your configuration file.

## Use Content (`Page`) Parameters

You can provide variables to be used by templates in individual content's [front matter][].

An example of this is used in the Hugo docs. Most of the pages benefit from having the table of contents provided, but sometimes the table of contents doesn't make a lot of sense. We've defined a `notoc` variable in our front matter that will prevent a table of contents from rendering when specifically set to `true`.

Here is the example front matter:

```
---
title: Roadmap
lastmod: 2017-03-05
date: 2013-11-18
notoc: true
---
```

Here is an example of corresponding code that could be used inside a `toc.html` [partial template][partials]:

{{< code file="layouts/partials/toc.html" download="toc.html" >}}
{{ if not .Params.notoc }}
<aside>
  <header>
    <a href="#{{.Title | urlize}}">
    <h3>{{.Title}}</h3>
    </a>
  </header>
  {{.TableOfContents}}
</aside>
<a href="#" id="toc-toggle"></a>
{{end}}
{{< /code >}}

We want the *default* behavior to be for pages to include a TOC unless otherwise specified. This template checks to make sure that the `notoc:` field in this page's front matter is not `true`.

## Use Site Configuration Parameters

You can arbitrarily define as many site-level parameters as you want in your [site's configuration file][config]. These parameters are globally available in your templates.

For instance, you might declare the following:

{{< code file="config.yaml" >}}
params:
  copyrighthtml: "Copyright &#xA9; 2017 John Doe. All Rights Reserved."
  twitteruser: "spf13"
  sidebarrecentlimit: 5
{{< /code >}}

Within a footer layout, you might then declare a `<footer>` that is only rendered if the `copyrighthtml` parameter is provided. If it *is* provided, you will then need to declare the string is safe to use via the [`safeHTML` function][safehtml] so that the HTML entity is not escaped again. This would let you easily update just your top-level config file each January 1st, instead of hunting through your templates.

```
{{if .Site.Params.copyrighthtml}}<footer>
<div class="text-center">{{.Site.Params.CopyrightHTML | safeHTML}}</div>
</footer>{{end}}
```

An alternative way of writing the "`if`" and then referencing the same value is to use [`with`][with] instead. `with` rebinds the context (`.`) within its scope and skips the block if the variable is absent:

{{< code file="layouts/partials/twitter.html" >}}
{{with .Site.Params.twitteruser}}
<div>
  <a href="https://twitter.com/{{.}}" rel="author">
  <img src="/images/twitter.png" width="48" height="48" title="Twitter: {{.}}" alt="Twitter"></a>
</div>
{{end}}
{{< /code >}}

Finally, you can pull "magic constants" out of your layouts as well. The following uses the [`first`][first] function, as well as the [`.RelPermalink`][relpermalink] page variable and the [`.Site.Pages`][sitevars] site variable.

```
<nav>
  <h1>Recent Posts</h1>
  <ul>
  {{- range first .Site.Params.SidebarRecentLimit .Site.Pages -}}
    <li><a href="{{.RelPermalink}}">{{.Title}}</a></li>
  {{- end -}}
  </ul>
</nav>
```

## Example: Show Only Upcoming Events

Go allows you to do more than what's shown here. Using Hugo's [`where` function][where] and Go built-ins, we can list only the items from `content/events/` whose date (set in a content file's [front matter][]) is in the future. The following is an example [partial template][partials]:

{{< code file="layouts/partials/upcoming-events.html" download="upcoming-events.html" >}}
<h4>Upcoming Events</h4>
<ul class="upcoming-events">
{{ range where .Data.Pages.ByDate "Section" "events" }}
  {{ if ge .Date.Unix .Now.Unix }}
    <li>
    <!-- add span for event type -->
      <span>{{ .Type | title }} —</span>
      {{ .Title }} on
    <!-- add span for event date -->
      <span>{{ .Date.Format "2 January at 3:04pm" }}</span>
      at {{ .Params.place }}
    </li>
  {{ end }}
{{ end }}
</ul>
{{< /code >}}


[`where` function]: /functions/where/
[config]: /getting-started/configuration/
[dotdoc]: http://golang.org/pkg/text/template/#hdr-Variables
[first]: /functions/first/
[front matter]: /content-management/front-matter/
[functions]: /functions/ "See the full list of Hugo's templating functions with a quick start reference guide and basic and advanced examples."
[Go html/template]: http://golang.org/pkg/html/template/ "Godocs references for Golang's html templating"
[gohtmltemplate]: http://golang.org/pkg/html/template/ "Godocs references for Golang's html templating"
[index]: /functions/index/
[math functions]: /functions/math/
[partials]: /templates/partials/ "Link to the partial templates page inside of the templating section of the Hugo docs"
[relpermalink]: /variables/page/
[safehtml]: /functions/safehtml/
[sitevars]: /variables/site/
[variables]: /variables/ "See the full extent of page-, site-, and other variables that Hugo make available to you in your templates."
[where]: /functions/where/
[with]: /functions/with/
[godocsindex]: http://golang.org/pkg/text/template/ "Godocs page for index function"
