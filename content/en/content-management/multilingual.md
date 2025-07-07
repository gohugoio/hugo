---
title: Multilingual mode
linkTitle: Multilingual
description: Localize your project for each language and region, including translations, images, dates, currencies, numbers, percentages, and collation sequence. Hugo's multilingual framework supports single-host and multihost configurations.
categories: []
keywords: []
aliases: [/content/multilingual/,/tutorials/create-a-multilingual-site/]
---

## Configuration

See [configure languages](/configuration/languages/).

## Translate your content

There are two ways to manage your content translations. Both ensure each page is assigned a language and is linked to its counterpart translations.

### Translation by file name

Considering the following example:

1. `/content/about.en.md`
1. `/content/about.fr.md`

The first file is assigned the English language and is linked to the second.
The second file is assigned the French language and is linked to the first.

Their language is __assigned__ according to the language code added as a __suffix to the file name__.

By having the same **path and base file name**, the content pieces are __linked__ together as translated pages.

> [!note]
> If a file has no language code, it will be assigned the default language.

### Translation by content directory

This system uses different content directories for each of the languages. Each language's `content` directory is set using the `contentDir` parameter.

{{< code-toggle file=hugo >}}
languages:
  en:
    weight: 10
    languageName: "English"
    contentDir: "content/english"
  fr:
    weight: 20
    languageName: "Français"
    contentDir: "content/french"
{{< /code-toggle >}}

The value of `contentDir` can be any valid path -- even absolute path references. The only restriction is that the content directories cannot overlap.

Considering the following example in conjunction with the configuration above:

1. `/content/english/about.md`
1. `/content/french/about.md`

The first file is assigned the English language and is linked to the second.
The second file is assigned the French language and is linked to the first.

Their language is __assigned__ according to the `content` directory they are __placed__ in.

By having the same **path and basename** (relative to their language `content` directory), the content pieces are __linked__ together as translated pages.

### Bypassing default linking

Any pages sharing the same `translationKey` set in front matter will be linked as translated pages regardless of basename or location.

Considering the following example:

1. `/content/about-us.en.md`
1. `/content/om.nn.md`
1. `/content/presentation/a-propos.fr.md`

{{< code-toggle file=hugo >}}
translationKey: "about"
{{< /code-toggle >}}

By setting the `translationKey` front matter parameter to `about` in all three pages, they will be __linked__ as translated pages.

### Localizing permalinks

Because paths and file names are used to handle linking, all translated pages will share the same URL (apart from the language subdirectory).

To localize URLs:

- For a regular page, set either [`slug`] or [`url`] in front matter
- For a section page, set [`url`] in front matter

For example, a French translation can have its own localized slug.

{{< code-toggle file=content/about.fr.md fm=true >}}
title: A Propos
slug: "a-propos"
{{< /code-toggle >}}

At render, Hugo will build both `/about/` and `/fr/a-propos/` without affecting the translation link.

### Page bundles

To avoid the burden of having to duplicate files, each Page Bundle inherits the resources of its linked translated pages' bundles except for the content files (Markdown files, HTML files etc.).

Therefore, from within a template, the page will have access to the files from all linked pages' bundles.

If, across the linked bundles, two or more files share the same basename, only one will be included and chosen as follows:

- File from current language bundle, if present.
- First file found across bundles by order of language `Weight`.

> [!note]
> Page Bundle resources follow the same language assignment logic as content files, both by file name (`image.jpg`, `image.fr.jpg`) and by directory (`english/about/header.jpg`, `french/about/header.jpg`).

## Reference translated content

To create a list of links to translated content, use a template similar to the following:

```go-html-template {file="layouts/_partials/i18nlist.html"}
{{ if .IsTranslated }}
<h4>{{ i18n "translations" }}</h4>
<ul>
  {{ range .Translations }}
  <li>
    <a href="{{ .RelPermalink }}">{{ .Language.Lang }}: {{ .LinkTitle }}{{ if .IsPage }} ({{ i18n "wordCount" . }}){{ end }}</a>
  </li>
  {{ end }}
</ul>
{{ end }}
```

The above can be put in a partial template then included in any template. It will not print anything if there are no translations for a given page.

The above also uses the [`i18n` function][i18func] described in the next section.

### List all available languages

`.AllTranslations` on a `Page` can be used to list all translations, including the page itself. On the home page it can be used to build a language navigator:

```go-html-template {file="layouts/_partials/allLanguages.html"}
<ul>
{{ range $.Site.Home.AllTranslations }}
<li><a href="{{ .RelPermalink }}">{{ .Language.LanguageName }}</a></li>
{{ end }}
</ul>
```

## Translation of strings

See the [`lang.Translate`] template function.

## Localization

The following localization examples assume your site's primary language is English, with translations to French and German.

{{< code-toggle file=hugo >}}
defaultContentLanguage = 'en'

[languages]
[languages.en]
contentDir = 'content/en'
languageName = 'English'
weight = 1
[languages.fr]
contentDir = 'content/fr'
languageName = 'Français'
weight = 2
[languages.de]
contentDir = 'content/de'
languageName = 'Deutsch'
weight = 3

{{< /code-toggle >}}

### Dates

With this front matter:

{{< code-toggle file=hugo >}}
date = 2021-11-03T12:34:56+01:00
{{< /code-toggle >}}

And this template code:

```go-html-template
{{ .Date | time.Format ":date_full" }}
```

The rendered page displays:

Language|Value
:--|:--
English|Wednesday, November 3, 2021
Français|mercredi 3 novembre 2021
Deutsch|Mittwoch, 3. November 2021

See [`time.Format`] for details.

### Currency

With this template code:

```go-html-template
{{ 512.5032 | lang.FormatCurrency 2 "USD" }}
```

The rendered page displays:

Language|Value
:--|:--
English|$512.50
Français|512,50 $US
Deutsch|512,50 $

See [lang.FormatCurrency] and [lang.FormatAccounting] for details.

### Numbers

With this template code:

```go-html-template
{{ 512.5032 | lang.FormatNumber 2 }}
```

The rendered page displays:

Language|Value
:--|:--
English|512.50
Français|512,50
Deutsch|512,50

See [lang.FormatNumber] and [lang.FormatNumberCustom] for details.

### Percentages

With this template code:

```go-html-template
{{ 512.5032 | lang.FormatPercent 2 }}
```

The rendered page displays:

Language|Value
:--|:--
English|512.50%
Français|512,50 %
Deutsch|512,50 %

See [lang.FormatPercent] for details.

## Menus

Localization of menu entries depends on how you define them:

- When you define menu entries [automatically] using the section pages menu, you must use translation tables to localize each entry.
- When you define menu entries [in front matter], they are already localized based on the front matter itself. If the front matter values are insufficient, use translation tables to localize each entry.
- When you define menu entries [in site configuration], you must create language-specific menu entries under each language key. If the names of the menu entries are insufficient, use translation tables to localize each entry.

### Create language-specific menu entries

#### Method 1 -- Use a single configuration file

For a simple menu with a small number of entries, use a single configuration file. For example:

{{< code-toggle file=hugo >}}
[languages.de]
languageCode = 'de-DE'
languageName = 'Deutsch'
weight = 1

[[languages.de.menus.main]]
name = 'Produkte'
pageRef = '/products'
weight = 10

[[languages.de.menus.main]]
name = 'Leistungen'
pageRef = '/services'
weight = 20

[languages.en]
languageCode = 'en-US'
languageName = 'English'
weight = 2

[[languages.en.menus.main]]
name = 'Products'
pageRef = '/products'
weight = 10

[[languages.en.menus.main]]
name = 'Services'
pageRef = '/services'
weight = 20
{{< /code-toggle >}}

#### Method 2 -- Use a configuration directory

With a more complex menu structure, create a [configuration directory] and split the menu entries into multiple files, one file per language. For example:

```text
config/
└── _default/
    ├── menus.de.toml
    ├── menus.en.toml
    └── hugo.toml
```

{{< code-toggle file=config/_default/menus.de >}}
[[main]]
name = 'Produkte'
pageRef = '/products'
weight = 10
[[main]]
name = 'Leistungen'
pageRef = '/services'
weight = 20
{{< /code-toggle >}}

{{< code-toggle file=config/_default/menus.en >}}
[[main]]
name = 'Products'
pageRef = '/products'
weight = 10
[[main]]
name = 'Services'
pageRef = '/services'
weight = 20
{{< /code-toggle >}}

### Use translation tables

When rendering the text that appears in menu each entry, the [example menu template] does this:

```go-html-template
{{ or (T .Identifier) .Name | safeHTML }}
```

It queries the translation table for the current language using the menu entry's `identifier` and returns the translated string. If the translation table does not exist, or if the `identifier` key is not present in the translation table, it falls back to `name`.

The `identifier` depends on how you define menu entries:

- If you define the menu entry [automatically] using the section pages menu, the `identifier` is the page's `.Section`.
- If you define the menu entry [in site configuration] or [in front matter], set the `identifier` property to the desired value.

For example, if you define menu entries in site configuration:

{{< code-toggle file=hugo >}}
[[menus.main]]
  identifier = 'products'
  name = 'Products'
  pageRef = '/products'
  weight = 10
[[menus.main]]
  identifier = 'services'
  name = 'Services'
  pageRef = '/services'
  weight = 20
{{< / code-toggle >}}

Create corresponding entries in the translation tables:

{{< code-toggle file=i18n/de >}}
products = 'Produkte'
services = 'Leistungen'
{{< / code-toggle >}}

## Missing translations

If a string does not have a translation for the current language, Hugo will use the value from the default language. If no default value is set, an empty string will be shown.

While translating a Hugo website, it can be handy to have a visual indicator of missing translations. The [`enableMissingTranslationPlaceholders` configuration option][config] will flag all untranslated strings with the placeholder `[i18n] identifier`, where `identifier` is the id of the missing translation.

> [!note]
> Hugo will generate your website with these missing translation placeholders. It might not be suitable for production environments.

For merging of content from other languages (i.e. missing content translations), see [lang.Merge].

To track down missing translation strings, run Hugo with the `--printI18nWarnings` flag:

```sh
hugo --printI18nWarnings | grep i18n
i18n|MISSING_TRANSLATION|en|wordCount
```

## Multilingual themes support

To support Multilingual mode in your themes, some considerations must be taken for the URLs in the templates. If there is more than one language, URLs must meet the following criteria:

- Come from the built-in `.Permalink` or `.RelPermalink`
- Be constructed with the [`relLangURL`] or [`absLangURL`] template function, or be prefixed with `{{ .LanguagePrefix }}`

If there is more than one language defined, the `LanguagePrefix` method will return `/en` (or whatever the current language is). If not enabled, it will be an empty string (and is therefore harmless for single-language Hugo websites).

## Generate multilingual content with `hugo new content`

If you organize content with translations in the same directory:

```sh
hugo new content post/test.en.md
hugo new content post/test.de.md
```

If you organize content with translations in different directories:

```sh
hugo new content content/en/post/test.md
hugo new content content/de/post/test.md
```

[`absLangURL`]: /functions/urls/abslangurl/
[`lang.Translate`]: /functions/lang/translate
[`relLangURL`]: /functions/urls/rellangurl/
[`slug`]: /content-management/urls/#slug
[`time.Format`]: /functions/time/format/
[`url`]: /content-management/urls/#url
[automatically]: /content-management/menus/#define-automatically
[config]: /configuration/
[configuration directory]: /configuration/introduction/#configuration-directory
[example menu template]: /templates/menu/#example
[i18func]: /functions/lang/translate/
[in front matter]: /content-management/menus/#define-in-front-matter
[in site configuration]: /content-management/menus/#define-in-site-configuration
[lang.FormatAccounting]: /functions/lang/formataccounting/
[lang.FormatCurrency]: /functions/lang/formatcurrency/
[lang.FormatNumber]: /functions/lang/formatnumber/
[lang.FormatNumberCustom]: /functions/lang/formatnumbercustom/
[lang.FormatPercent]: /functions/lang/formatpercent/
[lang.Merge]: /functions/lang/merge/
