---
date: 2016-01-02T21:21:00Z
menu:
  main:
    parent: content
next: /content/example
prev: /content/summaries
title: Multilingual Mode
weight: 68
toc: true
---

Since version 0.16, Hugo supports a native Multilingual mode. You can enable it with:

```
Multilingual: true
RenderLanguage: en
```

in your site configuration.

With the config above, all content, sitemap, RSS feeds, paginations
and taxonomy pages will be rendered under `/en`.

You will also need two config files (`config.en.yaml` and
`config.fr.yaml` for example) and you will need to run `hugo` twice,
to render each language's HTML.


### Translating your content

Translated articles are picked up by the name of the content files.

Example of translated articles:

1. `/content/about.en.md`
2. `/content/about.fr.md`

You can also have:

1. `/content/about.md`
2. `/content/about.fr.md`

in which case the config variable `DefaultContentLang` will be used to
affect the default language `about.md`.  This way, you can slowly
start to translate your current content without having to rename
everything.

If left unspecified, the value for `DefaultContentLang` defaults to
`en`.

By having the same _base file name_, the content pieces are linked
together as translated pieces. Only the content pieces in the language
defined by **.Site.RenderLanguage** will be rendered in a run of
`hugo`.  The translated content will be available in the
`.Page.Translations` so you can create links to the corresponding
translated pieces.


### Language switching links

A full example of the language-switching links would be:

```
{{if .Site.Multilingual}}
  {{if .IsPage}}
    {{ range $txLang := .Site.LinkLanguages }}
      {{if ne $lang $txLang}}
        {{if isset $.Translations $txLang}}
          <a href="{{ (index $.Translations $txLang).Permalink }}">{{ index $.Site.Data.i18n "langlinks" $txLang }}</a>
        {{else}}
          <a href="/{{$txLang}}">{{ index $.Site.Data.i18n "langlinks" $txLang }}</a>
        {{end}}
      {{end}}
    {{end}}
  {{end}}

  {{if .IsNode}}
    {{ range $txLang := .Site.LinkLanguages }}
      {{if ne $lang $txLang}}
        <a href="/{{$txLang}}">{{ index $.Site.Data.i18n "langlinks" $txLang }}</a>
      {{end}}
    {{end}}
  {{end}}
{{end}}
```

This makes use of the **.Site.LinkLanguages** variable to always
create links to the other available languages. You would define this
in your config as:

```
LinkLanguages:
 - fr
 - en
```

As you might notice, node pages link to the root of the other
available translations (`/en`), as those pages do not necessarily have
a translated counterpart.

Taxonomies (tags, categories) are completely segregated between
translations, will have their own tag clouds and list views.


### Multilingual Themes support

To support Multilingual mode in your themes, you only need to make
sure URLs defined manually (those not using `.Permalink` or `.URL`
variables) in your templates are prefixed with `{{
.Site.LanguagePrefix }}`. If `Multilingual` mode is enabled, this will
be rendered as `/en` (based on `RenderLanguage`). If not enabled, it
will be an empty string, so it is harmless for non-multilingual sites.
