---
title: Contribute to the Hugo Docs
linktitle: Contribute to the Hugo Docs
description: Documentation is an integral part of any open source project. The Hugo docs are as much a work in progress as the source it attempts to teach its users.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [contribute to hugo]
tags: [docs,documentation,community]
weight: 20
draft: false
slug:
aliases: [/docs-contribute/,/docscontrib/]
toc: true
needsreview: true
---

Documentation is an integral part of any open source project. The Hugo docs were completely reworked in anticipation of the release of v0.19, but there is always room for improvement.

## Edit Locally and Submit a Pull Request

## How Content is Ordered in the Hugo Docs

## Creating New Files from Archetypes

### New Default Content

### New Function

### New Showcase

### New Tutorial

## Code Block Shortcodes

### Input Code Block

### Output Code Block

### Example Site Code Block

## Blockquotes



## Admonition Short Codes

**Admonitions** are common directives in technical documentation. The most popular is that seen in [reStructuredTex Directives][sourceforge]. From the SourceForge documentation:

> Admonitions are specially marked "topics" that can appear anywhere an ordinary body element can. They contain arbitrary body elements. Typically, an admonition is rendered as an offset block in a document, sometimes outlined or shaded, with a title matching the admonition type. - [SourceForge][sourceforge]

{{% note "Admonitions are **NOT** Blockquotes" %}}
Previous versions of the Hugo documentation used [Markdown `<blockquote>` syntax](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet#blockquotes) to draw attention to content. This is not the [intended semantic use of the `<blockquote>` element](http://html5doctor.com/cite-and-blockquote-reloaded/).
{{% /note %}}

### Note Admonition Shortcode

Use the `note` shortcode when you want to draw attention to information subtly. `note` is intended to be less of an interruption in content than is `warning`.

#### Example `note` Admonition Shortcode Input

{{% input "example-note-with-heading.md" %}}
```golang
{{%/* note "Example Note Admonition" */%}}
Here is a piece of information I would like to draw your **attention** to.
{{%/* /note */%}}
```
{{% /input %}}

#### Examle `note` Admonition Shortcode Output

{{% note "Example Note Admonition" %}}
Here is a piece of information I would like to draw your **attention** to.
{{% /note %}}

### Warning Admonition Shortcode

Use the `warning` shortcode when you want to draw the user's attention

#### Example `warning` Admonition Shortcode Input

{{% input "example-note.md" %}}
```golang
{{%/* warning "Example Warning" */%}}
This is a warning, which should be reserved for *important* information like breaking changes, bad practices, etc.
{{%/* /warning */%}}
```
{{% /input %}}

#### Example `warning` Admonition Shortcode Output

{{% warning "Example Warning" %}}
This is a warning, which should be reserved for *important* information like breaking changes, bad practices, etc.
{{% /warning %}}

## Example Site Shortcodes

### Example File Shortcode

### Example Front Matter Shortcode

## Editorial Style Guide

The Hugo docs are not especially prescriptive in terms of grammar and usage. We encourage everyone to contribute, regardless of your writing style. **It's more important to contribute *some* documentation than no documentation at all**. That said, here are a few pointers to help the project maintain more consistency:

## How Content is Ordered in the Docs


## Be Mindful of Aliases

Use aliases sparingly. The following table shows a list of all the aliases used in the Hugo Docs. If you need to use an alias in your new content file's front matter, be sure to check here first.

{{< allaliases >}}

[sourceforge]: http://docutils.sourceforge.net/docs/ref/rst/directives.html#admonitions