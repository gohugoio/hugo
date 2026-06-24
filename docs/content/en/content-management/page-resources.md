---
title: Page resources
description: Use page resources to logically associate assets with a page.
categories: []
keywords: []
---

Page resources are only accessible from [page bundles][], those directories with `index.md` or`_index.md` files at their root. Page resources are only available to the page with which they are bundled.

In this example, `first-post` is a page bundle with access to 10 page resources including audio, data, documents, images, and video. Although `second-post` is also a page bundle, it has no page resources and is unable to directly access the page resources associated with `first-post`.

```tree
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

- [`Resources.ByType`][]
- [`Resources.Get`][]
- [`Resources.GetMatch`][]
- [`Resources.Match`][]

 Once you have captured a resource, use any of the applicable [`Resource`][] methods to return a value or perform an action.

The following examples assume this content structure:

```tree
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

The page resources' metadata is managed from the corresponding page's front matter with an array parameter named `resources`.

> [!NOTE]
> Resources of type `page` get `Title` etc. from their own front matter.

`src`
: (`string`) Required. A [glob pattern](g) matching one or more page resources by file path, relative to the page bundle. Matching is case-insensitive. When the pattern matches multiple resources, the same metadata is applied to each.

`name`
: (`string`) Sets the value returned by [`Name`][]. Supports the [`:counter`](#the-counter-placeholder-in-name-and-title) placeholder. After assignment, use `name`, not the original file path, with [`Resources.Get`][], [`Resources.Match`][], and [`Resources.GetMatch`][].

`title`
: (`string`) Sets the value returned by [`Title`][]. Supports the [`:counter`](#the-counter-placeholder-in-name-and-title) placeholder.

`params`
: (`map`) A map of custom key-value pairs. When multiple array entries match the same resource, their `params` maps are merged; later entries take precedence for duplicate keys.

### Resources metadata example

<!-- markdownlint-disable MD007 MD032 -->
{{< code-toggle file=content/example.md fm=true >}}
title: Application
date: 2018-01-25
resources:
  - src: images/sunset.jpg
    name: header
  - src: documents/photo_specs.pdf
    title: Photo Specifications
  - src: documents/guide.pdf
    title: Instruction Guide
  - src: documents/checklist.pdf
    title: Document Checklist
  - src: documents/payment.docx
    title: Proof of Payment
  - src: "**.pdf"
    name: pdf-file-:counter
    params:
      icon: pdf
  - src: "**.docx"
    params:
      icon: word
{{</ code-toggle >}}
<!-- markdownlint-enable MD007 MD032 -->

From the example above:

- `sunset.jpg` will receive a new `Name` and can now be found with `.GetMatch "header"`.
- `documents/photo_specs.pdf`, `documents/guide.pdf`, `documents/checklist.pdf`, and `documents/payment.docx` will get `Title` as set by `title`.
- All `PDF` files will get the `pdf` icon and a new `Name`. The `name` parameter contains a special placeholder [`:counter`](#the-counter-placeholder-in-name-and-title), so the `Name` will be `pdf-file-1`, `pdf-file-2`, `pdf-file-3`.
- All `.docx` files will get the `word` icon.

> [!NOTE]
> For `name` and `title`, the first matching array entry wins; later matches are ignored. For `params`, all matching entries contribute; later entries take precedence for duplicate keys. Place more specific `src` patterns before broader wildcards to control which `name` and `title` values are applied.

### The `:counter` placeholder in `name` and `title`

The `:counter` is a special placeholder recognized in `name` and `title` parameters `resources`.

Each unique `src` pattern maintains independent counters for `name` and `title`, each starting at 1 with the first matching resource.

For example, if a bundle has the resources `photo_specs.pdf`, `other_specs.pdf`, `guide.pdf` and `checklist.pdf`, and the front matter has specified the `resources` as:

{{< code-toggle file=content/inspections/engine/index.md fm=true >}}
title = 'Engine inspections'
[[resources]]
  src = '*specs.pdf'
  title = 'Specification #:counter'
[[resources]]
  src = '**.pdf'
  name = 'pdf-file-:counter.pdf'
{{</ code-toggle >}}

the `Name` and `Title` will be assigned to the resource files as follows:

| Resource file    | `Name`             | `Title`              |
|------------------|--------------------|----------------------|
| checklist.pdf    | `"pdf-file-1.pdf"` | `"checklist.pdf"`    |
| guide.pdf        | `"pdf-file-2.pdf"` | `"guide.pdf"`        |
| other\_specs.pdf | `"pdf-file-3.pdf"` | `"Specification #1"` |
| photo\_specs.pdf | `"pdf-file-4.pdf"` | `"Specification #2"` |

## Multilingual

By default, with a multilingual single-host project, Hugo does not duplicate shared page during the build.

> [!NOTE]
> This behavior is limited to Markdown content. Shared page resources for other [content formats][] are copied into each language bundle.

Consider this project configuration:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'de'
defaultContentLanguageInSubdir = true

[languages.de]
label = 'Deutsch'
locale = 'de-DE'
weight = 1

[languages.en]
label = 'English'
locale = 'en-US'
weight = 2
{{< /code-toggle >}}

And this content:

```tree
content/
└── my-bundle/
    ├── a.jpg     <-- shared page resource
    ├── b.jpg     <-- shared page resource
    ├── c.de.jpg
    ├── c.en.jpg
    ├── index.de.md
    └── index.en.md
```

Hugo places the shared resources in the page bundle for the default content language:

```tree
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

> [!IMPORTANT]
> To resolve Markdown link and image destinations to the correct location, you must use link and image render hooks that capture the page resource with the [`Resources.Get`][] method, and then invoke its [`RelPermalink`][] method.
>
> In its default configuration, Hugo automatically uses the [embedded link render hook][] and the [embedded image render hook][] for multilingual single-host projects, specifically when the [duplication of shared page resources][] feature is disabled. This is the default behavior for such projects. If custom link or image render hooks are defined by your project, modules, or themes, these will be used instead.
>
> You can also configure Hugo to `always` use the embedded link or image render hook, use it only as a `fallback`, or `never` use it. See [details][].

Although duplicating shared page resources is inefficient, you can enable this feature in your project configuration if desired:

{{< code-toggle file=hugo >}}
[markup.goldmark]
duplicateResourceFiles = true
{{< /code-toggle >}}

[`Name`]: /methods/resource/name/
[`RelPermalink`]: /methods/resource/relpermalink/
[`Resource`]: /methods/resource/
[`Resources.ByType`]: /methods/page/resources#bytype
[`Resources.GetMatch`]: /methods/page/resources#getmatch
[`Resources.Get`]: /methods/page/resources/#get
[`Resources.Match`]: /methods/page/resources#match
[`Title`]: /methods/resource/title/
[content formats]: /content-management/formats/
[details]: /configuration/markup/#renderhookslinkuseembedded
[duplication of shared page resources]: /configuration/markup/#duplicateresourcefiles
[embedded image render hook]: /render-hooks/images/#embedded
[embedded link render hook]: /render-hooks/links/#embedded
[page bundles]: /content-management/page-bundles/
