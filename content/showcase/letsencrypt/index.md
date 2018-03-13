---
title: "Let’s Encrypt"
date: 2018-03-13
description: "Showcase: Lessons learned from taking letsencrypt.org to Hugo."
siteURL: https://letsencrypt.org/
siteSource: https://github.com/letsencrypt/website
byline: "[bep](https://github.com/bep), Hugo Lead"
---

The **Let’s Encrypt website** has a common set of elements: A landing page and some other static info-pages, a document section, a blog, and a documentation section. Having it moved to Hugo was mostly motivated by a _simpler administration and Hugo's [multilingual support](/content-management/multilingual/)_. They already serve HTTPS to more than 60 million domains, and having the documentation available in more languages will increase that reach.[^1]

{{< tweet 971755920639307777 >}}

I helped them port the site from Jekyll to Hugo. There are usually very few surprises doing this. I know Hugo very well, but working on sites with a history usually comes up with something new.

That site is bookmarked in many browsers, so preserving the URLs was a must. Hugo's URL handling is very flexible, but there was one challenge. The website has a mix of standard and what we in Hugo call _ugly URLs_ (`https://letsencrypt.org/2017/12/07/looking-forward-to-2018.html`). In Hugo this is handled automatically, and you can turn it on globally or per language. But before Hugo `0.33` you could not configure it for parts of your site. You could set it manually for the relevant pages in front matter -- which is how it was done in Jekyll -- but that would be hard to manage, especially when you start to introduce translations. So, in [Hugo 0.33](https://gohugo.io/news/0.33-relnotes/) I added support for _ugly URLs_ per section and also `url` set in front matter for list pages (`https://letsencrypt.org/blog/`).

The lessons learned from this also lead to [disableLanguages](/content-management/multilingual/#disable-a-language) in Hugo `0.34` (a way to turn off languages during translation). And I also registered [this issue](https://github.com/gohugoio/hugo/issues/4463). Once fixed it will make it easier to handle partially translated sites.


[^1]: The work on getting the content translated is in progress.
