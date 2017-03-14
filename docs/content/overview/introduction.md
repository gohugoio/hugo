---
lastmod: 2016-08-14
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
a static site generator. Unlike other systems which dynamically build a page
every time a visitor requests one, Hugo does the building when you create
your content. Since websites are viewed far more often than they are
edited, Hugo is optimized for website viewing while providing a great
writing experience.

Sites built with Hugo are extremely fast and very secure. Hugo sites can
be hosted anywhere, including [Heroku][], [GoDaddy][], [DreamHost][],
[GitHub Pages][], [Netlify][], [Surge][], [Aerobatic][], [Firebase Hosting][],
[Google Cloud Storage][], [Amazon S3][] and [CloudFront][], and work well
with CDNs. Hugo sites run without dependencies on expensive runtimes
like Ruby, Python or PHP and without dependencies on any databases.

[Heroku]: https://www.heroku.com/
[GoDaddy]: https://www.godaddy.com/
[DreamHost]: http://www.dreamhost.com/
[GitLab]: https://about.gitlab.com
[GitHub Pages]: https://pages.github.com/
[Aerobatic]: https://www.aerobatic.com/
[Firebase Hosting]: https://firebase.google.com/docs/hosting/
[Google Cloud Storage]: http://cloud.google.com/storage/
[Amazon S3]: http://aws.amazon.com/s3/
[CloudFront]: http://aws.amazon.com/cloudfront/ "Amazon CloudFront"
[Surge]: https://surge.sh
[Netlify]: https://www.netlify.com

We think of Hugo as the ideal website creation tool. With nearly instant
build times and the ability to rebuild whenever a change is made, Hugo
provides a very fast feedback loop. This is essential when you are
designing websites, but also very useful when creating content.

## What makes Hugo different?

Web site generators render content into HTML files. Most are "dynamic
site generators." That means the HTTP
server (which is the program running on your website that the user's
browser talks to) runs the generator to create a new HTML file
each and every time a user wants to view a page.

Creating the page dynamically means that the computer hosting
the HTTP server has to have enough memory and CPU to effectively run
the generator around the clock. If not, then the user has to wait
in a queue for the page to be generated.

Nobody wants users to wait longer than needed, so the dynamic site
generators programmed their systems to cache the HTML files. When
a file is cached, a copy of it is temporarily stored on the computer.
It is much faster for the HTTP server to send that copy the next time
the page is requested than it is to generate it from scratch.

Hugo takes caching a step further. All HTML files are rendered on your
computer. You can review the files before you copy them to the computer
hosting the HTTP server. Since the HTML files aren't generated dynamically,
we say that Hugo is a "static site generator."

Not running a web site generator on your HTTP server has many benefits.
The most noticeable is performance - HTTP servers are very good at
sending files. So good that you can effectively serve the same number
of pages with a fraction of the memory and CPU needed for a dynamic site.

Hugo has two components to help you build and test your web site. The
one that you'll probably use most often is the built-in HTTP server.
When you run `hugo server`, Hugo renders all of your content into
HTML files and then runs an HTTP server on your computer so that you
can see what the pages look like.

The second component is used when you're ready to publish your web
site to the computer running your website. Running Hugo without any
actions will rebuild your entire web site using the `baseURL` setting
from your site's configuration file. That's required to have your page
links work properly with most hosting companies.

## How fast is Hugo?

{{% youtube CdiDYZ51a2o %}}

## What does Hugo do?

In technical terms, Hugo takes a source directory of files and
templates and uses these as input to create a complete website.

Hugo boasts the following features:

### General

  * Extremely fast build times (~1 ms per page)
  * Completely cross platform: Runs on <i class="fa fa-apple"></i>&nbsp;macOS, <i class="fa fa-linux"></i>&nbsp;Linux, <i class="fa fa-windows"></i>&nbsp;Windows, and more!
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

  * Native support for content written in [Markdown](/content/example/)
    * Support for other languages through _external helpers_, see [supported formats](/content/supported-formats)
  * Support for TOML, YAML and JSON metadata in [frontmatter](/content/front-matter/)
  * Completely [customizable homepage](/layout/homepage/)
  * Support for multiple [content types](/content/types/)
  * Automatic and user defined [summaries](/content/summaries/)
  * [Shortcodes](/extras/shortcodes/) to enable rich content inside of Markdown
  * ["Minutes to Read"](/layout/variables/) functionality
  * ["Wordcount"](/layout/variables/) functionality

### Additional Features

  * Integrated [Disqus](https://disqus.com/) comment support
  * Integrated [Google Analytics](https://google-analytics.com/) support
  * Automatic [RSS](/layout/rss/) creation
  * Support for [Go](http://golang.org/pkg/html/template/), [Amber](https://github.com/eknkc/amber) and [Ace](https://github.com/yosssi/ace) HTML templates
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
WordPress, my then website solution. With it, I couldn't create
content as efficiently as I wanted to.
It rendered slowly. It required me to be online to write
posts: plus its constant security updates and the horror stories of people's
hacked blogs! I hated how content for it was written only in HTML, instead of the much
simpler Markdown. Overall, I felt like WordPress got in my way
much more than it helped me. It kept
me from writing great content.

I looked at the existing static site generators
like [Jekyll][], [Middleman][] and [Nanoc][].
All had complicated installation dependencies and took far longer to render
my blog with its hundreds of posts than I felt was acceptable. I wanted
a framework to be able to give me rapid feedback while making changes to the
templates, and the 5+-minute render times were just too slow. In general,
they were also very blog minded and didn't have the ability to provide
other content types and flexible URLs.

[Jekyll]: http://jekyllrb.com/
[Middleman]: https://middlemanapp.com/
[Nanoc]: http://nanoc.ws/

I wanted to develop a fast and full-featured website framework without any
dependencies. The [Go language][] seemed to have all the features I needed
in a language. I began developing Hugo in Go and fell in love with the
language. I hope you will enjoy using Hugo (and contributing to it) as much
as I have writing it.

&mdash;Steve Francia (@spf13)

[Go language]: http://golang.org/ "The Go Programming Language"

## Next Steps

 * [Install Hugo](/overview/installing/)
 * [Quick start](/overview/quickstart/)
 * [Join the Mailing List](/community/mailing-list/)
 * [Star us on GitHub](https://github.com/spf13/hugo)
 * [Discussion Forum](http://discuss.gohugo.io/)
