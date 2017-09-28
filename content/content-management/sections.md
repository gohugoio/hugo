---
title: Content Sections
linktitle: Sections
description: Hugo supports content sections, which according to Hugo's default behavior, will reflect the structure of the rendered website.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [content management]
keywords: [lists,sections,content types,organization]
menu:
  docs:
    parent: "content-management"
    weight: 50
weight: 50	#rem
draft: false
aliases: [/content/sections/]
toc: true
---

{{% note %}}
This section is not updated with the new nested sections support in Hugo 0.24, see https://github.com/gohugoio/hugoDocs/issues/36
{{% /note %}}
{{% todo %}}
See above
{{% /todo %}}

Hugo believes that you organize your content with a purpose. The same structure that works to organize your source content is used to organize the rendered site (see [directory structure][]).

Following this pattern, Hugo uses the top level of your content organization as the content **section**.

The following example shows a content directory structure for a website that has three sections: "authors," "events," and "posts":

```
.
└── content
    ├── authors
    |   ├── _index.md     // <- example.com/authors/
    |   ├── john-doe.md   // <- example.com/authors/john-doe/
    |   └── jane-doe.md   // <- example.com/authors/jane-doe/
    └── events
    |   ├── _index.md     // <- example.com/events/
    |   ├── event-1.md    // <- example.com/events/event-1/
    |   ├── event-2.md    // <- example.com/events/event-2/
    |   └── event-3.md    // <- example.com/events/event-3/
    └── posts
    |   ├── _index.md    // <- example.com/posts/
    |   ├── post-1.md    // <- example.com/posts/post-1/
    |   ├── post-2.md    // <- example.com/posts/post-2/
    |   ├── post-3.md    // <- example.com/posts/post-3/
    |   ├── post-4.md    // <- example.com/posts/post-4/
    |   └── post-5.md    // <- example.com/posts/post-5/
```

## Content Section Lists

Hugo will automatically create pages for each section root that list all of the content in that section. See the documentation on [section templates][] for details on customizing the way these pages are rendered.

As of Hugo v0.18, section pages can also have a content file and front matter. These section content files must be placed in their corresponding section folder and named `_index.md` in order for Hugo to correctly render the front matter and content.

{{% warning "`index.md` vs `_index.md`" %}}
Hugo themes developed before v0.18 often used an `index.md`(i.e., without the leading underscore [`_`]) in a content section as a hack to emulate the behavior of `_index.md`. The hack may work...*sometimes*; however, the order of page rendering can be unpredictable in Hugo. What works now may fail to render appropriately as your site grows. It is **strongly advised** to use `_index.md` as content for your section index pages. **Note:** `_index.md`'s layout, as representative of a section, is a [list page template](/templates/section-templates/) and *not* a [single page template](/templates/single-page-templates/). If you want to alter the new default behavior for `_index.md`, configure `disableKinds` accordingly in your [site's configuration](/getting-started/configuration/).
{{% /warning %}}

## Content *Section* vs Content *Type*

By default, everything created within a section will use the [content type][] that matches the section name. For example, Hugo will assume that `posts/post-1.md` has a `posts` content type. If you are using an [archetype][] for your posts section, Hugo will generate front matter according to what it finds in `archetypes/posts.md`.

[archetype]: /content-management/archetypes/
[content type]: /content-management/types/
[directory structure]: /getting-started/directory-structure/
[section templates]: /templates/section-templates/


