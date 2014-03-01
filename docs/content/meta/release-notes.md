---
title: "Release Notes"
date: "2013-07-01"
aliases: ["/doc/release-notes/"]
groups: ["meta"]
groups_weight: 10
---
## **0.10.0** March 1, 2014
  * [Syntax highlighting](/extras/highlighting) powered by pygments (**slow**)
  * Ability to [sort content](/content/ordering) many more ways
  * Automatic [table of contents](/extras/toc) generation
  * Support for unicode urls, aliases and indexes
  * Configurable per-section [permalink](/extras/permalinks) pattern support
  * Support for [paired shortcodes](/extras/shortcodes)
  * Shipping with some [shortcodes](/extras/shortcodes) (highlight & figure)
  * Adding [canonify](/extras/urls) option to keep urls relative
  * A bunch of [additional template functions](/layouts/templatefunctions)
  * Watching very large sites now works on mac
  * RSS generation improved. Limited to 50 items by default, can limit further in [template](/layout/rss)
  * Boolean params now supported in [frontmatter](/content/front-matter)
  * Launched website [showcase](/showcase). Show off your own hugo site!
  * A bunch of [bug fixes](https://github.com/spf13/hugo/commits/master)

## **0.9.0** November 15, 2013
  * New [command based interface](/overview/usage) similar to git (hugo server -s ./ )
  * Amber template support
  * [Aliases](/extras/aliases) (redirects)
  * Support for top level pages (in addition to homepage)
  * Complete overhaul of the documentation site
  * Full Windows support
  * Better index support including [ordering by content weight](/content/ordering)
  * Add params to site config, available in .Site.Params from templates
  * Friendlier json support
  * Support for html & xml content (with frontmatter support)
  * Support for [summary](/content/summaries) content divider (&lt;!–more–>)
  * HTML in [summary](/content/summaries) (when using divider)
  * Added ["Minutes to Read"](/layout/variables) functionality
  * Support for a custom 404 page
  * Cleanup of how content organization is handled
  * Loads of unit and performance tests
  * Integration with travis ci
  * Static directory now watched and copied on any addition or modification
  * Support for relative permalinks
  * Fixed watching being triggered multiple times for the same event
  * Watch now ignores temp files (as created by Vim)
  * Configurable number of posts on [homepage](/layout/homepage/)
  * [Front matter](/content/front-matter) supports multiple types (int, string, date, float)
  * Indexes can now use a default template
  * Addition of truncated bool to content to determine if should show 'more' link
  * Support for [linkTitles](/layout/variables)
  * Better handling of most errors with directions on how to resolve
  * Support for more date / time formats
  * Support for go 1.2
  * Support for `first` in templates

## **0.8.0** August 2, 2013
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

