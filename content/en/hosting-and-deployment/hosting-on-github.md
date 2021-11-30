---
title: Host on GitHub
linktitle: Host on GitHub
description: Deploy Hugo as a GitHub Pages project or personal/organizational site and automate the whole process with Github Action Workflow
date: 2014-03-21
publishdate: 2014-03-21
categories: [hosting and deployment]
keywords: [github,git,deployment,hosting]
authors: [Spencer Lyon, Gunnar Morling]
menu:
  docs:
    parent: "hosting-and-deployment"
    weight: 30
weight: 30
sections_weight: 30
toc: true
aliases: [/tutorials/github-pages-blog/]
---

GitHub provides free and fast static hosting over SSL for personal, organization, or project pages directly from a GitHub repository via its [GitHub Pages service][] and automating development workflows and build with [GitHub Actions].

## Assumptions

1. You have Git 2.8 or greater [installed on your machine][installgit].
2. You have a GitHub account. [Signing up][ghsignup] for GitHub is free.
3. You have a ready-to-publish Hugo website or have at least completed the [Quick Start][].

## Types of GitHub Pages

There are two types of GitHub Pages:

- User/Organization Pages (`https://<USERNAME|ORGANIZATION>.github.io/`)
- Project Pages (`https://<USERNAME|ORGANIZATION>.github.io/<PROJECT>/`)

Please refer to the [GitHub Pages documentation][ghorgs] to decide which type of site you would like to create as it will determine which of the below methods to use.

## GitHub User or Organization Pages

As mentioned in the [GitHub Pages documentation][ghorgs], you can host a user/organization page in addition to project pages. Here are the key differences in GitHub Pages websites for Users and Organizations:

1. You must use a `<USERNAME>.github.io` to host your **generated** content
2. Content from the `main` branch will be used to publish your GitHub Pages site

This is a much simpler setup as your Hugo files and generated content are published into two different repositories.

## Build Hugo With GitHub Action

GitHub executes your software development workflows. Everytime you push your code on the Github repository, Github Actions will build the site automatically.

Create a file in `.github/workflows/gh-pages.yml` containing the following content (based on [actions-hugo](https://github.com/marketplace/actions/hugo-setup)):

```yml
name: github pages

on:
  push:
    branches:
      - main  # Set a branch to deploy
  pull_request:

jobs:
  deploy:
    runs-on: ubuntu-20.04
    steps:
      - uses: actions/checkout@v2
        with:
          submodules: true  # Fetch Hugo themes (true OR recursive)
          fetch-depth: 0    # Fetch all history for .GitInfo and .Lastmod

      - name: Setup Hugo
        uses: peaceiris/actions-hugo@v2
        with:
          hugo-version: 'latest'
          # extended: true

      - name: Build
        run: hugo --minify

      - name: Deploy
        uses: peaceiris/actions-gh-pages@v3
        if: github.ref == 'refs/heads/main'
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
```

For more advanced settings [actions-hugo](https://github.com/marketplace/actions/hugo-setup) and [actions-gh-pages](https://github.com/marketplace/actions/github-pages-action).

## Github pages setting
By default, the GitHub action pushes the generated content to the `gh-pages` branch. This means GitHub has to serve your `gh-pages` branch as a GitHub Pages branch. You can change this setting by going to Settings > GitHub Pages, and change the source branch to `gh-pages`.

## Change baseURL in config.toml
Don't forget to rename your `baseURL` in `config.toml` with the value `https://<USERNAME>.github.io` for your user repository or `https://<USERNAME>.github.io/<REPOSITORY_NAME>` for a project repository.

Unless this is present in your `config.toml`, your website won't work.

## Use a Custom Domain

If you'd like to use a custom domain for your GitHub Pages site, create a file `static/CNAME`. Your custom domain name should be the only contents inside `CNAME`. Since it's inside `static`, the published site will contain the CNAME file at the root of the published site, which is a requirement of GitHub Pages.

Refer to the [official documentation for custom domains][domains] for further information.

[config]: /getting-started/configuration/
[domains]: https://help.github.com/articles/using-a-custom-domain-with-github-pages/
[ghorgs]: https://help.github.com/articles/user-organization-and-project-pages/#user--organization-pages
[ghpfromdocs]: https://help.github.com/articles/configuring-a-publishing-source-for-github-pages/
[ghsignup]: https://github.com/join
[GitHub Pages service]: https://help.github.com/articles/what-is-github-pages/
[installgit]: https://git-scm.com/downloads
[orphan branch]: https://git-scm.com/docs/git-checkout/#Documentation/git-checkout.txt---orphanltnewbranchgt
[Quick Start]: /getting-started/quick-start/
[submodule]: https://github.com/blog/2104-working-with-submodules
[worktree feature]: https://git-scm.com/docs/git-worktree
[GitHub Actions]: https://docs.github.com/en/actions
