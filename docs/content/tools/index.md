---
date: 2015-09-12T10:40:31+02:00
title: Tools
weight: 120
---

This section highlights some projects around Hugo that are independently developed.
Those tools try to extend the functionality of our static site generator or help you to get started.


## Migration tools

Take a look a this list of migration tools if you currently use other blogging tools
like Jekyll or Wordpress but you intend to switch to Hugo instead. They'll take care to export
your content into Hugo-friendly formats.

### Jekyll

Alternatively you can follow the manual [migration guide]({{< relref "tutorials/migrate-from-jekyll.md" >}}).

- [JekyllToHugo](https://github.com/SenjinDarashiva/JekyllToHugo) - A Small script for converting Jekyll blog posts to a Hugo site.
- [ConvertToHugo](https://github.com/coderzh/ConvertToHugo) - Convert your blog from Jekyll to Hugo.


### Octopress

- [octohug](https://github.com/codebrane/octohug) - Octopress to Hugo migrator.

### Wordpress

- [wordpress-to-hugo-exporter](https://github.com/SchumacherFM/wordpress-to-hugo-exporter) - A one-click WordPress plugin that converts all posts, pages, taxonomies, metadata, and settings to Markdown and YAML which can be dropped into Hugo.

### Drupal

- [drupal2hugo](https://github.com/danapsimer/drupal2hugo) - Convert a Drupal site to Hugo.


### Blogger

- [blogimport](https://github.com/natefinch/blogimport) - A tool to import from Blogger posts to Hugo.

***

## Frontends

Do you prefer an graphical user interface over a text editor? Then give this frontends a try.

- [rango](https://github.com/stayradiated/rango) - A web frontend for Hugo. It's designed to make it easy to manage a small site, even for people with little computer experience.
- [enwrite](https://github.com/zzamboni/enwrite) - Evernote-powered statically-generated blogs and websites. Now posting to your blog or updating your website is as easy as writing a new note in Evernote!

***

## Search

A static site with a dynamic search function? Yes. Alternatively to embeddable scripts from Google or other search engines you can provide your visitors a custom search by indexing your content files directly.

- [Hugoidx](https://github.com/blevesearch/hugoidx) is an experimental application to create a search index. It's build on top of [Bleve](http://www.blevesearch.com/).
- This [Github Gist](https://gist.github.com/sebz/efddfc8fdcb6b480f567) contains simple workflow to create a search index for your static site. It uses a simple Grunt script to index your all of your content files and [lunr.js](http://lunrjs.com/) to serve the search results.

***

## Other

And for all the other small things around Hugo:

- [hugo-gallery](https://github.com/icecreammatt/hugo-gallery) lets you create an image gallery for Hugo sites.

***

> Do you know or maintain a similar project around Hugo? Feel free to open a
[pull request](https://github.com/spf13/hugo/pulls) on Github if you think it should to be added.
