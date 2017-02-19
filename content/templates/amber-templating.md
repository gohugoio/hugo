---
title: Amber Templating
linktitle:
description: In addition to Go templates and Ace templates, Hugo supports the powerful Amber templating from @eknkc.
godocref: https://godoc.org/github.com/eknkc/amber
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
weight:
categories: [templates]
tags: [amber, templating options]
draft: false
aliases: [/templates/amber/]
toc: false
notesforauthors:
---

Hugo also supports the Amber templating engine in addition to [Go templates][] and [Ace templates][] templates.

For template documentation, follow the links from the [Amber project][].

* Amber templates must be named with the Amber suffix; e.g., `list.amber`
* Partials in Amber or HTML can be included with the Amber template syntax:
    * `import ../partials/test.html `
    * `import ../partials/test_a.amber `

[Ace templates]: /templates/ace-templating/
[Amber project]: https://github.com/eknkc/amber
[Go templates]: /templates/go-template-primer/