# Content Reorganization

- [Changes to Existing Content Sections](#changes-to-existing-sections)
    - [Extras](#extras)
    - [Tutorials](#tutorials)
    - [Getting Started](#getting-started)
- [Content Organization: \(Site Navigation\](#content-organization-site-navigation)
- [Content Organization: \(Source\)](#content-organization-source)

## Changes to Existing Sections

The following is an *abbreviated* listing and only includes the *larger* changes to content organization. These changes do not include copy edits for consistent usage, which easily numbers in the thousands at this point.

### [Extras](http://gohugo.io/extras)

* This section no longer exists in the new documentation site
    * *Extras*, in the content world, is the equivalent of *miscellaneous* or *additional resources*. READ: "We don't have any idea of where to put this"
* Previous pages in extras are now in the following locations:
    * **Aliases** Incorporated into `/content-management/url-management/`
    * **Analytics** Incorporated into /templates/partial-templates/#built-in

### [Tutorials](http://gohugo.io/tutorials)

* Moved all installation guides to /getting-started/install-hugo/
    * Installing Hugo shouldn't be considered a separate tutorial
    * "Tutorials" is not an intuitive place for end-users to look for this kind of documentation
* All content moved from `/tutorials` edited to reflect a less tutorial-ish style of language (e.g., remove of lines starting with "In this tutorial...")
* Aliases added to new pages and in-page links updated throughout

## Content Organization: Site Navigation

The following is a list of weights for the newly restructure site architecture

## Content Organization: Source

**Updated 2017-02-21**

```
.
├── _index.md
├── about-hugo
│   ├── _index.md
│   ├── benefits-of-static.md
│   ├── hugo-features.md
│   ├── license.md
│   ├── roadmap.md
│   ├── what-is-hugo.md
│   └── why-i-built-hugo.md
├── commands
│   └── _index.md
├── content-management
│   ├── _index.md
│   ├── archetypes.md
│   ├── content-organization.md
│   ├── content-sections.md
│   ├── content-summaries.md
│   ├── content-types.md
│   ├── cross-references.md
│   ├── front-matter.md
│   ├── menus.md
│   ├── multilingual-mode.md
│   ├── shortcodes.md
│   ├── supported-content-formats.md
│   ├── table-of-contents.md
│   ├── taxonomies.md
│   └── url-management.md
├── contribute-to-hugo
│   ├── _index.md
│   ├── add-your-site-to-the-showcase.md
│   ├── contribute-to-hugo-development.md
│   └── contribute-to-the-hugo-docs.md
├── developer-tools
│   ├── _index.md
│   ├── migrate-to-hugo.md
│   └── syntax-highlighting.md
├── functions
│   ├── _index.md
│   ├── abslangurl.md
│   ├── absurl.md
│   ├── after.md
│   ├── apply.md
│   ├── base64decode.md
│   ├── base64encode.md
│   ├── chomp.md
│   ├── countrunes.md
│   ├── countwords.md
│   ├── dateformat.md
│   ├── default-function.md
│   ├── delimit.md
│   ├── dict.md
│   ├── echoparam.md
│   ├── emojify.md
│   ├── findre.md
│   ├── first.md
│   ├── get.md
│   ├── getenv.md
│   ├── getpage.md
│   ├── haschildren.md
│   ├── hasmenucurrent.md
│   ├── hasprefix.md
│   ├── highlight.md
│   ├── htmlescape.md
│   ├── htmlunescape.md
│   ├── humanize.md
│   ├── i18n.md
│   ├── imageconfig.md
│   ├── in.md
│   ├── index-function.md
│   ├── int.md
│   ├── intersect.md
│   ├── ismenucurrent.md
│   ├── isset.md
│   ├── jsonify.md
│   ├── last.md
│   ├── lower.md
│   ├── markdownify.md
│   ├── math.md
│   ├── md5.md
│   ├── param.md
│   ├── partialcached.md
│   ├── plainify.md
│   ├── pluralize.md
│   ├── printf.md
│   ├── querify.md
│   ├── range.md
│   ├── readdir.md
│   ├── readfile.md
│   ├── rel.md
│   ├── rellangurl.md
│   ├── relref.md
│   ├── relurl.md
│   ├── render.md
│   ├── replace.md
│   ├── safecss.md
│   ├── safehtml.md
│   ├── safehtmlattr.md
│   ├── safejs.md
│   ├── safeurl.md
│   ├── scratch.md
│   ├── seq.md
│   ├── sha1.md
│   ├── sha256.md
│   ├── shuffle.md
│   ├── singularize.md
│   ├── slice.md
│   ├── slicestr.md
│   ├── sort.md
│   ├── split.md
│   ├── string.md
│   ├── substr.md
│   ├── the-dot.md
│   ├── time.md
│   ├── title.md
│   ├── trim.md
│   ├── unix.md
│   ├── upper.md
│   ├── urlize.md
│   ├── where.md
│   └── with.md
├── getting-started
│   ├── _index.md
│   ├── basic-usage.md
│   ├── configuration.md
│   ├── directory-structure.md
│   ├── install-hugo.md
│   ├── quick-start.md
│   └── using-the-hugo-docs.md
├── hosting-and-deployment
│   ├── _index.md
│   ├── deployment-with-rsync.md
│   ├── deployment-with-wercker.md
│   ├── hosting-on-bitbucket.md
│   ├── hosting-on-github.md
│   └── hosting-on-gitlab.md
├── mailing-list.md
├── news
│   ├── _index.md
│   ├── press-and-articles.md
│   └── release-notes.md
├── showcase
│   ├── 2626info.md
│   ├── _index.md
│   ├── antzucaro.md
│   ├── appernetic.md
│   ├── arresteddevops.md
│   ├── asc.md
│   ├── astrochili.md
│   ├── aydoscom.md
│   ├── barricade.md
│   ├── bepsays.md
│   ├── bugtrackers.io.md
│   ├── camunda-blog.md
│   ├── camunda-docs.md
│   ├── cdnoverview.md
│   ├── chinese-grammar.md
│   ├── chingli.md
│   ├── chipsncookies.md
│   ├── christianmendoza.md
│   ├── cinegyopen.md
│   ├── clearhaus.md
│   ├── cloudshark.md
│   ├── coding-journal.md
│   ├── consequently.md
│   ├── ctlcompiled.md
│   ├── danmux.md
│   ├── datapipelinearchitect.md
│   ├── davidepetilli.md
│   ├── davidrallen.md
│   ├── davidyates.md
│   ├── devmonk.md
│   ├── dmitriid.com.md
│   ├── emilyhorsman.com.md
│   ├── esolia-com.md
│   ├── esolia-pro.md
│   ├── eurie.md
│   ├── fale.md
│   ├── fixatom.md
│   ├── fxsitecompat.md
│   ├── gntech.md
│   ├── gogb.md
│   ├── goin5minutes.md
│   ├── h10n.me.md
│   ├── hugo.md
│   ├── jamescampbell.md
│   ├── jorgennilsson.md
│   ├── kieranhealy.md
│   ├── klingt-net.md
│   ├── launchcode5.md
│   ├── leepenney.md
│   ├── leowkahman.md
│   ├── lk4d4.darth.io.md
│   ├── losslesslife.md
│   ├── mariosanchez.md
│   ├── mayan-edms.md
│   ├── michaelwhatcott.md
│   ├── mongodb-eng-journal.md
│   ├── mtbhomer.md
│   ├── nickoneill.md
│   ├── ninjaducks.in.md
│   ├── ninya.io.md
│   ├── nodesk.md
│   ├── novelist-xyz.md
│   ├── npf.md
│   ├── peteraba.md
│   ├── rahulrai.md
│   ├── rakutentech.md
│   ├── rdegges.md
│   ├── readtext.md
│   ├── richardsumilang.md
│   ├── rick-cogley-info.md
│   ├── ridingbytes.md
│   ├── robertbasic.md
│   ├── scottcwilson.md
│   ├── shapeshed.md
│   ├── shelan.md
│   ├── silvergeko.md
│   ├── softinio.md
│   ├── spf13.md
│   ├── steambap.md
│   ├── stefano.chiodino.md
│   ├── stou.md
│   ├── szymonkatra.md
│   ├── techmadeplain.md
│   ├── tendermint.md
│   ├── thecodeking.md
│   ├── thehome.md
│   ├── tutorialonfly.md
│   ├── ucsb.md
│   ├── upbeat.md
│   ├── vamp.md
│   ├── viglug.org.md
│   ├── vurt.co.md
│   ├── yslow-rules.md
│   ├── ysqi.md
│   └── yulinling.net.md
├── templates
│   ├── _index.md
│   ├── ace-templating.md
│   ├── amber-templating.md
│   ├── base-templates-and-blocks.md
│   ├── content-view-templates.md
│   ├── custom-404-page.md
│   ├── data-templates.md
│   ├── go-template-primer.md
│   ├── homepage-template.md
│   ├── list-and-section-templates.md
│   ├── local-file-templates.md
│   ├── menu-templates.md
│   ├── pagination.md
│   ├── partial-templates.md
│   ├── rss-templates.md
│   ├── shortcode-templates.md
│   ├── single-page-templates.md
│   ├── sitemap-template.md
│   ├── taxonomy-templates.md
│   └── template-debugging.md
├── themes
│   ├── _index.md
│   ├── creating-a-theme.md
│   ├── customizing-a-theme.md
│   ├── installing-and-using-themes.md
│   └── theme-showcase.md
├── troubleshooting
│   ├── _index.md
│   ├── accented-characters-in-urls.md
│   └── eof-error.md
├── tutorials
│   ├── _index.md
│   ├── create-a-multilingual-site.md
│   ├── creating-a-new-theme.md
│   └── migrate-from-jekyll-to-hugo.md
└── variables-and-params
    ├── _index.md
    ├── file-variables.md
    ├── page-variables.md
    ├── shortcode-git-and-hugo-variables.md
    ├── site-variables.md
    └── taxonomy-variables.md

15 directories, 264 files
```