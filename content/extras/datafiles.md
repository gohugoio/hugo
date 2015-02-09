---
aliases:
- /doc/scratch/
date: 2015-01-22
menu:
  main:
    parent: extras
next: /extras/highlighting
prev: /extras/scratch
title: Data Files
weight: 90
---

In addition to the [built-in variables](/templates/variables/) available from Hugo, you can specify your own custom data that can be accessed via templates or shortcodes.

Hugo supports loading data from [YAML](http://yaml.org/), [JSON](http://www.json.org/), and [TOML](https://github.com/toml-lang/toml) files located in the `data` directory.

**It even works with [LiveReload](/extras/livereload/).**

## The Data Folder

As explained in [Source Organization](/overview/source-directory/), the `data` folder is where you can store additional data for Hugo to use when generating your site. These files must be YAML, JSON or TOML files (using either the `.yml`, `.yaml`, `.json` or `toml` extension) and the data will be accessible as a `map` in `.Site.Data`.

**The keys in this map will be a dot chained set of _path_, _filename_ and _key_ in file (if applicable).** 

This is best explained with an example:

## Example: Jaco Pastorius' Solo Discography

[Jaco Pastorius](http://en.wikipedia.org/wiki/Jaco_Pastorius_discography) was a great bass player, but his solo discography is short enough to use as an example.

The example below uses TOML as file format.

Given the file `data/jazz/bass/jacopastorius.toml` with the content below:

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

This list can be accessed via `.Site.Data.jazz.bass.jacopastorius.discography`.

You can now render the list of recordings in a template:

```
<ul>
{{ range $.Site.Data.jazz.bass.jacopastorius.discography }}
  <li>{{ . }}</li>
{{ end }}
</ul>
```


