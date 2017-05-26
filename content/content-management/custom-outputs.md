---
title: Custom Outputs
linktitle: Custom Outputs
description: Hugo can output content in multiple formats to make quick work of
date: 2017-03-22
publishdate: 2017-03-22
lastmod: 2017-03-22
categories: [content management]
tags: ["amp","outputs"]
menu:
  docs:
    parent: "content-management"
    weight: 100
  quicklinks:
weight: 100	#rem
draft: false
aliases: [/extras/custom-output-types/]
toc: true
wip: true
---

{{% warning %}}
Custom output formats is a major feature being released with v20. The following copy is taken from the original proposal and spec and is therefore far from complete.
{{% /warning %}}

## Media Type

We add a  media type (also known as MIME type and content type). This is a two-part identifier for file formats and format contents transmitted on the Internet.

For Hugo's use cases, we use the top-level type name/subtype name + suffix. An example would be `application/json+json`.

Users can define their own media types by using them in an `Output Format` definition (see below).

The full set of media types will be registered in Go's `mime` package, so they will be recognised by Hugo's development server.

## Output Format

A `Page` in Hugo can be rendered to multiple representations on the file system: All will get an HTML page and some of them will get an RSS page (home page, sections etc.).

When we now create a more formal definition for these output representations, the built-ins mentioned above will be the standard set that can be extended.

So an `OutputFormat`:

```
OutputFormat:
   Name
   MediaType
   Path
   IsPlainText (default false)
   Protocol

  # And then some optional options
  NoUglyURLs
  URI # Turn index.x into somevalue.x (similar to `RSSUri` in Hugo `0.19`)
```

So:

* `Name`: The key.
* `Path` - defaults to "", which is the root. Multiple outputs to the same suffix must be separated with a path, ie. "amp" for AMP output.
* `IsPlainText`: Whether to parse the templates with `text/template` or `html/template`.
* `Protocol`: I.e. `webcal://` for calendar files. Defaults to the `baseURL` protocol.

## Standard Output Formats

So, according to the above, the current Hugo will look like this:

| Name        | MediaType           | Path  | IsPlainText
| -------------:|-------------| -----|-----|
| HTML     | text/html+html | ""  | false |
| RSS     | application/rss+xml | ""  | false |

## Layouts

The current situation (slightly simplified):

| Kind                     | Layouts
| ----------------:|:-------------|
| home                  | index.html, _default/list.html |
| section               | section/SECTION.html, SECTION/list.html, _default/section.html, _default/list.html |
| taxonomy           | taxonomy/SINGULAR.html,_default/taxonomy.html, _default/list.html |
| taxonomyTerm  |taxonomy/SINGULAR.terms.html, _default/terms.html|
| page                  | TYPE/LAYOUT.html, _default/LAYOUT.html, _default/single.html|

The above is what the Output Format `HTML` must resolve to.

So, let us make up some other Output Formats and see how that will have to look:

| Name        | MediaType           | Path  | IsPlainText
| -------------:|-------------| -----|-----|
| JSON     | application/json+json | ""  | true |
| AMP     |  text/html+html | amp  | false |

Both of the above can be handled if we add both `Name` and the `Suffix` to the mix. Let us use the home page as an example:

| Type                    | Layouts
| -----------:|:-------------|
| JSON          | index.json.json, index.json, _default/list.json.json, _default/list.json
| AMP            | index.amp.html, index.html,  _default/list.amp.html, _default/list.html

* The above adds the lower-case `Name` as a prefix to the lookup path.
* The above also assumes that it makes sense to edit the templates with the same suffix as the end result (.html, .json etc.).

TODO: RSS, 404 etc.

## Examples

`config.toml`:

```
# Add some custom output type definitions:
[[outputFormats]]
name = "Calendar"
mediaType = "text/calendar+ics"
protocol = "webcal://"
isPlainText = true

[[outputFormats]]
name = "JSON"
mediaType = "application/json" # Will get its file suffix from the sub-type, i.e. "json"
isPlainText = true

[[outputFormats]]
name = "AMP"
mediaType = "text/html"
path = "amp"

```

Note that Hugo well hard code a predefined list of the most common output types (not sure what that would be, suggestions welcome) with the obvious identifiers and sensible defaults: So whenever you want them, you can just say "json, yaml, amp ..." etc.

Page front matter:

```
title = "My Home Page"
outputs = ["html", "rss", "json", "calendar", "amp" ]
```

About the `outputs` in the page front matter:

* If none is provided, it defaults to the current behaviour (i.e. HTML for all pages and RSS for the list pages)
* If some are provided, no defaults will be added. So, if you want the plain HTML representation, you must be explicit. This way you can have the home page as JSON only if you want.
* The names used are case-insensitive and must match either a definition in `config.toml` or the standard set.

{{% note %}}
It should also be possible to set a list of default output formats in `config.toml`, avoiding the need to repeat the `outputs` list in thousands of pages, with a way to restrict each type to a set of pages (using `Kind`, probably).
{{% /note %}}
