---
title: Host on Netlify
description: Host your project on Netlify.
categories: []
keywords: []
aliases: [/hosting-and-deployment/hosting-on-netlify/]
---

Use these instructions to enable continuous deployment from a GitHub repository. The same general steps apply for other Git providers such as GitLab or Bitbucket.

{{% include "/_common/gitignore-public.md" %}}

## Prerequisites

Please complete the following tasks before continuing:

1. [Create](https://app.netlify.com/signup) a Netlify account.
1. [Log in](https://app.netlify.com/login) to your Netlify account.
1. [Create](https://github.com/signup) a GitHub account.
1. [Log in](https://github.com/login) to your GitHub account.
1. [Create](https://github.com/new) a GitHub repository for your project.
1. [Create](https://git-scm.com/docs/git-init) a local Git repository for your project with a [remote][] reference to your GitHub repository.
1. Create a Hugo project within your local Git repository and test it with the `hugo server` command.
1. Commit the changes to your local Git repository and push to your GitHub repository.

## Procedure

Step 1
: Create a `netlify.toml` file in the root of your project, adjusting the tool versions and time zone as needed.

  ```toml {file="netlify.toml" copy=true}
  [build.environment]
  GO_VERSION = "1.26.4"
  HUGO_VERSION = "0.163.3"
  NODE_VERSION = "24.16.0"
  TZ = "Europe/Oslo"

  [build]
  publish = "public"
  command = """\
    git config --global core.quotepath false && \
    hugo build --gc --minify --baseURL "${URL}"
    """
  ```

  If your project requires Dart Sass to transpile Sass to CSS, set the `DART_SASS_VERSION` and include the Dart Sass installation in the build step.

  ```toml {file="netlify.toml" copy=true}
  [build.environment]
  DART_SASS_VERSION = "1.101.0"
  GO_VERSION = "1.26.4"
  HUGO_VERSION = "0.163.3"
  NODE_VERSION = "24.16.0"
  TZ = "Europe/Oslo"

  [build]
  publish = "public"
  command = """\
    curl -sfLO "https://github.com/sass/dart-sass/releases/download/${DART_SASS_VERSION}/dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz" && \
    tar -C "${HOME}/.local" -xf "dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz" && \
    rm "dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz" && \
    export PATH="${HOME}/.local/dart-sass:${PATH}" && \
    git config --global core.quotepath false && \
    hugo build --gc --minify --baseURL "${URL}"
    """
  ```

Step 2
: In your project configuration, change the location of the image cache to the [`cacheDir`][] as shown below:

  {{< code-toggle file=hugo copy=true >}}
  [caches.images]
  dir = ':cacheDir/images'
  {{< /code-toggle >}}

  See [configure file caches][] for more information.

Step 3
: Commit the changes to your local Git repository and push to your GitHub repository.

Step 4
: In the upper right corner of the Netlify dashboard, press the **Add new project** button and select “Import an existing project".

  ![screen capture](netlify-01.png)

Step 5
: Connect to GitHub.

  ![screen capture](netlify-02.png)

Step 6
: Press the "Authorize Netlify" button to allow the Netlify application to access your GitHub account.

  ![screen capture](netlify-03.png)

Step 7
: Press the **Configure Netlify on GitHub** button.

  ![screen capture](netlify-04.png)

Step 8
: Select the GitHub account where you want to install the Netlify application.

  ![screen capture](netlify-05.png)

Step 9
: Authorize the Netlify application to access all repositories or only select repositories, then press the Install button.

  ![screen capture](netlify-06.png)

Your browser will be redirected to the Netlify dashboard.

Step 10
: Click on the name of the repository you wish to import.

  ![screen capture](netlify-07.png)

Step 11
: On the "Review configuration" page, enter a project name, leave the settings at their default values, then press the **Deploy** button.

  ![screen capture](netlify-08.png)

  ![screen capture](netlify-09.png)

Step 12
: When the deployment completes, click on the link to your published site.

  ![screen capture](netlify-10.png)

In the future, whenever you push a change from your local Git repository, Netlify will rebuild and deploy your site.

[`cacheDir`]: /configuration/all/#cachedir
[configure file caches]: /configuration/caches/
[remote]: https://git-scm.com/docs/git-remote
