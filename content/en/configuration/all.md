---
title: All settings
description: The complete list of Hugo configuration settings.
categories: []
keywords: []
weight: 20
aliases: [/getting-started/configuration/]
---

## Settings

archetypeDir
: (`string`) The designated directory for [archetypes](g). Default is `archetypes`. {{% module-mounts-note %}}

assetDir
: (`string`) The designated directory for [global resources](g). Default is `assets`. {{% module-mounts-note %}}

baseURL
: (`string`) The absolute URL of your published site including the protocol, host, path, and a trailing slash.

build
: See [configure build](/configuration/build/).

buildDrafts
: (`bool`) Whether to include draft content when building a site. Default is `false`.

buildExpired
: (`bool`) Whether to include expired content when building a site. Default is `false`.

buildFuture
: (`bool`) Whether to include future content when building a site. Default is `false`.

cacheDir
: (`string`) The designated cache directory. See&nbsp;[details](#cache-directory).

caches
: See [configure file caches](/configuration/caches/).

canonifyURLs
: (`bool`) See&nbsp;[details](/content-management/urls/#canonical-urls) before enabling this feature. Default is `false`.

capitalizeListTitles
: {{< new-in 0.123.3 />}}
: (`bool`) Whether to capitalize automatic list titles. Applicable to section, taxonomy, and term pages. Default is `true`. Use the [`titleCaseStyle`](#titlecasestyle) setting to configure capitalization rules.

cascade
: See [configure cascade](/configuration/cascade/).

cleanDestinationDir
: (`bool`) Whether to remove files from the site's destination directory that do not have corresponding files in the `static` directory during the build. Default is `false`.

contentDir
: (`string`) The designated directory for content files. Default is `content`. {{% module-mounts-note %}}

copyright
: (`string`) The copyright notice for a site, typically displayed in the footer.

dataDir
: (`string`) The designated directory for data files. Default is `data`. {{% module-mounts-note %}}

defaultContentLanguage
: (`string`) The project's default language key, conforming to the syntax described in [RFC 5646]. This value must match one of the defined language keys. Default is `en`.

defaultContentLanguageInSubdir
: (`bool`) Whether to publish the default language site to a subdirectory matching the `defaultContentLanguage`. Default is `false`.

defaultOutputFormat
: (`string`) The default output format for the site. If unspecified, the first available format in the defined order (by weight, then alphabetically) will be used.

deployment
: See [configure deployment](/configuration/deployment/).

disableAliases
: (`bool`) Whether to disable generation of alias redirects. Even if this option is enabled, the defined aliases will still be present on the page. This allows you to manage redirects separately, for example, by generating 301 redirects in an `.htaccess` file or a Netlify `_redirects` file using a custom output format. Default is `false`.

disableDefaultLanguageRedirect
: {{< new-in 0.140.0 />}}
: (`bool`) Whether to disable generation of the alias redirect to the default language when `DefaultContentLanguageInSubdir` is `true`. Default is `false`.

disableHugoGeneratorInject
: (`bool`) Whether to disable injection of a `<meta name="generator">` tag into the home page. Default is `false`.

disableKinds
: (`[]string`) A slice of page [kinds](g) to disable during the build process, any of `404`, `home`, `page`, `robotstxt`, `rss`, `section`, `sitemap`, `taxonomy`, or `term`.

disableLanguages
: (`[]string`) A slice of language keys representing the languages to disable during the build process. Although this is functional, consider using the [`disabled`] key under each language instead.

disableLiveReload
: (`bool`) Whether to disable automatic live reloading of the browser window. Default is `false`.

disablePathToLower
: (`bool`) Whether to disable transformation of page URLs to lower case.

enableEmoji
: (`bool`) Whether to allow emoji in Markdown. Default is `false`.

enableGitInfo
: (`bool`) For sites under Git version control, whether to enable the [`GitInfo`] object for each page. With the [default front matter configuration], the `Lastmod` method on a `Page` object will return the Git author date. Default is `false`.

enableMissingTranslationPlaceholders
: (`bool`) Whether to show a placeholder instead of the default value or an empty string if a translation is missing. Default is `false`.

enableRobotsTXT
: (`bool`) Whether to enable generation of a `robots.txt` file. Default is `false`.

environment
: (`string`) The build environment. Default is `production` when running `hugo` and `development` when running `hugo server`.

frontmatter
: See [configure front matter](/configuration/front-matter/).

hasCJKLanguage
: (`bool`) Whether to automatically detect [CJK](g) languages in content. Affects the values returned by the [`WordCount`] and [`FuzzyWordCount`] methods. Default is `false`.

HTTPCache
: See [configure HTTP cache](/configuration/http-cache/).

i18nDir
: (`string`) The designated directory for translation tables. Default is `i18n`. {{% module-mounts-note %}}

ignoreCache
: (`bool`) Whether to ignore the cache directory. Default is `false`.

ignoreFiles
: (`[]string`) A slice of [regular expressions](g) used to exclude specific files from a build. These expressions are matched against the absolute file path and apply to files within the `content`, `data`, and `i18n` directories. For more advanced file exclusion options, see the section on [module mounts].

ignoreLogs
: (`[]string`) A slice of message identifiers corresponding to warnings and errors you wish to suppress. See [`erroridf`] and [`warnidf`].

ignoreVendorPaths
: (`string`) A [glob](g) pattern matching the module paths to exclude from the `_vendor` directory.

imaging
: See [configure imaging](/configuration/imaging/).

languageCode
: (`string`) The site's language tag, conforming to the syntax described in [RFC 5646]. This value does not affect translations or localization. Hugo uses this value to populate:

  - The `language` element in the [embedded RSS template]
  - The `lang` attribute of the `html` element in the [embedded alias template]
  - The `og:locale` `meta` element in the [embedded Open Graph template]

  When present in the root of the configuration, this value is ignored if one or more language keys exists. Please specify this value independently for each language key.

languages
: See [configure languages](/configuration/languages/).

layoutDir
: (`string`) The designated directory for templates. Default is `layouts`. {{% module-mounts-note %}}

mainSections
: (`string` or `[]string`) The main sections of a site. If set, the [`MainSections`] method on the `Site` object returns the given sections, otherwise it returns the section with the most pages.

markup
: See [configure markup](/configuration/markup/).

mediaTypes
: See [configure media types](/configuration/media-types/).

menus
: See [configure menus](/configuration/menus/).

minify
: See [configure minify](/configuration/minify/).

module
: See [configure modules](/configuration/module/).

newContentEditor
: (`string`) The editor to use when creating new content.

noBuildLock
: (`bool`) Whether to disable creation of the `.hugo_build.lock` file. Default is `false`.

noChmod
: (`bool`) Whether to disable synchronization of file permission modes. Default is `false`.

noTimes
: (`bool`) Whether to disable synchronization of file modification times. Default is `false`.

outputFormats
: See [configure output formats](/configuration/output-formats/).

outputs
: See [configure outputs](/configuration/outputs/).

page
: See [configure page](/configuration/page/).

pagination
: See [configure pagination](/configuration/pagination/).

panicOnWarning
: (`bool`) Whether to panic on the first WARNING. Default is `false`.

params
: See [configure params](/configuration/params/).

permalinks
: See [configure permalinks](/configuration/permalinks/).

pluralizeListTitles
: (`bool`) Whether to pluralize automatic list titles. Applicable to section pages. Default is `true`.

printI18nWarnings
: (`bool`) Whether to log WARNINGs for each missing translation. Default is `false`.

printPathWarnings
: (`bool`) Whether to log WARNINGs when Hugo publishes two or more files to the same path. Default is `false`.

printUnusedTemplates
: (`bool`) Whether to log WARNINGs for each unused template. Default is `false`.

privacy
: See [configure privacy](/configuration/privacy/).

publishDir
: (`string`) The designated directory for publishing the site. Default is `public`.

refLinksErrorLevel
: (`string`) The logging error level to use when the `ref` and `relref` functions, methods, and shortcodes are unable to resolve a reference to a page. Either `ERROR` or `WARNING`. Any `ERROR` will fail the build. Default is `ERROR`.

refLinksNotFoundURL
: (`string`) The URL to return when the `ref` and `relref` functions, methods, and shortcodes are unable to resolve a reference to a page.

related
: See [configure related content](/configuration/related-content/).

relativeURLs
: (`bool`) See&nbsp;[details](/content-management/urls/#relative-urls) before enabling this feature. Default is `false`.

removePathAccents
: (`bool`) Whether to remove [non-spacing marks](https://www.compart.com/en/unicode/category/Mn) from [composite characters](https://en.wikipedia.org/wiki/Precomposed_character) in content paths. Default is `false`.

renderSegments
: {{< new-in 0.124.0 />}}
: (`[]string`) A slice of [segments](g) to render. If omitted, all segments are rendered. This option is typically set via a command-line flag, such as `hugo --renderSegments segment1,segment2`. The provided segment names must correspond to those defined in the [`segments`] configuration.

resourceDir
: (`string`) The designated directory for caching output from [asset pipelines](g). Default is `resources`.

security
: See [configure security](/configuration/security/).

sectionPagesMenu
: (`string`) When set, each top-level section will be added to the menu identified by the provided value. See&nbsp;[details](/content-management/menus/#define-automatically).

segments
: See [configure segments](/configuration/segments/).

server
: See [configure server](/configuration/server/).

services
: See [configure services](/configuration/services/).

sitemap
: See [configure sitemap](/configuration/sitemap/).

staticDir
: (`string`) The designated directory for static files. Default is `static`. {{% module-mounts-note %}}

summaryLength
: (`int`) Applicable to [automatic summaries], the minimum number of words returned by the [`Summary`] method on a `Page` object. The `Summary` method will return content truncated at the paragraph boundary closest to the specified `summaryLength`, but at least this minimum number of words.

taxonomies
: See [configure taxonomies](/configuration/taxonomies/).

templateMetrics
: (`bool`) Whether to print template execution metrics to the console. Default is `false`. See&nbsp;[details](/troubleshooting/performance/#template-metrics).

templateMetricsHints
: (`bool`) Whether to print template execution improvement hints to the console. Applicable when `templateMetrics` is `true`. Default is `false`. See&nbsp;[details](/troubleshooting/performance/#template-metrics).

theme
: (`string` or `[]string`) The [theme](g) to use. Multiple themes can be listed, with precedence given from left to right. See&nbsp;[details](/hugo-modules/theme-components/).

themesDir
: (`string`) The designated directory for themes. Default is `themes`.

timeout
: (`string`) The timeout for generating page content, either as a [duration] or in seconds. This timeout is used to prevent infinite recursion during content generation. You may need to increase this value if your pages take a long time to generate, for example, due to extensive image processing or reliance on remote content. Default is `60s`.

timeZone
: (`string`) The time zone used to parse dates without time zone offsets, including front matter date fields and values passed to the [`time.AsTime`] and [`time.Format`] template functions. The list of valid values may be system dependent, but should include `UTC`, `Local`, and any location in the [IANA Time Zone Database]. For example, `America/Los_Angeles` and `Europe/Oslo` are valid time zones.

title
: (`string`) The site title.

titleCaseStyle
: (`string`) The capitalization rules to follow when Hugo automatically generates a section title, or when using the [`strings.Title`] function. One of `ap`, `chicago`, `go`, `firstupper`, or `none`. Default is `ap`. See&nbsp;[details](#title-case-style).

uglyurls
: See [configure ugly URLs](/configuration/ugly-urls/).

## Cache directory

Hugo's file cache directory is configurable via the [`cacheDir`] configuration option or the `HUGO_CACHEDIR` environment variable. If neither is set, Hugo will use, in order of preference:

1. If running on Netlify: `/opt/build/cache/hugo_cache/`. This means that if you run your builds on Netlify, all caches configured with `:cacheDir` will be saved and restored on the next build. For other [CI/CD](g) vendors, please read their documentation. For an CircleCI example, see [this configuration].
1. In a `hugo_cache` directory below the OS user cache directory as defined by Go's [os.UserCacheDir] function. On Unix systems, per the [XDG base directory specification], this is `$XDG_CACHE_HOME` if non-empty, else `$HOME/.cache`. On MacOS, this is `$HOME/Library/Caches`. On Windows, this is`%LocalAppData%`. On Plan 9, this is `$home/lib/cache`.
1. In a  `hugo_cache_$USER` directory below the OS temp dir.

To determine the current `cacheDir`:

```sh
hugo config | grep cachedir
```

## Title case style

Hugo's [`titleCaseStyle`] setting governs capitalization for automatically generated section titles and the [`strings.Title`] function. By default, it follows the capitalization rules published in the Associated Press Stylebook. Change this setting to use other capitalization rules.

ap
: Use the capitalization rules published in the [Associated Press Stylebook]. This is the default.

chicago
: Use the capitalization rules published in the [Chicago Manual of Style].

go
: Capitalize the first letter of every word.

firstupper
: Capitalize the first letter of the first word.

none
: Disable transformation of automatic section titles, and disable the transformation performed by the `strings.Title` function. This is useful if you would prefer to manually capitalize section titles as needed, and to bypass opinionated theme usage of the `strings.Title` function.

## Localized settings

Some configuration settings, such as menus and custom parameters, can be defined separately for each language. See [configure languages](/configuration/languages/#localized-settings).

[`cacheDir`]: #cachedir
[`disabled`]: /configuration/languages/#disabled
[`erroridf`]: /functions/fmt/erroridf/
[`FuzzyWordCount`]: /methods/page/fuzzywordcount/
[`GitInfo`]: /methods/page/gitinfo/
[`MainSections`]: /methods/site/mainsections/
[`segments`]: /configuration/segments/
[`strings.Title`]: /functions/strings/title/
[`Summary`]: /methods/page/summary/
[`time.AsTime`]: /functions/time/astime/
[`time.Format`]: /functions/time/format/
[`titleCaseStyle`]: #titlecasestyle
[`warnidf`]: /functions/fmt/warnidf/
[`WordCount`]: /methods/page/wordcount/
[Associated Press Stylebook]: https://www.apstylebook.com/
[automatic summaries]: /content-management/summaries/#automatic-summary
[Chicago Manual of Style]: https://www.chicagomanualofstyle.org/home.html
[default front matter configuration]: /configuration/front-matter/
[duration]: https://pkg.go.dev/time#Duration
[embedded alias template]: <{{% eturl alias %}}>
[embedded Open Graph template]: <{{% eturl opengraph %}}>
[embedded RSS template]: <{{% eturl rss %}}>
[IANA Time Zone Database]: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
[module mounts]: /configuration/module/#mounts
[os.UserCacheDir]: https://pkg.go.dev/os#UserCacheDir
[RFC 5646]: https://datatracker.ietf.org/doc/html/rfc5646#section-2.1
[this configuration]: https://github.com/bep/hugo-sass-test/blob/6c3960a8f4b90e8938228688bc49bdcdd6b2d99e/.circleci/config.yml
[XDG base directory specification]: https://specifications.freedesktop.org/basedir-spec/latest/
