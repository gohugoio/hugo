---
title: Shortcodes
linktitle:
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
weight: 25
categories: [content management]
tags: [markdown,content,shortcodes]
draft: false
slug:
aliases: [/extras/shortcodes/]
toc: true
---

Hugo loves Markdown because of its simple content format, but there are times when markdown falls short. Often, content authors fall back on adding raw HTML (e.g., video `<iframes>`) to markdown content. We think blocks of raw HTML contradict the beautiful simplicity of markdown's syntax.

Hugo created **shortcodes** to circumvent these limitations.

A shortcode is a simple snippet inside a content file that Hugo will render using a predefined template. Note that shortcodes will not work in template files---if you need the type of drop-in functionality that shortcodes provide but in a template, you most likely want a [partial template][partialtemplatespage] instead.

In addition to cleaner markdown, shortcodes can be updated any time to reflect new classes, techniques, or standards. At the point of site generation, Hugo shortcodes will easily merge in your changes. You avoid a possibly complicated search and replace operation.

## Using Shortcodes

In your content files, a shortcode can be called by using the `{{%/* shortcodename parameters*/%}}` form. Shortcode parameters are space delimited. Parameters with internal spaces can be quoted.

The first word is always the name of the shortcode. Parameters follow the name.
Depending upon how the shortcode is defined, the parameters may be named,
positional or both (although you can't mixed parameter types in a single call).
The format for named parameters models that of HTML with the format
`name="value"`.

Some shortcodes use or require closing shortcodes. Like HTML, the opening and closing shortcodes match (name only) with the closing declaration prepended with a slash.

Example of a paired shortcode:

```golang
{{</* highlight go */>}} A bunch of code here {{</* /highlight */>}}
```

The examples above use two different delimiters, the difference being the `%` and the `<` character:

### Shortcodes with Markdown

The `%` characters indicates that the shortcode's inner content (`.Inner`) needs further processing by the page's rendering processor (i.e. Markdown), which would be needed to convert `**World**` to `<strong>World</strong>` in the following example:

```golang
{{%/* myshortcode */%}}Hello **World!**{{%/* /myshortcode */%}}
```

### Shortcodes without Markdown

The `<` character indicates that the shortcode's inner content doesn't need any further rendering, this will typically be pure HTML:

```golang
{{</* myshortcode */>}}<p>Hello <strong>World!</strong></p>{{</* /myshortcode */>}}
```

## Using the Built-in Shortcodes

Hugo ships with a set of predefined shortcodes the represent very common usage that would otherwise muddy your content with unnecessary markup.

### `figure`

`figure` is an extension of the image syntax in Markdown, which does not provide a shorthand for the more semantic [HTML5 `<figure>` element][figurelement].

`figure` can use the following named parameters:

* `src`
* `link`
* `title`
* `caption`
* `class`
* `attr` (i.e., attribution)
* `attrlink`
* `alt`

#### Example `figure` Input

{{% input "figure-input-example.md" %}}
```markdown
{{</* figure src="/media/spf13.jpg" title="Steve Francia" */>}}
```
{{% /input %}}

#### Example `figure` Output

{{% output "figure-output-example.html" %}}
```html
<figure>
  <img src="/media/spf13.jpg"  />
  <figcaption>
      <h4>Steve Francia</h4>
  </figcaption>
</figure>
```
{{% /output %}}

### `gist`

Bloggers often want to include GitHub gists when writing posts. Let's supposed we want to use the following [gist][examplegist]:

```html
https://gist.github.com/spf13/7896402
```

We can embed the gist in our content via username and gist ID pulled from the URL:

```golang
{{</* gist spf13 7896402 */>}}
```

If the gist contains several files and you want to quote just one of them, you can pass the filename (quoted) as an optional third argument:

```golang
{{</* gist spf13 7896402 "img.html" */>}}
```

### `highlight`

This shortcode will convert the source code provided into syntax highlighted
HTML. Read more on [highlighting](/extras/highlighting/). `highlight` takes exactly one required parameter of _language_ and requires a closing shortcode.

#### Example `highlight` Input

{{% input "highlight-shortcode.md" %}}
```golang
{{</* highlight html */>}}
<section id="main">
  <div>
   <h1 id="title">{{ .Title }}</h1>
    {{ range .Data.Pages }}
        {{ .Render "summary"}}
    {{ end }}
  </div>
</section>
{{</* /highlight */>}}
```
{{% /input %}}

#### Example `highlight` Output

{{% output "syntax-highlighted.html" %}}
```html
<span style="color: #f92672">&lt;section</span> <span style="color: #a6e22e">id=</span><span style="color: #e6db74">&quot;main&quot;</span><span style="color: #f92672">&gt;</span>
  <span style="color: #f92672">&lt;div&gt;</span>
   <span style="color: #f92672">&lt;h1</span> <span style="color: #a6e22e">id=</span><span style="color: #e6db74">&quot;title&quot;</span><span style="color: #f92672">&gt;</span>{{ .Title }}<span style="color: #f92672">&lt;/h1&gt;</span>
    {{ range .Data.Pages }}
        {{ .Render &quot;summary&quot;}}
    {{ end }}
  <span style="color: #f92672">&lt;/div&gt;</span>
<span style="color: #f92672">&lt;/section&gt;</span>
```
{{% /output %}}

### `instagram`

If you'd like to embed a photo from [Instagram][], all you need is photo ID from the URL:

```html
https://www.instagram.com/p/BMokmydjG-M/
```

Pass it to the shortcode:

```golang
{{</* instagram BMokmydjG-M */>}}
```

You also have the option to hide the caption:

```golang
{{</* instagram BMokmydjG-M hidecaption */>}}
```

### `ref` and `relref`

These shortcodes will look up the pages by their relative path (e.g., `blog/post.md`) or their logical name (`post.md`) and return the permalink (`ref`) or relative permalink (`relref`) for the found page.

`ref` and `relref` also make it possible to make fragmentary links that work for the header links generated by Hugo.

{{% note "More on Cross References" %}}
Read a more extensive description of `ref` and `relref` in the [cross references](/content-management/cross-references/) documentation.
{{% /note %}}

`ref` and `relref` take exactly one required parameter of _reference_, quoted and in position `0`.

#### Example `ref` and `relref` Input

```golang
[Neat]({{</* ref "blog/neat.md" */>}})
[Who]({{</* relref "about.md#who" */>}})
```

#### Example `ref` and `relref` Output

Assuming that standard Hugo pretty URLs are turned on.

```html
<a href="/blog/neat">Neat</a>
<a href="/about/#who:c28654c202e73453784cfd2c5ab356c0">Who</a>
```

### `speakerdeck`

To embed slides from [Speaker Deck][], click on "&lt;&#8239;/&gt;&nbsp;Embed" (under Share right next to the template on Speaker Deck) and copy the URL, e.g.:

```html
<script async class="speakerdeck-embed" data-id="4e8126e72d853c0060001f97" data-ratio="1.33333333333333" src="//speakerdeck.com/assets/embed.js"></script>
```

Extract the value from the field `data-id` and pass it to the shortcode:

{{% input "speakerdeck-example-input.md" %}}
```golang
{{</* speakerdeck 4e8126e72d853c0060001f97 */>}}
```
{{% /input %}}

### `tweet`

You want to include a single tweet into your blog post? Everything you need is the URL of the tweet. For example, let's say you want to include the following tweet from `https://twitter.com/spf13/status/666616452582129664`. Pass the tweet's ID from the URL as parameter to the shortcode as shown below:

#### Example `tweet` Input

```golang
{{</* tweet 666616452582129664 */>}}
```

### `vimeo`

Adding a video from [Vimeo][] is equivalent to the YouTube shortcode above. Extract the ID from the URL, e.g.:

* https://vimeo.com/channels/staffpicks/146022717

and pass it to the shortcode:

```golang
{{</* vimeo 146022717 */>}}
```

### `youtube`

This shortcode embeds a responsive video player for [YouTube videos][]. Only the ID of the video is required, e.g.:

* https://www.youtube.com/watch?v=w7Ft2ymGmfc

Copy the YouTube video ID that follows `v=` in the video's URL and pass it to the `youtube` shortcode:

```golang
{{</* youtube w7Ft2ymGmfc */>}}
```

Furthermore, you can autostart the embedded video by setting the `autostart` parameter to true. Remember that you can't mix named an unamed parameters. Assign the yet unamed video id to the parameter `id` like below too.

```golang
{{</* youtube id="w7Ft2ymGmfc" autoplay="true" */>}}
```

## Creating Custom Shortcodes

To learn more about creating your own shortcode templates, see the [shortcode template documentation][].

[contentmanagementsection]: /content-management/supported-content-formats/
[examplegist]: https://gist.github.com/spf13/7896402
[figureelement]: http://html5doctor.com/the-figure-figcaption-elements/ "An article from HTML5 doctor discussing the fig and figcaption elements."
[`figure` shortcode]: #figure
[Instagram]: https://www.instagram.com/
[pagevariables]: /variables-and-params/page-variables/
[partialtemplatespage]: /templates/partials-templates/
[Pygments]: http://pygments.org/
[projectorganizationsection]: /project-organization/directory-structure/
[shortcode template documentation]: /templates/shortcode-templates/
[Speaker Deck]: https://speakerdeck.com/
[templatessection]: /templates/
[Vimeo]: https://vimeo.com/
[YouTube Videos]: https://www.youtube.com/