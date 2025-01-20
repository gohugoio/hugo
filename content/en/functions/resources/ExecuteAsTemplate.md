---
title: resources.ExecuteAsTemplate
description: Returns a resource created from a Go template, parsed and executed with the given context.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/resources/FromString
  returnType: resource.Resource
  signatures: [resources.ExecuteAsTemplate TARGETPATH CONTEXT RESOURCE]
---

The `resources.ExecuteAsTemplate` function returns a resource created from a Go template, parsed and executed with the given context, caching the result using the target path as its cache key.

Hugo publishes the resource to the target path when you call its [`Publish`], [`Permalink`], or [`RelPermalink`] methods.

[`publish`]: /methods/resource/publish/
[`permalink`]: /methods/resource/permalink/
[`relpermalink`]: /methods/resource/relpermalink/

Let's say you have a CSS file that you wish to populate with values from your site configuration:

{{< code file=assets/css/template.css lang=go-html-template >}}
body {
  background-color: {{ site.Params.style.bg_color }};
  color: {{ site.Params.style.text_color }};
}
{{< /code >}}

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

{{< code file=public/css/main.css >}}
body {
  background-color: #fefefe;
  color: #222;
}
{{< /code >}}
