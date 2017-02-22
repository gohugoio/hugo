# Hugo Docs Strategy, Tactics, and Requirements

> **Note:** Before any of my fellow content strategists banish me to content strategy hell, know that I *know* this is a *very schlocky* version of the documentation required for a real content strategy.

**Updated 2017-02-21**

- [Assumptions](#assumptions)
- [Goals](#goals)
- [Audience](#audience)
- [Persona](#persona)
    - [End User: Developer](#end-user-developer)
    - [End User: Themes \(i.e. blogger/author/\)](#end-user-themes-ie-bloggerauthor)
- [Requirements](#requirements)
    - [Technical](#technical)
    - [SEO](#seo)
    - [Editorial/Content](#editorialcontent)
- [UX/UI](#uxui)
- [Author Experience \(AX\)](#author-experience-ax)
- [Analytics/Metrics](#analyticsmetrics)
- [Visual Design](#visual-design)

## Assumptions

> **Note**: These assumptions are empirical; i.e. the result of me spending a large (and potentially unhealthy) amount of time on the [Hugo Discussion Forum](https://discuss.gohugo.io). Google analytics *may* provide more quantitative insight into actual Hugo docs usage. I should also mention that these are *assumptions* and not *criticisms*. I **LOVE** Hugo.

* The current documentation is
    * confusing for new users
    * a common complaint in the Hugo forums ([example discussion 1][],[example discussion 2][])
    * lacks structure and therefore
        * doesn't scale, as demonstrated by [patch pages](http://gohugo.io/taxonomies/templates/) that seem out of place or require unnecessary drilldown)
        * is inconsistent in its terminology, style, and (sometimes) layout
        * limits the efficacy of Alogolia's document search feature through redundant content groups, headings, etc
        * does not leverage SEO for external search engines
    * does not leverage Hugo's more powerful feature (e.g., there is only *one* archetype); leveraging these features would help address the aforementioned shortcomings (i.e., scalability, consistency, and search)
    * assumes a higher level of Golang proficiency than is realistic for newcomers to the Golang programming language or to web development in general.
* If you don't make it *very easy* for authors to contribute to documentation correctly, they will inevitably contribute *incorrectly*;
    * content modeling is king
    * go DRY (e.g., with shortcodes)
    * set required metadata

## Goals

Hugo documentation should...

* reduce confusion surrounding Hugo concepts as `list`, `section`, `page`, `content type`, etc. and thereby
    * make it easier for new users to get up and running
    * create better consistency and scalability for Hugo-dependent projects (viz., http://themes.gohugo.io)
    * reduce frequency of questions surrounding said concepts in the Hugo Discuss Forum
* not require or assume any degree of Golang proficiency from end users;
    * that said, Hugo can&mdash;and *should*&mdash;act as a bridge for users interest in learning Golang (e.g., by including `godocref` as a default front matter field. See [`archetypes/functions.md`][functionarchetype].
* be easiest to expand and edit for *contributors**, but even easier to understand by *end users*.
* be equally accessible via mobile, tablet, desktop, *and* offline.
* not include an "extras" section because [this is the last place end users look to learn about Hugo](https://discuss.gohugo.io/t/site-with-different-lists-of-sections/5536/3). Instead all "extras" should be integrated into a new
* easily scaffold for future multilingual versions

## Audience

* Primary: Web developers interested in static site generators
* Secondary: Web publishers (bloggers, authors)
* Tertiary: Web developers interested in learning Golang

## Persona

### End User: Developer

* Limited proficiency in Git and DVCS
* No to little proficiency in Golang
* working proficiency in front-end development---HTML, CSS, JS---but not necessarily front-end build tools
* familiarity with at least one double-curly templating language (e.g., liquid, Twig, Swig, or Django)
* proficiency in the English language
    * proficiency in other languages (for future multilingual versions)

### End User: Themes (i.e. blogger/author/)

* Limited proficiency in the command line/prompt
* Proficiency in a supported content format (specifically markdown)
* Access to static hosting but with limited proficiency in basic deployments

## Requirements

### Technical

- [X] Built with Hugo
- [X] Performant (e.g., 80+ [Google Page Speed Score](https://developers.google.com/speed/pagespeed/insights/?url=https%3A%2F%2Fhugodocsconcept.netlify.com%2Fabout-hugo))
- [X] Front-end build tools for concatenation, minification
- [X] Browser compatibility: modern (i.e. Chrome, Edge, Firefox, Safari) and IE11
- [ ] AMP?
- [ ] CDN

### SEO

- [X] [Open Graph Protocol](http://ogp.me/)
- [X] [schema.org](http://schema.org)
- [ ] [JSON+LD](https://developers.google.com/schemas/formats/json-ld), [validated](https://search.google.com/structured-data/testing-tool)
- [X] Consistent heading structure
- [X] Semantic HTML5 elements (e.g., `article`, `main`, `aside`, `dl`)
- [X] SSL
- [ ] AMP?

### Editorial/Content

- [X] Basic style guide
- [X] Contribution guidelines (see [working draft on live site](https://hugodocsconcept.netlify.com/contribute-to-hugo/contribute-to-the-hugo-docs/))
- [X] Standardized content types (i.e, [see current archetypes](https://github.com/rdwatters/hugo-docs-concept/tree/master/themes/hugodocs/archetypes)
- [X] New content model, including taxonomies ([see tags page][tagspage])
- [ ] DRY. New shortcodes for repeat content (e.g., lists of aliases, page variables, site variables, and others)
- [X] New site architecture and content groupings
- [ ] Single sample website (include in docs source, [`/static/example`](https://github.com/rdwatters/hugo-docs-concept/tree/master/static/example)) for consistent code samples or in-page tutorials

#### [Content Strategy Statement](http://contentmarketinginstitute.com/2016/01/content-on-strategy-templates/)

> The Hugo documentation increases the Hugo user base and strengthens the Hugo community by providing intuitive, beginner-friendly content that makes visitors to the site feel excited and confident that Hugo is the ideal choice for all their static web publishing needs.

#### [Editorial Mission](http://contentmarketinginstitute.com/2015/10/statement-content-marketing/)

> The Hugo documentation is a joint effort between the Hugo maintainers and the open-source community. Hugo documentation is designed to promote Hugo, the world's fastest, friendliest, and most extensible static site generator. Hugo documentation is the primary vehicle by which the Hugo team reaches our target audiences. When visitors comes to our site, we want them to install Hugo, developer a new site in Hugo, and share their progress with the community at large.

## UX/UI

- [X] Copyable code blocks (via highlight.js, extended for hugo-specific keywords)
- [X] Dual in-page navigation (i.e. for site nav and in-page TOC)
- [X] Smooth scrolling
- [X] [RTD-style admonitions][admonitions] (see [example admonition shortcode](https://github.com/rdwatters/hugo-docs-concept/blob/master/layouts/shortcodes/note.html) and [examples on published site](http://localhost:1313/contribute-to-hugo/contribute-to-the-hugo-docs/#admonition-short-codes))
- [ ] Share buttons: Reddit, Twitter, LinkedIn, and "Copy Page Url" (the last of which has strongest utility for referencing in forums)

## Author Experience (AX)

- [X] Easy scaffolding of content types (CLI)
- [X] Type-based content storage model and scope (archetypes)

## Analytics/Metrics

- [X] Google Analytics
- [ ] Content groupings in GA to measure usage and define content gaps
- [ ] Automated reports

> **Note:** These are separate from usage statics re: Hugo downloads, `.Hugo.Generator`, etc.

## Visual Design

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

[admonitions]: http://docutils.sourceforge.net/docs/ref/rst/directives.html#admonitions
[designresources]: https://github.com/rdwatters/hugo-docs-concept/tree/master/dev-and-design-resources
[example discussion 1]: https://discuss.gohugo.io/t/frustrated-with-documentation/2810
[example discussion 2]: https://discuss.gohugo.io/t/documentation-restructure-and-design/1891
[functionarchetype]: https://github.com/rdwatters/hugo-docs-concept/blob/master/themes/hugodocs/archetypes/functions.md
[tagspage]: https://hugodocsconcept.netlify.com/tags/