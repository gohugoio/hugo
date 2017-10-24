---
title: Menu Templates
linktitle: Menu Templates
description: Menus are a powerful but simple feature for content management but can be easily manipulated in your templates to meet your design needs.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
keywords: [lists,sections,menus]
menu:
  docs:
    parent: "templates"
    weight: 130
weight: 130
sections_weight: 130
draft: false
aliases: [/templates/menus/]
toc: false
---

Hugo makes no assumptions about how your rendered HTML will be
structured. Instead, it provides all of the functions you will need to be
able to build your menu however you want.

The following is an example:

{{< code file="layouts/partials/sidebar.html" download="sidebar.html" >}}
<!-- sidebar start -->
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
            <li> <a href="https://github.com/gohugoio/hugo/issues" target="blank">Questions and Issues</a> </li>
            <li> <a href="#" target="blank">Edit this Page</a> </li>
        </ul>
        <!-- sidebar menu end-->
    </div>
</aside>
<!--sidebar end-->
{{< /code >}}

{{% note "`absLangURL` and `relLangURL`" %}}
Use the [`absLangUrl`](/functions/abslangurl) or [`relLangUrl`](/functions/rellangurl) functions if your theme makes use of the [multilingual feature](/content-management/multilingual/). In contrast to `absURL` and `relURL`, these two functions add the correct language prefix to the url.
{{% /note %}}

## Section Menu for Lazy Bloggers

To enable this menu, configure `sectionPagesMenu` in your site `config`:

```
sectionPagesMenu = "main"
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

{{% note %}}
The `identifier` *must* match the section name.
{{% /note %}}
