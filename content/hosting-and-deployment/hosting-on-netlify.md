---
title: Hosting on Netlify
linktitle: Hosting on Netlify
description: Netlify can host your Hugo site with CDN, continuous deployment, 1-click HTTPS, an admin GUI, and its own CLI.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-03-11
categories: [hosting and deployment]
tags: [netlify,hosting,deployment]
authors: [Ryan Watters, Seth MacLeod]
weight: 10
draft: false
aliases: []
toc: true
wip: true
---

## Assumptions

- Have an account with Github, GitLab, or Bitbucket
- Have completed the Quick Start or have a completed website ready for deployment

## Goals

We will connect a git repo to Netlify's continuous deployment service. Pushing a commit to your repo will automatically trigger Netlify's service and deploy your site.

## Create a Netlify account

Got to [netlify.com][netlify] and click on the signup button. Alternatively, you may go directly to their [signup page][netlifysignup].

![][1]

Choose how you would like to register your account. You will be able to connect to any service later regardless of what you choose now.

## Continuous Deployment

Click on the service that is hosting your repo.

![][2]

You will see a list of your repos. Click on the repo you wish to connect to Netlify.

![][3]

There are three settings on the Basic Settings tab. If you have multiple branches, you may choose which branch to deploy. Set the publish directory to `public` and the build command to `hugo`. After filling in the fields, click the button that says "Build your site."

![][4]

Your site is now configured for continuous deployment, and you may view your site once the build is complete.

![][5]

## Other Settings

Check out Netlify's settings page and read their documentation for further configuration. You can use custom domain names as well as free HTTPS.

[1]: /images/hosting-and-deployment/hosting-on-netlify/netlify-signup.png
[2]: /images/hosting-and-deployment/hosting-on-netlify/netlify-start.png
[3]: /images/hosting-and-deployment/hosting-on-netlify/netlify-start-repos.png
[4]: /images/hosting-and-deployment/hosting-on-netlify/netlify-configure-repo.png
[5]: /images/hosting-and-deployment/hosting-on-netlify/netlify-build-done.png

[netlify]: https://www.netlify.com/
[netlifysignup]: https://app.netlify.com/signup
