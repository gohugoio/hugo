---
title: Content Organization
linktitle: Organization
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
wip: true
---

## Organization of Content Source

In Hugo, your content should be organized in a manner that reflects the rendered website.

While Hugo supports content nested at any level, the top levels (i.e. `content/<DIRECTORIES>`) are special in Hugo and are considered the content [sections][]. Without any additional configuration, the following will just work:

```
.
└── content
    └── about
    |   └── _index.md  // <- http://yoursite.com/about/
    ├── post
    |   ├── firstpost.md   // <- http://yoursite.com/post/firstpost/
    |   ├── happy
    |   |   └── ness.md  // <- http://yoursite.com/post/happy/ness/
    |   └── secondpost.md  // <- http://yoursite.com/post/secondpost/
    └── quote
        ├── first.md       // <- http://yoursite.com/quote/first/
        └── second.md      // <- http://yoursite.com/quote/second/
```

## Path Breakdown in Hugo

The following demonstrates the relationships between your content organization and the output URL structure for your Hugo website at render. These examples assume you are using pretty URLs, which is the default behavior for Hugo. The examples also assume a key-value of `baseurl = "http://yoursite.com"` in your [site's configuration file][config].

### Section Index Page

`_index.md` has a special role in Hugo. It allows you to add front matter and content to your [list templates][lists] as of v0.18. These templates include those for [section templates][], [taxonomy templates][], [taxonomy terms templates][], and your [homepage template][].

You can keep one `_index.md` in each of your content sections. The following shows typical placement of an `_index.md` that would contain content and front matter for a `posts` section list page on a Hugo website:


```bash
.            url
.       ⊢------^------⊣
.        path    slug
.       ⊢--^-⊣⊢---^---⊣
.           filepath
.       ⊢------^------⊣
content/posts/_index.md
```

At build, this will output to the following destination with the associated values:

```bash

                     url ("/posts/")
                    ⊢-^-⊣
       baseurl      section ("posts")
⊢--------^---------⊣⊢-^-⊣
        permalink
⊢----------^-------------⊣
http://yoursite.com/posts/index.html
```

### Section Single Pages

Single content files in each of your sections are going to be rendered as [single page templates][singles]. Here is an example of a single `post` within `posts`:


```bash
                   path ("posts/my-first-hugo-post.md")
.       ⊢-----------^------------⊣
.      section        slug
.       ⊢-^-⊣⊢--------^----------⊣
content/posts/my-first-hugo-post.md
```

At the time Hugo renders your site, the content will be output to the following destination:

```bash

                               url ("/posts/my-first-hugo-post/")
                   ⊢------------^----------⊣
       baseurl     section     slug
⊢--------^--------⊣⊢-^--⊣⊢-------^---------⊣
                 permalink
⊢--------------------^---------------------⊣
http://yoursite.com/posts/my-first-hugo-post/index.html
```

### Section with Nested Directories

To continue the example, the following demonstrates destination paths for a file located at `content/events/chicago/lollapalooza.md` in the same site:


```bash
                    section
                    ⊢--^--⊣
                               url
                    ⊢-------------^------------⊣

      baseURL             path        slug
⊢--------^--------⊣ ⊢------^-----⊣⊢----^------⊣
                  permalink
⊢----------------------^-----------------------⊣
http://yoursite.com/events/chicago/lollapalooza/
```
## Path Properties Explained

#### `section`

A default content type is determined by a piece of content's section. `section` is determined by the location within the project's `content` directory. `section` *cannot* be specified or overridden in front matter.

#### `slug`

A content's `slug` is either `name.extension` or `name/`. The value for `slug` is determined by

* the name of the content file (e.g., `lollapalooza.md`) OR
* front matter overrides

#### `path`

A content's `path` is determined by the section's path to the file. The file `path`

* is based on the path to the content's location AND
* does not include the slug

#### `url`

The `url` is the relative URL for the piece of content. The `url`

* is based on the content's location within the directory structure OR
* is defined in front matter and *overrides all the above*

## Modifying Destinations for Content Source in Front Matter

Hugo believes that you organize your content with a purpose. The same structure that works to organize your source content is used to organize the rendered site. As displayed above, the organization of the source content will be mirrored in the destination.

Notice that the first level `about/` page URL was created using a directory named "about" with a single `_index.md` file inside.

There are times where you may need more control over your content. In these cases, there are fields that can be specified in the front matter to determine the destination of a specific piece of content.

The following items are defined in this order for a specific reason: latter items in the list will override earlier items, and not all of these items can be defined in front matter:

### `filename`

This isn't in the front matter, but is the actual name of the file minus the extension. This will be the name of the file in the destination (e.g., `content/posts/my-post.md` becomes `yoursite.com/posts/my-post/`).

### `slug`

When defined in the front matter, the `slug` can take the place of the filename for the destination.

{{% code file="content/posts/old-post.md" %}}
```yaml
---
title: New Post
slug: "new-post"
---
```
{{% /code %}}

This will render to the following destination:

```
yoursite.com/posts/new-post/
```

### `section`

`section` is determined by a content's location on disk and *cannot* be specified in the front matter. See [sections][] for more information.

### `type`

A content's `type` is also determined by its location on disk but, unlike `section`, it *can* be specified in the front matter. See [types][].

{{% code file="content/posts/my-post.md" %}}
```yaml
---
title: My Post
type: blog
---
```
{{% /code %}}

### `path`

`path` can be provided in the front matter. This will replace the actual path to the file on disk. Destination will create the destination with the same path, including the section.

### `url`

A complete URL can be provided. This will override all the above as it pertains to the end destination. This must be the path from the baseURL (starting with a `/``). When `url` is provided in the front matter, it will be used exactly. Using `url` will ignore the `--uglyURLs` setting.

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

[config]: /getting-started/configuration/
[formats]: /content-management/formats/
[front matter]: /content-management/front-matter/
[homepage template]: /templates/homepage/
[homepage]: /templates/homepage/
[lists]: /templates/lists/
[section templates]: /templates/section-templates/
[sections]: /content-management/sections/
[singles]: /templates/single-page-templates/
[taxonomy templates]: /templates/taxonomy-templates/
[taxonomy terms templates]: /templates/taxonomy-templates/
[types]: /content-management/types/
[urls]: /content-management/urls/
