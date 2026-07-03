---
title: Host on GitLab Pages
description: Host your project on GitLab Pages.
categories: []
keywords: []
aliases: [/hosting-and-deployment/hosting-on-gitlab/]
---

Use these instructions to enable continuous deployment from a GitLab repository to GitLab Pages.

{{% include "/_common/gitignore-public.md" %}}

## Prerequisites

Please complete the following tasks before continuing:

1. [Create](https://gitlab.com/users/sign_up) a GitLab account.
1. [Log in](https://gitlab.com/users/sign_in) to your GitLab account.
1. [Create](https://gitlab.com/projects/new) a GitLab repository for your project.
1. [Create](https://git-scm.com/docs/git-init) a local Git repository for your project with a [remote][] reference to your GitLab repository.
1. Create a Hugo project within your local Git repository and test it with the `hugo server` command.
1. Commit the changes to your local Git repository and push to your GitLab repository.

## BaseURL

The [`baseURL`][] in your project configuration must reflect the full URL of your GitLab Pages repository if you are using the default GitLab Pages URL (e.g., `https://<YourUsername>.gitlab.io/<your-hugo-site>/`) and not a custom domain.

## Procedure

Step 1
: Create a `.gitlab-ci.yml` file in the root of your project, adjusting the tool versions and time zone as needed.

  ```yaml {file=".gitlab-ci.yml" copy=true}
  variables:
    # Define tool versions
    DART_SASS_VERSION: 1.101.0
    GO_VERSION: 1.26.4
    HUGO_VERSION: 0.163.3
    NODE_VERSION: 24.16.0

    # Set the build timezone
    TZ: Europe/Oslo

    # Set the build cache directory
    HUGO_CACHEDIR: ${CI_PROJECT_DIR}/.cache/hugo

    # Set the repository clone and fetch strategy
    GIT_DEPTH: 0
    GIT_STRATEGY: clone
    GIT_SUBMODULE_STRATEGY: recursive
  cache:
    key: ${CI_COMMIT_REF_SLUG}
    fallback_keys:
      - ${CI_DEFAULT_BRANCH}
    paths:
      - .cache/hugo
  image:
    name: buildpack-deps:bookworm
  pages:
    stage: deploy
    script:
      - chmod a+x build.sh && ./build.sh
    artifacts:
      paths:
        - public
    rules:
      - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
  ```

Step 2
: Create a `build.sh` file in the root of your project.

  ```sh {file="build.sh" copy=true}
  #!/usr/bin/env bash

  #------------------------------------------------------------------------------
  # @file
  # Builds a Hugo project hosted on GitLab Pages.
  #------------------------------------------------------------------------------

  # Exit on error, undefined variables, or pipe failures
  set -euo pipefail

  # Perform cleanup
  cleanup() {
    if [[ -n "${build_temp_dir:-}" && -d "${build_temp_dir}" ]]; then
      rm -rf "${build_temp_dir}"
    fi
  }

  # Register the cleanup trap
  trap cleanup EXIT SIGINT SIGTERM

  main() {
    # Create a temporary directory for downloads
    build_temp_dir=$(mktemp -d)

    # Create a local tools directory
    mkdir -p "${HOME}/.local"

    # Install utilities
    echo "Installing utilities..."
    apt-get update > /dev/null
    apt-get install -y brotli > /dev/null

    # Install Dart Sass
    echo "Installing Dart Sass ${DART_SASS_VERSION}..."
    curl -sfLO --output-dir "${build_temp_dir}" "https://github.com/sass/dart-sass/releases/download/${DART_SASS_VERSION}/dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz"
    tar -C "${HOME}/.local" -xf "${build_temp_dir}/dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz"
    export PATH="${HOME}/.local/dart-sass:${PATH}"

    # Install Go
    if [[ -f "${CI_PROJECT_DIR}/go.mod" ]]; then
      echo "Installing Go ${GO_VERSION}..."
      curl -sfLO --output-dir "${build_temp_dir}" "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
      tar -C "${HOME}/.local" -xf "${build_temp_dir}/go${GO_VERSION}.linux-amd64.tar.gz"
      export PATH="${HOME}/.local/go/bin:${PATH}"
    fi

    # Install Hugo
    echo "Installing Hugo ${HUGO_VERSION}..."
    curl -sfLO --output-dir "${build_temp_dir}" "https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_linux-amd64.tar.gz"
    mkdir -p "${HOME}/.local/hugo"
    tar -C "${HOME}/.local/hugo" -xf "${build_temp_dir}/hugo_${HUGO_VERSION}_linux-amd64.tar.gz"
    export PATH="${HOME}/.local/hugo:${PATH}"

    # Install Node.js
    if [[ -f "${CI_PROJECT_DIR}/package-lock.json" ]]; then
      echo "Installing Node.js ${NODE_VERSION}..."
      curl -sfLO --output-dir "${build_temp_dir}" "https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-x64.tar.gz"
      tar -C "${HOME}/.local" -xf "${build_temp_dir}/node-v${NODE_VERSION}-linux-x64.tar.gz"
      export PATH="${HOME}/.local/node-v${NODE_VERSION}-linux-x64/bin:${PATH}"
    fi

    # Log tool versions
    echo "Logging tool versions..."
    command -v sass &> /dev/null && echo "Dart Sass: $(sass --version)" || echo "Dart Sass: not installed"
    command -v go &> /dev/null && echo "Go: $(go version)" || echo "Go: not installed"
    command -v hugo &> /dev/null && echo "Hugo: $(hugo version)" || echo "Hugo: not installed"
    command -v node &> /dev/null && echo "Node.js: $(node --version)" || echo "Node.js: not installed"

    # Configure Git
    echo "Configuring Git..."
    git config --global core.quotepath false

    # Fetch full Git history
    if [[ $(git rev-parse --is-shallow-repository) == true ]]; then
      echo "Fetching full Git history..."
      git fetch --unshallow
    fi

    # Initialize Git submodules
    if [[ -f .gitmodules ]]; then
      echo "Initializing Git submodules..."
      git submodule update --init --recursive
    fi

    # Install Node.js dependencies
    if [[ -f package-lock.json ]]; then
      echo "Installing Node.js dependencies..."
      npm ci
    fi

    # Build the project
    echo "Building the project..."
    hugo build --gc --minify

    # Compress published files
    echo "Compressing published files..."
    find public/ -type f -regextype posix-extended -regex '.+\.(cjs|css|html|js|json|mjs|svg|txt|xml)$' -print0 > "${build_temp_dir}/files.txt"
    xargs --null --max-procs=0 --max-args=1 brotli --quality=10 --force --keep < "${build_temp_dir}/files.txt"
    xargs --null --max-procs=0 --max-args=1 gzip -9 --force --keep < "${build_temp_dir}/files.txt"
  }

  main "$@"
  ```

Step 3
: In your project configuration, change the location of the image cache to the [`cacheDir`][] as shown below:

  {{< code-toggle file=hugo copy=true >}}
  [caches.images]
  dir = ':cacheDir/images'
  {{< /code-toggle >}}

  See [configure file caches][] for more information.

Step 4
: Commit the changes to your local Git repository and push to your GitLab repository.

Step 5
: From your GitLab repository, navigate to **Build**&nbsp;>&nbsp;**Pipelines** to follow the CI pipeline building your page.

Step 6
: When the pipeline has passed, your new website is available at `https://<YourUsername>.gitlab.io/<your-hugo-site>/`.

In the future, whenever you push a change from your local Git repository, GitLab Pages will rebuild and deploy your site.

[`baseURL`]: /configuration/all/#baseurl
[`cacheDir`]: /configuration/all/#cachedir
[configure file caches]: /configuration/caches/
[remote]: https://git-scm.com/docs/git-remote
