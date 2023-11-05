---
title: Contribute to documentation
linkTitle: Documentation
description: Documentation is an integral part of any open source project. The Hugo documentation is as much a work in progress as the source it attempts to cover.
categories: [contribute]
keywords: [docs,documentation,community, contribute]
menu:
  docs:
    parent: contribute
    weight: 30
toc: true
weight: 30
aliases: [/contribute/docs/]
---

## GitHub workflow

Step 1
: Fork the [documentation repository].

Step 2
: Clone your fork.

Step 3
: Create a new branch with a descriptive name.

```sh
git checkout -b fix/typos-site-variables
```

Step 4
: Make changes.

Step 5
: Commit your changes with a descriptive commit message, typically 50 characters or less. Included the "Closes" keyword if your change addresses one or more open [issues].

```sh
git commit -m "Fix typos on site variables page

Closes #1234
Closes #5678"
```

Step 5
: Push the new branch to your fork of the documentation repository.

Step 6
: Visit the [documentation repository] and create a pull request (PR).

[documentation repository]: https://github.com/gohugoio/hugoDocs/
[issues]: https://github.com/gohugoio/hugoDocs/issues

Step 7
: A project maintainer will review your PR, and may request changes. You may delete your branch after the maintainer merges your PR.

## Including sample code

{{% note %}}
Use this syntax to include shortcodes calls within your code samples:

`{{</*/* foo */*/>}}`\
`{{%/*/* foo */*/%}}`
{{% /note %}}

### Fenced code blocks

Include the language when using a fenced code block.

````text
```go-html-template
{{ if eq $foo "bar" }}
  {{ print "foo is bar" }}
{{ end }}
```
````

Rendered:

```go-html-template
{{ if eq $foo "bar" }}
  {{ print "foo is bar" }}
{{ end }}
```

### The code shortcode

Use the `code` shortcode to include the file name and a copy-to-clipboard button. This shortcode accepts these optional parameters:

copy
: (`bool`) If `true`, displays a copy-to-clipboard button. Default is `true`.

file
: (`string`) The file name to display. If you do not provide a `lang` parameter, the file extension determines the code language.

lang
: (`string`) The code language. Default is `text`.

````text
{{</* code file="layouts/_default_/single.html" */>}}
{{ if eq $foo "bar" }}
  {{ print "foo is bar" }}
{{ end }}
{{</* /code */>}}

````

Rendered:

{{< code file="layouts/_default_/single.html" >}}
{{ if eq $foo "bar" }}
  {{ print "foo is bar" }}
{{ end }}
{{< /code >}}

### The code-toggle shortcode

Use the `code-toggle` shortcode to display examples of site configuration, front matter, or data files. This shortcode accepts these optional parameters:

copy
: (`bool`) If `true`, displays a copy-to-clipboard button. Default is `true`.

file
: (`string`) The file name to display. Omit the file extension for site configuration and data file examples.

fm
: (`bool`) If `true`, displays the code as front matter. Default is `false`.

#### Site configuration example

```text
{{</* code-toggle file=hugo */>}}
baseURL = 'https://example.org'
languageCode = 'en-US'
title = "Example Site"
{{</* /code-toggle */>}}
```

Rendered:

{{< code-toggle file=hugo >}}
baseURL = 'https://example.org'
languageCode = 'en-US'
title = "Example Site"
{{< /code-toggle >}}

#### Front matter example

```text
{{</* code-toggle file="content/about.md" fm=true */>}}
title = "About"
date = 2023-04-02T12:47:24-07:00
draft = false
{{</* /code-toggle */>}}
```

Rendered:

{{< code-toggle file="content/about.md" fm=true >}}
title = "About"
date = 2023-04-02T12:47:24-07:00
draft = false
{{< /code-toggle >}}

## Admonitions

Use the `note` shortcode to draw attention to content. Use the `{{%/*  */%}}` notation when calling this shortcode.

```text
{{%/* note */%}}
This is **bold** text.
{{%/* /note */%}}
```

{{% note %}}
This is **bold** text.
{{% /note %}}
