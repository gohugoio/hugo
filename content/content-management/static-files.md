---
title: Static Files
description: "The `static` folder is where you place all your **static files**."
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

The `static` folder is where you place all your **static files**, e.g. stylesheets, JavaScript, images etc.

You can set the name of the static folder to use in your configuration file, for example `config.toml`.  From **Hugo 0.31** you can configure as many static directories as you need. All the files in all the static directories will form a union filesystem.

Example:

```toml
staticDir = ["static1", "static2"]
[languages]
[languages.no]
staticDir = ["staticDir_override", "static_no"]
baseURL = "https://example.no"
languageName = "Norsk"
weight = 1
title = "På norsk"

[languages.en]
staticDir2 = "static_en"
baseURL = "https://example.com"
languageName = "English"
weight = 2
title = "In English"
```

In the above, with no theme used:

* The English site will get its static files as a union of "static1", "static2" and "static_en". On file duplicates, the right-most version will win.
* The Norwegian site will get its static files as a union of "staticDir_override" and "static_no".

**Note:** The `2` `static2` (can be a number between 0 and 10) is added to tell Hugo that you want to **add** this directory to the global set of static directories. Using `staticDir` on the language level would replace the global value.


**Note:** The example above is a [multihost setup](/content-management/multilingual/#configure-multilingual-multihost). In a regular setup, all the static directories will be available to all sites.
