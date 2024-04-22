---
title: Content organization
linkTitle: Organization
description: Hugo assumes that the same structure that works to organize your source content is used to organize the rendered site.
categories: [content management,fundamentals]
keywords: [sections,content,organization,bundle,resources]
menu:
  docs:
    parent: content-management
    weight: 20
weight: 20
toc: true
aliases: [/content/sections/]
---

## Page bundles

Hugo `0.32` announced page-relative images and other resources packaged into `Page Bundles`.

These terms are connected, and you also need to read about [Page Resources](/content-management/page-resources) and [Image Processing](/content-management/image-processing) to get the full picture.

{{< imgproc "1-featured-content-bundles.png" "resize 300x" >}}
The illustration shows three bundles. Note that the home page bundle cannot contain other content pages, although other files (images etc.) are allowed.
{{< /imgproc >}}

{{% note %}}
The bundle documentation is a **work in progress**. We will publish more comprehensive docs about this soon.
{{% /note %}}

## Organization of content source

In Hugo, your content should be organized in a manner that reflects the rendered website.

While Hugo supports content nested at any level, the top levels (i.e. `content/<DIRECTORIES>`) are special in Hugo and are considered the content type used to determine layouts etc. To read more about sections, including how to nest them, see [sections].

Without any additional configuration, the following will automatically work:

```txt
.
└── content
    └── about
    |   └── index.md  // <- https://example.org/about/
    ├── posts
    |   ├── firstpost.md   // <- https://example.org/posts/firstpost/
    |   ├── happy
    |   |   └── ness.md  // <- https://example.org/posts/happy/ness/
    |   └── secondpost.md  // <- https://example.org/posts/secondpost/
    └── quote
        ├── first.md       // <- https://example.org/quote/first/
        └── second.md      // <- https://example.org/quote/second/
```

## Path breakdown in Hugo

The following demonstrates the relationships between your content organization and the output URL structure for your Hugo website when it renders. These examples assume you are [using pretty URLs][pretty], which is the default behavior for Hugo. The examples also assume a key-value of `baseURL = "https://example.org/"` in your [site's configuration file][config].

### Index pages: `_index.md`

`_index.md` has a special role in Hugo. It allows you to add front matter and content to your [list templates][lists]. These templates include those for [section templates], [taxonomy templates], [taxonomy terms templates], and your [homepage template].

{{% note %}}
**Tip:** You can get a reference to the content and metadata in `_index.md` using the [`.Site.GetPage` function](/methods/page/getpage).
{{% /note %}}

You can create one `_index.md` for your homepage and one in each of your content sections, taxonomies, and taxonomy terms. The following shows typical placement of an `_index.md` that would contain content and front matter for a `posts` section list page on a Hugo website:

```txt
.         url
.       ⊢--^-⊣
.        path    slug
.       ⊢--^-⊣⊢---^---⊣
.           file path
.       ⊢------^------⊣
content/posts/_index.md
```

At build, this will output to the following destination with the associated values:

```txt

                     url ("/posts/")
                    ⊢-^-⊣
       baseurl      section ("posts")
⊢--------^---------⊣⊢-^-⊣
        permalink
⊢----------^-------------⊣
https://example.org/posts/index.html
```

The [sections] can be nested as deeply as you want. The important thing to understand is that to make the section tree fully navigational, at least the lower-most section must include a content file. (i.e. `_index.md`).

### Single pages in sections

Single content files in each of your sections will be rendered as [single page templates][singles]. Here is an example of a single `post` within `posts`:

```txt
                   path ("posts/my-first-hugo-post.md")
.       ⊢-----------^------------⊣
.      section        slug
.       ⊢-^-⊣⊢--------^----------⊣
content/posts/my-first-hugo-post.md
```

When Hugo builds your site, the content will be output to the following destination:

```txt

                               url ("/posts/my-first-hugo-post/")
                   ⊢------------^----------⊣
       baseurl     section     slug
⊢--------^--------⊣⊢-^--⊣⊢-------^---------⊣
                 permalink
⊢--------------------^---------------------⊣
https://example.org/posts/my-first-hugo-post/index.html
```

## Paths explained

The following concepts provide more insight into the relationship between your project's organization and the default Hugo behavior when building output for the website.

### `section`

A default content type is determined by the section in which a content item is stored. `section` is determined by the location within the project's `content` directory. `section` *cannot* be specified or overridden in front matter.

### `slug`

The `slug` is the last segment of the URL path, defined by the file name and optionally overridden by a `slug` value in front matter. See [URL Management](/content-management/urls/#slug) for details.

### `path`

A content's `path` is determined by the section's path to the file. The file `path`

* is based on the path to the content's location AND
* does not include the slug

### `url`

The `url` is the entire URL path, defined by the file path and optionally overridden by a `url` value in front matter. See [URL Management](/content-management/urls/#slug) for details.

[config]: /getting-started/configuration/
[formats]: /content-management/formats/
[front matter]: /content-management/front-matter/
[getpage]: /methods/page/getpage/
[homepage template]: /templates/homepage/
[homepage]: /templates/homepage/
[lists]: /templates/lists/
[pretty]: /content-management/urls/#appearance
[section templates]: /templates/section-templates/
[sections]: /content-management/sections/
[singles]: /templates/single-page-templates/
[taxonomy templates]: /templates/taxonomy-templates/
[taxonomy terms templates]: /templates/taxonomy-templates/
[types]: /content-management/types/
[urls]: /content-management/urls/
