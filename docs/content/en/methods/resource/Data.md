---
title: Data
description: Returns supplemental information about resources produced by the resources.GetRemote, css.Build, and js.Build functions.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: map
    signatures: [RESOURCE.Data]
---

The `Data` method returns supplemental information about resources produced by the [`resources.GetRemote`][], [`css.Build`][], and [`js.Build`][] functions.

## Example

```go-html-template
{{ $url := "https://example.org/images/a.jpg" }}
{{ $opts := dict "responseHeaders" (slice "Server") }}
{{ with try (resources.GetRemote $url) }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else with .Value }}
    {{ with .Data }}
      {{ .ContentLength }} → 42764
      {{ .ContentType }} → image/jpeg
      {{ .Headers }} → map[Server:[Netlify]]
      {{ .Status }} → 200 OK
      {{ .StatusCode }} → 200
      {{ .TransferEncoding }} → []
    {{ end }}
  {{ else }}
    {{ errorf "Unable to get remote resource %q" $url }}
  {{ end }}
{{ end }}
```

## Methods

Use these methods on the `Data` object. Unless otherwise noted, these methods are applicable to resources returned by the `resources.GetRemote` function.

`Artifacts`
: {{< new-in 0.165.0 />}}
: (`slice`) Applicable to resources returned by the `css.Build` and `js.Build` functions, a slice of the additional output files published as part of the build, such as source maps and files emitted by esbuild's `file` loader. Each element of the slice provides `MediaType`, `Permalink`, and `RelPermalink` methods. See the Artifacts sections of the [`css.Build`][css artifacts] and [`js.Build`][js artifacts] documentation.

`ContentLength`
: (`int`) The content length in bytes.

`ContentType`
: (`string`) The content type.

`Headers`
: (`map[string][]string`) A map of response headers matching those requested in the [`responseHeaders`][] option passed to the `resources.GetRemote` function. The header name matching is case-insensitive. In most cases there will be one value per header key.

`Status`
: (`string`) The HTTP status text.

`StatusCode`
: (`int`) The HTTP status code.

`TransferEncoding`
: (`string`) The transfer encoding.

[`css.Build`]: /functions/css/build/
[`js.Build`]: /functions/js/build/
[`resources.GetRemote`]: /functions/resources/getremote/
[`responseHeaders`]: /functions/resources/getremote/#responseheaders
[css artifacts]: /functions/css/build/#artifacts
[js artifacts]: /functions/js/build/#artifacts
