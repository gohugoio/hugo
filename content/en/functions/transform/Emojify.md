---
title: transform.Emojify 
description: Runs a string through the Emoji emoticons processor.
categories: []
keywords: []
action:
  aliases: [emojify]
  related: []
  returnType: template.HTML
  signatures: [transform.Emojify INPUT]
aliases: [/functions/emojify]
---

`emojify` runs a passed string through the Emoji emoticons processor.

See the list of [emoji shortcodes] for available emoticons.

The `emojify` function can be called in your templates but not directly in your content files by default. For emojis in content files, set `enableEmoji` to `true` in your site's [configuration]. Then you can write emoji shorthand directly into your content files;


```text
I :heart: Hugo!
```

I :heart: Hugo!

[configuration]: /getting-started/configuration/
[emoji shortcodes]: /quick-reference/emojis/
[sc]: /templates/shortcode/
[scsource]: https://github.com/gohugoio/hugo/tree/master/docs/layouts/shortcodes
