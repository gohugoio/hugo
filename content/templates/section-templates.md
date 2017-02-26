---
title: Section Page Templates
linktitle: Section Page Templates
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [lists,sections]
weight: 40
draft: false
aliases: []
toc: true
needsreview: true
---

Section page templates are lists and therefore have all the variables and methods available to [list pages][lists].

{{% warning "Section Pages Pull Content from `_index.md`" %}}
To effectively leverage section page templates, you should first understand the Hugo [content organization](/content-management/content-organization/), and specifically the purpose of `_index.md`.
{{% /warning %}}

### Section Template Lookup Order

The [lookup order][lookup] for section pages

* /layouts/section/<SECTION>.html
* /layouts/\_default/section.html
* /layouts/\_default/list.html
* /themes/<THEME>/layouts/section/`SECTION`.html
* /themes/`THEME`/layouts/\_default/section.html
* /themes/`THEME`/layouts/\_default/list.html

Note that a sections list page can also have a content file with frontmatter,  see [Source Organization](/overview/source-directory/}}).

## `.Site.GetPage`

Every `Page` in Hugo has a `.Kind` attribute. `Kind` can easily be combined with [`where`](/functions/where/) in your templates to create kind-specific lists of content, but there are times where you may want to fetch the index page of a single section by the section's path.

[`.GetPage`](/function/getpage/) looks up an index page (i.e `_index.md`) of a given `Kind` and `path`. This method is only supported in section page templates but *may* support [single page templates][singlepages] in the future.

`.Site.GetPage` takes two arguments: `kind` and `kind value`.

The valid values for 'kind' are as follows:

1. `home`
2. `section`
3. `taxonomy`
4. `taxonomyTerm`

### `.Site.GetPage` Example

The `.Site.GetPage` example assumes the following project directory structure:

{{% code file="grab-blog-section-index-page-title.html" %}}
```golang
{{ with .Site.GetPage "section" "blog" }}{{ .Title }}{{ end }}
```
{{% /code %}}

`.Site.GetPage` will return `nil` if no `_index.md` page is found. If `content/blog/_index.md` does not exist, the template will output a blank section where `{{.Title}}` should have been in the preceding example.


[lists]: /templates/lists/
[lookup]: /templates/lookup-order/