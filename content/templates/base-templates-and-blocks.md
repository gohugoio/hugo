---
title: Base Templates and Blocks
linktitle:
description:
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [blocks,base,fundamentals]
weight: 20
draft: false
aliases: [/templates/blocks/]
toc: true
needsreview: true
---

Go 1.6 includes a powerful new keyword, `block`. This construct allows you to define the outer shell of your pages one or more master template(s), filling in or overriding portions as necessary.

## Base Template Lookup Order

This is the order Hugo searches for a base template:

1. /layouts/_current-path_/_template-name_-baseof.html, e.g. list-baseof.html.
2. /layouts/_current-path_/baseof.html
3. /layouts/_default/_template-name_-baseof.html e.g. list-baseof.html.
4. /layouts/_default/baseof.html

For each of the steps above, it will first look in the project, then, if theme is set, in the theme's layouts folder. Hugo picks the first base template found.

As an example, with a site using the theme `exampletheme`, when rendering the section list for the section `post`. Hugo picks the `section/post.html` as the template and this template has a `define` section that indicates it needs a base template. This is then the lookup order:

1. `/layouts/section/post-baseof.html`
2.  `/themes/exampletheme/layouts/section/post-baseof.html`
3.  `/layouts/section/baseof.html`
4. `/themes/exampletheme/layouts/section/baseof.html`
5.  `/layouts/_default/post-baseof.html`
6.  `/themes/exampletheme/layouts/_default/post-baseof.html`
7.   `/layouts/_default/baseof.html`
8. `/themes/exampletheme/layouts/_default/baseof.html`


## Define the base template

Let's define a simple base template (`_default/baseof.html`), a shell from which all our pages will start.

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

## Overriding the Base Template

Your [default list template](/templates/list/)---`_default/list.html`---will inherit all of the code defined in the base template. It could then implement its own "main" block from the base template above like so:

```html
<!-- Note the lack of Go's context "dot" when defining blocks -->
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

This replaces the contents of our (basically empty) "main" block with something useful for the list template. In this case, we didn't define a "title" block so the contents from our base template remain unchanged in lists.

In our [default single template](/templates/content/)---`_default/single.html`---let's implement both blocks:

```html
{{ define "title" }}
  {{ .Title }} &ndash; {{ .Site.Title }}
{{ end }}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
{{ end }}
```

This overrides both block areas from the base template with code unique to our single template.

