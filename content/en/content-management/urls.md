---
title: URL management
description: Control the structure and appearance of URLs through front matter entries and settings in your site configuration.
categories: []
keywords: []
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

Set the `slug` in front matter to override the last segment of the path. This front matter field is not applicable to `home`, `section`, `taxonomy`, or `term` pages.

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

> [!note]
> Hugo does not sanitize the `url` front matter field, allowing you to generate:
>
> - File paths that contain characters reserved by the operating system. For example, file paths on Windows may not contain any of these [reserved characters]. Hugo throws an error if a file path includes a character reserved by the current operating system.
> - URLs that contain disallowed characters. For example, the less than sign (`<`) is not allowed in a URL.

If you set both `slug` and `url` in front matter, the `url` value takes precedence.

#### Include a colon

{{< new-in 0.136.0 />}}

If you need to include a colon in the  `url` front matter field, escape it with backslash characters. Use one backslash if you wrap the string within single quotes, or use two backslashes if you wrap the string within double quotes. With YAML front matter, use a single backslash if you omit quotation marks.

For example, with this front matter:

{{< code-toggle file=content/example.md fm=true >}}
title: Example
url: "my\\:example"
{{< /code-toggle >}}

The resulting URL will be:

```text
https://example.org/my:example/
```

As described above, this will fail on Windows because the colon (`:`) is a reserved character.

#### File extensions

With this front matter:

{{< code-toggle file=content/posts/post-1.md fm=true >}}
title = 'My First Article'
url = 'articles/my-first-article'
{{< /code-toggle >}}

The resulting URL will be:

```text
https://example.org/articles/my-first-article/
```

If you include a file extension:

{{< code-toggle file=content/posts/post-1.md fm=true >}}
title = 'My First Article'
url = 'articles/my-first-article.html'
{{< /code-toggle >}}

The resulting URL will be:

```text
https://example.org/articles/my-first-article.html
```

#### Leading slashes

With monolingual sites, `url` values with or without a leading slash are relative to the [`baseURL`]. With multilingual sites, `url` values with a leading slash are relative to the `baseURL`, and  `url` values without a leading slash are relative to the `baseURL` plus the language prefix.

Site type|Front matter `url`|Resulting URL
:--|:--|:--
monolingual|`/about`|`https://example.org/about/`
monolingual|`about`|`https://example.org/about/`
multilingual|`/about`|`https://example.org/about/`
multilingual|`about`|`https://example.org/de/about/`

#### Permalinks tokens in front matter

{{< new-in 0.131.0 />}}

You can also use tokens when setting the `url` value. This is typically used in `cascade` sections:

{{< code-toggle file=content/foo/bar/_index.md fm=true >}}
title ="Bar"
[[cascade]]
  url = "/:sections[last]/:slug"
{{< /code-toggle >}}

Use any of these tokens:

{{% include "/_common/permalink-tokens.md" %}}

## Site configuration

### Permalinks

See [configure permalinks](/configuration/permalinks).

### Appearance

See [configure ugly URLs](/configuration/ugly-urls/).

### Post-processing

Hugo provides two mutually exclusive configuration options to alter URLs _after_ it renders a page.

#### Canonical URLs

> [!caution]
> This is a legacy configuration option, superseded by template functions and Markdown render hooks, and will likely be [removed in a future release].
{class="!mt-6"}

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

> [!caution]
> Do not enable this option unless you are creating a serverless site, navigable via the file system.
{class="!mt-6"}

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

{{< code-toggle file=content/posts/new-file-name.md fm=true >}}
aliases = ['/posts/previous-file-name']
{{< /code-toggle >}}

Each of these directory-relative aliases is equivalent to the site-relative alias above:

- `previous-file-name`
- `./previous-file-name`
- `../posts/previous-file-name`

You can create more than one alias to the current page:

{{< code-toggle file=content/posts/new-file-name.md fm=true >}}
aliases = ['previous-file-name','original-file-name']
{{< /code-toggle >}}

In a multilingual site, use a directory-relative alias, or include the language prefix with a site-relative alias:

{{< code-toggle file=content/posts/new-file-name.de.md fm=true >}}
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

```html {file="posts/previous-file-name/index.html"}
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
```

Collectively, the elements in the `head` section:

- Tell search engines that the new URL is canonical
- Tell search engines not to index the previous URL
- Tell the browser to redirect to the new URL

Hugo renders alias files before rendering pages. A new page with the previous file name will overwrite the alias, as expected.

### Customize

To override Hugo's embedded `alias` template, copy the [source code] to a file with the same name in the `layouts` directory. The template receives the following context:

Permalink
: The link to the page being aliased.

Page
: The Page data for the page being aliased.

[`baseURL`]: /configuration/all/#baseurl
[removed in a future release]: https://github.com/gohugoio/hugo/issues/4733
[reserved characters]: https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-file#naming-conventions
[source code]: <{{% eturl alias %}}>
