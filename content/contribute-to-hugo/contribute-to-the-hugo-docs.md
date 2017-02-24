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

<!-- ## Edit Locally and Submit a Pull Request

**IN DEVELOPMENT**

## How Content is Ordered in the Hugo Docs

**IN DEVELOPMENT** -->

## Creating New Content for the Hugo Docs

**IN DEVELOPMENT**


### Adding a New Function


### Adding a New Showcase

**IN DEVELOPMENT**

### Adding a New Tutorial

**IN DEVELOPMENT**

## Code Block Shortcode Examples

**IN DEVELOPMENT**

### Input Code Block

**IN DEVELOPMENT**

### Output Code Block

**IN DEVELOPMENT**

## Blockquotes



## Admonition Short Codes

**Admonitions** are common directives in technical documentation. The most popular is that seen in [reStructuredTex Directives][sourceforge]. From the SourceForge documentation:

> Admonitions are specially marked "topics" that can appear anywhere an ordinary body element can. They contain arbitrary body elements. Typically, an admonition is rendered as an offset block in a document, sometimes outlined or shaded, with a title matching the admonition type. - [SourceForge][sourceforge]


Both `note` and `warning` with a single, *optional* argument for the admonition title. If the title, a [positional parameter][shortcodeparams]

{{% note "Admonitions are **NOT** Blockquotes" %}}
Previous versions of the Hugo documentation used [Markdown `<blockquote>` syntax](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet#blockquotes) to draw attention to content. This is not the [intended semantic use of the `<blockquote>` element](http://html5doctor.com/cite-and-blockquote-reloaded/).
{{% /note %}}

### Note Admonition Shortcode

Use the `note` shortcode when you want to draw attention to information subtly. `note` is intended to be less of an interruption in content than is `warning`.

#### `note` Admonition Shortcode Input

{{% input file="note-with-heading.md" %}}
```golang
{{%/* note "Example Note Admonition" */%}}
Here is a piece of information I would like to draw your **attention** to.
{{%/* /note */%}}
```
{{% /input %}}

#### `note` Admonition Shortcode Output (Code)

{{% output "note-with-heading.html" %}}
```html
{{% note "Example Note Admonition" %}}
Here is a piece of information I would like to draw your **attention** to.
{{% /note %}}
```
{{% /output %}}

#### `note` Admonition Shortcode Display

{{% note "Example Note Admonition" %}}
Here is a piece of information I would like to draw your **attention** to.
{{% /note %}}

### Warning Admonition Examples

Use the `warning` shortcode when you want to draw the user's attention to something important. A good usage example is for announcing breaking changes for Hugo versions, known bugs, or templating gotchas.

#### `warning` Admonition Shortcode Input

{{% input file="warning-admonition-input.md" %}}
```golang
{{%/* warning "Example Warning" */%}}
This is a warning, which should be reserved for *important* information like breaking changes.
{{%/* /warning */%}}
```
{{% /input %}}

#### `warning` Admonition Shortcode Output

{{% output "warning-admonition-output.html" %}}
```html
{{% warning "Example Warning" %}}
This is a warning, which should be reserved for *important* information like breaking changes.
{{% /warning %}}
```
{{% /output %}}

#### `warning` Admonition Shortcode Display

{{% warning "Example Warning" %}}
This is a warning, which should be reserved for *important* information like breaking changes.
{{% /warning %}}

<!-- ## Example Site Shortcodes

### Example File Shortcode

### Example Front Matter Shortcode -->

## Editorial Style Guide

{{% note %}}
It's more important to contribute *some* documentation than no documentation at all. We need your help!
{{% /note %}}

The Hugo docs are not especially prescriptive in terms of grammar and usage. We encourage everyone to contribute regardless of your writing style. That said, here are a few gotchas when writing your documentation that, if observed, will create a more consistent documentation experience:

1. *Front matter* is two words.
2. *Homepage* is one word.
3. Be sure to add a `godocref` whenever possible to a new content file's front matter. We want to promote Hugo *and* Golang by demonstrating the inseparable wedding of the two.

## How Content is Ordered in the Docs

**IN DEVELOPMENT**

## Be Mindful of Aliases

Use aliases sparingly. The following table shows a list of all the aliases used in the Hugo Docs. If you need to use an alias in your new content file's front matter, be sure to check here first to prevent conflicts.

{{< allaliases >}}

[shortcodeparams]: content-management/shortcodes/#shortcodes-without-markdown
[sourceforge]: http://docutils.sourceforge.net/docs/ref/rst/directives.html#admonitions