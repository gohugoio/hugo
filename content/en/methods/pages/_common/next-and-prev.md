---
# Do not remove front matter.
---

Hugo determines the _next_ and _previous_ page by sorting the page collection according to this sorting hierarchy:

Field|Precedence|Sort direction
:--|:--|:--
[`weight`]|1|descending
[`date`]|2|descending
[`linkTitle`]|3|descending
[`path`]|4|descending

[`date`]: /methods/page/date/
[`weight`]: /methods/page/weight/
[`linkTitle`]: /methods/page/linktitle/
[`path`]: /methods/page/path/

The sorted page collection used to determine the _next_ and _previous_ page is independent of other page collections, which may lead to unexpected behavior.

For example, with this content structure:

```text
content/
├── pages/
│   ├── _index.md
│   ├── page-1.md   <-- front matter: weight = 10
│   ├── page-2.md   <-- front matter: weight = 20
│   └── page-3.md   <-- front matter: weight = 30
└── _index.md
```

And these templates:

{{< code file=layouts/_default/list.html >}}
{{ range .Pages.ByWeight}}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
{{< /code >}}

{{< code file=layouts/_default/single.html >}}
{{ $pages := .CurrentSection.Pages.ByWeight }}

{{ with $pages.Prev . }}
  <a href="{{ .RelPermalink }}">Previous</a>
{{ end }}

{{ with $pages.Next . }}
  <a href="{{ .RelPermalink }}">Next</a>
{{ end }}
{{< /code >}}

When you visit page-2:

- The `Prev` method points to page-3
- The `Next` method points to page-1

To reverse the meaning of _next_ and _previous_ you can chain the [`Reverse`] method to the page collection definition:

{{< code file=layouts/_default/single.html >}}
{{ $pages := .CurrentSection.Pages.ByWeight.Reverse }}

{{ with $pages.Prev . }}
  <a href="{{ .RelPermalink }}">Previous</a>
{{ end }}

{{ with $pages.Next . }}
  <a href="{{ .RelPermalink }}">Next</a>
{{ end }}
{{< /code >}}

[`Reverse`]: /methods/pages/reverse/
