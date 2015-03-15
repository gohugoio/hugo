---
date: 2014-05-14T02:36:37Z
menu:
  main:
    parent: extras
next: /extras/permalinks
prev: /extras/livereload
title: Menus
weight: 60
---

Hugo has a simple yet powerful menu system that permits content to be
placed in menus with a good degree of control without a lot of work. 

Some of the features of Hugo Menus:

* Place content in one or many menus
* Handle nested menus with unlimited depth
* Create menu entries without being attached to any content
* Distinguish active element (and active branch)

## What is a menu?

A menu is a named array of menu entries accessible on the site under
`.Site.Menus` by name. For example, if I have a menu called `main`, I would
access it via `.Site.Menus.main`.

A menu entry has the following properties:

* **Url**        string
* **Name**       string
* **Menu**       string
* **Identifier** string
* **Pre**        template.HTML
* **Post**       template.HTML
* **Weight**     int
* **Parent**     string
* **Children**   Menu

And the following functions:

* **HasChildren** bool

Additionally, there are some relevant functions available on the page:

* **IsMenuCurrent** (menu string, menuEntry *MenuEntry ) bool
* **HasMenuCurrent** (menu string, menuEntry *MenuEntry) bool


## Adding content to menus

Hugo supports a couple of different methods of adding a piece of content
to the front matter.

### Simple

If all you need to do is add an entry to a menu, the simple form works
well.

**A single menu:**

    ---
    menu: "main"
    ---

**Multiple menus:**

    ---
    menu: ["main", "footer"]
    ---


### Advanced

If more control is required, then the advanced approach gives you the
control you want. All of the menu entry properties listed above are
available.

    ---
    menu:
      main:
        parent: 'extras'
        weight: 20
    ---


## Adding (non-content) entries to a menu

You can also add entries to menus that aren’t attached to a piece of
content. This takes place in the sitewide [config file](/overview/configuration/).

Here’s an example `config.toml`:

    [[menu.main]]
        name = "about hugo"
        pre = "<i class='fa fa-heart'></i>"
        weight = -110
        identifier = "about"
        url = "/about/"
    [[menu.main]]
        name = "getting started"
        pre = "<i class='fa fa-road'></i>"
        weight = -100
        url = "/getting-started/"

And the equivalent example `config.yaml`:

    ---
    menu:
      main:
          - Name: "about hugo"
            Pre: "<i class='fa fa-heart'></i>"
            Weight: -110
            Identifier: "about"
            Url: "/about/"
          - Name: "getting started"
            Pre: "<i class='fa fa-road'></i>"
            Weight: -100
            Url: "/getting-started/"
    ---            


**NOTE:** The urls must be relative to the context root. If the `BaseUrl` is `http://example.com/mysite/`, then the urls in the menu must not include the context root `mysite`. 
  
## Nesting

All nesting of content is done via the `parent` field.

The parent of an entry should be the identifier of another entry.
Identifier should be unique (within a menu).

The following order is used to determine identity Identifier > Name >
LinkTitle > Title. This means that the title will be used unless
linktitle is present, etc. In practice Name and Identifier are never
displayed and only used to structure relationships.

In this example, the top level of the menu is defined in the config file
and all content entries are attached to one of these entries via the
`parent` field.

## Rendering menus

Hugo makes no assumptions about how your rendered HTML will be
structured. Instead, it provides all of the functions you will need to be
able to build your menu however you want. 


The following is an example:

    <!--sidebar start-->
    <aside>
        <div id="sidebar" class="nav-collapse">
            <!-- sidebar menu start-->
            <ul class="sidebar-menu">
              {{ $currentNode := . }}
              {{ range .Site.Menus.main }}
                  {{ if .HasChildren }}

                <li class="sub-menu{{if $currentNode.HasMenuCurrent "main" . }} active{{end}}">
                <a href="javascript:;" class="">
                    {{ .Pre }}
                    <span>{{ .Name }}</span>
                    <span class="menu-arrow arrow_carrot-right"></span>
                </a>
                <ul class="sub">
                    {{ range .Children }}
                    <li{{if $currentNode.IsMenuCurrent "main" . }} class="active"{{end}}><a href="{{.Url}}"> {{ .Name }} </a> </li>
                    {{ end }}
                </ul>
              {{else}}
                <li>
                <a class="" href="{{.Url}}">
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
