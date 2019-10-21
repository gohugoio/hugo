---
date: 2018-07-04
title: "Let’s celebrate Hugo’s 5th birthday"
description: "How a side project became one of the most popular frameworks for building websites."
categories: [blog]
author: bep
---

_By Bjørn Erik Pedersen ([@bepsays](https://twitter.com/bepsays) / [@bep](https://github.com/bep)), Hugo Lead_

**Five years ago today, [Steve Francia](https://github.com/spf13/) made his [first commit](https://github.com/gohugoio/hugo/commit/50a1d6f3f155ab837310e00ffb309a9199773c73
) on the Hugo project: "Hugo: A Fast and Flexible Static Site Generator built with love by spf13 in GoLang".**

Steve was writing that on a train commute to New York. I'm writing this article running Hugo `v0.43-DEV`, the preview version of the next Hugo release. The release is scheduled for Monday and adds a powerful [assets pipeline](https://github.com/gohugoio/hugo/issues/4854#issue-333062459), with SCSS/SASS support, assets minification, fingerprinting/subresource integrity, ad-hoc image processing and much more. 

**I cannot remember the last time I was this excited about a Hugo release. "Game changer" may be too strong, but it makes for a really nice integrated website design-workflow that, with Hugo's build speed, is hard to beat.**

{{< imgproc sunset Fill "600x300" >}}
Fetch and scale an image in the upcoming Hugo 0.43.
{{< /imgproc >}}

But that is a release for Monday. Now is a time to look at the current status of Hugo after the first five years.

## Hugo in Numbers

According to [BuiltWith](https://trends.builtwith.com/cms/Hugo), more than 29 000 live websites are built with Hugo. Of those, 390 are in the top 1 million. Wappalyzer [reports](https://www.wappalyzer.com/categories/static-site-generator) that Hugo serves almost 50% of the static sites.

Hugo is big in the [public sector](https://discourse.gohugo.io/t/hugo-in-public-administration/8792), with the US Government as a prominent user. Some examples are [vote.gov](https://vote.gov/) and [digital.gov](https://digital.gov/).

[Smashing Magazine](https://www.smashingmagazine.com/) is a big and very popular Hugo site. It is [reported](https://discourse.gohugo.io/t/smashing-magazine-s-redesign-powered-by-hugo-jamstack/5826/7) that they build their complex site with 7500 content pages in 13 seconds.

Some other example sites are [kubernetes.io](https://kubernetes.io/), [letsencrypt.org](https://gohugo.io/showcase/letsencrypt/), [support.1password.com](http://gohugo.io/showcase/1password-support/), [netlify.com](https://www.netlify.com), [litecoin.org](https://litecoin.org/), and [forestry.io](https://forestry.io/).


{{< imgproc graph-stars Fit "600x400" >}}
Number of GitHub stars in relation to the Hugo release dates.
{{< /imgproc >}}

More numbers:

* 26800+ [stars](https://github.com/gohugoio/hugo/stargazers) on GitHub. 
* 444+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors) to the Hugo source repository, 654+ to [Hugo Docs](https://github.com/gohugoio/hugoDocs/graphs/contributors). [@bep](https://github.com/bep) is the most active with around 70% of the current code base (line count).
* 235+ [themes](https://themes.gohugo.io/)
* 50% increase in the number of user sessions on the [gohugo.io](https://gohugo.io/) web sites the last 12 months.[^2]
* Hugo build release binaries for [a myriad](https://github.com/gohugoio/hugo/releases/tag/v0.42.2) of platforms. And since it can also be installed from Chocolatey on Windows, Brew on MacOs, Snap on Linux and `apt-get` on Debian/Ubuntu, it is impossible to give accurate download numbers. But the number is not small.

## Hugo Next

We're not finished with Hugo, but Hugo `0.43` very nicely wraps up the first five years. It started out as a small and fast static site generator. It is now [even faster](https://forestry.io/blog/hugo-vs-jekyll-benchmark/), and now so loaded with features that it has grown out of being just a "static site generator". It is a _framework for building websites_.

My interest in Hugo started on the Sunday when I moved my blog, [bepsays.com](https://bepsays.com/en/), twice. The second static generator choice of that day, Hugo, was a good match. I remember Steve being very enthusiastic about getting patches with fixes and enhancements, and I was eventually taken over by the simplicity and power of Go, the programming language, and started to implement lots of new features.

My goal with all of this, if there is one, is to get a single binary with native and really fast implementations of the complete stack I need for web development and content editing. The single binary takes most of the pain out of installation and upgrades (if you stick with the same binary, it will continue to just work for decades).

**With 0.43, we are almost there.** With that release, it should be possible to set up a Hugo-only project without any additional tools (Gulp, WebPack) for all aspects of website building. There will still be situations where those tools will still be needed, of course, but we will continue to fill the gaps in the feature set. 

Hugo has stuck with the sub-zero versions to signal active development, with a new main release every 5-6 weeks. But we take stability very seriously (breaking things add lots of support work, we don't like that) and most site upgrades are [smooth](https://twitter.com/tmmx/status/1006288444459503616). But we are closing in on the first major stable version.


### The Road to 1.0

We have some more technical tasks that needs to be done (there is ongoing work to get the page queries into a more consistent state, also a simpler `.GetPage` method), but also some cool new functionality. The following roadmap is taken from memory, and may not be complete, but should be a good indication of what's ahead.

Pages from "other data sources"
: Currently, in addition to Hugo's list pages, every URL must be backed by a content file (Markdown, HTML etc.). This covers most use cases, but we need a flexible way to generate pages from other data sources. Think product catalogues and similar.

Upgrade Blackfriday to v2
: [Blackfriday](https://github.com/russross/blackfriday) is the main content renderer in Hugo. It has been rewritten to a more flexible architecture, which should allow us to fix some of the current shortcomings.

We should be able to create a better and easier-to-use data structure from the rendered content: Summary, the content after the summary, being able to range over the footnotes and the ToC. Having ToC as a proper data structure also opens up a few other potential uses; using it as an index in [Related Content](https://gohugo.io/content-management/related/) would be one example.

This should also enable us to _do more_ with [Custom Output Formats](/templates/output-formats). It is already very powerful. GettyPubs are using it in [Quire](https://github.com/gettypubs/quire) to build [beautiful multi-platform publications](http://www.getty.edu/publications/digital/digitalpubs.html). But it can be improved. For rendering of content files, you are currently restricted to HTML. It would be great if we could configure alternative renderers per output format, such as LaTeX and EPUB.

Related to this is also to add a configurable "Markdown URL rewriter", which should make more portable URLs in Markdown, e.g. image links that work both when viewed on GitHub and your published site. 

### The Road to the Future

These are the items that first come to mind if you ask me to think even further ahead:

Dependency manager for Theme Components
: In Hugo `0.42` we added [Theme Components](/themes/theme-components/) and Theme Inheritance. With SCSS support in Hugo `0.43`, which also follows the same project/themes precedence order (add `_variables.scss` to your project, configure SASS colour variables in `config.toml`), we have a solid foundation for creating easy to use and extensible themes. But we are missing some infrastructure around this. We have a site with 235+ [themes](https://themes.gohugo.io/)[^themes] listed, but you currently need to do some added work to get the theme up and running for your site. In the Go world, we don't have NPM to use, which is a curse and a blessing, but I have some ideas about building a simple dependency manager into Hugo, modelled after how Go is doing it (`hugo install`). You should be able to configure what theme and theme components you want to use, and Hugo should handle the installation of the correct versions. This should make it easier for the user, but it would also enable community driven and even commercial "theme stores".


{{< imgproc graph-themes Fit "600x400" >}}
Number of Hugo themes on themes.gohugo.io in relation to the Hugo release dates.
{{< /imgproc >}}


The "New York Times on Hugo" Use Case
: There are recurring questions on the support forum from [really big sites](https://discourse.gohugo.io/t/transition-2m-posts-from-wordpress-to-hugo/12704) that want to move to Hugo. There are many [good reasons](https://www.netlify.com/blog/2016/05/18/9-reasons-your-site-should-be-static/) why they want this (security, cost-saving, EU regulations etc.). And while there have been reports about companies building 600 000 pages with Hugo on very powerful hardware, we will have to rethink the build model to make this usable. Keywords are: streaming builds, segmented builds, partial rebuilds. There are other site generators also talking about this. It should be possible, and my instinct tells me that it should be easier to do when your starting point is "really fast". But this is not a small weekend project for me, and I have already talked to several companies about sponsoring this.

Plugins
: A Theme Component could also be called a plugin. But there are several potential plugin hooks into Hugo's build pipeline: resource transformations, content rendering, etc. We will eventually get there, but we should do it without giving up too much of the Hugo speed and simplicity.


## Thanks

So, thanks to everyone who has contributed to getting Hugo where it is today. It is hard to single out individuals, but a big shout-out to all the Hugo experts and moderators helping out making [discourse.gohugo.io](https://discourse.gohugo.io/) a very active and possibly one of the best support forums out there.

And the last shout-out goes to two maintainers who have been there more or less from the start. [@digitalcraftsman](https://github.com/digitalcraftsman/) has been doing a fantastic job keeping the fast growing theme site and [repository](https://github.com/gohugoio/hugoThemes) in pristine condition. I have it on my watch list, but that is just out of curiosity. There is lots of activity, but it runs as clock work. [Anthony Fok](https://github.com/anthonyfok) has contributed with a variety of things but is most notable as the Linux expert on the team. He manages the Debian build and is the one to thank for up-to-date binaries on Debian and Ubuntu.

One final note: If you have not done so already, please visit [github.com/gohugoio/hugo](https://github.com/gohugoio/hugo) and push the "star button".

Gopher artwork by [Ashley McNamara](https://github.com/ashleymcnamara/gophers/) (licensed under [CC BY-NC-SA 4.0](https://creativecommons.org/licenses/by-nc-sa/4.0/)). Inspired by [Renee French](https://reneefrench.blogspot.com/).

[^2]: Numbers from Google Analytics. The Hugo websites are https://discourse.gohugo.io, https://gohugo.io and https://themes.gohugo.io. It is rumoured that when [Matt Biilman](https://twitter.com/biilmann?lang=en), CEO and Co-founder of Netlify, opened the first power bill after sponsoring Hugo's hosting, said: "Du må lave fis med mig, those Hugo sites have lots of web traffic!"
[^sgen]: That was at the time of writing this article. _Next_, a React based static site generator, has momentum and is closing in on Hugo's 2nd place. 
[^themes]: We pull all the themes from GitHub and build the theme site and 235 demo sites on Netlify in 4 minutes. That is impressive.
