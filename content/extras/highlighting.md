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
toc: true
---

Hugo provides the ability for you to highlight source code in _two different ways_ &mdash; either pre-processed server side from your content, or to defer the processing to the client side, using a JavaScript library. 

**The advantage of server side** is that it doesn’t depend on a JavaScript library and consequently works very well when read from an RSS feed. 

**The advantage of client side** is that it doesn’t cost anything when building your site and some of the highlighting scripts available cover more languages than Pygments does.

## Server-side

For the pre-processed approach, highlighting is performed by an external Python-based program called [Pygments](http://pygments.org/) and is triggered via an embedded Hugo shortcode (see example below). If Pygments is absent from the path, it will silently simply pass the content along unhighlighted.

### Pygments

If you have never worked with Pygments before, here is a brief primer:

+ Install Python from [python.org](https://www.python.org/downloads/). Version 2.7.x is already sufficient.
+ Run `pip install Pygments` in order to install Pygments. Once installed, Pygments gives you a command `pygmentize`. Make sure it sits in your PATH, otherwise Hugo cannot find it.

On Debian and Ubuntu systems, you may also install Pygments by running `sudo apt-get install python3-pygments`.

Hugo gives you two options that you can set with the variable `pygmentsuseclasses` (default `false`) in `config.toml` (or `config.yaml`).

1. Color-codes for highlighting keywords are directly inserted if `pygmentsuseclasses = false` (default). See in the example below. The color-codes depend on your choice of the `pygmentsstyle` (default `"monokai"`). You can explore the different color styles on [pygments.org](http://pygments.org/) after inserting some example code.
2. If you choose `pygmentsuseclasses = true`, Hugo includes class names in your code instead of color-codes. For class-names to be meaningful, you need to include a `.css`-file in your website representing your color-scheme. You can either generate this `.css`-files according to this [description](http://pygments.org/docs/cmdline/) or download the standard ones from the [GitHub pygments-css repository](https://github.com/richleland/pygments-css).

### Usage

Highlighting is carried out via the in-built shortcode `highlight`. `highlight` takes exactly one required parameter of language, and requires a closing shortcode. Note that `highlight` is _not_ used for client-side javascript highlighting. 

### Example

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

```
<span style="color: #f92672">&lt;section</span> <span style="color: #a6e22e">id=</span><span style="color: #e6db74">&quot;main&quot;</span><span style="color: #f92672">&gt;</span>
  <span style="color: #f92672">&lt;div&gt;</span>
    <span style="color: #f92672">&lt;h1</span> <span style="color: #a6e22e">id=</span><span style="color: #e6db74">&quot;title&quot;</span><span style="color: #f92672">&gt;</span>{{ .Title }}<span style="color: #f92672">&lt;/h1&gt;</span>
    {{ range .Data.Pages }}
      {{ .Render &quot;summary&quot;}}
    {{ end }}
  <span style="color: #f92672">&lt;/div&gt;</span>
<span style="color: #f92672">&lt;/section&gt;</span>
```

### Options

Options to control highlighting can be added as a quoted, comma separated key-value list as the second argument in the shortcode. The example below will highlight as language `go` with inline line numbers, with line number 2 and 3 highlighted.

```
{{</* highlight go "linenos=inline,hl_lines=2 3" */>}}
var a string
var b string
var c string
var d string
{{</* / highlight */>}}
```

Supported keywords:  `style`, `encoding`, `noclasses`, `hl_lines`, `linenos`. Note that `style` and `noclasses` will override the similar setting in the global config.

The keywords are the same you would using with Pygments from the command line, see the [Pygments doc](http://pygments.org/docs/) for more info.


### Disclaimers

 * Pygments is relatively slow and _causes a performance hit when building your site_, but Hugo has been designed to cache the results to disk. 
 * Languages available depends on your Pygments installation.

## Client-side

Alternatively, code highlighting can be done in client-side JavaScript.

Client-side syntax highlighting is very simple to add. You'll need to pick
a library and a corresponding theme. Some popular libraries are:

- [Highlight.js]
- [Prism]
- [Rainbow]
- [Syntax Highlighter]
- [Google Prettify]

### Highlight.js example

This example uses the popular [Highlight.js] library, hosted by [Yandex], a popular Russian search engine.

In your `./layouts/partials/` (or `./layouts/chrome/`) folder, depending on your specific theme, there will be a snippet that will be included in every generated HTML page, such as `header.html` or `header.includes.html`. Simply add the css and js to initialize [Highlight.js]:

~~~
<link rel="stylesheet" href="https://yandex.st/highlightjs/8.0/styles/default.min.css">
<script src="https://yandex.st/highlightjs/8.0/highlight.min.js"></script>
<script>hljs.initHighlightingOnLoad();</script>
~~~

### Prism example

Prism is another popular highlighter library, used on some major sites. Similar to Highlight.js, you simply load `prism.css` in your `<head>` via whatever Hugo partial template is creating that part of your pages, like so: 

```html
...
<link href="/css/prism.css" rel="stylesheet" />
...
```

... and add `prism.js` near the bottom of your `<body>` tag, again in whatever Hugo partial template is appropriate for your site or theme. 

```html
...
<script src="/js/prism.js"></script>
...
</body>
```

In this example, the local paths indicate that your own copy of these files are being added to the site, typically under `./static/`.

### Using Client-side highlighting

To use client-side highlighting, most of these javascript libraries expect your code to be wrapped in semantically correct `<code>` tags, with the language expressed in a class attribute on the `<code>` tag, such as `class="language-abc"`, where the `abc` is the code the highlighter script uses to represent that language. 

The script would be looking for classes like `language-go`, `language-html`, or `language-css`. If you look at the page's source, it would be marked up like so: 

~~~html
<pre>
<code class="language-css">
body {
  font-family: "Noto Sans", sans-serif;
}
</code>
</pre>
~~~

The markup in your content pages (e.g. `my-css-tutorial.md`) needs to look like the following, with the name of the language to be highlighted entered directly after the first "fence", in a fenced code block:

<pre><code class="language-css">&#126;&#126;&#126;css
body {
  font-family: "Noto Sans", sans-serif;
}
&#126;&#126;&#126;</code></pre>

When passed through the highlighter script, it would yield something like this output when viewed on your rendered page:

~~~css
body {
  font-family: "Noto Sans", sans-serif;
}
~~~

Please see individual libraries' documentation for how to implement each of the JavaScript-based libraries.

[Prism]: http://prismjs.com
[Highlight.js]: http://highlightjs.org/
[Rainbow]: http://craig.is/making/rainbows
[Syntax Highlighter]: http://alexgorbatchev.com/SyntaxHighlighter/
[Google Prettify]: https://code.google.com/p/google-code-prettify/
[Yandex]: http://yandex.ru/

