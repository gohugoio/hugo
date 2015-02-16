---
aliases:
- /doc/configuration/
date: 2013-07-01
linktitle: Configuration
menu:
  main:
    parent: getting started
next: /overview/source-directory
notoc: true
prev: /overview/usage
title: Configuring Hugo
weight: 40
---

The directory structure and templates provide the majority of the
configuration for a site. In fact, a config file isn't even needed for many
websites since the defaults follow commonly used patterns.

Hugo expects to find the config file in the root of the source directory and
will look there first for a `config.toml` file. If none is present, it will
then look for a `config.yaml` file, followed by a `config.json` file.

The config file is a site-wide config. The config file provides directions to
hugo on how to build the site as well as site-wide parameters and menus.

## Examples

The following is an example of a typical yaml config file:

    ---
    baseurl: "http://yoursite.example.com/"
    ...

The following is an example of a toml config file with some of the default values. Values under `[params]` will populate the `.Site.Params` variable for use in templates:

    contentdir = "content"
    layoutdir = "layouts"
    publishdir = "public"
    builddrafts = false
    baseurl = "http://yoursite.example.com/"
    canonifyurls = true

    [taxonomies]
      category = "categories"
      tag = "tags"
       
    [params]
      description = "Tesla's Awesome Hugo Site"
      author = "Nikola Tesla"

Here is a yaml configuration file which sets a few more options

    ---
    baseurl: "http://yoursite.example.com/"
    title: "Yoyodyne Widget Blogging"
    footnotereturnlinkcontents: "↩"
    permalinks:
      post: /:year/:month/:title/
    params:
      Subtitle: "Spinning the cogs in the widgets"
      AuthorName: "John Doe"
      GitHubUser: "spf13"
      ListOfFoo:
        - "foo1"
        - "foo2"
      SidebarRecentLimit: 5
    ...

## Configuration variables

Following is a list of Hugo-defined variables that you can configure and their current default values:

    ---
    archetypedir:               "archetype"
    # hostname (and path) to the root eg. http://spf13.com/
    baseurl:                    "" 
    # include content marked as draft
    buildDrafts:                false 
    # include content with datePublished in the future
    buildFuture:                false 
    canonifyUrls:               false
    # config file (default is path/config.yaml|json|toml)
    config:                     "config.toml"    
    contentdir:                 "content"
    dataDir:                    "data"
    defaultExtension:           "html"
    defaultLayout:              "post"
    # filesystem path to write files to
    destination:                ""    
    disableLiveReload:          false
    # Do not build RSS files
    disableRSS:                 false 
    # Do not build Sitemap file
    disableSitemap:             false 
    # edit new content with this editor, if provided
    editor:                     ""    
    footnoteAnchorPrefix:       ""
    footnoteReturnLinkContents: ""
    languageCode:               ""
    layoutdir:                  "layouts"
    # Enable Logging
    log:                        false 
    # Log File path (if set, logging enabled automatically)
    logFile:                    ""    
    # "yaml", "toml", "json"
    metaDataFormat:             "toml" 
    newContentEditor:           ""
    # Don't sync modification time of files
    noTimes:                    false 
    paginate:                   10
    paginatePath:               "page"
    permalinks:         
    # Pluralize titles in lists using inflect
    pluralizeListTitles:         true 
    publishdir:                 "public"
    # color-codes for highlighting derived from this style
    pygmentsStyle:              "monokai"
    # true: use pygments-css or false: color-codes directly
    pygmentsUseClasses:         false 
    sitemap:                    ""
    # filesystem path to read files relative from 
    source:                     ""    
    staticdir:                  "static"
    # display memory and timing of different steps of the program
    stepAnalysis:               false 
    # theme to use (located in /themes/THEMENAME/)
    theme:                      ""    
    title:                      ""
    # if true, use /filename.html instead of /filename/
    uglyUrls:                   false 
    # verbose output
    verbose:                    false 
    # verbose logging
    verboseLog:                 false 
    # watch filesystem for changes and recreate as needed
    watch:                      false 
    ---




## Configure Blackfriday rendering

[Blackfriday](https://github.com/russross/blackfriday) is the [Markdown](http://daringfireball.net/projects/markdown/) rendering engine used in Hugo. The Blackfriday configuration in Hugo is mostly a set of sane defaults that should fit most use cases.

But Hugo does expose some options---as listed in the table below, matched with the corresponding flag in the Blackfriday source ([html.go](https://github.com/russross/blackfriday/blob/master/html.go) and [markdown.go](https://github.com/russross/blackfriday/blob/master/markdown.go)):

<table class="table table-bordered">
<thead>
<tr>
<th>Flag</th><th>Default</th><th>Blackfriday flag</th>
</tr>
</thead>

<tbody>
<tr>
<td><code>angledQuotes</code></td>
<td><code>false</code></td>
<td><code>HTML_SMARTYPANTS_ANGLED_QUOTES</code></td>
</tr>
<tr>
<td class="purpose-title">Purpose:</td>
<td class="purpose-description" colspan="2">Enable smart angled double quotes <small>(e.g.&nbsp;<code>"Hugo"</code> renders to «Hugo» instead of “Hugo”)</small></td>
</tr>

<tr>
<td><code>fractions</code></td>
<td><code>true</code></td>
<td><code>HTML_SMARTYPANTS_FRACTIONS</code></td>
</tr>
<tr>
<td class="purpose-title">Purpose:</td>
<td class="purpose-description" colspan="2">Enable smart fractions
<small>(e.g.&nbsp;<code>5/12</code> renders to <sup>5</sup>&frasl;<sub>12</sub> (<code>&lt;sup&gt;5&lt;/sup&gt;&amp;frasl;&lt;sub&gt;12&lt;/sub&gt;</code>))
<strong>Caveat:</strong> Even with <code>fractions = false</code>,
Blackfriday would still convert 1/2, 1/4 and 3/4 to ½&nbsp;(<code>&amp;frac12;</code>),
¼&nbsp;(<code>&amp;frac14;</code>) and ¾&nbsp;(<code>&amp;frac34;</code>) respectively,
but only these three.</small></td>
</tr>

<tr>
<td><code>plainIdAnchors</code></td>
<td><code>false</code></td>
<td><code>FootnoteAnchorPrefix</code> and <code>HeaderIDSuffix</code></td>
</tr>
<tr>
<td class="purpose-title">Purpose:</td>
<td class="purpose-description" colspan="2">If <code>true</code>, then header and footnote IDs are generated without the document ID <small>(e.g.&nbsp;<code>#my-header</code> instead of <code>#my-header:bec3ed8ba720b9073ab75abcf3ba5d97</code>)</small></td>
</tr>

<tr>
<td><code>extensions</code></td>
<td><code>[]</code></td>
<td><code>EXTENSION_*</code></td>
</tr>
<tr>
<td class="purpose-title">Purpose:</td>
<td class="purpose-description" colspan="2">Use non-default additional extensions <small>(e.g.&nbsp;Add <code>"hardLineBreak"</code> to use <code>EXTENSION_HARD_LINE_BREAK</code>)</small></td>
</tr>
</tbody>
</table>


**Note** that these flags must be grouped under the `blackfriday` key and can be set on **both site and page level**. If set on page, it will override the site setting.  Example:

<table class="table">
<thead>
<tr>
<th>TOML</th><th>YAML</th>
</tr>
</thead>
<tbody>
<tr>
<td><pre><code>[blackfriday]
  angledQuotes = true
  fractions = false
  plainIdAnchors = true
  extensions = ["hardLineBreak"]
</code></pre></td>
<td><pre><code>blackfriday:
  angledQuotes: true
  fractions: false
  plainIdAnchors: true
  extensions:
    - hardLineBreak
</code></pre></td>
</tr>
</tbody>
</table>

## Notes

Config changes are not reflected with [LiveReload](/extras/livereload/).

Please restart `hugo server --watch` whenever you make a config change.
