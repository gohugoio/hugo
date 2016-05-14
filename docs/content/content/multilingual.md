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

Since version 0.16, Hugo supports a native Multilingual mode. You
define the languages to render as such:

```
Multilingual:
  en:
    weight: 1
    title: "My blog"
    params:
      linkedin: "english-link"
  fr:
    weight: 2

    title: "Mon blog"
    params:
      linkedin: "lien-francais"
    copyright: "Tout est miens"

copyright: "Everything is mine"
```

Anything not defined in a `lang:` block will fall back to the global
value for that key (like `copyright` in this example).

With the config above, all content, sitemap, RSS feeds, paginations
and taxonomy pages will be rendered under `/en` in English, and under
`/fr` in French.


### Translating your content

Translated articles are picked up by the name of the content files.

Example of translated articles:

1. `/content/about.en.md`
2. `/content/about.fr.md`

You can also have:

1. `/content/about.md`
2. `/content/about.fr.md`

in which case the config variable `DefaultContentLanguage` will be
used to affect the default language `about.md`.  This way, you can
slowly start to translate your current content without having to
rename everything.

If left unspecified, the value for `DefaultContentLanguage` defaults
to `en`.

By having the same _base file name_, the content pieces are linked
together as translated pieces. Only the content pieces in the language
defined by **.Site.CurrentLanguage** will be rendered in a run of
`hugo`.  The translated content will be available in the
`.Page.Translations` so you can create links to the corresponding
translated pieces.


### Language switching links

A full example of the language-switching links would be:

```
        {{if .Site.Multilingual}}
          {{if .IsPage}}
            {{ range $txLang := .Site.Languages }}
              {{if ne $lang $txLang}}
                {{if isset $.Translations $txLang}}
                  <a href="{{ (index $.Translations $txLang).Permalink }}">{{ i18n ( printf "language_switcher_%s" $txLang ) }}</a>
                {{else}}
                  <a href="/{{$txLang}}">{{ i18n ( printf "language_switcher_%s" $txLang ) }}</a>
                {{end}}
              {{end}}
            {{end}}
          {{end}}

          {{if .IsNode}}
            {{ range $txLang := .Site.Languages }}
              {{if ne $lang $txLang}}
                <a href="/{{$txLang}}">{{ i18n ( printf "language_switcher_%s" $txLang ) }}</a>
              {{end}}
            {{end}}
          {{end}}
        {{end}}
```

This makes use of the **.Site.Languages** variable to create links to
the other available languages.  The order in which the languages are
listed is defined by the `weight` attribute in each language under
`Multilingual`.

This will also require you to have some content in your `i18n/` files that would look like:

```
- id: language_switcher_en
  translation: "English"
- id: language_switcher_fr
  translation: "Français"
```

and a copy of this in translations for each language.

As you might notice, node pages link to the root of the other
available translations (`/en`), as those pages do not necessarily have
a translated counterpart.

Taxonomies (tags, categories) are completely segregated between
translations, will have their own tag clouds and list views.


### Translations of strings

Hugo uses [go-i18n](https://github.com/nicksnyder/go-i18n) to support
string translations.  Follow the link to find tools to manage your
translation workflows.

Translations are collected from the `themes/[name]/i18n/` folder, in
addition to the files present in `i18n/` at the root of your project.

From within your templates, use the `i18n` function as such:

```
{{ i18n "home" }}
```

to use a definition like this one in `i18n/en-US.yaml`:

```
- id: home
  translation: "Home"
```


### Multilingual Themes support

To support Multilingual mode in your themes, you only need to make
sure URLs defined manually (those not using `.Permalink` or `.URL`
variables) in your templates are prefixed with `{{
.Site.LanguagePrefix }}`. If `Multilingual` mode is enabled, the
`LanguagePrefix` variable will be rendered as `/en` (based on
`CurrentLanguage`). If not enabled, it will be an empty string, so it
is harmless for non-multilingual sites.
