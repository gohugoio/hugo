---
title: Static Files
description: "Files that get served **statically** (as-is, no modification) on the site root."
date: 2017-11-18
categories: [content management]
keywords: [source, directories]
menu:
  docs:
    parent: "content-management"
    weight: 130
weight: 130	#rem
aliases: [/static-files]
toc: true
---

By default, the `static/` directory in the site project is used for
all **static files** (e.g. stylesheets, JavaScript, images). The static files are served on the site root path (eg. if you have the file `static/image.png` you can access it using `http://{server-url}/image.png`, to include it in a document you can use `![Example image](/image.png) )`.

Hugo can be configured to look into a different directory, or even
**multiple directories** for such static files by configuring the
`staticDir` parameter in the [site config][]. All the files in all the
static directories will form a union filesystem.

This union filesystem will be served from your site root. So a file
`<SITE PROJECT>/static/me.png` will be accessible as
`<MY_BASEURL>/me.png`.

Here's an example of setting `staticDir` and `staticDir2` for a
multi-language site:

{{< code-toggle copy="false" file="config" >}}
staticDir = ["static1", "static2"]

[languages]
[languages.en]
staticDir2 = "static_en"
baseURL = "https://example.com"
languageName = "English"
weight = 2
title = "In English"
[languages.no]
staticDir = ["staticDir_override", "static_no"]
baseURL = "https://example.no"
languageName = "Norsk"
weight = 1
title = "På norsk"
{{</ code-toggle >}}

In the above, with no theme used:

- The English site will get its static files as a union of "static1",
  "static2" and "static_en". On file duplicates, the right-most
  version will win.
- The Norwegian site will get its static files as a union of
  "staticDir_override" and "static_no".

Note 1
: The **2** (can be a number between 0 and 10) in `staticDir2` is
  added to tell Hugo that you want to **add** this directory to the
  global set of static directories defined using `staticDir`. Using
  `staticDir` on the language level would replace the global value (as
  can be seen in the Norwegian site case).

Note 2
: The example above is a [multihost setup][]. In a regular setup, all
  the static directories will be available to all sites.


[site config]: /getting-started/configuration/#all-configuration-settings
[multihost setup]: /content-management/multilingual/#configure-multilingual-multihost
