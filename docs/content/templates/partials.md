---
aliases:
- /layout/chrome/
date: 2013-07-01
menu:
  main:
    parent: layout
next: /templates/rss
prev: /templates/views
title: Partial Templates
weight: 80
---

In practice, it's very convenient to split out common template portions into a
partial template that can be included anywhere. As you create the rest of your
templates, you will include templates from the /layout/partials directory.

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

By ensuring that we only reference [variables](/layout/variables/)
used for both nodes and pages, we can use the same partials for both.

## Partial vs Template 

Version v0.12 of Hugo introduced the `partial` call inside the template system.
This is a change to the way partials were handled previously inside the
template system. In earlier versions, Hugo didn’t treat partials specially, and
you could include a partial template with the `template` call in the standard
template language.

With the addition of the theme system in v0.11, it became apparent that a theme
& override aware partial was needed.

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

        <base href="{{ .Site.BaseUrl }}">
        <title> {{ .Title }} : spf13.com </title>
        <link rel="canonical" href="{{ .Permalink }}">
        {{ if .RSSlink }}<link href="{{ .RSSlink }}" rel="alternate" type="application/rss+xml" title="{{ .Title }}" />{{ end }}

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

**For examples of referencing these templates, see [single content
templates](/templates/content/), [list templates](/templates/list/) and [homepage templates](/templates/homepage/).**
