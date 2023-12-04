---
title: Host on AWS Amplify
description: Develop and deploy a cloud-powered web app with AWS Amplify.
categories: [hosting and deployment]
keywords: [hosting,amplify]
menu:
  docs:
    parent: hosting-and-deployment
toc: true
---

In this guide we'll walk through how to deploy and host your Hugo site using the [AWS Amplify Console](https://console.amplify.aws).

AWS Amplify is a combination of client library, CLI toolchain, and a Console for continuous deployment and hosting. The Amplify CLI and library allow developers to get up & running with full-stack cloud-powered applications with features like authentication, storage, serverless GraphQL or REST APIs, analytics, Lambda functions, & more. The Amplify Console provides continuous deployment and hosting for modern web apps (single page apps and static site generators). Continuous deployment allows developers to deploy updates to their web app on every code commit to their Git repository. Hosting includes features such as globally available CDNs, easy custom domain setup + HTTPS, feature branch deployments, and password protection.

## Pre-requisites

* [Sign up for an AWS Account](https://portal.aws.amazon.com/billing/signup?redirect_url=https%3A%2F%2Faws.amazon.com%2Fregistration-confirmation). There are no upfront charges or any term commitments to create an AWS account and signing up gives you immediate access to the AWS Free Tier.
* You have an account with GitHub, GitLab, or Bitbucket.
* You have completed the [Quick Start] or have a Hugo website you are ready to deploy and share with the world.

## Hosting

1. Log in to the [AWS Amplify Console](https://console.aws.amazon.com/amplify/home) and choose Get Started under Deploy.
   ![Hugo Amplify](/images/hosting-and-deployment/hosting-on-aws-amplify/amplify-gettingstarted.png)

1. Connect a branch from your GitHub, Bitbucket, GitLab, or AWS CodeCommit repository. Connecting your repository allows Amplify to deploy updates on every code commit to a branch.
   ![Hugo Amplify](/images/hosting-and-deployment/hosting-on-aws-amplify/amplify-connect-repo.gif)

1. Accept the default build settings. The Amplify Console automatically detects your Hugo build settings and output directory.
   ![Hugo Amplify](/images/hosting-and-deployment/hosting-on-aws-amplify/amplify-build-settings.png)

1. Review your changes and then choose **Save and deploy**. The Amplify Console will pull code from your repository, build changes to the backend and frontend, and deploy your build artifacts at `https://master.unique-id.amplifyapp.com`. Bonus: Screenshots of your app on different devices to find layout issues.

## Using a newer version of Hugo

If you need to use a different, perhaps newer, version of Hugo than the version currently supported by AWS Amplify:

1. Visit the [AWS Amplify Console](https://console.aws.amazon.com/amplify/home), and click the app you would like to modify
1. In the side navigation bar, Under App Settings, click **Build settings**
1. On the Build settings page, near the bottom, there is a section called **Build image settings**. Click **Edit**
1. Under **Live package updates**, click **Add package version override**
1. From the selection, click **Hugo** and ensure the version field says `latest`
1. Click **Save** to save the changes.

[Quick Start]: /getting-started/quick-start/
