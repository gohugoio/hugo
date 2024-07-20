---
title: Theme components
description: Hugo provides advanced theming support with Theme Components.
categories: [hugo modules]
keywords: [modules,themes]
menu:
  docs:
    parent: modules
    weight: 40
weight: 40
aliases: [/themes/customize/,/themes/customizing/]
toc: true
---

{{% note %}}
This section contain information that may be outdated and is in the process of being rewritten.
{{% /note %}}
Since Hugo `0.42` a project can configure a theme as a composite of as many theme components you need:

{{< code-toggle file=hugo >}}
theme = ["my-shortcodes", "base-theme", "hyde"]
{{< /code-toggle >}}

You can even nest this, and have the theme component itself include theme components in its own `hugo.toml` (theme inheritance).[^1]

The theme definition example above in `hugo.toml` creates a theme with 3 theme components with precedence from left to right.

For any given file, data entry, etc., Hugo will look first in the project and then in `my-shortcodes`, `base-theme`, and lastly `hyde`.

Hugo uses two different algorithms to merge the file systems, depending on the file type:

* For `i18n` and `data` files, Hugo merges deeply using the translation ID and data key inside the files.
* For `static`, `layouts` (templates), and `archetypes` files, these are merged on file level. So the left-most file will be chosen.

The name used in the `theme` definition above must match a folder in `/your-site/themes`, e.g. `/your-site/themes/my-shortcodes`. There are plans to improve on this and get a URL scheme so this can be resolved automatically.

Also note that a component that is part of a theme can have its own configuration file, e.g. `hugo.toml`. There are currently some restrictions to what a theme component can configure:

* `params` (global and per language)
* `menu` (global and per language)
* `outputformats` and `mediatypes`

The same rules apply here: The left-most parameter/menu etc. with the same ID will win. There are some hidden and experimental namespace support in the above, which we will work to improve in the future, but theme authors are encouraged to create their own namespaces to avoid naming conflicts.

[^1]: For themes hosted on the [Hugo Themes Showcase](https://themes.gohugo.io/) components need to be added as git submodules that point to the directory `exampleSite/themes`
