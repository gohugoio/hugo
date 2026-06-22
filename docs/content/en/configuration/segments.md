---
title: Configure segments
linkTitle: Segments
description: Configure your site for segmented rendering.
categories: []
keywords: []
---

> [!NOTE]
> The `segments` configuration applies only to segmented rendering. While it controls when content is rendered, it doesn't restrict access to Hugo's complete object graph (sites and pages), which remains fully available.

Segmented rendering offers several advantages:

- Faster builds: Process large sites more efficiently.
- Rapid development: Render only a subset of your site for quicker iteration.
- Scheduled rebuilds: Rebuild specific sections at different frequencies (e.g., home page and news hourly, full site weekly).
- Targeted output: Generate specific output formats (like JSON for search indexes).

## Segment definition

Each segment is defined by an `includes` key and an `excludes` key, both of which accept an array of filters.

A _filter_ is a collection of one or more conditions, represented as an item in the configuration array. A _condition_ compares a specific page [field](#fields) to a given [glob pattern](g).

### Evaluation rules

The evaluation logic adheres to three rules:

- All conditions within a single filter item must match for that filter to evaluate as true, creating an AND relationship.
- If the `includes` or `excludes` array contains multiple filters, only one filter needs to evaluate as true for the entire array to match, creating an OR relationship.
- The `excludes` array takes absolute precedence. If a page matches any filter in the `excludes` array, Hugo omits it from the segment regardless of whether it matches the `includes` array.

### Performance optimization

Using the `excludes` array to target sites or output formats allows Hugo to skip entire groups of pages during evaluation instead of checking every page. This optimization helps with performance in larger setups.

For example, excluding unwanted output formats is faster:

{{< code-toggle file=hugo >}}
[segments]
  [segments.segment1]
    [[segments.segment1.excludes]]
      output = '! json'
{{< /code-toggle >}}

Including only the desired output format is slower:

{{< code-toggle file=hugo >}}
[segments]
  [segments.segment1]
    [[segments.segment1.includes]]
      output = 'json'
{{< /code-toggle >}}

## Fields

`kind`
: (`string`) A [glob pattern](g) matching the [page kind](g). For example: `{taxonomy,term}`.

`lang`
: {{< deprecated-in 0.153.0 />}}
: Use [`sites`](#sites) instead.

`output`
: (`string`) A [glob pattern](g) matching the [output format](g) of the page. For example: `{html,json}`.

`path`
: (`string`) A [glob pattern](g) matching the page's [logical path](g). For example: `{/books,/books/**}`.

`sites`
: {{< new-in 0.153.0 />}}
: (`map`) A map to define [sites matrix](g).

## Targeting segments

To specify which segments Hugo builds, add the [`renderSegments`][] setting to your project configuration:

{{< code-toggle file=hugo >}}
renderSegments = ['segment1','segment2']
{{< /code-toggle >}}

Alternatively, pass the segment names directly to the `--renderSegments` command-line flag during a build:

```sh
hugo build --renderSegments segment1
```

You can target multiple segments by providing a comma-separated list:

```sh
hugo build --renderSegments segment1,segment2
```

## Example

<!--
To test the example below:

git clone --single-branch -b segmentation-test https://github.com/jmooring/hugo-testing segmentation-test
cd segmentation-test
rm -rf public/ && hugo build --renderSegments segment1 && tree public
-->

Consider a project with this content structure:

```tree
content/
├── books/
│   ├── _index.en.md
│   ├── _index.nb.md
│   ├── _index.nn.md
│   ├── book-1.en.md
│   ├── book-1.nb.md
│   └── book-1.nn.md
├── films/
│   ├── _index.en.md
│   ├── _index.nb.md
│   ├── _index.nn.md
│   ├── film-1.en.md
│   ├── film-1.nb.md
│   └── film-1.nn.md
├── _index.en.md
├── _index.nb.md
└── _index.nn.md
```

And this project configuration:

{{< code-toggle file=hugo >}}
baseURL                        = 'https://example.org/'
title                          = 'Segmentation'
defaultContentLanguage         = 'en'
defaultContentLanguageInSubdir = true

[languages.en]
  direction = 'ltr'
  label     = 'English'
  locale    = 'en-US'
  weight    = 1

[languages.nb]
  locale    = 'nb-NO'
  direction = 'ltr'
  label     = 'Bokmål'
  weight    = 2

[languages.nn]
  locale    = 'nn-NO'
  direction = 'ltr'
  label     = 'Norsk'
  weight    = 3

[segments]
  [segments.segment1]
    [[segments.segment1.excludes]]
      [segments.segment1.excludes.sites.matrix]
        languages = ['n*']
    [[segments.segment1.excludes]]
      output = 'rss'
      [segments.segment1.excludes.sites.matrix]
        languages = ['en']
    [[segments.segment1.includes]]
      kind = '{home,term,taxonomy}'
    [[segments.segment1.includes]]
      path = '{/books,/books/**}'

[taxonomies]
  tag = 'tags'
{{< /code-toggle >}}

When you run this command:

```sh
hugo build --renderSegments segment1
```

The published project has this structure:

```tree
public/
├── en/
│   ├── books/
│   │   ├── book-1/
│   │   │   └── index.html
│   │   └── index.html
│   ├── tags/
│   │   ├── tag-a/
│   │   │   └── index.html
│   │   ├── tag-b/
│   │   │   └── index.html
│   │   └── index.html
│   └── index.html
└── index.html
```

[`renderSegments`]: /configuration/all/#rendersegments
