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

Besides the [data files](/extras/datafiles/) available from Hugo, you can specify your own custom data that can be accessed via templates or shortcodes.

Hugo supports loading data from [YAML](http://yaml.org/), [JSON](http://www.json.org/), and [TOML](https://github.com/toml-lang/toml) files located in the `data` directory.

**It even works with [LiveReload](/extras/livereload/).**

## The Data Folder

As explained in [Source Organization](/overview/source-directory/), the `data` folder is where you can store additional data for Hugo to use when generating your site. These files must be YAML, JSON or TOML files (using either the `.yml`, `.yaml`, `.json` or `toml` extension) and the data will be accessible as a `map` in `.Site.Data`.

**The keys in this map will be a dot chained set of _path_, _filename_ and _key_ in file (if applicable).**

This is best explained with an example:

## The Future: getSQL

The outlook to support more sources is of course implementing SQL support.

Maybe adding two new CLI switches:

	--sqlDriver=mysql|postres|mssql
	--sqlSource=string|filename

#### `--sqlDriver`

specifies the driver to use which can be one from [https://github.com/golang/go/wiki/SQLDrivers](https://github.com/golang/go/wiki/SQLDrivers)

#### `--sqlSource`

You can either provide the connection string on the command file OR an existing file which contains the connection string.

How the connection string looks like depends heavily on the used driver. For MySQL:

	hugo --sqlDriver=mysql \
	--sqlSource=username:password@protocol(address)/dbname?param=value

or with a file name:

	hugo --sqlDriver=mysql --sqlSource=path/to/myCredentials.txt

The file myCredentials.txt contains the connection string: `username:password@protocol(address)/dbname?param=value` and nothing more!



```
$data := getSQL "SELECT id,artist,genre,title from musicTable"
```
