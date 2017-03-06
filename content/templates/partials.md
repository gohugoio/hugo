---
title: Partial Templates
linktitle: Partial Templates
description: Partials are smaller, context-aware components in your list and page templates that can be used economically to keep your templating DRY.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [lists,sections,partials]
weight: 90
draft: false
aliases: [/templates/partial/,/layout/chrome/,/extras/analytics/]
toc: true
wip: true
---

In practice, splitting out reusable template portions into **partial templates** to be included anywhere help keep your templating DRY.

## Partial Template Lookup Order

Partial templates---like [single page templates][singletemps] and [list page templates][listtemps]---have a specific lookup order. However, partials are simpler in that Hugo will only check in two places:

1. `layouts/partials/*<PARTIALNAME>.html`
2. `themes/<THEME>/layouts/partials/*<PARTIALNAME>.html`

This allows a theme's end user to copy a partial's contents into a file of the same name for [further customization][customize].

## Using Partials in your Templates

All partials for your Hugo project are located in a single `layouts/partials` directory. For better organization, you can create multiple subdirectories within `partials` as well:

```
.
└── layouts
    └── partials
        ├── footer
        │   ├── scripts.html
        │   └── site-footer.html
        ├── head
        │   ├── favicons.html
        │   ├── metadata.html
        │   ├── prerender.html
        │   └── twitter.html
        └── header
            ├── site-header.html
            └── site-nav.html
```

All partials are called within your templates using the following pattern:

```
{{ partial "<PATH>/<PARTIAL>.html" . }}
```

{{% note %}}
One of the most common mistakes with new Hugo users is failing to pass a context to the partial call. In the pattern above, note how "the dot" (`.`) is required as the second argument to give the partial context. You can read more about "the dot" in the [Go Template Primer](/templates/go-templates/).
{{% /note %}}

As shown in the above example directory structure, you can nest your directories within `partials` for better source organization. You only need to call the nested partial's path relative to the `partials` directory:

```golang
{{ partial "header/site-header.html" . }}
{{ partial "footer/scripts.html" . }}
```

{{% note %}}
Before v0.12, Hugo used the `template` call to include partial templates. When using Hugo v0.12 and newer, be sure to use the `{{ partial "<PATH>/<PARTIAL>.html" . }}` syntax. The old approach will still work but has fewer benefits.
{{% /note %}}

### Variable Scoping

The second argument in a partial call is the variable being passed down. The above examples are passing the `.`, which tells the template receiving the partial to apply the current [context][context].

This means the partial will *only* be able to access those variables. The partial is isolated and *has no access to the outer scope*. From within the partial, `$.Var` is equivalent to `.Var`.

### Cached Partials

The [`partialCached` template function][partialcached] can offer significant performance gains for complex templates that don't need to be re-rendered on every invocation. The simplest usage is as follows:

```
{{ partialCached "footer.html" . }}
```

You can also pass additional parameters to `partialCached` to create *variants* of the cached partial.

For example, you can tell Hugo to only render the partial `footer.html` once per section:

```
{{ partialCached "footer.html" . .Section }}
```

If you need to pass additional parameters to create unique variants, you can pass as many variant parameters as you need:

```
{{ partialCached "footer.html" . .Params.country .Params.province }}
```

Note that the variant parameters are not made available to the underlying partial template. They are only use to create a unique cache key.

### Example `header.html`

The following `header.html` partial template is used for [spf13.com](http://spf13.com/):

{{% code file="layouts/partials/header.html" download="header.html" %}}
```html
<!DOCTYPE html>
<html class="no-js" lang="en-US" prefix="og: http://ogp.me/ns# fb: http://ogp.me/ns/fb#">
<head>
    <meta charset="utf-8">

    {{ partial "meta.html" . }}

    <base href="{{ .Site.BaseURL }}">
    <title> {{ .Title }} : spf13.com </title>
    <link rel="canonical" href="{{ .Permalink }}">
    {{ if .RSSLink }}<link href="{{ .RSSLink }}" rel="alternate" type="application/rss+xml" title="{{ .Title }}" />{{ end }}

    {{ partial "head_includes.html" . }}
</head>
<body lang="en">
```
{{% /code %}}

{{% note %}}
The `header.html` example partial was built before the introduction of block templates to Hugo. Read more on [base templates and blocks](/templates/base/) for defining the outer chrome or shell of your master templates (i.e., your site's head, header, and footer). You can even combine blocks and partials for added flexibility.
{{% /note %}}

### Example `footer.html`

The following `footer.html` partial template is used for [spf13.com](http://spf13.com/):

{{% code file="layouts/partials/footer.html" download="footer.html" %}}
```html
<footer>
  <div>
    <p>
    &copy; 2013-14 Steve Francia.
    <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons Attribution">Some rights reserved</a>;
    please attribute properly and link back. Hosted by <a href="http://servergrove.com">ServerGrove</a>.
    </p>
  </div>
</footer>
<script type="text/javascript">

  var _gaq = _gaq || [];
  _gaq.push(['_setAccount', 'UA-XYSYXYSY-X']);
  _gaq.push(['_trackPageview']);

  (function() {
    var ga = document.createElement('script');
    ga.src = ('https:' == document.location.protocol ? 'https://ssl' :
        'http://www') + '.google-analytics.com/ga.js';
    ga.setAttribute('async', 'true');
    document.documentElement.firstChild.appendChild(ga);
  })();

</script>
</body>
</html>
```
{{% /code %}}


## Using Hugo's Internal Partial Templates

{{% warning %}}
While the following internal templates are called similar to partials, they do *not* observe the partial template lookup order.
{{% /warning %}}

### Google Analytics

Hugo ships with internal partial templates for Google Analytics tracking, including both synchronous and asynchronous tracking codes.

#### Configuring Google Analytics

Provide your tracking id in your configuration file:

```toml
googleAnalytics = "UA-123-45"
```

```yml
googleAnalytics: "UA-123-45"
```

#### Adding the Google Analytics Template

You can then include the Google Analytics internal partial in your templates:

```golang
{{ template "_internal/google_analytics.html" . }}
```


```golang
{{ template "_internal/google_analytics_async.html" . }}
```

### Disqus

Hugo also ships with a built-in partial for [Disqus comments][disqus], a popular commenting system for both static and dynamic websites. In order to effectively use Disqus, you will need to secure a Disqus "shortname" by [signing up for the free service][disqussignup].

#### Configuring Disqus

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

#### Adding the Disqus Template

To add Disqus, include the following line in templates where you want your comments to appear:

```golang
{{ template "_internal/disqus.html" . }}
```

#### Conditional Loading of Disqus Comments

Users have noticed that enabling Disqus comments when running the Hugo web server on `localhost` (i.e. via `hugo server`) causes the creation of unwanted discussions on the associated Disqus account.

You can create the following `layouts/partials/disqus.html`:

{{% code file="layouts/partials/disqus.html" download="disqus.html" %}}
```html
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
```
{{% /code %}}

The `if` statement skips the initialization of the Disqus comment injection when you are running on `localhost`.

You can then reference then render your custom analytics partial template as follows:

```golang
{{ partial "disqus.html" . }}
```

[context]: /templates/go-templates/ "The most easily overlooked concept to understand about Go templating is how the dot always refers to the current context."
[customize]: /themes/customizing/
[disqus]: https://disqus.com
[disqussignup]: https://disqus.com/profile/signup/
[listtemps]: /templates/lists/
[partialcached]: /functions/partialcached/
[singletemps]: /templates/single-page-templates/
[themes]: /themes/