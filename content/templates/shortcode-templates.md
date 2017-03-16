---
title: Creating Your Own Shortcode Templates
linktitle: Shortcode Templates
description: You can extend Hugo's built-in shortcodes by creating your own using the same templating syntax as that for single and list pages.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [shortcodes]
weight: 100
draft: false
aliases: []
toc: true
---

## Creating Custom Shortcodes

Hugo's built-in shortcodes cover many common, but not all, use cases. Luckily, Hugo provides the ability to easily create custom shortcodes to meet your website's needs. In this sense, you can think of shortcodes as the intermediary between [page and list templates][templates] and [basic content files][].

### File Placement

To create a shortcode, place an HTML template in the `layouts/shortcodes` directory of your [source organization][]. Consider the file name carefully since the shortcode name will mirror that of the file but without the `.html` extension. For example, `layouts/shortcodes/myshortcode.html` will be called with either `{{</* myshortcode /*/>}}` or `{{%/* myshortcode /*/%}}` depending on the type of parameters you choose.

### Deciding on Shortcode and Parameter Type

You can create the following types of shortcodes

* Positional parameters
* Named parameters
* Positional *or* named parameters (i.e, "flexible")
* Single-word shortcodes
* Nested

#### Positional Parameters

In shortcodes with positional parameters, the order of the parameters is important.

you can choose if the shortcode will use _positional parameters_, or _named parameters_, or _both_. A good rule of thumb is that if a shortcode has a single required value in the case of the `youtube` example below, then positional works very well. For more complex layouts with optional parameters, named parameters work best. Allowing both types of parameters is useful for complex layouts where you want to set default values that can be overridden.

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

If a closing shortcode is used, the variable `.Inner` will be populated with all of the content between the opening and closing shortcodes. If a closing shortcode is required, you can check the length of `.Inner` and provide a warning to the user.

A shortcode with `.Inner` content can be used without the inline content, and without the closing shortcode, by using the self-closing syntax:

```golang
{{</* innershortcode /*/>}}
```

The variable `.Params` contains the list of parameters in case you need to do more complicated things than `.Get`.  It is sometimes useful to provide a flexible shortcode that can take named or positional parameters. To meet this need, Hugo shortcodes have `.IsNamedParams`, a boolean available that can be used such as `{{ if .IsNamedParams }}...{{ else }}...{{ end }}`. See the [example Vimeo shortcode][vimeoexample] below for an example.

You can also use the variable `.Page` to access all the normal [page variables][pagevars].

A shortcodes can also be nested. In a nested shortcode, you can access the parent shortcode context with [`.Parent` variable][shortcodesvars]. This can be very useful for inheritance of common shortcode parameters from the root.

## Custom Shortcode Examples

The following are examples of the different types of shortcodes you can create via template files in `/layouts/shortcodes`.

### Single-word Example: `year`

Let's assume you would like to keep mentions of your copyright year current in your content files without having to continually review your markdown. Your goal is to be able to call the shortcode as follows:

```markdown
{{</* year */>}}
```

{{% code file="/layouts/shortcodes/year.html" %}}
```golang
{{ .Page.Now.Year }}
```
{{% /code %}}

### Single Positional Example: `youtube`

Embedded videos are a common addition to markdown content that can quickly become unsightly. The following is the code used by [Hugo's built-in YouTube shortcode][youtubeshortcode]:

```golang
{{</* youtube 09jf3ow9jfw */>}}
```

Would load the template at `/layouts/shortcodes/youtube.html`:

{{% code file="/layouts/shortcodes/youtube.html" %}}
```html
<div class="embed video-player">
<iframe class="youtube-player" type="text/html" width="640" height="385" src="http://www.youtube.com/embed/{{ index .Params 0 }}" allowfullscreen frameborder="0">
</iframe>
</div>
```
{{% /code %}}

{{% code file="youtube-embed.html" copy="false" %}}
```html
<div class="embed video-player">
    <iframe class="youtube-player" type="text/html"
        width="640" height="385"
        src="http://www.youtube.com/embed/09jf3ow9jfw"
        allowfullscreen frameborder="0">
    </iframe>
</div>
```
{{% /code %}}

### Single Named Example: `image`

Let's say you want to create your own `img` shortcode rather than use Hugo's built-in [`figure` shortcode][figure]. Your goal is to be able to call the shortcode as follows in your content files:

{{% code file="content-image.md" %}}
```golang
{{</* img src="/media/spf13.jpg" title="Steve Francia" */>}}
```
{{% /code %}}

You have created the shortcode at `/layouts/shortcodes/img.html`, which loads the following shortcode template:

{{% code file="/layouts/shortcodes/img.html" %}}
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
{{% /code %}}

Would be rendered as:

{{% code file="img-output.html" copy="false" %}}
```html
<figure >
    <img src="/media/spf13.jpg"  />
    <figcaption>
        <h4>Steve Francia</h4>
    </figcaption>
</figure>
```
{{% /code %}}

### Single Flexible Example: `vimeo`

```golang
{{</* vimeo 49718712 */>}}
{{</* vimeo id="49718712" class="flex-video" */>}}
```

Would load the template found at `/layouts/shortcodes/vimeo.html`:

{{% code file="/layouts/shortcodes/vimeo.html" %}}
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
{{% /code %}}

Would be rendered as:

{{% code file="vimeo-iframes.html" copy="false" %}}
```html
<div class="vimeo-container">
  <iframe src="//player.vimeo.com/video/49718712" allowfullscreen></iframe>
</div>
<div class="flex-video">
  <iframe src="//player.vimeo.com/video/49718712" allowfullscreen></iframe>
</div>
```
{{% /code %}}

### Paired Example: `highlight`

The following is taken from `highlight`, which is a [built-in shortcode][] that ships with Hugo.

{{% code file="highlight-example.md" %}}
```markdown
{{</* highlight html */>}}
  <html>
    <body> This HTML </body>
  </html>
{{</* /highlight */>}}
```
{{% /code %}}

The template for the `highlight` shortcode uses the following code, which is already included in Hugo:

```golang
{{ .Get 0 | highlight .Inner  }}
```

The rendered output of the HTML example code block will be as follows:

{{% code file="syntax-highlighted.html" copy="false" %}}
```html
<div class="highlight" style="background: #272822"><pre style="line-height: 125%"><span style="color: #f92672">&lt;html&gt;</span>
    <span style="color: #f92672">&lt;body&gt;</span> This HTML <span style="color: #f92672">&lt;/body&gt;</span>
<span style="color: #f92672">&lt;/html&gt;</span>
</pre></div>
```
{{% /code %}}

{{% note %}}
The preceding shortcode makes use of a Hugo-specific template function called `highlight`, which uses [Pygments](http://pygments.org) to add syntax highlighting to the example HTML code block. See the [developer tools page on syntax highlighting](/tools/syntax-highlighting/) for more information.
{{% /note %}}

### Nested Shortcode: Image Gallery

Hugo's [`.Parent` shortcode variable][parent] returns a boolean value depending on whether the shortcode in question is called within the context of a *parent* shortcode. This provides an inheritance model for common shortcode parameters.

The following example is contrived but demonstrates the concept. Assume you have a `gallery` shortcode that expects one named `class` parameter:

{{% code file="layouts/shortcodes/gallery.html" %}}
```html
<div class="{{.Get "class"}}">
  {{.Inner}}
</div>
```
{{% /code %}}

You also have an `image` shortcode with a single named `src` parameter that may or may not be called within `gallery`, as well as other shortcodes you've created for your site. When called inside a parent shortcode like `gallery`, you want the output HTML to inherit a CSS class from the value given to its parent:

{{% code file="layouts/shortcodes/image.html" %}}
```html
{{- $src := .Get "src" -}}
{{- with .Parent -}}
  <img src="{{$src}}" class="{{.Get "class"}}-image">
{{- else -}}
  <img src="{{$src}}">
{{- end }}
```
{{% /code %}}

You can then call your shortcode in your content as follows:

```markdown
{{</* gallery class="content-gallery" */>}}
  {{</* img src="/images/one.jpg" */>}}
  {{</* img src="/images/two.jpg" */>}}
{{</* /gallery */>}}
{{</* img src="/images/three.jpg" */>}}
```

This will output the following HTML. Note how the first two `image` shortcodes inherit the `class` value of `content-gallery` set with the call to the parent `gallery`, whereas the third `image` only uses `src`:

```html
<div class="content-gallery">
    <img src="/images/one.jpg" class="content-gallery-image">
    <img src="/images/two.jpg" class="content-gallery-image">
</div>
<img src="/images/three.jpg">
```

## More Shortcode Examples

More shortcode examples can be found in the [shortcodes directory for spf13.com][spfscs] and the [shortcodes directory for the Hugo docs][docsshortcodes].

[basic content files]: /content-management/formats/ "See how Hugo leverages markdown--and other supported formats--to create content for your website."
[built-in shortcode]: /content-management/shortcodes/
[source organization]: /getting-started/directory-structure/ "Learn how Hugo scaffolds new sites and what it expects to find in each of your directories."
[docsshortcodes]: https://github.com/spf13/hugo/tree/master/docs/layouts/shortcodes "See the shortcode source directory for the documentation site you're currently reading."
[figure]: /content-management/shortcodes/#figure
[pagevars]: /variables/page/ "See which variables you can leverage in your templating for page vs list templates."
[parent]: /variables/shortcodes/
[shortcodesvars]: /variables/shortcodes/ "Certain variables are specific to shortcodes, although most .Page variables can be accessed within your shortcode template."
[spfscs]: https://github.com/spf13/spf13.com/tree/master/layouts/shortcodes "See more examples of shortcodes by visiting the shortcode directory of the source for spf13.com, the blog of Hugo's creator, Steve Francia."
[templates]: /templates/ "The templates section of the Hugo docs."
[vimeoexample]: #single-flexible-example-vimeo
[youtubeshortcode]: /content-management/shortcodes/#youtube "See how to use Hugo's built-in YouTube shortcode."