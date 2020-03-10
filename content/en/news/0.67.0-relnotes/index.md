
---
date: 2020-03-09
title: "Hugo 0.67.0: Custom HTTP headers"
description: "This version brings Custom HTTP headers to the development server and exclude/include filters in Hugo Deploy."
categories: ["Releases"]
---

The two main items in Hugo 0.67.0 is custom HTTP header support in `hugo server` and include/exclude filters for [Hugo Deploy](https://gohugo.io/hosting-and-deployment/hugo-deploy/#readout).

Being able to [configure HTTP headers](https://gohugo.io/getting-started/configuration/#configure-server) in your development server means that you can now verify how your site behaves with the intended Content Security Policy settings etc., e.g.:

```toml
[server]
[[server.headers]]
for = "/**.html"

[server.headers.values]
X-Frame-Options = "DENY"
X-XSS-Protection = "1; mode=block"
X-Content-Type-Options = "nosniff"
Referrer-Policy = "strict-origin-when-cross-origin"
Content-Security-Policy = "script-src localhost:1313"
```

**Note:** This release also changes how raw HTML files inside /content is processed to be in line with the documentation. See [#7030](https://github.com/gohugoio/hugo/issues/7030).

This release represents **7 contributions by 4 contributors** to the main Hugo code base.[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@satotake](https://github.com/satotake), [@sams96](https://github.com/sams96), and [@davidejones](https://github.com/davidejones) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) and [@onedrawingperday](https://github.com/onedrawingperday) for their relentless work on keeping the themes site in pristine condition and to [@davidsneighbour](https://github.com/davidsneighbour) and [@kaushalmodi](https://github.com/kaushalmodi) for all the great work on the documentation site.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **5 contributions by 5 contributors**. A special thanks to [@bep](https://github.com/bep), [@psliwka](https://github.com/psliwka), [@digitalcraftsman](https://github.com/digitalcraftsman), and [@jasikpark](https://github.com/jasikpark) for their work on the documentation site.


Hugo now has:

* 42176+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 439+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 301+ [themes](http://themes.gohugo.io/)

## Enhancements

### Other

* Doument the server config [63393230](https://github.com/gohugoio/hugo/commit/63393230c9d3ba19ad182064787e3bfd7ecf82d8) [@bep](https://github.com/bep) 
* Support unComparable args of uniq/complement/in [8279d2e2](https://github.com/gohugoio/hugo/commit/8279d2e2271ee64725133d36a12d1d7e2158bffd) [@satotake](https://github.com/satotake) [#6105](https://github.com/gohugoio/hugo/issues/6105)
* Add HTTP header support for the dev server [10831444](https://github.com/gohugoio/hugo/commit/108314444b510bfc330ccac745dce7beccd52c91) [@bep](https://github.com/bep) [#7031](https://github.com/gohugoio/hugo/issues/7031)
* Add include and exclude support for remote [51e178a6](https://github.com/gohugoio/hugo/commit/51e178a6a28a3f305d89ebb489675743f80862ee) [@davidejones](https://github.com/davidejones) 

## Fixes

### Templates

* Fix error with unicode in file paths [c4fa2f07](https://github.com/gohugoio/hugo/commit/c4fa2f07996c7f1f4e257089a3c3c5b4c1339722) [@sams96](https://github.com/sams96) [#6996](https://github.com/gohugoio/hugo/issues/6996)

### Other

* Fix ambigous error on site.GetPage [6cceef65](https://github.com/gohugoio/hugo/commit/6cceef65c2f4b7c262bf67a249867658112b6de4) [@bep](https://github.com/bep) [#7016](https://github.com/gohugoio/hugo/issues/7016)
* Fix handling of HTML files without front matter [ffcb4aeb](https://github.com/gohugoio/hugo/commit/ffcb4aeb8e392a80da7cad0f1e03a4102efb24ec) [@bep](https://github.com/bep) [#7030](https://github.com/gohugoio/hugo/issues/7030)[#7028](https://github.com/gohugoio/hugo/issues/7028)[#6789](https://github.com/gohugoio/hugo/issues/6789)





