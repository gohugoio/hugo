
---
date: 2021-07-21
title: "Hugo 0.86.0: Cascade in Config"
description: "Hugo 0.86.0 adds cascade keyword to site config, much improved \"active menu item\" logic for section pages, and more."
categories: ["Releases"]
---

This release is a set of smaller fixes and improvements. Some of the more notable:

You can now have a top level [cascade](https://gohugo.io/content-management/front-matter#front-matter-cascade) (or one per language, if needed) section in your site configuration (e.g. `config.toml`). This way you can control default front matter values from outside of the content files.

Hugo's [Menu system](https://gohugo.io/content-management/menus/) works well, but hasn't been particularly easy to set the active menu state for section pages without a menu definition. We have had the option [Section Menu for Lazy Bloggers](https://gohugo.io/templates/menu-templates/#section-menu-for-lazy-bloggers). That helped for the common case, but we have now made it behave more sensible out of the box: `$page.HasMenuCurrent $sectionMenuEntry` will now always return true for any descendant of that section. To support this for menu definitions in the site config, we have added a new `pageRef` option on [MenuEntry](https://gohugo.io/variables/menus/#menu-entry-variables):

```toml
[[menus.main]]
pageRef = "blog"
# When pageRef is set, setting `url` is optional; it will be used as a fallback if the page is not found.
url = "/blog"
```

Set `pageRef` to a value that [site.GetPage](https://gohugo.io/functions/getpage/) understands, and the menu entry will be correctly connected to the page.

This release represents **14 contributions by 2 contributors** to the main Hugo code base.
Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs),
which has received **1 contributions by 1 contributors**.

Hugo now has:

* 53005+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 431+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 370+ [themes](http://themes.gohugo.io/)

## Enhancements

### Other

* Simplify "active menu" logic for section menus [d831d2fc](https://github.com/gohugoio/hugo/commit/d831d2fce8198fb814ea4d3d8c311db5c388d04c) [@bep](https://github.com/bep) [#8776](https://github.com/gohugoio/hugo/issues/8776)
* Make keepWhitespace = true default for HTML [c19f65f9](https://github.com/gohugoio/hugo/commit/c19f65f956739ab76c38222d48a3e461525e31af) [@bep](https://github.com/bep) [#8771](https://github.com/gohugoio/hugo/issues/8771)
* Make FileMeta a struct [022c4795](https://github.com/gohugoio/hugo/commit/022c4795510306e08a4aba31504ca382d41c7fac) [@bep](https://github.com/bep) [#8749](https://github.com/gohugoio/hugo/issues/8749)
* Add tabindex="0" to default <pre> wrapper [f27e5424](https://github.com/gohugoio/hugo/commit/f27e542442d19436f1428cc22bb03aca398d37a7) [@rhymes](https://github.com/rhymes) [#7194](https://github.com/gohugoio/hugo/issues/7194)
* Rename/reorder the hook methods [80566481](https://github.com/gohugoio/hugo/commit/805664818d0e1f95a3474271c2db3e5f49db26ba) [@bep](https://github.com/bep) [#8755](https://github.com/gohugoio/hugo/issues/8755)
* Support auto links in render hook [ee3d2bb1](https://github.com/gohugoio/hugo/commit/ee3d2bb1d3974584f47cde7c973fbd1ae1f512b6) [@bep](https://github.com/bep) [#8755](https://github.com/gohugoio/hugo/issues/8755)
* Adjust a test helper [eb2a5003](https://github.com/gohugoio/hugo/commit/eb2a500367780b07d67c301ce7c866e6b67aa687) [@bep](https://github.com/bep) 
* Add config.cascade [5cb52c23](https://github.com/gohugoio/hugo/commit/5cb52c23150032b3fdb211a095745c512369b463) [@bep](https://github.com/bep) [#8741](https://github.com/gohugoio/hugo/issues/8741)
* Regenerate image golden testdata [30eea391](https://github.com/gohugoio/hugo/commit/30eea3915b67f72611a3b2f4547146d4c6a96864) [@bep](https://github.com/bep) [#8729](https://github.com/gohugoio/hugo/issues/8729)

## Fixes

### Other

* Fix panic on invalid config in "hugo mod get" and similar [351ed0f5](https://github.com/gohugoio/hugo/commit/351ed0f569f96aff29b03925bf5154d80a164e00) [@bep](https://github.com/bep) [#8773](https://github.com/gohugoio/hugo/issues/8773)
* Fix Params case handling for menu items defined in site config [634481ba](https://github.com/gohugoio/hugo/commit/634481ba8cfcd865ba0d8811d8834f6af45663d7) [@bep](https://github.com/bep) [#8775](https://github.com/gohugoio/hugo/issues/8775)
* Fix default values when loading from config dir [ae6cf93c](https://github.com/gohugoio/hugo/commit/ae6cf93c84c3584b111f4b9fa3fb4e3f63d37915) [@bep](https://github.com/bep) [#8763](https://github.com/gohugoio/hugo/issues/8763)
* Fix the deprecation error/warn log levels [a70da2b7](https://github.com/gohugoio/hugo/commit/a70da2b74a6af0834cce9668cdb6acdb1c86a4c0) [@bep](https://github.com/bep) 
* Fix transparency problem when converting 32-bit images to WebP [8f40f34c](https://github.com/gohugoio/hugo/commit/8f40f34cd10a98598bb822ec633fd5d0ea64b612) [@bep](https://github.com/bep) [#8729](https://github.com/gohugoio/hugo/issues/8729)
