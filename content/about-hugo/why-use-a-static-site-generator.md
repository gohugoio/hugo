---
title: Why Use a Static Site Generator?
linktitle:
description:
date: 2016-02-01
publishdate: 2016-02-01
lastmod: 2016-02-01
weight: 20
draft: false
slug:
aliases: []
notes:
---

Website generators render content into HTML files. Most are "dynamic site generators." That means the HTTP server (i.e., the program that communicates with your end user's browser) runs the generator to create a new HTML file every time an end user requests a page.

Creating the page dynamically requires the HTTP server to have enough memory and CPU to effectively run the generator nonstop. If not, your end user will wait in a queue for the page to be generated.

To prevent unnecessary delays in delivering pages to end users, dynamic site generators were programmed to cache their HTML files. A cached page is a copy that is temporarily stored on the computer. Sending a cached copy is faster than generating a new page at the time of request because the majority of the work is already done.

Hugo and other static site generators take caching a step further. All HTML files are rendered on your computer. You can review the files before you copy them to the computer hosting the HTTP server. Since the HTML files aren't generated dynamically, we say that Hugo is a "static site generator."

Not running a web site generator on your HTTP server has many benefits. The most noticeable is performance---HTTP servers are very good at
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

## Further Reading

* ["Static Site Generators", O-Reilly][]
* [StaticGen: Top Open-Source Static Site Generators][]
* ["Top 10 Static Website Generators," Netlify blog][]



["Static Site Generators", O-Reilly]: /documents/oreilly-static-site-generators.pdf
["Top 10 Static Website Generators," Netlify blog]: https://www.netlify.com/blog/2016/05/02/top-ten-static-website-generators/
[StaticGen: Top Open-Source Static Site Generators]: https://www.staticgen.com/
