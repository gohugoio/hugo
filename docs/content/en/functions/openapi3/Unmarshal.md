---
title: openapi3.Unmarshal
description: Unmarshals the given resource into an OpenAPI 3 Description.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: openapi3.OpenAPIDocument
    signatures: ['openapi3.Unmarshal RESOURCE [OPTIONS]']
---

The resource passed to the `openapi3.Unmarshal` function must be an [OpenAPI Document][], typically in JSON or YAML format. This resource can be a [global resource](g) or a [remote resource](g).

This function automatically resolves and includes all external references, both local and remote, and returns a complete [OpenAPI Description][] that fully describes the surface of an API and its semantics.

## Options

The `openapi3.Unmarshal` function accepts an options map.

`getremote`
: {{< new-in 0.153.0 />}}
: (`map`) This is a map of the options for the [`resources.GetRemote`][] function, useful when an OpenAPI Document includes remote external references.

## Examples

### Remote resource

To work with a remote resource:

```go-html-template {copy=true}
{{ $api := "" }}
{{ $url := "https://petstore.swagger.io/v2/swagger.json" }}
{{ $opts := dict
  "headers" (dict "Authorization" "Bearer abcd")
}}
{{ with try (resources.GetRemote $url $opts) }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else with .Value }}
    {{ $api = openapi3.Unmarshal . (dict "getremote" $opts) }}
  {{ else }}
    {{ errorf "Unable to get remote resource %q" $url }}
  {{ end }}
{{ end }}
```

In the example above, the same HTTP Authorization header is used for both the initial remote request made by the `resources.GetRemote` function and for subsequent requests by the `openapi.Unmarshal` function as it retrieve remote external references.

### Global resource

To work with a global resource:

```go-html-template {copy=true}
{{ $api := "" }}
{{ $opts := dict
  "method" "post"
  "key" now.UnixNano
}}
{{ with resources.Get "api/petstore.json" }}
  {{ $api = openapi3.Unmarshal . (dict "getremote" $opts) }}
{{ end }}
```

For global resources, local external reference paths starting with `/` are resolved relative to the `assets` directory. All other local paths are resolved relative to the entry point. In the example above, local paths are resolved relative to `assets/api/petstore.json`.

## Inspection

> [!NOTE]
> The unmarshaled data structure is created with [`kin-openapi`][]. Many fields are structs or pointers (not maps), and therefore require accessors or other methods for indexing and iteration.
>
> For example, `Paths` is a pointer rather than a map; to iterate over the API paths, you must use the `.Paths.Map` accessor as shown in the example below.
>
> See the [`kin-openapi` godoc for OpenAPI 3][] for full type definitions.

To inspect the unmarshaled data structure:

```go-html-template {copy=true}
<pre>{{ debug.Dump $api }}</pre>
```

To list the GET and POST operations for each of the API paths:

```go-html-template {copy=true}
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

[OpenAPI Description]: https://swagger.io/specification/#openapi-description
[OpenAPI Document]: https://swagger.io/specification/#openapi-document
[`kin-openapi` godoc for OpenAPI 3]: https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3
[`kin-openapi`]: https://github.com/getkin/kin-openapi
[`resources.GetRemote`]: /functions/resources/getremote/#options
