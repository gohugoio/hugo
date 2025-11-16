---
title: Introduction
description: Configure your site using files, directories, and environment variables.
categories: []
keywords: []
weight: 10
---

## Sensible defaults

Hugo offers many configuration options, but its defaults are often sufficient. A new site requires only these settings:

{{< code-toggle file=hugo >}}
baseURL = 'https://example.org/'
languageCode = 'en-us'
title = 'My New Hugo Site'
{{< /code-toggle >}}

Only define settings that deviate from the defaults. A smaller configuration file is easier to read, understand, and debug. Keep your configuration concise.

> [!note]
> The best configuration file is a short configuration file.

## Configuration file

Create a site configuration file in the root of your project directory, naming it `hugo.toml`, `hugo.yaml`, or `hugo.json`, with that order of precedence.

```text
my-project/
└── hugo.toml
```

> [!note]
> For versions v0.109.0 and earlier, the site configuration file was named `config`. While you can still use this name, it's recommended to switch to the newer naming convention, `hugo`.

A simple example:

{{< code-toggle file=hugo >}}
baseURL = 'https://example.org/'
languageCode = 'en-us'
title = 'ABC Widgets, Inc.'
[params]
subtitle = 'The Best Widgets on Earth'
[params.contact]
email = 'info@example.org'
phone = '+1 202-555-1212'
{{< /code-toggle >}}

To use a different configuration file when building your site, use the `--config` flag:

```sh
hugo --config other.toml
```

Combine two or more configuration files, with left-to-right precedence:

```sh
hugo --config a.toml,b.yaml,c.json
```

> [!note]
> See the specifications for each file format: [TOML], [YAML], and [JSON].

## Configuration directory

Instead of a single site configuration file, split your configuration by [environment](g), root configuration key, and language. For example:

```text
my-project/
└── config/
    ├── _default/
    │   ├── hugo.toml
    │   ├── menus.en.toml
    │   ├── menus.de.toml
    │   └── params.toml
    └── production/
        └── params.toml
```

The root configuration keys are {{< root-configuration-keys >}}.

> [!note]
> You must define `cascade` tables in the root configuration file. You cannot define `cascade` tables in a dedicated file. See issue [#12899] for details.

[#12899]: https://github.com/gohugoio/hugo/issues/12899

### Omit the root key

When splitting the configuration by root key, omit the root key in the component file. For example, these are equivalent:

{{< code-toggle file=config/_default/hugo >}}
[params]
foo = 'bar'
{{< /code-toggle >}}

{{< code-toggle file=config/_default/params >}}
foo = 'bar'
{{< /code-toggle >}}

### Recursive parsing

Hugo parses the `config` directory recursively, allowing you to organize the files into subdirectories. For example:

```text
my-project/
└── config/
    └── _default/
        ├── navigation/
        │   ├── menus.de.toml
        │   └── menus.en.toml
        └── hugo.toml
```

### Example

```text
my-project/
└── config/
    ├── _default/
    │   ├── hugo.toml
    │   ├── menus.en.toml
    │   ├── menus.de.toml
    │   └── params.toml
    ├── production/
    │   ├── hugo.toml
    │   └── params.toml
    └── staging/
        ├── hugo.toml
        └── params.toml
```

Considering the structure above, when running `hugo --environment staging`, Hugo will use every setting from `config/_default` and merge `staging`'s on top of those.

Let's take an example to understand this better. Let's say you are using Google Analytics for your website. This requires you to specify a [Google tag ID] in your site configuration:

{{< code-toggle file=hugo >}}
[services.googleAnalytics]
ID = 'G-XXXXXXXXX'
{{< /code-toggle >}}

Now consider the following scenario:

1. You don't want to load the analytics code when running `hugo server`.
1. You want to use different Google tag IDs for your production and staging environments. For example:
    - `G-PPPPPPPPP` for production
    - `G-SSSSSSSSS` for staging

To satisfy these requirements, configure your site as follows:

1. `config/_default/hugo.toml`
    - Exclude the `services.googleAnalytics` section. This will prevent loading of the analytics code when you run `hugo server`.
    - By default, Hugo sets its `environment` to `development` when running `hugo server`. In the absence of a `config/development` directory, Hugo uses the `config/_default` directory.
1. `config/production/hugo.toml`
    - Include this section only:

      {{< code-toggle file=hugo >}}
      [services.googleAnalytics]
      ID = 'G-PPPPPPPPP'
      {{< /code-toggle >}}

    - You do not need to include other parameters in this file. Include only those parameters that are specific to your production environment. Hugo will merge these parameters with the default configuration.
    - By default, Hugo sets its `environment` to `production` when running `hugo`. The analytics code will use the `G-PPPPPPPPP` tag ID.

1. `config/staging/hugo.toml`

    - Include this section only:

      {{< code-toggle file=hugo >}}
      [services.googleAnalytics]
      ID = 'G-SSSSSSSSS'
      {{< /code-toggle >}}

    - You do not need to include other parameters in this file. Include only those parameters that are specific to your staging environment. Hugo will merge these parameters with the default configuration.
    - To build your staging site, run `hugo --environment staging`. The analytics code will use the `G-SSSSSSSSS` tag ID.

## Merge configuration settings

Hugo merges configuration settings from themes and modules, prioritizing the project's own settings. Given this simplified project structure with two themes:

```text
project/
├── themes/
│   ├── theme-a/
│   │   └── hugo.toml
│   └── theme-b/
│       └── hugo.toml
└── hugo.toml
```

and this project-level configuration:

{{< code-toggle file=hugo >}}
baseURL = 'https://example.org/'
languageCode = 'en-us'
title = 'My New Hugo Site'
theme = ['theme-a','theme-b']
{{< /code-toggle >}}

Hugo merges settings in this order:

1. Project configuration (`hugo.toml` in the project root)
1. `theme-a` configuration
1. `theme-b` configuration

The `_merge` setting within each top-level configuration key controls _which_ settings are merged and _how_ they are merged.

The value for `_merge` can be one of:

none
: No merge.

shallow
: Only add values for new keys.

deep
: Add values for new keys, merge existing.

Note that you don't need to be so verbose as in the default setup below; a `_merge` value higher up will be inherited if not set.

{{< code-toggle file=hugo dataKey="config_helpers.mergeStrategy" skipHeader=true />}}

## Environment variables

You can also configure settings using operating system environment variables:

```sh
export HUGO_BASEURL=https://example.org/
export HUGO_ENABLEGITINFO=true
export HUGO_ENVIRONMENT=staging
hugo
```

The above sets the [`baseURL`], [`enableGitInfo`], and [`environment`] configuration options and then builds your site.

> [!note]
> An environment variable takes precedence over the values set in the configuration file. This means that if you set a configuration value with both an environment variable and in the configuration file, the value in the environment variable will be used.

Environment variables simplify configuration for [CI/CD](g) deployments like GitHub Pages, GitLab Pages, and Netlify by allowing you to set values directly within their respective configuration and workflow files.

> [!note]
> Environment variable names must be prefixed with `HUGO_`.
>
> To set custom site parameters, prefix the name with `HUGO_PARAMS_`.

For snake_case variable names, the standard `HUGO_` prefix won't work. Hugo infers the delimiter from the first character following `HUGO`. This allows for variations like `HUGOxPARAMSxAPI_KEY=abcdefgh` using any [permitted delimiter].

In addition to configuring standard settings, environment variables may be used to override default values for certain internal settings:

DART_SASS_BINARY
: (`string`) The absolute path to the Dart Sass executable. By default, Hugo searches for the executable in each of the paths in the `PATH` environment variable.

HUGO_FILE_LOG_FORMAT
: (`string`) A format string for the file path, line number, and column number displayed when reporting errors, or when calling the `Position` method from a shortcode or Markdown render hook. Valid tokens are `:file`, `:line`, and `:col`. Default is `:file::line::col`.

HUGO_MEMORYLIMIT
: {{< new-in 0.123.0 />}}
: (`int`) The maximum amount of system memory, in gigabytes, that Hugo can use while rendering your site. Default is 25% of total system memory. Note that `HUGO_MEMORYLIMIT` is a "best effort" setting. Don't expect Hugo to build a million pages with only 1 GB of memory. You can get more information about how this behaves during the build by building with `hugo --logLevel info` and look for the `dynacache` label.

HUGO_NUMWORKERMULTIPLIER
: (`int`) The number of workers used in parallel processing. Default is the number of logical CPUs.

## Current configuration

Display the complete site configuration with:

```sh
hugo config
```

Display a specific configuration setting with:

```sh
hugo config | grep [key]
```

Display the configured file mounts with:

```sh
hugo config mounts
```

[`baseURL`]: /configuration/all#baseurl
[`enableGitInfo`]: /configuration/all#enablegitinfo
[`environment`]: /configuration/all#environment
[Google tag ID]: https://support.google.com/tagmanager/answer/12326985?hl=en
[JSON]: https://datatracker.ietf.org/doc/html/rfc7159
[permitted delimiter]: https://pubs.opengroup.org/onlinepubs/000095399/basedefs/xbd_chap08.html
[TOML]: https://toml.io/en/latest
[YAML]: https://yaml.org/spec/
