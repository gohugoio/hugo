---
title: Content view templates
description: Hugo can render alternative views of your content, useful in list and summary views.
categories: [templates]
keywords: []
menu:
  docs:
    parent: templates
    weight: 120
weight: 120
toc: true
aliases: [/templates/views/]
---

The following are common use cases for content views:

* You want content of every type to be shown on the home page but only with limited [summary views][summaries].
* You only want a bulleted list of your content in a [taxonomy template]. Views make this very straightforward by delegating the rendering of each different type of content to the content itself.

## Create a content view

To create a new view, create a template in each of your different content type directories with the view name. The following example contains an "li" view and a "summary" view for the `posts` and `project` content types. As you can see, these sit next to the [single template], `single.html`. You can even provide a specific view for a given type and continue to use the `_default/single.html` for the primary view.

```txt
  ▾ layouts/
    ▾ posts/
        li.html
        single.html
        summary.html
    ▾ project/
        li.html
        single.html
      summary.html
```

## Which template will be rendered?

The following is the lookup order for content views ordered by specificity.

1. `/layouts/<TYPE>/<VIEW>.html`
1. `/layouts/<SECTION>/<VIEW>.html`
1. `/layouts/_default/<VIEW>.html`
1. `/themes/<THEME>/layouts/<TYPE>/<VIEW>.html`
1. `/themes/<THEME>/layouts/<SECTION>/<VIEW>.html`
1. `/themes/<THEME>/layouts/_default/<VIEW>.html`

## Example: content view inside a list

### `list.html`

In this example, `.Render` is passed into the template to call the [render function][render]. `.Render` is a special function that instructs content to render itself with the view template provided as the first argument. In this case, the template is going to render the `summary.html` view that follows:

{{< code file=layouts/_default/list.html >}}
<main id="main">
  <div>
    <h1 id="title">{{ .Title }}</h1>
    {{ range .Pages }}
      {{ .Render "summary" }}
    {{ end }}
  </div>
</main>
{{< /code >}}

### `summary.html`

Hugo passes the `Page` object to the following `summary.html` view template.

{{< code file=layouts/_default/summary.html >}}
<article class="post">
  <header>
    <h2><a href="{{ .RelPermalink }}">{{ .Title }}</a></h2>
    <div class="post-meta">{{ .Date.Format "Mon, Jan 2, 2006" }} - {{ .FuzzyWordCount }} Words </div>
  </header>
  {{ .Summary }}
  <footer>
  <a href='{{ .RelPermalink }}'>Read&nbsp;more&nbsp;&raquo;</a>
  </footer>
</article>
{{< /code >}}

### `li.html`

Continuing on the previous example, we can change our render function to use a smaller `li.html` view by changing the argument in the call to the `.Render` function (i.e., `{{ .Render "li" }}`).

{{< code file=layouts/_default/li.html >}}
<li>
  <a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
  <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
</li>
{{< /code >}}

[render]: /methods/page/render/
[single template]: /templates/types/#single
[summaries]: /content-management/summaries/
[taxonomy template]: /templates/types/#taxonomy
