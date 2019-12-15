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
aliases: [/layouts/introduction/,/layout/introduction/, /templates/go-templates/]
toc: true
---

{{% note %}}
The following is only a primer on Go Templates. For an in-depth look into Go Templates, check the official [Go docs](https://golang.org/pkg/text/template/).
{{% /note %}}

Go Templates provide an extremely simple template language that adheres to the belief that only the most basic of logic belongs in the template or view layer.

## Basic Syntax

Go Templates are HTML files with the addition of [variables][variables] and [functions][functions]. Go Template variables and functions are accessible within `{{ }}`.

### Access a Predefined Variable

A _predefined variable_ could be a variable already existing in the
current scope (like the `.Title` example in the [Variables]({{< relref
"#variables" >}}) section below) or a custom variable (like the
`$address` example in that same section).


```go-html-template
{{ .Title }}
{{ $address }}
```

Parameters for functions are separated using spaces. The general syntax is:

```
{{ FUNCTION ARG1 ARG2 .. }}
```

The following example calls the `add` function with inputs of `1` and `2`:

```go-html-template
{{ add 1 2 }}
```

#### Methods and Fields are Accessed via dot Notation

Accessing the Page Parameter `bar` defined in a piece of content's [front matter][].

```go-html-template
{{ .Params.bar }}
```

#### Parentheses Can be Used to Group Items Together

```go-html-template
{{ if or (isset .Params "alt") (isset .Params "caption") }} Caption {{ end }}
```

## Variables {#variables}

Each Go Template gets a data object. In Hugo, each template is passed
a `Page`.  In the below example, `.Title` is one of the elements
accessible in that [`Page` variable][pagevars].

With the `Page` being the default scope of a template, the `Title`
element in current scope (`.` -- "the **dot**") is accessible simply
by the dot-prefix (`.Title`):

```go-html-template
<title>{{ .Title }}</title>
```

Values can also be stored in custom variables and referenced later:

{{% note %}}
The custom variables need to be prefixed with `$`.
{{% /note %}}

```go-html-template
{{ $address := "123 Main St." }}
{{ $address }}
```

{{% warning %}}
For Hugo v0.47 and older versions, variables defined inside `if`
conditionals and similar are not visible on the outside.
See [https://github.com/golang/go/issues/10608](https://github.com/golang/go/issues/10608).

Hugo has created a workaround for this issue in [Scratch](/functions/scratch).
{{% /warning %}}

For **Hugo v0.48** and newer, variables can be re-defined using the
new `=` operator (new in Go 1.11).

Below example will work only in these newer Hugo versions. The example
prints "Var is Hugo Home" on the home page, and "Var is Hugo Page" on
all other pages:

```go-html-template
{{ $var := "Hugo Page" }}
{{ if .IsHome }}
    {{ $var = "Hugo Home" }}
{{ end }}
Var is {{ $var }}
```

## Functions

Go Templates only ship with a few basic functions but also provide a mechanism for applications to extend the original set.

[Hugo template functions][functions] provide additional functionality specific to building websites. Functions are called by using their name followed by the required parameters separated by spaces. Template functions cannot be added without recompiling Hugo.

### Example 1: Adding Numbers

```go-html-template
{{ add 1 2 }}
<!-- prints 3 -->
```

### Example 2: Comparing Numbers

```go-html-template
{{ lt 1 2 }}
<!-- prints true (i.e., since 1 is less than 2) -->
```

Note that both examples make use of Go Template's [math functions][].

{{% note "Additional Boolean Operators" %}}
There are more boolean operators than those listed in the Hugo docs in the [Go Template documentation](https://golang.org/pkg/text/template/#hdr-Functions).
{{% /note %}}

## Includes

When including another template, you will need to pass it the data that it would
need to access.

{{% note %}}
To pass along the current context, please remember to include a trailing **dot**.
{{% /note %}}

The templates location will always be starting at the `layouts/` directory
within Hugo.

### Partial

The [`partial`][partials] function is used to include *partial* templates using
the syntax `{{ partial "<PATH>/<PARTIAL>.<EXTENSION>" . }}`.

Example of including a `layouts/partials/header.html` partial:

```go-html-template
{{ partial "header.html" . }}
```

### Template

The `template` function was used to include *partial* templates
in much older Hugo versions. Now it's useful only for calling
[*internal* templates][internal_templates]. The syntax is `{{ template
"_internal/<TEMPLATE>.<EXTENSION>" . }}`.

{{% note %}}
The available **internal** templates can be found
[here](https://github.com/gohugoio/hugo/tree/master/tpl/tplimpl/embedded/templates).
{{% /note %}}

Example of including the internal `opengraph.html` template:

```go-html-template
{{ template "_internal/opengraph.html" . }}
```

## Logic

Go Templates provide the most basic iteration and conditional logic.

### Iteration

The Go Templates make heavy use of `range` to iterate over a _map_,
_array_, or _slice_. The following are different examples of how to
use `range`.

#### Example 1: Using Context (`.`)

```go-html-template
{{ range $array }}
    {{ . }} <!-- The . represents an element in $array -->
{{ end }}
```

#### Example 2: Declaring a variable name for an array element's value

```go-html-template
{{ range $elem_val := $array }}
    {{ $elem_val }}
{{ end }}
```

#### Example 3: Declaring variable names for an array element's index _and_ value

For an array or slice, the first declared variable will map to each
element's index.

```go-html-template
{{ range $elem_index, $elem_val := $array }}
   {{ $elem_index }} -- {{ $elem_val }}
{{ end }}
```

#### Example 4: Declaring variable names for a map element's key _and_ value

For a map, the first declared variable will map to each map element's
key.

```go-html-template
{{ range $elem_key, $elem_val := $map }}
   {{ $elem_key }} -- {{ $elem_val }}
{{ end }}
```

#### Example 5: Conditional on empty _map_, _array_, or _slice_.

If the _map_, _array_, or _slice_ passed into the range is zero-length then the else statement is evaluated.

```go-html-template
{{ range $array }}
    {{ . }}
{{else}}
    <!-- This is only evaluated if $array is empty -->
{{ end }}
```

### Conditionals

`if`, `else`, `with`, `or`, and `and` provide the framework for handling conditional logic in Go Templates. Like `range`, each statement is closed with an `{{ end }}`.

Go Templates treat the following values as **false**:

- `false` (boolean)
- 0 (integer)
- any zero-length array, slice, map, or string

#### Example 1: `with`

It is common to write "if something exists, do this" kind of
statements using `with`.

{{% note %}}
`with` rebinds the context `.` within its scope (just like in `range`).
{{% /note %}}

It skips the block if the variable is absent, or if it evaluates to
"false" as explained above.

```go-html-template
{{ with .Params.title }}
    <h4>{{ . }}</h4>
{{ end }}
```

#### Example 2: `with` .. `else`

Below snippet uses the "description" front-matter parameter's value if
set, else uses the default `.Summary` [Page variable][pagevars]:


```go-html-template
{{ with .Param "description" }}
    {{ . }}
{{ else }}
    {{ .Summary }}
{{ end }}
```

See the [`.Param` function][param].

#### Example 3: `if`

An alternative (and a more verbose) way of writing `with` is using
`if`. Here, the `.` does not get rebinded.

Below example is "Example 1" rewritten using `if`:

```go-html-template
{{ if isset .Params "title" }}
    <h4>{{ index .Params "title" }}</h4>
{{ end }}
```

#### Example 4: `if` .. `else`

Below example is "Example 2" rewritten using `if` .. `else`, and using
[`isset` function][isset] + `.Params` variable (different from the
[`.Param` **function**][param]) instead:

```go-html-template
{{ if (isset .Params "description") }}
    {{ index .Params "description" }}
{{ else }}
    {{ .Summary }}
{{ end }}
```

#### Example 5: `if` .. `else if` .. `else`

Unlike `with`, `if` can contain `else if` clauses too.

```go-html-template
{{ if (isset .Params "description") }}
    {{ index .Params "description" }}
{{ else if (isset .Params "summary") }}
    {{ index .Params "summary" }}
{{ else }}
    {{ .Summary }}
{{ end }}
```

#### Example 6: `and` & `or`

```go-html-template
{{ if (and (or (isset .Params "title") (isset .Params "caption")) (isset .Params "attr")) }}
```

## Pipes

One of the most powerful components of Go Templates is the ability to stack actions one after another. This is done by using pipes. Borrowed from Unix pipes, the concept is simple: each pipeline's output becomes the input of the following pipe.

Because of the very simple syntax of Go Templates, the pipe is essential to being able to chain together function calls. One limitation of the pipes is that they can only work with a single value and that value becomes the last parameter of the next pipeline.

A few simple examples should help convey how to use the pipe.

### Example 1: `shuffle`

The following two examples are functionally the same:

```go-html-template
{{ shuffle (seq 1 5) }}
```


```go-html-template
{{ (seq 1 5) | shuffle }}
```

### Example 2: `index`

The following accesses the page parameter called "disqus_url" and escapes the HTML. This example also uses the [`index` function][index], which is built into Go Templates:

```go-html-template
{{ index .Params "disqus_url" | html }}
```

### Example 3: `or` with `isset`

```go-html-template
{{ if or (or (isset .Params "title") (isset .Params "caption")) (isset .Params "attr") }}
Stuff Here
{{ end }}
```

Could be rewritten as

```go-html-template
{{ if isset .Params "caption" | or isset .Params "title" | or isset .Params "attr" }}
Stuff Here
{{ end }}
```

### Example 4: Internet Explorer Conditional Comments {#ie-conditional-comments}

By default, Go Templates remove HTML comments from output. This has the unfortunate side effect of removing Internet Explorer conditional comments. As a workaround, use something like this:

```go-html-template
{{ "<!--[if lt IE 9]>" | safeHTML }}
  <script src="html5shiv.js"></script>
{{ "<![endif]-->" | safeHTML }}
```

Alternatively, you can use the backtick (`` ` ``) to quote the IE conditional comments, avoiding the tedious task of escaping every double quotes (`"`) inside, as demonstrated in the [examples](https://golang.org/pkg/text/template/#hdr-Examples) in the Go text/template documentation:

```go-html-template
{{ `<!--[if lt IE 7]><html class="no-js lt-ie9 lt-ie8 lt-ie7"><![endif]-->` | safeHTML }}
```

## Context (aka "the dot") {#the-dot}

The most easily overlooked concept to understand about Go Templates is
that `{{ . }}` always refers to the **current context**.

- In the top level of your template, this will be the data set made
  available to it.
- Inside of an iteration, however, it will have the value of the
  current item in the loop; i.e., `{{ . }}` will no longer refer to
  the data available to the entire page.

If you need to access page-level data (e.g., page params set in front
matter) from within the loop, you will likely want to do one of the
following:

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

For instance, the following Go Template will include the newlines and horizontal tab in its HTML output:

```go-html-template
<div>
  {{ .Title }}
</div>
```

Which will output:

```html
<div>
  Hello, World!
</div>
```

Leveraging the `-` in the following example will remove the extra white space surrounding the `.Title` variable and remove the newline:

```go-html-template
<div>
  {{- .Title -}}
</div>
```

Which then outputs:

```html
<div>Hello, World!</div>
```

Go considers the following characters _whitespace_:

* <kbd>space</kbd>
* horizontal <kbd>tab</kbd>
* carriage <kbd>return</kbd>
* newline

## Comments

In order to keep your templates organized and share information throughout your team, you may want to add comments to your templates. There are two ways to do that with Hugo.

### Go Templates comments

Go Templates support `{{/*` and `*/}}` to open and close a comment block. Nothing within that block will be rendered.

For example:

```go-html-template
Bonsoir, {{/* {{ add 0 + 2 }} */}}Eliott.
```

Will render `Bonsoir, Eliott.`, and not care about the syntax error (`add 0 + 2`) in the comment block.

### HTML comments

If you need to produce HTML comments from your templates, take a look at the [Internet Explorer conditional comments](#ie-conditional-comments) example. If you need variables to construct such HTML comments, just pipe `printf` to `safeHTML`. For example:

```go-html-template
{{ printf "<!-- Our website is named: %s -->" .Site.Title | safeHTML }}
```

#### HTML comments containing Go Templates

HTML comments are by default stripped, but their content is still evaluated. That means that although the HTML comment will never render any content to the final HTML pages, code contained within the comment may fail the build process.

{{% note %}}
Do **not** try to comment out Go Template code using HTML comments.
{{% /note %}}

```go-html-template
<!-- {{ $author := "Emma Goldman" }} was a great woman. -->
{{ $author }}
```

The templating engine will strip the content within the HTML comment, but will first evaluate any Go Template code if present within. So the above example will render `Emma Goldman`, as the `$author` variable got evaluated in the HTML comment. But the build would have failed if that code in the HTML comment had an error.

## Hugo Parameters

Hugo provides the option of passing values to your template layer through your [site configuration][config] (i.e. for site-wide values) or through the metadata of each specific piece of content (i.e. the [front matter][]). You can define any values of any type and use them however you want in your templates, as long as the values are supported by the [front matter format]({{< ref "front-matter.md#front-matter-formats" >}}).

## Use Content (`Page`) Parameters

You can provide variables to be used by templates in individual content's [front matter][].

An example of this is used in the Hugo docs. Most of the pages benefit from having the table of contents provided, but sometimes the table of contents doesn't make a lot of sense. We've defined a `notoc` variable in our front matter that will prevent a table of contents from rendering when specifically set to `true`.

Here is the example front matter (YAML):

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
{{ end }}
{{< /code >}}

We want the *default* behavior to be for pages to include a TOC unless otherwise specified. This template checks to make sure that the `notoc:` field in this page's front matter is not `true`.

## Use Site Configuration Parameters

You can arbitrarily define as many site-level parameters as you want in your [site's configuration file][config]. These parameters are globally available in your templates.

For instance, you might declare the following:

{{< code-toggle file="config" >}}
params:
  copyrighthtml: "Copyright &#xA9; 2017 John Doe. All Rights Reserved."
  twitteruser: "spf13"
  sidebarrecentlimit: 5
{{< /code >}}

Within a footer layout, you might then declare a `<footer>` that is only rendered if the `copyrighthtml` parameter is provided. If it *is* provided, you will then need to declare the string is safe to use via the [`safeHTML` function][safehtml] so that the HTML entity is not escaped again. This would let you easily update just your top-level config file each January 1st, instead of hunting through your templates.

```go-html-template
{{ if .Site.Params.copyrighthtml }}
    <footer>
        <div class="text-center">{{.Site.Params.CopyrightHTML | safeHTML}}</div>
    </footer>
{{ end }}
```

An alternative way of writing the "`if`" and then referencing the same value is to use [`with`][with] instead. `with` rebinds the context (`.`) within its scope and skips the block if the variable is absent:

{{< code file="layouts/partials/twitter.html" >}}
{{ with .Site.Params.twitteruser }}
    <div>
        <a href="https://twitter.com/{{.}}" rel="author">
        <img src="/images/twitter.png" width="48" height="48" title="Twitter: {{.}}" alt="Twitter"></a>
    </div>
{{ end }}
{{< /code >}}

Finally, you can pull "magic constants" out of your layouts as well. The following uses the [`first`][first] function, as well as the [`.RelPermalink`][relpermalink] page variable and the [`.Site.Pages`][sitevars] site variable.

```go-html-template
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
{{ range where .Pages.ByDate "Section" "events" }}
    {{ if ge .Date.Unix now.Unix }}
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
[dotdoc]: https://golang.org/pkg/text/template/#hdr-Variables
[first]: /functions/first/
[front matter]: /content-management/front-matter/
[functions]: /functions/ "See the full list of Hugo's templating functions with a quick start reference guide and basic and advanced examples."
[Go html/template]: https://golang.org/pkg/html/template/ "Godocs references for Go's html templating"
[gohtmltemplate]: https://golang.org/pkg/html/template/ "Godocs references for Go's html templating"
[index]: /functions/index/
[math functions]: /functions/math/
[partials]: /templates/partials/ "Link to the partial templates page inside of the templating section of the Hugo docs"
[internal_templates]: /templates/internal/
[relpermalink]: /variables/page/
[safehtml]: /functions/safehtml/
[sitevars]: /variables/site/
[pagevars]: /variables/page/
[variables]: /variables/ "See the full extent of page-, site-, and other variables that Hugo make available to you in your templates."
[where]: /functions/where/
[with]: /functions/with/
[godocsindex]: https://golang.org/pkg/text/template/ "Godocs page for index function"
[param]: /functions/param/
[isset]: /functions/isset/
