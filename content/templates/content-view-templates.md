---
title: Content View Templates
linktitle:
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [views]
weight: 70
draft: false
aliases: [/templates/views/]
toc: true
needsreview: true
---

In addition to the [single page content template][singletemplates], Hugo can render alternative views of your content. These are especially useful in [list templates][listtemplates].

Content views are appropriate for cases like the following:

* You want content of every type to be shown on the homepage but only with limited [summary views][summaries].
* You only want a bulleted list of your content on a [taxonomy list page][taxonomylists]. Views make this very straightforward by delegating the rendering of each different type of content to the content itself.

## Creating a Content View

To create a new view, simply create a template in each of your different
content type directories with the view name. In the following example, we
have created a "li" view and a "summary" view for our two content types
of post and project. As you can see, these sit next to the [single
content view](/templates/content/) template "single.html". You can even
provide a specific view for a given type and continue to use the
\_default/single.html for the primary view.

```bash
  ▾ layouts/
    ▾ post/
        li.html
        single.html
        summary.html
    ▾ project/
        li.html
        single.html
        summary.html
```

Hugo also has support for a default content template to be used in the event that a specific template has not been provided for that type. The default type works the same as the other types, but the directory must be called "_default". Content views can also be defined in the "_default" directory.


```bash
▾ layouts/
  ▾ _default/
      li.html
      single.html
      summary.html
```

## Which Template Will be Rendered?

Hugo uses a set of rules to figure out which template to use when rendering a specific page.

Hugo will use the following prioritized list. If a file isn’t present, then the next one in the list will be used. This enables you to craft specific layouts when you want to without creating more templates than necessary. For most sites only the \_default file at the end of the list will be needed.

* `/layouts/<TYPE>/<VIEW>.html`
* `/layouts/\_default/<VIEW>.html`
* `/themes/<THEME>/layouts/<TYPE>/<VIEW>.html`
* `/themes/<THEME>/layouts/\_default/<VIEW>.html`

## Example: Content View Inside a List

The following example demonstrates how to use content views inside of your [list page templates][listtemplates].

### `list.html`

In this example, `.Render` is passed into the template to call the [render function][]. `.Render` is a special function that instructs content to render itself with the view template provided as the first argument.

This `list.html` content view template is part of a larger `section.html` default template used for [spf13.com][spf]. ([See source on GitHub][spfsourcesection].)

{{% code file="layouts/_default/list.html" download="list.html" %}}
```
<section id="main">
  <div>
  <h1 id="title">{{ .Title }}</h1>
  {{ range .Data.Pages }}
    {{ .Render "summary"}}
  {{ end }}
  </div>
</section>
```
{{% /code %}}

### `summary.html`

Hugo will pass the entire page object to the view template. See [page
variables](/templates/variables/) for a complete list.

This `summary.html` content view template is used for [spf13.com][spf]. ([See source on GitHub][spfsourcesummary].)

{{% code file="layouts/_default/summary.html" download="summary.html" %}}
```html
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
```
{{% /code %}}

### `li.html`

Hugo will pass the entire page object to the view template. See [Page Variables][pagevars] for a complete list of variables Hugo makes available to you.

Continuing on the previous example, we can change our render function to use a smaller `li.html` view by changing the argument in the call to the `.Render` function (i.e., `{{ .Render "li" }}`).

This `li.html` content view template is used for [spf13.com][spf]. ([See source on GitHub][spfsourceli].)

{{% code file="layouts/_default/li.html" download="li.html" %}}
```html
<li>
  <a href="{{ .Permalink }}">{{ .Title }}</a>
  <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
</li>
```
{{% /code %}}

[listtemplates]: /templates/lists/
[pagevars]: /variables/page-variables/
[render function]: /functions/render/
[singletemplates]: /templates/single-page-templates/
[spf]: http://spf13.com
[spfsourceli]: https://github.com/spf13/spf13.com/blob/master/layouts/_default/li.html
[spfsourcesection]: https://github.com/spf13/spf13.com/blob/master/layouts/_default/section.html
[spfsourcesummary]: https://github.com/spf13/spf13.com/blob/master/layouts/_default/summary.html
[summaries]: /content-management/content-summaries/
[taxonomylists]: /templates/taxonomy-templates/