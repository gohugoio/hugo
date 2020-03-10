---
title: Multilingual Mode
linktitle: Multilingual and i18n
description: Hugo supports the creation of websites with multiple languages side by side.
date: 2017-01-10
publishdate: 2017-01-10
lastmod: 2017-01-10
categories: [content management]
keywords: [multilingual,i18n, internationalization]
menu:
  docs:
    parent: "content-management"
    weight: 150
weight: 150	#rem
draft: false
aliases: [/content/multilingual/,/tutorials/create-a-multilingual-site/]
toc: true
---

You should define the available languages in a `languages` section in your site configuration.

> Also See [Hugo Multilingual Part 1: Content translation](https://regisphilibert.com/blog/2018/08/hugo-multilingual-part-1-managing-content-translation/)

## Configure Languages

The following is an example of a site configuration for a multilingual Hugo project:

{{< code-toggle file="config" >}}
DefaultContentLanguage = "en"
copyright = "Everything is mine"

[params]
[params.navigation]
help  = "Help"

[languages]
[languages.en]
title = "My blog"
weight = 1
[languages.en.params]
linkedin = "https://linkedin.com/whoever"

[languages.fr]
title = "Mon blogue"
weight = 2
[languages.fr.params]
linkedin = "https://linkedin.com/fr/whoever"
[languages.fr.params.navigation]
help  = "Aide"

[languages.ar]
title = "مدونتي"
weight = 2
languagedirection = "rtl"
{{< /code-toggle >}}

Anything not defined in a `languages` block will fall back to the global value for that key (e.g., `copyright` for the English `en` language). This also works for `params`, as demonstrated with `help` above: You will get the value `Aide` in French and `Help` in all the languages without this parameter set.

With the configuration above, all content, sitemap, RSS feeds, paginations,
and taxonomy pages will be rendered below `/` in English (your default content language) and then below `/fr` in French.

When working with front matter `Params` in [single page templates][singles], omit the `params` in the key for the translation.

`defaultContentLanguage` sets the project's default language. If not set, the default language will be `en`.

If the default language needs to be rendererd below its own language code (`/en`) like the others, set `defaultContentLanguageInSubdir: true`.

Only the obvious non-global options can be overridden per language. Examples of global options are `baseURL`, `buildDrafts`, etc.

### Disable a Language

You can disable one or more languages. This can be useful when working on a new translation.

```toml
disableLanguages = ["fr", "ja"]
```

Note that you cannot disable the default content language.

We kept this as a standalone setting to make it easier to set via [OS environment](/getting-started/configuration/#configure-with-environment-variables):

```bash
HUGO_DISABLELANGUAGES="fr ja" hugo
```
If you have already a list of disabled languages in `config.toml`, you can enable them in development like this:

```bash
HUGO_DISABLELANGUAGES=" " hugo server
```


### Configure Multilingual Multihost

From **Hugo 0.31** we support multiple languages in a multihost configuration. See [this issue](https://github.com/gohugoio/hugo/issues/4027) for details.

This means that you can now configure a `baseURL` per `language`:


> If a `baseURL` is set on the `language` level, then all languages must have one and they must all be different.

Example:

{{< code-toggle file="config" >}}
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

```bash
public
├── en
└── fr
```

**All URLs (i.e `.Permalink` etc.) will be generated from that root. So the English home page above will have its `.Permalink` set to `https://example.com/`.**

When you run `hugo server` we will start multiple HTTP servers. You will typlically see something like this in the console:

```bash
Web Server is available at 127.0.0.1:1313 (bind address 127.0.0.1)
Web Server is available at 127.0.0.1:1314 (bind address 127.0.0.1)
Press Ctrl+C to stop
```

Live reload and `--navigateToChanged` between the servers work as expected.

### Taxonomies and Blackfriday

Taxonomies and [Blackfriday configuration][config] can also be set per language:


{{< code-toggle file="config" >}}
[Taxonomies]
tag = "tags"

[blackfriday]
angledQuotes = true
hrefTargetBlank = true

[languages]
[languages.en]
weight = 1
title = "English"
[languages.en.blackfriday]
angledQuotes = false

[languages.fr]
weight = 2
title = "Français"
[languages.fr.Taxonomies]
plaque = "plaques"
{{</ code-toggle >}}

## Translate Your Content

There are two ways to manage your content translations. Both ensure each page is assigned a language and is linked to its counterpart translations.

### Translation by filename

Considering the following example:

1. `/content/about.en.md`
2. `/content/about.fr.md`

The first file is assigned the English language and is linked to the second.
The second file is assigned the French language and is linked to the first.

Their language is __assigned__ according to the language code added as a __suffix to the filename__. 

By having the same **path and base filename**, the content pieces are __linked__ together as translated pages.

{{< note >}}
If a file has no language code, it will be assigned the default language.
{{</ note >}}

### Translation by content directory

This system uses different content directories for each of the languages. Each language's content directory is set using the `contentDir` param.

{{< code-toggle file="config" >}}

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

### Bypassing default linking.

Any pages sharing the same `translationKey` set in front matter will be linked as translated pages regardless of basename or location.

Considering the following example:

1. `/content/about-us.en.md`
2. `/content/om.nn.md`
3. `/content/presentation/a-propos.fr.md`

```yaml
# set in all three pages
translationKey: "about"
```

By setting the `translationKey` front matter param to `about` in all three pages, they will be __linked__ as translated pages.


### Localizing permalinks

Because paths and filenames are used to handle linking, all translated pages will share the same URL (apart from the language subdirectory).

To localize the URLs, the [`slug`]({{< ref "/content-management/organization/index.md#slug" >}}) or [`url`]({{< ref "/content-management/organization/index.md#url" >}}) front matter param can be set in any of the non-default language file. 

For example, a French translation (`content/about.fr.md`) can have its own localized slug.

{{< code-toggle >}}
Title: A Propos
slug: "a-propos"
{{< /code-toggle >}}


At render, Hugo will build both `/about/` and `/fr/a-propos/` while maintaining their translation linking.

{{% note %}}
If using `url`, remember to include the language part as well: `/fr/compagnie/a-propos/`.
{{%/ note %}}

### Page Bundles

To avoid the burden of having to duplicate files, each Page Bundle inherits the resources of its linked translated pages' bundles except for the content files (markdown files, html files etc...).

Therefore, from within a template, the page will have access to the files from all linked pages' bundles.

If, across the linked bundles, two or more files share the same basename, only one will be included and chosen as follows:

* File from current language bundle, if present.
* First file found across bundles by order of language `Weight`.

{{% note %}}
Page Bundle resources follow the same language assignment logic as content files, both by filename (`image.jpg`, `image.fr.jpg`) and by directory (`english/about/header.jpg`, `french/about/header.jpg`).
{{%/ note %}}

## Reference the Translated Content

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

The above can be put in a `partial` (i.e., inside `layouts/partials/`) and included in any template, whether a [single content page][contenttemplate] or the [homepage][]. It will not print anything if there are no translations for a given page.

The above also uses the [`i18n` function][i18func] described in the next section.

### List All Available Languages

`.AllTranslations` on a `Page` can be used to list all translations, including the page itself. On the home page it can be used to build a language navigator:


{{< code file="layouts/partials/allLanguages.html" >}}
<ul>
{{ range $.Site.Home.AllTranslations }}
<li><a href="{{ .Permalink }}">{{ .Language.LanguageName }}</a></li>
{{ end }}
</ul>
{{< /code >}}

## Translation of Strings

Hugo uses [go-i18n][] to support string translations. [See the project's source repository][go-i18n-source] to find tools that will help you manage your translation workflows.

Translations are collected from the `themes/<THEME>/i18n/` folder (built into the theme), as well as translations present in `i18n/` at the root of your project. In the `i18n`, the translations will be merged and take precedence over what is in the theme folder. Language files should be named according to [RFC 5646][] with names such as `en-US.toml`, `fr.toml`, etc.

{{% note %}}
From **Hugo 0.31** you no longer need to use a valid language code. It can be anything.

See https://github.com/gohugoio/hugo/issues/3564

{{% /note %}}

From within your templates, use the `i18n` function like this:

```
{{ i18n "home" }}
```

This uses a definition like this one in `i18n/en-US.toml`:

```
[home]
other = "Home"
```

Often you will want to use to the page variables in the translations strings. To do that, pass on the "." context when calling `i18n`:

```
{{ i18n "wordCount" . }}
```

This uses a definition like this one in `i18n/en-US.toml`:

```
[wordCount]
other = "This article has {{ .WordCount }} words."
```
An example of singular and plural form:

```
[readingTime]
one = "One minute to read"
other = "{{.Count}} minutes to read"
```
And then in the template:

```
{{ i18n "readingTime" .ReadingTime }}
```

## Customize Dates

At the time of this writing, Go does not yet have support for internationalized locales for dates, but if you do some work, you can simulate it. For example, if you want to use French month names, you can add a data file like ``data/mois.yaml`` with this content:

~~~yaml
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

...then index the non-English date names in your templates like so:

~~~html
<time class="post-date" datetime="{{ .Date.Format '2006-01-02T15:04:05Z07:00' | safeHTML }}">
  Article publié le {{ .Date.Day }} {{ index $.Site.Data.mois (printf "%d" .Date.Month) }} {{ .Date.Year }} (dernière modification le {{ .Lastmod.Day }} {{ index $.Site.Data.mois (printf "%d" .Lastmod.Month) }} {{ .Lastmod.Year }})
</time>
~~~

This technique extracts the day, month and year by specifying ``.Date.Day``, ``.Date.Month``, and ``.Date.Year``, and uses the month number as a key, when indexing the month name data file.

## Menus

You can define your menus for each language independently. Creating multilingual menus works just like [creating regular menus][menus], except they're defined in language-specific blocks in the configuration file:

```
defaultContentLanguage = "en"

[languages.en]
weight = 0
languageName = "English"

[[languages.en.menu.main]]
url    = "/"
name   = "Home"
weight = 0


[languages.de]
weight = 10
languageName = "Deutsch"

[[languages.de.menu.main]]
url    = "/"
name   = "Startseite"
weight = 0
```

The rendering of the main navigation works as usual. `.Site.Menus` will just contain the menu in the current language. Note that `absLangURL` below will link to the correct locale of your website. Without it, menu entries in all languages would link to the English version, since it's the default content language that resides in the root directory.

```
<ul>
    {{- $currentPage := . -}}
    {{ range .Site.Menus.main -}}
    <li class="{{ if $currentPage.IsMenuCurrent "main" . }}active{{ end }}">
        <a href="{{ .URL | absLangURL }}">{{ .Name }}</a>
    </li>
    {{- end }}
</ul>

```

## Missing Translations

If a string does not have a translation for the current language, Hugo will use the value from the default language. If no default value is set, an empty string will be shown.

While translating a Hugo website, it can be handy to have a visual indicator of missing translations. The [`enableMissingTranslationPlaceholders` configuration option][config] will flag all untranslated strings with the placeholder `[i18n] identifier`, where `identifier` is the id of the missing translation.

{{% note %}}
Hugo will generate your website with these missing translation placeholders. It might not be suitable for production environments.
{{% /note %}}

For merging of content from other languages (i.e. missing content translations), see [lang.Merge](/functions/lang.merge/).

To track down missing translation strings, run Hugo with the `--i18n-warnings` flag:

```
 hugo --i18n-warnings | grep i18n
i18n|MISSING_TRANSLATION|en|wordCount
```

## Multilingual Themes support

To support Multilingual mode in your themes, some considerations must be taken for the URLs in the templates. If there is more than one language, URLs must meet the following criteria:

* Come from the built-in `.Permalink` or `.RelPermalink`
* Be constructed with the [`relLangURL` template function][rellangurl] or the [`absLangURL` template function][abslangurl] **OR** be prefixed with `{{ .LanguagePrefix }}`

If there is more than one language defined, the `LanguagePrefix` variable will equal `/en` (or whatever your `CurrentLanguage` is). If not enabled, it will be an empty string (and is therefore harmless for single-language Hugo websites).

[abslangurl]: /functions/abslangurl
[config]: /getting-started/configuration/
[contenttemplate]: /templates/single-page-templates/
[go-i18n-source]: https://github.com/nicksnyder/go-i18n
[go-i18n]: https://github.com/nicksnyder/go-i18n
[homepage]: /templates/homepage/
[i18func]: /functions/i18n/
[menus]: /content-management/menus/
[rellangurl]: /functions/rellangurl
[RFC 5646]: https://tools.ietf.org/html/rfc5646
[singles]: /templates/single-page-templates/
