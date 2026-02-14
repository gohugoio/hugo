---
title: Configure sitemap
linkTitle: Sitemap
description: Configure the sitemap.
categories: []
keywords: []
---

These are the default sitemap configuration values. They apply to all pages unless overridden in front matter.

{{< code-toggle config=sitemap />}}

changefreq
: (`string`) How frequently a page is likely to change. Valid values are `always`, `hourly`, `daily`, `weekly`, `monthly`, `yearly`, and `never`. With the default value of `""` Hugo will omit this field from the sitemap. See&nbsp;[details](https://www.sitemaps.org/protocol.html#changefreqdef).

disable
: (`bool`) Whether to disable page inclusion. Default is `false`. Set to `true` in front matter to exclude the page.

filename
: (`string`) The name of the generated file. Default is `sitemap.xml`.

priority
: (`float`) The priority of a page relative to any other page on the site. Valid values range from 0.0 to 1.0. With the default value of `-1` Hugo will omit this field from the sitemap. See&nbsp;[details](https://www.sitemaps.org/protocol.html#prioritydef).
