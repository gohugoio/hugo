---
title: Configure Hugo
linktitle: Configuration
description: Often the default settings are good enough, but the config file can provide highly granular control over how your site is rendered.
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

The [directory structure][] of a Hugo website&mdash;or more precisely, the source organization of files containing the website's content and templates&mdash;provides most of the configuration information that Hugo needs in order to generate a finished website.

Because of Hugo's sensible defaults, many websites may not need a configuration file. Hugo is designed to recognize certain typical usage patterns.

## Configuration Lookup Order

Similar to the template [lookup order][], Hugo has a default set of rules for searching for a configuration file in the root of your website's source directory as a default behavior:

1. `./config.toml`
2. `./config.yaml`
3. `./config.json`

In your `config` file, you can direct Hugo as to how you want your website rendered, control your website's menus, and arbitrarily define site-wide parameters specific to your project.

## YAML Configuration

The following is a typical example of a YAML configuration file. The values nested under `params:` will populate the [`.Site.Params`][] variable for use in [templates][]:

{{< code file="config.yml">}}
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
{{< /code >}}

### All Variables, YAML

The following is the full list of Hugo-defined variables in an example YAML file. The values provided in this example represent the default values used by Hugo.

{{< code file="config.yml" download="config.yml" >}}
archetypeDir:               "archetypes"
# hostname (and path) to the root, e.g. http://spf13.com/
baseURL:                    ""
# include content marked as draft
buildDrafts:                false
# include content with publishdate in the future
buildFuture:                false
# include content already expired
buildExpired:               false
# enable this to make all relative URLs relative to content root. Note that this does not affect absolute URLs. See the "URL Management" page
relativeURLs:               false
canonifyURLs:               false
# config file (default is path/config.yaml|json|toml)
config:                     "config.toml"
contentDir:                 "content"
dataDir:                    "data"
defaultExtension:           "html"
defaultLayout:              "post"
# Missing translations will default to this content language
defaultContentLanguage:     "en"
# Renders the default content language in subdir, e.g. /en/. The root directory / will redirect to /en/
defaultContentLanguageInSubdir: false
disableLiveReload:          false
# Do not build RSS files
disableRSS:                 false
# Do not build Sitemap file
disableSitemap:             false
# Enable GitInfo feature
enableGitInfo:              false
# Build robots.txt file
enableRobotsTXT:            false
# Do not render 404 page
disable404:                 false
# Do not inject generator meta tag on homepage
disableHugoGeneratorInject: false
# Allows you to disable all page types and will render nothing related to 'kind';
# values = "page", "home", "section", "taxonomy", "taxonomyTerm", "RSS", "sitemap", "robotsTXT", "404"
disableKinds: []
# Do not make the url/path to lowercase
disablePathToLower:         false                   ""
# Enable Emoji emoticons support for page content; see emoji-cheat-sheet.com
enableEmoji:                false
# Show a placeholder instead of the default value or an empty string if a translation is missing
enableMissingTranslationPlaceholders: false
footnoteAnchorPrefix:       ""
footnoteReturnLinkContents: ""
# google analytics tracking id
googleAnalytics:            ""
# if true, auto-detect Chinese/Japanese/Korean Languages in the content. (.Summary and .WordCount can work properly in CJKLanguage)
hasCJKLanguage:             false
languageCode:               ""
# the length of text to show in a .Summary
summaryLength:              70
layoutDir:                  "layouts"
# Enable Logging
log:                        false
# Log File path (if set, logging enabled automatically)
logFile:                    ""
# "toml","yaml", or "json"
metaDataFormat:             "toml"
newContentEditor:           ""
# Don't sync permission mode of files
noChmod:                    false
# Don't sync modification time of files
noTimes:                    false
# Pagination
paginate:                   10
paginatePath:               "page"
# See "content-management/permalinks"
permalinks:
# Pluralize titles in lists using inflect
pluralizeListTitles:        true
# Preserve special characters in taxonomy names ("Gérard Depardieu" vs "Gerard Depardieu")
preserveTaxonomyNames:      false
# filesystem path to write files to
publishDir:                 "public"
# enables syntax guessing for code fences without specified language
pygmentsCodeFencesGuessSyntax: false
# color-codes for highlighting derived from this style
pygmentsStyle:              "monokai"
# true use pygments-css or false will color code directly
pygmentsUseClasses:         false
# maximum number of items in the RSS feed
rssLimit:                   15
# see "Section Menu for Lazy Bloggers", /templates/menu-templates for more info
SectionPagesMenu:           ""
# default sitemap configuration map
sitemap:
# filesystem path to read files relative from
source:                     ""
staticDir:                  "static"
# display memory and timing of different steps of the program
stepAnalysis:               false
# display metrics about template executions
templateMetrics:            false
# theme to use (located by default in /themes/THEMENAME/)
themesDir:                  "themes"
theme:                      ""
title:                      ""
# Title Case style guide for the title func and other automatic title casing in Hugo.
// Valid values are "AP" (default), "Chicago" and "Go" (which was what you had in Hugo <= 0.25.1).
// See https://www.apstylebook.com/ and http://www.chicagomanualofstyle.org/home.html
titleCaseStyle:             "AP"
# if true, use /filename.html instead of /filename/
uglyURLs:                   false
# verbose output
verbose:                    false
# verbose logging
verboseLog:                 false
# watch filesystem for changes and recreate as needed
watch:                      true
taxonomies:
  - category:               "categories"
  - tag:                    "tags"
{{< /code >}}

## TOML Configuration

The following is an example of a TOML configuration file. The values under `[params]` will populate the `.Site.Params` variable for use in [templates][]:

{{< code file="config.toml">}}
contentDir = "content"
layoutDir = "layouts"
publishDir = "public"
buildDrafts = false
baseURL = "https://yoursite.example.com/"
canonifyURLs = true
title = "My Hugo Site"

[taxonomies]
  category = "categories"
  tag = "tags"

[params]
  subtitle = "Hugo is Absurdly Fast!"
  author = "John Doe"
{{< /code >}}

### All Variables, TOML

The following is the full list of Hugo-defined variables in an example TOML file. The values provided in this example represent the default values used by Hugo.

{{< code file="config.toml" download="config.toml">}}
archetypeDir =                "archetypes"
# hostname (and path) to the root, e.g. http://spf13.com/
baseURL =                     ""
# include content marked as draft
buildDrafts =                 false
# include content with publishdate in the future
buildFuture =                 false
# include content already expired
buildExpired =                false
# enable this to make all relative URLs relative to content root. Note that this does not affect absolute URLs.
relativeURLs =                false
canonifyURLs =                false
# config file (default is path/config.yaml|json|toml)
config =                     "config.toml"
contentDir =                  "content"
dataDir =                     "data"
defaultExtension =            "html"
defaultLayout =               "post"
# Missing translations will default to this content language
defaultContentLanguage =      "en"
# Renders the default content language in subdir, e.g. /en/. The root directory / will redirect to /en/
defaultContentLanguageInSubdir =  false
disableLiveReload =           false
# Do not build RSS files
disableRSS =                  false
# Do not build Sitemap file
disableSitemap =              false
# Enable GitInfo feature
enableGitInfo =               false
# Build robots.txt file
enableRobotsTXT =             false
# Do not render 404 page
disable404 =                  false
# Do not inject generator meta tag on homepage
disableHugoGeneratorInject =  false
# Allows you to disable all page types and will render nothing related to 'kind';
# values = "page", "home", "section", "taxonomy", "taxonomyTerm", "RSS", "sitemap", "robotsTXT", "404"
disableKinds = []
# Do not make the url/path to lowercase
disablePathToLower =          false
# Enable Emoji emoticons support for page content; see emoji-cheat-sheet.com
enableEmoji =                 false
# Show a placeholder instead of the default value or an empty string if a translation is missing
enableMissingTranslationPlaceholders = false
footnoteAnchorPrefix =        ""
footnoteReturnLinkContents =  ""
# google analytics tracking id
googleAnalytics =             ""
# if true, auto-detect Chinese/Japanese/Korean Languages in the content. (.Summary and .WordCount can work properly in CJKLanguage)
hasCJKLanguage =              false
languageCode =                ""
# the length of text to show in a .Summary
summaryLength:              70
layoutDir =                   "layouts"
# Enable Logging
log =                         false
# Log File path (if set, logging enabled automatically)
logFile =
# maximum number of items in the RSS feed
rssLimit =                    15
# "toml","yaml", or "json"
metaDataFormat =              "toml"
newContentEditor =            ""
# Don't sync permission mode of files
noChmod =                     false
# Don't sync modification time of files
noTimes =                     false
# Pagination
paginate =                    10
paginatePath =                "page"
# See "content-management/permalinks"
permalinks =
# Pluralize titles in lists using inflect
pluralizeListTitles =         true
# Preserve special characters in taxonomy names ("Gérard Depardieu" vs "Gerard Depardieu")
preserveTaxonomyNames =       false
# filesystem path to write files to
publishDir =                  "public"
# enables syntax guessing for code fences without specified language
pygmentsCodeFencesGuessSyntax = false
# color-codes for highlighting derived from this style
pygmentsStyle =               "monokai"
# true: use pygments-css or false: color-codes directly
pygmentsUseClasses =          false
# see "Section Menu for Lazy Bloggers", /templates/menu-templates for more info
SectionPagesMenu =
# default sitemap configuration map
sitemap =
# filesystem path to read files relative from
source =                      ""
staticDir =                   "static"
# display memory and timing of different steps of the program
stepAnalysis =                false
# theme to use (located by default in /themes/THEMENAME/)
themesDir =                   "themes"
theme =                       ""
title =                       ""
# if true, use /filename.html instead of /filename/
uglyURLs =                    false
# verbose output
verbose =                     false
# verbose logging
verboseLog =                  false
# watch filesystem for changes and recreate as needed
watch =                       true
[taxonomies]
  category = "categories"
  tag = "tags"
{{< /code >}}

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

## Environmental Variables

In addition to the 3 config options already mentioned, configuration key-values can be defined through operating system environment variables.

For example, the following command will effectively set a website's title on Unix-like systems:

```
$ env HUGO_TITLE="Some Title" hugo
```

{{% note "Setting Environment Variables" %}}
Names must be prefixed with `HUGO_` and the configuration key must be set in uppercase when setting operating system environment variables.
{{% /note %}}

## Ignore Files When Rendering

The following statement inside `./config.toml` will cause Hugo to ignore files ending with `.foo` and `.boo` when rendering:

```
ignoreFiles = [ "\\.foo$", "\\.boo$" ]
```

The above is a list of regular expressions. Note that the backslash (`\`) character is escaped in this example to keep TOML happy.

## Configure Blackfriday

[Blackfriday](https://github.com/russross/blackfriday) is Hugo's built-in Markdown rendering engine.

Hugo typically configures Blackfriday with sane default values that should fit most use cases reasonably well.

However, if you have specific needs with respect to Markdown, Hugo exposes some of its Blackfriday behavior options for you to alter. The following table lists these Hugo options, paired with the corresponding flags from Blackfriday's source code ( [html.go](https://github.com/russross/blackfriday/blob/master/html.go) and [markdown.go](https://github.com/russross/blackfriday/blob/master/markdown.go)).

{{< readfile file="/content/readfiles/bfconfig.md" markdown="true" >}}

{{% note %}}
1. Blackfriday flags are *case sensitive* as of Hugo v0.15.
2. Blackfriday flags must be grouped under the `blackfriday` key and can be set on both the site level *and* the page level. Any setting on a page will override its respective site setting.
{{% /note %}}

{{< code file="bf-config.toml" >}}
[blackfriday]
  angledQuotes = true
  fractions = false
  plainIDAnchors = true
  extensions = ["hardLineBreak"]
{{< /code >}}

{{< code file="bf-config.yml" >}}
blackfriday:
  angledQuotes: true
  fractions: false
  plainIDAnchors: true
  extensions:
    - hardLineBreak
{{< /code >}}

## Configure Additional Output Formats

Hugo v0.20 introduced the ability to render your content to multiple output formats (e.g., to JSON, AMP html, or CSV). See [Output Formats][] for information on how to add these values to your Hugo project's configuration file.

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
