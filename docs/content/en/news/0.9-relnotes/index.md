---
date: 2013-11-16T04:52:32Z
description: "Hugo 0.9 is the most significant update to Hugo ever!  It contains contributions from dozens of contributors and represents hundreds of features, fixes and improvements."
title: "Hugo 0.9"
categories: ["Releases"]
---

This is the most significant update to Hugo ever!
It contains contributions from dozens of contributors and represents hundreds of features, fixes and improvements.

# Major New Features
- New command based interface similar to git (`hugo server -s ./`)
- Amber template support
- Full Windows support
- Better index support including ordering by content weight
- Add params to site config, available in `.Site.Params` from templates
- Support for html & xml content (with front matter support)
- Support for top level pages (in addition to homepage)

# Notable Fixes and Additions
- Friendlier json support
- Aliases (redirects)
- Support for summary content divider (`<!--more-->`)
- HTML & shortcodes supported in summary (when using divider)
- Complete overhaul of the documentation site
- Added "Minutes to Read" functionality
- Support for a custom 404 page
- Cleanup of how content organization is handled
- Loads of unit and performance tests
- Integration with Travis CI
- Static directory now watched and copied on any addition or modification
- Support for relative permalinks
- Fixed watching being triggered multiple times for the same event
- Watch now ignores temp files (as created by Vim)
- Configurable number of posts on homepage
- Front matter supports multiple types (int, string, date, float)
- Indexes can now use a default template
- Addition of truncated bool to content to determine if should show 'more' link
- Support for `linkTitles`
- Better handling of most errors with directions on how to resolve
- Support for more date / time formats
- Support for Go 1.2
- Loads more... see commit log for full list.
