---
title: Contribute to the Hugo Docs
linktitle: Documentation
description: Documentation is an integral part of any open source project. The Hugo docs are as much a work in progress as the source it attempts to cover.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [contribute]
keywords: [docs,documentation,community, contribute]
menu:
  docs:
    parent: "contribute"
    weight: 20
weight: 20
sections_weight: 20
draft: false
aliases: [/contribute/docs/]
toc: true
---

## Create Your Fork

It's best to make changes to the Hugo docs on your local machine to check for consistent visual styling. Make sure you've created a fork of [hugoDocs](https://github.com/gohugoio/hugoDocs) on GitHub and cloned the repository locally on your machine. For more information, you can see [GitHub's documentation on "forking"][ghforking] or follow along with [Hugo's development contribution guide][hugodev].

You can then create a separate branch for your additions. Be sure to choose a descriptive branch name that best fits the type of content. The following is an example of a branch name you might use for adding a new website to the showcase:

```
git checkout -b jon-doe-showcase-addition
```

## Add New Content

The Hugo docs make heavy use of Hugo's [archetypes][] feature. All content sections in Hugo documentation have an assigned archetype.

Adding new content to the Hugo docs follows the same pattern, regardless of the content section:

```
hugo new <DOCS-SECTION>/<new-content-lowercase>.md
```

### Add a New Function

Once you have cloned the Hugo repository, you can create a new function via the following command. Keep the file name lowercase.

```
hugo new functions/newfunction.md
```

The archetype for `functions` according to the Hugo docs is as follows:

{{< code file="archetypes/functions.md" >}}
{{< readfile file="/archetypes/functions.md">}}
{{< /code >}}

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

In the body of your function, expand the short description used in the front matter. Include as many examples as possible, and leverage the Hugo docs [`code` shortcode](#adding-code-blocks). If you are unable to add examples but would like to solicit help from the Hugo community, add `needsexample: true` to your front matter.

## Add Code Blocks

Code blocks are crucial for providing examples of Hugo's new features to end users of the Hugo docs. Whenever possible, create examples that you think Hugo users will be able to implement in their own projects.

### Standard Syntax

Across many pages on the Hugo docs, the typical triple-back-tick markdown syntax (```` ``` ````) is used. If you do not want to take the extra time to implement the following code block shortcodes, please use standard GitHub-flavored markdown. The Hugo docs use a version of [highlight.js](https://highlightjs.org/) with a specific set of languages.

Your options for languages are `xml`/`html`, `go`/`golang`, `md`/`markdown`/`mkd`, `handlebars`, `apache`, `toml`, `yaml`, `json`, `css`, `asciidoc`, `ruby`, `powershell`/`ps`, `scss`, `sh`/`zsh`/`bash`/`git`, `http`/`https`, and `javascript`/`js`.

````
```
<h1>Hello world!</h1>
```
````


### Code Block Shortcode

The Hugo documentation comes with a very robust shortcode for adding interactive code blocks.

{{% note %}}
With the `code` shortcodes, *you must include triple back ticks and a language declaration*. This was done by design so that the shortcode wrappers were easily added to legacy documentation and will be that much easier to remove if needed in future versions of the Hugo docs.
{{% /note %}}

### `code`

`code` is the Hugo docs shortcode you'll use most often. `code` requires has only one named parameter: `file`. Here is the pattern:

```
{{%/* code file="smart/file/name/with/path.html" download="download.html" copy="true" */%}}
A whole bunch of coding going on up in here!
{{%/* /code */%}}
```

The following are the arguments passed into `code`:


***`file`***
: the only *required* argument. `file` is needed for styling but also plays an important role in helping users create a mental model around Hugo's directory structure. Visually, this will be displayed as text in the top left of the code block.

`download`
: if omitted, this will have no effect on the rendered shortcode. When a value is added to `download`, it's used as the filename for a downloadable version of the code block.

`copy`
: a copy button is added automatically to all `code` shortcodes. If you want to keep the filename and styling of `code` but don't want to encourage readers to copy the code (e.g., a "Do not do" snippet in a tutorial), use `copy="false"`.

#### Example `code` Input

This example HTML code block tells Hugo users the following:

1. This file *could* live in `layouts/_default`, as demonstrated by `layouts/_default/single.html` as the value for `file`.
2. This snippet is complete enough to be downloaded and implemented in a Hugo project, as demonstrated by `download="single.html"`.

```
{{</* code file="layouts/_default/single.html" download="single.html" */>}}
{{ define "main" }}
<main>
    <article>
        <header>
            <h1>{{.Title}}</h1>
            {{with .Params.subtitle}}
            <span>{{.}}</span>
        </header>
        <div>
            {{.Content}}
        </div>
        <aside>
            {{.TableOfContents}}
        </aside>
    </article>
</main>
{{ end }}
{{</* /code */>}}
```

##### Example 'code' Display

The output of this example will render to the Hugo docs as follows:

{{< code file="layouts/_default/single.html" download="single.html" >}}
{{ define "main" }}
<main>
    <article>
        <header>
            <h1>{{.Title}}</h1>
            {{with .Params.subtitle}}
            <span>{{.}}</span>
        </header>
        <div>
            {{.Content}}
        </div>
        <aside>
            {{.TableOfContents}}
        </aside>
    </article>
</main>
{{ end }}
{{< /code >}}

<!-- #### Output Code Block

The `output` shortcode is almost identical to the `code` shortcode but only takes and requires `file`. The purpose of `output` is to show *rendered* HTML and therefore almost always follows another basic code block *or* and instance of the `code` shortcode:

```
{{%/* output file="posts/my-first-post/index.html" */%}}
<h1>This is my First Hugo Blog Post</h1>
<p>I am excited to be using Hugo.</p>
{{%/* /output */%}}
```

The preceding `output` example will render as follows to the Hugo docs:

{{< output file="posts/my-first-post/index.html" >}}
<h1>This is my First Hugo Blog Post</h1>
<p>I am excited to be using Hugo.</p>
{{< /output >}} -->

## Blockquotes

Blockquotes can be added to the Hugo documentation using [typical Markdown blockquote syntax][bqsyntax]:

```
> Without the threat of punishment, there is no joy in flight.
```

The preceding blockquote will render as follows in the Hugo docs:

> Without the threat of punishment, there is no joy in flight.

However, you can add a quick and easy `<cite>` element (added on the client via JavaScript) by separating your main blockquote and the citation with a hyphen with a single space on each side:

```
> Without the threat of punishment, there is no joy in flight. - [Kobo Abe](https://en.wikipedia.org/wiki/Kobo_Abe)
```

Which will render as follows in the Hugo docs:

> Without the threat of punishment, there is no joy in flight. - [Kobo Abe][abe]

{{% note "Blockquotes `!=` Admonitions" %}}
Previous versions of Hugo documentation used blockquotes to draw attention to text. This is *not* the [intended semantic use of `<blockquote>`](http://html5doctor.com/cite-and-blockquote-reloaded/). Use blockquotes when quoting. To note or warn your user of specific information, use the admonition shortcodes that follow.
{{% /note %}}

## Admonitions

**Admonitions** are common in technical documentation. The most popular is that seen in [reStructuredText Directives][sourceforge]. From the SourceForge documentation:

> Admonitions are specially marked "topics" that can appear anywhere an ordinary body element can. They contain arbitrary body elements. Typically, an admonition is rendered as an offset block in a document, sometimes outlined or shaded, with a title matching the admonition type. - [SourceForge][sourceforge]

The Hugo docs contain three admonitions: `note`, `tip`, and `warning`.

### `note` Admonition

Use the `note` shortcode when you want to draw attention to information subtly. `note` is intended to be less of an interruption in content than is `warning`.

#### Example `note` Input

{{< code file="note-with-heading.md" >}}
{{%/* note */%}}
Here is a piece of information I would like to draw your **attention** to.
{{%/* /note */%}}
{{< /code >}}

#### Example `note` Output

{{< output file="note-with-heading.html" >}}
{{% note %}}
Here is a piece of information I would like to draw your **attention** to.
{{% /note %}}
{{< /output >}}

#### Example `note` Display

{{% note %}}
Here is a piece of information I would like to draw your **attention** to.
{{% /note %}}

### `tip` Admonition

Use the `tip` shortcode when you want to give the reader advice. `tip`, like `note`, is intended to be less of an interruption in content than is `warning`.

#### Example `tip` Input

{{< code file="using-tip.md" >}}
{{%/* tip */%}}
Here's a bit of advice to improve your productivity with Hugo.
{{%/* /tip */%}}
{{< /code >}}

#### Example `tip` Output

{{< output file="tip-output.html" >}}
{{% tip %}}
Here's a bit of advice to improve your productivity with Hugo.
{{% /tip %}}
{{< /output >}}

#### Example `tip` Display

{{% tip %}}
Here's a bit of advice to improve your productivity with Hugo.
{{% /tip %}}

### `warning` Admonition

Use the `warning` shortcode when you want to draw the user's attention to something important. A good usage example is for articulating breaking changes in Hugo versions, known bugs, or templating "gotchas."

#### Example `warning` Input

{{< code file="warning-admonition-input.md" >}}
{{%/* warning */%}}
This is a warning, which should be reserved for *important* information like breaking changes.
{{%/* /warning */%}}
{{< /code >}}

#### Example `warning` Output

{{< output file="warning-admonition-output.html" >}}
{{% warning %}}
This is a warning, which should be reserved for *important* information like breaking changes.
{{% /warning %}}
{{< /output >}}

#### Example `warning` Display

{{% warning %}}
This is a warning, which should be reserved for *important* information like breaking changes.
{{% /warning %}}

{{% note "Pull Requests and Branches" %}}
Similar to [contributing to Hugo development](/contribute/development/), the Hugo team expects you to create a separate branch/fork when you make your contributions to the Hugo docs.
{{% /note %}}

[abe]: https://en.wikipedia.org/wiki/Kobo_Abe
[archetypes]: /content-management/archetypes/
[bqsyntax]: https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet#blockquotes
[charcount]: https://www.lettercount.com/
[`docs/static/images/showcase/`]: https://github.com/gohugoio/hugo/tree/master/docs/static/images/showcase/
[ghforking]: https://help.github.com/articles/fork-a-repo/
[hugodev]: /contribute/development/
[shortcodeparams]: content-management/shortcodes/#shortcodes-without-markdown
[sourceforge]: http://docutils.sourceforge.net/docs/ref/rst/directives.html#admonitions
[templating function]: /functions/
