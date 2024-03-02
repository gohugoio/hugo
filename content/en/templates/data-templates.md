---
title: Data templates
description: In addition to Hugo's built-in variables, you can specify your own custom data in templates or shortcodes that pull from both local and dynamic sources.
categories: [templates]
keywords: [data,dynamic,csv,json,toml,yaml,xml]
menu:
  docs:
    parent: templates
    weight: 150
weight: 150
toc: true
aliases: [/extras/datafiles/,/extras/datadrivencontent/,/doc/datafiles/]
---

Hugo supports loading data from YAML, JSON, XML, and TOML files located in the `data` directory at the root of your Hugo project.

{{< youtube FyPgSuwIMWQ >}}

## The data directory

The `data` directory should store additional data for Hugo to use when generating your site.

Data files are not for generating standalone pages. They should supplement content files by:

- Extending the content when the front matter fields grow out of control, or
- Showing a larger dataset in a template (see the example below).

In both cases, it's a good idea to outsource the data in their (own) files.

These files must be YAML, JSON, XML, or TOML files (using the `.yml`, `.yaml`, `.json`, `.xml`, or `.toml` extension). The data will be accessible as a `map` in the `.Site.Data` variable.

To access the data using the `site.Data.filename` notation, the file name must begin with an underscore or a Unicode letter, followed by zero or more underscores, Unicode letters, or Unicode digits. For example:

- `123.json` - Invalid
- `x123.json` - Valid
- `_123.json` - Valid

To access the data using the [`index`](/functions/collections/indexfunction) function, the file name is irrelevant. For example:

Data file|Template code
:--|:--
`123.json`|`{{ index .Site.Data "123" }}`
`x123.json`|`{{ index .Site.Data "x123" }}`
`_123.json`|`{{ index .Site.Data "_123" }}`
`x-123.json`|`{{ index .Site.Data "x-123" }}`

## Data files in themes

Data Files can also be used in themes.

However, note that the theme data files are merged with the project directory taking precedence. That is, Given two files with the same name and relative path, the data in the file in the root project `data` directory will override the data from the file in the `themes/<THEME>/data` directory *for keys that are duplicated*).

Therefore, theme authors should be careful not to include data files that could be easily overwritten by a user who decides to [customize a theme][customize]. For theme-specific data items that shouldn't be overridden, it can be wise to prefix the folder structure with a namespace; e.g. `mytheme/data/<THEME>/somekey/...`. To check if any such duplicate exists, run hugo with the `-v` flag.

The keys in the map created with data templates from data files will be a dot-chained set of `path`, `filename`, and `key` in the file (if applicable).

This is best explained with an example:

## Examples

### Jaco Pastorius' Solo Discography

[Jaco Pastorius](https://en.wikipedia.org/wiki/Jaco_Pastorius_discography) was a great bass player, but his solo discography is short enough to use as an example. [John Patitucci](https://en.wikipedia.org/wiki/John_Patitucci) is another bass giant.

The example below is a bit contrived, but it illustrates the flexibility of data Files. This example uses TOML as its file format with the two following data files:

* `data/jazz/bass/jacopastorius.toml`
* `data/jazz/bass/johnpatitucci.toml`

`jacopastorius.toml` contains the content below. `johnpatitucci.toml` contains a similar list:

{{< code-toggle file=data/jazz/bass/jacopastorius >}}
discography = [
"1974 - Modern American Music … Period! The Criteria Sessions",
"1974 - Jaco",
"1976 - Jaco Pastorius",
"1981 - Word of Mouth",
"1981 - The Birthday Concert (released in 1995)",
"1982 - Twins I & II (released in 1999)",
"1983 - Invitation",
"1986 - Broadway Blues (released in 1998)",
"1986 - Honestly Solo Live (released in 1990)",
"1986 - Live In Italy (released in 1991)",
"1986 - Heavy'n Jazz (released in 1992)",
"1991 - Live In New York City, Volumes 1-7.",
"1999 - Rare Collection (compilation)",
"2003 - Punk Jazz: The Jaco Pastorius Anthology (compilation)",
"2007 - The Essential Jaco Pastorius (compilation)"
]
{{< /code-toggle >}}

The list of bass players can be accessed via `.Site.Data.jazz.bass`, a single bass player by adding the file name without the suffix, e.g. `.Site.Data.jazz.bass.jacopastorius`.

You can now render the list of recordings for all the bass players in a template:

```go-html-template
{{ range $.Site.Data.jazz.bass }}
  {{ partial "artist.html" . }}
{{ end }}
```

And then in the `partials/artist.html`:

```go-html-template
<ul>
{{ range .discography }}
  <li>{{ . }}</li>
{{ end }}
</ul>
```

Discover a new favorite bass player? Just add another `.toml` file in the same directory.

### Accessing named values in a data file

Assume you have the following data structure in your `user0123` data file located directly in `data/`:

{{< code-toggle file=data/user0123 >}}
Name: User0123
"Short Description": "He is a **jolly good** fellow."
Achievements:
  - "Can create a Key, Value list from Data File"
  - "Learns Hugo"
  - "Reads documentation"
{{</ code-toggle >}}

You can use the following code to render the `Short Description` in your layout:

```go-html-template
<div>Short Description of {{ .Site.Data.user0123.Name }}: <p>{{ index .Site.Data.user0123 "Short Description" | markdownify }}</p></div>
```

Note the use of the [`markdownify`] function. This will send the description through the Markdown rendering engine.

## Remote data

Retrieve remote data using these template functions:

- [`resources.GetRemote`](/functions/resources/getremote) (recommended)
- [`data.GetCSV`](/functions/data/getcsv)
- [`data.GetJSON`](/functions/data/getjson)

## LiveReload with data files

There is no chance to trigger a [LiveReload] when the content of a URL changes. However, when a *local* file changes (i.e., `data/*` and `themes/<THEME>/data/*`), a LiveReload will be triggered. Symlinks are not supported. Note too that because downloading data takes a while, Hugo stops processing your Markdown files until the data download has been completed.

{{% note %}}
If you change any local file and the LiveReload is triggered, Hugo will read the data-driven (URL) content from the cache. If you have disabled the cache (i.e., by running the server with `hugo server --ignoreCache`), Hugo will re-download the content every time LiveReload triggers. This can create *huge* traffic. You may reach API limits quickly.
{{% /note %}}

## Examples of data-driven content

- Photo gallery JSON powered: [https://github.com/pcdummy/hugo-lightslider-example](https://github.com/pcdummy/hugo-lightslider-example).
- GitHub Starred Repositories [in a post](https://github.com/SchumacherFM/blog-cs/blob/master/content%2Fposts%2Fgithub-starred.md) using data-driven content in a [custom short code](https://github.com/SchumacherFM/blog-cs/blob/master/layouts%2Fshortcodes%2FghStarred.html).
- Importing exported social media data from popular services using [https://github.com/ttybitnik/diego](https://github.com/ttybitnik/diego).

## Specs for data formats

* [TOML Spec][toml]
* [YAML Spec][yaml]
* [JSON Spec][json]
* [CSV Spec][csv]
* [XML Spec][xml]

[config]: /getting-started/configuration/
[csv]: https://tools.ietf.org/html/rfc4180
[customize]: /hugo-modules/theme-components/
[json]: https://www.ecma-international.org/publications/files/ECMA-ST/ECMA-404.pdf
[LiveReload]: /getting-started/usage/#livereload
[lookup]: /templates/lookup-order/
[`markdownify`]: /functions/transform/markdownify/
[OAuth]: https://en.wikipedia.org/wiki/OAuth
[partials]: /templates/partials/
[toml]: https://toml.io/en/latest
[variadic]: https://en.wikipedia.org/wiki/Variadic_function
[vars]: /methods/
[yaml]: https://yaml.org/spec/
[xml]: https://www.w3.org/XML/
