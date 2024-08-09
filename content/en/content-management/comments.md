---
title: Comments
description: Hugo ships with an internal Disqus template, but this isn't the only commenting system that will work with your new Hugo website.
categories: [content management]
keywords: [sections,content,organization]
menu:
  docs:
    parent: content-management
    weight: 220
weight: 220
toc: true
aliases: [/extras/comments/]
---

Hugo ships with support for [Disqus](https://disqus.com/), a third-party service that provides comment and community capabilities to websites via JavaScript.

Your theme may already support Disqus, but if not, it is easy to add to your templates via [Hugo's built-in Disqus partial][disquspartial].

## Add Disqus

Hugo comes with all the code you need to load Disqus into your templates. Before adding Disqus to your site, you'll need to [set up an account][disqussetup].

### Configure Disqus

Disqus comments require you set a single value in your [site's configuration file][configuration] like so:

{{< code-toggle file=hugo >}}
[services.disqus]
shortname = 'your-disqus-shortname'
{{</ code-toggle >}}

For many websites, this is enough configuration. However, you also have the option to set the following in the [front matter] of a single content file:

* `disqus_identifier`
* `disqus_title`
* `disqus_url`

### Render Hugo's built-in Disqus partial template

Disqus has its own [internal template](/templates/embedded/#disqus) available, to render it add the following code where you want comments to appear:

```go-html-template
{{ template "_internal/disqus.html" . }}
```

## Alternatives

Commercial commenting systems:

- [Emote](https://emote.com/)
- [Graph Comment](https://graphcomment.com/)
- [Hyvor Talk](https://talk.hyvor.com/)
- [IntenseDebate](https://intensedebate.com/)
- [ReplyBox](https://getreplybox.com/)

Open-source commenting systems:

- [Cactus Comments](https://cactus.chat/docs/integrations/hugo/)
- [Comentario](https://gitlab.com/comentario/comentario/)
- [Comma](https://github.com/Dieterbe/comma/)
- [Commento](https://commento.io/)
- [Discourse](https://meta.discourse.org/t/embed-discourse-comments-on-another-website-via-javascript/31963)
- [Giscus](https://giscus.app/)
- [Isso](https://isso-comments.de/)
- [Remark42](https://remark42.com/)
- [Staticman](https://staticman.net/)
- [Talkyard](https://blog-comments.talkyard.io/)
- [Utterances](https://utteranc.es/)

[configuration]: /getting-started/configuration/
[disquspartial]: /templates/embedded/#disqus
[disqussetup]: https://disqus.com/profile/signup/
[forum]: https://discourse.gohugo.io
[front matter]: /content-management/front-matter/
[kaijuissue]: https://github.com/spf13/kaiju/issues/new
[issotutorial]: https://stiobhart.net/2017-02-24-isso-comments/
[partials]: /templates/partial/
[MongoDB]: https://www.mongodb.com/
