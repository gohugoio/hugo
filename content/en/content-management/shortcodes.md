---
title: Shortcodes
description: Shortcodes are simple snippets inside your content files calling built-in or custom templates.
categories: [content management]
keywords: [markdown,content,shortcodes]
menu:
  docs:
    parent: content-management
    weight: 100
toc: true
weight: 100
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

### Shortcodes with markdown

In Hugo `0.55` we changed how the `%` delimiter works. Shortcodes using the `%` as the outer-most delimiter will now be fully rendered when sent to the content renderer. They can be part of the generated table of contents, footnotes, etc.

If you want the old behavior, you can put the following line in the start of your shortcode template:

```go-html-template
{{ $_hugo_config := `{ "version": 1 }` }}
```

### Shortcodes without markdown

The `<` character indicates that the shortcode's inner content does *not* need further rendering. Often shortcodes without Markdown include internal HTML:

```go-html-template
{{</* myshortcode */>}}<p>Hello <strong>World!</strong></p>{{</* /myshortcode */>}}
```

### Nested shortcodes

You can call shortcodes within other shortcodes by creating your own templates that leverage the `.Parent` variable. `.Parent` allows you to check the context in which the shortcode is being called. See [Shortcode templates][sctemps].

## Use Hugo's built-in shortcodes

Hugo ships with a set of predefined shortcodes that represent very common usage. These shortcodes are provided for author convenience and to keep your Markdown content clean.

### `figure`

`figure` is an extension of the image syntax in Markdown, which does not provide a shorthand for the more semantic [HTML5 `<figure>` element][figureelement].

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
: Image caption.  Markdown within the value of `caption` will be rendered.

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

#### Example `figure` input

{{< code file="figure-input-example.md" >}}
{{</* figure src="elephant.jpg" title="An elephant at sunset" */>}}
{{< /code >}}

#### Example `figure` output

```html
<figure>
  <img src="elephant.jpg">
  <figcaption>An elephant at sunset</figcaption>
</figure>
```

### `gist`

To display a GitHub [gist] with this URL:

[gist]: https://docs.github.com/en/get-started/writing-on-github/editing-and-sharing-content-with-gists

```text
https://gist.github.com/user/50a7482715eac222e230d1e64dd9a89b
```

Include this in your markdown:

```text
{{</* gist user 50a7482715eac222e230d1e64dd9a89b */>}}
```

This will display all files in the gist alphabetically by file name.

{{< gist jmooring 50a7482715eac222e230d1e64dd9a89b >}}

To display a specific file within the gist:

```text
{{</* gist user 50a7482715eac222e230d1e64dd9a89b 1-template.html */>}}
```

Rendered:

{{< gist jmooring 50a7482715eac222e230d1e64dd9a89b 1-template.html >}}

### `highlight`

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

[highlighting options]: /functions/highlight/

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

### `instagram`

The `instagram` shortcode uses Facebook's **oEmbed Read** feature. The  Facebook [developer documentation] states:

- This permission or feature requires successful completion of the App Review process before your app can access live data. [Learn More]
- This permission or feature is only available with business verification. You may also need to sign additional contracts before your app can access data. [Learn More Here]

[developer documentation]: https://developers.facebook.com/docs/features-reference/oembed-read
[Learn More]: https://developers.facebook.com/docs/app-review
[Learn More Here]: https://developers.facebook.com/docs/development/release/business-verification

You must obtain an Access Token to use the `instagram` shortcode.

If your site configuration is private:

{{< code-toggle file="hugo" copy=false >}}
[services.instagram]
accessToken = 'xxx'
{{< /code-toggle >}}

If your site configuration is _not_ private, set the Access Token with an environment variable:

```text
HUGO_SERVICES_INSTAGRAM_ACCESSTOKEN=xxx hugo --gc --minify
```

{{% note %}}
If you are using a Client Access Token, you must combine the Access Token with your App ID using a pipe symbol (`APPID|ACCESSTOKEN`).
{{% /note %}}

To display an Instagram post with this URL:

```text
https://www.instagram.com/p/BWNjjyYFxVx/
```

Include this in your markdown:

```text
{{</* instagram BWNjjyYFxVx */>}}
```

### `param`

Gets a value from the current `Page's` parameters set in front matter, with a fallback to the site parameter value. It will log an `ERROR` if the parameter with the given key could not be found in either.

```bash
{{</* param testparam */>}}
```

Since `testparam` is a parameter defined in front matter of this page with the value `Hugo Rocks!`, the above will print:

{{< param testparam >}}

To access deeply nested parameters, use "dot syntax", e.g:

```bash
{{</* param "my.nested.param" */>}}
```

### `ref` and `relref`

These shortcodes will look up the pages by their relative path (e.g., `blog/post.md`) or their logical name (`post.md`) and return the permalink (`ref`) or relative permalink (`relref`) for the found page.

`ref` and `relref` also make it possible to make fragmentary links that work for the header links generated by Hugo.

{{% note %}}
Read a more extensive description of `ref` and `relref` in the [cross references](/content-management/cross-references/) documentation.
{{% /note %}}

`ref` and `relref` take exactly one required parameter of _reference_, quoted and in position `0`.

#### Example `ref` and `relref` input

```go-html-template
[Neat]({{</* ref "blog/neat.md" */>}})
[Who]({{</* relref "about.md#who" */>}})
```

#### Example `ref` and `relref` output

Assuming that standard Hugo pretty URLs are turned on.

```html
<a href="https://example.com/blog/neat">Neat</a>
<a href="/about/#who">Who</a>
```

### `tweet`

To display a Twitter post with this URL:

```txt
https://twitter.com/SanDiegoZoo/status/1453110110599868418
```

Include this in your markdown:

```text
{{</* tweet user="SanDiegoZoo" id="1453110110599868418" */>}}
```

Rendered:

{{< tweet user="SanDiegoZoo" id="1453110110599868418" >}}

### `vimeo`

To display a Vimeo video with this URL:

```text
https://vimeo.com/channels/staffpicks/55073825
```

Include this in your markdown:

```text
{{</* vimeo 55073825 */>}}
```

Rendered:

{{< vimeo 55073825 >}}

{{% note %}}
If you want to further customize the visual styling of the YouTube or Vimeo output, add a `class` named parameter when calling the shortcode. The new `class` will be added to the `<div>` that wraps the `<iframe>` *and* will remove the inline styles. Note that you will need to call the `id` as a named parameter as well. You can also give the vimeo video a descriptive title with `title`.

```go
{{</* vimeo id="146022717" class="my-vimeo-wrapper-class" title="My vimeo video" */>}}
```
{{% /note %}}

### `youtube`

The `youtube` shortcode embeds a responsive video player for [YouTube videos]. Only the ID of the video is required, e.g.:

```txt
https://www.youtube.com/watch?v=w7Ft2ymGmfc
```

#### Example `youtube` input

Copy the YouTube video ID that follows `v=` in the video's URL and pass it to the `youtube` shortcode:

{{< code file="example-youtube-input.md" >}}
{{</* youtube w7Ft2ymGmfc */>}}
{{< /code >}}

Furthermore, you can automatically start playback of the embedded video by setting the `autoplay` parameter to `true`. Remember that you can't mix named and unnamed parameters, so you'll need to assign the yet unnamed video ID to the parameter `id`:


{{< code file="example-youtube-input-with-autoplay.md" >}}
{{</* youtube id="w7Ft2ymGmfc" autoplay="true" */>}}
{{< /code >}}

For [accessibility reasons](https://dequeuniversity.com/tips/provide-iframe-titles), it's best to provide a title for your YouTube video.  You  can do this using the shortcode by providing a `title` parameter. If no title is provided, a default of "YouTube Video" will be used.

{{< code file="example-youtube-input-with-title.md" >}}
{{</* youtube id="w7Ft2ymGmfc" title="A New Hugo Site in Under Two Minutes" */>}}
{{< /code >}}

#### Example `youtube` output

Using the preceding `youtube` example, the following HTML will be added to your rendered website's markup:

{{< code file="example-youtube-output.html" >}}
{{< youtube id="w7Ft2ymGmfc" autoplay="true" >}}
{{< /code >}}

#### Example `youtube` display

Using the preceding `youtube` example (without `autoplay="true"`), the following simulates the displayed experience for visitors to your website. Naturally, the final display will be contingent on your style sheets and surrounding markup. The video is also include in the [Quick Start of the Hugo documentation][quickstart].

{{< youtube w7Ft2ymGmfc >}}

## Privacy configuration

To learn how to configure your Hugo site to meet the new EU privacy regulation, see [Hugo and the GDPR].

## Create custom shortcodes

To learn more about creating custom shortcodes, see the [shortcode template documentation].

[`figure` shortcode]: #figure
[contentmanagementsection]: /content-management/formats/
[examplegist]: https://gist.github.com/spf13/7896402
[figureelement]: https://html5doctor.com/the-figure-figcaption-elements/ "An article from HTML5 doctor discussing the fig and figcaption elements."
[Hugo and the GDPR]: /about/hugo-and-gdpr/
[Instagram]: https://www.instagram.com/
[pagevariables]: /variables/page/
[partials]: /templates/partials/
[quickstart]: /getting-started/quick-start/
[sctemps]: /templates/shortcode-templates/
[scvars]: /variables/shortcodes/
[shortcode template documentation]: /templates/shortcode-templates/
[templatessection]: /templates/
[Vimeo]: https://vimeo.com/
[YouTube Videos]: https://www.youtube.com/
[YouTube Input shortcode]: #youtube
