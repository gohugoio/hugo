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

Hugo publishes the resource to the target path when you call its [`Publish`], [`Permalink`], or [`RelPermalink`] methods.

[`publish`]: /methods/resource/publish/
[`permalink`]: /methods/resource/permalink/
[`relpermalink`]: /methods/resource/relpermalink/

Let's say you need to publish a file named "site.json" in the root of your `public` directory, containing the build date, the Hugo version used to build the site, and the date that the content was last modified. For example:

```json
{
  "build_date": "2025-07-08T13:12:19-07:00",
  "hugo_version": "0.148.0",
  "last_modified": "2025-07-07T22:09:13-07:00"
}
```

Place this in your baseof.html template:

```go-html-template
{{ if .IsHome }}
  {{ $rfc3339 := "2006-01-02T15:04:05Z07:00" }}
  {{ $m := dict
    "hugo_version" hugo.Version
    "build_date" (now.Format $rfc3339)
    "last_modified" (site.Lastmod.Format $rfc3339)
  }}
  {{ $json := jsonify $m }}
  {{ $r := resources.FromString "site.json" $json }}
  {{ $r.Publish }}
{{ end }}
```

The example above:

1. Creates a map with the relevant key-value pairs using the [`dict`] function
1. Encodes the map as a JSON string using the [`jsonify`] function
1. Creates a resource from the JSON string using the `resources.FromString` function
1. Publishes the file to the root of the `public` directory using the resource's `.Publish` method

Combine `resources.FromString` with [`resources.ExecuteAsTemplate`] if your string contains template actions. Rewriting the example above:

```go-html-template
{{ if .IsHome }}
  {{ $string := `
    {{ $rfc3339 := "2006-01-02T15:04:05Z07:00" }}
    {{ $m := dict
      "hugo_version" hugo.Version
      "build_date" (now.Format $rfc3339)
      "last_modified" (site.Lastmod.Format $rfc3339)
    }}
    {{ $json := jsonify $m }}
    `
  }}
  {{ $r := resources.FromString "" $string }}
  {{ $r = $r | resources.ExecuteAsTemplate "site.json" . }}
  {{ $r.Publish }}
{{ end }}
```

[`dict`]: /functions/collections/dictionary/
[`jsonify`]: /functions/encoding/jsonify/
[`resources.ExecuteAsTemplate`]: /functions/resources/executeastemplate/
