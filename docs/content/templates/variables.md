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



**.Content** The content itself, defined below the front matter.<br>
**.Data** The data specific to this type of page.<br>
**.Date** The date the page is associated with.<br>
**.Description** The description for the page.<br>
**.Draft** A boolean, `true` if the content is marked as a draft in the front matter.<br>
**.ExpiryDate** The date where the content is scheduled to expire on.<br>
**.FuzzyWordCount** The approximate number of words in the content.<br>
**.Hugo** See [Hugo Variables]({{< relref "#hugo-variables" >}}) below.<br>
**.IsHome** True if this is the home page.<br>
**.IsNode** Always false for regular content pages.<br>
**.IsPage** Always true for regular content pages.<br>
**.IsTranslated** Whether there are any translations to display.<br>
**.Keywords** The meta keywords for this content.<br>
**.Kind** What *kind* of page is this: is one of *page, home, section, taxonomy or taxonomyTerm.* There are also *RSS, sitemap, robotsTXT and 404*, but these will only available during rendering of that kind of page, and not available in any of the `Pages` collections.<br>
**.Lang** Language taken from the language extension notation.<br>
**.Language** A language object that points to this the language's definition in the site config.<br>
**.Lastmod** The date the content was last modified.<br>
**.LinkTitle** Access when creating links to this content. Will use `linktitle` if set in front matter, else `title`.<br>
**.Next** Pointer to the following content (based on pub date).<br>
**.NextInSection** Pointer to the following content within the same section (based on pub date)<br>
**.Pages** a collection of associated pages. This will be nil for regular content pages. This is an alias for **.Data.Pages**.<br>
**.Permalink** The Permanent link for this page.<br>
**.Prev** Pointer to the previous content (based on pub date).<br>
**.PrevInSection** Pointer to the previous content within the same section (based on pub date). For example, `{{if .PrevInSection}}{{.PrevInSection.Permalink}}{{end}}`.<br>
**.PublishDate** The date the content is published on.<br>
**.RSSLink** Link to the taxonomies' RSS link.<br>
**.RawContent** Raw Markdown content without the metadata header. Useful with [remarkjs.com](http://remarkjs.com)<br>
**.ReadingTime** The estimated time it takes to read the content in minutes.<br>
**.Ref** Returns the permalink for a given reference.  Example: `.Ref "sample.md"`. See [cross-references]({{% ref "extras/crossreferences.md" %}}). Does not handle in-page fragments correctly.<br>
**.RelPermalink** The Relative permanent link for this page.<br>
**.RelRef** Returns the relative permalink for a given reference.  Example: `RelRef "sample.md"`. See [cross-references]({{% ref "extras/crossreferences.md" %}}). Does not handle in-page fragments.<br>
**.Section** The [section](/content/sections/) this content belongs to.<br>
**.Site** See [Site Variables]({{< relref "#site-variables" >}}) below.<br>
**.Summary** A generated summary of the content for easily showing a snippet in a summary view. Note that the breakpoint can be set manually by inserting <code>&lt;!&#x2d;&#x2d;more&#x2d;&#x2d;&gt;</code> at the appropriate place in the content page.  See [Summaries](/content/summaries/) for more details.<br>
**.TableOfContents** The rendered table of contents for this content.<br>
**.Title**  The title for this page.<br>
**.Translations** A list of translated versions of the current page. See [Multilingual]({{< relref "content/multilingual.md" >}}) for more info.<br>
**.Truncated** A boolean, `true` if the `.Summary` is truncated.  Useful for showing a "Read more..." link only if necessary.  See [Summaries](/content/summaries/) for more details.<br>
**.Type** The content [type](/content/types/) (e.g. post).<br>
**.URL** The relative URL for this page. Note that if `URL` is set directly in frontmatter, that URL is returned as-is.<br>
**.UniqueID** The MD5-checksum of the content file's path<br>
**.Weight** Assigned weight (in the front matter) to this content, used in sorting.<br>
**.WordCount** The number of words in the content.<br>

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

In Hugo you can declare params both for the site and the individual page. A
common use case is to have a general value for the site and a more specific
value for some of the pages (i.e. a header image):

```
{{ $.Param "header_image" }}
```

The `.Param` method provides a way to resolve a single value whether it's
in a page parameter or a site parameter.

When frontmatter contains nested fields, like:

```
---
author:
  given_name: John
  family_name: Feminella
  display_name: John Feminella
---
```

then `.Param` can access them by concatenating the field names together with a
dot:

```
{{ $.Param "author.display_name" }}
```

If your frontmatter contains a top-level key that is ambiguous with a nested
key, as in the following case,

```
---
favorites.flavor: vanilla
favorites:
  flavor: chocolate
---
```

then the top-level key will be preferred. In the previous example, this

```
{{ $.Param "favorites.flavor" }}
```

will print `vanilla`, not `chocolate`.

### Taxonomy Terms Page Variables

[Taxonomy Terms](/templates/terms/) pages are of the type `Page` and have the following additional variables. These are available in `layouts/_defaults/terms.html` for example.

**.Data.Singular** The singular name of the taxonomy<br>
**.Data.Plural** The plural name of the taxonomy<br>
**.Data.Pages** the list of pages in this taxonomy<br>
**.Data.Terms** The taxonomy itself<br>
**.Data.Terms.Alphabetical** The Terms alphabetized<br>
**.Data.Terms.ByCount** The Terms ordered by popularity<br>

The last two can also be reversed: **.Data.Terms.Alphabetical.Reverse**, **.Data.Terms.ByCount.Reverse**.

### Taxonomies elsewhere

The **.Site.Taxonomies** variable holds all taxonomies defines site-wide.  It is a map of the taxonomy name to a list of its values. For example: "tags" -> ["tag1", "tag2", "tag3"]. Each value, though, is not a string but rather a [Taxonomy variable](#the-taxonomy-variable).

#### The Taxonomy variable

The Taxonomy variable, available as **.Site.Taxonomies.tags** for example, contains the list of tags (values) and, for each of those, their corresponding content pages.

## Site Variables

Also available is `.Site` which has the following:

**.Site.BaseURL** The base URL for the site as defined in the site configuration file.<br>
**.Site.RSSLink** The URL for the site RSS.<br>
**.Site.Taxonomies** The [taxonomies](/taxonomies/usage/) for the entire site.  Replaces the now-obsolete `.Site.Indexes` since v0.11. Also see section [Taxonomies elsewhere](#taxonomies-elsewhere).<br>
**.Site.Pages** Array of all content ordered by Date, newest first. `.Site.Pages` replaced `.Site.Recent`, which is no longer supported. This array contains only the pages in the current language.<br>
**.Site.AllPages** Array of all pages regardless of their translation.<br>
**.Site.Params** A container holding the values from the `params` section of your site configuration file. For example, a TOML config file might look like this:

    baseURL = "http://yoursite.example.com/"

    [params]
      description = "Tesla's Awesome Hugo Site"
      author = "Nikola Tesla"
**.Site.Sections** Top level directories of the site.<br>
**.Site.Files** All of the source files of the site.<br>
**.Site.Menus** All of the menus in the site.<br>
**.Site.Title** A string representing the title of the site.<br>
**.Site.Author** A map of the authors as defined in the site configuration.<br>
**.Site.LanguageCode** A string representing the language as defined in the site configuration. This is mostly used to populate the RSS feeds with the right language code.<br>
**.Site.DisqusShortname** A string representing the shortname of the Disqus shortcode as defined in the site configuration.<br>
**.Site.GoogleAnalytics** A string representing your tracking code for Google Analytics as defined in the site configuration.<br>
**.Site.Copyright** A string representing the copyright of your web site as defined in the site configuration.<br>
**.Site.LastChange** A string representing the date/time of the most recent change to your site, based on the [`date` variable]({{< ref "content/front-matter.md#required-variables" >}}) in the front matter of your content pages.<br>
**.Site.Permalinks** A string to override the default permalink format. Defined in the site configuration.<br>
**.Site.BuildDrafts** A boolean (Default: false) to indicate whether to build drafts. Defined in the site configuration.<br>
**.Site.Data**  Custom data, see [Data Files](/extras/datafiles/).<br>
**.Site.IsMultiLingual** Whether there are more than one language in this site.<br> See [Multilingual]({{< relref "content/multilingual.md" >}}) for more info.<br>
**.Site.Language** This indicates which language you are currently rendering the website for.  This is an object with the attributes set in your language definition in your site config.<br>
**.Site.Language.Lang** The language code of the current locale, e.g. `en`.<br>
**.Site.Language.Weight** The weight that defines the order in the `.Site.Languages` list.<br>
**.Site.Language.LanguageName** The full language name, e.g. `English`.<br>
**.Site.LanguagePrefix** This can be used to prefix  theURLs with whats needed to point to the correct language. It will even work when only one language defined. See also the functions [absLangURL and relLangURL]({{< relref "templates/functions.md#abslangurl-rellangurl" >}}).<br>
**.Site.Languages** An ordered list (ordered by defined weight) of languages.<br>
**.Site.RegularPages** A shortcut to the *regular page* collection. Equivalent to `where .Site.Pages "Kind" "page"`.<br>

## File Variables

The `.File` variable gives you additional information of a page.

> **Note:** `.File` is only accessible on *Pages* that has a content page attached to it.

Available are the following attributes:

**.File.Path** The original relative path of the page, e.g. `content/posts/foo.en.md`<br>
**.File.LogicalName** The name of the content file that represents a page, e.g. `foo.en.md`<br>
**.File.TranslationBaseName** The filename without extension or optional language identifier, e.g. `foo`<br>
**.File.BaseFileName** The filename without extension, e.g. `foo.en`<br>
**.File.Ext** or **.File.Extension** The file extension of the content file, e.g. `md`<br>
**.File.Lang** The language associated with the given file if [Multilingual]({{< relref "content/multilingual.md" >}}) is enabled, e.g. `en`<br>
**.File.Dir** Given the path `content/posts/dir1/dir2/`, the relative directory path of the content file will be returned, e.g. `posts/dir1/dir2/`<br>

## Hugo Variables

Also available is `.Hugo` which has the following:

**.Hugo.Generator** Meta tag for the version of Hugo that generated the site. Highly recommended to be included by default in all theme headers so we can start to track the usage and popularity of Hugo. Unlike other variables it outputs a **complete** HTML tag, e.g. `<meta name="generator" content="Hugo 0.15" />`<br>
**.Hugo.Version** The current version of the Hugo binary you are using e.g. `0.13-DEV`<br>
**.Hugo.CommitHash** The git commit hash of the current Hugo binary e.g. `0e8bed9ccffba0df554728b46c5bbf6d78ae5247`<br>
**.Hugo.BuildDate** The compile date of the current Hugo binary formatted with RFC 3339 e.g. `2002-10-02T10:00:00-05:00`<br>
