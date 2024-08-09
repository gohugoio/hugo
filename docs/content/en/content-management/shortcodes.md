---
title: Shortcodes
description: Shortcodes are simple snippets inside your content files calling built-in or custom templates.
categories: [content management]
keywords: [markdown,content,shortcodes]
menu:
  docs:
    parent: content-management
    weight: 100
weight: 100
toc: true
aliases: [/extras/shortcodes/]
testparam: "Hugo Rocks!"
---

## What a shortcode is

Hugo loves Markdown because of its simple content format, but there are times when Markdown falls short. Often, content authors are forced to add raw HTML (e.g., video `<iframe>`'s) to Markdown content. We think this contradicts the beautiful simplicity of Markdown's syntax.

Hugo created **shortcodes** to circumvent these limitations.

A shortcode is a simple snippet inside a content file that Hugo will render using a predefined template. Note that shortcodes will not work in template files. If you need the type of drop-in functionality that shortcodes provide but in a template, you most likely want a [partial template][partials] instead.

In addition to cleaner Markdown, shortcodes can be updated any time to reflect new classes, techniques, or standards. At the point of site generation, Hugo shortcodes will easily merge in your changes. You avoid a possibly complicated search and replace operation.

## Use shortcodes

{{< youtube 2xkNJL4gJ9E >}}

In your content files, a shortcode can be called by calling `{{%/* shortcodename arguments */%}}`. Shortcode arguments are space delimited, and arguments with internal spaces must be quoted.

The first word in the shortcode declaration is always the name of the shortcode. Arguments follow the name. Depending upon how the shortcode is defined, the arguments may be named, positional, or both, although you can't mix argument types in a single call. The format for named arguments models that of HTML with the format `name="value"`.

Some shortcodes use or require closing shortcodes. Again like HTML, the opening and closing shortcodes match (name only) with the closing declaration, which is prepended with a slash.

Here are two examples of paired shortcodes:

```go-html-template
{{%/* mdshortcode */%}}Stuff to `process` in the *center*.{{%/* /mdshortcode */%}}
```

```go-html-template
{{</* highlight go */>}} A bunch of code here {{</* /highlight */>}}
```

The examples above use two different delimiters, the difference being the `%` character in the first and the `<>` characters in the second.

### Shortcodes with raw string arguments

You can pass multiple lines as arguments to a shortcode by using raw string literals:

```go-html-template
{{</*  myshortcode `This is some <b>HTML</b>,
and a new line with a "quoted string".` */>}}
```

### Shortcodes with Markdown

Shortcodes using the `%` as the outer-most delimiter will be fully rendered when sent to the content renderer. This means that the rendered output from a shortcode can be part of the page's table of contents, footnotes, etc.

### Shortcodes without Markdown

The `<` character indicates that the shortcode's inner content does *not* need further rendering. Often shortcodes without Markdown include internal HTML:

```go-html-template
{{</* myshortcode */>}}<p>Hello <strong>World!</strong></p>{{</* /myshortcode */>}}
```

### Nested shortcodes

You can call shortcodes within other shortcodes by creating your own templates that leverage the `.Parent` method. `.Parent` allows you to check the context in which the shortcode is being called. See [Shortcode templates][sctemps].

## Embedded shortcodes

Use these embedded shortcodes as needed.

### figure

{{% note %}}
To override Hugo's embedded `figure` shortcode, copy the [source code] to a file with the same name in the layouts/shortcodes directory.

[source code]: {{% eturl figure %}}
{{% /note %}}

The `figure` shortcode can use the following named arguments:

src
: URL of the image to be displayed.

link
: If the image needs to be hyperlinked, URL of the destination.

target
: Optional `target` attribute for the URL if `link` argument is set.

rel
: Optional `rel` attribute for the URL if `link` argument is set.

alt
: Alternate text for the image if the image cannot be displayed.

title
: Image title.

caption
: Image caption. Markdown within the value of `caption` will be rendered.

class
: `class` attribute of the HTML `figure` tag.

height
: `height` attribute of the image.

width
: `width` attribute of the image.

loading
: `loading` attribute of the image.

attr
: Image attribution text. Markdown within the value of `attr` will be rendered.

attrlink
: If the attribution text needs to be hyperlinked, URL of the destination.

Example usage:

```text
{{</* figure src="elephant.jpg" title="An elephant at sunset" */>}}
```

Rendered:

```html
<figure>
  <img src="elephant.jpg">
  <figcaption><h4>An elephant at sunset</h4></figcaption>
</figure>
```

### gist

{{% note %}}
To override Hugo's embedded `gist` shortcode, copy the [source code] to a file with the same name in the layouts/shortcodes directory.

[source code]: {{% eturl gist %}}
{{% /note %}}

To display a GitHub [gist] with this URL:

[gist]: https://docs.github.com/en/get-started/writing-on-github/editing-and-sharing-content-with-gists

```text
https://gist.github.com/user/50a7482715eac222e230d1e64dd9a89b
```

Include this in your Markdown:

```text
{{</* gist user 50a7482715eac222e230d1e64dd9a89b */>}}
```

This will display all files in the gist alphabetically by file name.

{{< gist jmooring 23932424365401ffa5e9d9810102a477 >}}

To display a specific file within the gist:

```text
{{</* gist user 23932424365401ffa5e9d9810102a477 list.html */>}}
```

Rendered:

{{< gist jmooring 23932424365401ffa5e9d9810102a477 list.html >}}

### highlight

{{% note %}}
To override Hugo's embedded `highlight` shortcode, copy the [source code] to a file with the same name in the layouts/shortcodes directory.

[source code]: {{% eturl highlight %}}
{{% /note %}}

To display a highlighted code sample:

```text
{{</* highlight go-html-template */>}}
{{ range .Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
{{</* /highlight */>}}
```

Rendered:

{{< highlight go-html-template >}}
{{ range .Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
{{< /highlight >}}

To specify one or more [highlighting options], include a quotation-encapsulated, comma-separated list:

[highlighting options]: /functions/transform/highlight/

```text
{{</* highlight go-html-template "lineNos=inline, lineNoStart=42" */>}}
{{ range .Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
{{</* /highlight */>}}
```

Rendered:

{{< highlight go-html-template "lineNos=inline, lineNoStart=42" >}}
{{ range .Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
{{< /highlight >}}

### instagram

{{% note %}}
To override Hugo's embedded `instagram` shortcode, copy the [source code] to a file with the same name in the layouts/shortcodes directory.

[source code]: {{% eturl instagram %}}
{{% /note %}}

To display an Instagram post with this URL:

```text
https://www.instagram.com/p/CxOWiQNP2MO/
```

Include this in your Markdown:

```text
{{</* instagram CxOWiQNP2MO */>}}
```

Rendered:

{{< instagram CxOWiQNP2MO >}}

### param

{{% note %}}
To override Hugo's embedded `param` shortcode, copy the [source code] to a file with the same name in the layouts/shortcodes directory.

[source code]: {{% eturl param %}}
{{% /note %}}

The `param` shortcode renders a parameter from the page's front matter, falling back to a site parameter of the same name. The shortcode throws an error if the parameter does not exist.

Example usage:

```text
{{</* param testparam */>}}
```

Access nested values by [chaining] the [identifiers]:

[chaining]: /getting-started/glossary/#chain
[identifiers]: /getting-started/glossary/#identifier

```text
{{</* param my.nested.param */>}}
```

### ref

{{% note %}}
To override Hugo's embedded `ref` shortcode, copy the [source code] to a file with the same name in the layouts/shortcodes directory.

Always use the `{{%/* */%}}` notation when calling this shortcode.

[source code]: {{% eturl ref %}}
{{% /note %}}

The `ref` shortcode returns the permalink of the given page reference.

Example usage:

```text
[Post 1]({{%/* ref "/posts/post-1" */%}})
[Post 1]({{%/* ref "/posts/post-1.md" */%}})
[Post 1]({{%/* ref "/posts/post-1#foo" */%}})
[Post 1]({{%/* ref "/posts/post-1.md#foo" */%}})
```

Rendered:

```html
<a href="http://example.org/posts/post-1/">Post 1</a>
<a href="http://example.org/posts/post-1/">Post 1</a>
<a href="http://example.org/posts/post-1/#foo">Post 1</a>
<a href="http://example.org/posts/post-1/#foo">Post 1</a>
```

### relref

{{% note %}}
To override Hugo's embedded `relref` shortcode, copy the [source code] to a file with the same name in the layouts/shortcodes directory.

Always use the `{{%/* */%}}` notation when calling this shortcode.

[source code]: {{% eturl relref %}}
{{% /note %}}

The `relref` shortcode returns the permalink of the given page reference.

Example usage:

```text
[Post 1]({{%/* relref "/posts/post-1" */%}})
[Post 1]({{%/* relref "/posts/post-1.md" */%}})
[Post 1]({{%/* relref "/posts/post-1#foo" */%}})
[Post 1]({{%/* relref "/posts/post-1.md#foo" */%}})
```

Rendered:

```html
<a href="/posts/post-1/">Post 1</a>
<a href="/posts/post-1/">Post 1</a>
<a href="/posts/post-1/#foo">Post 1</a>
<a href="/posts/post-1/#foo">Post 1</a>
```

### twitter

{{% note %}}
To override Hugo's embedded `twitter` shortcode, copy the [source code] to a file with the same name in the layouts/shortcodes directory.

You may call the `twitter` shortcode by using its `tweet` alias.

[source code]: {{% eturl twitter %}}
{{% /note %}}

To display a Twitter post with this URL:

```txt
https://x.com/SanDiegoZoo/status/1453110110599868418
```

Include this in your Markdown:

```text
{{</* twitter user="SanDiegoZoo" id="1453110110599868418" */>}}
```

Rendered:

{{< twitter user="SanDiegoZoo" id="1453110110599868418" >}}

### vimeo

{{% note %}}
To override Hugo's embedded `vimeo` shortcode, copy the [source code] to a file with the same name in the layouts/shortcodes directory.

[source code]: {{% eturl vimeo %}}
{{% /note %}}

To display a Vimeo video with this URL:

```text
https://vimeo.com/channels/staffpicks/55073825
```

Include this in your Markdown:

```text
{{</* vimeo 55073825 */>}}
```

Rendered:

{{< vimeo 55073825 >}}

{{% note %}}
If you want to further customize the visual styling, add a `class` argument when calling the shortcode. The new `class` will be added to the `<div>` that wraps the `<iframe>` *and* will remove the inline styles. Note that you will need to call the `id` as a named argument as well. You can also give the vimeo video a descriptive title with `title`.

```go
{{</* vimeo id="146022717" class="my-vimeo-wrapper-class" title="My vimeo video" */>}}
```
{{% /note %}}

### youtube

{{% note %}}
To override Hugo's embedded `youtube` shortcode, copy the [source code] to a file with the same name in the layouts/shortcodes directory.

[source code]: {{% eturl youtube %}}
{{% /note %}}

To display a YouTube video with this URL:

```text
https://www.youtube.com/watch?v=0RKpf3rK57I
```

Include this in your Markdown:

```text
{{</* youtube 0RKpf3rK57I */>}}
```

Rendered:

{{< youtube 0RKpf3rK57I >}}

The youtube shortcode accepts these named arguments:

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

## Privacy configuration

To learn how to configure your Hugo site to meet the new EU privacy regulation, see [privacy protections].

## Create custom shortcodes

To learn more about creating custom shortcodes, see the [shortcode template documentation].

[privacy protections]: /about/privacy/
[partials]: /templates/partial/
[quickstart]: /getting-started/quick-start/
[sctemps]: /templates/shortcode/
[shortcode template documentation]: /templates/shortcode/
[Vimeo]: https://vimeo.com/
[YouTube Videos]: https://www.youtube.com/
