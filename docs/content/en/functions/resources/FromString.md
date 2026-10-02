---
title: resources.FromString
description: Returns a resource created from a string.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: resource.Resource
    signatures: [resources.FromString TARGETPATH STRING]
---

The `resources.FromString` function returns a resource created from a string, caching the result using the target path as its cache key.

For example, to create and publish a `security.txt` file from your site configuration:

```go-html-template {file="layouts/baseof.html"}
{{ $content := printf "Contact: mailto:%s\n" site.Params.email }}
{{ with resources.FromString ".well-known/security.txt" $content }}
  {{ .Publish }}
{{ end }}
```

To publish within a pipeline, use the [`resources.Publish`][] function:

```go-html-template {file="layouts/baseof.html"}
{{ $content := printf "Contact: mailto:%s\n" site.Params.email }}
{{ resources.FromString ".well-known/security.txt" $content | resources.Publish }}
```

The [`Permalink`][] and [`RelPermalink`][] methods also publish the resource.

Combine `resources.FromString` with [`resources.ExecuteAsTemplate`][] if your string contains template actions:

```go-html-template {file="layouts/baseof.html"}
{{ $string := `Contact: mailto:{{ site.Params.email }}
Expires: {{ (now.AddDate 1 0 0).UTC.Format "2006-01-02T15:04:05Z" }}
` }}
{{ $r := resources.FromString "" $string }}
{{ $r = $r | resources.ExecuteAsTemplate ".well-known/security.txt" . }}
{{ $r.Publish }}
```

[`Permalink`]: /methods/resource/permalink/
[`RelPermalink`]: /methods/resource/relpermalink/
[`resources.ExecuteAsTemplate`]: /functions/resources/executeastemplate/
[`resources.Publish`]: /functions/resources/publish/
