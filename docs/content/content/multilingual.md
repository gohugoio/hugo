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
Hugo supports multiple languages side-by-side (added in `Hugo 0.17`). Define the available languages in a `Languages` section in your top-level `config.yaml` (or equivalent).

Example:

```
DefaultContentLanguage = "en"

Languages:
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
value for that key (like `copyright` for the English (`en`) language in this example).

With the config above, all content, sitemap, RSS feeds, paginations
and taxonomy pages will be rendered below `/` in English (your default content language), and below `/fr` in French.

If you want all of the languages to be put below their respective language code, enable `DefaultContentLanguageInSubdir: true` in your configuration.

Only the obvious non-global options can be overridden per language. Examples of global options are `BaseURL`, `BuildDrafts`, etc.

Taxonomies and Blackfriday configuration can also be set per language, example:

```
[Taxonomies]
tag = "tags"

[blackfriday]
angledQuotes = true
hrefTargetBlank = true

[Languages]
[Languages.en]
weight = 1
title = "English"
[Languages.en.blackfriday]
angledQuotes = false

[Languages.fr]
weight = 2
title = "Français"
[Languages.fr.Taxonomies]
plaque = "plaques"
```


### Translating your content

Translated articles are identified by the name of the content file.

Example of translated articles:

1. `/content/about.en.md`
2. `/content/about.fr.md`

You can also have:

1. `/content/about.md`
2. `/content/about.fr.md`

In which case the config variable `DefaultContentLanguage` will be used to affect the default language `about.md`.  This way, you can
slowly start to translate your current content without having to rename everything.

If left unspecified, the value for `DefaultContentLanguage` defaults to `en`.

By having the same _base file name_, the content pieces are linked together as translated pieces. 

### Link to translated content

To create a list of links to translated content, use a template similar to this:

```
{{ if .IsTranslated }}
<h4>{{ i18n "translations" }}</h4>
<ul>
    {{ range .Translations }}
    <li>
        <a href="{{ .Permalink }}">{{ .Lang }}: {{ .Title }}{{ if .IsPage }} ({{ i18n "wordCount" . }}){{ end }}</a>
    </li>
    {{ end}}
</ul>
{{ end }}
```
The above can be put in a `partial` and included in any template, be it for a content page or the home page.  It will not print anything if there are no translations for a given page, or if it is -- in the case of the home page, section listing etc. -- a site with only one language.

The above also uses the `i8n` func, see [Translation of strings](#translation-of-strings).

### Translation of strings

Hugo uses [go-i18n](https://github.com/nicksnyder/go-i18n) to support string translations.  Follow the link to find tools to manage your translation workflows.

Translations are collected from the `themes/[name]/i18n/` folder (built into the theme), as well as translations present in `i18n/` at the root of your project.  In the `i18n`, the translations will be merged and take precedence over what is in the theme folder.  Language files should be named according to RFC 5646  with names such as `en-US.yaml`, `fr.yaml`, etc.

From within your templates, use the `i18n` function like this:

```
{{ i18n "home" }}
```
This uses a definition like this one in `i18n/en-US.yaml`:
```
- id: home
  translation: "Home"
```

Often you will want to use to the page variables in the translations strings. To do that, pass on the "." context when calling `18n`:

```
{{ i18n "wordCount" . }}
```
This uses a definition like this one in `i18n/en-US.yaml`:
```
- id: wordCount
  translation: "This article has {{ .WordCount }} words."
```
To track down missing translation strings, run Hugo with the `--i18n-warnings` flag:

```bash
 hugo --i18n-warnings | grep i18n
i18n|MISSING_TRANSLATION|en|wordCount
```

### Multilingual Themes support

To support Multilingual mode in your themes, some considerations must be taken for the URLs in the templates. If there are more than one language, URLs  must either  come from the built-in `.Permalink` or `.URL`, be constructed with `relLangURL` or `absLangURL` template funcs -- or prefixed with `{{.LanguagePrefix }}`.

If there are more than one language defined, the`LanguagePrefix` variable will equal `"/en"` (or whatever your `CurrentLanguage` is). If not enabled, it will be an empty string, so it is harmless for single-language sites.


