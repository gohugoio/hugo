---
title: Host on Azure Static Web Apps
description: Host your project on Azure Static Web Apps.
categories: []
keywords: []
aliases: [/hosting-and-deployment/hosting-on-azure-static-web-apps/]
---

Use these instructions to enable continuous deployment from a GitHub repository. The same general steps apply for other Git providers such as GitLab or Bitbucket.

{{% include "/_common/gitignore-public.md" %}}

## Prerequisites

Please complete the following tasks before continuing:

1. [Create](https://signup.live.com/) a Microsoft account.
1. [Create](https://azure.microsoft.com/free/) an Azure account.
1. [Log in](https://portal.azure.com/) to the Azure Portal.
1. [Create](https://github.com/signup) a GitHub account.
1. [Log in](https://github.com/login) to your GitHub account.
1. [Create](https://github.com/new) a GitHub repository for your project.
1. [Create](https://git-scm.com/docs/git-init) a local Git repository for your project with a [remote][] reference to your GitHub repository.
1. Create a Hugo project within your local Git repository and test it with the `hugo server` command.
1. Commit the changes to your local Git repository and push to your GitHub repository.

## Procedure

Step 1
: Create an Azure Static Web App.

  1. Go to [Static Web Apps][] in the Azure Portal.
  1. Press the **Create** button.

      ![screen capture](azure-01.png)

  1. Under **Project details**, select your Subscription and choose or create a Resource Group.

      ![screen capture](azure-02.png)

  1. Under **Static Web App details**, enter a Name for your site.

      ![screen capture](azure-03.png)

  1. Under **Hosting plan**, select **Free**.

      ![screen capture](azure-04.png)

  1. Under **Deployment details**, select **Other** as the deployment source, then press the **Review + create** button. This allows deployment using a GitHub Actions token without auto-generating default workflow files.

      ![screen capture](azure-05.png)

  1. Wait for the validation to complete, then press the **Create** button.

      ![screen capture](azure-06.png)

  1. Once the deployment is complete, press the **Go to resource** button.

      ![screen capture](azure-07.png)

  1. Copy the assigned URL to your clipboard.

      ![screen capture](azure-08.png)

  1. In the project configuration file in the root of your local Git repository, set the [`baseURL`][] to the assigned URL as shown below.

      {{< code-toggle file=hugo >}}
      baseURL = 'https://salmon-desert-04c512910.4.azurestaticapps.net/'
      locale  = 'en-US'
      title   = 'Hosting Test - Azure'
      {{< /code-toggle >}}

  1. Click the **Manage deployment token** link at the top of the page, and copy the deployment token to your clipboard.

      ![screen capture](azure-09.png)

Step 2
: Add the deployment token to GitHub Secrets.

  1. Go to your GitHub repository.
  1. Navigate to **Settings** > **Secrets and variables** > **Actions**.
  1. Click the **New repository secret** button.
  1. Enter `AZURE_STATIC_WEB_APPS_API_TOKEN` for the Name.
  1. Paste the deployment token into the Secret field.
  1. Press the **Add secret** button.

Step 3
: Create a `hugo.yaml` file in the `.github/workflows` directory, adjusting the tool versions and time zone as needed.

  ```yaml {file=".github/workflows/hugo.yaml" copy=true}
  name: Build and deploy
  env:
    # Define tool versions
    DART_SASS_VERSION: 1.104.0
    GO_VERSION: 1.27.0
    HUGO_VERSION: 0.166.0
    NODE_VERSION: 24.20.0

    # Set the build time zone
    TZ: Europe/Oslo
  on:
    push:
      branches:
        - main
    workflow_dispatch:
  permissions:
    contents: read
  concurrency:
    group: deployment
    cancel-in-progress: false
  defaults:
    run:
      shell: bash
  jobs:
    build:
      runs-on: ubuntu-latest
      steps:
        - name: Checkout
          uses: actions/checkout@v7
          with:
            submodules: recursive
            fetch-depth: 0
            lfs: false

        - name: Create a local tools directory
          run: |
            mkdir -p "${HOME}/.local"

        - name: Install Go
          if: hashFiles('go.mod') != ''
          uses: actions/setup-go@v7
          with:
            go-version: ${{ env.GO_VERSION }}
            cache: false

        - name: Install Node.js
          if: hashFiles('package-lock.json') != ''
          uses: actions/setup-node@v7
          with:
            node-version: ${{ env.NODE_VERSION }}

        - name: Install Dart Sass
          run: |
            echo "Installing Dart Sass ${DART_SASS_VERSION}..."
            curl -sfL --output-dir "${{ runner.temp }}" -O "https://github.com/sass/dart-sass/releases/download/${DART_SASS_VERSION}/dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz"
            tar -C "${HOME}/.local" -xf "${{ runner.temp }}/dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz"
            echo "${HOME}/.local/dart-sass" >> "${GITHUB_PATH}"

        - name: Install Hugo
          run: |
            echo "Installing Hugo ${HUGO_VERSION}..."
            curl -sfL --output-dir "${{ runner.temp }}" -O "https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_linux-amd64.tar.gz"
            mkdir "${HOME}/.local/hugo"
            tar -C "${HOME}/.local/hugo" -xf "${{ runner.temp }}/hugo_${HUGO_VERSION}_linux-amd64.tar.gz"
            echo "${HOME}/.local/hugo" >> "${GITHUB_PATH}"

        - name: Log tool versions
          run: |
            echo "Logging tool versions..."
            command -v sass &> /dev/null && echo "Dart Sass: $(sass --version)" || echo "Dart Sass: not installed"
            command -v go &> /dev/null && echo "Go: $(go version)" || echo "Go: not installed"
            command -v hugo &> /dev/null && echo "Hugo: $(hugo version)" || echo "Hugo: not installed"
            command -v node &> /dev/null && echo "Node.js: $(node --version)" || echo "Node.js: not installed"

        - name: Configure Git
          run: |
            echo "Configuring Git..."
            git config --global core.quotepath false

        - name: Fetch full Git history
          run: |
            if [[ $(git rev-parse --is-shallow-repository) == true ]]; then
              echo "Fetching full Git history..."
              git fetch --unshallow
            fi

        - name: Initialize Git submodules
          run: |
            if [[ -f .gitmodules ]]; then
              echo "Initializing Git submodules..."
              git submodule update --init --recursive
            fi

        - name: Install Node.js dependencies
          run: |
            if [[ -f package-lock.json ]]; then
              echo "Installing Node.js dependencies..."
              npm ci
            fi

        - name: Cache restore
          id: cache-restore
          uses: actions/cache/restore@v6
          with:
            path: ${{ runner.temp }}/.cache/hugo
            key: hugo-${{ github.run_id }}
            restore-keys: hugo-

        - name: Build
          run: |
            echo "Building the project..."
            hugo build \
              --gc \
              --minify \
              --cacheDir "${{ runner.temp }}/.cache/hugo"

        - name: Cache save
          uses: actions/cache/save@v6
          with:
            path: ${{ runner.temp }}/.cache/hugo
            key: ${{ steps.cache-restore.outputs.cache-primary-key }}

        - name: Upload build artifact
          uses: actions/upload-artifact@v7
          with:
            name: build-artifact
            path: public
            retention-days: 1
    deploy:
      needs: build
      runs-on: ubuntu-latest
      steps:
        - name: Download build artifact
          uses: actions/download-artifact@v8
          with:
            name: build-artifact
            path: public

        - name: Create Azure Static Web Apps config
          run: |
            cat << 'EOF' > staticwebapp.config.json
            {
              "responseOverrides": {
                "404": {
                  "rewrite": "/404.html",
                  "statusCode": 404
                }
              }
            }
            EOF

        - name: Setup Node.js
          uses: actions/setup-node@v7
          with:
            node-version: ${{ env.NODE_VERSION }}

        - name: Install SWA CLI
          run: npm install -g @azure/static-web-apps-cli --no-fund --no-audit --quiet

        - name: Deploy
          env:
            SWA_CLI_DEPLOYMENT_TOKEN: ${{ secrets.AZURE_STATIC_WEB_APPS_API_TOKEN }}
          run: swa deploy ./public --env production --api-location "" --swa-config-location ./
  ```

Step 4
: In the project configuration file in the root of your local Git repository, set the location of the image cache to the [`cacheDir`][] as shown below.

  {{< code-toggle file=hugo copy=true >}}
  [caches.images]
  dir = ':cacheDir/images'
  {{< /code-toggle >}}

  See [configure file caches][] for more information.

Step 5
: Commit the changes to your local Git repository and push to your GitHub repository.

Step 6
: From GitHub's main menu, choose **Actions**. You will see something like this:

  ![screen capture](azure-10.png)

Step 7
: When GitHub has finished building and deploying your site, the color of the status indicator will change to green.

  ![screen capture](azure-11.png)

In the future, whenever you push a change from your local Git repository, GitHub will rebuild and deploy your site.

## Related resources

For more information on hosting and managing your site with Azure Static Web Apps, consult the official documentation:

- [General documentation][]
- [Custom domain setup][]

[Custom domain setup]: https://learn.microsoft.com/en-us/azure/static-web-apps/custom-domain-external
[General documentation]: https://learn.microsoft.com/en-us/azure/static-web-apps/overview
[Static Web Apps]: https://portal.azure.com/#browse/Microsoft.Web%2FStaticSites
[`baseURL`]: /configuration/all/#baseurl
[`cacheDir`]: /configuration/all/#cachedir
[configure file caches]: /configuration/caches/
[remote]: https://git-scm.com/docs/git-remote
