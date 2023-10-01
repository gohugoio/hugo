---
title: strings.Trim
linkTitle: trim
description: Returns a slice of a passed string with all leading and trailing characters from cutset removed.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [title]
  returnType: string
  signatures: [strings.Trim INPUT CUTSET]
relatedFunctions:
  - strings.Chomp
  - strings.Trim
  - strings.TrimLeft
  - strings.TrimPrefix
  - strings.TrimRight
  - strings.TrimSuffix
aliases: [/functions/trim]
---

```go-html-template
{{ trim "++Batman--" "+-" }} → "Batman"
```

`trim` *requires* the second argument, which tells the function specifically what to remove from the first argument. There is no default value for the second argument, so **the following usage will not work**:

```go-html-template
{{ trim .Inner }}
```

Instead, the following example tells `trim` to remove extra new lines from the content contained in the [shortcode `.Inner` variable][shortcodevars]:

```go-html-template
{{ trim .Inner "\n" }}
```

{{% note %}}
Go templates also provide a simple [method for trimming whitespace](/templates/introduction/#whitespace) from either side of a Go tag by including a hyphen (`-`).
{{% /note %}}


[shortcodevars]: /variables/shortcodes/
