---
title: Config
description: Returns a subset of the site configuration.
categories: []
keywords: []
action:
  related: []
  returnType: page.SiteConfig
  signatures: [SITE.Config]
toc: true
---

The `Config` method on a `Site` object provides access to a subset of the site configuration, specifically the `services` and `privacy` keys.

## Services

These are the default service settings, typically used by Hugo's built-in templates and shortcodes.

{{< code-toggle config=services />}}

For example, to use Hugo's built-in Google Analytics template you must add a [Google tag ID]:

[Google tag ID]: https://support.google.com/tagmanager/answer/12326985?hl=en

{{< code-toggle file=hugo >}}
[services.googleAnalytics]
id = 'G-XXXXXXXXX'
{{< /code-toggle >}}

To access this value from a template:

```go-html-template
{{ .Site.Config.Services.GoogleAnalytics.ID }} → G-XXXXXXXXX
```

You must capitalize each identifier as shown above.

## Privacy

These are the default privacy settings, typically used by Hugo's built-in templates and shortcodes:

{{< code-toggle config=privacy />}}

For example, to disable usage of the built-in YouTube shortcode:

{{< code-toggle file=hugo >}}
[privacy.youtube]
disable = true
{{< /code-toggle >}}

To access this value from a template:

```go-html-template
{{ .Site.Config.Privacy.YouTube.Disable }} → true
```

You must capitalize each identifier as shown above.
