---
aliases:
- /layout/related/
date: 2014-12-13
linktitle: "Related documents"
menu:
  main:
    parent: layout
next: /taxonomies/overview
notoc: true
prev: /templates/404
title: Related Documents Templates
weight: 100
---

You can add in your page the following code, in order to enable the page suggestion 
according the keyword used in the front matter of the currently rendering page.

## related.html
Example of a partial to use

    <div class="post-related">
      <section id="related">
        <h4>For more reading ...</h4>
        <ul id="list">
          {{ range first 5 .RelatedPages }}
          <li>
            <a href="{{ .Permalink }}">{{ .Title }} - {{ .ReadingTime }} Minutes</a>
          </li>
          {{ end }}
        </ul>
      </section>
    </div>