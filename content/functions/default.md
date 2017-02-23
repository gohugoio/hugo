---
title: default
linktitle: default
description:
qref: "Returns a default value if a value is not set when checked."
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [defaults]
categories: [functions]
toc: false
draft: false
aliases: [/functions/default/]
---

Checks whether a given value is set and returns a default value if it is not. *Set* in this context means different things depending on date type:

* non-zero for numeric types and times
* non-zero length for strings, arrays, slices, and maps
* any boolean or struct value
* non-nil for any other types

`default` function examples reference the following content page:

{{% input "content/posts/default-function-example.md" %}}
```yaml
---
title: Sane Defaults
seo_title:
date: 2017-02-18
font:
oldparam: The default function helps make your templating DRYer.
newparam:
---
```
{{% /input %}}

`default` can be written in more than one way:

```golang
{{ index .Params "font" | default "Roboto" }}
{{ default "Roboto" (index .Params "font") }}
```

Both of the above `default` function calls return `Roboto`.

A `default` value, however, does not need to be hard coded like the previous example. The `default` value can be a variable or pulled directly from the front matter using dot notation:

{{% input "variable-as-default-value.html" nocopy %}}
```golang
{{$old := .Params.oldparam }}
<p>{{ .Params.newparam | default $old }}</p>
```
{{% /input %}}

Which would return:

```
<p>The default function helps make your templating DRYer.</p>
```

And then using dot notation

{{% input "dot-notation-default-value.html" %}}
```golang
<title>{{ .Params.seo_title | default .Title }}</title>
```
{{% /input %}}

Which would return

{{% output "dot-notation-default-return-value.html" %}}
```html
<title>Sane Defaults</title>
```
{{% /output %}}

The following have equivalent return values but are far less terse. This demonstrates the utility of `default`:

Using `if`:

{{% input "if-instead-of-default.html" nocopy %}}
```golang
<title>{{if .Params.seo_title}}{{.Params.seo_title}}{{else}}{{.Title}}{{end}}</title>
=> Sane Defaults
```
{{% /input %}}

Using `with`:

{{% input "with-instead-of-default.html" nocopy %}}
```golang
<title>{{with .Params.seo_title}}{{.}}{{else}}{{.Title}}{{end}}</title>
=> Sane Defaults
```
{{% /input %}}





