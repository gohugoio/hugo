---
title: Shortcodes
linktitle:
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
weight: 25
tags: [markdown,content]
categories: [content management]
draft: false
slug:
aliases: [/extras/shortcodes/]
toc: true
notes:
---

Hugo uses Markdown for its simple content format. However, there are a lot of things that Markdown doesn’t support well.

We are unwilling to accept being constrained by our simple format. Also unacceptable is writing raw HTML in our Markdown every time we want to include unsupported content such as a video. To do so is in complete opposition to the intent of using a bare-bones format for our content and utilizing templates to apply styling for display.

To avoid both of these limitations, Hugo created **shortcodes**.

A shortcode is a simple snippet inside a content file that Hugo will render using a predefined template. Note that shortcodes will not work in template files---if you need a functionality like that in a template, you most likely want a [partial template][] instead.

Another benefit is, you can update your shortcode with any related new classes or techniques, and upon generation, Hugo will easily merge in your changes. You avoid a possibly complicated search and replace operation.

## Using a shortcode

In your content files, a shortcode can be called by using the `{{%/* name parameters*/%}}` form. Shortcode parameters are space delimited. Parameters with spaces can be quoted.

The first word is always the name of the shortcode. Parameters follow the name.
Depending upon how the shortcode is defined, the parameters may be named,
positional or both (although you can't mixed parameter types in a single call).
The format for named parameters models that of HTML with the format
`name="value"`.

Some shortcodes use or require closing shortcodes. Like HTML, the opening and closing shortcodes match (name only), the closing being prepended with a slash.

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

## Built-in Shortcodes

Hugo ships with a set of predefined shortcodes.

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

### `figure`

`figure` is simply an extension of the image capabilities present with Markdown.`figure` provides the ability to add captions, CSS classes, alt text, links etc.

`figure` can use the following named parameters:

 * src
 * link
 * title
 * caption
 * class
 * attr (attribution)
 * attrlink
 * alt

#### Example `figure` Input

```markdown
{{</* figure src="/media/spf13.jpg" title="Steve Francia" */>}}
```

#### Example `figure` Output

    <figure>
        <img src="/media/spf13.jpg"  />
        <figcaption>
            <h4>Steve Francia</h4>
        </figcaption>
    </figure>

### `ref` and `relref`

These shortcodes will look up the pages by their relative path (e.g., `blog/post.md`) or their logical name (`post.md`) and return the permalink (`ref`) or relative permalink (`relref`) for the found page.

`ref` and `relref` also make it possible to make fragmentary links that work for the header links generated by Hugo.

{{% note "More on Cross References" %}}
Read a more extensive description of `ref` and `relref` in the [cross-references][] documentation.
{{% /note %}}

`ref` and `relref` take exactly one required parameter of _reference_.

#### Example `ref` and `relref` Input

```markdown
[Neat]({{</* ref "blog/neat.md" */>}})
[Who]({{</* relref "about.md#who" */>}})
```

#### Example `ref` and `relref` Output

Assuming that standard Hugo pretty URLs are turned on.

```markdown
<a href="/blog/neat">Neat</a>
<a href="/about/#who:c28654c202e73453784cfd2c5ab356c0">Who</a>
```

### `tweet` (Twitter)

You want to include a single tweet into your blog post? Everything you need is the URL of the tweet. For example, let's say you want to include the following tweet from `https://twitter.com/spf13/status/666616452582129664`. Pass the tweet's ID from the URL as parameter to the shortcode as shown below:

#### Example `tweet` Input

```markdown
{{</* tweet 666616452582129664 */>}}
```

### `youtube`

This shortcode embeds a responsive video player for [YouTube videos][]. Only the ID of the video is required, e.g.:

* https://www.youtube.com/watch?v=w7Ft2ymGmfc

Copy the YouTube video ID that follows `v=` in the video's URL and pass it to the `youtube` shortcode:

```markdown
{{</* youtube w7Ft2ymGmfc */>}}
```

Furthermore, you can autostart the embedded video by setting the `autostart` parameter to true. Remember that you can't mix named an unamed parameters. Assign the yet unamed video id to the parameter `id` like below too.

```markdown
{{</* youtube id="w7Ft2ymGmfc" autoplay="true" */>}}
```

### `vimeo`

Adding a video from [Vimeo][] is equivalent to the YouTube shortcode above. Extract the ID from the URL, e.g.:

* https://vimeo.com/channels/staffpicks/146022717

and pass it to the shortcode:

```golang
{{</* vimeo 146022717 */>}}
```

### `gist` (GitHub)

Including code snippets with GitHub gists while writing a tutorial is common situation bloggers face. With a given URL of the gist, e.g.:

* https://gist.github.com/spf13/7896402

pass the owner and the ID of the gist to the shortcode:

```
{{</* gist spf13 7896402 */>}}
```

If the gist contains several files and you want to quote just one of them, you can pass the filename (quoted) as an optional third argument:

```golang
{{</* gist spf13 7896402 "img.html" */>}}
```

### `speakerdeck`

To embed slides from [Speaker Deck][], click on "&lt;&#8239;/&gt;&nbsp;Embed" (under Share right next to the template on Speaker Deck) and copy the URL, e.g.:

    <script async class="speakerdeck-embed" data-id="4e8126e72d853c0060001f97" data-ratio="1.33333333333333" src="//speakerdeck.com/assets/embed.js"></script>

Extract the value from the field `data-id` and pass it to the shortcode:

    {{</* speakerdeck 4e8126e72d853c0060001f97 */>}}

### Instagram

If you'd like to embed photo from [Instagram](https://www.instagram.com/), all you need is photo ID from the URL, e. g.:

* https://www.instagram.com/p/BMokmydjG-M/

Pass it to the shortcode:

```golang
{{</* instagram BMokmydjG-M */>}}
```

Optionally, hide caption:

```golang
{{</* instagram BMokmydjG-M hidecaption */>}}
```

## Creating your own shortcodes

To create a shortcode, place a template in the `layouts/shortcodes` directory of your [source organization][]. The template name will be the name of the shortcode. In creating a shortcode, you can choose if the shortcode will use _positional parameters_, or _named parameters_, or _both_. A good rule of thumb is that if a shortcode has a single required value in the case of the `youtube` example below, then positional works very well. For more complex layouts with optional parameters, named parameters work best. Allowing both types of parameters is useful for complex layouts where you want to set default values that can be overridden.

### Accessing Parameters

To access a parameter in any shortcode, use the `.Get` method. Whether you pass a key (string) or a number to the `.Get` method depends on whether you are accessing a named or positional parameter, respectively.

To access a parameter by name, the `.Get` method followed by the named parameter as a quoted string. Named parameters are less terse but do not require that a content author be mindful of the order of parameters.

```golang
{{ .Get "class" }}
```

To access a parameter by position, the `.Get` method can be used, keeping in mind that the first positional parameter within the shortcode declaration starts at `0`:

```golang
{{ .Get 0 }}
```

`with` is great when the output depends on a parameter being set:

```golang
{{ with .Get "class"}} class="{{.}}"{{ end }}
```

`.Get` can also be used to check if a parameter has been provided. This is
most helpful when the condition depends on either of the values, or both:

```golang
{{ or .Get "title" | .Get "alt" | if }} alt="{{ with .Get "alt"}}{{.}}{{else}}{{.Get "title"}}{{end}}"{{ end }}
```

If a closing shortcode is used, the variable `.Inner` will be populated with all
of the content between the opening and closing shortcodes. If a closing
shortcode is required, you can check the length of `.Inner` and provide a warning
to the user.

A shortcode with `.Inner` content can be used without the inline content, and without the closing shortcode, by using the self-closing syntax:

```markdown
{{</* innershortcode /*/>}}
```

The variable `.Params` contains the list of parameters in case you need to do more complicated things than `.Get`.  It is sometimes useful to provide a flexible shortcode that can take named or positional parameters. To meet this need, Hugo shortcodes have `.IsNamedParams`, a boolean available that can be used such as `{{ if .IsNamedParams }}...{{ else }}...{{ end }}`. See the `Single Flexible Example` below for an example.

You can also use the variable `.Page` to access all the normal [page variables][].

A shortcodes can be nested. In a nested shortcode you can access the parent shortcode context with `.Parent`. This can be very useful for inheritance of common shortcode parameters from the root.

### Single Positional Example: `youtube`

```html
{{</* youtube 09jf3ow9jfw */>}}
```

Would load the template at `/layouts/shortcodes/youtube.html`:

{{% input "/layouts/shortcodes/youtube.html" %}}
```html
<div class="embed video-player">
<iframe class="youtube-player" type="text/html" width="640" height="385" src="http://www.youtube.com/embed/{{ index .Params 0 }}" allowfullscreen frameborder="0">
</iframe>
</div>
```
{{% /input %}}

This would be rendered as:

{{% output "youtube-embed.html" %}}
```html
<div class="embed video-player">
    <iframe class="youtube-player" type="text/html"
        width="640" height="385"
        src="http://www.youtube.com/embed/09jf3ow9jfw"
        allowfullscreen frameborder="0">
    </iframe>
</div>
```
{{% /output %}}

### Single Named Example: `image`

{{% input "content-image.md" %}}
```markdown
{{</* img src="/media/spf13.jpg" title="Steve Francia" */>}}
```
{{% /input %}}

Would load the template at `/layouts/shortcodes/img.html`:

{{% input "/layouts/shortcodes/img.html" %}}
```html
<!-- image -->
<figure {{ with .Get "class" }}class="{{.}}"{{ end }}>
    {{ with .Get "link"}}<a href="{{.}}">{{ end }}
        <img src="{{ .Get "src" }}" {{ if or (.Get "alt") (.Get "caption") }}alt="{{ with .Get "alt"}}{{.}}{{else}}{{ .Get "caption" }}{{ end }}"{{ end }} />
    {{ if .Get "link"}}</a>{{ end }}
    {{ if or (or (.Get "title") (.Get "caption")) (.Get "attr")}}
    <figcaption>{{ if isset .Params "title" }}
        <h4>{{ .Get "title" }}</h4>{{ end }}
        {{ if or (.Get "caption") (.Get "attr")}}<p>
        {{ .Get "caption" }}
        {{ with .Get "attrlink"}}<a href="{{.}}"> {{ end }}
            {{ .Get "attr" }}
        {{ if .Get "attrlink"}}</a> {{ end }}
        </p> {{ end }}
    </figcaption>
    {{ end }}
</figure>
<!-- image -->
```
{{% /input %}}

Would be rendered as:

{{% output "figure.html" %}}
```html
<figure >
    <img src="/media/spf13.jpg"  />
    <figcaption>
        <h4>Steve Francia</h4>
    </figcaption>
</figure>
```
{{% /output %}}

### Single Flexible Example: `vimeo`

```markdown
{{</* vimeo 49718712 */>}}
{{</* vimeo id="49718712" class="flex-video" */>}}
```

Would load the template found at `/layouts/shortcodes/vimeo.html`:

{{% input "/layouts/shortcodes/vimeo.html" %}}
```html
{{ if .IsNamedParams }}
  <div class="{{ if .Get "class" }}{{ .Get "class" }}{{ else }}vimeo-container{{ end }}">
    <iframe src="//player.vimeo.com/video/{{ .Get "id" }}" allowfullscreen></iframe>
  </div>
{{ else }}
  <div class="{{ if len .Params | eq 2 }}{{ .Get 1 }}{{ else }}vimeo-container{{ end }}">
    <iframe src="//player.vimeo.com/video/{{ .Get 0 }}" allowfullscreen></iframe>
  </div>
{{ end }}
```
{{% /input %}}

Would be rendered as:

{{% output "vimeo-iframes.html" %}}
```html
<div class="vimeo-container">
  <iframe src="//player.vimeo.com/video/49718712" allowfullscreen></iframe>
</div>
<div class="flex-video">
  <iframe src="//player.vimeo.com/video/49718712" allowfullscreen></iframe>
</div>
```
{{% /output %}}

### Paired Example: `highlight`

*Hugo already ships with the `highlight` shortcode*

{{% input "highlight-example.md" %}}
```markdown
{{</* highlight html */>}}
<html>
    <body> This HTML </body>
</html>
{{</* /highlight */>}}
```
{{% /input %}}

The template for this utilizes the following code (already included in Hugo)

```golang
{{ .Get 0 | highlight .Inner  }}
```

And will be rendered as:

{{% output "syntax-highlighted.html" %}}
```html
<div class="highlight" style="background: #272822"><pre style="line-height: 125%"><span style="color: #f92672">&lt;html&gt;</span>
    <span style="color: #f92672">&lt;body&gt;</span> This HTML <span style="color: #f92672">&lt;/body&gt;</span>
<span style="color: #f92672">&lt;/html&gt;</span>
</pre></div>
```
{{% /output %}}

{{% note %}}
The preceding template makes use of a Hugo-specific template function
called `highlight`, which uses [Pygments](http://pygments.org/) to add the highlighting code.
{{% /note %}}

### Simple Single-word Example: `year`

Let's assume you would like to have a shortcode to be replaced by the current year in your Markdown content files, for a license or copyright statement. Your goal is to be able to call the following shortcode in a content file:

```markdown
{{</* year */>}}
```

{{% input "/layouts/shortcodes/year.html" %}}
```golang
{{ .Page.Now.Year }}
```
{{% /input %}}

More shortcode examples can be found in the [shortcodes directory for spf13.com][].

[cross-references]: /content-management/cross-references/
[page variables]: /variables-and-params/page-variables/
[shortcodes directory for spf13.com]: https://github.com/spf13/spf13.com/tree/master/layouts/shortcodes
[source organization]: /project-organization/directory-structure/
[Speaker Deck]: https://speakerdeck.com/
[Vimeo]: https://vimeo.com/
[YouTube Videos]: https://www.youtube.com/