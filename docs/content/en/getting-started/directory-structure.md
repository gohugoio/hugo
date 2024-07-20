---
title: Directory structure
description: Each Hugo project is a directory, with subdirectories that contribute to the content, structure, behavior, and presentation of your site.
categories: [getting started,fundamentals]
keywords: [source, organization, directories]
menu:
  docs:
    parent: getting-started
    weight: 30
weight: 30
toc: true
aliases: [/overview/source-directory/]
---

## Site skeleton

Hugo generates a project skeleton when you create a new site. For example, this command:

```sh
hugo new site my-site
```

Creates this directory structure:

```txt
my-site/
├── archetypes/
│   └── default.md
├── assets/
├── content/
├── data/
├── i18n/
├── layouts/
├── static/
├── themes/
└── hugo.toml         <-- site configuration
```

Depending on requirements, you may wish to organize your site configuration into subdirectories:

```txt
my-site/
├── archetypes/
│   └── default.md
├── assets/
├── config/           <-- site configuration
│   └── _default/
│       └── hugo.toml
├── content/
├── data/
├── i18n/
├── layouts/
├── static/
└── themes/
```

When you build your site, Hugo creates a `public` directory, and typically a `resources` directory as well:

```txt
my-site/
├── archetypes/
│   └── default.md
├── assets/
├── config/       
│   └── _default/
│       └── hugo.toml
├── content/
├── data/
├── i18n/
├── layouts/
├── public/       <-- created when you build your site
├── resources/    <-- created when you build your site
├── static/
└── themes/
```

## Directories

Each of the subdirectories contributes to the content, structure, behavior, or presentation of your site.

###### archetypes

The `archetypes` directory contains templates for new content. See&nbsp;[details](/content-management/archetypes/).

###### assets

The `assets` directory contains global resources typically passed through an asset pipeline. This includes resources such as images, CSS, Sass, JavaScript, and TypeScript. See&nbsp;[details](/hugo-pipes/introduction/).

###### config

The `config` directory contains your site configuration, possibly split into multiple subdirectories and files. For projects with minimal configuration or projects that do not need to behave differently in different environments, a single configuration file named `hugo.toml` in the root of the project is sufficient. See&nbsp;[details](/getting-started/configuration/#configuration-directory).

###### content

The `content` directory contains the markup files (typically Markdown) and page resources that comprise the content of your site. See&nbsp;[details](/content-management/organization/).

###### data

The `data` directory contains data files (JSON, TOML, YAML, or XML) that augment content, configuration, localization, and navigation. See&nbsp;[details](/content-management/data-sources/).

###### i18n

The `i18n` directory contains translation tables for multilingual sites. See&nbsp;[details](/content-management/multilingual/).

###### layouts

The layouts directory contains templates to transform content, data, and resources into a complete website. See&nbsp;[details](/templates/).

###### public

The `public` directory contains the published website, generated when you run the `hugo` or `hugo server` commands. Hugo recreates this directory and its content as needed. See&nbsp;[details](/getting-started/usage/#build-your-site).

###### resources

The `resources` directory contains cached output from Hugo's asset pipelines, generated when you run the `hugo` or `hugo server` commands. By default this cache directory includes CSS and images. Hugo recreates this directory and its content as needed.

###### static

The `static` directory contains files that will be copied to the public directory when you build your site. For example: `favicon.ico`, `robots.txt`, and files that verify site ownership. Before the introduction of [page bundles](/getting-started/glossary/#page-bundle) and [asset pipelines](/hugo-pipes/introduction/), the `static` directory was also used for images, CSS, and JavaScript.

###### themes

The `themes` directory contains one or more [themes](/getting-started/glossary/#theme), each in its own subdirectory.

## Union file system

Hugo creates a union file system, allowing you to mount two or more directories to the same location. For example, let's say your home directory contains a Hugo project in one directory, and shared content in another:

```text
home/
└── user/
    ├── my-site/            
    │   ├── content/
    │   │   ├── books/
    │   │   │   ├── _index.md
    │   │   │   ├── book-1.md
    │   │   │   └── book-2.md
    │   │   └── _index.md
    │   ├── themes/
    │   │   └── my-theme/
    │   └── hugo.toml
    └── shared-content/     
        └── films/
            ├── _index.md
            ├── film-1.md
            └── film-2.md
```

You can include the shared content when you build your site using mounts. In your site configuration:

{{< code-toggle file=hugo >}}
[[module.mounts]]
source = 'content'
target = 'content'

[[module.mounts]]
source = '/home/user/shared-content'
target = 'content'
{{< /code-toggle >}}

{{% note %}}
When you overlay one directory on top of another, you must mount both directories.

Hugo does not follow symbolic links. If you need the functionality provided by symbolic links, use Hugo's union file system instead.
{{% /note %}}

After mounting, the union file system has this structure:

```text
home/
└── user/
    └── my-site/
        ├── content/
        │   ├── books/
        │   │   ├── _index.md
        │   │   ├── book-1.md
        │   │   └── book-2.md
        │   ├── films/
        │   │   ├── _index.md
        │   │   ├── film-1.md
        │   │   └── film-2.md
        │   └── _index.md
        ├── themes/
        │   └── my-theme/
        └── hugo.toml
```

{{% note %}}
When two or more files have the same path, the order of precedence follows the order of the mounts. For example, if the shared content directory contains `books/book-1.md`, it will be ignored because the project's content directory was mounted first.
{{% /note %}}

You can mount directories to `archetypes`, `assets`, `content`, `data`, `i18n`, `layouts`, and `static`. See&nbsp;[details](/hugo-modules/configuration/#module-configuration-mounts).

You can also mount directories from Git repositories using Hugo Modules. See&nbsp;[details](/hugo-modules/).

## Theme skeleton

Hugo generates a functional theme skeleton when you create a new theme. For example, this command:

```text
hugo new theme my-theme
```

Creates this directory structure (subdirectories not shown):

```text
my-theme/
├── archetypes/
├── assets/
├── content/
├── data/
├── i18n/
├── layouts/
├── static/
├── LICENSE
├── README.md
├── hugo.toml
└── theme.toml
```

Using the union file system described above, Hugo mounts each of these directories to the corresponding location in the project. When two files have the same path, the file in the project directory takes precedence. This allows you, for example, to override a theme's template by placing a copy in the same location within the project directory.

If you are simultaneously using components from two or more themes or modules, and there's a path collision, the first mount takes precedence.
