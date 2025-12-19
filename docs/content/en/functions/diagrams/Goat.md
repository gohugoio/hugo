---
title: diagrams.Goat
description: Returns an SVGDiagram object created from the given GoAT markup and options.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: diagrams.SVGDiagram
    signatures: [diagrams.Goat MARKUP]
---

Useful in a [code block render hook], the `diagrams.Goat` function returns an SVGDiagram object created from the given [GoAT] markup.

## Methods

The SVGDiagram object has the following methods:

Inner
: (`template.HTML`) Returns the SVG child elements without a wrapping `svg` element, allowing you to create your own wrapper.

Wrapped
: (`template.HTML`) Returns the SVG child elements wrapped in an `svg` element.

Width
: (`int`) Returns the width of the rendered diagram, in pixels.

Height
: (`int`) Returns the height of the rendered diagram, in pixels.

## GoAT Diagrams

Hugo natively supports GoAT diagrams with an [embedded code block render hook].

This Markdown:

````text
```goat
.---.     .-.       .-.       .-.     .---.
| A +--->| 1 |<--->| 2 |<--->| 3 |<---+ B |
'---'     '-'       '+'       '+'     '---'
```
````

Is rendered to:

```html
<div class="goat svg-container">
  <svg xmlns="http://www.w3.org/2000/svg" font-family="Menlo,Lucida Console,monospace" viewBox="0 0 352 57">
    ...
  </svg>
</div>
```

Which appears in your browser as:

```goat {class="mw6-ns"}
.---.     .-.       .-.       .-.     .---.
| A +--->| 1 |<--->| 2 |<--->| 3 |<---+ B |
'---'     '-'       '+'       '+'     '---'
```

To customize rendering, override Hugo's [embedded code block render hook] for GoAT diagrams.

## Code block render hook

By way of example, let's create a code block render hook to render GoAT diagrams as `figure` elements with an optional caption.

```go-html-template {file="layouts/_markup/render-codeblock-goat.html"}
{{ $caption := or .Attributes.caption "" }}
{{ $class := or .Attributes.class "diagram" }}
{{ $id := or .Attributes.id (printf "diagram-%d" (add 1 .Ordinal)) }}

<figure id="{{ $id }}">
  {{ with diagrams.Goat (trim .Inner "\n\r") }}
    <svg class="{{ $class }}" width="{{ .Width }}" height="{{ .Height }}"  xmlns="http://www.w3.org/2000/svg" version="1.1">
      {{ .Inner }}
    </svg>
  {{ end }}
  <figcaption>{{ $caption }}</figcaption>
</figure>
```

This Markdown:

````text {file="content/example.md" }
```goat {class="foo" caption="Diagram 1: Example"}
.---.     .-.       .-.       .-.     .---.
| A +--->| 1 |<--->| 2 |<--->| 3 |<---+ B |
'---'     '-'       '+'       '+'     '---'
```
````

Is rendered to:

```html
<figure id="diagram-1">
  <svg class="foo" width="272" height="57" xmlns="http://www.w3.org/2000/svg" version="1.1">
    ...
  </svg>
  <figcaption>Diagram 1: Example</figcaption>
</figure>
```

Use CSS to style the SVG as needed:

```css
svg.foo {
  font-family: "Segoe UI","Noto Sans",Helvetica,Arial,sans-serif
}
```

[code block render hook]: /render-hooks/code-blocks/
[embedded code block render hook]: <{{% eturl render-codeblock-goat %}}>
[GoAT]: https://github.com/bep/goat
