---
title: RegularPagesRecursive
description: Returns a collection of regular pages within the current section, and regular pages within all descendant sections.
categories: []
keywords: []
action:
  related:
    - methods/page/Pages
    - methods/page/RegularPages
  returnType: page.Pages
  signatures: [PAGE.RegularPagesRecursive]
---

The `RegularPagesRecursive` method on a `Page` object is available to these [page kinds]: `home`, `section`, `taxonomy`, and `term`. The templates for these page kinds receive a page [collection] in [context].

Range through the page collection in your template:

```go-html-template
{{ range .RegularPagesRecursive.ByTitle }}
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

When rendering the home page, the `RegularPagesRecursive` method returns:

    contact.md
    lessons/grading-policy.md
    legal.md
    lessons/lesson-plan.md
    lessons/lesson-2/part-1.md
    lessons/lesson-1/part-1.md
    lessons/lesson-2/part-2.md
    lessons/lesson-1/part-2.md
    lessons/lesson-2/resources/task-list.md
    lessons/lesson-2/resources/worksheet.md

When rendering the lessons page, the `RegularPagesRecursive` method returns:

    lessons/grading-policy.md
    lessons/lesson-plan.md
    lessons/lesson-2/part-1.md
    lessons/lesson-1/part-1.md
    lessons/lesson-2/part-2.md
    lessons/lesson-1/part-2.md
    lessons/lesson-2/resources/task-list.md
    lessons/lesson-2/resources/worksheet.md

When rendering lesson-1, the `RegularPagesRecursive` method returns:

    lessons/lesson-1/part-1.md
    lessons/lesson-1/part-2.md

When rendering lesson-2, the `RegularPagesRecursive` method returns:

    lessons/lesson-2/part-1.md
    lessons/lesson-2/part-2.md
    lessons/lesson-2/resources/task-list.md
    lessons/lesson-2/resources/worksheet.md

{{% note %}}
The `RegularPagesRecursive` method in not available on a `Site` object.
{{% /note %}}

[collection]: /getting-started/glossary/#collection
[context]: /getting-started/glossary/#context
[page kinds]: /getting-started/glossary/#page-kind
