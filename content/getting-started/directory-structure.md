---
title: Directory Structure
linktitle: Directory Structure
description: Hugo's CLI scaffolds a project directory structure and then takes that single directory and uses it as the input for creating a complete website.
date: 2017-01-02
publishdate: 2017-02-01
lastmod: 2017-03-09
categories: [project organization]
tags: [source, organization, directories,fundamentals]
weight: 50
draft: false
aliases: [/overview/source-directory/]
toc: true
---

Hugo takes a single directory and uses it as the input for creating a complete
website.

## New Site Scaffolding

Running `hugo new site` from the command line will create a directory structure with the following elements:

```bash
.
├── archetypes
├── config.toml
├── content
├── data
├── layouts
├── static
└── themes
```


## Directory Structure Explained

The following is a high-level overview of each of the directories with links to each of their respective sections with in the Hugo docs.

[`archetypes`](/content-management/archetypes/)
: You can create new content files in Hugo using the `hugo new` command.
By default, hugo will create new content files with at least `date`, `title` (inferred from the file name), and `draft = true`. This saves time and promotes consistency for sites using multiple content types. You can create your own [archetypes][] with custom preconfigured front matter fields as well.

[`config.toml`](/getting-started/configuration/)
: Every Hugo project should have a configuration file in TOML, YAML, or JSON format at the root. Many sites may need little to no configuration, but Hugo ships with a large number of [configuration directives][] for more granular directions on how you want Hugo to build your website.

[`content`][]
: All content for your website will live inside this directory. Each top-level folder in Hugo is considered a [content section][]. For example, if your site has three main sections---`blog`, `articles`, and `tutorials`---you will have three directories at `content/blog`, `content/articles`, and `content/tutorials`. Hugo uses sections to assign default [content types][].

[`data`](/templates/data-templates/)
: This directory is used to store configuration files that can be
used by Hugo when generating your website. You can write these files in YAML, JSON, or TOML format. In addition to the files you add to this folder, you can also create [data templates][] that pull from dynamic content.

[`layouts`][]
: stores templates in the form of `.html` files that specify how views of your content will be rendered into a static website. Templates include [list pages][lists], your [homepage][], [taxonomy templates][], [partials][], [single page templates][singles], and more.

`static`
: stores all the static content for your future website: images, CSS, JavaScript, etc. Note that when Hugo build your site, all assets inside your static directory are copied over as-is.


## Example Hugo Project Directory

The following is an example of a typical Hugo project directory:

```bash
.
├── config.toml
├── archetypes
|   └── default.md
├── content
|   ├── post
|   |   ├── _index.md
|   |   ├── post-01.md
|   |   └── post-02.md
|   └── quote
|   |   ├── quote-01.md
|   |   └── quote-02.md
├── data
├── i18n
├── layouts
|   ├── _default
|   |   ├── single.html
|   |   └── list.html
|   ├── partials
|   |   ├── header.html
|   |   └── footer.html
|   ├── taxonomies
|   |   ├── category.html
|   |   ├── post.html
|   |   ├── quote.html
|   |   └── tag.html
|   ├── post
|   |   ├── li.html
|   |   ├── single.html
|   |   └── summary.html
|   ├── quote
|   |   ├── li.html
|   |   ├── single.html
|   |   └── summary.html
|   ├── shortcodes
|   |   ├── img.html
|   |   ├── vimeo.html
|   |   └── youtube.html
|   ├── index.html
|   └── sitemap.xml
├── themes
|   ├── hyde
|   └── doc
└── static
    ├── css
    ├── images
    └── js
```

The above directory structure tells us a lot about this project.

The rendered website

* has two different [types of content][types]: `posts` and `quotes`.
* applies two different [taxonomies][] to the content: `categories` and `tags`
* displays content in 3 different views: a list, a summary, and a full-page view

## Homepage and List Page Content

Since v0.18, [everything in Hugo is a `Page`][bepsays]. This means list pages and the homepage can have associated content files---i.e. `_index.md`---that contains page metadata (i.e., front matter) and content. This model allows you to include list-specific front matter via `.Params` and also means that list templates (e.g., `layouts/_default/list.html`) also have access to all [page variables][pagevars].

Using the above example, let's assume you have the following in `content/post/_index.md`:

{{% code file="content/post/_index.md" %}}
```yaml
---
title: My Golang Journey
date: 2017-03-23
publishdate: 2017-03-24
---

I decided to start learning Golang in March 2017.

Follow my journey through this new blog.
```
{{% /code %}}

You can now access this `_index.md` content in a [list template][lists]:

{{% code file="layouts/_default/list.html" %}}
```html
{{ define "main" }}
<main class="main">
    <article>
        <header>
            <h1>{{.Title}}</h1>
        </header>
        {{.Content}}
    </article>
    <ul class="section-contents">
    {{ range .Data.Pages }}
        <li>
            <a href="{{.Permalink}}">{{.Date.Format "2006-01-02"}} | {{.Title}}</a
        </li>
    {{ end }}
    </ul>
</main>
{{ end }}
```
{{% /code %}}

This will output the following HTML:

{{% code file="yoursite.com/post/index.html" copy="false" %}}
```html
<!--all your baseof.html code-->
<main class="main">
    <article>
        <header>
            <h1>My Golang Journey</h1>
        </header>
        <p>I decided to start learning Golang in March 2017.</p>
        <p>Follow my journey through this new blog.</p>
    </article>
    <ul class="section-contents">
        <li><a href="/post/post-01/">Post 1</a></li>
        <li><a href="/post/post-02/">Post 2</a></li>
    </ul>
</main>
<!--all your other baseof.html code-->
```
{{% /code %}}

### List Pages Without `_index.md`

You do *not* have to create an `_index.md` file for every list page (i.e. section, taxonomy, taxonomy terms, etc) or the homepage. If Hugo does not find an `_index.md` within the respective content section when rendering a [list template][lists], the page will be created but with no `{{.Content}}` and only the default values for `.Title` etc.

Using this same `layouts/_default/list.html` template and applying it to the the `quotes` section above will render the following output. Note that `quotes` does not have an `_index.md` file to pull from:

{{% code file="yoursite.com/quote/index.html" copy="false" %}}
```html
<!--baseof.html code-->
<main class="main">
    <article>
        <header>
            <h1>Quotes</h1>
        </header>
    </article>
    <ul class="section-contents">
        <li><a href="https://yoursite.com/quote/quotes-01/">Quote 1</a></li>
        <li><a href="https://yoursite.com/quote/quotes-02/">Quote 2</a></li>
    </ul>
</main>
<!--baseof.html code-->
```
{{% /code %}}

{{% note %}}
The default behavior of Hugo is to pluralize list titles; hence the inflection of the `quote` section to "Quotes" when called with the `.Title` [page variable](/variables/page/). You can change this via the `pluralizeListTitles` directive in your [site configuration](/getting-started/configuration/).
{{% /note %}}

[archetypes]: /content-management/archetypes/
[bepsays]: http://bepsays.com/en/2016/12/19/hugo-018/
[configuration directives]: /getting-started/configuration/#all-variables-yaml
[`content`]: /content-management/organization/
[content section]: /content-management/sections/
[content types]: /content-management/types/
[data templates]: /templates/data-templates/
[homepage]: /templates/homepage-templates/
[`layouts`]: /templates/
[lists]: /templates/list/
[pagevars]: /variables/page/
[partials]: /templates/partials/
[singles]: /templates/single-page-templates/
[taxonomies]: /content-management/taxonomies/
[taxonomy templates]: /templates/taxonomy-templates/
[types]: /content-management/types/