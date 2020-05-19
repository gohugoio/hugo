---
title: Host on Tiiny.Host
linktitle: Host on Tiiny.Host
description: You can use tiiny.host to deploy your static hugo website, no registration required.
date: 2020-05-19
publishdate: 2020-05-19
lastmod: 2017-03-15
categories: [hosting and deployment]
keywords: [hosting, tiiny]
authors: [Elston Baretto]
menu:
  docs:
    parent: "hosting-and-deployment"
    weight: 20
weight: 20
sections_weight: 20
draft: false
toc: true
aliases: []
---

## Assumptions

You have completed the [Quick Start][] or have a completed Hugo website ready for deployment.

## Deploy

### Step 1: Prepare your build

Build your site by running this command in your project's root directory:

```
hugo
```

This will generate a publishable version of your site in the `./public` folder.

### Step 2: Deploy

Zip your `./public` folder and head over to https://tiiny.host.

Here, enter a subdomain, upload your zip file and click `Launch`.

And that's it, you're done! There is no need to register and your site will be live for 7 days or can be extended on the Pro plan.

## Reference links

* [tiiny.host FAQ](https://tiiny.host/help/)

[Quick Start]: /getting-started/quick-start/