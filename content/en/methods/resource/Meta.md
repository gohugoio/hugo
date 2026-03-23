---
title: Meta
description: Applicable to images, returns an object containing Exif, IPTC, and XMP metadata for supported image formats.
categories: []
keywords: ['metadata']
params:
  functions_and_methods:
    returnType: meta.MetaInfo
    signatures: [RESOURCE.Meta]
---

{{< new-in 0.155.3 />}}

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

The `Meta` method on an image `Resource` object returns an object containing [Exif][Exif_Definition], [IPTC][IPTC_Definition], and [XMP][XMP_Definition] metadata.

While Hugo classifies many file types as images, only certain formats support metadata extraction. Supported formats include AVIF, BMP, GIF, HEIC, HEIF, JPEG, PNG, TIFF, and WebP.

> [!note]
> Metadata is not preserved during image transformation. Use this method with the _original_ image resource to extract metadata from supported formats.

## Usage

Use the [`reflect.IsImageResourceWithMeta`][] function to verify that a resource supports metadata extraction before calling the `Meta` method.

```go-html-template
{{ with resources.GetMatch "images/featured.*" }}
  {{ if reflect.IsImageResourceWithMeta . }}
    {{ with .Meta }}
      {{ .Date.Format "2006-01-02" }}
    {{ end }}
  {{ end }}
{{ end }}
```

## Methods

### Date

(`time.Time`) Returns the image creation date/time. Format with the [`time.Format`][] function.

### Lat

(`float64`) Returns the GPS latitude in degrees from Exif metadata, with a fallback to XMP metadata.

### Long

(`float64`) Returns the GPS longitude in degrees from Exif metadata, with a fallback to XMP metadata.

### Orientation

(`int`) Returns the value of the Exif `Orientation` tag, one of eight possible values.

Value|Description
:--|:--
`1`|Horizontal (normal)
`2`|Mirrored horizontal
`3`|Rotated 180 degrees
`4`|Mirrored vertical
`5`|Mirrored horizontal and rotated 270 degrees clockwise
`6`|Rotated 90 degrees clockwise
`7`|Mirrored horizontal and rotated 90 degrees clockwise
`8`|Rotated 270 degrees clockwise
{class="!mt-0"}

> [!tip]
> Use the [`images.AutoOrient`][] image filter to rotate and flip an image as needed per its Exif orientation tag

### Exif

(`meta.Tags`) Returns a collection of available Exif fields for this image. Availability is determined by the [`sources`][] setting and specific fields are managed via the [`fields`][] setting, both of which are managed in your project configuration.

### IPTC

(`meta.Tags`) Returns a collection of available IPTC fields for this image. Availability is determined by the [`sources`][] setting and specific fields are managed via the [`fields`][] setting, both of which are managed in your project configuration.

### XMP

(`meta.Tags`) Returns a collection of available XMP fields for this image. Availability is determined by the [`sources`][] setting and specific fields are managed via the [`fields`][] setting, both of which are managed in your project configuration.

## Examples

To list the creation date, latitude, longitude, and orientation:

```go-html-template
{{ with resources.GetMatch "images/featured.*" }}
  {{ if reflect.IsImageResourceWithMeta . }}
    {{ with .Meta }}
      <pre>
        {{ printf "%-25s %v\n" "Date" .Date }}
        {{ printf "%-25s %v\n" "Latitude" .Lat }}
        {{ printf "%-25s %v\n" "Longitude" .Long }}
        {{ printf "%-25s %v\n" "Orientation" .Orientation }}
      </pre>
    {{ end }}
  {{ end }}
{{ end }}
```

{{% include "/_common/functions/reflect/image-reflection-functions.md" %}}

[`fields`]: /configuration/imaging/#fields
[`images.AutoOrient`]: /functions/images/autoorient/
[`reflect.IsImageResourceWithMeta`]: /functions/reflect/isimageresourcewithmeta/
[`sources`]: /configuration/imaging/#sources
[`time.Format`]: /functions/time/format/
[Exif_Definition]: https://en.wikipedia.org/wiki/Exif
[IPTC_Definition]: https://en.wikipedia.org/wiki/IPTC_Information_Interchange_Model
[XMP_Definition]: https://en.wikipedia.org/wiki/Extensible_Metadata_Platform
