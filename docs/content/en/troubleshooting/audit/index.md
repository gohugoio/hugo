---
title: Site audit
linkTitle: Audit
description: Run this audit before deploying your production site.
categories: [troubleshooting]
keywords: []
menu:
  docs:
    parent: troubleshooting
    weight: 20
weight: 20
---

There are several conditions that can produce errors in your published site which are not detected during the build. Run this audit before your final build.

{{< code copy=true >}}
HUGO_MINIFY_TDEWOLFF_HTML_KEEPCOMMENTS=true HUGO_ENABLEMISSINGTRANSLATIONPLACEHOLDERS=true hugo && grep -inorE "<\!-- raw HTML omitted -->|ZgotmplZ|\[i18n\]|\(<nil>\)|(&lt;nil&gt;)|hahahugo" public/
{{< /code >}}

_Tested with GNU Bash 5.1 and GNU grep 3.7._

## Example output

![site audit terminal output](screen-capture.png)

## Explanation

### Environment variables

`HUGO_MINIFY_TDEWOLFF_HTML_KEEPCOMMENTS=true`
: Retain HTML comments even if minification is enabled. This takes precedence over `minify.tdewolff.html.keepComments` in the site configuration. If you minify without keeping HTML comments when performing this audit, you will not be able to detect when raw HTML has been omitted.

`HUGO_ENABLEMISSINGTRANSLATIONPLACEHOLDERS=true`
: Show a placeholder instead of the default value or an empty string if a translation is missing. This takes precedence over `enableMissingTranslationPlaceholders` in the site configuration.

### Grep options

`-i, --ignore-case`
: Ignore case distinctions in patterns and input data, so that characters that differ only in case match each other.

`-n, --line-number`
: Prefix each line of output with the 1-based line number within its input file.

`-o, --only-matching`
: Print only the matched (non-empty) parts of a matching line, with each such part on a separate output line.

`-r, --recursive`
: Read all files under each directory, recursively, following symbolic links only if they are on the command line.

`-E, --extended-regexp`
: Interpret PATTERNS as extended regular expressions.

### Patterns

`<!-- raw HTML omitted -->`
: By default, Hugo strips raw HTML from your Markdown prior to rendering, and leaves this HTML comment in its place.

`ZgotmplZ`
: ZgotmplZ is a special value that indicates that unsafe content reached a CSS or URL context at runtime. See&nbsp;[details].

[details]: https://pkg.go.dev/html/template

`[i18n]`
: This is the placeholder produced instead of the default value or an empty string if a translation is missing.

`(<nil>)`
: This string will appear in the rendered HTML when passing a nil value to the `printf` function.

`(&lt;nil&gt;)`
: Same as above when the value returned from the `printf` function has not been passed through `safeHTML`.

`HAHAHUGO`
: Under certain conditions a rendered shortcode may include all or a portion of the string H&#xfeff;AHAHUGOSHORTCODE in either uppercase or lowercase. This is difficult to detect in all circumstances, but a case-insensitive search of the output for `HAHAHUGO` is likely to catch the majority of cases without producing false positives.
