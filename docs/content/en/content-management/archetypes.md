---
title: Archetypes
description: An archetype is a template for new content.
categories: [content management]
keywords: [archetypes,generators,metadata,front matter]
menu:
  docs:
    parent: content-management
    weight: 140
  quicklinks:
weight: 140
toc: true
aliases: [/content/archetypes/]
---

## Overview

A content file consists of [front matter] and markup. The markup is typically Markdown, but Hugo also supports other [content formats]. Front matter can be TOML, YAML, or JSON.

The `hugo new content` command creates a new file in the `content` directory, using an archetype as a template. This is the default archetype:

{{< code-toggle file=archetypes/default.md fm=true >}}
title = '{{ replace .File.ContentBaseName `-` ` ` | title }}'
date = '{{ .Date }}'
draft = true
{{< /code-toggle >}}

When you create new content, Hugo evaluates the [template actions] within the archetype. For example:

```sh
hugo new content posts/my-first-post.md
```

With the default archetype shown above, Hugo creates this content file:

{{< code-toggle file=content/posts/my-first-post.md fm=true >}}
title = 'My First Post'
date = '2023-08-24T11:49:46-07:00'
draft = true
{{< /code-toggle >}}

You can create an archetype for one or more [content types]. For example, use one archetype for posts, and use the default archetype for everything else:

```text
archetypes/
├── default.md
└── posts.md
```

## Lookup order

Hugo looks for archetypes in the `archetypes` directory in the root of your project, falling back to the `archetypes` directory in themes or installed modules. An archetype for a specific content type takes precedence over the default archetype.

For example, with this command:

```sh
hugo new content posts/my-first-post.md
```

The archetype lookup order is:

1. archetypes/posts.md
1. archetypes/default.md
1. themes/my-theme/archetypes/posts.md
1. themes/my-theme/archetypes/default.md

If none of these exists, Hugo uses a built-in default archetype.

## Functions and context

You can use any [template function] within an archetype. As shown above, the default archetype uses the [`replace`](/functions/strings/replace) function to replace hyphens with spaces when populating the title in front matter.

Archetypes receive the following [context]:

Date
: (`string`) The current date and time, formatted in compliance with RFC3339.

File
: (`hugolib.fileInfo`) Returns file information for the current page. See [details](/methods/page/file).

Type
: (`string`) The [content type] inferred from the top-level directory name, or as specified by the `--kind` flag passed to the `hugo new content` command.

[content type]: /getting-started/glossary#content-type

Site
: (`page.Site`) The current site object. See [details](/methods/site/).

## Alternate date format

To insert date and time with an alternate format, use the [`time.Now`] function:

[`time.Now`]: /functions/time/now/

{{< code-toggle file=archetypes/default.md fm=true >}}
title = '{{ replace .File.ContentBaseName `-` ` ` | title }}'
date = '{{ time.Now.Format "2006-01-02" }}'
draft = true
{{< /code-toggle >}}

## Include content

Although typically used as a front matter template, you can also use an archetype to populate content.

For example, in a documentation site you might have a section (content type) for functions. Every page within this section should follow the same format: a brief description, the function signature, examples, and notes. We can pre-populate the page to remind content authors of the standard format.

{{< code file=archetypes/functions.md >}}
---
date: '{{ .Date }}'
draft: true
title: '{{ replace .File.ContentBaseName `-` ` ` | title }}'
---

A brief description of what the function does, using simple present tense in the third person singular form. For example:

`someFunction` returns the string `s` repeated `n` times.

## Signature

```text
func someFunction(s string, n int) string
```

## Examples

One or more practical examples, each within a fenced code block.

## Notes

Additional information to clarify as needed.
{{< /code >}}

Although you can include [template actions] within the content body, remember that Hugo evaluates these once---at the time of content creation. In most cases, place template actions in a [template] where Hugo evaluates the actions every time you [build](/getting-started/glossary/#build) the site.

## Leaf bundles

You can also create archetypes for [leaf bundles](/getting-started/glossary/#leaf-bundle).

For example, in a photography site you might have a section (content type) for galleries. Each gallery is leaf bundle with content and images.

Create an archetype for galleries:

```text
archetypes/
├── galleries/
│   ├── images/
│   │   └── .gitkeep
│   └── index.md      <-- same format as default.md
└── default.md
```

Subdirectories within an archetype must contain at least one file. Without a file, Hugo will not create the subdirectory when you create new content. The name and size of the file are irrelevant. The example above includes a&nbsp;`.gitkeep` file, an empty file commonly used to preserve otherwise empty directories in a Git repository.

To create a new gallery:

```sh
hugo new galleries/bryce-canyon
```

This produces:

```text
content/
├── galleries/
│   └── bryce-canyon/
│       ├── images/
│       │   └── .gitkeep
│       └── index.md
└── _index.md
```

## Use alternate archetype

Use the `--kind` command line flag to specify an alternate archetype when creating content.

For example, let's say your site has two sections: articles and tutorials. Create an archetype for each content type:

```text
archetypes/
├── articles.md
├── default.md
└── tutorials.md
```

To create an article using the articles archetype:

```sh
hugo new content articles/something.md
```

To create an article using the tutorials archetype:

```sh
hugo new content --kind tutorials articles/something.md
```

[content formats]: /getting-started/glossary/#content-format
[content types]: /getting-started/glossary/#content-type
[context]: /getting-started/glossary/#context
[front matter]: /getting-started/glossary/#front-matter
[template actions]: /getting-started/glossary/#template-action
[template]: /getting-started/glossary/#template
[template function]: /getting-started/glossary/#function
