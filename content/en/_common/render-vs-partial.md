---
_comment: Do not remove front matter.
---

`Render` method|`partial` function
:--|:--
By default, the `Page` object is the context. You may pass an optional `CONTEXT` argument to replace it, allowing you to pass a combination of objects, slices, maps, and scalars.|You must specify the context, allowing you to pass a combination of objects, slices, maps, and scalars.
Hugo resolves the template automatically via the [template lookup order][], and can target any page kind, content type, logical path, language, or output format.|Hugo does not consider the current page kind, content type, logical path, language, or output format when searching for a matching template.
Templates may reside at any level within the `layouts` directory.|Templates must reside within the `layouts/_partials` directory.
There is no cached variant.|The [`partialCached`][] function is a cached variant.

[`partialCached`]: /functions/partials/includeCached/
[template lookup order]: /templates/lookup-order/
