---
title: emojify
description: Runs a string through the Emoji emoticons processor.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings,emojis]
signature: ["emojify INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
---

`emoji` runs a passed string through the Emoji emoticons processor.

See the [Emoji cheat sheet][emojis] for available emoticons.

The `emojify` function can be called in your templates but not directly in your content files by default. For emojis in content files, set `enableEmoji` to `true` in your site's [configuration][config]. Then you can write emoji shorthand directly into your content files; e.g. <code>I :</code><code>heart</code><code>: Hugo!</code>:

I :heart: Hugo!


[config]: /getting-started/configuration/
[emojis]: https://www.webfx.com/tools/emoji-cheat-sheet/
[sc]: /templates/shortcode-templates/
[scsource]: https://github.com/gohugoio/hugo/tree/master/docs/layouts/shortcodes
