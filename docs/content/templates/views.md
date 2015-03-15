---
aliases:
- /templates/views/
date: 2013-07-01
menu:
  main:
    parent: layout
next: /templates/partials
prev: /templates/terms
title: Content Views
weight: 70
---

In addition to the [single content template](/templates/content/), Hugo can render alternative views of
your content. These are especially useful in [list templates](/templates/list/).

For example you may want content of every type to be shown on the
homepage, but only a summary view of it there. Perhaps on a taxonomy
list page you would only want a bulleted list of your content. Views
make this very straightforward by delegating the rendering of each
different type of content to the content itself.


## Creating a content view

To create a new view, simply create a template in each of your different
content type directories with the view name. In the following example, we
have created a "li" view and a "summary" view for our two content types
of post and project. As you can see, these sit next to the [single
content view](/templates/content/) template "single.html". You can even
provide a specific view for a given type and continue to use the
\_default/single.html for the primary view.

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
works the same as the other types, but the directory must be called "_default".
Content views can also be defined in the "_default" directory.


    ▾ layouts/
      ▾ _default/
          li.html
          single.html
          summary.html


## Which Template will be rendered?
Hugo uses a set of rules to figure out which template to use when
rendering a specific page.

Hugo will use the following prioritized list. If a file isn’t present,
then the next one in the list will be used. This enables you to craft
specific layouts when you want to without creating more templates
than necessary. For most sites only the \_default file at the end of
the list will be needed.

* /layouts/`TYPE`/`VIEW`.html
* /layouts/\_default/`VIEW`.html
* /themes/`THEME`/layouts/`TYPE`/`VIEW`.html
* /themes/`THEME`/layouts/\_default/`view`.html


## Example using views

### rendering view inside of a list

Using the summary view (defined below) inside of a ([list
templates](/templates/list/)).

    <section id="main">
    <div>
    <h1 id="title">{{ .Title }}</h1>
    {{ range .Data.Pages }}
    {{ .Render "summary"}}
    {{ end }}
    </div>
    </section>

In the above example, you will notice that we have called `.Render` and passed in
which view to render the content with. `.Render` is a special function available on
a content which tells the content to render itself with the provided view template.
In this example, we are not using the li view. To use this we would
change the render line to `{{ .Render "li" }}`.


### li.html

Hugo will pass the entire page object to the view template. See [page
variables](/templates/variables/) for a complete list.

This content template is used for [spf13.com](http://spf13.com/).

    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>

### summary.html

Hugo will pass the entire page object to the view template. See [page
variables](/templates/variables/) for a complete list.

This content template is used for [spf13.com](http://spf13.com/).

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


