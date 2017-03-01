---
title: Configuration
linktitle:
description: Hugo is designed to make enough assumptions that often configuration is unnecessary. However, a site config file can include precise directions to Hugo on how you want to render your website.
date: 2017-01-02
publishdate: 2017-01-02
lastmod: 2017-01-02
categories: [project organization]
tags: [configuration,fundamentals,toml,yaml,json]
weight: 60
draft: false
aliases: [/overview/source-directory/]
toc: true
notesforauthors:
---

The [directory structure][] of a Hugo website&mdash;or more precisely, the source organization of files containing the website's content and templates&mdash;provides most of the configuration information that Hugo needs in order to statically generate a finished website.

Because of Hugo's preference for sane defaults, many websites may not need a configuration file. Hugo is designed to recognize certain typical usage patterns (and even expects them by default).


## Configuration Lookup Order

Hugo searches for a configuration file in the root of your website's source directory as a default behavior. First, it looks for a `./config.toml` file. If that's not present, it will seek a `./config.yaml` file, followed by a `./config.json` file.

In this `config` file, you can direct to Hugo as to how it should render your website, control your website's menus, and arbitrarily define site-wide parameters specific to your project.

## YAML Configuration

The following is a typical example of a YAML configuration file. Note the document opens with 3 hyphens and closes with 3 periods. The values nested under `params:` will populate the [`.Site.Params`][] variable for use in [templates][]:

{{% code file="config.yml"%}}
```yaml
---
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
...
```
{{% /code %}}

### All Variables, YAML

The following is the full list of Hugo-defined variables in an example YAML file. The values provided in this example represent the default values used by Hugo.

{{% code file="config.yml" download="config.yml" %}}
```yaml
---
archetypeDir:               "archetypes"
# hostname (and path) to the root, e.g. http://spf13.com/
baseURL:                    ""
# include content marked as draft
buildDrafts:                false
# include content with publishdate in the future
buildFuture:                false
# include content already expired
buildExpired:               false
# enable this to make all relative URLs relative to content root. Note that this does not affect absolute URLs.
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
disablePathToLower:         false
# edit new content with this editor, if provided
editor:                     ""
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
# theme to use (located by default in /themes/THEMENAME/)
themesDir:                  "themes"
theme:                      ""
title:                      ""
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
---
```
{{% /code %}}

## TOML Configuration

The following is an example of a TOML configuration file. The values under `[params]` will populate the `.Site.Params` variable for use in [templates][]:

```toml
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
```

### All Variables, TOML

The following is the full list of Hugo-defined variables in an example TOML file. The values provided in this example represent the default values used by Hugo.

{{% code file="config.toml" download="config.toml"%}}
```toml
+++
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
# edit new content with this editor, if provided
editor =                      ""
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
+++
```
{{% /code %}}

## Configuration Through Environmental Variables

In addition to the 3 configuration file options already mentioned, website configuration can be accomplished through operating system environment variables.

For example, the following command will effectively set a website's title on Unix-like systems:

```bash
$ env HUGO_TITLE="Some Title" hugo
```

{{% note "Setting Environment Variables" %}}
Names must be prefixed with `HUGO_` when setting environment variables through operating system environment variables.
{{% /note %}}

## Ignoring Files When Rendering

The following statement inside `./config.toml` will cause Hugo to ignore files
ending with `.foo` and `.boo` when rendering:

```toml
ignoreFiles = [ "\\.foo$", "\\.boo$" ]
```
The above is a list of regular expressions. Note that the backslash (`\`) character is escaped in this example to keep TOML happy.

## Configuring Blackfriday Rendering

[Blackfriday](https://github.com/russross/blackfriday) is Hugo's built-in
[Markdown](http://daringfireball.net/projects/markdown/) rendering engine.

Hugo typically configures Blackfriday with sane default values. These defaults should fit most use cases reasonably well.

However, if you have specific needs with respect to Markdown, Hugo exposes some of its Blackfriday behavior options for you to alter. The following table lists these Hugo options, paired with the corresponding flags from Blackfriday's source code ( [html.go](https://github.com/russross/blackfriday/blob/master/html.go) and [markdown.go](https://github.com/russross/blackfriday/blob/master/markdown.go)).

**WIP: WORKING ON MOVING BF CONFIG FROM TABLE TO DL**

`taskLists`
: default: **`true`**<br>
    Blackfriday flag: **``**<br>
    Purpose: `false` turns off GitHub-style automatic task/TODO list generation

`smartypants`
: default: **`true`** <br>
    Blackfriday flag: **`HTML_USE_SMARTYPANTS`** <br>
    Purpose: `false` disables smart punctuation substitutions, including smart quotes, smart dashes, smart fractions, etc. If `true`, it may be fine-tuned with the `angledQuotes`, `fractions`, `smartDashes`, and `latexDashes` flags (see below).

`angledQuotes`
: default: **`false`**<br>
    Blackfriday flag: **`HTML_SMARTYPANTS_ANGLED_QUOTES`**<br>
    Purpose: `true` enables smart, angled double quotes. Example: "Hugo" renders to renders to «Hugo» instead of “Hugo”.

`fractions`
: default: **`true`**<br>
    Blackfriday flag: **`HTML_SMARTYPANTS_FRACTIONS`** <br>
    Purpose: <code>false</code> disables smart fractions.<br>
    Example: `5/12` renders to <sup>5</sup>&frasl;<sub>12</sub>(<code>&lt;sup&gt;5&lt;/sup&gt;&amp;frasl;&lt;sub&gt;12&lt;/sub&gt;</code>).<br> <strong>Caveat:</strong> Even with <code>fractions = false</code>, Blackfriday still converts `1/2`, `1/4`, and `3/4` respectively to ½ (<code>&amp;frac12;</code>), ¼ (<code>&amp;frac14;</code>) and ¾ (<code>&amp;frac34;</code>), but only these three.</small>

`smartDashes`
: default: **`true`** <br>
    Blackfriday flag: **`HTML_SMARTY_DASHES`** <br>
    Purpose: `false` disables smart dashes; i.e., the conversion of multiple hyphens into an en dash or em dash. If `true`, its behavior can be modified with the `latexDashes` flag below.

`latexDashes`
: default: **`true`** <br>
    Blackfriday flag: **`HTML_SMARTYPANTS_LATEX_DASHES`** <br>
    Purpose: `false` disables LaTeX-style smart dashes and selects conventional smart dashes. Assuming `smartDashes`: <br>
    If `true`, `--` is translated into &ndash; (`&ndash;`), whereas `---` is translated into &mdash; (`&mdash;`). <br>
    However, *spaced* single hyphen between two words is translated into an en&nbsp;dash&mdash; e.g., "`12 June - 3 July`" becomes `12 June ndash; 3 July` upon rendering.

`hrefTargetBlank`
: default: **`false`** <br>
    Blackfriday flag: **`HTML_HREF_TARGET_BLANK`** <br>
    Purpose: `true` opens external links in a new window or tab.


{{< bfconfig >}}

{{% note %}}
1. Blackfriday flags are *case sensitive* as of Hugo v0.15.
2. Blackfriday flags must be grouped under the `blackfriday` key and can be set on both the site level *and* the page level. Any setting on a page will override the site setting there. See [site configuration for more information](/content-management/front-matter/#override-global-blackfriday-configuration).
{{% /note %}}

{{% code file="bf-config.toml" %}}
```toml
[blackfriday]
  angledQuotes = true
  fractions = false
  plainIDAnchors = true
  extensions = ["hardLineBreak"]
```
{{% /code %}}

{{% code file="bf-config.yml" %}}
```yaml
blackfriday:
  angledQuotes: true
  fractions: false
  plainIDAnchors: true
  extensions:
    - hardLineBreak
```
{{% /code %}}

## Specs for Configuration Formats

* [TOML Spec][toml]
* [YAML Spec][yaml]
* [JSON Spec][json]

[`.Site.Params`]: /variables/
[directory structure]: /project-organization/directory-structure
[json]: /documents/ecma-404-json-spec.pdf
[templates]: /templates/
[toml]: https://github.com/toml-lang/toml
[yaml]: http://yaml.org/spec/