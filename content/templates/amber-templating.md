---
title: Amber Templating
linktitle:
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
weight:
tags: [amber,layout]
categories: [templates]
draft: false
slug: [/templates/amber/]
aliases:
toc: false
notes:
---

Amber templates are another template type which Hugo supports, in addition to [Go templates][] and [Ace templates][] templates.

For template documentation, follow the links from the [Amber project][].

* Amber templates must be named with the amber-suffix, e.g. `list.amber`
* Partials in Amber or HTML can be included with the Amber template syntax:
    * `import ../partials/test.html `
    * `import ../partials/test_a.amber `

[Ace templates]: /templates/ace-templating/
[Amber project]: https://github.com/eknkc/amber
[Go templates]: /templates/go-template-primer/

