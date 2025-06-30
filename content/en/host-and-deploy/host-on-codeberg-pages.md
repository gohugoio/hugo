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

To use Codeberg's Woodpecker CI, you need to have or [request] access to it, as well as add a `.woodpecker.yaml` file in the root of your project. A template and additional instructions are available in the official [examples repository].

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

Two examples of such a file are provided below.

The first file should work for automatically building your website from the source branch (`main` in this case) and committing the result to the target branch (`pages`). Without changes, this file should make your built website accessible under `https://<YourUsername>.codeberg.page/<YourWebsiteRepositoryName>/`:

```yaml {file=".forgejo/workflows/hugo.yaml" copy=true}
name: Deploy Hugo site to Pages

on:
  # Runs on pushes targeting the default branch
  push:
    branches:
      # If you want to build from a different branch, change it here.
      - main
  # Allows you to run this workflow manually from the Actions tab
  workflow_dispatch:

jobs:
  build:
    # You can find the list of available runners on https://codeberg.org/actions/meta, or run one yourself.
    runs-on: codeberg-tiny-lazy
    container:
      # Specify "hugomods/hugo:exts" if you want to always use the latest version of Hugo for building.
      image: "hugomods/hugo:exts-0.147.9"
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
      - name: Checkout the target branch and clean it up
        # If you want to commit to a branch other than "pages", change the two references below, as well as the reference in the last step. 
        run: |
          git checkout pages || git switch --orphan pages && \
          rm -Rfv $(ls -A | egrep -v '^(\.git|LICENSE)$')
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

The second file implements a more complex scenario: having your website sources in one repository and the resulting static website in another repository (in this case, `pages`). If you want Codeberg to make your website available at the root of your pages subdomain (`https://<YourUsername>.codeberg.page/`), you have to push that website to the default branch of your repository named `pages`.

Since this action involves more than one repository, it will require a bit more preparation:
1. Create the target repository. Name it `pages`.
2. Generate a new SSH key. *Do not* use any of your own SSH keys for this, but generate one for this specific task only. On Linux, BSD, and, likely, other operating systems, you can open a terminal emulator and run the following command to generate the key:
   ```shell
   ssh-keygen -f pagesbuild -P ""
   ```
   This will generate two files in your current directory: `pagesbuild` (private key) and `pagesbuild.pub` (public key).
3. Add the newly generated public key as a deploy key to your `pages` repository: navigate to its Settings, click on "Deploy keys" in the left menu, click the "Add deploy key" button, give it a name (e.g. "Actions deploy key"), paste the contents of the **public** key file (`pagesbuild.pub`) to the Content field, tick the "Enable write access" checkbox, then submit the form.
4. Navigate back to your source repository settings, expand the "Actions" menu and click on "Secrets". Then click "Add Secret", enter "DEPLOY_KEY" as the secret name and paste the contents of the newly generated **private** key file (`pagesbuild`) into the Value field.
5. Navigate to the "Variables" submenu of the "Actions" menu and add the following variables:

    | Name                | Value                                                                            |
    |---------------------|----------------------------------------------------------------------------------|
    | `TARGET_REPOSITORY` | `<YourUsername>/pages`                                                           |
    | `TARGET_BRANCH`     | `main` (enter the default branch name of the `pages` repo here)                  |
    | `SSH_KNOWN_HOSTS`   | (paste the output you get by running `ssh-keyscan codeberg.org` in the terminal) |

Once you've done all of the above, commit the following file to your repository as `.forgejo/workflows/hugo.yaml`. As you can see, the `deploy` job of this workflow is slightly different from the file above:

```yaml {file=".forgejo/workflows/hugo.yaml" copy=true}
name: Deploy Hugo site to Pages

on:
  # Runs on pushes targeting the default branch
  push:
    branches:
      # If you want to build from a different branch, change it here.
      - main
  # Allows you to run this workflow manually from the Actions tab
  workflow_dispatch:

jobs:
  build:
    runs-on: codeberg-tiny-lazy
    container:
      # Specify "hugomods/hugo:exts" if you want to always use the latest version of Hugo for building.
      image: "hugomods/hugo:exts-0.147.9"
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
            --minify \
            --source ${PWD} \
            --destination ${PWD}/public/
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
          repository: ${{ vars.TARGET_REPOSITORY }}
          ref: ${{ vars.TARGET_BRANCH }}
          submodules: recursive
          fetch-depth: 0
          ssh-key: ${{ secrets.DEPLOY_KEY }}
          ssh-known-hosts: ${{ vars.SSH_KNOWN_HOSTS }}
      - name: Remove all files
        run: |
          rm -Rfv $(ls -A | egrep -v '^(\.git|LICENSE)$')
      - name: Download generated files
        uses: https://code.forgejo.org/actions/download-artifact@v3
        with:
          name: Generated files
      - name: Commit and push the website
        run: |
          git config user.email codeberg-ci && \
          git config user.name "Codeberg CI" && \
          git add -v . && \
          git commit -v --allow-empty --message "Codeberg build for ${GITHUB_SHA}" && \
          git push -v origin ${{ vars.TARGET_BRANCH }}
```

Once you commit one of the two files to your website source repository, you should see your first automated build firing up pretty soon. You can also trigger it manually by navigating to the **Actions** section of your repository web page, choosing **hugo.yaml** on the left and clicking on **Run workflow**.

## Other resources

- [Codeberg Pages](https://codeberg.page/)
- [Codeberg Pages official documentation](https://docs.codeberg.org/codeberg-pages/)
