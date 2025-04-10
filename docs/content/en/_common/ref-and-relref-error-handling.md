---
_comment: Do not remove front matter.
---

By default, Hugo will throw an error and fail the build if it cannot resolve the path. You can change this to a warning in your site configuration, and specify a URL to return when the path cannot be resolved.

{{< code-toggle file=hugo >}}
refLinksErrorLevel = 'warning'
refLinksNotFoundURL = '/some/other/url'
{{< /code-toggle >}}
