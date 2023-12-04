---
title: urls.AbsURL 
description: Returns an absolute URL.
categories: []
keywords: []
action:
  aliases: [absURL]
  related:
    - functions/urls/AbsLangURL
    - functions/urls/RelLangURL
    - functions/urls/RelURL
  returnType: string
  signatures: [urls.AbsURL INPUT]
aliases: [/functions/absurl]
---

With multilingual configurations, use the [`absLangURL`] function instead. The URL returned by this function depends on:

- Whether the input begins with a slash
- The `baseURL` in site configuration

### Input does not begin with a slash

If the input does not begin with a slash, the resulting URL will be correct regardless of the `baseURL`.

With `baseURL = https://example.org/`

```go-html-template
{{ absURL "" }}           →   https://example.org/
{{ absURL "articles" }}   →   https://example.org/articles
{{ absURL "style.css" }}  →   https://example.org/style.css
```

With `baseURL = https://example.org/docs/`

```go-html-template
{{ absURL "" }}           →   https://example.org/docs/
{{ absURL "articles" }}   →   https://example.org/docs/articles
{{ absURL "style.css" }}  →   https://example.org/docs/style.css
```

#### Input begins with a slash

If the input begins with a slash, the resulting URL will be incorrect when the `baseURL` includes a subdirectory. With a leading slash, the function returns a URL relative to the protocol+host section of the `baseURL`.

With `baseURL = https://example.org/`

```go-html-template
{{ absURL "/" }}          →   https://example.org/
{{ absURL "/articles" }}  →   https://example.org/articles
{{ absURL "/style.css" }} →   https://example.org/style.css
```

With `baseURL = https://example.org/docs/`

```go-html-template
{{ absURL "/" }}          →   https://example.org/
{{ absURL "/articles" }}  →   https://example.org/articles
{{ absURL "/style.css" }} →   https://example.org/style.css
```

{{% note %}}
The last three examples are not desirable in most situations. As a best practice, never include a leading slash when using this function.
{{% /note %}}

[`absLangURL`]: /functions/urls/abslangurl/
