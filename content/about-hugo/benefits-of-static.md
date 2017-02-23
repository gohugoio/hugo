---
title: The Benefits of Static Site Generators
linktitle: The Benefits of Static
description: Learn why static site generators have become such a popular option for developers. Benefits include performance, security, ease of use, and exportability of content.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [ssg,static,performance, security]
weight: 30
draft: false
aliases: []
toc: false
---

Website generators render content into HTML files. Most are "dynamic site generators." That means the HTTP server (i.e., the program that communicates with your end user's browser) runs the generator to create a new HTML file every time an end user requests a page.

Creating the page dynamically requires the HTTP server to have enough memory and CPU to effectively run the generator nonstop. If not, your end user will wait in a queue for the page to be generated.

To prevent unnecessary delays in delivering pages to end users, dynamic site generators were programmed to cache their HTML files. A cached page is a copy that is temporarily stored on the computer. Sending a cached copy is faster than generating a new page at the time of request because the majority of the work is already done.

Hugo and other static site generators take caching a step further. All HTML files are rendered on your computer. You can review the files before you copy them to the computer hosting the HTTP server. Since the HTML files aren't generated dynamically, we say that Hugo is a "static site generator."

Not running a website generator on your HTTP server has many benefits. The most noticeable is performance---HTTP servers are very good at sending files. So good that you can effectively serve the same number of pages with a fraction of the memory and CPU needed for a dynamic site.

## Resources on Static Site Generators

* ["An Introduction to Static Site Generators", David Walsh][]
* ["Static Site Generators", O-Reilly][]
* [StaticGen: Top Open-Source Static Site Generators (GitHub Stars)][]
* ["Top 10 Static Website Generators," Netlify blog][]


["An Introduction to Static Site Generators", David Walsh]: https://davidwalsh.name/introduction-static-site-generators
["Static Site Generators", O-Reilly]: /documents/oreilly-static-site-generators.pdf
["Top 10 Static Website Generators," Netlify blog]: https://www.netlify.com/blog/2016/05/02/top-ten-static-website-generators/
[StaticGen: Top Open-Source Static Site Generators (GitHub Stars)]: https://www.staticgen.com/
