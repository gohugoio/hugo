---
title: Archetypes
linktitle: Archetypes
description: Archetypes are templates used when creating new content.
date: 2017-02-01
publishdate: 2017-02-01
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

## What are Archetypes?

**Archetypes** are content template files in the [archetypes directory][] of your project that contain preconfigured [front matter][] and possibly also a content disposition for your website's [content types][]. These will be used when you run `hugo new`.


The `hugo new` uses the `content-section` to find the most suitable archetype template in your project. If your project does not contain any archetype files, it will also look in the theme.

{{< code file="archetype-example.sh" >}}
hugo new posts/my-first-post.md
{{< /code >}}

The above will create a new content file in `content/posts/my-first-post.md` using the first archetype file found of these:

1. `archetypes/posts.md`
2. `archetypes/default.md`
3. `themes/my-theme/archetypes/posts.md`
4. `themes/my-theme/archetypes/default.md`

The last two list items are only applicable if you use a theme and it uses the `my-theme` theme name as an example.

## Create a New Archetype Template

A fictional example for the section `newsletter` and the archetype file `archetypes/newsletter.md`. Create a new file in `archetypes/newsletter.md` and open it in a text editor.

{{< code file="archetypes/newsletter.md" >}}
---
title: "{{ replace .Name "-" " " | title }}"
date: {{ .Date }}
draft: true
---

**Insert Lead paragraph here.**

## New Cool Posts

{{ range first 10 ( where .Site.RegularPages "Type" "cool" ) }}
* {{ .Title }}
{{ end }}
{{< /code >}}

When you create a new newsletter with:

```bash
hugo new newsletter/the-latest-cool.stuff.md
```

It will create a new newsletter type of content file based on the archetype template.

**Note:** the site will only be built if the `.Site` is in use in the archetype file, and this can be time consuming for big sites.

The above _newsletter type archetype_ illustrates the possibilities: The full Hugo `.Site` and all of Hugo&#39;s template funcs can be used in the archetype file.


## Directory based archetypes

Since Hugo `0.49` you can use complete directories as archetype templates. Given this archetype directory:

```bash
archetypes
├── default.md
└── post-bundle
    ├── bio.md
    ├── images
    │   └── featured.jpg
    └── index.md
```

```bash
hugo new --kind post-bundle posts/my-post
```

Will create a new folder in `/content/posts/my-post` with the same set of files as in the `post-bundle` archetypes folder. All content files (`index.md` etc.) can contain template logic, and will receive the correct `.Site` for the content's language.



[archetypes directory]: /getting-started/directory-structure/
[content types]: /content-management/types/
[front matter]: /content-management/front-matter/
