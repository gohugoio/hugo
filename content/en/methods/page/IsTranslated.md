---
title: IsTranslated
description: Reports whether the given page has one or more translations.
categories: []
keywords: []
action:
  related:
   - methods/page/Translations
   - methods/page/AllTranslations
   - methods/page/TranslationKey
  returnType: bool
  signatures: [PAGE.IsTranslated]
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
{{< /code-toggle >}}

And this content:

```text
content/
├── de/
│   ├── books/
│   │   └── book-1.md
│   └── _index.md
├── en/
│   ├── books/
│   │   ├── book-1.md
│   │   └── book-2.md
│   └── _index.md
└── _index.md
```

When rendering content/en/books/book-1.md:

```go-html-template
{{ .IsTranslated }} → true
```

When rendering content/en/books/book-2.md:

```go-html-template
{{ .IsTranslated }} → false
```
