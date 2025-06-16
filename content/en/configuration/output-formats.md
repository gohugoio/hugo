---
title: Configure output formats
linkTitle: Output formats
description: Configure output formats.
categories: []
keywords: []
---

{{% glossary-term "output format" %}}

You can output a page in as many formats as you want. Define an infinite number of output formats, provided they each resolve to a unique file system path.

This is the default output format configuration in tabular form:

{{< datatable
  "config"
  "outputFormats"
  "_key"
  "mediaType"
  "weight"
  "baseName"
  "isHTML"
  "isPlainText"
  "noUgly"
  "notAlternative"
  "path"
  "permalinkable"
  "protocol"
  "rel"
  "root"
  "ugly"
>}}

## Default configuration

The following is the default configuration that matches the table above:

{{< code-toggle config=outputFormats />}}

baseName
: (`string`) The base name of the published file. Default is `index`.

isHTML
: (`bool`) Whether to classify the output format as HTML. Hugo uses this value to determine when to create alias redirects and when to inject the LiveReload script. Default is `false`.

isPlainText
: (`bool`) Whether to parse templates for this output format with Go's [text/template] package instead of the [html/template] package. Default is `false`.

mediaType
: (`string`) The [media type](g) of the published file. This must match one of the [configured media types].

notAlternative
: (`bool`) Whether to exclude this output format from the values returned by the [`AlternativeOutputFormats`] method on a `Page` object. Default is `false`.

noUgly
: (`bool`) Whether to disable ugly URLs for this output format when [`uglyURLs`] are enabled in your site configuration. Default is `false`.

path
: (`string`) The published file's directory path, relative to the root of the publish directory. If not specified, the file will be published using its content path.

permalinkable
: (`bool`) Whether to return the rendering output format rather than main output format when invoking the [`Permalink`] and [`RelPermalink`] methods on a `Page` object. See&nbsp;[details](#link-to-output-formats). Enabled by default for the `html` and `amp` output formats. Default is `false`.

protocol
: (`string`) The protocol (scheme) of the URL for this output format. For example, `https://` or `webcal://`. Default is the scheme of the [`baseURL`] parameter in your site configuration, typically `https://`.

rel
: (`string`) If provided, you can assign this value to `rel` attributes in `link` elements when iterating over output formats in your templates. Default is `alternate`.

root
: (`bool`) Whether to publish files to the root of the publish directory. Default is `false`.

ugly
: (`bool`) Whether to enable uglyURLs for this output format when `uglyURLs` is `false` in your site configuration. Default is `false`.

weight
: (`int`) When set to a non-zero value, Hugo uses the `weight` as the first criteria when sorting output formats, falling back to the name of the output format. Lighter items float to the top, while heavier items sink to the bottom. Hugo renders output formats sequentially based on the sort order. Default is `0`, except for the `html` output format, which has a default weight of `10`.

## Modify an output format

You can modify any of the default output formats. For example, to prioritize `json` rendering over `html` rendering, when both are generated, adjust the [`weight`](#weight):

{{< code-toggle file=hugo >}}
[outputFormats.json]
weight = 1
[outputFormats.html]
weight = 2
{{< /code-toggle >}}

The example above shows that when you modify a default content format, you only need to define the properties that differ from their default values.

## Create an output format

You can create new output formats as needed. For example, you may wish to create an output format to support Atom feeds.

### Step 1

Output formats require a specified media type. Because Atom feeds use `application/atom+xml`, which is not one of the [default media types], you must create it first.

{{< code-toggle file=hugo >}}
[mediaTypes.'application/atom+xml']
suffixes = ['atom']
{{< /code-toggle >}}

See [configure media types] for more information.

### Step 2

Create a new output format:

{{< code-toggle file=hugo >}}
[outputFormats.atom]
mediaType = 'application/atom+xml'
noUgly = true
{{< /code-toggle >}}

Note that we use the default settings for all other output format properties.

### Step 3

Specify the page [kinds](g) for which to render this output format:

{{< code-toggle file=hugo >}}
[outputs]
home = ['html', 'rss', 'atom']
section = ['html', 'rss', 'atom']
taxonomy = ['html', 'rss', 'atom']
term = ['html', 'rss', 'atom']
{{< /code-toggle >}}

See [configure outputs] for more information.

### Step 4

Create a template to render the output format. Since Atom feeds are lists, you need to create a list template. Consult the [template lookup order] to find the correct template path:

```text
layouts/list.atom.atom
```

We leave writing the template code as an exercise for you. Aim for a result similar to the [embedded RSS template].

## List output formats

To access output formats, each `Page` object provides two methods: [`OutputFormats`] (for all formats, including the current one) and [`AlternativeOutputFormats`]. Use `AlternativeOutputFormats` to create a link `rel` list within your site's `head` element, as shown below:

```go-html-template
{{ range .AlternativeOutputFormats }}
  <link rel="{{ .Rel }}" type="{{ .MediaType.Type }}" href="{{ .Permalink | safeURL }}">
{{ end }}
```

## Link to output formats

By default, a `Page` object's [`Permalink`] and [`RelPermalink`] methods return the URL of the [primary output format](g), typically `html`. This behavior remains consistent regardless of the template used.

For example, in `page.json.json`, you'll see:

```go-html-template
{{ .RelPermalink }} → /that-page/
{{ with .OutputFormats.Get "json" }}
  {{ .RelPermalink }} → /that-page/index.json
{{ end }}
```

To make these methods return the URL of the _current_ template's output format, you must set the [`permalinkable`] setting to `true` for that format.

With `permalinkable` set to true for `json` in the same `page.json.json` template:

```go-html-template
{{ .RelPermalink }} → /that-page/index.json
{{ with .OutputFormats.Get "html" }}
  {{ .RelPermalink }} → /that-page/
{{ end }}
```

## Template lookup order

Each output format requires a template conforming to the [template lookup order].

For the highest specificity in the template lookup order, include the page kind, output format, and suffix in the file name:

```text
[page kind].[output format].[suffix]
```

For example, for section pages:

Output format|Template path
:--|:--
`html`|`layouts/section.html.html`
`json`|`layouts/section.json.json`
`rss`|`layouts/section.rss.xml`

[`AlternativeOutputFormats`]: /methods/page/alternativeoutputformats/
[`OutputFormats`]: /methods/page/outputformats/
[`Permalink`]: /methods/page/permalink/
[`RelPermalink`]: /methods/page/relpermalink/
[`baseURL`]: /configuration/all/#baseurl
[`permalinkable`]: #permalinkable
[`uglyURLs`]: /configuration/ugly-urls/
[configure media types]: /configuration/media-types/
[configure outputs]: /configuration/outputs/
[configured media types]: /configuration/media-types/
[default media types]: /configuration/media-types/
[embedded RSS template]: {{% eturl rss %}}
[html/template]: https://pkg.go.dev/html/template
[template lookup order]: /templates/lookup-order/
[text/template]: https://pkg.go.dev/text/template
