---
date: 2013-07-01
linktitle: Types
menu:
  main:
    parent: content
next: /content/archetypes
prev: /content/sections
title: Content Types
weight: 40
---

Hugo has full support for different types of content. A content type can have a
unique set of meta data, template and can be automatically created by the new
command through using content [archetypes](/content/archetypes/).

A good example of when multiple types are needed is to look at [Tumblr](https://www.tumblr.com/). A piece
of content could be a photo, quote or post, each with different meta data and
rendered differently.

## Assigning a content type

Hugo assumes that your site will be organized into [sections](/content/sections/)
and each section will use the corresponding type. If you are taking advantage of
this, then each new piece of content you place into a section will automatically
inherit the type.

Alternatively, you can set the type in the meta data under the key "`type`".


## Creating new content of a specific type

Hugo has the ability to create a new content file and populate the front matter
with the data set corresponding to that type. Hugo does this by utilizing
[archetypes](/content/archetypes/).

To create a new piece of content, use:

    hugo new relative/path/to/content.md

For example, if I wanted to create a new post inside the post section, I would type:

    hugo new post/my-newest-post.md


## Defining a content type

Creating a new content type is easy in Hugo. You simply provide the templates and archetype
that the new type will use. You only need to define the templates, archetypes and/or views
unique to that content type. Hugo will fall back to using the general templates and default archetype
whenever a specific file is not present.

*Remember, all of the following are optional:*

### Create Type Directory
Create a directory with the name of the type in `layouts`. Type is always singular.  *E.g. `/layouts/post`*.

### Create single template
Create a file called `single.html` inside your directory. *E.g. `/layouts/post/single.html`*.

### Create list template
Create a file called `list.html` inside your directory. *E.g. `/layouts/post/list.html`*.

### Create views
Many sites support rendering content in a few different ways, for instance,
a single page view and a summary view to be used when displaying a list
of contents on a single page. Hugo makes no assumptions here about how you want
to display your content, and will support as many different views of a content
type as your site requires. All that is required for these additional views is
that a template exists in each layouts/`TYPE` directory with the same name.

### Create a corresponding archetype

Create a file called <code><em>type</em>.md</code> in the `/archetypes` directory. *E.g. `/archetypes/post.md`*.

More details about archetypes can be found at the [archetypes docs](/content/archetypes/).
