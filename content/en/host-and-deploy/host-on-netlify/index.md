---
title: Host on Netlify
description: Host your site on Netlify.
categories: []
keywords: []
aliases: [/hosting-and-deployment/hosting-on-netlify/]
---

Use these instructions to enable continuous deployment from a GitHub repository. The same general steps apply if you are using Azure DevOps, Bitbucket, or GitLab for version control.

## Prerequisites

Please complete the following tasks before continuing:

1. [Create](https://app.netlify.com/signup) a Netlify account
1. [Log in](https://app.netlify.com/login) to your Netlify account
1. [Create](https://github.com/signup) a GitHub account
1. [Log in](https://github.com/login) to your GitHub account
1. [Create](https://github.com/new) a GitHub repository for your project
1. [Create](https://git-scm.com/docs/git-init) a local Git repository for your project with a [remote](https://git-scm.com/docs/git-remote) reference to your GitHub repository
1. Create a Hugo site within your local Git repository and test it with the `hugo server` command
1. Commit the changes to your local Git repository and push to your GitHub repository.

## Procedure

Step 1
: Log in to your Netlify account, navigate to the Sites page, press the **Add new site** button, and choose "Import an existing project" from the dropdown menu.

Step 2
: Select your deployment method.

  ![screen capture](netlify-step-02.png)

Step 3
: Authorize Netlify to connect with your GitHub account by pressing the **Authorize Netlify** button.

  ![screen capture](netlify-step-03.png)

Step 4
: Press the **Configure Netlify on GitHub** button.

  ![screen capture](netlify-step-04.png)

Step 5
: Install the Netlify app by selecting your GitHub account.

  ![screen capture](netlify-step-05.png)

Step 6
: Press the **Install** button.

  ![screen capture](netlify-step-06.png)

Step 7
: Click on the site's repository from the list.

  ![screen capture](netlify-step-07.png)

Step 8
: Set the site name and branch from which to deploy.

  ![screen capture](netlify-step-08.png)

Step 9
: Define the build settings, press the **Add environment variables** button, then press the **New variable** button.

  ![screen capture](netlify-step-09.png)

Step 10
: Create a new environment variable named `HUGO_VERSION` and set the value to the [latest version](https://github.com/gohugoio/hugo/releases/latest).

  ![screen capture](netlify-step-10.png)

Step 11
: Press the "Deploy my new site" button at the bottom of the page.

  ![screen capture](netlify-step-11.png)

Step 12
: At the bottom of the screen, wait for the deploy to complete, then click on the deploy log entry.

  ![screen capture](netlify-step-12.png)

Step 13
: Press the **Open production deploy** button to view the live site.

  ![screen capture](netlify-step-13.png)

## Configuration file

In the procedure above we configured our site using the Netlify user interface. Most site owners find it easier to use a configuration file checked into source control.

Create a new file named `netlify.toml` in the root of your project directory. In its simplest form, the configuration file might look like this:

```toml {file="netlify.toml"}
[build.environment]
GO_VERSION = "1.24.5"
HUGO_VERSION = "0.148.2"
NODE_VERSION = "22.18.0"
TZ = "Europe/Oslo"

[build]
publish = "public"
command = """\
  git config core.quotepath false && \
  hugo --gc --minify --baseURL "${URL}"
  """
```

If your site requires Dart Sass to transpile Sass to CSS, the configuration file should look something like this:

```toml {file="netlify.toml"}
[build.environment]
DART_SASS_VERSION = "1.90.0"
GO_VERSION = "1.24.5"
HUGO_VERSION = "0.148.2"
NODE_VERSION = "22.18.0"
TZ = "Europe/Oslo"

[build]
publish = "public"
command = """\
  curl -sLJO "https://github.com/sass/dart-sass/releases/download/${DART_SASS_VERSION}/dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz" && \
  tar -C "${HOME}/.local" -xf "dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz" && \
  rm "dart-sass-${DART_SASS_VERSION}-linux-x64.tar.gz" && \
  export PATH="${HOME}/.local/dart-sass:${PATH}" && \
  git config core.quotepath false && \
  hugo --gc --minify --baseURL "${URL}"
  """
```
