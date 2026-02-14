---
title: rune
---

A _rune_ is a way to represent a single character as a number. In Hugo and Go, text is stored as a sequence of bytes. However, while a basic letter like `x` uses only one byte, a single character such as the German `ü` is made up of multiple bytes. A _rune_ represents the entire character as one single value, no matter how many bytes it takes to store it.

  Technically, a _rune_ is just another name for a 32-bit [_integer_](g). It stores the Unicode [code point][], which is the official number assigned to that specific character.

  When you want to manipulate text character-by-character rather than by raw data size, you are working with _runes_. You write a _rune_ in a [_template_](g) using a [_rune literal_](g), such as `'x'`, `'\n'`, or `'ü'`.

  [code point]: https://en.wikipedia.org/wiki/Code_point
