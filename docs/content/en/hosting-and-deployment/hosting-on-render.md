---
title: Host on Render
linktitle: Host on Render
description: Host your Hugo site for free with Render's global CDN, fully-managed SSL and auto deploys from GitHub.
date: 2019-06-06
publishdate: 2019-06-06
lastmod: 2020-01-01
categories: [hosting and deployment]
keywords: [render,hosting,deployment]
authors: [Anurag Goel]
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

## Introduction

[Render](https://render.com) is a fully-managed cloud platform where you can host static sites, backend APIs, databases, cron jobs, and all your other apps in one place.

Static sites are **completely free** on Render and include the following:

- Continuous, automatic builds & deploys from [GitHub](https://render.com/docs/github) and [GitLab](https://render.com/docs/gitlab).
- Automatic SSL certificates through [Let's Encrypt](https://letsencrypt.org).
- Instant cache invalidation with a lightning fast, global CDN.
- Unlimited collaborators.
- Unlimited [custom domains](https://render.com/docs/custom-domains).
- Automatic [Brotli compression](https://en.wikipedia.org/wiki/Brotli) for faster sites.
- Native HTTP/2 support.
- [Pull Request Previews](https://render.com/docs/pull-request-previews).
- Automatic HTTP → HTTPS redirects.
- Custom URL redirects and rewrites.

## Assumptions

* You have an account with GitHub or GitLab.
* You have completed the [Quick Start][] or have a Hugo website you are ready to deploy and share with the world.
* You have a Render account. You can sign up at https://render.com/register.

## Deployment

You can set up a Hugo site on Render in two quick steps:

1. Create a new **Web Service** on Render, and give Render permission to access your GitHub/Gitlab repo.
2. Use the following values during creation:

  Field                | Value
  -------------------  |  -------------------
 **Environment**       | `Static Site`
 **Build Command**     | `hugo --gc --minify` (or your own build command)
 **Publish Directory** | `public` (or your own output directory)

That's it! Your site will be live on your Render URL (which looks like `yoursite.onrender.com`) as soon as the build is done.

## Continuous Deploys

Now that Render is connected to your repo, it will **automatically build and publish your site** any time you push to your GitHub/Gitlab.

You can choose to disable auto deploys under the **Settings** section for your site and deploy it manually from the Render dashboard.

## CDN and Cache Invalidation

Render hosts your site on a global, lightning fast CDN which ensures the fastest possible download times for all your users across the globe.

Every deploy automatically and instantly invalidates the CDN cache, so your users can always access the latest content on your site.

## Custom Domains

Add your own domains to your site easily using Render's [custom domains](https://render.com/docs/custom-domains) guide.

## Pull Request Previews

With Pull Request (PR) previews, you can visualize changes introduced in a pull request instead of simply relying on code reviews.

Once enabled, every PR for your site will automatically generate a new static site based on the code in the PR. It will have its own URL, and it will be deleted automatically when the PR is closed.

Read more about [Pull Request Previews](https://render.com/docs/pull-request-previews) on Render.

## Hugo Themes

Render automatically downloads all Git submodules defined in your Git repo on every build. This way Hugo themes added as submodules work as expected.

## Support

Chat with Render developers at https://render.com/chat or email `support@render.com` if you need help.


[Quick Start]: /getting-started/quick-start/
