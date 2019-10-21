---
title: Lists of Content in Hugo
linktitle: List Page Templates
description: Lists have a specific meaning and usage in Hugo when it comes to rendering your site homepage, section page, taxonomy list, or taxonomy terms list.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
keywords: [lists,sections,rss,taxonomies,terms]
menu:
  docs:
    parent: "templates"
    weight: 22
weight: 22
sections_weight: 22
draft: false
aliases: [/templates/list/,/layout/indexes/]
toc: true
---

## What is a List Page Template?

{{< youtube 8b2YTSMdMps >}}

A list page template is a template used to render multiple pieces of content in a single HTML page. The exception to this rule is the homepage, which is still a list but has its own [dedicated template][homepage].

Hugo uses the term *list* in its truest sense; i.e. a sequential arrangement of material, especially in alphabetical or numerical order. Hugo uses list templates on any output HTML page where content is traditionally listed:

* [Taxonomy terms pages][taxterms]
* [Taxonomy list pages][taxlists]
* [Section list pages][sectiontemps]
* [RSS][rss]

For template lookup order, see [Template Lookup](/templates/lookup-order/).

The idea of a list page comes from the [hierarchical mental model of the web][mentalmodel] and is best demonstrated visually:

[![Image demonstrating a hierarchical website sitemap.](/images/site-hierarchy.svg)](/images/site-hierarchy.svg)

## List Defaults

### Default Templates

Since section lists and taxonomy lists (N.B., *not* [taxonomy terms lists][taxterms]) are both *lists* with regards to their templates, both have the same terminating default of `_default/list.html` or `themes/<THEME>/layouts/_default/list.html` in their lookup order. In addition, both [section lists][sectiontemps] and [taxonomy lists][taxlists] have their own default list templates in `_default`.

See [Template Lookup Order](/templates/lookup-order/) for the complete reference.

## Add Content and Front Matter to List Pages

Since v0.18, [everything in Hugo is a `Page`][bepsays]. This means list pages and the homepage can have associated content files (i.e. `_index.md`) that contain page metadata (i.e., front matter) and content.

This new model allows you to include list-specific front matter via `.Params` and also means that list templates (e.g., `layouts/_default/list.html`) have access to all [page variables][pagevars].

{{% note %}}
It is important to note that all `_index.md` content files will render according to a *list* template and not according to a [single page template](/templates/single-page-templates/).
{{% /note %}}

### Example Project Directory

The following is an example of a typical Hugo project directory's content:

```
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

{{< code file="content/posts/_index.md" >}}
---
title: My Go Journey
date: 2017-03-23
publishdate: 2017-03-24
---

I decided to start learning Go in March 2017.

Follow my journey through this new blog.
{{< /code >}}

You can now access this `_index.md`'s' content in your list template:

{{< code file="layouts/_default/list.html" download="list.html" >}}
{{ define "main" }}
<main>
    <article>
        <header>
            <h1>{{.Title}}</h1>
        </header>
        <!-- "{{.Content}}" pulls from the markdown content of the corresponding _index.md -->
        {{.Content}}
    </article>
    <ul>
    <!-- Ranges through content/posts/*.md -->
    {{ range .Pages }}
        <li>
            <a href="{{.Permalink}}">{{.Date.Format "2006-01-02"}} | {{.Title}}</a>
        </li>
    {{ end }}
    </ul>
</main>
{{ end }}
{{< /code >}}

This above will output the following HTML:

{{< code file="example.com/posts/index.html" copy="false" >}}
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

### List Pages Without `_index.md`

You do *not* have to create an `_index.md` file for every list page (i.e. section, taxonomy, taxonomy terms, etc) or the homepage. If Hugo does not find an `_index.md` within the respective content section when rendering a list template, the page will be created but with no `{{.Content}}` and only the default values for `.Title` etc.

Using this same `layouts/_default/list.html` template and applying it to the `quotes` section above will render the following output. Note that `quotes` does not have an `_index.md` file to pull from:

{{< code file="example.com/quote/index.html" copy="false" >}}
<!--baseof-->
<main>
    <article>
        <header>
        <!-- Hugo assumes that .Title is the name of the section since there is no _index.md content file from which to pull a "title:" field -->
            <h1>Quotes</h1>
        </header>
    </article>
    <ul>
        <li><a href="https://example.com/quote/quotes-01/">Quote 1</a></li>
        <li><a href="https://example.com/quote/quotes-02/">Quote 2</a></li>
    </ul>
</main>
<!--baseof-->
{{< /code >}}

{{% note %}}
The default behavior of Hugo is to pluralize list titles; hence the inflection of the `quote` section to "Quotes" when called with the `.Title` [page variable](/variables/page/). You can change this via the `pluralizeListTitles` directive in your [site configuration](/getting-started/configuration/).
{{% /note %}}

## Example List Templates

### Section Template

This list template has been modified slightly from a template originally used in [spf13.com](http://spf13.com/). It makes use of [partial templates][partials] for the chrome of the rendered page rather than using a [base template][base]. The examples that follow also use the [content view templates][views] `li.html` or `summary.html`.

{{< code file="layouts/section/posts.html" >}}
{{ partial "header.html" . }}
{{ partial "subheader.html" . }}
<main>
  <div>
   <h1>{{ .Title }}</h1>
        <ul>
        <!-- Renders the li.html content view for each content/posts/*.md -->
            {{ range .Pages }}
                {{ .Render "li"}}
            {{ end }}
        </ul>
  </div>
</main>
{{ partial "footer.html" . }}
{{< /code >}}

### Taxonomy Template

{{< code file="layouts/_default/taxonomy.html" download="taxonomy.html" >}}
{{ define "main" }}
<main>
  <div>
   <h1>{{ .Title }}</h1>
   <!-- ranges through each of the content files associated with a particular taxonomy term and renders the summary.html content view -->
    {{ range .Pages }}
        {{ .Render "summary"}}
    {{ end }}
  </div>
</main>
{{ end }}
{{< /code >}}

## Order Content

Hugo lists render the content based on metadata you provide in [front matter][]. In addition to sane defaults, Hugo also ships with multiple methods to make quick work of ordering content inside list templates:

### Default: Weight > Date > LinkTitle > FilePath

{{< code file="layouts/partials/default-order.html" >}}
<ul>
    {{ range .Pages }}
        <li>
            <h1><a href="{{ .Permalink }}">{{ .Title }}</a></h1>
            <time>{{ .Date.Format "Mon, Jan 2, 2006" }}</time>
        </li>
    {{ end }}
</ul>
{{< /code >}}

### By Weight

Lower weight gets higher precedence. So content with lower weight will come first.

{{< code file="layouts/partials/by-weight.html" >}}
<ul>
    {{ range .Pages.ByWeight }}
        <li>
            <h1><a href="{{ .Permalink }}">{{ .Title }}</a></h1>
            <time>{{ .Date.Format "Mon, Jan 2, 2006" }}</time>
        </li>
    {{ end }}
</ul>
{{< /code >}}

### By Date

{{< code file="layouts/partials/by-date.html" >}}
<ul>
    <!-- orders content according to the "date" field in front matter -->
    {{ range .Pages.ByDate }}
        <li>
            <h1><a href="{{ .Permalink }}">{{ .Title }}</a></h1>
            <time>{{ .Date.Format "Mon, Jan 2, 2006" }}</time>
        </li>
    {{ end }}
</ul>
{{< /code >}}

### By Publish Date

{{< code file="layouts/partials/by-publish-date.html" >}}
<ul>
    <!-- orders content according to the "publishdate" field in front matter -->
    {{ range .Pages.ByPublishDate }}
        <li>
            <h1><a href="{{ .Permalink }}">{{ .Title }}</a></h1>
            <time>{{ .Date.Format "Mon, Jan 2, 2006" }}</time>
        </li>
    {{ end }}
</ul>
{{< /code >}}

### By Expiration Date

{{< code file="layouts/partials/by-expiry-date.html" >}}
<ul>
    {{ range .Pages.ByExpiryDate }}
        <li>
            <h1><a href="{{ .Permalink }}">{{ .Title }}</a></h1>
            <time>{{ .Date.Format "Mon, Jan 2, 2006" }}</time>
        </li>
    {{ end }}
</ul>
{{< /code >}}

### By Last Modified Date

{{< code file="layouts/partials/by-last-mod.html" >}}
<ul>
    <!-- orders content according to the "lastmod" field in front matter -->
    {{ range .Pages.ByLastmod }}
        <li>
            <h1><a href="{{ .Permalink }}">{{ .Title }}</a></h1>
            <time>{{ .Date.Format "Mon, Jan 2, 2006" }}</time>
        </li>
    {{ end }}
</ul>
{{< /code >}}

### By Length

{{< code file="layouts/partials/by-length.html" >}}
<ul>
    <!-- orders content according to content length in ascending order (i.e., the shortest content will be listed first) -->
    {{ range .Pages.ByLength }}
        <li>
            <h1><a href="{{ .Permalink }}">{{ .Title }}</a></h1>
            <time>{{ .Date.Format "Mon, Jan 2, 2006" }}</time>
        </li>
    {{ end }}
</ul>
{{< /code >}}

### By Title

{{< code file="layouts/partials/by-title.html" >}}
<ul>
    <!-- ranges through content in ascending order according to the "title" field set in front matter -->
    {{ range .Pages.ByTitle }}
        <li>
            <h1><a href="{{ .Permalink }}">{{ .Title }}</a></h1>
            <time>{{ .Date.Format "Mon, Jan 2, 2006" }}</time>
        </li>
    {{ end }}
</ul>
{{< /code >}}

### By Link Title

{{< code file="layouts/partials/by-link-title.html" >}}
<ul>
    <!-- ranges through content in ascending order according to the "linktitle" field in front matter. If a "linktitle" field is not set, the range will start with content that only has a "title" field and use that value for .LinkTitle -->
    {{ range .Pages.ByLinkTitle }}
        <li>
            <h1><a href="{{ .Permalink }}">{{ .LinkTitle }}</a></h1>
            <time>{{ .Date.Format "Mon, Jan 2, 2006" }}</time>
        </li>
    {{ end }}
</ul>
{{< /code >}}

### By Parameter

Order based on the specified front matter parameter. Content that does not have the specified front matter field  will use the site's `.Site.Params` default. If the parameter is not found at all in some entries, those entries will appear together at the end of the ordering.

{{< code file="layouts/partials/by-rating.html" >}}
<!-- Ranges through content according to the "rating" field set in front matter -->
{{ range (.Pages.ByParam "rating") }}
  <!-- ... -->
{{ end }}
{{< /code >}}

If the targeted front matter field is nested beneath another field, you can access the field using dot notation.

{{< code file="layouts/partials/by-nested-param.html" >}}
{{ range (.Pages.ByParam "author.last_name") }}
  <!-- ... -->
{{ end }}
{{< /code >}}

### Reverse Order

Reversing order can be applied to any of the above methods. The following uses `ByDate` as an example:

{{< code file="layouts/partials/by-date-reverse.html" >}}
<ul>
    {{ range .Pages.ByDate.Reverse }}
        <li>
            <h1><a href="{{ .Permalink }}">{{ .Title }}</a></h1>
            <time>{{ .Date.Format "Mon, Jan 2, 2006" }}</time>
        </li>
    {{ end }}
</ul>
{{< /code >}}

## Group Content

Hugo provides some functions for grouping pages by Section, Type, Date, etc.

### By Page Field

{{< code file="layouts/partials/by-page-field.html" >}}
<!-- Groups content according to content section. The ".Key" in this instance will be the section's title. -->
{{ range .Pages.GroupBy "Section" }}
<h3>{{ .Key }}</h3>
<ul>
    {{ range .Pages }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}
</ul>
{{ end }}
{{< /code >}}

In the above example, you may want `{{.Title}}` to point the `title` field you have added to your `_index.md` file instead. You can access this value using the [`.GetPage` function][getpage]:

{{< code file="layouts/partials/by-page-field.html" >}}
<!-- Groups content according to content section.-->
{{ range .Pages.GroupBy "Section" }}
<!-- Checks for existence of _index.md for a section; if available, pulls from "title" in front matter -->
{{ with $.Site.GetPage "section" .Key }}
<h3>{{.Title}}</h3>
{{ else }}
<!-- If no _index.md is available, ".Key" defaults to the section title and filters to title casing -->
<h3>{{ .Key | title }}</h3>
{{ end }}
<ul>
    {{ range .Pages }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}
</ul>
{{ end }}
{{< /code >}}

### By Date

{{< code file="layouts/partials/by-page-date.html" >}}
<!-- Groups content by month according to the "date" field in front matter -->
{{ range .Pages.GroupByDate "2006-01" }}
<h3>{{ .Key }}</h3>
<ul>
    {{ range .Pages }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}
</ul>
{{ end }}
{{< /code >}}

### By Publish Date

{{< code file="layouts/partials/by-page-publish-date.html" >}}
<!-- Groups content by month according to the "publishdate" field in front matter -->
{{ range .Pages.GroupByPublishDate "2006-01" }}
<h3>{{ .Key }}</h3>
<ul>
    {{ range .Pages }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .PublishDate.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}
</ul>
{{ end }}
{{< /code >}}

### By Page Parameter

{{< code file="layouts/partials/by-page-param.html" >}}
<!-- Groups content according to the "param_key" field in front matter -->
{{ range .Pages.GroupByParam "param_key" }}
<h3>{{ .Key }}</h3>
<ul>
    {{ range .Pages }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}
</ul>
{{ end }}
{{< /code >}}

### By Page Parameter in Date Format

The following template takes grouping by `date` a step further and uses Go's layout string. See the [`Format` function][] for more examples of how to use Go's layout string to format dates in Hugo.

{{< code file="layouts/partials/by-page-param-as-date.html" >}}
<!-- Groups content by month according to the "param_key" field in front matter -->
{{ range .Pages.GroupByParamDate "param_key" "2006-01" }}
<h3>{{ .Key }}</h3>
<ul>
    {{ range .Pages }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}
</ul>
{{ end }}
{{< /code >}}

### Reverse Key Order

Ordering of groups is performed by keys in alphanumeric order (A–Z, 1–100) and in reverse chronological order (i.e., with the newest first) for dates.

While these are logical defaults, they are not always the desired order. There are two different syntaxes to change Hugo's default ordering for groups, both of which work the same way.

#### 1. Adding the Reverse Method

```
{{ range (.Pages.GroupBy "Section").Reverse }}
```

```
{{ range (.Pages.GroupByDate "2006-01").Reverse }}
```

#### 2. Providing the Alternate Direction

```
{{ range .Pages.GroupByDate "2006-01" "asc" }}
```

```
{{ range .Pages.GroupBy "Section" "desc" }}
```

### Order Within Groups

Because Grouping returns a `{{.Key}}` and a slice of pages, all of the ordering methods listed above are available.

Here is the ordering for the example that follows:

1. Content is grouped by month according to the `date` field in front matter.
2. Groups are listed in ascending order (i.e., the oldest groups first)
3. Pages within each respective group are ordered alphabetically according to the `title`.

{{< code file="layouts/partials/by-group-by-page.html" >}}
{{ range .Pages.GroupByDate "2006-01" "asc" }}
<h3>{{ .Key }}</h3>
<ul>
    {{ range .Pages.ByTitle }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}
</ul>
{{ end }}
{{< /code >}}

## Filtering and Limiting Lists {#filtering-and-limiting-lists}

Sometimes you only want to list a subset of the available content. A
common is to only display posts from [**main sections**][mainsections]
on the blog's homepage.

See the documentation on [`where` function][wherefunction] and
[`first` function][firstfunction] for further details.

[base]: /templates/base/
[bepsays]: http://bepsays.com/en/2016/12/19/hugo-018/
[directorystructure]: /getting-started/directory-structure/
[`Format` function]: /functions/format/
[front matter]: /content-management/front-matter/
[getpage]: /functions/getpage/
[homepage]: /templates/homepage/
[homepage]: /templates/homepage/
[mentalmodel]: http://webstyleguide.com/wsg3/3-information-architecture/3-site-structure.html
[pagevars]: /variables/page/
[partials]: /templates/partials/
[RSS 2.0]: http://cyber.law.harvard.edu/rss/rss.html "RSS 2.0 Specification"
[rss]: /templates/rss/
[sections]: /content-management/sections/
[sectiontemps]: /templates/section-templates/
[sitevars]: /variables/site/
[taxlists]: /templates/taxonomy-templates/#taxonomy-list-templates
[taxterms]: /templates/taxonomy-templates/#taxonomy-terms-templates
[taxvars]: /variables/taxonomy/
[views]: /templates/views/
[wherefunction]: /functions/where/
[firstfunction]: /functions/first/
[mainsections]: /functions/where/#mainsections
