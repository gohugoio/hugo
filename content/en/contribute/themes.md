---
title: Add Your Hugo Theme to the Showcase
linktitle: Themes
description: If you've built a Hugo theme and want to contribute back to the Hugo Community, add your theme to the Hugo Showcase.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-27
categories: [contribute]
keywords: [contribute,themes,design]
authors: [digitalcraftsman]
menu:
  docs:
    parent: "contribute"
    weight: 30
weight: 30
sections_weight: 30
draft: false
aliases: [/contribute/theme/]
wip: true
toc: true
---

A collection of all themes created by the Hugo community, including screenshots and demos, can be found at <https://themes.gohugo.io>. Every theme in this list will automatically be added to the theme site. Theme updates aren't scheduled but usually happen at least once a week.

## Adding a theme to the list

1. Create your theme using `hugo new theme <THEMENAME>`;
2. Test your theme against <https://github.com/gohugoio/hugoBasicExample> \*
3. Add a `theme.toml` file to the root of the theme with all required metadata
4. Add a descriptive `README.md` to the root of the theme source
5. Add `/images/screenshot.png` and `/images/tn.png`

\* If your theme doesn't fit into the `Hugo Basic Example` site, we encourage theme authors to supply a self-contained Hugo site in `/exampleSite`.

{{% note %}}
The folder name here---`exampleSite`---is important, as this folder will be picked up and used by the script that generates the Hugo Theme Site. It mirrors the root directory of a Hugo website and allows you to add custom content, assets, and a `config` file with preset values.
{{% /note %}}

See the [Hugo Artist theme's exampleSite][artistexample] for a good example.

{{% note %}}
Please make your example site's content is as neutral as possible. We hope this goes without saying.
{{% /note %}}

## Theme Requirements

In order to add your theme to the Hugo Themes Showcase, the following requirements need to be met:

1. `theme.toml` with all required fields
2. Images for thumbnail and screenshot
3. A good README file instructions for users
4. Addition to the hugoThemes GitHub repository

### Add Your Theme to the Repo

The easiest way to add your theme is to [open up a new issue in the theme repository][themeissuenew] with a link to the theme's repository on GitHub.

### Create a `theme.toml` File

`theme.toml` contains metadata about the theme and its creator and should be created automatically when running the `hugo new theme`. The auto-generated file is provided here as well for easy downloading:

{{< code file="theme.toml" download="theme.toml" >}}
name = ""
license = "MIT"
licenselink = "https://github.com/<YOURNAME>/<YOURTHEME>/blob/master/LICENSE.md"
description = ""
homepage = "https://example.com/"
tags = []
features = []
min_version = 0.19

[author]
  name = ""
  homepage = ""

# If porting an existing theme
[original]
  name = ""
  homepage = ""
  repo = ""
{{< /code >}}

The following fields are required:

```
name = "Hyde"
license = "MIT"
licenselink = "https://github.com/spf13/hyde/blob/master/LICENSE.md"
description = "An elegant open source and mobile first theme"
homepage = "http://siteforthistheme.com/"
tags = ["blog", "company"]
features = ["blog"]
min_version = 0.13

[author]
    name = "spf13"
    homepage = "http://spf13.com/"

# If porting an existing theme
[original]
    author = "mdo"
    homepage = "http://hyde.getpoole.com/"
    repo = "https://www.github.com/mdo/hyde"
```

{{% note %}}
1. This is different from the `theme.toml` file created by `hugo new theme` in Hugo versions before v0.14.
2. Only `theme.toml` is accepted; ie. not `theme.yaml` and `theme.json`.
{{% /note %}}

### Images

Screenshots are used for previews in the Hugo Theme Gallery. Make sure that they have the right dimensions:

* Thumbnail should be 900px × 600px
* Screenshot should be 1500px × 1000px
* Media must be located in
    * <THEMEDIR>/images/screenshot.png</code>
    * <THEMEDIR>/images/tn.png</code>

Additional media may be provided in the same directory.

### Create a README File

Your theme's README file should be written in markdown and saved at the root of your theme's directory structure. Your `README.md` serves as

1. Content for your theme's details page at <https://themes.gohugo.io>
2. General information about the theme in your GitHub repository (i.e., it's usual purpose)

#### Example `README.md`

You can download the following `README.md` as an outline:

{{< code file="README.md" download="README.md" >}}

# Theme Title

**Need input from @digitalcraftsman on what could be added to this file.**




{{< /code >}}

{{% note "Screenshots in your `README.md`"%}}
If you add screenshots to the README, please make use of absolute file paths instead of relative ones like `/images/screenshot.png`. Relative paths work great on GitHub but they don't correspond to the directory structure of [themes.gohugo.io](http://themes.gohugo.io/). Therefore, browsers will not be able to display screenshots on the theme site under the given (relative) path.
{{% /note %}}

[artistexample]: https://github.com/digitalcraftsman/hugo-artists-theme/tree/master/exampleSite
[themeissuenew]: https://github.com/gohugoio/hugoThemes/issues/new
