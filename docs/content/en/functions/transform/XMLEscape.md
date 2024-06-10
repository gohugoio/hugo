---
title: transform.XMLEscape
description: Returns the given string, removing disallowed characters then escaping the result to its XML equivalent.
categories: []
keywords: []
action:
  aliases: []
  related: []
  returnType: string
  signatures: [transform.XMLEscape INPUT]
---

{{< new-in 0.121.0 >}}

The `transform.XMLEscape` function removes [disallowed characters] as defined in the XML specification, then escapes the result by replacing the following characters with [HTML entities]:

- `"` → `&#34;`
- `'` → `&#39;`
- `&` → `&amp;`
- `<` → `&lt;`
- `>` → `&gt;`
- `\t` → `&#x9;`
- `\n` → `&#xA;`
- `\r` → `&#xD;`

For example:

```go-html-template
{{ transform.XMLEscape "<p>abc</p>" }} → &lt;p&gt;abc&lt;/p&gt;
```

When using `transform.XMLEscape` in a template rendered by Go's [html/template] package, declare the string to be safe HTML to avoid double escaping. For example, in an RSS template:

{{< code file="layouts/_default/rss.xml" >}}
<description>{{ .Summary | transform.XMLEscape | safeHTML }}</description>
{{< /code >}}

[disallowed characters]: https://www.w3.org/TR/xml/#charsets
[html entities]: https://developer.mozilla.org/en-us/docs/glossary/entity
[html/template]: https://pkg.go.dev/html/template
