---
lastmod: 2015-08-04
date: 2014-05-14T02:36:37Z
toc: true
menu:
  main:
    parent: extras
next: /extras/pagination
prev: /extras/livereload
title: Menus
---

Hugo has a simple yet powerful menu system that permits content to be
placed in menus with a good degree of control without a lot of work.


*TIP:* If all you want is a simple menu for your sections, see [Section Menu for "the Lazy Blogger"]({{< relref "#section-menu-for-the-lazy-blogger" >}}).

Some of the features of Hugo Menus:

* Place content in one or many menus
* Handle nested menus with unlimited depth
* Create menu entries without being attached to any content
* Distinguish active element (and active branch)

## What is a menu?

A menu is a named array of menu entries accessible on the site under
`.Site.Menus` by name. For example, if I have a menu called `main`, I would
access it via `.Site.Menus.main`.

If you make use of the [multilingual feature]({{< relref "content/multilingual.md#menus">}}) you can define menus language independent.

A menu entry has the following properties:

* `URL        string`
* `Name       string`
* `Menu       string`
* `Identifier string`
* `Pre        template.HTML`
* `Post       template.HTML`
* `Weight     int`
* `Parent     string`
* `Children   Menu`

And the following functions:

* `HasChildren() bool`

Additionally, the `Page` object has two functions, which can be used when rendering menus:

* `IsMenuCurrent (menu string, menuEntry *MenuEntry ) bool`
* `HasMenuCurrent** (menu string, menuEntry *MenuEntry) bool`

See [Menu Functions](#menu-functions) for explanations of these functions, and [Rendering Nested Menus](#rendering-nested-menus) for an example of their use.

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
            URL: "/about/"
          - Name: "getting started"
            Pre: "<i class='fa fa-road'></i>"
            Weight: -100
            URL: "/getting-started/"
    ---


**NOTE:** The URLs must be relative to the context root. If the `baseURL` is `http://example.com/mysite/`, then the URLs in the menu must not include the context root `mysite`. Using an absolute URL will overide the baseURL. If the `URL` is `http://subdomain.example.com/`, the output will be `http://subdomain.example.com`.

## Nesting

All nesting of content is done via the `parent` field.

The parent of an entry should be the identifier of another entry.
Identifier should be unique (within a menu).

The following order is used to determine an Identifier:

> Name > LinkTitle > Title.

This means that the title will be used unless
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

> **Note**: use the `absLangURL` or `relLangURL` if your theme makes use of the [multilingual feature]({{< relref "content/multilingual.md" >}}). In contrast to `absURL` and `relURL` it adds the correct language prefix to the url. [Read more]({{< relref "templates/functions.md#urls" >}}).

## Section Menu for "the Lazy Blogger"

To enable this menu, add this to your site config, i.e. `config.toml`:

```
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




## Menu Functions

Suppose you have the menu structure shown below.

```
    [menu.main]
    │
    ├───colour
    │   │
    │   ├───warm
    │   │   ├───orange
    │   │   ├───red
    │   │   └───yellow
    │   │
    │   └───cool
    │       ├───blue
    │       ├───green
    │       └───purple
    │   
    └───tool
        ├───hammer
        ├───shovel
        └───saw
```

For each menu item, you can determine:

* If the menu item has any children: `.HasChildren()`
* If the menu item is a parent of the page you are currently rendering: `.Page.HasMenuCurrent()`
* If the menu item **is** the page you are currently rendering: `.Page.IsMenuCurrent()`

For example, if you are currently rendering the page `/colour/warm`, the values of `.HasChildren`, `HasMenuCurrent`, and `IsMenuCurrent` would be as shown below:

```

    [menu.main]           |                |               |             |
    │                     |                |               |             |
    ├───colour            | HasMenuCurrent |               | HasChildren |
    │   ├───<< WARM >>    |                | IsMenuCurrent | HasChildren |
    │   │   ├───orange    |                |               |             |
    │   │   ├───red       |                |               |             |
    │   │   └───yellow    |                |               |             |
    │   └───cool          |                |               | HasChildren |
    │       ├───blue      |                |               |             |
    │       ├───green     |                |               |             |
    │       └───purple    |                |               |             |
    └───tool              |                |               | HasChildren |
        ├───hammer        |                |               |             |
        ├───shovel        |                |               |             |
        └───saw           |                |               |             |
```

## Rendering nested menus

Hugo supports nested menus with as many levels as you like.

Nested menus can be rendered using a recursive partial template, such as the example below.

```
<!-- layouts/index.html, layouts/_default/single.html, ... -->
<h1>{{ .Title }}</h1>
<!-- Put this line in your main template, at the place where you want to
     render the menu. -->
{{ partial "menu_include.html" . }} 
```

```
<!-- layouts/partials/menu_include.html -->
{{ partial "menu_recursive.html" (dict "menu" .Site.Menus.main "page" . "site" .Site) }}
```

```
<!-- layouts/partials/menu_recursive.html -->
{{ $page := .page }}
{{ $site := .site }}
<ul>
{{ range .menu }}
  {{ $is := $page.IsMenuCurrent "main" . }}
  {{ $has := $page.HasMenuCurrent "main" . }}
  {{ if .HasChildren }} 
    <li>
      <a href="{{ .URL }}">
        {{ .Name }}
        {{ if $is }}[Is]{{ end }}
        {{ if $has }}[Has]{{ end }}
        {{ if .HasChildren }}[Children]{{ end }}
      </a>
        <!-- If the menu item has children, include this partial template again (recursively) -->
        {{ partial "menu_recursive.html" (dict "menu" .Children "page" $page "site" $site) }}
    </li>
  {{ else }}
    <li>
      <a href="{{ .URL }}">
        {{ .Name }}
        {{ if $is }}[Is]{{ end }}
        {{ if $has }}[Has]{{ end }}
        {{ if .HasChildren }}[Children]{{ end }}
      </a>
    </li>
  {{ end }}
{{ end }}
</ul>
```

This example code renders the words `[Is]`, `[Has]`, and `[Children]` to demonstrate how the `IsMenuCurrent()`, `HasMenuCurrent()`, and `HasChildren()` functions work.

You can customise this example to implement features such as:

* Highlight the current item, by applying a CSS style:

        <a href="{{ .URL }}"{{ if $is }} class="active"{{ end }}>
                            ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

* Highlight parents of the current item, by applying a CSS style:

        <a href="{{ .URL }}"{{ if $has }} class="parent-active"{{ end }}>
                            ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

* Only render sub-menus for parents of the current menu item, and the current menu item itself:

        {{ if or $is $has }}
            {{ partial "menu_recursive.html" (dict "menu" .Children "page" $page "site" $site) }}
        {{ end }}

A working example, implementing these features, is shown below.

```
<!-- layouts/partials/menu_recursive.html -->
{{ $page := .page }}
{{ $site := .site }}
<ul>
<!-- Menu items sorted alphabetically by name -->
{{ range .menu.ByName }}
  {{ $is := $page.IsMenuCurrent "main" . }}
  {{ $has := $page.HasMenuCurrent "main" . }}
  {{ if .HasChildren }} 
    <li>
      <a href="{{ .URL }}" class="{{ if $is }} active{{ end }}{{ if $has }} parent-active{{ end }}">
        {{ .Name }}
        <!-- Show a » symbol if there is a sub-menu we haven't rendered -->
        {{ if not (or $is $has) }}»{{ end }}
      </a>
      <!-- Only render sub-menu for parent items and the current item -->
      {{ if or $is $has }}
        {{ partial "menu_recursive.html" (dict "menu" .Children "page" $page "site" $site) }}
      {{ end }}
    </li>
  {{ else }}
    <li>
      <a href="{{ .URL }}" class="{{ if $is }}active{{end}}">{{ .Name }}</a>
    </li>
  {{ end }}
{{ end }}
</ul>
```
