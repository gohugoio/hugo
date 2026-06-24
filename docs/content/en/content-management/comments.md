---
title: Comments
description: Hugo ships with an embedded Disqus partial, but this isn't the only commenting system that will work with your new Hugo website.
categories: []
keywords: []
aliases: [/extras/comments/]
---

Hugo ships with support for [Disqus][], a third-party service that provides comment and community capabilities to websites via JavaScript.

Your theme may already support Disqus, but if not, it is easy to add to your templates via Hugo's [embedded partial][].

## Add Disqus

Hugo comes with all the code you need to load Disqus into your templates. Before adding Disqus to your site, you'll need to [set up an account][].

### Configure Disqus

Disqus comments require you set a single value in your project configuration:

{{< code-toggle file=hugo >}}
[services.disqus]
shortname = 'your-disqus-shortname'
{{</ code-toggle >}}

For many websites, this is enough configuration. However, you also have the option to set the following in the front matter of a single content file:

- `params.disqus_identifier`
- `params.disqus_title`
- `params.disqus_url`

### Render Hugo's embedded Disqus partial

To render it, add the following code where you want comments to appear:

```go-html-template
{{ partial "disqus.html" . }}
```

## Alternatives

Commercial commenting systems:

- [Commentix][]
- [Emote][]
- [FastComments][]
- [Graph Comment][]
- [Hyvor Talk][]
- [IntenseDebate][]
- [ReplyBox][]

Open-source commenting systems:

- [Cactus Comments][]
- [Comentario][]
- [Comma][]
- [Discourse][]
- [Giscus][]
- [Isso][]
- [Remark42][]
- [Staticman][]
- [Talkyard][]
- [Utterances][]
- [Zoomment][]

[Cactus Comments]: https://cactus.chat/docs/integrations/hugo/
[Comentario]: https://gitlab.com/comentario/comentario/
[Comma]: https://github.com/Dieterbe/comma/
[Commentix]: https://www.commentix.com/
[Discourse]: https://meta.discourse.org/t/embed-discourse-comments-on-another-website-via-javascript/31963
[Disqus]: https://disqus.com/
[Emote]: https://emote.com/
[FastComments]: https://fastcomments.com/commenting-system-for-hugo
[Giscus]: https://giscus.app/
[Graph Comment]: https://graphcomment.com/
[Hyvor Talk]: https://talk.hyvor.com/
[IntenseDebate]: https://intensedebate.com/
[Isso]: https://isso-comments.de/
[Remark42]: https://remark42.com/
[ReplyBox]: https://getreplybox.com/
[Staticman]: https://staticman.net/
[Talkyard]: https://blog-comments.talkyard.io/
[Utterances]: https://utteranc.es/
[Zoomment]: https://zoomment.com/
[embedded partial]: /templates/embedded/#disqus
[set up an account]: https://disqus.com/profile/signup/
