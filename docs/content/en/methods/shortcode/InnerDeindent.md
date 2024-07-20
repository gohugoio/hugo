---
title: InnerDeindent
description: Returns the content between opening and closing shortcode tags, with indentation removed, applicable when the shortcode call includes a closing tag. 
categories: []
keywords: []
action:
  related:
    - methods/shortcode/Inner
  returnType: template.HTML
  signatures: [SHORTCODE.InnerDeindent]
---

Similar to the [`Inner`] method, `InnerDeindent` returns the content between opening and closing shortcode tags. However, with `InnerDeindent`, indentation before the content is removed.

This allows us to effectively bypass the rules governing [indentation] as provided in the [CommonMark] specification.

Consider this Markdown, an unordered list with a small gallery of thumbnail images within each list item:

{{< code file=content/about.md lang=md >}}
- Gallery one

    {{</* gallery */>}}
    ![kitten a](thumbnails/a.jpg)
    ![kitten b](thumbnails/b.jpg)
    {{</* /gallery */>}}

- Gallery two

    {{</* gallery */>}}
    ![kitten c](thumbnails/c.jpg)
    ![kitten d](thumbnails/d.jpg)
    {{</* /gallery */>}}
{{< /code >}}

In the example above, notice that the content between the opening and closing shortcode tags is indented by four spaces. Per the CommonMark specification, this is treated as an indented code block.

With this shortcode, calling `Inner` instead of `InnerDeindent`:

{{< code file=layouts/shortcodes/gallery.html  >}}
<div class="gallery">
  {{ trim .Inner "\r\n" | .Page.RenderString }}
</div>
{{< /code >}}

Hugo renders the Markdown to:

```html
<ul>
  <li>
    <p>Gallery one</p>
    <div class="gallery">
      <pre><code>![kitten a](images/a.jpg)
      ![kitten b](images/b.jpg)
      </code></pre>
    </div>
  </li>
  <li>
    <p>Gallery two</p>
    <div class="gallery">
      <pre><code>![kitten c](images/c.jpg)
      ![kitten d](images/d.jpg)
      </code></pre>
    </div>
  </li>
</ul>
```

Although technically correct per the CommonMark specification, this is not what we want. If we remove the indentation using the `InnerDeindent` method:

{{< code file=layouts/shortcodes/gallery.html  >}}
<div class="gallery">
  {{ trim .InnerDeindent "\r\n" | .Page.RenderString }}
</div>
{{< /code >}}

Hugo renders the Markdown to:

```html
<ul>
  <li>
    <p>Gallery one</p>
    <div class="gallery">
      <img src="images/a.jpg" alt="kitten a">
      <img src="images/b.jpg" alt="kitten b">
    </div>
  </li>
  <li>
    <p>Gallery two</p>
    <div class="gallery">
      <img src="images/c.jpg" alt="kitten c">
      <img src="images/d.jpg" alt="kitten d">
    </div>
  </li>
</ul>
```

[commonmark]: https://commonmark.org/
[indentation]: https://spec.commonmark.org/0.30/#indented-code-blocks
[`Inner`]: /methods/shortcode/inner/
