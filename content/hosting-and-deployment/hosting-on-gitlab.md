---
author: Riku-Pekka Silvola
lastmod: 2016-06-23
date: 2016-06-23
linktitle: Hosting on GitLab
toc: true
menu:
  main:
    parent: tutorials
next: /tutorials/how-to-contribute-to-hugo/
prev: /tutorials/github-pages-blog
title: Hosting on GitLab Pages
weight: 10
---
# Continuous deployment with GitLab

## Introduction

In this tutorial, we will use [GitLab](https://gitlab.com/) to build, deploy, and host a [Hugo](https://gohugo.io/) site. With Hugo and GitLab, this is incredibly easy.

It is assumed that you know how to use git for version control and have a GitLab account, and that you have gone through the [quickstart guide]({{< relref "overview/quickstart.md" >}}) and already have a Hugo site on your local machine.


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

## Push Hugo site to GitLab
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

## Wait for your page to be built
That's it! You can now follow the CI agent building your page at https://gitlab.com/YourUsername/your-hugo-site/pipelines.
After the build has passed, your new website is available at https://YourUsername.gitlab.io/your-hugo-site/

## Suggested next steps

GitLab supports using custom CNAME's and TLS certificates, but this is out of the scope of this tutorial. For more details on GitLab Pages, see [https://about.gitlab.com/2016/04/07/gitlab-pages-setup/](https://about.gitlab.com/2016/04/07/gitlab-pages-setup/)
