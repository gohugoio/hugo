---
title: encoding.Base64Decode
description: Returns the base64 decoding of the given content.
categories: []
keywords: []
action:
  aliases: [base64Decode]
  related:
    - functions/encoding/Base64Encode
  returnType: string
  signatures: [encoding.Base64Decode INPUT]
aliases: [/functions/base64Decode]
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
{{ $u := "https://api.github.com/repos/gohugoio/hugo/readme" }}
{{ with resources.GetRemote $u }}
  {{ with .Err }}
    {{ errorf "%s" . }}
  {{ else }}
    {{ with . | transform.Unmarshal }}
      {{ .content | base64Decode | markdownify }}
    {{ end }}
  {{ end }}
{{ else }}
  {{ errorf "Unable to get remote resource %q" $u }}
{{ end }}
```
