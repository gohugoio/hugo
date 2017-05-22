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


## Other Considerations

The Git integrations should be fairly performant, but it does add some time to the build, which depends somewhat on the Git history size.

The accuracy of data depends on the underlying local git respository.  If the local repository is a *shallow clone*, then any file that hasn't been modified in the truncated history will default to data in the oldest commit.  In particular, if the respository has been cloned using `--depth=1` then every file will the exact same `GitInfo` data -- that of the only commit in the repository.

In particular, many CI/CD systems such as [travis-ci.org](https://travis-ci.org) default to a clone depth of 50 which is unlikely to be deep enough.  You can explicitly add back the missing history using using `git fetch --unshallow` or [make the initial checkout deeper](https://docs.travis-ci.com/user/customizing-the-build#Git-Clone-Depth).

