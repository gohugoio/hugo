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

Runs the string through the Emoji emoticons processor. The result will be declared as "safe" to prevent Go templates from filtering it.

See the [Emoji cheat sheet][emojis] for available emoticons. The `emojify` function can be called in your templates but not content.

```
{{ "I :heart: Hugo" | emojify }}
```

However, emojis are most often seen inline. The following is a very simple [partial template][partials] you can use to add emojis directly to your content files. It is also the `emo` shortcode used for the Hugo docs.

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

The output of which is...

I {{< emo ":heart:" >}} Hugo!

[emojis]: http://www.emoji-cheat-sheet.com/
[partials]: /templates/partials/