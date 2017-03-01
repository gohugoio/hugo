---
title: Base Templates and Blocks
linktitle:
description: The base and block constructs allow you to define the outer shell of your master templates (i.e., the chrome of the page) in a syntax that allows for easy extending and overwriting.
godocref: https://golang.org/pkg/text/template/#example_Template_block
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [blocks,base,fundamentals]
weight: 20
draft: false
aliases: [/templates/blocks/,/templates/base-templates-and-blocks/]
toc: true
wip: true
---

Go 1.6 includes a powerful new keyword, `block`. This construct allows you to define the outer shell of your pages one or more master template(s), filling in or overriding portions as necessary.

## Base Template Lookup Order

This is the order Hugo searches for a base template:

1. `/layouts/<CURRENTPATH>/<TEMPLATENAME>-baseof.html`
2. `/layouts/<CURRENTPATH>/baseof.html`
3. `/layouts/_default/<TEMPLATENAME>-baseof.html`
4. `/layouts/_default/baseof.html`

As an example, with a site using the theme `exampletheme`, when rendering the section list for the section `post`. Hugo picks the `section/post.html` as the template and this template has a `define` section that indicates it needs a base template. This is then the lookup order:

1. `/layouts/section/post-baseof.html`
2. `/themes/<THEME>/layouts/section/post-baseof.html`
3. `/layouts/section/baseof.html`
4. `/themes/<THEME>/layouts/section/baseof.html`
5. `/layouts/_default/post-baseof.html`
6. `/themes/<THEME>/layouts/_default/post-baseof.html`
7. `/layouts/_default/baseof.html`
8. `/themes/<THEME>/layouts/_default/baseof.html`


## Defining the Base Template

The following defines a simple base template at `_default/baseof.html`). As a default template, it is the shell from which all our pages will start unless a more specific `*baseof.html` is defined.

{{% code file="layouts/_default/baseof.html" download="baseof.html" %}}
```html
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>{{ block "title" . }}
      <!-- Blocks may include default content. -->
      {{ .Site.Title }}
    {{ end }}</title>
  </head>
  <body>
    <!-- Code that all your templates share, like a header -->

    {{ block "main" . }}
      <!-- The part of the page that begins to differ between templates -->
    {{ end }}

    <!-- More shared code, perhaps a footer -->
  </body>
</html>
```
{{% /code %}}

## Overriding the Base Template

From the above base template, you can define a [default list template][hugolists]. The default list template will inherit all of the code defined above and can then implement its own `"main"` block from:

{{% code file="layouts/_default/list.html" download="list.html" %}}
```html
{{ define "main" }}
  <h1>Posts</h1>
  {{ range .Data.Pages }}
    <article>
      <h2>{{ .Title }}</h2>
      {{ .Content }}
    </article>
  {{ end }}
{{ end }}
```
{{% /code %}}

{{% note "No Go Context \"Dot\" in Block Definitions" %}}
When using the `define` keyword, you do *not* need to use Go templates context reference (i.e., 'The Dot"). (Read more on ["The Dot" in the Go Template Primer](/templates/go-templates/).)
{{% /note %}}

This replaces the contents of our (basically empty) "main" block with something useful for the list template. In this case, we didn't define a `"title"`` block, so the contents from our base template remain unchanged in lists.

{{% warning %}}
Code that you put outside the block definitions *can* break your layout. This even includes HTML comments. For example:

```html
<!-- Harmless comment..that will break your layout at build -->
{{ define "main" }}
...your code here
{{ end }}
```
[See this thread from the discussion forums](https://discuss.gohugo.io/t/baseof-html-block-templates-and-list-types-results-in-empty-pages/5612/6)
{{% /warning %}}

The following shows how you can override both the `"main"` and `"title"` block areas from the base template with code unique to your [default single page template][singletemplate]:

{{% code file="layouts/_default/single.html" download="single.html" %}}
```html
{{ define "title" }}
  <!-- This will override the default value set in baseof.html; i.e., "{{.Site.Title}}" in the original example-->
  {{ .Title }} &ndash; {{ .Site.Title }}
{{ end }}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
{{ end }}
```
{{% /code %}}

[hugolists]: /templates/lists
[singletemplate]: /templates/single-page-templates/