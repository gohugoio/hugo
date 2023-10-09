---
title: Local file templates
description: Hugo's `readDir` and `readFile` functions make it easy to traverse your project's directory structure and write file contents to your templates.
categories: [templates]
keywords: [files,directories]
menu:
  docs:
    parent: templates
    weight: 180
weight: 180
aliases: [/extras/localfiles/,/templates/local-files/]
toc: true
---

## Traverse local files

With Hugo's [`readDir`][readDir] and [`readFile`][readFile] template functions, you can traverse your website's files on your server.

## Use `readDir`

The [`readDir` function][readDir] returns an array of [`os.FileInfo`][osfileinfo]. It takes the file's `path` as a single string argument. This path can be to any directory of your website (i.e., as found on your server's file system).

Whether the path is absolute or relative does not matter because---at least for `readDir`---the root of your website (typically `./public/`) in effect becomes both:

1. The file system root
2. The current working directory

## Use `readFile`

The [`readfile` function][readFile] reads a file from disk and converts it into a string to be manipulated by other Hugo functions or added as-is. `readFile` takes the file, including path, as an argument passed to the function.

To use the `readFile` function in your templates, make sure the path is relative to your *Hugo project's root directory*:

```go-html-template
{{ readFile "/content/templates/local-file-templates" }}
```

### `readFile` Example: Add a Project File to Content

As `readFile` is a function, it is only available to you in your templates and not your content. However, we can create a simple [shortcode template][sct] that calls `readFile`, passes the first argument through the function, and then allows an optional second argument to send the file through the Markdown processor. The pattern for adding this shortcode to your content will be as follows:

```go-html-template
{{</* readfile file="/path/to/local/file.txt" markdown="true" */>}}
```

{{% note %}}
If you are going to create [custom shortcodes](/templates/shortcode-templates/) with `readFile` for a theme, note that usage of the shortcode will refer to the project root and *not* your `themes` directory.
{{% /note %}}

[called directly in the Hugo docs]: https://github.com/gohugoio/hugoDocs/blob/master/content/en/templates/files.md
[osfileinfo]: https://golang.org/pkg/os/#FileInfo
[readDir]: /functions/readdir/
[readFile]: /functions/readfile/
[sc]: /content-management/shortcodes/
[sct]: /templates/shortcode-templates/
[readfilesource]: https://github.com/gohugoio/hugoDocs/blob/master/layouts/shortcodes/readfile.html
[testfile]: https://github.com/gohugoio/hugoDocs/blob/master/content/en/readfiles/testing.txt
