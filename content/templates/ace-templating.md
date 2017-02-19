---
title: Ace Templating
linktitle:
description: In addition to Go templates and Amber templates, Hugo supports the powerful Ace templating from @yosssi.
godocref: https://godoc.org/github.com/yosssi/ace
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
weight:
tags: []
categories: [amber, templating options]
draft: false
aliases: []
toc: false
notesforauthors:
---

In addition to [Go templates][] and [Amber templates][], Hugo supports the powerful Ace templates.

For template documentation, follow the links from the [Ace project](https://github.com/yosssi/ace).

* Ace templates must be named with the ace-suffix; e.g., `list.ace`
* It's possible to use both Go templates and Ace templates side by side and even include one into the other
* Full Go template syntax support, including all the useful helper [template functions][]
* Partials can be included both with the Ace and the Go template syntax. For example, the following two will have the same output in Ace:
    * `= include partials/foo.html .`
    * `{{ partial "foo" . }}`

One noticeable difference between Ace and the other templating engines in Hugo is [Ace's inheritance support through base and inner templates][].

In Hugo, the Ace base template will be chosen with the same rule set as for [Go base templates][].

```
.:
index.ace

./blog:
single.ace
baseof.ace

./_default:
baseof.ace  list.ace  single.ace  single-baseof.ace
```

{{% note %}}
Note that the `html` suffix is needed, even if the filename is suffixed `ace`. This does not work from inside a theme ([issue #763](https://github.com/spf13/hugo/issues/763)).
{{% /note %}}

Some examples for the layout files above:

* Home page: `./index.ace` +  `./_default/baseof.ace`
* Single page in the `blog` section: `./blog/single.ace` +  `./blog/baseof.ace`
* Single page in another section: `./_default/single.ace` +  `./_default/single-baseof.ace`
* Taxonomy page in any section: `./_default/list.ace` +  `./_default/baseof.ace`

{{% note %}}
In most cases, one `baseof.ace` in `_default` will suffice. An Ace template without a reference to a base section (e.g., `= content`) will be handled as a standalone template.
{{% /note %}}

[Ace's inheritance support through base and inner templates]: https://github.com/yosssi/ace/tree/master/examples/base_inner_template
[Amber templates]: /templates/amber-templating/
[template functions]: /functions/
[Go templates]: /templates/go-template-primer/
[Go base templates]: /templates/base-templates-and-blocks/