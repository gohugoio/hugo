---
title: GetPage
description: Returns a Page object from the given path. 
categories: []
keywords: []
action:
  related:
    - methods/site/GetPage
  returnType: page.Page
  signatures: [PAGE.GetPage PATH]
aliases: [/functions/getpage]
---

The `GetPage` method is also available on a `Site` object. See&nbsp;[details].

[details]: /methods/site/getpage/

When using the `GetPage` method on the `Page` object, specify a path relative to the current directory or relative to the content directory.

If Hugo cannot resolve the path to a page, the method returns nil. If the path is ambiguous, Hugo throws an error and fails the build.

Consider this content structure:

```text
content/
├── works/
│   ├── paintings/
│   │   ├── _index.md
│   │   ├── starry-night.md
│   │   └── the-mona-lisa.md
│   ├── sculptures/
│   │   ├── _index.md
│   │   ├── david.md
│   │   └── the-thinker.md
│   └── _index.md
└── _index.md
```

The examples below depict the result of rendering works/paintings/the-mona-lisa.md:

{{< code file=layouts/works/single.html >}}
{{ with .GetPage "starry-night" }}
  {{ .Title }} → Starry Night
{{ end }}

{{ with .GetPage "./starry-night" }}
  {{ .Title }} → Starry Night
{{ end }}

{{ with .GetPage "../paintings/starry-night" }}
  {{ .Title }} → Starry Night
{{ end }}

{{ with .GetPage "/works/paintings/starry-night" }}
  {{ .Title }} → Starry Night
{{ end }}

{{ with .GetPage "../sculptures/david" }}
  {{ .Title }} → David
{{ end }}

{{ with .GetPage "/works/sculptures/david" }}
  {{ .Title }} → David
{{ end }}
{{< /code >}}
