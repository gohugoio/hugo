---
title: Menus
description: Create menus by defining entries, localizing each entry, and rendering the resulting data structure.
categories: []
keywords: []
aliases: [/extras/menus/]
---

## Overview

To create a menu for your site:

1. Define the menu entries
1. [Localize](multilingual/#menus) each entry
1. Render the menu with a [template]

Create multiple menus, either flat or nested. For example, create a main menu for the header, and a separate menu for the footer.

There are three ways to define menu entries:

1. Automatically
1. In front matter
1. In site configuration

> [!note]
> Although you can use these methods in combination when defining a menu, the menu will be easier to conceptualize and maintain if you use one method throughout the site.

## Define automatically

To automatically define a menu entry for each top-level [section](g) of your site, enable the section pages menu in your site configuration.

{{< code-toggle file=hugo >}}
sectionPagesMenu = "main"
{{< /code-toggle >}}

This creates a menu structure that you can access with `site.Menus.main` in your templates. See [menu templates] for details.

## Define in front matter

To add a page to the "main" menu:

{{< code-toggle file=content/about.md fm=true >}}
title = 'About'
menus = 'main'
{{< /code-toggle >}}

Access the entry with `site.Menus.main` in your templates. See [menu templates] for details.

To add a page to the "main" and "footer" menus:

{{< code-toggle file=content/contact.md fm=true >}}
title = 'Contact'
menus = ['main','footer']
{{< /code-toggle >}}

Access the entry with `site.Menus.main` and `site.Menus.footer` in your templates. See [menu templates] for details.

> [!note]
> The configuration key in the examples above is `menus`. The `menu` (singular) configuration key is an alias for `menus`.

### Properties

Use these properties when defining menu entries in front matter:

{{% include "/_common/menu-entry-properties.md" %}}

### Example

This front matter menu entry demonstrates some of the available properties:

{{< code-toggle file=content/products/software.md fm=true >}}
title = 'Software'
[menus.main]
parent = 'Products'
weight = 20
pre = '<i class="fa-solid fa-code"></i>'
[menus.main.params]
class = 'center'
{{< /code-toggle >}}

Access the entry with `site.Menus.main` in your templates. See [menu templates] for details.

## Define in site configuration

See [configure menus](/configuration/menus/).

## Localize

Hugo provides two methods to localize your menu entries. See [multilingual].

## Render

See [menu templates].

[menu templates]: /templates/menu/
[multilingual]: /content-management/multilingual/#menus
[template]: /templates/menu/
