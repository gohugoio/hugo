---
title: TranslationKey
description: Returns the translation key of the given page.
categories: []
keywords: []
action:
  related:
   - methods/page/Translations
   - methods/page/AllTranslations
   - methods/page/IsTranslated
  returnType: string
  signatures: [PAGE.TranslationKey]
---

The translation key creates a relationship between all translations of a given page. The translation key is derived from the file path, or from the `translationKey` parameter if defined in front matter.

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
│   │   ├── buch-1.md
│   │   └── book-2.md
│   └── _index.md
├── en/
│   ├── books/
│   │   ├── book-1.md
│   │   └── book-2.md
│   └── _index.md
└── _index.md
```

And this front matter:

{{< code-toggle file=content/en/books/book-1.md fm=true >}}
title = 'Book 1'
translationKey = 'foo'
{{< /code-toggle >}}

{{< code-toggle file=content/de/books/buch-1.md fm=true >}}
title = 'Buch 1'
translationKey = 'foo'
{{< /code-toggle >}}

When rendering either either of the pages above:

```go-html-template
{{ .TranslationKey }} → page/foo
```

If the front matter of Book 2, in both languages, does not include a translation key:

```go-html-template
{{ .TranslationKey }} → page/books/book-2
```
