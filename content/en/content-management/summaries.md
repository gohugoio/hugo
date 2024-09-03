---
title: Content summaries
linkTitle: Summaries
description: Create and render content summaries.
categories: [content management]
keywords: [summaries,abstracts,read more]
menu:
  docs:
    parent: content-management
    weight: 160
weight: 160
toc: true
aliases: [/content/summaries/,/content-management/content-summaries/]
---

<!-- Do not remove the manual summary divider below. -->
<!-- If you do, you will break its first literal usage on this page. -->
<!--more-->

You can define a content summary manually, in front matter, or automatically. A manual content summary takes precedence over a front matter summary, and a front matter summary takes precedence over an automatic summary.

Review the [comparison table](#comparison) below to understand the characteristics of each summary type.

## Manual summary

Use a `<!--more-->` divider to indicate the end of the content summary. Hugo will not render the summary divider itself.

{{< code file=content/sample.md >}}
+++
title: 'Example'
date: 2024-05-26T09:10:33-07:00
+++

Thénardier was not mistaken. The man was sitting there, and letting
Cosette get somewhat rested.

<!--more-->

The inn-keeper walked round the brushwood and presented himself
abruptly to the eyes of those whom he was in search of.
{{< /code >}}

When using the Emacs Org Mode [content format], use a `# more` divider to indicate the end of the content summary.

[content format]: /content-management/formats/

## Front matter summary

Use front matter to define a summary independent of content.

{{< code file=content/sample.md >}}
+++
title: 'Example'
date: 2024-05-26T09:10:33-07:00
summary: 'Learn more about _Les Misérables_ by Victor Hugo.'
+++

Thénardier was not mistaken. The man was sitting there, and letting
Cosette get somewhat rested. The inn-keeper walked round the
brushwood and presented himself abruptly to the eyes of those whom
he was in search of.
{{< /code >}}

## Automatic summary

If you have not defined the summary manually or in front matter, Hugo automatically defines the summary based on the [`summaryLength`] in your site configuration.

[`summaryLength`]: /getting-started/configuration/#summarylength

{{< code file=content/sample.md >}}
+++
title: 'Example'
date: 2024-05-26T09:10:33-07:00
+++

Thénardier was not mistaken. The man was sitting there, and letting
Cosette get somewhat rested. The inn-keeper walked round the
brushwood and presented himself abruptly to the eyes of those whom
he was in search of.
{{< /code >}}

For example, with a `summaryLength` of 10, the automatic summary will be:

```text
Thénardier was not mistaken. The man was sitting there, and letting
Cosette get somewhat rested.
```

Note that the `summaryLength` is an approximate number of words.

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
