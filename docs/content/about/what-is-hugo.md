---
title: What is Hugo
linktitle: What is Hugo
description: Hugo is a fast and modern static site generator written in Go and designed to make website creation fun again.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
layout: single
menu:
  main:
    parent: "About Hugo"
    weight: 10
weight: 10
sections_weight: 10
draft: false
aliases: [/overview/introduction/,/about/why-i-built-hugo/]
toc: true
---

Hugo is a general-purpose website framework. Technically speaking, Hugo is a [static site generator][]. Unlike systems that dynamically build a page with each visitor request, Hugo builds pages when you create or update your content. Since websites are viewed far more often than they are edited, Hugo is designed to provide an optimal viewing experience for your website's end users and an ideal writing experience for website authors.

Websites built with Hugo are extremely fast and secure. Hugo sites can be hosted anywhere, including [Netlify][], [Heroku][], [GoDaddy][], [DreamHost][], [GitHub Pages][], [Surge][], [Aerobatic][], [Firebase][], [Google Cloud Storage][], [Amazon S3][], [Rackspace][], [Azure][], and [CloudFront][] and work well with CDNs. Hugo sites run without the need for a database or dependencies on expensive runtimes like Ruby, Python, or PHP.

We think of Hugo as the ideal website creation tool. Hugo provides nearly instant build times and the ability to rebuild whenever a change is made, which is invaluable when you are designing websites and creating content.

## How Fast is Hugo?

{{< youtube "CdiDYZ51a2o" >}}

## What Does Hugo Do?

In technical terms, Hugo takes a source directory of files and templates and uses these as input to create a complete website.

## Who Should Use Hugo?

Hugo is for people that prefer writing in a text editor over a browser.

Hugo is for people who want to hand code their own website without worrying about setting up complicated runtimes, dependencies and databases.

Hugo is for people building a blog, a company site, a portfolio site, documentation, a single landing page, or a website with thousands of pages.

## Why I Built Hugo

I wrote Hugo ultimately for a few reasons. First, I was disappointed with WordPress, my then website solution. With it, I couldn't create content as efficiently as I wanted to.

WordPress sites rendered slowly and required I be online to write posts. WordPress is known for its constant security updates and horror stories of hacked blogs. I hated writing content in HTML instead of the much simpler Markdown. Overall, I felt like WordPress hindered more than helped. It kept me from writing great content.

I looked at existing static site generators like [Jekyll][], [Middleman][], and [Nanoc][]. All had complicated installation dependencies and an unacceptably long time to render my blog with hundreds of posts. I wanted a framework that would give me rapid feedback while making changes to the templates, and the 5-minute-plus render times were just too slow. In general, the static site generators were very blog-minded and didn't provide for other content types and flexible URLs.

I wanted to develop a fast and full-featured website framework without any dependencies. The [Go language][] seemed to have all the features I needed. I began developing Hugo in Go and fell in love with the language. I hope you will enjoy using Hugo (and [contributing to it][]) as much as I have writing it.

&#8213;Steve Francia ([@spf13][])


[@spf13]: https://twitter.com/@spf13
[Aerobatic]: https://www.aerobatic.com/
[Amazon S3]: http://aws.amazon.com/s3/
[Azure]: https://blogs.msdn.microsoft.com/acoat/2016/01/28/publish-a-static-web-site-using-azure-web-apps/
[CloudFront]: http://aws.amazon.com/cloudfront/ "Amazon CloudFront"
[contributing to it]: https://github.com/spf13/hugo
[DreamHost]: http://www.dreamhost.com/
[Firebase]: https://firebase.google.com/docs/hosting/ "Firebase static hosting"
[GitHub Pages]: https://pages.github.com/
[GitLab]: https://about.gitlab.com
[Go language]: https://golang.org/
[GoDaddy]: https://www.godaddy.com/ "Godaddy.com Hosting"
[Google Cloud Storage]: http://cloud.google.com/storage/
[Heroku]: https://www.heroku.com/
[Jekyll]: http://jekyllrb.com/
[Jekyll]: https://jekyllrb.com/
[Middleman]: https://middlemanapp.com/
[Middleman]: https://middlemanapp.com/
[Nanoc]: http://nanoc.ws/
[Nanoc]: https://nanoc.ws/
[Netlify]: https://netlify.com
[rackspace]: https://www.rackspace.com/cloud/files
[static site generator]: /about/benefits/
[Rackspace]: https://www.rackspace.com/cloud/files
[static site generator]: /about/benefits/
[Surge]: https://surge.sh