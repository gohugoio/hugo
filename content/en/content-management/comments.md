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

* [Staticman](https://staticman.net/)
* [Talkyard](https://www.talkyard.io/blog-comments) (Open source, & serverless hosting)
* [IntenseDebate](https://intensedebate.com/)
* [Graph Comment][]
* [Muut](https://muut.com/)
* [isso](https://posativ.org/isso/) (Self-hosted, Python)
    * [Tutorial on Implementing Isso with Hugo][issotutorial]
* [Utterances](https://utteranc.es/) (Open source, GitHub comments widget built on GitHub issues)
* [Remark](https://github.com/umputun/remark) (Open source, Golang, Easy to run docker)
* [Commento](https://commento.io/) (Open Source, available as a service, local install, or docker image)
* [JustComments](https://just-comments.com) (Open Source, available as a service, can be self-hosted)


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
