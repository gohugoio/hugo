---
title: Archetypes
linktitle: Archetypes
description: Archetypes allow you to create new instances of content types and set default parameters from the command line.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [archetypes,generators,metadata,front matter]
categories: ["content management"]
weight: 70
draft: false
aliases: [/content/archetypes/,/content-management/content-archetypes/]
toc: true
notesforauthors:
---

## What are Archetypes?

**Archetypes** are content files in the [archetypes directory][] of your project that contain preconfigured [front matter][] for your website's [content types][]. Archetypes facilitate consistent metadata across your website content and allow content authors to quickly generate instances of a content type via the `hugo new` command.

Hugo's generator assumes your working directory is the content folder at the root of your project. Hugo is able to infer the appropriate archetype by assuming the content type from the content section passed to the CLI command:

```bash
hugo new [content-section/file-name.md]
```

We can use this pattern to create a new `.md` file in the `posts` section of the [example site][]:

{{% input file="archetype-example.sh" %}}
```bash
hugo new posts/my-first-post.md
```
{{% /input %}}

{{% note "Override Content Type in a New File" %}}
To override the content type Hugo infers from `[content-section]`, add the `--kind` flag to the end of the `hugo new` command.
{{% /note %}}

Running this command in a new site that does not have default or custom archetypes will create the following file:

{{% output "content/posts/my-first-post.md" %}}
```toml
+++
date = "2017-02-01T19:20:04-07:00"
title = my first post
draft = true
+++
```
{{% /output %}}

Note that if you do not already have a `posts` directory, Hugo will create both `content/posts/` and `content/posts/my-first-post.md`.

`date` and `title` are the variables that ship with Hugo and are therefore included in *all* content files created with the Hugo CLI. `title` is generated from the new content's filename. `date` is generated in [RFC 3339 format][] by way of Golang's [`now()`][] function, which returns the current time. The third variable, `draft = true` is not carried over into your default archetype but has been added as a convenience to the Hugo's internal/base archetype.

Three variables per content file are often not enough for effective content management of larger websites. Luckily, Hugo provides a simple mechanism for extending the number of variables through default and custom archetypes.

## Lookup Order for Archetypes

Similar to the lookup order for [templates in the `layouts` directory][], Hugo looks for a default file before falling back on the base/internal archetype. For the `my-first-post.md` example, Hugo looks for the new content's archetype file in the following order:

1. `archetypes/posts.md`
2. `archetypes/default.md`
3. `themes/theme-name/archetypes/posts.md`
4. `themes/theme-name/archetypes/default.md`
5. `_internal` (i.e., `title` and `date`)

{{% note "Using a Theme Archetype" %}}
If you wish to use archetypes that ship with a theme, `theme` must be specified in your [configuration file](/project-organization/configuration/).
{{% /note %}}

## Choosing Your Front Matter Format

By default, `hugo new` content files include front matter in the TOML format regardless of the format used in `archetypes/*md`.

You can specify a different default format in your site [configuration file][] file using the `metaDataFormat` directive. Possible values are `toml`, `yaml`, and `json`.

## Default Archetypes

Default archetypes are convenient if your content's front matter stays consistent across multiple [content sections][].

### Creating the Default Archetype

The [example site][] includes `tags` and `categories` as [taxonomies][]. If we assume that all content files will require these two key-values, we can create a `default.md` archetype that *extends* Hugo's base archetype. In this example, we are including "golang" and "hugo" as tags and "web development" as a category.

{{% input file="archetypes/default.md" %}}
```toml
+++
tags = ["golang", "hugo"]
categories = ["web development"]
+++
```
{{% /input %}}

{{% warning "EOL Characters in Text Editors"%}}
If you get an `EOF error` when using `hugo new`, add a carriage return after the closing `+++` or `---` for your TOML or YAML front matter, respectively. (See [troubleshooting](/troubleshooting/eof-error/).)
{{% /warning %}}

### Using the Default Archetype

With an `archetypes/default.md` in place, we can use the CLI to create a new post in the `posts` content section:

{{% input file="new-post-from-default.sh" %}}
```bash
$ hugo new posts/my-new-post.md
```
{{% /input %}}

Hugo then creates a new markdown file with the following front matter:

{{% output "content/posts/my-new-post.md" %}}
```toml
+++
categories = ["web development"]
date = "2017-02-01T19:20:04-07:00"
tags = ["golang", "hugo"]
title = "my new post"
+++
```
{{% /output %}}

We see that the `title` and `date` key-values have been added in addition to the `tags` and `categories` key-values from `archetypes/default.md`.

{{% note "Ordering of Front Matter" %}}
You may notice that content files created with `hugo new` do not observe the order of the key-values specified in your archetype files and instead list your front matter alphabetically. This is a [known issue](https://github.com/spf13/hugo/issues/452).
{{% /note %}}

### Example Site Default Archetype

The following is the default archetype used in the [example site][].

{{< exfile "static/example/archetypes/default.md" "yaml">}}

## Custom Archetypes

Suppose the example site's `posts` section requires more sophisticated front matter than what has been specified in `archetypes/default.md`. We can create a custom archetype for our posts at `archetypes/posts.md` that includes the full set of front matter.

### Creating a Custom Archetype

{{% input file="archetypes/posts.md"%}}
```toml
+++
description = ""
tags = ""
categories = ""
+++
```
{{% /input %}}

### Using a Custom Archetype

With an `archetypes/posts.md` in place, we can use the CLI to create a new posts with custom `posts` metadata in the `posts` content section:

{{% input file="new-post-from-custom.sh" %}}
```bash
$ hugo new posts/post-from-custom.md
```
{{% /input %}}

This time, Hugo recognizes our custom `archetypes/posts.md` archetype and uses it instead of `archetypes/default.md`. The generated file will now include the full list of front matter parameters, as well as the base archetype's `title` and `date`.

{{% output "content/posts/post-from-custom.md" %}}
```toml
+++
categories = ""
date = 2017-02-13T17:24:43-08:00
description = ""
tags = ""
title = post from custom
+++
```
{{% /output %}}

### Example Site Custom Archetype

`musicians` in the [example site] require more sophisticated front matter than what has been specified in `archetypes/default.md`. We can create a custom archetype for musicians at `archetypes/musicians.md` that includes the full set of front matter.

The following is the `musicians` archetype from the [example site][]:

{{< exfile "static/example/archetypes/musicians.md" "yaml" >}}


[archetypes directory]: /project-organization/directory-structure/
[`now()`]: http://golang.org/pkg/time/#Now
[configuration file]: /project-organization/configuration/
[content sections]: /sections/
[content types]: /content-management/content-types/
[example site]: /getting-started/using-the-hugo-docs/#example-site
[front matter]: /content-management/front-matter/
[RFC 3339 format]: https://www.ietf.org/rfc/rfc3339.txt
[taxonomies]: /content-management/taxonomies/
[templates in the `layouts` directory]: /templates/base-templates-and-blocks/
[templates]: /templates/
