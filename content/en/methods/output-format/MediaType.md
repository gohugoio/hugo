---
title: MediaType
description: Returns the media type of the given output format.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: media.Type
    signatures: [OUTPUTFORMAT.MediaType]
---

{{% include "/_common/methods/output-formats/to-use-this-method.md" %}}

```go-html-template
{{ with .Site.Home.OutputFormats.Get "rss" }}
  {{ with .MediaType }}
    {{ .Type }}       → application/rss+xml
    {{ .MainType }}   → application
    {{ .SubType }}    → rss
  {{ end }}
{{ end }}
```

## Methods

### MainType

(`string`) Returns the main type of the output format's media type.

### SubType

(`string`) Returns the subtype of the current format's media type.

### Type

(`string`) Returns the the current format's media type.
