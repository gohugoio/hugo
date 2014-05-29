---
date: 2014-05-26
linktitle: Comments
menu:
  main:
    parent: extras
next: /extras/livereload
prev: /extras/builders
title: Comments in Hugo
weight: 14
---

As Hugo is a static site generator, the content produced is static and
doesn’t interact with the users. The most common interaction people ask
for is comment capability.

Hugo ships with support for [disqus](http://disqus.com), a third party
service that provides comment and community capabilities to website via
javascript.

Your theme may already support disqus, but even it if doesn’t it is easy
to add.

# Disqus Support

## Adding Disqus to a template

Hugo comes with all the code you would need to include load disqus.
Simply include the following line where you want your comments to appear

    {{ template "_internal/disqus.html" . }}


## Configuring Disqus

That template requires you to set a single value in your site config file, eg. config.yaml.

    disqusShortname = "XYW"

Additionally you can optionally set the following in the front matter
for a given piece of content

 * **disqus_identifier**
 * **disqus_title**
 * **disqus_url**

# Alternatives

A few alternatives exist to Disqus.

* [Intense Debate](http://intensedebate.com/)
* [LiveFyre](http://livefyre.com/)
* [Moot](http://muut.com)
* [Kaiju](http://github.com/spf13/kaiju)


[Kaiju](http://github.com/spf13/kaiju) is a open source project started
by [spf13](http://spf13.com) (Hugo’s author) to bring easy and fast real
time discussions to the web.

Written using Go, Socket.io and MongoDB it is very fast and easy to
deploy.

It is in early development but shows promise.. If you have interest
please help by contributing whether via a pull request, an issue or even
just a tweet. Everything helps.

