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
aliases: [/content/multilingual/,/content-management/multilingual/,/tutorials/create-a-multilingual-site/]
toc: true
---

You should define the available languages in a `languages` section in your site configuration.

## Configure Languages

The following is an example of a TOML site configuration for a multilingual Hugo project:

{{< code file="config.toml" download="config.toml" >}}
DefaultContentLanguage = "en"
copyright = "Everything is mine"

[params.navigation]
help  = "Help"

[languages]
[languages.en]
title = "My blog"
weight = 1
linkedin = "english-link"

[languages.fr]
copyright = "Tout est à moi"
title = "Mon blog"
weight = 2
linkedin = "lien-francais"

# skip params key for front matter
[languages.fr.navigation]
help  = "Aide"
{{< /code >}}

Anything not defined in a `[languages]` block will fall back to the global
value for that key (e.g., `copyright` for the English [`en`] language).

With the configuration above, all content, sitemap, RSS feeds, paginations,
and taxonomy pages will be rendered below `/` in English (your default content language) and then below `/fr` in French.

When working with front matter `Params` in [single page templates][singles], omit the `params` in the key for the translation.

If you want all of the languages to be put below their respective language code, enable `defaultContentLanguageInSubdir: true`.

Only the obvious non-global options can be overridden per language. Examples of global options are `baseURL`, `buildDrafts`, etc.

## Configure Multilingual Multihost

From **Hugo 0.31** we support multiple languages in a multihost configuration. See [this issue](https://github.com/gohugoio/hugo/issues/4027) for details.

This means that you can now configure a `baseURL` per `language`:


> If a `baseURL` is set on the `language` level, then all languages must have one and they must all be different.

Example:

```bash
[languages]
[languages.no]
baseURL = "https://example.no"
languageName = "Norsk"
weight = 1
title = "På norsk"

[languages.en]
baseURL = "https://example.com"
languageName = "English"
weight = 2
title = "In English"
```

With the above, the two sites will be generated into `public` with their own root:

```bash
public
├── en
└── no
```

**All URLs (i.e `.Permalink` etc.) will be generated from that root. So the English home page above will have its `.Permalink` set to `https://example.com/`.**

When you run `hugo server` we will start multiple HTTP servers. You will typlically see something like this in the console:

```bash
Web Server is available at 127.0.0.1:1313 (bind address 127.0.0.1)
Web Server is available at 127.0.0.1:1314 (bind address 127.0.0.1)
Press Ctrl+C to stop
```

Live reload and `--navigateToChanged` between the servers work as expected.

## Taxonomies and Blackfriday

Taxonomies and [Blackfriday configuration][config] can also be set per language:


{{< code file="bf-config.toml" >}}
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
{{< /code >}}

## Translate Your Content

Translated articles are identified by the name of the content file.

### Examples of Translated Articles

1. `/content/about.en.md`
2. `/content/about.fr.md`

In this example, the `about.md` will be assigned the configured `defaultContentLanguage`. 

1. `/content/about.md`
2. `/content/about.fr.md`

This way, you can slowly start to translate your current content without having to rename everything. If left unspecified, the default value for `defaultContentLanguage` is `en`.

By having the same **directory and base filename**, the content pieces are linked together as translated pieces.

You can also set the key used to link the translations explicitly in front matter:

```yaml
translationKey: "my-story"
```


{{% note %}}
**Before Hugo 0.31**, the file's directory was not considered when looking for translations. This did not work when you named all of your content files, say, `index.md`. Now we use the full content path.
{{% /note %}}

If you need distinct URLs per language, you can set the slug in the non-default language file. For example, you can define a custom slug for a French translation in the front matter of `content/about.fr.md` as follows:

```yaml
slug: "a-propos"

```

At render, Hugo will build both `/about/` and `/a-propos/` as properly linked translated pages.


## Link to Translated Content

To create a list of links to translated content, use a template similar to the following:

{{< code file="layouts/partials/i18nlist.html" >}}
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
{{< /code >}}

The above can be put in a `partial` (i.e., inside `layouts/partials/`) and included in any template, be it for a [single content page][contenttemplate] or the [homepage][]. It will not print anything if there are no translations for a given page.

The above also uses the [`i18n` function][i18func] described in the next section.

## List All Available Languages

`.AllTranslations` on a `Page` can be used to list all translations, including itself. Called on the home page it can be used to build a language navigator:


{{< code file="layouts/partials/allLanguages.html" >}}
<ul>
{{ range $.Site.Home.AllTranslations }}
<li><a href="{{ .}}">{{ .Language.LanguageName }}</a></li>
{{ end }}
</ul>
{{< /code >}}

## Translation of Strings

Hugo uses [go-i18n][] to support string translations. [See the project's source repository][go-i18n-source] to find tools that will help you manage your translation workflows.

Translations are collected from the `themes/<THEME>/i18n/` folder (built into the theme), as well as translations present in `i18n/` at the root of your project. In the `i18n`, the translations will be merged and take precedence over what is in the theme folder. Language files should be named according to [RFC 5646][] with names such as `en-US.toml`, `fr.toml`, etc.

{{% note %}}
From **Hugo 0.31** you no longer need to use a valid language code. It _can be_ anything.

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
one = "One minute read"
other = "{{.Count}} minutes read"
```
And then in the template:

```
{{ i18n "readingTime" .ReadingTime }}
```
To track down missing translation strings, run Hugo with the `--i18n-warnings` flag:

```
 hugo --i18n-warnings | grep i18n
i18n|MISSING_TRANSLATION|en|wordCount
```

## Customize Dates

At the time of this writing, Golang does not yet have support for internationalized locales, but if you do some work, you can simulate it. For example, if you want to use French month names, you can add a data file like ``data/mois.yaml`` with this content:

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

... then index the non-English date names in your templates like so:

~~~html
<time class="post-date" datetime="{{ .Date.Format "2006-01-02T15:04:05Z07:00" | safeHTML }}">
  Article publié le {{ .Date.Day }} {{ index $.Site.Data.mois (printf "%d" .Date.Month) }} {{ .Date.Year }} (dernière modification le {{ .Lastmod.Day }} {{ index $.Site.Data.mois (printf "%d" .Lastmod.Month) }} {{ .Lastmod.Year }})
</time>
~~~

This technique extracts the day, month and year by specifying ``.Date.Day``, ``.Date.Month``, and ``.Date.Year``, and uses the month number as a key, when indexing the month name data file.

## Menus

You can define your menus for each language independently. The [creation of a menu][menus] works analogous to earlier versions of Hugo, except that they have to be defined in their language-specific block in the configuration file:

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

The rendering of the main navigation works as usual. `.Site.Menus` will just contain the menu of the current language. Pay attention to the generation of the menu links. `absLangURL` takes care that you link to the correct locale of your website. Otherwise, both menu entries would link to the English version as the default content language that resides in the root directory.

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

## Missing translations

If a string does not have a translation for the current language, Hugo will use the value from the default language. If no default value is set, an empty string will be shown.

While translating a Hugo website, it can be handy to have a visual indicator of missing translations. The [`enableMissingTranslationPlaceholders` configuration option][config] will flag all untranslated strings with the placeholder `[i18n] identifier`, where `identifier` is the id of the missing translation.

{{% note %}}
Hugo will generate your website with these missing translation placeholders. It might not be suited for production environments.
{{% /note %}}

## Multilingual Themes support

To support Multilingual mode in your themes, some considerations must be taken for the URLs in the templates. If there is more than one language, URLs must meet the following criteria:

* Come from the built-in `.Permalink` or `.URL`
* Be constructed with
    * The [`relLangURL` template function][rellangurl] or the [`absLangURL` template function][abslangurl] **OR**
    * Prefixed with `{{ .LanguagePrefix }}`

If there is more than one language defined, the `LanguagePrefix` variable will equal `/en` (or whatever your `CurrentLanguage` is). If not enabled, it will be an empty string and is therefore harmless for single-language Hugo websites.

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
