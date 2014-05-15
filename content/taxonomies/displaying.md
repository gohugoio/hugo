---
title:  "Rendering Taxonomies"
date: "2013-07-01"
linktitle: "Displaying"
aliases: ["/indexes/displaying/"]
weight: 20
menu:
  main:
    parent: 'taxonomy'
---

## Rendering index values assigned to this content

Within your content templates you may wish to display 
the indexes that that piece of content is assigned to.

Because we are leveraging the front matter system to 
define indexes for content, the indexes assigned to 
each content piece are located in the usual place 
(.Params.`plural`)

### Example

    <ul id="tags">
      {{ range .Params.tags }}
        <li><a href="tags/{{ . | urlize }}">{{ . }}</a> </li>
      {{ end }}
    </ul>

## Rendering a Site's Indexes

If you wish to display the list of all keys for an index you can find retrieve
them from the `.Site` variable which is available on every page.

This may take the form of a tag cloud, a menu or simply a list.

The following example displays all tag keys:

### Example

    <ul id="all-tags">
      {{ range .Site.Indexes.tags }}
        <li><a href="/tags/{{ .Name | urlize }}">{{ .Name }}</a></li>  
      {{ end }}
    </ul>

## Creating a menu based on indexes

Hugo can generate menus based on indexes by iterating and
nesting the index keys. This can be used to build a hierarchy
of content within your site.

To have hugo create the menu, simply create a template in chrome
called menu.html, then include it using the 
`{{ template "chrome/menu.html" . }}` syntax.



### Example complete menu.html file
This example will list all indexes, each of their keys and all the content assigned to each key.

    <section id="menu">
      <ul>
        {{ range $indexname, $index := .Site.Indexes }}
          <li><a href="/{{ $indexname | urlize }}">{{ $indexname }}</a> 
            <ul> 
              {{ range $key, $value := $index }}
              <li> {{ $key }} </li>
                    <ul>
                    {{ range $value.Pages }}
                        <li hugo-nav="{{ .RelPermalink}}"><a href="{{ .Permalink}}"> {{ .LinkTitle }} </a> </li>
                    {{ end }}
                    </ul>
              {{ end }}
            </ul>
          </li> 
        {{ end }}
      </ul>
    </section>

### menu.html using a single index
It is more likely that you would want to use a single index for navigation.
In this example we are using the `groups` index for our menu.

    <section id="menu">
        <ul>
            {{ range $key, $index := .Site.Indexes.groups }}
            <li> {{ $key }} </li>
            <ul>
                {{ range $index.Pages }}
                <li hugo-nav="{{ .RelPermalink}}"><a href="{{ .Permalink}}"> {{ .LinkTitle }} </a> </li>
                {{ end }}
            </ul>
            {{ end }}
        </ul>
    </section>


### menu.html using a single index ordered by Popularity

    <section id="menu">
        <ul>
            {{ range .Site.Indexes.groups.ByCount }}
            <li> {{ .Name }} </li>
            <ul>
                {{ range .Pages }}
                <li hugo-nav="{{ .RelPermalink}}"><a href="{{ .Permalink}}"> {{ .LinkTitle }} </a> </li>
                {{ end }}
            </ul>
            {{ end }}
        </ul>
    </section>
