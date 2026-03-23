---
title: strings.ReplacePairs
description: Returns a copy of a string with multiple replacements performed in a single pass, using a slice of old and new string pairs.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: string
    signatures: ['strings.ReplacePairs OLD NEW [OLD NEW ...] INPUT']
---

{{< new-in 0.158.0 />}}

Use the `strings.ReplacePairs` function to perform multiple replacements on a string in a single operation. This approach is faster than sequentially calling the [`strings.Replace`][] function.

Replacing strings sequentially requires multiple function calls and variable re-assignments.

```go-html-template
{{ $s := "aabbcc" }}
{{ $s = strings.Replace $s "a" "x" }}
{{ $s = strings.Replace $s "b" "y" }}
{{ $s = strings.Replace $s "c" "z" }}
{{ $s }} → xxyyzz
```

Using `strings.ReplacePairs` produces the same result with fewer function calls in less time.

```go-html-template
{{ "aabbcc" | strings.ReplacePairs "a" "x" "b" "y" "c" "z" }} → xxyyzz
```

Pairs may also be passed as a single slice:

```go-html-template
{{ $pairs := slice
  "a" "x"
  "b" "y"
  "c" "z"
}}
{{ "aabbcc" | strings.ReplacePairs $pairs }} → xxyyzz
```

## Examples

Observe that replacements are not applied recursively because the function scans the string only once.

```go-html-template
{{ $pairs := slice
  "a" "b"
  "b" "c"
}}
{{ "a" | strings.ReplacePairs $pairs }} → b
```

Apply the first match when multiple old strings could match at the same position.

```go-html-template
{{ $pairs := slice
  "app" "pear"
  "apple" "orange"
}}
{{ "apple" | strings.ReplacePairs $pairs }} → pearle
```

Delete specific strings by providing an empty string as the second value in a pair.

```go-html-template
{{ $pairs := slice "b" "" }}
{{ "abc" | strings.ReplacePairs $pairs }} → ac
```

## Edge cases

The table below outlines how the function handles various input scenarios.

Scenario|Result
:--|:--
Fewer than two arguments|Error
Odd number of slice elements|Error
Empty slice|Returns the input string
Empty input string|Returns an empty string
Empty old string|Returns the input string [interleaved](g) with the new string

## Performance

While `strings.Replace` and `strings.ReplacePairs` can produce the same results, they handle data differently. Choosing the right one can noticeably reduce the time Hugo takes to build your project.

### Single pass vs. multiple passes

When using `strings.Replace`, Hugo must scan the text from start to finish to find a match. If you chain three replacements together, Hugo performs three separate passes over the entire string.

The `strings.ReplacePairs` function is more efficient because it performs a single pass. Hugo looks through the text once and applies all replacements simultaneously.

### Caching

Unlike `strings.Replace`, which performs a direct substitution, `strings.ReplacePairs` requires an initialization step to prepare the single-pass replacement logic. To make this efficient, Hugo manages this logic using a cache:

- During the initial call, Hugo initializes and stores the logic for that specific set of pairs.
- During subsequent calls, Hugo retrieves the stored logic, skipping the initialization step and reducing the duration of the call.

### Choosing the right function

The efficiency of `strings.ReplacePairs` increases as the text gets longer or the number of pairs grows. Consider these scenarios when deciding which function to use:

- For a single replacement on a short string like a title, `strings.Replace` is efficient.
- For multiple replacements or long strings like a long-form article, `strings.ReplacePairs` is much faster.

For a document with about 8000 characters, which is roughly the length of a long-form article, `strings.ReplacePairs` outperforms five sequential `strings.Replace` calls during the initial call. Once cached, it is the faster choice for almost any situation with two or more pairs.

[`strings.Replace`]: /functions/strings/replace/
