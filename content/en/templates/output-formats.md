---
title: Custom Output Formats
linktitle: Custom Output Formats
description: Hugo can output content in multiple formats, including calendar events, e-book formats, Google AMP, and JSON search indexes, or any custom text format.
date: 2017-03-22
publishdate: 2017-03-22
lastmod: 2019-12-11
categories: [templates]
keywords: ["amp","outputs","rss"]
menu:
  docs:
    parent: "templates"
    weight: 18
weight: 18
sections_weight: 18
draft: false
aliases: [/templates/outputs/,/extras/output-formats/,/content-management/custom-outputs/]
toc: true
---

This page describes how to properly configure your site with the media types and output formats, as well as where to create your templates for your custom outputs.

## Media Types

A [media type][] (also known as *MIME type* and *content type*) is a two-part identifier for file formats and format contents transmitted on the Internet.

This is the full set of built-in media types in Hugo:

{{< datatable "media" "types" "type" "suffixes" >}}

**Note:**

* It is possible to add custom media types or change the defaults; e.g., if you want to change the suffix for `text/html` to `asp`.
* `Suffixes` are the values that will be used for URLs and filenames for that media type in Hugo.
* The `Type` is the identifier that must be used when defining new/custom `Output Formats` (see below).
* The full set of media types will be registered in Hugo's built-in development server to make sure they are recognized by the browser.

To add or modify a media type, define it in a `mediaTypes` section in your [site configuration][config], either for all sites or for a given language.

{{< code-toggle file="config" >}}
[mediaTypes]
  [mediaTypes."text/enriched"]
  suffixes = ["enr"]
  [mediaTypes."text/html"]
  suffixes = ["asp"]
{{</ code-toggle >}}

The above example adds one new media type, `text/enriched`, and changes the suffix for the built-in `text/html` media type.

**Note:** these media types are configured for **your output formats**. If you want to redefine one of Hugo's default output formats (e.g. `HTML`), you also need to redefine the media type. So, if you want to change the suffix of the `HTML` output format from `html` (default) to `htm`:

```toml
[mediaTypes]
[mediaTypes."text/html"]
suffixes = ["htm"]

# Redefine HTML to update its media type.
[outputFormats]
[outputFormats.HTML]
mediaType = "text/html"
```

**Note** that for the above to work, you also need to add an `outputs` definition in your site config.

## Output Format Definitions

Given a media type and some additional configuration, you get an **Output Format**.

This is the full set of Hugo's built-in output formats:

{{< datatable "output" "formats" "name" "mediaType" "path" "baseName" "rel" "protocol" "isPlainText" "isHTML" "noUgly" "permalinkable" >}}

* A page can be output in as many output formats as you want, and you can have an infinite amount of output formats defined **as long as they resolve to a unique path on the file system**. In the above table, the best example of this is `AMP` vs. `HTML`. `AMP` has the value `amp` for `Path` so it doesn't overwrite the `HTML` version; e.g. we can now have both `/index.html` and `/amp/index.html`.
* The `MediaType` must match the `Type` of an already defined media type.
* You can define new output formats or redefine built-in output formats; e.g., if you want to put `AMP` pages in a different path.

To add or modify an output format, define it in an `outputFormats` section in your site's [configuration file](/getting-started/configuration/), either for all sites or for a given language.

{{< code-toggle file="config" >}}
[outputFormats.MyEnrichedFormat]
mediaType = "text/enriched"
baseName = "myindex"
isPlainText = true
protocol = "bep://"
{{</ code-toggle >}}

The above example is fictional, but if used for the homepage on a site with `baseURL` `https://example.org`, it will produce a plain text homepage with the URL `bep://example.org/myindex.enr`.

### Configure Output Formats

The following is the full list of configuration options for output formats and their default values:

`name`
: the output format identifier. This is used to define what output format(s) you want for your pages.

`mediaType`
: this must match the `Type` of a defined media type.

`path`
: sub path to save the output files.

`baseName`
: the base filename for the list filenames (homepage, etc.). **Default:** `index`.

`rel`
: can be used to create `rel` values in `link` tags. **Default:** `alternate`.

`protocol`
: will replace the "http://" or "https://" in your `baseURL` for this output format.

`isPlainText`
: use Go's plain text templates parser for the templates. **Default:** `false`.

`isHTML`
: used in situations only relevant for `HTML`-type formats; e.g., page aliases.

`noUgly`
: used to turn off ugly URLs If `uglyURLs` is set to `true` in your site. **Default:** `false`.

`notAlternative`
: enable if it doesn't make sense to include this format in an `AlternativeOutputFormats` format listing on `Page` (e.g., with `CSS`). Note that we use the term *alternative* and not *alternate* here, as it does not necessarily replace the other format. **Default:** `false`.

`permalinkable`
: make `.Permalink` and `.RelPermalink` return the rendering Output Format rather than main ([see below](#link-to-output-formats)). This is enabled by default for `HTML` and `AMP`. **Default:** `false`.

## Output Formats for Pages

A `Page` in Hugo can be rendered to multiple *output formats* on the file
system.

### Default Output Formats
Every `Page` has a [`Kind`][page_kinds] attribute, and the default Output
Formats are set based on that.

| Kind           | Default Output Formats |
|--------------- |----------------------- |
| `page`         | HTML                   |
| `home`         | HTML, RSS              |
| `section`      | HTML, RSS              |
| `taxonomyTerm` | HTML, RSS              |
| `taxonomy`     | HTML, RSS              |

### Customizing Output Formats

This can be changed by defining an `outputs` list of output formats in either
the `Page` front matter or in the site configuration (either for all sites or
per language).

Example from site config file:

{{< code-toggle file="config" >}}
[outputs]
  home = ["HTML", "AMP", "RSS"]
  page = ["HTML"]
{{</ code-toggle >}}


Note that in the above examples, the *output formats* for `section`,
`taxonomyTerm` and `taxonomy` will stay at their default value `["HTML",
"RSS"]`.

* The `outputs` definition is per [`Page` `Kind`][page_kinds] (`page`, `home`, `section`, `taxonomy`, or `taxonomyTerm`).
* The names (e.g. `HTML`, `AMP`) used must match the `Name` of a defined *Output Format*.
  * These names are case insensitive.
* These can be overridden per `Page` in the front matter of content files.

The following is an example of `YAML` front matter in a content file that defines output formats for the rendered `Page`:

```yaml
---
date: "2016-03-19"
outputs:
- html
- amp
- json
---
```

##  List Output formats

Each `Page` has both an `.OutputFormats` (all formats, including the current) and an `.AlternativeOutputFormats` variable, the latter of which is useful for creating a `link rel` list in your site's `<head>`:

```go-html-template
{{ range .AlternativeOutputFormats -}}
<link rel="{{ .Rel }}" type="{{ .MediaType.Type }}" href="{{ .Permalink | safeURL }}">
{{ end -}}
```

## Link to Output Formats

`.Permalink` and `.RelPermalink` on `Page` will return the first output format defined for that page (usually `HTML` if nothing else is defined). This is regardless of the template file they are being called from.

__from `single.json.json`:__
```go-html-template
{{ .RelPermalink }} > /that-page/
{{ with  .OutputFormats.Get "json" -}}
{{ .RelPermalink }} > /that-page/index.json
{{- end }}
```

In order for them to return the output format of the current template file instead, the given output format should have its `permalinkable` setting set to true.

__Same template file as above with json output format's `permalinkable` set to true:__

```go-html-template
{{ .RelPermalink }} > /that-page/index.json
{{ with  .OutputFormats.Get "html" -}}
{{ .RelPermalink }} > /that-page/
{{- end }}
```

From content files, you can use the [`ref` or `relref` shortcodes](/content-management/shortcodes/#ref-and-relref):

```go-html-template
[Neat]({{</* ref "blog/neat.md" "amp" */>}})
[Who]({{</* relref "about.md#who" "amp" */>}})
```

## Templates for Your Output Formats

A new output format needs a corresponding template in order to render anything useful.

{{% note %}}
The key distinction for Hugo versions 0.20 and newer is that Hugo looks at an output format's `Name` and MediaType's `Suffixes` when choosing the template used to render a given `Page`.
{{% /note %}}

The following table shows examples of different output formats, the suffix used, and Hugo's respective template [lookup order][]. All of the examples in the table can:

* Use a [base template][base].
* Include [partial templates][partials]

{{< datatable "output" "layouts" "Example" "OutputFormat" "Suffix" "Template Lookup Order" >}}

Hugo will now also detect the media type and output format of partials, if possible, and use that information to decide if the partial should be parsed as a plain text template or not.

Hugo will look for the name given, so you can name it whatever you want. But if you want it treated as plain text, you should use the file suffix and, if needed, the name of the Output Format. The pattern is as follows:

```
[partial name].[OutputFormat].[suffix]
```

The partial below is a plain text template (Output Format is `CSV`, and since this is the only output format with the suffix `csv`, we don't need to include the Output Format's `Name`):

```
{{ partial "mytextpartial.csv" . }}
```

[base]: /templates/base/
[config]: /getting-started/configuration/
[lookup order]: /templates/lookup/
[media type]: https://en.wikipedia.org/wiki/Media_type
[partials]: /templates/partials/
[page_kinds]: /templates/section-templates/#page-kinds
