---
title: Single Page Templates
linktitle:
description: The primary view of content in Hugo is the single view. Hugo will render every Markdown file provided with a corresponding single template.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-04-06
categories: [templates]
keywords: [page]
menu:
  docs:
    parent: "templates"
    weight: 60
weight: 60
sections_weight: 60
draft: false
aliases: [/layout/content/]
toc: true
---

## Single Page Template Lookup Order

You can specify a [content's `type`][content type] and `layout` in a single content file's [front matter][]. However, you cannot specify `section` because this is determined based on file location (see [content section][section]).

Hugo assumes your content section and content type are the same unless you tell Hugo otherwise by providing a `type` directly in the front matter of a content file. This is why #1 and #3 come before #2 and #4, respectively, in the following lookup order. Values in angle brackets (`<>`) are variable.

1. `/layouts/<TYPE>/<LAYOUT>.html`
2. `/layouts/<SECTION>/<LAYOUT>.html`
3. `/layouts/<TYPE>/single.html`
4. `/layouts/<SECTION>/single.html`
5. `/layouts/_default/single.html`
6. `/themes/<THEME>/layouts/<TYPE>/<LAYOUT>.html`
7. `/themes/<THEME>/layouts/<SECTION>/<LAYOUT>.html`
8. `/themes/<THEME>/layouts/<TYPE>/single.html`
9. `/themes/<THEME>/layouts/<SECTION>/single.html`
10. `/themes/<THEME>/layouts/_default/single.html`

{{< youtube ZYQ5k0RQzmo >}}

## Example Single Page Templates

Content pages are of the type `page` and will therefore have all the [page variables][pagevars] and [site variables][] available to use in their templates.

### `post/single.html`

This single page template makes use of Hugo [base templates][], the [`.Format` function][] for dates, the [`.WordCount` page variable][pagevars], and ranges through the single content's specific [taxonomies][pagetaxonomy]. [`with`][] is also used to check whether the taxonomies are set in the front matter.

{{< code file="layouts/post/single.html" download="single.html" >}}
{{ define "main" }}
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
      <h5 id="wordcount"> {{ .WordCount }} Words </h5>
    </section>
    {{ with .Params.topics }}
    <ul id="topics">
      {{ range . }}
        <li><a href="{{ "topics" | absURL}}{{ . | urlize }}">{{ . }}</a> </li>
      {{ end }}
    </ul>
    {{ end }}
    {{ with .Params.tags }}
    <ul id="tags">
      {{ range . }}
        <li> <a href="{{ "tags" | absURL }}{{ . | urlize }}">{{ . }}</a> </li>
      {{ end }}
    </ul>
    {{ end }}
    </div>
    <div>
        {{ with .PrevInSection }}
          <a class="previous" href="{{.Permalink}}"> {{.Title}}</a>
        {{ end }}
        {{ with .NextInSection }}
          <a class="next" href="{{.Permalink}}"> {{.Title}}</a>
        {{ end }}
    </div>
</aside>
{{ end }}
{{< /code >}}

To easily generate new instances of a content type (e.g., new `.md` files in a section like `project/`) with preconfigured front matter, use [content archetypes][archetypes].

[archetypes]: /content-management/archetypes/
[base templates]: /templates/base/
[config]: /getting-started/configuration/
[content type]: /content-management/types/
[directory structure]: /getting-started/directory-structure/
[dry]: https://en.wikipedia.org/wiki/Don%27t_repeat_yourself
[`.Format` function]: /functions/format/
[front matter]: /content-management/front-matter/
[pagetaxonomy]: /templates/taxonomy-templates/#displaying-a-single-piece-of-content-s-taxonomies
[pagevars]: /variables/page/
[partials]: /templates/partials/
[section]: /content-management/sections/
[site variables]: /variables/site/
[spf13]: http://spf13.com/
[`with`]: /functions/with/
