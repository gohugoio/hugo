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

<dl>
<dt><code>:year</code></dt><dd>the 4-digit year</dd>
<dt><code>:month</code></dt><dd>the 2-digit zero-padded month (01, 02, …, 12)</dd>
<dt><code>:monthname</code></dt><dd>the name of the month in English (“January”, “February”, …)</dd>
<dt><code>:day</code></dt><dd>the 2-digit zero-padded day (01, 02, …, 31)</dd>
<dt><code>:weekday</code></dt><dd>the 1-digit day of the week (Sunday = 0)</dd>
<dt><code>:weekdayname</code></dt><dd>the name of the day of the week in English (“Sunday”, “Monday”, …)</dd>
<dt><code>:yearday</code></dt><dd>the 1- to 3-digit day of the year, in the range [1,365] or [1,366]</dd>
<dt><code>:section</code></dt><dd>the content’s section</dd>
<dt><code>:title</code></dt><dd>the content’s title</dd>
<dt><code>:slug</code></dt><dd>the content’s slug (or title if no slug)</dd>
<dt><code>:filename</code></dt><dd>the content’s filename (without extension)</dd>
</dl>
