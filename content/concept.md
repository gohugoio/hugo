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

### Introduction

The claims made in this strategic document are largely *empirical* and pulled from two major sources:

* My experience starting 18 months ago as a new Hugo user.
* Conversations with fellow Hugo users and noted trends within the [Discussion Forum][forum].

{{% warning "Disclaimer" %}}
WIP. Before any of my fellow content strategists banish me to content strategy hell, know that I *know* this is a *schlocky* version of a true strategic document. It'll get better. I promise.
{{% /warning %}}

## Strategy, Tactics, and Requirements

### Assumptions

The current Hugo documentation

* is confusing for new users
* is a common complaint in the Hugo forums ([forum discussion 1][ex1], [forum discussion 2][ex2])
* lacks structure and is therefore
    * unscalable, as demonstrated by patch pages (e.g. [here][patch1] and [here][patch2]) that seem out of place, require unnecessary drilldown, or duplicate content in other areas of the docs, thus requiring duplicative efforts to update
    * inconsistent in its terminology, style, and (sometimes) layout
    * limited in effective use of Alogolia's document search (i.e., because of redundant content grouping, headings, etc)
    * difficult to optimize for external search engines (SEO)
* does not leverage Hugo's more powerful feature (e.g., there is only *one* archetype); leveraging these features would help address the aforementioned shortcomings (i.e., scalability, consistency, and search)
* assumes a higher level of Golang proficiency than is realistic for newcomers to static site generators or general web development. A prime example is the sparsity of basic and advanced code samples for templating functions, some of which may still be wholly undocumented.

### Goals

New Hugo documentation should...

* reduce confusion surrounding Hugo concepts; e.g., `list`, `section`, `page`, `kind`, and `content type` with the intention of
    * making it easier for new users to get up and running
    * creating better consistency and scalability for Hugo-dependent projects (viz., [themes.gohugo.io][hugothemes])
    * reducing the frequency of beginner-level questions in the [Hugo Discussion Forum][forum]
* not require, or assume, any degree of Golang proficiency from end users;
    * that said, Hugo can&mdash;and *should*&mdash;act as a bridge for users interested in learning Golang. A implementationn example of this strategy point is the inclusion of `godocref:` as a default front matter field for all function and template pages. See [`archetypes/functions.md`][functionarchetype].
* be easiest to expand and edit for *contributors** but even easier to understand by *end users*.
     If you don't make it *very easy* for authors to contribute to documentation correctly, they will inevitably contribute *incorrectly*;
        * Content modeling is king
        * Go DRY (e.g., by leveraging shortcodes whenever possible)
        * Set required metadata (e.g., via section-specific *archetypes*)
        * Develop contribution guidelines for both development *and* documentation
* be equally accessible via mobile, tablet, desktop, and offline.
* avoid "miscellaneous" sections (e.g.,"Extras"). [This is the last place end users look to get up and running with Hugo](https://discuss.gohugo.io/t/site-with-different-lists-of-sections/5536/3). All content in miscellaneous sections should be edited and incorporated into more logical content groupings (e.g., with the goal of removing *extras* entirely).
* easily scaffold for potential i18n/multilingual versions.

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
* access to static hosting;
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
- [ ] [JSON+LD](https://developers.google.com/schemas/formats/json-ld), [validated](https://search.google.com/structured-data/testing-tool)
- [X] Consistent heading structure
- [X] Semantic HTML5 elements (e.g., `article`, `main`, `aside`, `dl`)
- [X] SSL
- [ ] AMP?
- [ ] 301s [^1]

#### Accessibility

- [ ] Aria roles
- [ ] Alt text for all images

#### Editorial and Content

- [ ] Basic style guide
    - The style guide should facilitate a more consistent UX for the site but not be so complex as to deter documentation contributors
- [X] Contribution guidelines (see [WIP on live site](https://hugodocsconcept.netlify.com/contribute/documentation/))
- [X] Standardized content types (see [WIP archetypes in source](https://github.com/rdwatters/hugo-docs-concept/tree/master/themes/hugodocs/archetypes)
- [X] New content model, including taxonomies ([see tags page][tagspage])
- [ ] DRY. New shortcodes for repeat content (e.g., lists of aliases, page variables, site variables, and others)
- [X] New site architecture and content groupings
- [ ] Single sample website (include in docs source, [`/static/example`](https://github.com/rdwatters/hugo-docs-concept/tree/master/static/example)) for consistent code samples or in-page tutorials

#### Content Strategy Statement

[What is this?](http://contentmarketinginstitute.com/2016/01/content-on-strategy-templates/)

> The Hugo documentation increases the Hugo user base and strengthens the Hugo community by providing intuitive, beginner-friendly, regularly updated usage guides. Hugo documentation makes visitors feel excited and confident that Hugo is the ideal choice for static website development.

#### Editorial Mission

[What is this?](http://contentmarketinginstitute.com/2015/10/statement-content-marketing/)

> The Hugo documentation is a joint effort between the Hugo maintainers and the open-source community. Hugo documentation is designed to promote Hugo, the world's fastest, friendliest, and most extensible static site generator. Hugo documentation is the primary vehicle by which the Hugo team reaches its target audiences. When visitors comes to the Hugo documentation, we want them to install Hugo, develop a new static website with our tool, and share their progress and insights with the Hugo community at large.

## UX/UI

- [X] Copyable code blocks (via highlight.js, extended for hugo-specific keywords)
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
- [X] Custom iconography
- [X] Design assets versioned with source ([see design resources directory][designresources])
- [X] [WCAG color contrast requirements](http://webaim.org/blog/wcag-2-0-and-link-colors/)
- [X] [Sass Guidelines for Source Organization](https://sass-guidelin.es/)
    - [X] Abstracted color palette
    - [X] Abstracted typefaces (multiple open-source fonts available)


## Content Changes

The following is an *abbreviated* listing of *substantive* changes made to the current documentation's source content and organization. Sections here are ordered according to the current site navigation. The changes delimited here do not include copy edits for consistent or preferred usage, improvements in semantics, etc, all of which easily numbers in the thousands, likely more.

### Download Hugo

This is no longer a site navigation link and is instead a button along with "File an Issue" and "Discuss Hugo" at the bottom of the sidebar.

### Site Showcase

Site showcase has stayed more or less as is, including styling, etc. However...
 * The showcase archetype has changed for simplicity.
 * To keep compatibility, all [showcase content files][showcasefiles] have been edited to reflect the new content type. This will also be updated in the ["docs" page of the contribute section](/contribute/documentation/)

### Press & Articles

* The press and article pages has been moved under "News" along with "Release Notes". Also, this whole section is lower on the navigation because it's less frequently visited---I'm assuming---than just about everything on the site.
* Like everything else, I've kept up with changes to the docs upstream on GitHub, but in this case, I also includes a [half dozen *new* articles as well](/news/press-and-articles/).

### About Hugo

### Getting Started

* The [Quick Start][] has been completely updated for more consistent heading structure, etc. Also, **I may delete the "deployment" section of the Quick Start** since this a) adds unnecessary length, making the guide less "quick" and b) detracts from the new "hosting and deployment" section, which offers better advice, and c) is redundant with [Hosting on Github](https://hugodocsconcept.netlify.com/hosting-and-deployment/hosting-on-github/). For example, the Quick Start didn't mention that files already written to public are not necessarily erased at build time. This can cause problems with drafts. I think the other options&mdash;e.g. Arjen's Wercker tutorial&mdash;are more viable and represent better practices for newcomers to Hugo. If future versions of Hugo include baked-in deployment features, I think it's worth reconsidering adding the deployment step back to the Quick Start.

### Content

* This section has been renamed "Content Management" to facilitate elimination of the ["extras"](http://gohugo.io) section. **Note**: this section does *not* include any templating. The convention is `content-management/shortcodes.md` (for explanation and usage) and then `content-management/shortcode-templates.md`

### Themes

### Templates

* Reworked considerably. Page titles have all been changed to reflect their obvious connection to *templating*.
* "Lookup order" page added. The order of the template pages in the main navigation is now such that it could be seen as a sequence of pages showing how to learn templating. Hence the primer, lookup order, and base templates as the first three pages in this section.
* Shortcodes, menus, pagination

### Taxonomies


### Extras

* This section no longer exists in the new documentation site
    * *Extras*, in the content world, is the equivalent of *miscellaneous* or *additional resources*. It's an area that's been tacked onto site navigation to accommodate a *seemingly* disparate set of new features. In other words, READ: "We don't have any idea of where to put this"
* *Extras* pages:
    * **Aliases** Incorporated into [URL Management](/content-management/urls/)
    * **Analytics** Incorporated into [built-in partials](/templates/partials/#using-hugos-built-in-partials)
    * **Builders** This has been removed completely since it has no real added value. The three "builders" mentioned (`new site`, `new theme`, and `new <content>`) are all well-delimited in their respective pages, which is where end users expect to find this type of information in the first place.
    * **Comments** Incorporated into [content management](/content-management/comments/) for content-related pieces and mentioned in [partials](/templates/partials/) for implementation.
    * **Cross-References** Added as its own page under Content Management (`/content-management/cross-references/`)

### Community

### Tutorials

* Original page: <http://gohugo.io/tutorials>
* All installation guides have been consolidated under [/getting-started/installing/]
    * Installing Hugo shouldn't be considered a separate tutorial
    * "Tutorials" is not an intuitive place for end-users to look for this kind of documentation
* All content moved from `/tutorials` edited to reflect a less tutorial-ish style of language (e.g., remove of lines starting with "In this tutorial...")
* Aliases added to new pages and in-page links updated throughout
* Michael Henderson's "Creating a Theme" website ([current][],[concept][]) has been copy edited and content edited to include the new code block shortcodes.

### Troubleshooting

### Tools

### Hugo Cmd Reference

This hasn't been touched. I'll make the necessary style changes once/if the site is integrated into the Hugo GH repo since these pages are pulled automatically using Viper.

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
* Base Templates And Blocks - 20
* Lists in Hugo - 25
* Homepage Template - 30
* Section Templates - 40
* Taxonomy Templates - 50
* Single Page Templates - 60
* Content View Templates - 70
* Data Templates - 80
* Partial Templates - 90
* Shortcode Templates - 100
* Local File Templates - 110
* Custom 404 Page - 120
* Menu Templates - 130
* Pagination - 140
* RSS Templates - 150
* Sitemap Template - 160
* Additional Templating Languages - 170
* Template Debugging - 180

### "functions" Ordering (`.OrderByTitle`)

### "variables-and-Params" Ordering (`.OrderByWeight`)

* _index.md - 01
* Site Variables - 10
* Page Variables - 20
* Taxonomy Variables - 30
* File Variables - 40
* Other Variables - 50

### "hosting-and-deployment" Ordering (`.OrderByWeight`)

* _index.md - 01
* Deployment with Rsync - 10
* Deployment with Wercker - 20
* Hosting on Bitbucket - 30
* Hosting on Github - 40
* Hosting on Gitlab - 50

### "themes" Ordering (`.OrderByWeight`)

* _index.md - 01
* Installing and Using Themes - 10
* Customizing a Theme - 20
* Creating a Theme - 30
* Theme Showcase - 40

### "site-showcase" Ordering (`.OrderByPublishDate`)


### "Troubleshooting" Ordering

**Ordered by title**

## Current Content (Source)

```markdown
{{< readfile file="content/tree.txt" >}}
```

## Proposed Schedule for Hugo Docs Release

If the Hugo Team finds the improvements to the Hugo documentation acceptable, I've proposed the following schedule for leveraging the new Hugo documentation.

1. **2017-02-26** Release to Gitter Channel for dev review
2. **2017-03-05** Post in Discussion Forum for Hugo user feedback
3. **2017-03-6** Pull request/add to Hugo Rep
3. **2017-??-??** Add to Hugo repo for release with v19?

[^1]: As this point, the URL structure has changed considerably. I've been fastidious about adding aliases wherever possible and trying to retain URLs for related content on the current site if applicable. That said, the [current list of aliases is quite large](/contribute/documentation/#be-mindful-of-aliases).

[admonitions]: http://docutils.sourceforge.net/docs/ref/rst/directives.html#admonitions
[designresources]: https://github.com/rdwatters/hugo-docs-concept/tree/master/dev-and-design-resources
[ex1]: https://discuss.gohugo.io/t/frustrated-with-documentation/2810
[ex2]: https://discuss.gohugo.io/t/documentation-restructure-and-design/1891
[forum]: https://discuss.gohugo.io
[functionarchetype]: https://github.com/rdwatters/hugo-docs-concept/blob/master/themes/hugodocs/archetypes/functions.md
[hugothemes]: http://themes.gohugo.io
[patch1]: http://gohugo.io/taxonomies/templates/
[patch2]: https://github.com/spf13/hugo/commit/eaabecf586fd0375585e27c752e05dd8cb4c72b4
[Quick Start]: https://hugodocsconcept.netlify.com/getting-started/quick-start/
[showcasefiles]:
[tagspage]: https://hugodocsconcept.netlify.com/tags/