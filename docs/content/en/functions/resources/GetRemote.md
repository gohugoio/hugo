---
title: resources.GetRemote
description: Returns a remote resource from the given URL, or nil if none found.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: resource.Resource
    signatures: ['resources.GetRemote URL [OPTIONS]']
---

{{< new-in 0.141.0 >}}
The `Err` method on the returned resource was removed in v0.141.0.

Use the [`try`] statement instead, as shown in the [error handling] example below.

{{< /new-in >}}

```go-html-template
{{ $url := "https://example.org/images/a.jpg" }}
{{ with try (resources.GetRemote $url) }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else with .Value }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ else }}
    {{ errorf "Unable to get remote resource %q" $url }}
  {{ end }}
{{ end }}
```

## Options

The `resources.GetRemote` function takes an optional map of options.

body
: (`string`) The data you want to transmit to the server.

headers
: (`map[string][]string`) The collection of key-value pairs that provide additional information about the request.

key
: (`string`) The cache key. Hugo derives the default value from the URL and options map. See [caching](#caching).

method
: (`string`) The action to perform on the requested resource, typically one of `GET`, `POST`, or `HEAD`.

responseHeaders
: {{< new-in 0.143.0 />}}
: (`[]string`) The headers to extract from the server's response, accessible through the resource's [`Data.Headers`] method. Header name matching is case-insensitive.

[`Data.Headers`]: /methods/resource/data/#headers

## Options examples

> [!note]
> For brevity, the examples below do not include [error handling].

To include a header:

```go-html-template
{{ $url := "https://example.org/api" }}
{{ $opts := dict
  "headers" (dict "Authorization" "Bearer abcd")
}}
{{ $resource := resources.GetRemote $url $opts }}
```

To specify more than one value for the same header key, use a slice:

```go-html-template
{{ $url := "https://example.org/api" }}
{{ $opts := dict
  "headers" (dict "X-List" (slice "a" "b" "c"))
}}
{{ $resource := resources.GetRemote $url $opts }}
```

To post data:

```go-html-template
{{ $url := "https://example.org/api" }}
{{ $opts := dict
  "method" "post"
  "body" `{"complete": true}` 
  "headers" (dict  "Content-Type" "application/json")
}}
{{ $resource := resources.GetRemote $url $opts }}
```

To override the default cache key:

```go-html-template
{{ $url := "https://example.org/images/a.jpg" }}
{{ $opts := dict 
  "key" (print $url (now.Format "2006-01-02"))
}}
{{ $resource := resources.GetRemote $url $opts }}
```

To extract specific headers from the server's response:

```go-html-template
{{ $url := "https://example.org/images/a.jpg" }}
{{ $opts := dict
  "method" "HEAD"
  "responseHeaders" (slice "X-Frame-Options" "Server")
}}
{{ $resource := resources.GetRemote $url $opts }}
```

## Remote data

When retrieving remote data, use the [`transform.Unmarshal`] function to [unmarshal](g) the response.

[`transform.Unmarshal`]: /functions/transform/unmarshal/

```go-html-template
{{ $data := dict }}
{{ $url := "https://example.org/books.json" }}
{{ with try (resources.GetRemote $url) }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else with .Value }}
    {{ $data = . | transform.Unmarshal }}
  {{ else }}
    {{ errorf "Unable to get remote resource %q" $url }}
  {{ end }}
{{ end }}
```

> [!note]
> When retrieving remote data, a misconfigured server may send a response header with an incorrect [Content-Type]. For example, the server may set the Content-Type header to `application/octet-stream` instead of `application/json`.
>
> In these cases, pass the resource `Content` through the `transform.Unmarshal` function instead of passing the resource itself. For example, in the above, do this instead:
>
> `{{ $data = .Content | transform.Unmarshal }}`

## Error handling

Use the [`try`] statement to capture HTTP request errors. If you do not handle the error yourself, Hugo will fail the build.

> [!note]
> Hugo does not classify an HTTP response with status code 404 as an error. In this case `resources.GetRemote` returns nil.

```go-html-template
{{ $url := "https://broken-example.org/images/a.jpg" }}
{{ with try (resources.GetRemote $url) }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else with .Value }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ else }}
    {{ errorf "Unable to get remote resource %q" $url }}
  {{ end }}
{{ end }}
```

To log an error as a warning instead of an error:

```go-html-template
{{ $url := "https://broken-example.org/images/a.jpg" }}
{{ with try (resources.GetRemote $url) }}
  {{ with .Err }}
    {{ warnf "%s" . }}
  {{ else with .Value }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ else }}
    {{ warnf "Unable to get remote resource %q" $url }}
  {{ end }}
{{ end }}
```

## HTTP response

The [`Data`] method on a resource returned by the `resources.GetRemote` function returns information from the HTTP response.

[`Data`]: /methods/resource/data/

## Caching

Resources returned from `resources.GetRemote` are cached to disk. See [configure file caches] for details.

By default, Hugo derives the cache key from the arguments passed to the function. Override the cache key by setting a `key` in the options map. Use this approach to have more control over how often Hugo fetches a remote resource.

```go-html-template
{{ $url := "https://example.org/images/a.jpg" }}
{{ $cacheKey := print $url (now.Format "2006-01-02") }}
{{ $opts := dict "key" $cacheKey }}
{{ $resource := resources.GetRemote $url $opts }}
```

## Security

To protect against malicious intent, the `resources.GetRemote` function inspects the server response including:

- The [Content-Type] in the response header
- The file extension, if any
- The content itself

If Hugo is unable to resolve the media type to an entry in its [allowlist], the function throws an error:

```text
ERROR error calling resources.GetRemote: failed to resolve media type...
```

For example, you will see the error above if you attempt to download an executable.

Although the allowlist contains entries for common media types, you may encounter situations where Hugo is unable to resolve the media type of a file that you know to be safe. In these situations, edit your site configuration to add the media type to the allowlist. For example:

{{< code-toggle file=hugo >}}
[security.http]
mediaTypes = ['^image/avif$','^application/vnd\.api\+json$']
{{< /code-toggle >}}

Note that the entry above is:

- An _addition_ to the allowlist; it does not _replace_ the allowlist
- An array of [regular expressions](g)

[allowlist]: https://en.wikipedia.org/wiki/Whitelist
[Content-Type]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type

[`try`]: /functions/go-template/try
[configure file caches]: /configuration/caches/
[error handling]: #error-handling
