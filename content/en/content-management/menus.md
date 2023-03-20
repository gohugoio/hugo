---
title: Menus
linkTitle: Menus
description:  Create menus by defining entries, localizing each entry, and rendering the resulting data structure.
categories: [content management]
keywords: [menus]
menu:
  docs:
    parent: content-management
    weight: 190
toc: true
weight: 190
aliases: [/extras/menus/]
---

## Overview

To create a menu for your site:

1. Define the menu entries
2. [Localize] each entry
3. Render the menu with a [template]

Create multiple menus, either flat or nested. For example, create a main menu for the header, and a separate menu for the footer.

There are three ways to define menu entries:

1. Automatically
1. In front matter
1. In site configuration

{{% note %}}
Although you can use these methods in combination when defining a menu, the menu will be easier to conceptualize and maintain if you use one method throughout the site.
{{% /note %}}

## Define automatically

To automatically define menu entries for each top-level section of your site, enable the section pages menu in your site configuration.
q
{{< code-toggle file="config" copy=false >}}
sectionPagesMenu = "main"
{{< /code-toggle >}}

This creates a menu structure that you can access with `site.Menus.main` in your templates. See [menu templates] for details.

## Define in front matter

To add a page to the "main" menu:

{{< code-toggle file="content/about.md" copy=false fm=true >}}
title = 'About'
menu = 'main'
{{< /code-toggle >}}

Access the entry with `site.Menus.main` in your templates. See [menu templates] for details.

To add a page to the "main" and "footer" menus:

{{< code-toggle file="content/contact.md" copy=false fm=true >}}
title = 'Contact'
menu = ['main','footer']
{{< /code-toggle >}}

Access the entry with `site.Menus.main` and `site.Menus.footer` in your templates. See [menu templates] for details.

### Properties {#properties-front-matter}

Use these properties when defining menu entries in front matter:

identifier
: (string) Required when two or more menu entries have the same `name`, or when localizing the `name` using translation tables. Must start with a letter, followed by letters, digits, or underscores.

name
: (string) The text to display when rendering the menu entry.

params
: (map) User-defined properties for the menu entry.

parent
: (string) The `identifier` of the parent menu entry. If `identifier` is not defined, use `name`. Required for child entries in a nested menu.

post
: (string) The HTML to append when rendering the menu entry.

pre
: (string) The HTML to prepend when rendering the menu entry.

title
: (string) The HTML `title` attribute of the rendered menu entry.

weight
: (int) A non-zero integer indicating the entry's position relative the root of the menu, or to its parent for a child entry. Lighter entries float to the top, while heavier entries sink to the bottom.

### Example {#example-front-matter}

This front matter menu entry demonstrates some of the available properties:

{{< code-toggle file="content/products/software.md" copy=false fm=true >}}
title = 'Software'
[menu.main]
parent = 'Products'
weight = 20
pre = '<i class="fa-solid fa-code"></i>'
[menu.main.params]
class = 'center'
{{< /code-toggle >}}

Access the entry with `site.Menus.main` in your templates. See [menu templates] for details.


## Define in site configuration

To define entries for the "main" menu:

{{< code-toggle file="config" copy=false >}}
[[menu.main]]
name = 'Home'
pageRef = '/'
weight = 10

[[menu.main]]
name = 'Products'
pageRef = '/products'
weight = 20

[[menu.main]]
name = 'Services'
pageRef = '/services'
weight = 30
{{< /code-toggle >}}

This creates a menu structure that you can access with `site.Menus.main` in your templates. See [menu templates] for details.

To define entries for the "footer" menu:

{{< code-toggle file="config" copy=false >}}
[[menu.footer]]
name = 'Terms'
pageRef = '/terms'
weight = 10

[[menu.footer]]
name = 'Privacy'
pageRef = '/privacy'
weight = 20
{{< /code-toggle >}}

This creates a menu structure that you can access with `site.Menus.footer` in your templates. See [menu templates] for details.

### Properties {#properties-site-configuration}

Each menu entry defined in site configuration requires two or more properties:

- Specify `name` and `pageRef` for internal links
- Specify `name` and `url` for external links

pageRef
: (string) The file path of the target page, relative to the `content` directory. Required for *internal* links.

url
: (string) Required for *external* links.

{{% note %}}
The [properties] available to entries defined in front matter are also available to entries defined in site configuration.

[properties]: /content-management/menus/#properties-front-matter
{{% /note %}}

### Example {#example-site-configuration}

This nested menu demonstrates some of the available properties:

{{< code-toggle file="config" copy=false >}}
[[menu.main]]
name = 'Products'
pageRef = '/products'
weight = 10

[[menu.main]]
name = 'Hardware'
pageRef = '/products/hardware'
parent = 'Products'
weight = 1

[[menu.main]]
name = 'Software'
pageRef = '/products/software'
parent = 'Products'
weight = 2

[[menu.main]]
name = 'Services'
pageRef = '/services'
weight = 20

[[menu.main]]
name = 'Hugo'
pre = '<i class="fa fa-heart"></i>'
url = 'https://gohugo.io/'
weight = 30
[menu.main.params]
rel = 'external'
{{< /code-toggle >}}

This creates a menu structure that you can access with `site.Menus.main` in your templates. See [menu templates] for details.

## Localize

Hugo provides two methods to localize your menu entries. See [multilingual].

## Render

See [menu templates].

[Localize]: /content-management/multilingual/#menus
[menu templates]: /templates/menu-templates/
[multilingual]: /content-management/multilingual/#menus
[template]: /templates/menu-templates/
