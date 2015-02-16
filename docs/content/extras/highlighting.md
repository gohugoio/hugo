---
aliases:
- /extras/highlight/
date: 2013-07-01
menu:
  main:
    parent: extras
next: /extras/toc
prev: /extras/shortcodes
title: Syntax Highlighting
weight: 90
---

Hugo provides the ability for you to highlight source code in two different
ways &mdash; either pre-processed server side from your content, or to defer
the processing to the client side, using a JavaScript library. The advantage of
server side is that it doesn’t depend on a JavaScript library and consequently
works very well when read from an RSS feed. The advantage of client side is that
it doesn’t cost anything when building your site and some of the highlighting
scripts available cover more languages than Pygments does.

## Server-side

For the pre-processed approach, Highlighting is performed by an external
Python-based program called [Pygments](http://pygments.org/) and is triggered
via an embedded shortcode (see example below). If Pygments is absent from the path, it will silently simply pass the content along unhighlighted.

### Pygments

If you have never worked with Pygments before, here is a brief primer.

+ Install Python from [python.org](https://www.python.org/downloads/). Version 2.7.x is already sufficient.
+ Run `pip install Pygments` in order to install Pygments. Once installed, Pygments gives you a command `pygmentize`. Make sure it sits in your PATH, otherwise Hugo cannot find it.

On Debian and Ubuntu systems, you may also install Pygments by running `sudo apt-get install python3-pygments`.

Hugo gives you two options that you can set with the variable `pygmentsuseclasses` (default `false`) in `config.toml` (or `config.yaml`).

1. Color-codes for highlighting keywords are directly inserted if `pygmentsuseclasses = false` (default). See in the example below. The color-codes depend on your choice of the `pygmentsstyle` (default `"monokai"`). You can explore the different color styles on [pygments.org](http://pygments.org/) after inserting some example code.
2. If you choose `pygmentsuseclasses = true`, Hugo includes class names in your code instead of color-codes. For class-names to be meaningful, you need to include a `.css`-file in your website representing your color-scheme. You can either generate this `.css`-files according to this [description](http://pygments.org/docs/cmdline/) or download the standard ones from the [GitHub pygments-css repository](https://github.com/richleland/pygments-css).

### Usage
Highlighting is carried out via the in-built shortcode `highlight`. `highlight` takes exactly one required parameter of language and requires a
closing shortcode.

### Example
If you want to highlight code, you need to either fence the code with ``` according to GitHub Flavored Markdown or preceed each line with 4 spaces to identify each line as a line of code.

Not doing either will result in the text being rendered as HTML. This will prevent Pygments highlighting from working.

```
{{</* highlight html */>}}
<section id="main">
  <div>
    <h1 id="title">{{ .Title }}</h1>
    {{ range .Data.Pages }}
        {{ .Render "summary"}}
    {{ end }}
  </div>
</section>
{{</* /highlight */>}}
```

### Example Output

    <span style="color: #f92672">&lt;section</span> <span style="color: #a6e22e">id=</span><span style="color: #e6db74">&quot;main&quot;</span><span style="color: #f92672">&gt;</span>
      <span style="color: #f92672">&lt;div&gt;</span>
       <span style="color: #f92672">&lt;h1</span> <span style="color: #a6e22e">id=</span><span style="color: #e6db74">&quot;title&quot;</span><span style="color: #f92672">&gt;</span>{{ .Title }}<span style="color: #f92672">&lt;/h1&gt;</span>
        {{ range .Data.Pages }}
            {{ .Render &quot;summary&quot;}}
        {{ end }}
      <span style="color: #f92672">&lt;/div&gt;</span>
    <span style="color: #f92672">&lt;/section&gt;</span>


### Disclaimers

 * **Warning:** Pygments is relatively slow. Expect much longer build times when using server-side highlighting.
 * Languages available depends on your Pygments installation.
 * We have sought to have the simplest interface possible, which consequently
limits configuration. An ambitious user is encouraged to extend the current
functionality to offer more customization.


## Client-side

Alternatively, code highlighting can be done in client-side JavaScript.

Client-side syntax highlighting is very simple to add. You'll need to pick
a library and a corresponding theme. Some popular libraries are:

- [Highlight.js]
- [Rainbow]
- [Syntax Highlighter]
- [Google Prettify]

This example uses the popular [Highlight.js] library, hosted by [Yandex], a
popular Russian search engine.

In your `./layouts/partials/` (or `./layouts/chrome/`) folder, depending on your specific theme, there
will be a snippet that will be included in every generated HTML page, such
as `header.html` or `header.includes.html`. Simply add:

    <link rel="stylesheet" href="https://yandex.st/highlightjs/8.0/styles/default.min.css">
    <script src="https://yandex.st/highlightjs/8.0/highlight.min.js"></script>
    <script>hljs.initHighlightingOnLoad();</script>

You can of course use your own copy of these files, typically in `./static/`.

[Highlight.js]: http://highlightjs.org/
[Rainbow]: http://craig.is/making/rainbows
[Syntax Highlighter]: http://alexgorbatchev.com/SyntaxHighlighter/
[Google Prettify]: https://code.google.com/p/google-code-prettify/
[Yandex]: http://yandex.ru/

Please see individual libraries documentation for how to implement the JavaScript-based libraries.
