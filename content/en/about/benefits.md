---
title: The Benefits of Static Site Generators
linktitle: The Benefits of Static
description: Improved performance, security and ease of use are just a few of the reasons static site generators are so appealing.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
keywords: [ssg,static,performance,security]
menu:
  docs:
    parent: "about"
    weight: 30
weight: 30
sections_weight: 30
draft: false
aliases: []
toc: false
---

The purpose of website generators is to render content into HTML files. Most are "dynamic site generators." That means the HTTP server---i.e., the program that sends files to the browser to be viewed---runs the generator to create a new HTML file every time an end user requests a page.

Over time, dynamic site generators were programmed to cache their HTML files to prevent unnecessary delays in delivering pages to end users. A cached page is a static version of a web page.

Hugo takes caching a step further and all HTML files are rendered on your computer. You can review the files locally before copying them to the computer hosting the HTTP server. Since the HTML files aren't generated dynamically, we say that Hugo is a *static site generator*.

This has many benefits. The most noticeable is performance. HTTP servers are *very* good at sending files---so good, in fact, that you can effectively serve the same number of pages with a fraction of the memory and CPU needed for a dynamic site.

## More on Static Site Generators

* ["An Introduction to Static Site Generators", David Walsh][]
* ["Hugo vs. Wordpress page load speed comparison: Hugo leaves WordPress in its dust", GettingThingsTech][hugovwordpress]
* ["Static Site Generators", O'Reilly][]
* [StaticGen: Top Open-Source Static Site Generators (GitHub Stars)][]
* ["Top 10 Static Website Generators", Netlify blog][]
* ["The Resurgence of Static", dotCMS][dotcms]


["An Introduction to Static Site Generators", David Walsh]: https://davidwalsh.name/introduction-static-site-generators
["Static Site Generators", O'Reilly]: https://www.oreilly.com/web-platform/free/files/static-site-generators.pdf
["Top 10 Static Website Generators", Netlify blog]: https://www.netlify.com/blog/2016/05/02/top-ten-static-website-generators/
[hugovwordpress]: https://gettingthingstech.com/hugo-vs.-wordpress-page-load-speed-comparison-hugo-leaves-wordpress-in-its-dust/
[StaticGen: Top Open-Source Static Site Generators (GitHub Stars)]: https://www.staticgen.com/
[dotcms]: https://dotcms.com/blog/post/the-resurgence-of-static
