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

## All Configuration Settings

The following is the full list of Hugo-defined variables with their default
value in parentheses. Users may choose to override those values in their site
config file(s).

archetypeDir ("archetypes")
: The directory where Hugo finds archetype files (content templates).

assetDir ("assets")
: The directory where Hugo finds asset files used in [Hugo Pipes](/hugo-pipes/).

baseURL
: Hostname (and path) to the root, e.g. http://bep.is/

blackfriday
: See [Configure Blackfriday](/getting-started/configuration/#configure-blackfriday)

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
: The directory from where Hugo reads content files.

dataDir ("data")
: The directory from where Hugo reads data files.

defaultContentLanguage ("en")
: Content without language indicator will default to this language.

defaultContentLanguageInSubdir (false)
: Render the default content language in subdir, e.g. `content/en/`. The site root `/` will then redirect to `/en/`.

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

enableInlineShortcodes
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
: The site's language code.

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

menu
: See [Add Non-content Entries to a Menu](/content-management/menus/#add-non-content-entries-to-a-menu).

metaDataFormat ("toml")
: Front matter meta-data format. Valid values: `"toml"`, `"yaml"`, or `"json"`.

newContentEditor ("")
: The editor to use when creating new content.

noChmod (false)
: Don't sync permission mode of files.

noTimes (false)
: Don't sync modification time of files.

paginate (10)
: Default number of pages per page in [pagination](/templates/pagination/).

paginatePath ("page")
: The path element used during pagination (https://example.com/page/2).

permalinks
: See [Content Management](/content-management/urls/#permalinks).

pluralizeListTitles (true)
: Pluralize titles in lists.

preserveTaxonomyNames (false)
: Preserve special characters in taxonomy names ("Gérard Depardieu" vs "Gerard Depardieu").

publishDir ("public")
: The directory to where Hugo will write the final static site (the HTML files etc.).

pygmentsCodeFencesGuessSyntax (false)
: Enable syntax guessing for code fences without specified language.

pygmentsStyle ("monokai")
: Color-theme or style for syntax highlighting. See [Pygments Color Themes](https://help.farbox.com/pygments.html).

pygmentsUseClasses (false)
: Enable using external CSS for syntax highlighting.

related
: See [Related Content](/content-management/related/#configure-related-content).

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
: A directory or a list of directories from where Hugo reads [static files][static-files].

stepAnalysis (false)
: Display memory and timing of different steps of the program.

summaryLength (70)
: The length of text to show in a [`.Summary`](/content-management/summaries/#hugo-defined-automatic-summary-splitting).

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
  post: /:year/:month/:title/
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
{{% /note %}}

{{< todo >}}
Test and document setting params via JSON env var.
{{< /todo >}}

## Ignore Files When Rendering

The following statement inside `./config.toml` will cause Hugo to ignore files ending with `.foo` and `.boo` when rendering:

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

## Configure Blackfriday

[Blackfriday](https://github.com/russross/blackfriday) is Hugo's built-in Markdown rendering engine.

Hugo typically configures Blackfriday with sane default values that should fit most use cases reasonably well.

However, if you have specific needs with respect to Markdown, Hugo exposes some of its Blackfriday behavior options for you to alter. The following table lists these Hugo options, paired with the corresponding flags from Blackfriday's source code ( [html.go](https://github.com/russross/blackfriday/blob/master/html.go) and [markdown.go](https://github.com/russross/blackfriday/blob/master/markdown.go)).

{{< readfile file="/content/en/readfiles/bfconfig.md" markdown="true" >}}

{{% note %}}
1. Blackfriday flags are *case sensitive* as of Hugo v0.15.
2. Blackfriday flags must be grouped under the `blackfriday` key and can be set on both the site level *and* the page level. Any setting on a page will override its respective site setting.
{{% /note %}}

{{< code-toggle file="config" >}}
[blackfriday]
  angledQuotes = true
  fractions = false
  plainIDAnchors = true
  extensions = ["hardLineBreak"]
{{< /code-toggle >}}

## Configure Additional Output Formats

Hugo v0.20 introduced the ability to render your content to multiple output formats (e.g., to JSON, AMP html, or CSV). See [Output Formats][] for information on how to add these values to your Hugo project's configuration file.

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
```


You can override any of these cache setting in your own `config.toml`. 

### The keywords explained

:cacheDir
: This is the value of the `cacheDir` config option if set (can also be set via OS env variable `HUGO_CACHEDIR`). It will fall back to `/opt/build/cache/hugo_cache/` on Netlify, or a `hugo_cache` directory below the OS temp dir for the others. This means that if you run your builds on Netlify, all caches configured with `:cacheDir` will be saved and restored on the next build. For other CI vendors, please read their documentation. For an CircleCI example, see [this configuration](https://github.com/bep/hugo-sass-test/blob/6c3960a8f4b90e8938228688bc49bdcdd6b2d99e/.circleci/config.yml).

`:project`

The base directory name of the current Hugo project. This means that, in its default setting, every project will have separated file caches, which means that when you do `hugo --gc` you will not touch files related to other Hugo projects running on the same PC.

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
[yaml]: http://yaml.org/spec/
[static-files]: /content-management/static-files/
