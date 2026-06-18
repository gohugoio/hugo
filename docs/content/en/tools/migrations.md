---
title: Migrate to Hugo
linkTitle: Migrations
description: A list of community-developed tools for migrating from your existing static site generator or content management system to Hugo.
categories: []
keywords: []
weight: 40
aliases: [/developer-tools/migrations/, /developer-tools/migrated/]
---

This section highlights some independently developed projects related to Hugo. These tools extend functionality or help you to get started.

Take a look at this list of migration tools if you currently use other blogging tools like Jekyll or WordPress but intend to switch to Hugo instead. They'll help you export your content into Hugo-friendly formats.

## Jekyll

Alternatively, you can use the [Jekyll import command](/commands/hugo_import_jekyll/).

[JekyllToHugo][]
: A Small script for converting Jekyll blog posts to a Hugo site.

[ConvertToHugo][]
: Convert your blog from Jekyll to Hugo.

## Octopress

[octohug][]
: Octopress to Hugo migrator.

## DokuWiki

[dokuwiki-to-hugo][]
: Migrates your DokuWiki source pages from [DokuWiki syntax][] to Hugo Markdown syntax. Includes extras like the TODO plugin. Written with extensibility in mind using Python 3. Also generates a TOML header for each page. Designed to copy-paste the wiki directory into your `content` directory.

## WordPress

[wordpress-to-hugo-exporter][]
: A one-click WordPress plugin that converts all posts, pages, taxonomies, metadata, and settings to Markdown and YAML which can be dropped into Hugo. (Note: If you have trouble using this plugin, you can [export your site for Jekyll][] and use Hugo's built-in Jekyll converter listed above.)

[blog2md][]
: Works with [exported xml](https://en.support.wordpress.com/export/) file of your free YOUR-TLD.wordpress.com website. It also saves approved comments to `YOUR-POST-NAME-comments.md` file along with posts.

[wordhugopress][]
: A small utility written in Java that exports the entire WordPress site from the database and resource (e.g., images) files stored locally or remotely. Therefore, migration from the backup files is possible. Supports merging multiple WordPress sites into a single Hugo site.

[wp2hugo][]
: A Go-based CLI tool to migrate WordPress websites to Hugo. It preserves original URLs, GUIDs, image URLs, code highlights, tables of contents, and WordPress navigation categories. It migrates WordPress custom post types, custom taxonomies, custom fields, and page hierarchy. It supports translated WordPress blogs via Polylang or WPML. It imports a WordPress media library database with original titles and dates. The tool can download all media or only media inserted into pages from the original server. It converts WordPress shortcodes and Gutenberg blocks to Hugo shortcodes including galleries, images, audio, YouTube embeds, Gists, and Google Maps.

## Medium

[medium2md][]
: A simple Medium to Hugo exporter able to import stories in one command, including front matter.

[medium-to-hugo][]
: A CLI tool written in Go to export medium posts into a Hugo-compatible Markdown format. Tags and images are included. All images will be downloaded locally and linked appropriately.

## Tumblr

[tumblr-importr][]
: An importer that uses the Tumblr API to create a Hugo static site.

[tumblr2hugomarkdown][]
: Export all your Tumblr content to Hugo Markdown files with preserved original formatting.

[Tumblr to Hugo][]
: A migration tool that converts each of your Tumblr posts to a content file with a proper title and path. It also generates a CSV file to help you set up URL redirects.

## Drupal

[drupal2hugo][]
: Convert a Drupal site to Hugo.

## Joomla

[hugojoomla][]
: This utility written in Java takes a Joomla database and converts all the content into Markdown files. It changes any URLs that are in Joomla's internal format and converts them to a suitable form.

## Blogger

[blogimport][]
: A tool to import from Blogger posts to Hugo.

[blogger-to-hugo][]
: Another tool to import Blogger posts to Hugo. It also downloads embedded images so they will be stored locally.

[blog2md][]
: Works with [exported xml](https://support.google.com/blogger/answer/41387?hl=en) file of your YOUR-TLD.blogspot.com website. It also saves comments to `YOUR-POST-NAME-comments.md` file along with posts.

[BloggerToHugo][]
: Yet another tool to import Blogger posts to Hugo. For Windows platform only, and .NET Framework 4.5 is required. See README.md before using this tool.

[blogger2hugo][]
: Converts a Blogger backup file (`.atom`) from [Google Takeout][] to Markdown (`.md`) files. The tool generates output compatible with the Hugo `content/` structure.

## Contentful

[contentful-hugo][]
: A tool to create content-files for Hugo from content on [Contentful][].

## BlogML

[BlogML2Hugo][]
: A tool that helps you convert BlogML xml file to Hugo Markdown files. Users need to take care of links to attachments and images by themselves. This helps the blogs that export BlogML files (e.g. BlogEngine.NET) transform to hugo sites easily.

[BlogML2Hugo]: https://github.com/jijiechen/BlogML2Hugo
[BloggerToHugo]: https://github.com/huanlin/blogger-to-hugo
[Contentful]: https://www.contentful.com/
[ConvertToHugo]: https://github.com/coderzh/ConvertToHugo
[DokuWiki syntax]: https://www.dokuwiki.org/wiki:syntax
[Google Takeout]: https://takeout.google.com/takeout/custom/blogger?hl=en
[JekyllToHugo]: https://github.com/fredrikloch/JekyllToHugo
[Tumblr to Hugo]: https://github.com/jipiboily/tumblr-to-hugo
[blog2md]: https://github.com/palaniraja/blog2md
[blogger-to-hugo]: https://pypi.org/project/blogger-to-hugo/
[blogger2hugo]: https://github.com/noorkhafidzin/blogger2hugo
[blogimport]: https://github.com/natefinch/blogimport
[contentful-hugo]: https://github.com/ModiiMedia/contentful-hugo
[dokuwiki-to-hugo]: https://github.com/wgroeneveld/dokuwiki-to-hugo
[drupal2hugo]: https://github.com/danapsimer/drupal2hugo
[export your site for Jekyll]: https://wordpress.org/plugins/jekyll-exporter/
[hugojoomla]: https://github.com/davetcc/hugojoomla
[medium-to-hugo]: https://github.com/bgadrian/medium-to-hugo
[medium2md]: https://github.com/gautamdhameja/medium-2-md
[octohug]: https://github.com/codebrane/octohug
[tumblr-importr]: https://github.com/carlmjohnson/tumblr-importr
[tumblr2hugomarkdown]: https://github.com/Wysie/tumblr2hugomarkdown
[wordhugopress]: https://github.com/nantipov/wordhugopress
[wordpress-to-hugo-exporter]: https://github.com/SchumacherFM/wordpress-to-hugo-exporter
[wp2hugo]: https://github.com/ashishb/wp2hugo
