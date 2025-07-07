---
title: Inner
description: Returns the content between opening and closing shortcode tags, applicable when the shortcode call includes a closing tag.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: template.HTML
    signatures: [SHORTCODE.Inner]
---

This content:

```text {file="content/services.md"}
{{</* card title="Product Design" */>}}
We design the **best** widgets in the world.
{{</* /card */>}}
```

With this shortcode:

```go-html-template {file="layouts/_shortcodes/card.html"}
<div class="card">
  {{ with .Get "title" }}
    <div class="card-title">{{ . }}</div>
  {{ end }}
  <div class="card-content">
    {{ .Inner | strings.TrimSpace }}
  </div>
</div>
```

Is rendered to:

```html
<div class="card">
  <div class="card-title">Product Design</div>
  <div class="card-content">
    We design the **best** widgets in the world.
  </div>
</div>
```

> [!note]
> Content between opening and closing shortcode tags may include leading and/or trailing newlines, depending on placement within the Markdown. Use the [`strings.TrimSpace`] function as shown above to remove carriage returns and newlines.

> [!note]
> In the example above, the value returned by `Inner` is Markdown, but it was rendered as plain text. Use either of the following approaches to render Markdown to HTML.

## Use RenderString

Let's modify the example above to pass the value returned by `Inner` through the [`RenderString`] method on the `Page` object:

```go-html-template {file="layouts/_shortcodes/card.html"}
<div class="card">
  {{ with .Get "title" }}
    <div class="card-title">{{ . }}</div>
  {{ end }}
  <div class="card-content">
    {{ .Inner | strings.TrimSpace | .Page.RenderString }}
  </div>
</div>
```

Hugo renders this to:

```html
<div class="card">
  <div class="card-title">Product design</div>
  <div class="card-content">
    We produce the <strong>best</strong> widgets in the world.
  </div>
</div>
```

You can use the [`markdownify`] function instead of the `RenderString` method, but the latter is more flexible. See&nbsp;[details].

## Alternative notation

Instead of calling the shortcode with the `{{</* */>}}` notation, use the `{{%/* */%}}` notation:

```text {file="content/services.md"}
{{%/* card title="Product Design" */%}}
We design the **best** widgets in the world.
{{%/* /card */%}}
```

When you use the `{{%/* */%}}` notation, Hugo renders the entire shortcode as Markdown, requiring the following changes.

First, configure the renderer to allow raw HTML within Markdown:

{{< code-toggle file=hugo >}}
[markup.goldmark.renderer]
unsafe = true
{{< /code-toggle >}}

This configuration is not unsafe if _you_ control the content. Read more about Hugo's [security model].

Second, because we are rendering the entire shortcode as Markdown, we must adhere to the rules governing [indentation] and inclusion of [raw HTML blocks] as provided in the [CommonMark] specification.

```go-html-template {file="layouts/_shortcodes/card.html"}
<div class="card">
  {{ with .Get "title" }}
  <div class="card-title">{{ . }}</div>
  {{ end }}
  <div class="card-content">

  {{ .Inner | strings.TrimSpace }}
  </div>
</div>
```

The difference between this and the previous example is subtle but required. Note the change in indentation, the addition of a blank line, and removal of the `RenderString` method.

```diff
--- layouts/_shortcodes/a.html
+++ layouts/_shortcodes/b.html
@@ -1,8 +1,9 @@
 <div class="card">
   {{ with .Get "title" }}
-    <div class="card-title">{{ . }}</div>
+  <div class="card-title">{{ . }}</div>
   {{ end }}
   <div class="card-content">
-    {{ .Inner | strings.TrimSpace | .Page.RenderString }}
+
+  {{ .Inner | strings.TrimSpace }}
   </div>
 </div>
```

> [!note]
> Don't process the `Inner` value with `RenderString` or `markdownify` when using [Markdown notation] to call the shortcode.

[`markdownify`]: /functions/transform/markdownify/
[`RenderString`]: /methods/page/renderstring/
[`strings.TrimSpace`]: /functions/strings/trimspace/
[CommonMark]: https://spec.commonmark.org/current/
[details]: /methods/page/renderstring/
[indentation]: https://spec.commonmark.org/current/#indented-code-blocks
[Markdown notation]: /content-management/shortcodes/#notation
[raw HTML blocks]: https://spec.commonmark.org/current/#html-blocks
[security model]: /about/security/
