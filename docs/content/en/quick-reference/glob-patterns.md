---
title: Glob patterns
description: A quick reference guide to glob pattern syntax and matching rules for wildcards, character sets, and delimiters, featuring illustrative examples.
categories: []
keywords: []
---

{{% glossary-term "glob pattern" %}}

The table below details the supported glob pattern syntax and its matching behavior. Each example illustrates a specific match type, the pattern used, and the expected boolean result when evaluated against a test string.

| Match type | Glob pattern | Test string | Match? |
| :--- | :--- | :--- | :--- |
| Simple wildcard | `a/*.md` | `a/page.md` | true |
| Literal match | `'a/*.md'` | `a/*.md` | true |
| Single-level wildcard | `a/*/page.md` | `a/b/page.md` | true |
| Single-level wildcard | `a/*/page.md` | `a/b/c/page.md` | false |
| Multi-level wildcard | `a/**/page.md` | `a/b/c/page.md` | true |
| Single character | `file.???` | `file.txt` | true |
| Single character | `file.???` | `file.js` | false |
| Delimiter exclusion | `?at` | `f/at` | false |
| Character list | `f.[jt]xt` | `f.txt` | true |
| Negated list | `f.[!j]xt` | `f.txt` | true |
| Character range | `f.[a-c].txt` | `f.b.txt` | true |
| Character range | `f.[a-c].txt` | `f.z.txt` | false |
| Negated range | `f.[!a-c].txt` | `f.z.txt` | true |
| Pattern alternates | `*.{jpg,png}` | `logo.png` | true |
| No match | `*.{jpg,png}` | `logo.webp` | false |

The matching logic follows these rules:

- Standard wildcard (`*`) matches any character except for a delimiter.
- Super wildcard (`**`) matches any character including delimiters.
- Single character (`?`) matches exactly one character, excluding delimiters.
- Negation (`!`) matches any character except those specified in a list or range when used inside brackets.
- Character ranges (`[a-z]`) match any single character within the specified range.

The delimiter is a slash (`/`), except when matching semantic version strings, where the delimiter is a dot&nbsp;(`.`).
