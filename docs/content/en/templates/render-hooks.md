---
title: "Markdown Render Hooks"
linkTitle: "Render Hooks"
description: "Render Hooks allow custom templates to override markdown rendering functionality."
date: 2017-03-11
categories: [templates]
keywords: [markdown]
toc: true
menu:
  docs:
    title: "Markdown Render Hooks"
    parent: "templates"
    weight: 20
---

Note that this is only supported with the [Goldmark](/getting-started/configuration-markup#goldmark) renderer.

You can override certain parts of the default Markdown rendering to HTML by creating templates with base names `render-{kind}` in `layouts/_default/_markup`.

You can also create type/section specific hooks in `layouts/[type/section]/_markup`, e.g.: `layouts/blog/_markup`.

The hook kinds currently supported are:

* `image`
* `link`
* `heading`
* `codeblock`{{< new-in "0.93.0" >}}
* `list`{{< new-in "0.99.0" >}}
* `listitem`{{< new-in "0.99.0" >}}

You can define [Output-Format-](/templates/output-formats) and [language-](/content-management/multilingual/)specific templates if needed. Your `layouts` folder may look like this:

```text
layouts/
└── _default/
    └── _markup/
        ├── render-codeblock-bash.html
        ├── render-codeblock.html
        ├── render-heading.html
        ├── render-image.html
        ├── render-image.rss.xml
        └── render-link.html
```

Some use cases for the above:

* Resolve link references using `.GetPage`. This would make links portable as you could translate `./my-post.md` (and similar constructs that would work on GitHub) into `/blog/2019/01/01/my-post/` etc.
* Add `target=_blank` to external links.
* Resolve and [process](/content-management/image-processing/) images.
* Add [header links](https://remysharp.com/2014/08/08/automatic-permalinks-for-blog-posts).

## Render Hooks for Headings, Links and Images

The `render-link` and `render-image` templates will receive this context:

Page
: The [Page](/variables/page/) being rendered.

Destination
: The URL.

Title
: The title attribute.

Text
: The rendered (HTML) link text.

PlainText
: The plain variant of the above.

The `render-heading` template will receive this context:

Page
: The [Page](/variables/page/) being rendered.

Level
: The header level (1--6)

Anchor
: An auto-generated html id unique to the header within the page

Text
: The rendered (HTML) text.

PlainText
: The plain variant of the above.

Attributes (map)
: A map of attributes (e.g. `id`, `class`). Note that this will currently always be empty for links.

The `render-image` templates will also receive:

IsBlock {{< new-in "0.108.0" >}}
: Returns true if this is a standalone image and the config option [markup.goldmark.parser.wrapStandAloneImageWithinParagraph](/getting-started/configuration-markup/#goldmark) is disabled.

Ordinal  {{< new-in "0.108.0" >}}
: Zero-based ordinal for all the images in the current document.


### Link with title Markdown example

```md
[Text](https://www.gohugo.io "Title")
```

Here is a code example for how the render-link.html template could look:

{{< code file="layouts/_default/_markup/render-link.html" >}}
<a href="{{ .Destination | safeURL }}"{{ with .Title}} title="{{ . }}"{{ end }}{{ if strings.HasPrefix .Destination "http" }} target="_blank" rel="noopener"{{ end }}>{{ .Text | safeHTML }}</a>
{{< /code >}}

### Image Markdown example

```md
![Text](https://gohugo.io/images/hugo-logo-wide.svg "Title")
```

Here is a code example for how the render-image.html template could look:

{{< code file="layouts/_default/_markup/render-image.html" >}}
<p class="md__image">
  <img src="{{ .Destination | safeURL }}" alt="{{ .Text }}" {{ with .Title}} title="{{ . }}"{{ end }} />
</p>
{{< /code >}}

### Heading link example

Given this template file

{{< code file="layouts/_default/_markup/render-heading.html" >}}
<h{{ .Level }} id="{{ .Anchor | safeURL }}">{{ .Text | safeHTML }} <a href="#{{ .Anchor | safeURL }}">¶</a></h{{ .Level }}>
{{< /code >}}

And this markdown

```md
### Section A
```

The rendered html will be

```html
<h3 id="section-a">Section A <a href="#section-a">¶</a></h3>
```

## Render Hooks for Code Blocks

{{< new-in "0.93.0" >}}

You can add a hook template for either all code blocks or for a specific type/language (`bash` in the example below):

```goat { class="black f7" }
layouts
└── _default
    └── _markup
        └── render-codeblock.html
        └── render-codeblock-bash.html
```

The default behavior for these code blocks is to do [Code Highlighting](/content-management/syntax-highlighting/#highlighting-in-code-fences), but since you can pass attributes to these code blocks, they can be used for almost anything. One example would be the built-in [GoAT Diagrams](/content-management/diagrams/#goat-diagrams-ascii) or this [Mermaid Diagram Code Block Hook](/content-management/diagrams/#mermaid-diagrams) example.

The context (the ".") you receive in a code block template contains:

Type (string)
: The type of code block. This will be the programming language, e.g. `bash`, when doing code highlighting.

Attributes (map)
: Attributes passed in from Markdown (e.g. `{ attrName1=attrValue1 attrName2="attr Value 2" }`).

Options (map)
: Chroma highlighting processing options. This will only be filled if `Type` is a known [Chroma Lexer](/content-management/syntax-highlighting/#list-of-chroma-highlighting-languages).

Inner (string)
: The text between the code fences.

Ordinal (integer)
: Zero-based ordinal for all code blocks in the current document.

Page
: The owning `Page`.

Position
: Useful in error logging as it prints the filename and position (linenumber, column), e.g. `{{ errorf "error in code block: %s" .Position }}`.

## Render Hooks for Lists and List Items

{{< new-in "0.99.0" >}}

You can add a hook template for lists and list items e.g. for different output types

```goat { class="black f7" }
layouts
└── _default
    └── _markup
        └── render-listitem.html
        └── render-listitem.json
        └── render-list.html
        └── render-list.json
```

The `render-list` template will receive this context:

Page
: The [Page](/variables/page/) being rendered.

Text
: The rendered (HTML) list content.

PlainText
: The plain variant of the above.

IsOrdered (bool)
: If this is an ordered list.

Parent 
: The Parent node of the list.

Attributes (map) {{< new-in "0.82.0" >}}
: A map of attributes (e.g. `id`, `class`)


The `render-listitem` template will receive this context:

Page
: The [Page](/variables/page/) being rendered.

Text
: The rendered (HTML) list content.

PlainText
: The plain variant of the above.

IsFirst (bool)
: If this is the first item in the list.

IsLast (bool)
: If this is the last item in the list.

Parent 
: The Parent node of the item, the list node.

### ListItem rendered as JSON-LD example:

```md
1. Do This
2. Then That
```

Here is a code example for how the render-listitem.json template could look:

{{< code file="layouts/_default/_markup/render-list.html" >}}
{{- if eq .Parent.IsOrdered true -}}
{
    "@type": "HowToStep",
    "text": "{{ .Text | plainify}}"
}{{ if not .IsLast }},{{ end }}{{ print "\n"}}
{{- else -}}
{{- if not .IsLast -}}
    {{ printf "\"%s\",\n" .Text | plainify}}
{{- else -}}
    {{ printf "\"%s\"" .Text | plainify}}
{{- end -}}
{{- end -}}
{{< /code >}}


The rendered html will be

```js
{
    "@type": "HowToStep",
    "text": "Do This"
},
{
    "@type": "HowToStep",
    "text": "Then That"
}
```