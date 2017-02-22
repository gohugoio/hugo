---
title: Comments
linktitle: Comments
description: Hugo ships with an internal Disqus template that requires negligible effort to implement, but Disqus isn't the only commenting system that will work with your new Hugo website.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [sections,content,organization, fundamentals]
categories: [project organization]
weight: 140
draft: false
aliases: [/extras/comments/]
toc: true
needsreview: true
---

<!-- pulled from extras/comments -->

As Hugo is a static site generator, the content produced is static and doesn’t interact with the users. The most common interaction people ask for is comment capability.

Hugo ships with support for [Disqus](https://disqus.com/), a third-party service that provides comment and community capabilities to website via JavaScript.

Your theme may already support Disqus, but even it if doesn’t, it is easy to add.

## Disqus

### Adding Disqus to a Template

Hugo comes with all the code you would need to include load Disqus. Simply include the following line where you want your comments to appear:

```golang
{{ template "_internal/disqus.html" . }}
```

### Configuring Disqus

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

### Conditional Loading of Disqus Comments

Users have noticed that enabling Disqus comments when running the Hugo web server on `localhost` (i.e. via `hugo server`) causes the creation of unwanted discussions on the associated Disqus account. In order to prevent this, a slightly tweaked partial template is required. So, rather than using the built-in `"_internal/disqus.html"` template referenced above, create a template in `layouts/partials` that looks like the following:

{{% input "layouts/partials/disqus.html" %}}
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
{{% /input %}}

The `if` statement skips the initialization of the Disqus comment injection when you are running on `localhost`.

You can then reference the partial template:

{{% input "disqus-reference.html" %}}
```golang
{{ partial "disqus.html" . }}
```
{{% /input %}}

## Alternatives

There are a few alternatives to commenting on static sites for those who do not want to use Hugo's built-in Disqus support:

* [Static Man](https://staticman.net/)
* [IntenseDebate](http://intensedebate.com/)
* [Livefyre](http://livefyre.com/)
* [Muut](http://muut.com/)
* [多说](http://duoshuo.com/) ([Duoshuo](http://duoshuo.com/), popular in China)
* [isso](http://posativ.org/isso/) (Self-hosted, Python)
* [Kaiju](https://github.com/spf13/kaiju)

## Kaiju

[Kaiju](https://github.com/spf13/kaiju) is an open-source project started by [spf13](http://spf13.com/) (Hugo’s author) to bring easy and fast real time discussions to the web.

Written using Go, Socket.io and MongoDB, it is very fast and easy to deploy.

It is in early development but shows promise. If you have interest, please help by contributing whether via a pull request, an issue or even just a tweet. Everything helps.

## Discourse

Additionally, you may recognize [Discourse](http://www.discourse.org) as the system that powers the [Hugo Discussion Forum](http://discuss.gohugo.io).