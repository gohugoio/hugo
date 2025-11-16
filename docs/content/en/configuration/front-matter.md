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

Hugo provides the following [tokens](g) to help you configure your front matter:

`:default`
: The default ordered sequence of date fields.

`:fileModTime`
: The file's last modification timestamp.

`:filename`
: Extracts the date from the file name, provided the file name begins with a date in one of the following formats:

  - `YYYY-MM-DD`
  - `YYYY-MM-DD-HH-MM-SS` {{< new-in 0.148.0 />}}

  Within the `YYYY-MM-DD-HH-MM-SS` format, the date and time values may be separated by any character including a space (e.g., `2025-02-01T14-30-00`).

  Hugo resolves the extracted date to the [`timeZone`] defined in your site configuration, falling back to the system time zone. After extracting the date, Hugo uses the remaining part of the file name to generate the page's [`slug`], but only if you haven't already specified a slug in the page's front matter.

  For example, if you name your file `2025-02-01-article.md`, Hugo will set the date to `2025-02-01` and the slug to `article`.

`:git`
: The Git author date for the file's last revision. To enable access to the Git author date, set [`enableGitInfo`] to `true`, or use the `--enableGitInfo` flag when building your site.

## Example

Consider this site configuration:

{{< code-toggle file=hugo >}}
[frontmatter]
date = [':filename', ':default']
lastmod = ['lastmod', ':fileModTime']
{{< /code-toggle >}}

To determine `date`, Hugo tries to extract the date from the file name, falling back to the default ordered sequence of date fields.

To determine `lastmod`, Hugo looks for a `lastmod` field in front matter, falling back to the file's last modification timestamp.

[`enableGitInfo`]: /configuration/all/#enablegitinfo
[`slug`]: /content-management/front-matter/#slug
[`timeZone`]: /configuration/all/#timezone
