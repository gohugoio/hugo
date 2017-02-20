---
title: Hosting on GitLab
linktitle: Hosting on GitLab
description:
date: 2016-06-23
publishdate: 2016-06-23
lastmod: 2016-06-23
categories: [hosting and deployment]
tags: [hosting,deployment,git,gitlab]
authors: [Riku-Pekka Silvola]
weight:
draft: false
toc: true
needsreview: true
aliases: [/tutorials/hosting-on-gitlab/]
notesforauthors:
---


[GitLab](https://gitlab.com/) makes it incredibly easy to build, deploy, and host your Hugo website.

## Assumptions

* Working familiarity with Git for version control
* Completion of the [Quick Start][quickstart]
* A [GitLab account](https://gitlab.com/users/sign_in)
* A Hugo website on your local machine that you are ready to publish

## Create .gitlab-ci.yml

```bash
cd your-hugo-site
```

In the root directory of your Hugo site, create a `.gitlab-ci.yml` file. The `.gitlab-ci.yml` configures the GitLab CI on how to build your page. Simply add the content below.

```yml
image: publysher/hugo

pages:
  script:
  - hugo
  artifacts:
    paths:
    - public
  only:
  - master
```

## Push Your Hugo Website to GitLab

Next up, create a new repository on GitLab. It is *not* necessary to set the repository public. In addition, you might want to add `/public` to your .gitignore file, as there is no need to push compiled assets to GitLab.

```bash
# initialize new git repository
git init

# add /public directory to our .gitignore file
echo "/public" >> .gitignore

# commit and push code to master branch
git add .
git commit -m "Initial commit"
git remote add origin https://gitlab.com/YourUsername/your-hugo-site.git
git push -u origin master
```

## Wait for Your Page to be Built

That's it! You can now follow the CI agent building your page at https://gitlab.com/YourUsername/your-hugo-site/pipelines.

After the build has passed, your new website is available at `https://YourUsername.gitlab.io/your-hugo-site/`.

## Next Steps

GitLab supports using custom CNAME's and TLS certificates. For more details on GitLab Pages, see [https://about.gitlab.com/2016/04/07/gitlab-pages-setup/](https://about.gitlab.com/2016/04/07/gitlab-pages-setup/).

[quickstart]: /getting-started/quick-start/