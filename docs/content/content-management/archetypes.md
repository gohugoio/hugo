---
title: Archetypes
linktitle: Archetypes
description: Archetypes allow you to create new instances of content types and set default parameters from the command line.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
keywords: [archetypes,generators,metadata,front matter]
categories: ["content management"]
menu:
  docs:
    parent: "content-management"
    weight: 70
  quicklinks:
weight: 70	#rem
draft: false
aliases: [/content/archetypes/]
toc: true
---

{{% note %}}
This section is outdated, see https://github.com/gohugoio/hugoDocs/issues/11
{{% /note %}}
{{% todo %}}
See above
{{% /todo %}}

## What are Archetypes?

**Archetypes** are content files in the [archetypes directory][] of your project that contain preconfigured [front matter][] for your website's [content types][]. Archetypes facilitate consistent metadata across your website content and allow content authors to quickly generate instances of a content type via the `hugo new` command.

{{< youtube bcme8AzVh6o >}}

The `hugo new` generator for archetypes assumes your working directory is the content folder at the root of your project. Hugo is able to infer the appropriate archetype by assuming the content type from the content section passed to the CLI command:

```
hugo new <content-section>/<file-name.md>
```

We can use this pattern to create a new `.md` file in the `posts` section:

{{< code file="archetype-example.sh" >}}
hugo new posts/my-first-post.md
{{< /code >}}

{{% note "Override Content Type in a New File" %}}
To override the content type Hugo infers from `[content-section]`, add the `--kind` flag to the end of the `hugo new` command.
{{% /note %}}

Running this command in a new site that does not have default or custom archetypes will create the following file:

{{< output file="content/posts/my-first-post.md" >}}
+++
date = "2017-02-01T19:20:04-07:00"
title = "my first post"
draft = true
+++
{{< /output >}}

{{% note %}}
In this example, if you do not already have a `content/posts` directory, Hugo will create both `content/posts/` and `content/posts/my-first-post.md` for you.
{{% /note %}}

The  auto-populated fields are worth examining:

* `title` is generated from the new content's filename (i.e. in this case, `my-first-post` becomes `"my first post"`)
* `date` and `title` are the variables that ship with Hugo and are therefore included in *all* content files created with the Hugo CLI. `date` is generated in [RFC 3339 format][] by way of Go's [`now()`][] function, which returns the current time.
* The third variable, `draft = true`, is *not* inherited by your default or custom archetypes but is included in Hugo's automatically scaffolded `default.md` archetype for convenience.

Three variables per content file are often not enough for effective content management of larger websites. Luckily, Hugo provides a simple mechanism for extending the number of variables through custom archetypes, as well as default archetypes to keep content creation DRY.

## Lookup Order for Archetypes

Similar to the [lookup order for templates][lookup] in your `layouts` directory, Hugo looks for a section- or type-specific archetype, then a default archetype, and finally an internal archetype that ships with Hugo. For example, Hugo will look for an archetype for `content/posts/my-first-post.md` in the following order:

1. `archetypes/posts.md`
2. `archetypes/default.md`
3. `themes/<THEME>/archetypes/posts.md`
4. `themes/<THEME>/archetypes/default.md` (Auto-generated with `hugo new site`)

{{% note "Using a Theme Archetype" %}}
If you wish to use archetypes that ship with a theme, the `theme` field must be specified in your [configuration file](/getting-started/configuration/).
{{% /note %}}

## Choose Your Archetype's Front Matter Format

By default, `hugo new` content files include front matter in the TOML format regardless of the format used in `archetypes/*.md`.

You can specify a different default format in your site [configuration file][] file using the `metaDataFormat` directive. Possible values are `toml`, `yaml`, and `json`.

## Default Archetypes

Default archetypes are convenient if your content's front matter stays consistent across multiple [content sections][sections].

### Create the Default Archetype

When you create a new Hugo project using `hugo new site`, you'll notice that Hugo has already scaffolded a file at `archetypes/default.md`.

The following examples are from a site that's using `tags` and `categories` as [taxonomies][]. If we assume that all content files will require these two key-values, we can create a `default.md` archetype that *extends* Hugo's base archetype. In this example, we are including "golang" and "hugo" as tags and "web development" as a category.

{{< code file="archetypes/default.md" >}}
+++
tags = ["golang", "hugo"]
categories = ["web development"]
+++
{{< /code >}}

{{% warning "EOL Characters in Text Editors"%}}
If you get an `EOF error` when using `hugo new`, add a carriage return after the closing `+++` or `---` for your TOML or YAML front matter, respectively. (See the [troubleshooting article on EOF errors](/troubleshooting/eof-error/) for more information.)
{{% /warning %}}

### Use the Default Archetype

With an `archetypes/default.md` in place, we can use the CLI to create a new post in the `posts` content section:

{{< code file="new-post-from-default.sh" >}}
$ hugo new posts/my-new-post.md
{{< /code >}}

Hugo then creates a new markdown file with the following front matter:

{{< output file="content/posts/my-new-post.md" >}}
+++
categories = ["web development"]
date = "2017-02-01T19:20:04-07:00"
tags = ["golang", "hugo"]
title = "my new post"
+++
{{< /output >}}

We see that the `title` and `date` key-values have been added in addition to the `tags` and `categories` key-values from `archetypes/default.md`.

{{% note "Ordering of Front Matter" %}}
You may notice that content files created with `hugo new` do not respect the order of the key-values specified in your archetype files. This is a [known issue](https://github.com/gohugoio/hugo/issues/452).
{{% /note %}}

## Custom Archetypes

Suppose your site's `posts` section requires more sophisticated front matter than what has been specified in `archetypes/default.md`. You can create a custom archetype for your posts at `archetypes/posts.md` that includes the full set of front matter to be added to the two default archetypes fields.

### Create a Custom Archetype

{{< code file="archetypes/posts.md">}}
+++
description = ""
tags = ""
categories = ""
+++
{{< /code >}}

### Use a Custom Archetype

With an `archetypes/posts.md` in place, you can use the Hugo CLI to create a new post with your preconfigured front matter in the `posts` content section:

{{< code file="new-post-from-custom.sh" >}}
$ hugo new posts/post-from-custom.md
{{< /code >}}

This time, Hugo recognizes our custom `archetypes/posts.md` archetype and uses it instead of `archetypes/default.md`. The generated file will now include the full list of front matter parameters, as well as the base archetype's `title` and `date`:

{{< output file="content/posts/post-from-custom-archetype.md" >}}
+++
categories = ""
date = 2017-02-13T17:24:43-08:00
description = ""
tags = ""
title = "post from custom archetype"
+++
{{< /output >}}

### Hugo Docs Custom Archetype

As an example of archetypes in practice, the following is the `functions` archetype from the Hugo docs:

{{< code file="archetypes/functions.md" >}}
{{< readfile file="/archetypes/functions.md" >}}
{{< /code >}}

{{% note %}}
The preceding archetype is kept up to date with every Hugo build by using Hugo's [`readFile` function](/functions/readfile/). For similar examples, see [Local File Templates](/templates/files/).
{{% /note %}}

[archetypes directory]: /getting-started/directory-structure/
[`now()`]: http://golang.org/pkg/time/#Now
[configuration file]: /getting-started/configuration/
[sections]: /content-management/sections/
[content types]: /content-management/types/
[front matter]: /content-management/front-matter/
[RFC 3339 format]: https://www.ietf.org/rfc/rfc3339.txt
[taxonomies]: /content-management/taxonomies/
[lookup]: /templates/lookup/
[templates]: /templates/
