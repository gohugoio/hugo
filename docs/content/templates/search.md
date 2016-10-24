---
aliases:
- /layout/rss/
lastmod: 2016-04-02
date: 2016-04-02
linktitle: Search
menu:
  main:
    parent: layout
next: /templates/sitemap
prev: /templates/rss
title: Search template
weight: 90
---

A single search template is used to generate the `search.json` file. Hugo automatically comes with this template file. **No work is needed on the users' part unless they want to customize `search.json`, e.g. by adding more metadata for the search results.**

The search template is of the **type "node"** and has access to all the [node variables](/layout/variables/).


## Hugo's search.json template

	{{- $.Scratch.Add "search" slice -}}
	{{- range .Site.Pages -}}
	{{- $.Scratch.Add "search" (dict "title" .Title "ref" .Permalink "tags" .Params.tags "categories" .Params.categories "content" .Plain) -}}
	{{- end -}}
	{{- $.Scratch.Get "search" | jsonify -}}

First of all, the template initializes a global slice using `.Scratch` called `search`. While ranging over all pages of your website Hugo creates a `dict` with the metadata of a page. Each `dict` will be appended one by one to the slice from the beginning. Finally, we convert our data structure into valid JSON with the help of the `jsonify` template function.


## Adding a custom search template

Perhaps you want to customize the metadata of a search result. This can be done with our own `search.json` template that overwrites Hugo's default one. Just create a new template at `layouts/_default/search.json`.


## How to exclude a page type

Excluding certain page types requires just a small modification of the default template. Instead of ranging over all pages we filter out all unwanted pages by type using the `where` template function:


	{{- range where .Site.Pages "Type" "not in"  (slice "page" "CONTENT TYPE") -}}
	{{ end }}

**Don't forget** to replace `CONTENT TYPE` with it's real value.


## Configuration

By default the search index will be accissble under `www.example.com/search.json`. You can change the destination with the `searchURI` variable in the config file. In order to link the `search.json` within a template, i.e. to use [lunr.js](//lunrjs.com/), you can access its destination with the `.Site.SearchIndexLink` variable.