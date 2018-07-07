
---
date: 2018-02-15
title: "Hugo 0.36.1: One Bugfix"
description: "Fixes a multi-thread image processing issue."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---
	This release fixes a multi-thread issue when reprocessing and reusing images across pages. When doing something like this with the same image from a partial used in, say, both the home page and the single page:

```bash
{{ with $img }}
{{ $big := .Fill "1024x512 top" }}
{{ $small := $big.Resize "512x" }}
{{ end }}
```

There would be timing issues making Hugo in some cases trying to process the same image twice at the same time.

You would experience errors of type:

```bash
png: invalid format: not enough pixel data
```

This commit fixes that by adding a mutex per image. This should also improve the performance, slightly, as it avoids duplicate work.

The current workaround before this fix is to always operate on the original:

```bash
{{ with $img }}
{{ $big := .Fill "1024x512 top" }}
{{ $small := .Fill "512x256 top" }}
{{ end }}
```
This error was rare (no reports on GitHub or the discussion forum), but very hard to debug for the end user.

* Fix multi-threaded image processing issue [d8fdffb5](https://github.com/gohugoio/hugo/commit/d8fdffb55268464d54558d6f9cd3874b612dc7c7) [@bep](https://github.com/bep) [#4404](https://github.com/gohugoio/hugo/issues/4404)
* Improve error message in .Render [08521dac](https://github.com/gohugoio/hugo/commit/08521dac8323403933a8fd11acfd16930af5f17d) [@bep](https://github.com/bep) 
* Bump Travis/Snapcraft to Go 1.9.4 [fc23a80f](https://github.com/gohugoio/hugo/commit/fc23a80ffd3878b9ba9a160ce37e0e1d8703faf3) [@bep](https://github.com/bep) 
* Improve error processing error message [2851af02](https://github.com/gohugoio/hugo/commit/2851af0225cdf6c4e47058979cd22949ed6d1fc0) [@bep](https://github.com/bep) 
