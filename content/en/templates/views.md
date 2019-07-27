---
title: Content View Templates
# linktitle: Content Views
description: Hugo can render alternative views of your content, which is especially useful in list and summary views.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
keywords: [views]
menu:
  docs:
    parent: "templates"
    weight: 70
weight: 70
sections_weight: 70
draft: false
aliases: []
toc: true
---

These alternative **content views** are especially useful in [list templates][lists].

The following are common use cases for content views:

* You want content of every type to be shown on the homepage but only with limited [summary views][summaries].
* You only want a bulleted list of your content on a [taxonomy list page][taxonomylists]. Views make this very straightforward by delegating the rendering of each different type of content to the content itself.

## Create a Content View

To create a new view, create a template in each of your different content type directories with the view name. The following example contains an "li" view and a "summary" view for the `posts` and `project` content types. As you can see, these sit next to the [single content view][single] template, `single.html`. You can even provide a specific view for a given type and continue to use the `_default/single.html` for the primary view.

```
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

Hugo also has support for a default content template to be used in the event that a specific content view template has not been provided for that type. Content views can also be defined in the `_default` directory and will work the same as list and single templates who eventually trickle down to the `_default` directory as a matter of the lookup order.


```
▾ layouts/
  ▾ _default/
      li.html
      single.html
      summary.html
```

## Which Template Will be Rendered?

The following is the [lookup order][lookup] for content views:

1. `/layouts/<TYPE>/<VIEW>.html`
2. `/layouts/_default/<VIEW>.html`
3. `/themes/<THEME>/layouts/<TYPE>/<VIEW>.html`
4. `/themes/<THEME>/layouts/_default/<VIEW>.html`

## Example: Content View Inside a List

The following example demonstrates how to use content views inside of your [list templates][lists].

### `list.html`

In this example, `.Render` is passed into the template to call the [render function][render]. `.Render` is a special function that instructs content to render itself with the view template provided as the first argument. In this case, the template is going to render the `summary.html` view that follows:

{{< code file="layouts/_default/list.html" download="list.html" >}}
<main id="main">
  <div>
  <h1 id="title">{{ .Title }}</h1>
  {{ range .Pages }}
    {{ .Render "summary"}}
  {{ end }}
  </div>
</main>
{{< /code >}}

### `summary.html`

Hugo will pass the entire page object to the following `summary.html` view template. (See [Page Variables][pagevars] for a complete list.)

{{< code file="layouts/_default/summary.html" download="summary.html" >}}
<article class="post">
  <header>
    <h2><a href='{{ .Permalink }}'> {{ .Title }}</a> </h2>
    <div class="post-meta">{{ .Date.Format "Mon, Jan 2, 2006" }} - {{ .FuzzyWordCount }} Words </div>
  </header>
  {{ .Summary }}
  <footer>
  <a href='{{ .Permalink }}'><nobr>Read more →</nobr></a>
  </footer>
</article>
{{< /code >}}

### `li.html`

Continuing on the previous example, we can change our render function to use a smaller `li.html` view by changing the argument in the call to the `.Render` function (i.e., `{{ .Render "li" }}`).

{{< code file="layouts/_default/li.html" download="li.html" >}}
<li>
  <a href="{{ .Permalink }}">{{ .Title }}</a>
  <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
</li>
{{< /code >}}

[lists]: /templates/lists/
[lookup]: /templates/lookup-order/
[pagevars]: /variables/page/
[render]: /functions/render/
[single]: /templates/single-page-templates/
[spf]: https://spf13.com
[spfsourceli]: https://github.com/spf13/spf13.com/blob/master/layouts/_default/li.html
[spfsourcesection]: https://github.com/spf13/spf13.com/blob/master/layouts/_default/section.html
[spfsourcesummary]: https://github.com/spf13/spf13.com/blob/master/layouts/_default/summary.html
[summaries]: /content-management/summaries/
[taxonomylists]: /templates/taxonomy-templates/
