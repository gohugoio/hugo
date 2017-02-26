---
title: Partial Templates
linktitle: Partial Templates
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [lists,sections,partials]
weight: 90
draft: false
aliases: [/templates/partials/,/layout/chrome/,/extras/analytics/]
toc: true
notesforauthors:
---

## Partials vs Templates

In practice, it's very convenient to split out common template portions into a
partial template that can be included anywhere. As you create the rest of your
templates, you will include templates from the `/layouts/partials/` directory
or from arbitrary partial subdirectories like `/layouts/partials/post/tag/`.

Partials are especially important for [themes][] because they give theme users an opportunity to [overwrite just a small portion of a theme][customize] while maintaining compatibility with the theme's upstream.

Theme developers may want to include a few partials with empty HTML
files in the theme just so end users have an easy place to inject their
customized content.

I've found it helpful to include a header and footer template in
partials so I can include those in all the full page layouts.  There is
nothing special about header.html and footer.html other than they seem
like good names to use for inclusion in your other templates.

```bash
▾ layouts/
  ▾ partials/
      header.html
      footer.html
```

## Partial vs Template

Version v0.12 of Hugo introduced the `partial` call inside the template system.
This is a change to the way partials were handled previously inside the
template system. In earlier versions, Hugo didn’t treat partials specially, and
you could include a partial template with the `template` call in the standard
template language.

With the addition of the theme system in v0.11, it became apparent that a theme & override-aware partial was needed.

When using Hugo v0.12 and above, please use the `partial` call (and leave out the “partial/” path). The old approach would still work, but wouldn’t benefit from the ability to have users override the partial theme file with local layouts.

## Example header.html

This header template is used for [spf13.com](http://spf13.com/):

{{% code file="layouts/partials/header.html" %}}
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

## Example footer.html

This footer template is used for [spf13.com](http://spf13.com/):

{{% code file="layouts/partials/footer.html" %}}
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

To reference a partial template stored in a subfolder, e.g. `/layouts/partials/post/tag/list.html`, call it this way:

```
{{ partial "post/tag/list" . }}
```

Note that the subdirectories you create under /layouts/partials can be named whatever you like.

For more examples of referencing these templates, see [single content templates](/templates/content/), [list templates](/templates/list/) and [homepage templates](/templates/homepage/).

## Variable scoping

As you might have noticed, `partial` calls receive two parameters.

1. The first is the name of the partial and determines the file
location to be read.
2. The second is the variables to be passed down to the partial.

This means that the partial will _only_ be able to access those variables. It is isolated and has no access to the outer scope. From within the partial, `$.Var` is equivalent to `.Var`.

## Cached Partials

The `partialCached` template function can offer significant performance gains for complex templates that don't need to be rerendered upon every invocation. The simplest usage is as follows:

```
{{ partialCached "footer.html" . }}
```

You can also pass additional parameters to `partialCached` to create *variants* of the cached partial. For example, say you have a complex partial that should be identical when rendered for pages within the same section. You could use a variant based upon section so that the partial is only rendered once per section:

```
{{ partialCached "footer.html" . .Section }}
```

If you need to pass additional parameters to create unique variants, you can pass as many variant parameters as you need:

```
{{ partialCached "footer.html" . .Params.country .Params.province }}
```

Note that the variant parameters are not made available to the underlying partial template. They are only use to create a unique cache key.

## Using Hugo's Built-in Partials

### Google Analytics

Hugo ships with prebuilt internal partial templates for Google Analytics tracking, including both synchronous and asynchronous tracking codes.

<!-- pulled from extras/analytics -->

#### Configuring Google Analytics

Provide your tracking id in your configuration file, e.g. config.yaml.

```toml
googleAnalytics = "UA-123-45"
```

#### Google Analytics Example

Include the internal template in your templates like so:

{{% code file="call-ga.html" %}}
```golang
{{ template "_internal/google_analytics.html" . }}
```
{{% /code %}}


{{% code file="call-ga-async.html" %}}
```golang
{{ template "_internal/google_analytics_async.html" . }}
```
{{% /code %}}

<!-- pulled from extras/comments -->
### Disqus

Hugo also ships with a built-in partial for [Disqus comments][disqus], a popular commenting system for both static and dynamic websites.

#### Adding Disqus to a Template

Hugo comes with all the code you would need to include load Disqus. Simply include the following line where you want your comments to appear:

```golang
{{ template "_internal/disqus.html" . }}
```

#### Configuring Disqus

That template requires you to set a single value in your site `config`:

```toml
disqusShortname = "yourdiscussshortname"
```

Or with a `config.yml`:

```yaml
disqusShortname: "yourdiscussshortname"
```

You also have the option to set the following in the front matter for a given piece of content:

* `disqus_identifier`
* `disqus_title`
* `disqus_url`

#### Conditional Loading of Disqus Comments

Users have noticed that enabling Disqus comments when running the Hugo web server on `localhost` (i.e. via `hugo server`) causes the creation of unwanted discussions on the associated Disqus account. In order to prevent this, a slightly tweaked partial template is required. So, rather than using the built-in `"_internal/disqus.html"` template referenced above, create a template in `layouts/partials` that looks like the following:

{{% code file="layouts/partials/disqus.html" %}}
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

You can then reference the partial template:

{{% code file="disqus-reference.html" %}}
```golang
{{ partial "disqus.html" . }}
```
{{% /code %}}

[themes]: /themes/
[customize]: /themes/customizing/
[disqus]: https://disqus.com