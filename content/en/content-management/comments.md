---
title: Comments
linktitle: Comments
description: Hugo ships with an internal Disqus template, but this isn't the only commenting system that will work with your new Hugo website.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-03-09
keywords: [sections,content,organization]
categories: [project organization, fundamentals]
menu:
  docs:
    parent: "content-management"
    weight: 140
weight: 140	#rem
draft: false
aliases: [/extras/comments/]
toc: true
---

Hugo ships with support for [Disqus](https://disqus.com/), a third-party service that provides comment and community capabilities to websites via JavaScript.

Your theme may already support Disqus, but if not, it is easy to add to your templates via [Hugo's built-in Disqus partial][disquspartial].

## Add Disqus

Hugo comes with all the code you need to load Disqus into your templates. Before adding Disqus to your site, you'll need to [set up an account][disqussetup].

### Configure Disqus

Disqus comments require you set a single value in your [site's configuration file][configuration] like so:

{{< code-toggle copy="false" >}}
disqusShortname = "yourdiscussshortname"
{{</ code-toggle >}}

For many websites, this is enough configuration. However, you also have the option to set the following in the [front matter][] of a single content file:

* `disqus_identifier`
* `disqus_title`
* `disqus_url`

### Render Hugo's Built-in Disqus Partial Template

Disqus has its own [internal template](https://gohugo.io/templates/internal/#disqus) available, to render it add the following code where you want comments to appear:

```
{{ template "_internal/disqus.html" . }}
```

## Comments Alternatives

There are a few alternatives to commenting on static sites for those who do not want to use Disqus:

* [Static Man](https://staticman.net/)
* [Talkyard](https://www.talkyard.io/blog-comments) (Open source, & serverless hosting)
* [txtpen](https://txtpen.github.io/hn/)
* [IntenseDebate](http://intensedebate.com/)
* [Graph Comment][]
* [Muut](http://muut.com/)
* [isso](http://posativ.org/isso/) (Self-hosted, Python)
    * [Tutorial on Implementing Isso with Hugo][issotutorial]
* [Utterances](https://utteranc.es/) (Open source, Github comments widget built on Github issues)
* [Remark](https://github.com/umputun/remark) (Open source, Golang, Easy to run docker)

<!-- I don't think this is worth including in the documentation since it seems that Steve is no longer supporting or developing this project. rdwatters - 2017-02-29.-->
<!-- * [Kaiju](https://github.com/spf13/kaiju) -->

<!-- ## Kaiju

[Kaiju](https://github.com/spf13/kaiju) is an open-source project started by [spf13](http://spf13.com/) (Hugo’s author) to bring easy and fast real time discussions to the web.

Written using Go, Socket.io, and [MongoDB][], Kaiju is very fast and easy to deploy.

It is in early development but shows promise. If you have interest, please help by contributing via pull request, [opening an issue in the Kaiju GitHub repository][kaijuissue], or [Tweeting about it][tweet]. Every bit helps. -->

[configuration]: /getting-started/configuration/
[disquspartial]: /templates/partials/#disqus
[disqussetup]: https://disqus.com/profile/signup/
[forum]: https://discourse.gohugo.io
[front matter]: /content-management/front-matter/
[Graph Comment]: https://graphcomment.com/
[kaijuissue]: https://github.com/spf13/kaiju/issues/new
[issotutorial]: https://stiobhart.net/2017-02-24-isso-comments/
[partials]: /templates/partials/
[MongoDB]: https://www.mongodb.com/
[tweet]: https://twitter.com/spf13
