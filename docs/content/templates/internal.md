---
title: Internal Templates
linktitle: Internal Templates
description: Hugo ships with a group of boilerplate templates that cover the most common use cases for static websites.
date: 2017-03-06
publishdate: 2017-03-06
lastmod: 2017-03-06
categories: [templates]
#tags: [internal, analytics,]
menu:
  docs:
    parent: "templates"
    weight: 168
weight: 168
sections_weight: 168
draft: false
aliases: []
toc: true
wip: true
---
<!-- reference: https://discourse.gohugo.io/t/lookup-order-for-partials/5705/6
code: https://github.com/gohugoio/hugo/blob/e445c35d6a0c7f5fc2f90f31226cd1d46e048bbc/tpl/template_embedded.go#L147 -->

{{% warning %}}
While the following internal templates are called similar to partials, they do *not* observe the partial template lookup order.
{{% /warning %}}

## Google Analytics

Hugo ships with internal templates for Google Analytics tracking, including both synchronous and asynchronous tracking codes.

### Configure Google Analytics

Provide your tracking id in your configuration file:

```toml
googleAnalytics = "UA-123-45"
```

```yml
googleAnalytics: "UA-123-45"
```

### Use the Google Analytics Template

You can then include the Google Analytics internal template:

```golang
{{ template "_internal/google_analytics.html" . }}
```


```golang
{{ template "_internal/google_analytics_async.html" . }}
```

## Disqus

Hugo also ships with an internal template for Disqus, a popular commenting system for both static and dynamic websites. In order to use Disqus, you will need to secure a Disqus "shortname" by [signing up for the free service](https://disqus.com/profile/signup/).

### Configure Disqus

To use Hugo's Disqus template, you first need to set a single value in your site's `config.toml` or `config.yml`:

```toml
disqusShortname = "yourdiscussshortname"
```

```yaml
disqusShortname: "yourdiscussshortname"
```

You also have the option to set the following in the front matter for a given piece of content:

* `disqus_identifier`
* `disqus_title`
* `disqus_url`

### Use the Disqus Template

To add Disqus, include the following line in templates where you want your comments to appear:

```golang
{{ template "_internal/disqus.html" . }}
```

Be aware that this template will not load Disqus when you are previewing your website locally. When running on `localhost` (i.e. via `hugo server`), initialization of Disqus is skipped to avoid creating unwanted discussions on your Disqus account.
