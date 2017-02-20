---
title: Why I Built Hugo
linktitle:
description: Learn why Steve Francia (@spf13) decided to ditch his previous content publishing workflow and develop Hugo in Golang.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
weight: 30
draft: false
slug:
aliases: []
toc: false
notesforauthors:
---

I wrote Hugo ultimately for a few reasons. First, I was disappointed with WordPress, my then website solution. With it, I couldn't create content as efficiently as I wanted to.

WordPress sites rendered slowly and required I be online to write posts. WordPress is known for its constant security updates and horror stories of hacked blogs. I hated writing content in HTML instead of the much simpler Markdown. Overall, I felt like WordPress hindered more than helped. It kept me from writing great content.

I looked at existing static site generators like [Jekyll][], [Middleman][], and [Nanoc][]. All had complicated installation dependencies and an unacceptably long time to render my blog with hundreds of posts. I wanted a framework that would give me rapid feedback while making changes to the templates, and the 5-minute-plus render times were just too slow. In general, the static site generators were very blog-minded and didn't provide for other content types and flexible URLs.

I wanted to develop a fast and full-featured website framework without any dependencies. The [Go language][] seemed to have all the features I needed. I began developing Hugo in Go and fell in love with the language. I hope you will enjoy using Hugo (and [contributing to it][]) as much as I have writing it.

&#8213;Steve Francia ([@spf13][])

[contributing to it]: https://github.com/spf13/hugo
[Go language]: https://golang.org/
[Jekyll]: https://jekyllrb.com/
[Middleman]: https://middlemanapp.com/
[Nanoc]: https://nanoc.ws/
[@spf13]: https://twitter.com/@spf13
