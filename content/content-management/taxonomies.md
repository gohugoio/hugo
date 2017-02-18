---
title: Taxonomies
linktitle:
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [taxonomies,metadata,front matter,terms]
categories: [content management]
weight: 60
draft: false
slug:
aliases: [/taxonomies/overview/,/taxonomies/usage/,/indexes/overview/,/doc/indexes/,/extras/indexes]
toc: true
notes:
---

## What is a Taxonomy?

Hugo includes support for user-defined groupings of content called **taxonomies**. Taxonomies are classifications that demonstrate logical relationships between content.

{{% note %}}
Taxonomies were previously known as *indexes* in Hugo before v0.11.
{{% /note %}}

### Taxonomy Definitions

In order to effectively work with taxonomies in Hugo, it's important to first understand the language used to described different facets of usage.

* **Taxonomy:** A categorization that can be used to classify content
* **Term:** A key within that taxonomy
* **Value:** A piece of content assigned to that Term

## Default Values and URLs

Hugo natively supports taxonomies, which means there are architectural patterns and default values baked into Hugo's core. Luckily, Hugo limits these default behaviors to those that save you time as you develop your site.

### Default Taxonomies

Hugo ships with *tags* and *categories* as default taxonomies. These taxonomies are common to many website systems (e.g., WordPress, Drupal, Jekyll). Unlike these systems, Hugo makes it trivial to customize the taxonomies you will be using for your website. In addition to *tags* and *categories*, a *series* of posts for a blog is another common example of taxonomies.

### Pages Generated

When taxonomies are used---and [taxonomy templates][] are provided---Hugo will automatically create a taxonomy page listing all of the taxonomy's terms and individual pages for all content associated with the term. For example, a `categories` taxonomy will create the following pages:

* A single page at `yoursite.com/categories/` that lists all the [terms within the taxonomy][]
* [Individual taxonomy list pages][] (e.g., `/categories/development/`) for each of the terms that shows a listing of all pages marked as part of that taxonomy within any content file's [front matter][]

## Configuring Taxonomies

Taxonomies must be defined in your [website configuration][] before they can be
used throughout the site. You need to provide both the plural and
singular labels for each taxonomy.

Here is an example configuration in TOML and YAML
that specifies three taxonomies (the default two, plus `series`).

Notice the format is `singular key = "plural value"` for TOML and `singular key: "plural value"` for YAML:

### TOML Configuration

```toml
[taxonomies]
  tag = "tags"
  category = "categories"
  series = "series"
```

### YAML Configuration

```yaml
taxonomies:
  tag: "tags"
  category: "categories"
  series: "series"
```

### Preserving Taxonomy Values

By default, taxonomy names are hyphenated, lower-cased, normalized, and then fixed and title-ized on the archive page.

However, if you want to have a taxonomy value with special characters such as `Gérard Depardieu` instead of `Gerard Depardieu`, you need to set the value for `preserveTaxonomyNames` in your [site configuration](/overview/configuration/) to `true`. Hugo will then preserve special characters in taxonomy values but will still titleize the values for titles and normalize them in URLs.

Note that if you use `preserveTaxonomyNames` and intend to manually construct URLs to the archive pages, you will need to pass the taxonomy values through the [`urlize` template function][].

## Adding Taxonomies to Content

Once a taxonomy is defined at the site level, any piece of content can be assigned to it, regardless of [content type][] or [content section][].

Assigning content to a taxonomy is done in the [front matter][]. Simply create a variable with the *plural* name of the taxonomy and assign all terms you want to apply to the instance of the content type.

{{% note %}}
If you would like the ability to quickly generate content files with preconfigured taxonomies or terms, read the docs on [Hugo archetypes](/content-management/archetypes/).
{{% /note %}}

### TOML Front Matter Example

```toml
+++
title = "Hugo: A fast and flexible static site generator"
tags = [ "Development", "Go", "fast", "Blogging" ]
categories = [ "Development" ]
series = [ "Go Web Dev" ]
slug = "hugo"
project_url = "https://github.com/spf13/hugo"
+++
```

### YAML Front Matter Example

```yaml
+++
title: "Hugo: A fast and flexible static site generator"
tags: ["Development", "Go", "fast", "Blogging"]
categories: ["Development"]
categories: ["Go Web Dev"]
slug: "hugo"
project_url: "https://github.com/spf13/hugo"
+++
```

### JSON Front Matter Example

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
    "project_url": "https://github.com/spf13/hugo"
}
```

## Add Content File with Front Matter

See [project organization][].

## Example Taxonomy

For example, if I was writing about movies, I may want the following
taxonomies:

* Actors
* Directors
* Studios
* Genre
* Year
* Awards

I would then specify in each movie’s front matter the specific terms for each of those taxonomies. Hugo would then automatically create pages for each Actor, Director, Studio, Genre, Year and Award listing all of the Movies that matched that specific Actor, Director, etc.

### Taxonomy Organization

Let’s use an example to demonstrate the different labels in action.
From the perspective of the taxonomy, it could be visualized as:

```
Actor                    <- Taxonomy
    Bruce Willis         <- Term
        The Six Sense    <- Content
        Unbreakable      <- Content
        Moonrise Kingdom <- Content
    Samuel L. Jackson    <- Term
        Unbreakable      <- Content
        The Avengers     <- Content
        xXx              <- Content
```

From the perspective of the content, it would appear differently, although the data and labels used are the same:

```
Unbreakable                 <- Content
    Actors                  <- Taxonomy
        Bruce Willis        <- Term
        Samuel L. Jackson   <- Term
    Director                <- Taxonomy
        M. Night Shyamalan  <- Term
    ...
Moonrise Kingdom            <- Content
    Actors                  <- Taxonomy
        Bruce Willis        <- Term
        Bill Murray         <- Term
    Director                <- Taxonomy
        Wes Anderson        <- Term
    ...
```

[`urlize` template function]: /functions/urlize/
[content section]: /content-section/
[content type]: /content-type/
[front matter]: /content-management/front-matter/
[Individual taxonomy list pages]: /templates/taxonomy-templates/#taxonomy-templates
[project organization]: /project-organization/
[taxonomy templates]: /templates/taxonomy-templates/
[terms within the taxonomy]: /templates/taxonomy-templates/#terms-templates
[website configuration]: /project-organization/configuration/