---
title: Host on GitLab Pages
description: Host your site on GitLab Pages.
categories: []
keywords: []
aliases: [/hosting-and-deployment/hosting-on-gitlab/]
---

## Assumptions

- Working familiarity with Git for version control
- Completion of the Hugo [Quick Start]
- A [GitLab account](https://gitlab.com/users/sign_in)
- A Hugo website on your local machine that you are ready to publish

## BaseURL

The `baseURL` in your [site configuration](/configuration/) must reflect the full URL of your GitLab pages repository if you are using the default GitLab Pages URL (e.g., `https://<YourUsername>.gitlab.io/<your-hugo-site>/`) and not a custom domain.

## Configure GitLab CI/CD

Define your [CI/CD](g) jobs by creating a `.gitlab-ci.yml` file in the root of your project.

```yaml {file=".gitlab-ci.yml" copy=true}
variables:
  # Application versions
  DART_SASS_VERSION: 1.90.0
  HUGO_VERSION: 0.148.2
  NODE_VERSION: 22.18.0
  # Git
  GIT_DEPTH: 0
  GIT_STRATEGY: clone
  GIT_SUBMODULE_STRATEGY: recursive
  # Time zone
  TZ: Europe/Oslo

image:
  name: golang:1.24.5-bookworm

pages:
  stage: deploy
  script:
    # Create directory for user-specific executable files
    - echo "Creating directory for user-specific executable files..."
    - mkdir -p "${HOME}/.local"

    # Install utilities
    - echo "Installing utilities..."
    - apt-get update
    - apt-get install -y brotli xz-utils zstd

    # Install Dart Sass
    - echo "Installing Dart Sass ${DART_SASS_VERSION}..."
    - curl -sLJO "https://github.com/sass/dart-sass/releases/download/${DART_SASS_VERSION}/dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz"
    - tar -C "${HOME}/.local" -xf "dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz"
    - rm "dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz"
    - export PATH="${HOME}/.local/dart-sass:${PATH}"

    # Install Hugo
    - echo "Installing Hugo ${HUGO_VERSION}..."
    - curl -sLJO "https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_linux-amd64.tar.gz"
    - mkdir "${HOME}/.local/hugo"
    - tar -C "${HOME}/.local/hugo" -xf "hugo_extended_${HUGO_VERSION}_linux-amd64.tar.gz"
    - rm "hugo_extended_${HUGO_VERSION}_linux-amd64.tar.gz"
    - export PATH="${HOME}/.local/hugo:${PATH}"

    # Install Node.js
    - echo "Installing Node.js ${NODE_VERSION}..."
    - curl -sLJO "https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-x64.tar.xz"
    - tar -C "${HOME}/.local" -xf "node-v${NODE_VERSION}-linux-x64.tar.xz"
    - rm "node-v${NODE_VERSION}-linux-x64.tar.xz"
    - export PATH="${HOME}/.local/node-v${NODE_VERSION}-linux-x64/bin:${PATH}"

    # Verify installations
    - echo "Verifying installations..."
    - "echo Dart Sass: $(sass --version)"
    - "echo Go: $(go version)"
    - "echo Hugo: $(hugo version)"
    - "echo Node.js: $(node --version)"
    - "echo brotli: $(brotli --version)"
    - "echo xz: $(xz --version)"
    - "echo zstd: $(zstd --version)"

    # Install Node.js dependencies
    - echo "Installing Node.js dependencies..."
    - "[[ -f package-lock.json || -f npm-shrinkwrap.json ]] && npm ci --prefer-offline || true"

    # Configure Git
    - echo "Configuring Git..."
    - git config core.quotepath false

    # Build site
    - echo "Building site..."
    - hugo --gc --minify --baseURL "${CI_PAGES_URL}"

    # Compress published files
    - echo "Compressing published files..."
    - find public/ -type f -regextype posix-extended -regex '.+\.(css|html|js|json|mjs|svg|txt|xml)$' -print0 > files.txt
    - time xargs --null --max-procs=0 --max-args=1 brotli --quality=10 --force --keep < files.txt
    - time xargs --null --max-procs=0 --max-args=1 gzip -9 --force --keep < files.txt
  artifacts:
    paths:
      - public
  rules:
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
```

## Push your Hugo website to GitLab

Next, create a new repository on GitLab. It is not necessary to make the repository public. In addition, you might want to add `/public` to your .gitignore file, as there is no need to push compiled assets to GitLab or keep your output website in version control.

```sh
# initialize new git repository
git init

# add /public directory to our .gitignore file
echo "/public" >> .gitignore

# commit and push code to master branch
git add .
git commit -m "Initial commit"
git remote add origin https://gitlab.com/YourUsername/your-hugo-site.git
git push -u origin master
```

## Wait for your page to build

That's it! You can now follow the CI agent building your page at `https://gitlab.com/<YourUsername>/<your-hugo-site>/pipelines`.

After the build has passed, your new website is available at `https://<YourUsername>.gitlab.io/<your-hugo-site>/`.

## Next steps

GitLab supports using custom CNAME's and TLS certificates. For more details on GitLab Pages, see the [GitLab Pages setup documentation](https://about.gitlab.com/2016/04/07/gitlab-pages-setup/).

[Quick Start]: /getting-started/quick-start/
