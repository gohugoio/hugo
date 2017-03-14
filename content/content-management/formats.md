---
title: Supported Content Formats
linktitle: Formats
description: Markdown is natively supported in Hugo and is parsed by the feature-rich and incredibly speed Blackfriday parse. Hugo also provides support for additional syntaxes (eg, Asciidoc) via external helpers.
date: 2017-01-10
publishdate: 2017-01-10
lastmod: 2017-01-10
categories: [content management]
tags: [markdown,asciidoc,mmark,content format]
weight: 20
draft: false
aliases: [/content/markdown-extras/,/content/supported-formats/,/doc/supported-formats/]
toc: true
---

## Markdown

Markdown is the native content format for Hugo and is rendered using the excellent [Blackfriday project][blackfriday], a blazingly fast parser written in Golang.

{{% note "Deeply Nested Lists" %}}
Before you begin writing your content in markdown, Blackfriday has a known issue [(#329)](https://github.com/russross/blackfriday/issues/329) with handling deeply nested lists. Luckily, there is an easy workaround. Use 4-spaces (i.e., <kbd>tab</kbd>) rather than 2-space indentations.
{{% /note %}}

## Configuring Markdown Rendering

You can configure multiple aspects of Blackfriday as show in the following list. See the docs on [Configuration][config] for the full list of explicit directions you can give to Hugo when rendering your site.

{{< readfile file="content/readfiles/bfconfig.md" markdown="true" >}}

## Extending Markdown

Hugo provides some convenient methods for extending markdown.

### Task Lists

Hugo supports [GitHub-styled task lists (i.e., TODO lists)][gfmtasks] for the Blackfriday markdown renderer. If you do not want to use this feature, you can disable it in your configuration.

#### Example Task List Input

{{% code file="content/my-to-do-list.md" %}}
```markdown
- [ ] a task list item
- [ ] list syntax required
- [ ] incomplete
- [x] completed
```
{{% /code %}}

#### Example Task List Output

The preceding markdown produces the following HTML in your rendered website:

{{% output file="my-to-do-list.html" %}}
```html
<ul class="task-list">
    <li><input type="checkbox" disabled="" class="task-list-item"> a task list item</li>
    <li><input type="checkbox" disabled="" class="task-list-item"> list syntax required</li>
    <li><input type="checkbox" disabled="" class="task-list-item"> incomplete</li>
    <li><input type="checkbox" checked="" disabled="" class="task-list-item"> completed</li>
</ul>
```
{{% /output %}}

#### Example Task List Display

The following shows how the example task list will look to the end users of your website. Note that visual styling of lists is up to you. This list has been styled according to [the Hugo Docs stylesheet][hugocss].

- [ ] a task list item
- [ ] list syntax required
- [ ] incomplete
- [x] completed

### Shortcodes

If you write in markdown and find yourself frequently embedding your content with raw HTML, Hugo provides built-in shortcodes functionality to act as the intermediary between your content and templating. This is one of the most powerful features in Hugo and allows you to essentially create your own markdown extensions very quickly.

See [Shortcodes][sc] for usage, particularly for the built-in shortcodes that ship with Hugo, and [Shortcode Templating][sct] to learn how to build your own.

### Code Blocks

Hugo supports GitHub-flavored markdown's use of triple back ticks, as well as provides a special [`highlight` nested shortcode][hlsc] to render syntax highlighting via [Pygments][]. For usage examples and a complete explanation, see the [syntax highlighting documentation][hl] in [developer tools][].

## Additional Content Formats

Since 0.14, Hugo has defined a new concept called _external helpers_. It means that you can write your content using [Asciidoc][ascii], [reStructuredText][rest], or Emacs org-mode. If you have files with associated extensions, Hugo will call external commands to generate the content. ([See the Hugo source code for external helpers][helperssource]).

For example, for Asciidoc files, Hugo will try to call the `asciidoctor` or `asciidoc` command. This means that you will have to install the associated tool on your machine to be able to use these formats. ([See the Asciidoctor docs for installation instructions](http://asciidoctor.org/docs/install-toolchain/)).

To use these formats, just use the standard extension and the front matter exactly as you would do with natively supported `.md` files.

{{% warning "Performance of External Helpers" %}}
Because additional formats are external commands---with the exception of org mode---generation performance will rely heavily on the performance of the external tool you are using. As this feature is still in its infancy, feedback is welcome.
{{% /warning %}}

## Markdown Learning Resources

If you are unfamiliar with markdown syntax, it can easily be learned within a single sitting. The following are excellent resources to get you up and running:

* [Daring Fireball: Markdown, John Gruber (Creator of Markdown)][fireball]
* [Markdown Cheatsheet, Adam Pritchard][mdcheatsheet]
* [Markdown Tutorial (Interactive), Garen Torikian][mdtutorial]

[ascii]: http://asciidoc.org/
[bfconfig]: /getting-started/configuration/#configuring-blackfriday-rendering
[blackfriday]: https://github.com/russross/blackfriday
[config]: /getting-started/configuration/
[developer tools]: /tools/
[fireball]: https://daringfireball.net/projects/markdown/
[gfmtasks]: https://guides.github.com/features/mastering-markdown/#syntax
[helperssource]: https://github.com/spf13/hugo/blob/77c60a3440806067109347d04eb5368b65ea0fe8/helpers/general.go#L65
[hl]: /tools/syntax-highlighting/
[hlsc]: /content-management/shortcodes/#highlight
[hugocss]: /css/style.min.css
[mdcheatsheet]: https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet
[mdtutorial]: http://www.markdowntutorial.com/
[org]: http://orgmode.org/
[Pygments]: http://pygments.org/
[rest]: http://docutils.sourceforge.net/rst.html
[sc]: /content-management/shortcodes/
[sct]: /templates/shortcode-templates/