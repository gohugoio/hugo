---
aliases:
- /doc/templates/amber
- /layout/templates/amber
- /layout/amber/
date: 2015-07-20
linktitle: Amber templates
menu:
  main:
    parent: layout
next: /templates/functions
prev: /templates/go-templates
title: Amber Templates
weight: 18
---

Amber templates are another template type which Hugo supports, in addition to [Go templates](/templates/go-templates) and [Ace templates](/templates/ace) templates.

For template documentation, follow the links from the [Amber project](https://github.com/eknkc/amber)

* Amber templates must be named with the amber-suffix, e.g. `list.amber`
* Partials in Amber or HTML can be included with the Amber template syntax:
	* `import ../partials/test.html `
	* `import ../partials/test_a.amber `


