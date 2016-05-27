---
aliases:
- /doc/variables/
- /layout/variables/
lastmod: 2015-12-08
date: 2013-07-01
linktitle: Variables
menu:
  main:
    parent: layout
next: /templates/content
prev: /templates/functions
title: Template Variables
weight: 20
toc: true
---

Hugo makes a set of values available to the templates. Go templates are context based. The following
are available in the context for the templates.

## Page Variables

The following is a list of most of the accessible variables which can be
defined for a piece of content. Many of these will be defined in the front
matter, content or derived from file location.

**See also:** [Scratch](/extras/scratch) for page-scoped writable variables.

**.Title**  The title for the content.<br>
**.Content** The content itself, defined below the front matter.<br>
**.Summary** A generated summary of the content for easily showing a snippet in a summary view. Note that the breakpoint can be set manually by inserting <code>&lt;!&#x2d;&#x2d;more&#x2d;&#x2d;&gt;</code> at the appropriate place in the content page.  See [Summaries](/content/summaries/) for more details.<br>
**.Truncated** A boolean, `true` if the `.Summary` is truncated.  Useful for showing a "Read more..." link only if necessary.  See [Summaries](/content/summaries/) for more details.<br>
**.Description** The description for the content.<br>
**.Keywords** The meta keywords for this content.<br>
**.Date** The date the content is associated with.<br>
**.Lastmod** The date the content was last modified.<br>
**.PublishDate** The date the content is published on.<br>
**.Type** The content [type](/content/types/) (e.g. post).<br>
**.Section** The [section](/content/sections/) this content belongs to.<br>
**.Permalink** The Permanent link for this page.<br>
**.RelPermalink** The Relative permanent link for this page.<br>
**.LinkTitle** Access when creating links to this content. Will use `linktitle` if set in front matter, else `title`.<br>
**.RSSLink** Link to the taxonomies' RSS link.<br>
**.TableOfContents** The rendered table of contents for this content.<br>
**.Prev** Pointer to the previous content (based on pub date).<br>
**.Next** Pointer to the following content (based on pub date).<br>
**.PrevInSection** Pointer to the previous content within the same section (based on pub date). For example, `{{if .PrevInSection}}{{.PrevInSection.Permalink}}{{end}}`.<br>
**.NextInSection** Pointer to the following content within the same section (based on pub date)<br>
**.FuzzyWordCount** The approximate number of words in the content.<br>
**.WordCount** The number of words in the content.<br>
**.ReadingTime** The estimated time it takes to read the content in minutes.<br>
**.Weight** Assigned weight (in the front matter) to this content, used in sorting.<br>
**.RawContent** Raw Markdown content without the metadata header. Useful with [remarkjs.com](http://remarkjs.com)<br>
**.UniqueID** The MD5-checksum of the page's filename
**.Draft** A boolean, `true` if the content is marked as a draft in the front matter.<br>
**.IsNode** Always false for pages.<br>
**.IsPage** Always true for page.<br>
**.Site** See [Site Variables]({{< relref "#site-variables" >}}) below.<br>
**.Hugo** See [Hugo Variables]({{< relref "#hugo-variables" >}}) below.<br>

## Page Params

Any other value defined in the front matter, including taxonomies, will be made available under `.Params`.
For example, the *tags* and *categories* taxonomies are accessed with:

* **.Params.tags**
* **.Params.categories**

**All Params are only accessible using all lowercase characters.**

This is particularly useful for the introduction of user defined fields in content files. For example, a Hugo website on book reviews could have in the front matter of <code>/content/review/book01.md</code>

    ---
	...
    affiliatelink: "http://www.my-book-link.here"
	recommendedby: "my Mother"
    ---

Which would then be accessible to a template at `/themes/yourtheme/layouts/review/single.html` through `.Params.affiliatelink` and `.Params.recommendedby`, respectively. Two common situations where these could be introduced are as a value of a certain attribute (like `href=""` below) or by itself to be displayed. Sample syntaxes include:

    <h3><a href={{ printf "%s" $.Params.affiliatelink }}>Buy this book</a></h3>
	<p>It was recommended by {{ .Params.recommendedby }}.</p>

which would render

    <h3><a href="http://www.my-book-link.here">Buy this book</a></h3>
	<p>It was recommended by my Mother.</p>

**See also:** [Archetypes]({{% ref "content/archetypes.md" %}}) for consistency of `Params` across pieces of content.

### Param method
In Hugo you can declare params both for the site and the individual page.  A common use case is to have a general value for the site and a more specific value for some of the pages (i.e. an image).

With the `Param` method the most specific value will be selected for you, and it is safe to use it in any template (it's defined on both Page and Node):

```
$.Param "image"
```

## Node Variables
In Hugo, a node is any page not rendered directly by a content file. This
includes taxonomies, lists and the homepage.

**See also:** [Scratch](/extras/scratch) for global node variables.

**.Title**  The title for the content.<br>
**.Date** The date the content is published on.<br>
**.Lastmod** The date the content was last modified.<br>
**.Permalink** The Permanent link for this node<br>
**.URL** The relative URL for this node.<br>
**.Ref(ref)** Returns the permalink for `ref`. See [cross-references]({{% ref "extras/crossreferences.md" %}}). Does not handle in-page fragments correctly.<br>
**.RelRef(ref)** Returns the relative permalink for `ref`. See [cross-references]({{% ref "extras/crossreferences.md" %}}). Does not handle in-page fragments correctly.<br>
**.RSSLink** Link to the taxonomies' RSS link.<br>
**.Data** The data specific to this type of node.<br>
**.IsHome** True if the node is the site home page.<br>
**.IsNode** Always true for nodes.<br>
**.IsPage** Always false for nodes.<br>
**.Site** See [Site Variables]({{< relref "#site-variables" >}}) below.<br>
**.Hugo** See [Hugo Variables]({{< relref "#hugo-variables" >}}) below.<br>

### Taxonomy Term Variables

[Taxonomy Terms](/templates/terms/) pages are of the type "node" and have the following additional variables.

**.Data.Singular** The singular name of the taxonomy<br>
**.Data.Plural** The plural name of the taxonomy<br>
**.Data.Pages** the list of pages in this taxonomy<br>
**.Data.Terms** The taxonomy itself<br>
**.Data.Terms.Alphabetical** The Terms alphabetized<br>
**.Data.Terms.ByCount** The Terms ordered by popularity<br>

The last two can also be reversed: **.Data.Terms.Alphabetical.Reverse**, **.Data.Terms.ByCount.Reverse**.

## Site Variables

Also available is `.Site` which has the following:

**.Site.BaseURL** The base URL for the site as defined in the site configuration file.<br>
**.Site.RSSLink** The URL for the site RSS.<br>
**.Site.Taxonomies** The [taxonomies](/taxonomies/usage/) for the entire site.  Replaces the now-obsolete `.Site.Indexes` since v0.11.<br>
**.Site.Pages** Array of all content ordered by Date, newest first.  Replaces the now-deprecated `.Site.Recent` starting v0.13.<br>
**.Site.Params** A container holding the values from the `params` section of your site configuration file. For example, a TOML config file might look like this:

    baseurl = "http://yoursite.example.com/"

    [params]
      description = "Tesla's Awesome Hugo Site"
      author = "Nikola Tesla"
**.Site.Sections** Top level directories of the site.<br>
**.Site.Files** All of the source files of the site.<br>
**.Site.Menus** All of the menus in the site.<br>
**.Site.Title** A string representing the title of the site.<br>
**.Site.Author** A map of the authors as defined in the site configuration.<br>
**.Site.LanguageCode** A string representing the language as defined in the site configuration.<br>
**.Site.DisqusShortname** A string representing the shortname of the Disqus shortcode as defined in the site configuration.<br>
**.Site.GoogleAnalytics** A string representing your tracking code for Google Analytics as defined in the site configuration.<br>
**.Site.Copyright** A string representing the copyright of your web site as defined in the site configuration.<br>
**.Site.LastChange** A string representing the date/time of the most recent change to your site, based on the [`date` variable]({{< ref "content/front-matter.md#required-variables" >}}) in the front matter of your content pages.<br>
**.Site.Permalinks** A string to override the default permalink format. Defined in the site configuration.<br>
**.Site.BuildDrafts** A boolean (Default: false) to indicate whether to build drafts. Defined in the site configuration.<br>
**.Site.Data**  Custom data, see [Data Files](/extras/datafiles/).<br>

## Hugo Variables

Also available is `.Hugo` which has the following:

**.Hugo.Generator** Meta tag for the version of Hugo that generated the site. Highly recommended to be included by default in all theme headers so we can start to track the usage and popularity of Hugo. Unlike other variables it outputs a **complete** HTML tag, e.g. `<meta name="generator" content="Hugo 0.15" />`<br>
**.Hugo.Version** The current version of the Hugo binary you are using e.g. `0.13-DEV`<br>
**.Hugo.CommitHash** The git commit hash of the current Hugo binary e.g. `0e8bed9ccffba0df554728b46c5bbf6d78ae5247`<br>
**.Hugo.BuildDate** The compile date of the current Hugo binary formatted with RFC 3339 e.g. `2002-10-02T10:00:00-05:00`<br>
