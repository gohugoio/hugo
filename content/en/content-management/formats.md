---
title: Supported Content Formats
linktitle: Supported Content Formats
description: Both HTML and Markdown are supported content formats.
date: 2017-01-10
publishdate: 2017-01-10
lastmod: 2017-04-06
categories: [content management]
keywords: [markdown,asciidoc,mmark,pandoc,content format]
menu:
  docs:
    parent: "content-management"
    weight: 20
weight: 20	#rem
draft: false
aliases: [/content/markdown-extras/,/content/supported-formats/,/doc/supported-formats/,/tutorials/mathjax/]
toc: true
---

**Markdown is the main content format** and comes in two flavours:  The excellent [Blackfriday project][blackfriday] (name your files `*.md` or set `markup = "markdown"` in front matter) or its fork [Mmark][mmark] (name your files `*.mmark` or set `markup = "mmark"` in front matter), both very fast markdown engines written in Go.

For Emacs users, [goorgeous](https://github.com/chaseadamsio/goorgeous) provides built-in native support for Org-mode  (name your files `*.org` or set `markup = "org"` in front matter)

But in many situations, plain HTML is what you want. Just name your files with `.html` or `.htm` extension inside your content folder. Note that if you want your HTML files to have a layout, they need front matter. It can be empty, but it has to be there:

```html
---
title: "This is a content file in HTML"
---

<div>
  <h1>Hello, Hugo!</h1>
</div>
```

{{% note "Deeply Nested Lists" %}}
Before you begin writing your content in markdown, Blackfriday has a known issue [(#329)](https://github.com/russross/blackfriday/issues/329) with handling deeply nested lists. Luckily, there is an easy workaround. Use 4-spaces (i.e., <kbd>tab</kbd>) rather than 2-space indentations.
{{% /note %}}

## Configure BlackFriday Markdown Rendering

You can configure multiple aspects of Blackfriday as show in the following list. See the docs on [Configuration][config] for the full list of explicit directions you can give to Hugo when rendering your site.

{{< readfile file="/content/en/readfiles/bfconfig.md" markdown="true" >}}

## Extend Markdown

Hugo provides some convenient methods for extending markdown.

### Task Lists

Hugo supports [GitHub-styled task lists (i.e., TODO lists)][gfmtasks] for the Blackfriday markdown renderer. If you do not want to use this feature, you can disable it in your configuration.

#### Example Task List Input

{{< code file="content/my-to-do-list.md" >}}
- [ ] a task list item
- [ ] list syntax required
- [ ] incomplete
- [x] completed
{{< /code >}}

#### Example Task List Output

The preceding markdown produces the following HTML in your rendered website:

```
<ul class="task-list">
    <li><input type="checkbox" disabled="" class="task-list-item"> a task list item</li>
    <li><input type="checkbox" disabled="" class="task-list-item"> list syntax required</li>
    <li><input type="checkbox" disabled="" class="task-list-item"> incomplete</li>
    <li><input type="checkbox" checked="" disabled="" class="task-list-item"> completed</li>
</ul>
```

#### Example Task List Display

The following shows how the example task list will look to the end users of your website. Note that visual styling of lists is up to you. This list has been styled according to [the Hugo Docs stylesheet][hugocss].

- [ ] a task list item
- [ ] list syntax required
- [ ] incomplete
- [x] completed

### Emojis

To add emojis directly to content, set `enableEmoji` to `true` in your [site configuration][config]. To use emojis in templates or shortcodes, see [`emojify` function][].

For a full list of emojis, see the [Emoji cheat sheet][emojis].

### Shortcodes

If you write in Markdown and find yourself frequently embedding your content with raw HTML, Hugo provides built-in shortcodes functionality. This is one of the most powerful features in Hugo and allows you to create your own Markdown extensions very quickly.

See [Shortcodes][sc] for usage, particularly for the built-in shortcodes that ship with Hugo, and [Shortcode Templating][sct] to learn how to build your own.

### Code Blocks

Hugo supports GitHub-flavored markdown's use of triple back ticks, as well as provides a special [`highlight` shortcode][hlsc], and syntax highlights those code blocks natively using *Chroma*. Users also have an option to use *Pygments* instead. See the [Syntax Highlighting][hl] section for details.

## Mmark

Mmark is a [fork of BlackFriday][mmark] and markdown superset that is well suited for writing [IETF documentation][ietf]. You can see examples of the syntax in the [Mmark GitHub repository][mmarkgh] or the full syntax on [Miek Gieben's website][].

### Use Mmark

As Hugo ships with Mmark, using the syntax is as easy as changing the extension of your content files from `.md` to `.mmark`.

In the event that you want to only use Mmark in specific files, you can also define the Mmark syntax in your content's front matter:

```
---
title: My Post
date: 2017-04-01
markup: mmark
---
```

{{% warning %}}
Thare are some features not available in Mmark; one example being that shortcodes are not translated when used in an included `.mmark` file ([#3131](https://github.com/gohugoio/hugo/issues/3137)), and `EXTENSION_ABBREVIATION` ([#1970](https://github.com/gohugoio/hugo/issues/1970)) and the aforementioned GFM todo lists ([#2270](https://github.com/gohugoio/hugo/issues/2270)) are not fully supported. Contributions are welcome.
{{% /warning %}}

## MathJax with Hugo

[MathJax](http://www.mathjax.org/) is a JavaScript library that allows the display of mathematical expressions described via a LaTeX-style syntax in the HTML (or Markdown) source of a web page. As it is a pure a JavaScript library, getting it to work within Hugo is fairly straightforward, but does have some oddities that will be discussed here.

This is not an introduction into actually using MathJax to render typeset mathematics on your website. Instead, this page is a collection of tips and hints for one way to get MathJax working on a website built with Hugo.

### Enable MathJax

The first step is to enable MathJax on pages that you would like to have typeset math. There are multiple ways to do this (adventurous readers can consult the [Loading and Configuring](http://docs.mathjax.org/en/latest/configuration.html) section of the MathJax documentation for additional methods of including MathJax), but the easiest way is to use the secure MathJax CDN by include a `<script>` tag for the officially recommended secure CDN ([cdn.js.com](https://cdnjs.com)):

{{< code file="add-mathjax-to-page.html" >}}
<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.1/MathJax.js?config=TeX-AMS-MML_HTMLorMML">
</script>
{{< /code >}}

One way to ensure that this code is included in all pages is to put it in one of the templates that live in the `layouts/partials/` directory. For example, I have included this in the bottom of my template `footer.html` because I know that the footer will be included in every page of my website.

### Options and Features

MathJax is a stable open-source library with many features. I encourage the interested reader to view the [MathJax Documentation](http://docs.mathjax.org/en/latest/index.html), specifically the sections on [Basic Usage](http://docs.mathjax.org/en/latest/index.html#basic-usage) and [MathJax Configuration Options](http://docs.mathjax.org/en/latest/index.html#mathjax-configuration-options).

### Issues with Markdown

{{% note %}}
The following issues with Markdown assume you are using `.md` for content and BlackFriday for parsing. Using [Mmark](#mmark) as your content format will obviate the need for the following workarounds.

When using Mmark with MathJax, use `displayMath: [['$$','$$'], ['\\[','\\]']]`. See the [Mmark `README.md`](https://github.com/miekg/mmark/wiki/Syntax#math-blocks) for more information. In addition to MathJax, Mmark has been shown to work well with [KaTeX](https://github.com/Khan/KaTeX). See this [related blog post from a Hugo user](http://nosubstance.me/post/a-great-toolset-for-static-blogging/).
{{% /note %}}

After enabling MathJax, any math entered between proper markers (see the [MathJax documentation][mathjaxdocs]) will be processed and typeset in the web page. One issue that comes up, however, with Markdown is that the underscore character (`_`) is interpreted by Markdown as a way to wrap text in `emph` blocks while LaTeX (MathJax) interprets the underscore as a way to create a subscript. This "double speak" of the underscore can result in some unexpected and unwanted behavior.

### Solution

There are multiple ways to remedy this problem. One solution is to simply escape each underscore in your math code by entering `\_` instead of `_`. This can become quite tedious if the equations you are entering are full of subscripts.

Another option is to tell Markdown to treat the MathJax code as verbatim code and not process it. One way to do this is to wrap the math expression inside a `<div>` `</div>` block. Markdown would ignore these sections and they would get passed directly on to MathJax and processed correctly. This works great for display style mathematics, but for inline math expressions the line break induced by the `<div>` is not acceptable. The syntax for instructing Markdown to treat inline text as verbatim is by wrapping it in backticks (`` ` ``). You might have noticed, however, that the text included in between backticks is rendered differently than standard text (on this site these are items highlighted in red). To get around this problem, we could create a new CSS entry that would apply standard styling to all inline verbatim text that includes MathJax code. Below I will show the HTML and CSS source that would accomplish this (note this solution was adapted from [this blog post](http://doswa.com/2011/07/20/mathjax-in-markdown.html)---all credit goes to the original author).

{{< code file="mathjax-markdown-solution.html" >}}
<script type="text/x-mathjax-config">
MathJax.Hub.Config({
  tex2jax: {
    inlineMath: [['$','$'], ['\\(','\\)']],
    displayMath: [['$$','$$'], ['\[','\]']],
    processEscapes: true,
    processEnvironments: true,
    skipTags: ['script', 'noscript', 'style', 'textarea', 'pre'],
    TeX: { equationNumbers: { autoNumber: "AMS" },
         extensions: ["AMSmath.js", "AMSsymbols.js"] }
  }
});
</script>

<script type="text/x-mathjax-config">
  MathJax.Hub.Queue(function() {
    // Fix <code> tags after MathJax finishes running. This is a
    // hack to overcome a shortcoming of Markdown. Discussion at
    // https://github.com/mojombo/jekyll/issues/199
    var all = MathJax.Hub.getAllJax(), i;
    for(i = 0; i < all.length; i += 1) {
        all[i].SourceElement().parentNode.className += ' has-jax';
    }
});
</script>
{{< /code >}}



As before, this content should be included in the HTML source of each page that will be using MathJax. The next code snippet contains the CSS that is used to have verbatim MathJax blocks render with the same font style as the body of the page.

{{< code file="mathjax-style.css" >}}
code.has-jax {
    font: inherit;
    font-size: 100%;
    background: inherit;
    border: inherit;
    color: #515151;
}
{{< /code >}}

In the CSS snippet, notice the line `color: #515151;`. `#515151` is the value assigned to the `color` attribute of the `body` class in my CSS. In order for the equations to fit in with the body of a web page, this value should be the same as the color of the body.

### Usage

With this setup, everything is in place for a natural usage of MathJax on pages generated using Hugo. In order to include inline mathematics, just put LaTeX code in between `` `$ TeX Code $` `` or `` `\( TeX Code \)` ``. To include display style mathematics, just put LaTeX code in between `<div>$$TeX Code$$</div>`. All the math will be properly typeset and displayed within your Hugo generated web page!

## Additional Formats Through External Helpers

Hugo has a new concept called _external helpers_. It means that you can write your content using [Asciidoc][ascii], [reStructuredText][rest], or [pandoc]. If you have files with associated extensions, Hugo will call external commands to generate the content. ([See the Hugo source code for external helpers][helperssource].)

For example, for Asciidoc files, Hugo will try to call the `asciidoctor` or `asciidoc` command. This means that you will have to install the associated tool on your machine to be able to use these formats. ([See the Asciidoctor docs for installation instructions](http://asciidoctor.org/docs/install-toolchain/)).

To use these formats, just use the standard extension and the front matter exactly as you would do with natively supported `.md` files.

Hugo passes reasonable default arguments to these external helpers by default:

- `asciidoc`: `--no-header-footer --safe -`
- `asciidoctor`: `--no-header-footer --safe --trace -`
- `rst2html`: `--leave-comments --initial-header-level=2`
- `pandoc`: `--mathjax`

{{% warning "Performance of External Helpers" %}}
Because additional formats are external commands generation performance will rely heavily on the performance of the external tool you are using. As this feature is still in its infancy, feedback is welcome.
{{% /warning %}}

## Learn Markdown

Markdown syntax is simple enough to learn in a single sitting. The following are excellent resources to get you up and running:

* [Daring Fireball: Markdown, John Gruber (Creator of Markdown)][fireball]
* [Markdown Cheatsheet, Adam Pritchard][mdcheatsheet]
* [Markdown Tutorial (Interactive), Garen Torikian][mdtutorial]
* [The Markdown Guide, Matt Cone][mdguide]

[`emojify` function]: /functions/emojify/
[ascii]: http://asciidoctor.org/
[bfconfig]: /getting-started/configuration/#configuring-blackfriday-rendering
[blackfriday]: https://github.com/russross/blackfriday
[mmark]: https://github.com/miekg/mmark
[config]: /getting-started/configuration/
[developer tools]: /tools/
[emojis]: https://www.webpagefx.com/tools/emoji-cheat-sheet/
[fireball]: https://daringfireball.net/projects/markdown/
[gfmtasks]: https://guides.github.com/features/mastering-markdown/#syntax
[helperssource]: https://github.com/gohugoio/hugo/blob/77c60a3440806067109347d04eb5368b65ea0fe8/helpers/general.go#L65
[hl]: /content-management/syntax-highlighting/
[hlsc]: /content-management/shortcodes/#highlight
[hugocss]: /css/style.css
[ietf]: https://tools.ietf.org/html/
[mathjaxdocs]: https://docs.mathjax.org/en/latest/
[mdcheatsheet]: https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet
[mdguide]: https://www.markdownguide.org/
[mdtutorial]: http://www.markdowntutorial.com/
[Miek Gieben's website]: https://miek.nl/2016/march/05/mmark-syntax-document/
[mmark]: https://github.com/miekg/mmark
[mmarkgh]: https://github.com/miekg/mmark/wiki/Syntax
[org]: http://orgmode.org/
[pandoc]: http://www.pandoc.org/
[Pygments]: http://pygments.org/
[rest]: http://docutils.sourceforge.net/rst.html
[sc]: /content-management/shortcodes/
[sct]: /templates/shortcode-templates/
