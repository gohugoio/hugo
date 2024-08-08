---
title: strings.Diff
description: Returns an anchored diff of the two texts OLD and NEW in the unified diff format. If OLD and NEW are identical, returns an empty string.
categories: []
keywords: []
action:
  related: []
  returnType: string
  signatures: [strings.Diff OLDNAME OLD NEWNAME NEW]
---

{{< new-in 0.125.0 >}}

Use `strings.Diff` to compare two strings and render a highlighted diff:

```go-html-template
{{ $want := `
<p>The product of 6 and 7 is 42.</p>
<p>The product of 7 and 6 is 42.</p>
`}}

{{ $got := `
<p>The product of 6 and 7 is 42.</p>
<p>The product of 7 and 6 is 13.</p>
`}}

{{ $diff := strings.Diff "want" $want "got" $got }}
{{ transform.Highlight $diff "diff" }}
```

Rendered:

![screen capture](diff-screen-capture.png)
