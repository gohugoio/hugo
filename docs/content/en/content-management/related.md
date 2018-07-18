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

{{% note %}}
We currently do not index **Page content**. We thought we would release something that will make most people happy before we start solving [Sherlock's last case](https://github.com/joearms/sherlock).
{{% /note %}}

## List Related Content

To list up to 5 related pages is as simple as including something similar to this partial in your single page template:

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


{{% note %}}
Read [this blog article](https://regisphilibert.com/blog/2018/04/hugo-optmized-relashionships-with-related-content/) for a great explanation of more advanced usage of this feature.
{{% /note %}}

The full set of methods available on the page lists can be seen in this Go interface:

```go
// A PageGenealogist finds related pages in a page collection. This interface is implemented
// by Pages and PageGroup, which makes it available as `{{ .RegularPages.Related . }}` etc.
type PageGenealogist interface {

	// Template example:
	// {{ $related := .RegularPages.Related . }}
	Related(doc related.Document) (Pages, error)

	// Template example:
	// {{ $related := .RegularPages.RelatedIndices . "tags" "date" }}
	RelatedIndices(doc related.Document, indices ...interface{}) (Pages, error)

	// Template example:
	// {{ $related := .RegularPages.RelatedTo ( keyVals "tags" "hugo" "rocks")  ( keyVals "date" .Date ) }}
	RelatedTo(args ...types.KeyValues) (Pages, error)
}
```
## Configure Related Content
Hugo provides a sensible default configuration of Related Content, but you can fine-tune this in your configuration, on the global or language level if needed.

{{% note %}}
If you add a `related` config section, you need to add a complete configuration. It is not possible to just set, say, `includeNewer` and use the rest  from the Hugo defaults.
{{% /note %}}

Below is a sample `config.toml` section:

```
[related]

# Only include matches with rank >= threshold. This is a normalized rank between 0 and 100.
threshold = 80

# To get stable "See also" sections we, by default, exclude newer related pages.
includeNewer = false

# Will lower case keywords in both queries and in the indexes.
toLower = false

[[related.indices]]
name = "keywords"
weight = 150
[[related.indices]]
name  = "author"
toLower = true
weight = 30
[[related.indices]]
name  = "tags"
weight = 100
[[related.indices]]
name  = "date"
weight = 10
pattern = "2006"
```
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







