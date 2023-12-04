---
title: AllTranslations
description: Returns all translation of the given page, including the given page. 
categories: []
keywords: []
action:
  related:
   - methods/page/Translations
   - methods/page/IsTranslated
   - methods/page/TranslationKey
  returnType: page.Pages
  signatures: [PAGE.AllTranslations]
---

With this site configuration:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'en'

[languages.en]
contentDir = 'content/en'
languageCode = 'en-US'
languageName = 'English'
weight = 1

[languages.de]
contentDir = 'content/de'
languageCode = 'de-DE'
languageName = 'Deutsch'
weight = 2

[languages.fr]
contentDir = 'content/fr'
languageCode = 'fr-FR'
languageName = 'Français'
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
{{ with .AllTranslations }}
  <ul>
    {{ range . }}
      {{ $lang := .Language.LanguageName}}
      <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }} ({{ $lang }})</a></li>
    {{ end }}
  </ul>
{{ end }}
```

Hugo will render this list on the "Book 1" page of each site:

```html
<ul>
  <li><a href="/books/book-1/">Book 1 (English)</a></li>
  <li><a href="/de/books/book-1/">Book 1 (Deutsch)</a></li>
  <li><a href="/fr/books/book-1/">Book 1 (Français)</a></li>
</ul>
```

On the "Book 2" page of the English and German sites, Hugo will render this:

```html
<ul>
  <li><a href="/books/book-1/">Book 1 (English)</a></li>
  <li><a href="/de/books/book-1/">Book 1 (Deutsch)</a></li>
</ul>
```
