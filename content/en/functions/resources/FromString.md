---
title: resources.FromString
description: Creates a resource from a string.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/resources/ExecuteAsTemplate
  returnType: resource.Resource
  signatures: [resources.FromString TARGETPATH STRING]
---

Hugo publishes the resource to the target path when you call its`.Publish`, `.Permalink`, or `.RelPermalink` method. The resource is cached, using the target path as the cache key.

Let's say you need to publish a file named "site.json" in the root of your public directory, containing the build date, the Hugo version used to build the site, and the date that the content was last modified. For example:

```json
{
  "build_date": "2023-10-03T10:50:40-07:00",
  "hugo_version": "0.120.0",
  "last_modified": "2023-10-02T15:21:27-07:00"
}
```

Place this in your baseof.html template:

```go-html-template
{{ if .IsHome }}
  {{ $rfc3339 := "2006-01-02T15:04:05Z07:00" }}
  {{ $m := dict
    "hugo_version" hugo.Version
    "build_date" (now.Format $rfc3339)
    "last_modified" (site.LastChange.Format $rfc3339)
  }}
  {{ $json := jsonify $m }}
  {{ $r := resources.FromString "site.json" $json }}
  {{ $r.Publish }}
{{ end }}
```

The example above:

1. Creates a map with the relevant key/value pairs using the [`dict`] function
2. Encodes the map as a JSON string using the [`jsonify`] function
3. Creates a resource from the JSON string using the `resources.FromString` function
4. Publishes the file to the root of the public directory using the resource's `.Publish` method

Combine `resources.FromString` with [`resources.ExecuteAsTemplate`] if your string contains template actions. Rewriting the example above:

```go-html-template
{{ if .IsHome }}
  {{ $string := `
    {{ $rfc3339 := "2006-01-02T15:04:05Z07:00" }}
    {{ $m := dict
      "hugo_version" hugo.Version
      "build_date" (now.Format $rfc3339)
      "last_modified" (site.LastChange.Format $rfc3339)
    }}
    {{ $json := jsonify $m }}
    `
  }}
  {{ $r := resources.FromString "" $string }}
  {{ $r = $r | resources.ExecuteAsTemplate "site.json" . }}
  {{ $r.Publish }}
{{ end }}
```

[`dict`]: /functions/collections/dictionary
[`jsonify`]: /functions/encoding/jsonify
[`resources.ExecuteAsTemplate`]: /functions/resources/executeastemplate
