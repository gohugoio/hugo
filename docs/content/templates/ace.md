---
aliases:
- /doc/templates/ace
- /layout/templates/ace
- /layout/ace/
lastmod: 2015-08-04
date: 2014-04-20
linktitle: Ace templates
menu:
  main:
    parent: layout
next: /templates/functions
prev: /templates/go-templates
title: Ace Templates
weight: 17
---

In addition to [Go templates](/templates/go-templates) and [Amber](/templates/amber) templates, Hugo supports the powerful Ace templates.

For template documentation, follow the links from the [Ace project](https://github.com/yosssi/ace).

* Ace templates must be named with the ace-suffix, e.g. `list.ace`
* It's possible to use both Go templates and Ace templates side-by-side, and include one into the other
* Full Go template syntax support, including all the useful helper funcs
* Partials can be included both with the Ace and the Go template syntax:
	* `= include partials/foo.html .`[^ace-theme]
	* `{{ partial "foo" . }}`


One noticeable difference between Ace and the others is the inheritance support through [base and inner templates](https://github.com/yosssi/ace/tree/master/examples/base_inner_template).

In Hugo the base template will be chosen with the same ruleset as for [Go templates](/templates/blocks/).


.:
index.ace

./blog:
single.ace
baseof.ace

./_default:
baseof.ace  list.ace  single.ace  single-baseof.ace
```

Some examples for the layout files above:

* Home page: `./index.ace` +  `./_default/baseof.ace`
* Single page in the `blog` section: `./blog/single.ace` +  `./blog/baseof.ace`
* Single page in another section: `./_default/single.ace` +  `./_default/single-baseof.ace`
* Taxonomy page in any section: `./_default/list.ace` +  `./_default/baseof.ace`

**Note:** In most cases one `baseof.ace` in `_default` will suffice.
**Note:** An Ace template without a reference to a base section, e.g. `= content`, will be handled as a standalone template.


[^ace-theme]: Note that the `html` suffix is needed, even if the filename is suffixed `ace`. This does not work from inside a theme, see [issue 763](https://github.com/spf13/hugo/issues/763).

