---
lastmod: 2015-12-24
date: 2014-05-26
linktitle: Builders
menu:
  main:
    parent: extras
next: /extras/comments
prev: /extras/analytics
title: Hugo Builders
---

Hugo provides the functionality to quickly get a site, theme or page
started.


## New Site

Want to get a site built quickly?

{{< nohighlight >}}$ hugo new site <i>path/to/site</i>
{{< /nohighlight >}}

Hugo will create all the needed directories and files to get started
quickly.

Hugo will only touch the files and create the directories (in the right
places), [configuration](/overview/configuration/) and content are up to
you... but luckily we have builders for content (see below).

## New Theme

Want to design a new theme?

    $ hugo new theme THEME_NAME

Run from your working directory, this will create a new theme with all
the needed files in your themes directory. Hugo will provide you with a
license and theme.toml file with most of the work done for you.

Follow the [Theme Creation Guide](/themes/creation/) once the builder is
done.

## New Content

You will use this builder the most of all. Every time you want to create
a new piece of content, the content builder will get you started right.

Leveraging [content archetypes](/content/archetypes/) the content builder
will not only insert the current date and appropriate metadata, but it
will pre-populate values based on the content type.

    $ hugo new relative/path/to/content

This assumes it is being run from your working directory and the content
path starts from your content directory. Now, Hugo watches your content directory by default and rebuilds your entire website if any change occurs.
