---
title: RenderShortcodes
description: Renders all shortcodes in the content of the given page, preserving the surrounding markup.
categories: []
keywords: []
action:
  related:
    - methods/page/Content
    - methods/page/Summary
    - methods/page/ContentWithoutSummary
    - methods/page/RawContent
    - methods/page/Plain
    - methods/page/PlainWords
    - methods/page/RenderString
  returnType: template.HTML
  signatures: [PAGE.RenderShortcodes]
toc: true
---

{{< new-in 0.117.0 >}}

Use this method in shortcode templates to compose a page from multiple content files, while preserving a global context for footnotes and the table of contents.

For example:

{{< code file=layouts/shortcodes/include.html >}}
{{ with site.GetPage (.Get 0) }}
  {{ .RenderShortcodes }}
{{ end }}
{{< /code >}}

Then call the shortcode in your Markdown:

{{< code file=content/about.md lang=md >}}
{{%/* include "/snippets/services.md" */%}}
{{%/* include "/snippets/values.md" */%}}
{{%/* include "/snippets/leadership.md" */%}}
{{< /code >}}

Each of the included Markdown files can contain calls to other shortcodes.

## Shortcode notation

In the example above it's important to understand the difference between the two delimiters used when calling a shortcode:

- `{{</* myshortcode */>}}` tells Hugo that the rendered shortcode does not need further processing. For example, the shortcode content is HTML.
- `{{%/* myshortcode */%}}` tells Hugo that the rendered shortcode needs further processing. For example, the shortcode content is Markdown.

Use the latter for the "include" shortcode described above.

## Further explanation

To understand what is returned by the `RenderShortcodes` method, consider this content file

{{< code file=content/about.md lang=text >}}
+++
title = 'About'
date = 2023-10-07T12:28:33-07:00
+++

{{</* ref "privacy" */>}}

An *emphasized* word.
{{< /code >}}

With this template code:

```go-html-template
{{ $p := site.GetPage "/about" }}
{{ $p.RenderShortcodes }}
```

Hugo renders this:;

```html
https://example.org/privacy/

An *emphasized* word.
```

Note that the shortcode within the content file was rendered, but the surrounding Markdown was preserved.
