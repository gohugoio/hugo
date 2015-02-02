---
aliases:
- /doc/organization/
date: 2013-07-01
linktitle: Organization
menu:
  main:
    parent: content
next: /content/front-matter
prev: /overview/source-directory
title: Content Organization
weight: 10
---

Hugo uses Markdown files with headers commonly called the *front matter*. Hugo
respects the organization that you provide for your content to minimize any
extra configuration, though this can be overridden by additional configuration
in the front matter.

## Organization

In Hugo, the content should be arranged in the same way they are intended for
the rendered website. Without any additional configuration, the following will
just work. Hugo supports content nested at any level. The top level is special
in Hugo and is used as the [section](/content/sections/).

    .
    └── content
        ├── post
        |   ├── firstpost.md   // <- http://1.com/post/firstpost/
        |   ├── happy
        |   |   └── ness.md    // <- http://1.com/post/happy/ness/
        |   └── secondpost.md  // <- http://1.com/post/secondpost/
        └── quote
            ├── first.md       // <- http://1.com/quote/first/
            └── second.md      // <- http://1.com/quote/second/

**Here's the same organization run with `hugo --uglyUrls`**

    .
    └── content
        ├── post
        |   ├── firstpost.md   // <- http://1.com/post/firstpost.html
        |   ├── happy
        |   |   └── ness.md    // <- http://1.com/post/happy/ness.html
        |   └── secondpost.md  // <- http://1.com/post/secondpost.html
        └── quote
            ├── first.md       // <- http://1.com/quote/first.html
            └── second.md      // <- http://1.com/quote/second.html

## Destinations

Hugo believes that you organize your content with a purpose. The same structure
that works to organize your source content is used to organize the rendered
site. As displayed above, the organization of the source content will be
mirrored in the destination.

There are times when one would need more control over their content. In these
cases, there are a variety of things that can be specified in the front matter
to determine the destination of a specific piece of content.

The following items are defined in order, latter items in the list will override
earlier settings.

### filename
This isn't in the front matter, but is the actual name of the file minus the
extension. This will be the name of the file in the destination.

### slug
Defined in the front matter, the slug can take the place of the filename for the
destination.

### filepath
The actual path to the file on disk. Destination will create the destination
with the same path. Includes [section](/content/sections/).

### section
section can be provided in the front matter overriding the section derived from
the source content location on disk. See [section](/content/sections/).

### path
path can be provided in the front matter. This will replace the actual
path to the file on disk. Destination will create the destination with the same
path. Includes [section](/content/sections/).

### url
A complete URL can be provided. This will override all the above as it pertains
to the end destination. This must be the path from the baseurl (starting with a "/").
When a url is provided, it will be used exactly. Using url will ignore the
`--uglyUrls` setting.


## Path breakdown in Hugo

### Content

    .             path           slug
    .       ⊢-------^----⊣ ⊢------^-------⊣
    content/extras/indexes/category-example/index.html


    .       section              slug
    .       ⊢--^--⊣        ⊢------^-------⊣
    content/extras/indexes/category-example/index.html


    .       section  slug
    .       ⊢--^--⊣⊢--^--⊣
    content/extras/indexes/index.html

### Destination


               permalink
    ⊢--------------^-------------⊣
    http://spf13.com/projects/hugo


       baseUrl       section  slug
    ⊢-----^--------⊣ ⊢--^---⊣ ⊢-^⊣
    http://spf13.com/projects/hugo


       baseUrl       section          slug
    ⊢-----^--------⊣ ⊢--^--⊣        ⊢--^--⊣
    http://spf13.com/extras/indexes/example


       baseUrl            path       slug
    ⊢-----^--------⊣ ⊢------^-----⊣ ⊢--^--⊣
    http://spf13.com/extras/indexes/example


       baseUrl            url
    ⊢-----^--------⊣ ⊢-----^-----⊣
    http://spf13.com/projects/hugo


       baseUrl               url
    ⊢-----^--------⊣ ⊢--------^-----------⊣
    http://spf13.com/extras/indexes/example



**section** = which type the content is by default

* based on content location 
* front matter overrides

**slug** = name.ext or name/

* based on content-name.md 
* front matter overrides

**path** = section + path to file excluding slug

* based on path to content location


**url** = relative URL

* defined in front matter
* overrides all the above

