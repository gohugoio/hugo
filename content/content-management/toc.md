---
title: Table of Contents
linktitle:
description: Hugo can automatically parse Markdown content and create a Table of Contents you can leverage in your templates to guide readers to sections of longer pages.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [content management]
tags: [table of contents, toc]
weight: 130
draft: false
aliases: [/extras/toc/,/content-management/toc/]
toc: false
wip: true
---

Hugo can automatically parse Markdown content and create a Table of Contents you can leverage in your templates to guide readers to sections of longer pages.

{{% note "TOC Heading Levels are Fixed" %}}
Currently, the {{.TableOfContents}} [page variable](/variables/page-variables/) is fixed in its behavior; i.e., you do not have the option to set the heading level at which the TOC renders. This is a [known issue (#1778)](https://github.com/spf13/hugo/issues/1778), and as always, [contributions are welcome](/contribute/development/).
{{% /note %}}

## Usage

Create your markdown the way you normally would with the appropriate headers. Here is some example content:

```md
<!-- Your front matter up here -->

## Introduction

One morning, when Gregor Samsa woke from troubled dreams, he found himself transformed in his bed into a horrible vermin.

## My Heading

He lay on his armour-like back, and if he lifted his head a little he could see his brown belly, slightly domed and divided by arches into stiff sections. The bedding was hardly able to cover it and seemed ready to slide off any moment.

His many legs, pitifully thin compared with the size of the rest of him, waved about helplessly as he looked. "What's happened to me? " he thought. It wasn't a dream. His room, a proper human room although a little too small, lay peacefully between its four familiar walls.

### My Subheading

A collection of textile samples lay spread out on the table - Samsa was a travelling salesman - and above it there hung a picture that he had recently cut out of an illustrated magazine and housed in a nice, gilded frame. It showed a lady fitted out with a fur hat and fur boa who sat upright, raising a heavy fur muff that covered the whole of her lower arm towards the viewer. Gregor then turned to look out the window at the dull weather. Drops
```

Hugo will take this Markdown and create a table of contents from `## Introuduction`, `## My Heading`, and `### My Subheading`stored in the [content variable](/variables/page-variables/) `.TableOfContents`.

## Template Example

This is example code of a [single.html template](/templates/single-page-templates/).

```golang
{{ partial "header.html" . }}
    <aside id="toc" class="well col-md-4 col-sm-6">
    {{ .TableOfContents }}
    </aside>
    <h1>{{ .Title }}</h1>
    {{ .Content }}
{{ partial "footer.html" . }}
```
