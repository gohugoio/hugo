---
title: Docs Concept
linktitle: Docs Concept
description: Notes on Hugo docs overhaul.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-24
weight: 01
categories: []
tags: []
draft: false
aliases: []
toc: true
---

## Background

The claims made in this document are largely *empirical* and pulled from two major sources:

* My experience starting 18 months ago as a new Hugo user.
* Conversations with fellow Hugo users and noted trends within the [Discussion Forum][forum].

I am submitting these changes because I have been given [explicit permission to do so](https://discuss.gohugo.io/t/roadmap-to-hugo-v1-0/2278/34). Project kickoff was in November 2016 as a result of [this forum thread on a proposed source re-organization](https://discuss.gohugo.io/t/proposed-source-organization-for-hugo-docs-concept/4506). See also [this Discuss thread from @bep re: a new restructure and redesign](https://discuss.gohugo.io/t/documentation-restructure-and-design/1891/10).

{{% note "What is this document, Ryan?" %}}
At the bottom of this document is a section that starts with the [most important questions for the Hugo team](#ideal-workflow) and continues with a very high-level description of an *ideal* content strategy process. This has only been included to provide insight for those unfamiliar with basic content strategy..
{{% /note %}}

{{% warning "Disclaimer" %}}
WIP. Before any colleagues banish me to content strategy hell, know that I *know* this is a *schlocky* version of a true strategic document. It'll get better. I promise.
{{% /warning %}}

## Strategy, Tactics, and Requirements

### Assumptions

The biggest shortcomings of the current documentation are more a matter of *content* than visual design.

Additionally, current Hugo documentation

* is confusing for new users
* is a common complaint in the Hugo forums ([forum discussion 1][ex1], [forum discussion 2][ex2])
* lacks structure and is therefore
    * unscalable, as demonstrated by patch pages (e.g. [here][patch1] and [here][patch2]) that seem out of place, require unnecessary drill-down, or duplicate content in other areas of the docs, thus requiring duplicative efforts to keep updated
    * inconsistent in its terminology, style, and (sometimes) layout
    * limited in effective use of Algolia's document search (i.e., due to redundant content grouping, headings, etc)
    * difficult to optimize for external search engines (SEO)
* does not leverage Hugo's more powerful features (e.g., only *one* archetype); leveraging these features would help address the aforementioned shortcomings (i.e., scalability, consistency, and search)
* assumes a higher level of Golang proficiency than is realistic for newcomers to static site generators or general web development; a prime example is the sparsity of basic and advanced code samples for templating functions, some of which may still be wholly undocumented
* lacks well-defined contribution and editorial guidelines for those interested in submitting or editing documentation
* is not up to date
* often unusable on smart phones or other mobile devices
* needs reworking of its URL structure to represent a more intuitive mental model for end users ([Discuss][ds1])

### Goals

New Hugo documentation should...

* reduce confusion surrounding core Hugo concepts; e.g., `list`, `section`, `page`, `kind`, and `content type` with the intention of
    * making it easier for new users to get up and running
    * creating better consistency and scalability for Hugo-dependent projects (viz., [themes.gohugo.io][hugothemes])
    * reducing the frequency of beginner-level questions in the [Hugo Discussion Forum][forum]
* not require, or assume, any degree of Golang proficiency from end users;
    * that said, Hugo can&mdash;and *should*&mdash;act as a bridge for users interested in learning Golang. An implemented example of this strategic point is the inclusion of `godocref:` as a default front matter field for all function and template pages. See [`archetypes/functions.md`][functionarchetype].
* be easiest to expand and edit for *contributors** but even easier to understand by *end users*
     * if you don't make it *very easy* for authors to contribute to documentation correctly, they will inevitably contribute *incorrectly*;
        * content modeling is king
        * go DRY (e.g., via shortcodes)
        * set required metadata (e.g., via section-specific *archetypes*)
        * develop contribution guidelines for both development *and* documentation
* be equally accessible via mobile, tablet, desktop, and offline.
* avoid "miscellaneous" sections (e.g.,"Extras"); [these catch-all sections are the last place end users look to get up and running](https://discuss.gohugo.io/t/site-with-different-lists-of-sections/5536/3)
    * all content in miscellaneous sections should be edited and incorporated into more logical content groupings with the goal of removing miscellaneous sections entirely.
* easily scaffold for potential *i18n/multilingual versions*.

### Audiences

* Primary: Web developers interested in static site generators
* Secondary: Web publishers (bloggers, authors) and hobbyists
* Tertiary: Web developers, both novice and professional, interested in learning Golang

### Persona

{{% note %}}
This is far from an inadequate persona exercise, but I think it helped my mental model as I worked through the existing docs.
{{% /note %}}

#### End User: SSG Developer

The SSG developer has

* basic proficiency in Git and DVCS
* no to little proficiency in Golang
* working proficiency in front-end development---HTML, CSS, JS---but not necessarily front-end build tools
* basic familiarity with at least one double-curly templating language (e.g., liquid, Twig, Swig, or Django)
* proficiency in the English language for the current version of the documentation
    * proficiency in other languages (for future multilingual versions)

#### End User: Themes (i.e. blogger/author/hobbyist)

The themes end user has

* limited proficiency in the command line/prompt
* proficiency in one of the [supported content formats](https://hugodocsconcept.netlify.com/content-management/formats/)(specifically markdown)
* access to static hosting
* limited proficiency in deploying a static website

### Requirements

The following are high-level requirements for the documentation site.

#### Technical

- [X] Built with Hugo
- [X] Performant (e.g., 80+ [Google Page Speed Score](https://developers.google.com/speed/pagespeed/insights/?url=https%3A%2F%2Fhugodocsconcept.netlify.com%2Fabout-hugo))
- [X] Front-end build tools for concatenation, minification, of static assets
- [X] Browser compatibility: modern (i.e. Chrome, Edge, Firefox, Safari) and IE11
- [ ] CDN
- [ ] AMP?

#### SEO

- [X] [Open Graph Protocol](http://ogp.me/)
- [X] [schema.org](http://schema.org)
- [X] [JSON+LD](https://developers.google.com/schemas/formats/json-ld), [validated](https://search.google.com/structured-data/testing-tool)
- [X] Consistent heading structure
- [X] Semantic HTML5 elements (e.g., `article`, `main`, `aside`, `dl`)
- [X] SSL
- [ ] AMP?
- [ ] 301s [^1]

#### Accessibility

- [X] Aria roles
- [X] Alt text for all images
- [X] `title` attribute for links

#### Editorial and Content

- [X] Basic style guide
    - The style guide should facilitate a more consistent UX for the site but not be so complex as to deter documentation contributors
- [X] Contribution guidelines (see [WIP on live site](https://hugodocsconcept.netlify.com/contribute/documentation/))
- [X] Standardized content types (see [WIP archetypes in source](https://github.com/rdwatters/hugo-docs-concept/tree/master/themes/hugodocs/archetypes)
- [X] New content model, including taxonomies ([see tags page][tagspage])
- [X] DRY. New shortcodes for repeat content (e.g., lists of aliases, page variables, site variables, and others)
- [X] New site architecture and content groupings

#### Content Strategy Statement

[What is this?](http://contentmarketinginstitute.com/2016/01/content-on-strategy-templates/)

> The Hugo documentation increases the Hugo user base and strengthens the Hugo community by providing intuitive, beginner-friendly, regularly updated usage guides. Hugo documentation makes visitors feel excited and confident that Hugo is the ideal choice for static website development.

#### Editorial Mission

[What is this?](http://contentmarketinginstitute.com/2015/10/statement-content-marketing/)

> The Hugo documentation is a joint effort between the Hugo maintainers and the open-source community. Hugo documentation is designed to promote Hugo, the world's fastest, friendliest, and most extensible static site generator. Hugo documentation is the primary vehicle by which the Hugo team reaches its target audiences. When visitors comes to the Hugo documentation, we want them to install Hugo, develop a new static website with our tool, and share their progress and insights with the Hugo community at large.

## UX/UI

- [X] Copyable and/or downloadable code blocks (via highlight.js and clipboard.js, extended for hugo-specific keywords); see [this specific request from @bep for this feature](https://discuss.gohugo.io/t/proposed-source-organization-for-hugo-docs-concept/4506/10?u=rdwatters)
- [X] Dual in-page navigation (i.e. site nav *and* in-page TOC)
- [X] Smooth scrolling
- [X] [RTD-style admonitions][admonitions] (see [example admonition shortcode](https://github.com/rdwatters/hugo-docs-concept/blob/master/layouts/shortcodes/note.html) and [examples on published site](/contribute/documentation/#admonition-short-codes))
- [ ] Share buttons: Reddit, Twitter, LinkedIn, and "Copy Page Url"; the last of these provides the strongest utility for docs references in the Hugo forums

## Author Experience (AX)

- [X] Easy scaffolding of content types (i.e., via Hugo CLI [`hugo new`])
- [X] Type-based content storage model and scope (i.e, via archetypes)

## Analytics/Metrics

- [X] Google Analytics
- [ ] Content groupings (GA) to measure usage, behavior flow, and define content gaps
- [ ] Automated reports (GA)
    - [ ] Potential for using Google Sheets CSV for published content on the docs

{{% note %}}
The preceding analytics and metrics are separate from usage statics re: Hugo downloads, `.Hugo.Generator`, etc.
{{% /note %}}

## Visual Design and Front-end Development

- [X] Clean typography with open-source font
    - [X] Optimal line length (50-80 characters)
    - [X] Consistent vertical rhythm
- [X] Responsive
    - [X] Flexbox
    - [X] Typography (via ems)
- [X] Custom iconography (i.e., via Fontello rather than FontAwesome bloat)
- [X] Design assets versioned with source ([see design resources directory][designresources])
- [X] [WCAG color contrast requirements](http://webaim.org/blog/wcag-2-0-and-link-colors/)
- [X] [Sass Guidelines for Source Organization](https://sass-guidelin.es/)
    - [X] Abstracted color palette
    - [X] Abstracted typefaces (multiple open-source fonts available)


## Annotated Content Changes

The following is an *abbreviated* listing of *substantive* changes made to the current documentation's source content and organization. Sections here are ordered according to the current site navigation. The changes delimited here do *not* include copy edits for consistent or preferred usage, improvements in semantics, etc, all of which easily numbers in the thousands, likely more.

### Download Hugo

This is no longer a site navigation link and is instead a button along with "File an Issue" and "Discuss Hugo" at the bottom of the sidebar.

### Site Showcase

Site showcase has stayed more or less as is, including styling, etc. However...
 * The showcase archetype has been refined for simplicity and now has new [contribution guidelines](/contribute/documentation).
 * To keep compatibility, all [showcase content files][showcasefiles] have been edited to reflect the new content type. This will also be updated in the ["docs" page of the contribute section](/contribute/documentation/)

### Press & Articles

* The press and article pages have been moved under "News" along with "Release Notes". This section is lower on the navigation because it's less frequently visited than other areas (my assumption).
* Like everything else, I've kept up with changes to the docs upstream on GitHub, but in this case, I also included a [half dozen *new* articles as well](/news/press-and-articles/).

### About Hugo

* Content is more or less the same, but I've cleaned up a lot of the language, and copy edited for consistency throughout. I've also added in some extra frills (e.g. resources that extoll and teach more about the benefits of SSGs on [/about/benefits](/about/benefits/)).
* Release notes are now in the "News" section, although I'm still unsure on this decision. I can gladly move this back into "About" to give it a higher degree of discoverability in the menu.

### Getting Started

* **[[UPDATE 2017-03-12]]**
    * The Quick Start needs to be completely reworked; this could include a walkthrough of a new default theme (#)
* ~~The [Quick Start][] has been completely updated for more consistent heading structure, etc. Also, **I may delete the "deployment" section of the Quick Start** since this a) adds unnecessary length, making the guide less "quick" and b) detracts from the new "hosting and deployment" section, which offers better advice, and c) is redundant with [Hosting on Github](https://hugodocsconcept.netlify.com/hosting-and-deployment/hosting-on-github/). For example, the Quick Start didn't mention that files already written to public are not necessarily erased at build time. This can cause problems with drafts. I think the other options&mdash;e.g. Arjen's Wercker tutorial&mdash;are more viable and represent better practices for newcomers to Hugo. If future versions of Hugo include baked-in deployment features, I think it's worth reconsidering adding the deployment step back to the Quick Start.~~

### Content

* This section has been renamed "Content Management" to facilitate elimination of the ["extras"](http://gohugo.io) section. **Note**: this section does *not* include any templating. The convention is `content-management/concept.md` (for explanation and usage) `templates/concept-templates.md` (for examples, functions, etc), and then `variables/concept-variables.md`.
    * That said, I'm working on refactoring a series of shortcodes for variables so that it's only a matter of referencing them once and having them update everywhere.

### Themes

Themes section organization has only changed slightly in that the 6 content pages have been consolidated to just 4.

* "Installing a theme" and "Using a Theme" have been combined since one largely dovetails with the other. The current [using a theme page](http://gohugo.io/themes/usage/) is pretty skimpy. An alias for `themes/usage` has been set up accordingly.
*

### Templates

* Reworked considerably. Page titles have all been changed to reflect their obvious connection to *templating*.
* "Lookup order" page added. Placement of this page in the site navigation is significant in that the order of the template pages are no sequentially ordered as to when they should be learned; i.e.---
    * Go template primer
    * Lookup order
    * Base templates
    * Hugo Lists (introduces the lists concept [i.e. sections, taxonomies, etc]); this includes one of multiple forthcoming visualizations for Hugo architecture
    * Rendering Hugo Lists (i.e., ordering, grouping, etc)
* Shortcodes and menus (templating), pagination, data, traversing local files, data-driven content, and data files have all moved out from "Extras" and into templating. *Note that there is only one stylesheet in the local example now*..
* **2017-02-26**. I am currently working on a new example for `readDir`.

### Taxonomies

*Taxonomies* is no longer an independent section. Similar to shortcodes and menus, taxonomies is broken into two equal pages: one under Content Management, and the other under Templates.


### Extras

* This section no longer exists in the new documentation site
    * *Extras*, in the content world, is the equivalent of *miscellaneous* or *additional resources*. It's an area that's been tacked onto site navigation to accommodate a *seemingly* disparate set of new features. In other words, READ: "We don't have any idea of where to put this"
* *Extras* pages:
    * Aliases. Incorporated into [URL Management](/content-management/urls/)
    * Analytics. Incorporated into [built-in partials](/templates/partials/#using-hugos-built-in-partials)
    * Builders. This has been removed completely since it has no real added value. The three "builders" mentioned (`new site`, `new theme`, and `new <content>`) are all well-delimited in their respective pages, which is where end users expect to find this type of information in the first place.
    * Comments. Incorporated into [content management](/content-management/comments/) for content-related pieces and mentioned in [partials](/templates/partials/) for implementation.
    * Cross-References. Added as its own page under Content Management (`/content-management/cross-references/`)
    * Custom robots.txt. Incorporated into [/templates/robots/](/templates/robots/)
    * Data Files and Data-Driven Content. Combined and incorporated into [/templates/data-templates/](/templates/data-templates/)
    * GitInfo. Incorporated into [/variables/git/](/variables/git/)
    * LiveReload. This doesn't really merits its own page. It's mentioned in features, usage, about, and elsewhere.
    * Menus. This is broken into [/content-management/menus/](/content-management/menus/) and [/templates/menu-templates/](/templates/menu-templates/)
    * Pagination. Now in [/templates/pagination/](/templates/pagination/)
    * Permalinks. Now a heading/subsection of [/content-management/urls/](/content-management/urls/
    * Scratch. Has it's own devoted function page at [/functions/scratch/](/functions/scratch/), and is therefore in the [Functions Quick Reference](/functions/). Also mentioned in an admonition in [/variables/page/](/variables/page/).
    * Shortcodes. Now split into two pages, one a [/content-management/shortcodes/](/content-management/shortcodes/) and the templating portion (i.e. "create your own shortcodes") at [/templates/shortcode-templates/]
    * URLS. Now combined with permalinks and others as a heading/subsection of [/content-management/urls/](/content-management/urls/)
    * Syntax Highlighting. The shortcode is featured and explained with usage examples at [/content-management/shortcodes/](/content-management/shortcodes/), as well as expaned upon in it's own page under [/tools/syntax-highlighting/](/tools/syntax-highlighting/). I did this under the assumption that *developers* are most interested in adding code blocks to their content.
    * Table of Contents. This is now it's own page under [/content-management/toc](/content-management/toc/) and referenced in [/variables/page-variables](/variables/page/).
    * Traversing Local Files. This is now split into [/templates/files/](/templates/files/) and variables delimited at [/variables/files/](/variables/files/).

### Community

The "Community" section has been removed as a site navigation item because `/contribute` is now it's own section.
    * There are now more calls than ever for contributing to Hugo throughout the Hugo docs.
    * The most important changes to "Community" (now [/contribute/](/contribute/)) are that @digitalcraftsman's tutorial on contributing to Hugo can be found under [/contribute/development/](/contribute/development/), and a brand-new page for *contributing to documentation*, including examples of shortcodes used throughout the docs, etc, can be found at [/contribute/documentation/]. This is a **VERY IMPORTANT* change since it's instructions, archetypes, and guidelines like this that will make the documentation site scale more easily.

### Tutorials

* Original page: <http://gohugo.io/tutorials>
* All installation guides have been consolidated under [/getting-started/installing/](/getting-started/installing/)
    * Why? Installing Hugo shouldn't be considered a separate tutorial
    * "Tutorials" is not an intuitive place for end users to look for this kind of documentation
* All content moved from `/tutorials` has been edited to reflect a less tutorial-ish style of language (e.g., removal of lines starting with "In this tutorial...")
* Aliases added to new pages and in-page links updated throughout
* Remaining Tutorials
    * Are these worth keeping in their entirety if they reflect (sometimes much) older versions of Hugo?
    * Michael Henderson's "Creating a Theme" website ([current][],[concept][]) has been copy edited and content edited to include the new code block shortcodes. Michael did an *amazing* job with this tutorial, and it must have taken him *forever*, but much of the information included in the tutorial is now spread throughout the documentation in more appropriate places. Also, because this is an older tutorial, some of the paradigms aren't quite as up to date.
    * Rick Cogley's still needs to be copy edited a bit, but overall looks good. That said, this tutorial was put together before Hugo began implementing it's international features.
* **Guidelines for New Tutorials**
    * **[[Update 2017-03-12]]** "Tutorials" as a site section has been completely removed. The remaining three articles: 1. Creating a Multilingual Site, 2. Creating a New Theme, and 3. Migrate from Jekyll to Hugo were a) duplicative, b) outdated, or c) contradictory (specifically w/r/t Hugo-specific terms). Full-length tutorials can be a very difficult thing to maintain within documentation; it's my position that it's better to provide examples and explanations throughout the documentation and then allow the community to publish their own tutorials (i.e., with links to the tutorials on the Articles Page).
    * ~~To keep the content in tutorials maintainable, it's important to set standards on what should be contained within said tutorials when published directly to the Hugo docs. (Of course, listing beginning-to-end tutorials in other areas of the website [i.e., press and articles] is a very good idea). The following pieces of information should be omitted from full-text tutorials in the Hugo docs because they are better delimited and kept current in other areas:~~
        * ~~Explanations directory structure or content organization~~
        * ~~Explanations of content formats (namely, `.md`) or front matter~~
        * ~~Explanations of how to set up hosting, deployments, or automated deployments (although these make excellent additions to the "Hosting and Deployments" section)~~
        * ~~"Using Hugo Shortcodes with Google Sheets or Data-Driven Content" is a better tutorial example than "Getting Up and Running with Hugo" or "Deploying Your Hugo Website with an Apache Server"~~

### Troubleshooting

* This section still only contains the same two troubleshooting content pages from the current site.
    * Both pages have been copy edited, and the markdown has been cleaned up for consistency.

### Tools

* Rather than a single-page list, a full "developer tools" section is part of the main navigation and includes the following pages:
    * Migrate to Hugo. List of project-descriptions of community-developed migration tools
    * Syntax highlighting. This builds on the syntax highlighting shortcode used in [/shortcodes/#highlight](/shortcodes/#highlight).
    * Starter Kits. Only two items for now, but this should remain a community-aggregated (and edited) list of kits developed to help new users get up and running.
    * Frontends. Same frontends material previously under "tools." Copy edited for consistency.
    * Editor Plug-ins. Same editor plug-ins material found in current documentation. Copy edited for consistency.
    * Search. Same search material under "tools" in current documentation. Copy edited for consistency.
    * Other projects. This might be worth restructuring since I'm not a fan of catch-all sections or pages.

### Hugo Cmd Reference

This hasn't been touched. I'll make the necessary style changes once/if the site is integrated into the Hugo GH repo. I believe these pages are pulled automatically using Viper.
    * **2017-02-26 Update** Pulled latest version of generated command files after noticing [this commit from BEP](https://github.com/spf13/hugo/tree/a8a8249f67b3725b68316f5e22410554405d1e5e)

### Issues & Help

This is no longer a site navigation link and is instead a button along with "Download" and "Discuss Hugo".

## Content Ordering: Site Navigation

The following shows weights and ordering for the newly restructured site architecture.

### "about-hugo" Ordering (`.OrderByWeight`)

_index - 01
What Is Hugo - 10
Hugo Features - 20
The Benefits Of Static - 30
Why I Built Hugo - 40
Roadmap - 50
Apache License - 60

### "getting-started" Ordering (`.OrderByWeight`)

* _index.md - 01
* Quick Start - 10
* Using the Hugo Docs - 20
* Install from Source - 30
* Install on Linux - 40
* Install on Mac - 50
* Install on Pc - 60
* Basic Usage - 70
* Directory Structure - 80
* Configuration - 90

### "content-management" Ordering (`.OrderByWeight`)

* _index.md - 01
* Content Organization - 10
* Supported Content Formats - 20
* Front Matter - 30
* Shortcodes - 40
* Sections - 50
* Authors - 55 (in draft status)
* Content Types - 60
* Archetypes - 70
* Taxonomies - 80
* Content Summaries - 90
* Cross References - 100
* URL Management - 110
* Menus - 120
* Table of Contents - 130
* Comments - 140
* Multilingual Mode - 150

### "templates" Ordering (`.OrderByWeight`)

* _index - 01
* Go Template Primer - 10
* Template Lookup Order -15
* Base Templates And Blocks - 20
* Lists in Hugo - 22
* Ordering and Grouping Hugo Lists - 27
* Homepage Template - 30
* Section Templates - 40
* Taxonomy Templates - 50
* Single Page Templates - 60
* Content View Templates - 70
* Data Templates - 80
* Partial Templates - 90
* Shortcode Templates - 100
* Local File Templates - 110
* 404 Page - 120
* Menu Templates - 130
* Pagination - 140
* RSS Templates - 150
* Sitemap Template - 160
* Robots.txt - 165
* Internal Templates - 168
* Additional Templating Languages - 170
* Template Debugging - 180

### "functions" Ordering (`.OrderByTitle`)

* _index.md - 01 (i.e., "Functions Quick Reference")

### "variables" Ordering (`.OrderByWeight`)

* _index.md - 01
* Site Variables - 10
* Page Variables - 20
* Taxonomy Variables - 30
* File Variables - 40
* Menu Variables - 50
* Hugo Variables - 60
* Git Variables - 70
* Sitemap Variables - 80

### "themes" Ordering (`.OrderByWeight`)

* _index.md - 01
* Installing and Using Themes - 10
* Customizing a Theme - 20
* Creating a Theme - 30
* Theme Showcase - 40

### "hosting-and-deployment" Ordering (`.OrderByWeight`)

* _index.md - 01
* Hosting on Netlify - 10
* Hosting on Github - 20
* Hosting on Gitlab - 30
* Hosting on Bitbucket - 40
* Deployment with Wercker - 50
* Deployment with Rsync - 60


### "tools" Ordering

* _index.md - 01
* migrate to hugo - 10
* syntax highlighting - 20
* starter kits - 30
* frontends - 40
* editor plug-ins - 50
* search - 60
* other (community projects) - 70

### "showcase" Ordering (`.OrderByPublishDate.Reverse`)


### "troubleshooting" Ordering (`.OrderByTitle`)


## Current Content (Source)

```markdown
{{< readfile file="content/tree.txt" >}}
```

## Ideal Workflow

The following demonstrates the *ideal* sequence for content strategy and is only included in this document for those who are unfamiliar with some of these concepts.

### 1. Defining Goals and Audiences

* What are the goals for this product?
* Who are the targeted users/audiences?
    * Personal development

### 2. Content Stragegy Development

* What is the editorial mission statement?
    * What are we promising to target audiences to help them achiece *their* goals? (We [at least ostensibly] serve *their* needs first)
    * What have we failed to deliver to end users with the current website in helping them accomplish their goals?
* What is the content strategy statement? I.e., now that end users know what they're getting...
    * How do we accomplish *our* goals?
    * [Core model][]
        * What content do we already have?
        * How did users get to our site?
        * What do we want them to do after they leave?
        * How do their actions serve to execute on our primary strategy?

### 3. Content Strategy Implementation (i.e. *tactics*)

* Content typing
    * Define content types and attributes (in this case, archetypes with fixed front matter)
* [Content modeling][]
    * What are the defined relationships between the content types?
    * Basic IA
* Build out all content as instances of content types. (The general rule is that *more structure = more freedom*)
* Content grouping (this may or may not flip positions with the last step)
* Site architecture and IA
* Define opportunities for repurposing content across multiple areas, COPE, and scaling

### 4. Website Development

* Begin development of templating and UX considerations
* Guidelines
    * Style guide (editorial, not visual)
    * Contribution guidelines (including descriptions of types to make it easier on authors/contributors)
* Visual Design


## Proposed Schedule for Hugo Docs Release

If the Hugo Team finds the improvements to the Hugo documentation acceptable, I've proposed the following schedule for releasing the new Hugo documentation.

1. **2017-02-27** Release to Gitter Channel for dev review
2. **2017-03-05** Post in Discussion Forum for Hugo user feedback
3. **2017-03-6** Pull request/add to Hugo Rep
3. **2017-??-??** Add to Hugo repo for release with v19?

[^1]: At this point, not too much of the URL structure has changed that considerably. I've been fastidious about adding aliases wherever possible and trying to retain URLs whenever still applicable. That said, the [current list of aliases is quite large](/contribute/documentation/#be-mindful-of-aliases).

[admonitions]: http://docutils.sourceforge.net/docs/ref/rst/directives.html#admonitions
[Content Modeling]: https://alistapart.com/article/content-modelling-a-master-skill
[Core Model]: https://alistapart.com/article/the-core-model-designing-inside-out-for-better-results
[designresources]: https://github.com/rdwatters/hugo-docs-concept/tree/master/dev-and-design-resources
[ds1]: https://discuss.gohugo.io/t/hugo-docs-domain/5657
[ex1]: https://discuss.gohugo.io/t/frustrated-with-documentation/2810
[ex2]: https://discuss.gohugo.io/t/documentation-restructure-and-design/1891
[forum]: https://discuss.gohugo.io
[functionarchetype]: https://github.com/rdwatters/hugo-docs-concept/blob/master/themes/hugodocs/archetypes/functions.md
[hugothemes]: http://themes.gohugo.io
[patch1]: http://gohugo.io/taxonomies/templates/
[patch2]: https://github.com/spf13/hugo/commit/eaabecf586fd0375585e27c752e05dd8cb4c72b4
[Quick Start]: https://hugodocsconcept.netlify.com/getting-started/quick-start/
[showcasefiles]: https://github.com/rdwatters/hugo-docs-concept/tree/master/content/showcase
[tagspage]: https://hugodocsconcept.netlify.com/tags/