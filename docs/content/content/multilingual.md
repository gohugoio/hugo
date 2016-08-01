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

Since version 0.17, Hugo supports a native Multilingual mode. In your
top-level `config.yaml` (or equivalent), you define the available
languages in a `Multilingual` section such as:

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

Anything not defined in a `[lang]:` block will fall back to the global
value for that key (like `copyright` for the `en` lang in this
example).

With the config above, all content, sitemap, RSS feeds, paginations
and taxonomy pages will be rendered under `/en` in English, and under
`/fr` in French.

Only those keys are read under `Multilingual`: `weight`, `title`,
`author`, `social`, `languageCode`, `copyright`, `disqusShortname`,
`params` (which can contain a map of several other keys).


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

Here is a simple example if all your pages are translated:

```
{{if .IsPage}}
  {{ range $txLang := .Site.Languages }}
    {{if isset $.Translations $txLang}}
      <a href="{{ (index $.Translations $txLang).Permalink }}">{{ i18n ( printf "language_switcher_%s" $txLang ) }}</a>
    {{end}}
  {{end}}
{{end}}

{{if .IsNode}}
  {{ range $txLang := .Site.Languages }}
    <a href="/{{$txLang}}">{{ i18n ( printf "language_switcher_%s" $txLang ) }}</a>
  {{end}}
{{end}}
```

This is a more complete example. It handles missing translations and will support non-multilingual sites. Better for theme authors:

```
{{if .Site.Multilingual}}
  {{if .IsPage}}
    {{ range $txLang := .Site.Languages }}
      {{if isset $.Translations $txLang}}
        <a href="{{ (index $.Translations $txLang).Permalink }}">{{ i18n ( printf "language_switcher_%s" $txLang ) }}</a>
      {{else}}
        <a href="/{{$txLang}}">{{ i18n ( printf "language_switcher_%s" $txLang ) }}</a>
      {{end}}
    {{end}}
  {{end}}

  {{if .IsNode}}
    {{ range $txLang := .Site.Languages }}
      <a href="/{{$txLang}}">{{ i18n ( printf "language_switcher_%s" $txLang ) }}</a>
    {{end}}
  {{end}}
{{end}}
```

This makes use of the **.Site.Languages** variable to create links to
the other available languages.  The order in which the languages are
listed is defined by the `weight` attribute in each language under
`Multilingual`.

This will also require you to have some content in your `i18n/` files
(see below) that would look like:

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
translations and will have their own tag clouds and list views.


### Translation of strings

Hugo uses [go-i18n](https://github.com/nicksnyder/go-i18n) to support
string translations.  Follow the link to find tools to manage your
translation workflows.

Translations are collected from the `themes/[name]/i18n/` folder
(built into the theme), as well as translations present in `i18n/` at
the root of your project.  In the `i18n`, the translations will be
merged and take precedence over what is in the theme folder.  Files in
there follow RFC 5646 and should be named something like `en-US.yaml`,
`fr.yaml`, etc..

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
`LanguagePrefix` variable will equal `"/en"` (or whatever your
`CurrentLanguage` is). If not enabled, it will be an empty string, so
it is harmless for non-multilingual sites.


### Multilingual index.html and 404.html

To redirect your users to their closest language, drop an `index.html`
in `/static` of your site, with the following content (tailored to
your needs) to redirect based on their browser's language:

```
<html><head>
<meta http-equiv="refresh" content="1;url=/en" /><!-- just in case JS doesn't work -->
<script>
lang = window.navigator.language.substr(0, 2);
if (lang == "fr") {
    window.location = "/fr";
} else {
    window.location = "/en";
}

/* or simply:
window.location = "/en";
*/
</script></head><body></body></html>
```

An even simpler version will always redirect your users to a given language:

```
<html><head>
<meta http-equiv="refresh" content="0;url=/en" />
</head><body></body></html>
```

You can do something similar with your `404.html` page, as you don't
know the language of someone arriving at a non-existing page.  You
could inspect the prefix of the navigator path in Javascript or use
the browser's language detection like above.


### Sitemaps

As sitemaps are generated once per language and live in
`[lang]/sitemap.xml`. Write this content in `static/sitemap.xml` to
link all your sitemaps together:

```
<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
   <sitemap>
      <loc>https://example.com/en/sitemap.xml</loc>
   </sitemap>
   <sitemap>
      <loc>https://example.com/fr/sitemap.xml</loc>
   </sitemap>
</sitemapindex>
```

and explicitly list all the languages you want referenced.
