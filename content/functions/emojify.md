---
title: emojify
linktitle:
description: Runs a string through the Emoji emoticons processor.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
#tags: [strings,emojis]
ns:
signature: ["emojify INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
---

`emoji` runs a passed string through the Emoji emoticons processor. The result will be declared as "safe" to prevent Go templates from filtering it as unsafe HTML.

See the [Emoji cheat sheet][emojis] for available emoticons.

The `emojify` function can be called in your templates but not directly in your content files by default. To add emojis directly into content without further templating or shortcode considerations, set `enableEmoji` to `true` in your site's [configuration][config].

## Example

When enabled, you can write emoji shorthand directly into your content files; e.g. <code>I :</code><code>heart</code><code>: Hugo!</code>:

I :heart: Hugo!


[config]: /getting-started/configuration/
[emojis]: http://www.emoji-cheat-sheet.com/
[sc]: /templates/shortcode-templates/
[scsource]: https://github.com/spf13/hugo/tree/master/docs/layouts/shortcodes
