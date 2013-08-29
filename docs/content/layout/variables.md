---
title: "Variables"
date: "2013-07-01"
aliases: ["/doc/variables/"]
---

Hugo makes a set of values available to the templates. Go templates are context based. The following
are available in the context for the templates.

## Page Variables

**.Title**  The title for the content.<br>
**.Description** The description for the content.<br>
**.Keywords** The meta keywords for this content.<br>
**.Date** The date the content is published on.<br>
**.Indexes** These will use the field name of the plural form of the index (see tags and categories above)<br>
**.Permalink** The Permanent link for this page.<br>
**.FuzzyWordCount** The approximate number of words in the content.<br>
**.RSSLink** Link to the indexes' rss link <br>
**.Prev** Pointer to the previous content (based on pub date)<br>
**.Next** Pointer to the following content (based on pub date)<br>
**.Site** See site variables below<br>
**.Content** The content itself, defined below the front matter.<br>
**.Summary** A generated summary of the content for easily showing a snippet in a summary view.<br>

Any value defined in the front matter, including indexes will be made available under `.Params`.
Take for example I'm using tags and categories as my indexes. The following would be how I would access them:

**.Params.Tags** <br>
**.Params.Categories** <br>

## Node Variables
In Hugo a node is any page not rendered directly by a content file. This
includes indexes, lists and the homepage.

**.Title**  The title for the content.<br>
**.Date** The date the content is published on.<br>
**.Data** The data specific to this type of node.<br>
**.Permalink** The Permanent link for this node<br>
**.Url** The relative url for this node.<br>
**.RSSLink** Link to the indexes' rss link <br>
**.Site** See site variables below<br>

## Site Variables

Also available is `.Site` which has the following:

**.Site.BaseUrl** The base URL for the site as defined in the config.json file.<br>
**.Site.Indexes** The names of the indexes of the site.<br>
**.Site.LastChange** The date of the last change of the most recent content.<br>
**.Site.Recent** Array of all content ordered by Date, newest first<br>

