---
title: Content summaries
linkTitle: Summaries
description: Create and render content summaries.
categories: []
keywords: []
aliases: [/content/summaries/,/content-management/content-summaries/]
---

<!-- Do not remove the manual summary divider below. -->
<!-- If you do, you will break its first literal usage on this page. -->

<!--more-->

You can define a summary manually, in front matter, or automatically. A manual summary takes precedence over a front matter summary, and a front matter summary takes precedence over an automatic summary.

Review the [comparison table](#comparison) below to understand the characteristics of each summary type.

## Manual summary

Use a `<!--more-->` divider to indicate the end of the summary. Hugo will not render the summary divider itself.

```text {file="content/example.md"}
+++
title: 'Example'
date: 2024-05-26T09:10:33-07:00
+++

This is the first paragraph.

<!--more-->

This is the second paragraph.
```

> [!NOTE]
> Place the summary divider on its own line. Do not place it inline with other content.

Correct placement:

```text {file="content/example.md"}
---
title: 'Example'
---

This is an example of **strong text** in a sentence. This is another sentence.

<!--more-->

This is another paragraph.
```

Incorrect placement:

```text {file="content/example.md"}
---
title: 'Example'
---

This is an example of **strong text** <!--more--> in a sentence. This is another sentence.

This is another paragraph.
```

When using the Emacs Org Mode [content format], use a `# more` divider to indicate the end of the summary.

[content format]: /content-management/formats/

## Front matter summary

Use front matter to define a summary independent of content.

```text {file="content/example.md"}
+++
title: 'Example'
date: 2024-05-26T09:10:33-07:00
summary: 'This summary is independent of the content.'
+++

This is the first paragraph.

This is the second paragraph.
```

## Automatic summary

If you do not define the summary manually or in front matter, Hugo automatically defines the summary based on the [`summaryLength`] in your site configuration.

[`summaryLength`]: /configuration/all/#summarylength

```text {file="content/example.md"}
+++
title: 'Example'
date: 2024-05-26T09:10:33-07:00
+++

This is the first paragraph.

This is the second paragraph.

This is the third paragraph.
```

For example, with a `summaryLength` of 7, the automatic summary will be:

```html
<p>This is the first paragraph.</p>
<p>This is the second paragraph.</p>
```

> [!warning]
> Automatic `.Summary` may cut block tags (e.g., `blockquote`) in the middle when `summaryLength` is reached, causing the browser to recover the end tag (the end tag will be inserted before the parent's end tag), resulting in unexpected rendering behavior. To avoid this, wrap `.Summary` in a `<div>`; alternatively, wrap it together with the heading tag using `<section>`. You can avoid this entirely by using a manual summary. See issue [#14044] for details.

## Comparison

Each summary type has different characteristics:

Type|Precedence|Renders markdown|Renders shortcodes|Wraps single lines with `<p>`
:--|:-:|:-:|:-:|:-:
Manual|1|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:
Front&nbsp;matter|2|:heavy_check_mark:|:x:|:x:
Automatic|3|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:

## Rendering

Render the summary in a template by calling the [`Summary`] method on a `Page` object.

[`Summary`]: /methods/page/summary

```go-html-template
{{ range site.RegularPages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  <div class="summary">
    {{ .Summary }}
    {{ if .Truncated }}
      <a href="{{ .RelPermalink }}">More ...</a>
    {{ end }}
  </div>
{{ end }}
```

## Alternative

Instead of calling the `Summary` method on a `Page` object, use the [`strings.Truncate`] function for granular control of the summary length. For example:

[`strings.Truncate`]: /functions/strings/truncate/

```go-html-template
{{ range site.RegularPages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  <div class="summary">
    {{ .Content | strings.Truncate 42 }}
  </div>
{{ end }}
```

[#14044]: https://github.com/gohugoio/hugo/issues/14044
