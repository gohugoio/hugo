---
title: Host on GitLab
linktitle: Host on GitLab
description: GitLab makes it incredibly easy to build, deploy, and host your Hugo website via their free GitLab Pages service, which provides native support for Hugo.
date: 2016-06-23
publishdate: 2016-06-23
lastmod: 2017-11-16
categories: [hosting and deployment]
keywords: [hosting,deployment,git,gitlab]
authors: [Riku-Pekka Silvola]
menu:
  docs:
    parent: "hosting-and-deployment"
    weight: 40
weight: 40
sections_weight: 40
draft: false
toc: true
wip: false
aliases: [/tutorials/hosting-on-gitlab/]
---

[GitLab](https://gitlab.com/) makes it incredibly easy to build, deploy, and host your Hugo website via their free GitLab Pages service, which provides [native support for Hugo, as well as numerous other static site generators](https://gitlab.com/pages/hugo).

## Assumptions

* Working familiarity with Git for version control
* Completion of the Hugo [Quick Start][]
* A [GitLab account](https://gitlab.com/users/sign_in)
* A Hugo website on your local machine that you are ready to publish

## Create .gitlab-ci.yml

```
cd your-hugo-site
```

In the root directory of your Hugo site, create a `.gitlab-ci.yml` file. The `.gitlab-ci.yml` configures the GitLab CI on how to build your page. Simply add the content below.

{{< code file=".gitlab-ci.yml" >}}
image: monachus/hugo

variables:
  GIT_SUBMODULE_STRATEGY: recursive

pages:
  script:
  - hugo
  artifacts:
    paths:
    - public
  only:
  - master
{{< /code >}}

## Push Your Hugo Website to GitLab

Next, create a new repository on GitLab. It is *not* necessary to make the repository public. In addition, you might want to add `/public` to your .gitignore file, as there is no need to push compiled assets to GitLab or keep your output website in version control.

```
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

## Wait for Your Page to Build

That's it! You can now follow the CI agent building your page at `https://gitlab.com/<YourUsername>/<your-hugo-site>/pipelines`.

After the build has passed, your new website is available at `https://<YourUsername>.gitlab.io/<your-hugo-site>/`.

## Next Steps

GitLab supports using custom CNAME's and TLS certificates. For more details on GitLab Pages, see the [GitLab Pages setup documentation](https://about.gitlab.com/2016/04/07/gitlab-pages-setup/).

[Quick Start]: /getting-started/quick-start/
