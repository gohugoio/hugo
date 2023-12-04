---
title: transform.Remarshal
description: Marshals a string of serialized data, or a map, into a string of serialized data in the specified format.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/encoding/Jsonify
    - functions/transform/Unmarshal
  returnType: string
  signatures: [transform.Remarshal FORMAT INPUT]
aliases: [/functions/transform.remarshal]
---

The format must be one of `json`, `toml`, `yaml`, or `xml`. If the input is a string of serialized data, it must be valid JSON, TOML, YAML, or XML.

{{% note %}}
This function is primarily a helper for Hugo's documentation, used to convert configuration and front matter examples to JSON, TOML, and YAML.

This is not a general purpose converter, and may change without notice if required for Hugo's documentation site.
{{% /note %}}

Example 1
: Convert a string of TOML to JSON.

```go-html-template
{{ $s := `
  baseURL = 'https://example.org/'
  languageCode = 'en-US'
  title = 'ABC Widgets'
`}}
<pre>{{ transform.Remarshal "json" $s }}</pre>
```

Resulting HTML:

```html
<pre>{
   &#34;baseURL&#34;: &#34;https://example.org/&#34;,
   &#34;languageCode&#34;: &#34;en-US&#34;,
   &#34;title&#34;: &#34;ABC Widgets&#34;
}
</pre>
```

Rendered in browser:

```text
{
   "baseURL": "https://example.org/",
   "languageCode": "en-US",
   "title": "ABC Widgets"
}
```

Example 2
: Convert a map to YAML.

```go-html-template
{{ $m := dict
  "a" "Hugo rocks!"
  "b" (dict "question" "What is 6x7?" "answer" 42)
  "c" (slice "foo" "bar")
}}
<pre>{{ transform.Remarshal "yaml" $m }}</pre>
```

Resulting HTML:

```html
<pre>a: Hugo rocks!
b:
  answer: 42
  question: What is 6x7?
c:
- foo
- bar
</pre>
```

Rendered in browser:

```text
a: Hugo rocks!
b:
  answer: 42
  question: What is 6x7?
c:
- foo
- bar
```
