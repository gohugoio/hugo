---
title: Host on GitLab Pages
description: Host your site on GitLab Pages.
categories: []
keywords: []
aliases: [/hosting-and-deployment/hosting-on-gitlab/]
---

Use these instructions to enable continuous deployment from a GitLab repository.

{{% include "/_common/gitignore-public.md" %}}

## Prerequisites

Please complete the following tasks before continuing:

1. [Create](https://gitlab.com/users/sign_up) a GitLab account
1. [Log in](https://gitlab.com/users/sign_in) to your GitLab account
1. [Create](https://gitlab.com/projects/new) a GitLab repository for your project
1. [Create](https://git-scm.com/docs/git-init) a local Git repository for your project with a [remote][] reference to your GitLab repository
1. Create a Hugo project within your local Git repository and test it with the `hugo server` command
1. Commit the changes to your local Git repository and push to your GitLab repository

## BaseURL

The [`baseURL`][] in your project configuration must reflect the full URL of your GitLab Pages repository if you are using the default GitLab Pages URL (e.g., `https://<YourUsername>.gitlab.io/<your-hugo-site>/`) and not a custom domain.

## Procedure

Step 1
: Create a `.gitlab-ci.yml` file in the root of your project.

  ```yaml {file=".gitlab-ci.yml" copy=true}
  variables:
    # Application versions
    DART_SASS_VERSION: 1.101.0
    HUGO_VERSION: 0.163.2
    NODE_VERSION: 24.16.0
    # Git
    GIT_DEPTH: 0
    GIT_STRATEGY: clone
    GIT_SUBMODULE_STRATEGY: recursive
    # Time zone
    TZ: Europe/Oslo

  image:
    name: golang:1.26.4-bookworm

  pages:
    stage: deploy
    script:
      - |
        # Create directory for user-specific executable files
        echo "Creating directory for user-specific executable files..."
        mkdir -p "${HOME}/.local"

        # Install utilities
        echo "Installing utilities..."
        apt-get update
        apt-get install -y brotli xz-utils zstd

        # Install Dart Sass
        echo "Installing Dart Sass ${DART_SASS_VERSION}..."
        curl -sLJO "https://github.com/sass/dart-sass/releases/download/${DART_SASS_VERSION}/dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz"
        tar -C "${HOME}/.local" -xf "dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz"
        rm "dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz"
        export PATH="${HOME}/.local/dart-sass:${PATH}"

        # Install Hugo
        echo "Installing Hugo ${HUGO_VERSION}..."
        curl -sLJO "https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_linux-amd64.tar.gz"
        mkdir -p "${HOME}/.local/hugo"
        tar -C "${HOME}/.local/hugo" -xf "hugo_${HUGO_VERSION}_linux-amd64.tar.gz"
        rm "hugo_${HUGO_VERSION}_linux-amd64.tar.gz"
        export PATH="${HOME}/.local/hugo:${PATH}"

        # Install Node.js
        echo "Installing Node.js ${NODE_VERSION}..."
        curl -sLJO "https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-x64.tar.xz"
        tar -C "${HOME}/.local" -xf "node-v${NODE_VERSION}-linux-x64.tar.xz"
        rm "node-v${NODE_VERSION}-linux-x64.tar.xz"
        export PATH="${HOME}/.local/node-v${NODE_VERSION}-linux-x64/bin:${PATH}"

        # Verify installations
        echo "Verifying installations..."
        echo "Dart Sass: $(sass --version)"
        echo "Go: $(go version)"
        echo "Hugo: $(hugo version)"
        echo "Node.js: $(node --version)"
        echo "brotli: $(brotli --version)"
        echo "xz: $(xz --version)"
        echo "zstd: $(zstd --version)"

        # Install Node.js dependencies
        echo "Installing Node.js dependencies..."
        [[ -f package-lock.json || -f npm-shrinkwrap.json ]] && npm ci --prefer-offline || true

        # Configure Git
        echo "Configuring Git..."
        git config --global core.quotepath false

        # Build site
        echo "Building site..."
        hugo --gc --minify --baseURL "${CI_PAGES_URL}"

        # Compress published files
        echo "Compressing published files..."
        find public/ -type f -regextype posix-extended -regex '.+\.(css|html|js|json|mjs|svg|txt|xml)$' -print0 > files.txt
        time xargs --null --max-procs=0 --max-args=1 brotli --quality=10 --force --keep < files.txt
        time xargs --null --max-procs=0 --max-args=1 gzip -9 --force --keep < files.txt
    artifacts:
      paths:
        - public
    rules:
      - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
  ```

Step 2
: Commit the changes to your local Git repository and push to your GitLab repository.

Step 3
: From your GitLab repository, navigate to **Build**&nbsp;>&nbsp;**Pipelines** to follow the CI pipeline building your page.

Step 4
: When the pipeline has passed, your new website is available at `https://<YourUsername>.gitlab.io/<your-hugo-site>/`.

In the future, whenever you push a change from your local Git repository, GitLab Pages will rebuild and deploy your site.

## Other resources

- [GitLab Pages documentation][]

[GitLab Pages documentation]: https://docs.gitlab.com/user/project/pages/
[`baseURL`]: /configuration/all/#baseurl
[remote]: https://git-scm.com/docs/git-remote
