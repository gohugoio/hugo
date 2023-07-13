---
title: Related content
description: List related content in "See Also" sections.
categories: [content management]
keywords: [content]
menu:
  docs:
    parent: content-management
    weight: 110
toc: true
weight: 110
aliases: [/content/related/,/related/]
---

Hugo uses a set of factors to identify a page's related content based on front matter parameters. This can be tuned to the desired set of indices and parameters or left to Hugo's default [Related Content configuration](#configure-related-content).

## List related content

To list up to 5 related pages (which share the same _date_ or _keyword_ parameters) is as simple as including something similar to this partial in your single page template:

{{< code file="layouts/partials/related.html" >}}
{{ $related := .Site.RegularPages.Related . | first 5 }}
{{ with $related }}
<h3>See Also</h3>
<ul>
 {{ range . }}
 <li><a href="{{ .RelPermalink }}">{{ .Title }}</a></li>
 {{ end }}
</ul>
{{ end }}
{{< /code >}}

The `Related` method takes one argument which may be a `Page` or a options map. The options map have these options:

indices
: The indices to search in.

document
: The document to search for related content for.

namedSlices
: The keywords to search for.

fragments
: Fragments holds a a list of special keywords that is used for indices configured as type "fragments". This will match the fragment identifiers of the documents.

A fictional example using all of the above options:

```go-html-template
{{ $page := . }}
{{ $opts := 
  "indices" (slice "tags" "keywords")
  "document" $page
  "namedSlices" (slice (keyVals "tags" "hugo" "rocks") (keyVals "date" $page.Date))
  "fragments" (slice "heading-1" "heading-2")
}}
```

{{% note %}}
We improved and simplified this feature in Hugo 0.111.0. Before this we had 3 different methods: `Related`, `RelatedTo` and `RelatedIndicies`. Now we have only one method: `Related`. The old methods are still available but deprecated. Also see [this blog article](https://regisphilibert.com/blog/2018/04/hugo-optmized-relashionships-with-related-content/) for a great explanation of more advanced usage of this feature.
{{% /note %}}

## Index content headings in related content

{{< new-in "0.111.0" >}}

Hugo can index the headings in your content and use this to find related content. You can enable this by adding a index of type `fragments` to your `related` configuration:

{{< code-toggle file="hugo" copy=false >}}
[related]
threshold    = 20
includeNewer = true
toLower      = false
[[related.indices]]
name        = "fragmentrefs"
type        = "fragments"
applyFilter = false
weight      = 80
{{< /code-toggle >}}

* The `name` maps to a optional front matter slice attribute that can be used to link from the page level down to the fragment/heading level.
* If `applyFilter`is enabled, the `.HeadingsFiltered` on each page in the result will reflect the filtered headings. This is useful if you want to show the headings in the related content listing:

```go-html-template
{{ $related := .Site.RegularPages.Related . | first 5 }}
{{ with $related }}
  <h2>See Also</h2>
  <ul>
    {{ range $i, $p := . }}
      <li>
        <a href="{{ .RelPermalink }}">{{ .Title }}</a>
        {{ with .HeadingsFiltered }}
          <ul>
            {{ range . }}
              {{ $link := printf "%s#%s" $p.RelPermalink .ID | safeURL }}
              <li>
                <a href="{{ $link }}">{{ .Title }}</a>
              </li>
            {{ end }}
          </ul>
        {{ end }}
      </li>
    {{ end }}
  </ul>
{{ end }}
```

## Configure related content

Hugo provides a sensible default configuration of Related Content, but you can fine-tune this in your configuration, on the global or language level if needed.

### Default configuration

Without any `related` configuration set on the project, Hugo's Related Content methods will use the following.

{{< code-toggle file="hugo" >}}
related:
  threshold: 80
  includeNewer: false
  toLower: false
  indices:
  - name: keywords
    weight: 100
  - name: date
    weight: 10
{{< /code-toggle >}}

Note that if you have configured `tags` as a taxonomy, `tags` will also be added to the default configuration above with the weight of `80`.

Custom configuration should be set using the same syntax.

{{% note %}}
If you add a `related` configuration section, you need to add a complete configuration. It is not possible to just set, say, `includeNewer` and use the rest  from the Hugo defaults.
{{% /note %}}

### Top level configuration options

threshold
:  A value between 0-100. Lower value will give more, but maybe not so relevant, matches.

includeNewer
:  Set to true to include **pages newer than the current page** in the related content listing. This will mean that the output for older posts may change as new related content gets added.

toLower
: Set to true to lower case keywords in both the indexes and the queries. This may give more accurate results at a slight performance penalty. Note that this can also be set per index.

### Configuration options per index

name
:  The index name. This value maps directly to a page parameter. Hugo supports string values (`author` in the example) and lists (`tags`, `keywords` etc.) and time and date objects.

type
: {{< new-in "0.111.0" >}}. One of `basic`(default) or `fragments`.

applyFilter
: {{< new-in "0.111.0" >}}. Apply a `type` specific filter to the result of a search. This is currently only used for the `fragments` type.

weight
: An integer weight that indicates _how important_ this parameter is relative to the other parameters.  It can be 0, which has the effect of turning this index off, or even negative. Test with different values to see what fits your content best.


cardinalityThreshold (default 0)
: {{< new-in "0.111.0" >}}. A percentage (0-100) used to remove common keywords from the index. As an example, setting this to 50 will remove all keywords that are used in more than 50% of the documents in the index.

pattern
: This is currently only relevant for dates. When listing related content, we may want to list content that is also close in time. Setting "2006" (default value for date indexes) as the pattern for a date index will add weight to pages published in the same year. For busier blogs, "200601" (year and month) may be a better default.

toLower
: See above.

## Performance considerations

**Fast is Hugo's middle name** and we would not have released this feature had it not been blistering fast.

This feature has been in the back log and requested by many for a long time. The development got this recent kick start from this Twitter thread:

{{< tweet user="scott_lowe" id="898398437527363585" >}}

Scott S. Lowe removed the "Related Content" section built using the `intersect` template function on tags, and the build time dropped from 30 seconds to less than 2 seconds on his 1700 content page sized blog.

He should now be able to add an improved version of that "Related Content" section without giving up the fast live-reloads. But it's worth noting that:

* If you don't use any of the `Related` methods, you will not use the Relate Content feature, and performance will be the same as before.
* Calling `.RegularPages.Related` etc. will create one inverted index, also sometimes named posting list, that will be reused for any lookups in that same page collection. Doing that in addition to, as an example, calling `.Pages.Related` will work as expected, but will create one additional inverted index. This should still be very fast, but worth having in mind, especially for bigger sites.

{{% note %}}
We currently do not index **Page content**. We thought we would release something that will make most people happy before we start solving [Sherlock's last case](https://github.com/joearms/sherlock).
{{% /note %}}
