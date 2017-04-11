---
title: Custom Output Formats
linktitle: Output Formats
description: Hugo can output content in multiple formats to make quick work of
date: 2017-03-22
publishdate: 2017-03-22
lastmod: 2017-03-22
categories: [content management]
tags: ["amp","outputs"]
menu:
  main:
    parent: "Content Management"
    weight: 18
weight: 18
draft: false
aliases: [/templates/outputs/,/extras/output-formats/,/doc/output-formats/,/doc/custom-output/]
toc: true
---

Hugo `0.20` introduced the powerful feature **Custom Output Formats**; Hugo isn't just that "static HTML with an added RSS feed" anymore. _Say hello_ to calendars, e-book formats, Google AMP, and JSON search indexes, to name a few.

This page describes how to properly configure your site with the media types and output formats, as well as where to create your templates for your custom outputs.

## Media Types

A [media type][] (also known as *MIME type* and *content type*) is a two-part identifier for file formats and format contents transmitted on the Internet.

This is the full set of built-in media types in Hugo:

{{< datatable "media" "types" "Type" "Suffix" >}}

**Note:**

* It is possible to add custom media types or change the defaults; e.g., if want to change the suffix for `text/html` to `asp`.
* The `Suffix` is the value that will be used for URLs and filenames for that media type in Hugo.
* The `Type` is the identifier that must be used when defining new/custom `Output Formats` (see below).
* The full set of media types will be registered in Hugo's built-in development server to make sure they are recognized by the browser.

To add or modify a media type, define it in a `mediaTypes` section in your [site configuration][config], either for all sites or for a given language.

Example in `config.toml`:

```toml
[mediaTypes]
[mediaTypes."text/enriched"]
suffix = "enr"
[mediaTypes."text/html"]
suffix = "asp"
```

The above example adds one new media type, `text/enriched`, and changes the suffix for the built-in `text/html` media type.

## Output Formats

Given a media type and some additional configuration, you get an `Output Format`:

This is the full set of Hugo's built-in output formats:

{{< datatable "output" "formats" "Name" "MediaType" "Path" "BaseName" "Rel" "Protocol" "IsPlainText" "IsHTML" "NoUgly">}}

* A page can be output in as many output formats as you want, and you can have an infinite amount of output formats defined **as long as they resolve to a unique path on the file system**. In the above table, the best example of this is `AMP` vs. `HTML`. `AMP` has the value `amp` for `Path` so it doesn't overwrite the `HTML` version; e.g. we can now have both `/index.html` and `/amp/index.html`.
* The `MediaType` must match the `Type` of an already defined media type.
* You can define new output formats or redefine built-in output formats; e.g., if you want to put `AMP` pages in a different path.

To add or modify an output format, define it in an `outputFormats` section in your site's [configuration file](/templates/configuration/), either for all sites or for a given language.

```toml
[outputFormats.MyEnrichedFormat]
mediaType = "text/enriched"
baseName = "myindex"
isPlainText = true
protocol = "bep://"
```

The above example is fictional, but if used for the homepage on a site with `baseURL` `http://example.org`, it will produce a plain text homepage with the URL `bep://example.org/myindex.enr`.

### Configuring Output Formats

The following is the full list of configuration options for output formats and their default values:

`Name`
: the output format identifier. This is used to define what output format(s) you want for your pages.

`MediaType`
: this must match the `Type` of a defined media type.

`Path`
: sub path to save the output files.

`BaseName`
: the base filename for the list filenames (homepage, etc.). **Default:** `index`.

`Rel`
: can be used to create `rel` values in `link` tags. **Default:** `alternate`.

`Protocol`
: will replace the "http://" or "https://" in your `baseURL` for this output format.

`IsPlainText`
: use Go's plain text templates parser for the templates. **Default:** `false`.

`IsHTML`
: used in situations only relevant for `HTML`-type formats; e.g., page aliases.

`NoUgly`
: used to turn off ugly URLs If `uglyURLs` is set to `true` in your site. **Default:** `false`.

`NotAlternative`
: enable if it doesn't make sense to include this format in an `AlternativeOutputFormats` format listing on `Page` (e.g., with `CSS`). Note that we use the term *alternative* and not *alternate* here, as it does not necessarily replace the other format. **Default:** `false`.

## Output Formats for Pages

A `Page` in Hugo can be rendered to multiple representations on the file system. By default, all pages will render as `HTML` with some of them also as `RSS` (homepage, sections, etc.).

This can be changed by defining an `outputs` list of output formats in either the `Page` front matter or in the site configuration (either for all sites or per language).

Example from site `config.toml`:

```toml
[outputs]
  home = ["HTML", "AMP", "RSS"]
  page = ["HTML"]
```

Example from site `config.yml`:

```yml
outputs:
  home: ["HTML", "AMP", "RSS"]
  page: ["HTML"]
```


* The output definition is per `Page` `Kind` (i.e, `page`, `home`, `section`, `taxonomy`, or `taxonomyTerm`).
* The names used must match the `Name` of a defined `Output Format`.
* Any `Kind` without a definition will default to `HTML`.
* These can be overridden per `Page` in the front matter of content files.
* Output formats are case insensitive.

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

## Linking to Output Formats

Each `Page` has both an `.OutputFormats` (all formats, including the current) and an `.AlternativeOutputFormats` variable, the latter of which is useful for creating a `link rel` list in your site's `<head>`:

```
{{ range .AlternativeOutputFormats -}}
<link rel="{{ .Rel }}" type="{{ .MediaType.Type }}" href="{{ .Permalink | safeURL }}">
{{ end -}}
```

Note that `.Permalink` and `.RelPermalink` on `Page` will return the first output format defined for that page (usually `HTML` if nothing else is defined).

This is how you link to a given output format:

```
{{ with  .OutputFormats.Get "json" -}}
<a href="{{ .Permalink }}">{{ .Name }}</a>
{{- end }}
```

From content files, you can use the [`ref` or `relref` shortcodes](/content-management/shortcodes/#ref-and-relref):

```
[Neat]({{</* ref "blog/neat.md" "amp" */>}})
[Who]({{</* relref "about.md#who" "amp" */>}})
```

## Templates for Your Output Formats

A new output format needs needs a corresponding template in order to render anything useful.

{{% note %}}
The key distinction for Hugo versions 0.20 and newer is that Hugo looks at an output format's `Name` and MediaType's `Suffix` when choosing the template used to render a given `Page`.
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

The partial below is a plain text template (Outpuf Format is `CSV`, and since this is the only output format with the suffix `csv`, we don't need to include the Output Format's `Name`):

```
{{ partial "mytextpartial.csv" . }}
```

{{% note %}}
Plain text partials can currently only be included in plain text templates, and vice versa. See [this issue](https://github.com/spf13/hugo/issues/3273) for some background.
{{% /note %}}

[base]: /templates/base/
[config]: /getting-started/configuration/
[lookup order]: /templates/lookup/
[media type]: https://en.wikipedia.org/wiki/Media_type
[partials]: /templates/partials/