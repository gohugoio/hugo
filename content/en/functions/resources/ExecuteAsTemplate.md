---
title: resources.ExecuteAsTemplate
description: Creates a resource from a Go template, parsed and executed with the given context.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/resources/FromString
  returnType: resource.Resource
  signatures: [resources.ExecuteAsTemplate TARGETPATH CONTEXT RESOURCE]
---

Hugo publishes the resource to the target path when you call its`.Publish`, `.Permalink`, or `.RelPermalink` method. The resource is cached, using the target path as the cache key.

Let's say you have a CSS file that you wish to populate with values from your site configuration:

{{< code lang=go-html-template file=assets/css/template.css >}}
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
2. Executes the resource as a template, passing the current page in context
3. Publishes the resource to css/main.css

The result is:

{{< code file="public/css/main.css" >}}
body {
  background-color: #fefefe;
  color: #222;
}
{{< /code >}}
