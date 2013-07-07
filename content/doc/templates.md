---
title: "Templates"
Pubdate: "2013-07-01"
---

Hugo uses the excellent golang html/template library for it's template engine. It is an extremely
lightweight engine that provides a very small amount of logic. In our 
experience that it is just the right amount of logic to be able to create a good static website

This document will not cover how to use golang templates, but the [golang docs](http://golang.org/pkg/html/template/)
provide a good introduction.

### Template roles

There are 5 different kinds of templates that Hugo works with.

#### index.html
This file must exist in the layouts directory. It is the template used to render the 
homepage of your site.

#### rss.xml
This file must exist in the layouts directory. It will be used to render all rss documents.
The one provided in the example application will generate an ATOM format. 

*Important: Hugo will automatically add the following header line to this file.*

    <?xml version="1.0" encoding="utf-8" standalone="yes" ?>

#### Indexes
An index is a page that list multiple pieces of content. If you think of a typical blog, the tag 
pages are good examples of indexes.


#### Content Type(s)
Hugo supports multiple types of content. Another way of looking at this is that Hugo has the ability
to render content in a variety of ways as determined by the type.

#### Chrome
Chrome is simply the decoration of your site. It's not a requirement to have this, but in practice
it's very convenient. Hugo doesn't know anything about Chrome, it's simply a convention that you may
likely find beneficial. As you create the rest of your templates you will include templates from the 
/layout/chrome directory. I've found it helpful to include a header and footer template 
in Chrome so I can include those in the other full page layouts (index.html, indexes/ type/single.html).

### Adding a new content type

Adding a type is easy.

**Step 1:**
Create a directory with the name of the type in layouts.Type is always singular.  *Eg /layouts/post*.

**Step 2:**
Create a file called single.html inside your directory. *Eg /layouts/post/single.html*.

**Step 3:**
Create a file with the same name as your directory in /layouts/indexes/. *Eg /layouts/index/post.html*.

**Step 4:**
Many sites support rendering content in a few different ways, for instance a single page view and a 
summary view to be used when displaying a list of contents on a single page. Hugo makes no assumptions
here about how you want to display your content, and will support as many different views of a content
type as your site requires. All that is required for these additional views is that a template
exists in each layout/type directory with the same name.

For these, reviewing this example site will be very helpful in order to understand how these types work.

