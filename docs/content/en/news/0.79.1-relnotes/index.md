
---
date: 2020-12-19
title: "Hugo 0.79.1: One Security Patch for Hugo on Windows"
description: "Disallow running of e.g. Pandoc in the current directory."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

Hugo depends on Go's `os/exec` for certain features, e.g. for rendering of Pandoc documents if these binaries are found in the system `%PATH%` on Windows. However, if a malicious file with the same name (`exe` or `bat`) was found in the current working directory at the time of running `hugo`, the malicious command would be invoked instead of the system one.

Windows users who ran `hugo` inside untrusted Hugo sites were affected.

The origin of this issue comes from Go, see https://github.com/golang/go/issues/38736

We have fixed this in Hugo by [using](https://github.com/gohugoio/hugo/commit/4a8267d64a40564aced0695bca05249da17b0eab) a patched version of `exec.LookPath` from https://github.com/cli/safeexec (thanks to [@mislav](https://github.com/mislav) for the implementation).

Thanks to [@Ry0taK](https://github.com/Ry0taK) for the bug report.


