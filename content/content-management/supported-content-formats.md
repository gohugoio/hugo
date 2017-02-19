---
title: Supported Content Formats
linktitle:
description: Hugo uses the Blackfriday markdown parser for content files but also provides support for additional syntaxes (eg, Asciidoc) via external helpers.
date: 2017-01-10
publishdate: 2017-01-10
lastmod: 2017-01-10
categories: [content management]
tags: [markdown,asciidoc,mmark,content format]
weight: 10
draft: false
slug:
aliases: [/content/markdown-extras/,/content/supported-formats/,/content/markdown/]
toc: true
notesforauthors:
---

## Markdown

Markdown is the natively supported content format for Hugo and is rendered using the excellent [Blackfriday project][], a markdown parser written in Golang.

{{% note "Deeply Nested Lists" %}}
Blackfriday has a known issue [(#329)](https://github.com/russross/blackfriday/issues/329) with handling deeply nested lists, but there is a workaround. If you write lists in markdown, use 4 spaces (i.e., <kbd>tab</kbd>) rather than 2 to delimit nesting of lists.
{{% /note %}}

## Additional Content Formats

Since 0.14, Hugo has defined a new concept called _external helpers_. This means you can write your content using [Asciidoc][] or [reStructuredText][]. If you have files with associated extensions, Hugo will call external commands to generate the content ([see Hugo source code][]).

For example, for Asciidoc files, Hugo will try to call the **asciidoctor** or **asciidoc** command. This means that you will have to install the associated tool on your machine to be able to use these formats.

To use these formats, just use the standard extension and the front matter exactly as you would do with natively supported `.md` files.

{{% note "Performance of External Helpers" %}}
Because these are external commands, generation performance for your preferred content format will heavily depend on the performance of the external tool used. As this feature is still in its infancy, feedback is especially welcome.
{{% /note %}}

## Extending Markdown

Hugo provides some convenient methods for extending markdown.

### Task Lists

Hugo supports GitHub styled task lists (TODO lists) for the Blackfriday markdown renderer. If you do not want to use this feature, you can disable it in the See [Blackfriday config](/overview/configuration/#configure-blackfriday-rendering) for how to turn it off.

#### Task List

```markdown
- [ ] a task list item
- [ ] list syntax required
- [ ] incomplete
- [x] completed
```

Renders as:

- [ ] a task list item
- [ ] list syntax required
- [ ] incomplete
- [x] completed

And produces this HTML:

```html
<ul class="task-list">
    <li><input type="checkbox" disabled="" class="task-list-item"> a task list item</li>
    <li><input type="checkbox" disabled="" class="task-list-item"> list syntax required</li>
    <li><input type="checkbox" disabled="" class="task-list-item"> incomplete</li>
    <li><input type="checkbox" checked="" disabled="" class="task-list-item"> completed</li>
</ul>
```

### Shortcodes

If you write in markdown and find yourself frequently embedding your content with raw HTML, Hugo provides built-in [shortcodes][] functionality to act as the intermediary between your content and templating.

### Code Blocks

Hugo supports GitHub-flavored markdown's use of triple back ticks, as well as provides a special [`highlight` nested shortcode][] to render syntax highlighting via [Pygments][]. For usage examples and a complete explanation, see the [syntax highlighting documentation][] in [developer tools][].

## Markdown Learning Resources

* [Markdown Tutorial][]
* [Daring Fireball: Markdown, John Gruber][]

[`highlight` nested shortcode]: /content-management/shortcodes/#highlight
[AsciiDoc]: http://asciidoc.org/
[Blackfriday project]: https://github.com/russross/blackfriday
[Daring Fireball: Markdown, John Gruber]: https://daringfireball.net/projects/markdown/
[developer tools]: /developer-tools/
[Markdown Tutorial]: http://www.markdowntutorial.com/
[Pygments]: http://pygments.org/
[reStructuredText]: http://docutils.sourceforge.net/rst.html
[see Hugo source code]: https://github.com/spf13/hugo/blob/77c60a3440806067109347d04eb5368b65ea0fe8/helpers/general.go#L65
[shortcodes]: /content-management/shortcodes/
[syntax highlighting documentation]: /developer-tools/syntax-highlighting/