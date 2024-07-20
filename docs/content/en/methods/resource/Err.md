---
title: Err
description: Applicable to resources returned by the resources.GetRemote function, returns an error message if the HTTP request fails, else nil. 
categories: []
keywords: []
action:
  related:
    - functions/resources/GetRemote
    - methods/resource/Data
  returnType: resource.resourceError
  signatures: [RESOURCE.Err]
---

The `Err` method on a resource returned by the [`resources.GetRemote`] function returns an error message if the HTTP request fails, else nil. If you do not handle the error yourself, Hugo will fail the build.

[`resources.GetRemote`]: /functions/resources/getremote/

In this example we send an HTTP request to a nonexistent domain:

```go-html-template
{{ $url := "https://broken-example.org/images/a.jpg" }}
{{ with resources.GetRemote $url }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ else }}
  {{ errorf "Unable to get remote resource %q" $url }}
{{ end }}
```

The code above captures the error from the HTTP request, then fails the build:

```text
ERROR error calling resources.GetRemote: Get "https://broken-example.org/images/a.jpg": dial tcp: lookup broken-example.org on 127.0.0.53:53: no such host
```

To log an error as a warning instead of an error:

```go-html-template
{{ $url := "https://broken-example.org/images/a.jpg" }}
{{ with resources.GetRemote $url }}
  {{ with .Err }}
    {{ warnf "%s" . }}
  {{ else }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ else }}
  {{ errorf "Unable to get remote resource %q" $url }}
{{ end }}
```

{{% note %}}
An HTTP response with a 404 status code is not an HTTP request error. To handle 404 status codes, code defensively using the nested `with-else-end` construct as shown above.
{{% /note %}}
