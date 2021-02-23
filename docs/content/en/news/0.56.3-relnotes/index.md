
---
date: 2019-07-31
title: "Hugo 0.56.3: A couple of Bug Fixes"
description: "This version fixes a couple of bugs introduced in 0.56.0."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

This is a bug-fix release with a couple of important fixes. After getting feedback about the new **Hugo Modules** feature, this release also adds some minor improvements:

It adds support for overlapping file mounts, even for the filesystems where we walk down the directory structure. One relevant example that is fixed by this release:

```toml
[module]
[[module.mounts]]
source="content1"
target="content"
[[module.mounts]]
source="content2"
target="content/docs"
```

The above is obviously both common and very useful. This was never an issue with the situations where you load a specific file/directory (e.g. `resources.Get "a/b/c/d/sunset.jpg"`).

User feedback also told us that these file mounts were a little hard to debug, so we added a new command that prints the configured mounts as a JSON:

```bash
hugo config mounts
```

* hugolib: Fix bundle header clone logic [0e086785](https://github.com/gohugoio/hugo/commit/0e086785fa4be8086256e9d7de6cda78e18d00ee) [@bep](https://github.com/bep) [#6136](https://github.com/gohugoio/hugo/issues/6136)
* docs: Regenerate CLI docs [02b947ea](https://github.com/gohugoio/hugo/commit/02b947eaa3cc68404180d796a2f7119dce074539) [@bep](https://github.com/bep) 
* commands: Add "hugo config mounts" command [d7c233af](https://github.com/gohugoio/hugo/commit/d7c233afee6a16b1947f60b7e5450e40612997bb) [@bep](https://github.com/bep) [#6144](https://github.com/gohugoio/hugo/issues/6144)
* commands: Cleanup the hugo config command [45ee8a7a](https://github.com/gohugoio/hugo/commit/45ee8a7a52213bf394c7f41a72be78084ddc789a) [@bep](https://github.com/bep) [#6144](https://github.com/gohugoio/hugo/issues/6144)
* Move the mount duplicate filter to the modules package [4b6c5eba](https://github.com/gohugoio/hugo/commit/4b6c5eba306e6e69f3dd07a6c102bfc8040b38c9) [@bep](https://github.com/bep) 
* Allow overlap in module mounts [edf9f0a3](https://github.com/gohugoio/hugo/commit/edf9f0a354e5eaa556f8faed70b5243b7273b35c) [@bep](https://github.com/bep) [#6146](https://github.com/gohugoio/hugo/issues/6146)
* Fix self-mounts on the main project [36220851](https://github.com/gohugoio/hugo/commit/36220851e4ed7fc3fa78aa250d001d5f922210e7) [@bep](https://github.com/bep) [#6143](https://github.com/gohugoio/hugo/issues/6143)







