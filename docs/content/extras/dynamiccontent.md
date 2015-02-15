---
aliases:
- /doc/dynamiccontent/
date: 2015-02-14
menu:
  main:
    parent: extras
next: /extras/highlighting
prev: /extras/datafiles
title: Dynamic Content
weight: 91
---

Dynamic content with a static site generator? Yes it is possible!

In addition to the [data files](/extras/datafiles/) feature, we have also
implemented the feature "Dynamic Content", which lets you load
any [JSON](http://www.json.org/) or
[CSV](http://en.wikipedia.org/wiki/Comma-separated_values) file
from nearly any resource.

"Dynamic Content" currently consists of two functions, `getJson`
and `getCsv`, which are available in **all template files**.

## Implementation details

### Calling the functions with an URL

In any HTML template or Markdown document call the functions like:


	{{ $dataJ := getJson "url" }}
	{{ $dataC := getCsv "separator" "url" }}


or if you use a prefix or postfix for the URL the functions
accept [variadic arguments](http://en.wikipedia.org/wiki/Variadic_function):

	{{ $dataJ := getJson "url prefix" "arg1" "arg2" "arg n" }}
	{{ $dataC := getCsv  "separator" "url prefix" "arg1" "arg2" "arg n" }}

The separator for `getCsv` must be put on the first position and can be
only one character long.

All passed arguments will be joined to the final URL, example:

	{{ $urlPre := "https://api.github.com" }}
	{{ $gistJ := getJson $urlPre "/users/GITHUB_USERNAME/gists" }}

will resolve internally to:

	{{ $gistJ := getJson "https://api.github.com/users/GITHUB_USERNAME/gists" }}

Eventually you can range over the array. This example will output the
first 5 Github gists for a user:

	<ul>
		{{ $urlPre := "https://api.github.com" }}
		{{ $gistJ := getJson $urlPre "/users/GITHUB_USERNAME/gists" }}
		{{range first 5 $gistJ }}
			{{ if .public }}
				<li><a href="{{ .html_url }}" target="_blank">{{.description}}</a></li>
			{{ end }}
		{{end}}
	</ul>


### Example for CSV files

For `getCsv` the one character long separator must be placed on the
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
	  {{ range $i, $r := getCsv $sep $url }}
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
The variable `$TMPDIR` will be resolved to your system dependent
temporary directory.

With the command line flag `--cacheDir` you can specify any folder on
your system as a caching directory.

If you don't like caching at all, you can fully disable to read from the
cache with the command line flag `--ignoreCache`. However Hugo will always
write, on each build of the site, to the cache folder (silent backup).

### Authentication when using REST URLs

Currently you can only use those authentication methods that can
be put into an URL. [OAuth](http://en.wikipedia.org/wiki/OAuth) or
other authentication methods are not implemented.

### Loading local files

To load local files with the two functions `getJson` and `getCsv` the
source files must reside within Hugos working directory. The file
extension does not matter but the content.

It applies the same output logic as in the topic: *Calling the functions with an URL*.

## Live reload

There is no chance to trigger a [LiveReload](/extras/livereload/) when
the content of an URL changes. However when a local JSON/CSV file changes
then a live reload will be triggered of course. Symlinks not supported.

**URLs and Live reload**: If you change any local file and the live reload
got triggered Hugo will either read the URL content from the cache or, if
you have disabled the cache, Hugo will re-download the content.
This can create huge traffic and you may also reach API limits quickly.

As downloading of content takes a while, Hugo stops with processing
your markdown files until the content has been downloaded.

## Examples

- Photo gallery JSON powered: [https://github.com/pcdummy/hugo-lightslider-example](https://github.com/pcdummy/hugo-lightslider-example)
- Github Starred Repositories [in a posts](https://github.com/SchumacherFM/blog-cs/blob/master/content%2Fposts%2Fgithub-starred.md) with the related [short code](https://github.com/SchumacherFM/blog-cs/blob/master/layouts%2Fshortcodes%2FghStarred.html).
- more?
