---
aliases:
- /indexes/displaying/
date: 2013-07-01
linktitle: Displaying
menu:
  main:
    parent: taxonomy
next: /taxonomies/templates
prev: /taxonomies/usage
title: Displaying Taxonomies
weight: 20
---

There are four common ways you can display the data in your
taxonomies in addition to the automatic taxonomy pages created by hugo
using the [list templates](/templates/list/):

1. For a given piece of content, you can list the terms attached
2. For a given piece of content, you can list other content with the same
   term
3. You can list all terms for a taxonomy
4. You can list all taxonomies (with their terms)

## 1. Displaying taxonomy terms assigned to this content

Within your content templates, you may wish to display
the taxonomies that that piece of content is assigned to.

Because we are leveraging the front matter system to
define taxonomies for content, the taxonomies assigned to
each content piece are located in the usual place
(.Params.`plural`).

### Example

    <ul id="tags">
      {{ range .Params.tags }}
        <li><a href="tags/{{ . | urlize }}">{{ . }}</a> </li>
      {{ end }}
    </ul>

## 2. Listing content with the same taxonomy term

First, you may be asking why you would use this. If you are using a
taxonomy for something like a series of posts, this is exactly how you
would do it. It’s also an quick and dirty way to show some related
content.


### Example

    <ul>
      {{ range .Site.Taxonomies.series.golang }}
        <li><a href="{{ .Url }}">{{ .Name }}</a></li>
      {{ end }}
    </ul>

## 3. Listing all content in a given taxonomy

This would be very useful in a sidebar as “featured content”. You could
even have different sections of “featured content” by assigning
different terms to the content.

### Example

    <section id="menu">
        <ul>
            {{ range $key, $taxonomy := .Site.Taxonomies.featured }}
            <li> {{ $key }} </li>
            <ul>
                {{ range $taxonomy.Pages }}
                <li hugo-nav="{{ .RelPermalink}}"><a href="{{ .Permalink}}"> {{ .LinkTitle }} </a> </li>
                {{ end }}
            </ul>
            {{ end }}
        </ul>
    </section>


## 4. Rendering a Site's Taxonomies

If you wish to display the list of all keys for a taxonomy, you can find retrieve
them from the `.Site` variable which is available on every page.

This may take the form of a tag cloud, a menu or simply a list.

The following example displays all tag keys:

### Example

    <ul id="all-tags">
      {{ range $name, $taxonomy := .Site.Taxonomies.tags }}
        <li><a href="/tags/{{ $name | urlize }}">{{ $name }}</a></li>
      {{ end }}
    </ul>

### Complete Example
This example will list all taxonomies, each of their keys and all the content assigned to each key.

    <section>
      <ul>
        {{ range $taxonomyname, $taxonomy := .Site.Taxonomies }}
          <li><a href="/{{ $taxonomyname | urlize }}">{{ $taxonomyname }}</a>
            <ul>
              {{ range $key, $value := $taxonomy }}
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

