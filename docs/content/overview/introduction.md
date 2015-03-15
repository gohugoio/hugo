---
date: 2013-07-01
linktitle: Introduction
menu:
  main:
    parent: getting started
next: /overview/quickstart
title: Introduction to Hugo
weight: 5
---

## What is Hugo?

Hugo is a general-purpose website framework. Technically speaking, Hugo is
a static site generator. This means that, unlike systems like WordPress,
Ghost and Drupal, which run on your web server expensively building a page
every time a visitor requests one, Hugo does the building when you create
your content. Since websites are viewed far more often than they are
edited, Hugo is optimized for website viewing while providing a great
writing experience.

Sites built with Hugo are extremely fast and very secure. Hugo sites can
be hosted anywhere, including [Heroku][], [GoDaddy][], [DreamHost][],
[GitHub Pages][], [Amazon S3][] and [CloudFront][], and work well with CDNs.
Hugo sites run without dependencies on expensive runtimes like Ruby,
Python or PHP and without dependencies on any databases.

[Heroku]: https://www.heroku.com/
[GoDaddy]: https://www.godaddy.com/
[DreamHost]: http://www.dreamhost.com/
[GitHub Pages]: https://pages.github.com/
[Amazon S3]: http://aws.amazon.com/s3/
[CloudFront]: http://aws.amazon.com/cloudfront/ "Amazon CloudFront"

We think of Hugo as the ideal website creation tool. With nearly instant
build times and the ability to rebuild whenever a change is made, Hugo
provides a very fast feedback loop. This is essential when you are
designing websites, but also very useful when creating content.

## How fast is Hugo?

{{% youtube CdiDYZ51a2o %}}

## What does Hugo do?

In technical terms, Hugo takes a source directory of Markdown files and
templates and uses these as input to create a complete website.

Hugo boasts the following features:

### General

  * Extremely fast build times (~1&nbsp;ms per page)
  * Completely cross platform: Runs on <i class="fa fa-apple"></i>&nbsp;Mac OS&nbsp;X, <i class="fa fa-linux"></i>&nbsp;Linux, <i class="fa fa-windows"></i>&nbsp;Windows, and more!
  * Easy [installation](/overview/installing/)
  * Render changes [on the fly](/overview/usage/) with [LiveReload](/extras/livereload/) as you develop
  * Complete theme support
  * Host your site anywhere

### Organization

  * Straightforward [organization](/content/organization/)
  * Support for [website sections](/content/sections/)
  * Completely customizable [URLs](/extras/urls/)
  * Support for configurable [taxonomies](/taxonomies/overview/) which includes categories and tags.  Create your own custom organization of content
  * Ability to [sort content](/content/ordering/) as you desire
  * Automatic [table of contents](/extras/toc/) generation
  * Dynamic menu creation
  * [Pretty URLs](/extras/urls/) support
  * [Permalink](/extras/permalinks/) pattern support
  * [Aliases](/extras/aliases/) (redirects)

### Content

  * Content written in [Markdown](/content/example/)
  * Support for TOML, YAML and JSON metadata in [frontmatter](/content/front-matter/)
  * Completely [customizable homepage](/layout/homepage/)
  * Support for multiple [content types](/content/types/)
  * Automatic and user defined [summaries](/content/summaries/)
  * [Shortcodes](/extras/shortcodes/) to enable rich content inside of Markdown
  * ["Minutes to Read"](/layout/variables/) functionality
  * ["Wordcount"](/layout/variables/) functionality

### Additional Features

  * Integrated [Disqus](https://disqus.com/) comment support
  * Automatic [RSS](/layout/rss/) creation
  * Support for [Go](http://golang.org/pkg/html/template/), [Amber](https://github.com/eknkc/amber) and [Ace](http://ace.yoss.si/) HTML templates
  * Syntax [highlighting](/extras/highlighting/) powered by [Pygments](http://pygments.org/)

See what's coming next in the [roadmap](/meta/roadmap/).

## Who should use Hugo?

Hugo is for people that prefer writing in a text editor over
a browser.

Hugo is for people who want to hand code their own website without
worrying about setting up complicated runtimes, dependencies and
databases.

Hugo is for people building a blog, company site, portfolio, tumblog,
documentation, single page site or a site with thousands of
pages.

## Why did you write Hugo?

I wrote Hugo ultimately for a few reasons. First, I was disappointed with
WordPress, my then website solution. It rendered slowly. I couldn't create
content as efficiently as I wanted to and needed to be online to write
posts. The constant security updates and the horror stories of people's
hacked blogs. I hated how content was written in HTML instead of the much
simpler Markdown. Overall, I felt like it got in my way more than it helped
me from writing great content.

I looked at existing static site generators like [Jekyll][], [Middleman][] and [nanoc][].
All had complicated dependencies to install and took far longer to render
my blog with hundreds of posts than I felt was acceptable. I wanted
a framework to be able to get rapid feedback while making changes to the
templates, and the 5+-minute render times was just too slow. In general,
they were also very blog minded and didn't have the ability to have
different content types and flexible URLs.

[Jekyll]: http://jekyllrb.com/
[Middleman]: https://middlemanapp.com/
[nanoc]: http://nanoc.ws/

I wanted to develop a fast and full-featured website framework without
dependencies. The [Go language][] seemed to have all of the features I needed
in a language. I began developing Hugo in Go and fell in love with the
language. I hope you will enjoy using (and contributing to) Hugo as much
as I have writing it.

[Go language]: http://golang.org/ "The Go Programming Language"

## Next Steps

 * [Install Hugo](/overview/installing/)
 * [Quick start](/overview/quickstart/)
 * [Join the Mailing List](/community/mailing-list/)
 * [Star us on GitHub](https://github.com/spf13/hugo)
 * [Discussion Forum](http://discuss.gohugo.io/)

