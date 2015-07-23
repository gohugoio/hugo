---
aliases:
- /doc/datadrivencontent/
date: 2015-02-14
menu:
  main:
    parent: extras
next: /extras/highlighting
prev: /extras/datafiles
title: Data-driven Content
weight: 91
toc: true
---

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

Each downloaded URL will be cached in the default folder `$TMPDIR/hugo_cache/`.
The variable `$TMPDIR` will be resolved to your system-dependent
temporary directory.

With the command-line flag `--cacheDir`, you can specify any folder on
your system as a caching directory.

If you don't like caching at all, you can fully disable to read from the
cache with the command line flag `--ignoreCache`. However, Hugo will always
write, on each build of the site, to the cache folder (silent backup).

### Authentication when using REST URLs

Currently, you can only use those authentication methods that can
be put into an URL. [OAuth](http://en.wikipedia.org/wiki/OAuth) or
other authentication methods are not implemented.

### Loading local files

To load local files with the two functions `getJSON` and `getCSV`, the
source files must reside within Hugo's working directory. The file
extension does not matter but the content.

It applies the same output logic as in the topic: *Calling the functions with an URL*.

## LiveReload

There is no chance to trigger a [LiveReload](/extras/livereload/) when
the content of an URL changes. However, when a local JSON/CSV file changes,
then a LiveReload will be triggered of course. Symlinks not supported.

**URLs and LiveReload**: If you change any local file and the LiveReload
got triggered, Hugo will either read the URL content from the cache or, if
you have disabled the cache, Hugo will re-download the content.
This can create huge traffic and you may also reach API limits quickly.

As downloading of content takes a while, Hugo stops with processing
your Markdown files until the content has been downloaded.

## Examples

- Photo gallery JSON powered: [https://github.com/pcdummy/hugo-lightslider-example](https://github.com/pcdummy/hugo-lightslider-example)
- GitHub Starred Repositories [in a posts](https://github.com/SchumacherFM/blog-cs/blob/master/content%2Fposts%2Fgithub-starred.md) with the related [short code](https://github.com/SchumacherFM/blog-cs/blob/master/layouts%2Fshortcodes%2FghStarred.html).
- More?  Please tell us!
