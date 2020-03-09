---
title: Migrate to Hugo
linktitle: Migrations
description: A list of community-developed tools for migrating from your existing static site generator or content management system to Hugo.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
keywords: [migrations,jekyll,wordpress,drupal,ghost,contentful]
menu:
  docs:
    parent: "tools"
    weight: 10
weight: 10
sections_weight: 10
draft: false
aliases: [/developer-tools/migrations/,/developer-tools/migrated/]
toc: true
---

This section highlights some projects around Hugo that are independently developed. These tools try to extend the functionality of our static site generator or help you to get started.

{{% note %}}
Do you know or maintain a similar project around Hugo? Feel free to open a [pull request](https://github.com/gohugoio/hugo/pulls) on GitHub if you think it should be added.
{{% /note %}}

Take a look at this list of migration tools if you currently use other blogging tools like Jekyll or WordPress but intend to switch to Hugo instead. They'll take care to export your content into Hugo-friendly formats.

## Jekyll

Alternatively, you can use the new [Jekyll import command](/commands/hugo_import_jekyll/).

- [JekyllToHugo](https://github.com/SenjinDarashiva/JekyllToHugo) - A Small script for converting Jekyll blog posts to a Hugo site.
- [ConvertToHugo](https://github.com/coderzh/ConvertToHugo) - Convert your blog from Jekyll to Hugo.

## Ghost

- [ghostToHugo](https://github.com/jbarone/ghostToHugo) - Convert Ghost blog posts and export them to Hugo.

## Octopress

- [octohug](https://github.com/codebrane/octohug) - Octopress to Hugo migrator.

## DokuWiki

- [dokuwiki-to-hugo](https://github.com/wgroeneveld/dokuwiki-to-hugo) - Migrates your dokuwiki source pages from [DokuWiki syntax](https://www.dokuwiki.org/wiki:syntax) to Hugo Markdown syntax. Includes extra's like the TODO plugin. Written with extensibility in mind using python 3. Also generates a TOML header for each page. Designed to copypaste the wiki directory into your /content directory.

## WordPress

- [wordpress-to-hugo-exporter](https://github.com/SchumacherFM/wordpress-to-hugo-exporter) - A one-click WordPress plugin that converts all posts, pages, taxonomies, metadata, and settings to Markdown and YAML which can be dropped into Hugo. (Note: If you have trouble using this plugin, you can [export your site for Jekyll](https://wordpress.org/plugins/jekyll-exporter/) and use Hugo's built in Jekyll converter listed above.)
- [exitwp-for-hugo](https://github.com/wooni005/exitwp-for-hugo) - A python script which works with the xml export from Wordpress and converts Wordpress pages and posts to Markdown and YAML for hugo.
- [blog2md](https://github.com/palaniraja/blog2md) - Works with [exported xml](https://en.support.wordpress.com/export/) file of your free YOUR-TLD.wordpress.com website. It also saves approved comments to `YOUR-POST-NAME-comments.md` file along with posts.
- [wordhugopress](https://github.com/nantipov/wordhugopress) - A small utility written in Java, exports the entire WordPress site from the database and resource (e.g. images) files stored locally or remotelly. Therefore, migration from the backup files is possible. Supports merging of the multiple WordPress sites into a single Hugo one.

## Medium

- [medium2md](https://github.com/gautamdhameja/medium-2-md) - A simple Medium to Hugo exporter able to import stories in one command, including Front Matter.
- [medium-to-hugo](https://github.com/bgadrian/medium-to-hugo) - CLI tool written in Go to export medium posts into a Hugo compatible Markdown format. Tags and images are included. All images will be downloaded locally and linked appropriately. 

## Tumblr

- [tumblr-importr](https://github.com/carlmjohnson/tumblr-importr) - An importer that uses the Tumblr API to create a Hugo static site.
- [tumblr2hugomarkdown](https://github.com/Wysie/tumblr2hugomarkdown) - Export all your Tumblr content to Hugo Markdown files with preserved original formatting.
- [Tumblr to Hugo](https://github.com/jipiboily/tumblr-to-hugo) - A migration tool that converts each of your Tumblr posts to a content file with a proper title and path. Furthermore, "Tumblr to Hugo" creates a CSV file with the original URL and the new path on Hugo, to help you setup the redirections.

## Drupal

- [drupal2hugo](https://github.com/danapsimer/drupal2hugo) - Convert a Drupal site to Hugo.

## Joomla

- [hugojoomla](https://github.com/davetcc/hugojoomla) - This utility written in Java takes a Joomla database and converts all the content into Markdown files. It changes any URLs that are in Joomla's internal format and converts them to a suitable form.

## Blogger

- [blogimport](https://github.com/natefinch/blogimport) - A tool to import from Blogger posts to Hugo.
- [blogger-to-hugo](https://bitbucket.org/petraszd/blogger-to-hugo) - Another tool to import Blogger posts to Hugo. It also downloads embedded images so they will be stored locally.
- [blog2md](https://github.com/palaniraja/blog2md) - Works with [exported xml](https://support.google.com/blogger/answer/41387?hl=en) file of your YOUR-TLD.blogspot.com website. It also saves comments to `YOUR-POST-NAME-comments.md` file along with posts. 
- [BloggerToHugo](https://github.com/huanlin/blogger-to-hugo) - Yet another tool to import Blogger posts to Hugo. For Windows platform only, and .NET Framework 4.5 is required. See README.md before using this tool.

## Contentful

- [contentful2hugo](https://github.com/ArnoNuyts/contentful2hugo) - A tool to create content-files for Hugo from content on [Contentful](https://www.contentful.com/).


## BlogML

- [BlogML2Hugo](https://github.com/jijiechen/BlogML2Hugo) - A tool that helps you convert BlogML xml file to Hugo markdown files. Users need to take care of links to attachments and images by themselves. This helps the blogs that export BlogML files (e.g. BlogEngine.NET) tramsform to hugo sites easily.
