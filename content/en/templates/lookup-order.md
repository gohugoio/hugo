---
title: Template lookup order
linkTitle: Lookup order
description: Hugo uses the rules below to select a template for a given page, starting from the most specific.
categories: [templates,fundamentals]
keywords: []
menu:
  docs:
    parent: templates
    weight: 40
weight: 40
toc: true
---

## Lookup rules

Hugo takes the parameters listed below into consideration when choosing a template for a given page. The templates are ordered by specificity. This should feel natural, but look at the table below for concrete examples of the different parameter variations.

Kind
: The page `Kind` (the home page is one). See the example tables below per kind. This also determines if it is a **single page** (i.e. a regular content page. We then look for a template in `_default/single.html` for HTML) or a **list page** (section listings, home page, taxonomy lists, taxonomy terms. We then look for a template in `_default/list.html` for HTML).

Layout
: Can be set in front matter.

Output Format
: See [Custom Output Formats](/templates/output-formats). An output format has both a `name` (e.g. `rss`, `amp`, `html`) and a `suffix` (e.g. `xml`, `html`). We prefer matches with both (e.g. `index.amp.html`), but look for less specific templates.

Note that if the output format's Media Type has more than one suffix defined, only the first is considered.

Language
: We will consider a language tag in the template name. If the site language is `fr`, `index.fr.amp.html` will win over `index.amp.html`, but `index.amp.html` will be chosen before `index.fr.html`.

Type
: Is value of `type` if set in front matter, else it is the name of the root section (e.g. "blog"). It will always have a value, so if not set, the value is "page".

Section
: Is relevant for `section`, `taxonomy` and `term` types.

{{% note %}}
Templates can live in either the project's or the themes' layout folders, and the most specific templates will be chosen. Hugo will interleave the lookups listed below, finding the most specific one either in the project or themes.
{{% /note %}}

## Target a template

You cannot change the lookup order to target a content page, but you can change a content page to target a template. Specify `type`, `layout`, or both in front matter.

Consider this content structure:

```text
content/
├── about.md
└── contact.md
```

Files in the root of the content directory have a [content type] of `page`. To render these pages with a unique template, create a matching subdirectory:

[content type]: /getting-started/glossary/#content-type

```text
layouts/
└── page/
    └── single.html
```

But the contact page probably has a form and requires a different template. In the front matter specify `layout`:

{{< code-toggle file=content/contact.md >}}
title = 'Contact'
layout = 'contact'
{{< /code-toggle >}}

Then create the template for the contact page:

```text
layouts/
└── page/
    └── contact.html  <-- renders contact.md
    └── single.html   <-- renders about.md
```

As a content type, the word `page` is vague. Perhaps `miscellaneous` would be better. Add `type` to the front matter of each page:

{{< code-toggle file=content/about.md >}}
title = 'About'
type = 'miscellaneous'
{{< /code-toggle >}}

{{< code-toggle file=content/contact.md >}}
title = 'Contact'
type = 'miscellaneous'
layout = 'contact'
{{< /code-toggle >}}

Now place the layouts in the corresponding directory:

```text
layouts/
└── miscellaneous/
    └── contact.html  <-- renders contact.md
    └── single.html   <-- renders about.md
```

## Home templates

These template paths are sorted by specificity in descending order. The least specific path is at the bottom of each list.

{{< datatable-filtered "output" "layouts" "Kind == home" "Example" "OutputFormat" "Suffix" "Template Lookup Order" >}}

## Single templates

These template paths are sorted by specificity in descending order. The least specific path is at the bottom of each list.

{{< datatable-filtered "output" "layouts" "Kind == page" "Example" "OutputFormat" "Suffix" "Template Lookup Order" >}}

## Section templates

These template paths are sorted by specificity in descending order. The least specific path is at the bottom of each list.

{{< datatable-filtered "output" "layouts" "Kind == section" "Example" "OutputFormat" "Suffix" "Template Lookup Order" >}}

## Taxonomy templates

These template paths are sorted by specificity in descending order. The least specific path is at the bottom of each list.

The examples below assume the following site configuration:

{{< code-toggle file=hugo >}}
[taxonomies]
category = 'categories'
{{< /code-toggle >}}

{{< datatable-filtered "output" "layouts" "Kind == taxonomy" "Example" "OutputFormat" "Suffix" "Template Lookup Order" >}}

## Term templates

These template paths are sorted by specificity in descending order. The least specific path is at the bottom of each list.

The examples below assume the following site configuration:

{{< code-toggle file=hugo >}}
[taxonomies]
category = 'categories'
{{< /code-toggle >}}

{{< datatable-filtered "output" "layouts" "Kind == term" "Example" "OutputFormat" "Suffix" "Template Lookup Order" >}}

## RSS templates

These template paths are sorted by specificity in descending order. The least specific path is at the bottom of each list.

The examples below assume the following site configuration:

{{< code-toggle file=hugo >}}
[taxonomies]
category = 'categories'
{{< /code-toggle >}}

{{< datatable-filtered "output" "layouts" "OutputFormat == rss" "Example" "OutputFormat" "Suffix" "Template Lookup Order" >}}
