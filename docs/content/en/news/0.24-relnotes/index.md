---
date: 2017-06-17T17:53:58-04:00
categories: ["Releases"]
description: "The Revival of the Archetypes!"
link: ""
title: "Hugo 0.24"
draft: false
author: bep
aliases: [/0-24/]
---

> "A feature that could be the name of the next Indiana Jones movie deserves its own release," says [@bep](https://github.com/bep).

Hugo now handles the **archetype files as Go templates**. This means that the issues with sorting and lost comments are long gone. This also means that you will have to supply all values, including title and date. But this also opens up a lot of new windows.

A fictional example for the section `newsletter` and the archetype file `archetypes/newsletter.md`:

```
---
title: "{{ replace .TranslationBaseName "-" " " | title }}"
date: {{ .Date }}
draft: true
---

**Insert Lead paragraph here.**

<!--more-->

## New Cool Posts

{{ range first 10 ( where .Site.RegularPages "Type" "cool" ) }}
* {{ .Title }}
{{ end }}
```

And then create a new post with:

```
hugo new newsletter/the-latest-cool.stuff.md
```

**Note:** the site will only be built if the `.Site` is in use in the archetype file, and this can be time consuming for big sites.

**Hot Tip:** If you set the `newContentEditor` configuration variable to an editor on your `PATH`, the newly created article will be opened.

The above _newsletter type archetype_ illustrates the possibilities: The full Hugo `.Site` and all of Hugo&#39;s template funcs can be used in the archetype file.

**Also, Hugo now supports archetype files for all content formats, not just markdown.**

Hugo now has:

* 17839&#43; [stars](https://github.com/gohugoio/hugo/stargazers)
* 493&#43; [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 166&#43; [themes](http://themes.gohugo.io/)

## Notes

* Archetype files now need to be complete, including `title` and `date`.
* The `-f` (format) flag in `hugo new` is removed: Now use the archetype files as is.

## Enhancements

* Support extension-less media types. The motivation behind this change is to support Netlify&#39;s `_redirects` files, so we can generate server-side redirects for the Hugo docs site. See [this commit](https://github.com/gohugoio/hugoDocs/commit/c1ab9894e8292e0a74c43bbca2263b1fb3840f9e) to see how we configured that. [0f40e1fa](https://github.com/gohugoio/hugo/commit/0f40e1fadfca2276f65adefa6d7d5d63aef9160a) [@bep](https://github.com/bep) [#3614](https://github.com/gohugoio/hugo/issues/3614) 
* Add `disableAliases` [516e6c6d](https://github.com/gohugoio/hugo/commit/516e6c6dc5733cdaf985317d58eedbc6ec0ef2f7) [@bep](https://github.com/bep) [#3613](https://github.com/gohugoio/hugo/issues/3613) 
* Support non-md files as archetype files [19f2e729](https://github.com/gohugoio/hugo/commit/19f2e729135af700c5d4aa06e7b3540e6d4847fd) [@bep](https://github.com/bep) [#3597](https://github.com/gohugoio/hugo/issues/3597) [#3618](https://github.com/gohugoio/hugo/issues/3618) 
* Identify extension-less text types as text [c43b512b](https://github.com/gohugoio/hugo/commit/c43b512b4700f76ac77f12d632bb030c3a241393) [@bep](https://github.com/bep) [#3614](https://github.com/gohugoio/hugo/issues/3614) 
* Add `.Site` to the archetype templates [662e12f3](https://github.com/gohugoio/hugo/commit/662e12f348a638a6fcc92a416ee7f7c2a7ef8792) [@bep](https://github.com/bep) [#1629](https://github.com/gohugoio/hugo/issues/1629) 
* Use archetype template as-is as a Go template [422057f6](https://github.com/gohugoio/hugo/commit/422057f60709696bbbd1c38c9ead2bf114d47e31) [@bep](https://github.com/bep) [#452](https://github.com/gohugoio/hugo/issues/452) [#1629](https://github.com/gohugoio/hugo/issues/1629) 
* Update links to new discuss URL [4aa12390](https://github.com/gohugoio/hugo/commit/4aa1239070bb9d4324d3582f3e809b702a59d3ac) [@bep](https://github.com/bep) 

## Fixes

* Fix error handling for `JSON` front matter [fb53987a](https://github.com/gohugoio/hugo/commit/fb53987a4ff2acb9da8dec6ec7b11924d37352ce) [@bep](https://github.com/bep) [#3610](https://github.com/gohugoio/hugo/issues/3610) 
* Fix handling of quoted brackets in `JSON` front matter [3183b9a2](https://github.com/gohugoio/hugo/commit/3183b9a29d8adac962fbc73f79b04542f4c4c55d) [@bep](https://github.com/bep) [#3511](https://github.com/gohugoio/hugo/issues/3511) 