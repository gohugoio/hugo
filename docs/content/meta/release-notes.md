---
aliases:
- /doc/release-notes/
- /meta/release-notes/
date: 2013-07-01
menu:
  main:
    parent: about
title: Release Notes
weight: 10
---

## **0.13.0** Feb 21, 2015

The v0.13.0 release is the largest Hugo release to date. The release introduced
some long sought after features (pagination, sequencing, data loading, tons of
template improvements) as well as major internal improvements. In addition to
the code changes, the Hugo community has grown significantly and now has over
3000 stars on github, 134 contributors, 24 themes and 1000s of happy users.

This release represents **448 contributions by 65 contributors**

A special shout out to [@bep](https://github.com/bep) and
[@anthonyfok](https://github.com/anthonyfok) for their new role as Hugo
maintainers and their tremendous contributions this release.

### New major features
* Support for [data files](/extras/datafiles/) in [YAML](http://yaml.org/),
  [JSON](http://www.json.org/), or [TOML](https://github.com/toml-lang/toml)
  located in the `data` directory ([#885][])
* Support for [dynamic content](/extras/dynamiccontent/) by loading JSON & CSV
  from remote sources via GetJson and GetCsv in short codes or other layout
  files ([#748][])
* [Pagination support](/extras/pagination/) for home page, sections and
  taxonomies ([#750][])
* Universal sequencing support
    * A new, generic Next/Prev functionality is added to all lists of pages
      (sections, taxonomies, etc.)
    * Add in-section [Next/Prev](/templates/variables/) content pointers
* `Scratch` -- [a "scratchpad"](/extras/scratch) for your node- and page-scoped
  variables
* [Cross Reference](/extras/crossreferences/) support to easily link documents
  together with the ref and relref shortcodes. 
* [Ace](http://ace.yoss.si/) template engine support ([#541][])
* A new [shortcode](/extras/shortcodes/) token of `{{</* */>}}` (raw HTML)
  alongside the existing `{{%/* */%}}` (Markdown)
* A top level `Hugo` variable (on Page & Node) is added with various build
  information
* Several new ways to order and group content:
    * `ByPublishDate`
    * `GroupByPublishDate(format, order)`
    * `GroupByParam(key, order)`
    * `GroupByParamDate(key, format, order)`
* Hugo has undergone a major refactoring, with a new handler system and a
  generic file system. This sounds and is technical, but will pave the way for
  new features and make Hugo even speedier

### Notable enhancements to existing features

* The [shortcode](/extras/shortcodes/) handling is rewritten for speed and
  better error messages. 
* Several improvements to the [template functions](/templates/functions/):
    * `where` is now even more powerful and accepts SQL-like syntax with the
      operators `==`, `eq`; `!=`, `<>`, `ne`; `>=`, `ge`; `>`, `gt`; `<=`,
      `le`; `<`, `lt`; `in`, `not in`
    * `where` template function now also accepts dot chaining key argument
      (e.g. `"Params.foo.bar"`)
* New template functions:
    * `apply`
    * `chomp`
    * `delimit`
    * `sort`
    * `markdownify`
    * `in` and `intersect` 
    * `trim`
    * `replace`
    * `dateFormat`
* Several [configurable improvements related to Markdown
  rendering](/overview/configuration/#configure-blackfriday-rendering:a66b35d20295cb764719ac8bd35837ec):
    * Configuration of footnote rendering
    * Optional support for smart angled quotes, e.g. `"Hugo"` → «Hugo»
    * Enable descriptive header IDs
* URLs in XML output is now correctly canonified ([#725][], [#728][], and part
  of [#789][])

### Other improvements

* Internal change to use byte buffer pool significantly lowering memory usage
  and providing measurable performance improvements overall
* Changes to docs:
    * A new [Troubleshooting](/troubleshooting/overview/) section is added
    * It's now searchable through Google Custom Search ([#753][])
    * Some new great tutorials:
        * [Automated deployments with
          Wercker](/tutorials/automated-deployments/)
        * [Creating a new theme](/tutorials/creating-a-new-theme/)
* [`hugo new`](/content/archetypes/) now copies the content in addition to the front matter
* Improved unit test coverage
* Fixed a lot of Windows-related path issues
* Improved error messages for template and rendering errors
* Enabled soft LiveReload of CSS and images ([#490][])
* Various fixes in RSS feed generation ([#789][])
* `HasMenuCurrent` and `IsMenuCurrent` is now supported on Nodes
* A bunch of [bug fixes](https://github.com/spf13/hugo/commits/master)

[#490]: https://github.com/spf13/hugo/pull/490 "Pull Request #490: Livereload CSS and images without browser refresh"
[#541]: https://github.com/spf13/hugo/pull/541 "Pull Request #541: Add Ace template engine support"
[#725]: https://github.com/spf13/hugo/issues/725 "Issue #725: CanonifyUrls does not canonicalize urls in RSS"
[#728]: https://github.com/spf13/hugo/issues/728 "Pull Request #728: Add ability to canonify URLs in rendered XML output."
[#748]: https://github.com/spf13/hugo/issues/748 "Feature: GetJson and GetJson in short codes or other layout files"
[#750]: https://github.com/spf13/hugo/issues/750 "Pull Request: Add pagination support for home page, sections and taxonomies"
[#753]: https://github.com/spf13/hugo/issues/753 "Add search to documentation"
[#789]: https://github.com/spf13/hugo/issues/789 "Issue #789: RSS feeds do not validate"
[#885]: https://github.com/spf13/hugo/issues/885 "Feature/datadir"


## **0.12.0** Sept 1, 2014

A lot has happened since Hugo v0.11.0 was released. Most of the work has been
focused on polishing the theme engine and adding critical functionality to the
templates.

This release represents over 90 code commits from 28 different contributors.

  * 10 [new themes](https://github.com/spf13/hugoThemes) created by the community
  * Fully themable [Partials](/templates/partials/)
  * [404 template](/templates/404/) support in themes
  * [Shortcode](/extras/shortcodes/) support in themes
  * [Views](/templates/views/) support in themes
  * Inner [shortcode](/extras/shortcodes/) content now treated as Markdown
  * Support for header ids in Markdown (# Header {#myid})
  * [Where](/templates/list/) template function to filter lists of content, taxonomies, etc.
  * [GroupBy](/templates/list/) & [GroupByDate](/templates/list/) methods to group pages
  * Taxonomy [pages list](/taxonomies/methods/) now sortable, filterable, limitable & groupable
  * General cleanup to taxonomies & documentation to make it more clear and consistent
  * [Showcase](/showcase/) returned and has been expanded
  * Pretty links now always have trailing slashes
  * [BaseUrl](/overview/configuration/) can now include a subdirectory
  * Better feedback about draft & future post rendering
  * A variety of improvements to [the website](http://gohugo.io/)

## **0.11.0** May 28, 2014

This release represents over 110 code commits from 29 different contributors.

  * Considerably faster... about 3 - 4x faster on average
  * [LiveReload](/extras/livereload/). Hugo will automatically reload the browser when the build is complete
  * Theme engine w/[Theme Repository](https://github.com/spf13/hugoThemes)
  * [Menu system](/extras/menus/) with support for active page
  * [Builders](/extras/builders/) to quickly create a new site, content or theme
  * [XML sitemap](/templates/sitemap/) generation
  * [Integrated Disqus](/extras/comments/) support
  * Streamlined [template organization](/templates/overview/)
  * [Brand new docs site](http://gohugo.io/)
  * Support for publishDate which allows for posts to be dated in the future
  * More [sort](/content/ordering/) options
  * Logging support
  * Much better error handling
  * More informative verbose output
  * Renamed Indexes > [Taxonomies](/taxonomies/overview/)
  * Renamed Chrome > [Partials](/templates/partials/)

## **0.10.0** March 1, 2014

This release represents over 110 code commits from 29 different contributors.

  * [Syntax highlighting](/extras/highlighting/) powered by pygments (**slow**)
  * Ability to [sort content](/content/ordering/) many more ways
  * Automatic [table of contents](/extras/toc/) generation
  * Support for Unicode URLs, aliases and indexes
  * Configurable per-section [permalink](/extras/permalinks/) pattern support
  * Support for [paired shortcodes](/extras/shortcodes/)
  * Shipping with some [shortcodes](/extras/shortcodes/) (highlight & figure)
  * Adding [canonify](/extras/urls/) option to keep urls relative
  * A bunch of [additional template functions](/layout/functions/)
  * Watching very large sites now works on Mac
  * RSS generation improved. Limited to 50 items by default, can limit further in [template](/layout/rss/)
  * Boolean params now supported in [frontmatter](/content/front-matter/)
  * Launched website [showcase](/showcase/). Show off your own hugo site!
  * A bunch of [bug fixes](https://github.com/spf13/hugo/commits/master)

## **0.9.0** November 15, 2013

This release represents over 220 code commits from 22 different contributors.

  * New [command based interface](/overview/usage/) similar to git (`hugo server -s ./`)
  * Amber template support
  * [Aliases](/extras/aliases/) (redirects)
  * Support for top level pages (in addition to homepage)
  * Complete overhaul of the documentation site
  * Full Windows support
  * Better index support including [ordering by content weight](/content/ordering/)
  * Add params to site config, available in .Site.Params from templates
  * Friendlier json support
  * Support for html & xml content (with frontmatter support)
  * Support for [summary](/content/summaries/) content divider (<code>&lt;!&#45;&#45;more&#45;&#45;&gt;</code>)
  * HTML in [summary](/content/summaries/) (when using divider)
  * Added ["Minutes to Read"](/layout/variables/) functionality
  * Support for a custom 404 page
  * Cleanup of how content organization is handled
  * Loads of unit and performance tests
  * Integration with travis ci
  * Static directory now watched and copied on any addition or modification
  * Support for relative permalinks
  * Fixed watching being triggered multiple times for the same event
  * Watch now ignores temp files (as created by Vim)
  * Configurable number of posts on [homepage](/layout/homepage/)
  * [Front matter](/content/front-matter/) supports multiple types (int, string, date, float)
  * Indexes can now use a default template
  * Addition of truncated bool to content to determine if should show 'more' link
  * Support for [linkTitles](/layout/variables/)
  * Better handling of most errors with directions on how to resolve
  * Support for more date / time formats
  * Support for go 1.2
  * Support for `first` in templates

## **0.8.0** August 2, 2013

This release represents over 65 code commits from 6 different contributors.

  * Added support for pretty urls (filename/index.html vs filename.html)
  * Hugo supports a destination directory
  * Will efficiently sync content in static to destination directory
  * Cleaned up options.. now with support for short and long options
  * Added support for TOML
  * Added support for YAML
  * Added support for Previous & Next
  * Added support for indexes for the indexes
  * Better Windows compatibility
  * Support for series
  * Adding verbose output
  * Loads of bugfixes

## **0.7.0** July 4, 2013
  * Hugo now includes a simple server
  * First public release

## **0.6.0** July 2, 2013
  * Hugo includes an example documentation site which it builds

## **0.5.0** June 25, 2013
  * Hugo is quite usable and able to build spf13.com

