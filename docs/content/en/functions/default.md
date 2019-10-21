---
title: default
description: Allows setting a default value that can be returned if a first value is not set.
qref: "Returns a default value if a value is not set when checked."
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
keywords: [defaults]
categories: [functions]
menu:
  docs:
    parent: "functions"
toc:
signature: ["default DEFAULT INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
needsexamples: false
---

`default` checks whether a given value is set and returns a default value if it is not. *Set* in this context means different things depending on the data type:

* non-zero for numeric types and times
* non-zero length for strings, arrays, slices, and maps
* any boolean or struct value
* non-nil for any other types

`default` function examples reference the following content page:

{{< code file="content/posts/default-function-example.md" >}}
---
title: Sane Defaults
seo_title:
date: 2017-02-18
font:
oldparam: The default function helps make your templating DRYer.
newparam:
---
{{< /code >}}

`default` can be written in more than one way:

```
{{ index .Params "font" | default "Roboto" }}
{{ default "Roboto" (index .Params "font") }}
```

Both of the above `default` function calls return `Roboto`.

A `default` value, however, does not need to be hard coded like the previous example. The `default` value can be a variable or pulled directly from the front matter using dot notation:

{{< code file="variable-as-default-value.html" nocopy="true" >}}
{{$old := .Params.oldparam }}
<p>{{ .Params.newparam | default $old }}</p>
{{< /code >}}

Which would return:

```
<p>The default function helps make your templating DRYer.</p>
```

And then using dot notation

{{< code file="dot-notation-default-value.html" >}}
<title>{{ .Params.seo_title | default .Title }}</title>
{{< /code >}}

Which would return

{{< output file="dot-notation-default-return-value.html" >}}
<title>Sane Defaults</title>
{{< /output >}}

The following have equivalent return values but are far less terse. This demonstrates the utility of `default`:

Using `if`:

{{< code file="if-instead-of-default.html" nocopy="true" >}}
<title>{{if .Params.seo_title}}{{.Params.seo_title}}{{else}}{{.Title}}{{end}}</title>
=> Sane Defaults
{{< /code >}}

Using `with`:

{{< code file="with-instead-of-default.html" nocopy="true" >}}
<title>{{with .Params.seo_title}}{{.}}{{else}}{{.Title}}{{end}}</title>
=> Sane Defaults
{{< /code >}}
