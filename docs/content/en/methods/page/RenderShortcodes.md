---
title: RenderShortcodes
description: Renders all shortcodes in the content of the given page, preserving the surrounding markup.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: template.HTML
    signatures: [PAGE.RenderShortcodes]
---

Use this method in _shortcode_ templates to compose a page from multiple content files, while preserving a global context for footnotes and the table of contents.

For example:

```go-html-template {file="layouts/_shortcodes/include.html" copy=true}
{{ with .Get 0 }}
  {{ with $.Page.GetPage . }}
    {{- .RenderShortcodes }}
  {{ else }}
    {{ errorf "The %q shortcode was unable to find %q. See %s" $.Name . $.Position }}
  {{ end }}
{{ else }}
  {{ errorf "The %q shortcode requires a positional parameter indicating the logical path of the file to include. See %s" .Name .Position }}
{{ end }}
```

Then call the shortcode in your Markdown:

```text {file="content/about.md"}
{{%/* include "/snippets/services" */%}}
{{%/* include "/snippets/values" */%}}
{{%/* include "/snippets/leadership" */%}}
```

Each of the included Markdown files can contain calls to other shortcodes.

## Shortcode notation

In the example above it's important to understand the difference between the two delimiters used when calling a shortcode:

- `{{</* myshortcode */>}}` tells Hugo that the rendered shortcode does not need further processing. For example, the shortcode content is HTML.
- `{{%/* myshortcode */%}}` tells Hugo that the rendered shortcode needs further processing. For example, the shortcode content is Markdown.

Use the latter for the "include" shortcode described above.

## Further explanation

To understand what is returned by the `RenderShortcodes` method, consider this content file

```text {file="content/about.md"}
+++
title = 'About'
date = 2023-10-07T12:28:33-07:00
+++

{{</* ref "privacy" */>}}

An *emphasized* word.
```

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

## Limitations

The primary use case for `.RenderShortcodes` is inclusion of Markdown content. If you try to use `.RenderShortcodes` inside `HTML` blocks when inside Markdown, you will get a warning similar to this:

```text
WARN .RenderShortcodes detected inside HTML block in "/content/mypost.md"; this may not be what you intended ...
```

The above warning can be turned off is this is what you really want.
