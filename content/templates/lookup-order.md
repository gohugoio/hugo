---
title: Template Lookup Order
linktitle: Template Lookup Order
description: The lookup order is a prioritized list used by Hugo as it traverses your files looking for the appropriate template to render your content.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-25
categories: [templates]
tags: [lookup,fundamentals]
weight: 15
draft: false
aliases: [/templates/lookup/]
wip: true
---

Before creating your templates, it's important to know how Hugo looks for files within your project's [directory structure][].

Hugo uses a prioritized list called the **lookup order** as it traverses your files *looking* for the appropriate template to render your content.

The template lookup order is an inverted cascade: if template A isn’t present or specified, Hugo will look to template B. If template B isn't present or specified, Hugo will look for template C...and so on until it reaches the `layouts/_default/` directory for your project, or in the case of themes, `themes/<THEME>/layouts/_default/`. In many ways, the lookup order is similar to the [control mechanism of a switch statement (i.e. without fallthrough)][switch] seen in many programming languages.

The power of the lookup order is that it enables you to craft specific layouts as needed without creating more templating than necessary, thereby keeping your templating [DRY][].

{{% note %}}
Most Hugo websites will only need the default template files at the end of the lookup order (i.e. `_default/*.html`).
{{% /note %}}

See examples of the lookup order for each of the Hugo template types:

* [Homepage Template][home]
* [Base Templates][base]
* [Section Page Templates][sectionlookup]
* [Taxonomy List Templates][taxonomylookup]
* [Taxonomy Terms Templates][termslookup]
* [Single Page Templates][singlelookup]
* [RSS Templates][rsslookup]

[base]: /templates/base/#base-template-lookup-order
[directory structure]: /getting-started/directory-structure/
[DRY]: https://en.wikipedia.org/wiki/Don%27t_repeat_yourself
[home]: /templates/homepage/#homepage-template-lookup-order
[rsslookup]: /templates/rss/#rss-template-lookup-order
[sectionlookup]: /templates/section-templates/#section-template-lookup-order
[singlelookup]: templates/single-page-templates/#single-page-template-lookup-order
[switch]: https://en.wikipedia.org/wiki/Switch_statement#Fallthrough
[taxonomylookup]: /templates/taxonomy-templates/#taxonomy-list-template-lookup-order
[termslookup]: /templates/taxonomy-templates/#taxonomy-terms-template-lookup-order