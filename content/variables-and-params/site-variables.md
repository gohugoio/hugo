---
title: Site Variables
linktitle: Site Variables
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
weight:
categories: [variables and params]
tags: [global,site]
draft: false
slug:
aliases: []
toc: false
notesforauthors:
---

Also available is `.Site` which has the following:

* `.Site.AllPages` Array of all pages regardless of their translation.
* `.Site.Author` A map of the authors as defined in the site configuration.
* `.Site.BaseURL` The base URL for the site as defined in the site configuration file.
* `.Site.BuildDrafts` A boolean (Default: false) to indicate whether to build drafts. Defined in the site configuration.
* `.Site.Copyright` A string representing the copyright of your web site as defined in the site configuration.
* `.Site.Data`  Custom data, see [Data Files](/extras/datafiles/).
* `.Site.DisqusShortname` A string representing the shortname of the Disqus shortcode as defined in the site configuration.
* `.Site.Files` All of the source files of the site.
* `.Site.GoogleAnalytics` A string representing your tracking code for Google Analytics as defined in the site configuration.
* `.Site.IsMultiLingual` Whether there are more than one language in this site. See [Multilingual](/content-management/multilingual-mode/) for more information.
* `.Site.Language.Lang` The language code of the current locale, e.g. `en`.
* `.Site.Language.LanguageName` The full language name, e.g. `English`.
* `.Site.Language.Weight` The weight that defines the order in the `.Site.Languages` list.
* `.Site.Language` This indicates which language you are currently rendering the website for.  This is an object with the attributes set in your language definition in your site config.
* `.Site.LanguageCode` A string representing the language as defined in the site configuration. This is mostly used to populate the RSS feeds with the right language code.
* `.Site.LanguagePrefix` This can be used to prefix  theURLs with whats needed to point to the correct language. It will even work when only one language defined. See also the functions [absLangURL](/functions/abslangurl/) and [relLangURL](/functions/rellangurl).
* `.Site.Languages` An ordered list (ordered by defined weight) of languages.
* `.Site.LastChange` A string representing the date/time of the most recent change to your site, based on the [`date` variable in the front matter](/content-management/front-matter) of your content pages.
* `.Site.Menus` All of the menus in the site.
* `.Site.Pages` Array of all content ordered by Date, newest first.  Replaces the now-deprecated `.Site.Recent` starting v0.13. This array contains only the pages in the current language.
* `.Site.Permalinks` A string to override the default [permalink](/content-management/url-management/) format defined in the [site configuration](/getting-started/configuration/).
* `.Site.RegularPages` A shortcut to the *regular page* collection. Equivalent to `where .Site.Pages "Kind" "page"`.
* `.Site.RSSLink` The URL for the site RSS.
* `.Site.Sections` Top level directories of the site.
* `.Site.Taxonomies` The [taxonomies](/taxonomies/usage/) for the entire site.  Replaces the now-obsolete `.Site.Indexes` since v0.11. Also see section [Taxonomies elsewhere](#taxonomies-elsewhere).
* `.Site.Title` A string representing the title of the site.

`.Site.Params` is a container holding the values from the `params` section of your site configuration file. For example, a TOML config file might look like this:

```toml
baseURL = "http://yoursite.example.com/"

[params]
  description = "Tesla's Awesome Hugo Site"
  author = "Nikola Tesla"
```

