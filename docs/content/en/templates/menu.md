---
title: Menus
description: Create templates to render one or more menus.
categories: [templates]
keywords: []
menu:
  docs:
    identifier: templates-menu
    parent: templates
    weight: 180
weight: 180
toc: true
aliases: [/templates/menus/,/templates/menu-templates/]
---

## Overview

After [defining menu entries], use [menu methods] to render a menu.

Three factors determine how to render a menu:

1. The method used to define the menu entries: [automatic], [in front matter], or [in site configuration]
1. The menu structure: flat or nested
1. The method used to [localize the menu entries]: site configuration or translation tables

The example below handles every combination.

## Example

This partial template recursively "walks" a menu structure, rendering a localized, accessible nested list.

{{< code file=layouts/partials/menu.html copy=true >}}
{{- $page := .page }}
{{- $menuID := .menuID }}

{{- with index site.Menus $menuID }}
  <nav>
    <ul>
      {{- partial "inline/menu/walk.html" (dict "page" $page "menuEntries" .) }}
    </ul>
  </nav>
{{- end }}

{{- define "partials/inline/menu/walk.html" }}
  {{- $page := .page }}
  {{- range .menuEntries }}
    {{- $attrs := dict "href" .URL }}
    {{- if $page.IsMenuCurrent .Menu . }}
      {{- $attrs = merge $attrs (dict "class" "active" "aria-current" "page") }}
    {{- else if $page.HasMenuCurrent .Menu .}}
      {{- $attrs = merge $attrs (dict "class" "ancestor" "aria-current" "true") }}
    {{- end }}
    {{- $name := .Name }}
    {{- with .Identifier }}
      {{- with T . }}
        {{- $name = . }}
      {{- end }}
    {{- end }}
    <li>
      <a
        {{- range $k, $v := $attrs }}
          {{- with $v }}
            {{- printf " %s=%q" $k $v | safeHTMLAttr }}
          {{- end }}
        {{- end -}}
      >{{ $name }}</a>
      {{- with .Children }}
        <ul>
          {{- partial "inline/menu/walk.html" (dict "page" $page "menuEntries" .) }}
        </ul>
      {{- end }}
    </li>
  {{- end }}
{{- end }}
{{< /code >}}

Call the partial above, passing a menu ID and the current page in context.

{{< code file=layouts/_default/single.html >}}
{{ partial "menu.html" (dict "menuID" "main" "page" .) }}
{{ partial "menu.html" (dict "menuID" "footer" "page" .) }}
{{< /code >}}

## Page references

Regardless of how you [define menu entries], an entry associated with a page has access to page context.

This simplistic example renders a page parameter named `version` next to each entry's `name`. Code defensively using `with` or `if` to handle entries where (a) the entry points to an external resource, or (b) the `version` parameter is not defined.

{{< code file=layouts/_default/single.html >}}
{{- range site.Menus.main }}
  <a href="{{ .URL }}">
    {{ .Name }}
    {{- with .Page }}
      {{- with .Params.version -}}
        ({{ . }})
      {{- end }}
    {{- end }}
  </a>
{{- end }}
{{< /code >}}

## Menu entry parameters

When you define menu entries [in site configuration] or [in front matter], you can include a `params` key as shown in these examples:

- [Menu entry defined in site configuration]
- [Menu entry defined in front matter]

This simplistic example renders a `class` attribute for each anchor element. Code defensively using `with` or `if` to handle entries where `params.class` is not defined.

{{< code file=layouts/partials/menu.html >}}
{{- range site.Menus.main }}
  <a {{ with .Params.class -}} class="{{ . }}" {{ end -}} href="{{ .URL }}">
    {{ .Name }}
  </a>
{{- end }}
{{< /code >}}

## Localize

Hugo provides two methods to localize your menu entries. See [multilingual].

[automatic]: /content-management/menus/#define-automatically
[define menu entries]: /content-management/menus/
[defining menu entries]: /content-management/menus/
[in front matter]: /content-management/menus/#define-in-front-matter
[in site configuration]: /content-management/menus/#define-in-site-configuration
[localize the menu entries]: /content-management/multilingual/#menus
[menu entry defined in front matter]: /content-management/menus/#example-front-matter
[menu entry defined in site configuration]: /content-management/menus/#example-site-configuration
[menu methods]: /methods/menu/
[multilingual]: /content-management/multilingual/#menus
