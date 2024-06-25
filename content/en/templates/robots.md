---
title: robots.txt template
linkTitle: robots.txt templates
description: Hugo can generate a customized robots.txt in the same way as any other template.
categories: [templates]
keywords: []
menu:
  docs:
    parent: templates
    weight: 170
weight: 170
aliases: [/extras/robots-txt/]
---

To generate a robots.txt file from a template, change the [site configuration]:

{{< code-toggle file=hugo >}}
enableRobotsTXT = true
{{< /code-toggle >}}

By default, Hugo generates robots.txt using an [embedded template].

[embedded template]: {{% eturl robots %}}

```text
User-agent: *
```

Search engines that honor the Robots Exclusion Protocol will interpret this as permission to crawl everything on the site.

## robots.txt template lookup order

You may overwrite the internal template with a custom template. Hugo selects the template using this lookup order:

1. `/layouts/robots.txt`
2. `/themes/<THEME>/layouts/robots.txt`

## robots.txt template example

{{< code file=layouts/robots.txt >}}
User-agent: *
{{ range .Pages }}
Disallow: {{ .RelPermalink }}
{{ end }}
{{< /code >}}

This template creates a robots.txt file with a `Disallow` directive for each page on the site. Search engines that honor the Robots Exclusion Protocol will not crawl any page on the site.

{{% note %}}
To create a robots.txt file without using a template:

1. Set `enableRobotsTXT` to `false` in the site configuration.
2. Create a robots.txt file in the `static` directory.

Remember that Hugo copies everything in the [static directory][static] to the root of `publishDir` (typically `public`) when you build your site.

[static]: /getting-started/directory-structure/
{{% /note %}}

[site configuration]: /getting-started/configuration/
