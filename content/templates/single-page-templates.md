---
title: Single Page Templates
linktitle:
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [page]
weight: 60
draft: false
aliases: [/layout/content/]
toc: true
---

The primary view of content in Hugo is the single view. Hugo will render every Markdown file provided with a corresponding single template.

## Introduction to the Template Lookup Order

{{< lookupexplanation >}}

## Lookup Order for Single Page Templates

You can specify `type` (i.e., [content type][]) and `layout` in a single content file's [front matter][]. However, you cannot specify `section` because this is determined based on file location (see [content section][section]).

Hugo assumes your content section and content type are the same unless you tell Hugo otherwise by providing a `type` directly in the front matter of a content file. This is why #1 and #3 come before #2 and #4, respectively, in the following lookup order. Values in angle brackets (`<>`) are variable.

1. `/layouts/<TYPE>/<LAYOUT>.html`
2. `/layouts/<SECTION>/<LAYOUT>.html`
3. `/layouts/<TYPE>/single.html`
4. `/layouts/<SECTION>/single.html`
5. `/layouts/_default/single.html`
6. `/themes/<THEME>/layouts/<TYPE>/<LAYOUT.html`
7. `/themes/<THEME>/layouts/<SECTION/LAYOUT.html`
8. `/themes/<THEME>/layouts/<TYPE>/single.html`
9. `/themes/<THEME>/layouts/<SECTION>/single.html`
10. `/themes/<THEME>/layouts/_default/single.html`

## Single Page Lookup Examples

The following examples assume two things:

1. The project is using the theme `mytheme`, which would be specified as `theme: mytheme` or `theme = "mytheme` in the project's [`config.toml` or `config.yaml`][config], respectively.
2. The layouts and content directories for the project are as follows:

```bash
.
├── content
│   ├── events
│   │   ├── _index.md
│   │   └── my-first-event.md
│   └── posts
│       ├── my-first-post.md
│       └── my-second-post.md
├── layouts
│   ├── _default
│   │   └── single.html
│   ├── posts
│   │   └── single.html
│   └── reviews
│       └── reviewarticle.html
└── themes
    └── mytheme
        └── layouts
            ├── _default
            │   ├── list.html
            │   └── single.html
            └── posts
                ├── list.html
                └── single.html
```


Now we can look at the front matter for the three single-page content (i.e.`.md`) files.

{{% note "Three Content Pages but *Four* Markdown Files?" %}}
`_index.md` may seem like a single page of content but is actually a specific `kind` in Hugo. Whereas `my-first-post.md`, `my-second-post.md`, and `my-first-event.md` are all of kind `page`, all `_index.md` files in a Hugo project are of kind `section` and therefore do not submit themselves to the *single* page template lookup. Instead, `events/_index.md` will render according to its [section template](templates/section-templates/) and respective lookup order.
{{% /note %}}

### `my-first-post.md`

{{% input "content/posts/my-first-post.md" %}}
```yaml
---
title: My First Post
date: 2017-02-19
description: This is my first post.
---
```
{{% /input %}}

When it comes time for Hugo to render the content to the page, it will go through the single page template lookup order until it finds what it needs for `my-first-post.md`:

1. <span class="no">`/layouts/UNSPECIFIED/UNSPECIFIED.html`</span>
2. <span class="no">`/layouts/posts/UNSPECIFIED.html`</span>
3. <span class="no">`/layouts/UNSPECIFIED/single.html`</span>
4. <span class="yes">`/layouts/posts/single.html`</span>
  <br>**BREAK**
5. <span class="na">`/layouts/_default/single.html`</span>
6. <span class="na">`/themes/mytheme/layouts/UNSPECIFIED/UNSPECIFIED.html`</span>
7. <span class="na">`/themes/mytheme/layouts/posts/UNSPECIFIED.html`</span>
8. <span class="na">`/themes/mytheme/layouts/UNSPECIFIED/single.html`</span>
9. <span class="na">`/themes/mytheme/layouts/posts/single.html`</span>
10. <span class="na">`/themes/mytheme/layouts/_default/single.html`</span>

Notice the term `UNSPECIFIED` rather than `UNDEFINED`. If you don't tell Hugo the specific type and layout, it makes assumptions based on sane defaults. `my-first-post.md` does not specify a content `type` in its front matter. Therefore, Hugo assumes the content `type` and `section` (i.e. `posts`, which is defined by file location) are one in the same. ([Read more on sections][section].)

`my-first-post.md` also does not specify a `layout` in its front matter. Therefore, Hugo assumes that `my-first-post.md`, which is of type `page` and a *single* piece of content, should default to the next occurrence of a `single.html` template in the lookup (#4).

### `my-second-post.md`

{{% input "content/posts/my-second-post.md" %}}
```yaml
---
title: My Second Post
date: 2017-02-21
description: This is my second post.
type: review
layout: reviewarticle
---
```
{{% /input %}}

Here is the way Hugo's traverses the single-page lookup order for `my-second-post.md`:

1. <span class="yes">`/layouts/review/reviewarticle.html`</span>
  <br>**BREAK**
2. <span class="na">`/layouts/posts/reviewarticle.html`</span>
3. <span class="na">`/layouts/review/single.html`</span>
4. <span class="na">`/layouts/posts/single.html`</span>
5. <span class="na">`/layouts/_default/single.html`</span>
6. <span class="na">`/themes/mytheme/layouts/review/reviewarticle.html`</span>
7. <span class="na">`/themes/mytheme/layouts/posts/reviewarticle.html`</span>
8. <span class="na">`/themes/mytheme/layouts/review/single.html`</span>
9. <span class="na">`/themes/mytheme/layouts/posts/single.html`</span>
10. <span class="na">`/themes/mytheme/layouts/_default/single.html`</span>

The front matter in `my-second-post.md` specifies the content `type` (i.e. `review`) as well as the `layout` (i.e. `reviewarticle`). Hugo finds the layout it needs at the top level of the lookup (#1) and does not continue to search through the other templates.

{{% note "Type and not Types" %}}
Notice that the directory for the template for `my-second-post.md` is `review` and not `reviews`. This is because *type is always singular*.
{{% /note%}}

### `my-first-event.md`

{{% input "content/events/my-first-event.md" %}}
```yaml
---
title: My First
date: 2017-02-21
description: This is an upcoming event..
---
```
{{% /input %}}

Here is the way Hugo's traverses the single-page lookup order for `my-first-event.md`:

1. <span class="no">`/layouts/UNSPECIFIED/UNSPECIFIED.html`</span>
2. <span class="no">`/layouts/events/UNSPECIFIED.html`</span>
3. <span class="no">`/layouts/UNSPECIFIED/single.html`</span>
4. <span class="no">`/layouts/events/single.html`</span>
5. <span class="yes">`/layouts/_default/single.html`</span>
<br>**BREAK**
6. <span class="na">`/themes/mytheme/layouts/UNSPECIFIED/UNSPECIFIED.html`</span>
7. <span class="na">`/themes/mytheme/layouts/events/UNSPECIFIED.html`</span>
8. <span class="na">`/themes/mytheme/layouts/UNSPECIFIED/single.html`</span>
9. <span class="na">`/themes/mytheme/layouts/events/single.html`</span>
10. <span class="na">`/themes/mytheme/layouts/_default/single.html`</span>

{{% note %}}
`my-first-event.md` is significant because it demonstrates the role of the lookup order in Hugo themes. Both the root project directory *and* the `mytheme` themes directory have a file at `_default/single.html`. Understanding this order allows you to [customize Hugo themes](/themes/customizing-a-theme/) by creating template files with identical names in your project directory that step in front of theme template files in the lookup. This allows you to customize the look and feel of your website while maintaining compatibility with the theme's upstream.
{{% /note %}}

## Example Single Page Templates

Content pages are of the type `page` and will therefore have all the [page variables][] and [site variables][] available to use in their templates.

### `post/single.html`

This content template is used for [spf13.com][spf13]. It makes use of [partial templates][partials]:

{{% input "layouts/post/single.html" %}}
```html
{{ partial "header.html" . }}
{{ partial "subheader.html" . }}
{{ $baseURL := .Site.BaseURL }}
<section id="main">
  <h1 id="title">{{ .Title }}</h1>
  <div>
        <article id="content">
           {{ .Content }}
        </article>
  </div>
</section>

<aside id="meta">
    <div>
    <section>
      <h4 id="date"> {{ .Date.Format "Mon Jan 2, 2006" }} </h4>
      <h5 id="wc"> {{ .FuzzyWordCount }} Words </h5>
    </section>
    <ul id="categories">
      {{ range .Params.topics }}
        <li><a href="{{ $baseURL }}/topics/{{ . | urlize }}">{{ . }}</a> </li>
      {{ end }}
    </ul>
    <ul id="tags">
      {{ range .Params.tags }}
        <li> <a href="{{ $baseURL }}/tags/{{ . | urlize }}">{{ . }}</a> </li>
      {{ end }}
    </ul>
    </div>
    <div>
        {{ if .Prev }}
          <a class="previous" href="{{.Prev.Permalink}}"> {{.Prev.Title}}</a>
        {{ end }}
        {{ if .Next }}
          <a class="next" href="{{.Next.Permalink}}"> {{.Next.Title}}</a>
        {{ end }}
    </div>
</aside>
{{ partial "disqus.html" . }}
{{ partial "footer.html" . }}
```
{{% /input %}}

### `project/single.html`

This content template is also used for [spf13.com][spf13] and makes use of [partial templates][partials]:

{{% input "project/single.html" %}}
```html
  {{ partial "header.html" . }}
  {{ partial "subheader.html" . }}
  {{ $baseURL := .Site.BaseURL }}

  <section id="main">
    <h1 id="title">{{ .Title }}</h1>
    <div>
          <article id="content">
             {{ .Content }}
          </article>
    </div>
  </section>

  <aside id="meta">
      <div>
      <section>
        <h4 id="date"> {{ .Date.Format "Mon Jan 2, 2006" }} </h4>
        <h5 id="wc"> {{ .FuzzyWordCount }} Words </h5>
      </section>
      <ul id="categories">
        {{ range .Params.topics }}
        <li><a href="{{ $baseURL }}/topics/{{ . | urlize }}">{{ . }}</a> </li>
        {{ end }}
      </ul>
      <ul id="tags">
        {{ range .Params.tags }}
          <li> <a href="{{ $baseURL }}/tags/{{ . | urlize }}">{{ . }}</a> </li>
        {{ end }}
      </ul>
      </div>
  </aside>

  {{if isset .Params "project_url" }}
  <div id="ribbon">
      <a href="{{ index .Params "project_url" }}" rel="me">Fork me on GitHub</a>
  </div>
  {{ end }}

  {{ partial "footer.html" . }}
```
{{% /input %}}

Notice how `project/single.html` uses an additional parameter unique to this template. This doesn't need to be defined ahead of time. The key can wait to be used in the template if present in the content file's front matter.

To easily generate new instances of this content type (e.g., new `.md` files in `project/`) with preconfigured front matter, use [content archetypes][archetypes].

[archetypes]: /content-management/archetypes/
[config]: /getting-started/configuration/
[content type]: /content-management/content-types/
[directory structure]: /getting-started/directory-structure/
[dry]: https://en.wikipedia.org/wiki/Don%27t_repeat_yourself
[front matter]: /content-management/front-matter/
[page variables]: /variables-and-parms/page-variables/
[partials]: /templates/partial-templates/
[section]: /content-management/sections/
[site variables]: /variables-and-params/site-variables/
[spf13]: http://spf13.com/