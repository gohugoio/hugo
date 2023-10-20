---
title: transform.Emojify 
linkTitle: emojify
description: Runs a string through the Emoji emoticons processor.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [emojify]
  returnType: template.HTML
  signatures: [transform.Emojify INPUT]
namespace: transform
relatedFunctions: []
aliases: [/functions/emojify]
---

`emojify` runs a passed string through the Emoji emoticons processor.

See the [Emoji cheat sheet][emojis] for available emoticons.

The `emojify` function can be called in your templates but not directly in your content files by default. For emojis in content files, set `enableEmoji` to `true` in your site's [configuration]. Then you can write emoji shorthand directly into your content files; e.g. <code>I :</code><code>heart</code><code>: Hugo!</code>:

I :heart: Hugo!


[configuration]: /getting-started/configuration/
[emojis]: https://www.webfx.com/tools/emoji-cheat-sheet/
[sc]: /templates/shortcode-templates/
[scsource]: https://github.com/gohugoio/hugo/tree/master/docs/layouts/shortcodes
