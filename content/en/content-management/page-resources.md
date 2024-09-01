---
title: Page resources
description: Use page resources to logically associate assets with a page.
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

## Examples

Use any of these methods on a `Page` object to capture page resources:

 - [`Resources.ByType`]
 - [`Resources.Get`]
 - [`Resources.GetMatch`]
 - [`Resources.Match`]

 Once you have captured a resource, use any of the applicable [`Resource`] methods to return a value or perform an action. 

[`Resource`]: /methods/resource
[`Resources.ByType`]: /methods/page/resources#bytype
[`Resources.GetMatch`]: /methods/page/resources#getmatch
[`Resources.Get`]: /methods/page/resources#get
[`Resources.Match`]: /methods/page/resources#match

The following examples assume this content structure:

```text
content/
└── example/
    ├── data/
    │  └── books.json   <-- page resource
    ├── images/
    │  ├── a.jpg        <-- page resource
    │  └── b.jpg        <-- page resource
    ├── snippets/
    │  └── text.md      <-- page resource
    └── index.md
```

Render a single image, and throw an error if the file does not exist:

```go-html-template
{{ $path := "images/a.jpg" }}
{{ with .Resources.Get $path }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ else }}
  {{ errorf "Unable to get page resource %q" $path }}
{{ end }}
```

Render all images, resized to 300 px wide:

```go-html-template
{{ range .Resources.ByType "image" }}
  {{ with .Resize "300x" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

Render the markdown snippet:

```go-html-template
{{ with .Resources.Get "snippets/text.md" }}
  {{ .Content }}
{{ end }}
```

List the titles in the data file, and throw an error if the file does not exist.

```go-html-template
{{ $path := "data/books.json" }}
{{ with .Resources.Get $path }}
  {{ with . | transform.Unmarshal }}
    <p>Books:</p>
    <ul>
      {{ range . }}
        <li>{{ .title }}</li>
      {{ end }}
    </ul>
  {{ end }}
{{ else }}
  {{ errorf "Unable to get page resource %q" $path }}
{{ end }}
```

## Metadata

The page resources' metadata is managed from the corresponding page's front matter with an array/table parameter named `resources`. You can batch assign values using [wildcards](https://tldp.org/LDP/GNU-Linux-Tools-Summary/html/x11655.htm).

{{% note %}}
Resources of type `page` get `Title` etc. from their own front matter.
{{% /note %}}

name
: (`string`) Sets the value returned in `Name`.

{{% note %}}
The methods `Match`, `Get` and `GetMatch` use `Name` to match the resources.
{{% /note %}}

title
: (`string`) Sets the value returned in `Title`

params
: (`map`) A map of custom key-value pairs.

### Resources metadata example

{{< code-toggle file=content/example.md fm=true >}}
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
The order matters; only the first set values of the `title`, `name` and `params` keys will be used. Consecutive parameters will be set only for the ones not already set. In the above example, `.Params.icon` is first set to `"photo"` in `src = "documents/photo_specs.pdf"`. So that would not get overridden to `"pdf"` by the later set `src = "**.pdf"` rule.
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
