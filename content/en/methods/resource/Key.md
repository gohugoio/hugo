---
title: Key
description: Returns the unique key for the given resource, equivalent to its publishing path.
draft: true
categories: []
keywords: []
action:
  related:
    - methods/resource/Permalink
    - methods/resource/RelPermalink
    - methods/resource/Publish
  returnType: string
  signatures: [RESOURCE.Key]
---

By way of example, consider this site configuration:

{{< code-toggle file=hugo >}}
baseURL = 'https://example.org/docs/'
{{< /code-toggle >}}

And this template:

```go-html-template
  {{ with resources.Get "images/a.jpg" }}
    {{ with resources.Copy "foo/bar/b.jpg" . }}
      {{ .Key }} → foo/bar/b.jpg

      {{ .Name }} → images/a.jpg
      {{ .Title }} → images/a.jpg

      {{ .RelPermalink }} → /docs/foo/bar/b.jpg
    {{ end }}
  {{ end }}
```

We used the [`resources.Copy`] function to change the publishing path. The `Key` method returns the updated path, but note that it is different than the value returned by [`RelPermalink`]. The `RelPermalink` value includes the subdirectory segment of the `baseURL` in the site configuration.

The `Key` method is useful if you need to get the resource's publishing path without publishing the resource. Unlike the `Permalink`, `RelPermalink`, or `Publish` methods, calling `Key` will not publish the resource.


{{% include "methods/resource/_common/global-page-remote-resources.md" %}}

[`Permalink`]: /methods/resource/permalink/
[`RelPermalink`]: /methods/resource/relpermalink/
[`resources.Copy`]: /functions/resources/copy/
