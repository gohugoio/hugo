---
title: transform.Unmarshal
description: "`transform.Unmarshal` (alias `unmarshal`) parses the input and converts it into a map or an array. Supported formats are JSON, TOML, YAML, XML and CSV."
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
signature: ["RESOURCE or STRING | transform.Unmarshal [OPTIONS]"]
---

The function accepts either a `Resource` created in [Hugo Pipes](/hugo-pipes/) or via [Page Bundles](/content-management/page-bundles/), or simply a string. The two examples below will produce the same map:

```go-html-template
{{ $greetings := "hello = \"Hello Hugo\"" | transform.Unmarshal }}`
```

```go-html-template
{{ $greetings := "hello = \"Hello Hugo\"" | resources.FromString "data/greetings.toml" | transform.Unmarshal }}
```

In both the above examples, you get a map you can work with:

```go-html-template
{{ $greetings.hello }}
```

The above prints `Hello Hugo`.

## CSV options

Unmarshal with CSV as input has some options you can set:

delimiter
: The delimiter used, default is `,`.

comment
: The comment character used in the CSV. If set, lines beginning with the comment character without preceding whitespace are ignored.:

Example:

```go-html-template
{{ $csv := "a;b;c" | transform.Unmarshal (dict "delimiter" ";") }}
```

## XML data

As a convenience, Hugo allows you to access XML data in the same way that you access JSON, TOML, and YAML: you do not need to specify the root node when accessing the data.

To get the contents of `<title>` in the document below, you use `{{ .message.title }}`:

```xml
<root>
    <message>
        <title>Hugo rocks!</title>
        <description>Thanks for using Hugo</description>
    </message>
</root>
```

The following example lists the items of an RSS feed:

```go-html-template
{{ with resources.GetRemote "https://example.com/rss.xml" | transform.Unmarshal }}
    {{ range .channel.item }}
        <strong>{{ .title | plainify | htmlUnescape }}</strong><br />
        <p>{{ .description | plainify | htmlUnescape }}</p>
        {{ $link := .link | plainify | htmlUnescape }}
        <a href="{{ $link }}">{{ $link }}</a><br />
        <hr>
    {{ end }}
{{ end }}
```
