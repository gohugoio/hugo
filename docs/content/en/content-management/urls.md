---
title: URL Management
linktitle: URL Management
description: Hugo supports permalinks, aliases, link canonicalization, and multiple options for handling relative vs absolute URLs.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-03-09
keywords: [aliases,redirects,permalinks,urls]
categories: [content management]
menu:
  docs:
    parent: "content-management"
    weight: 110
weight: 110	#rem
draft: false
aliases: [/extras/permalinks/,/extras/aliases/,/extras/urls/,/doc/redirects/,/doc/alias/,/doc/aliases/]
toc: true
---

## Permalinks

The default Hugo target directory for your built website is `public/`. However, you can change this value by specifying a different `publishDir` in your [site configuration][config]. The directories created at build time for a section reflect the position of the content's directory within the `content` folder and namespace matching its layout within the `contentdir` hierarchy.

The `permalinks` option in your [site configuration][config] allows you to adjust the directory paths (i.e., the URLs) on a per-section basis. This will change where the files are written to and will change the page's internal "canonical" location, such that template references to `.RelPermalink` will honor the adjustments made as a result of the mappings in this option.

{{% note "Default Publish and Content Folders" %}}
These examples use the default values for `publishDir` and `contentDir`; i.e., `public` and `content`, respectively. You can override the default values in your [site's `config` file](/getting-started/configuration/).
{{% /note %}}

For example, if one of your [sections][] is called `posts` and you want to adjust the canonical path to be hierarchical based on the year, month, and post title, you could set up the following configurations in YAML and TOML, respectively.

### Permalinks Configuration Example

{{< code-toggle file="config" copy="false" >}}
permalinks:
  posts: /:year/:month/:title/
{{< /code-toggle >}}

Only the content under `posts/` will have the new URL structure. For example, the file `content/posts/sample-entry.md` with `date: 2017-02-27T19:20:00-05:00` in its front matter will render to `public/2017/02/sample-entry/index.html` at build time and therefore be reachable at `https://example.com/2017/02/sample-entry/`.

If the standard date-based permalink configuration does not meet your needs, you can also format URL segments using [Go time formatting directives](https://golang.org/pkg/time/#Time.Format). For example, a URL structure with two digit years and month and day digits without zero padding can be accomplished with:

{{< code-toggle file="config" copy="false" >}}
permalinks:
  posts: /:06/:1/:2/:title/
{{< /code-toggle >}}

You can also configure permalinks of taxonomies with the same syntax, by using the plural form of the taxonomy instead of the section. You will probably only want to use the configuration values `:slug` or `:title`.

### Permalink Configuration Values

The following is a list of values that can be used in a `permalink` definition in your site `config` file. All references to time are dependent on the content's date.

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
: the content's sections hierarchy

`:title`
: the content's title

`:slug`
: the content's slug (or title if no slug is provided in the front matter)

`:filename`
: the content's filename (without extension)

Additionally, a Go time format string prefixed with `:` may be used.

## Aliases

Aliases can be used to create redirects to your page from other URLs.

Aliases comes in two forms:

1. Starting with a `/` meaning they are relative to the `BaseURL`, e.g. `/posts/my-blogpost/`
2. They are relative to the `Page` they're defined in, e.g. `my-blogpost` or even something like `../blog/my-blogpost` (new in Hugo 0.55).

### Example: Aliases

Let's assume you create a new piece of content at `content/posts/my-awesome-blog-post.md`. The content is a revision of your previous post at `content/posts/my-original-url.md`. You can create an `aliases` field in the front matter of your new `my-awesome-blog-post.md` where you can add previous paths. The following examples show how to create this field in TOML and YAML front matter, respectively.

#### TOML Front Matter

{{< code file="content/posts/my-awesome-post.md" copy="false" >}}
+++
aliases = [
    "/posts/my-original-url/",
    "/2010/01/01/even-earlier-url.html"
]
+++
{{< /code >}}

#### YAML Front Matter

{{< code file="content/posts/my-awesome-post.md" copy="false" >}}
---
aliases:
    - /posts/my-original-url/
    - /2010/01/01/even-earlier-url.html
---
{{< /code >}}

Now when you visit any of the locations specified in aliases---i.e., *assuming the same site domain*---you'll be redirected to the page they are specified on. For example, a visitor to `example.com/posts/my-original-url/` will be immediately redirected to `example.com/posts/my-awesome-post/`.

### Example: Aliases in Multilingual

On [multilingual sites][multilingual], each translation of a post can have unique aliases. To use the same alias across multiple languages, prefix it with the language code.

In `/posts/my-new-post.es.md`:

```
---
aliases:
    - /es/posts/my-original-post/
---
```

From Hugo 0.55 you can also have page-relative aliases, so ` /es/posts/my-original-post/` can be simplified to the more portable `my-original-post/`

### How Hugo Aliases Work

When aliases are specified, Hugo creates a directory to match the alias entry. Inside the directory, Hugo creates an `.html` file specifying the canonical URL for the page and the new redirect target.

For example, a content file at `posts/my-intended-url.md` with the following in the front matter:

```
---
title: My New post
aliases: [/posts/my-old-url/]
---
```

Assuming a `baseURL` of `example.com`, the contents of the auto-generated alias `.html` found at `https://example.com/posts/my-old-url/` will contain the following:

```
<!DOCTYPE html>
<html>
  <head>
    <title>https://example.com/posts/my-intended-url</title>
    <link rel="canonical" href="https://example.com/posts/my-intended-url"/>
    <meta name="robots" content="noindex">
    <meta http-equiv="content-type" content="text/html; charset=utf-8"/>
    <meta http-equiv="refresh" content="0; url=https://example.com/posts/my-intended-url"/>
  </head>
</html>
```

The `http-equiv="refresh"` line is what performs the redirect, in 0 seconds in this case. If an end user of your website goes to `https://example.com/posts/my-old-url`, they will now be automatically redirected to the newer, correct URL. The addition of `<meta name="robots" content="noindex">` lets search engine bots know that they should not crawl and index your new alias page.

### Customize
You may customize this alias page by creating an `alias.html` template in the
layouts folder of your site (i.e., `layouts/alias.html`). In this case, the data passed to the template is

`Permalink`
: the link to the page being aliased

`Page`
: the Page data for the page being aliased

### Important Behaviors of Aliases

1. Hugo makes no assumptions about aliases. They also do not change based
on your UglyURLs setting. You need to provide absolute paths to your web root
and the complete filename or directory.
2. Aliases are rendered *before* any content are rendered and therefore will be overwritten by any content with the same location.

## Pretty URLs

Hugo's default behavior is to render your content with "pretty" URLs. No non-standard server-side configuration is required for these pretty URLs to work.

The following demonstrates the concept:

```
content/posts/_index.md
=> example.com/posts/index.html
content/posts/post-1.md
=> example.com/posts/post-1/
```

## Ugly URLs

If you would like to have what are often referred to as "ugly URLs" (e.g., example.com/urls.html), set `uglyurls = true` or `uglyurls: true` in your site's `config.toml` or `config.yaml`, respectively. You can also use the `--uglyURLs=true` [flag from the command line][usage] with `hugo` or `hugo server`.

If you want a specific piece of content to have an exact URL, you can specify this in the [front matter][] under the `url` key. The following are examples of the same content directory and what the eventual URL structure will be when Hugo runs with its default behavior.

See [Content Organization][contentorg] for more details on paths.

```
.
└── content
    └── about
    |   └── _index.md  // <- https://example.com/about/
    ├── posts
    |   ├── firstpost.md   // <- https://example.com/posts/firstpost/
    |   ├── happy
    |   |   └── ness.md  // <- https://example.com/posts/happy/ness/
    |   └── secondpost.md  // <- https://example.com/posts/secondpost/
    └── quote
        ├── first.md       // <- https://example.com/quote/first/
        └── second.md      // <- https://example.com/quote/second/
```

Here's the same organization run with `hugo --uglyURLs`:

```
.
└── content
    └── about
    |   └── _index.md  // <- https://example.com/about.html
    ├── posts
    |   ├── firstpost.md   // <- https://example.com/posts/firstpost.html
    |   ├── happy
    |   |   └── ness.md    // <- https://example.com/posts/happy/ness.html
    |   └── secondpost.md  // <- https://example.com/posts/secondpost.html
    └── quote
        ├── first.md       // <- https://example.com/quote/first.html
        └── second.md      // <- https://example.com/quote/second.html
```


## Canonicalization

By default, all relative URLs encountered in the input are left unmodified, e.g. `/css/foo.css` would stay as `/css/foo.css`. The `canonifyURLs` field in your site `config` has a default value of `false`.

By setting `canonifyURLs` to `true`, all relative URLs would instead be *canonicalized* using `baseURL`.  For example, assuming you have `baseURL = https://example.com/`, the relative URL `/css/foo.css` would be turned into the absolute URL `https://example.com/css/foo.css`.

Benefits of canonicalization include fixing all URLs to be absolute, which may aid with some parsing tasks. Note, however, that all modern browsers handle this on the client without issue.

Benefits of non-canonicalization include being able to have scheme-relative resource inclusion; e.g., so that `http` vs `https` can be decided according to how the page was retrieved.

{{% note "`canonifyURLs` default change" %}}
In the May 2014 release of Hugo v0.11, the default value of `canonifyURLs` was switched from `true` to `false`, which we think is the better default and should continue to be the case going forward. Please verify and adjust your website accordingly if you are upgrading from v0.10 or older versions.
{{% /note %}}

To find out the current value of `canonifyURLs` for your website, you may use the handy `hugo config` command added in v0.13.

```
hugo config | grep -i canon
```

Or, if you are on Windows and do not have `grep` installed:

```
hugo config | FINDSTR /I canon
```

## Set URL in Front Matter

In addition to specifying permalink values in your site configuration for different content sections, Hugo provides even more granular control for individual pieces of content.

Both `slug` and `url` can be defined in individual front matter. For more information on content destinations at build time, see [Content Organization][contentorg].

From Hugo 0.55, you can use URLs relative to the current site context (the language), which makes it simpler to maintain. For a Japanese translation, both of the following examples would get the same URL:

```markdown
---
title: "Custom URL!"
url: "/jp/custom/foo"
---
```

```markdown
---
title: "Custom URL!"
url: "custom/foo"
---
```


## Relative URLs

By default, all relative URLs are left unchanged by Hugo, which can be problematic when you want to make your site browsable from a local file system.

Setting `relativeURLs` to `true` in your [site configuration][config] will cause Hugo to rewrite all relative URLs to be relative to the current content.

For example, if your `/posts/first/` page contains a link to `/about/`, Hugo will rewrite the URL to `../../about/`.

[config]: /getting-started/configuration/
[contentorg]: /content-management/organization/
[front matter]: /content-management/front-matter/
[multilingual]: /content-management/multilingual/
[sections]: /content-management/sections/
[usage]: /getting-started/usage/
