---
title: Page resources
description: Page resources -- images, other pages, documents, etc. -- have page-relative URLs and their own metadata.
categories: [content management]
keywords: [bundle,content,resources]
menu:
  docs:
    parent: content-management
    weight: 80
weight: 80
toc: true
---

Page resources are only accessible from [page bundles](/content-management/page-bundles), those directories with `index.md` or
`_index.md` files at their root. Page resources are only available to the
page with which they are bundled.

In this example, `first-post` is a page bundle with access to 10 page resources including audio, data, documents, images, and video. Although `second-post` is also a page bundle, it has no page resources and is unable to directly access the page resources associated with `first-post`.

```text
content
└── post
    ├── first-post
    │   ├── images
    │   │   ├── a.jpg
    │   │   ├── b.jpg
    │   │   └── c.jpg
    │   ├── index.md (root of page bundle)
    │   ├── latest.html
    │   ├── manual.json
    │   ├── notice.md
    │   ├── office.mp3
    │   ├── pocket.mp4
    │   ├── rating.pdf
    │   └── safety.txt
    └── second-post
        └── index.md (root of page bundle)
```

## Properties

ResourceType
: The main type of the resource's [Media Type](/templates/output-formats/#media-types). For example, a file of MIME type `image/jpeg` has the ResourceType `image`. A `Page` will have `ResourceType` with value `page`.

Name
: Default value is the file name (relative to the owning page). Can be set in front matter.

Title
: Default value is the same as `.Name`. Can be set in front matter.

Permalink
: The absolute URL to the resource. Resources of type `page` will have no value.

RelPermalink
: The relative URL to the resource. Resources of type `page` will have no value.

Content
: The content of the resource itself. For most resources, this returns a string
with the contents of the file. Use this to create inline resources.

```go-html-template
{{ with .Resources.GetMatch "script.js" }}
  <script>{{ .Content | safeJS }}</script>
{{ end }}

{{ with .Resources.GetMatch "style.css" }}
  <style>{{ .Content | safeCSS }}</style>
{{ end }}

{{ with .Resources.GetMatch "img.png" }}
  <img src="data:{{ .MediaType.Type }};base64,{{ .Content | base64Encode }}">
{{ end }}
```

MediaType.Type
: The media type (formerly known as a MIME type) of the resource (e.g., `image/jpeg`).

MediaType.MainType
: The main type of the resource's media type (e.g., `image`).

MediaType.SubType
: The subtype of the resource's type (e.g., `jpeg`). This may or may not correspond to the file suffix.

MediaType.Suffixes
: A slice of possible file suffixes for the resource's media type (e.g., `[jpg jpeg jpe jif jfif]`).

## Methods

ByType
: Returns the page resources of the given type.

```go-html-template
{{ .Resources.ByType "image" }}
```
Match
: Returns all the page resources (as a slice) whose `Name` matches the given Glob pattern ([examples](https://github.com/gobwas/glob/blob/master/readme.md)). The matching is case-insensitive.

```go-html-template
{{ .Resources.Match "images/*" }}
```

GetMatch
: Same as `Match` but will return the first match.

### Pattern matching

```go
// Using Match/GetMatch to find this images/sunset.jpg ?
.Resources.Match "images/sun*" ✅
.Resources.Match "**/sunset.jpg" ✅
.Resources.Match "images/*.jpg" ✅
.Resources.Match "**.jpg" ✅
.Resources.Match "*" 🚫
.Resources.Match "sunset.jpg" 🚫
.Resources.Match "*sunset.jpg" 🚫
```

## Metadata

The page resources' metadata is managed from the corresponding page's front matter with an array/table parameter named `resources`. You can batch assign values using [wildcards](https://tldp.org/LDP/GNU-Linux-Tools-Summary/html/x11655.htm).

{{% note %}}
Resources of type `page` get `Title` etc. from their own front matter.
{{% /note %}}

name
: Sets the value returned in `Name`.

{{% note %}}
The methods `Match`, `Get` and `GetMatch` use `Name` to match the resources.
{{% /note %}}

title
: Sets the value returned in `Title`

params
: A map of custom key-value pairs.

### Resources metadata example

{{< code-toggle >}}
title: Application
date : 2018-01-25
resources :
- src : "images/sunset.jpg"
  name : "header"
- src : "documents/photo_specs.pdf"
  title : "Photo Specifications"
  params:
    icon : "photo"
- src : "documents/guide.pdf"
  title : "Instruction Guide"
- src : "documents/checklist.pdf"
  title : "Document Checklist"
- src : "documents/payment.docx"
  title : "Proof of Payment"
- src : "**.pdf"
  name : "pdf-file-:counter"
  params :
    icon : "pdf"
- src : "**.docx"
  params :
    icon : "word"
{{</ code-toggle >}}

From the example above:

- `sunset.jpg` will receive a new `Name` and can now be found with `.GetMatch "header"`.
- `documents/photo_specs.pdf` will get the `photo` icon.
- `documents/checklist.pdf`, `documents/guide.pdf` and `documents/payment.docx` will get `Title` as set by `title`.
- Every `PDF` in the bundle except `documents/photo_specs.pdf` will get the `pdf` icon.
- All `PDF` files will get a new `Name`. The `name` parameter contains a special placeholder [`:counter`](#the-counter-placeholder-in-name-and-title), so the `Name` will be `pdf-file-1`, `pdf-file-2`, `pdf-file-3`.
- Every docx in the bundle will receive the `word` icon.

{{% note %}}
The __order matters__ --- Only the **first set** values of the `title`, `name` and `params`-**keys** will be used. Consecutive parameters will be set only for the ones not already set. In the above example, `.Params.icon` is first set to `"photo"` in `src = "documents/photo_specs.pdf"`. So that would not get overridden to `"pdf"` by the later set `src = "**.pdf"` rule.
{{% /note %}}

### The `:counter` placeholder in `name` and `title`

The `:counter` is a special placeholder recognized in `name` and `title` parameters `resources`.

The counter starts at 1 the first time they are used in either `name` or `title`.

For example, if a bundle has the resources `photo_specs.pdf`, `other_specs.pdf`, `guide.pdf` and `checklist.pdf`, and the front matter has specified the `resources` as:

{{< code-toggle file=content/inspections/engine/index.md fm=true >}}
title = 'Engine inspections'
[[resources]]
  src = "*specs.pdf"
  title = "Specification #:counter"
[[resources]]
  src = "**.pdf"
  name = "pdf-file-:counter"
{{</ code-toggle >}}

the `Name` and `Title` will be assigned to the resource files as follows:

| Resource file     | `Name`            | `Title`               |
|-------------------|-------------------|-----------------------|
| checklist.pdf     | `"pdf-file-1.pdf` | `"checklist.pdf"`     |
| guide.pdf         | `"pdf-file-2.pdf` | `"guide.pdf"`         |
| other\_specs.pdf  | `"pdf-file-3.pdf` | `"Specification #1"` |
| photo\_specs.pdf  | `"pdf-file-4.pdf` | `"Specification #2"` |

## Multilingual

{{< new-in 0.123.0 >}}

By default, with a multilingual single-host site, Hugo does not duplicate shared page resources when building the site.

{{% note %}}
This behavior is limited to Markdown content. Shared page resources for other [content formats] are copied into each language bundle.

[content formats]: /content-management/formats/
{{% /note %}}

Consider this site configuration:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'de'
defaultContentLanguageInSubdir = true

[languages.de]
languageCode = 'de-DE'
languageName = 'Deutsch'
weight = 1

[languages.en]
languageCode = 'en-US'
languageName = 'English'
weight = 2
{{< /code-toggle >}}

And this content:

```text
content/
└── my-bundle/
    ├── a.jpg     <-- shared page resource
    ├── b.jpg     <-- shared page resource
    ├── c.de.jpg
    ├── c.en.jpg
    ├── index.de.md
    └── index.en.md
```

With v0.122.0 and earlier, Hugo duplicated the shared page resources, creating copies for each language:

```text
public/
├── de/
│   ├── my-bundle/
│   │   ├── a.jpg     <-- shared page resource
│   │   ├── b.jpg     <-- shared page resource
│   │   ├── c.de.jpg
│   │   └── index.html
│   └── index.html
├── en/
│   ├── my-bundle/
│   │   ├── a.jpg     <-- shared page resource (duplicate)
│   │   ├── b.jpg     <-- shared page resource (duplicate)
│   │   ├── c.en.jpg
│   │   └── index.html
│   └── index.html
└── index.html

```

With v0.123.0 and later, Hugo places the shared resources in the page bundle for the default content language:

```text
public/
├── de/
│   ├── my-bundle/
│   │   ├── a.jpg     <-- shared page resource
│   │   ├── b.jpg     <-- shared page resource
│   │   ├── c.de.jpg
│   │   └── index.html
│   └── index.html
├── en/
│   ├── my-bundle/
│   │   ├── c.en.jpg
│   │   └── index.html
│   └── index.html
└── index.html
```

This approach reduces build times, storage requirements, bandwidth consumption, and deployment times, ultimately reducing cost.

{{% note %}}
To resolve Markdown link and image destinations to the correct location, you must use link and image render hooks that capture the page resource with the [`Resources.Get`] method, and then invoke its [`RelPermalink`] method.

By default, with multilingual single-host sites, Hugo enables its [embedded link render hook] and [embedded image render hook] to resolve Markdown link and image destinations.

You may override the embedded render hooks as needed, provided they capture the resource as described above.

[embedded link render hook]: /render-hooks/links/#default
[embedded image render hook]: /render-hooks/images/#default
[`Resources.Get`]: /methods/page/resources/#get
[`RelPermalink`]: /methods/resource/relpermalink/
{{% /note %}}

Although duplicating shared page resources is inefficient, you can enable this feature in your site configuration if desired:

{{< code-toggle file=hugo >}}
[markup.goldmark]
duplicateResourceFiles = true
{{< /code-toggle >}}
