---
title: Configure Hugo
linktitle: Configuration
description: How to configure your Hugo site.
categories: [getting started,fundamentals]
keywords: [configuration,toml,yaml,json]
menu:
  docs:
    parent: getting-started
    weight: 60
weight: 60
aliases: [/overview/source-directory/,/overview/configuration/]
toc: true
---

## Configuration File

Hugo uses the `config.toml`, `config.yaml`, or `config.json` (if found in the
site root) as the default site config file.

The user can choose to override that default with one or more site config files
using the command-line `--config` switch.

Examples:

```txt
hugo --config debugconfig.toml
hugo --config a.toml,b.toml,c.toml
```

{{% note %}}
Multiple site config files can be specified as a comma-separated string to the `--config` switch.
{{% /note %}}

## hugo.toml vs config.toml

In Hugo 0.110.0 we changed the default config base filename to `hugo`, e.g. `hugo.toml`. We will still look for `config.toml` etc., but we recommend you eventually rename it (but you need to wait if you want to support older Hugo versions). The main reason we're doing this is to make it easier for code editors and build tools to identify this as a Hugo configuration file and project.

{{< new-in "0.110.0" >}}

## Configuration Directory

In addition to using a single site config file, one can use the `configDir` directory (default to `config/`) to maintain easier organization and environment specific settings.

- Each file represents a configuration root object, such as `params.toml` for `[Params]`, `menu(s).toml` for `[Menu]`, `languages.toml` for `[Languages]` etc...
- Each file's content must be top-level, for example:

{{< code-toggle file="config" >}}
[Params]
  foo = "bar"
{{< /code-toggle >}}

{{< code-toggle file="params" >}}
foo = "bar"
{{< /code-toggle >}}

- Each directory holds a group of files containing settings unique to an environment.
- Files can be localized to become language specific.


```txt
├── config
│   ├── _default
│   │   ├── config.toml
│   │   ├── languages.toml
│   │   ├── menus.en.toml
│   │   ├── menus.zh.toml
│   │   └── params.toml
│   ├── production
│   │   ├── config.toml
│   │   └── params.toml
│   └── staging
│       ├── config.toml
│       └── params.toml
```

Considering the structure above, when running `hugo --environment staging`, Hugo will use every setting from `config/_default` and merge `staging`'s on top of those.

Let's take an example to understand this better. Let's say you are using Google Analytics for your website. This requires you to specify `googleAnalytics = "G-XXXXXXXX"` in `config.toml`. Now consider the following scenario:
- You don't want the Analytics code to be loaded in development i.e. in your `localhost`
- You want to use separate googleAnalytics IDs for your production & staging environments (say):
  - `G-PPPPPPPP` for production
  - `G-SSSSSSSS` for staging

This is how you need to configure your `config.toml` files considering the above scenario:
1. In `_default/config.toml` you don't need to mention `googleAnalytics` parameter at all. This ensures that no Google Analytics code is loaded in your development server i.e. when you run `hugo serve`. This works since, by default Hugo sets `Environment=development` when you run `hugo serve` which uses the config files from `_default` folder
2. In `production/config.toml` you just need to have one line:

    ```googleAnalytics = "G-PPPPPPPP"```

    You don't need to mention all other parameters like `title`, `baseURL`, `theme` etc. again in this config file. You need to mention only those parameters which are different or new for the production environment. This is due to the fact that Hugo is going to __merge__ this on top of `_default/config.toml`. Now when you run `hugo` (build command), by default hugo sets `Environment=production`, so the `G-PPPPPPPP` analytics code will be there in your production website
3. Similarly in `staging/config.toml` you just need to have one line:

    ```googleAnalytics = "G-SSSSSSSS"```

    Now you need to tell Hugo that you are using the staging environment. So your build command should be `hugo --environment staging` which will load the `G-SSSSSSSS` analytics code in your staging website

{{% note %}}
Default environments are __development__ with `hugo server` and __production__ with `hugo`.
{{%/ note %}}

## Merge Configuration from Themes

The configuration value for `_merge` can be one of:

none
: No merge.

shallow
: Only add values for new keys.

deep
: Add values for new keys, merge existing.

Note that you don't need to be so verbose as in the default setup below; a `_merge` value higher up will be inherited if not set.

{{< code-toggle config="mergeStrategy" skipHeader=true />}}

## All Configuration Settings

The following is the full list of Hugo-defined variables. Users may choose to override those values in their site
config file(s).

### archetypeDir

**Default value:** "archetypes"

The directory where Hugo finds archetype files (content templates). {{% module-mounts-note %}}

### assetDir

**Default value:** "assets"

The directory where Hugo finds asset files used in [Hugo Pipes](/hugo-pipes/). {{% module-mounts-note %}}

### baseURL

Hostname (and path) to the root, e.g. https://bep.is/

### build

See [Configure Build](#configure-build)

### buildDrafts (false)

**Default value:** false

Include drafts when building.

### buildExpired

**Default value:** false

Include content already expired.

### buildFuture

**Default value:** false

Include content with publishdate in the future.

### caches

See [Configure File Caches](#configure-file-caches)

### cascade

Pass down default configuration values (front matter) to pages in the content tree. The options in site config is the same as in page front matter, see [Front Matter Cascade](/content-management/front-matter#front-matter-cascade).

### canonifyURLs

**Default value:** false

Enable to turn relative URLs into absolute.

### cleanDestinationDir

**Default value:** false

When building, removes files from destination not found in static directories.

### contentDir

**Default value:** "content"

The directory from where Hugo reads content files. {{% module-mounts-note %}}

### copyright

**Default value:** ""

Copyright notice for your site, typically displayed in the footer.

### dataDir

**Default value:** "data"

The directory from where Hugo reads data files. {{% module-mounts-note %}}

### defaultContentLanguage

**Default value:** "en"

Content without language indicator will default to this language.

### defaultContentLanguageInSubdir

**Default value:**  false

Render the default content language in subdir, e.g. `content/en/`. The site root `/` will then redirect to `/en/`.

### disableAliases

**Default value:**  false

Will disable generation of alias redirects. Note that even if `disableAliases` is set, the aliases themselves are preserved on the page. The motivation with this is to be able to generate 301 redirects in an `.htaccess`, a Netlify `_redirects` file or similar using a custom output format.

### disableHugoGeneratorInject

**Default value:**  false

Hugo will, by default, inject a generator meta tag in the HTML head on the _home page only_. You can turn it off, but we would really appreciate if you don't, as this is a good way to watch Hugo's popularity on the rise.

### disableKinds

**Default value:**  []

Enable disabling of all pages of the specified *Kinds*. Allowed values in this list: `"page"`, `"home"`, `"section"`, `"taxonomy"`, `"term"`, `"RSS"`, `"sitemap"`, `"robotsTXT"`, `"404"`.

### disableLiveReload

**Default value:**  false

Disable automatic live reloading of browser window.

### disablePathToLower

**Default value:**  false

Do not convert the url/path to lowercase.

### enableEmoji

**Default value:**  false

Enable Emoji emoticons support for page content; see the [Emoji Cheat Sheet](https://www.webpagefx.com/tools/emoji-cheat-sheet/).

### enableGitInfo

**Default value:**  false

Enable `.GitInfo` object for each page (if the Hugo site is versioned by Git). This will then update the `Lastmod` parameter for each page using the last git commit date for that content file.

### enableInlineShortcodes

**Default value:**  false

Enable inline shortcode support. See [Inline Shortcodes](/templates/shortcode-templates/#inline-shortcodes).

### enableMissingTranslationPlaceholders

**Default value:**  false

Show a placeholder instead of the default value or an empty string if a translation is missing.

### enableRobotsTXT

**Default value:**  false

Enable generation of `robots.txt` file.

### frontmatter

See [Front matter Configuration](#configure-front-matter).

### googleAnalytics

**Default value:**  ""

Google Analytics tracking ID.

### hasCJKLanguage

**Default value:** false

If true, auto-detect Chinese/Japanese/Korean Languages in the content. This will make `.Summary` and `.WordCount` behave correctly for CJK languages.

### imaging

See [Image Processing Config](/content-management/image-processing/#imaging-configuration).

### languageCode

**Default value:**  ""

A language tag as defined by [RFC 5646](https://datatracker.ietf.org/doc/html/rfc5646). This value is used to populate:

- The `<language>` element in the internal [RSS template](https://github.com/gohugoio/hugo/blob/master/tpl/tplimpl/embedded/templates/_default/rss.xml)
- The `lang` attribute of the `<html>` element in the internal [alias template](https://github.com/gohugoio/hugo/blob/master/tpl/tplimpl/embedded/templates/alias.html)

### languages

See [Configure Languages](/content-management/multilingual/#configure-languages).

### disableLanguages

See [Disable a Language](/content-management/multilingual/#disable-a-language)

### markup

See [Configure Markup](/getting-started/configuration-markup).

### mediaTypes

See [Configure Media Types](/templates/output-formats/#media-types).

### menus

See [Menus](/content-management/menus/#define-in-site-configuration).

### minify

See [Configure Minify](#configure-minify)

### module

Module config see [Module Config](/hugo-modules/configuration/).

### newContentEditor

**Default value:** ""

The editor to use when creating new content.

### noChmod

**Default value:** false

Don't sync permission mode of files.

### noTimes

**Default value:** false

Don't sync modification time of files.

### outputFormats

See [Configure Output Formats](#configure-additional-output-formats).

### paginate

**Default value:** 10

Default number of elements per page in [pagination](/templates/pagination/).

### paginatePath

**Default value:** "page"

The path element used during pagination (`https://example.com/page/2`).

### permalinks

See [Content Management](/content-management/urls/#permalinks).

### pluralizeListTitles

**Default value:** true

Pluralize titles in lists.

### publishDir

**Default value:** "public"

The directory to where Hugo will write the final static site (the HTML files etc.).

### related

: See [Related Content](/content-management/related/#configure-related-content).

### relativeURLs

**Default value:** false

Enable this to make all relative URLs relative to content root. Note that this does not affect absolute URLs.

### refLinksErrorLevel

**Default value:** "ERROR"

When using `ref` or `relref` to resolve page links and a link cannot be resolved, it will be logged with this log level. Valid values are `ERROR` (default) or `WARNING`. Any `ERROR` will fail the build (`exit -1`).

### refLinksNotFoundURL

URL to be used as a placeholder when a page reference cannot be found in `ref` or `relref`. Is used as-is.

### removePathAccents

**Default value:** false

Removes [non-spacing marks](https://www.compart.com/en/unicode/category/Mn) from [composite characters](https://en.wikipedia.org/wiki/Precomposed_character) in content paths.

```text
content/post/hügó.md --> https://example.org/post/hugo/
```

### rssLimit

**Default value:** -1 (unlimited)

Maximum number of items in the RSS feed.

### sectionPagesMenu

See [Menus](/content-management/menus/#define-in-site-configuration).

### security

See [Security Policy](/about/security-model/#security-policy)

### sitemap

Default [sitemap configuration](/templates/sitemap-template/#configuration).

### summaryLength

**Default value:** 70

The length of text in words to show in a [`.Summary`](/content-management/summaries/#automatic-summary-splitting).

### taxonomies

See [Configure Taxonomies](/content-management/taxonomies#configure-taxonomies).

### theme

: See [Module Config](/hugo-modules/configuration/#module-config-imports) for how to import a theme.

### themesDir

**Default value:**  "themes"

The directory where Hugo reads the themes from.

### timeout

**Default value:** "30s"

Timeout for generating page contents, specified as a [duration](https://pkg.go.dev/time#Duration) or in milliseconds. *Note:*&nbsp;this is used to bail out of recursive content generation. You might need to raise this limit if your pages are slow to generate (e.g., because they require large image processing or depend on remote contents).

### timeZone

The time zone (or location), e.g. `Europe/Oslo`,  used to parse front matter dates without such information and in the [`time` function](/functions/time/). The list of valid values may be system dependent, but should include `UTC`, `Local`, and any location in the [IANA Time Zone database](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).

### title

Site title.

### titleCaseStyle

**Default value:**  "AP"

See [Configure Title Case](#configure-title-case)

### uglyURLs

**Default value:** false

When enabled, creates URL of the form `/filename.html` instead of `/filename/`.

### watch

**Default value:** false

Watch filesystem for changes and recreate as needed.

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

## Configure Build

The `build` configuration section contains global build-related configuration options.

{{< code-toggle file="config">}}
[build]
useResourceCacheWhen="fallback"
writeStats = false
noJSConfigInAssets = false
{{< /code-toggle >}}


useResourceCacheWhen
: When to use the cached resources in `/resources/_gen` for PostCSS and ToCSS. Valid values are `never`, `always` and `fallback`. The last value means that the cache will be tried if PostCSS/extended version is not available.

writeStats
: When enabled, a file named `hugo_stats.json` will be written to your project root with some aggregated data about the build, e.g. list of HTML entities published to be used to do [CSS pruning](/hugo-pipes/postprocess/#css-purging-with-postcss). If you're only using this for the production build, you should consider placing it below [config/production](/getting-started/configuration/#configuration-directory). It's also worth mentioning that, due to the nature of the partial server builds, new HTML entities will be added when you add or change them while the server is running, but the old values will not be removed until you restart the server or run a regular `hugo` build.

**Note** that the prime use case for this is purging of unused CSS; it is built for speed and there may be false positives (e.g., detection of HTML elements that are not HTML elements).

noJSConfigInAssets
: Turn off writing a `jsconfig.json` into your `/assets` folder with mapping of imports from running [js.Build](https://gohugo.io/hugo-pipes/js). This file is intended to help with intellisense/navigation inside code editors such as [VS Code](https://code.visualstudio.com/). Note that if you do not use `js.Build`, no file will be written.

## Configure Server

This is only relevant when running `hugo server`, and it allows to set HTTP headers during development, which allows you to test out your Content Security Policy and similar. The configuration format matches [Netlify's](https://docs.netlify.com/routing/headers/#syntax-for-the-netlify-configuration-file) with slightly more powerful [Glob matching](https://github.com/gobwas/glob):


{{< code-toggle file="config">}}
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


{{< code-toggle file="config/development/server">}}
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

{{< code-toggle file="config/development/server">}}
[[redirects]]
from = "/myspa/**"
to = "/myspa/"
status = 200
force = false
{{< /code-toggle >}}

Setting `force=true` will make a redirect even if there is existing content in the path. Note that before Hugo 0.76  `force` was the default behavior, but this is inline with how Netlify does it.

## 404 Server Error Page {#_404-server-error-page}

{{< new-in "0.103.0" >}}

Hugo will, by default, render all 404 errors when running `hugo server` with the `404.html` template. Note that if you have already added one or more redirects to your [Server Config](#configure-server), you need to add the 404 redirect explicitly, e.g:

{{< code-toggle file="config/development/server" copy=false >}}
[[redirects]]
from   = "/**"
to     = "/404.html"
status = 404
{{< /code-toggle >}}

## Configure Title Case

Set `titleCaseStyle` to specify the title style used by the [title](/functions/title/) template function and the automatic section titles in Hugo. It defaults to [AP Stylebook](https://www.apstylebook.com/) for title casing, but you can also set it to `Chicago` or `Go` (every word starts with a capital letter).

## Configuration Environment Variables

HUGO_NUMWORKERMULTIPLIER
: Can be set to increase or reduce the number of workers used in parallel processing in Hugo. If not set, the number of logical CPUs will be used.

## Configuration Lookup Order

Similar to the template [lookup order], Hugo has a default set of rules for searching for a configuration file in the root of your website's source directory as a default behavior:

1. `./config.toml`
2. `./config.yaml`
3. `./config.json`

In your `config` file, you can direct Hugo as to how you want your website rendered, control your website's menus, and arbitrarily define site-wide parameters specific to your project.


## Example Configuration

The following is a typical example of a configuration file. The values nested under `params:` will populate the [`.Site.Params`] variable for use in [templates]:

{{< code-toggle file="config">}}
baseURL: "https://yoursite.example.com/"
title: "My Hugo Site"
permalinks:
  posts: /:year/:month/:title/
params:
  Subtitle: "Hugo is Absurdly Fast!"
  AuthorName: "Jon Doe"
  GitHubUser: "spf13"
  ListOfFoo:
    - "foo1"
    - "foo2"
  SidebarRecentLimit: 5
{{< /code-toggle >}}

## Configure with Environment Variables

In addition to the 3 config options already mentioned, configuration key-values can be defined through operating system environment variables.

For example, the following command will effectively set a website's title on Unix-like systems:

```txt
$ env HUGO_TITLE="Some Title" hugo
```

This is really useful if you use a service such as Netlify to deploy your site. Look at the Hugo docs [Netlify configuration file](https://github.com/gohugoio/hugoDocs/blob/master/netlify.toml) for an example.

{{% note "Setting Environment Variables" %}}
Names must be prefixed with `HUGO_` and the configuration key must be set in uppercase when setting operating system environment variables.

To set config params, prefix the name with `HUGO_PARAMS_`
{{% /note %}}

If you are using snake_cased variable names, the above will not work. Hugo determines the delimiter to use by the first character after `HUGO`. This allows you to define environment variables on the form `HUGOxPARAMSxAPI_KEY=abcdefgh`, using any [allowed](https://stackoverflow.com/questions/2821043/allowed-characters-in-linux-environment-variable-names#:~:text=So%20names%20may%20contain%20any,not%20begin%20with%20a%20digit.) delimiter.

{{< todo >}}
Test and document setting params via JSON env var.
{{< /todo >}}

## Ignore Content and Data Files when Rendering

**Note:** This works, but we recommend you use the newer and more powerful [includeFiles and excludeFiles](https://gohugo.io/hugo-modules/configuration/#module-config-mounts) mount options.

To exclude specific files from the `content` and `data` directories when rendering your site, set `ignoreFiles` to one or more regular expressions to match against the absolute file path.

To ignore files ending with `.foo` or `.boo`:

{{< code-toggle copy=false >}}
ignoreFiles = ['\.foo$', '\.boo$']
{{< /code-toggle >}}

To ignore a file using the absolute file path:

{{< code-toggle copy=false >}}
ignoreFiles = ['^/home/user/project/content/test\.md$']
{{< /code-toggle >}}

## Configure Front Matter

### Configure Dates

Dates are important in Hugo, and you can configure how Hugo assigns dates to your content pages. You do this by adding a `frontmatter` section to your `config.toml`.

The default configuration is:

{{< code-toggle file="config" >}}
[frontmatter]
date = ["date", "publishDate", "lastmod"]
lastmod = [":git", "lastmod", "date", "publishDate"]
publishDate = ["publishDate", "date"]
expiryDate = ["expiryDate"]
{{< /code-toggle >}}

If you, as an example, have a non-standard date parameter in some of your content, you can override the setting for `date`:

{{< code-toggle file="config" >}}
[frontmatter]
date = ["myDate", ":default"]
{{< /code-toggle >}}

The `:default` is a shortcut to the default settings. The above will set `.Date` to the date value in `myDate` if present, if not we will look in `date`,`publishDate`, `lastmod` and pick the first valid date.

In the list to the right, values starting with ":" are date handlers with a special meaning (see below). The others are just names of date parameters (case insensitive) in your front matter configuration.  Also note that Hugo have some built-in aliases to the above: `lastmod` => `modified`, `publishDate` => `pubdate`, `published` and `expiryDate` => `unpublishdate`. With that, as an example, using `pubDate` as a date in front matter, will, by default, be assigned to `.PublishDate`.

The special date handlers are:


`:fileModTime`
: Fetches the date from the content file's last modification timestamp.

An example:

{{< code-toggle file="config" >}}
[frontmatter]
lastmod = ["lastmod", ":fileModTime", ":default"]
{{< /code-toggle >}}


The above will try first to extract the value for `.Lastmod` starting with the `lastmod` front matter parameter, then the content file's modification timestamp. The last, `:default` should not be needed here, but Hugo will finally look for a valid date in `:git`, `date` and then `publishDate`.


`:filename`
: Fetches the date from the content file's filename. For example, `2018-02-22-mypage.md` will extract the date `2018-02-22`. Also, if `slug` is not set, `mypage` will be used as the value for `.Slug`.

An example:

{{< code-toggle file="config" >}}
[frontmatter]
date  = [":filename", ":default"]
{{< /code-toggle >}}

The above will try first to extract the value for `.Date` from the filename, then it will look in front matter parameters `date`, `publishDate` and lastly `lastmod`.


`:git`
: This is the Git author date for the last revision of this content file. This will only be set if `--enableGitInfo` is set or `enableGitInfo = true` is set in site config.

## Configure Additional Output Formats

Hugo v0.20 introduced the ability to render your content to multiple output formats (e.g., to JSON, AMP html, or CSV). See [Output Formats] for information on how to add these values to your Hugo project's configuration file.

## Configure Minify

Default configuration:

{{< code-toggle config="minify" />}}

## Configure File Caches

Since Hugo 0.52 you can configure more than just the `cacheDir`. This is the default configuration:

{{< code-toggle >}}
[caches]
[caches.getjson]
dir = ":cacheDir/:project"
maxAge = -1
[caches.getcsv]
dir = ":cacheDir/:project"
maxAge = -1
[caches.getresource]
dir = ":cacheDir/:project"
maxAge = -1
[caches.images]
dir = ":resourceDir/_gen"
maxAge = -1
[caches.assets]
dir = ":resourceDir/_gen"
maxAge = -1
[caches.modules]
dir = ":cacheDir/modules"
maxAge = -1
{{< /code-toggle >}}

You can override any of these cache settings in your own `config.toml`.

### The keywords explained

`:cacheDir`
: This is the value of the `cacheDir` config option if set (can also be set via OS env variable `HUGO_CACHEDIR`). It will fall back to `/opt/build/cache/hugo_cache/` on Netlify, or a `hugo_cache` directory below the OS temp dir for the others. This means that if you run your builds on Netlify, all caches configured with `:cacheDir` will be saved and restored on the next build. For other CI vendors, please read their documentation. For an CircleCI example, see [this configuration](https://github.com/bep/hugo-sass-test/blob/6c3960a8f4b90e8938228688bc49bdcdd6b2d99e/.circleci/config.yml).

`:project`
: The base directory name of the current Hugo project. This means that, in its default setting, every project will have separated file caches, which means that when you do `hugo --gc` you will not touch files related to other Hugo projects running on the same PC.

`:resourceDir`
: This is the value of the `resourceDir` config option.

maxAge
: This is the duration before a cache entry will be evicted, -1 means forever and 0 effectively turns that particular cache off. Uses Go's `time.Duration`, so valid values are `"10s"` (10 seconds), `"10m"` (10 minutes) and `"10h"` (10 hours).

dir
: The absolute path to where the files for this cache will be stored. Allowed starting placeholders are `:cacheDir` and `:resourceDir` (see above).

## Configuration Format Specs

- [TOML Spec][toml]
- [YAML Spec][yaml]
- [JSON Spec][json]

[`.Site.Params`]: /variables/site/
[directory structure]: /getting-started/directory-structure
[json]: https://www.ecma-international.org/publications/files/ECMA-ST/ECMA-404.pdf "Specification for JSON, JavaScript Object Notation"
[lookup order]: /templates/lookup-order/
[Output Formats]: /templates/output-formats/
[templates]: /templates/
[toml]: https://github.com/toml-lang/toml
[yaml]: https://yaml.org/spec/
[static-files]: /content-management/static-files/
