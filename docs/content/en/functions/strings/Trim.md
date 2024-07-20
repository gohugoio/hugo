---
title: strings.Trim
description: Returns the given string, removing leading and trailing characters specified in the cutset.
categories: []
keywords: []
action:
  aliases: [trim]
  related:
    - functions/strings/Chomp
    - functions/strings/TrimLeft
    - functions/strings/TrimPrefix
    - functions/strings/TrimRight
    - functions/strings/TrimSuffix
  returnType: string
  signatures: [strings.Trim INPUT CUTSET]
aliases: [/functions/trim]
---

```go-html-template
{{ trim "++foo--" "+-" }} → foo
```

To remove leading and trailing newline characters and carriage returns:

```go-html-template
{{ trim "\nfoo\n" "\n\r" }} → foo
{{ trim "\n\nfoo\n\n" "\n\r" }} → foo

{{ trim "\r\nfoo\r\n" "\n\r" }} → foo
{{ trim "\r\n\r\nfoo\r\n\r\n" "\n\r" }} → foo
```

The `strings.Trim` function is commonly used in shortcodes to remove leading and trailing newlines characters and carriage returns from the content within the opening and closing shortcode tags.

For example, with this Markdown:

```text
{{</* my-shortcode */>}}
Able was I ere I saw Elba.
{{</* /my-shortcode */>}}
```

The value of `.Inner` in the shortcode template is:

```text
\nAble was I ere I saw Elba.\n
```

If authored on a Windows system the value of `.Inner` might, depending on the editor configuration, be:

```text
\r\nAble was I ere I saw Elba.\r\n
```

This construct is common in shortcode templates:

```go-html-template
{{ trim .Inner "\n\r" }}
```
