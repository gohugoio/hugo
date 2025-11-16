---
title: resources.ExecuteAsTemplate
description: Returns a resource created from a Go template, parsed and executed with the given context.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: resource.Resource
    signatures: [resources.ExecuteAsTemplate TARGETPATH CONTEXT RESOURCE]
---

The `resources.ExecuteAsTemplate` function returns a resource created from a Go template, parsed and executed with the given context, caching the result using the target path as its cache key.

Hugo publishes the resource to the target path when you call its [`Publish`], [`Permalink`], or [`RelPermalink`] methods.

Let's say you have a CSS file that you wish to populate with values from your site configuration:

```go-html-template {file="assets/css/template.css"}
body {
  background-color: {{ site.Params.style.bg_color }};
  color: {{ site.Params.style.text_color }};
}
```

And your site configuration contains:

{{< code-toggle file=hugo >}}
[params.style]
bg_color = '#fefefe'
text_color = '#222'
{{< /code-toggle >}}

Place this in your baseof.html template:

```go-html-template
{{ with resources.Get "css/template.css" }}
  {{ with resources.ExecuteAsTemplate "css/main.css" $ . }}
    <link rel="stylesheet" href="{{ .RelPermalink }}">
  {{ end }}
{{ end }}
```

The example above:

1. Captures the template as a resource
1. Executes the resource as a template, passing the current page in context
1. Publishes the resource to css/main.css

The result is:

```css {file="public/css/main.css"}
body {
  background-color: #fefefe;
  color: #222;
}
```

[`publish`]: /methods/resource/publish/
[`permalink`]: /methods/resource/permalink/
[`relpermalink`]: /methods/resource/relpermalink/
