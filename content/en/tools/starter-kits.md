---
title: Starter Kits
linktitle: Starter Kits
description: A list of community-developed projects designed to help you get up and running with Hugo.
date: 2017-02-22
publishdate: 2017-02-01
lastmod: 2018-08-11
keywords: [starters,assets,pipeline]
menu:
  docs:
    parent: "tools"
    weight: 30
weight: 30
sections_weight: 30
draft: false
aliases: [/developer-tools/migrations/,/developer-tools/migrated/]
toc: false
---

Know of a Hugo-related starter kit that isn't mentioned here? [Please add it to the list.][addkit]

{{% note "Starter Kits are Not Maintained by the Hugo Team"%}}
The following starter kits are developed by active members of the Hugo community. If you find yourself having issues with any of the projects, it's best to file an issue directly with the project's maintainer(s).
{{% /note %}}

* [Hugo Wrapper][hugow]. Hugo Wrapper is a POSIX-style shell script which acts as a wrapper to download and run Hugo binary for your platform. It can be executed in variety of [Operating Systems][hugow-test] and [Command Shells][hugow-test].
* [Victor Hugo][]. Victor Hugo is a Hugo boilerplate for creating truly epic websites using Webpack as an asset pipeline. Victor Hugo uses post-css and Babel for CSS and JavaScript, respectively, and is actively maintained.
* [GOHUGO AMP][]. GoHugo AMP is a starter theme that aims to make it easy to adopt [Google's AMP Project][amp]. The starter kit comes with 40+ shortcodes and partials plus automatic structured data. The project also includes a [separate site with extensive documentation][gohugodocs].
* [Blaupause][]. Blaupause is a developer-friendly Hugo starter kit based on Gulp tasks. It comes ES6-ready with several helpers for SVG and fonts and basic structure for HTML, SCSS, and JavaScript.
* [hugulp][]. hugulp is a tool to optimize the assets of a Hugo website. The main idea is to recreate the famous Ruby on Rails Asset Pipeline, which minifies, concatenates and fingerprints the assets used in your website.
* [Atlas][]. Atlas is a Hugo boilerplate designed to speed up development with support for Netlify, Hugo Pipes, SCSS & more. It's actively maintained and contributions are always welcome.


[addkit]: https://github.com/gohugoio/hugo/edit/master/docs/content/en/tools/starter-kits.md
[amp]: https://www.ampproject.org/
[Blaupause]: https://github.com/fspoettel/blaupause
[GOHUGO AMP]: https://github.com/wildhaber/gohugo-amp
[gohugodocs]: https://gohugo-amp.gohugohq.com/
[hugow]: https://github.com/khos2ow/hugo-wrapper
[hugow-test]: https://github.com/khos2ow/hugo-wrapper#tested-on
[hugulp]: https://github.com/jbrodriguez/hugulp
[Victor Hugo]: https://github.com/netlify/victor-hugo
[Atlas]: https://github.com/indigotree/atlas
