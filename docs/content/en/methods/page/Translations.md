---
title: Translations
description: Returns all translations of the given page, excluding the current language, sorted by language weight then language name.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.Pages
    signatures: [PAGE.Translations]
---

With this project configuration:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'en'

[languages.en]
contentDir = 'content/en'
label = 'English'
locale = 'en-US'
weight = 1

[languages.de]
contentDir = 'content/de'
label = 'Deutsch'
locale = 'de-DE'
weight = 2

[languages.fr]
contentDir = 'content/fr'
label = 'Français'
locale = 'fr-FR'
weight = 3
{{< /code-toggle >}}

And this content:

```text
content/
├── de/
│   ├── books/
│   │   ├── book-1.md
│   │   └── book-2.md
│   └── _index.md
├── en/
│   ├── books/
│   │   ├── book-1.md
│   │   └── book-2.md
│   └── _index.md
├── fr/
│   ├── books/
│   │   └── book-1.md
│   └── _index.md
└── _index.md
```

And this template:

```go-html-template
{{ with .Translations }}
  <ul>
    {{ range . }}
      <li>
        <a href="{{ .RelPermalink }}" hreflang="{{ .Language.Locale }}">{{ .LinkTitle }} ({{ or .Language.Label .Language.Name }})</a>
      </li>
    {{ end }}
  </ul>
{{ end }}
```

Hugo will render this list on the "Book 1" page of the English site:

```html
<ul>
  <li><a href="/de/books/book-1/" hreflang="de-DE">Book 1 (Deutsch)</a></li>
  <li><a href="/fr/books/book-1/" hreflang="fr-FR">Book 1 (Français)</a></li>
</ul>
```

Hugo will render this list on the "Book 2" page of the English site:

```html
<ul>
  <li><a href="/de/books/book-1/" hreflang="de-DE">Book 1 (Deutsch)</a></li>
</ul>
```
