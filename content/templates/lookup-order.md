---
title: Hugo's Lookup Order
linktitle: Template Lookup Order
description: The lookup order is a prioritized list used by Hugo as it traverses your files looking for the appropriate corresponding file to render your content.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-05-25
categories: [templates,fundamentals]
keywords: [lookup]
menu:
  docs:
    parent: "templates"
    weight: 15
  quicklinks:
weight: 15
sections_weight: 15
draft: false
aliases: [/templates/lookup/]
toc: true
---

Before creating your templates, it's important to know how Hugo looks for files within your project's [directory structure][].

Hugo uses a prioritized list called the **lookup order** as it traverses your `layouts` folder in your Hugo project *looking* for the appropriate template to render your content.

The template lookup order is an inverted cascade: if template A isn’t present or specified, Hugo will look to template B. If template B isn't present or specified, Hugo will look for template C...and so on until it reaches the `_default/` directory for your project or theme. In many ways, the lookup order is similar to the programming concept of a [switch statement without fallthrough][switch].

The power of the lookup order is that it enables you to craft specific layouts and keep your templating [DRY][].

{{% note %}}
Most Hugo websites will only need the default template files at the end of the lookup order (i.e. `_default/*.html`).
{{% /note %}}

## Lookup Orders

The respective lookup order for each of Hugo's templates has been defined throughout the Hugo docs:

* [Homepage Template][home]
* [Base Templates][base]
* [Section Page Templates][sectionlookup]
* [Taxonomy List Templates][taxonomylookup]
* [Taxonomy Terms Templates][termslookup]
* [Single Page Templates][singlelookup]
* [RSS Templates][rsslookup]
* [Shortcode Templates][sclookup]

## Template Lookup Examples

The lookup order is best illustrated through examples. The following shows you the process Hugo uses for finding the appropriate template to render your [single page templates][], but the concept holds true for all templates in Hugo.

1. The project is using the theme `mytheme` (specified in the project's [configuration][config]).
2. The layouts and content directories for the project are as follows:

```
.
├── content
│   ├── events
│   │   ├── _index.md
│   │   └── my-first-event.md
│   └── posts
│       ├── my-first-post.md
│       └── my-second-post.md
├── layouts
│   ├── _default
│   │   └── single.html
│   ├── posts
│   │   └── single.html
│   └── reviews
│       └── reviewarticle.html
└── themes
    └── mytheme
        └── layouts
            ├── _default
            │   ├── list.html
            │   └── single.html
            └── posts
                ├── list.html
                └── single.html
```


Now we can look at the front matter for the three content (i.e.`.md`) files.

{{% note  %}}
Only three of the four markdown files in the above project are subject to the *single* page lookup order. `_index.md` is a specific `kind` in Hugo. Whereas `my-first-post.md`, `my-second-post.md`, and `my-first-event.md` are all of kind `page`, all `_index.md` files in a Hugo project are used to add content and front matter to [list pages](/templates/lists/). In this example, `events/_index.md` will render according to its [section template](/templates/section-templates/) and respective lookup order.
{{% /note %}}

### Example: `my-first-post.md`

{{< code file="content/posts/my-first-post.md" copy="false" >}}
---
title: My First Post
date: 2017-02-19
description: This is my first post.
---
{{< /code >}}

When building your site, Hugo will go through the lookup order until it finds what it needs for `my-first-post.md`:

1. ~~`/layouts/UNSPECIFIED/UNSPECIFIED.html`~~
2. ~~`/layouts/posts/UNSPECIFIED.html`~~
3. ~~`/layouts/UNSPECIFIED/single.html`~~
4. <span class="yes">`/layouts/posts/single.html`</span>
  <br><span class="break">BREAK</span>
5. <span class="na">`/layouts/_default/single.html`</span>
6. <span class="na">`/themes/<THEME>/layouts/UNSPECIFIED/UNSPECIFIED.html`</span>
7. <span class="na">`/themes/<THEME>/layouts/posts/UNSPECIFIED.html`</span>
8. <span class="na">`/themes/<THEME>/layouts/UNSPECIFIED/single.html`</span>
9. <span class="na">`/themes/<THEME>/layouts/posts/single.html`</span>
10. <span class="na">`/themes/<THEME>/layouts/_default/single.html`</span>

Notice the term `UNSPECIFIED` rather than `UNDEFINED`. If you don't tell Hugo the specific type and layout, it makes assumptions based on sane defaults. `my-first-post.md` does not specify a content `type` in its front matter. Therefore, Hugo assumes the content `type` and `section` (i.e. `posts`, which is defined by file location) are one in the same. ([Read more on sections][sections].)

`my-first-post.md` also does not specify a `layout` in its front matter. Therefore, Hugo assumes that `my-first-post.md`, which is of type `page` and a *single* piece of content, should default to the next occurrence of a `single.html` template in the lookup (#4).

### Example: `my-second-post.md`

{{< code file="content/posts/my-second-post.md" copy="false" >}}
---
title: My Second Post
date: 2017-02-21
description: This is my second post.
type: review
layout: reviewarticle
---
{{< /code >}}

Here is the way Hugo traverses the single-page lookup order for `my-second-post.md`:

1. <span class="yes">`/layouts/review/reviewarticle.html`</span>
  <br><span class="break">BREAK</span>
2. <span class="na">`/layouts/posts/reviewarticle.html`</span>
3. <span class="na">`/layouts/review/single.html`</span>
4. <span class="na">`/layouts/posts/single.html`</span>
5. <span class="na">`/layouts/_default/single.html`</span>
6. <span class="na">`/themes/<THEME>/layouts/review/reviewarticle.html`</span>
7. <span class="na">`/themes/<THEME>/layouts/posts/reviewarticle.html`</span>
8. <span class="na">`/themes/<THEME>/layouts/review/single.html`</span>
9. <span class="na">`/themes/<THEME>/layouts/posts/single.html`</span>
10. <span class="na">`/themes/<THEME>/layouts/_default/single.html`</span>

The front matter in `my-second-post.md` specifies the content `type` (i.e. `review`) as well as the `layout` (i.e. `reviewarticle`). Hugo finds the layout it needs at the top level of the lookup (#1) and does not continue to search through the other templates.

{{% note "Type and not Types" %}}
Notice that the directory for the template for `my-second-post.md` is `review` and not `reviews`. This is because *type is always singular when defined in front matter*.
{{% /note%}}

### Example: `my-first-event.md`

{{< code file="content/events/my-first-event.md" copy="false" >}}
---
title: My First
date: 2017-02-21
description: This is an upcoming event..
---
{{< /code >}}

Here is the way Hugo traverses the single-page lookup order for `my-first-event.md`:

1. ~~`/layouts/UNSPECIFIED/UNSPECIFIED.html`~~
2. ~~`/layouts/events/UNSPECIFIED.html`~~
3. ~~`/layouts/UNSPECIFIED/single.html`~~
4. ~~`/layouts/events/single.html`~~
5. <span class="yes">`/layouts/_default/single.html`</span>
<br><span class="break">BREAK</span>
6. <span class="na">`/themes/<THEME>/layouts/UNSPECIFIED/UNSPECIFIED.html`</span>
7. <span class="na">`/themes/<THEME>/layouts/events/UNSPECIFIED.html`</span>
8. <span class="na">`/themes/<THEME>/layouts/UNSPECIFIED/single.html`</span>
9. <span class="na">`/themes/<THEME>/layouts/events/single.html`</span>
10. <span class="na">`/themes/<THEME>/layouts/_default/single.html`</span>


{{% note %}}
`my-first-event.md` is significant because it demonstrates the role of the lookup order in Hugo themes. Both the root project directory *and* the `mytheme` themes directory have a file at `_default/single.html`. Understanding this order allows you to [customize Hugo themes](/themes/customizing/) by creating template files with identical names in your project directory that step in front of theme template files in the lookup. This allows you to customize the look and feel of your website while maintaining compatibility with the theme's upstream.
{{% /note %}}

[base]: /templates/base/#base-template-lookup-order
[config]: /getting-started/configuration/
[directory structure]: /getting-started/directory-structure/
[DRY]: https://en.wikipedia.org/wiki/Don%27t_repeat_yourself
[home]: /templates/homepage/#homepage-template-lookup-order
[rsslookup]: /templates/rss/#rss-template-lookup-order
[sclookup]: /templates/shortcode-templates/#shortcode-template-lookup-order
[sections]: /content-management/sections/
[sectionlookup]: /templates/section-templates/#section-template-lookup-order
[single page templates]: /templates/single-page-templates/
[singlelookup]: /templates/single-page-templates/#single-page-template-lookup-order
[switch]: https://en.wikipedia.org/wiki/Switch_statement#Fallthrough
[taxonomylookup]: /templates/taxonomy-templates/#taxonomy-list-template-lookup-order
[termslookup]: /templates/taxonomy-templates/#taxonomy-terms-template-lookup-order
