---
title: Hugo Modules
linktitle: Hugo Modules Overview
description: How to use Hugo Modules.
date: 2017-02-01
publishdate: 2017-02-01
menu:
  docs:
    parent: "modules"
    weight: 01
weight: 01
sections_weight: 01
categories: [hugo modules]
keywords: [themes,modules]
draft: false
aliases: [/themes/overview/,/themes/]
toc: true
---

**Hugo Modules** are the core building blocks in Hugo. A _module_ can be your main project or a smaller module providing one or more of the 7 component types defined in Hugo: **static**, **content**, **layouts**, **data**, **assets**, **i18n**, and **archetypes**.

You can combine modules in any combination you like, and even mount directories from non-Hugo projects, forming a big, virtual union file system.

Hugo Modules is powered by Go Modules. For more information about Go Modules, see:

- [https://github.com/golang/go/wiki/Modules](https://github.com/golang/go/wiki/Modules)
- [https://blog.golang.org/using-go-modules](https://blog.golang.org/using-go-modules)

This is all very much brand new and there are only a few example projects around:

- [https://github.com/bep/docuapi](https://github.com/bep/docuapi) is a theme that has been ported to Hugo Modules while testing this feature. It is a good example of a non-Hugo-project mounted into Hugo’s folder structure. It even shows a JS Bundler implementation in regular Go templates.
- [https://github.com/bep/my-modular-site](https://github.com/bep/my-modular-site) is a very simple site used for testing.
