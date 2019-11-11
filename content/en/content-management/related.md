---
title: Related Content
description: List related content in "See Also" sections.
date: 2017-09-05
categories: [content management]
keywords: [content]
menu:
  docs:
    parent: "content-management"
    weight: 40
weight: 30
draft: false
aliases: [/content/related/,/related/]
toc: true
---


Hugo uses a set of factors to identify a page's related content based on Front Matter parameters. This can be tuned to the desired set of indices and parameters or left to Hugo's default [Related Content configuration](#configure-related-content).

## List Related Content


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

### Methods

Here is the list of "Related" methods available on a page collection such `.RegularPages`.

#### .Related PAGE
Returns a collection of pages related the given one.

```
{{ $related := .Site.RegularPages.Related . }}
```

#### .RelatedIndices PAGE INDICE1 [INDICE2 ...]
Returns a collection of pages related to a given one restricted to a list of indices.

```
{{ $related := .Site.RegularPages.RelatedIndices . "tags" "date" }}
```

#### .RelatedTo KEYVALS [KEYVALS2 ...]
Returns a collection of pages related together by a set of indices and their match.

In order to build those set and pass them as argument, one must use the `keyVals` function where the first argument would be the `indice` and the consecutive ones its potential `matches`.

```
{{ $related := .Site.RegularPages.RelatedTo ( keyVals "tags" "hugo" "rocks")  ( keyVals "date" .Date ) }}
```

{{% note %}}
Read [this blog article](https://regisphilibert.com/blog/2018/04/hugo-optmized-relashionships-with-related-content/) for a great explanation of more advanced usage of this feature.
{{% /note %}}

## Configure Related Content
Hugo provides a sensible default configuration of Related Content, but you can fine-tune this in your configuration, on the global or language level if needed.

### Default configuration

Without any `related` configuration set on the project, Hugo's Related Content methods will use the following.

```yaml
related:
  threshold: 80
  includeNewer: false
  toLower: false
  indices:
  - name: keywords
    weight: 100
  - name: date
    weight: 10
```

Custom configuration should be set using the same syntax.

{{% note %}}
If you add a `related` config section, you need to add a complete configuration. It is not possible to just set, say, `includeNewer` and use the rest  from the Hugo defaults.
{{% /note %}}

### Top Level Config Options

threshold
:  A value between 0-100. Lower value will give more, but maybe not so relevant, matches.

includeNewer
:  Set to true to include **pages newer than the current page** in the related content listing. This will mean that the output for older posts may change as new related content gets added.

toLower
: Set to true to lower case keywords in both the indexes and the queries. This may give more accurate results at a slight performance penalty. Note that this can also be set per index.

### Config Options per Index

name
:  The index name. This value maps directly to a page param. Hugo supports string values (`author` in the example) and lists (`tags`, `keywords` etc.) and time and date objects.

weight
: An integer weight that indicates _how important_ this parameter is relative to the other parameters.  It can be 0, which has the effect of turning this index off, or even negative. Test with different values to see what fits your content best.

pattern
: This is currently only relevant for dates. When listing related content, we may want to list content that is also close in time. Setting "2006" (default value for date indexes) as the pattern for a date index will add weight to pages published in the same year. For busier blogs, "200601" (year and month) may be a better default.

toLower
: See above.

## Performance Considerations

**Fast is Hugo's middle name** and we would not have released this feature had it not been blistering fast.

This feature has been in the back log and requested by many for a long time. The development got this recent kick start from this Twitter thread:

{{< tweet 898398437527363585 >}}

Scott S. Lowe removed the "Related Content" section built using the `intersect` template function on tags, and the build time dropped from 30 seconds to less than 2 seconds on his 1700 content page sized blog.

He should now be able to add an improved version of that "Related Content" section without giving up the fast live-reloads. But it's worth noting that:

* If you don't use any of the `Related` methods, you will not use the Relate Content feature, and performance will be the same as before.
* Calling `.RegularPages.Related` etc. will create one inverted index, also sometimes named posting list, that will be reused for any lookups in that same page collection. Doing that in addition to, as an example, calling `.Pages.Related` will work as expected, but will create one additional inverted index. This should still be very fast, but worth having in mind, especially for bigger sites.

{{% note %}}
We currently do not index **Page content**. We thought we would release something that will make most people happy before we start solving [Sherlock's last case](https://github.com/joearms/sherlock).
{{% /note %}}
