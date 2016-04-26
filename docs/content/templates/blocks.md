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

Go 1.6 includes a powerful new keyword, `block`. This construct allows you to define the outer shell of your pages in a single "base" template, filling in or overriding portions as necessary.

## Define the base template

Let's define a simple base template, a shell from which all our pages will start. To find a base template, Hugo searches the same paths and file names as it does for [Ace templates]({{< relref "templates/ace.md" >}}), just with files suffixed `.html` rather than `.ace`.

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
