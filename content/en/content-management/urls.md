---
title: URL Management
description: Control the structure and appearance of URLs through front matter entries and settings in your site configuration.
categories: [content management]
keywords: [aliases,redirects,permalinks,urls]
menu:
  docs:
    parent: content-management
    weight: 180
toc: true
weight: 180
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

{{< code-toggle file="content/posts/post-1.md" copy=false fm=true >}}
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

{{< code-toggle file="content/posts/post-1.md" copy=false fm=true >}}
title = 'My First Article'
url = '/articles/my-first-article'
{{< /code-toggle >}}

The resulting URL will be:

```text
https://example.org/articles/my-first-article/
```

If you include a file extension:

{{< code-toggle file="content/posts/post-1.md" copy=false fm=true >}}
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

## Site configuration

### Permalinks

In your site configuration, set a URL pattern for regular pages within a top-level section. This is recursive, affecting descendant regular pages.

{{% note %}}
The  `permalinks` defined in your site configuration do not apply to section pages. To adjust the URL for section pages, set `url` in front matter.
{{% /note %}}

#### Examples {#permalinks-examples}

With this content structure:

```text
content/
├── posts/
│   ├── _index.md
│   ├── post-1.md
│   └── post-2.md
└── _index.md
```

Create a date-based hierarchy, recursively, for regular pages within the `posts` section:

{{< code-toggle file="config" copy=false >}}
[permalinks]
  posts = '/posts/:year/:month/:title/'
{{< /code-toggle >}}

The structure of the published site will be:

```text
public/
├── posts/
│   ├── 2023/
│   │   └── 03/
│   │       ├── post-1/
│   │       │   └── index.html
│   │       └── post-2/
│   │           └── index.html
│   └── index.html
├── favicon.ico
└── index.html
```

To create a date-based hierarchy for regular pages in the content root:

{{< code-toggle file="config" copy=false >}}
[permalinks]
  '/' = '/:year/:month/:title/'
{{< /code-toggle >}}

{{% note %}}
A URL pattern defined for the content root is not recursive.
{{% /note %}}

Use the same approach with taxonomies. For example, to omit the taxonomy segment of the URL:

{{< code-toggle file="config" copy=false >}}
[permalinks]
  'tags' = '/:title/'
{{< /code-toggle >}}

Front matter `url` values take precedence over URL patterns defined in `permalinks`.

#### Tokens

Use these tokens when defining the URL pattern. The `date` field in front matter determines the value of time-related tokens.

`:year`
: the 4-digit year

`:month`
: the 2-digit month

`:monthname`
: the name of the month

`:day`
: the 2-digit day

`:weekday`
: the 1-digit day of the week (Sunday = 0)

`:weekdayname`
: the name of the day of the week

`:yearday`
: the 1- to 3-digit day of the year

`:section`
: the content's section

`:sections`
: the content's sections hierarchy. You can use a selection of the sections using _slice syntax_: `:sections[1:]` includes all but the first, `:sections[:last]` includes all but the last, `:sections[last]` includes only the last, `:sections[1:2]` includes section 2 and 3. Note that this slice access will not throw any out-of-bounds errors, so you don't have to be exact.

`:title`
: the content's title

`:slug`
: the content's slug (or title if no slug is provided in the front matter)

`:slugorfilename`
: the content's slug (or filename if no slug is provided in the front matter)

`:filename`
: the content's filename (without extension)

For time-related values, you can also use the layout string components defined in Go's [time package]. For example:

[time package]: https://pkg.go.dev/time#pkg-constants

{{< code-toggle file="config" copy=false >}}
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

{{< code-toggle file="config" copy=false >}}
uglyURLs = true
{{< /code-toggle >}}

### Post-processing

Hugo provides two mutually exclusive configuration options to alter URLs _after_ it renders a page.

#### Canonical URLs

{{% note %}}
This is a legacy configuration option, superseded by template functions and markdown render hooks, and will likely be [removed in a future release].

[removed in a future release]: https://github.com/gohugoio/hugo/issues/4733
{{% /note %}}

If enabled, Hugo performs a search and replace _after_ it renders the page. It searches for site-relative URLs (those with a leading slash) associated with `action`, `href`, `src`, `srcset`, and `url` attributes. It then prepends the `baseURL` to create absolute URLs.

```text
<a href="/about"> → <a href="https://example.org/about/">
<img src="/a.gif"> → <img src="https://example.org/a.gif">
```

This is an imperfect, brute force approach that can affect content as well as HTML attributes. As noted above, this is a legacy configuration option that will likely be removed in a future release.

To enable:

{{< code-toggle file="config" copy=false >}}
canonifyURLs = true
{{< /code-toggle >}}

#### Relative URLs

{{% note %}}
Do not enable this option unless you are creating a serverless site, navigable via the file system.
{{% /note %}}

If enabled, Hugo performs a search and replace _after_ it renders the page. It searches for site-relative URLs (those with a leading slash) associated with `action`, `href`, `src`, `srcset`, and `url` attributes. It then transforms the URL to be relative to the current page.

For example, when rendering `content/posts/post-1`:

```text
<a href="/about"> → <a href="../../about">
<img src="/a.gif"> → <img src="../../a.gif">
```

This is an imperfect, brute force approach that can affect content as well as HTML attributes. As noted above, do not enable this option unless you are creating a serverless site.

To enable:

{{< code-toggle file="config" copy=false >}}
relativeURLs = true
{{< /code-toggle >}}

## Aliases

Create redirects from old URLs to new URLs with aliases:

- An alias with a leading slash is relative to the `baseURL`
- An alias without a leading slash is relative to the current directory

### Examples {#alias-examples}

Change the file name of an existing page, and create an alias from the previous URL to the new URL:

{{< code-toggle file="content/posts/new-file-name.md" copy=false >}}
aliases = ['/posts/previous-file-name']
{{< /code-toggle >}}

Each of these directory-relative aliases is equivalent to the site-relative alias above:

- `previous-file-name`
- `./previous-file-name`
- `../posts/previous-file-name`

You can create more than one alias to the current page:

{{< code-toggle file="content/posts/new-file-name.md" copy=false >}}
aliases = ['previous-file-name','original-file-name']
{{< /code-toggle >}}

In a multilingual site, use a directory-relative alias, or include the language prefix with a site-relative alias:

{{< code-toggle file="content/posts/new-file-name.de.md" copy=false >}}
aliases = ['/de/posts/previous-file-name']
{{< /code-toggle >}}

### How Aliases Work

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

{{< code file="posts/previous-file-name/index.html" copy=false >}}
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

Create a new template (`layouts/alias.html`) to customize the content of the alias files. The template receives the following context:

`Permalink`
: the link to the page being aliased

`Page`
: the Page data for the page being aliased
