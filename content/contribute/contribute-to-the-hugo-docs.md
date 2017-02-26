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

Documentation is a critical component of any open-source project. The Hugo docs were completely reworked in anticipation of the release of v0.19, but there is always room for improvement.

## Create Your Fork

First, make sure that you created a [fork](https://help.github.com/articles/fork-a-repo/) of Hugo on Github and cloned the fork locally on your computer. Next, create a separate branch for your additions. Note that you can choose a different descriptive branch name that best fits the type of content you're trying to submit:

```git
git checkout -b showcase-addition
```

## Adding a New Content Page

The Hugo docs are built using Hugo and therefore make heavy use of Hugo's [archetype][] feature to easily scaffold new instances of content types. All [content sections][] in Hugo documentation have an assigned archetypes ([see source][archsource]).


### Adding a New Function

### Adding a New Tutorial

Once you have cloned the Hugo repository, you can create a new function via the following command. For functions, title the new function in lowercase.

```
hugo new tutorials/newfunction.md
```

The archetype for the `functions` content type is as follows:

{{% code file="archetypes/functions.md" %}}
```yaml
{{< readfile file="archetypes/functions.md">}}
```
{{% /code %}}


### Adding a New Showcase

Once you have cloned the Hugo repository, you can create add your site as a new showcase content file via the following command. Name the markdown file accordingly:

```
hugo new tutorials/my-showcase-addition.md
```

The archetype for the `showcase` content type is as follows:

{{% code file="archetypes/showcase.md" %}}
```yaml
{{< readfile file="archetypes/showcase.md">}}
```
{{% /code %}}

Add at least values for `sitelink`, `title`,  `description`, and a path for `image`.

#### Add an Image for the Showcase

We need to create the thumbnail of your website. Give your thumbnail a name like `my-hugo-site-name.png`. Save it under [`docs/static/images/showcase/`][].

{{% warning "Showcase Image Size" %}}
It's important that the image you use for your showcase submission has the required dimensions of 600px &times; 400px or the site will not render appropriately. Be sure to optimize your image as a matter of best practice. If you're looking for a quick online optimization tool, check out [Compressor](https://compressor.io/).
{{% /warning %}}

### Adding a New Tutorial

Once you have cloned the Hugo repository, you can create a new tutorial via the following command. Name the markdown file accordingly:

```
hugo new tutorials/my-new-tutorial.md
```

The archetype for the `tutorials` content type is as follows:

{{% code file="archetypes/tutorials.md" %}}
```yaml
{{< readfile file="archetypes/tutorials.md">}}
```
{{% /code %}}

## Adding Code Blocks to Hugo Docs

Code blocks are crucial for providing examples of Hugo's new features to end users of the Hugo docs. Whenever possible, create examples that you think Hugo users will be able to implement in their own projects.

### Standard Code Block Syntax

Across all pages on the Hugo docs, the typical triple-back-tick markdown syntax is used. If you do not want to take the extra time to implement the following code block shortcodes, please use standard Github-flavored markdown. The Hugo docs use a version of [highlight.js](https://highlightjs.org/) that's been modified for specific Hugo keywords.

Your options for languages are `xml`/`html`, `go`/`golang`, `md`/`markdown`/`mkd`, `handlebars`, `apache`, `toml`, `yaml`, `json`, `css`, `asciidoc`, `ruby`, `powershell`/`ps`, `scss`, `sh`/`zsh`/`bash`/`git`, `http`/`https`, and `javascript`/`js`.

````html
```html
<h1>Hello world!</h1>
```
````

### Code Block Shortcodes

The Hugo documentation comes with very robust shortcodes to help you add interactive code snippets.

{{% note %}}
With both the `code` and the `output` shortcodes, *you still need to include the triple back ticks and language declaration*. This was done by design so that the shortcode wrappers were easily added to legacy documentation and will be that much easier to remove if needed in future versions of the docs. We assume that the triple-back-tick syntax will live longer than the shortcode.
{{% /note %}}

#### Input Code Block

The first shortcode is the one you'll use most often, `code`. `code`, like all code block shortcodes, requires at least a single `file` named parameter. Here is the signature:

````golang
{{%/* code file="smart/file/name/with/path.html" download="download.html" copy="true" */%}}
```language
A whole bunch of coding going on up in here! Boo-yah!
```
{{%/* /code */%}}
````

Let's go through each of the three arguments passed into `code`:

1. `file`. This is the only ***required*** argument for this shortcode. `file` is needed for styling but also plays an important role in helping users create a mental model around Hugo's directory structure. Visually, this will be displayed as text in the top bar of the text-editor design you've seen throughout the docs. Note that you always want to end the file with an extension, but never use more than a single period in the filename. For example, instead of `./archetypes/default`, use `archetypes/default.md`. The file extension is important for displaying the correct icon.
2. `download`. If omitted entirely, this will have no effect on the rendered shortcode. When a value is added to `download`, it's used as the filename for a downloadable version of the code block.
3. `copy`. All `code` instances add a copy button to the bottom right automatically. However, there are times where you may not want to encourage your end user to copy a code block but still want to keep consistent styling for the filename (e.g., if you adding a "Do not do" code in a tutorial). If you want to turn off the copy functionality of `code`, you can add `copy="false"`.

##### Example of `Code`

Here is an HTML code block we want to tell the user lives in the `layouts/_default` directory. We are going to also make the block downloadable because it works as a standalone file:

````html
{{%/* code file="layouts/_default/single.html" download="single.html" */%}}
```html
{{ define "main" }}
<main class="main">
    <article class="content">
        <header>
            <h1>{{.Title}}</h1>
            {{with .Params.subtitle}}
            <span class="subtitle">{{.}}</span>
        </header>
        <div class="body-copy">
            {{.Content}}
        </div>
        <aside class="toc">
            {{.TableOfContents}}
        </aside>
    </article>
</main>
{{ end }}
```
{{%/* /code */%}}
````

The output of this example will render to the Hugo docs as follows:

{{% code file="layouts/_default/single.html" download="single.html" %}}
```html
{{ define "main" }}
<main class="main">
    <article class="content">
        <header>
            <h1>{{.Title}}</h1>
            {{with .Params.subtitle}}
            <span class="subtitle">{{.}}</span>
        </header>
        <div class="body-copy">
            {{.Content}}
        </div>
        <aside class="toc">
            {{.TableOfContents}}
        </aside>
    </article>
</main>
{{ end }}
```
{{% /code %}}

#### Output Code Block

The `output` shortcode is almost identical to the `code` shortcode but doesn't take any more arguments than the required `file`. The purpose is to demonstrate what the output, or *rendered*, HTML will look after Hugo builds its templates:

````html
{{%/* output file="post/my-first-post/index.html" */%}}
```html
<h1>This is my First Hugo Blog Post</h1>
<p>I am excited to be using Hugo.</p>
```
{{%/* /output */%}}
````

The preceding `output` example will render as follows to the Hugo docs:

{{% output file="post/my-first-post/index.html" %}}
```html
<h1>This is my First Hugo Blog Post</h1>
<p>I am excited to be using Hugo.</p>
```
{{% /output %}}

## Blockquotes

Blockquotes can be


{{% note "Blockquotes `!=` Admonitions" %}}
Previous versions of the Hugo documentation used [Markdown `<blockquote>` syntax](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet#blockquotes) to draw attention to content. This is [*not* the intended semantic use of the `<blockquote>` element](http://html5doctor.com/cite-and-blockquote-reloaded/). Use blockquotes when quoting actual text. To note or warn your user of specific information, use the admonition shortcodes that follow.
{{% /note %}}

## Admonition Short Codes

**Admonitions** are common directives in technical documentation. The most popular is that seen in [reStructuredTex Directives][sourceforge]. From the SourceForge documentation:

> Admonitions are specially marked "topics" that can appear anywhere an ordinary body element can. They contain arbitrary body elements. Typically, an admonition is rendered as an offset block in a document, sometimes outlined or shaded, with a title matching the admonition type. - [SourceForge][sourceforge]

Both `note` and `warning` use a single, *optional* argument for the admonition title, which accepts markdown syntax as well. If the title, a [positional parameter][shortcodeparams] in quotes is missing, the default behavior of the `note` and `warning` shortcodes will be to display the text "Note" and "Warning", respectively.

### Note Admonition Shortcode

Use the `note` shortcode when you want to draw attention to information subtly. `note` is intended to be less of an interruption in content than is `warning`.

#### `note` Admonition Shortcode Input

{{% code file="note-with-heading.md" %}}
```golang
{{%/* note "Example Note Admonition" */%}}
Here is a piece of information I would like to draw your **attention** to.
{{%/* /note */%}}
```
{{% /code %}}

#### `note` Admonition Shortcode Output (Code)

{{% output file="note-with-heading.html" %}}
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

{{% code file="warning-admonition-input.md" %}}
```golang
{{%/* warning "Example Warning Admonition" */%}}
This is a warning, which should be reserved for *important* information like breaking changes.
{{%/* /warning */%}}
```
{{% /code %}}

#### `warning` Admonition Shortcode Output

{{% output file="warning-admonition-output.html" %}}
```html
{{% warning "Example Warning Admonition" %}}
This is a warning, which should be reserved for *important* information like breaking changes.
{{% /warning %}}
```
{{% /output %}}

#### `warning` Admonition Shortcode Display

{{% warning "Example Warning Admonition" %}}
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

## Pages Needing Code Examples

Examples

{{< needsexamples >}}

## How Content is Ordered in the Docs

**IN DEVELOPMENT**

## Be Mindful of Aliases

Use aliases sparingly. The following table shows a list of all the aliases used in the Hugo Docs. If you need to use an alias in your new content file's front matter, be sure to check here first to prevent conflicts.

{{< allaliases >}}

[archsource]: https://github.com/spf13/hugo/tree/master/docs/archetypes
[archetype]: /content-management/archetypes/
[shortcodeparams]: content-management/shortcodes/#shortcodes-without-markdown
[sourceforge]: http://docutils.sourceforge.net/docs/ref/rst/directives.html#admonitions