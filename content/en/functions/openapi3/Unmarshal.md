---
title: openapi3.Unmarshal
description: Unmarshals the given resource into an OpenAPI 3 document.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: openapi3.OpenAPIDocument
    signatures: ['openapi3.Unmarshal RESOURCE']
---

Use the `openapi3.Unmarshal` function with [global resources](g), [page resources](g), or [remote resources](g).

[OpenAPI]: https://www.openapis.org/

For example, to work with a remote [OpenAPI] definition:

```go-html-template
{{ $url := "https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/examples/v3.0/petstore.json" }}
{{ $api := "" }}
{{ with try (resources.GetRemote $url) }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else with .Value }}
    {{ $api = . | openapi3.Unmarshal }}
  {{ else }}
    {{ errorf "Unable to get remote resource %q" $url }}
  {{ end }}
{{ end }}
```

To inspect the data structure:

```go-html-template
<pre>{{ debug.Dump $api }}</pre>
```

To list the GET and POST operations for each of the API paths:

```go-html-template
{{ range $path, $details := $api.Paths.Map }}
  <p>{{ $path }}</p>
  <dl>
    {{ with $details.Get }}
      <dt>GET</dt>
      <dd>{{ .Summary }}</dd>
    {{ end }}
    {{ with $details.Post }}
      <dt>POST</dt>
      <dd>{{ .Summary }}</dd>
    {{ end }}
  </dl>
{{ end }}
```

> [!warning]
> The unmarshaled data structure is created with [`kin-openapi`](https://github.com/getkin/kin-openapi). Many fields are structs or pointers (not maps), and therefore require accessors or other methods for indexing and iteration. 
> For example, prior to [`kin-openapi` v0.122.0](https://github.com/getkin/kin-openapi#v01220) / [Hugo v0.121.0](https://github.com/gohugoio/hugo/releases/tag/v0.121.0), `Paths` was a map (so `.Paths` was iterable) and it is now a pointer (and requires the `.Paths.Map` accessor, as in the example above).
> See the [`kin-openapi` godoc for OpenAPI 3](https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3) for full type definitions.

Hugo renders this to:

```html
<p>/pets</p>
<dl>
  <dt>GET</dt>
  <dd>List all pets</dd>
  <dt>POST</dt>
  <dd>Create a pet</dd>
</dl>
<p>/pets/{petId}</p>
<dl>
  <dt>GET</dt>
  <dd>Info for a specific pet</dd>
</dl>
```
