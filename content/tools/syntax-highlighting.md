---
title: Syntax Highlighting
linktitle:
description: Hugo provides server-side syntax highlighting via Pygments and, like most static site generators, works very well with client-side (JavaScript) syntax highlighting libraries as well.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
keywords: [highlighting,pygments,code blocks,syntax]
categories: [developer tools]
menu:
  docs:
    parent: "tools"
    weight: 20
weight: 20
sections_weight: 20
draft: false
aliases: [/extras/highlighting/,/extras/highlight/]
toc: true
---

Hugo can highlight source code in _two different ways_&mdash;either pre-processed server side from your content or to defer the processing to the client side, using a JavaScript library.

## Server-side

For the pre-processed approach, highlighting is performed by an external Python-based program called [Pygments](http://pygments.org/) and is triggered via an embedded Hugo shortcode (see [example](#example-highlight-shortcode-input) below). If Pygments is absent from the path, it will silently simply pass the content along without highlighting.

### Server-side Advantages

The advantages of server-side syntax highlighting are that it doesn’t depend on a JavaScript library and, consequently, works very well when read from an RSS feed.

### Pygments

If you have never worked with Pygments before, here is a brief primer:

+ Install Python from [python.org](https://www.python.org/downloads/). Version 2.7.x is already sufficient.
+ Run `pip install Pygments` in order to install Pygments. Once installed, Pygments gives you a command `pygmentize`. Make sure it sits in your PATH; otherwise, Hugo will not be able to find and use it.

On Debian and Ubuntu systems, you may also install Pygments by running `sudo apt-get install python3-pygments`.

Hugo gives you two options that you can set with the variable `pygmentsuseclasses` (default `false`) in your [site configuration](/getting-started/configuration/).

1. Color codes for highlighting keywords are directly inserted if `pygmentsuseclasses = false` (default). The color codes depend on your choice of the `pygmentsstyle` (default = `"monokai"`). You can explore the different color styles on [pygments.org](http://pygments.org/) after inserting some example code.
2. If you choose `pygmentsuseclasses = true`, Hugo includes class names in your code instead of color codes. For class-names to be meaningful, you need to include a `.css` file in your website representing your color scheme. You can either generate this `.css` files according to the [description from the Pygments documentation](http://pygments.org/docs/cmdline/) or download the one of the many pre-built color schemes from [Pygment's GitHub css repository](https://github.com/richleland/pygments-css).

### Server-side Usage

Highlighting is carried out via the [built-in shortcode](/content-management/shortcodes/) `highlight`. `highlight` takes exactly one required parameter for the programming language to be highlighted and requires a closing shortcode. Note that `highlight` is *not* used for client-side javascript highlighting.

### Example `highlight` Shortcode Input

{{< code file="example-highlight-shortcode-input.md" >}}
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
{{< /code >}}

### Example `highlight` Shortcode Output

{{< output file="example-highlight-shortcode-output.html" >}}
<span style="color: #f92672">&lt;section</span> <span style="color: #a6e22e">id=</span><span style="color: #e6db74">&quot;main&quot;</span><span style="color: #f92672">&gt;</span>
  <span style="color: #f92672">&lt;div&gt;</span>
    <span style="color: #f92672">&lt;h1</span> <span style="color: #a6e22e">id=</span><span style="color: #e6db74">&quot;title&quot;</span><span style="color: #f92672">&gt;</span>{{ .Title }}<span style="color: #f92672">&lt;/h1&gt;</span>
    {{ range .Data.Pages }}
      {{ .Render &quot;summary&quot;}}
    {{ end }}
  <span style="color: #f92672">&lt;/div&gt;</span>
<span style="color: #f92672">&lt;/section&gt;</span>
{{< /output >}}

### Options

Options for controlling highlighting can be added in the second argument as a quoted, comma-separated key-value list. The example below will syntax highlight in `go` with inline line numbers and line numbers 2 and 3 highlighted.

```
{{</* highlight go "linenos=inline,hl_lines=2 3" */>}}
var a string
var b string
var c string
var d string
{{</* / highlight */>}}
```

The `highlight` shortcode includes the following supported keywords:

* `style`
* `encoding`
* `noclasses`
* `hl_lines`
* `linenos`

Note that `style` and `noclasses` will override the similar setting in the [global config](/getting-started/configuration/).

The keywords in the `highlight` shortcode mirror those of Pygments from the command line. See the [Pygments documentation](http://pygments.org/docs/) for more information.

### Code Fences

It is also possible to add syntax highlighting with GitHub flavored code fences. To enable this, set the `PygmentsCodeFences` to `true` in Hugo's [configuration file](/getting-started/configuration/);

```
```
<section id="main">
  <div>
    <h1 id="title">{{ .Title }}</h1>
    {{ range .Data.Pages }}
      {{ .Render "summary"}}
    {{ end }}
  </div>
</section>
```
```

{{% note "Disclaimers on Pygments" %}}
* Pygments is relatively slow and _causes a performance hit when building your site_, but Hugo has been designed to cache the results to disk.
* The caching can be turned off by setting the `--ignoreCache` flag to `true`.
* The languages available for highlighting depend on your Pygments installation.
{{% /note %}}

## Client-side

Alternatively, code highlighting can be applied to your code blocks in client-side JavaScript.

Client-side syntax highlighting is very simple to add. You'll need to pick
a library and a corresponding theme. Some popular libraries are:

- [Highlight.js]
- [Prism]
- [Rainbow]
- [Syntax Highlighter]
- [Google Prettify]

### Client-side Advantages

The advantages of client-side syntax highlighting are that it doesn’t cost anything when building your site, and some of the highlighting scripts available cover more languages than Pygments does.

### Highlight.js Example

This example uses the popular [Highlight.js] library, hosted by [Yandex], a popular Russian search engine.

In your `./layouts/partials/` (or `./layouts/chrome/`) folder, depending on your specific theme, there will be a snippet that will be included in every generated HTML page, such as `header.html` or `header.includes.html`. Simply add the css and js to initialize [Highlight.js]:

```
<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/highlight.js/9.6.0/styles/default.min.css">
<script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/9.6.0/highlight.min.js"></script>
<script>hljs.initHighlightingOnLoad();</script>
```

### Prism example

Prism is another popular highlighter library and is used on some major sites.
The [download section of the prism.js website][prismdownload] is simple to use and affords you a high degree of customization to pick only the languages you'll be using on your site.

Similar to Highlight.js, you simply load `prism.css` in your `<head>` via whatever Hugo partial template is creating that part of your pages:

```
...
<link href="/css/prism.css" rel="stylesheet" />
...
```

Add `prism.js` near the bottom of your `<body>` tag in whatever Hugo partial template is appropriate for your site or theme.

```
...
<script src="/js/prism.js"></script>
</body>
```

In this example, the local paths indicate that your downloaded copy of these files are being added to the site, typically under `./static/css/` and `./static/js/`, respectively.

### Client-side Usage

To use client-side highlighting, most of these javascript libraries expect your code to be wrapped in semantically correct `<code>` elements with language-specific class attributes. For example, a code block for HTML would have `class="language-html"`.

The client-side highlighting script therefore looks for programming language classes according to this convention: `language-go`, `language-html`, `language-css`, `language-bash`, etc. If you look at the page's source, you might see something like the following:

```
<pre>
  <code class="language-css">
  body {
    font-family: "Noto Sans", sans-serif;
  }
  </code>
</pre>
```

If you are using markdown, your content pages needs to use the following syntax, with the name of the language to be highlighted entered directly after the first "fence." A fenced code block can be noted by opening and closing triple tilde <kbd>~</kbd> or triple back ticks <kbd>`</kbd>:

{{< nohighlight >}}
~~~css
body {
  font-family: "Noto Sans", sans-serif;
}
~~~
{{< /nohighlight >}}

Here is the same example but with triple back ticks to denote the fenced code block:

{{< nohighlight >}}
```css
body {
  font-family: "Noto Sans", sans-serif;
}
```
{{< /nohighlight >}}

Passing the above examples through the highlighter script would yield the following markup:

{{< nohighlight >}}
&lt;pre&gt;&lt;code class="language-css hljs"&gt;;&lt;span class="hljs-selector-tag"&gt;body&lt;/span&gt; {
  &lt;span class="hljs-attribute"&gt;font-family&lt;/span&gt;: &ltspan class="hljs-string"&gt;"Noto Sans"&lt;/span&gt;, sans-serif;
}
{{< /nohighlight >}}

In the case of the coding color scheme used by the Hugo docs, the resulting output would then look like the following to the website's end users:

```
body {
  font-family: "Noto Sans", sans-serif;
}
```

Please see individual libraries' documentation for how to implement each of the JavaScript-based libraries.

[Prism]: http://prismjs.com
[prismdownload]: http://prismjs.com/download.html
[Highlight.js]: http://highlightjs.org/
[Rainbow]: http://craig.is/making/rainbows
[Syntax Highlighter]: http://alexgorbatchev.com/SyntaxHighlighter/
[Google Prettify]: https://github.com/google/code-prettify
[Yandex]: http://yandex.ru/
