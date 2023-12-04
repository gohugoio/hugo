---
title: urls.RelRef
description: Returns the relative permalink to a page at the given path.
categories: []
keywords: []
action:
  aliases: [relref]
  related:
    - functions/urls/Ref
    - methods/page/Ref
    - methods/page/RelRef
  returnType: string
  signatures:
    - urls.RelRef PAGE PATH
    - urls.RelRef PAGE OPTIONS
aliases: [/functions/relref]
---

The first argument is the context of the page from which to resolve relative paths, typically the current page.

The second argument is a path to a page, with or without a file extension, with or without an anchor. A path without a leading `/` is first resolved relative to the given context, then to the remainder of the site. Alternatively, provide an [options map](#options) instead of a path.
.
```go-html-template
{{ relref . "about" }}
{{ relref . "about#anchor" }}
{{ relref . "about.md" }}
{{ relref . "about.md#anchor" }}
{{ relref . "#anchor" }}
{{ relref . "/blog/my-post" }}
{{ relref . "/blog/my-post.md" }}
```

The permalink returned is relative to the protocol+host portion of the baseURL specified in the site configuration. For example:

Code|baseURL|Permalink
:--|:--|:--
`{{ relref . "/about" }}`|`https://example.org/`|`/about/`
`{{ relref . "/about" }}`|`https://example.org/x/`|`/x/about/`

## Options

Instead of specifying a path, you can also provide an options map:

path
: (`string`) The path to the page, relative to the content directory. Required.

lang
: (`string`) The language (site) to search for the page. Default is the current language. Optional.

outputFormat
: (`string`) The output format to search for the page. Default is the current output format. Optional.

To return the relative permalink to another language version of a page:

```go-html-template
{{ relref . (dict "path" "about.md" "lang" "fr") }}
```

To return the relative permalink to another Output Format of a page:

```go-html-template
{{ relref . (dict "path" "about.md" "outputFormat" "rss") }}
```

By default, Hugo will throw an error and fail the build if it cannot resolve the path. You can change this to a warning in your site configuration, and specify a URL to return when the path cannot be resolved.

{{< code-toggle file=hugo >}}
refLinksErrorLevel = 'warning'
refLinksNotFoundURL = '/some/other/url'
{{< /code-toggle >}}
