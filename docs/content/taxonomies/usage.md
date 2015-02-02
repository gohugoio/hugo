---
date: 2014-05-26
linktitle: Usage
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
<td><pre><code>[taxonomies]
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

**taxonomy values are case insensitive**

### Front Matter Example (in TOML)

    +++
    title = "Hugo: A fast and flexible static site generator"
    tags = [ "Development", "Go", "fast", "Blogging" ]
    categories = [ "Development" ]
    series = [ "Go Web Dev" ]
    slug = "hugo"
    project_url = "https://github.com/spf13/hugo"
    +++

### Front Matter Example (in JSON)

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
        "project_url": "https://github.com/spf13/hugo"
    }
