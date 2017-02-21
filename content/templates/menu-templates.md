---
title: Menu Templates
linktitle: Menu Templates
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
weight:
categories: [templates]
tags: [lists,sections,menus]
draft: false
slug:
aliases: [/templates/menus/]
toc: false
needsreview: true
---


Hugo makes no assumptions about how your rendered HTML will be
structured. Instead, it provides all of the functions you will need to be
able to build your menu however you want.


The following is an example:

```html
<!--sidebar start-->
<aside>
    <div id="sidebar" class="nav-collapse">
        <!-- sidebar menu start-->
        <ul class="sidebar-menu">
          {{ $currentPage := . }}
          {{ range .Site.Menus.main }}
              {{ if .HasChildren }}

            <li class="sub-menu{{if $currentPage.HasMenuCurrent "main" . }} active{{end}}">
            <a href="javascript:;" class="">
                {{ .Pre }}
                <span>{{ .Name }}</span>
                <span class="menu-arrow arrow_carrot-right"></span>
            </a>
            <ul class="sub">
                {{ range .Children }}
                <li{{if $currentPage.IsMenuCurrent "main" . }} class="active"{{end}}><a href="{{.URL}}"> {{ .Name }} </a> </li>
                {{ end }}
            </ul>
          {{else}}
            <li>
            <a href="{{.URL}}">
                {{ .Pre }}
                <span>{{ .Name }}</span>
            </a>
          {{end}}
          </li>
          {{end}}
            <li> <a href="https://github.com/spf13/hugo/issues" target="blank">Questions and Issues</a> </li>
            <li> <a href="#" target="blank">Edit this Page</a> </li>
        </ul>
        <!-- sidebar menu end-->
    </div>
</aside>
<!--sidebar end-->
```

{{% note "`absLangURL` and `relLangURL`" %}}
Use the `absLangURL` or `relLangURL` if your theme makes use of the [multilingual feature](/content-management/multilingual-mode/). In contrast to `absURL` and `relURL`, these two functions add the correct language prefix to the url. [Read more](/functions/abslangurl).
{{% /note %}}

## Section Menu for "the Lazy Blogger"

To enable this menu, add the following to your site `config`:

```toml
SectionPagesMenu = "main"
```

The menu name can be anything, but take a note of what it is.

This will create a menu with all the sections as menu items and all the sections' pages as "shadow-members". The _shadow_ implies that the pages isn't represented by a menu-item themselves, but this enables you to create a top-level menu like this:

```
  <nav class="sidebar-nav">
        {{ $currentPage := . }}
        {{ range .Site.Menus.main }}
        <a class="sidebar-nav-item{{if or ($currentPage.IsMenuCurrent "main" .) ($currentPage.HasMenuCurrent "main" .) }} active{{end}}" href="{{.URL}}">{{ .Name }}</a>
        {{ end }}
    </nav>

```

In the above, the menu item is marked as active if on the current section's list page or on a page in that section.

The above is all that's needed. But if you want custom menu items, e.g. changing weight or name, you can define them manually in the site config, i.e. `config.toml`:

```
 [[menu.main]]
        name = "This is the blog section"
        weight = -110
        identifier = "blog"
        url = "/blog/"

```

**Note** that the `identifier` must match the section name.