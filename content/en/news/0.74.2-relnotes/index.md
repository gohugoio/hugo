
---
date: 2020-07-17
title: "Hugo 0.74.2: A couple of Bug Fixes"
description: "This version fixes a couple of bugs introduced in 0.74.0."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

	

Add .Defines to js.Build options [35011bcb](https://github.com/gohugoio/hugo/commit/35011bcb26b6fcfcbd77dc05aa8246ca45b2c2ba) [@bep](https://github.com/bep) [#7489](https://github.com/gohugoio/hugo/issues/7489)

This is needed to import `react` as a library, e.g.:

```
{{ $jsx := resources.Get "index.jsx" }}
{{ $options := dict "defines" (dict "process.env.NODE_ENV" "\"development\"") }}
{{ $js := $jsx | js.Build $options }}
```


