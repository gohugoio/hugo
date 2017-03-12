---
title: Custom 404 Page
linktitle: 404 Page
description: If you know how to create a single page template, you have unlimited options for creating a custom 404.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [404, page not found]
weight: 120
draft: false
aliases: [/templates/404/]
toc: false
wip: true
---

When using Hugo with [GitHub Pages](http://pages.github.com/), you can provide your own template for a [custom 404 error page](https://help.github.com/articles/custom-404-pages/) by creating a 404.html template file in your `/layouts` folder. When Hugo generates your site, the `404.html` file will be placed in the root.

404 pages will have all the regular [page variables][pagevars] available to use in the templates.

In addition to the standard page variables, the 404 page has access to all site content accessible from `.Data.Pages`.

```bash
▾ layouts/
    404.html
```

## 404.html

This is a basic example of a 404.html template:

{{% code file="404.html"%}}
```html
{{ partial "header.html" . }}
{{ partial "subheader.html" . }}

<section id="main">
  <div>
   <h1 id="title">{{ .Title }}</h1>
  </div>
</section>

{{ partial "footer.html" . }}
```
{{% /code %}}

## Automatic Loading

Your 404.html file can be set to load automatically when a visitor enters a mistaken URL path, dependent upon the web serving environment you are using. For example:

* _GitHub Pages_ - it's automatic.
* _Apache_ - one way is to specify `ErrorDocument 404 /404.html` in an `.htaccess` file in the root of your site.
* _Nginx_ - you might specify `error_page   404  =  /404.html;` in your `nginx.conf` file.
* _Amazon AWS S3_ - when setting a bucket up for static web serving, you can specify the error file.
* _Caddy Server_ - using `errors { 404 /404.html }`. [Details here](https://caddyserver.com/docs/errors)

[pagevars]: /variables/page/