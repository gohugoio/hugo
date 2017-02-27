---
title: Content Organization
linktitle: Content Organization
description: Hugo assumes that the same structure that works to organize your source content is used to organize the rendered site.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [content management]
tags: [sections,content,organization,fundamentals]
weight: 10
draft: false
aliases: [/content-management/sections,/content/sections/]
toc: true
needsreview: true
---

## Introduction

Hugo uses files (see [Hugo's supported content formats][formats]) with headers called [front matter][]. By default, Hugo assumes the same structure that works to organize your content should be used to organize your rendered website. This is done in an effort to reduce configuration. However, this convention can be overridden through additional configuration in the front matter, as well as through Hugo's extensive features related to [URL management][urls].

## Organizing Content Source

In Hugo, the content should be organized in a manner that reflects the rendered website. Without any additional configuration, the following will just work. Hugo supports content nested at any level, but the top level (i.e. `content/<directories>*``) is special in Hugo and is considered the content [section][].

## Destinations

Hugo believes that you organize your content with a purpose. The same structure that works to organize your source content is used to organize the rendered site. As displayed above, the organization of the source content will be mirrored in the destination.

Notice that the first level `about/` page URL was created using a directory named "about" with a single `_index.md` file inside. Find out more about `_index.md` specifically in [content for the homepage and other list pages](https://gohugo.io/overview/source-directory#content-for-home-page-and-other-list-pages).

There are times when one would need more control over their content. In these cases, there are a variety of things that can be specified in the front matter to determine the destination of a specific piece of content.

The following items are defined in order; latter items in the list will override earlier settings.

### `filename`

This isn't in the front matter, but is the actual name of the file minus the extension. This will be the name of the file in the destination.

### `slug`

Defined in the front matter, the `slug` can take the place of the filename for the destination.

### `filepath`

The actual path to the file on disk. Destination will create the destination with the same path. Includes [section](/content/sections/).

### `section`

`section` is determined by its location on disk and *cannot* be specified in the front matter. See [section](/content/sections/).

### `type`

`type` is also determined by its location on disk but, unlike `section`, it *can* be specified in the front matter. See [type](/content/types/).

### `path`

`path` can be provided in the front matter. This will replace the actual path to the file on disk. Destination will create the destination with the same path. Includes [section](/content/sections/).

### `url`

A complete URL can be provided. This will override all the above as it pertains to the end destination. This must be the path from the baseURL (starting with a "/"). When a `url` is provided, it will be used exactly. Using `url` will ignore the `--uglyURLs` setting.


## Path Breakdown in Hugo

### Content

```bash
.             path           slug
.       ⊢-------^----⊣ ⊢------^-------⊣
content/extras/indexes/category-example/index.html
```

```bash
.       section              slug
.       ⊢--^--⊣        ⊢------^-------⊣
content/extras/indexes/category-example/index.html
```

```bash
.       section  slug
.       ⊢--^--⊣⊢--^--⊣
content/extras/indexes/index.html
```

### Destination

```bash
           permalink
⊢--------------^-------------⊣
http://spf13.com/projects/hugo
```

```bash
   baseURL       section  slug
⊢-----^--------⊣ ⊢--^---⊣ ⊢-^⊣
http://spf13.com/projects/hugo
```

```bash
   baseURL       section          slug
⊢-----^--------⊣ ⊢--^--⊣        ⊢--^--⊣
http://spf13.com/extras/indexes/example
```

```bash
   baseURL            path       slug
⊢-----^--------⊣ ⊢------^-----⊣ ⊢--^--⊣
http://spf13.com/extras/indexes/example
```

```bash
   baseURL            url
⊢-----^--------⊣ ⊢-----^-----⊣
http://spf13.com/projects/hugo
```

```bash
   baseURL               url
⊢-----^--------⊣ ⊢--------^-----------⊣
http://spf13.com/extras/indexes/example
```

#### `section`

A section is the content type the piece of content is assigned to by default. `section` is determined by the following:

* content location within the project's directory structure
* front matter overrides

#### `slug`

A content's `slug` is either `name.extension` or `name/`. `slug` is determined by the following:

* the name of the content file (e.g., `content-name.md`)
* front matter overrides

#### `path`

A content's `path` is determined by the section's path to the file. `path`

* is based on the path to the content's location
* excludes the slug

#### `url`

The `url` is the relative URL for the piece of content. The `url`

* is defined in front matter
* overrides all the above

## \_index.md and "Everything is a Page"

As of version v0.18, Hugo now treats "[everything as a page](http://bepsays.com/en/2016/12/19/hugo-018/)". This allows you to add content and front matter to any page, including list pages like [sections][sectiontemplates], [taxonomy list pages][taxonomytemplates], [taxonomy terms pages](/templates/terms/) and even to potential "special case" pages like the [homepage][].

In order to take advantage of this behavior, you need to do a few things.

1. Create an `_index.md` file that contains the front matter and content you would like to apply.

2. Place the `_index.md` file in the correct place in the directory structure.

3. Ensure that the respective template is configured to display `{{ .Content }}` if you wish for the content of the `_index.md` file to be rendered on the respective page.

### How `_index.md` Works

Before continuing, it's important to know that this page must reference certain templates to describe how the \_index.md page will be rendered. Hugo has a multitude of possible templates that can be used and placed in various places (think theme templates for instance). For simplicity/brevity the default/top level template location will be used to refer to the entire range of places the template can be placed.

If this is confusing or you are unfamiliar with Hugo's template hierarchy, visit the various template pages listed below. You may need to find the 'active' template responsible for any particular page on your own site by going through the template hierarchy and matching it to your particular setup/theme you are using.

- [Homepage template](/templates/homepage/)
- [Content List templates](/templates/list/)
- [Single Content templates](/templates/content/)
- [Taxonomy Terms templates](/templates/terms/)

Now that you've got a handle on templates lets recap some Hugo basics to understand how to use an \_index.md file with a List page.

1. Sections and Taxonomies are 'List' pages, NOT single pages.
2. List pages are rendered using the template heirarchy found in the [Content - List Template](http://localhost:1313/templates/list/) docs.
3. The homepage, though technically a list page, can have [it's own template](/templates/homepage/) at layouts/index.html rather than \_default/list.html. Many themes exploit this behavior so you are likely to encounter this specific use case.
4. Taxonomy terms pages are "lists of metadata" and not lists of content and therefore [have their own templates](/templates/terms/).

Let's put all this information together:

* `_index.md` files are used in list pages, terms pages, or the homepage and are *not* rendered as single pages or with [single page templates][singles].

{{% note %}}
All pages, including List pages, can have front matter and front matter can have markdown content. Thus, `_index.md` files are the way to _provide_ front matter *and* content to the respective list, terms, and homepage templates.
{{% /note %}}


Here are a couple of examples to make it clearer...

```
| \_index.md location                 | Page affected             | Rendered by                   |
| -------------------                 | ------------              | -----------                   |
| /content/post/\_index.md            | site.com/post/            | /layouts/section/post.html    |
| /content/categories/hugo/\_index.md | site.com/categories/hugo/ | /layouts/taxonomy/hugo.html   |
```

### Why `_index.md` Files are Used

With a Single page such as a post it's possible to add the front matter and content directly into the .md page itself. With List/Terms/Homepages this is not possible so \_index.md files can be used to provide that front matter/content to them.

### How to Display Content From `_index.md`

From the information above it should follow that content within an \_index.md file won't be rendered in its own Single Page, instead it'll be made available to the respective list, terms, Homepage.

To **_actually render that content_** you need to ensure that the relevant template responsible for rendering the List/Terms/Homepage contains (at least) `{{ .Content }}`.

This is the way to actually display the content within the \_index.md file on the List/Terms/Homepage.

A very simple example is shown in the following default section list page:

{{% code file="layouts/_default/section.html" download="section.html" %}}
```html
{{ define "main" }}
  <main>
      {{ .Content }}
          <ul class="contents">
          {{ range .Paginator.Pages }}
              <li>{{.Title}}
                  <div>
                    {{ partial "summary.html" . }}
                  </div>
              </li>
          {{ end }}
          </ul>
      {{ partial "pagination.html" . }}
  </main>
{{ end }}
```
{{% /code %}}

You can see `{{ .Content }}` just after the `<main>` element. For this particular example, the content of the \_index.md file will show before the main list of summaries.

### Where to Organize `_index.md` Files

To add content and front matter to the homepage, a section, a taxonomy or a taxonomy terms listing, add a markdown file with the base name \_index on the relevant place on the file system.

```bash
└── content
    ├── _index.md
    ├── categories
    │   ├── _index.md
    │   └── photo
    │       └── _index.md
    ├── post
    │   ├── _index.md
    │   └── firstpost.md
    └── tags
        ├── _index.md
        └── hugo
            └── _index.md
```

In the above example, `_index.md` pages have been added to each section and taxonomy.

An `_index.md` file has also been added in the top level 'content' directory.

### Where to Place `_index.md` for the Homepage Template

Hugo themes are designed to use the 'content' directory as the root of the website, so adding an `_index.md` file here (like has been done in the example above) is how you would add front matter and content to the homepage.

[front matter]: /content-management/front-matter/
[homepage]: /templates/homepage/
[section]: /content-management/section/
[formats]: /content-management/formats/
[singles]: /templates/single-page-templates/
[urls]: /content-management/urls/
