---
title: URL Management
linktitle: URL Management
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [aliases,redirects,permalinks,urls]
categories: [content management]
weight: 110
draft: false
aliases: [/extras/permalinks/,/extras/aliases/,/content-management/permalinks-and-redirects/,/extras/urls/,/doc/redirects/,/doc/alias/,/doc/aliases/]
toc: true
needsreview: true
---

## Permalinks

By default, content is laid out into the target `publishdir` (public)
namespace matching its layout within the `contentdir` hierarchy. The `permalinks` [site configuration][] option allows you to adjust this on a per-section basis. This will change where the files are written to and will change the page's internal "canonical" location, such that template references to `.RelPermalink` will honor the adjustments made as a result of the mappings in this option.

For example, if one of your [sections][] is called `post` and you want to adjust the canonical path to be hierarchical based on the year and month, you could set up the following configurations in YAML and TOML, respectively.

### YAML Permalinks Configuration Example

```yaml
permalinks:
  post: /:year/:month/:title/
```

### TOML Permalinks Configuration Example

```toml
[permalinks]
  post = "/:year/:month/:title/"
```

Only the content under `post/` will have the new URL structure. For example, the file `content/post/sample-entry` with `date: 2013-11-18T19:20:00-05:00` in its front matter will render to `public/2013/11/sample-entry/index.html` at build time and therefore be reachable at `http://yoursite.example.com/2013/11/sample-entry/`.

### Permalink Configuration Values

The following is a list of values that can be used in a `permalink` definition in your site `config` file. All references to time are dependent on the content's date.

* `:year` = the 4-digit year
* `:month` = the 2-digit month
* `:monthname` = the name of the month
* `:day` = the 2-digit day
* `:weekday` = the 1-digit day of the week (Sunday = 0)
* `:weekdayname` = the name of the day of the week
* `:yearday` = the 1- to 3-digit day of the year
* `:section` = the content's section
* `:title` = the content's title
* `:slug` = the content's slug (or title if no slug)
* `:filename` = the content's filename (without extension)

## Aliases

For people migrating existing published content to Hugo, there's a good chance you need a mechanism to handle redirecting old URLs.

Luckily, redirects can be handled easily with _aliases_ in Hugo.

## Example

Given a post on your current Hugo site, with a path of:

``content/posts/my-awesome-blog-post.md``

... you create an "aliases" section in the frontmatter of your post, and add previous paths to that.

### TOML frontmatter

```toml
+++
        ...
aliases = [
    "/posts/my-original-url/",
    "/2010/01/01/even-earlier-url.html"
]
        ...
+++
```

### YAML frontmatter

```yaml
---
        ...
aliases:
    - /posts/my-original-url/
    - /2010/01/01/even-earlier-url.html
        ...
---
```

Now when you visit any of the locations specified in aliases, _assuming the same site domain_, you'll be redirected to the page they are specified on.

## Important Behaviors

1. *Hugo makes no assumptions about aliases. They also don't change based
on your UglyURLs setting. You need to provide absolute path to your webroot
and the complete filename or directory.*

2. *Aliases are rendered prior to any content and will be overwritten by
any content with the same location.*

## Multilingual example

On [multilingual sites](/content-management/multilingual/), each translation of a post can have unique aliases. To use the same alias across multiple languages, prefix it with the language code.

In `/posts/my-new-post.es.md`:

```yaml
---
aliases:
    - /es/posts/my-original-post/
---
```

## How Hugo Aliases Work

When aliases are specified, Hugo creates a physical folder structure to match the alias entry, and, an html file specifying the canonical URL for the page, and a redirect target.

Assuming a baseURL of `mysite.tld`, the contents of the html file will look something like:

```html
<!DOCTYPE html>
<html>
  <head>
    <title>http://mysite.tld/posts/my-original-url</title>
    <link rel="canonical" href="http://mysite.tld/posts/my-original-url"/>
    <meta http-equiv="content-type" content="text/html; charset=utf-8"/>
    <meta http-equiv="refresh" content="0; url=http://mysite.tld/posts/my-original-url"/>
  </head>
</html>
```

The `http-equiv="refresh"` line is what performs the redirect, in 0 seconds in this case.

## Customizing

You may customize this alias page by creating an alias.html template in the
layouts folder of your site.  In this case, the data passed to the template is

* Permalink - the link to the page being aliased
* Page - the Page data for the page being aliased

## Pretty URLs

By default, Hugo renders your content with "pretty" URLs. For example,
content created at `/content/extras/urls.md` will be rendered at
`/public/extras/urls/index.html` according to Hugo's default behavior after running the `hugo` CLI build command. No non-standard server-side
configuration is required for these pretty URLs to work.

## Ugly URLs

If you would like to have what we call "ugly URLs" (e.g.,&nbsp;http://example.com/extras/urls.html), set `uglyurls = true` or `uglyurls: true` to your site-wide `config.toml` or `config.yaml`, respectively. You can also use the `--uglyURLs=true` [flag from the command line][].

If you want a specific piece of content to have an exact URL, you can specify this in the front matter under the `url` key. See [Content Organization][] for more details.

## Canonicalization

By default, all relative URLs encountered in the input are left unmodified, e.g. `/css/foo.css` would stay as `/css/foo.css`, i.e. `canonifyURLs` defaults to `false`.

By setting `canonifyURLs` to `true`, all relative URLs would instead be *canonicalized* using `baseURL`.  For example, assuming you have `baseURL = http://yoursite.example.com/` defined in the site-wide `config.toml`, the relative URL `/css/foo.css` would be turned into the absolute URL `http://yoursite.example.com/css/foo.css`.

Benefits of canonicalization include fixing all URLs to be absolute, which may aid with some parsing tasks.  Note though that all real browsers handle this client-side without issues.

Benefits of non-canonicalization include being able to have resource inclusion be scheme-relative, so that http vs https can be decided based on how this page was retrieved.

{{% note "`canonifyURLs` default change" %}}
In the May 2014 release of Hugo v0.11, the default value of `canonifyURLs` was switched from `true` to `false`, which we think is the better default and should continue to be the case going forward. So, please verify and adjust your website accordingly if you are upgrading from v0.10 or older versions.
{{% /note %}}

To find out the current value of `canonifyURLs` for your website, you may use the handy `hugo config` command added in v0.13.

```bash
hugo config | grep -i canon
```

Or, if you are on Windows and do not have `grep` installed:

```
hugo config | FINDSTR /I canon
```

## Relative URLs

By default, all relative URLs are left unchanged by Hugo, which can be problematic when you want to make your site browsable from a local file system.

Setting `relativeURLs` to `true` in the site configuration will cause Hugo to rewrite all relative URLs to be relative to the current content.

For example, if the `/post/first/` page contained a link with a relative URL of `/about/`, Hugo would rewrite that URL to `../../about/`.

[Content Organization]: /content-management/content-organization/
[flag from the command line]: /getting-started/basic-usage/
[sections]: /content-management/sections/
[site configuration]: /project-organization/configuration/