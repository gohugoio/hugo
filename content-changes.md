# Annotated Content Reorganization

- [Changes to Existing Content Sections](#changes-to-existing-sections)
    - [Extras](#extras)
    - [Tutorials](#tutorials)
    - [Getting Started](#getting-started)
- [Content Organization: \(Site Navigation\](#content-organization-site-navigation)
- [Content Organization: \(Source\)](#content-organization-source)

## Changes to Existing Sections

The following is an *abbreviated* listing and only includes the *larger* changes to content organization. Everything is ordered according to [existing site structure and site navigation](http://gohugo.io/overview/introduction/). These changes do not include copy edits for consistent usage, which easily numbers in the thousands, if not more.

### Download Hugo

This is no longer a site navigation link and is instead a button along with "File and Issue" and "Discuss Hugo".

### Site Showcase

### Press & Articles

### About Hugo

### Getting Started

* The [Quick Start][] has been completely updated for more consistent heading structure, etc. Also, **I may delete the "deployment" section of the Quick Start** since this a) adds unnecessary length, making the guide less "quick" and b) detracts from the new "hosting and deployment" section, which offers better advice, and c) is redundant with [Hosting on Github](https://hugodocsconcept.netlify.com/hosting-and-deployment/hosting-on-github/). For example, the Quick Start didn't mention that files already written to public are not necessarily erased at build time. This can cause problems with drafts. I think the other options&mdash;e.g. Arjen's Wercker tutorial&mdash;are more viable and represent better practices for newcomers to Hugo. If future versions of Hugo include baked-in deployment features, I think it's worth reconsidering adding the deployment step back to the Quick Start.

### Content

### Themes

### Templates

### Taxonomies


### [Extras](http://gohugo.io/extras)

* This section no longer exists in the new documentation site
    * *Extras*, in the content world, is the equivalent of *miscellaneous* or *additional resources*. READ: "We don't have any idea of where to put this"
* Previous pages in extras are now in the following locations:
    * **Aliases** Incorporated into `/content-management/url-management/`
    * **Analytics** Incorporated into /templates/partial-templates/#built-in

### Community

### [Tutorials](http://gohugo.io/tutorials)

* Moved all installation guides to /getting-started/install-hugo/
    * Installing Hugo shouldn't be considered a separate tutorial
    * "Tutorials" is not an intuitive place for end-users to look for this kind of documentation
* All content moved from `/tutorials` edited to reflect a less tutorial-ish style of language (e.g., remove of lines starting with "In this tutorial...")
* Aliases added to new pages and in-page links updated throughout

### Troubleshooting

### Tools

### Hugo Cmd Reference

This hasn't been touched. I'll make the necessary style changes once/if the site is integrated into the Hugo GH repo since these pages are pulled automatically using Viper.

### Issues & Help

This is no longer a site navigation link and is instead a button along with "Download" and "Discuss Hugo".

## Content Organization: Site Navigation

The following shows weights and ordering for the newly restructured site architecture.

### "About Hugo" Ordering (- weight)

_index - 01
what is hugo - 10
hugo features - 20
the benefits of static - 30
why i built hugo - 40
roadmap - 50
apache license - 60

### "Getting Started" Ordering (- weight)

* _index.md - 01
* quick start - 10
* using the hugo docs - 20
* install from source - 30
* install on linux - 40
* install on mac - 50
* install on pc - 60
* basic usage - 70
* directory structure - 80
* configuration - 90

### "Content Management" Ordering (- weight)

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

### "Templates" Ordering (- weight)

* _index - 01
* go template primer - 10
* base templates and blocks - 20
* lists in Hugo - 25
* homepage template - 30
* section templates - 40
* taxonomy templates - 50
* single page templates - 60
* content view templates - 70
* data templates - 80
* partial templates - 90
* shortcode templates - 100
* local file templates - 110
* custom 404 page - 120
* menu templates - 130
* pagination - 140
* rss templates - 150
* sitemap template - 160
* additional templating languages - 170
* template debugging - 180

## "Functions" Ordering

**Ordered by title (note that `.Title` is all lowercase, whereas `.Linktitle` is used for proper casing)**

### "Variables and Params" Ordering (- weight)

* _index.md - 01
* site variables - 10
* page variables - 20
* taxonomy variables - 30
* file variables - 40
* shortcode git and huge variables - 50

### "Hosting and Deployment" Ordering (- weight)

* _index.md - 01
* deployment with rsync - 10
* deployment with wercker - 20
* hosting on bitbucket - 30
* hosting on github - 40
* hosting on gitlab - 50

### "Themes" Ordering (- weight)

* _index.md - 01
* installing and using themes - 10
* customizing a theme - 20
* creating a them - 30
* theme showcase - 30

### "Site Showcase" Ordering

**Ordered by `.PublishDate`**

### "Themes" Ordering (- weight)

* _index.md - 01
* Installing and Using Themes - 10
* Customizing a Theme - 20
* Creating a Theme - 30
* Theme Showcase - 40

### "Troubleshooting" Ordering

**Ordered by title**

## Content Organization: Source

**[See tree.md at the root of this repository](tree.md).**



[Quick Start]: https://hugodocsconcept.netlify.com/getting-started/quick-start/