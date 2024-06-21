---
title: end
description: Terminates if, with, range, block, and define statements.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/go-template/block
    - functions/go-template/define
    - functions/go-template/if
    - functions/go-template/range
    - functions/go-template/with
  returnType:
  signatures: [end]
---

Use with the [`if`] statement:

```go-html-template
{{ $var := "foo" }}
{{ if $var }}
  {{ $var }} → foo
{{ end }}
```

Use with the [`with`] statement:

```go-html-template
{{ $var := "foo" }}
{{ with $var }}
  {{ . }} → foo
{{ end }}
```

Use with the [`range`] statement:

```go-html-template
{{ $var := slice 1 2 3 }}
{{ range $var }}
  {{ . }} → 1 2 3 
{{ end }}
```

Use with the [`block`] statement:

```go-html-template
{{ block "main" . }}{{ end }}
```

Use with the [`define`] statement:

```go-html-template
{{ define "main" }}
  {{ print "this is the main section" }}
{{ end }}
```

{{% include "functions/go-template/_common/text-template.md" %}}

[`block`]: /functions/go-template/block/
[`define`]: /functions/go-template/define/
[`if`]: /functions/go-template/if/
[`range`]: /functions/go-template/range/
[`with`]: /functions/go-template/with/
