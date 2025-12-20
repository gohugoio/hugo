---
_comment: Do not remove front matter.
---

## PageInner details

{{< new-in 0.125.0 />}}

The primary use case for `PageInner` is to resolve links and [page resources](g) relative to an included `Page`. For example, create an "include" shortcode to compose a page from multiple content files, while preserving a global context for footnotes and the table of contents:

```go-html-template {file="layouts/_shortcodes/include.html" copy=true}
{{ with .Get 0 }}
  {{ with $.Page.GetPage . }}
    {{- .RenderShortcodes }}
  {{ else }}
    {{ errorf "The %q shortcode was unable to find %q. See %s" $.Name . $.Position }}
  {{ end }}
{{ else }}
  {{ errorf "The %q shortcode requires a positional parameter indicating the logical path of the file to include. See %s" .Name .Position }}
{{ end }}
```

Then call the shortcode in your Markdown:

```text {file="content/posts/post-1.md"}
{{%/* include "/posts/post-2" */%}}
```

Any render hook triggered while rendering `/posts/post-2` will get:

- `/posts/post-1` when calling `Page`
- `/posts/post-2` when calling `PageInner`

`PageInner` falls back to the value of `Page` if not relevant, and always returns a value.

> [!note]
> The `PageInner` method is only relevant for shortcodes that invoke the [`RenderShortcodes`] method, and you must call the shortcode using [Markdown notation].

As a practical example, Hugo's embedded link and image render hooks use the `PageInner` method to resolve markdown link and image destinations. See the source code for each:

- [Embedded link render hook]
- [Embedded image render hook]

[`RenderShortcodes`]: /methods/page/rendershortcodes/
[Markdown notation]: /content-management/shortcodes/#notation
[Embedded link render hook]: <{{% eturl render-link %}}>
[Embedded image render hook]: <{{% eturl render-image %}}>
