---
title: llms.txt template
linkTitle: llms.txt templates
description: Hugo provides embedded templates for generating llms.txt and llms-full.txt files for AI agents and language models.
categories: []
keywords: []
weight: 200
---

Hugo provides embedded templates for generating [llms.txt] files. These files list the pages on your site in a markdown format that AI agents and language models can consume directly, without needing to parse HTML.

The `llms.txt` file provides a curated listing of pages with descriptions, organized by section. The `llms-full.txt` file contains the full plain text content of every page on your site.

## Configuration

To enable the llms.txt output formats, add them to your project configuration:

{{< code-toggle file=hugo >}}
outputFormats:
  LLMS:
    baseName: llms
    mediaType: text/plain
    isPlainText: true
    root: true
    rel: alternate
  LLMSFull:
    baseName: llms-full
    mediaType: text/plain
    isPlainText: true
    root: true
    rel: alternate
outputs:
  home:
    - HTML
    - RSS
    - LLMS
    - LLMSFull
{{< /code-toggle >}}

This configuration instructs Hugo to generate both `llms.txt` and `llms-full.txt` at the root of your site when building the home page.

## llms.txt

The [embedded template] generates a markdown file listing all regular pages on your site, grouped by section, with each page's title, permalink, and description.

```text
# Site Title

Site description from your site configuration.

## Posts
- [First Post](https://example.com/posts/first/): A description of the first post.
- [Second Post](https://example.com/posts/second/): A description of the second post.

## Projects
- [Project Alpha](https://example.com/projects/alpha/): What Project Alpha is about.
```

Pages with `draft: true` in front matter are excluded. To set a site-wide description, use the `description` key in your site configuration.

## llms-full.txt

The [embedded template] generates a plain text file containing the full content of every regular page, concatenated and separated with horizontal rules:

```text
# Site Title

---

# First Post

Full plain text content of the first post...

---

# Second Post

Full plain text content of the second post...
```

## Template lookup order

You may overwrite the internal templates with custom templates. Hugo selects the template using this lookup order:

For `llms.txt`:
1. `/layouts/llms.txt`
1. `/themes/<THEME>/layouts/llms.txt`

For `llms-full.txt`:
1. `/layouts/llms-full.txt`
1. `/themes/<THEME>/layouts/llms-full.txt`

## Custom template example

This template generates a simple llms.txt with only pages that have the `llms` tag:

```text {file="layouts/llms.txt"}
# {{ .Site.Title }}
{{ range where .Site.RegularPages "Params.tags" "intersect" (slice "llms") }}
- [{{ .Title }}]({{ .Permalink }}){{ with .Description }}: {{ . }}{{ end }}
{{ end }}
```

[embedded template]: <{{% eturl llms %}}>
[llms.txt]: https://llmstxt.org/
