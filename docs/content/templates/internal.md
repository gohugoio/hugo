---
title: Internal Templates
linktitle: Internal Templates
description: Hugo ships with a group of boilerplate templates that cover the most common use cases for static websites.
date: 2017-03-06
publishdate: 2017-03-06
lastmod: 2017-03-06
categories: [templates]
keywords: [internal, analytics,]
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

```
googleAnalytics = "UA-123-45"
```

```
googleAnalytics: "UA-123-45"
```

### Use the Google Analytics Template

You can then include the Google Analytics internal template:

```
{{ template "_internal/google_analytics.html" . }}
```


```
{{ template "_internal/google_analytics_async.html" . }}
```

## Disqus

Hugo also ships with an internal template for [Disqus comments][disqus], a popular commenting system for both static and dynamic websites. In order to effectively use Disqus, you will need to secure a Disqus "shortname" by [signing up for the free service][disqussignup].

### Configure Disqus

To use Hugo's Disqus template, you first need to set a single value in your site's `config.toml` or `config.yml`:

```
disqusShortname = "yourdiscussshortname"
```

```
disqusShortname: "yourdiscussshortname"
```

You also have the option to set the following in the front matter for a given piece of content:

* `disqus_identifier`
* `disqus_title`
* `disqus_url`

### Use the Disqus Template

To add Disqus, include the following line in templates where you want your comments to appear:

```
{{ template "_internal/disqus.html" . }}
```

### Conditional Loading of Disqus Comments

Users have noticed that enabling Disqus comments when running the Hugo web server on `localhost` (i.e. via `hugo server`) causes the creation of unwanted discussions on the associated Disqus account.

You can create the following `layouts/partials/disqus.html`:

{{< code file="layouts/partials/disqus.html" download="disqus.html" >}}
<div id="disqus_thread"></div>
<script type="text/javascript">

(function() {
    // Don't ever inject Disqus on localhost--it creates unwanted
    // discussions from 'localhost:1313' on your Disqus account...
    if (window.location.hostname == "localhost")
        return;

    var dsq = document.createElement('script'); dsq.type = 'text/javascript'; dsq.async = true;
    var disqus_shortname = '{{ .Site.DisqusShortname }}';
    dsq.src = '//' + disqus_shortname + '.disqus.com/embed.js';
    (document.getElementsByTagName('head')[0] || document.getElementsByTagName('body')[0]).appendChild(dsq);
})();
</script>
<noscript>Please enable JavaScript to view the <a href="http://disqus.com/?ref_noscript">comments powered by Disqus.</a></noscript>
<a href="http://disqus.com/" class="dsq-brlink">comments powered by <span class="logo-disqus">Disqus</span></a>
{{< /code >}}

The `if` statement skips the initialization of the Disqus comment injection when you are running on `localhost`.

You can then render your custom Disqus partial template as follows:

```
{{ partial "disqus.html" . }}
```

## The Internal Templates

* `_internal/disqus.html`
* `_internal/google_news.html`
* `_internal/google_analytics.html`
* `_internal/google_analytics_async.html`
* `_internal/opengraph.html`
* `_internal/pagination.html`
* `_internal/schema.html`
* `_internal/twitter_cards.html`

[disqus]: https://disqus.com
[disqussignup]: https://disqus.com/profile/signup/
