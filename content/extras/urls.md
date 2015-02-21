---
aliases:
- /doc/urls/
date: 2014-01-03
menu:
  main:
    parent: extras
next: /community/mailing-list
notoc: true
prev: /extras/toc
title: URLs
weight: 110
---

## Pretty URLs

By default, Hugo create content with 'pretty' URLs. For example,
content created at `/content/extras/urls.md` will be rendered at
`/public/extras/urls/index.html`, thus accessible from the browser
at http://example.com/extras/urls/.  No non-standard server-side
configuration is required for these pretty URLs to work.

If you would like to have what we call "ugly URLs",
e.g.&nbsp;http://example.com/extras/urls.html, you are in luck.
Hugo supports the ability to create your entire site with ugly URLs.
Simply add `uglyurls = true` to your site-wide `config.toml`,
or use the `--uglyUrls=true` flag on the command line.

If you want a specific piece of content to have an exact URL, you can
specify this in the front matter under the `url` key. See [Content
Organization](/content/organization/) for more details. 

## Canonicalization

By default, all relative URLs encountered in the input are left unmodified,
e.g. `/css/foo.css` would stay as `/css/foo.css`,
i.e. `canonifyurls` defaults to `false`.

By setting `canonifyurls` to `true`, all relative URLs would instead
be *canonicalized* using `baseurl`.  For example, assuming you have
`baseurl = http://yoursite.example.com/` defined in the site-wide
`config.toml`, the relative URL `/css/foo.css` would be turned into
the absolute URL `http://yoursite.example.com/css/foo.css`.

Benefits of canonicalization include fixing all URLs to be absolute, which may
aid with some parsing tasks.  Note though that all real browsers handle this
client-side without issues.

Benefits of non-canonicalization include being able to have resource inclusion
be scheme-relative, so that http vs https can be decided based on how this
page was retrieved.

> Note: In the May 2014 release of Hugo v0.11, the default value of `canonifyurls` was switched from `true` to `false`, which we think is the better default and should continue to be the case going forward. So, please verify and adjust your website accordingly if you are upgrading from v0.10 or older versions.

To find out the current value of `canonifyurls` for your website, you may use the handy `hugo config` command added in v0.13:

    hugo config | grep -i canon

Or, if you are on Windows and do not have `grep` installed:

    hugo config | FINDSTR /I canon

