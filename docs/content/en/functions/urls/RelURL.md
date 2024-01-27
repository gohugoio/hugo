---
title: urls.RelURL
description: Returns a relative URL.
categories: []
keywords: []
action:
  aliases: [relURL]
  related:
    - functions/urls/AbsLangURL
    - functions/urls/AbsURL 
    - functions/urls/RelLangURL
  returnType: string
  signatures: [urls.RelURL INPUT]
aliases: [/functions/relurl]
---

With multilingual configurations, use the [`relLangURL`] function instead. The URL returned by this function depends on:

- Whether the input begins with a slash
- The `baseURL` in site configuration

### Input does not begin with a slash

If the input does not begin with a slash, the resulting URL will be correct regardless of the `baseURL`.

With `baseURL = https://example.org/`

```go-html-template
{{ relURL "" }}           →   /
{{ relURL "articles" }}   →   /articles
{{ relURL "style.css" }}  →   /style.css
```

With `baseURL = https://example.org/docs/`

```go-html-template
{{ relURL "" }}           →   /docs/
{{ relURL "articles" }}   →   /docs/articles
{{ relURL "style.css" }}  →   /docs/style.css
```

#### Input begins with a slash

If the input begins with a slash, the resulting URL will be incorrect when the `baseURL` includes a subdirectory. With a leading slash, the function returns a URL relative to the protocol+host section of the `baseURL`.

With `baseURL = https://example.org/`

```go-html-template
{{ relURL "/" }}          →   /
{{ relURL "/articles" }}  →   /articles
{{ relURL "style.css" }}  →   /style.css
```

With `baseURL = https://example.org/docs/`

```go-html-template
{{ relURL "/" }}          →   /
{{ relURL "/articles" }}  →   /articles
{{ relURL "/style.css" }} →   /style.css
```

{{% note %}}
The last three examples are not desirable in most situations. As a best practice, never include a leading slash when using this function.
{{% /note %}}

[`relLangURL`]: /functions/urls/rellangurl/
