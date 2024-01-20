---
title: Page
description: Returns the Page object of the given page.
categories: []
keywords: []
action:
  related: []
  returnType: page.Page
  signatures: [PAGE.Page]
---

This is a convenience method, useful within partial templates that are called from both [shortcodes] and page templates.

{{< code file=layouts/shortcodes/foo.html  >}}
{{ partial "my-partial.html" . }}
{{< /code >}}

When the shortcode calls the partial, it passes the current [context] (the dot). The context includes identifiers such as `Page`, `Params`, `Inner`, and `Name`.

{{< code file=layouts/_default/single.html  >}}
{{ partial "my-partial.html" . }}
{{< /code >}}

When the page template calls the partial, it also passes the current context (the dot). But in this case, the dot _is_ the `Page` object.

{{< code file=layouts/partials/my-partial.html  >}}
The page title is: {{ .Page.Title }}
{{< /code >}}

To handle both scenarios, the partial template must be able to access the `Page` object with `Page.Page`.

{{% note %}}
And yes, that means you can do `.Page.Page.Page.Page.Title` too.

But don't.
{{% /note %}}


[context]: getting-started/glossary/#context
[shortcodes]: /getting-started/glossary/#shortcode
