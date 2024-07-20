---
title: Data
description: Applicable to resources returned by the resources.GetRemote function, returns information from the HTTP response.
categories: []
keywords: []
action:
  related:
    - functions/resources/GetRemote
    - methods/resource/Err
  returnType: map
  signatures: [RESOURCE.Data]
---

The `Data` method on a resource returned by the [`resources.GetRemote`] function returns information from the HTTP response.

[`resources.GetRemote`]: /functions/resources/getremote/

```go-html-template
{{ $url := "https://example.org/images/a.jpg" }}
{{ with resources.GetRemote $url }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else }}
    {{ with .Data }}
      {{ .ContentLength }} → 42764
      {{ .ContentType }} → image/jpeg
      {{ .Status }} → 200 OK
      {{ .StatusCode }} → 200
      {{ .TransferEncoding }} → []
    {{ end }}
  {{ end }}
{{ else }}
  {{ errorf "Unable to get remote resource %q" $url }}
{{ end }}
```

ContentLength
: (`int`) The content length in bytes.

ContentType
: (`string`) The content type.

Status
: (`string`) The HTTP status text.

StatusCode
: (`int`) The HTTP status code.

TransferEncoding
: (`string`) The transfer encoding.


[`resources.GetRemote`]: /functions/resources/getremote/
