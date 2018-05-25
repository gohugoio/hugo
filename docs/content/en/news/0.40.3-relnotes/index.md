
---
date: 2018-05-09
title: "Hugo 0.40.3: One Bug Fix"
description: "Fixes a rare, but possible Content truncation issue."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

	Hugo `0.40.3` fixes a possible `.Content` truncation issue introduced in `0.40.1` [90d0d830](https://github.com/gohugoio/hugo/commit/90d0d83097a20a3f521ffc1f5a54a2fbfaf14ce2) [@bep](https://github.com/bep) [#4706](https://github.com/gohugoio/hugo/issues/4706). This should be very rare. It has been reported by only one user on a synthetic site. We have tested a number of big sites that does not show this problem with `0.40.2`, but this is serious enough to warrant a patch release.
