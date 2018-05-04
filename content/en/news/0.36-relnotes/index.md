
---
date: 2018-02-05
title: "Hugo 0.36: Smart Image Cropping!"
description: "Hugo 0.36 announces smart image cropping and some important bug fixes."
categories: ["Releases"]
---

Hugo `0.36` announces **smart cropping** of images, using the [library](https://github.com/muesli/smartcrop) created by [muesli](https://github.com/muesli). We will work with him to improve this even more in the future, but this is now the default used when cropping images in Hugo.

Go [here](http://hugotest.bep.is/resourcemeta/smartcrop/) for a list of examples.

This release represents **7 contributions by 3 contributors** to the main Hugo code base.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **9 contributions by 4 contributors**. A special thanks to [@bep](https://github.com/bep), [@Jibec](https://github.com/Jibec), [@Nick-Rivera](https://github.com/Nick-Rivera), and [@kaushalmodi](https://github.com/kaushalmodi) for their work on the documentation site.


Hugo now has:

* 23100+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 448+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 197+ [themes](http://themes.gohugo.io/)

## Notes
Hugo now defaults to **smart crop** when cropping images, if you don't specify it when calling `.Fill`.

You can get the old default by adding this to your `config.toml`:

```toml
[imaging]
anchor = "center"
```
Also, we have removed the superflous anchor name from the processed filenames that does not use this anchor, so it can be wise to run `hugo --gc` once to remove unused images.

## Enhancements
* Add smart cropping [722086b4](https://github.com/gohugoio/hugo/commit/722086b4ed3e77d1aba6724474bec06d08e7de06) [@bep](https://github.com/bep) [#4375](https://github.com/gohugoio/hugo/issues/4375)

## Fixes
* Ensure site templates can override theme templates [084cf419](https://github.com/gohugoio/hugo/commit/084cf4191b3c1e7590a4223fd9251019ef5d4c21) [@moorereason](https://github.com/moorereason) [#3505](https://github.com/gohugoio/hugo/issues/3505)
* Add additional test to `TestTemplateLookupOrder` [fc06d5c1](https://github.com/gohugoio/hugo/commit/fc06d5c18bb1e47f90f0297aa8121ee0775e047d) [@moorereason](https://github.com/moorereason) [#3505](https://github.com/gohugoio/hugo/issues/3505)
* Fix broken `TestTemplateLookupOrder` [9a367d9d](https://github.com/gohugoio/hugo/commit/9a367d9d06db6f6cf22121d0397c464ae36e7089) [@moorereason](https://github.com/moorereason) 
* Fix JSON array-based data file handling regression [4402c077](https://github.com/gohugoio/hugo/commit/4402c077754991df19c3bbab0c4a671dcfdc192c) [@vassudanagunta](https://github.com/vassudanagunta) [#4361](https://github.com/gohugoio/hugo/issues/4361)
* Increase data directory test coverage [4743de0d](https://github.com/gohugoio/hugo/commit/4743de0d3c7564fc06972074e903d5502d204353) [@vassudanagunta](https://github.com/vassudanagunta) [#4138](https://github.com/gohugoio/hugo/issues/4138)






