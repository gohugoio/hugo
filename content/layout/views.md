---
title: "Content Views"
date: "2013-07-01"
weight: 70
menu:
  main:
    parent: 'layout'
---

In addition to the [single content view](/layout/content/), Hugo can render alternative views of
your content. These are especially useful in [index](/layout/indexes) templates.

To create a new view simple create a template in each of your different content
type directories with the view name. In the following example we have created a
"li" view and a "summary" view for our two content types of post and project. As
you can see these sit next to the [single content view](/layout/content)
template "single.html"

    ▾ layouts/
      ▾ post/
          li.html
          single.html
          summary.html
      ▾ project/
          li.html
          single.html
          summary.html

Hugo also has support for a default content template to be used in the event
that a specific template has not been provided for that type. The default type
works the same as the other types but the directory must be called "_default".
Content views can also be defined in the "_default" directory.


    ▾ layouts/
      ▾ _default/
          li.html
          single.html
          summary.html


Hugo will pass the entire page object to the view template. See [page
variables](/layout/variables) for a complete list.

## Example li.html
This content template is used for [spf13.com](http://spf13.com).

    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>

## Example summary.html
This content template is used for [spf13.com](http://spf13.com).

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


## Example render of view
Using the summary view inside of another ([index](/layout/index)) template.

    <section id="main">
    <div>
    <h1 id="title">{{ .Title }}</h1>
    {{ range .Data.Pages }}
    {{ .Render "summary"}}
    {{ end }}
    </div>
    </section>

In the above example you will notice that we have called .Render and passed in
which view to render the content with. Render is a special function available on
a content which tells the content to render itself with the provided view template.
