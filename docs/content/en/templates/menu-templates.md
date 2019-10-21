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
    title: "how to use menus in templates"
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
    <ul>
        {{ $currentPage := . }}
        {{ range .Site.Menus.main }}
            {{ if .HasChildren }}
                <li class="{{ if $currentPage.HasMenuCurrent "main" . }}active{{ end }}">
                    <a href="#">
                        {{ .Pre }}
                        <span>{{ .Name }}</span>
                    </a>
                </li>
                <ul class="sub-menu">
                    {{ range .Children }}
                        <li class="{{ if $currentPage.IsMenuCurrent "main" . }}active{{ end }}">
                            <a href="{{ .URL }}">{{ .Name }}</a>
                        </li>
                    {{ end }}
                </ul>
            {{ else }}
                <li>
                    <a href="{{ .URL }}">
                        {{ .Pre }}
                        <span>{{ .Name }}</span>
                    </a>
                </li>
            {{ end }}
        {{ end }}
        <li>
            <a href="#" target="_blank">Hardcoded Link 1</a>
        </li>
        <li>
            <a href="#" target="_blank">Hardcoded Link 2</a>
        </li>
    </ul>
</aside>
{{< /code >}}

{{% note "`absLangURL` and `relLangURL`" %}}
Use the [`absLangURL`](/functions/abslangurl) or [`relLangURL`](/functions/rellangurl) functions if your theme makes use of the [multilingual feature](/content-management/multilingual/). In contrast to `absURL` and `relURL`, these two functions add the correct language prefix to the url.
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
    <a class="sidebar-nav-item{{if or ($currentPage.IsMenuCurrent "main" .) ($currentPage.HasMenuCurrent "main" .) }} active{{end}}" href="{{ .URL }}" title="{{ .Title }}">{{ .Name }}</a>
    {{ end }}
</nav>
```

In the above, the menu item is marked as active if on the current section's list page or on a page in that section.


## Site Config menus

The above is all that's needed. But if you want custom menu items, e.g. changing weight, name, or link title attribute, you can define them manually in the site config file:

{{< code-toggle file="config" >}}
[[menu.main]]
    name = "This is the blog section"
    title = "blog section"
    weight = -110
    identifier = "blog"
    url = "/blog/"
{{</ code-toggle >}}

{{% note %}}
The `identifier` *must* match the section name.
{{% /note %}}

## Menu Entries from the Page's front matter

It's also possible to create menu entries from the page (i.e. the `.md`-file).

Here is a `yaml` example:

```
---
title: Menu Templates
linktitle: Menu Templates
menu:
  docs:
    title: "how to use menus in templates"
    parent: "templates"
    weight: 130
---
...
```

{{% note %}}
You can define more than one menu. It also doesn't have to be a complex value,
`menu` can also be a string, an array of strings, or an array of complex values
like in the example above.
{{% /note %}}

### Using .Page in Menus

If you use the front matter method of defining menu entries, you'll get access to the `.Page` variable.
This allows to use every variable that's reachable from the [page variable](/variables/page/).

This variable is only set when the menu entry is defined in the page's front matter.
Menu entries from the site config don't know anything about `.Page`.

That's why you have to use the go template's `with` keyword or something similar in your templating language.

Here's an example:

```
<nav class="sidebar-nav">
  {{ range .Site.Menus.main }}
    <a href="{{ .URL }}" title="{{ .Title }}">
      {{- .Name -}}
      {{- with .Page -}}
        <span class="date">
        {{- dateFormat " (2006-01-02)" .Date -}}
        </span>
      {{- end -}}
    </a>
  {{ end }}
</nav>
```
