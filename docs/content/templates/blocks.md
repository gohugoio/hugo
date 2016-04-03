---
date: 2016-03-29T21:26:20-05:00
menu:
  main:
    parent: x
title: Blocks
weight: 5
---

Go 1.6 includes a powerful new keyword, `block`, which provides templates with a nice way to share common code and override only the parts that need changing. This construct allows you to define the outer shell of your pages in a single base template.

## Base

Let's define a very simple base template:

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

Your [default list template](/templates/list) (`_default/list.html`) could implement its own _main_ block from the base template above like so:

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

This replaces the contents of our (basically empty) _main_ block with something useful for the list template.

Since we didn't define a _title_ block, the contents from our base template remain unchanged.

In our [the default single template](/templates/content) (`_default/single.html`), let's implement both block:

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

To find a base Go template, Hugo searches the same paths and file names as it does for [Ace templates](/templates/ace), just with files with `.html` rather than `.ace`.
