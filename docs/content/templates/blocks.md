---
date: 2016-03-29T21:26:20-05:00
menu:
  main:
    parent: layout
prev: /templates/views/
next: /templates/partials/
title: Block Templates
weight: 80
---

The `block` keyword in Go templates allows you to define the outer shell of your pages one or more master template(s), filling in or overriding portions as necessary.

## Base template lookup

In version `0.20` Hugo introduced custom [Output Formats]({{< relref "extras/output-formats.md" >}}), all of which can have their own templates that also can use a base template if needed.

This introduced two new terms relevant in the lookup of the templates, the media type's `Suffix` and the output format's `Name`. 

Given the above, Hugo tries to use the most specific base tamplate it finds:

1. /layouts/_current-path_/_template-name_-baseof.[output-format].[suffix], e.g. list-baseof.amp.html.
1. /layouts/_current-path_/_template-name_-baseof.[suffix], e.g. list-baseof.html.
2. /layouts/_current-path_/baseof.[output-format].[suffix], e.g baseof.amp.html
2. /layouts/_current-path_/baseof.[suffix], e.g baseof.html
3. /layouts/_default/_template-name_-baseof.[output-format].[suffix] e.g. list-baseof.amp.html.
3. /layouts/_default/_template-name_-baseof.[suffix], e.g. list-baseof.html.
4. /layouts/_default/baseof.[output-format].[suffix]
4. /layouts/_default/baseof.[suffix]

For each of the steps above, it will first look in the project, then, if theme is set, in the theme's layouts folder. Hugo picks the first base template found.

As an example, with a site using the theme `exampletheme`, when rendering the section list for the section `post` for the output format `Calendar`. Hugo picks the `section/post.calendar.ics` as the template and this template has a `define` section that indicates it needs a base template. This is then the lookup order:

1. `/layouts/section/post-baseof.calendar.ics`
1. `/layouts/section/post-baseof.ics`
2.  `/themes/exampletheme/layouts/section/post-baseof.calendar.ics`
2.  `/themes/exampletheme/layouts/section/post-baseof.ics`
3.  `/layouts/section/baseof.calendar.ics`
3.  `/layouts/section/baseof.ics`
4. `/themes/exampletheme/layouts/section/baseof.calendar.ics`
4. `/themes/exampletheme/layouts/section/baseof.ics`
5.  `/layouts/_default/post-baseof.calendar.ics`
5.  `/layouts/_default/post-baseof.ics`
6.  `/themes/exampletheme/layouts/_default/post-baseof.calendar.ics`
6.  `/themes/exampletheme/layouts/_default/post-baseof.ics`
7.   `/layouts/_default/baseof.calendar.ics`
7.   `/layouts/_default/baseof.ics`
8. `/themes/exampletheme/layouts/_default/baseof.calendar.ics`
8. `/themes/exampletheme/layouts/_default/baseof.ics`


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

## Overriding the base

Your [default list template]({{< relref "templates/list.md" >}}) (`_default/list.html`) will inherit all of the code defined in the base template. It could then implement its own "main" block from the base template above like so:

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

In our [default single template]({{< relref "templates/content.md" >}}) (`_default/single.html`), let's implement both blocks:

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
