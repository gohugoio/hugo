---
title: Ancestors
description: Returns a collection of Page objects, one for each ancestor section of the given page. 
categories: []
keywords: []
action:
  related:
    - methods/page/CurrentSection
    - methods/page/FirstSection
    - methods/page/InSection
    - methods/page/IsAncestor
    - methods/page/IsDescendant
    - methods/page/Parent
    - methods/page/Sections
  returnType: page.Pages
  signatures: [PAGE.Ancestors]
---

{{< new-in 0.109.0 >}}

{{% include "methods/page/_common/definition-of-section.md" %}}

With this content structure:

```text
content/
├── auctions/
│   ├── 2023-11/
│   │   ├── _index.md     <-- front matter: weight = 202311
│   │   ├── auction-1.md
│   │   └── auction-2.md
│   ├── 2023-12/
│   │   ├── _index.md     <-- front matter: weight = 202312
│   │   ├── auction-3.md
│   │   └── auction-4.md
│   ├── _index.md         <-- front matter: weight = 30
│   ├── bidding.md
│   └── payment.md
├── books/
│   ├── _index.md         <-- front matter: weight = 10
│   ├── book-1.md
│   └── book-2.md
├── films/
│   ├── _index.md         <-- front matter: weight = 20
│   ├── film-1.md
│   └── film-2.md
└── _index.md
```

And this template:

```go-html-template
{{ range .Ancestors }}
  <a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
{{ end }}
```

On the November 2023 auctions page, Hugo renders:

```html
<a href="/auctions/2023-11/">Auctions in November 2023</a>
<a href="/auctions/">Auctions</a>
<a href="/">Home</a>
```

In the example above, notice that Hugo orders the ancestors from closest to furthest. This makes breadcrumb navigation simple:

```go-html-template
<nav aria-label="breadcrumb" class="breadcrumb">
  <ol>
    {{ range .Ancestors.Reverse }}
      <li>
        <a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
      </li>
    {{ end }}
    <li class="active">
      <a aria-current="page" href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
    </li>
  </ol>
</nav>
```

With some CSS, the code above renders something like this, where each breadcrumb links to its page:

```text
Home > Auctions > Auctions in November 2023 > Auction 1
```
