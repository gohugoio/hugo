---
title: .Param
description: Returns a page parameter, falling back to a site parameter if present.
signature: ['.Param KEY']
categories: [functions]
keywords: ['front matter', 'params']
menu:
  docs:
    parent: 'functions'
aliases: []
---

The `.Param` method on `.Page` looks for the given `KEY` in page parameters, and returns the corresponding value. If it cannot find the `KEY` in page parameters, it looks for the `KEY` in site parameters. If it cannot find the `KEY` in either location, the `.Param` method returns `nil`. 

Site and theme developers commonly set parameters at the site level, allowing content authors to override those parameters at the page level.

For example, to show a table of contents on every page, but allow authors to hide the table of contents as needed:

**Configuration**

{{< code-toggle file="config" copy=false >}}
[params]
display_toc = true
{{< /code-toggle >}}

**Content**

{{< code-toggle file="content/example.md" fm=true copy=false >}}
title = 'Example'
date = 2023-01-01
draft = false
display_toc = false
{{< /code-toggle >}}

**Template**

{{< code file="layouts/_default/single.html" copy="false" >}}
{{ if .Param "display_toc" }}
  {{ .TableOfContents }}
{{ end }}
{{< /code >}}

The `.Param` method returns the value associated with the given `KEY`, regardless of whether the value is truthy or falsy. If you need to ignore falsy values, use this construct instead:

{{< code file="layouts/_default/single.html" copy="false" >}}
{{ or .Params.foo site.Params.foo }}
{{< /code >}}
