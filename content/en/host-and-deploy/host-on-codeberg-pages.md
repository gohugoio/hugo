---
title: Host on Codeberg Pages
description: Host your site on Codeberg Pages.
categories: []
keywords: []
aliases: [/hosting-and-deployment/hosting-on-codeberg/]
---

## Assumptions

- Working familiarity with [Git] for version control
- Completion of the Hugo [Quick Start]
- A [Codeberg account]
- A Hugo website on your local machine that you are ready to publish

[Codeberg account]: https://codeberg.org/user/login/
[Git]: https://git-scm.com/
[Quick Start]: /getting-started/quick-start/

Any and all mentions of `<YourUsername>` refer to your actual Codeberg username and must be substituted accordingly. Likewise, `<YourWebsite>` represents your actual website name.

## BaseURL

The [`baseURL`] in your site configuration must reflect the full URL provided by Codeberg Pages if using the default address (e.g. `https://<YourUsername>.codeberg.page/`). If you want to use another domain, follow the instructions in the [custom domain section] of the official documentation.

[`baseURL`]: /configuration/all/#baseurl
[custom domain section]: https://docs.codeberg.org/codeberg-pages/using-custom-domain/

For more details regarding the URL of your deployed website, refer to Codeberg Pages' [quickstart instructions].

[quickstart instructions]: https://codeberg.page/

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

## Automated deployment using Woodpecker CI

There are two methods you can use to deploy your Hugo website to Codeberg automatically. These are: Woodpecker CI and Forgejo Actions.

To use Codeberg's Woodperker CI, you need to have or [request] access to it, as well as add a `.woodpecker.yaml` file in the root of your project. A template and additional instructions are available in the official [examples repository].

[request]: https://codeberg.org/Codeberg-e.V./requests/issues/new?template=ISSUE_TEMPLATE%2fWoodpecker-CI.yaml
[examples repository]: https://codeberg.org/Codeberg-CI/examples/src/branch/main/Hugo/.woodpecker.yaml

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

Your project will then be built and deployed by Codeberg's Woodpecker CI.

## Automated deployment using Forgejo Actions

The other way to deploy your website to Codeberg pages automatically is to make use of Forgejo Actions. Actions need a _runner_ to work, and Codeberg has [great documentation] on how to set one up yourself. However, Codeberg provides a [handful of humble runners] themselves (they say this feature is in "open alpha"), which actually seem powerful enough to build at least relatively simple websites.

[great documentation]: https://docs.codeberg.org/ci/actions/
[handful of humble runners]: https://codeberg.org/actions/meta

To deploy your website this way, you don't need to request any access. All you need to do is enable actions in your repository settings (see the documentation link above) and add a workflow configuration file, for example, `hugo.yaml`, to the `.forgejo/workflows/` directory in your website's source repository.

An example of such a file is provided below. It should work for automatically building your website from the `main` branch, and committing the result to the `pages` branch to have it accessible under `https://<YourUsername>.codeberg.page/<YourWebsiteRepositoryName>/`.

**Please note** however that this is a slightly different approach than the one described above, where you use a separate `pages` repository and the resulting website is available at `https://<YourUsername>.codeberg.page/` directly. This file should be a good starting point to explore that path though:

```yaml {file=".forgejo/workflows/hugo.yaml" copy=true}
name: Deploy Hugo site to Pages

on:
  # Runs on pushes targeting the default branch
  push:
    branches:
      - main
  # Allows you to run this workflow manually from the Actions tab
  workflow_dispatch:

jobs:
  build:
    # You can find the list of available runners on https://codeberg.org/actions/meta, or run one yourself.
    runs-on: codeberg-tiny-lazy
    container:
      image: "hugomods/hugo:exts-0.147.3"
    steps:
      - name: Clone the repository
        uses: https://code.forgejo.org/actions/checkout@v4
        with:
          submodules: recursive
          fetch-depth: 0
      - name: Generate static files with Hugo
        env:
          # For maximum backward compatibility with Hugo modules
          HUGO_ENVIRONMENT: production
          HUGO_ENV: production
        run: |
          hugo \
            --gc \
            --minify
      - name: Upload generated files
        uses: https://code.forgejo.org/actions/upload-artifact@v3
        with:
          name: Generated files
          path: public/
  deploy:
    needs: [ build ]
    runs-on: codeberg-tiny-lazy
    steps:
      - name: Clone the repository
        uses: https://code.forgejo.org/actions/checkout@v4
        with:
          submodules: recursive
          fetch-depth: 0
      - name: Checkout the pages branch and clean it up
        run: |
          git checkout pages || git switch --orphan pages && \
          rm -Rf $(ls -A | grep -v ^\.git$)
      - name: Download generated files
        uses: https://code.forgejo.org/actions/download-artifact@v3
        with:
          name: Generated files
      - name: Publish the website
        run: |
          git config user.email codeberg-ci && \
          git config user.name "Codeberg CI" && \
          git add . && \
          git commit --allow-empty --message "Codeberg build for ${GITHUB_SHA}" && \
          git push origin pages
```

Once you commit this file to your website source repository, you should see a build firing up pretty soon. You can also trigger it manually by navigating to the **Actions** section of your repository web page, choosing **hugo.yaml** on the left and clicking on **Run workflow**.

## Other resources

- [Codeberg Pages](https://codeberg.page/)
- [Codeberg Pages official documentation](https://docs.codeberg.org/codeberg-pages/)
