---
author: "Rick Cogley"
date: 2015-07-08
linktitle: Multilingual Site
menu:
  main:
    parent: tutorials
prev: /tutorials/migrate-from-jekyll
title: Create a Multilingual Site
weight: 10
---

## Introduction

Hugo allows you to create a multilingual site from its built-in tools. This tutorial will show one way to do it, and assumes:

* You already know the basics about creating a Hugo site
* You have a separate domain name for each language
* You'll use `/data` files for some translation strings
* You'll use single, combined `layout` and `static` folders
* You'll use a subfolder for each language under `content` and `public`

## Site Configs

Create your site configs in the root of your repository, for example for an English and Japanese site.

**English Config `config_en.toml`**:

~~~toml
baseurl = "http://acme.com/"
title = "Acme Inc."
contentdir = "content/en"
publishdir = "public/en"
...
[params]
    locale = "en-US"
~~~

**Japanese Config `config_ja.toml`**:

~~~toml
baseurl = "http://acme.jp/"
title = "有限会社アクミー"
contentdir = "content/ja"
publishdir = "public/ja"
...
[params]
    locale = "ja-JP"
~~~

If you had more domains and languages, you would just create more config files. The standard `config.toml` is what Hugo will run as a default, but since we're creating language-specific ones, you'll need to specify each config file when running `hugo server` or just `hugo` before deploying.

## Prep Translation Strings in `/data`

Create `.yaml` (or `.json` or `.toml`) files for each language, under `/data/translations`.

**English Strings `en-US.yaml`**:

~~~yaml
topslogan: Acme Inc.
topsubslogan: You'll love us
...
~~~

**Japanese Strings `ja-JP.yaml`**:

~~~yaml
topslogan: 有限会社アクミー
topsubslogan: キット勝つぞ
...
~~~

In some cases, where there is more complex formatting within the strings you want to show, it might be better to employ some conditional logic in your template, to display a block of html per language.

## Reference Strings in templates

Now you can reference the strings in your templates. One way is to do it like in this `layouts/index.html`, leveraging the fact that you have the locale set:

~~~html
<!DOCTYPE html>
<html lang="{{ .Site.Params.locale }}">
...
  <head>
    <meta charset="utf-8">
    <title>{{ if eq .Site.Params.locale "en-US" }}{{ if .IsHome }}Welcome to {{ end }}{{ end }}{{ .Title }}{{ if eq .Site.Params.locale "ja-JP" }}{{ if .IsHome }}へようこそ{{ end }}{{ end }}{{ if ne .Title .Site.Title }} : {{ .Site.Title }}{{ end }}</title>
    ...
  </head>
  <body>
    <div class="container">
      <h1 class="header">{{ ( index $.Site.Data.translations $.Site.Params.locale ).topslogan }}</h1>
      <h3 class="subheader">{{ ( index $.Site.Data.translations $.Site.Params.locale ).topsubslogan }}</h3>
    </div>
  </body>
</html>
~~~

The above shows both techniques, using an `if eq` and `else if eq` to check the locale, and using `index` to pull strings from the data file that matches the locale set in the site's config file.

## Customize Dates

At the time of this writing, Golang does not yet have support for internationalized locales, but if you do some work, you can simulate it. For example, if you want to use French month names, you can add a data file like ``data/mois.yaml`` with this content:

~~~toml
1: "janvier"
2: "février"
3: "mars"
4: "avril"
5: "mai"
6: "juin"
7: "juillet"
8: "août"
9: "septembre"
10: "octobre"
11: "novembre"
12: "décembre"
~~~

... then index the non-English date names in your templates like so:

~~~html
<time class="post-date" datetime="{{ .Date.Format "2006-01-02T15:04:05Z07:00" | safeHTML }}">
  Article publié le {{ .Date.Day }} {{ index $.Site.Data.mois (printf "%d" .Date.Month) }} {{ .Date.Year }} (dernière modification le {{ .Lastmod.Day }} {{ index $.Site.Data.mois (printf "%d" .Lastmod.Month) }} {{ .Lastmod.Year }})
</time>
~~~

This technique extracts the day, month and year by specifying ``.Date.Day``, ``.Date.Month``, and ``.Date.Year``, and uses the month number as a key, when indexing the month name data file.  

## Create Multilingual Content

Now you can create markdown content in your languages, in the `content/en` and `content/ja` folders. The frontmatter stays the same on the key side, but the values would be set in each of the languages.

## Run Hugo Server or Deploy Commands

Once you have things set up, you can run `hugo server` or `hugo` before deploying. You can create scripts to do it, or as shell functions. Here are sample basic `zsh` functions:

**Live Reload with `hugo server`**:

~~~shell
function hugoserver-com {
  cd /Users/me/dev/mainsite
  hugo server --buildDrafts --watch --verbose --source="/Users/me/dev/mainsite" --config="/Users/me/dev/mainsite/config_en.toml" --port=1377
}
function hugoserver-jp {
  cd /Users/me/dev/mainsite
  hugo server --buildDrafts --watch --verbose --source="/Users/me/dev/mainsite" --config="/Users/me/dev/mainsite/config_ja.toml" --port=1399
}
~~~

**Deploy with `hugo` and `rsync`**:

~~~shell
function hugodeploy-acmecom {
    rm -rf /tmp/acme.com
    hugo --config="/Users/me/dev/mainsite/config_en.toml" -s /Users/me/dev/mainsite/ -d /tmp/acme.com
    rsync -avze "ssh -p 22" --delete /tmp/acme.com/ me@mywebhost.com:/home/me/webapps/acme_com_site
}

function hugodeploy-acmejp {
    rm -rf /tmp/acme.jp
    hugo --config="/Users/me/dev/mainsite/config_ja.toml" -s /Users/me/dev/mainsite/ -d /tmp/acme.jp
    rsync -avze "ssh -p 22" --delete /tmp/acme.jp/ me@mywebhost.com:/home/me/webapps/acme_jp_site
}
~~~

Adjust to fit your situation, setting dns, your webserver config, and other settings as appropriate.
