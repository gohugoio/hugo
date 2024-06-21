---
title: Pages
description: Returns a collection of regular pages within the current section, and section pages of immediate descendant sections.
categories: []
keywords: []
action:
  related:
    - methods/page/RegularPages
    - methods/page/RegularPagesRecursive
  returnType: page.Pages
  signatures: [PAGE.Pages]
---

The `Pages` method on a `Page` object is available to these [page kinds]: `home`, `section`, `taxonomy`, and `term`. The templates for these page kinds receive a page [collection] in [context].

Range through the page collection in your template:

```go-html-template
{{ range .Pages.ByTitle }}
  <h2><a href="{{ .RelPermalink }}">{{ .Title }}</a></h2>
{{ end }}
```

Consider this content structure:

```text
content/
├── lessons/
│   ├── lesson-1/
│   │   ├── _index.md
│   │   ├── part-1.md
│   │   └── part-2.md
│   ├── lesson-2/
│   │   ├── resources/
│   │   │   ├── task-list.md
│   │   │   └── worksheet.md
│   │   ├── _index.md
│   │   ├── part-1.md
│   │   └── part-2.md
│   ├── _index.md
│   ├── grading-policy.md
│   └── lesson-plan.md
├── _index.md
├── contact.md
└── legal.md
```

When rendering the home page, the `Pages` method returns:

    contact.md
    legal.md
    lessons/_index.md

When rendering the lessons page, the `Pages` method returns:

    lessons/grading-policy.md
    lessons/lesson-plan.md
    lessons/lesson-1/_index.md
    lessons/lesson-2/_index.md

When rendering lesson-1, the `Pages` method returns:

    lessons/lesson-1/part-1.md
    lessons/lesson-1/part-2.md

When rendering lesson-2, the `Pages` method returns:

    lessons/lesson-2/part-1.md
    lessons/lesson-2/part-2.md
    lessons/lesson-2/resources/task-list.md
    lessons/lesson-2/resources/worksheet.md

In the last example, the collection includes pages in the resources subdirectory. That directory is not a [section]---it does not contain an _index.md file. Its contents are part of the lesson-2 section.

{{% note %}}
When used with a `Site` object, the `Pages` method recursively returns all pages within the site. See&nbsp;[details].

[details]: /methods/site/pages/
{{% /note %}}

```go-html-template
{{ range .Site.Pages.ByTitle }}
  <h2><a href="{{ .RelPermalink }}">{{ .Title }}</a></h2>
{{ end }}
```

[collection]: /getting-started/glossary/#collection
[context]: /getting-started/glossary/#context
[page kinds]: /getting-started/glossary/#page-kind
[section]: /getting-started/glossary/#section
