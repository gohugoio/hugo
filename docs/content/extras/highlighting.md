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
weight: 50
---

Hugo provides the ability for you to highlight source code in two different
ways &mdash; either pre-processed server side from your content, or to defer
the processing to the client side, using a JavaScript library. The advantage of
server side is that it doesn’t depend on a JavaScript library and consequently
works very well when read from an RSS feed. The advantage of client side is that
it doesn’t cost anything when building your site and some of the highlighting 
scripts available cover more languages than Pygments does.

For the pre-processed approach, Highlighting is performed by an external
Python-based program called [Pygments](http://pygments.org) and is triggered
via an embedded shortcode. If Pygments is absent from the path, it will
silently simply pass the content along unhighlighted.

## Server-side

### Disclaimers

 * **Warning:** Pygments is relatively slow. Expect much longer build times when using server-side highlighting.
 * Languages available depends on your Pygments installation.
 * Styles are inline in order to be supported in syndicated content when references
to style sheets are not carried over.
 * We have sought to have the simplest interface possible, which consequently
limits configuration. An ambitious user is encouraged to extend the current
functionality to offer more customization.
* You can change appearance with config options `pygmentsstyle`(default
`"monokai"`) and `pygmentsuseclasses`(defaut `false`).

### Usage
Highlight takes exactly one required parameter of language and requires a
closing shortcode.

### Example
The example has an extra space between the “{{” and “%” characters to prevent rendering here.

    {{ % highlight html %}}
    <section id="main">
      <div>
       <h1 id="title">{{ .Title }}</h1>
        {{ range .Data.Pages }}
            {{ .Render "summary"}}
        {{ end }}
      </div>
    </section>
    {{ % /highlight %}}


### Example Output

    <span style="color: #f92672">&lt;section</span> <span style="color: #a6e22e">id=</span><span style="color: #e6db74">&quot;main&quot;</span><span style="color: #f92672">&gt;</span>
      <span style="color: #f92672">&lt;div&gt;</span>
       <span style="color: #f92672">&lt;h1</span> <span style="color: #a6e22e">id=</span><span style="color: #e6db74">&quot;title&quot;</span><span style="color: #f92672">&gt;</span>{{ .Title }}<span style="color: #f92672">&lt;/h1&gt;</span>
        {{ range .Data.Pages }}
            {{ .Render &quot;summary&quot;}}
        {{ end }}
      <span style="color: #f92672">&lt;/div&gt;</span>
    <span style="color: #f92672">&lt;/section&gt;</span>

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
