---
title: relURL
description: Given a string, prepends the relative URL according to a page's position in the project directory structure.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [urls]
signature: ["relURL INPUT"]
workson: []
hugoversion:
relatedfuncs: [absURL]
deprecated: false
aliases: []
---

Both `absURL` and `relURL` consider the configured value of `baseURL` in your site's [`config` file][configuration]. Given a `baseURL` set to `https://example.com/hugo/`:

```
{{ "mystyle.css" | absURL }} → "https://example.com/hugo/mystyle.css"
{{ "mystyle.css" | relURL }} → "/hugo/mystyle.css"
{{ "http://gohugo.io/" | relURL }} →  "http://gohugo.io/"
{{ "http://gohugo.io/" | absURL }} →  "http://gohugo.io/"
```

The last two examples may look strange but can be very useful. For example, the following shows how to use `absURL` in [JSON-LD structured data for SEO][jsonld] where some of your images for a piece of content may or may not be hosted locally:

{{< code file="layouts/partials/schemaorg-metadata.html" download="schemaorg-metadata.html" >}}
<script type="application/ld+json">
{
    "@context" : "http://schema.org",
    "@type" : "BlogPosting",
    "image" : {{ apply .Params.images "absURL" "." }}
}
</script>
{{< /code >}}

The above uses the [apply function][] and also exposes how the Go template parser JSON-encodes objects inside `<script>` tags. See [the safeJS template function][safejs] for examples of how to tell Hugo not to escape strings inside of such tags.

{{% note "Ending Slash" %}}
`absURL` and `relURL` are smart about missing slashes, but they will *not* add a closing slash to a URL if it is not present.
{{% /note %}}

[apply function]: /functions/apply/
[configuration]: /getting-started/configuration/
[jsonld]: https://developers.google.com/search/docs/guides/intro-structured-data
[safejs]: /functions/safejs
