---
aliases:
- /layout/chrome/
lastmod: 2016-01-01
date: 2013-07-01
menu:
  main:
    parent: layout
next: /templates/rss
prev: /templates/blocks/
title: Partial Templates
weight: 80
toc: true
---

In practice, it's very convenient to split out common template portions into a
partial template that can be included anywhere. As you create the rest of your
templates, you will include templates from the ``/layouts/partials` directory
or from arbitrary subdirectories like `/layouts/partials/post/tag`.

Partials are especially important for themes as it gives users an opportunity
to overwrite just a small part of your theme, while maintaining future compatibility.

Theme developers may want to include a few partials with empty HTML
files in the theme just so end users have an easy place to inject their
customized content.

I've found it helpful to include a header and footer template in
partials so I can include those in all the full page layouts.  There is
nothing special about header.html and footer.html other than they seem
like good names to use for inclusion in your other templates.

    ▾ layouts/
      ▾ partials/
          header.html
          footer.html

## Partial vs Template

Version v0.12 of Hugo introduced the `partial` call inside the template system.
This is a change to the way partials were handled previously inside the
template system. In earlier versions, Hugo didn’t treat partials specially, and
you could include a partial template with the `template` call in the standard
template language.

With the addition of the theme system in v0.11, it became apparent that a theme
& override-aware partial was needed.

When using Hugo v0.12 and above, please use the `partial` call (and leave out
the “partial/” path). The old approach would still work, but wouldn’t benefit from
the ability to have users override the partial theme file with local layouts.

## Example header.html
This header template is used for [spf13.com](http://spf13.com/):

    <!DOCTYPE html>
    <html class="no-js" lang="en-US" prefix="og: http://ogp.me/ns# fb: http://ogp.me/ns/fb#">
    <head>
        <meta charset="utf-8">

        {{ partial "meta.html" . }}

        <base href="{{ .Site.BaseURL }}">
        <title> {{ .Title }} : spf13.com </title>
        <link rel="canonical" href="{{ .Permalink }}">
        {{ if .RSSLink }}<link href="{{ .RSSLink }}" rel="alternate" type="application/rss+xml" title="{{ .Title }}" />{{ end }}

        {{ partial "head_includes.html" . }}
    </head>
    <body lang="en">

## Example footer.html
This footer template is used for [spf13.com](http://spf13.com/):

    <footer>
      <div>
        <p>
        &copy; 2013-14 Steve Francia.
        <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons Attribution">Some rights reserved</a>;
        please attribute properly and link back. Hosted by <a href="http://servergrove.com">ServerGrove</a>.
        </p>
      </div>
    </footer>
    <script type="text/javascript">

      var _gaq = _gaq || [];
      _gaq.push(['_setAccount', 'UA-XYSYXYSY-X']);
      _gaq.push(['_trackPageview']);

      (function() {
        var ga = document.createElement('script');
        ga.src = ('https:' == document.location.protocol ? 'https://ssl' :
            'http://www') + '.google-analytics.com/ga.js';
        ga.setAttribute('async', 'true');
        document.documentElement.firstChild.appendChild(ga);
      })();

    </script>
    </body>
    </html>

To reference a partial template stored in a subfolder, e.g. `/layouts/partials/post/tag/list.html`, call it this way:

     {{ partial "post/tag/list" . }}

Note that the subdirectories you create under /layouts/partials can be named whatever you like.

For more examples of referencing these templates, see
[single content templates](/templates/content/),
[list templates](/templates/list/) and
[homepage templates](/templates/homepage/).


## Variable scoping

As you might have noticed, `partial` calls receive two parameters.

1. The first is the name of the partial and determines the file
location to be read.
2. The second is the variables to be passed down to the partial.

This means that the partial will _only_ be able to access those variables. It is
isolated and has no access to the outer scope. From within the
partial, `$.Var` is equivalent to `.Var`

## Cached Partials

The `partialCached` template function can offer significant performance gains
for complex templates that don't need to be rerendered upon every invocation.
The simplest usage is as follows:

    {{ partialCached "footer.html" . }}

You can also pass additional parameters to `partialCached` to create *variants* of the cached partial.
For example, say you have a complex partial that should be identical when rendered for pages within the same section.
You could use a variant based upon section so that the partial is only rendered once per section:

    {{ partialCached "footer.html" . .Section }}

If you need to pass additional parameters to create unique variants,
you can pass as many variant parameters as you need:

    {{ partialCached "footer.html" . .Params.country .Params.province }}

Note that the variant parameters are not made available to the underlying partial template.
They are only use to create a unique cache key.
