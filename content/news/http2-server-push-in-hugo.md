---
title: "HTTP/2 Server Push in Hugo"
date: 2017-07-24T18:36:00+02:00
description: >
    As every page in Hugo can be output to multiple formats, it is easy to create Netlify's _redirects and _headers files on the fly.
categories: [blog]
keywords: []
slug: "http2-server-push-in-hugo"
aliases: []
author: bep
images:
- images/gohugoio-card-1.png
---

**Netlify** recently announced support for [HTTP/2 server push](https://www.netlify.com/blog/2017/07/18/http/2-server-push-on-netlify/), and we have now added it to the **gohugo.io** sites for the main `CSS` and `JS` bundles, along with server-side 301 redirect support. 

If you navigate to https://gohugo.io and look in the Chrome developer network console, you should now see `Push` as the new source ("Initiator") for the `CSS` and `JSS`:

{{< figure src="/images/blog/hugo-http2-push.png" caption="Network log for https://gohugo.io" >}}

**Setting up this in Hugo was easy:**

## 1. Configure Netlify Output Formats

Add a new custom media type and two new output formats to `config.toml`. For more on output formats in Hugo, see [Custom Output Formats](/templates/output-formats/).
```bash
[outputs]
home = [ "HTML", "RSS", "REDIR", "HEADERS" ]

[mediaTypes]
[mediaTypes."text/netlify"]
suffix = ""
delimiter = ""

[outputFormats]
[outputFormats.REDIR]
mediatype = "text/netlify"
baseName = "_redirects"
isPlainText = true
notAlternative = true
[outputFormats.HEADERS]
mediatype = "text/netlify"
baseName = "_headers"
isPlainText = true
notAlternative = true
```
## 2. Add Template For the _headers File 

Add `layouts/index.headers`:

```bash
/*
  X-Frame-Options: DENY
  X-XSS-Protection: 1; mode=block
  X-Content-Type-Options: nosniff
  Referrer-Policy: origin-when-cross-origin
*/
  Link: <{{ "dist/app.bundle.js" | relURL }}>; rel=preload; as=script
  Link: <{{ "dist/main.css" | relURL }}>; rel=preload; as=style
```
The template above creates both a security header definition and a HTTP/2 server push configuration.

Also note that this is a template for the home page, so the full `Page` with its `Site` and many variables are available. You can also use `partial` to include other templates.




## 3. Add Template For the _redirects File 
Add `layouts/index.redir`:
```bash
# Netlify redirects. See https://www.netlify.com/docs/redirects/
{{  range $p := .Site.Pages -}}
{{ range .Aliases }}
{{  . | printf "%-35s" }}	{{ $p.RelPermalink -}}
{{ end -}}
{{- end -}}
```
The template above creates 301 redirects for your [aliases](/content-management/urls/#aliases), so you will probably want to turn off aliases in your `config.toml`: `disableAliases = true`.

