---
title: Lists of content in Hugo
linkTitle: List templates
description: Lists have a specific meaning and usage in Hugo when it comes to rendering your site homepage, section page, taxonomy list, or taxonomy terms list.
categories: [templates]
keywords: [lists,sections,rss,taxonomies,terms]
menu:
  docs:
    parent: templates
    weight: 60
weight: 60
toc: true
aliases: [/templates/list/,/layout/indexes/]
---

## What is a list page template?

{{< youtube 8b2YTSMdMps >}}

A list page template is a template used to render multiple pieces of content in a single HTML page. The exception to this rule is the homepage, which is still a list but has its own [dedicated template][homepage].

Hugo uses the term *list* in its truest sense; i.e. a sequential arrangement of material, especially in alphabetical or numerical order. Hugo uses list templates on any output HTML page where content is traditionally listed:

* [Home page](/templates/homepage)
* [Section pages](/templates/section-templates)
* [Taxonomy pages](/templates/taxonomy-templates)
* [Taxonomy term pages](/templates/taxonomy-templates)
* [RSS feeds](/templates/rss)
* [Sitemaps](/templates/sitemap-template)

For template lookup order, see [Template Lookup](/templates/lookup-order/).

The idea of a list page comes from the [hierarchical mental model of the web][mentalmodel] and is best demonstrated visually:

[![Image demonstrating a hierarchical website sitemap.](site-hierarchy.svg)](site-hierarchy.svg)

## List defaults

### Default templates

Since section lists and taxonomy lists (N.B., *not* [taxonomy terms lists][taxterms]) are both *lists* with regards to their templates, both have the same terminating default of `_default/list.html` or `themes/<THEME>/layouts/_default/list.html` in their lookup order. In addition, both [section lists][sectiontemps] and [taxonomy lists][taxlists] have their own default list templates in `_default`.

See [Template Lookup Order](/templates/lookup-order/) for the complete reference.

## Add content and front matter to list pages

Since v0.18, [everything in Hugo is a `Page`][bepsays]. This means list pages and the homepage can have associated content files (i.e. `_index.md`) that contain page metadata (i.e., front matter) and content.

This new model allows you to include list-specific front matter via `.Params` and also means that list templates (e.g., `layouts/_default/list.html`) have access to all [page variables][pagevars].

{{% note %}}
It is important to note that all `_index.md` content files will render according to a *list* template and not according to a [single page template](/templates/single-page-templates/).
{{% /note %}}

### Example project directory

The following is an example of a typical Hugo project directory's content:

```txt
.
...
├── content
|   ├── posts
|   |   ├── _index.md
|   |   ├── post-01.md
|   |   └── post-02.md
|   └── quote
|   |   ├── quote-01.md
|   |   └── quote-02.md
...
```

Using the above example, let's assume you have the following in `content/posts/_index.md`:

{{< code file=content/posts/_index.md >}}
---
title: My Go Journey
date: 2017-03-23
publishdate: 2017-03-24
---

I decided to start learning Go in March 2017.

Follow my journey through this new blog.
{{< /code >}}

You can now access this `_index.md`'s' content in your list template:

{{< code file=layouts/_default/list.html >}}
{{ define "main" }}
  <main>
    <article>
      <header>
        <h1>{{ .Title }}</h1>
      </header>
      <!-- "{{ .Content }}" pulls from the markdown content of the corresponding _index.md -->
      {{ .Content }}
    </article>
    <ul>
      <!-- Ranges through content/posts/*.md -->
      {{ range .Pages }}
        <li>
          <a href="{{ .Permalink }}">{{ .Date.Format "2006-01-02" }} | {{ .Title }}</a>
        </li>
      {{ end }}
    </ul>
  </main>
{{ end }}
{{< /code >}}

This above will output the following HTML:

{{< code file=example.com/posts/index.html >}}
<!--top of your baseof code-->
<main>
  <article>
    <header>
      <h1>My Go Journey</h1>
    </header>
    <p>I decided to start learning Go in March 2017.</p>
    <p>Follow my journey through this new blog.</p>
  </article>
  <ul>
    <li><a href="/posts/post-01/">Post 1</a></li>
    <li><a href="/posts/post-02/">Post 2</a></li>
  </ul>
</main>
<!--bottom of your baseof-->
{{< /code >}}

### List pages without `_index.md`

You do *not* have to create an `_index.md` file for every list page (i.e. section, taxonomy, taxonomy terms, etc) or the homepage. If Hugo does not find an `_index.md` within the respective content section when rendering a list template, the page will be created but with no `{{ .Content }}` and only the default values for `.Title` etc.

Using this same `layouts/_default/list.html` template and applying it to the `quotes` section above will render the following output. Note that `quotes` does not have an `_index.md` file to pull from:

{{< code file=example.com/quote/index.html >}}
<!--baseof-->
<main>
  <article>
    <header>
      <!-- Hugo assumes that .Title is the name of the section since there is no _index.md content file from which to pull a "title:" field -->
      <h1>Quotes</h1>
    </header>
  </article>
  <ul>
    <li><a href="https://example.org/quote/quotes-01/">Quote 1</a></li>
    <li><a href="https://example.org/quote/quotes-02/">Quote 2</a></li>
  </ul>
</main>
<!--baseof-->
{{< /code >}}

{{% note %}}
The default behavior of Hugo is to pluralize list titles; hence the inflection of the `quote` section to "Quotes" when called with the `.Title` [page variable](/variables/page/). You can change this via the `pluralizeListTitles` directive in your [site configuration](/getting-started/configuration/).
{{% /note %}}

## Example list templates

### Section template

This list template has been modified slightly from a template originally used in [spf13.com](https://spf13.com/). It makes use of [partial templates][partials] for the chrome of the rendered page rather than using a [base template][base]. The examples that follow also use the [content view templates][views] `li.html` or `summary.html`.

{{< code file=layouts/section/posts.html >}}
{{ partial "header.html" . }}
{{ partial "subheader.html" . }}
<main>
  <div>
    <h1>{{ .Title }}</h1>
    <ul>
      <!-- Renders the li.html content view for each content/posts/*.md -->
      {{ range .Pages }}
        {{ .Render "li" }}
      {{ end }}
    </ul>
  </div>
</main>
{{ partial "footer.html" . }}
{{< /code >}}

### Taxonomy template

{{< code file=layouts/_default/taxonomy.html >}}
{{ define "main" }}
<main>
  <div>
    <h1>{{ .Title }}</h1>
    <!-- ranges through each of the content files associated with a particular taxonomy term and renders the summary.html content view -->
    {{ range .Pages }}
      {{ .Render "summary" }}
    {{ end }}
  </div>
</main>
{{ end }}
{{< /code >}}

## Sort content

By default, Hugo sorts page collections by:

1. Page [weight]
2. Page [date] (descending)
3. Page [linkTitle], falling back to page [title]
4. Page file path if the page is backed by a file

[date]: /methods/page/date
[weight]: /methods/page/weight
[linkTitle]: /methods/page/linktitle
[title]: /methods/page/title

Change the sort order using any of the methods below.

{{< list-pages-in-section path=/methods/pages filter=methods_pages_sort filterType=include titlePrefix=. omitElementIDs=true >}}

## Group content

Group your content by field, parameter, or date using any of the methods below.

{{< list-pages-in-section path=/methods/pages filter=methods_pages_group filterType=include titlePrefix=. omitElementIDs=true >}}

## Filtering and limiting lists

Sometimes you only want to list a subset of the available content. A
common is to only display posts from [main sections]
on the blog's homepage.

See the documentation on [`where`] and
[`first`] for further details.

[base]: /templates/base/
[bepsays]: https://bepsays.com/en/2016/12/19/hugo-018/
[directorystructure]: /getting-started/directory-structure/
[`Format` function]: /methods/time/format/
[front matter]: /content-management/front-matter/
[getpage]: /methods/page/getpage/
[homepage]: /templates/homepage/
[mentalmodel]: https://webstyleguide.com/wsg3/3-information-architecture/3-site-structure.html
[pagevars]: /variables/page/
[partials]: /templates/partials/
[RSS 2.0]: https://cyber.harvard.edu/rss/rss.html
[rss]: /templates/rss/
[sections]: /content-management/sections/
[sectiontemps]: /templates/section-templates/
[sitevars]: /variables/site/
[taxlists]: /templates/taxonomy-templates/#taxonomy-list-templates
[taxterms]: /templates/taxonomy-templates/#taxonomy-terms-templates
[taxvars]: /variables/taxonomy/
[views]: /templates/views/
[`where`]: /functions/collections/where
[`first`]: /functions/collections/first/
[main sections]: /functions/collections/where#mainsections
[`time.Format`]: /functions/time/format
