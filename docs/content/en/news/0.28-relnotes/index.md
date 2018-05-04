
---
date: 2017-09-25
title: "Hugo 0.28: High-speed Syntax Highlighting!"
description: "Chroma is the new default syntax highlighter in Hugo."
categories: ["Releases"]
images:
- images/blog/hugo-28-poster.png
---

	Hugo `0.28` brings **blistering fast and native syntax highlighting** from [Chroma](https://github.com/alecthomas/chroma) ([fb33d828](https://github.com/gohugoio/hugo/commit/fb33d8286d78a78a74deb44355b621852a1c4033) [@bep](https://github.com/bep) [#3888](https://github.com/gohugoio/hugo/issues/3888)). A big thank you to [Alec Thomas](https://github.com/alecthomas) for taking on this massive task of porting the popular python highlighter Pygments to Go.

Hugo has been caching the highlighter output to disk, so for repeated builds it has been fine, but this little snippet, showing a build of the [gohugo.io](https://gohugo.io/) site without cache and with both Pygments and Chroma, should illustrate the improvement:

{{< asciicast Lc5iwTVny2kuUC8lqvNnL6oDU >}}

See the [Updated Documentation](https://gohugo.io/content-management/syntax-highlighting/) for more information about how this works.

Worth mentioning is also the `liveReloadPort`  flag on `hugo server`, which makes it possible to do "live reloads" behind a HTTPS proxy, which makes for very cool remote customer demos.

One example would be a Hugo server running behind a [ngrok](https://ngrok.com) tunnel:

```bash
ngrok http 1313
```
Then start the Hugo server with:

```bash
hugo server -b https://youridhere.ngrok.io --appendPort=false --liveReloadPort=443 --navigateToChanged
```

The `navigateToChanged` flag is slightly unrelated, but it is super cool ...

This release represents **15 contributions by 2 contributors** to the main Hugo code base.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs), 
which has received **9 contributions by 7 contributors**. A special thanks to [@bep](https://github.com/bep), [@i-give-up](https://github.com/i-give-up), [@muhajirframe](https://github.com/muhajirframe), and [@icannotfly](https://github.com/icannotfly) for their work on the documentation site.

Hugo now has:

* 19771+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 454+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 180+ [themes](http://themes.gohugo.io/)

## Notes
* Hugo now uses Chroma as new default syntax highlighter. This should in most cases work out-of-the box or with very little adjustments. But if you want to continue to use Pygments, set `pygmentsUseClassic=true` in your site config.
* We now add a set of "no cache" headers to the responses for `hugo server`, which makes the most sense in most development scenarios. Run with `hugo server --noHTTPCache=false` to get the old behaviour.

## Enhancements

### Templates

* Add `urls.Parse` function [81ed5647](https://github.com/gohugoio/hugo/commit/81ed564793609a32be20a569cc15da2cc02dd734) [@moorereason](https://github.com/moorereason) [#3849](https://github.com/gohugoio/hugo/issues/3849)
* Add `math.Ceil`, `Floor`, and `Round` [19c59104](https://github.com/gohugoio/hugo/commit/19c5910485242838d6678c2aacd8501f7e646a53) [@moorereason](https://github.com/moorereason) [#3883](https://github.com/gohugoio/hugo/issues/3883)

### Other

* Use Chroma as new default syntax highlighter [fb33d828](https://github.com/gohugoio/hugo/commit/fb33d8286d78a78a74deb44355b621852a1c4033) [@bep](https://github.com/bep) [#3888](https://github.com/gohugoio/hugo/issues/3888)
* Trim newlines in the hightlight shortcode [0d29a0f7](https://github.com/gohugoio/hugo/commit/0d29a0f7819e8d73149701052c29f090cd6db42b) [@bep](https://github.com/bep) [#3898](https://github.com/gohugoio/hugo/issues/3898)
* Update `goorgeous` [b8fabce2](https://github.com/gohugoio/hugo/commit/b8fabce217fcb52e3f273491bef95c7977058732) [@bep](https://github.com/bep) [#3899](https://github.com/gohugoio/hugo/issues/3899)
* Add `liveReloadPort` flag to server [b1804776](https://github.com/gohugoio/hugo/commit/b180477631555824a06293053e2b6e63c5f07361) [@bep](https://github.com/bep) [#3882](https://github.com/gohugoio/hugo/issues/3882)
* Add `noHTTPCache` flag to hugo server (default on) [0b34af21](https://github.com/gohugoio/hugo/commit/0b34af216154367af7f53ce93d44e6b3d58c3f34) [@bep](https://github.com/bep) [#3897](https://github.com/gohugoio/hugo/issues/3897)
* Make `noHTTPCache` default on [80c7ea60](https://github.com/gohugoio/hugo/commit/80c7ea60a0e0f488563a6b7311f3d4c23457aac7) [@bep](https://github.com/bep) [#3897](https://github.com/gohugoio/hugo/issues/3897)

