---
title: 搜索你的 Hugo 网站
linktitle: 搜索
description: 为你最新添加的 Hugo 网站看一些开源和商业搜索选择。
date: 2018-05-29
publishdate: 2018-05-29
lastmod: 2018-05-29
categories: [developer tools]
keywords: [search,tools]
menu:
  docs:
    parent: "tools"
    weight: 60
weight: 60
sections_weight: 60
draft: false
aliases: []
toc: true
---

静态网站还有动态搜索功能？是的。作为可选方案，来自 Google 或者其它搜索引擎的嵌入式脚本，可以给你的访客提供一个自定义的直接基于你的内容文件索引的搜素。

* [GitHub Gist for Hugo Workflow](https://gist.github.com/sebz/efddfc8fdcb6b480f567). Gist 包含一个为你的网站创建索引的简单流程。它使用简单的 Grunt 脚本索引你所有的内容文件，并且 [lunr.js](http://lunrjs.com/) 会提供搜索结果。
* [hugo-lunr](https://www.npmjs.com/package/hugo-lunr). 一个使用 [lunr.js](http://lunrjs.com/) 为你的 Hugo 静态网站添加搜索的简单方法。Hugo-lunr 将会给你的 Hugo 项目中任意 html 和 markdown 文件创建一个索引文件。
* [hugo-lunr-zh](https://www.npmjs.com/package/hugo-lunr-zh). 有点像 Hugo-lunr，但是 Hugo-lunr-zh 能帮助你分割中文关键字。
* [Github Gist for Fuse.js integration](https://gist.github.com/eddiewebb/735feb48f50f0ddd65ae5606a1cb41ae). 该 gist 显示如何借助 Hugo 已有的构建时间在客户端通过 [Fuse.js](http://fusejs.io/) 生成可搜索的 JSON 索引。尽管 gist 使用 Fuse.js 进行模糊匹配，任何有能力读取 JSON 索引的客户端工具都可以运行。不需要 Hugo 以外， npm、grunt 或者其它构建时工具。
* [hugo-search-index](https://www.npmjs.com/package/hugo-search-index). 一个包含实现了搜索的 Gulp 任务和内置浏览器脚本的库。Gulp 从项目中的 markdown 文件生成索引。

## 商业搜索服务

* [Algolia](https://www.algolia.com/) 的搜索 API 使得在你的应用和网站中提供一个很好的搜索体验变得简单。Algolia 搜索提供托管的全文本、数字化、分面以及地理定位搜索。
