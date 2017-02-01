---
title: Introduction to Hugo
linktitle:
description: An introduction to Hugo and the impetus for its creation.
qr_description:
qr_returns:
date: 2016-11-01
publishdate: 2016-11-01
lastmod: 2016-11-01
weight: 10
draft: false
type:
layout:
slug:
aliases: [/overview/introduction]
notes:
---

## What is Hugo?

Hugo is a general-purpose website framework. Technically speaking, Hugo is
a static site generator. Unlike systems that dynamically build a page
with each visitor request, Hugo builds pages when you create
your content. Since websites are viewed far more often than they are
edited, Hugo is optimized for website viewing for the end users of your site and an ideal writing experience for authors.

Sites built with Hugo are extremely fast and secure. Hugo sites can
be hosted anywhere, including [Heroku][], [GoDaddy][], [DreamHost][],
[GitHub Pages][], [Surge][], [Aerobatic][], [Firebase Hosting][],
[Google Cloud Storage][], [Amazon S3][], and [CloudFront][], and work well
with CDNs. Hugo sites run without the need for a database or dependencies on expensive runtimes like Ruby, Python, or PHP.

We think of Hugo as the ideal website creation tool. With nearly instant
build times and the ability to rebuild whenever a change is made, Hugo
provides a very fast feedback loop. This is essential when you are
designing websites, but also very useful when creating content.

## Why a Static Site Generator?

Website generators render content into HTML files. Most are "dynamic site generators." That means the HTTP server (ie, the program that communicates with your end user's browser) runs the generator to create a new HTML file every time an end user requests a page.

Creating the page dynamically requires the HTTP server have enough memory and CPU to effectively run the generator around the clock. If not, your end user will wait in a queue for the page to be generated.

To prevent unnecessary delays in delivering pages to end users, dynamic site generators programmed their systems to cache the HTML files. A cached page is a copy that is temporarily stored on the computer. Sending a cached copy is faster than generating a new page at the time of request because the majority of the work is already done.

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


## What does Hugo do?

In technical terms, Hugo takes a source directory of files and
templates and uses these as input to create a complete website.

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

It rendered slowly and required me to be online to write posts. Plus, its constant security updates and the horror stories of people's hacked blogs! I hated how content was only written in HTML instead of the much simpler Markdown. Overall, I felt like WordPress hindered more than helped. It kept me from writing great content.

I looked at existing static site generators like [Jekyll][], [Middleman][] and [Nanoc][]. All had complicated installation dependencies and an unacceptably long time to render my blog with hundreds of posts. I wanted a framework that would give me rapid feedback while making changes to the templates, and the 5-minute-plus render times were just too slow. In general, the static site generators were very blog-minded and didn't provide for other content types and flexible URLs.

I wanted to develop a fast and full-featured website framework without any
dependencies. The [Go language][] seemed to have all the features I needed
in a language. I began developing Hugo in Go and fell in love with the
language. I hope you will enjoy using Hugo (and contributing to it) as much
as I have writing it.

&mdash;Steve Francia ([@spf13][])

[@spf13]: https://twitter.com/spf13
[Aerobatic]: https://www.aerobatic.com/
[Amazon S3]: http://aws.amazon.com/s3/
[CloudFront]: http://aws.amazon.com/cloudfront/ "Amazon CloudFront"
[DreamHost]: http://www.dreamhost.com/
[Firebase Hosting]: https://firebase.google.com/docs/hosting/
[GitHub Pages]: https://pages.github.com/
[GitLab]: https://about.gitlab.com
[Go language]: http://golang.org/ "The Go Programming Language"
[GoDaddy]: https://www.godaddy.com/
[Google Cloud Storage]: http://cloud.google.com/storage/
[Heroku]: https://www.heroku.com/
[Jekyll]: http://jekyllrb.com/
[Middleman]: https://middlemanapp.com/
[Nanoc]: http://nanoc.ws/
[Surge]: https://surge.sh

## Next Steps

 * [Install Hugo](/overview/installing/)
 * [Quick start](/overview/quickstart/)
 * [Join the Mailing List](/community/mailing-list/)
 * [Star us on GitHub](https://github.com/spf13/hugo)
 * [Discussion Forum](http://discuss.gohugo.io/)


