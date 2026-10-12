---
title: encoding.Base64Decode
description: Returns the base64 decoding of the given content.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [base64Decode]
    returnType: string
    signatures: [encoding.Base64Decode INPUT]
---

```go-html-template
{{ "SHVnbw==" | base64Decode }} → Hugo
```

Use the `base64Decode` function to decode responses from APIs. For example, the result of this call to GitHub's API contains the base64-encoded representation of the repository's README file:

```text
https://api.github.com/repos/gohugoio/hugo/readme
```

To retrieve and render the content:

```go-html-template
{{ $url := "https://api.github.com/repos/gohugoio/hugo/readme" }}
{{ with try (resources.GetRemote $url) }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else with .Value }}
    {{ with . | transform.Unmarshal }}
      {{ .content | base64Decode | markdownify }}
    {{ end }}
  {{ else }}
    {{ errorf "Unable to get remote resource %q" $url }}
  {{ end }}
{{ end }}
```
