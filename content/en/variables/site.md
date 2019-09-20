---
title: Site Variables
linktitle: Site Variables
description: Many, but not all, site-wide variables are defined in your site's configuration. However, Hugo provides a number of built-in variables for convenient access to global values in your templates.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [variables and params]
keywords: [global,site]
draft: false
menu:
  docs:
    parent: "variables"
    weight: 10
weight: 10
sections_weight: 10
aliases: [/variables/site-variables/]
toc: true
---

The following is a list of site-level (aka "global") variables. Many of these variables are defined in your site's [configuration file][config], whereas others are built into Hugo's core for convenient usage in your templates.

## Site Variables List

.Site.AllPages
: array of all pages, regardless of their translation.

.Site.Author
: a map of the authors as defined in the site configuration.

.Site.BaseURL
: the base URL for the site as defined in the site configuration.

.Site.BuildDrafts
: a boolean (default: `false`) to indicate whether to build drafts as defined in the site configuration.

.Site.Copyright
: a string representing the copyright of your website as defined in the site configuration.

.Site.Data
: custom data, see [Data Templates](/templates/data-templates/).

.Site.DisqusShortname
: a string representing the shortname of the Disqus shortcode as defined in the site configuration.

.Site.GoogleAnalytics
: a string representing your tracking code for Google Analytics as defined in the site configuration.

.Site.Home
: reference to the homepage's [page object](https://gohugo.io/variables/page/)

.Site.IsMultiLingual
: whether there are more than one language in this site. See [Multilingual](/content-management/multilingual/) for more information.

.Site.IsServer
: a boolean to indicate if the site is being served with Hugo's built-in server. See [`hugo server`](/commands/hugo_server/) for more information.

.Site.Language.Lang
: the language code of the current locale (e.g., `en`).

.Site.Language.LanguageName
: the full language name (e.g. `English`).

.Site.Language.Weight
: the weight that defines the order in the `.Site.Languages` list.

.Site.Language
: indicates the language currently being used to render the website. This object's attributes are set in site configurations' language definition.

.Site.LanguageCode
: a string representing the language as defined in the site configuration. This is mostly used to populate the RSS feeds with the right language code.

.Site.LanguagePrefix
: this can be used to prefix URLs to point to the correct language. It will even work when only one defined language. See also the functions [absLangURL](/functions/abslangurl/) and [relLangURL](/functions/rellangurl).

.Site.Languages
: an ordered list (ordered by defined weight) of languages.

.Site.LastChange
: a string representing the date/time of the most recent change to your site. This string is based on the [`date` variable in the front matter](/content-management/front-matter) of your content pages.

.Site.Menus
: all of the menus in the site.

.Site.Pages
: array of all content ordered by Date with the newest first. This array contains only the pages in the current language. See [`.Site.Pages`](#site-pages).

.Site.RegularPages
: a shortcut to the *regular* page collection. `.Site.RegularPages` is equivalent to `where .Site.Pages "Kind" "page"`. See [`.Site.Pages`](#site-pages).

.Site.Sections
: top-level directories of the site.

.Site.Taxonomies
: the [taxonomies](/taxonomies/usage/) for the entire site.  Replaces the now-obsolete `.Site.Indexes` since v0.11. Also see section [Taxonomies elsewhere](#taxonomies-elsewhere).

.Site.Title
: a string representing the title of the site.

## The `.Site.Params` Variable

`.Site.Params` is a container holding the values from the `params` section of your site configuration.

### Example: `.Site.Params`

The following `config.[yaml|toml|json]` defines a site-wide param for `description`:

{{< code-toggle file="config" >}}
baseURL = "https://yoursite.example.com/"

[params]
  description = "Tesla's Awesome Hugo Site"
  author = "Nikola Tesla"
{{</ code-toggle >}}

You can use `.Site.Params` in a [partial template](/templates/partials/) to call the default site description:

{{< code file="layouts/partials/head.html" >}}
<meta name="description" content="{{if .IsHome}}{{ $.Site.Params.description }}{{else}}{{.Description}}{{end}}" />
{{< /code >}}

## The `.Site.Pages` Variable {#site-pages}

### `.Site.Pages` compared to `.Pages`

{{< readfile file="/content/en/readfiles/pages-vs-site-pages.md" markdown="true" >}}




[config]: /getting-started/configuration/
