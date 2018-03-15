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


## All Configuration Settings

The following is the full list of Hugo-defined variables with its default value in parens.

archetypeDir ("archetypes")
: The directory where Hugo finds archetype files (content templates).

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

canonifyURLs (false)
: Enable to turn relative URLs into absolute.

config ("config.toml")
: Config file (default is path/config.yaml|json|toml).

contentDir ("content")
: The directory from where Hugo reads content files.

dataDir ("data")
: The directory from where Hugo reads data files.

defaultContentLanguage ("en")
: Content without language indicator will default to this language.

defaultContentLanguageInSubdir (false)
: Renders the default content language in subdir, e.g. /en/. The root directory / will redirect to /en/.

disableHugoGeneratorInject (false)
: Hugo will, by default, inject a generator meta tag in the HTML head on the _home page only_. You can turn it off, but we would really appreciate if you don't, as this is a good way to watch Hugo's popularity on the rise.

disableKinds ([])
: Allows you to disable all page types and will render nothing related to 'kind'. Allowed values are "page", "home", "section", "taxonomy", "taxonomyTerm", "RSS", "sitemap", "robotsTXT", "404".

disableLiveReload (false)
: Turn off automatic live reloading of browser window.

disablePathToLower (false)
: Do not make the url/path to lowercase.

enableEmoji (false)
: Enable Emoji emoticons support for page content; see emoji-cheat-sheet.com.

enableGitInfo (false)
: If the Hugo site is versioned by Git, you will then get a `.GitInfo` object per page, and `Lastmod` will get updated by the last commit date for content.	

enableMissingTranslationPlaceholders (false)
: Show a placeholder instead of the default value or an empty string if a translation is missing

enableRobotsTXT (false)
: When enabled, Hugo will generate a `robots.txt` file.

footnoteAnchorPrefix ("")
: A prefix for your footnote anchors.

footnoteReturnLinkContents ("")
: A return link for your footnote.

googleAnalytics ("")
: google analytics tracking id

hasCJKLanguage (false)
: If true, auto-detect Chinese/Japanese/Korean Languages in the content. This will make `.Summary` and `.WordCount` behave correctly in CJK languages.

imaging
: See [Image Processing Config](/content-management/image-processing/#image-processing-config).

languages
: See [Configure Languages](/content-management/multilingual/#configure-languages).

languageCode ("")
: The site's language code.

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
: "toml","yaml", or "json"

newContentEditor ("")
: The editor to use when creating new content.

noChmod (false)
: Don't sync permission mode of files.

noTimes (false)
: Don't sync modification time of files

paginate (10)
: Default number of pages per page in pagination.

paginatePath ("page")
: The path element used during pagination (http://example.com/page/2).

permalinks
: See [Content Management](/content-management/urls/#permalinks)

pluralizeListTitles (true)
: Pluralize titles in lists using inflect.

preserveTaxonomyNames (false)
: Preserve special characters in taxonomy names ("Gérard Depardieu" vs "Gerard Depardieu").

publishDir ("public")
: The directory to where Hugo will write the final static site (the HTML files etc.).

pygmentsCodeFencesGuessSyntax (false)
: Enables syntax guessing for code fences without specified language.

pygmentsStyle ("monokai")
: Color-codes for highlighting derived from this style. See https://help.farbox.com/pygments.html

pygmentsUseClasses (false)
: Enable to use external CSS for code highlighting.

related
: See [Related Content](/content-management/related/#configure-related-content).

relativeURLs (false)
: Enable this to make all relative URLs relative to content root. Note that this does not affect absolute URLs.

rssLimit (unlimited)
: Maximum number of items in the RSS feed.

sectionPagesMenu ("")(
: See ["Section Menu for Lazy Bloggers"](/templates/menu-templates/#section-menu-for-lazy-bloggers).

sitemap
: Default sitemap configuration.

staticDir ("static")
: Relative directory from where Hugo reads static files.

stepAnalysis (false)
: Display memory and timing of different steps of the program.

summaryLength (70)
: The length of text to show in a `.Summary`.

taxonomies
: See [Configure Taxonomies](content-management/taxonomies#configure-taxonomies)

theme ("")
: Theme to use (located by default in /themes/THEMENAME/)

themesDir ("themes")
: The directory where Hugo reads the themes from.

title ("")
: Site title.

uglyURLs (false)
: When enabled creates URL on the form `/filename.html` instead of `/filename/`

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
