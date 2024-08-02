---
title: URL management
description: Control the structure and appearance of URLs through front matter entries and settings in your site configuration.
categories: [content management]
keywords: [aliases,redirects,permalinks,urls]
menu:
  docs:
    parent: content-management
    weight: 180
weight: 180
toc: true
aliases: [/extras/permalinks/,/extras/aliases/,/extras/urls/,/doc/redirects/,/doc/alias/,/doc/aliases/]
---

## Overview

By default, when Hugo renders a page, the resulting URL matches the file path within the `content` directory. For example:

```text
content/posts/post-1.md → https://example.org/posts/post-1/
```

You can change the structure and appearance of URLs with front matter values and site configuration options.

## Front matter

### `slug`

Set the `slug` in front matter to override the last segment of the path. The `slug` value does not affect section pages.

{{< code-toggle file=content/posts/post-1.md fm=true >}}
title = 'My First Post'
slug = 'my-first-post'
{{< /code-toggle >}}

The resulting URL will be:

```text
https://example.org/posts/my-first-post/
```

### `url`

Set the `url` in front matter to override the entire path. Use this with either regular pages or section pages.

With this front matter:

{{< code-toggle file=content/posts/post-1.md fm=true >}}
title = 'My First Article'
url = '/articles/my-first-article'
{{< /code-toggle >}}

The resulting URL will be:

```text
https://example.org/articles/my-first-article/
```

If you include a file extension:

{{< code-toggle file=content/posts/post-1.md fm=true >}}
title = 'My First Article'
url = '/articles/my-first-article.html'
{{< /code-toggle >}}

The resulting URL will be:

```text
https://example.org/articles/my-first-article.html
```

In a monolingual site, a `url` value with or without a leading slash is relative to the `baseURL`.

In a multilingual site:

- A `url` value with a leading slash is relative to the `baseURL`.
- A `url` value without a leading slash is relative to the `baseURL` plus the language prefix.

Site type|Front matter `url`|Resulting URL
:--|:--|:--
monolingual|`/about`|`https://example.org/about/`
monolingual|`about`|`https://example.org/about/`
multilingual|`/about`|`https://example.org/about/`
multilingual|`about`|`https://example.org/de/about/`

If you set both `slug` and `url` in front matter, the `url` value takes precedence.

#### Permalinks tokens in front matter

{{< new-in "0.131.0" >}}

You can also use [Permalinks tokens](#tokens) in the `url` front matter. This is typically used in `cascade` sections:

{{< code-toggle file=content/foo/bar/_index.md fm=true >}}
title ="Bar"
[[cascade]]
  url = "/:sections[last]/:slug"
{{< /code-toggle >}}


## Site configuration

### Permalinks

In your site configuration, define a URL pattern for each top-level section. Each URL pattern can target a given language and/or page kind.

Front matter `url` values override the URL patterns defined in the `permalinks` section of your site configuration.

#### Monolingual examples {#permalinks-monolingual-examples}

With this content structure:

```text
content/
├── posts/
│   ├── bash-in-slow-motion.md
│   └── tls-in-a-nutshell.md
├── tutorials/
│   ├── git-for-beginners.md
│   └── javascript-bundling-with-hugo.md
└── _index.md
```

Render tutorials under "training", and render the posts under "articles" with a date-base hierarchy:

{{< code-toggle file=hugo >}}
[permalinks.page]
posts = '/articles/:year/:month/:slug/'
tutorials = '/training/:slug/'
[permalinks.section]
posts = '/articles/'
tutorials = '/training/'
{{< /code-toggle >}}

The structure of the published site will be:

```text
public/
├── articles/
│   ├── 2023/
│   │   ├── 04/
│   │   │   └── bash-in-slow-motion/
│   │   │       └── index.html
│   │   └── 06/
│   │       └── tls-in-a-nutshell/
│   │           └── index.html
│   └── index.html
├── training/
│   ├── git-for-beginners/
│   │   └── index.html
│   ├── javascript-bundling-with-hugo/
│   │   └── index.html
│   └── index.html
└── index.html
```

To create a date-based hierarchy for regular pages in the content root:

{{< code-toggle file=hugo >}}
[permalinks.page]
"/" = "/:year/:month/:slug/"
{{< /code-toggle >}}

Use the same approach with taxonomy terms. For example, to omit the taxonomy segment of the URL:

{{< code-toggle file=hugo >}}
[permalinks.term]
'tags' = '/:slug/'
{{< /code-toggle >}}

#### Multilingual example {#permalinks-multilingual-example}

Use the `permalinks` configuration as a component of your localization strategy.

With this content structure:

```text
content/
├── en/
│   ├── books/
│   │   ├── les-miserables.md
│   │   └── the-hunchback-of-notre-dame.md
│   └── _index.md
└── es/
    ├── books/
    │   ├── les-miserables.md
    │   └── the-hunchback-of-notre-dame.md
    └── _index.md
```

And this site configuration:

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = true

[languages.en]
contentDir = 'content/en'
languageCode = 'en-US'
languageDirection = 'ltr'
languageName = 'English'
weight = 1

[languages.en.permalinks.page]
books = "/books/:slug/"

[languages.en.permalinks.section]
books = "/books/"

[languages.es]
contentDir = 'content/es'
languageCode = 'es-ES'
languageDirection = 'ltr'
languageName = 'Español'
weight = 2

[languages.es.permalinks.page]
books = "/libros/:slug/"

[languages.es.permalinks.section]
books = "/libros/"
{{< /code-toggle >}}

The structure of the published site will be:

```text
public/
├── en/
│   ├── books/
│   │   ├── les-miserables/
│   │   │   └── index.html
│   │   ├── the-hunchback-of-notre-dame/
│   │   │   └── index.html
│   │   └── index.html
│   └── index.html
├── es/
│   ├── libros/
│   │   ├── les-miserables/
│   │   │   └── index.html
│   │   ├── the-hunchback-of-notre-dame/
│   │   │   └── index.html
│   │   └── index.html
│   └── index.html
└── index.html
````

#### Tokens

Use these tokens when defining the URL pattern. These can both be used in the `permalinks` configuration and in the front matter [url](#permalinks-tokens-in-front-matter).

`:filename`

`:year`
: The 4-digit year as defined in the front matter `date` field.

`:month`
: The 2-digit month as defined in the front matter `date` field.

`:monthname`
: The name of the month as defined in the front matter `date` field.

`:day`
: The 2-digit day as defined in the front matter `date` field.

`:weekday`
: The 1-digit day of the week as defined in the front matter `date` field  (Sunday = 0).

`:weekdayname`
: The name of the day of the week as defined in the front matter `date` field.

`:yearday`
: The 1- to 3-digit day of the year as defined in the front matter `date` field.

`:section`
: The content's section.

`:sections`
: The content's sections hierarchy. You can use a selection of the sections using _slice syntax_: `:sections[1:]` includes all but the first, `:sections[:last]` includes all but the last, `:sections[last]` includes only the last, `:sections[1:2]` includes section 2 and 3. Note that this slice access will not throw any out-of-bounds errors, so you don't have to be exact.

`:title`
: The title as defined in front matter, else the automatic title. Hugo generates titles automatically for section, taxonomy, and term pages that are not backed by a file.

`:slug`
: The slug as defined in front matter, else the title as defined in front matter, else the automatic title. Hugo generates titles automatically for section, taxonomy, and term pages that are not backed by a file.

`:filename`
: The content's file name without extension, applicable to the `page` page kind.

`:slugorfilename`
: The slug as defined in front matter, else the content's file name without extension, applicable to the `page` page kind.

For time-related values, you can also use the layout string components defined in Go's [time package]. For example:

[time package]: https://pkg.go.dev/time#pkg-constants

{{< code-toggle file=hugo >}}
permalinks:
  posts: /:06/:1/:2/:title/
{{< /code-toggle >}}

### Appearance

The appearance of a URL is either ugly or pretty.

Type|Path|URL
:--|:--|:--
ugly|content/about.md|`https://example.org/about.html`
pretty|content/about.md|`https://example.org/about/`

By default, Hugo produces pretty URLs. To generate ugly URLs, change your site configuration:

{{< code-toggle file=hugo >}}
uglyURLs = true
{{< /code-toggle >}}

### Post-processing

Hugo provides two mutually exclusive configuration options to alter URLs _after_ it renders a page.

#### Canonical URLs

{{% note %}}
This is a legacy configuration option, superseded by template functions and Markdown render hooks, and will likely be [removed in a future release].

[removed in a future release]: https://github.com/gohugoio/hugo/issues/4733
{{% /note %}}

If enabled, Hugo performs a search and replace _after_ it renders the page. It searches for site-relative URLs (those with a leading slash) associated with `action`, `href`, `src`, `srcset`, and `url` attributes. It then prepends the `baseURL` to create absolute URLs.

```html
<a href="/about"> → <a href="https://example.org/about/">
<img src="/a.gif"> → <img src="https://example.org/a.gif">
```

This is an imperfect, brute force approach that can affect content as well as HTML attributes. As noted above, this is a legacy configuration option that will likely be removed in a future release.

To enable:

{{< code-toggle file=hugo >}}
canonifyURLs = true
{{< /code-toggle >}}

#### Relative URLs

{{% note %}}
Do not enable this option unless you are creating a serverless site, navigable via the file system.
{{% /note %}}

If enabled, Hugo performs a search and replace _after_ it renders the page. It searches for site-relative URLs (those with a leading slash) associated with `action`, `href`, `src`, `srcset`, and `url` attributes. It then transforms the URL to be relative to the current page.

For example, when rendering `content/posts/post-1`:

```html
<a href="/about"> → <a href="../../about">
<img src="/a.gif"> → <img src="../../a.gif">
```

This is an imperfect, brute force approach that can affect content as well as HTML attributes. As noted above, do not enable this option unless you are creating a serverless site.

To enable:

{{< code-toggle file=hugo >}}
relativeURLs = true
{{< /code-toggle >}}

## Aliases

Create redirects from old URLs to new URLs with aliases:

- An alias with a leading slash is relative to the `baseURL`
- An alias without a leading slash is relative to the current directory

### Examples {#alias-examples}

Change the file name of an existing page, and create an alias from the previous URL to the new URL:

{{< code-toggle file=content/posts/new-file-name.md >}}
aliases = ['/posts/previous-file-name']
{{< /code-toggle >}}

Each of these directory-relative aliases is equivalent to the site-relative alias above:

- `previous-file-name`
- `./previous-file-name`
- `../posts/previous-file-name`

You can create more than one alias to the current page:

{{< code-toggle file=content/posts/new-file-name.md >}}
aliases = ['previous-file-name','original-file-name']
{{< /code-toggle >}}

In a multilingual site, use a directory-relative alias, or include the language prefix with a site-relative alias:

{{< code-toggle file=content/posts/new-file-name.de.md >}}
aliases = ['/de/posts/previous-file-name']
{{< /code-toggle >}}

### How aliases work

Using the first example above, Hugo generates the following site structure:

```text
public/
├── posts/
│   ├── new-file-name/
│   │   └── index.html
│   ├── previous-file-name/
│   │   └── index.html
│   └── index.html
└── index.html
```

The alias from the previous URL to the new URL is a client-side redirect:

{{< code file=posts/previous-file-name/index.html >}}
<!DOCTYPE html>
<html lang="en-us">
  <head>
    <title>https://example.org/posts/new-file-name/</title>
    <link rel="canonical" href="https://example.org/posts/new-file-name/">
    <meta name="robots" content="noindex">
    <meta charset="utf-8">
    <meta http-equiv="refresh" content="0; url=https://example.org/posts/new-file-name/">
  </head>
</html>
{{< /code >}}

Collectively, the elements in the `head` section:

- Tell search engines that the new URL is canonical
- Tell search engines not to index the previous URL
- Tell the browser to redirect to the new URL

Hugo renders alias files before rendering pages. A new page with the previous file name will overwrite the alias, as expected.

### Customize

To override Hugo's embedded `alias` template, copy the [source code] to a file with the same name in the layouts directory. The template receives the following context:

Permalink
: The link to the page being aliased.

Page
: The Page data for the page being aliased.

[source code]: {{% eturl alias %}}
