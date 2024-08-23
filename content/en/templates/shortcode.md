---
title: Create your own shortcodes
linkTitle: Shortcode templates
description: You can extend Hugo's embedded shortcodes by creating your own using the same templating syntax as that for single and list pages.
categories: [templates]
keywords: []
menu:
  docs:
    parent: templates
    weight: 130
weight: 130
aliases: [/templates/shortcode-templates/]
toc: true
---

Shortcodes are a means to consolidate templating into small, reusable snippets that you can embed directly inside your content.

{{% note %}}
Hugo also ships with embedded shortcodes for common use cases. (See [Content Management: Shortcodes](/content-management/shortcodes/).)
{{% /note %}}

## Create custom shortcodes

Hugo's embedded shortcodes cover many common, but not all, use cases. Luckily, Hugo provides the ability to easily create custom shortcodes to meet your website's needs.

{{< youtube Eu4zSaKOY4A >}}

### File location

To create a shortcode, place an HTML template in the `layouts/shortcodes` directory. Consider the file name carefully since the shortcode name will mirror that of the file but without the `.html` extension. For example, `layouts/shortcodes/myshortcode.html` will be called with either `{{</* myshortcode /*/>}}` or `{{%/* myshortcode /*/%}}`.

You can organize your shortcodes in subdirectories, e.g. in `layouts/shortcodes/boxes`. These shortcodes would then be accessible with their relative path, e.g:

```go-html-template
{{</* boxes/square */>}}
```

Note the forward slash.

### Template lookup order

Hugo selects shortcode templates based on the shortcode name, the current output format, and the current language. The examples below are sorted by specificity in descending order. The least specific path is at the bottom of the list.

Shortcode name|Output format|Language|Template path
:--|:--|:--|:--
foo|html|en|layouts/shortcodes/foo.en.html
foo|html|en|layouts/shortcodes/foo.html.html
foo|html|en|layouts/shortcodes/foo.html
foo|html|en|layouts/shortcodes/foo.html.en.html

Shortcode name|Output format|Language|Template path
:--|:--|:--|:--
foo|rss|en|layouts/shortcodes/foo.en.xml
foo|rss|en|layouts/shortcodes/foo.rss.xml
foo|rss|en|layouts/shortcodes/foo.en.html
foo|rss|en|layouts/shortcodes/foo.rss.en.xml
foo|rss|en|layouts/shortcodes/foo.xml
foo|rss|en|layouts/shortcodes/foo.html.en.html
foo|rss|en|layouts/shortcodes/foo.html.html
foo|rss|en|layouts/shortcodes/foo.html

Note that templates provided by a theme or module always take precedence.

### Positional vs. named arguments

You can create shortcodes using the following types of arguments:

* Positional arguments
* Named arguments
* Positional *or* named arguments

In shortcodes with positional arguments, the order of the arguments is important. If a shortcode has a single required value, positional arguments require less typing from content authors.

For more complex layouts with multiple or optional arguments, named arguments work best. While less terse, named arguments require less memorization from a content author and can be added in a shortcode declaration in any order.

Allowing both types of arguments is useful for complex layouts where you want to set default values that can be easily overridden by users.

### Access arguments

All shortcode arguments can be accessed via the `.Get` method. Whether you pass a string or a number to the `.Get` method depends on whether you are accessing a named or positional argument, respectively.

To access an argument by name, use the `.Get` method followed by the named argument as a quoted string:

```go-html-template
{{ .Get "class" }}
```

To access an argument by position, use the `.Get` followed by a numeric position, keeping in mind that positional arguments are zero-indexed:

```go-html-template
{{ .Get 0 }}
```

For the second position, you would just use:

```go-html-template
{{ .Get 1 }}
```

`with` is great when the output depends on a argument being set:

```go-html-template
{{ with .Get "class" }} class="{{ . }}"{{ end }}
```

`.Get` can also be used to check if a argument has been provided. This is
most helpful when the condition depends on either of the values, or both:

```go-html-template
{{ if or (.Get "title") (.Get "alt") }} alt="{{ with .Get "alt" }}{{ . }}{{ else }}{{ .Get "title" }}{{ end }}"{{ end }}
```

#### `.Inner`

The `.Inner` method returns the content between the opening and closing shortcode tags. To check if `.Inner` returns anything other than whitespace:

```go-html-template
{{ if strings.ContainsNonSpace .Inner }}
  Inner is not empty
{{ end }}
```

{{% note %}}
Any shortcode that calls the `.Inner` method must be closed or self-closed. To call a shortcode using the self-closing syntax.

```go-html-template
{{</* innershortcode /*/>}}
```

{{% /note %}}

#### `.Params`

The `.Params` method in shortcodes returns the arguments passed to the shortcode for more complicated use cases. You can also access higher-scoped arguments with the following logic:

$.Params
: these are the arguments passed directly into the shortcode declaration (e.g., a YouTube video ID)

$.Page.Params
: refers to the page's parameters; the "page" in this case refers to the content file in which the shortcode is declared (e.g., a `shortcode_color` field in a content's front matter could be accessed via `$.Page.Params.shortcode_color`).

$.Site.Params
: refers to parameters defined in your site configuration.

#### `.IsNamedParams`

The `.IsNamedParams` method checks whether the shortcode declaration uses named arguments and returns a boolean value.

For example, you could create an `image` shortcode that can take either a `src` named argument or the first positional argument, depending on the preference of the content's author. Let's assume the `image` shortcode is called as follows:

```go-html-template
{{</* image src="images/my-image.jpg" */>}}
```

You could then include the following as part of your shortcode templating:

```go-html-template
{{ if .IsNamedParams }}
  <img src="{{ .Get "src" }}" alt="">
{{ else }}
  <img src="{{ .Get 0 }}" alt="">
{{ end }}
```

See the [example Vimeo shortcode][vimeoexample] below for `.IsNamedParams` in action.

{{% note %}}
While you can create shortcode templates that accept both positional and named arguments, you *cannot* declare shortcodes in content with a mix of argument types. Therefore, a shortcode declared like `{{</* image src="images/my-image.jpg" "This is my alt text" */>}}` will return an error.
{{% /note %}}

Shortcodes can also be nested. In a nested shortcode, you can access the parent shortcode context with the [`.Parent`] shortcode method. This can be very useful for inheritance from the root.

### Checking for existence

You can check if a specific shortcode is used on a page by calling `.HasShortcode` in that page template, providing the name of the shortcode. This is useful when you want to include specific scripts or styles in the header that are only used by that shortcode.

## Custom shortcode examples

The following are examples of the different types of shortcodes you can create via shortcode template files in `/layouts/shortcodes`.

### Single-word example: `year`

Let's assume you would like to keep mentions of your copyright year current in your content files without having to continually review your Markdown. Your goal is to be able to call the shortcode as follows:

```go-html-template
{{</* year */>}}
```

{{< code file=layouts/shortcodes/year.html >}}
{{ now.Format "2006" }}
{{< /code >}}

### Single positional example: `youtube`

Embedded videos are a common addition to Markdown content. The following is the code used by [Hugo's built-in YouTube shortcode][youtubeshortcode]:

```go-html-template
{{</* youtube 09jf3ow9jfw */>}}
```

Would load the template at `/layouts/shortcodes/youtube.html`:

{{< code file=layouts/shortcodes/youtube.html >}}
<div class="embed video-player">
<iframe class="youtube-player" type="text/html" width="640" height="385" src="https://www.youtube.com/embed/{{ index .Params 0 }}" allowfullscreen frameborder="0">
</iframe>
</div>
{{< /code >}}

{{< code file=youtube-embed.html >}}
<div class="embed video-player">
    <iframe class="youtube-player" type="text/html"
        width="640" height="385"
        src="https://www.youtube.com/embed/09jf3ow9jfw"
        allowfullscreen frameborder="0">
    </iframe>
</div>
{{< /code >}}

### Single named example: `image`

Let's say you want to create your own `img` shortcode rather than use Hugo's built-in [`figure` shortcode][figure]. Your goal is to be able to call the shortcode as follows in your content files:

{{< code file=content-image.md >}}
{{</* img src="/media/spf13.jpg" title="Steve Francia" */>}}
{{< /code >}}

You have created the shortcode at `/layouts/shortcodes/img.html`, which loads the following shortcode template:

{{< code file=layouts/shortcodes/img.html >}}
<!-- image -->
<figure {{ with .Get "class" }}class="{{ . }}"{{ end }}>
  {{ with .Get "link" }}<a href="{{ . }}">{{ end }}
    <img src="{{ .Get "src" }}" {{ if or (.Get "alt") (.Get "caption") }}alt="{{ with .Get "alt" }}{{ . }}{{ else }}{{ .Get "caption" }}{{ end }}"{{ end }} />
    {{ if .Get "link" }}</a>{{ end }}
    {{ if or (or (.Get "title") (.Get "caption")) (.Get "attr") }}
      <figcaption>{{ if isset .Params "title" }}
        <h4>{{ .Get "title" }}</h4>{{ end }}
        {{ if or (.Get "caption") (.Get "attr") }}<p>
        {{ .Get "caption" }}
        {{ with .Get "attrlink" }}<a href="{{ . }}"> {{ end }}
          {{ .Get "attr" }}
        {{ if .Get "attrlink" }}</a> {{ end }}
        </p> {{ end }}
      </figcaption>
  {{ end }}
</figure>
<!-- image -->
{{< /code >}}

Would be rendered as:

{{< code file=img-output.html >}}
<figure>
  <img src="/media/spf13.jpg"  />
  <figcaption>
      <h4>Steve Francia</h4>
  </figcaption>
</figure>
{{< /code >}}

### Single flexible example: `vimeo`

```go-html-template
{{</* vimeo 49718712 */>}}
{{</* vimeo id="49718712" class="flex-video" */>}}
```

Would load the template found at `/layouts/shortcodes/vimeo.html`:

{{< code file=layouts/shortcodes/vimeo.html >}}
{{ if .IsNamedParams }}
  <div class="{{ if .Get "class" }}{{ .Get "class" }}{{ else }}vimeo-container{{ end }}">
    <iframe src="https://player.vimeo.com/video/{{ .Get "id" }}" allowfullscreen></iframe>
  </div>
{{ else }}
  <div class="{{ if len .Params | eq 2 }}{{ .Get 1 }}{{ else }}vimeo-container{{ end }}">
    <iframe src="https://player.vimeo.com/video/{{ .Get 0 }}" allowfullscreen></iframe>
  </div>
{{ end }}
{{< /code >}}

Would be rendered as:

{{< code file=vimeo-iframes.html >}}
<div class="vimeo-container">
  <iframe src="https://player.vimeo.com/video/49718712" allowfullscreen></iframe>
</div>
<div class="flex-video">
  <iframe src="https://player.vimeo.com/video/49718712" allowfullscreen></iframe>
</div>
{{< /code >}}

### Paired example: `highlight`

The following is taken from `highlight`, which is a [built-in shortcode] that ships with Hugo.

{{< code file=highlight-example.md >}}
{{</* highlight html */>}}
  <html>
    <body> This HTML </body>
  </html>
{{</* /highlight */>}}
{{< /code >}}

The template for the `highlight` shortcode uses the following code, which is already included in Hugo:

```go-html-template
{{ .Get 0 | highlight .Inner }}
```

The rendered output of the HTML example code block will be as follows:

{{< code file=syntax-highlighted.html >}}
<div class="highlight" style="background: #272822"><pre style="line-height: 125%"><span style="color: #f92672">&lt;html&gt;</span>
    <span style="color: #f92672">&lt;body&gt;</span> This HTML <span style="color: #f92672">&lt;/body&gt;</span>
<span style="color: #f92672">&lt;/html&gt;</span>
</pre></div>
{{< /code >}}

### Nested shortcode: image gallery

Hugo's [`.Parent`] shortcode method provides access to the parent shortcode context when the shortcode in question is called within the context of a parent shortcode. This provides an inheritance model.

The following example is contrived but demonstrates the concept. Assume you have a `gallery` shortcode that expects one named `class` argument:

{{< code file=layouts/shortcodes/gallery.html >}}
<div class="{{ .Get "class" }}">
  {{ .Inner }}
</div>
{{< /code >}}

You also have an `img` shortcode with a single named `src` argument that you want to call inside of `gallery` and other shortcodes, so that the parent defines the context of each `img`:

{{< code file=layouts/shortcodes/img.html >}}
{{- $src := .Get "src" -}}
{{- with .Parent -}}
  <img src="{{ $src }}" class="{{ .Get "class" }}-image">
{{- else -}}
  <img src="{{ $src }}">
{{- end -}}
{{< /code >}}

You can then call your shortcode in your content as follows:

```go-html-template
{{</* gallery class="content-gallery" */>}}
  {{</* img src="/images/one.jpg" */>}}
  {{</* img src="/images/two.jpg" */>}}
{{</* /gallery */>}}
{{</* img src="/images/three.jpg" */>}}
```

This will output the following HTML. Note how the first two `img` shortcodes inherit the `class` value of `content-gallery` set with the call to the parent `gallery`, whereas the third `img` only uses `src`:

```html
<div class="content-gallery">
    <img src="/images/one.jpg" class="content-gallery-image">
    <img src="/images/two.jpg" class="content-gallery-image">
</div>
<img src="/images/three.jpg">
```

## Error handling in shortcodes

Use the [`errorf`] template function with the [`Name`] and [`Position`] shortcode methods to generate useful error messages:

{{< code file=layouts/shortcodes/greeting.html >}}
{{ with .Get "name" }}
  <p>Hello, my name is {{ . }}.</p>
{{ else }}
  {{ errorf "The %q shortcode requires a 'name' argument. See %s" .Name .Position }}
{{ end }}
{{< /code >}}

When the above fails, you will see an `ERROR` message such as:

```sh
ERROR The "greeting" shortcode requires a 'name' argument. See "/home/user/project/content/_index.md:12:1"
```

## Inline shortcodes

You can also implement your shortcodes inline -- e.g. where you use them in the content file. This can be useful for scripting that you only need in one place.

This feature is disabled by default, but can be enabled in your site configuration:

{{< code-toggle file=hugo >}}
[security]
enableInlineShortcodes = true
{{< /code-toggle >}}

It is disabled by default for security reasons. The security model used by Hugo's template handling assumes that template authors are trusted, but that the content files are not, so the templates are injection-safe from malformed input data. But in most situations you have full control over the content, too, and then `enableInlineShortcodes = true` would be considered safe. But it's something to be aware of: It allows ad-hoc [Go Text templates](https://golang.org/pkg/text/template/) to be executed from the content files.

And once enabled, you can do this in your content files:

 ```go-html-template
 {{</* time.inline */>}}{{ now }}{{</* /time.inline */>}}
 ```

The above will print the current date and time.

Note that an inline shortcode's inner content is parsed and executed as a Go text template with the same context as a regular shortcode template.

This means that the current page can be accessed via `.Page.Title` etc. This also means that there are no concept of "nested inline shortcodes".

The same inline shortcode can be reused later in the same content file, with different arguments if needed, using the self-closing syntax:

 ```go-html-template
{{</* time.inline /*/>}}
```

[`.Parent`]: /methods/shortcode/parent/
[`errorf`]: /functions/fmt/errorf/
[`Name`]: /methods/shortcode/name/
[`Position`]: /methods/shortcode/position/
[built-in shortcode]: /content-management/shortcodes/
[figure]: /content-management/shortcodes/#figure
[lookup order]: /templates/lookup-order/
[source organization]: /getting-started/directory-structure/
[vimeoexample]: #single-flexible-example-vimeo
[youtubeshortcode]: /content-management/shortcodes/#youtube
