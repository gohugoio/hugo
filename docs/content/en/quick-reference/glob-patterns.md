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
| Simple wildcard | `images/*.jpg` | `images/a.jpg` | true |
| Literal match | `images/a\*.jpg` | `images/a*.jpg` | true |
| Single-level wildcard | `images/*/a.jpg` | `images/foo/a.jpg` | true |
| Single-level wildcard | `images/*/a.jpg` | `images/foo/bar/a.jpg` | false |
| Multi-level wildcard | `images/**/a.jpg` | `images/foo/bar/a.jpg` | true |
| Multi-level wildcard | `images/**/a.jpg` | `images/a.jpg` | false |
| Single character | `image.???` | `image.jpg` | true |
| Single character | `image.???` | `image.avif` | false |
| Delimiter exclusion | `?at` | `f/at` | false |
| Character list | `images/a.[jp]pg` | `images/a.jpg` | true |
| Negated list | `images/a.[!p]pg` | `images/a.jpg` | true |
| Character range | `images/a-[a-c].jpg` | `images/a-b.jpg` | true |
| Character range | `images/a-[a-c].jpg` | `images/a-z.jpg` | false |
| Negated range | `images/a-[!a-c].jpg` | `images/a-z.jpg` | true |
| Pattern alternates | `images/*.{jpg,png}` | `images/logo.png` | true |
| No match | `images/*.{jpg,png}` | `images/logo.webp` | false |

The matching logic follows these rules:

- Standard wildcard (`*`) matches any character except for a delimiter.
- Super wildcard (`**`) matches any character including delimiters, but when placed between two delimiters it requires at least one intervening character; it does not match zero directories.
- Single character (`?`) matches exactly one character, excluding delimiters.
- Negation (`!`) matches any character except those specified in a list or range when used inside brackets.
- Character ranges (`[a-z]`) match any single character within the specified range.

The delimiter is a slash (`/`), except when matching semantic version strings, where the delimiter is a dot&nbsp;(`.`).
