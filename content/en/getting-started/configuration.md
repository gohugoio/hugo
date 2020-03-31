---
title: Configure Hugo
linktitle: Configuration
description: How to configure your Hugo site.
date: 2013-07-01
publishdate: 2017-01-02
lastmod: 2017-03-05
categories: [getting started,fundamentals]
keywords: [configuration,toml,yaml,json]
menu:
  docs:
    parent: "getting-started"
    weight: 60
weight: 60
sections_weight: 60
draft: false
aliases: [/overview/source-directory/,/overview/configuration/]
toc: true
---


## Configuration File

Hugo uses the `config.toml`, `config.yaml`, or `config.json` (if found in the
site root) as the default site config file.

The user can choose to override that default with one or more site config files
using the command line `--config` switch.

Examples:

```
hugo --config debugconfig.toml
hugo --config a.toml,b.toml,c.toml
```

{{% note %}}
Multiple site config files can be specified as a comma-separated string to the `--config` switch.
{{% /note %}}

{{< todo >}}TODO: distinct config.toml and others (the root object files){{< /todo >}}

## Configuration Directory

In addition to using a single site config file, one can use the `configDir` directory (default to `config/`) to maintain easier organization and environment specific settings.

- Each file represents a configuration root object, such as `Params`, `Menus`, `Languages` etc...
- Each directory holds a group of files containing settings unique to an environment.
- Files can be localized to become language specific.


```
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

Considering the structure above, when running `hugo --environment staging`, Hugo will use every settings from `config/_default` and merge `staging`'s on top of those.
{{% note %}}
Default environments are __development__ with `hugo serve` and __production__ with `hugo`.
{{%/ note %}}
## All Configuration Settings

The following is the full list of Hugo-defined variables with their default
value in parentheses. Users may choose to override those values in their site
config file(s).

archetypeDir ("archetypes")
: The directory where Hugo finds archetype files (content templates). {{% module-mounts-note %}}

assetDir ("assets")
: The directory where Hugo finds asset files used in [Hugo Pipes](/hugo-pipes/). {{% module-mounts-note %}}

baseURL
: Hostname (and path) to the root, e.g. https://bep.is/

blackfriday
: See [Configure Blackfriday](/getting-started/configuration-markup#blackfriday)

build
: See [Configure Build](#configure-build)

buildDrafts (false)
: Include drafts when building.

buildExpired  (false)
: Include content already expired.

buildFuture (false)
: Include content with publishdate in the future.

caches
: See [Configure File Caches](#configure-file-caches)

canonifyURLs (false)
: Enable to turn relative URLs into absolute.

contentDir ("content")
: The directory from where Hugo reads content files. {{% module-mounts-note %}}

dataDir ("data")
: The directory from where Hugo reads data files. {{% module-mounts-note %}}

defaultContentLanguage ("en")
: Content without language indicator will default to this language.

defaultContentLanguageInSubdir (false)
: Render the default content language in subdir, e.g. `content/en/`. The site root `/` will then redirect to `/en/`.

disableAliases (false)
: Will disable generation of alias redirects. Note that even if `disableAliases` is set, the aliases themselves are preserved on the page. The motivation with this is to be able to generate 301 redirects in an `.htaccess`, a Netlify `_redirects` file or similar using a custom output format.

disableHugoGeneratorInject (false)
: Hugo will, by default, inject a generator meta tag in the HTML head on the _home page only_. You can turn it off, but we would really appreciate if you don't, as this is a good way to watch Hugo's popularity on the rise.

disableKinds ([])
: Enable disabling of all pages of the specified *Kinds*. Allowed values in this list: `"page"`, `"home"`, `"section"`, `"taxonomy"`, `"taxonomyTerm"`, `"RSS"`, `"sitemap"`, `"robotsTXT"`, `"404"`.

disableLiveReload (false)
: Disable automatic live reloading of browser window.

disablePathToLower (false)
: Do not convert the url/path to lowercase.

enableEmoji (false)
: Enable Emoji emoticons support for page content; see the [Emoji Cheat Sheet](https://www.webpagefx.com/tools/emoji-cheat-sheet/).

enableGitInfo (false)
: Enable `.GitInfo` object for each page (if the Hugo site is versioned by Git). This will then update the `Lastmod` parameter for each page using the last git commit date for that content file.

enableInlineShortcodes (false)
: Enable inline shortcode support. See [Inline Shortcodes](/templates/shortcode-templates/#inline-shortcodes).

enableMissingTranslationPlaceholders (false)
: Show a placeholder instead of the default value or an empty string if a translation is missing.

enableRobotsTXT (false)
: Enable generation of `robots.txt` file.

frontmatter

: See [Front matter Configuration](#configure-front-matter).

footnoteAnchorPrefix ("")
: Prefix for footnote anchors.

footnoteReturnLinkContents ("")
: Text to display for footnote return links.

googleAnalytics ("")
: Google Analytics tracking ID.

hasCJKLanguage (false)
: If true, auto-detect Chinese/Japanese/Korean Languages in the content. This will make `.Summary` and `.WordCount` behave correctly for CJK languages.

imaging
: See [Image Processing Config](/content-management/image-processing/#image-processing-config).

languages
: See [Configure Languages](/content-management/multilingual/#configure-languages).

languageCode ("")
: The site's language code. It is used in the default [RSS template](/templates/rss/#configure-rss) and can be useful for [multi-lingual sites](/content-management/multilingual/#configure-multilingual-multihost).

languageName ("")
: The site's language name.

disableLanguages
: See [Disable a Language](/content-management/multilingual/#disable-a-language)

layoutDir ("layouts")
: The directory from where Hugo reads layouts (templates).

log (false)
: Enable logging.

logFile ("")
: Log File path (if set, logging enabled automatically).

markup
: See [Configure Markup](/getting-started/configuration-markup).{{< new-in "0.60.0" >}}

menu
: See [Add Non-content Entries to a Menu](/content-management/menus/#add-non-content-entries-to-a-menu).

minify
: See [Configure Minify](#configure-minify)

module
: Module config see [Module Config](/hugo-modules/configuration/).{{< new-in "0.56.0" >}}

newContentEditor ("")
: The editor to use when creating new content.

noChmod (false)
: Don't sync permission mode of files.

noTimes (false)
: Don't sync modification time of files.

paginate (10)
: Default number of elements per page in [pagination](/templates/pagination/).

paginatePath ("page")
: The path element used during pagination (https://example.com/page/2).

permalinks
: See [Content Management](/content-management/urls/#permalinks).

pluralizeListTitles (true)
: Pluralize titles in lists.

publishDir ("public")
: The directory to where Hugo will write the final static site (the HTML files etc.).

related
: See [Related Content](/content-management/related/#configure-related-content).{{< new-in "0.27" >}}

relativeURLs (false)
: Enable this to make all relative URLs relative to content root. Note that this does not affect absolute URLs.

refLinksErrorLevel ("ERROR")
: When using `ref` or `relref` to resolve page links and a link cannot resolved, it will be logged with this logg level. Valid values are `ERROR` (default) or `WARNING`. Any `ERROR` will fail the build (`exit -1`).

refLinksNotFoundURL
: URL to be used as a placeholder when a page reference cannot be found in `ref` or `relref`. Is used as-is.

rssLimit (unlimited)
: Maximum number of items in the RSS feed.

sectionPagesMenu ("")
: See ["Section Menu for Lazy Bloggers"](/templates/menu-templates/#section-menu-for-lazy-bloggers).

sitemap
: Default [sitemap configuration](/templates/sitemap-template/#configure-sitemap-xml).

staticDir ("static")
: A directory or a list of directories from where Hugo reads [static files][static-files]. {{% module-mounts-note %}}

summaryLength (70)
: The length of text in words to show in a [`.Summary`](/content-management/summaries/#hugo-defined-automatic-summary-splitting).

taxonomies
: See [Configure Taxonomies](/content-management/taxonomies#configure-taxonomies).

theme ("")
: Theme to use (located by default in `/themes/THEMENAME/`).

themesDir ("themes")
: The directory where Hugo reads the themes from.

timeout (10000)
: Timeout for generating page contents, in milliseconds (defaults to 10&nbsp;seconds). *Note:* this is used to bail out of recursive content generation, if your pages are slow to generate (e.g., because they require large image processing or depend on remote contents) you might need to raise this limit.

title ("")
: Site title.

titleCaseStyle ("AP")
: See [Configure Title Case](#configure-title-case)

uglyURLs (false)
: When enabled, creates URL of the form `/filename.html` instead of `/filename/`.

verbose (false)
: Enable verbose output.

verboseLog (false)
: Enable verbose logging.

watch (false)
: Watch filesystem for changes and recreate as needed.

{{% note %}}
If you are developing your site on a \*nix machine, here is a handy shortcut for finding a configuration option from the command line:
```
cd ~/sites/yourhugosite
hugo config | grep emoji
```

which shows output like

```
enableemoji: true
```
{{% /note %}}

## Configure Build

{{< new-in "0.66.0" >}}

The `build` configuration section contains global build-realated configuration options.

{{< code-toggle file="config">}}
[build]
useResourceCacheWhen="fallback"
{{< /code-toggle >}}


useResourceCacheWhen
: When to use the cached resources in `/resources/_gen` for PostCSS and ToCSS. Valid values are `never`, `always` and `fallback`. The last value means that the cache will be tried if PostCSS/extended version is not available.

## Configure Server

{{< new-in "0.67.0" >}}

This is only relevant when running `hugo server`, and it allows to set HTTP headers during development, which allows you to test out your Content Security Policy and similar. The configuration format matches [Netlify's](https://docs.netlify.com/routing/headers/#syntax-for-the-netlify-configuration-file) with slighly more powerful [Glob matching](https://github.com/gobwas/glob):


{{< code-toggle file="config">}}
[server]
[[server.headers]]
for = "/**.html"

[server.headers.values]
X-Frame-Options = "DENY"
X-XSS-Protection = "1; mode=block"
X-Content-Type-Options = "nosniff"
Referrer-Policy = "strict-origin-when-cross-origin"
Content-Security-Policy = "script-src localhost:1313"
{{< /code-toggle >}}

Since this is is "devlopment only", it may make sense to put it below the `development` environment:


{{< code-toggle file="config/development/server">}}
[[headers]]
for = "/**.html"

[headers.values]
X-Frame-Options = "DENY"
X-XSS-Protection = "1; mode=block"
X-Content-Type-Options = "nosniff"
Referrer-Policy = "strict-origin-when-cross-origin"
Content-Security-Policy = "script-src localhost:1313"
{{< /code-toggle >}}




## Configure Title Case

Set `titleCaseStyle` to specify the title style used by the [title](/functions/title/) template function and the automatic section titles in Hugo. It defaults to [AP Stylebook](https://www.apstylebook.com/) for title casing, but you can also set it to `Chicago` or `Go` (every word starts with a capital letter).

## Configuration Environment Variables

HUGO_NUMWORKERMULTIPLIER
: Can be set to increase or reduce the number of workers used in parallel processing in Hugo. If not set, the number of logical CPUs will be used.

## Configuration Lookup Order

Similar to the template [lookup order][], Hugo has a default set of rules for searching for a configuration file in the root of your website's source directory as a default behavior:

1. `./config.toml`
2. `./config.yaml`
3. `./config.json`

In your `config` file, you can direct Hugo as to how you want your website rendered, control your website's menus, and arbitrarily define site-wide parameters specific to your project.


## Example Configuration

The following is a typical example of a configuration file. The values nested under `params:` will populate the [`.Site.Params`][] variable for use in [templates][]:

{{< code-toggle file="config">}}
baseURL: "https://yoursite.example.com/"
title: "My Hugo Site"
footnoteReturnLinkContents: "↩"
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

```
$ env HUGO_TITLE="Some Title" hugo
```

This is really useful if you use a service such as Netlify to deploy your site. Look at the Hugo docs [Netlify configuration file](https://github.com/gohugoio/hugoDocs/blob/master/netlify.toml) for an example.

{{% note "Setting Environment Variables" %}}
Names must be prefixed with `HUGO_` and the configuration key must be set in uppercase when setting operating system environment variables.

To set config params, prefix the name with `HUGO_PARAMS_`
{{% /note %}}

{{< todo >}}
Test and document setting params via JSON env var.
{{< /todo >}}

## Ignore Content Files When Rendering

The following statement inside `./config.toml` will cause Hugo to ignore content files ending with `.foo` and `.boo` when rendering:

```
ignoreFiles = [ "\\.foo$", "\\.boo$" ]
```

The above is a list of regular expressions. Note that the backslash (`\`) character is escaped in this example to keep TOML happy.

## Configure Front Matter

### Configure Dates

Dates are important in Hugo, and you can configure how Hugo assigns dates to your content pages. You do this by adding a `frontmatter` section to your `config.toml`.


The default configuration is:

```toml
[frontmatter]
date = ["date", "publishDate", "lastmod"]
lastmod = [":git", "lastmod", "date", "publishDate"]
publishDate = ["publishDate", "date"]
expiryDate = ["expiryDate"]
```

If you, as an example, have a non-standard date parameter in some of your content, you can override the setting for `date`:

 ```toml
[frontmatter]
date = ["myDate", ":default"]
```

The `:default` is a shortcut to the default settings. The above will set `.Date` to the date value in `myDate` if present, if not we will look in `date`,`publishDate`, `lastmod` and pick the first valid date.

In the list to the right, values starting with ":" are date handlers with a special meaning (see below). The others are just names of date parameters (case insensitive) in your front matter configuration.  Also note that Hugo have some built-in aliases to the above: `lastmod` => `modified`, `publishDate` => `pubdate`, `published` and `expiryDate` => `unpublishdate`. With that, as an example, using `pubDate` as a date in front matter, will, by default, be assigned to `.PublishDate`.

The special date handlers are:


`:fileModTime`
: Fetches the date from the content file's last modification timestamp.

An example:

 ```toml
[frontmatter]
lastmod = ["lastmod", ":fileModTime", ":default"]
```


The above will try first to extract the value for `.Lastmod` starting with the `lastmod` front matter parameter, then the content file's modification timestamp. The last, `:default` should not be needed here, but Hugo will finally look for a valid date in `:git`, `date` and then `publishDate`.


`:filename`
: Fetches the date from the content file's filename. For example, `2018-02-22-mypage.md` will extract the date `2018-02-22`. Also, if `slug` is not set, `mypage` will be used as the value for `.Slug`.

An example:

```toml
[frontmatter]
date  = [":filename", ":default"]
```

The above will try first to extract the value for `.Date` from the filename, then it will look in front matter parameters `date`, `publishDate` and lastly `lastmod`.


`:git`
: This is the Git author date for the last revision of this content file. This will only be set if `--enableGitInfo` is set or `enableGitInfo = true` is set in site config.

## Configure Additional Output Formats

Hugo v0.20 introduced the ability to render your content to multiple output formats (e.g., to JSON, AMP html, or CSV). See [Output Formats][] for information on how to add these values to your Hugo project's configuration file.

## Configure Minify

{{< new-in "0.68.0" >}}

Default configuration:

{{< code-toggle config="minify" />}}

## Configure File Caches

Since Hugo 0.52 you can configure more than just the `cacheDir`. This is the default configuration:

```toml
[caches]
[caches.getjson]
dir = ":cacheDir/:project"
maxAge = -1
[caches.getcsv]
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
```

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

* [TOML Spec][toml]
* [YAML Spec][yaml]
* [JSON Spec][json]

[`.Site.Params`]: /variables/site/
[directory structure]: /getting-started/directory-structure
[json]: https://www.ecma-international.org/publications/files/ECMA-ST/ECMA-404.pdf "Specification for JSON, JavaScript Object Notation"
[lookup order]: /templates/lookup-order/
[Output Formats]: /templates/output-formats/
[templates]: /templates/
[toml]: https://github.com/toml-lang/toml
[yaml]: https://yaml.org/spec/
[static-files]: /content-management/static-files/
