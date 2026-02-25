---
title: Meta
description: Applicable to JPEG, PNG, TIFF, and WebP images, returns an object containing Exif, IPTC, and XMP metadata.
categories: []
keywords: ['metadata']
params:
  functions_and_methods:
    returnType: meta.MetaInfo
    signatures: [RESOURCE.Meta]
---

{{< new-in 0.155.3 />}}

{{% include "/_common/methods/resource/global-page-remote-resources.md" %}}

Applicable to JPEG, PNG, TIFF, and WebP images, the `Meta` method on an image `Resource` object returns an object containing [Exif][Exif_Definition], [IPTC][IPTC_Definition], and [XMP][XMP_Definition] metadata.

To extract Exif metadata only, use the [`Exif`] method instead.

> [!note]
> Metadata is not preserved during image transformation. Use this method with the _original_ image resource to extract metadata from JPEG, PNG, TIFF, and WebP images.

## Methods

### Date

(`time.Time`) Returns the image creation date/time. Format with the [`time.Format`] function.

### Lat

(`float64`) Returns the GPS latitude in degrees from Exif metadata, with a fallback to XMP metadata.

### Long

(`float64`) Returns the GPS longitude in degrees from Exif metadata, with a fallback to XMP metadata.

### Orientation

(`int`) Returns the value of the Exif `Orientation` tag, one of eight possible values:

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
> Use the [`images.AutoOrient`] image filter to rotate and flip an image as needed per its Exif orientation tag

### Exif

(`meta.Tags`) Returns a collection of available Exif fields for this image. Availability is determined by the [`sources`][] setting and specific fields are managed via the [`fields`][] setting, both of which are managed in your project configuration.

### IPTC

(`meta.Tags`) Returns a collection of available IPTC fields for this image. Availability is determined by the [`sources`][] setting and specific fields are managed via the [`fields`][] setting, both of which are managed in your project configuration.

### XMP

(`meta.Tags`) Returns a collection of available XMP fields for this image. Availability is determined by the [`sources`][] setting and specific fields are managed via the [`fields`][] setting, both of which are managed in your project configuration.

## Examples

To list the creation date, latitude, longitude, and orientation:

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ with .Meta }}
    <pre>
      {{ printf "%-25s %v\n" "Date" .Date }}
      {{ printf "%-25s %v\n" "Latitude" .Lat }}
      {{ printf "%-25s %v\n" "Longitude" .Long }}
      {{ printf "%-25s %v\n" "Orientation" .Orientation }}
    </pre>
  {{ end }}
{{ end }}
```

To list the available Exif fields:

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ with .Meta }}
    <pre>
      {{ range $k, $v := .Exif -}}
        {{ printf "%-25s %v\n" $k $v }}
      {{ end }}
    </pre>
  {{ end }}
{{ end }}
```

To list the available IPTC fields:

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ with .Meta }}
    <pre>
      {{ range $k, $v := .IPTC -}}
        {{ printf "%-25s %v\n" $k $v }}
      {{ end }}
    </pre>
  {{ end }}
{{ end }}
```

To list the available XMP fields:

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ with .Meta }}
    <pre>
      {{ range $k, $v := .XMP -}}
        {{ printf "%-25s %v\n" $k $v }}
      {{ end }}
    </pre>
  {{ end }}
{{ end }}
```

To list the available Exif, IPTC, and XMP fields together:

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ with .Meta }}
    <pre>
      {{ range $k, $v := merge .Exif .IPTC .XMP -}}
        {{ printf "%-25s %v\n" $k $v }}
      {{ end }}
    </pre>
  {{ end }}
{{ end }}
```

To list specific values:

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ with .Meta }}
    <pre>
      {{ with .Exif.ApertureValue }}{{ printf "%-25s %v\n" "ApertureValue" . }}{{ end }}
      {{ with .Exif.BrightnessValue }}{{ printf "%-25s %v\n" "BrightnessValue" . }}{{ end }}

      {{ with .IPTC.Headline }}{{ printf "%-25s %v\n" "Headline" . }}{{ end }}
      {{ with index .IPTC "Province-State" }}{{ printf "%-25s %v\n" "Province-State" . }}{{ end }}

      {{ with .XMP.Creator }}{{ printf "%-25s %v\n" "Creator" . }}{{ end }}
      {{ with .XMP.Subject }}{{ printf "%-25s %v\n" "Subject" . }}{{ end }}
    </pre>
  {{ end }}
{{ end }}
```

[`Exif`]: /methods/resource/exif/
[`fields`]: /configuration/imaging/#fields
[`images.AutoOrient`]: /functions/images/autoorient/
[`sources`]: /configuration/imaging/#sources
[`time.Format`]: /functions/time/format/
[Exif_Definition]: https://en.wikipedia.org/wiki/Exif
[IPTC_Definition]: https://en.wikipedia.org/wiki/IPTC_Information_Interchange_Model
[XMP_Definition]: https://en.wikipedia.org/wiki/Extensible_Metadata_Platform
