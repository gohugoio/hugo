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

{{< code-toggle file="config" >}}
googleAnalytics = "UA-123-45"
{{</ code-toggle >}}

### Use the Google Analytics Template

You can then include the Google Analytics internal template:

```
{{ template "_internal/google_analytics.html" . }}
```


```
{{ template "_internal/google_analytics_async.html" . }}
```

A `.Site.GoogleAnalytics` variable is also exposed from the config.

## Disqus

Hugo also ships with an internal template for [Disqus comments][disqus], a popular commenting system for both static and dynamic websites. In order to effectively use Disqus, you will need to secure a Disqus "shortname" by [signing up for the free service][disqussignup].

### Configure Disqus

To use Hugo's Disqus template, you first need to set a single value in your site's `config.toml` or `config.yml`:

{{< code-toggle file="config" >}}
disqusShortname = "yourdiscussshortname"
{{</ code-toggle >}}

You also have the option to set the following in the front matter for a given piece of content:

* `disqus_identifier`
* `disqus_title`
* `disqus_url`

### Use the Disqus Template

To add Disqus, include the following line in templates where you want your comments to appear:

```
{{ template "_internal/disqus.html" . }}
```

A `.Site.DisqusShortname` variable is also exposed from the config.

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
<noscript>Please enable JavaScript to view the <a href="https://disqus.com/?ref_noscript">comments powered by Disqus.</a></noscript>
<a href="https://disqus.com/" class="dsq-brlink">comments powered by <span class="logo-disqus">Disqus</span></a>
{{< /code >}}

The `if` statement skips the initialization of the Disqus comment injection when you are running on `localhost`.

You can then render your custom Disqus partial template as follows:

```
{{ partial "disqus.html" . }}
```

## Open Graph
An internal template for the [Open Graph protocol](https://ogp.me/), metadata that enables a page to become a rich object in a social graph.
This format is used for Facebook and some other sites.

### Configure Open Graph

Hugo's Open Graph template is configured using a mix of configuration variables and [front-matter](/content-management/front-matter/) on individual pages.

{{< code-toggle file="config" >}}
[params]
  title = "My cool site"
  images = ["site-feature-image.jpg"]
  description = "Text about my cool site"
[taxonomies]
  series = "series"
{{</ code-toggle >}}

{{< code-toggle file="content/blog/my-post" >}}
title = "Post title"
description = "Text about this post"
date = "2006-01-02"
images = ["post-cover.png"]
audio = []
videos = []
series = []
tags = []
{{</ code-toggle >}}

Hugo uses the page title and description for the title and description metadata.
The first 6 URLs from the `images` array are used for image metadata.

Various optional metadata can also be set:

- Date, published date, and last modified data are used to set the published time metadata if specified.
- `audio` and `videos` are URL arrays like `images` for the audio and video metadata tags, respectively.
- The first 6 `tags` on the page are used for the tags metadata.
- The `series` taxonomy is used to specify related "see also" pages by placing them in the same series.

If using YouTube this will produce a og:video tag like `<meta property="og:video" content="url">`. If using a YouTube link make sure this is in **https://www.youtube.com/v/NlXVWtgLNjY** not __https://www.youtube.com/watch?v=NlXVWtgLNjY__

### Use the Open Graph Template

To add Open Graph metadata, include the following line between the `<head>` tags in your templates:

```
{{ template "_internal/opengraph.html" . }}
```

## Twitter Cards

An internal template for [Twitter Cards](https://developer.twitter.com/en/docs/tweets/optimize-with-cards/overview/abouts-cards),
metadata used to attach rich media to Tweets linking to your site.

### Configure Twitter Cards

Hugo's Twitter Card template is configured using a mix of configuration variables and [front-matter](/content-management/front-matter/) on individual pages.

{{< code-toggle file="config" >}}
[params]
  images = ["site-feature-image.jpg"]
  description = "Text about my cool site"
{{</ code-toggle >}}

{{< code-toggle file="content/blog/my-post" >}}
title = "Post title"
description = "Text about this post"
images = ["post-cover.png"]
{{</ code-toggle >}}

If `images` aren't specified in the page front-matter, then hugo searches for [image page resources](/content-management/image-processing/) with `feature`, `cover`, or `thumbnail` in their name.
If no image resources with those names are found, the images defined in the [site config](/getting-started/configuration/) are used instead.
If no images are found at all, then an image-less Twitter `summary` card is used instead of `summary_large_image`.

Hugo uses the page title and description for the card's title and description fields. The page summary is used if no description is given.

### Use the Twitter Cards Template

To add Twitter card metadata, include the following line between the `<head>` tags in your templates:

```
{{ template "_internal/twitter_cards.html" . }}
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
