---
title: urls.AbsLangURL
linkTitle: absLangURL
description: Returns an absolute URL with a language prefix, if any.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [absLangURL]
  returnType: template.HTML
  signatures: [urls.AbsLangURL INPUT]
relatedFunctions:
  - urls.AbsLangURL
  - urls.AbsURL 
  - urls.RelLangURL
  - urls.RelURL
aliases: [/functions/abslangurl]
---

Use this function with both monolingual and multilingual configurations. The URL returned by this function depends on:

- Whether the input begins with a slash
- The `baseURL` in site configuration
- The language prefix, if any

In examples that follow, the project is multilingual with content in both Español (`es`) and English (`en`). The default language is Español. The returned values are from the English site.

### Input does not begin with a slash

If the input does not begin with a slash, the resulting URL will be correct regardless of the `baseURL`.

With `baseURL = https://example.org/`

```go-html-template
{{ absLangURL "" }}           →   https://example.org/en/
{{ absLangURL "articles" }}   →   https://example.org/en/articles
{{ absLangURL "style.css" }}  →   https://example.org/en/style.css
```

With `baseURL = https://example.org/docs/`

```go-html-template
{{ absLangURL "" }}           →   https://example.org/docs/en/
{{ absLangURL "articles" }}   →   https://example.org/docs/en/articles
{{ absLangURL "style.css" }}  →   https://example.org/docs/en/style.css
```

### Input begins with a slash

If the input begins with a slash, the resulting URL will be incorrect when the `baseURL` includes a subdirectory. With a leading slash, the function returns a URL relative to the protocol+host section of the `baseURL`.

With `baseURL = https://example.org/`

```go-html-template
{{ absLangURL "/" }}          →   https://example.org/en/
{{ absLangURL "/articles" }}  →   https://example.org/en/articles
{{ absLangURL "/style.css" }} →   https://example.org/en/style.css
```

With `baseURL = https://example.org/docs/`

```go-html-template
{{ absLangURL "/" }}          →   https://example.org/en/
{{ absLangURL "/articles" }}  →   https://example.org/en/articles
{{ absLangURL "/style.css" }} →   https://example.org/en/style.css
```

{{% note %}}
The last three examples are not desirable in most situations. As a best practice, never include a leading slash when using this function.
{{% /note %}}
