---
title: URL management
description: Control the structure and appearance of URLs through front matter entries and settings in your project configuration.
categories: []
keywords: []
aliases: [/extras/permalinks/,/extras/aliases/,/extras/urls/,/doc/redirects/,/doc/alias/,/doc/aliases/]
---

## Overview

By default, when Hugo renders a page, the resulting URL matches the file path within the `content` directory. For example:

```text
content/posts/post-1.md → https://example.org/posts/post-1/
```

You can change the structure and appearance of URLs with front matter values and project configuration options.

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
> - File paths that contain characters reserved by the operating system. For example, file paths on Windows may not contain any of these [reserved characters][]. Hugo throws an error if a file path includes a character reserved by the current operating system.
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

With monolingual sites, `url` values with or without a leading slash are relative to the [`baseURL`][]. With multilingual sites, `url` values with a leading slash are relative to the `baseURL`, and  `url` values without a leading slash are relative to the `baseURL` plus the language prefix.

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

## Project configuration

### Permalinks

See [configure permalinks](/configuration/permalinks).

### Appearance

See [configure ugly URLs](/configuration/ugly-urls/).

### Post-processing

Hugo provides two mutually exclusive configuration options to alter URLs _after_ it renders a page.

#### Canonical URLs

> [!caution]
> This is a legacy configuration option, superseded by template functions and Markdown render hooks, and will likely be [removed in a future release][].
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

Aliases allow you to redirect old URLs to new URLs. This is essential for preventing broken links and ensuring that existing bookmarks or external links continue to function when you rename or move content.

### Defining aliases

To add redirects to a page, list the previous paths in the [`aliases`][aliases_field] field in your front matter. Hugo resolves these to [server-relative](g) paths during the build process, accounting for the [`baseURL`][] and [content dimension](g) prefixes such as language, version, or role.

{{< code-toggle file=content/examples/example-1.en.md fm=true >}}
title = 'Example 1'
date = 2025-02-02
aliases = ['/old-url', 'old-name', '../old/path']
{{< /code-toggle >}}

As shown in the example above, you can use [site-relative](g) paths or [page-relative](g) paths. Page-relative paths can also include directory traversal. Using the file `content/examples/example-1.en.md` as a reference point, here is how Hugo interprets those different path types:

Path type|Alias|Server-relative path
:--|:--|:--
site-relative|`/old-url`|`/en/old-url/`
page-relative|`old-name`|`/en/examples/old-name/`
page-relative|`../old/path`|`/en/old/path/`

### Redirection methods

There are two ways to implement aliases depending on your hosting environment and preferences: client-side redirection and server-side redirection.

> [!note]
> Alias data is only generated for [output formats](g) where both [`isHTML`][] and [`permalinkable`][] are `true`. This affects both the creation of client-side redirect files and the results returned by the [`Aliases`][aliases_method] method used in server-side redirection.

#### Client-side redirection

By default, Hugo uses client-side redirection, generating a small HTML file for every alias. This file contains a `meta http-equiv="refresh"` tag that instructs the browser to navigate to the new URL. This approach is portable across all hosting providers.

When using this method, Hugo creates a physical directory and an `index.html` file at each alias location. For example, if a page at `content/posts/new.md` has a page-relative alias of `old-path`, a file is generated at `public/posts/old-path/index.html`.

Unless you provide a custom layout, Hugo uses its [embedded alias template][] to generate the redirect files:

```go-html-template
<!DOCTYPE html>
<html lang="{{ site.Language.LanguageCode }}">
  <head>
    <title>{{ .Permalink }}</title>
    {{ with .OutputFormats.Canonical }}<link rel="{{ .Rel }}" href="{{ .Permalink }}">{{ end }}
    <meta charset="utf-8">
    <meta http-equiv="refresh" content="0; url={{ .Permalink }}">
  </head>
</html>
```

To override this, create a file named `alias.html` in your `layouts` directory. This custom template has access to the following context:

`Permalink`
: (`string`) The absolute URL of the destination page.

`Page`
: (`page.Page`) The full `Page` object of the destination.

#### Server-side redirection

Alternatively, you can implement server-side redirection by using the [`Aliases`][aliases_method] method on a `Page` object to generate a single configuration file that the web server processes. This method is more efficient because the redirect happens at the HTTP header level before any page content is processed, whereas a meta refresh requires the browser to download and parse the HTML body before acting. Additionally, server-side redirection improves build and deployment times because Hugo doesn't need to write a physical directory and HTML file for every alias.

To implement this, you typically create a single template to generate the necessary rules for your specific host or server. Common examples include:

- A `_redirects` file for hosting services such as Cloudflare, GitLab Pages, and Netlify.
- An `.htaccess` file for web servers such as Apache and LiteSpeed.

See the [`Aliases`][aliases_method] method page for a complete example of how to iterate through pages to generate these rules.

If you implement server-side redirects, you should disable the generation of individual HTML files by setting [`disableAliases`][] to `true` in your project configuration. This setting only prevents the generation of the physical HTML files; the `Aliases` method on a `Page` object remains available for use in your configuration templates.

[`baseURL`]: /configuration/all/#baseurl
[`disableAliases`]: /configuration/all/#disablealiases
[`isHTML`]: /configuration/output-formats/#ishtml
[`permalinkable`]: /configuration/output-formats/#permalinkable
[aliases_field]: /content-management/front-matter/#aliases
[aliases_method]: /methods/page/aliases/
[embedded alias template]: <{{% eturl alias %}}>
[removed in a future release]: https://github.com/gohugoio/hugo/issues/4733
[reserved characters]: https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-file#naming-conventions
