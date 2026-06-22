---
title: Configure services
linkTitle: Services
description: Configure embedded templates.
categories: []
keywords: []
---

Hugo provides [embedded templates](g) to simplify site and content creation. Some of these templates are configurable. For example, the embedded Google Analytics template requires a Google tag ID.

This is the default configuration:

{{< code-toggle config=services />}}

`disqus.shortname`
: (`string`) The `shortname` used with the Disqus commenting system. See [details][disqus]. To access this value from a template:

  ```go-html-template
  {{ .Site.Config.Services.Disqus.Shortname }}
  ```

`googleAnalytics.id`
: (`string`) The Google tag ID for Google Analytics 4 properties. See [details][google-analytics]. To access this value from a template:

  ```go-html-template
  {{ .Site.Config.Services.GoogleAnalytics.ID }}
  ```

`rss.limit`
: (`int`) The maximum number of items to include in an RSS feed. Set to `-1` for no limit. Default is `-1`. See [details][rss]. To access this value from a template:

  ```go-html-template
  {{ .Site.Config.Services.RSS.Limit }}
  ```

`x.disableInlineCSS`
: (`bool`) Whether to disable the inline CSS rendered by the embedded `x` shortode. See [details][privacy]. Default is `false`. To access this value from a template:

  ```go-html-template
  {{ .Site.Config.Services.X.DisableInlineCSS }}
  ```

[disqus]: /templates/embedded/#disqus
[google-analytics]: /templates/embedded/#google-analytics
[privacy]: /shortcodes/x/#privacy
[rss]: /templates/rss/
