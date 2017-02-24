---
title: Local File Templates
linktitle: Local File Templates
description: Hugo's `readerDir` and `readFile` functions make it easy to traverse your project's directory structure and write file contents to your templates.
godocref: https://golang.org/pkg/os/#FileInfo
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [files,directories]
weight: 110
draft: false
aliases: [/extras/localfiles/,/templates/files/]
toc: true
needsreview: true
---

## Traversing Local Files

With Hugo's [`readDir` function][], you can traverse your website's files on your server.

## Using _readDir_

The `readDir` function returns an array of [`os.FileInfo`](https://golang.org/pkg/os/#FileInfo). It takes the file's `path` as a single string argument. This path can be to any directory of your website (i.e. as found on your server's file system).

Whether the path is absolute or relative does not matter because---at least for `readDir`---the root of your website (typically `./public/`) in effect becomes both:

1. The file system root
2. The current working directory

## Example Shortcode: List Directory's Files

So, let's create a new shortcode using `readDir`:

{{% input file="layouts/shortcodes/directoryindex.html" %}}<pre><code>{{< readfile "layouts/shortcodes/directoryindex.html" >}}</code></pre>{{% /input %}}

For the files in any given directory, this shortcode usefully lists the files' basenames and sizes and also creates a link to each of them.

This shortcode [has already been included in this very website][]. So, let's list some of its CSS files. (If you click on their names, you can reveal the contents.)

{{< directoryindex path="/static/css" pathURL="/css" >}}

The following is the [shortcode declaration][shortcodes] used to render the above output:

```html
{{</* directoryindex path="/static/css" pathURL="/css" */>}}
```

{{% note "Slashes are Important" %}}
The initial slash `/` in `pathURL` is important. Otherwise, `pathURL` becomes relative to the current web page.
{{% /note %}}

[has already been included in this very website]: https://github.com/spf13/hugo/blob/master/docs/layouts/shortcodes/directoryindex.html
[`readDir` function]: /functions/readdir/
[shortcodes]: /content-management/shortcodes/