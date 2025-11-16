---
title: Configure ugly URLs
linkTitle: Ugly URLs
description: Configure ugly URLs.
categories: []
keywords: []
---

{{% glossary-term "ugly url" %}}&nbsp;For example:

```text
https://example.org/section/article.html
```

In its default configuration, Hugo generates [pretty URLs](g). For example:

```text
https://example.org/section/article/
```

This is the default configuration:

{{< code-toggle config=uglyURLs />}}

To generate ugly URLs for the entire site:

{{< code-toggle file=hugo >}}
uglyURLs = true
{{< /code-toggle >}}

To generate ugly URLs for specific sections of your site:

{{< code-toggle file=hugo >}}
[uglyURLs]
books = true
films = false
{{< /code-toggle >}}
