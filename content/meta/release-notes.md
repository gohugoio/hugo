---
title: "Release Notes"
date: "2013-07-01"
aliases: ["/doc/release-notes/"]
groups: ["meta"]
groups_weight: 10
---

* **0.9.0** HEAD
  * New command based interface similar to git (hugo server -s ./ )
  * Added support for aliases (redirects)
  * Cleanup of how content organization is handled
  * Support for top level pages (other than homepage)
  * Loads of unit and performance tests
  * Integration with travis ci
  * Complete overhaul of the documentation site
  * Full Windows support
  * Support for ordering pages by weight
  * More support for indexes
  * Support for relative permalinks
  * Fixed watching being triggered multiple times for the same event
  * Watch now ignores temp files (as created by Vim)
  * Front matter supports multiple types (int, string, date, float)
  * Support for summary content divider (&lt;!--more-->)
  * Ability to highlight current page.
  * Added default templates for indexes
  * Support for Amber templates
  * Support for more date / time formats
  * Support for go 1.2
* **0.8.0** August 2, 2013
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
* **0.7.0** July 4, 2013
  * Hugo now includes a simple server
  * First public release
* **0.6.0** July 2, 2013
  * Hugo includes an example documentation site which it builds
* **0.5.0** June 25, 2013
  * Hugo is quite usable and able to build spf13.com

