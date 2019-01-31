---
title: Host on AWS Amplify
linktitle: Host on AWS Amplify
description: Netlify can host your Hugo site with CDN, continuous deployment, 1-click HTTPS, an admin GUI, and its own CLI.
date: 2018-01-31
publishdate: 2018-01-31
lastmod: 2018-01-31
categories: [hosting and deployment]
keywords: [amplify,hosting,deployment]
authors: [Nikhil Swaminathan]
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

In this guide we'll walk through how to deploy and host your Hugo site using the [AWS Amplify Console](https://console.amplify.aws).

AWS Amplify is a combination of client library, CLI toolchain, and a Console for continuous deployment and hosting. The Amplify CLI and library allow developers to get up & running with full-stack cloud-powered applications with features like authentication, storage, serverless GraphQL or REST APIs, analytics, Lambda functions, & more. The Amplify Console provides continuous deployment and hosting for modern web apps (single page apps and static site generators). Continuous deployment allows developers to deploy updates to their web app on every code commit to their Git repository. Hosting includes features such as globally available CDNs, easy custom domain setup + HTTPS, feature branch deployments, and password protection.

## Pre-requisites

* [Sign up for an AWS Account](https://portal.aws.amazon.com/billing/signup?redirect_url=https%3A%2F%2Faws.amazon.com%2Fregistration-confirmation). There are no upfront charges or any term commitments to create an AWS account and signing up gives you immediate access to the AWS Free Tier.
* You have an account with GitHub, GitLab, or Bitbucket.
* You have completed the [Quick Start][] or have a Hugo website you are ready to deploy and share with the world.
