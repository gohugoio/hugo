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

In your content files, a shortcode can be called by calling `{{%/* shortcodename parameters */%}}`. Shortcode parameters are space delimited, and parameters with internal spaces can be quoted.

The first word in the shortcode declaration is always the name of the shortcode. Parameters follow the name. Depending upon how the shortcode is defined, the parameters may be named, positional, or both, although you can't mix parameter types in a single call. The format for named parameters models that of HTML with the format `name="value"`.

Some shortcodes use or require closing shortcodes. Again like HTML, the opening and closing shortcodes match (name only) with the closing declaration, which is prepended with a slash.

Here are two examples of paired shortcodes:

```go-html-template
{{%/* mdshortcode */%}}Stuff to `process` in the *center*.{{%/* /mdshortcode */%}}
```

```go-html-template
{{</* highlight go */>}} A bunch of code here {{</* /highlight */>}}
```

The examples above use two different delimiters, the difference being the `%` character in the first and the `<>` characters in the second.

### Shortcodes with raw string parameters

You can pass multiple lines as parameters to a shortcode by using raw string literals:

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

The `figure` shortcode can use the following named parameters:

src
: URL of the image to be displayed.

link
: If the image needs to be hyperlinked, URL of the destination.

target
: Optional `target` attribute for the URL if `link` parameter is set.

rel
: Optional `rel` attribute for the URL if `link` parameter is set.

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
https://twitter.com/SanDiegoZoo/status/1453110110599868418
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
If you want to further customize the visual styling of the YouTube or Vimeo output, add a `class` parameter when calling the shortcode. The new `class` will be added to the `<div>` that wraps the `<iframe>` *and* will remove the inline styles. Note that you will need to call the `id` as a named parameter as well. You can also give the vimeo video a descriptive title with `title`.

```go
{{</* vimeo id="146022717" class="my-vimeo-wrapper-class" title="My vimeo video" */>}}
```
{{% /note %}}

### youtube

{{% note %}}
To override Hugo's embedded `vimeo` shortcode, copy the [source code] to a file with the same name in the layouts/shortcodes directory.

[source code]: {{% eturl vimeo %}}
{{% /note %}}

The `youtube` shortcode embeds a responsive video player for [YouTube videos]. Only the ID of the video is required, e.g.:

```txt
https://www.youtube.com/watch?v=w7Ft2ymGmfc
```

#### Example `youtube` input

Copy the YouTube video ID that follows `v=` in the video's URL and pass it to the `youtube` shortcode:

{{< code file=example-youtube-input.md >}}
{{</* youtube w7Ft2ymGmfc */>}}
{{< /code >}}

Furthermore, you can automatically start playback of the embedded video by setting the `autoplay` parameter to `true`. Remember that you can't mix named and unnamed parameters, so you'll need to assign the yet unnamed video ID to the parameter `id`:

{{< code file=example-youtube-input-with-autoplay.md >}}
{{</* youtube id="w7Ft2ymGmfc" autoplay="true" */>}}
{{< /code >}}

For [accessibility reasons](https://dequeuniversity.com/tips/provide-iframe-titles), it's best to provide a title for your YouTube video. You  can do this using the shortcode by providing a `title` parameter. If no title is provided, a default of "YouTube Video" will be used.

{{< code file=example-youtube-input-with-title.md >}}
{{</* youtube id="w7Ft2ymGmfc" title="A New Hugo Site in Under Two Minutes" */>}}
{{< /code >}}

#### Example `youtube` output

Using the preceding `youtube` example, the following HTML will be added to your rendered website's markup:

{{< code file=example-youtube-output.html >}}
{{< youtube id="w7Ft2ymGmfc" autoplay="true" >}}
{{< /code >}}

#### Example `youtube` display

Using the preceding `youtube` example (without `autoplay="true"`), the following simulates the displayed experience for visitors to your website. Naturally, the final display will be contingent on your style sheets and surrounding markup. The video is also include in the [Quick Start of the Hugo documentation][quickstart].

{{< youtube w7Ft2ymGmfc >}}

## Privacy configuration

To learn how to configure your Hugo site to meet the new EU privacy regulation, see [Hugo and the GDPR].

## Create custom shortcodes

To learn more about creating custom shortcodes, see the [shortcode template documentation].

[Hugo and the GDPR]: /about/hugo-and-gdpr/
[partials]: /templates/partials/
[quickstart]: /getting-started/quick-start/
[sctemps]: /templates/shortcode-templates/
[shortcode template documentation]: /templates/shortcode-templates/
[Vimeo]: https://vimeo.com/
[YouTube Videos]: https://www.youtube.com/
