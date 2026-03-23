---
_comment: Do not remove front matter.
---

## Image operations

Use these functions to determine which operations Hugo supports for a given resource. While Hugo classifies a variety of file types as image resources, its ability to process them or extract metadata varies by format.

- [`reflect.IsImageResource`][]: {{% get-page-desc "/functions/reflect/isimageresource" %}}
- [`reflect.IsImageResourceProcessable`][]: {{% get-page-desc "/functions/reflect/isimageresourceprocessable" %}}
- [`reflect.IsImageResourceWithMeta`][]: {{% get-page-desc "/functions/reflect/isimageresourcewithmeta" %}}

The table below shows the values these functions return for various file formats. Use it to determine which checks are required before calling specific methods in your templates.

|Format|IsImageResource|IsImageResourceProcessable|IsImageResourceWithMeta|
|:-----|:--------------|:-------------------------|:----------------------|
|AVIF  |true           |**false**                 |true                   |
|BMP   |true           |true                      |true                   |
|GIF   |true           |true                      |true                   |
|HEIC  |true           |**false**                 |true                   |
|HEIF  |true           |**false**                 |true                   |
|ICO   |true           |**false**                 |**false**              |
|JPEG  |true           |true                      |true                   |
|PNG   |true           |true                      |true                   |
|SVG   |true           |**false**                 |**false**              |
|TIFF  |true           |true                      |true                   |
|WebP  |true           |true                      |true                   |

This contrived example demonstrates how to iterate through resources and use these functions to apply the appropriate handling for each image format.

```go-html-template
{{ range resources.Match "**" }}
  {{ if reflect.IsImageResource . }}
    {{ if reflect.IsImageResourceProcessable . }}
      {{ with .Process "resize 300x webp" }}
        <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
      {{ end }}
    {{ else if reflect.IsImageResourceWithMeta . }}
      <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
    {{ else }}
      <img src="{{ .RelPermalink }}" alt="">
    {{ end }}
  {{ end }}
{{ end }}
```

[`reflect.IsImageResource`]: /functions/reflect/isimageresource/
[`reflect.IsImageResourceProcessable`]: /functions/reflect/isimageresourceprocessable/
[`reflect.IsImageResourceWithMeta`]: /functions/reflect/isimageresourcewithmeta/
