---
title: Parent
description: Returns the Page object of the parent section of the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/Ancestors
    - methods/page/CurrentSection
    - methods/page/FirstSection
    - methods/page/InSection
    - methods/page/IsAncestor
    - methods/page/IsDescendant
    - methods/page/Sections
  returnType: page.Page
  signatures: [PAGE.Parent]
---

{{% include "methods/page/_common/definition-of-section.md" %}}

{{% note %}}
The parent section of a regular page is the [current section].

[current section]: /methods/page/currentsection/
{{% /note %}}

Consider this content structure:

```text
content/
├── auctions/
│   ├── 2023-11/
│   │   ├── _index.md     <-- parent: auctions
│   │   ├── auction-1.md
│   │   └── auction-2.md  <-- parent: 2023-11
│   ├── 2023-12/
│   │   ├── _index.md     
│   │   ├── auction-3.md
│   │   └── auction-4.md
│   ├── _index.md         <-- parent: home
│   ├── bidding.md
│   └── payment.md        <-- parent: auctions
├── books/
│   ├── _index.md         <-- parent: home
│   ├── book-1.md
│   └── book-2.md         <-- parent: books
├── films/
│   ├── _index.md         <-- parent: home 
│   ├── film-1.md
│   └── film-2.md         <-- parent: films
└── _index.md             <-- parent: nil
```

In the example above, note the parent section of the home page is nil. Code defensively by verifying existence of the parent section before calling methods on its `Page` object. To create a link to the parent section page of the current page:

```go-html-template
{{ with .Parent }}
  <a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
{{ end }}
```
