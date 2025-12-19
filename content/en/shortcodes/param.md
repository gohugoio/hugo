---
title: Param shortcode
linkTitle: Param
description: Insert a parameter from front matter or site configuration into your content using the param shortcode.
categories: []
keywords: []
---

> [!note]
> To override Hugo's embedded `param` shortcode, copy the [source code] to a file with the same name in the `layouts/_shortcodes` directory.

The `param` shortcode renders a parameter from front matter, falling back to a site parameter of the same name. The shortcode throws an error if the parameter does not exist.

```text {file="content/example.md"}
---
title: Example
date: 2025-01-15T23:29:46-08:00
params:
  color: red
  size: medium
---

We found a {{%/* param "color" */%}} shirt.
```

Hugo renders this to:

```html
<p>We found a red shirt.</p>
```

Access nested values by [chaining](g) the [identifiers](g):

```text
{{%/* param my.nested.param */%}}
```

[source code]: <{{% eturl param %}}>
