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

disqus.shortname
: (`string`) The `shortname` used with the Disqus commenting system. See&nbsp;[details](/templates/embedded/#disqus). To access this value from a template:

  ```go-html-template
  {{ .Site.Config.Services.Disqus.Shortname }}
  ```

googleAnalytics.id
: (`string`) The Google tag ID for Google Analytics 4 properties. See&nbsp;[details](/templates/embedded/#google-analytics). To access this value from a template:

  ```go-html-template
  {{ .Site.Config.Services.GoogleAnalytics.ID }}
  ```

instagram.accessToken <!-- TODO: Remove when no longer in docs.yaml -->
: (`string`) Do not use. Deprecated in [v0.123.0]. The embedded `instagram` shortcode no longer uses this setting.

instagram.disableInlineCSS <!-- TODO: Remove when no longer in docs.yaml -->
: (`bool`) Do not use. Deprecated in [v0.123.0]. The embedded `instagram` shortcode no longer uses this setting.

rss.limit
: (`int`) The maximum number of items to include in an RSS feed. Set to `-1` for no limit. Default is `-1`. See&nbsp;[details](/templates/rss/). To access this value from a template:

  ```go-html-template
  {{ .Site.Config.Services.RSS.Limit }}
  ```

twitter.disableInlineCSS <!-- TODO: Remove when no longer in docs.yaml -->
: (`bool`) Do not use. Deprecated in [v0.141.0]. Use the `x` shortcode instead.

x.disableInlineCSS
: (`bool`) Whether to disable the inline CSS rendered by the embedded `x` shortode. See&nbsp;[details](/shortcodes/x/#privacy). Default is `false`. To access this value from a template:

  ```go-html-template
  {{ .Site.Config.Services.X.DisableInlineCSS }}

[v0.141.0]: https://github.com/gohugoio/hugo/releases/tag/v0.141.0
[v0.123.0]: https://github.com/gohugoio/hugo/releases/tag/v0.123.0
