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

Besides the [data files](/extras/datafiles/) feature, we have also implemented the feature "Dynamic Content"
which lets you load any [JSON](http://www.json.org/) or [CSV](http://en.wikipedia.org/wiki/Comma-separated_values)
file from nearly any resource.

"Dynamic Content" consists at the moment of two functions `getJson` and `getCsv` which are available in
**all template files**.

## Implementation details

### Calling the functions with an URL

In any template file call the functions like:

```
	{{ $dataJ := getJson "url" }}
	{{ $dataC := getCsv "separator" "url" }}
```

or if you use a prefix or postfix for the URL the functions
accepts [variadic arguments](http://en.wikipedia.org/wiki/Variadic_function):

```
	{{ $dataJ := getJson "url prefix" "arg1" "arg2" "arg n" }}
	{{ $dataC := getCsv  "separator" "url prefix" "arg1" "arg2" "arg n" }}
```

The separator for `getCsv` must be put on the first position and can be only one character long.

All passed arguments will be joined to the final URL, example:

```
	{{ $urlPre := "https://api.github.com" }}
	{{ $gistJ := getJson $urlPre "/users/GITHUB_USERNAME/gists" }}
```

will resolve internally to:

```
	{{ $gistJ := getJson "https://api.github.com/users/GITHUB_USERNAME/gists" }}
```

Eventually you can range or the map/array/slice. This example will output the first 5 Github gists for a user:

```
      <ul>
	{{ $urlPre := "https://api.github.com" }}
	{{ $gistJ := getJson $urlPre "/users/GITHUB_USERNAME/gists" }}
	{{range first 5 $gistJ }}
	  {{ if .public }}
	    <li><a href="{{ .html_url }}" target="_blank">{{.description}}</a></li>
	  {{ end }}
	{{end}}
      </ul>
```

### Example for CSV files

For `getCsv` the one character long separator must be placed on the first position followed by the URL.

```
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
```

### Caching of URLs

Each downloaded URL will be cached in the default folder `$TMPDIR/hugo_cache/`. The variable `$TMPDIR` will be
resolved to your system dependent temporary directory.

With the command line flag `--cacheDir` you can specify any folder on your system as a caching directory.

If you don't like caching at all you can fully disabled to read from the cache with the command line
flag `--ignoreCache`. But hugo will still write on each build of the site to the cache folder (silent backup).

### Authentication when using REST URLs

At the moment you can only use those authentication methods which can be put into an URL.
OAuth or other stuff is not implemented.

### Loading local files

To load local files with the two functions `getJson` and `getCsv` the source files must reside within
Hugos working directory. The file extension does not matter but the content.

It applies the same logic as in the topic: *Calling the functions with an URL*.

## Live reload

There is not chance to trigger a [LiveReload](/extras/livereload/) when the content of an URL changes.
But when a local JSON/CSV file changes then of course a live reload will be triggered. Symlinks not supported.

**URLs and Live reload**: If you change any local file and the live reload got trigger Hugo will
either read the URL content from the cache or if you have disabled the cache Hugo will re-download the content.
This can create a huge traffic and also you may reach API limits quickly.

As downloading of content takes a while Hugo stops with processing your markdown files until the content
has been downloaded.

## The Future:

### YAML and TOML

If the community demands the implementation of *getYaml* [YAML](http://yaml.org/) or
*getToml* [TOML](https://github.com/toml-lang/toml) these functions will for sure follow.

### getSql

The outlook to support more sources is of course implementing SQL support.

Maybe adding a new CLI option:

	--sqlSource=path/to/filename.ext

#### `--sqlSource`

The file must start with `[mysql|postres|mssql|...]_whatever.ext`

The part until the first underscore specifies the driver to use which can be one
from [https://github.com/golang/go/wiki/SQLDrivers](https://github.com/golang/go/wiki/SQLDrivers).

The file itself contains only the connection string and no other comments or characters.

How the connection string looks like depends heavily on the used driver. For MySQL:

	hugo --sqlSource=path/to/mysql_Credentials.txt

The file `mysql_Credentials.txt` contains the connection string:
`username:password@protocol(address)/dbname?param=value` and nothing more!

In your template you can process as with the `getCsv` function:

```
$data := getSql "SELECT id,artist,genre,title from musicTable"
```

Abusing `getSql` with [DML](http://en.wikipedia.org/wiki/Data_manipulation_language) or
[DDL](http://en.wikipedia.org/wiki/Data_definition_language) statements is up to you.
