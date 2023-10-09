---
title: Multilingual mode
linkTitle: Multilingual
description: Hugo supports the creation of websites with multiple languages side by side.
categories: [content management]
keywords: [multilingual,i18n, internationalization]
menu:
  docs:
    parent: content-management
    weight: 230
toc: true
weight: 230
aliases: [/content/multilingual/,/tutorials/create-a-multilingual-site/]
---

You should define the available languages in a `languages` section in your site configuration.

Also See [Hugo Multilingual Part 1: Content translation].

## Configure languages

This is the default language configuration:

{{< code-toggle config="languages" />}}

This is an example of a site configuration for a multilingual project. Any key not defined in a `languages` object will fall back to the global value in the root of your site configuration.

{{< code-toggle file="hugo" >}}
defaultContentLanguage = 'de'
defaultContentLanguageInSubdir = true

[languages.de]
contentDir = 'content/de'
disabled = false
languageCode = 'de-DE'
languageDirection = 'ltr'
languageName = 'Deutsch'
title = 'Projekt Dokumentation'
weight = 1

[languages.de.params]
subtitle = 'Referenz, Tutorials und Erklärungen'

[languages.en]
contentDir = 'content/en'
disabled = false
languageCode = 'en-US'
languageDirection = 'ltr'
languageName = 'English'
title = 'Project Documentation'
weight = 2

[languages.en.params]
subtitle = 'Reference, Tutorials, and Explanations'
{{< /code-toggle >}}

`defaultContentLanguage`
: (`string`) The project's default language tag as defined by [RFC 5646]. Must be lower case, and must match one of the defined language keys. Default is `en`. Examples:

- `en`
- `en-gb`
- `pt-br`

`defaultContentLanguageInSubdir`
: (`bool`)  If `true`, Hugo renders the default language site in a subdirectory matching the `defaultContentLanguage`. Default is `false`.

`contentDir`
: (`string`) The content directory for this language. Omit if [translating by file name].

`disabled`
: (`bool`) If `true`, Hugo will not render content for this language. Default is `false`.

`languageCode`
: (`string`) The language tag as defined by [RFC 5646]. This value may include upper and lower case characters, hyphens or underscores, and does not affect localization or URLs. Hugo uses this value to populate the `language` element in the [built-in RSS template], and the `lang` attribute of the `html` element in the [built-in alias template]. Examples:

- `en`
- `en-GB`
- `pt-BR`

`languageDirection`
: (`string`) The language direction, either left-to-right (`ltr`) or right-to-left (`rtl`). Use this value in your templates with the global [`dir`] HTML attribute.

`languageName`
: (`string`) The language name, typically used when rendering a language switcher.

`title`
: (`string`) The language title. When set, this overrides the site title for this language.

`weight`
: (`int`) The language weight. When set to a non-zero value, this is the primary sort criteria for this language.

[`dir`]: https://developer.mozilla.org/en-US/docs/Web/HTML/Global_attributes/dir
[built-in RSS template]: https://github.com/gohugoio/hugo/blob/master/tpl/tplimpl/embedded/templates/_default/rss.xml
[built-in alias template]: https://github.com/gohugoio/hugo/blob/master/tpl/tplimpl/embedded/templates/alias.html
[RFC 5646]: https://datatracker.ietf.org/doc/html/rfc5646
[translating by file name]: #translation-by-file-name

### Changes in Hugo 0.112.0

{{< new-in "0.112.0" >}}

In Hugo `v0.112.0` we consolidated all configuration options, and improved how the languages and their parameters are merged with the main configuration. But while testing this on Hugo sites out there, we received some error reports and reverted some of the changes in favor of deprecation warnings:

1. `site.Language.Params` is deprecated. Use `site.Params` directly.
1. Adding custom parameters to the top level language configuration is deprecated. Define custom parameters within `languages.xx.params`. See `color` in the example below.

{{< code-toggle file=hugo copy=false >}}

title = "My blog"
languageCode = "en-us"

[languages]
[languages.sv]
title = "Min blogg"
languageCode = "sv"
[languages.en.params]
color = "blue"
{{< /code-toggle >}}

In the example above, all settings except `color` below `params` map to predefined configuration options in Hugo for the site and its language, and should be accessed via the documented accessors:

```go-html-template
{{ site.Title }}
{{ site.LanguageCode }}
{{ site.Params.color }}
```

### Disable a language

To disable a language within a `languages` object in your site configuration:

{{< code-toggle file="hugo" copy=false >}}
[languages.es]
disabled = true
{{< /code-toggle >}}

To disable one or more languages in the root of your site configuration:

{{< code-toggle file="hugo" copy=false >}}
disableLanguages = ["es", "fr"]
{{< /code-toggle >}}

To disable one or more languages using an environment variable:

```bash
HUGO_DISABLELANGUAGES="es fr" hugo
```

Note that you cannot disable the default content language.

### Configure multilingual multihost

From **Hugo 0.31** we support multiple languages in a multihost configuration. See [this issue](https://github.com/gohugoio/hugo/issues/4027) for details.

This means that you can now configure a `baseURL` per `language`:

{{% note %}}
If a `baseURL` is set on the `language` level, then all languages must have one and they must all be different.
{{% /note %}}

Example:

{{< code-toggle file="hugo" >}}
[languages]
[languages.fr]
baseURL = "https://example.fr"
languageName = "Français"
weight = 1
title = "En Français"

[languages.en]
baseURL = "https://example.com"
languageName = "English"
weight = 2
title = "In English"
{{</ code-toggle >}}

With the above, the two sites will be generated into `public` with their own root:

```text
public
├── en
└── fr
```

**All URLs (i.e `.Permalink` etc.) will be generated from that root. So the English home page above will have its `.Permalink` set to `https://example.com/`.**

When you run `hugo server` we will start multiple HTTP servers. You will typically see something like this in the console:

```text
Web Server is available at 127.0.0.1:1313 (bind address 127.0.0.1)
Web Server is available at 127.0.0.1:1314 (bind address 127.0.0.1)
Press Ctrl+C to stop
```

Live reload and `--navigateToChanged` between the servers work as expected.

## Translate your content

There are two ways to manage your content translations. Both ensure each page is assigned a language and is linked to its counterpart translations.

### Translation by file name

Considering the following example:

1. `/content/about.en.md`
2. `/content/about.fr.md`

The first file is assigned the English language and is linked to the second.
The second file is assigned the French language and is linked to the first.

Their language is __assigned__ according to the language code added as a __suffix to the file name__.

By having the same **path and base file name**, the content pieces are __linked__ together as translated pages.

{{% note %}}
If a file has no language code, it will be assigned the default language.
{{% /note %}}

### Translation by content directory

This system uses different content directories for each of the languages. Each language's content directory is set using the `contentDir` parameter.

{{< code-toggle file="hugo" >}}
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
2. `/content/french/about.md`

The first file is assigned the English language and is linked to the second.
The second file is assigned the French language and is linked to the first.

Their language is __assigned__ according to the content directory they are __placed__ in.

By having the same **path and basename** (relative to their language content directory), the content pieces are __linked__ together as translated pages.

### Bypassing default linking

Any pages sharing the same `translationKey` set in front matter will be linked as translated pages regardless of basename or location.

Considering the following example:

1. `/content/about-us.en.md`
2. `/content/om.nn.md`
3. `/content/presentation/a-propos.fr.md`

{{< code-toggle >}}
translationKey: "about"
{{< /code-toggle >}}

By setting the `translationKey` front matter parameter to `about` in all three pages, they will be __linked__ as translated pages.

### Localizing permalinks

Because paths and file names are used to handle linking, all translated pages will share the same URL (apart from the language subdirectory).

To localize URLs:

- For a regular page, set either [`slug`] or [`url`] in front matter
- For a section page, set [`url`] in front matter

[`slug`]: /content-management/urls/#slug
[`url`]: /content-management/urls/#url

For example, a French translation can have its own localized slug.

{{< code-toggle file="content/about.fr.md" fm=true copy=false >}}
title: A Propos
slug: "a-propos"
{{< /code-toggle >}}

At render, Hugo will build both `/about/` and `/fr/a-propos/` without affecting the translation link.

### Page bundles

To avoid the burden of having to duplicate files, each Page Bundle inherits the resources of its linked translated pages' bundles except for the content files (Markdown files, HTML files etc...).

Therefore, from within a template, the page will have access to the files from all linked pages' bundles.

If, across the linked bundles, two or more files share the same basename, only one will be included and chosen as follows:

* File from current language bundle, if present.
* First file found across bundles by order of language `Weight`.

{{% note %}}
Page Bundle resources follow the same language assignment logic as content files, both by file name (`image.jpg`, `image.fr.jpg`) and by directory (`english/about/header.jpg`, `french/about/header.jpg`).
{{%/ note %}}

## Reference translated content

To create a list of links to translated content, use a template similar to the following:

{{< code file="layouts/partials/i18nlist.html" >}}
{{ if .IsTranslated }}
<h4>{{ i18n "translations" }}</h4>
<ul>
  {{ range .Translations }}
  <li>
    <a href="{{ .Permalink }}">{{ .Lang }}: {{ .Title }}{{ if .IsPage }} ({{ i18n "wordCount" . }}){{ end }}</a>
  </li>
  {{ end }}
</ul>
{{ end }}
{{< /code >}}

The above can be put in a `partial` (i.e., inside `layouts/partials/`) and included in any template, whether a [single content page][contenttemplate] or the [homepage]. It will not print anything if there are no translations for a given page.

The above also uses the [`i18n` function][i18func] described in the next section.

### List all available languages

`.AllTranslations` on a `Page` can be used to list all translations, including the page itself. On the home page it can be used to build a language navigator:

{{< code file="layouts/partials/allLanguages.html" >}}
<ul>
{{ range $.Site.Home.AllTranslations }}
<li><a href="{{ .Permalink }}">{{ .Language.LanguageName }}</a></li>
{{ end }}
</ul>
{{< /code >}}

## Translation of strings

Hugo uses [go-i18n] to support string translations. [See the project's source repository][go-i18n-source] to find tools that will help you manage your translation workflows.

Translations are collected from the `themes/<THEME>/i18n/` folder (built into the theme), as well as translations present in `i18n/` at the root of your project. In the `i18n`, the translations will be merged and take precedence over what is in the theme folder. Language files should be named according to [RFC 5646] with names such as `en-US.toml`, `fr.toml`, etc.

Artificial languages with private use subtags as defined in [RFC 5646 &#167; 2.2.7](https://datatracker.ietf.org/doc/html/rfc5646#section-2.2.7) are also supported. You may omit the `art-x-` prefix for brevity. For example:

```text
art-x-hugolang
hugolang
```

Private use subtags must not exceed 8 alphanumeric characters.

### Query basic translation

From within your templates, use the `i18n` function like this:

```go-html-template
{{ i18n "home" }}
```

The function will search for the `"home"` id:

{{< code-toggle file="i18n/en-US" >}}
[home]
other = "Home"
{{< /code-toggle >}}

The result will be

```text
Home
```

### Query a flexible translation with variables

Often you will want to use the page variables in the translation strings. To do so, pass the `.` context when calling `i18n`:

```go-html-template
{{ i18n "wordCount" . }}
```

The function will pass the `.` context to the `"wordCount"` id:

{{< code-toggle file="i18n/en-US" >}}
[wordCount]
other = "This article has {{ .WordCount }} words."
{{< /code-toggle >}}

Assume `.WordCount` in the context has value is 101. The result will be:

```text
This article has 101 words.
```

### Query a singular/plural translation

To enable pluralization when translating, pass a map with a numeric `.Count` property to the `i18n` function. The example below uses `.ReadingTime` variable which has a built-in `.Count` property.

```go-html-template
{{ i18n "readingTime" .ReadingTime }}
```

The function will read `.Count` from `.ReadingTime` and evaluate whether the number is singular (`one`) or plural (`other`). After that, it will pass to `readingTime` id in `i18n/en-US.toml` file:

{{< code-toggle file="i18n/en-US" >}}
[readingTime]
one = "One minute to read"
other = "{{ .Count }} minutes to read"
{{< /code-toggle >}}

Assuming `.ReadingTime.Count` in the context has value is 525600. The result will be:

```text
525600 minutes to read
```

If `.ReadingTime.Count` in the context has value is 1. The result is:

```text
One minute to read
```

In case you need to pass a custom data: (`(dict "Count" numeric_value_only)` is minimum requirement)

```go-html-template
{{ i18n "readingTime" (dict "Count" 25 "FirstArgument" true "SecondArgument" false "Etc" "so on, so far") }}
```

## Localization

The following localization examples assume your site's primary language is English, with translations to French and German.

{{< code-toggle file="hugo" >}}
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

{{< code-toggle >}}
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

See [time.Format] for details.

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

Localization of menu entries depends on the how you define them:

- When you define menu entries [automatically] using the section pages menu, you must use translation tables to localize each entry.
- When you define menu entries [in front matter], they are already localized based on the front matter itself. If the front matter values are insufficient, use translation tables to localize each entry.
- When you define menu entries [in site configuration], you can (a) use translation tables, or (b) create language-specific menu entries under each language key.

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

{{< code-toggle file="hugo" copy=false >}}
[[menu.main]]
  identifier = 'products'
  name = 'Products'
  pageRef = '/products'
  weight = 10
[[menu.main]]
  identifier = 'services'
  name = 'Services'
  pageRef = '/services'
  weight = 20
{{< / code-toggle >}}

Create corresponding entries in the translation tables:

{{< code-toggle file="i18n/de" copy=false >}}
products = 'Produkte'
services = 'Leistungen'
{{< / code-toggle >}}

[example menu template]: /templates/menu-templates/#example
[automatically]: /content-management/menus/#define-automatically
[in front matter]: /content-management/menus/#define-in-front-matter
[in site configuration]: /content-management/menus/#define-in-site-configuration

### Create language-specific menu entries

For example:

{{< code-toggle file="hugo" copy=false >}}
[languages.de]
languageCode = 'de-DE'
languageName = 'Deutsch'
weight = 1

[[languages.de.menu.main]]
name = 'Produkte'
pageRef = '/products'
weight = 10

[[languages.de.menu.main]]
name = 'Leistungen'
pageRef = '/services'
weight = 20

[languages.en]
languageCode = 'en-US'
languageName = 'English'
weight = 2

[[languages.en.menu.main]]
name = 'Products'
pageRef = '/products'
weight = 10

[[languages.en.menu.main]]
name = 'Services'
pageRef = '/services'
weight = 20
{{< /code-toggle >}}

For a simple menu with two languages, these menu entries are easy to create and maintain. For a larger menu, or with more than two languages, using translation tables as described above is preferable.

## Missing translations

If a string does not have a translation for the current language, Hugo will use the value from the default language. If no default value is set, an empty string will be shown.

While translating a Hugo website, it can be handy to have a visual indicator of missing translations. The [`enableMissingTranslationPlaceholders` configuration option][config] will flag all untranslated strings with the placeholder `[i18n] identifier`, where `identifier` is the id of the missing translation.

{{% note %}}
Hugo will generate your website with these missing translation placeholders. It might not be suitable for production environments.
{{% /note %}}

For merging of content from other languages (i.e. missing content translations), see [lang.Merge].

To track down missing translation strings, run Hugo with the `--printI18nWarnings` flag:

```bash
hugo --printI18nWarnings | grep i18n
i18n|MISSING_TRANSLATION|en|wordCount
```

## Multilingual themes support

To support Multilingual mode in your themes, some considerations must be taken for the URLs in the templates. If there is more than one language, URLs must meet the following criteria:

* Come from the built-in `.Permalink` or `.RelPermalink`
* Be constructed with the [`relLangURL` template function][rellangurl] or the [`absLangURL` template function][abslangurl] **OR** be prefixed with `{{ .LanguagePrefix }}`

If there is more than one language defined, the `LanguagePrefix` variable will equal `/en` (or whatever your `CurrentLanguage` is). If not enabled, it will be an empty string (and is therefore harmless for single-language Hugo websites).


## Generate multilingual content with `hugo new content`

If you organize content with translations in the same directory:

```text
hugo new content post/test.en.md
hugo new content post/test.de.md
```

If you organize content with translations in different directories:

```text
hugo new content content/en/post/test.md
hugo new content content/de/post/test.md
```

[abslangurl]: /functions/abslangurl
[config]: /getting-started/configuration/
[contenttemplate]: /templates/single-page-templates/
[go-i18n-source]: https://github.com/nicksnyder/go-i18n
[go-i18n]: https://github.com/nicksnyder/go-i18n
[homepage]: /templates/homepage/
[Hugo Multilingual Part 1: Content translation]: https://regisphilibert.com/blog/2018/08/hugo-multilingual-part-1-managing-content-translation/
[i18func]: /functions/i18n/
[lang.FormatAccounting]: /functions/lang
[lang.FormatCurrency]: /functions/lang
[lang.FormatNumber]: /functions/lang
[lang.FormatNumberCustom]: /functions/lang
[lang.FormatPercent]: /functions/lang
[lang.Merge]: /functions/lang.merge/
[menus]: /content-management/menus/
[OS environment]: /getting-started/configuration/#configure-with-environment-variables
[rellangurl]: /functions/rellangurl
[RFC 5646]: https://tools.ietf.org/html/rfc5646
[single page templates]: /templates/single-page-templates/
[time.Format]: /functions/dateformat
