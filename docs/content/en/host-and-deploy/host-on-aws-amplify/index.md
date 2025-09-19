---
title: Host on AWS Amplify
description: Host your site on AWS Amplify.
categories: []
keywords: []
aliases: [/hosting-and-deployment/hosting-on-aws-amplify/]
---

Use these instructions to enable continuous deployment from a GitHub repository. The same general steps apply if you are using GitLab for version control.

## Prerequisites

Please complete the following tasks before continuing:

1. [Create](https://aws.amazon.com/resources/create-account/) an AWS account
1. [Log in](https://console.aws.amazon.com/) to your AWS account
1. [Create](https://github.com/signup) a GitHub account
1. [Log in](https://github.com/login) to your GitHub account
1. [Create](https://github.com/new) a GitHub repository for your project
1. [Create](https://git-scm.com/docs/git-init) a local Git repository for your project with a [remote](https://git-scm.com/docs/git-remote) reference to your GitHub repository
1. Create a Hugo site within your local Git repository and test it with the `hugo server` command
1. Commit the changes to your local Git repository and push to your GitHub repository.

## Procedure

This procedure will enable continuous deployment from a GitHub repository. The procedure is essentially the same if you are using GitLab or Bitbucket.

Step 1
: Create a file named `amplify.yml` in the root of your project.

  ```sh
  touch amplify.yml
  ```

Step 2
: Copy and paste the YAML below into the file you created. Change the application versions and time zone as needed.

  ```yaml {file="amplify.yml" copy=true}
  version: 1
  env:
    variables:
      # Application versions
      DART_SASS_VERSION: 1.90.0
      GO_VERSION: 1.24.5
      HUGO_VERSION: 0.148.2
      # Time zone
      TZ: Europe/Oslo
      # Cache
      HUGO_CACHEDIR: ${PWD}/.hugo
      NPM_CONFIG_CACHE: ${PWD}/.npm
  frontend:
    phases:
      preBuild:
        commands:
          # Create directory for user-specific executable files
          - echo "Creating directory for user-specific executable files..."
          - mkdir -p "${HOME}/.local"

          # Install Dart Sass
          - echo "Installing Dart Sass ${DART_SASS_VERSION}..."
          - curl -sLJO "https://github.com/sass/dart-sass/releases/download/${DART_SASS_VERSION}/dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz"
          - tar -C "${HOME}/.local" -xf "dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz"
          - rm "dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz"
          - export PATH="${HOME}/.local/dart-sass:${PATH}"

          # Install Go
          - echo "Installing Go ${GO_VERSION}..."
          - curl -sLJO "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
          - tar -C "${HOME}/.local" -xf "go${GO_VERSION}.linux-amd64.tar.gz"
          - rm "go${GO_VERSION}.linux-amd64.tar.gz"
          - export PATH="${HOME}/.local/go/bin:${PATH}"

          # Install Hugo
          - echo "Installing Hugo ${HUGO_VERSION}..."
          - curl -sLJO "https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_linux-amd64.tar.gz"
          - mkdir "${HOME}/.local/hugo"
          - tar -C "${HOME}/.local/hugo" -xf "hugo_extended_${HUGO_VERSION}_linux-amd64.tar.gz"
          - rm "hugo_extended_${HUGO_VERSION}_linux-amd64.tar.gz"
          - export PATH="${HOME}/.local/hugo:${PATH}"

          # Verify installations
          - echo "Verifying installations..."
          - "echo Dart Sass: $(sass --version)"
          - "echo Go: $(go version)"
          - "echo Hugo: $(hugo version)"
          - "echo Node.js: $(node --version)"

          # Install Node.js dependencies
          - echo "Installing Node.js dependencies..."
          - "[[ -f package-lock.json || -f npm-shrinkwrap.json ]] && npm ci --prefer-offline || true"

          # Configure Git
          - echo "Configuring Git..."
          - git config core.quotepath false
      build:
        commands:
          - echo "Building site..."
          - hugo --gc --minify
    artifacts:
      baseDirectory: public
      files:
        - '**/*'
    cache:
      paths:
        - ${HUGO_CACHEDIR}/**/*
        - ${NPM_CONFIG_CACHE}/**/*
  ```

Step 3
: Commit and push the change to your GitHub repository.

  ```sh
  git add -A
  git commit -m "Create amplify.yml"
  git push
  ```

Step 4
: Log in to your AWS account, navigate to the [Amplify Console], then press the  **Deploy an app** button.

Step 5
: Choose a source code provider, then press the **Next** button.

  ![screen capture](amplify-step-05.png)

Step 6
: Authorize AWS Amplify to access your GitHub account.

  ![screen capture](amplify-step-06.png)

Step 7
: Select your personal account or relevant organization.

  ![screen capture](amplify-step-07.png)

Step 8
: Authorize access to one or more repositories.

  ![screen capture](amplify-step-08.png)

Step 9
: Select a repository and branch, then press the **Next** button.

  ![screen capture](amplify-step-09.png)

Step 10
: On the "App settings" page, scroll to the bottom then press the **Next** button. Amplify reads the `amplify.yml` file you created in Steps 1-3 instead of using the values on this page.

Step 11
: On the "Review" page, scroll to the bottom then press the **Save and deploy** button.

Step 12
: When your site has finished deploying, press the **Visit deployed URL** button to view your published site.

  ![screen capture](amplify-step-11.png)

[Amplify Console]: https://console.aws.amazon.com/amplify/apps
