---
title: "URLs"
date: "2014-01-03"
aliases:
  - "/doc/urls/"
groups: ["extras"]
groups_weight: 40
---
By default, all relative URLs encountered in the input will be canonicalized
using `baseurl`, so that a link `/css/foo.css` becomes
`http://yoursite.example.com/css/foo.css`.

Setting `canonifyurls` to `false` will prevent this canonicalization.

Benefits of canonicalization include fixing all URLs to be absolute, which may
aid with some parsing tasks.  Note though that all real browsers handle this
client-side without issues.

Benefits of non-canonicalization include being able to have resource inclusion
be scheme-relative, so that http vs https can be decided based on how this
page was retrieved.
