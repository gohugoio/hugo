---
title: Directory structure
description: An overview of Hugo's directory structure.
categories: []
keywords: []
weight: 30
---

Each Hugo project is a directory, with subdirectories that contribute to  content, structure, behavior, and presentation.

## Project skeleton

Hugo generates a project skeleton when you create a new project. For example, this command:

```sh
hugo new project my-project
```

Creates this directory structure:

```txt
my-project/
├── archetypes/
│   └── default.md
├── assets/
├── content/
├── data/
├── i18n/
├── layouts/
├── static/
├── themes/
└── hugo.toml         <-- project configuration
```

Depending on requirements, you may wish to organize your project configuration into subdirectories:

```txt
my-project/
├── archetypes/
│   └── default.md
├── assets/
├── config/           <-- project configuration
│   └── _default/
│       └── hugo.toml
├── content/
├── data/
├── i18n/
├── layouts/
├── static/
└── themes/
```

When you build your project, Hugo creates a `public` directory, and typically a `resources` directory as well:

```txt
my-project/
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
├── public/       <-- created when you build your project
├── resources/    <-- created when you build your project
├── static/
└── themes/
```

## Directories

Each of the subdirectories contributes to content, structure, behavior, or presentation.

`archetypes`
: The `archetypes` directory contains templates for new content. See [details](/content-management/archetypes/).

`assets`
: The `assets` directory contains global resources typically passed through an asset pipeline. This includes resources such as images, CSS, Sass, JavaScript, and TypeScript. See [details](/hugo-pipes/introduction/).

`config`
: The `config` directory contains your project configuration, possibly split into multiple subdirectories and files. For projects with minimal configuration or projects that do not need to behave differently in different environments, a single configuration file named `hugo.toml` in the root of the project is sufficient. See [details][configuration-directory].

`content`
: The `content` directory contains the markup files (typically Markdown) and page resources that comprise the content of your project. See [details](/content-management/organization/).

`data`
: The `data` directory contains data files (JSON, TOML, YAML, or XML) that augment content, configuration, localization, and navigation. See [details](/content-management/data-sources/).

`i18n`
: The `i18n` directory contains translation tables for multilingual projects. See [details](/content-management/multilingual/).

`layouts`
: The `layouts` directory contains templates to transform content, data, and resources into a complete website. See [details](/templates/).

`public`
: The `public` directory contains the published website, generated when you run the `hugo build` or `hugo server` commands. Hugo recreates this directory and its content as needed. See [details][build-your-project].

`resources`
: The `resources` directory contains cached output from Hugo's asset pipelines, generated when you run the `hugo build` or `hugo server` commands. By default this cache directory includes CSS and images. Hugo recreates this directory and its content as needed.

`static`
: The `static` directory contains files that will be copied to the `public` directory when you build your project. For example: `favicon.ico`, `robots.txt`, and files that verify website ownership. Before the introduction of [page bundles](g) and [asset pipelines][], the `static` directory was also used for images, CSS, and JavaScript.

`themes`
: The `themes` directory contains one or more [themes](g), each in its own subdirectory.

## Unified file system

Hugo creates a [unified file system](g), allowing you to mount two or more directories to the same location. For example, let's say your home directory contains a Hugo project in one directory, and shared content in another:

```tree
home/
└── user/
    ├── my-project/
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

You can include the shared content using mounts. In your project configuration:

{{< code-toggle file=hugo >}}
[[module.mounts]]
source = 'content'
target = 'content'

[[module.mounts]]
source = '/home/user/shared-content'
target = 'content'
{{< /code-toggle >}}

> [!NOTE]
> Defining a custom mount replaces the default mounting for that [component](g). To overlay an external directory on top of the project default, you must explicitly mount both.
>
> Hugo does not follow symbolic links. If you need the functionality provided by symbolic links, use Hugo's unified file system instead.

After mounting, the unified file system has this structure:

```tree
home/
└── user/
    └── my-project/
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

When two or more files share the same path, the version in the highest layer takes precedence. In the example above, if the `shared-content` directory contains `books/book-1.md`, it is ignored because the project's `content` directory is the first (highest) mount.

You can mount directories to `archetypes`, `assets`, `content`, `data`, `i18n`, `layouts`, and `static`. See [details][mounts].

You can also mount directories from Git repositories using modules. See [details](/hugo-modules/).

## Theme skeleton

Hugo generates a functional theme skeleton when you create a new theme. For example, this command:

```sh
hugo new theme my-theme
```

Creates this directory structure (subdirectories not shown):

```tree
my-theme/
├── archetypes/
├── assets/
├── content/
├── data/
├── i18n/
├── layouts/
├── static/
└── hugo.toml
```

Using the unified file system described above, Hugo mounts each of these directories to the corresponding location in the project. When two files have the same path, the file in the project directory takes precedence. This allows you, for example, to override a theme's template by placing a copy in the same location within the project directory.

If you are simultaneously using components from two or more themes or modules, and there's a path collision, the first mount takes precedence.

[asset pipelines]: /hugo-pipes/introduction/
[build-your-project]: /getting-started/usage/#build-your-project
[configuration-directory]: /configuration/introduction/#configuration-directory
[mounts]: /configuration/module/#mounts
