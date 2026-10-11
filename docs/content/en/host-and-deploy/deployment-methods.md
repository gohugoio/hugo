---
title: Deployment methods
description: Choose a method to deploy your site.
categories: []
keywords: []
weight: 10
---

You can deploy your site in one of two ways: push your source code to a Git repository and let the hosting platform build and deploy it, or build it locally and transfer the files to your host. Choose a method below, then follow the guide for your tool or platform.

## Push to a Git repository

Push your source code to a Git repository, and the hosting platform builds and deploys your site each time you push. Choose this method if you use Git and want [CI/CD](g) without managing a server. Follow the guide for your hosting platform:

- [AWS Amplify][]
- [Azure Static Web Apps][]
- [Cloudflare][]
- [Codeberg Pages][]
- [Firebase][]
- [GitHub Pages][]
- [GitLab Pages][]
- [Netlify][]
- [Render][]
- [SourceHut Pages][]
- [Vercel][]

## Build locally and transfer

Build your site on your local machine or in your own environment, then transfer the generated files to your host with a command-line tool or a graphical client such as an FTP client. Choose this method if you deploy to a web server, a shared hosting account, or a cloud storage bucket, or if you do not want to use Git. Follow the guide for your transfer tool:

- [FTP][]
- [Hugo][]
- [Rclone][]
- [Rsync][]

[AWS Amplify]: /host-and-deploy/deploy-to-aws-amplify/
[Azure Static Web Apps]: /host-and-deploy/deploy-to-azure-static-web-apps/
[Cloudflare]: /host-and-deploy/deploy-to-cloudflare/
[Codeberg Pages]: /host-and-deploy/deploy-to-codeberg-pages/
[Firebase]: /host-and-deploy/deploy-to-firebase/
[GitHub Pages]: /host-and-deploy/deploy-to-github-pages/
[GitLab Pages]: /host-and-deploy/deploy-to-gitlab-pages/
[Netlify]: /host-and-deploy/deploy-to-netlify/
[Render]: /host-and-deploy/deploy-to-render/
[SourceHut Pages]: /host-and-deploy/deploy-to-sourcehut-pages/
[Vercel]: /host-and-deploy/deploy-to-vercel/
[FTP]: /host-and-deploy/deploy-with-ftp/
[Hugo]: /host-and-deploy/deploy-with-hugo/
[Rclone]: /host-and-deploy/deploy-with-rclone/
[Rsync]: /host-and-deploy/deploy-with-rsync/
