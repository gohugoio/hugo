---
title: Code block render hooks
linkTitle: Code blocks
description: Create a code block render hook to override the rendering of Markdown code blocks to HTML.
categories: [render hooks]
keywords: []
menu:
  docs:
    parent: render-hooks
    weight: 40
weight: 40
toc: true
---

## Markdown

This Markdown example contains a fenced code block:

{{< code file=content/example.md lang=text >}}
```bash {class="my-class" id="my-codeblock" lineNos=inline tabWidth=2}
declare a=1
echo "$a"
exit
```
{{< /code >}}

A fenced code block consists of:

- A leading [code fence]
- An optional [info string]
- A code sample
- A trailing code fence

[code fence]: https://spec.commonmark.org/0.31.2/#code-fence
[info string]: https://spec.commonmark.org/0.31.2/#info-string

In the previous example, the info string contains:

- The language of the code sample (the first word)
- An optional space-delimited or comma-delimited list of attributes (everything within braces)

The attributes in the info string can be generic attributes or highlighting options.

In the example above, the _generic attributes_ are `class` and `id`. In the absence of special handling within a code block render hook, Hugo adds each generic attribute to the HTML element surrounding the rendered code block. Consistent with its content security model, Hugo removes HTML event attributes such as `onclick` and `onmouseover`. Generic attributes are typically global HTML attributes, but you may include custom attributes as well.

In the example above, the _highlighting options_ are `lineNos` and `tabWidth`. Hugo uses the [Chroma] syntax highlighter to render the code sample. You can control the appearance of the rendered code by specifying one or more [highlighting options].

[Chroma]: https://github.com/alecthomas/chroma/
[highlighting options]: /functions/transform/highlight/#options

{{% note %}}
Although `style` is a global HTML attribute, when used in an info string it is a highlighting option.
{{% /note %}}

## Context

Code block render hook templates receive the following [context]:

[context]: /getting-started/glossary/#context

###### Attributes

(`map`) The generic attributes from the info string.

###### Inner

(`string`) The content between the leading and trailing code fences, excluding the info string.

###### Options

(`map`) The highlighting options from the info string.

###### Ordinal

(`int`) The zero-based ordinal of the code block on the page.

###### Page

(`page`) A reference to the current page.

###### PageInner

{{< new-in 0.125.0 >}}

(`page`) A reference to a page nested via the [`RenderShortcodes`] method. [See details](#pageinner-details).

[`RenderShortcodes`]: /methods/page/rendershortcodes

###### Position

(`text.Position`) The position of the code block within the page content.

###### Type

(`string`) The first word of the info string, typically the code language.

## Examples

In its default configuration, Hugo renders fenced code blocks by passing the code sample through the Chroma syntax highlighter and wrapping the result. To create a render hook that does the same thing:

[CommonMark specification]: https://spec.commonmark.org/current/

{{< code file=layouts/_default/_markup/render-codeblock.html copy=true >}}
{{ $result := transform.HighlightCodeBlock . }}
{{ $result.Wrapped }}
{{< /code >}}

Although you can use one template with conditional logic to control the behavior on a per-language basis, you can also create language-specific templates.

```text
layouts/
└── _default/
    └── _markup/
        ├── render-codeblock-mermaid.html
        ├── render-codeblock-python.html
        └── render-codeblock.html
```

For example, to create a code block render hook to render [Mermaid] diagrams:

[Mermaid]: https://mermaid.js.org/

{{< code file=layouts/_default/_markup/render-codeblock-mermaid.html copy=true >}}
<pre class="mermaid">
  {{- .Inner | safeHTML }}
</pre>
{{ .Page.Store.Set "hasMermaid" true }}
{{< /code >}}

Then include this snippet at the bottom of the your base template:

{{< code file=layouts/_default/baseof.html copy=true >}}
{{ if .Store.Get "hasMermaid" }}
  <script type="module">
    import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.esm.min.mjs';
    mermaid.initialize({ startOnLoad: true });
  </script>
{{ end }}
{{< /code >}}

See the [diagrams] page for details.

[diagrams]: /content-management/diagrams/#mermaid-diagrams

## Embedded

Hugo includes an [embedded code block render hook] to render [GoAT diagrams].

[embedded code block render hook]: {{% eturl render-codeblock-goat %}}
[GoAT diagrams]: /content-management/diagrams/#goat-diagrams-ascii

{{% include "/render-hooks/_common/pageinner.md" %}}
