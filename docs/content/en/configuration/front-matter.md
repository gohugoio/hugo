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
[`Date`][]|Returns the date of the given page.
[`ExpiryDate`][]|Returns the expiry date of the given page.
[`Lastmod`][]|Returns the last modification date of the given page.
[`PublishDate`][]|Returns the publish date of the given page.

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

  Hugo resolves the extracted date to the [`timeZone`][] defined in your project configuration, falling back to the system time zone. Hugo also derives the page [`slug`][] from the remaining file name, unless the page already defines a `slug` in its front matter.

  Slug inference only occurs when `:filename` is the winning date source. If an earlier entry in the list provides a valid date, Hugo skips `:filename` entirely. For example, with `date = ["date", ":filename"]`, a page that defines `date` in its front matter will use that value, and the slug will not be inferred from the file name.

  For example, if you name your file `2025-02-01-article.md`, Hugo will set the date to `2025-02-01` and the slug to `article`.

`:git`
: The Git author date for the file's last revision. To enable access to the Git author date, set [`enableGitInfo`][] to `true`.

## Example

Consider this project configuration:

{{< code-toggle file=hugo >}}
[frontmatter]
date = [':filename', ':default']
publishDate = [':filename', ':default']
lastmod = ['lastmod', ':fileModTime']
{{< /code-toggle >}}

To determine `date` and `publishDate`, Hugo tries to extract the value from the file name, falling back to the default ordered sequence of date fields.

To determine `lastmod`, Hugo looks for a `lastmod` field in front matter, falling back to the file's last modification timestamp.

[`Date`]: /methods/page/date/
[`ExpiryDate`]: /methods/page/expirydate/
[`Lastmod`]: /methods/page/lastmod/
[`PublishDate`]: /methods/page/publishdate/
[`enableGitInfo`]: /configuration/all/#enablegitinfo
[`slug`]: /content-management/front-matter/#slug
[`timeZone`]: /configuration/all/#timezone
