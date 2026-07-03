---
title: Host on AWS Amplify
description: Host your project on AWS Amplify.
categories: []
keywords: []
aliases: [/hosting-and-deployment/hosting-on-aws-amplify/]
---

Use these instructions to enable continuous deployment from a GitHub repository. The same general steps apply for other Git providers such as GitLab or Bitbucket.

{{% include "/_common/gitignore-public.md" %}}

## Prerequisites

Please complete the following tasks before continuing:

1. [Create](https://aws.amazon.com/resources/create-account/) an AWS account.
1. [Log in](https://console.aws.amazon.com/) to your AWS account.
1. [Create](https://github.com/signup) a GitHub account.
1. [Log in](https://github.com/login) to your GitHub account.
1. [Create](https://github.com/new) a GitHub repository for your project.
1. [Create](https://git-scm.com/docs/git-init) a local Git repository for your project with a [remote][] reference to your GitHub repository.
1. Create a Hugo project within your local Git repository and test it with the `hugo server` command.
1. Commit the changes to your local Git repository and push to your GitHub repository.

## Procedure

Step 1
: Create an `amplify.yml` file in the root of your project, adjusting the tool versions and time zone as needed.

  ```yaml {file="amplify.yml" copy=true}
  version: 1
  env:
    variables:
      # Define tool versions
      DART_SASS_VERSION: 1.101.0
      GO_VERSION: 1.26.4
      HUGO_VERSION: 0.163.3
      NODE_VERSION: 24.16.0

      # Set the build time zone
      TZ: Europe/Oslo

      # Set the build cache directory
      HUGO_CACHEDIR: ${PWD}/.cache/hugo
  frontend:
    phases:
      preBuild:
        commands:
          # Create a temporary directory for downloads
          - build_temp_dir=$(mktemp -d)

          # Create a local tools directory
          - mkdir -p "${HOME}/.local"

          # Install Dart Sass
          - |
            echo "Installing Dart Sass ${DART_SASS_VERSION}..."
            curl -sfL --output-dir "${build_temp_dir}" -O "https://github.com/sass/dart-sass/releases/download/${DART_SASS_VERSION}/dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz"
            tar -C "${HOME}/.local" -xf "${build_temp_dir}/dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz"
            export PATH="${HOME}/.local/dart-sass:${PATH}"

          # Install Go
          - |
            if [[ -f "go.mod" ]]; then
              echo "Installing Go ${GO_VERSION}..."
              curl -sfL --output-dir "${build_temp_dir}" -O "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
              tar -C "${HOME}/.local" -xf "${build_temp_dir}/go${GO_VERSION}.linux-amd64.tar.gz"
              export PATH="${HOME}/.local/go/bin:${PATH}"
            fi

          # Install Hugo
          - |
            echo "Installing Hugo ${HUGO_VERSION}..."
            curl -sfL --output-dir "${build_temp_dir}" -O "https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_linux-amd64.tar.gz"
            mkdir -p "${HOME}/.local/hugo"
            tar -C "${HOME}/.local/hugo" -xf "${build_temp_dir}/hugo_${HUGO_VERSION}_linux-amd64.tar.gz"
            export PATH="${HOME}/.local/hugo:${PATH}"

          # Install Node.js
          - |
            if [[ -f "package-lock.json" ]]; then
              echo "Installing Node.js ${NODE_VERSION}..."
              curl -sfL --output-dir "${build_temp_dir}" -O "https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-x64.tar.gz"
              tar -C "${HOME}/.local" -xf "${build_temp_dir}/node-v${NODE_VERSION}-linux-x64.tar.gz"
              export PATH="${HOME}/.local/node-v${NODE_VERSION}-linux-x64/bin:${PATH}"
            fi

          # Log tool versions
          - |
            echo "Logging tool versions..."
            command -v sass &> /dev/null && echo "Dart Sass: $(sass --version)" || echo "Dart Sass: not installed"
            command -v go &> /dev/null && echo "Go: $(go version)" || echo "Go: not installed"
            command -v hugo &> /dev/null && echo "Hugo: $(hugo version)" || echo "Hugo: not installed"
            command -v node &> /dev/null && echo "Node.js: $(node --version)" || echo "Node.js: not installed"

          # Configure Git
          - |
            echo "Configuring Git..."
            git config --global core.quotepath false

          # Fetch full Git history
          - |
            if [[ $(git rev-parse --is-shallow-repository) == true ]]; then
              echo "Fetching full Git history..."
              git fetch --unshallow
            fi

          # Initialize Git submodules
          - |
            if [[ -f .gitmodules ]]; then
              echo "Initializing Git submodules..."
              git submodule update --init --recursive
            fi

          # Install Node.js dependencies
          - |
            if [[ -f package-lock.json ]]; then
              echo "Installing Node.js dependencies..."
              npm ci
            fi
      build:
        commands:
          - echo "Building the project..."
          - hugo build --gc --minify
    artifacts:
      baseDirectory: public
      files:
        - '**/*'
    cache:
      paths:
        - .cache/hugo/**/*
  ```

Step 2
: In your project configuration, change the location of the image cache to the [`cacheDir`][] as shown below:

  {{< code-toggle file=hugo copy=true >}}
  [caches.images]
  dir = ':cacheDir/images'
  {{< /code-toggle >}}

  See [configure file caches][] for more information.

Step 3
: Commit and push the change to your GitHub repository.

  ```sh
  git add -A
  git commit -m "Create amplify.yml"
  git push
  ```

Step 4
: Log in to your AWS account, navigate to the [Amplify Console][], then press the  **Deploy an app** button.

Step 5
: Choose a source code provider, then press the **Next** button.

  ![screen capture](amplify-01.png)

Step 6
: Authorize AWS Amplify to access your GitHub account.

  ![screen capture](amplify-02.png)

Step 7
: Select your personal account or relevant organization.

  ![screen capture](amplify-03.png)

Step 8
: Authorize access to one or more repositories.

  ![screen capture](amplify-04.png)

Step 9
: Select a repository and branch, then press the **Next** button.

  ![screen capture](amplify-05.png)

Step 10
: On the "App settings" page, scroll to the bottom then press the **Next** button. Amplify reads the `amplify.yml` file you created in Steps 1-3 instead of using the values on this page.

Step 11
: On the "Review" page, scroll to the bottom then press the **Save and deploy** button.

Step 12
: When your site has finished deploying, press the **Visit deployed URL** button to view your published site.

  ![screen capture](amplify-06.png)

[Amplify Console]: https://console.aws.amazon.com/amplify/apps
[`cacheDir`]: /configuration/all/#cachedir
[configure file caches]: /configuration/caches/
[remote]: https://git-scm.com/docs/git-remote
