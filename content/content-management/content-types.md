---
title: Content Types
linktitle: Content Types
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [content management]
tags: [lists,sections,content types,types,organization]
weight: 60
draft: false
aliases: [/content/types]
toc: true
needsreview: true
---

Hugo provides full support for sites with multiple content types. A **content type** can have a unique set of metadata (i.e., [front matter][]) or customized [template][] and can be created by the `hugo new` command via content [archetypes][].

A good example of when multiple types are needed is to look at [Tumblr][]. A piece of "content" could be a photo, quote or post, each with different meta data and rendered differently.

## Assigning a content type

Hugo assumes that your site will be organized into [sections][] and each section will use the corresponding type. If you are taking advantage of this, then each new piece of content you place into a section will automatically inherit the type.

Alternatively, you can set the content type in a content file's metadata (i.e., [front matter][]) under the key "`type`".

## Creating New Content of a Specific Type

You can manually add files to your content directories, but Hugo has the ability to create and populate a new content file with preconfigured front matter via [archetypes][].

## Defining a Content Type

Creating a new content type is easy in Hugo. You simply provide the templates and archetype that the new type will use. You only need to define the templates, archetypes and/or views unique to that content type. Hugo will fall back to using the general templates and default archetype whenever a specific file is not present.

{{% note "Declaring Content Types" %}}
Remember, all of the following are *optional*. If you do not specifically declare content types in your front matter or develop specific layouts for content types, Hugo is smart enough to infer the content type from the file path and content section (see [content sections](/content-management/sections/)).
{{% /note %}}

### Create Type Layout Directory

Create a directory with the name of the type in `/layouts`. Type is always singular; e.g., even if your content directory is `content/posts`, create `/layouts/post/`.

### Create a Single Template

Create a file called `single.html` inside your directory. *E.g. `/layouts/post/single.html`*.

### Create a List Template

Create a file called `post.html` inside the section lists template directory, `/layouts/section`; e.g., `/layouts/section/post.html`.

### Create views

Many sites support rendering content in a few different ways, for instance, a single page view and a summary view to be used when displaying a [list of section contents][]. Hugo makes no assumptions about how you want to display your content and will support as many different views of a content type as your site requires. All that is required for these additional views is that a template exists in each `/layouts/TYPE` directory with the same name.

### Create A Corresponding Archetype

Create an archetype file for your type at `yourtype.md` in the `/archetypes` directory; e.g., `/archetypes/post.md`.

More details about archetypes can be found in the [archetypes documentation][].

[archetypes]: /content-management/archetypes/
[archetypes documentation]: /content-management/archetypes/
[sections]: /content-management/sections/
[front matter]: /content-management/front-matter/
[list of section contents]: /templates/section-templates/
[template]: /templates/
[Tumblr]: https://www.tumblr.com/