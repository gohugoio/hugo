---
title: BundleType
description: Returns the bundle type of the given page, or an empty string if the page is not a page bundle.
categories: []
keywords: []
action:
  related: []
  returnType: string
  signatures: [PAGE.BundleType]
---

A page bundle is a directory that encapsulates both content and associated [resources]. There are two types of page bundles: [leaf bundles] and [branch bundles]. See&nbsp;[details](/content-management/page-bundles/).

The `BundleType` method on a `Page` object returns `branch` for branch bundles, `leaf` for leaf bundles, and an empty string if the page is not a page bundle.

```text
content/
├── films/
│   ├── film-1/
│   │   ├── a.jpg
│   │   └── index.md  <-- leaf bundle
│   ├── _index.md     <-- branch bundle
│   ├── b.jpg
│   ├── film-2.md
│   └── film-3.md
└── _index.md         <-- branch bundle
```

To get the value within a template:

```go-html-template
{{ .BundleType }}
```

[resources]: /getting-started/glossary/#resource
[leaf bundles]: /getting-started/glossary/#leaf-bundle
[branch bundles]: /getting-started/glossary/#branch-bundle
