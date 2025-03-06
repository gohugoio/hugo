---
title: Host on SourceHut Pages
description: Host your site on SourceHut Pages.
categories: []
keywords: []
aliases: [/hosting-and-deployment/hosting-on-sourcehut/]
---

## Assumptions

- Working familiarity with [Git] or [Mercurial] for version control
- Completion of the Hugo [Quick Start]
- A [SourceHut account]
- A Hugo website on your local machine that you are ready to publish

[Git]: https://git-scm.com/
[Mercurial]: https://www.mercurial-scm.org/
[SourceHut account]: https://meta.sr.ht/login
[Quick Start]: /getting-started/quick-start/

Any and all mentions of `<YourUsername>` refer to your actual SourceHut username and must be substituted accordingly.

## BaseURL

The [`baseURL`] in your site configuration must reflect the full URL provided by SourceHut Pages if you are using the default address (e.g. `https://<YourUsername>.srht.site/`). If you want to use another domain, check the [custom domain section] of the official documentation.

[`baseURL`]: /configuration/all/#baseurl
[custom domain section]: https://srht.site/custom-domains

## Manual deployment

This method does not require a paid account. To proceed you will need to create a [SourceHut personal access token] and install and configure the [hut]CLI tool:

[SourceHut personal access token]: https://meta.sr.ht/oauth2/personal-token/
[hut]: https://sr.ht/~xenrox/hut/

```sh
hugo
tar -C public -cvz . > site.tar.gz
hut init
hut pages publish -d <YourUsername>.srht.site site.tar.gz
```

A TLS certificate will be automatically obtained for you, and your new website will be available at `https://<YourUsername>.srht.site/` (or the provided custom domain).

## Automated deployment

This method requires a paid account and relies on the SourceHut build system.

First, define your [build manifest] by creating a `.build.yml` file in the root of your project. The following is a bare-bones template:

[build manifest]: https://man.sr.ht/builds.sr.ht/#build-manifests

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
    hugo
    tar -C public -cvz . > ../site.tar.gz
- upload: |
    hut pages publish -d $site site.tar.gz
```

Now what's left is creating a repository titled `<YourUsername>.srht.site` (or your custom domain, if applicable) and pushing your local project. Here's an example using Git:

```sh
# initialize new git repository
git init

# add /public directory to our .gitignore file
echo "/public" >> .gitignore

# commit and push code to main branch
git add .
git commit -m "Initial commit"
git remote add origin https://git.sr.ht/~<YourUsername>/<YourUsername>.srht.site
git push -u origin main
```

You can now follow the build progress of your page at `https://builds.sr.ht/`.

After the build has passed, a TLS certificate will be automatically obtained for you and your new website will be available at `https://<YourUsername>.srht.site/` (or the provided custom domain).

## Other resources

- [SourceHut Pages](https://srht.site/)
- [SourceHut Builds user manual](https://man.sr.ht/builds.sr.ht/)
