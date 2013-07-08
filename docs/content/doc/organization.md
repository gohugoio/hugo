---
title: "Organization"
Pubdate: "2013-07-01"
---

Hugo uses markdown files with headers commonly called the front matter. Hugo respects the organization
that you provide for your content to minimize any extra configuration, though this can be overridden
by additional configuration in the front matter.

## Organization
In Hugo the content should be arranged in the same way they are intended for the rendered website.
Without any additional configuration the following will just work.

    .
    └── content
        ├── post
        |   ├── firstpost.md   // <- http://site.com/post/firstpost.html
        |   └── secondpost.md  // <- http://site.com/post/secondpost.html
        └── quote
            ├── first.md       // <- http://site.com/quote/first.html
            └── second.md      // <- http://site.com/quote/second.html

