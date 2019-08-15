
---
date: 2019-08-15
title: "Hugo 0.57.1: A couple of Bug Fixes"
description: "This version fixes a couple of bugs introduced in 0.57.0."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

	

This is a bug-fix release with a couple of important fixes.

* hugolib: Fix draft etc. handling of _index.md pages [6ccf50ea](https://github.com/gohugoio/hugo/commit/6ccf50ea7bb291bcbe1d56a4d697a6fd57a9c629) [@bep](https://github.com/bep) [#6222](https://github.com/gohugoio/hugo/issues/6222)[#6210](https://github.com/gohugoio/hugo/issues/6210)
* Fix mainSections logic [67524c99](https://github.com/gohugoio/hugo/commit/67524c993623871626f0f22e6a2ac705a816a959) [@bep](https://github.com/bep) [#6217](https://github.com/gohugoio/hugo/issues/6217)
* Fix live reload mount logic with sub paths [952a3194](https://github.com/gohugoio/hugo/commit/952a3194962dd91f87e5bd227a1591b00c39ff05) [@bep](https://github.com/bep) [#6209](https://github.com/gohugoio/hugo/issues/6209)
* modules: Disable "auto tidy" for now [321418f2](https://github.com/gohugoio/hugo/commit/321418f22a4a94b87f01e1403a2f4a71106461fb) [@bep](https://github.com/bep) [#6115](https://github.com/gohugoio/hugo/issues/6115)
* hugolib: Recover and log panics in content init [7f3aab5a](https://github.com/gohugoio/hugo/commit/7f3aab5ac283ecfc7029b680d4c0a34920e728c8) [@bep](https://github.com/bep) [#6210](https://github.com/gohugoio/hugo/issues/6210)
* hugolib: Add some outputs tests [028b9926](https://github.com/gohugoio/hugo/commit/028b992611209b241b1f55def8d47f9188038dc3) [@bep](https://github.com/bep) [#6210](https://github.com/gohugoio/hugo/issues/6210)
* hugolib: Fix taxonomies vs expired [9475f61a](https://github.com/gohugoio/hugo/commit/9475f61a377fcf23f910cbfd4ddca59261326665) [@bep](https://github.com/bep) [#6213](https://github.com/gohugoio/hugo/issues/6213)
* commands: Make sure the hugo field is always initialized before it's used [ea9261e8](https://github.com/gohugoio/hugo/commit/ea9261e856c13c1d4ae05fcca08766d410b4b65c) [@vazrupe](https://github.com/vazrupe) [#6193](https://github.com/gohugoio/hugo/issues/6193)



