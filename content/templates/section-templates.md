---
title: Section Page Templates
linktitle: Section Page Templates
description: Templates used for section pages are lists and therefore have all the variables and methods available to list pages.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [lists,sections]
weight: 40
draft: false
aliases: [/templates/sections/]
toc: true
wip: true
---

Templates used for section pages are *lists* and therefore have all the variables and methods available to [list pages][lists].

{{% note "Section Pages Pull Content from `_index.md`" %}}
To effectively leverage section page templates, you should first understand Hugo's [content organization][contentorg] and, specifically, the purpose of `_index.md` for adding content and front matter to section and other list pages.
{{% /note %}}

## Section Template Lookup Order

The [lookup order][lookup] for section pages is as follows:

1. `/layouts/section/<SECTION>.html`
2. `/layouts/<SECTION>/list.html`
2. `/layouts/_default/section.html`
3. `/layouts/_default/list.html`
4. `/themes/<THEME>/layouts/section/<SECTION>.html`
5. `/themes/<THEME>/layouts/<SECTION>/list.html`
5. `/themes/<THEME>/layouts/_default/section.html`
6. `/themes/<THEME>/layouts/_default/list.html`

## `.Site.GetPage` with Sections

Every `Page` in Hugo has a `.Kind` attribute. `Kind` can easily be combined with the [`where` function][where] in your templates to create kind-specific lists of content. This method is ideal for creating lists, but there are times where you may want to fetch just the index page of a single section via the section's path.

The [`.GetPage` function][getpage] looks up an index page of a given `Kind` and `path`.

{{% note %}}
`.GetPage` is only supported in section page templates but *may* be supported in [single page templates](/templates/single-page-templates/) in the future.
{{% /note %}}

You can call `.Site.GetPage` with two arguments: `kind` and `kind value`.

These are the valid values for 'kind':

1. `home`
2. `section`
3. `taxonomy`
4. `taxonomyTerm`

### Example: Using `.Site.GetPage`

The `.Site.GetPage` example that follows assumes the following project directory structure:

```bash
.
└── content
    ├── blog
    │   ├── _index.md # "title: My Hugo Blog" in the front matter
    │   ├── post-1.md
    │   ├── post-2.md
    │   └── post-3.md
    └── events #Note there is no _index.md file in "events"
        ├── event-1.md
        └── event-2.md
```

`.Site.GetPage` will return `nil` if no `_index.md` page is found. Therefore, if `content/blog/_index.md` does not exist, the template will output the section name:

```html
<h1>{{ with .Site.GetPage "section" "blog" }}{{ .Title }}{{ end }}</h1>
```

Since `blog` has a section index page with front matter at `content/blog/_index.md`, the above code will return the following result:

```html
<h1>My Hugo Blog</h1>
```

If we try the same code with the `events` section, however, Hugo will default to the section title because there is no `content/events/_index.md` from which to pull content and front matter:

```html
<h1>{{ with .Site.GetPage "section" "events" }}{{ .Title }}{{ end }}</h1>
```

Which then returns

```html
<h1>Events</h1>
```

## Nested Sections

**Content forthcoming**

[contentorg]: /content-management/organization/
[getpage]: /functions/getpage/
[lists]: /templates/lists/
[lookup]: /templates/lookup-order/
[where]: /functions/where/