---
title: Contribute to the Hugo Docs
linktitle: Docs
description: Documentation is an integral part of any open source project. The Hugo docs are as much a work in progress as the source it attempts to teach its users.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [contribute]
tags: [docs,documentation,community, contribute]
weight: 20
draft: false
aliases: [/contribute/docs/]
toc: true
---

Documentation is a critical component of any open-source project. The Hugo docs were completely reworked for the release of v0.20, but there is always room for improvement.

## Create Your Fork

It's best to make changes to the Hugo docs on your local machine to check for consistent visual styling. Make sure you've created a fork of Hugo on GitHub and cloned the repository locally on your machine. For more information, you can use the [GitHub docs for "forking"][ghforking] or see [Hugo's extensive development contribution guide][hugodev].

You can then create a separate branch for your additions. Note that you can choose a different descriptive branch name that best fits the type of content. The following is an example of a branch name you might use for adding a new website to the showcase:

```git
git checkout -b jon-doe-showcase-addition
```

## Adding New Content

The Hugo docs make heavy use of Hugo's [archetypes][] feature to easily scaffold new instances of content types. All content sections in Hugo documentation have an assigned archetype. You can [see the Hugo docs archetype source][archsource] for more clarity.

Adding new content to the Hugo docs follows the same pattern, regardless of the content section:

```
hugo new <DOCS-SECTION>/<new-content-lowercase>.md
```

{{% note "`title:`, `date:`, and Field Order" %}}
`title` and `date` fields are added automatically when using archetypes via `hugo new`. Do not be worried if the order of the new file's front matter fields on your local machine is different than that of the examples provided in the Hugo docs. This is a known issue [(#452)](https://github.com/spf13/hugo/issues/452).
{{% /note %}}

### Adding a New Function

Once you have cloned the Hugo repository, you can create a new function via the following command. Keep the file name lowercase.

```
hugo new functions/newfunction.md
```

The archetype for the `functions` content type is as follows:

{{% code file="archetypes/functions.md" %}}
```yaml
{{< readfile file="archetypes/functions.md">}}
```
{{% /code %}}

#### New Function Required Fields

Here is a review of the front matter fields automatically generated for you using `hugo new functions/*`:

***`title`***
: this will be auto-populated in all lowercase when you use `hugo new` generator.

***`linktitle`***
: the function's actual casing (e.g., `replaceRE` rather than `replacere`).

***`description`***
: a brief description used to populate the [Functions Quick Reference](/functions/).

`categories`
: currently auto-populated with 'functions` for future-proofing and portability reasons only; ignore this field.

`tags`
: only if you think it will help end users find other related functions

`signature`
: this is a signature/syntax definition for calling the function (e.g., `apply SEQUENCE FUNCTION [PARAM...]`).

`workson`
: acceptable values include `lists`,`taxonomies`, `terms`, `groups`, and `files`.

`hugoversion`
: the version of Hugo that will ship with this new function.

`relatedfuncs`
: other [templating functions][] you feel are related to your new function to help fellow Hugo users.

`{{.Content}}`
: an extended description of the new function; examples are not only welcomed but encouraged.

In the body of you function, expand the short description used in the front matter. Include as many examples as possible, and leverage the Hugo docs [code shortcodes](#adding-code-blocks). If you are unable to add examples but would like to solicit help from the Hugo community, add `needsexample: true` to your front matter.

### Adding to the Showcase

Once you have cloned the Hugo repository, you can add your Hugo website as a new showcase content file via the following command. Name the markdown file accordingly:

```
hugo new tutorials/my-hugo-showcase-website.md
```

The archetype for the `showcase` content type is as follows:

{{% code file="archetypes/showcase.md" %}}
```yaml
{{< readfile file="archetypes/showcase.md">}}
```
{{% /code %}}

#### Showcase Required Fields

`sitelink`
: the *full* URL to your website.

`title`
: the `<title>` of your website.

`description`
: a general description of your website, preferably < 180 characters.

`image`
: the image (filename only) you want to associate with your website on the Showcase page. The image should be 450px &times; 300px.

We also appreciate the addition of the remaining fields, specially `sourcelink` and `license` if you are willing to share your hard work with the open-source community. `tags` is optional, but we recommend adding at least 2 to 3 tags to improve discoverability.

#### Add an Image for the Showcase

We need to create the thumbnail of your website. Give your thumbnail a name like `my-hugo-site-name.png`. Save it under [`docs/static/images/showcase/`][].

{{% warning "Showcase Image Size" %}}
It's important that the image you use for your showcase submission has the required dimensions of 600px &times; 400px or the site will not render appropriately. Be sure to optimize your image as a matter of best practice. [Compressor](https://compressor.io/) offers a simple drag-and-drop GUI for optimizing your images.
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

## Adding Code Blocks

Code blocks are crucial for providing examples of Hugo's new features to end users of the Hugo docs. Whenever possible, create examples that you think Hugo users will be able to implement in their own projects.

### Standard Syntax

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
With both `code` and `output` shortcodes, *you must include triple back ticks and language declaration*. This was done by design so that the shortcode wrappers were easily added to legacy documentation and will be that much easier to remove if needed in future versions of the Hugo docs. We assume that the triple-back-tick syntax will live longer than our current, pretty shortcode. {{< emo ":smile:" >}}
{{% /note %}}

#### `code`

`code` is the code block shortcode you'll use most often. `code` requires at least a single `file` named parameter. Here is the signature:

````golang
{{%/* code file="smart/file/name/with/path.html" download="download.html" copy="true" */%}}
```language
A whole bunch of coding going on up in here! Boo-yah!
```
{{%/* /code */%}}
````

These are the arguments passed into `code`

***`file`***
: the only *required* argument. `file` is needed for styling but also plays an important role in helping users create a mental model around Hugo's directory structure. Visually, this will be displayed as text in the top bar of the text-editor. Always end the value with an extension. For example, instead of `./public/section/`, use `public/section/index.html`. The file extension is used to display the icon.

`download`
: if omitted, this will have no effect on the rendered shortcode. When a value is added to `download`, it's used as the filename for a downloadable version of the code block.

`copy`
: a copy button is added automatically to all `code` (i.e., default value  = `true`). If you want to keep the filename and styling of `code` but don't want to encourage readers to copy the code (e.g., a "Do not do" snippet in a tutorial), pass the `copy` argument as `copy="false"`.

##### Example `code` Input

Here is an HTML code block in a case where we want to show the user of the Hugo docs the following:

1. This type of file *could* live in `layouts/_default`.
2. This snippet is complete enough that it might be worth downloading as a standalone file.

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

##### Example 'code' Display

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

Cool, right?

#### Output Code Block

The `output` shortcode is almost identical to the `code` shortcode but only takes and requires `file`. The purpose of `output` is to show *rendered* HTML and therefore almost always follows another basic code block *or* and instance of the `code` shortcode:

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

Blockquotes can be added to the Hugo documentation using [typical Markdown blockquote syntax][bqsyntax]:

```markdown
> Without the threat of punishment, there is no joy in flight.
```

The preceding blockquote will render as follows in the Hugo docs:

> Without the threat of punishment, there is no joy in flight.

However, you can add a quick and easy `<cite>` element (added on the client via JavaScript) by separating your main blockquote and the citation with ` - `:

```markdown
> Without the threat of punishment, there is no joy in flight. - [Kobo Abe](https://en.wikipedia.org/wiki/K%C5%8Db%C5%8D_Abe)
```

Which will render as follows on the Hugo docs:

> Without the threat of punishment, there is no joy in flight. - [Kobo Abe][abe]

{{% note "Blockquotes `!=` Admonitions" %}}
Previous versions of Hugo documentation used blockquotes to draw attention to text. This is *not* the [intended semantic use of `<blockquote>`](http://html5doctor.com/cite-and-blockquote-reloaded/). Use blockquotes when quoting. To note or warn your user of specific information, use the admonition shortcodes that follow.
{{% /note %}}

## Admonitions

**Admonitions** are common in technical documentation. The most popular is that seen in [reStructuredTex Directives][sourceforge]. From the SourceForge documentation:

> Admonitions are specially marked "topics" that can appear anywhere an ordinary body element can. They contain arbitrary body elements. Typically, an admonition is rendered as an offset block in a document, sometimes outlined or shaded, with a title matching the admonition type. - [SourceForge][sourceforge]

Both `note` and `warning` use a single, *optional* argument for the admonition title. You can use markdown syntax in this title if you would like. If the title, a [positional parameter in quotes][shortcodeparams] is missing, the default behavior of the `note` and `warning` shortcodes will be to display the text "Note" and "Warning", respectively.

### `note` Admonition

Use the `note` shortcode when you want to draw attention to information subtly. `note` is intended to be less of an interruption in content than is `warning`.

#### Example `note` Input

{{% code file="note-with-heading.md" %}}
```golang
{{%/* note "Example Note Admonition" */%}}
Here is a piece of information I would like to draw your **attention** to.
{{%/* /note */%}}
```
{{% /code %}}

#### Example `note` Output

{{% output file="note-with-heading.html" %}}
```html
{{% note "Example Note Admonition" %}}
Here is a piece of information I would like to draw your **attention** to.
{{% /note %}}
```
{{% /output %}}

#### Example `note` Display

{{% note "Example Note Admonition" %}}
Here is a piece of information I would like to draw your **attention** to.
{{% /note %}}

### `warning` Admonition

Use the `warning` shortcode when you want to draw the user's attention to something important. A good usage example is for announcing breaking changes for Hugo versions, known bugs, or templating gotchas.

#### Example `warning` Input

{{% code file="warning-admonition-input.md" %}}
```golang
{{%/* warning "Example Warning Admonition" */%}}
This is a warning, which should be reserved for *important* information like breaking changes.
{{%/* /warning */%}}
```
{{% /code %}}

#### Example `warning` Output

{{% output file="warning-admonition-output.html" %}}
```html
{{% warning "Example Warning Admonition" %}}
This is a warning, which should be reserved for *important* information like breaking changes.
{{% /warning %}}
```
{{% /output %}}

#### Example `warning` Display

{{% warning "Example Warning Admonition" %}}
This is a warning, which should be reserved for *important* information like breaking changes.
{{% /warning %}}

## Editorial Style Guide

{{% note %}}
It's more important to contribute *some* documentation than no documentation at all. We need your help!
{{% /note %}}

The Hugo docs are not especially prescriptive in terms of grammar and usage. We encourage everyone to contribute regardless of your writing style. That said, here are a few gotchas when writing your documentation that, if observed, will create a more consistent documentation experience for the Hugo community:

1. *Front matter* and *file system* are two words; *Homepage* is one word.
3. Add a `godocref` value to the front matter of content files whenever possible. We want to promote Hugo *and* Golang by demonstrating the inseparable wedding of the two.

## Ask for Code Examples

Sometimes you want to contribute to the docs but don't have enough time to provide lengthy examples. If you want to flag a piece of content as needing more examples, add the following field to your front matter:

```
needsexamples: true
```

## Places to Start

The preceding `needsexamples` field is used to generate the following list of flagged content. Links will take you directly to the edit URL for the file within the GitHub GUI in the event that you are not comfortable cloning and editing the repository locally.

{{< needsexamples >}}

{{% note "Pull Requests and Branches" %}}
Similar to [contributing to Hugo development](/contribute/development/), the Hugo team expects you to create a separate branch/fork when you make your generous contributions to the Hugo docs.
{{% /note %}}

[abe]: https://en.wikipedia.org/wiki/K%C5%8Db%C5%8D_Abe
[archetypes]: /content-management/archetypes/
[archsource]: https://github.com/spf13/hugo/tree/master/docs/archetypes
[bqsyntax]: https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet#blockquotes
[charcount]: http://www.lettercount.com/
[ghforking]: https://help.github.com/articles/fork-a-repo/
[hugodev]: /contribute/development/
[shortcodeparams]: content-management/shortcodes/#shortcodes-without-markdown
[sourceforge]: http://docutils.sourceforge.net/docs/ref/rst/directives.html#admonitions
[templating function]: /functions/