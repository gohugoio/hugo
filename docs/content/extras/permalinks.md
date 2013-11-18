---
title: "Permalinks"
date: "2013-11-18"
aliases:
  - "/doc/permalinks/"
groups: ["extras"]
groups_weight: 30
---

By default, content is laid out into the target `publishdir` (public)
namespace matching its layout within the `contentdir` hierarchy.
The `permalinks` site configuration option allows you to adjust this on a
per-section basis.
This will change where the files are written to and will change the page's
internal "canonical" location, such that template references to
`.RelPermalink` will honour the adjustments made as a result of the mappings
in this option.

For instance, if one of your sections is called `post` and you want to adjust
the canonical path to be hierarchical based on the year and month, then you
might use:

```yaml
permalinks:
  post: /:year/:month/:title/
```

Only the content under `post/` will be so rewritten.
A file named `content/post/sample-entry` which contains a line
`date: 2013-11-18T19:20:00-05:00` might end up with the rendered page
appearing at `public/2013/11/sample-entry/index.html` and be reachable via
the URL <http://yoursite.example.com/2013/11/sample-entry/>.

