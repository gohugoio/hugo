---
title: YouTube
description: Embed a YouTube video in your content using the youtube shortcode.
categories: [shortcodes]
keywords: []
menu:
  docs:
    parent: shortcodes
    weight:
weight:
toc: true
---

{{% note %}}
To override Hugo's embedded `youtube` shortcode, copy the [source code] to a file with the same name in the `layouts/shortcodes` directory.

[source code]: {{% eturl youtube %}}
{{% /note %}}

## Example

To display a YouTube video with this URL:

```text
https://www.youtube.com/watch?v=0RKpf3rK57I
```

Include this in your Markdown:

```text
{{</* youtube 0RKpf3rK57I */>}}
```

Hugo renders this to:

{{< youtube 0RKpf3rK57I >}}

## Arguments

id
: (`string`) The video `id`. Optional if the `id` is provided as a positional argument as shown in the example above.

allowFullScreen
{{< new-in 0.125.0 >}}
: (`bool`) Whether the `iframe` element can activate full screen mode. Default is `true`.

autoplay
 {{< new-in 0.125.0 >}}
: (`bool`) Whether to automatically play the video. Forces `mute` to `true`. Default is `false`.

class
: (`string`) The `class` attribute of the wrapping `div` element. When specified, removes the `style` attributes from the `iframe` element and its wrapping `div` element.

controls
{{< new-in 0.125.0 >}}
: (`bool`) Whether to display the video controls. Default is `true`.

end
{{< new-in 0.125.0 >}}
: (`int`) The time, measured in seconds from the start of the video, when the player should stop playing the video.

loading
{{< new-in 0.125.0 >}}
: (`string`) The loading attribute of the `iframe` element, either `eager` or `lazy`. Default is `eager`.

loop
{{< new-in 0.125.0 >}}
: (`bool`) Whether to indefinitely repeat the video. Ignores the `start` and `end` arguments after the first play.  Default is `false`.

mute
{{< new-in 0.125.0 >}}
: (`bool`) Whether to mute the video. Always `true` when `autoplay` is `true`. Default is `false`.

start
{{< new-in 0.125.0 >}}
: (`int`) The time, measured in seconds from the start of the video, when the player should start playing the video.

title
: (`string`) The `title` attribute of the `iframe` element. Defaults to "YouTube video".

Example using some of the above:

```text
{{</* youtube id=0RKpf3rK57I start=30 end=60 loading=lazy */>}}
```

## Privacy

Adjust the relevant privacy settings in your site configuration.

{{< code-toggle config=privacy.youTube />}}

disable
: (`bool`) Whether to disable the shortcode. Default is `false`.

privacyEnhanced
: (`bool`) Whether to block YouTube from storing information about visitors on your website unless the user plays the embedded video. Default is `false`.
