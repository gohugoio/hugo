---
title: Translations
description: Returns all translations of the given page, excluding the current language.  
categories: []
keywords: []
action:
  related:
   - methods/page/AllTranslations
   - methods/page/IsTranslated
   - methods/page/TranslationKey
  returnType: page.Pages
  signatures: [PAGE.Translations]
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
{{ with .Translations }}
  <ul>
    {{ range . }}
      <li>
        <a href="{{ .RelPermalink }}" hreflang="{{ .Language.LanguageCode }}">{{ .LinkTitle }} ({{ or .Language.LanguageName .Language.Lang }})</a>
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
