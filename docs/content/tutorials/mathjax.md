---
author: Spencer Lyon
date: 2014-03-20
menu:
  main:
    parent: tutorials
next: /tutorials/migrate-from-jekyll
prev: /tutorials/github-pages-blog
title: MathJax Support
weight: 10
---

## What is MathJax?

[MathJax](http://www.mathjax.org/) is a JavaScript library that allows the display of mathematical expressions described via a LaTeX-style syntax in the HTML (or Markdown) source of a web page. As it is a pure a JavaScript library, getting it to work within Hugo is fairly straightforward, but does have some oddities that will be discussed here.

This is not an introduction into actually using MathJax to render typeset mathematics on your website. Instead, this page is a collection of tips and hints for one way to get MathJax working on a website built with Hugo.

## Enabling MathJax

The first step is to enable MathJax on pages that you would like to have typeset math. There are multiple ways to do this (adventerous readers can consult the [Loading and Configuring](http://docs.mathjax.org/en/latest/configuration.html) section of the MathJax documentation for additional methods of including MathJax), but the easiest way is to use the secure MathJax CDN by including the following HTML snippet in the source of a page:

    <script type="text/javascript"
      src="https://cdn.mathjax.org/mathjax/latest/MathJax.js?config=TeX-AMS-MML_HTMLorMML">
    </script>

One way to ensure that this code is included in all pages is to put it in one of the templates that live in the `layouts/partials/` directory. For example, I have included this in the bottom of my template `footer.html` because I know that the footer will be included in every page of my website.

### Options and Features

MathJax is a stable open-source library with many features. I encourage the interested reader to view the [MathJax Documentation](http://docs.mathjax.org/en/latest/index.html), specifically the sections on [Basic Usage](http://docs.mathjax.org/en/latest/index.html#basic-usage) and [MathJax Configuration Options](http://docs.mathjax.org/en/latest/index.html#mathjax-configuration-options).

## Issues with Markdown

After enabling MathJax, any math entered in-between proper markers (see documentation) will be processed and typeset in the web page. One issue that comes up, however, with Markdown is that the underscore character (`_`) is interpreted by Markdown as a way to wrap text in `emph` blocks while LaTeX (MathJax) interprets the underscore as a way to create a subscript. This "double speak" of the underscore can result in some unexpected and unwanted behavior.

### Solution

There are multiple ways to remedy this problem. One solution is to simply escape each underscore in your math code by entering `\_` instead of `_`. This can become quite tedious if the equations you are entering are full of subscripts.

Another option is to tell Markdown to treat the MathJax code as verbatim code and not process it. One way to do this is to wrap the math expression inside a `<div>` `</div>` block. Markdown would ignore these sections and they would get passed directly on to MathJax and processed correctly. This works great for display style mathematics, but for inline math expressions the line break induced by the `<div>` is not acceptable. The syntax for instructing Markdown to treat inline text as verbatim is by wrapping it in backticks (`` ` ``). You might have noticed, however, that the text included in between backticks is rendered differently than standard text (on this site these are items highlighted in red). To get around this problem, we could create a new CSS entry that would apply standard styling to all inline verbatim text that includes MathJax code. Below I will show the HTML and CSS source that would accomplish this (note this solution was adapted from [this blog post](http://doswa.com/2011/07/20/mathjax-in-markdown.html)---all credit goes to the original author).

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

As before, this content should be included in the HTML source of each page that will be using MathJax. The next code snippet contains the CSS that is used to have verbatim MathJax blocks render with the same font style as the body of the page.


    code.has-jax {font: inherit;
                  font-size: 100%;
                  background: inherit;
                  border: inherit;
                  color: #515151;}

In the CSS snippet, notice the line `color: #515151;`. `#515151` is the value assigned to the `color` attribute of the `body` class in my CSS. In order for the equations to fit in with the body of a web page, this value should be the same as the color of the body.

### Usage

With this setup, everything is in place for a natural usage of MathJax on pages generated using Hugo. In order to include inline mathematics, just put LaTeX code in between `` `$ TeX Code $` `` or `` `\( TeX Code \)` ``. To include display style mathematics, just put LaTeX code in between `<div>$$TeX Code$$</div>`. All the math will be properly typeset and displayed within your Hugo generated web page!
