---
title: Draft
description: Reports whether the given page is a draft as defined in front matter.
categories: []
keywords: []
action:
  related: []
  returnType: bool
  signatures: [PAGE.Draft]
---

By default, Hugo does not publish draft pages when you build your site. To include draft pages when you build your site, use the `--buildDrafts` command line flag.

{{< code-toggle file=content/posts/post-1.md fm=true >}}
title = 'Post 1'
draft = true
{{< /code-toggle >}}

```go-html-template
{{ .Draft }} → true
```
