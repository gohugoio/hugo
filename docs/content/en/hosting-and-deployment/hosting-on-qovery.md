---
title: Host on Qovery
linktitle: Host on Qovery
description: Fully-managed cloud platform that runs on your AWS, GCP, Azure and Digital Ocean account where you can host static sites, backend APIs, databases, cron jobs, and all your other apps in one place.
date: 2021-03-18
publishdate: 2021-03-18
lastmod: 2021-03-18
categories: [hosting and deployment]
keywords: [render,hosting,deployment]
authors: [Arnaud Jeannin]
menu:
  docs:
    parent: "hosting-and-deployment"
    weight: 10
weight: 10
sections_weight: 10
draft: false
aliases: []
toc: true
---

[Qovery](https://qovery.com) is a fully-managed cloud platform that runs on your AWS, GCP, Azure and Digital Ocean account where you can host static sites, backend APIs, databases, cron jobs, and all your other apps in one place.

Static sites are **completely free** on Qovery and include the following:

- Continuous, automatic builds & deploys from GitHub, Bitbucket, and GitLab.
- Automatic SSL certificates through [Let's Encrypt](https://letsencrypt.org).
- Free managed PostgreSQL.
- Free SSD storage.
- Unlimited collaborators.
- Unlimited [custom domains](https://docs.qovery.com/guides/getting-started/setting-custom-domain/).

## Prerequisites

This guide assumes you already have a Hugo project to deploy. If you need a project, use the [Quick Start][] to get started.

## Setup

Follow the procedure below to set up a Hugo on Qovery:

### 1. Create a Qovery account.

Visit the [Qovery dashboard](https://console.qovery.com) to create an account if you don't already have one.

### 2. Create a project

Click on "Create a new project" and give a name to your project. 

Click on "Next".

### 3. Add an application

Click on "Create an application" then choose "I have an application" and select your GitHub or GitLab repository where your Hugo site is located.

Click on "Next".

Skip adding services for static website.

Click on "Deploy".

## Deploy

Your app should be deployed. You can see the status in real time by clicking on deployment logs.

## Continuous deploys

Now that Qovery is connected to your repo, it will **automatically build and publish your site** any time you push to GitHub.

## Custom domains

Add your own domains to your site easily using Qovery's [custom domains](https://docs.qovery.com/guides/getting-started/setting-custom-domain/) guide.

## Support

Chat with Qovery developers on [Discord](https://discord.qovery.com) if you need help.

[Quick Start]: /getting-started/quick-start/

