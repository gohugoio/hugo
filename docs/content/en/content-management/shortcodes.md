---
title: Shortcodes
linktitle:
description: Shortcodes are simple snippets inside your content files calling built-in or custom templates.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2019-11-07
menu:
  docs:
    parent: "content-management"
    weight: 35
weight: 35	#rem
categories: [content management]
keywords: [markdown,content,shortcodes]
draft: false
aliases: [/extras/shortcodes/]
testparam: "Hugo Rocks!"
toc: true
---

## What a Shortcode is

Hugo loves Markdown because of its simple content format, but there are times when Markdown falls short. Often, content authors are forced to add raw HTML (e.g., video `<iframes>`) to Markdown content. We think this contradicts the beautiful simplicity of Markdown's syntax.

Hugo created **shortcodes** to circumvent these limitations.

A shortcode is a simple snippet inside a content file that Hugo will render using a predefined template. Note that shortcodes will not work in template files. If you need the type of drop-in functionality that shortcodes provide but in a template, you most likely want a [partial template][partials] instead.

In addition to cleaner Markdown, shortcodes can be updated any time to reflect new classes, techniques, or standards. At the point of site generation, Hugo shortcodes will easily merge in your changes. You avoid a possibly complicated search and replace operation.

## Use Shortcodes

{{< youtube 2xkNJL4gJ9E >}}

In your content files, a shortcode can be called by calling `{{%/* shortcodename parameters */%}}`. Shortcode parameters are space delimited, and parameters with internal spaces can be quoted.

The first word in the shortcode declaration is always the name of the shortcode. Parameters follow the name. Depending upon how the shortcode is defined, the parameters may be named, positional, or both, although you can't mix parameter types in a single call. The format for named parameters models that of HTML with the format `name="value"`.

Some shortcodes use or require closing shortcodes. Again like HTML, the opening and closing shortcodes match (name only) with the closing declaration, which is prepended with a slash.

Here are two examples of paired shortcodes:

```
{{%/* mdshortcode */%}}Stuff to `process` in the *center*.{{%/* /mdshortcode */%}}
```

```
{{</* highlight go */>}} A bunch of code here {{</* /highlight */>}}
```

The examples above use two different delimiters, the difference being the `%` character in the first and the `<>` characters in the second.

### Shortcodes with raw string parameters

{{< new-in "0.64.1" >}}

You can pass multiple lines as parameters to a shortcode by using raw string literals:

```
{{</*  myshortcode `This is some <b>HTML</b>,
and a new line with a "quoted string".` */>}}
```

### Shortcodes with Markdown

In Hugo `0.55` we changed how the `%` delimiter works. Shortcodes using the `%` as the outer-most delimiter will now be fully rendered when sent to the content renderer (e.g. Blackfriday for Markdown), meaning they can be part of the generated table of contents, footnotes, etc.

If you want the old behavior, you can put the following line in the start of your shortcode template:

```
{{ $_hugo_config := `{ "version": 1 }` }}
```


### Shortcodes Without Markdown

The `<` character indicates that the shortcode's inner content does *not* need further rendering. Often shortcodes without markdown include internal HTML:

```
{{</* myshortcode */>}}<p>Hello <strong>World!</strong></p>{{</* /myshortcode */>}}
```

### Nested Shortcodes

You can call shortcodes within other shortcodes by creating your own templates that leverage the `.Parent` variable. `.Parent` allows you to check the context in which the shortcode is being called. See [Shortcode templates][sctemps].

## Use Hugo's Built-in Shortcodes

Hugo ships with a set of predefined shortcodes that represent very common usage. These shortcodes are provided for author convenience and to keep your markdown content clean.

### `figure`

`figure` is an extension of the image syntax in markdown, which does not provide a shorthand for the more semantic [HTML5 `<figure>` element][figureelement].

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

attr
: Image attribution text. Markdown within the value of `attr` will be rendered.

attrlink
: If the attribution text needs to be hyperlinked, URL of the destination.

#### Example `figure` Input

{{< code file="figure-input-example.md" >}}
{{</* figure src="/media/spf13.jpg" title="Steve Francia" */>}}
{{< /code >}}

#### Example `figure` Output

{{< output file="figure-output-example.html" >}}
<figure>
  <img src="/media/spf13.jpg"  />
  <figcaption>
      <h4>Steve Francia</h4>
  </figcaption>
</figure>
{{< /output >}}

### `gist`

Bloggers often want to include GitHub gists when writing posts. Let's suppose we want to use the [gist at the following url][examplegist]:

```
https://gist.github.com/spf13/7896402
```

We can embed the gist in our content via username and gist ID pulled from the URL:

```
{{</* gist spf13 7896402 */>}}
```

#### Example `gist` Input

If the gist contains several files and you want to quote just one of them, you can pass the filename (quoted) as an optional third argument:

{{< code file="gist-input.md" >}}
{{</* gist spf13 7896402 "img.html" */>}}
{{< /code >}}

#### Example `gist` Output

{{< output file="gist-output.html" >}}
{{< gist spf13 7896402 >}}
{{< /output >}}

#### Example `gist` Display

To demonstrate the remarkably efficiency of Hugo's shortcode feature, we have embedded the `spf13` `gist` example in this page. The following simulates the experience for visitors to your website. Naturally, the final display will be contingent on your stylesheets and surrounding markup.

{{< gist spf13 7896402 >}}

### `highlight`

This shortcode will convert the source code provided into syntax-highlighted HTML. Read more on [highlighting](/tools/syntax-highlighting/). `highlight` takes exactly one required `language` parameter and requires a closing shortcode.

#### Example `highlight` Input

{{< code file="content/tutorials/learn-html.md" >}}
{{</* highlight html */>}}
<section id="main">
  <div>
   <h1 id="title">{{ .Title }}</h1>
    {{ range .Pages }}
        {{ .Render "summary"}}
    {{ end }}
  </div>
</section>
{{</* /highlight */>}}
{{< /code >}}

#### Example `highlight` Output

The `highlight` shortcode example above would produce the following HTML when the site is rendered:

{{< output file="tutorials/learn-html/index.html" >}}
<span style="color: #f92672">&lt;section</span> <span style="color: #a6e22e">id=</span><span style="color: #e6db74">&quot;main&quot;</span><span style="color: #f92672">&gt;</span>
  <span style="color: #f92672">&lt;div&gt;</span>
   <span style="color: #f92672">&lt;h1</span> <span style="color: #a6e22e">id=</span><span style="color: #e6db74">&quot;title&quot;</span><span style="color: #f92672">&gt;</span>{{ .Title }}<span style="color: #f92672">&lt;/h1&gt;</span>
    {{ range .Pages }}
        {{ .Render &quot;summary&quot;}}
    {{ end }}
  <span style="color: #f92672">&lt;/div&gt;</span>
<span style="color: #f92672">&lt;/section&gt;</span>
{{< /output >}}

{{% note "More on Syntax Highlighting" %}}
To see even more options for adding syntax-highlighted code blocks to your website, see [Syntax Highlighting in Developer Tools](/tools/syntax-highlighting/).
{{% /note %}}

### `instagram`

If you'd like to embed a photo from [Instagram][], you only need the photo's ID. You can discern an Instagram photo ID from the URL:

```
https://www.instagram.com/p/BWNjjyYFxVx/
```

#### Example `instagram` Input

{{< code file="instagram-input.md" >}}
{{</* instagram BWNjjyYFxVx */>}}
{{< /code >}}

You also have the option to hide the caption:

{{< code file="instagram-input-hide-caption.md" >}}
{{</* instagram BWNjjyYFxVx hidecaption */>}}
{{< /code >}}

#### Example `instagram` Output

By adding the preceding `hidecaption` example, the following HTML will be added to your rendered website's markup:

{{< output file="instagram-hide-caption-output.html" >}}
{{< instagram BWNjjyYFxVx hidecaption >}}
{{< /output >}}

#### Example `instagram` Display

Using the preceding `instagram` with `hidecaption` example above, the following simulates the displayed experience for visitors to your website. Naturally, the final display will be contingent on your stylesheets and surrounding markup.

{{< instagram BWNjjyYFxVx hidecaption >}}


### `param`

Gets a value from the current `Page's` params set in front matter, with a fall back to the site param value. It will log an `ERROR` if the param with the given key could not be found in either.

```bash
{{</* param testparam */>}}
```

Since `testparam` is a param defined in front matter of this page with the value `Hugo Rocks!`, the above will print:

{{< param testparam >}}

To access deeply nested params, use "dot syntax", e.g:

```bash
{{</* param "my.nested.param" */>}}
```

### `ref` and `relref`

These shortcodes will look up the pages by their relative path (e.g., `blog/post.md`) or their logical name (`post.md`) and return the permalink (`ref`) or relative permalink (`relref`) for the found page.

`ref` and `relref` also make it possible to make fragmentary links that work for the header links generated by Hugo.

{{% note "More on Cross References" %}}
Read a more extensive description of `ref` and `relref` in the [cross references](/content-management/cross-references/) documentation.
{{% /note %}}

`ref` and `relref` take exactly one required parameter of _reference_, quoted and in position `0`.

#### Example `ref` and `relref` Input

```
[Neat]({{</* ref "blog/neat.md" */>}})
[Who]({{</* relref "about.md#who" */>}})
```

#### Example `ref` and `relref` Output

Assuming that standard Hugo pretty URLs are turned on.

```
<a href="/blog/neat">Neat</a>
<a href="/about/#who:c28654c202e73453784cfd2c5ab356c0">Who</a>
```

### `tweet`

You want to include a single tweet into your blog post? Everything you need is the URL of the tweet:

```
https://twitter.com/spf13/status/877500564405444608
```

#### Example `tweet` Input

Pass the tweet's ID from the URL as a parameter to the `tweet` shortcode:

{{< code file="example-tweet-input.md" >}}
{{</* tweet 877500564405444608 */>}}
{{< /code >}}

#### Example `tweet` Output

Using the preceding `tweet` example, the following HTML will be added to your rendered website's markup:

{{< output file="example-tweet-output.html" >}}
{{< tweet 877500564405444608 >}}
{{< /output >}}

#### Example `tweet` Display

Using the preceding `tweet` example, the following simulates the displayed experience for visitors to your website. Naturally, the final display will be contingent on your stylesheets and surrounding markup.

{{< tweet 877500564405444608 >}}

### `vimeo`

Adding a video from [Vimeo][] is equivalent to the [YouTube Input shortcode][].

```
https://vimeo.com/channels/staffpicks/146022717
```

#### Example `vimeo` Input

Extract the ID from the video's URL and pass it to the `vimeo` shortcode:

{{< code file="example-vimeo-input.md" >}}
{{</* vimeo 146022717 */>}}
{{< /code >}}

#### Example `vimeo` Output

Using the preceding `vimeo` example, the following HTML will be added to your rendered website's markup:

{{< output file="example-vimeo-output.html" >}}
{{< vimeo 146022717 >}}
{{< /output >}}

{{% tip %}}
If you want to further customize the visual styling of the YouTube or Vimeo output, add a `class` named parameter when calling the shortcode. The new `class` will be added to the `<div>` that wraps the `<iframe>` *and* will remove the inline styles. Note that you will need to call the `id` as a named parameter as well. You can also give the vimeo video a descriptive title with `title`. 

```
{{</* vimeo id="146022717" class="my-vimeo-wrapper-class" title="My vimeo video" */>}}
```
{{% /tip %}}

#### Example `vimeo` Display

Using the preceding `vimeo` example, the following simulates the displayed experience for visitors to your website. Naturally, the final display will be contingent on your stylesheets and surrounding markup.

{{< vimeo 146022717 >}}

### `youtube`

The `youtube` shortcode embeds a responsive video player for [YouTube videos][]. Only the ID of the video is required, e.g.:

```
https://www.youtube.com/watch?v=w7Ft2ymGmfc
```


#### Example `youtube` Input

Copy the YouTube video ID that follows `v=` in the video's URL and pass it to the `youtube` shortcode:

{{< code file="example-youtube-input.md" >}}
{{</* youtube w7Ft2ymGmfc */>}}
{{< /code >}}

Furthermore, you can automatically start playback of the embedded video by setting the `autoplay` parameter to `true`. Remember that you can't mix named and unnamed parameters, so you'll need to assign the yet unnamed video id to the parameter `id`:


{{< code file="example-youtube-input-with-autoplay.md" >}}
{{</* youtube id="w7Ft2ymGmfc" autoplay="true" */>}}
{{< /code >}}

#### Example `youtube` Output

Using the preceding `youtube` example, the following HTML will be added to your rendered website's markup:

{{< code file="example-youtube-output.html" >}}
{{< youtube id="w7Ft2ymGmfc" autoplay="true" >}}
{{< /code >}}

#### Example `youtube` Display

Using the preceding `youtube` example (without `autoplay="true"`), the following simulates the displayed experience for visitors to your website. Naturally, the final display will be contingent on your stylesheets and surrounding markup. The video is also include in the [Quick Start of the Hugo documentation][quickstart].

{{< youtube w7Ft2ymGmfc >}}

## Privacy Config

To learn how to configure your Hugo site to meet the new EU privacy regulation, see [Hugo and the GDPR][].

## Create Custom Shortcodes

To learn more about creating custom shortcodes, see the [shortcode template documentation][].

[`figure` shortcode]: #figure
[contentmanagementsection]: /content-management/formats/
[examplegist]: https://gist.github.com/spf13/7896402
[figureelement]: http://html5doctor.com/the-figure-figcaption-elements/ "An article from HTML5 doctor discussing the fig and figcaption elements."
[Hugo and the GDPR]: /about/hugo-and-gdpr/
[Instagram]: https://www.instagram.com/
[pagevariables]: /variables/page/
[partials]: /templates/partials/
[Pygments]: http://pygments.org/
[quickstart]: /getting-started/quick-start/
[sctemps]: /templates/shortcode-templates/
[scvars]: /variables/shortcodes/
[shortcode template documentation]: /templates/shortcode-templates/
[templatessection]: /templates/
[Vimeo]: https://vimeo.com/
[YouTube Videos]: https://www.youtube.com/
[YouTube Input shortcode]: #youtube
