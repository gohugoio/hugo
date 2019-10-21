---
title: Host-Agnostic Deploys with Nanobox
linktitle: Host-Agnostic Deploys with Nanobox
description: Easily deploy Hugo to AWS, DigitalOcean, Google, Azure, and more...
date: 2017-08-24
publishdate: 2017-08-24
lastmod: 2017-08-24
categories: [hosting and deployment]
keywords: [nanobox,deployment,hosting,aws,digitalocean,azure,google,linode]
authors: [Steve Domino]
menu:
  docs:
    parent: "hosting-and-deployment"
    weight: 05
weight: 05
sections_weight: 05
draft: false
aliases: [/tutorials/deployment-with-nanobox/]
toc: true
---

![hugo with nanobox](/images/hosting-and-deployment/deployment-with-nanobox/hugo-with-nanobox.png)

Nanobox provides an entire end-to-end workflow for developing and deploying applications. Using Nanobox to deploy also means you'll use it to develop your application.

{{% note %}}
If you're already using Nanobox and just need deployment instructions, you can skip to [Deploying Hugo with Nanobox](#deploying-hugo-with-nanobox)
{{% /note %}}


## What You'll Need

With Nanobox you don't need to worry about having Go or Hugo installed. They'll be installed as part of the development environment created for you.

To get started you'll just need the following three items:

* [A Nanobox Account](https://nanobox.io) - Signup is free
* [Nanobox Desktop](https://dashboard.nanobox.io/download) - The free desktop development tool
* An account with a hosting provider such as:
  - [AWS](https://docs.nanobox.io/providers/hosting-accounts/aws/)
  - [Google](https://docs.nanobox.io/providers/hosting-accounts/gcp/)
  - [Azure](https://docs.nanobox.io/providers/hosting-accounts/azure/)
  - [DigitalOcean](https://docs.nanobox.io/providers/hosting-accounts/digitalocean/)
  - [Linode](https://docs.nanobox.io/providers/hosting-accounts/linode/)
  - [More...](https://docs.nanobox.io/providers/hosting-accounts/)
  - [Roll Your Own](https://docs.nanobox.io/providers/create/)

### Before You Begin

There are a few things to get out of the way before diving into the guide. To deploy, you'll need to make sure you have connected a host account to your Nanobox account, and launched a new application.

#### Connect a Host Account

Nanobox lets you choose where to host your application (AWS, DigitalOcean, Google, Azure, etc.). In the [Hosting Accounts](https://dashboard.nanobox.io/provider_accounts) section of your Nanobox dashboard [link your Nanobox account with your host](https://docs.nanobox.io/providers/hosting-accounts/).

#### Launch a New Application on Nanobox

[Launching a new app on Nanobox](https://docs.nanobox.io/workflow/launch-app/) is very simple. Navigate to [Launch New App](https://dashboard.nanobox.io/apps/new) in the dashboard, and follow the steps there. You'll be asked to name your app, and to select a host and region.

With those out of the way you're ready to get started!


## Getting Started

{{% note %}}
If you already have a functioning Hugo app, you can skip to [Configure Hugo to run with Nanobox](#configure-hugo-to-run-with-nanobox)
{{% /note %}}

To get started, all you'll need is an empty project directory. Create a directory wherever you want your application to live and `cd` into it:

`mkdir path/to/project && cd path/to/project`

### Configure Hugo to run with Nanobox

Nanobox uses a simple config file known as a [boxfile.yml](https://docs.nanobox.io/boxfile/) to describe your application's infrastructure. In the root of your project add the following `boxfile.yml`:

{{< code file="boxfile.yml" >}}
run.config:

  # use the static engine
  engine: static
  engine.config:

    # tell the engine where to serve static assets from
    rel_dir: public

  # enable file watching for live reload
  fs_watch: true

  # install hugo
  extra_steps:
    - bash ./install.sh

deploy.config:

  # generate site on deploy
  extra_steps:
    - hugo

{{< /code >}}

{{% note %}}
If you already have a functioning Hugo app, after adding the boxfile, you can skip to [Deploying Hugo with Nanobox](#deploying-hugo-with-nanobox).
{{% /note %}}

### Installing Hugo

Nanobox uses Docker to create instant, isolated, development environments. Because of this, you'll need to make sure that during development you have Hugo available.

Do this by adding a custom install script at the root of your project that will install Hugo automatically for you:

{{< code file="install.sh" >}}

#!/bin/bash

if [[ ! -f /data/bin/hugo ]]; then
  cd /tmp
  wget https://github.com/gohugoio/hugo/releases/download/v0.31.1/hugo_0.31.1_Linux-64bit.tar.gz
  tar -xzf hugo_0.31.1_Linux-64bit.tar.gz
  mv hugo /data/bin/hugo
  cd -
  rm -rf /tmp/*
fi

{{< /code >}}

{{% note %}}
If the install script fails during `nanobox run` you may need to make it executable with `chmod +x install.sh`
{{% /note %}}
{{% note %}}
Make sure to check the version of Hugo you have installed and update the install script to match.
{{% /note %}}

### Generating a New Hugo App

You'll generate your new application from inside the Nanobox VM (this is why you don't need to worry about having Go or Hugo installed).

Run the following command to drop into a Nanobox console (inside the VM) where your codebase is mounted:

```
nanobox run
```

![nanobox run](/images/hosting-and-deployment/deployment-with-nanobox/nanobox-run.png)

Once inside the console use the following steps to create a new Hugo application:

```
# cd into the /tmp dir to create an app
cd /tmp

# generate the hugo app
hugo new site app

# cd back into the /app dir
cd -

# copy the generated app into the project
shopt -s dotglob
cp -a /tmp/app/* .
```

### Install a theme

`cd` into the `themes` directory and clone the `nanobox-hugo-theme` repo:

```
cd themes
git clone https://github.com/sdomino/nanobox-hugo-theme
```

To use the theme *either* copy the entire `config.toml` that comes with the theme, or just add the theme to your existing `config.toml`

```
# copy the config.toml that comes with the theme
cp ./themes/nanobox-hugo-theme/config.toml config.toml

# or, add it to your existing config.toml
theme = "nanobox-hugo-theme"
```

{{% note %}}
It is not intended that you use the `nanobox-hugo-theme` as your actual theme. It's simply a theme to start with and should be replaced.
{{% /note %}}

### View Your App

To view your application simply run the following command from a Nanobox console:

```
hugo server --bind="0.0.0.0" --baseUrl=$APP_IP
```

![hugo server](/images/hosting-and-deployment/deployment-with-nanobox/hugo-server.png)

With that you should be able to visit your app at the given IP:1313 address

{{% note %}}
You can [add a custom DNS alias](https://docs.nanobox.io/cli/dns/#add) to make it easier to access your app. Run `nanobox dns add local hugo.dev`. After starting your server, visit your app at [hugo.dev:1313](http://hugo.dev:1313)
{{% /note %}}

### Develop, Develop, Develop

{{% note %}}
IMPORTANT: One issue we are aware of, and actively investigating, is livereload. Currently, livereload does not work when developing Hugo applications with Nanobox.
{{% /note %}}

With Hugo installed you're ready to go. Develop Hugo like you would normally (using all the generators, etc.). Once your app is ready to deploy, run `hugo` to generate your static assets and get ready to deploy!


## Deploying Hugo with Nanobox

{{% note %}}
If you haven't already, make sure to [connect a hosting account](#connect-a-host-account) to your Nanobox account, and [launch a new application](#launch-a-new-application-on-nanobox) in the Dashboard.
{{% /note %}}

To deploy your application to Nanobox you simply need to [link your local codebase](https://docs.nanobox.io/workflow/deploy-code/#add-your-live-app-as-a-remote) to an application you've created on Nanobox. That is done with the following command:

```
nanobox remote add <your-app-name>
```

{{% note %}}
You may be prompted to login using your ***Nanobox credentials*** at this time
{{% /note %}}

### Stage Your Application (optional)

Nanobox gives you the ability to [simulate your production environment locally](https://docs.nanobox.io/workflow/deploy-code/#preview-locally). While staging is optional it's always recommended, so there's no reason not to!

To stage your app simply run:

```
nanobox deploy dry-run
```

Now visit your application with the IP address provided.

![nanobox deploy dry-run](/images/hosting-and-deployment/deployment-with-nanobox/nanobox-deploy-dry-run.png)

### Deploy Your Application

Once everything checks out and you're [ready to deploy](https://docs.nanobox.io/workflow/deploy-code/#deploy-to-production), simply run:

```
nanobox deploy
```

Within minutes you're Hugo app will be deployed to your host and humming along smoothly. That's it!
