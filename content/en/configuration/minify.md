---
title: Configure minify
linkTitle: Minify
description: Configure minify.
categories: []
keywords: []
---

This is the default configuration:

{{< code-toggle config=minify />}}

See the [tdewolff/minify] project page for details, but note the following:

- `css.inline` is for internal use. Changing this setting has no effect.
- `css.keepCSS2` has been deprecated. Use `css.version` instead.
- `html.keepConditionalComments` has been deprecated. Use `html.keepSpecialComments` instead.
- `svg.inline` is for internal use. Changing this setting has no effect.

[tdewolff/minify]: https://github.com/tdewolff/minify
