---
title: transform.Emojify
description: Returns the given string with emoji shortcodes replaced by their corresponding emoji characters.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [emojify]
    returnType: any
    signatures: [transform.Emojify INPUT]
---

If the argument is of type `template.HTML`, returns `template.HTML`, else returns a `string`.

See the list of [emoji shortcodes][] for available emoticons.

The `emojify` function can be called in your templates but not directly in your content files by default. For emojis in content files, set [`enableEmoji`][] to `true` in your project configuration. Then you can write emoji shorthand directly into your content files;

```md
I :heart: Hugo!
```

I :heart: Hugo!

[`enableEmoji`]: /configuration/all/#enableemoji
[emoji shortcodes]: /quick-reference/emojis/
