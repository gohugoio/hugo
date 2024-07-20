---
# Do not remove front matter.
---

## PageInner details

{{< new-in 0.125.0 >}}

The primary use case for `PageInner` is to resolve links and [page resources] relative to an included `Page`. For example, create an "include" shortcode to compose a page from multiple content files, while preserving a global context for footnotes and the table of contents:

{{< code file=layouts/shortcodes/include.html >}}
{{ with site.GetPage (.Get 0) }}
  {{ .RenderShortcodes }}
{{ end }}
{{< /code >}}

Then call the shortcode in your Markdown:

{{< code file=content/posts/p1.md >}}
{{%/* include "/posts/p2" */%}}
{{< /code >}}

Any render hook triggered while rendering `/posts/p2` will get:

- `/posts/p1` when calling `Page`
- `/posts/p2` when calling `PageInner`

`PageInner` falls back to the value of `Page` if not relevant, and always returns a value.

{{% note %}}
The `PageInner` method is only relevant for shortcodes that invoke the [`RenderShortcodes`] method, and you must call the shortcode using the `{{%/*..*/%}}` notation.

[`RenderShortcodes`]: /methods/page/rendershortcodes/
{{% /note %}}

As a practical example, Hugo's embedded link and image render hooks use the `PageInner` method to resolve markdown link and image destinations. See the source code for each:

- [Embedded link render hook]({{% eturl render-link %}})
- [Embedded image render hook]({{% eturl render-image %}})

[`RenderShortcodes`]: /methods/page/rendershortcodes/
[page resources]: /getting-started/glossary/#page-resource
