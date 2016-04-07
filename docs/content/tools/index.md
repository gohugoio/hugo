---
lastmod: 2015-12-30
date: 2015-09-12T10:40:31+02:00
title: Tools
weight: 120
---

This section highlights some projects around Hugo that are independently developed. These tools try to extend the functionality of our static site generator or help you to get started.


## Migration

Take a look at this list of migration tools if you currently use other blogging tools like Jekyll or WordPress but intend to switch to Hugo instead. They'll take care to export
your content into Hugo-friendly formats.

### Jekyll

Alternatively, you can follow the manual [migration guide]({{< relref "tutorials/migrate-from-jekyll.md" >}}) or use the new [Jekyll import command]({{< relref "commands/hugo_import_jekyll.md" >}}).

- [JekyllToHugo](https://github.com/SenjinDarashiva/JekyllToHugo) - A Small script for converting Jekyll blog posts to a Hugo site.
- [ConvertToHugo](https://github.com/coderzh/ConvertToHugo) - Convert your blog from Jekyll to Hugo.

### Ghost

- [ghostToHugo](https://github.com/jbarone/ghostToHugo) - Convert Ghost blog posts and export them to Hugo.

### Octopress

- [octohug](https://github.com/codebrane/octohug) - Octopress to Hugo migrator.

### WordPress

- [wordpress-to-hugo-exporter](https://github.com/SchumacherFM/wordpress-to-hugo-exporter) - A one-click WordPress plugin that converts all posts, pages, taxonomies, metadata, and settings to Markdown and YAML which can be dropped into Hugo. (Note: If you have trouble using this plugin, you can [export your site for Jekyll](https://wordpress.org/plugins/jekyll-exporter/) and use Hugo's built in Jekyll converter listed above.)

### Drupal

- [drupal2hugo](https://github.com/danapsimer/drupal2hugo) - Convert a Drupal site to Hugo.

### Blogger

- [blogimport](https://github.com/natefinch/blogimport) - A tool to import from Blogger posts to Hugo.

### Contentful

- [contentful2hugo](https://github.com/ArnoNuyts/contentful2hugo) - A tool to create content-files for Hugo from content on [Contentful](https://www.contentful.com/).

----

## Deployment

If you don't want to use [Wercker for automated deployments]({{< relref "tutorials/automated-deployments.md" >}}), give these tools a try to get your content to the public:

- [hugomac](https://github.com/nickoneill/hugomac) - Hugomac is an OS&nbsp;X menubar app to publish your blog directly to Amazon S3. No command line is needed.
- [hugo-lambda](https://github.com/ryansb/hugo-lambda) - A wrapper around the Hugo static site generator to have it run in AWS Lambda whenever new (Markdown or other) content is uploaded.
- [hugodeploy](https://github.com/mindok/hugodeploy) - Simple SFTP deployment tool for static websites (e.g. created by Hugo) with optional minification.

----

## Frontends

Do you prefer an graphical user interface over a text editor? Then give these frontends a try:

- [rango](https://github.com/stayradiated/rango) - A web frontend for Hugo. It's designed to make it easy to manage a small site, even for people with little computer experience.
- [enwrite](https://github.com/zzamboni/enwrite) - Evernote-powered statically-generated blogs and websites. Now posting to your blog or updating your website is as easy as writing a new note in Evernote!
- [caddy-hugo](https://github.com/hacdias/caddy-hugo) - This is an add-on for [Caddy](https://caddyserver.com/) which wants to deliver a good UI to edit the content of the website.
- [Hugopit](https://github.com/sjardim/Hugopit) - A web-based editor for Hugo build on top of [Cockpit CMS](http://www.getcockpit.com/).

----

## Editor plugins

If you still want to use an editor, look at these plugins to automate your workflow:

### Sublime Text

- [Hugofy](https://github.com/akmittal/Hugofy) - Hugofy is a plugin for Sublime Text 3 to make life easier to use Hugo static site generator.

----

## Search

A static site with a dynamic search function? Yes. Alternatively to embeddable scripts from Google or other search engines you can provide your visitors a custom search by indexing your content files directly.

- [Hugoidx](https://github.com/blevesearch/hugoidx) is an experimental application to create a search index. It's build on top of [Bleve](http://www.blevesearch.com/).
- This [GitHub Gist](https://gist.github.com/sebz/efddfc8fdcb6b480f567) contains simple workflow to create a search index for your static site. It uses a simple Grunt script to index all your content files and [lunr.js](http://lunrjs.com/) to serve the search results.
- [hugo-lunr](https://www.npmjs.com/package/hugo-lunr) - A simple way to add site search to your static Hugo site using [lunr.js](http://lunrjs.com/). Hugo-lunr will create an index file of any html and markdown documents in your Hugo project.

----

## Commercial Services

- [Appernetic.io](https://appernetic.io) is a Hugo Static Site Generator as a Service that is easy to use for non-technical users.
Features: inline PageDown editor, visual tree view, image upload and digital asset management with Cloudinary, site preview, continuous integration with GitHub, atomic deploy and hosting, Git and Hugo integration, autosave, custom domain, project syncing, theme cloning and management. Developers have complete control over the source code and can manage it with GitHub’s deceptively simple workflow.

- [Netlify.com](https://www.netlify.com), builds, deploy & hosts your static site or app (Hugo, Jekyll etc). Build, deploy and host your static site or app with a drag and drop interface and automatic deploys from GitHub or Bitbucket.
Features: global CDN, atomic deploys, ultra fast DNS, instant cache invalidation, high availability, automated hosting, Git integration, form submission hooks, authentication providers, custom domain. Developers have complete control over the source code and can manage it with GitHub’s or Bitbuckets deceptively simple workflow.

----

## Other

And for all the other small things around Hugo:

- [hugo-gallery](https://github.com/icecreammatt/hugo-gallery) lets you create an image gallery for Hugo sites.
- [flickr-hugo-embed](https://github.com/nikhilm/flickr-hugo-embed) prints shortcodes to embed a set of images from an album on Flickr into Hugo.

----

> Do you know or maintain a similar project around Hugo? Feel free to open a
[pull request](https://github.com/spf13/hugo/pulls) on GitHub if you think it should be added.
