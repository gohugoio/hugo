---
title: Shortcode templates
description: Create custom shortcodes to simplify and standardize content creation.
categories: []
keywords: []
weight: 120
aliases: [/templates/shortcode-templates/]
---

{{< newtemplatesystem >}}


> [!note]
> Before creating custom shortcodes, please review the [shortcodes] page in the [content management] section. Understanding the usage details will help you design and create better templates.

## Introduction

Hugo provides [embedded shortcodes] for many common tasks, but you'll likely need to create your own for more specific needs. Some examples of custom shortcodes you might develop include:

- Audio players
- Video players
- Image galleries
- Diagrams
- Maps
- Tables
- And many other custom elements

## Directory structure

Create shortcode templates within the `layouts/_shortcodes` directory, either at its root or organized into subdirectories.

```text
layouts/
└── _shortcodes/
    ├── diagrams/
    │   ├── kroki.html
    │   └── plotly.html
    ├── media/
    │   ├── audio.html
    │   ├── gallery.html
    │   └── video.html
    ├── capture.html
    ├── column.html
    ├── include.html
    └── row.html
```

When calling a shortcode in a subdirectory, specify its path relative to the `_shortcode` directory, excluding the file extension.

```text
{{</* media/audio path=/audio/podcast/episode-42.mp3 */>}}
```

## Lookup order

Hugo selects shortcode templates based on the shortcode name, the current output format, and the current language. The examples below are sorted by specificity in descending order. The least specific path is at the bottom of the list.

Shortcode name|Output format|Language|Template path
:--|:--|:--|:--
foo|html|en|`layouts/_shortcodes/foo.en.html`
foo|html|en|`layouts/_shortcodes/foo.html.html`
foo|html|en|`layouts/_shortcodes/foo.html`
foo|html|en|`layouts/_shortcodes/foo.html.en.html`

Shortcode name|Output format|Language|Template path
:--|:--|:--|:--
foo|json|en|`layouts/_shortcodes/foo.en.json`
foo|json|en|`layouts/_shortcodes/foo.json`
foo|json|en|`layouts/_shortcodes/foo.json.json`
foo|json|en|`layouts/_shortcodes/foo.json.en.json`

## Methods

Use these methods in your shortcode templates. Refer to each methods's documentation for details and examples.

{{% list-pages-in-section path=/methods/shortcode %}}

## Examples

These examples range in complexity from simple to moderately advanced, with some simplified for clarity.

### Insert year

Create a shortcode to insert the current year:

```go-html-template {file="layouts/_shortcodes/year.html"}
{{- now.Format "2006" -}}
```

Then call the shortcode from within your markup:

```text {file="content/example.md"}
This is {{</* year */>}}, and look at how far we've come.
```

This shortcode can be used inline or as a block on its own line. If a shortcode might be used inline, remove the surrounding [whitespace] by using [template action](g) delimiters with hyphens.

### Insert image

This example assumes the following content structure, where `content/example/index.md` is a [page bundle](g) containing one or more [page resources](g).

```text
content/
├── example/
│   ├── a.jpg
│   └── index.md
└── _index.md
```

Create a shortcode to capture an image as a page resource, resize it to the given width, convert it to the WebP format, and add an `alt` attribute:

```go-html-template {file="layouts/_shortcodes/image.html"}
{{- with .Page.Resources.Get (.Get "path") }}
  {{- with .Process (printf "resize %dx wepb" ($.Get "width")) -}}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="{{ $.Get "alt" }}">
  {{- end }}
{{- end -}}
```

Then call the shortcode from within your markup:

```text {file="content/example/index.md"}
{{</* image path=a.jpg width=300 alt="A white kitten" */>}}
```

The example above uses:

- The [`with`] statement to rebind the [context](g) after each successful operation
- The [`Get`] method to retrieve arguments by name
- The `$` to access the template context

> [!note]
> Make sure that you thoroughly understand the concept of context. The most common templating errors made by new users relate to context.
>
> Read more about context in the [introduction to templating].

### Insert image with error handling

The previous example, while functional, silently fails if the image is missing, and does not gracefully exit if a required argument is missing. We'll add error handling to address these issues:

```go-html-template {file="layouts/_shortcodes/image.html"}
{{- with .Get "path" }}
  {{- with $r := $.Page.Resources.Get ($.Get "path") }}
    {{- with $.Get "width" }}
      {{- with $r.Process (printf "resize %dx wepb" ($.Get "width" )) }}
        {{- $alt := or ($.Get "alt") "" -}}
        <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="{{ $alt }}">
      {{- end }}
    {{- else }}
      {{- errorf "The %q shortcode requires a 'width' argument: see %s" $.Name $.Position }}
    {{- end }}
  {{- else }}
    {{- warnf "The %q shortcode was unable to find %s: see %s" $.Name ($.Get "path") $.Position }}
  {{- end }}
{{- else }}
  {{- errorf "The %q shortcode requires a 'path' argument: see %s" .Name .Position }}
{{- end -}}
```

This template throws an error and gracefully fails the build if the author neglected to provide a `path` or `width` argument, and it emits a warning if it cannot find the image at the specified path. If the author does not provide an `alt` argument, the `alt` attribute is set to an empty string.

The [`Name`] and [`Position`] methods provide helpful context for errors and warnings. For example, a missing `width` argument causes the shortcode to throw this error:

```text
ERROR The "image" shortcode requires a 'width' argument: see "/home/user/project/content/example/index.md:7:1"
```

### Positional arguments

Shortcode arguments can be [named or positional]. We used named arguments previously; let's explore positional arguments. Here's the named argument version of our example:

```text {file="content/example/index.md"}
{{</* image path=a.jpg width=300 alt="A white kitten" */>}}
```

Here's how to call it with positional arguments:

```text {file="content/example/index.md"}
{{</* image a.jpg 300 "A white kitten" */>}}
```

Using the `Get` method with zero-indexed keys, we'll initialize variables with descriptive names in our template:

```go-html-template {file="layouts/_shortcodes/image.html"}
{{ $path := .Get 0 }}
{{ $width := .Get 1 }}
{{ $alt := .Get 2 }}
```

> [!note]
> Positional arguments work well for frequently used shortcodes with one or two arguments. Since you'll use them often, the argument order will be easy to remember. For less frequently used shortcodes, or those with more than two arguments, named arguments improve readability and reduce the chance of errors.

### Named and positional arguments

You can create a shortcode that will accept both named and positional arguments, but not at the same time. Use the [`IsNamedParams`] method to determine whether the shortcode call used named or positional arguments:

```go-html-template {file="layouts/_shortcodes/image.html"}
{{ $path := cond (.IsNamedParams) (.Get "path") (.Get 0) }}
{{ $width := cond (.IsNamedParams) (.Get "width") (.Get 1) }}
{{ $alt := cond (.IsNamedParams) (.Get "alt") (.Get 2) }}
```

This example uses the `cond` alias for the [`compare.Conditional`] function to get the argument by name if `IsNamedParams` returns `true`, otherwise get the argument by position.

### Argument collection

Use the [`Params`] method to access the arguments as a collection.

When using named arguments, the `Params` method returns a map:

```text {file="content/example/index.md"}
{{</* image path=a.jpg width=300 alt="A white kitten" */>}}
```

```go-html-template {file="layouts/_shortcodes/image.html"}
{{ .Params.path }} → a.jpg
{{ .Params.width }} → 300
{{ .Params.alt }} → A white kitten
```

 When using positional arguments, the `Params` method returns a slice:

```text {file="content/example/index.md"}
{{</* image a.jpg 300 "A white kitten" */>}}
```

```go-html-template {file="layouts/_shortcodes/image.html"}
{{ index .Params 0 }} → a.jpg
{{ index .Params 1 }} → 300
{{ index .Params 1 }} → A white kitten
```

Combine the `Params` method with the [`collections.IsSet`] function to determine if a parameter is set, even if its value is falsy.

### Inner content

Extract the content enclosed within shortcode tags using the [`Inner`] method. This example demonstrates how to pass both content and a title to a shortcode. The shortcode then generates a `div` element containing an `h2` element (displaying the title) and the provided content.

```text {file="content/example.md"}
{{</* contrived title="A Contrived Example" */>}}
This is a **bold** word, and this is an _emphasized_ word.
{{</* /contrived  */>}}
```

```go-html-template {file="layouts/_shortcodes/contrived.html"}
<div class="contrived">
  <h2>{{ .Get "title" }}</h2>
  {{ .Inner | .Page.RenderString }}
</div>
```

The preceding example called the shortcode using [standard notation], requiring us to process the inner content with the [`RenderString`] method to convert the Markdown to HTML. This conversion is unnecessary when calling a shortcode using [Markdown notation].

### Nesting

The  [`Parent`] method provides access to the parent shortcode context when the shortcode in question is called within the context of a parent shortcode. This provides an inheritance model.

The following example is contrived but demonstrates the concept. Assume you have a `gallery` shortcode that expects one named `class` argument:

```go-html-template {file="layouts/_shortcodes/gallery.html"}
<div class="{{ .Get "class" }}">
  {{ .Inner }}
</div>
```

You also have an `img` shortcode with a single named `src` argument that you want to call inside of `gallery` and other shortcodes, so that the parent defines the context of each `img`:

```go-html-template {file="layouts/_shortcodes/img.html"}
{{ $src := .Get "src" }}
{{ with .Parent }}
  <img src="{{ $src }}" class="{{ .Get "class" }}-image">
{{ else }}
  <img src="{{ $src }}">
{{ end }}
```

You can then call your shortcode in your content as follows:

```text {file="content/example.md"}
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

### Other examples

For guidance, consider examining Hugo's embedded shortcodes. The source code, available on [GitHub], can provide a useful model.

## Detection

The [`HasShortcode`] method allows you to check if a specific shortcode has been called on a page. For example, consider a custom audio shortcode:

```text {file="content/example.md"}
{{</* audio src=/audio/test.mp3 */>}}
```

You can use the `HasShortcode` method in your base template to conditionally load CSS if the audio shortcode was used on the page:

```go-html-template {file="layouts/baseof.html"}
<head>
  ...
  {{ if .HasShortcode "audio" }}
    <link rel="stylesheet" src="/css/audio.css">
  {{ end }}
  ...
</head>
```

[`collections.IsSet`]: /functions/collections/isset/
[`compare.Conditional`]: /functions/compare/conditional/
[`Get`]: /methods/shortcode/get/
[`HasShortcode`]: /methods/page/hasshortcode/
[`Inner`]: /methods/shortcode/inner/
[`IsNamedParams`]: /methods/shortcode/isnamedparams/
[`Name`]: /methods/shortcode/name/
[`Params`]: /methods/shortcode/params/
[`Parent`]: /methods/shortcode/parent/
[`Position`]: /methods/shortcode/position/
[`RenderString`]: /methods/page/renderstring/
[`with`]: /functions/go-template/with/
[content management]: /content-management/shortcodes/
[embedded shortcodes]: /shortcodes/
[GitHub]: https://github.com/gohugoio/hugo/tree/master/tpl/tplimpl/embedded/templates/_shortcodes
[introduction to templating]: /templates/introduction/
[Markdown notation]: /content-management/shortcodes/#markdown-notation
[named or positional]: /content-management/shortcodes/#arguments
[shortcodes]: /content-management/shortcodes/
[standard notation]: /content-management/shortcodes/#standard-notation
[whitespace]: /templates/introduction/#whitespace
