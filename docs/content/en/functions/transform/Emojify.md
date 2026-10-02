---
title: transform.Emojify
description: Returns the given string with emoji shortcodes replaced by their corresponding emoji characters.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [emojify]
    returnType: template.HTML
    signatures: [transform.Emojify INPUT]
aliases: [/functions/emojify]
---

See the list of [emoji shortcodes][] for available emoticons.

The `emojify` function can be called in your templates but not directly in your content files by default. For emojis in content files, set [`enableEmoji`][] to `true` in your project configuration. Then you can write emoji shorthand directly into your content files;

```md
I :heart: Hugo!
```

I :heart: Hugo!

[`enableEmoji`]: /configuration/all/#enableemoji
[emoji shortcodes]: /quick-reference/emojis/
