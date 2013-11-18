---
title: "Go Templates"
date: "2013-07-01"
groups: ["layout"]
groups_weight: 80
draft: true
---

Hugo uses the excellent [golang][] [html/template][gohtmltemplate] library for
its template engine.
It is an extremely lightweight engine that provides a very small amount of
logic.
In our experience that it is just the right amount of logic to be able to
create a good static website.

This is a brief primer on using go templates. The [golang docs][gohtmltemplate]
provide more details.

In your top-level configuration file (eg, `config.yaml`) you can define site
parameters, which are values which will be available to you in chrome.

For instance, you might declare:

```yaml
params:
  CopyrightHTML: "Copyright &#xA9; 2013 John Doe. All Rights Reserved."
  TwitterUser: "spf13"
  SidebarRecentLimit: 5
```

Within a footer layout, you might then declare a `<footer>` which is only
provided if the `CopyrightHTML` parameter is provided, and if it is given,
you would declare it to be HTML-safe, so that the HTML entity is not escaped
again.  This would let you easily update just your top-level config file each
January 1st, instead of hunting through your templates.

```
{{if .Site.Params.CopyrightHTML}}<footer>
<div class="text-center">{{.Site.Params.CopyrightHTML | safeHtml}}</div>
</footer>{{end}}
```

An alternative way of writing the "if" and then referencing the same value is
to use "with" instead, which rebinds the pointer `.` within its scope, and
elides the scope if the variable is absent:

```
{{with .Site.Params.TwitterUser}}<span class="twitter">
<a href="https://twitter.com/{{.}}" rel="author">
<img src="/images/twitter.png" width="48" height="48" title="Twitter: {{.}}"
 alt="Twitter"></a>
</span>{{end}}
```

Finally, if you want to pull "magic constants" out of your layouts, you can do
so, such as in this example:

```
<nav class="recent">
  <h1>Recent Posts</h1>
  <ul>{{range first .Site.Params.SidebarRecentLimit .Site.Recent}}
    <li><a href="{{.RelPermalink}}">{{.Title}}</a></li>
  {{end}}</ul>
</nav>
```


[golang]: <http://golang.org/>
[gohtmltemplate]: <http://golang.org/pkg/html/template/>
