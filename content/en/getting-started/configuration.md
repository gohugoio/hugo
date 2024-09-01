---
title: Configure Hugo
linkTitle: Configuration
description: How to configure your Hugo site.
categories: [getting started,fundamentals]
keywords: [configuration,toml,yaml,json]
menu:
  docs:
    parent: getting-started
    weight: 40
weight: 40
toc: true
aliases: [/overview/source-directory/,/overview/configuration/]
---

## Configuration file

Create a site configuration file in the root of your project directory, naming it `hugo.toml`, `hugo.yaml`, or `hugo.json`, with that order of precedence.

```text
my-project/
└── hugo.toml
```

{{% note %}}
With v0.109.0 and earlier the basename of the site configuration file was `config` instead of `hugo`. You can use either, but should transition to the new naming convention when practical.
{{% /note %}}

A simple example:

{{< code-toggle file=hugo >}}
baseURL = 'https://example.org/'
languageCode = 'en-us'
title = 'ABC Widgets, Inc.'
[params]
subtitle = 'The Best Widgets on Earth'
[params.contact]
email = 'info@example.org'
phone = '+1 202-555-1212'
{{< /code-toggle >}}

To use a different configuration file when building your site, use the `--config` flag:

```sh
hugo --config other.toml
```

Combine two or more configuration files, with left-to-right precedence:

```sh
hugo --config a.toml,b.yaml,c.json
```

{{% note %}}
See the specifications for each file format: [TOML], [YAML], and [JSON].

[TOML]: https://toml.io/en/latest
[YAML]: https://yaml.org/spec/
[JSON]: https://datatracker.ietf.org/doc/html/rfc7159
{{% /note %}}

## Configuration directory

Instead of a single site configuration file, split your configuration by [environment], root configuration key, and language. For example:

[environment]: /getting-started/glossary/#environment

```text
my-project/
└── config/
    ├── _default/
    │   ├── hugo.toml
    │   ├── menus.en.toml
    │   ├── menus.de.toml
    │   └── params.toml
    └── production/
        └── params.toml
```

The root configuration keys are `build`, `caches`, `cascade`, `deployment`, `frontmatter`, `imaging`, `languages`, `markup`, `mediatypes`, `menus`, `minify`, `module`, `outputformats`, `outputs`, `params`, `permalinks`, `privacy`, `related`, `security`, `segments`, `server`, `services`, `sitemap`, and `taxonomies`.

### Omit the root key

When splitting the configuration by root key, omit the root key in the given file. For example, these are equivalent:

{{< code-toggle file=hugo >}}
[params]
foo = 'bar'
{{< /code-toggle >}}

{{< code-toggle file=params >}}
foo = 'bar'
{{< /code-toggle >}}

### Recursive parsing

Hugo parses the `config` directory recursively, allowing you to organize the files into subdirectories. For example:

```text
my-project/
└── config/
    └── _default/
        ├── navigation/
        │   ├── menus.de.toml
        │   └── menus.en.toml
        └── hugo.toml
```

### Example

```text
my-project/
└── config/
    ├── _default/
    │   ├── hugo.toml
    │   ├── menus.en.toml
    │   ├── menus.de.toml
    │   └── params.toml
    ├── production/
    │   ├── hugo.toml
    │   └── params.toml
    └── staging/
        ├── hugo.toml
        └── params.toml
```

Considering the structure above, when running `hugo --environment staging`, Hugo will use every setting from `config/_default` and merge `staging`'s on top of those.

Let's take an example to understand this better. Let's say you are using Google Analytics for your website. This requires you to specify a [Google tag ID] in your site configuration:

[Google tag ID]: https://support.google.com/tagmanager/answer/12326985?hl=en

{{< code-toggle file=hugo copy=false >}}
[services.googleAnalytics]
ID = 'G-XXXXXXXXX'
{{< /code-toggle >}}

Now consider the following scenario:

1. You don't want to load the analytics code when running `hugo server`.
2. You want to use different Google tag IDs for your production and staging environments. For example:

    - `G-PPPPPPPPP` for production
    - `G-SSSSSSSSS` for staging

To satisfy these requirements, configure your site as follows:

1. `config/_default/hugo.toml`

    Exclude the `services.googleAnalytics` section. This will prevent loading of the analytics code when you run `hugo server`.

    By default, Hugo sets its `environment` to `development` when running `hugo server`. In the absence of a `config/development` directory, Hugo uses the `config/_default` directory.

2. `config/production/hugo.toml`

    Include this section only:

    {{< code-toggle file=hugo copy=false >}}
    [services.googleAnalytics]
    ID = 'G-PPPPPPPPP'
    {{< /code-toggle >}}

    You do not need to include other parameters in this file. Include only those parameters that are specific to your production environment. Hugo will merge these parameters with the default configuration.

    By default, Hugo sets its `environment` to `production` when running `hugo`. The analytics code will use the `G-PPPPPPPPP` tag ID.

3. `config/staging/hugo.toml`

    Include this section only:

    {{< code-toggle file=hugo copy=false >}}
    [services.googleAnalytics]
    ID = 'G-SSSSSSSSS'
    {{< /code-toggle >}}

    You do not need to include other parameters in this file. Include only those parameters that are specific to your staging environment. Hugo will merge these parameters with the default configuration.

    To build your staging site, run `hugo --environment staging`. The analytics code will use the `G-SSSSSSSSS` tag ID.

## Merge configuration from themes

The configuration value for `_merge` can be one of:

none
: No merge.

shallow
: Only add values for new keys.

deep
: Add values for new keys, merge existing.

Note that you don't need to be so verbose as in the default setup below; a `_merge` value higher up will be inherited if not set.

{{< code-toggle file=hugo dataKey="config_helpers.mergeStrategy" skipHeader=true />}}

## All configuration settings

###### archetypeDir

(`string`) The directory where Hugo finds archetype files (content templates). Default is `archetypes`. {{% module-mounts-note %}}

###### assetDir

(`string`) The directory where Hugo finds asset files used in [Hugo Pipes](/hugo-pipes/). Default is `assets`. {{% module-mounts-note %}}

###### baseURL

(`string`) The absolute URL (protocol, host, path, and trailing slash) of your published site (e.g., `https://www.example.org/docs/`).

###### build

See [Configure Build](#configure-build).

###### buildDrafts

(`bool`) Include drafts when building. Default is `false`.

###### buildExpired

(`bool`) Include content already expired. Default is `false`.

###### buildFuture

(`bool`) Include content with a future publication date. Default is `false`.

###### caches

See [Configure File Caches](#configure-file-caches).

###### capitalizeListTitles

{{< new-in 0.123.3 >}}

(`bool`) Whether to capitalize automatic list titles. Applicable to section, taxonomy, and term pages. Default is `true`. You can change the capitalization style in your site configuration to one of `ap`, `chicago`, `go`, `firstupper`, or `none`. See [details].

[details]: /getting-started/configuration/#configure-title-case

###### cascade

Pass down default configuration values (front matter) to pages in the content tree. The options in site config is the same as in page front matter, see [Front Matter Cascade](/content-management/front-matter#cascade).

{{% note %}}
For a website in a single language, define the `[[cascade]]` in [Front Matter](/content-management/front-matter#cascade). For a multilingual website, define the `[[cascade]]` in [Site Config](/getting-started/configuration/#cascade).

To remain consistent and prevent unexpected behavior, do not mix these strategies.
{{% /note %}}

###### canonifyURLs

(`bool`) See [details](/content-management/urls/#canonical-urls) before enabling this feature. Default is `false`.

###### cleanDestinationDir

(`bool`) When building, removes files from destination not found in static directories. Default is `false`.

###### contentDir

(`string`) The directory from where Hugo reads content files.  Default is `content`. {{% module-mounts-note %}}

###### copyright

(`string`) Copyright notice for your site, typically displayed in the footer.

###### dataDir

(`string`) The directory from where Hugo reads data files. Default is `data`. {{% module-mounts-note %}}

###### defaultContentLanguage

(`string`) Content without language indicator will default to this language. Default is `en`.

###### defaultContentLanguageInSubdir

(`bool`) Render the default content language in subdir, e.g. `content/en/`. The site root `/` will then redirect to `/en/`. Default is `false`.

###### disableAliases

(`bool`) Will disable generation of alias redirects. Note that even if `disableAliases` is set, the aliases themselves are preserved on the page. The motivation with this is to be able to generate 301 redirects in an `.htaccess`, a Netlify `_redirects` file or similar using a custom output format. Default is `false`.

###### disableHugoGeneratorInject

(`bool`) Hugo will, by default, inject a generator meta tag in the HTML head on the _home page only_. You can turn it off, but we would really appreciate if you don't, as this is a good way to watch Hugo's popularity on the rise. Default is `false`.

###### disableKinds

(`string slice`) Disable rendering of the specified page [kinds], any of `404`, `home`, `page`, `robotstxt`, `rss`, `section`, `sitemap`, `taxonomy`, or `term`.

[kinds]: /getting-started/glossary/#page-kind

###### disableLiveReload

(`bool`) Disable automatic live reloading of browser window. Default is `false`.

###### disablePathToLower

(`bool`) Do not convert the url/path to lowercase. Default is `false`.

###### enableEmoji

(`bool`) Enable Emoji emoticons support for page content; see the [emoji shortcode quick reference guide](/quick-reference/emojis/). Default is `false`.

###### enableGitInfo

(`bool`) Enable `.GitInfo` object for each page (if the Hugo site is versioned by Git). This will then update the `Lastmod` parameter for each page using the last git commit date for that content file. Default is `false`.

###### enableMissingTranslationPlaceholders

(`bool`) Show a placeholder instead of the default value or an empty string if a translation is missing. Default is `false`.

###### enableRobotsTXT

(`bool`) Enable generation of `robots.txt` file. Default is `false`.

###### frontmatter

See [Front matter Configuration](#configure-front-matter).

###### hasCJKLanguage

(`bool`) If true, auto-detect Chinese/Japanese/Korean Languages in the content. This will make `.Summary` and `.WordCount` behave correctly for CJK languages. Default is `false`.

###### imaging

See [image processing configuration](/content-management/image-processing/#imaging-configuration).

###### languageCode

(`string`) A language tag as defined by [RFC 5646](https://datatracker.ietf.org/doc/html/rfc5646). This value is used to populate:

- The `<language>` element in the embedded [RSS template]({{% eturl rss %}})
- The `lang` attribute of the `<html>` element in the embedded [alias template]({{% eturl alias %}})
- The `og:locale` `meta` element in the embedded [Open Graph template]({{% eturl opengraph %}})

When present in the root of the configuration, this value is ignored if one or more language keys exists. Please specify this value independently for each language key.

###### languages

See [Configure Languages](/content-management/multilingual/#configure-languages).

###### disableLanguages

See [Disable a Language](/content-management/multilingual/#disable-a-language)

###### markup

See [Configure Markup](/getting-started/configuration-markup).

###### mediaTypes

See [Configure Media Types](/templates/output-formats/#media-types).

###### menus

See [Menus](/content-management/menus/#define-in-site-configuration).

###### minify

See [Configure Minify](#configure-minify).

###### module

Module configuration see [module configuration](/hugo-modules/configuration/).

###### newContentEditor

(`string`) The editor to use when creating new content.

###### noChmod

(`bool`) Don't sync permission mode of files. Default is `false`.

###### noTimes

(`bool`) Don't sync modification time of files. Default is `false`.

###### outputFormats

See [custom output formats].

###### page

See [configure page](#configure-page).

###### pagination

See [configure pagination](/templates/pagination/#configuration).

###### permalinks

See [Content Management](/content-management/urls/#permalinks).

###### pluralizeListTitles

(`bool`) Whether to pluralize automatic list titles. Applicable to section pages. Default is `true`.

###### publishDir

(`string`) The directory to where Hugo will write the final static site (the HTML files etc.). Default is `public`.

###### refLinksErrorLevel

(`string`) When using `ref` or `relref` to resolve page links and a link cannot be resolved, it will be logged with this log level. Valid values are `ERROR` (default) or `WARNING`. Any `ERROR` will fail the build (`exit -1`).  Default is `ERROR`.

###### refLinksNotFoundURL

(`string`) URL to be used as a placeholder when a page reference cannot be found in `ref` or `relref`. Is used as-is.

###### related

See [Related Content](/content-management/related/#configure-related-content).

###### relativeURLs

(`bool`) See [details](/content-management/urls/#relative-urls) before enabling this feature. Default is `false`.

###### renderSegments

{{< new-in 0.124.0 >}}

(`string slice`) A list of segments to render. If not set, everything will be rendered. This is more commonly set in a CLI flag, e.g. `hugo --renderSegments segment1,segment2`. The segment names must match the names in the [segments](#configure-segments) configuration.

###### removePathAccents

(`bool`) Removes [non-spacing marks](https://www.compart.com/en/unicode/category/Mn) from [composite characters](https://en.wikipedia.org/wiki/Precomposed_character) in content paths. Default is `false`.

```text
content/post/hügó.md → https://example.org/post/hugo/
```

###### sectionPagesMenu

See [Menus](/content-management/menus/#define-automatically).

###### security

See [Security Policy](/about/security/#security-policy).

###### segments

See [Segments](#configure-segments).

###### sitemap

Default [sitemap configuration](/templates/sitemap/#configuration).

###### summaryLength

(`int`) Applicable to automatic summaries, the approximate number of words to render when calling the [`Summary`] method on a `Page` object. Default is `70`.

[`Summary`]: /methods/page/summary/

###### taxonomies

See [Configure Taxonomies](/content-management/taxonomies#configure-taxonomies).

###### theme

See [module configuration](/hugo-modules/configuration/#module-configuration-imports) for how to import a theme.

###### themesDir

(`string`) The directory where Hugo reads the themes from. Default is `themes`.

###### timeout

(`string`) Timeout for generating page contents, specified as a [duration](https://pkg.go.dev/time#Duration) or in seconds. *Note:*&nbsp;this is used to bail out of recursive content generation. You might need to raise this limit if your pages are slow to generate (e.g., because they require large image processing or depend on remote contents). Default is `30s`.

###### timeZone

(`string`) The time zone (or location), e.g. `Europe/Oslo`, used to parse front matter dates without such information and in the [`time`] function. The list of valid values may be system dependent, but should include `UTC`, `Local`, and any location in the [IANA Time Zone database](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).

###### title

(`string`) Site title.

###### titleCaseStyle

(`string`) Default is `ap`. See [Configure Title Case](#configure-title-case).

###### uglyURLs

(`bool`) When enabled, creates URL of the form `/filename.html` instead of `/filename/`. Default is `false`.

###### watch

(`bool`) Watch filesystem for changes and recreate as needed. Default is `false`.

{{% note %}}
If you are developing your site on a \*nix machine, here is a handy shortcut for finding a configuration option from the command line:
```txt
cd ~/sites/yourhugosite
hugo config | grep emoji
```

which shows output like

```txt
enableemoji: true
```
{{% /note %}}

## Configure page

{{< new-in 0.133.0 >}}

These methods on a `Page` object navigate to the next or previous page within a page collection, relative to the current page:

- [Next](/methods/page/next/)
- [NextInSection](/methods/page/nextinsection/)
- [Prev](/methods/page/prev/)
- [PrevInSection](/methods/page/previnsection/)

Hugo determines the _next_ and _previous_ page by sorting a page collection according to this sorting hierarchy:

Field|Precedence|Sort direction
:--|:--|:--
[`weight`]|1|descending
[`date`]|2|descending
[`linkTitle`]|3|descending
[`path`]|4|descending

[`date`]: /methods/page/date/
[`weight`]: /methods/page/weight/
[`linkTitle`]: /methods/page/linktitle/
[`path`]: /methods/page/path/

The sort direction in the table above corresponds to these default site configuration values:

{{< code-toggle config=page />}}

To sort all fields in ascending order:

{{< code-toggle file=hugo >}}
[page]
  nextPrevInSectionSortOrder = 'asc'
  nextPrevSortOrder = 'asc'
{{< /code-toggle >}}

{{% note %}}
These settings do not apply to the [`Next`] or [`Prev`] methods on a `Pages` object.

[`Next`]: /methods/pages/next
[`Prev`]: /methods/pages/next
{{% /note %}}

## Configure build

The `build` configuration section contains global build-related configuration options.

{{< code-toggle config=build />}}

buildStats {{< new-in 0.115.1 >}}
: When enabled, creates a `hugo_stats.json` file in the root of your project. This file contains arrays of the `class` attributes, `id` attributes, and tags of every HTML element within your published site. Use this file as data source when [removing unused CSS] from your site. This process is also known as pruning, purging, or tree shaking.

[removing unused CSS]: /hugo-pipes/postprocess/#css-purging-with-postcss

Exclude `class` attributes, `id` attributes, or tags from `hugo_stats.json` with the `disableClasses`, `disableIDs`, and `disableTags` keys.

{{% note %}}
With v0.115.0 and earlier this feature was enabled by setting `writeStats` to `true`. Although still functional, the `writeStats` key will be deprecated in a future release.

Given that CSS purging is typically limited to production builds, place the `buildStats` object below [config/production].

[config/production]: /getting-started/configuration/#configuration-directory

Built for speed, there may be "false positive" detections (e.g., HTML elements that are not HTML elements) while parsing the published site. These "false positives" are infrequent and inconsequential.
{{% /note %}}

Due to the nature of partial server builds, new HTML entities are added while the server is running, but old values will not be removed until you restart the server or run a regular `hugo` build.

cachebusters
: See [Configure Cache Busters](#configure-cache-busters)

noJSConfigInAssets
: Turn off writing a `jsconfig.json` into your `/assets` folder with mapping of imports from running [js.Build](/hugo-pipes/js). This file is intended to help with intellisense/navigation inside code editors such as [VS Code](https://code.visualstudio.com/). Note that if you do not use `js.Build`, no file will be written.

useResourceCacheWhen
: When to use the cached resources in `/resources/_gen` for PostCSS and ToCSS. Valid values are `never`, `always` and `fallback`. The last value means that the cache will be tried if PostCSS/extended version is not available.

## Configure cache busters

{{< new-in 0.112.0 >}}

The `build.cachebusters` configuration option was added to support development using Tailwind 3.x's JIT compiler where a `build` configuration may look like this:

{{< code-toggle file=hugo >}}
[build]
  [build.buildStats]
    enable = true
  [[build.cachebusters]]
    source = "assets/watching/hugo_stats\\.json"
    target = "styles\\.css"
  [[build.cachebusters]]
    source = "(postcss|tailwind)\\.config\\.js"
    target = "css"
  [[build.cachebusters]]
    source = "assets/.*\\.(js|ts|jsx|tsx)"
    target = "js"
  [[build.cachebusters]]
    source = "assets/.*\\.(.*)$"
    target = "$1"
{{< /code-toggle >}}

When `buildStats` {{< new-in 0.115.1 >}} is enabled, Hugo writes a `hugo_stats.json` file on each build with HTML classes etc. that's used in the rendered output. Changes to this file will trigger a rebuild of the `styles.css` file. You also need to add `hugo_stats.json` to Hugo's server watcher. See [Hugo Starter Tailwind Basic](https://github.com/bep/hugo-starter-tailwind-basic) for a running example.

source
: A regexp matching file(s) relative to one of the virtual component directories in Hugo, typically `assets/...`.

target
: A regexp matching the keys in the resource cache that should be expired when `source` changes. You can use the matching regexp groups from `source` in the expression, e.g. `$1`.

## Configure server

This is only relevant when running `hugo server`, and it allows to set HTTP headers during development, which allows you to test out your Content Security Policy and similar. The configuration format matches [Netlify's](https://docs.netlify.com/routing/headers/#syntax-for-the-netlify-configuration-file) with slightly more powerful [Glob matching](https://github.com/gobwas/glob):

{{< code-toggle file=hugo >}}
[server]
[[server.headers]]
for = "/**"

[server.headers.values]
X-Frame-Options = "DENY"
X-XSS-Protection = "1; mode=block"
X-Content-Type-Options = "nosniff"
Referrer-Policy = "strict-origin-when-cross-origin"
Content-Security-Policy = "script-src localhost:1313"
{{< /code-toggle >}}

Since this is "development only", it may make sense to put it below the `development` environment:

{{< code-toggle file=config/development/server >}}
[[headers]]
for = "/**"

[headers.values]
X-Frame-Options = "DENY"
X-XSS-Protection = "1; mode=block"
X-Content-Type-Options = "nosniff"
Referrer-Policy = "strict-origin-when-cross-origin"
Content-Security-Policy = "script-src localhost:1313"
{{< /code-toggle >}}

You can also specify simple redirects rules for the server. The syntax is again similar to Netlify's.

Note that a `status` code of 200 will trigger a [URL rewrite](https://docs.netlify.com/routing/redirects/rewrites-proxies/), which is what you want in SPA situations, e.g:

{{< code-toggle file=config/development/server >}}
[[redirects]]
from = "/myspa/**"
to = "/myspa/"
status = 200
force = false
{{< /code-toggle >}}

Setting `force=true` will make a redirect even if there is existing content in the path. Note that before Hugo 0.76 `force` was the default behavior, but this is inline with how Netlify does it.

## 404 server error page {#_404-server-error-page}

Hugo will, by default, render all 404 errors when running `hugo server` with the `404.html` template. Note that if you have already added one or more redirects to your [server configuration](#configure-server), you need to add the 404 redirect explicitly, e.g:

{{< code-toggle file=config/development/server >}}
[[redirects]]
from   = "/**"
to     = "/404.html"
status = 404
{{< /code-toggle >}}

With a multilingual site, define the redirect for the default content language last:

{{< code-toggle file=config/development/server >}}
defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = false
[[redirects]]
from = '/fr/**'
to = '/fr/404.html'
status = 404

[[redirects]] # Default language must be last.
from = '/**'
to = '/404.html'
status = 404
{{< /code-toggle >}}

If you are serving the default content language from a subdirectory:

{{< code-toggle file=config/development/server >}}
defaultContentLanguage = 'en'
defaultContentLanguageInSubdir = true
[[redirects]]
from = '/fr/**'
to = '/fr/404.html'
status = 404

[[redirects]] # Default language must be last.
from = '/**'
to = '/en/404.html'
status = 404
{{< /code-toggle >}}


## Configure title case

By default, Hugo follows the capitalization rules published in the [Associated Press Stylebook] when creating automatic section titles, and when transforming strings with the [`strings.Title`] function.

Change this behavior by setting `titleCaseStyle` in your site configuration to any of the values below:

ap
: Use the capitalization rules published in the [Associated Press Stylebook].

chicago
: Use the capitalization rules published in the [Chicago Manual of Style].

go
: Capitalize the first letter of every word.

firstupper
: Capitalize the first letter of the first word.

none
: Disable transformation of automatic section titles, and disable the transformation performed by the `strings.Title` function. This is useful if you would prefer to manually capitalize section titles as needed, and to bypass opinionated theme usage of the `strings.Title` function.

[`strings.Title`]: /functions/strings/title/
[Associated Press Stylebook]: https://www.apstylebook.com/
[Chicago Manual of Style]: https://www.chicagomanualofstyle.org/home.html
[site configuration]: /getting-started/configuration/#configure-title-case

## Configuration environment variables

DART_SASS_BINARY
: (`string`) The absolute path to the Dart Sass executable. By default, Hugo searches for the executable in each of the paths in the `PATH` environment variable.

HUGO_ENVIRONMENT
: (`string`) Overrides the default [environment], typically one of `development`, `staging`, or `production`.

[environment]: /getting-started/glossary/#environment

HUGO_FILE_LOG_FORMAT
: (`string`) A format string for the file path, line number, and column number displayed when reporting errors, or when calling the `Position` method from a shortcode or Markdown render hook. Valid tokens are `:file`, `:line`, and `:col`. Default is `:file::line::col`.

{{< new-in 0.123.0 >}}

HUGO_MEMORYLIMIT
: (`int`) The maximum amount of system memory, in gigabytes, that Hugo can use while rendering your site. Default is 25% of total system memory.

HUGO_NUMWORKERMULTIPLIER
: (`int`) The number of workers used in parallel processing. Default is the number of logical CPUs.

## Configure with environment variables

In addition to the 3 configuration options already mentioned, configuration key-values can be defined through operating system environment variables.

For example, the following command will effectively set a website's title on Unix-like systems:

```txt
$ env HUGO_TITLE="Some Title" hugo
```

This is really useful if you use a service such as Netlify to deploy your site. Look at the Hugo docs [Netlify configuration file](https://github.com/gohugoio/hugoDocs/blob/master/netlify.toml) for an example.

{{% note %}}
Names must be prefixed with `HUGO_` and the configuration key must be set in uppercase when setting operating system environment variables.

To set configuration parameters, prefix the name with `HUGO_PARAMS_`
{{% /note %}}

If you are using snake_cased variable names, the above will not work. Hugo determines the delimiter to use by the first character after `HUGO`. This allows you to define environment variables on the form `HUGOxPARAMSxAPI_KEY=abcdefgh`, using any [allowed](https://stackoverflow.com/questions/2821043/allowed-characters-in-linux-environment-variable-names#:~:text=So%20names%20may%20contain%20any,not%20begin%20with%20a%20digit.) delimiter.

## Ignore content and data files when rendering

{{% note %}}
This works, but we recommend you use the newer and more powerful [includeFiles and excludeFiles](/hugo-modules/configuration/#module-configuration-mounts) mount options.
{{% /note %}}

To exclude specific files from the `content`, `data`, and `i18n` directories when rendering your site, set `ignoreFiles` to one or more regular expressions to match against the absolute file path.

To ignore files ending with `.foo` or `.boo`:

{{< code-toggle file=hugo >}}
ignoreFiles = ['\.foo$', '\.boo$']
{{< /code-toggle >}}

To ignore a file using the absolute file path:

{{< code-toggle file=hugo >}}
ignoreFiles = ['^/home/user/project/content/test\.md$']
{{< /code-toggle >}}

## Configure front matter

### Configure dates

Dates are important in Hugo, and you can configure how Hugo assigns dates to your content pages. You do this by adding a `frontmatter` section to your `hugo.toml`.

The default configuration is:

{{< code-toggle config=frontmatter />}}

If you, as an example, have a non-standard date parameter in some of your content, you can override the setting for `date`:

{{< code-toggle file=hugo >}}
[frontmatter]
date = ["myDate", ":default"]
{{< /code-toggle >}}

The `:default` is a shortcut to the default settings. The above will set `.Date` to the date value in `myDate` if present, if not we will look in `date`,`publishDate`, `lastmod` and pick the first valid date.

In the list to the right, values starting with ":" are date handlers with a special meaning (see below). The others are just names of date parameters (case insensitive) in your front matter configuration. Also note that Hugo have some built-in aliases to the above: `lastmod` => `modified`, `publishDate` => `pubdate`, `published` and `expiryDate` => `unpublishdate`. With that, as an example, using `pubDate` as a date in front matter, will, by default, be assigned to `.PublishDate`.

The special date handlers are:

`:fileModTime`
: Fetches the date from the content file's last modification timestamp.

An example:

{{< code-toggle file=hugo >}}
[frontmatter]
lastmod = ["lastmod", ":fileModTime", ":default"]
{{< /code-toggle >}}

The above will try first to extract the value for `.Lastmod` starting with the `lastmod` front matter parameter, then the content file's modification timestamp. The last, `:default` should not be needed here, but Hugo will finally look for a valid date in `:git`, `date` and then `publishDate`.

`:filename`
: Fetches the date from the content file's file name. For example, `2018-02-22-mypage.md` will extract the date `2018-02-22`. Also, if `slug` is not set, `mypage` will be used as the value for `.Slug`.

An example:

{{< code-toggle file=hugo >}}
[frontmatter]
date  = [":filename", ":default"]
{{< /code-toggle >}}

The above will try first to extract the value for `.Date` from the file name, then it will look in front matter parameters `date`, `publishDate` and lastly `lastmod`.

`:git`
: This is the Git author date for the last revision of this content file. This will only be set if `--enableGitInfo` is set or `enableGitInfo = true` is set in site configuration.

## Configure minify

See the [tdewolff/minify] project page for details.

[tdewolff/minify]: https://github.com/tdewolff/minify

Default configuration:

{{< code-toggle config=minify />}}

## Configure file caches

Since Hugo 0.52 you can configure more than just the `cacheDir`. This is the default configuration:

{{< code-toggle config=caches />}}

You can override any of these cache settings in your own `hugo.toml`.

### The keywords explained

cacheDir
: (`string`) See [Configure cacheDir](#configure-cachedir).

project
: (`string`) The base directory name of the current Hugo project. This means that, in its default setting, every project will have separated file caches, which means that when you do `hugo --gc` you will not touch files related to other Hugo projects running on the same PC.

resourceDir
: (`string`) This is the value of the `resourceDir` configuration option.

maxAge
: (`string`) This is the duration before a cache entry will be evicted, -1 means forever and 0 effectively turns that particular cache off. Uses Go's `time.Duration`, so valid values are `"10s"` (10 seconds), `"10m"` (10 minutes) and `"10h"` (10 hours).

dir
: (`string`) The absolute path to where the files for this cache will be stored. Allowed starting placeholders are `:cacheDir` and `:resourceDir` (see above).

## Configure cacheDir

This is the directory where Hugo by default will store its file caches. See [Configure File Caches](#configure-file-caches).

This can be set using the `cacheDir` config option or via the OS environment variable `HUGO_CACHEDIR`.

If this is not set, Hugo will use, in order of preference:

1. If running on Netlify: `/opt/build/cache/hugo_cache/`. This means that if you run your builds on Netlify, all caches configured with `:cacheDir` will be saved and restored on the next build. For other CI vendors, please read their documentation. For an CircleCI example, see [this configuration](https://github.com/bep/hugo-sass-test/blob/6c3960a8f4b90e8938228688bc49bdcdd6b2d99e/.circleci/config.yml).
1. In a `hugo_cache` directory below the OS user cache directory as defined by Go's [os.UserCacheDir](https://pkg.go.dev/os#UserCacheDir). On Unix systems, this is `$XDG_CACHE_HOME` as specified by [basedir-spec-latest](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) if non-empty, else `$HOME/.cache`. On MacOS, this is `$HOME/Library/Caches`. On Windows, this is`%LocalAppData%`. On Plan 9, this is `$home/lib/cache`. {{< new-in 0.116.0 >}}
1. In a  `hugo_cache_$USER` directory below the OS temp dir.

If you want to know the current value of `cacheDir`, you can run `hugo config`, e.g: `hugo config | grep cachedir`.

[`time`]: /functions/time/astime/
[`.Site.Params`]: /method/site/params/
[directory structure]: /getting-started/directory-structure/
[lookup order]: /templates/lookup-order/
[custom output formats]: /templates/output-formats/
[templates]: /templates/
[static-files]: /content-management/static-files/


## Configure HTTP cache

{{< new-in 0.127.0 >}}

Note that this configuration is currently only relevant when using the [resources.GetRemote] function.

The caching in Hugo is layered:

```goat {.w-40}
 .-----------.
|  dynacache  |
 '-----+-----'
       |
       v
 .----------.
| HTTP cache |
 '-----+----'
       |
       v
 .----------.
| file cache |
 '-----+----'
```

Dynacache
: A in memory LRU cache that gets evicted on changes, [Cache Buster](#configure-cache-busters) matches and in low memory situations.

HTTP Cache
: Enables HTTP cache behavior (RFC 9111) for remote resources. This works best for resources with properly set up HTTP cache headers. The HTTP cache uses the [file cache] to store and serve cached resources.

File Cache
: See [file cache].

The default HTTP cache disables everything:

{{< code-toggle config=HTTPCache />}}

caching
: Enabled RFC 9111 cache behavior _for_ a configured set of resources. Stale resources will be refreshed from the [file cache] even if their configured TTL isn't reached.

polling
: Enables polling _for_ a set of resources. Note that you can enable polling for resources even if HTTP caching is disabled. This setting is only used when in watch mode (e.g. `hugo server`). When a changed resource is detected, that change triggers a rebuild of pages using that resource.

[resources.GetRemote]: /functions/resources/getremote/
[file cache]: #configure-file-caches

## Configure segments

{{< new-in 0.124.0 >}}

{{% note %}}
The `segments` configuration is currently only used to configure partitioned rendering.
This feature is only about what gets rendered when, Hugo's entire object graph (sites and pages) is
always available.
{{% /note %}}

* Each segment consists of zero or more `exclude` filters and zero or more `include` filters.
* Each filter consists of one or more field Glob matchers.
* Each filter in a section (`exclude` or `include`) is ORed together, each matcher in a filter is ANDed together.

The fields that can be used in the filters are:

path
: The logical page [path].

lang
: The [page language].

kind
: The [kind] of the page.

output
: The [output format] of the page.

It is recommended to put coarse grained filters (e.g. for language and output format) in the excludes section, e.g.:

{{< code-toggle file=hugo >}}
[segments.segment1]
  [[segments.segment1.excludes]]
    lang = "n*"
  [[segments.segment1.excludes]]
    lang   = "en"
    output = "rss"
  [[segments.segment1.includes]]
    kind = "{home,term,taxonomy}"
  [[segments.segment1.includes]]
    path = "{/docs,/docs/**}"
{{< /code-toggle >}}

With the above you can render only the pages in `segment1` by configuring the [renderSegments](#rendersegments) or setting the `--renderSegments` flag:

```bash
hugo --renderSegments segment1
```

Multiple segments can be configured, and the `--renderSegments` flag can take a comma separated list of segments.

Some use cases for this feature:

* Splitting builds of big sites.
* Enable faster builds during development by only rendering a subset of the site.
* Partial rebuilds, e.g. render the home page and the "news section" every hour, render the entire site once a week.
* Render only e.g. the JSON output format to push to e.g. a search index.
  
[path]: /methods/page/path/
[page language]: /methods/page/language/
[kind]: /getting-started/glossary/#page-kind
[output format]: /getting-started/glossary/#output-format
[type]: /getting-started/glossary/#content-type
