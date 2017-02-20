---
title: Data Templates
linktitle:
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
weight:
categories: [templates]
tags: [data,dynamic,csv,json,toml,yaml]
draft: false
aliases: [/extras/datafiles/,/extras/datadrivencontent/,/doc/datafiles/]
toc: false
needsreview: true
---

<!-- begin data files -->

In addition to the [built-in variables](/templates/variables/) available from Hugo, you can specify your own custom data that can be accessed via templates or shortcodes.

Hugo supports loading data from [YAML](http://yaml.org/), [JSON](http://www.json.org/), and [TOML](https://github.com/toml-lang/toml) files located in the `data` directory.

**It even works with [LiveReload](/extras/livereload/).**

Data Files can also be used in [themes](/themes/overview/), but note: If the same `key` is used in both the main data folder and in the theme's data folder, the main one will win. So, for theme authors,  for theme specific data items that shouldn't be overridden, it can be wise to prefix the folder structure with a namespace, e.g. `mytheme/data/mytheme/somekey/...`. To check if any such duplicate exists, run hugo with the `-v` flag, e.g. `hugo -v`.

## The Data Folder

The `data` folder is where you can store additional data for Hugo to use when generating your site. Data files aren't used to generate standalone pages - rather they're meant supplemental to the content files. This feature can extend the content in case your frontmatter would grow immensely. Or perhaps your want to show a larger dataset in a template (see example below). In both cases it's a good idea to outsource the data in their own file.

These files must be YAML, JSON or TOML files (using either the `.yml`, `.yaml`, `.json` or `toml` extension) and the data will be accessible as a `map` in `.Site.Data`.

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

Data-driven content with a static site generator? Yes, it is possible!

In addition to the [data files](/extras/datafiles/) feature, we have also
implemented the feature "Data-driven Content", which lets you load
any [JSON](http://www.json.org/) or
[CSV](http://en.wikipedia.org/wiki/Comma-separated_values) file
from nearly any resource.

"Data-driven Content" currently consists of two functions, `getJSON`
and `getCSV`, which are available in **all template files**.

## Implementation details

### Calling the functions with an URL

In any HTML template or Markdown document, call the functions like this:

    {{ $dataJ := getJSON "url" }}
    {{ $dataC := getCSV "separator" "url" }}

or, if you use a prefix or postfix for the URL, the functions
accept [variadic arguments](http://en.wikipedia.org/wiki/Variadic_function):

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

For `getCSV`, the one-character long separator must be placed in the
first position followed by the URL.

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

The expression `{{index $r number}}` must be used to output the nth-column from
the current row.

### Caching of URLs

Each downloaded URL will be cached in the default folder `$TMPDIR/hugo_cache/`. The variable `$TMPDIR` will be resolved to your system-dependent temporary directory.

With the command-line flag `--cacheDir`, you can specify any folder on your system as a caching directory.

You can also set `cacheDir` in the main configuration file.

If you don't like caching at all, you can fully disable caching with the command line flag `--ignoreCache`.

### Authentication when using REST URLs

Currently, you can only use those authentication methods that can be put into an URL. [OAuth](http://en.wikipedia.org/wiki/OAuth) or other authentication methods are not implemented.

### Loading local files

To load local files with the two functions `getJSON` and `getCSV`, the source files must reside within Hugo's working directory. The file extension does not matter but the content does.

It applies the same output logic as in the topic: *Calling the functions with an URL*.

## LiveReload

There is no chance to trigger a [LiveReload](/extras/livereload/) when the content of an URL changes. However, when a local JSON/CSV file changes, then a LiveReload will be triggered of course. Symlinks are not supported.

{{% note "URLs and LiveReload" %}}
If you change any local file and the LiveReload is triggered, Hugo will either read the URL content from the cache or, if you have disabled the cache, Hugo will re-download the content. This can create huge traffic and you may also reach API limits quickly.
{{% /note %}}

As downloading of content takes a while, Hugo stops processing
your Markdown files until the content has been downloaded.

## Examples

- Photo gallery JSON powered: [https://github.com/pcdummy/hugo-lightslider-example](https://github.com/pcdummy/hugo-lightslider-example)
- GitHub Starred Repositories [in a posts](https://github.com/SchumacherFM/blog-cs/blob/master/content%2Fposts%2Fgithub-starred.md) with the related [short code](https://github.com/SchumacherFM/blog-cs/blob/master/layouts%2Fshortcodes%2FghStarred.html).
- More?  Please tell us!