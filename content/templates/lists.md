---
title: Introduction to Lists in Hugo
linktitle: Hugo Lists Introduction
description: Lists have a specific meaning and usage in Hugo. If you want to know how to render your site homepage, section homepage, taxonomy list, or taxonomy terms list, start here.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [lists,sections,rss,taxonomies,terms]
weight: 22
draft: false
aliases: [/templates/lists-in-hugo/]
toc: true
needsreview: true
---

## What is a List Page Template?

A list page template is a template used to render multiple pieces of content in a single HTML page (with the exception of the homepage, which has a [dedicated template][homepage]).

Hugo uses the term *list* in its truest sense: a sequential arrangement of material, especially in alphabetical or numerical order. Hugo uses list templates on any output HTML page where content is being listed (e.g., [taxonomies][], [sections][], and [RSS][]). The idea of a list page comes from the [hierarchical mental model of the web][mentalmodel] and is best demonstrated visually:

```yaml
---
aliases:
- /doc/using-index-md/
lastmod: 2017-02-22
date: 2017-02-22
linktitle: Using _index.md
menu:
  main:
    parent: content
prev: /content/example
next: /themes/overview
notoc: true
title: Using _index.md
weight: 70
---
```

## \_index.md and 'Everything is a Page'

As of version v0.18 Hugo now treats '[everything as a page](http://bepsays.com/en/2016/12/19/hugo-018/)'. This allows you to add content and front matter to any page - including List pages like [Sections](/content/sections/), [Taxonomies](/taxonomies/overview/), [Taxonomy Terms pages](/templates/terms/) and even to potential 'special case' pages like the [Home page](/templates/homepage/).

In order to take advantage of this behaviour you need to do a few things.

1. Create an \_index.md file that contains the front matter and content you would like to apply.

2. Place the \_index.md file in the correct place in the directory structure.

3. Ensure that the respective template is configured to display `{{ .Content }}` if you wish for the content of the \_index.md file to be rendered on the respective page.

### How \_index.md pages work

Before continuing it's important to know that this page must reference certain templates to describe how the \_index.md page will be rendered. Hugo has a multitude of possible templates that can be used and placed in various places (think theme templates for instance). For simplicity/brevity the default/top level template location will be used to refer to the entire range of places the template can be placed.

If this is confusing or you are unfamiliar with Hugo's template hierarchy, visit the various template pages listed below. You may need to find the 'active' template responsible for any particular page on your own site by going through the template hierarchy and matching it to your particular setup/theme you are using.

- [Home page template](/templates/homepage/)
- [Content List templates](/templates/list/)
- [Single Content templates](/templates/content/)
- [Taxonomy Terms templates](/templates/terms/)

Now that you've got a handle on templates lets recap some Hugo basics to understand how to use an \_index.md file with a List page.

1. Sections and Taxonomies are 'List' pages, NOT single pages.
2. List pages are rendered using the template heirarchy found in the [Content - List Template](http://localhost:1313/templates/list/) docs.
3. The Home page, though technically a List page, can have [it's own template](/templates/homepage/) at layouts/index.html rather than \_default/list.html. Many themes exploit this behaviour so you are likely to encounter this specific use case.
4. Taxonomy terms pages are 'lists of metadata' not lists of content, so [have their own templates](/templates/terms/).

Let's put all this information together:

> **\_index.md files used in List pages, Terms pages or the Home page are NOT rendered as single pages or with Single Content templates.**

> **All pages, including List pages, can have front matter and front matter can have markdown content - meaning \_index.md files are the way to _provide_ front matter and content to the respective List/Terms/Home page.**

Here are a couple of examples to make it clearer...

```
| \_index.md location                 | Page affected             | Rendered by                   |
| -------------------                 | ------------              | -----------                   |
| /content/post/\_index.md            | site.com/post/            | /layouts/section/post.html    |
| /content/categories/hugo/\_index.md | site.com/categories/hugo/ | /layouts/taxonomy/hugo.html   |
```

## Why \_index.md files are used

With a Single page such as a post it's possible to add the front matter and content directly into the .md page itself. With List/Terms/Home pages this is not possible so \_index.md files can be used to provide that front matter/content to them.

## How to display content from \_index.md files

From the information above it should follow that content within an \_index.md file won't be rendered in its own Single Page, instead it'll be made available to the respective list, terms, Home page.

To **_actually render that content_** you need to ensure that the relevant template responsible for rendering the List/Terms/Home page contains (at least) `{{ .Content }}`.

This is the way to actually display the content within the \_index.md file on the List/Terms/Home page.

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

## Where to Organize `\_index.md` Files

To add content and front matter to the home page, a section, a taxonomy or a taxonomy terms listing, add a markdown file with the base name \_index on the relevant place on the file system.

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

In the above example \_index.md pages have been added to each section/taxonomy.

An `_index.md` file has also been added in the top level 'content' directory.

### Where to place \_index.md for the Homepage Template

Hugo themes are designed to use the 'content' directory as the root of the website, so adding an \_index.md file here (like has been done in the example above) is how you would add front matter/content to the home page.

## List Defaults

### Default Templates

Since section lists and taxonomy lists (N.B., *not* [taxonomy terms lists][]) are both *lists* with regards to their templates, both of these templates have the same terminating default of `_default/list.html`---or `themes/<MYTHEME>/layouts/_default/list.html` in the case of a themed project---in their *lookup orders*. In addition, both [section lists][sections] and [taxonomy lists][taxonomies] have their own default list templates in `_default`:

#### Default Section Templates

1. `layouts/section/<SECTIONNAME>.html`
2. `layouts/section/list.html`
3. `layouts/_default/section.html`
4. `layouts/_default/list.html`

### Understanding `.Data.Pages`


### Taxonomy Lists

A Taxonomy will be rendered at /`PLURAL`/`TERM`/ (e.g.&nbsp;http://spf13.com/topics/golang/) from:

* /layouts/taxonomy/`SINGULAR`.html (e.g.&nbsp;`/layouts/taxonomy/topic.html`)
* /layouts/\_default/taxonomy.html
* /layouts/\_default/list.html
* /themes/`THEME`/layouts/taxonomy/`SINGULAR`.html
* /themes/`THEME`/layouts/\_default/taxonomy.html
* /themes/`THEME`/layouts/\_default/list.html

Note that a taxonomy list page can also have a content file with front matter,  see [Source Organization](/overview/source-directory/).

### Taxonomy RSS

A Taxonomy’s RSS will be rendered at /`PLURAL`/`TERM`/index.xml (e.g.&nbsp;http://spf13.com/topics/golang/index.xml)

*Hugo ships with its own [RSS 2.0][] template. In most cases this will
be sufficient, and an RSS template will not need to be provided by the
user.*

Hugo provides the ability for you to define any RSS type you wish, and
can have different RSS files for each section and taxonomy.

* /layouts/taxonomy/`SINGULAR`.rss.xml
* /layouts/\_default/rss.xml
* /themes/`THEME`/layouts/taxonomy/`SINGULAR`.rss.xml
* /themes/`THEME`/layouts/\_default/rss.xml

## Variables

A list page is a `Page` and have all the [page variables](/templates/variables/)
and [site variables](/templates/variables/) available to use in the templates.

Taxonomy pages will additionally have:

**.Data.`Singular`** The taxonomy itself.<br>



{{% note %}}
If `where` or `first` receives invalid input or a field name that doesn’t exist, it will return an error and stop site generation. `where` and `first` also work on taxonomy list templates *and* taxonomy terms templates. (See [Taxonomy Templates](/templates/taxonomy-templates/).)
{{% /note %}}


[directorystructure]: /getting-started/directory-structure/
[homepage]: /templates/homepage-template/
[homepage]: /templates/homepage-template/
[limitkeyword]: https://www.techonthenet.com/sql/select_limit.php
[mentalmodel]: http://webstyleguide.com/wsg3/3-information-architecture/3-site-structure.html
[partials]: /templates/partial-templates/
[RSS 2.0]: http://cyber.law.harvard.edu/rss/rss.html "RSS 2.0 Specification"
[RSS]: /templates/rss-templates/
[sections]: /content-management/sections/
[sections]: /templates/section-templates
[singlepages]: /templates/single-page-templates/
[taxonomies]: /templates/taxonomy-templates/#taxonomy-list-templates/
[taxonomy terms lists]: /templates/#taxonomy-terms-templates/
