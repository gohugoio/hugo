---
title: rune literal
reference: https://go.dev/ref/spec#Rune_literals
---

A _rune literal_ is the textual representation of a [_rune_](g) within a [_template_](g). It consists of a character sequence enclosed in single quotes, such as `'x'`, `'\n'`, or `'ü'`.

  Unlike [_interpreted string literals_](g) or [_raw string literals_](g), which represent a sequence of characters, a _rune literal_ represents a single [_integer_](g) value identifying a Unicode [code point][]. Within the quotes, any character may appear except a newline or an unescaped single quote. Multi-character sequences starting with a backslash (`\`) can be used to encode specific values, such as `\n` for a newline or `\u00FC` for the letter `ü`.

  [code point]: https://en.wikipedia.org/wiki/Code_point
