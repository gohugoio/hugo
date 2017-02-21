<!-- MarkdownTOC -->

- [Assumptions][assumptions]
- [Goals][goals]
- [Audience][audience]
- [Persona][persona]
    - [End User: Developer][end-user-developer]
    - [End User: Themes \(i.e. blogger/author/\)][end-user-themes-ie-bloggerauthor]
- [Requirements][requirements]
    - [Technical][technical]
    - [SEO][seo]
    - [Editorial/Content][editorialcontent]
- [UX][ux]
- [Author Experience \(AX\)][author-experience-ax]
- [Visual Design][visual-design]

<!-- /MarkdownTOC -->

**Updated 2017-02-21**

This is a *very* schlocky version of the documentation I'd put together in my professional life. That said, I think the following pieces are still important and should provide some insight as to how I've approached reworking the Hugo documentation over the last three months.

<a name="assumptions"></a>
## Assumptions

> **Note**: These assumptions are empirical; i.e. the result of me spending a large (and potentially unhealthy) amount of time on the [Hugo Discussion Forum](https://discuss.gohugo.io).

* The current documentation is...
    * confusing for new users
    * a common complaint in the Hugo forums [example discussion 1][],[example discussion 2][]
    * Lacks structure and therefore...
        * doesn't scale (as demonstrated by [patch pages](http://gohugo.io/taxonomies/templates/) that seem out of place or require unnecessary drilldown).
        * is inconsistent in its terminology, style, and (sometimes) layout
        * limits the efficacy of Alogolia's document search feature
    * Does not itself leverage Hugo's more powerful feature (e.g., there is only *one* archetype and five shortcodes for what is ultimately a complex documentation site).
        * Leveraging these features would help address the aforementioned shortcomings (i.e., scalability, consistency, and search)
    * Assumes a higher level of Golang proficiency than is realistic for newcomers to the language or web development in general.
* If you don't make it *very easy* for authors to contribute to documentation correctly, they will inevitably contribute *incorrectly*.
    * Content modeling is king
    * Go DRY (e.g. with shortcodes)
    * Require metadata

<a name="goals"></a>
## Goals

Hugo documentation should...

* reduce confusion surrounding `list` vs `section` vs `page` vs `content type`,etc., and thereby
    * make it easier for new users to get up and running
    * create better consistency and scalability Hugo-dependent projects (namely, http://themes.gohugo.io)
    * reduce frequency of such questions in the Hugo Discuss Forum
* not require any degree of Golang proficiency from end users.
    * That said, Hugo *can* act as a bridge for those interested in learning Golang (e.g., by including `godocref` as a default front matter field in [`archetypes/functions.md`](https://github.com/rdwatters/hugo-docs-concept/blob/master/themes/hugodocs/archetypes/functions.md)).
* be easy to expand and edit for contributors Editing and expanding documentation should be easiest for *contributors**, whereas usage of documentation should be easiest for *end users*.
* be equally accessible via mobile, tablet, desktop, *and* offline.
* not include an "extras" section because [this is the last place end users look to learn about Hugo](https://discuss.gohugo.io/t/site-with-different-lists-of-sections/5536/3). Instead all "extras" should be integrated into a new
* easily scaffold for future multilingual versions

<a name="audience"></a>
## Audience

* Primary: Web developers interested in static site generators
* Secondary: Web publishers (bloggers, authors)
* Tertiary: Web developers interested in learning Golang

<a name="persona"></a>
## Persona

<a name="end-user-developer"></a>
### End User: Developer

* Limited proficiency in Git and DVCS
* No to little proficiency in Golang
* working proficiency in front-end development---HTML, CSS, JS---but not necessarily front-end build tools
* familiarity with at least one double-curly templating language (e.g., liquid, Twig, Swig, or Django)
* proficiency in the English language
    * proficiency in other languages (for future multilingual versions)

<a name="end-user-themes-ie-bloggerauthor"></a>
### End User: Themes (i.e. blogger/author/)

* Limited proficiency in the command line/prompt
* Proficiency in a supported content format (specifically markdown)
* Access to static hosting but with limited proficiency in basic deployments

<a name="requirements"></a>
## Requirements

<a name="technical"></a>
### Technical

- [X] Built with Hugo
- [X] Performant (e.g., 80+ on [Google Page Speed Score](https://developers.google.com/speed/pagespeed/insights/))
- [X] Front-end build tools for concatenation, minification
- [ ] CDN

<a name="seo"></a>
### SEO

- [X] [Open Graph Protocol](http://ogp.me/)
- [X] [schema.org](http://schema.org)
- [ ] [JSON+LD](https://developers.google.com/schemas/formats/json-ld), [validated](https://search.google.com/structured-data/testing-tool)
- [X] Consistent heading structure
- [X] Semantic HTML5 elements (e.g., `article`, `main`, `aside`, `dl`)
- [ ] SSL
- [ ] AMP

<a name="editorialcontent"></a>
### Editorial/Content

- [X] Basic style guide
- [X] Contribution guidelines (see [working draft on live site](https://hugodocsconcept.netlify.com/contribute-to-hugo/contribute-to-the-hugo-docs/))
- [X] Standardized content types (i.e, [see current archetypes](https://github.com/rdwatters/hugo-docs-concept/tree/master/themes/hugodocs/archetypes)
- [X] New Content Model, including taxonomies ([see tags page]())
- [ ] DRY. New shortcodes for repeat content (e.g., list of aliases, list of page variables)
- [X] New site architecture and content groupings
- [ ] Examples pulling from a single sample website (including in docs source) for consistent code samples or in-page tutorials

#### [Content Strategy Statement](http://contentmarketinginstitute.com/2016/01/content-on-strategy-templates/)

> The Hugo documentation increases the Hugo user base and strengthens the Hugo community by providing intuitive, beginner-friendly content that makes visitors to the site feel excited and confident that Hugo is the ideal choice for all their static web publishing needs.


#### [Editorial Mission](http://contentmarketinginstitute.com/2015/10/statement-content-marketing/)

> The Hugo documentation is a joint effort between the Hugo maintainers and the open-source community. Hugo documentation is designed to promote Hugo, the world's fastest, friendliest, and most extensible static site generator. Hugo documentation is the primary vehicle by which the Hugo team reaches our target audiences. When visitors comes to our site, we want them to install Hugo, developer a new site in Hugo, and share their progress with the community at large.

<a name="ux"></a>
## UX

- [ ] Share buttons
- [ ] Copy-page links
- [X] Copyable code blocks (via highlight.js, extended for hugo-specific keywords)
- [X] Dual in-page navigation
- [X] Smooth scrolling

<a name="author-experience-ax"></a>
## Author Experience (AX)

- [X] Easy scaffolding of content types (CLI)
- [X] Type-based content storage model and scope (archetypes)

<a name="visual-design"></a>
## Visual Design

- [X] Clean typography with open-source font
    - [X] Optimal line length (50-80 characters)
    - [X] Consistent vertical rhythm
- [X] Responsive
    - [X] Flexbox
    - [X] Typography (via ems)
- [X] Custom iconography
- [X] Design assets versioned with source ([see design resources directory](https://github.com/rdwatters/hugo-docs-concept/tree/master/dev-and-design-resources))
- [X] [WCAG color contrast requirements](http://webaim.org/blog/wcag-2-0-and-link-colors/)
- [X] [Sass Guidelines for Source Organization](https://sass-guidelin.es/)
    - [X] Abstracted color palette
    - [X] Abstracted typefaces (multiple open-source fonts available)

[example discussion 1]: https://discuss.gohugo.io/t/frustrated-with-documentation/2810
[example discussion 2]: https://discuss.gohugo.io/t/documentation-restructure-and-design/1891
