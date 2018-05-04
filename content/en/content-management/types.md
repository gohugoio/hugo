---
title: Content Types
linktitle: Types
description: Hugo supports sites with multiple content types and assumes your site will be organized into sections, where each section represents the corresponding type.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [content management]
keywords: [lists,sections,content types,types,organization]
menu:
  docs:
    parent: "content-management"
    weight: 60
weight: 60	#rem
draft: false
aliases: [/content/types]
toc: true
---

A **content type** can have a unique set of metadata (i.e., [front matter][]) or customized [template][] and can be created by the `hugo new` command via [archetypes][].

## What is a Content Type

[Tumblr][] is a good example of a website with multiple content types. A piece of "content" could be a photo, quote, or a post, each with different sets of metadata and different visual rendering.

## Assign a Content Type

Hugo assumes that your site will be organized into [sections][] and each section represents a corresponding type. This is to reduce the amount of configuration necessary for new Hugo projects.

If you are taking advantage of this default behavior, each new piece of content you place into a section will automatically inherit the type. Therefore a new file created at `content/posts/new-post.md` will automatically be assigned the type `posts`. Alternatively, you can set the content type in a content file's [front matter][] in the field "`type`".

## Create New Content of a Specific Type

You can manually add files to your content directories, but Hugo can create and populate a new content file with preconfigured front matter via [archetypes][].

## Define a Content Type

Creating a new content type is easy. You simply define the templates and archetype unique to your new content type, or Hugo will use defaults.


{{% note "Declaring Content Types" %}}
Remember, all of the following are *optional*. If you do not specifically declare content types in your front matter or develop specific layouts for content types, Hugo is smart enough to assume the content type from the file path and section. (See [Content Sections](/content-management/sections/) for more information.)
{{% /note %}}

The following examples take you stepwise through creating a new type layout for a content file that contains the following front matter:

{{< code file="content/events/my-first-event.md" copy="false" >}}
+++
title = My First Event
date = "2016-06-24T19:20:04-07:00"
description = "Today is my 36th birthday. How time flies."
type = "event"
layout = "birthday"
+++
{{< /code >}}

By default, Hugo assumes `*.md` under `events` is of the `events` content type. However, we have specified that this particular file at `content/events/my-first-event.md` is of type `event` and should render using the `birthday` layout.

### Create a Type Layout Directory

Create a directory with the name of the type in `/layouts`. For creating these custom layouts, **type is always singular**; e.g., `events => event` and `posts => post`.

For this example, you need to create `layouts/event/birthday.html`.

{{% note %}}
If you have multiple content files in your `events` directory that are of the `special` type and you don't want to define the `layout` specifically for each piece of content, you can create a layout at `layouts/special/single.html` to observe the [single page template lookup order](/templates/single-page-templates/).
{{% /note %}}

{{% warning %}}
With the "everything is a page" data model introduced in v0.18 (see [Content Organization](/content-management/organization/)), you can use `_index.md` in content directories to add both content and front matter to [list pages](/templates/lists/). However, `type` and `layout` declared in the front matter of `_index.md` are *not* currently respected at build time as of v0.19. This is a known issue [(#3005)](https://github.com/gohugoio/hugo/issues/3005).
{{% /warning %}}

### Create Views

Many sites support rendering content in a few different ways; e.g., a single page view and a summary view to be used when displaying a [list of section contents][sectiontemplates].

Hugo limits assumptions about how you want to display your content to an intuitive set of sane defaults and will support as many different views of a content type as your site requires. All that is required for these additional views is that a template exists in each `/layouts/<TYPE>` directory with the same name.

### Custom Content Type Template Lookup Order

The lookup order for the `content/events/my-first-event.md` templates would be as follows:

* `layouts/event/birthday.html`
* `layouts/event/single.html`
* `layouts/events/single.html`
* `layouts/_default/single.html`

### Create a Corresponding Archetype

We can then create a custom archetype with preconfigured front matter at `event.md` in the `/archetypes` directory; i.e. `archetypes/event.md`.

Read [Archetypes][archetypes] for more information on archetype usage with `hugo new`.

[archetypes]: /content-management/archetypes/
[front matter]: /content-management/front-matter/
[sectiontemplates]: /templates/section-templates/
[sections]: /content-management/sections/
[template]: /templates/
[Tumblr]: https://www.tumblr.com/
