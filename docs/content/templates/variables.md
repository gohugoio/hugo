---
title: "Template Variables"
linktitle: "Variables"
date: "2013-07-01"
aliases: ["/doc/variables/","/layout/variables/"]
weight: 20
menu:
  main:
    parent: 'layout'
prev: "/templates/functions"
next: "/templates/content"
---

Hugo makes a set of values available to the templates. Go templates are context based. The following
are available in the context for the templates.

## Page Variables

The following is a list of most of the accessible variables which can be
defined for a piece of content. Many of these will be defined in the front
matter, content or derived from file location.

**.Title**  The title for the content.<br>
**.Content** The content itself, defined below the front matter.<br>
**.Summary** A generated summary of the content for easily showing a snippet in a summary view.<br>
**.Description** The description for the content.<br>
**.Keywords** The meta keywords for this content.<br>
**.Date** The date the content is associated with.<br>
**.PublishDate** The date the content is published on.<br>
**.Type** The content [type](/content/types/) (eg. post)<br>
**.Section** The [section](/content/sections/) this content belongs to<br>
**.Permalink** The Permanent link for this page.<br>
**.RelPermalink** The Relative permanent link for this page.<br>
**.LinkTitle** Access when creating links to this content. Will use linktitle if set in front-matter, else title<br>
**.Indexes** These will use the field name of the plural form of the index (see tags and categories above)<br>
**.RSSLink** Link to the indexes' rss link <br>
**.TableOfContents** The rendered table of contents for this content<br>
**.Prev** Pointer to the previous content (based on pub date)<br>
**.Next** Pointer to the following content (based on pub date)<br>
**.FuzzyWordCount** The approximate number of words in the content.<br>
**.WordCount** The number of words in the content.<br>
**.ReadingTime** The estimated time it takes to read the content in minutes.<br>
**.Weight** Assigned weight (in the front matter) to this content, used in sorting.<br>
**.Site** See site variables below<br>

## Page Params

Any other value defined in the front matter, including indexes will be made available under `.Params`.
Take for example I'm using tags and categories as my indexes. The following would be how I would access them:

**.Params.tags** <br>
**.Params.categories** <br>
<br>
**All Params are only accessible using all lowercase characters**<br>

## Node Variables
In Hugo a node is any page not rendered directly by a content file. This
includes indexes, lists and the homepage.

**.Title**  The title for the content.<br>
**.Date** The date the content is published on.<br>
**.Permalink** The Permanent link for this node<br>
**.Url** The relative url for this node.<br>
**.RSSLink** Link to the indexes' rss link <br>
**.Data** The data specific to this type of node.<br>
**.Site** See site variables below<br>

## Site Variables

Also available is `.Site` which has the following:

**.Site.BaseUrl** The base URL for the site as defined in the config.json file.<br>
**.Site.Indexes** The indexes for the entire site.<br>
**.Site.LastChange** The date of the last change of the most recent content.<br>
**.Site.Recent** Array of all content ordered by Date, newest first.<br>
**.Site.Params** A container holding the values from `params` in your site configuration file.<br>
