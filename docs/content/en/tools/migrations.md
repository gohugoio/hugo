---
title: Migrate to Hugo
linkTitle: Migrations
description: A list of community-developed tools for migrating from your existing static site generator or content management system to Hugo.
categories: [developer tools]
keywords: [migrations,jekyll,wordpress,drupal,ghost,contentful]
menu:
  docs:
    parent: developer-tools
    weight: 50
weight: 50
toc: true
aliases: [/developer-tools/migrations/, /developer-tools/migrated/]
---

This section highlights some independently developed projects related to Hugo. These tools extend functionality or help you to get started.

Take a look at this list of migration tools if you currently use other blogging tools like Jekyll or WordPress but intend to switch to Hugo instead. They'll help you export your content into Hugo-friendly formats.

## Jekyll

Alternatively, you can use the [Jekyll import command](/commands/hugo_import_jekyll/).

[JekyllToHugo](https://github.com/fredrikloch/JekyllToHugo)
: A Small script for converting Jekyll blog posts to a Hugo site.

[ConvertToHugo](https://github.com/coderzh/ConvertToHugo)
: Convert your blog from Jekyll to Hugo.

## Octopress

[octohug](https://github.com/codebrane/octohug)
: Octopress to Hugo migrator.

## DokuWiki

[dokuwiki-to-hugo](https://github.com/wgroeneveld/dokuwiki-to-hugo)
: Migrates your DokuWiki source pages from [DokuWiki syntax](https://www.dokuwiki.org/wiki:syntax) to Hugo Markdown syntax. Includes extras like the TODO plugin. Written with extensibility in mind using Python 3. Also generates a TOML header for each page. Designed to copy-paste the wiki directory into your /content directory.

## WordPress

[wordpress-to-hugo-exporter](https://github.com/SchumacherFM/wordpress-to-hugo-exporter)
: A one-click WordPress plugin that converts all posts, pages, taxonomies, metadata, and settings to Markdown and YAML which can be dropped into Hugo. (Note: If you have trouble using this plugin, you can [export your site for Jekyll](https://wordpress.org/plugins/jekyll-exporter/) and use Hugo's built-in Jekyll converter listed above.)

[blog2md](https://github.com/palaniraja/blog2md)
: Works with [exported xml](https://en.support.wordpress.com/export/) file of your free YOUR-TLD.wordpress.com website. It also saves approved comments to `YOUR-POST-NAME-comments.md` file along with posts.

[wordhugopress](https://github.com/nantipov/wordhugopress)
: A small utility written in Java that exports the entire WordPress site from the database and resource (e.g., images) files stored locally or remotely. Therefore, migration from the backup files is possible. Supports merging multiple WordPress sites into a single Hugo site.

[wp2hugo](https://github.com/ashishb/wp2hugo)
: A Go-based CLI tool to migrate WordPress website to Hugo while preserving original URLs, GUIDs (for feeds), image URLs, code highlights, table of contents, YouTube embeds, Google Maps embeds, and original WordPress navigation categories. 

## Medium

[medium2md](https://github.com/gautamdhameja/medium-2-md)
: A simple Medium to Hugo exporter able to import stories in one command, including front matter.

[medium-to-hugo](https://github.com/bgadrian/medium-to-hugo)
: A CLI tool written in Go to export medium posts into a Hugo-compatible Markdown format. Tags and images are included. All images will be downloaded locally and linked appropriately.

## Tumblr

[tumblr-importr](https://github.com/carlmjohnson/tumblr-importr)
: An importer that uses the Tumblr API to create a Hugo static site.

[tumblr2hugomarkdown](https://github.com/Wysie/tumblr2hugomarkdown)
: Export all your Tumblr content to Hugo Markdown files with preserved original formatting.

[Tumblr to Hugo](https://github.com/jipiboily/tumblr-to-hugo)
: A migration tool that converts each of your Tumblr posts to a content file with a proper title and path. It also generates a  CSV file to help you set up URL redirects.

## Drupal

[drupal2hugo](https://github.com/danapsimer/drupal2hugo)
: Convert a Drupal site to Hugo.

## Joomla

[hugojoomla](https://github.com/davetcc/hugojoomla)
: This utility written in Java takes a Joomla database and converts all the content into Markdown files. It changes any URLs that are in Joomla's internal format and converts them to a suitable form.

## Blogger

[blogimport](https://github.com/natefinch/blogimport)
: A tool to import from Blogger posts to Hugo.

[blogger-to-hugo](https://pypi.org/project/blogger-to-hugo/)
: Another tool to import Blogger posts to Hugo. It also downloads embedded images so they will be stored locally.

[blog2md](https://github.com/palaniraja/blog2md)
: Works with [exported xml](https://support.google.com/blogger/answer/41387?hl=en) file of your YOUR-TLD.blogspot.com website. It also saves comments to `YOUR-POST-NAME-comments.md` file along with posts.

[BloggerToHugo](https://github.com/huanlin/blogger-to-hugo)
: Yet another tool to import Blogger posts to Hugo. For Windows platform only, and .NET Framework 4.5 is required. See README.md before using this tool.

## Contentful

[contentful-hugo](https://github.com/ModiiMedia/contentful-hugo)
: A tool to create content-files for Hugo from content on [Contentful](https://www.contentful.com/).

## BlogML

[BlogML2Hugo](https://github.com/jijiechen/BlogML2Hugo)
: A tool that helps you convert BlogML xml file to Hugo Markdown files. Users need to take care of links to attachments and images by themselves. This helps the blogs that export BlogML files (e.g. BlogEngine.NET) transform to hugo sites easily.
