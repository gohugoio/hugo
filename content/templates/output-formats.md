---
aliases:
- /doc/output-formats/
- /doc/custom-output/
date: 2017-03-22T08:20:13+01:00
menu:
  main:
    parent: extras
title: Output Formats
weight: 5
toc: true
---

Hugo `0.20` introduced the powerful feature **Custom Output Formats**; Hugo isn't just that "static HTML with an added RSS feed" anymore. _Say hello_ to calendars, e-book formats, Google AMP, and JSON search indexes, to name a few.

This page describes how to properly configure your site with the media types and output formats you need.

## Media Types

A [media type](https://en.wikipedia.org/wiki/Media_type) (also known as MIME type and content type) is a two-part identifier for file formats and format contents transmitted on the Internet.

This is the full set of built-in media types in Hugo:

{{< datatable "media" "types" "Type" "Suffix" >}}

**Note:**

* It is possible to add custom media types or change the defaults (if you, say, want to change the suffix to `asp` for `text/html`).
* The `Suffix` is the value that will be used for URLs and filenames for that media type in Hugo.
* The `Type` is the identifier that must be used when defining new `Output Formats` (see below). 
* The `Delimiter` defaults to ".", but can be changed or even blanked out to support, as an example, Netlify's `_redirect` files.
* The full set of media types will be registered in Hugo's built-in development server to make sure they are recognized by the browser.

To add or modify a media type, define it in a `mediaTypes` section in your site config (either for all sites or for a given language). 

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
Given a media type and some additional configuration, you get an `Output Format`.

This is the full set of built-in output formats in Hugo:

{{< datatable "output" "formats" "Name" "MediaType" "Path" "BaseName" "Rel" "Protocol" "IsPlainText" "IsHTML" "NoUgly" "NotAlternative">}}

**Note:**

* A page can be output in as many output formats as you want, and you can have an infinite amount of output formats defined, as long as _they resolve to a unique path on the file system_. In the table above, the best example of this is `AMP` vs. `HTML`: We have given `AMP` a value for `Path` so it doesn't overwrite the `HTML` version, i.e. we can now have both `/index.html` and `/amp/index.html`.
* The `MediaType` must match the `Type` of an already defined media type (see above).
* You can define new or redefine built-in output formats (if you, as an example, want to put `AMP` pages in a different path).

To add or modify a media type, define it in a `outputFormats` section in your site config (either for all sites or for a given language). 

Example in `config.toml`:

```toml
[outputFormats.MyEnrichedFormat]
mediaType = "text/enriched" 
baseName = "myindex"
isPlainText = true
protocol = "bep://"
```

The above example is fictional, but if used for the home page on a site with `baseURL` `http://example.org`, it will produce a plain text home page with the URL `bep://example.org/myindex.enr`.

All the available configuration options for output formats and their default values:

Field | Description 
--- | ---
**Name** | The output format identifier. This is used to define what output format(s) you want for your pages. 
**MediaType**|This must match the `Type` of a defined media type. |
**Path** | Sub path to save the output files.
**BaseName** | The base filename for the list filenames (home page etc.). **Default:** _index_.
**Rel** | Can be used to create `rel` values in `link` tags. **Default:** _alternate_.
**Protocol** | Will replace the "http://" or "https://" in your `baseURL` for this output format.
**IsPlainText** | Use Go's plain text templates parser for the templates. **Default:** _false_.
**IsHTML** | Used in situations only relevant for `HTML` type of formats, page aliases being one example.|
**NoUgly** | If `uglyURLs` is enabled globally, this can be used to turn it off for a given output format. **Default:** _false_.
**NotAlternative** |  Enable if it doesn't make sense to include this format in an the `.AlternativeOutputFormats` format listing on `Page`, `CSS` being one good example. Note that we use the term "alternative" and not "alternate" here, as it does not necessarily replace the other format, it is an alternative representation.  **Default:** _false_.


## Output Formats for your pages

A `Page` in Hugo can be rendered to multiple representations on the file system: In its default configuration all will get an `HTML` page and some of them will get an `RSS` page (home page, sections etc.).

This can be changed by defining an `outputs` list of output formats in either the `Page` front matter or in the site configuration (either for all sites or per language).

Example from site config in `config.toml`:

```toml
 [outputs]
 home = [ "HTML", "AMP", "RSS"]
 page = [ "HTML"]
 ```
 Note: 
 
 * The output definition is per `Page` `Kind`(`page`, `home`, `section`, `taxonomy`, `taxonomyTerm`). 
 * The names used must match the `Name` of a defined `Output Format`.
 * Any `Kind` without a definition will get `HTML`.
 * These can be overriden per `Page` in front matter (see below).
 * When `outputs` is specified, only the formats defined in outputs will be rendered

A `Page` with `YAML` front matter defining some output formats for that `Page`:

```yaml
---
 date: "2016-03-19"
 outputs:
 - html
 - amp
 - json
 ---
 ```
 Note:
 
 * The names used for the output formats are case insensitive.
 * The first output format in the list will act as the default.
 * The default output format is used when generating links to other pages in menus, etc.
 
## Link to Output Formats
 
 `Page` has both  `.OutputFormats` (all formats including the current) and `.AlternativeOutputFormats`, the latter useful for creating a `link rel` list in your `head` section:
 
 ```
 {{ range .AlternativeOutputFormats -}}
 <link rel="{{ .Rel }}" type="{{ .MediaType.Type }}" href="{{ .Permalink | safeURL }}">
 {{ end -}}
  ```
  
Note that `.Permalink` on `RelPermalink` on `Page` will return the first output format defined for that page (usually `HTML` if nothing else is defined). 

This is how you link to a given output format:

```
{{ with  .OutputFormats.Get "json" -}}
<a href="{{ .Permalink }}">{{ .Name }}</a>
{{- end }}
```
From content files, you can use the `ref` or `relref` shortcodes:

```
[Neat]({{</* ref "blog/neat.md" "amp" */>}})
[Who]({{</* relref "about.md#who" "amp" */>}})
```

## Templates for your Output Formats

Of course, for a new Output Format to render anything useful, we need a template for it.

**The fundamental thing to understand about this is that we in `Hugo 0.20` now also look at Output Format´s `Name` and MediaType´s `Suffix` when we choose the templates to use to render a given `Page`.**

And with so many possible variations, this is best explained with some examples:


{{< datatable "output" "layouts" "Example" "OutputFormat" "Suffix" "Template Lookup Order" >}}

**Note:**

* All of the above examples can use a base template, see [Blocks]({{< relref "templates/blocks.md" >}}).
* All of the above examples can also include partials.

Hugo will now also detect the media type and output format of partials, if possible, and use that information to decide if the partial should be parsed as a plain text template or not.

Hugo will look for the name given, so you can name it whatever you want. But if you want it treated as plain text, you should use the file suffix and, if needed, the name of the Output Format (`[partial name].[OutputFormat].[suffix])`.

The partial below is a plain text template (Output Format is `CSV`, and since this is the only output format with the suffix `csv`, we don't need to include the Output Format's `Name`):

```
{{ partial "mytextpartial.csv" . }}
```

Also note that plain text partials can currently only be included in plain text templates, and vice versa. See [this issue](https://github.com/gohugoio/hugo/issues/3273) for some background.




