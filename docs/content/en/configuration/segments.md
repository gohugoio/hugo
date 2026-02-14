---
title: Configure segments
linkTitle: Segments
description: Configure your site for segmented rendering.
categories: []
keywords: []
---

> [!note]
> The `segments` configuration applies only to segmented rendering. While it controls when content is rendered, it doesn't restrict access to Hugo's complete object graph (sites and pages), which remains fully available.

Segmented rendering offers several advantages:

- Faster builds: Process large sites more efficiently.
- Rapid development: Render only a subset of your site for quicker iteration.
- Scheduled rebuilds: Rebuild specific sections at different frequencies (e.g., home page and news hourly, full site weekly).
- Targeted output: Generate specific output formats (like JSON for search indexes).

## Segment definition

Each segment is defined by include and exclude filters:

- Filters: Each segment has zero or more exclude filters and zero or more include filters.
- Matchers: Each filter contains one or more field [glob pattern](g) matchers.
- Logic: Matchers within a filter use AND logic. Filters within a section (include or exclude) use OR logic.

## Filter fields

Available fields for filtering:

kind
: (`string`) A [glob pattern](g) matching the [page kind](g). For example: `{taxonomy,term}`.

sites
: {{< new-in 0.153.0 />}}
: (`map`) A map to define [sites matrix](g).

output
: (`string`) A [glob pattern](g) matching the [output format](g) of the page. For example: `{html,json}`.

path
: (`string`) A [glob pattern](g) matching the page's [logical path](g). For example: `{/books,/books/**}`.

## Example

Place broad filters, such as those for language or output format, in the excludes section. For example:

{{< code-toggle file=hugo >}}
[segments.segment1]
  [[segments.segment1.excludes]]
    lang = "n*"
  [[segments.segment1.excludes]]
    lang   = "en"
    output = "rss"
  [[segments.segment1.includes]]
    kind = "{home,term,taxonomy}"
  [[segments.segment1.includes]]
    path = "{/docs,/docs/**}"
{{< /code-toggle >}}

## Rendering segments

Render specific segments using the [`renderSegments`] configuration or the `--renderSegments` flag:

```bash
hugo --renderSegments segment1
```

You can configure multiple segments and use a comma-separated list with `--renderSegments` to render them all.

```bash
hugo --renderSegments segment1,segment2
```

[`renderSegments`]: /configuration/all/#rendersegments
