---
title: ReadingTime
description: Returns the estimated reading time, in minutes, for the given page.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: int
    signatures: [PAGE.ReadingTime]
---

Hugo calculates the estimated reading time by dividing the number of words in the content by a reading speed of 212 words per minute.

> [!NOTE]
> For content in [CJK](g) languages, set [`hasCJKLanguage`][] to `true` in your project configuration. When enabled, Hugo applies CJK word counting rules and a reading speed of 500 words per minute to pages containing CJK characters. To override this behavior on a given page, set the [`isCJKLanguage`][] field in its front matter.

```go-html-template
{{ printf "Estimated reading time: %d minutes" .ReadingTime }}
```

Reading speed varies by language. Create language-specific estimated reading times on your multilingual project using site parameters.

{{< code-toggle file=hugo >}}
[languages]
  [languages.de]
    contentDir = 'content/de'
    label = 'Deutsch'
    locale = 'de-DE'
    weight = 2
    [languages.de.params]
    reading_speed = 179
  [languages.en]
    contentDir = 'content/en'
    label = 'English'
    locale = 'en-US'
    weight = 1
    [languages.en.params]
      reading_speed = 228
{{< /code-toggle >}}

Then in your template:

```go-html-template
{{ $readingTime := div (float .WordCount) .Site.Params.reading_speed }}
{{ $readingTime = math.Ceil $readingTime }}
```

We cast the `.WordCount` to a float to obtain a float when we divide by the reading speed. Then round up to the nearest integer.

[`hasCJKLanguage`]: /configuration/all/#hascjklanguage
[`isCJKLanguage`]: /content-management/front-matter/#iscjklanguage
