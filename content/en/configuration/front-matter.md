---
title: Configure front matter
linkTitle: Front matter
description: Configure front matter.
categories: []
keywords: []
---

## Dates

There are four methods on a `Page` object that return a date.

Method|Description
:--|:--
[`Date`]|Returns the date of the given page.
[`ExpiryDate`]|Returns the expiry date of the given page.
[`Lastmod`]|Returns the last modification date of the given page.
[`PublishDate`]|Returns the publish date of the given page.

[`Date`]: /methods/page/date
[`ExpiryDate`]: /methods/page/expirydate
[`Lastmod`]: /methods/page/lastmod
[`PublishDate`]: /methods/page/publishdate

Hugo determines the values to return based on this configuration:

{{< code-toggle config=frontmatter />}}

The `ExpiryDate` method, for example, returns the `expirydate` value if it exists, otherwise it returns `unpublishdate`.

You can also use custom date parameters:

{{< code-toggle file=hugo >}}
[frontmatter]
date = ["myDate", "date"]
{{< /code-toggle >}}

In the example above, the `Date` method returns the `myDate` value if it exists, otherwise it returns `date`.

To fall back to the default sequence of dates, use the `:default` token:

{{< code-toggle file=hugo >}}
[frontmatter]
date = ["myDate", ":default"]
{{< /code-toggle >}}

In the example above, the `Date` method returns the `myDate` value if it exists, otherwise it returns the first valid date from `date`, `publishdate`, `pubdate`, `published`, `lastmod`, and `modified`.

## Aliases

Some of the front matter fields have aliases.

Front matter field|Aliases
:--|:--
`expiryDate`|`unpublishdate`
`lastmod`|`modified`
`publishDate`|`pubdate`, `published`

The default front matter configuration includes these aliases.

## Tokens

Hugo provides several [tokens](g) to assist with front matter configuration.

Token|Description
:--|:--
`:default`|The default ordered sequence of date fields.
`:fileModTime`|The file's last modification timestamp.
`:filename`|The date from the file name, if present.
`:git`|The Git author date for the file's last revision.

When Hugo extracts a date from a file name, it uses the rest of the file name to generate the page's [`slug`], but only if a slug isn't already specified in the page's front matter. For example, given the name `2025-02-01-article.md`, Hugo will set the `date` to `2025-02-01` and the `slug` to `article`.

[`slug`]: /content-management/front-matter/#slug

To enable access to the Git author date, set [`enableGitInfo`] to `true`, or use\
the `--enableGitInfo` flag when building your site.

[`enableGitInfo`]: /configuration/all/#enablegitinfo

Consider this example:

{{< code-toggle file=hugo >}}
[frontmatter]
date = [':filename', ':default']
lastmod = ['lastmod', ':fileModTime']
{{< /code-toggle >}}

To determine `date`, Hugo tries to extract the date from the file name, falling back to the default ordered sequence of date fields.

To determine `lastmod`, Hugo looks for a `lastmod` field in front matter, falling back to the file's last modification timestamp.
