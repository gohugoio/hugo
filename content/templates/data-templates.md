---
title: Data Templates
linktitle:
description: In addition to the Hugo's many built-in variables, you can specify your own custom data in templates or shortcodes that pull from both local and dynamic sources.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [data,dynamic,csv,json,toml,yaml]
weight: 80
draft: false
aliases: [/extras/datafiles/,/extras/datadrivencontent/,/doc/datafiles/]
toc: true
wip: true
---

<!-- begin data files -->

In addition to the [built-in variables][vars] available from Hugo, you can specify your own custom data that can be accessed via templates or shortcodes.

Hugo supports loading data from YAML, JSON, and TOML files located in the `data` directory in the root of your Hugo project.


## The Data Folder

The `data` folder is where you can store additional data for Hugo to use when generating your site. Data files aren't used to generate standalone pages - rather they're meant supplemental to the content files. This feature can extend the content in case your front matter would grow immensely. Or perhaps your want to show a larger dataset in a template (see example below). In both cases it's a good idea to outsource the data in their own file.

These files must be YAML, JSON or TOML files (using either the `.yml`, `.yaml`, `.json` or `toml` extension). The data will be accessible as a `map` in the `.Site.Data` variable.

## Data Files in Themes

Data Files can also be used in [Hugo themes][themes] but note that theme data files follow the same logic as other template files in the [Hugo lookup order][lookup] (i.e., give two files with the same name and relative path, the file in the root project `data` directory will override the file in the `themes/<THEME>/data` directory).

Therefore, theme authors should take care to not include data files that could be easily overwritten by a user who decides to [customize a theme][customize]. for theme specific data items that shouldn't be overridden, it can be wise to prefix the folder structure with a namespace, e.g. `mytheme/data/<THEME>/somekey/...`. To check if any such duplicate exists, run hugo with the `-v` flag.

**The keys in this map will be a dot chained set of _path_, _filename_ and _key_ in file (if applicable).**

This is best explained with an example:

## Example: Jaco Pastorius' Solo Discography

[Jaco Pastorius](http://en.wikipedia.org/wiki/Jaco_Pastorius_discography) was a great bass player, but his solo discography is short enough to use as an example. [John Patitucci](http://en.wikipedia.org/wiki/John_Patitucci) is another bass giant.

The example below is a bit constructed, but it illustrates the flexibility of Data Files. It uses TOML as file format.

Given the files:

* `data/jazz/bass/jacopastorius.toml`
* `data/jazz/bass/johnpatitucci.toml`

`jacopastorius.toml` contains the content below, `johnpatitucci.toml` contains a similar list:

```
discography = [
"1974 – Modern American Music … Period! The Criteria Sessions",
"1974 – Jaco",
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
```

The list of bass players can be accessed via `.Site.Data.jazz.bass`, a single bass player by adding the filename without the suffix, e.g. `.Site.Data.jazz.bass.jacopastorius`.

You can now render the list of recordings for all the bass players in a template:

```
{{ range $.Site.Data.jazz.bass }}
   {{ partial "artist.html" . }}
{{ end }}
```

And then in `partial/artist.html`:

```
<ul>
{{ range .discography }}
  <li>{{ . }}</li>
{{ end }}
</ul>
```

Discover a new favourite bass player? Just add another TOML-file.

## Example: Accessing named values in a Data File

Assuming you have the following YAML structure to your `User0123.yml` Data File located directly in `data/`

```
Name: User0123
"Short Description": "He is a **jolly good** fellow."
Achievements:
  - "Can create a Key, Value list from Data File"
  - "Learns Hugo"
  - "Reads documentation"
```

To render the `Short Description` in your `layout` File following code is required.

```
<div>Short Description of {{.Site.Data.User0123.Name}}: <p>{{ index .Site.Data.User0123 "Short Description" | markdownify }}</p></div>
```

Note the use of the `markdownify` template function. This will send the description through the Blackfriday Markdown rendering engine.

<!-- begin "Data-drive Content" page -->

## Data-Driven Content

Data-driven content with a static site generator? Yes, it is possible!

In addition to the [data files](/extras/datafiles/) feature, we have also
implemented the feature "Data-driven Content", which lets you load
any [JSON](http://www.json.org/) or
[CSV](http://en.wikipedia.org/wiki/Comma-separated_values) file
from nearly any resource.

"Data-driven Content" currently consists of two functions, `getJSON`
and `getCSV`, which are available in **all template files**.

## Implementation details

### Calling the Functions with a URL

In any HTML template or Markdown document, call the functions like this:

```golang
{{ $dataJ := getJSON "url" }}
{{ $dataC := getCSV "separator" "url" }}
```

If you use a prefix or postfix for the URL, the functions
accept [variadic arguments][variadic]:

    {{ $dataJ := getJSON "url prefix" "arg1" "arg2" "arg n" }}
    {{ $dataC := getCSV  "separator" "url prefix" "arg1" "arg2" "arg n" }}

The separator for `getCSV` must be put in the first position and can only
be one character long.

All passed arguments will be joined to the final URL; for example:

    {{ $urlPre := "https://api.github.com" }}
    {{ $gistJ := getJSON $urlPre "/users/GITHUB_USERNAME/gists" }}

will resolve internally to:

    {{ $gistJ := getJSON "https://api.github.com/users/GITHUB_USERNAME/gists" }}

Finally, you can range over an array. This example will output the
first 5 gists for a GitHub user:

    <ul>
      {{ $urlPre := "https://api.github.com" }}
      {{ $gistJ := getJSON $urlPre "/users/GITHUB_USERNAME/gists" }}
      {{ range first 5 $gistJ }}
        {{ if .public }}
          <li><a href="{{ .html_url }}" target="_blank">{{ .description }}</a></li>
        {{ end }}
      {{ end }}
    </ul>

### Example for CSV files

For `getCSV`, the one-character-long separator must be placed in the
first position followed by the URL. The following is an example of creating an HTML table in a [partial template][partials] from a published CSV:

{{% code file="layouts/partials/get-csv.html" %}}
```html
  <table>
    <thead>
      <tr>
      <th>Name</th>
      <th>Position</th>
      <th>Salary</th>
      </tr>
    </thead>
    <tbody>
    {{ $url := "http://a-big-corp.com/finance/employee-salaries.csv" }}
    {{ $sep := "," }}
    {{ range $i, $r := getCSV $sep $url }}
      <tr>
        <td>{{ index $r 0 }}</td>
        <td>{{ index $r 1 }}</td>
        <td>{{ index $r 2 }}</td>
      </tr>
    {{ end }}
    </tbody>
  </table>
```
{{% /code %}}

The expression `{{index $r number}}` must be used to output the nth-column from
the current row.

### Caching of URLs

Each downloaded URL will be cached in the default folder `$TMPDIR/hugo_cache/`. The variable `$TMPDIR` will be resolved to your system-dependent temporary directory.

With the command-line flag `--cacheDir`, you can specify any folder on your system as a caching directory.

You can also set `cacheDir` in the main configuration file.

If you don't like caching at all, you can fully disable caching with the command line flag `--ignoreCache`.

### Authentication When Using REST URLs

Currently, you can only use those authentication methods that can be put into an URL. [OAuth][] and other authentication methods are not implemented.

### Loading Local files

To load local files with `getJSON` and `getCSV`, the source files must reside within Hugo's working directory. The file extension does not matter, but the content does.

It applies the same output logic as above in [Calling the Functions with a URL](#calling-the-functions-with-a-url).

## LiveReload with Data Files

There is no chance to trigger a [LiveReload][] when the content of a URL changes. However, when a *local* file changes (i.e., `data/*` and `themes/<THEME>/data/*`), a LiveReload will be triggered. Symlinks are not supported. Note too that because downloading of data takes a while, Hugo stops processing your Markdown files until the data download has completed.

{{% warning "URL Data and LiveReload" %}}
If you change any local file and the LiveReload is triggered, Hugo will read the data-driven (URL) content from the cache. If you have disabled the cache (i.e., by running the server with `hugo server --ignoreCache`), Hugo will re-download the content every time LiveReload triggers. This can create *huge* traffic. You may reach API limits quickly.
{{% /warning %}}

## Examples

- Photo gallery JSON powered: [https://github.com/pcdummy/hugo-lightslider-example](https://github.com/pcdummy/hugo-lightslider-example)
- GitHub Starred Repositories [in a posts](https://github.com/SchumacherFM/blog-cs/blob/master/content%2Fposts%2Fgithub-starred.md) with the related [short code](https://github.com/SchumacherFM/blog-cs/blob/master/layouts%2Fshortcodes%2FghStarred.html).
- More?  Please tell us!

## Specs for Data Formats

* [TOML Spec][toml]
* [YAML Spec][yaml]
* [JSON Spec][json]
* [CSV Spec][csv]

[csv]: https://tools.ietf.org/html/rfc4180
[customize]: /themes/customizing/
[lookup]: /templates/lookup-order/
[json]: /documents/ecma-404-json-spec.pdf
[LiveReload]: /getting-started/usage/#livereload
[OAuth]: http://en.wikipedia.org/wiki/OAuth
[partials]: /templates/partials/
[themes]: /themes/
[toml]: https://github.com/toml-lang/toml
[yaml]: http://yaml.org/spec/
[variadic]: http://en.wikipedia.org/wiki/Variadic_function
[vars]: /variables/