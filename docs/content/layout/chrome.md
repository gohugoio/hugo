---
title: "Chrome Templates"
date: "2013-07-01"
weight: 80
menu:
  main:
    parent: 'layout'
---
Chrome is a convention to create templates that are used by the other templates
throughout the site. There is nothing special about the name "chrome", feel free
to provide and use your own.

It's not a requirement to have this, but in practice it's very convenient. Hugo doesn't
know anything about Chrome, it's simply a convention that you may likely find
beneficial. As you create the rest of your templates you will include templates
from the /layout/chrome directory.

I've found it helpful to include a header and footer template in Chrome so I can
include those in the other full page layouts (index.html, indexes/
type/single.html).  There is nothing special about header.html and footer.html
other than they seem like good names to use for inclusion in your other
templates.

    ▾ layouts/
      ▾ chrome/
          header.html
          footer.html

By ensuring that we only reference [variables](/layout/variables/) variables
used for both nodes and pages we can use the same chrome for both.

## example header.html
This header template is used for [spf13.com](http://spf13.com).

    <!DOCTYPE html>
    <html class="no-js" lang="en-US" prefix="og: http://ogp.me/ns# fb: http://ogp.me/ns/fb#">
    <head>
        <meta charset="utf-8">

        {{ template "chrome/meta.html" . }}

        <base href="{{ .Site.BaseUrl }}">
        <title> {{ .Title }} : spf13.com </title>
        <link rel="canonical" href="{{ .Permalink }}">
        {{ if .RSSlink }}<link href="{{ .RSSlink }}" rel="alternate" type="application/rss+xml" title="{{ .Title }}" />{{ end }}

        {{ template "chrome/head_includes.html" . }}
    </head>
    <body lang="en">

## example footer.html
This header template is used for [spf13.com](http://spf13.com).

    <footer>
      <div>
        <p>
        &copy; 2013 Steve Francia.
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

**For examples of referencing these templates, see [content
templates](/layout/content/) and [homepage templates](/layout/homepage/)**
