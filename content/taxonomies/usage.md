---
lastmod: 2015-12-23
date: 2014-05-26
linktitle: Usage
toc: true
menu:
  main:
    parent: taxonomy
next: /taxonomies/displaying
prev: /taxonomies/overview
title: Using Taxonomies
weight: 15
---

## Defining taxonomies for a site

Taxonomies must be defined in the site configuration before they can be
used throughout the site. You need to provide both the plural and
singular labels for each taxonomy.

Here is an example configuration in TOML and YAML
that specifies three taxonomies (the default two, plus `series`).

Notice the format is <code><strong>singular key</strong> = &quot;<em>plural value</em>&quot;</code> for TOML,
or <code><strong>singular key</strong>: &quot;<em>plural value</em>&quot;</code> for YAML:

<table class="table">
<thead>
<tr>
<th>config.toml excerpt:</th><th>config.yaml excerpt:</th>
</tr>
</thead>
<tbody>
<tr valign="top">
<td><pre><code class="language-toml">[taxonomies]
tag = "tags"
category = "categories"
series = "series"
</code></pre></td>
<td><pre><code class="language-yaml">taxonomies:
  tag: "tags"
  category: "categories"
  series: "series"
</code></pre></td>
</tr>
</tbody>
</table>

## Assigning taxonomy values to content

Once an taxonomy is defined at the site level, any piece of content
can be assigned to it regardless of content type or section.

Assigning content to an taxonomy is done in the front matter.
Simply create a variable with the *plural* name of the taxonomy
and assign all terms you want to apply to this content.

## Preserving taxonomy values

By default, taxonomy names are hyphenated, lower-cased and normalized, and then
fixed and titleized on the archive page.

However, if you want to have a taxonomy value with special characters
such as `Gérard Depardieu` instead of `Gerard Depardieu`,
you need to set the `preserveTaxonomyNames` [site configuration](/overview/configuration/) variable to `true`.
Hugo will then preserve special characters in taxonomy values
but will still titleize the values for titles and normalize them in URLs.

Note that if you use `preserveTaxonomyNames` and intend to manually construct URLs to the archive pages,
you will need to pass the taxonomy values through the `urlize` template function.

## Front Matter Example (in TOML)

```toml
+++
title = "Hugo: A fast and flexible static site generator"
tags = [ "Development", "Go", "fast", "Blogging" ]
categories = [ "Development" ]
series = [ "Go Web Dev" ]
slug = "hugo"
project_url = "https://github.com/gohugoio/hugo"
+++
```

## Front Matter Example (in JSON)

```json
{
    "title": "Hugo: A fast and flexible static site generator",
    "tags": [
        "Development",
        "Go",
        "fast",
        "Blogging"
    ],
    "categories" : [
        "Development"
    ],
    "series" : [
        "Go Web Dev"
    ],
    "slug": "hugo",
    "project_url": "https://github.com/gohugoio/hugo"
}
```

## Add content file with frontmatter

See [Source Organization]({{< relref "overview/source-directory.md#content-for-home-page-and-other-list-pages" >}}).
