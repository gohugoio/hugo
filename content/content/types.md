---
title: "Content Types"
date: "2013-07-01"
linktitle: "Types"
menu:
    main:
        parent: 'content'
weight: 30
---

Hugo has full support for multiple content types each with its own set
of meta data and template. A good example of when multiple types are
needed is to look at Tumblr. A piece of content could be a photo, quote
or post, each with different meta data and rendered differently.

## Defining a content type

Creating a new content type is easy in Hugo. You simply provide the
templates that the new type will use.

It is essential to provide the single render view template as well as a
list view template.

### Step 1: Create Type Directory
Create a directory with the name of the type in layouts.Type is always singular.  *Eg /layouts/post*.

### Step 2: Create template
Create a file called single.html inside your directory. *Eg /layouts/post/single.html*.

### Step 3: Create list template
Create a file with the same name as your directory in /layouts/indexes/. *Eg /layouts/indexes/post.html*.

### Step 4: Create views
Many sites support rendering content in a few different ways, for
instance a single page view and a summary view to be used when displaying a list
of contents on a single page. Hugo makes no assumptions here about how you want
to display your content, and will support as many different views of a content
type as your site requires. All that is required for these additional views is
that a template exists in each layout/type directory with the same name.

For these, reviewing this example site will be very helpful in order to understand how these types work.

## Assigning a content type

Hugo assumes that your site will be organized into [sections](/content/sections)
and each section will use the corresponding type. If you are taking advantage of
this then each new piece of content you place into a section will automatically
inherit the type.

Alternatively you can set the type in the meta data under the key "type".
