---
title: Menus
description:  Create menus by defining entries, localizing each entry, and rendering the resulting data structure.
categories: [content management]
keywords: [menus]
menu:
  docs:
    parent: content-management
    weight: 190
weight: 190
toc: true
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

To automatically define a menu entry for each top-level [section] of your site, enable the section pages menu in your site configuration.

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

{{% note %}}
The configuration key in the examples above is `menus`. The `menu` (singular) configuration key is an alias for `menus`.
{{% /note %}}

### Properties {#properties-front-matter}

Use these properties when defining menu entries in front matter:

identifier
: (`string`) Required when two or more menu entries have the same `name`, or when localizing the `name` using translation tables. Must start with a letter, followed by letters, digits, or underscores.

name
: (`string`) The text to display when rendering the menu entry.

params
: (`map`) User-defined properties for the menu entry.

parent
: (`string`) The `identifier` of the parent menu entry. If `identifier` is not defined, use `name`. Required for child entries in a nested menu.

post
: (`string`) The HTML to append when rendering the menu entry.

pre
: (`string`) The HTML to prepend when rendering the menu entry.

title
: (`string`) The HTML `title` attribute of the rendered menu entry.

weight
: (`int`) A non-zero integer indicating the entry's position relative the root of the menu, or to its parent for a child entry. Lighter entries float to the top, while heavier entries sink to the bottom.

### Example {#example-front-matter}

This front matter menu entry demonstrates some of the available properties:

{{< code-toggle file=content/products/software.md fm=true >}}
title = 'Software'
[[menus.main]]
parent = 'Products'
weight = 20
pre = '<i class="fa-solid fa-code"></i>'
[menus.main.params]
class = 'center'
{{< /code-toggle >}}

Access the entry with `site.Menus.main` in your templates. See [menu templates] for details.

## Define in site configuration

To define entries for the "main" menu:

{{< code-toggle file=hugo >}}
[[menus.main]]
name = 'Home'
pageRef = '/'
weight = 10

[[menus.main]]
name = 'Products'
pageRef = '/products'
weight = 20

[[menus.main]]
name = 'Services'
pageRef = '/services'
weight = 30
{{< /code-toggle >}}

This creates a menu structure that you can access with `site.Menus.main` in your templates. See [menu templates] for details.

To define entries for the "footer" menu:

{{< code-toggle file=hugo >}}
[[menus.footer]]
name = 'Terms'
pageRef = '/terms'
weight = 10

[[menus.footer]]
name = 'Privacy'
pageRef = '/privacy'
weight = 20
{{< /code-toggle >}}

This creates a menu structure that you can access with `site.Menus.footer` in your templates. See [menu templates] for details.

{{% note %}}
The configuration key in the examples above is `menus`. The `menu` (singular) configuration key is an alias for `menus`.
{{% /note %}}

### Properties {#properties-site-configuration}

{{% note %}}
The [properties available to entries defined in front matter] are also available to entries defined in site configuration.

[properties available to entries defined in front matter]: /content-management/menus/#properties-front-matter
{{% /note %}}

Each menu entry defined in site configuration requires two or more properties:

- Specify `name` and `pageRef` for internal links
- Specify `name` and `url` for external links

pageRef
: (`string`) The logical path of the target page, relative to the `content` directory. Omit language code and file extension. Required for *internal* links.

Kind|pageRef
:--|:--
home|`/`
page|`/books/book-1`
section|`/books`
taxonomy|`/tags`
term|`/tags/foo`

url
: (`string`) Required for *external* links.

### Example {#example-site-configuration}

This nested menu demonstrates some of the available properties:

{{< code-toggle file=hugo >}}
[[menus.main]]
name = 'Products'
pageRef = '/products'
weight = 10

[[menus.main]]
name = 'Hardware'
pageRef = '/products/hardware'
parent = 'Products'
weight = 1

[[menus.main]]
name = 'Software'
pageRef = '/products/software'
parent = 'Products'
weight = 2

[[menus.main]]
name = 'Services'
pageRef = '/services'
weight = 20

[[menus.main]]
name = 'Hugo'
pre = '<i class="fa fa-heart"></i>'
url = 'https://gohugo.io/'
weight = 30
[menus.main.params]
rel = 'external'
{{< /code-toggle >}}

This creates a menu structure that you can access with `site.Menus.main` in your templates. See [menu templates] for details.

## Localize

Hugo provides two methods to localize your menu entries. See [multilingual].

## Render

See [menu templates].

[localize]: /content-management/multilingual/#menus
[menu templates]: /templates/menu/
[multilingual]: /content-management/multilingual/#menus
[section]: /getting-started/glossary/#section
[template]: /templates/menu/
