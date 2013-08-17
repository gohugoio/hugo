---
title: "Sections"
date: "2013-07-01"
---

Hugo thinks that you organize your content with a purpose. The same structure
that works to organize your source content is used to organize the rendered
site ( [see organization](/content/organization) ). Following this pattern Hugo
uses the top level of your content organization as **the Section**.

The following example site uses two sections, "post" and "quote".

    .
    └── content
        ├── post
        |   ├── firstpost.md       // <- http://site.com/post/firstpost/
        |   ├── happy
        |   |   └── happiness.md   // <- http://site.com/happy/happiness/
        |   └── secondpost.md      // <- http://site.com/post/secondpost/
        └── quote
            ├── first.md           // <- http://site.com/quote/first/
            └── second.md          // <- http://site.com/quote/second/


*Regardless of location on disk, the section can be provided in the front matter
which will affect the destination location*.

## Sections and Types

By default everything created within a section will use the content type
that matches the section name.

Section defined in the front matter have the same impact.

To change the type of a given piece of content simply define the type
in the front matter.

If a layout for a given type hasn't been provided a default type template will
be used instead provided is exists.




