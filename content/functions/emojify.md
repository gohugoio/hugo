---
title: emojify
linktitle:
description:
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [strings,emojis]
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
---

Runs the string through the Emoji emoticons processor. The result will be declared as "safe" to prevent Go templates from filtering it as unsafe HTML.

See the [Emoji cheat sheet][emojis] for available emoticons.

```
{{ "I :heart: Hugo" | emojify }}
```

The `emojify` function can be called in your templates but not directly in your content files. However, emojis are most often seen inline. The following is a very simple [shortcode template][sc] you can use to add emojis quickly while you write content. It is also the `emo` shortcode used for the Hugo docs. ([See Hugo Docs Shortcodes Source][scsource].)

{{% code file="layouts/shortcodes/emo.html" download="emo.html" %}}
```golang
{{< readfile file="layouts/shortcodes/emo.html" >}}
```
{{% /code %}}

You can then call the shortcode directly in your content using the following:

{{% code file="content/functions/emojify.md" %}}
```golang
I {{</* emo ":heart:" */>}} Hugo!
```
{{% /code %}}

The preceding use of the `emo` shortcode called in content will display as follows:

I {{< emo ":heart:" >}} Hugo!

[emojis]: http://www.emoji-cheat-sheet.com/
[sc]: /templates/shortcode-templates/
[scsource]: https://github.com/spf13/hugo/tree/master/docs/layouts/shortcodes