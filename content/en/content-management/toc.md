---
title: Table of Contents
linktitle:
description: Hugo can automatically parse Markdown content and create a Table of Contents you can use in your templates.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [content management]
keywords: [table of contents, toc]
menu:
  docs:
    parent: "content-management"
    weight: 130
weight: 130	#rem
draft: false
aliases: [/extras/toc/]
toc: true
---

{{% note "TOC Heading Levels are Fixed" %}}

Previously, there was no out-of-the-box way to specify which heading levels you want the TOC to render. [See the related GitHub discussion (#1778)](https://github.com/gohugoio/hugo/issues/1778). As such, the resulting `<nav id="TableOfContents"><ul></ul></nav>` was going to start at `<h1>` when pulling from `{{.Content}}`.

Hugo [v0.60.0](https://github.com/gohugoio/hugo/releases/tag/v0.60.0) made a switch to [Goldmark](https://github.com/yuin/goldmark/) as the default library for Markdown which has improved and configurable implementation of TOC. Take a look at [how to configure TOC](/getting-started/configuration-markup/#table-of-contents) for Goldmark renderer.

{{% /note %}}

## Usage

Create your markdown the way you normally would with the appropriate headings. Here is some example content:

```
<!-- Your front matter up here -->

## Introduction

One morning, when Gregor Samsa woke from troubled dreams, he found himself transformed in his bed into a horrible vermin.

## My Heading

He lay on his armour-like back, and if he lifted his head a little he could see his brown belly, slightly domed and divided by arches into stiff sections. The bedding was hardly able to cover it and seemed ready to slide off any moment.

### My Subheading

A collection of textile samples lay spread out on the table - Samsa was a travelling salesman - and above it there hung a picture that he had recently cut out of an illustrated magazine and housed in a nice, gilded frame. It showed a lady fitted out with a fur hat and fur boa who sat upright, raising a heavy fur muff that covered the whole of her lower arm towards the viewer. Gregor then turned to look out the window at the dull weather. Drops
```

Hugo will take this Markdown and create a table of contents from `## Introduction`, `## My Heading`, and `### My Subheading` and then store it in the [page variable][pagevars]`.TableOfContents`.

The built-in `.TableOfContents` variables outputs a `<nav id="TableOfContents">` element with a child `<ul>`, whose child `<li>` elements begin with appropriate HTML headings. See [the available settings](/getting-started/configuration-markup/#table-of-contents) to configure what heading levels you want to include in TOC.

{{% note "Table of contents not available for MMark" %}}
Hugo documents created in the [MMark](/content-management/formats/#mmark) Markdown dialect do not currently display TOCs. TOCs are, however, compatible with all other supported Markdown formats.
{{% /note %}}

## Template Example: Basic TOC

The following is an example of a very basic [single page template][]:

{{< code file="layout/_default/single.html" download="single.html" >}}
{{ define "main" }}
<main>
    <article>
    <header>
        <h1>{{ .Title }}</h1>
    </header>
        {{ .Content }}
    </article>
    <aside>
        {{ .TableOfContents }}
    </aside>
</main>
{{ end }}
{{< /code >}}

## Template Example: TOC Partial

The following is a [partial template][partials] that adds slightly more logic for page-level control over your table of contents. It assumes you are using a `toc` field in your content's [front matter][] that, unless specifically set to `false`, will add a TOC to any page with a `.WordCount` (see [Page Variables][pagevars]) greater than 400. This example also demonstrates how to use [conditionals][] in your templating:

{{< code file="layouts/partials/toc.html" download="toc.html" >}}
{{ if and (gt .WordCount 400 ) (.Params.toc) }}
<aside>
    <header>
    <h2>{{.Title}}</h2>
    </header>
    {{.TableOfContents}}
</aside>
{{ end }}
{{< /code >}}

{{% note %}}
With the preceding example, even pages with > 400 words *and* `toc` not set to `false` will not render a table of contents if there are no headings in the page for the `{{.TableOfContents}}` variable to pull from.
{{% /note %}}

[conditionals]: /templates/introduction/#conditionals
[front matter]: /content-management/front-matter/
[pagevars]: /variables/page/
[partials]: /templates/partials/
[single page template]: /templates/single-page-templates/
