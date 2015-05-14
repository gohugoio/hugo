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

Dynamic content with a static site generator? Yes, it is possible!

In addition to the [data files](/extras/datafiles/) feature, we have also
implemented the feature "Dynamic Content", which lets you load
any [JSON](http://www.json.org/) or
[CSV](http://en.wikipedia.org/wiki/Comma-separated_values) file
from nearly any resource.

"Dynamic Content" currently consists of two functions, `getJSON`
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
- Github Starred Repositories [in a posts](https://github.com/SchumacherFM/blog-cs/blob/master/content%2Fposts%2Fgithub-starred.md) with the related [short code](https://github.com/SchumacherFM/blog-cs/blob/master/layouts%2Fshortcodes%2FghStarred.html).
- more?

## Implementation details for `getSql`

Hugo is now capable to connect to any MSSQL, MySQL, PostgreSQL server and of course sqlite3 databases.

### Supported Databases

The link refers to the GoLang based driver.

-  [MySQL](https://github.com/go-sql-driver/mysql/)
-  [Postgres](https://github.com/lib/pq)
-  [MSSQL](https://github.com/denisenkom/go-mssqldb)
-  [Sqlite3](https://github.com/mattn/go-sqlite3)

The default built-in driver is MySQL. Other drivers needs to be compiled by yourself.

To enable all drivers:

```
$ cd $GOPATH/src/github.com/spf13/hugo
$ go build -tags alldb .
```

To enable specific driver/s:

```
$ cd $GOPATH/src/github.com/spf13/hugo
$ go build -tags driverName1 driverName2 driverNameN .
```

Driver name can be: [mssql](https://github.com/SchumacherFM/hugo/blob/dynamicPagesWithGetSql/hugosql/mssql.go#L1),
[postgres](https://github.com/SchumacherFM/hugo/blob/dynamicPagesWithGetSql/hugosql%2Fpostgres.go#L1) or
[sqlite3](https://github.com/SchumacherFM/hugo/blob/dynamicPagesWithGetSql/hugosql%2Fsqlite.go#L1)

Maybe all drivers are included in the future. @todo Discuss this and maybe remove this self-compile section.

### Hugo Configuration

To tell Hugo where to find the database credentials you have two possibilities:

#### Via command line argument `--sqlSource`

`--sqlSource` requires the path to the file which contains the data source name (DSN).

```
$ hugo --sqlSource=/path/to/file.txt
```

#### Via environment variable `HUGO_SQL_SOURCE`

To [set the variable](http://askubuntu.com/questions/58814/how-do-i-add-environment-variables)
in e.g. the bash shell you can type:

```
$ export HUGO_SQL_SOURCE='data source name'
```

Env var is the abbreviation for environment variable.

### Data Source Name (DSN) and Driver name configuration

To make Hugo aware of which driver to use you must prepend the driver name at the
beginning of the DSN followed by an underscore character as separator.
Driver names are always lowercase: `driverName_dataSourceName`. The file extension does not matter.

To see all supported drivers run `$ hugo -h`.

Examples for a file:

```
$ hugo --sqlSource=/path/to/music_collection_dsn.txt
```

Examples for the env var:

```
$ export HUGO_SQL_SOURCE='mysql_dsn...'
$ export HUGO_SQL_SOURCE='mssql_dsn...'
$ export HUGO_SQL_SOURCE='postgres_dsn...'
$ export HUGO_SQL_SOURCE='sqlite3_dsn...'
```

The content of the file or the env var is the data source name as explained
in the documentation of each driver:

- MySQL [DSN](https://github.com/go-sql-driver/mysql/#dsn-data-source-name)
- Postgres [DSN](http://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters)
- MSSQL [DSN](https://github.com/denisenkom/go-mssqldb#connection-parameters)
- Sqlite3 [DSN](https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go#L14)

A quick DSN overview for each driver:

- mssql: `server=localhost;user id=sa;password=admin123;port=1411`
- sqlite3: `./path/to/foo.db`
- postgres: `postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full`
- mysql: `username:passw0rd@tcp(localhost:3306)/databaseName`

```
$ export HUGO_SQL_SOURCE='mysql_username:passw0rd@tcp(localhost:3306)/databaseName'
```

A file contains: `mysql_username:passw0rd@tcp(localhost:3306)/databaseName` for the MySQL driver.

Setting both values `--sqlSource` and the env var `HUGO_SQL_SOURCE`; the env var will be applied.

### getSql

`getSql` is the only SQL-able function name which can be called in every template.

`getSql` accepts multiple string arguments and returns an array with all rows from the query.

```
{{ range _, $r := getSql "SELECT * FROM gopher_locations"  }}
    ...
{{end}}
```

If you would like to use a *dynamic* query:

```
{{ $city := "Sydney" }}
{{ range _, $r := getSql "SELECT * FROM gopher_locations WHERE city=\"" $city "\""  }}
    ...
{{end}}
```

**Heads up:** There is no protection from [SQL injections](https://www.owasp.org/index.php/SQL_Injection).
You cannot have line breaks in the query parts or anywhere else.

If you would like to easily read longer queries you can put that query into a file
and provide the path to the file as an argument to `getSql`. See the following example.

File `demo_query.sql` selects data from Magento product flat table using demo data:

```
SELECT
  entity_id,
  name,
  sku,
  url_path,
  price,
  updated_at
FROM `catalog_product_flat_1`
WHERE type_id = "configurable" AND price > 80
```

This advanced `getSql` example of a short code can even sum up the price using
the [Scratch](http://gohugo.io/extras/scratch/) feature.

```
<table border="1">
  {{ $.Scratch.Set "totalSum" 0 }}
  {{ range $i, $r := getSql $myConfigGlobalPath "demo_query.sql"  }}
    {{ if eq $i 0 }}
      <thead>
	<tr>
	    <th>ID</th>
	    <th>Name</th>
	    <th>Sku</th>
	    <th>Price</th>
	    <th>Updated</th>
	</tr>
      </thead>
      <tbody>
    {{end}}
    <tr>
	<td>{{ $r.Int "entity_id" | printf "%09d" }}</td>
	<td><a href="/{{ $r.Column "url_path" }}">{{ $r.Column "name" }}</a></td>
	<td>{{ $r.Column "sku" }}</td>
	{{ $p := $r.Float "price" }}
	<td>{{ $p | printf "%.2f" }}€</td>
	{{ $.Scratch.Add "totalSum" $p }}
	<td>{{ $r.DateTime "updated_at" "2006-01-02 15:04:05.999999" | dateFormat "02/Jan/2006" }}
	  <br> {{ $r.Column "updated_at" }}
	</td>
    </tr>
  {{ end }}
  </tbody>
  <tfoot>
    <tr>
	<th>&nbsp;</th>
	<th>&nbsp;</th>
	<th>Total:</th>
	<th>{{ $.Scratch.Get "totalSum" }}€</th>
	<th>&nbsp;</th>
    </tr>
  </tfoot>
</table>
```

- `printf "%09d"` Takes care that an integer value is at minimum 9 numbers long and filled with zeros from left.
- `printf "%.2f"` rounds the price to two decimals.


#### Row functions

For each row `$r` you can use additional functions to retrieve the value from a column.

- `$r.Column "columnName"` gets the string value of a column.
- `$r.Columns` returns an array of all column names.
- `$r.JoinValues "Separator" "columnName1" "columnName2" "columnNameN"` joins the value of n-columns together using
the first argument as a separator. The separator can have nearly any length. If you pass just a `*` as second
argument then all columns will be joined: `$r.JoinValues "Separator" "*"`.
- `$r.JoinColumns "Separator"` joins all column names using a separator.
- `$r.Int "columnName"` gets the integer value of a column. On error returns 0.
- `$r.Float "columnName"` gets the floating point number of a column. On error returns 0.
- `$r.DateTime "columnName" "layout"` parses the column string according to layout into the
time object. On error returns 0000-00-00. More [info](http://golang.org/pkg/time/#example_Parse).

`$r.Int` and `$r.Float` can be perfectly used in conjunction with [printf](http://golang.org/pkg/fmt/).

To output all columns at once you should use the functions `$r.JoinValues` and `$r.JoinColumns`,
it would look like:

```
<table border="1">
  {{ range $i, $r := getSql "./demo_query.sql"  }}
    {{ if eq $i 0 }}
	<thead>
	    <tr>
		<th>{{ $r.JoinColumns "</th><th>" | safeHtml }}</th>
	    </tr>
	</thead>
     <tbody>
    {{ end }}
	<tr>
	    <td>{{ $r.JoinValues "</td><td>" "*" | safeHtml }}</td>
	</tr>
  {{ end }}
  </tbody>
</table>
```

**Heads up**: You cannot iterate over the `$r` variable.

#### Allowed SQL commands

The only possible SQL command is the `SELECT` statement. The validation checks if `SELECT` can
be recognized at the beginning of each SQL statement. Using hacks for workarounds is up to you.
