---
title: Create Your Own Shortcodes
linktitle: Shortcode Templates
description: You can extend Hugo's built-in shortcodes by creating your own using the same templating syntax as that for single and list pages.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
keywords: [shortcodes,templates]
menu:
  docs:
    parent: "templates"
    weight: 100
weight: 100
sections_weight: 100
draft: false
aliases: []
toc: true
---

Shortcodes are a means to consolidate templating into small, reusable snippets that you can embed directly inside of your content. In this sense, you can think of shortcodes as the intermediary between [page and list templates][templates] and [basic content files][].

{{% note %}}
Hugo also ships with built-in shortcodes for common use cases. (See [Content Management: Shortcodes](/content-management/shortcodes/).)
{{% /note %}}

## Create Custom Shortcodes

Hugo's built-in shortcodes cover many common, but not all, use cases. Luckily, Hugo provides the ability to easily create custom shortcodes to meet your website's needs.

{{< youtube Eu4zSaKOY4A >}}

### File Location

To create a shortcode, place an HTML template in the `layouts/shortcodes` directory of your [source organization][]. Consider the file name carefully since the shortcode name will mirror that of the file but without the `.html` extension. For example, `layouts/shortcodes/myshortcode.html` will be called with either `{{</* myshortcode /*/>}}` or `{{%/* myshortcode /*/%}}` depending on the type of parameters you choose.

You can organize your shortcodes in subfolders, e.g. in `layouts/shortcodes/boxes`. These shortcodes would then be accessible with their relative path, e.g:

```
{{</* boxes/square */>}}
```

Note the forward slash.

### Shortcode Template Lookup Order

Shortcode templates have a simple [lookup order][]:

1. `/layouts/shortcodes/<SHORTCODE>.html`
2. `/themes/<THEME>/layouts/shortcodes/<SHORTCODE>.html`

### Positional vs Named Parameters

You can create shortcodes using the following types of parameters:

* Positional parameters
* Named parameters
* Positional *or* named parameters (i.e, "flexible")

In shortcodes with positional parameters, the order of the parameters is important. If a shortcode has a single required value (e.g., the `youtube` shortcode below), positional parameters work very well and require less typing from content authors.

For more complex layouts with multiple or optional parameters, named parameters work best. While less terse, named parameters require less memorization from a content author and can be added in a shortcode declaration in any order.

Allowing both types of parameters (i.e., a "flexible" shortcode) is useful for complex layouts where you want to set default values that can be easily overridden by users.

### Access Parameters

All shortcode parameters can be accessed via the `.Get` method. Whether you pass a key (i.e., string) or a number to the `.Get` method depends on whether you are accessing a named or positional parameter, respectively.

To access a parameter by name, use the `.Get` method followed by the named parameter as a quoted string:

```
{{ .Get "class" }}
```

To access a parameter by position, use the `.Get` followed by a numeric position, keeping in mind that positional parameters are zero-indexed:

```
{{ .Get 0 }}
```

For the second position, you would just use: 

```
{{ .Get 1 }}
```

`with` is great when the output depends on a parameter being set:

```
{{ with .Get "class"}} class="{{.}}"{{ end }}
```

`.Get` can also be used to check if a parameter has been provided. This is
most helpful when the condition depends on either of the values, or both:

```
{{ or .Get "title" | .Get "alt" | if }} alt="{{ with .Get "alt"}}{{.}}{{else}}{{.Get "title"}}{{end}}"{{ end }}
```

#### `.Inner`

If a closing shortcode is used, the `.Inner` variable will be populated with all of the content between the opening and closing shortcodes. If a closing shortcode is required, you can check the length of `.Inner` as an indicator of its existence.

A shortcode with content declared via the `.Inner` variable can also be declared without the inline content and without the closing shortcode by using the self-closing syntax:

```
{{</* innershortcode /*/>}}
```

#### `.Params`

The `.Params` variable in shortcodes contains the list parameters passed to shortcode for more complicated use cases. You can also access higher-scoped parameters with the following logic:

`$.Params`
: these are the parameters passed directly into the shortcode declaration (e.g., a YouTube video ID)

`$.Page.Params`
: refers to the page's params; the "page" in this case refers to the content file in which the shortcode is declared (e.g., a `shortcode_color` field in a content's front matter could be accessed via `$.Page.Params.shortcode_color`).

`$.Page.Site.Params`
: refers to global variables as defined in your [site's configuration file][config].

#### `.IsNamedParams`

The `.IsNamedParams` variable checks whether the shortcode declaration uses named parameters and returns a boolean value.

For example, you could create an `image` shortcode that can take either a `src` named parameter or the first positional parameter, depending on the preference of the content's author. Let's assume the `image` shortcode is called as follows:

```
{{</* image src="images/my-image.jpg"*/>}}
```

You could then include the following as part of your shortcode templating:

```
{{ if .IsNamedParams }}
<img src="{{.Get "src" }}" alt="">
{{ else }}
<img src="{{.Get 0}}" alt="">
{{ end }}
```

See the [example Vimeo shortcode][vimeoexample] below for `.IsNamedParams` in action.

{{% warning %}}
While you can create shortcode templates that accept both positional and named parameters, you *cannot* declare shortcodes in content with a mix of parameter types. Therefore, a shortcode declared like `{{</* image src="images/my-image.jpg" "This is my alt text" */>}}` will return an error.
{{% /warning %}}

You can also use the variable `.Page` to access all the normal [page variables][pagevars].

A shortcodes can also be nested. In a nested shortcode, you can access the parent shortcode context with [`.Parent` variable][shortcodesvars]. This can be very useful for inheritance of common shortcode parameters from the root.

### Checking for Existence

You can check if a specific shortcode is used on a page by calling `.HasShortcode` in that page template, providing the name of the shortcode. This is sometimes useful when you want to include specific scripts or styles in the header that are only used by that shortcode.

## Custom Shortcode Examples

The following are examples of the different types of shortcodes you can create via shortcode template files in `/layouts/shortcodes`.

### Single-word Example: `year`

Let's assume you would like to keep mentions of your copyright year current in your content files without having to continually review your markdown. Your goal is to be able to call the shortcode as follows:

```
{{</* year */>}}
```

{{< code file="/layouts/shortcodes/year.html" >}}
{{ now.Format "2006" }}
{{< /code >}}

### Single Positional Example: `youtube`

Embedded videos are a common addition to markdown content that can quickly become unsightly. The following is the code used by [Hugo's built-in YouTube shortcode][youtubeshortcode]:

```
{{</* youtube 09jf3ow9jfw */>}}
```

Would load the template at `/layouts/shortcodes/youtube.html`:

{{< code file="/layouts/shortcodes/youtube.html" >}}
<div class="embed video-player">
<iframe class="youtube-player" type="text/html" width="640" height="385" src="https://www.youtube.com/embed/{{ index .Params 0 }}" allowfullscreen frameborder="0">
</iframe>
</div>
{{< /code >}}

{{< code file="youtube-embed.html" copy="false" >}}
<div class="embed video-player">
    <iframe class="youtube-player" type="text/html"
        width="640" height="385"
        src="https://www.youtube.com/embed/09jf3ow9jfw"
        allowfullscreen frameborder="0">
    </iframe>
</div>
{{< /code >}}

### Single Named Example: `image`

Let's say you want to create your own `img` shortcode rather than use Hugo's built-in [`figure` shortcode][figure]. Your goal is to be able to call the shortcode as follows in your content files:

{{< code file="content-image.md" >}}
{{</* img src="/media/spf13.jpg" title="Steve Francia" */>}}
{{< /code >}}

You have created the shortcode at `/layouts/shortcodes/img.html`, which loads the following shortcode template:

{{< code file="/layouts/shortcodes/img.html" >}}
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
{{< /code >}}

Would be rendered as:

{{< code file="img-output.html" copy="false" >}}
<figure>
    <img src="/media/spf13.jpg"  />
    <figcaption>
        <h4>Steve Francia</h4>
    </figcaption>
</figure>
{{< /code >}}

### Single Flexible Example: `vimeo`

```
{{</* vimeo 49718712 */>}}
{{</* vimeo id="49718712" class="flex-video" */>}}
```

Would load the template found at `/layouts/shortcodes/vimeo.html`:

{{< code file="/layouts/shortcodes/vimeo.html" >}}
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

{{< code file="vimeo-iframes.html" copy="false" >}}
<div class="vimeo-container">
  <iframe src="https://player.vimeo.com/video/49718712" allowfullscreen></iframe>
</div>
<div class="flex-video">
  <iframe src="https://player.vimeo.com/video/49718712" allowfullscreen></iframe>
</div>
{{< /code >}}

### Paired Example: `highlight`

The following is taken from `highlight`, which is a [built-in shortcode][] that ships with Hugo.

{{< code file="highlight-example.md" >}}
{{</* highlight html */>}}
  <html>
    <body> This HTML </body>
  </html>
{{</* /highlight */>}}
{{< /code >}}

The template for the `highlight` shortcode uses the following code, which is already included in Hugo:

```
{{ .Get 0 | highlight .Inner  }}
```

The rendered output of the HTML example code block will be as follows:

{{< code file="syntax-highlighted.html" copy="false" >}}
<div class="highlight" style="background: #272822"><pre style="line-height: 125%"><span style="color: #f92672">&lt;html&gt;</span>
    <span style="color: #f92672">&lt;body&gt;</span> This HTML <span style="color: #f92672">&lt;/body&gt;</span>
<span style="color: #f92672">&lt;/html&gt;</span>
</pre></div>
{{< /code >}}

### Nested Shortcode: Image Gallery

Hugo's [`.Parent` shortcode variable][parent] returns a boolean value depending on whether the shortcode in question is called within the context of a *parent* shortcode. This provides an inheritance model for common shortcode parameters.

The following example is contrived but demonstrates the concept. Assume you have a `gallery` shortcode that expects one named `class` parameter:

{{< code file="layouts/shortcodes/gallery.html" >}}
<div class="{{.Get "class"}}">
  {{.Inner}}
</div>
{{< /code >}}

You also have an `img` shortcode with a single named `src` parameter that you want to call inside of `gallery` and other shortcodes, so that the parent defines the context of each `img`:

{{< code file="layouts/shortcodes/img.html" >}}
{{- $src := .Get "src" -}}
{{- with .Parent -}}
  <img src="{{$src}}" class="{{.Get "class"}}-image">
{{- else -}}
  <img src="{{$src}}">
{{- end }}
{{< /code >}}

You can then call your shortcode in your content as follows:

```
{{</* gallery class="content-gallery" */>}}
  {{</* img src="/images/one.jpg" */>}}
  {{</* img src="/images/two.jpg" */>}}
{{</* /gallery */>}}
{{</* img src="/images/three.jpg" */>}}
```

This will output the following HTML. Note how the first two `img` shortcodes inherit the `class` value of `content-gallery` set with the call to the parent `gallery`, whereas the third `img` only uses `src`:

```
<div class="content-gallery">
    <img src="/images/one.jpg" class="content-gallery-image">
    <img src="/images/two.jpg" class="content-gallery-image">
</div>
<img src="/images/three.jpg">
```


## Error Handling in Shortcodes

Use the [errorf](/functions/errorf) template func and [.Position](/variables/shortcodes/) variable to get useful error messages in shortcodes:

```bash
{{ with .Get "name" }}
{{ else }}
{{ errorf "missing value for param 'name': %s" .Position }}
{{ end }}
```

When the above fails, you will see an `ERROR` log similar to the below:

```bash
ERROR 2018/11/07 10:05:55 missing value for param name: "/Users/bep/dev/go/gohugoio/hugo/docs/content/en/variables/shortcodes.md:32:1"
```

## More Shortcode Examples

More shortcode examples can be found in the [shortcodes directory for spf13.com][spfscs] and the [shortcodes directory for the Hugo docs][docsshortcodes].


## Inline Shortcodes

Since Hugo 0.52, you can implement your shortcodes inline -- e.g. where you use them in the content file. This can be useful for scripting that you only need in one place.

This feature is disabled by default, but can be enabled in your site config:

{{< code-toggle file="config">}}
enableInlineShortcodes = true
{{< /code-toggle >}}

It is disabled by default for security reasons. The security model used by Hugo's template handling assumes that template authors are trusted, but that the content files are not, so the templates are injection-safe from malformed input data. But in most situations you have full control over the content, too, and then `enableInlineShortcodes = true` would be considered safe. But it's something to be aware of: It allows ad-hoc [Go Text templates](https://golang.org/pkg/text/template/) to be executed from the content files.

And once enabled, you can do this in your content files:

 ```go-text-template
 {{</* time.inline */>}}{{ now }}{{</* /time.inline */>}}
 ```

The above will print the current date and time.

 Note that an inline shortcode's inner content is parsed and executed as a Go text template with the same context as a regular shortcode template.

This means that the current page can be accessed via `.Page.Title` etc. This also means that there are no concept of "nested inline shortcodes".

The same inline shortcode can be reused later in the same content file, with different params if needed, using the self-closing syntax:

 ```go-text-template
{{</* time.inline /*/>}}
```
	

[basic content files]: /content-management/formats/ "See how Hugo leverages markdown--and other supported formats--to create content for your website."
[built-in shortcode]: /content-management/shortcodes/
[config]: /getting-started/configuration/ "Learn more about Hugo's built-in configuration variables as well as how to us your site's configuration file to include global key-values that can be used throughout your rendered website."
[Content Management: Shortcodes]: /content-management/shortcodes/#using-hugo-s-built-in-shortcodes "Check this section if you are not familiar with the definition of what a shortcode is or if you are unfamiliar with how to use Hugo's built-in shortcodes in your content files."
[source organization]: /getting-started/directory-structure/#directory-structure-explained "Learn how Hugo scaffolds new sites and what it expects to find in each of your directories."
[docsshortcodes]: https://github.com/gohugoio/hugo/tree/master/docs/layouts/shortcodes "See the shortcode source directory for the documentation site you're currently reading."
[figure]: /content-management/shortcodes/#figure
[hugosc]: /content-management/shortcodes/#using-hugo-s-built-in-shortcodes
[lookup order]: /templates/lookup-order/ "See the order in which Hugo traverses your template files to decide where and how to render your content at build time"
[pagevars]: /variables/page/ "See which variables you can leverage in your templating for page vs list templates."
[parent]: /variables/shortcodes/
[shortcodesvars]: /variables/shortcodes/ "Certain variables are specific to shortcodes, although most .Page variables can be accessed within your shortcode template."
[spfscs]: https://github.com/spf13/spf13.com/tree/master/layouts/shortcodes "See more examples of shortcodes by visiting the shortcode directory of the source for spf13.com, the blog of Hugo's creator, Steve Francia."
[templates]: /templates/ "The templates section of the Hugo docs."
[vimeoexample]: #single-flexible-example-vimeo
[youtubeshortcode]: /content-management/shortcodes/#youtube "See how to use Hugo's built-in YouTube shortcode."
