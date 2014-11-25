---
aliases:
- /doc/permalinks/
date: 2013-11-18
menu:
  main:
    parent: extras
next: /extras/shortcodes
notoc: true
prev: /extras/menus
title: Permalinks
weight: 70
---

By default, content is laid out into the target `publishdir` (public)
namespace matching its layout within the `contentdir` hierarchy.
The `permalinks` site configuration option allows you to adjust this on a
per-section basis.
This will change where the files are written to and will change the page's
internal "canonical" location, such that template references to
`.RelPermalink` will honour the adjustments made as a result of the mappings
in this option.

For instance, if one of your sections is called `post`, and you want to adjust
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

The following is a list of values that can be used in a permalink definition.
All references to time are dependent on the content's date.

  * **:year** the 4-digit year
  * **:month** the 2-digit month
  * **:monthname** the name of the month
  * **:day** the 2-digit day
  * **:weekday** the 1-digit day of the week (Sunday = 0)
  * **:weekdayname** the name of the day of the week
  * **:yearday** the 1- to 3-digit day of the year
  * **:section** the content's section
  * **:title** the content's title
  * **:slug** the content's slug (or title if no slug)
  * **:filename** the content's filename (without extension)

