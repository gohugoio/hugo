---
title: Host on Codeberg Pages
description: Codeberg Pages is a hosting service available to Codeberg users which allows you to publish your Hugo website, as well as websites managed by other static site generators.
categories: [hosting and deployment]
keywords: [hosting,codeberg]
menu:
  docs:
    parent: hosting-and-deployment
toc: true
aliases: [/tutorials/hosting-on-codeberg/]
---

## Assumptions

* Working familiarity with [Git](https://git-scm.com/) for version control
* Completion of the Hugo [Quick Start]
* A [Codeberg account](https://codeberg.org/user/login)
* A Hugo website on your local machine that you are ready to publish

Any and all mentions of `<YourUsername>` refer to your actual Codeberg username and must be substituted accordingly. Likewise, `<YourWebsite>` represents your actual website name.

## BaseURL

The `baseURL` in your [site configuration](/getting-started/configuration/) must reflect the full URL provided by Codeberg Pages if using the default address (e.g. `https://<YourUsername>.codeberg.page/`). If you want to use another domain, follow the instructions in the [custom domain section](https://docs.codeberg.org/codeberg-pages/using-custom-domain/) of the official documentation.

For more details regarding the URL of your deployed website, refer to Codeberg Pages' [quickstart instructions](https://codeberg.page/).

## Manual deployment

Create a public repository on your Codeberg account titled `pages` or create a branch of the same name in an existing public repository. Finally, push the contents of Hugo's output directory (by default, `public`) to it. Here's an example:

```sh
# build the website
hugo

# access the output directory
cd public

# initialize new git repository
git init

# commit and push code to main branch
git add .
git commit -m "Initial commit"
git remote add origin https://codeberg.org/<YourUsername>/pages.git
git push -u origin main
```

## Automated deployment

In order to automatically deploy your Hugo website, you need to have or [request](https://codeberg.org/Codeberg-e.V./requests/issues/new?template=ISSUE_TEMPLATE%2fWoodpecker-CI.yaml) access to Codeberg's CI, as well as add a `.woodpecker.yaml` file in the root of your project. A template and additional instructions are available in the official [examples repository](https://codeberg.org/Codeberg-CI/examples/src/branch/main/Hugo/.woodpecker.yaml).

In this case, you must create a public repository on Codeberg (e.g. `<YourWebsite>`) and push your local project to it. Here's an example:

```sh
# initialize new git repository
git init

# add /public directory to our .gitignore file
echo "/public" >> .gitignore

# commit and push code to main branch
git add .
git commit -m "Initial commit"
git remote add origin https://codeberg.org/<YourUsername>/<YourWebsite>.git
git push -u origin main
```

Your project will then be built and deployed by Codeberg's CI.

## Other resources

* [Codeberg Pages](https://codeberg.page/)
* [Codeberg Pages official documentation](https://docs.codeberg.org/codeberg-pages/)

[Quick Start]: /getting-started/quick-start/
