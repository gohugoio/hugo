---
aliases:
- /doc/localfiles/
lastmod: 2016-09-12
date: 2015-06-12
menu:
  main:
    parent: extras
next: /extras/urls
notoc: true
prev: /extras/toc
title: Traversing Local Files
---
## Traversing Local Files

Using Hugo's function `readDir`,
you can traverse your web site's files on your server.
## Using _readDir_

The `readDir` function returns an array
of [`os.FileInfo`](https://golang.org/pkg/os/#FileInfo).
It takes a single, string argument: a path.
This path can be to any directory of your web site
(as found on your server's filesystem).

Whether the path is absolute or relative makes no difference,
because&mdash;at least for `readDir`&mdash;the root of your web site (typically `./public/`)
in effect becomes both:

1. The filesystem root; and
1. The current working directory.

## New Shortcode

So, let's create a new shortcode using `readDir`:

**layouts/shortcodes/directoryindex.html**
```html
{{< readfile "layouts/shortcodes/directoryindex.html" >}}
```
For the files in any given directory,
this shortcode usefully lists their basenames and sizes,
while providing links to them.

Already&mdash;actually&mdash;this shortcode
has been included in this very web site.
So, let's list some of its CSS files.
(If you click on their names, you can reveal the contents.)
{{<   directoryindex path="/static/css" pathURL="/css"   >}}
<br />
This is the call that rendered the above output:
```html
{{</* directoryindex path="/static/css" pathURL="/css" */>}}
```
By the way,
regarding the pathURL argument, the initial slash `/` is important.
Otherwise, it becomes relative to the current web page.
