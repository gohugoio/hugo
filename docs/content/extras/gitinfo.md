---
aliases:
- /doc/gitinfo/
lastmod: 2016-12-11
date: 2016-12-11
menu:
  main:
    parent: extras
next: /extras/livereload
prev: /extras/datadrivencontent
title: GitInfo
---

Hugo provides a way to integrate Git data into your site.


## Prerequisites

1. The Hugo site must be in a Git-enabled directory.
1. The Git executable must be installed and in your system `PATH`.
1. Enable the GitInfo feature in Hugo by using `--enableGitInfo` on the command
   line or by setting `enableGitInfo` to `true` in your site configuration.

## The GitInfo Object

The `GitInfo` object contains the following fields:

AbbreviatedHash
: abbreviated commit hash, e.g. `866cbcc`

AuthorName
: author name, respecting `.mailmap`

AuthorEmail
: author email address, respecting `.mailmap`

AuthorDate
: the author date

Hash
: commit hash, e.g. `866cbccdab588b9908887ffd3b4f2667e94090c3`

Subject
: commit message subject, e.g. `tpl: Add custom index function`


## Performance Considerations

The Git integrations should be fairly performant, but it does add some time to the build, which depends somewhat on the Git history size.

