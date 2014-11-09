---
date: 2013-07-01
linktitle: Quickstart
menu:
  main:
    parent: getting started
next: /overview/installing
prev: /overview/introduction
title: Hugo Quickstart Guide
weight: 10
---

_This quickstart depends on features introduced in hugo v0.11. If you
have an earlier version of hugo you will need to [upgrade](/overview/installing/) before
proceeding._

## Step 1. Install Hugo

Go to [Hugo Releases](https://github.com/spf13/hugo/releases) and download the
appropriate version for your os and architecture.

Save it somewhere specific as we will be using it in the next step.

More complete instructions are available at [Installing Hugo](/overview/installing/)

## Step 2. Have Hugo Create a site for you

Hugo has the ability to create a skeleton site.

    hugo new site /path/to/site

For the rest of the operations we will be executing all commands from within the site directory

    cd /path/to/site

The new site will have the following structure

      ▸ archetypes/
      ▸ content/
      ▸ layouts/
      ▸ static/
        config.toml

Currently the site doesn’t have any content, nor is it configured.

## Step 3. Create Some Content

Hugo also has the ability to create content for you.

    hugo new about.md

A new file is now created in `content/` with the following contents

    +++
    draft = true
    title = "about"
    date = 2014-05-20T10:04:31Z
    +++

Notice the date is automatically set to the moment you created the content.

Place some content in this file below the `+++` in the Markdown format.

For example you could put this

    ## A headline

    Some Content

For fun, let’s create another piece of content and place some Markdown in it as well.

    hugo new post/first.md

The new file is located at `content/post/first.md`

We still lack any templates to tell us how to display the content.

## Step 4. Install some themes

Hugo has rich theme support and a growing set of themes to choose from:

    git clone --recursive https://github.com/spf13/hugoThemes themes

## Step 5. Run Hugo

Hugo contains its own high performance web server. Simply run `hugo
server` and Hugo will find an available port and run a server with
your content:

    hugo server --theme=hyde --buildDrafts
    2 pages created
    0 tags created
    0 categories created
    in 5 ms
    Serving pages from exampleHugoSite/public
    Web Server is available at http://localhost:1313
    Press ctrl+c to stop

We specified two options here:

 * `--theme` to pick which theme;
 * `--buildDrafts` because we want to display our content, both set to draft status.

To learn about what other options hugo has, run:

    hugo help

To learn about the server options:

    hugo help server

## Step 6. Edit Content

Not only can Hugo run a server, but it can also watch your files for
changes and automatically rebuild your site. Hugo will then
communicate with your browser and automatically reload any open page.
This even works in mobile browsers.

Stop the Hugo process by hitting ctrl+c. Then run the following:

    hugo server --theme=hyde --buildDrafts --watch
    2 pages created
    0 tags created
    0 categories created
    in 5 ms
    Watching for changes in exampleHugoSite/content
    Serving pages from exampleHugoSite/public
    Web Server is available at http://localhost:1313
    Press ctrl+c to stop

Open your [favorite editor](http://vim.spf13.com), edit and save your content and watch as Hugo rebuilds and reloads automatically.

It’s especially productive to leave a browser open on a second monitor
and just glance at it whenever you save. You don’t even need to tab to
your browser. Hugo is so fast that the new site will be there before
you can look at the browser in most cases.

Change and save this file. Notice what happened in your terminal.

    Change detected, rebuilding site

    2 pages created
    0 tags created
    0 categories created
    in 5 ms

## Step 7. Have fun

The best way to learn something is to play with it.

Things to try:

 * Add a [new content file](/content/organization/)
 * Create a [new section](/content/sections/)
 * Modify [a template](/layout/templates/)
 * Create content with [TOML front matter](/content/front-matter/)
 * Define your own field in [front matter](/content/front-matter/)
 * Display that [field in the template](/layout/variables/)
 * Create a [new content type](/content/types/)
