---
title: Host on SourceHut Pages
description: Host your site on SourceHut Pages.
categories: []
keywords: []
aliases: [/hosting-and-deployment/hosting-on-sourcehut/]
---

Use these instructions to host your site on SourceHut Pages using either manual deployment or the SourceHut build system.

{{% include "/_common/gitignore-public.md" %}}

## Prerequisites

- Working familiarity with [Git][] or [Mercurial][] for version control
- Completion of the Hugo [Quick Start][]
- A [SourceHut account][]
- A Hugo website on your local machine that you are ready to publish

Any and all mentions of `<YourUsername>` refer to your actual SourceHut username and must be substituted accordingly.

## BaseURL

The [`baseURL`][] in your project configuration must reflect the full URL provided by SourceHut Pages if you are using the default address (e.g. `https://<YourUsername>.srht.site/`). If you want to use another domain, check the [custom domain section][] of the official documentation.

## Manual deployment

This method does not require a paid account. To proceed you will need to create a [SourceHut personal access token][] and install and configure the [hut][] CLI tool:

```sh
hugo build
tar -C public -cvz . > site.tar.gz
hut init
hut pages publish -d <YourUsername>.srht.site site.tar.gz
```

A TLS certificate will be automatically obtained for you, and your new website will be available at `https://<YourUsername>.srht.site/` (or the provided custom domain).

## Automated deployment

This method requires a paid account and relies on the SourceHut build system.

First, define your [build manifest][] by creating a `.build.yml` file in the root of your project. The following is a bare-bones template:

```yaml {file=".build.yml" copy=true}
image: alpine/edge
packages:
  - hugo
  - hut
oauth: pages.sr.ht/PAGES:RW
environment:
  site: <YourUsername>.srht.site
tasks:
- package: |
    cd $site
    hugo build
    tar -C public -cvz . > ../site.tar.gz
- upload: |
    hut pages publish -d $site site.tar.gz
```

If your site requires [Dart Sass][] to transpile Sass to CSS, set the DART_SASS_VERSION to the [latest version number][] and include the Dart Sass installation lines before running the Hugo build step. Note that for Alpine, the `linux-x64-musl` version is used.

```yaml {file=".build.yml" copy=true}
image: alpine/edge
packages:
  - hugo
  - hut
  - curl # For Dart Sass installation
oauth: pages.sr.ht/PAGES:RW
environment:
  site: <YourUsername>.srht.site
tasks:
- package: |
    DART_SASS_VERSION=1.101.0
    mkdir -p $HOME/.local
    curl -L https://github.com/sass/dart-sass/releases/download/${DART_SASS_VERSION}/dart-sass-${DART_SASS_VERSION}-linux-x64-musl.tar.gz -o dart-sass.tar.gz
    tar -xzf dart-sass.tar.gz -C $HOME/.local
    rm dart-sass.tar.gz
    chmod -R +x $HOME/.local/dart-sass/src
    export PATH="$HOME/.local/dart-sass:$PATH"
    sass --version # Verify installation
    cd $site
    hugo build
    tar -C public -cvz . > ../site.tar.gz
- upload: |
    hut pages publish -d $site site.tar.gz
```

Create a repository titled `<YourUsername>.srht.site` (or your custom domain, if applicable) and push your local project to the repository.

You can now follow the build progress of your page at `https://builds.sr.ht/`.

After the build has passed, a TLS certificate will be automatically obtained for you and your new website will be available at `https://<YourUsername>.srht.site/` (or the provided custom domain).

## Other resources

- [SourceHut Pages][]
- [SourceHut Builds user manual][]

[Dart Sass]: https://gohugo.io/functions/css/sass/#dart-sass
[Git]: https://git-scm.com/
[Mercurial]: https://www.mercurial-scm.org/
[Quick Start]: /getting-started/quick-start/
[SourceHut Builds user manual]: https://man.sr.ht/builds.sr.ht/
[SourceHut Pages]: https://srht.site/
[SourceHut account]: https://meta.sr.ht/login
[SourceHut personal access token]: https://meta.sr.ht/oauth2/personal-token
[`baseURL`]: /configuration/all/#baseurl
[build manifest]: https://man.sr.ht/builds.sr.ht/#build-manifests
[custom domain section]: https://srht.site/custom-domains
[hut]: https://sr.ht/~xenrox/hut/
[latest version number]: https://github.com/sass/dart-sass/releases
