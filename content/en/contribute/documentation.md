---
title: Documentation
description: Help us to improve the documentation by identifying issues and suggesting changes.
categories: []
keywords: []
aliases: [/contribute/docs/]
---

## Introduction

We welcome corrections and improvements to the documentation. The documentation lives in a separate repository from the main project. To contribute:

- For corrections and improvements to existing documentation, submit issues and pull requests to the [documentation repository].
- For documentation of new features, include the documentation changes in your pull request to the [project repository].

## Guidelines

### Style

Follow Google's [developer documentation style guide].

[developer documentation style guide]: https://developers.google.com/style

### Markdown

Adhere to these Markdown conventions:

- Use [ATX] headings (levels 2-4), not [setext] headings.
- Use [fenced code blocks], not [indented code blocks].
- Use hyphens, not asterisks, for unordered [list items].
- Use the [note shortcode] instead of blockquotes or bold text for emphasis.
- Do not mix [raw HTML] within Markdown.
- Do not use bold text in place of a heading or description term (`dt`).
- Remove consecutive blank lines (limit to two).
- Remove trailing spaces.

### Glossary

[Glossary] terms are defined on individual pages, providing a central repository for definitions, though these pages are not directly linked from the site.

Definitions must be complete sentences, with the first sentence defining the term. Italicize the first occurrence of the term and any referenced glossary terms for consistency.

Link to glossary terms using this syntax: `[term](g)`

Term lookups are case-insensitive, ignore formatting, and support singular and plural forms. For example, all of these variations will link to the same glossary term:

```text
[global resource](g)
[Global Resource](g)
[Global Resources](g)
[`Global Resources`](g)
```

Use the glossary-term shortcode to insert a term definition:

```text
{{%/* glossary-term "global resource" */%}}
```

### Terminology

Link to the [glossary] as needed and use terms consistently. Pay particular attention to:

- "front matter" (two words, except when referring to the configuration key)
- "home page" (two words)
- "website" (one word)
- "standalone" (one word, no hyphen)
- "map" (instead of "dictionary")
- "flag" (instead of "option" for command-line flags)
- "client side" (noun), "client-side" (adjective)
- "Markdown" (capitalized)
- "open-source" (hyphenated adjective)

### Titles and headings

- Use sentence-style capitalization.
- Avoid formatted strings.
- Keep them concise.

### Writing style

Use active voice and present tense wherever possible.

No → With Hugo you can build a static site.\
Yes → Build a static site with Hugo.

No → This will cause Hugo to generate HTML files in the `public` directory.\
Yes → Hugo generates HTML files in the `public` directory.

Use second person instead of third person.

No → Users should exercise caution when deleting files.\
Better → You must be cautious when deleting files.\
Best → Be cautious when deleting files.

Minimize adverbs.

No → Hugo is extremely fast.\
Yes → Hugo is fast.

{{< note >}}
"It's an adverb, Sam. It's a lazy tool of a weak mind." (Outbreak, 1995).
{{< /note >}}

### Function and method descriptions

Start descriptions in the functions and methods sections with "Returns" or "Reports whether" (for boolean values).

[functions]: /functions
[methods]: /methods

### File paths and names

Enclose directory names, file names, and file paths in backticks, except when used in:

- Page titles
- Section headings (h1-h6)
- Definition list terms
- The `description` field in front matter

### Miscellaneous

Other best practices:

- Introduce lists with a sentence or phrase, not directly under a heading.
- Avoid bold text; use the note shortcode for emphasis.
- Do not put description terms (`dt`) in backticks unless syntactically necessary.
- Do not use Hugo's `ref` or `relref` shortcodes.
- Prioritize current best practices over multiple options or historical information.
- Use short, focused code examples.
- Use [basic english] where possible for a global audience.

[basic english]: https://simple.wikipedia.org/wiki/Basic_English

## Related content

When available, the "See also" sidebar on this site displays related pages using Hugo's [related content] feature, based on front matter keywords. We ensure keyword accuracy by validating them against `data/keywords.yaml` during the build process. If a keyword is not found, you'll be alerted and must either modify the keyword or update the data file. This validation process helps to refine the related content for better results.

[related content]: /content-management/related-content/

## Code examples

Indent code by two spaces. With examples of template code, add spaces around the action delimiters:

```go-html-template
{{ if eq $foo $bar }}
  {{ fmt.Printf "%s is %s" $foo $bar }}
{{ end }}
```

### Fenced code blocks

Always specify the language:

````text
```go-html-template
{{ if eq $foo "bar" }}
  {{ print "foo is bar" }}
{{ end }}
```
````

To include a filename header and copy-to-clipboard button:

````text
```go-html-template {file="layouts/partials/foo.html" copy=true}
{{ if eq $foo "bar" }}
  {{ print "foo is bar" }}
{{ end }}
```
````

### Shortcode calls

Use this syntax :

````text
```text
{{</*/* foo */*/>}}
{{%/*/* foo */*/%}}
```
````

### Site configuration

Use the [code-toggle shortcode] to include site configuration examples:

```text
{{</* code-toggle file=hugo */>}}
baseURL = 'https://example.org/'
languageCode = 'en-US'
title = 'My Site'
{{</* /code-toggle */>}}
```

### Front matter

Use the [code-toggle shortcode] to include front matter examples:

```text
{{</* code-toggle file=content/posts/my-first-post.md fm=true */>}}
title = 'My first post'
date = 2023-11-09T12:56:07-08:00
draft = false
{{</* /code-toggle */>}}
```

## Shortcodes

These shortcodes are commonly used throughout the documentation. Other shortcodes are available for specialized use.

### code-toggle

Use the `code-toggle` shortcode to display examples of site configuration, front matter, or data files. This shortcode takes these arguments:

config
: (`string`) The section of `site.Data.docs.config` to render.

copy
: (`bool`) Whether to display a copy-to-clipboard button. Default is `false`.

datakey:
: (`string`) The section of `site.Data.docs` to render.

file
: (`string`) The file name to display above the rendered code. Omit the file extension for site configuration examples.

fm
: (`bool`) Whether to render the code as front matter. Default is `false`.

skipHeader
: (`bool`) Whether to omit top-level key(s) when rendering a section of `site.Data.docs.config`.

```text
{{</* code-toggle file=hugo copy=true */>}}
baseURL = 'https://example.org/'
languageCode = 'en-US'
title = 'My Site'
{{</* /code-toggle */>}}
```

### deprecated-in

Use the `deprecated-in` shortcode to indicate that a feature is deprecated:

```text
{{</* deprecated-in 0.144.0 */>}}

Use [`hugo.IsServer`] instead.

[`hugo.IsServer`]: /functions/hugo/isserver/
{{</* /deprecated-in */>}}
```

### eturl

Use the embedded template URL (`eturl`) shortcode to insert an absolute URL to the source code for an embedded template. The shortcode takes a single argument, the base file name of the template (omit the file extension).

```text
This is a link to the [embedded alias template].

[embedded alias template]: {{%/* eturl alias */%}}
```

### glossary-term

Use the `glossary-term` shortcode to insert the definition of the given glossary term.

```text
{{%/* glossary-term scalar */%}}
```

### include

Use the `include` shortcode to include content from another page.

```text
{{%/* include "_common/glob-patterns.md" */%}}
```

### new-in

Use the `new-in` shortcode to indicate a new feature:

```text
{{</* new-in 0.144.0 /*/>}}
```

You can also include details:

```text
{{</* new-in 0.144.0 */>}}
This is a new feature.
{{</* /new-in */>}}
```

### note

Use the `note` shortcode to call attention to important content:

```text
{{</* note */>}}
Use the [`math.Mod`] function to control...

[`math.Mod`]: /functions/math/mod/
{{</* /note */>}}
```

## New features

Use the [`new-in`](#new-in) shortcode to indicate a new feature:

```text
{{</* new-in 0.144.0 */>}}
```

The "new in" label will be hidden if the specified version is older than a predefined threshold, based on differences in major and minor versions. See&nbsp;[details](https://github.com/gohugoio/hugoDocs/blob/master/_vendor/github.com/gohugoio/gohugoioTheme/layouts/shortcodes/new-in.html).

## Deprecated features

Use the [`deprecated-in`](#deprecated-in) shortcode to indicate that a feature is deprecated:

```text
{{</* deprecated-in 0.144.0 */>}}
Use [`hugo.IsServer`] instead.

[`hugo.IsServer`]: /functions/hugo/isserver/
{{</* /deprecated-in */>}}
```

When deprecating a function or method, add something like this to front matter:

{{< code-toggle file=content/something/foo.md fm=true >}}
expiryDate: 2027-02-17 # deprecated 2025-02-17 in v0.144.0
{{< /code-toggle >}}

Set the `expiryDate` to two years from the date of deprecation, and add a brief front matter comment to explain the setting.

## GitHub workflow

{{< note >}}
This section assumes that you have a working knowledge of Git and GitHub, and are comfortable working on the command line.
{{< /note >}}

Use this workflow to create and submit pull requests.

### Step 1

Fork the [documentation repository].

### Step 2

Clone your fork.

### Step 3

Create a new branch with a descriptive name that includes the corresponding issue number, if any:

```sh
git checkout -b restructure-foo-page-99999
```

### Step 4

Make changes.

### Step 5

Build the site locally to preview your changes.

### Step 6

Commit your changes with a descriptive commit message:

- Provide a summary on the first line, typically 50 characters or less, followed by a blank line.
- Optionally, provide a detailed description where each line is 80 characters or less, followed by a blank line.
- Optionally, add one or more "Fixes" or "Closes" keywords, each on its own line, referencing the [issues] addressed by this change.

For example:

```text
git commit -m "Restructure the taxonomy page

This restructures the taxonomy page by splitting topics into logical
sections, each with one or more examples.

Fixes #9999
Closes #9998"
```

### Step 7

Push the new branch to your fork of the documentation repository.

### Step 8

Visit the [documentation repository] and create a pull request (PR).

### Step 9

A project maintainer will review your PR and may request changes. You may delete your branch after the maintainer merges your PR.

[ATX]: https://spec.commonmark.org/0.30/#atx-headings
[Microsoft Writing Style Guide]: https://learn.microsoft.com/en-us/style-guide/welcome/
[`glossary-term`]: #glossary-term
[basic english]: https://simple.wikipedia.org/wiki/Basic_English
[code examples]: #code-examples
[code-toggle shortcode]: #code-toggle
[documentation repository]: https://github.com/gohugoio/hugoDocs/
[fenced code blocks]: https://spec.commonmark.org/0.30/#fenced-code-blocks
[glossary]: /quick-reference/glossary/
[indented code blocks]: https://spec.commonmark.org/0.30/#indented-code-blocks
[issues]: https://github.com/gohugoio/hugoDocs/issues
[list items]: https://spec.commonmark.org/0.30/#list-items
[note shortcode]: #note
[project repository]: https://github.com/gohugoio/hugo
[raw HTML]: https://spec.commonmark.org/0.30/#raw-html
[setext]: https://spec.commonmark.org/0.30/#setext-heading
