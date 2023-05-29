---
title: Host on GitLab
description: GitLab makes it easy to build, deploy, and host your Hugo website via their free GitLab Pages service, which provides native support for Hugo.
categories: [hosting and deployment]
keywords: [hosting,deployment,git,gitlab]
menu:
  docs:
    parent: hosting-and-deployment
    weight: 40
weight: 40
toc: true
aliases: [/tutorials/hosting-on-gitlab/]
---

## Assumptions

* Working familiarity with Git for version control
* Completion of the Hugo [Quick Start]
* A [GitLab account](https://gitlab.com/users/sign_in)
* A Hugo website on your local machine that you are ready to publish

## BaseURL

The `baseURL` in your [site configuration](/getting-started/configuration/) must reflect the full URL of your GitLab pages repository if you are using the default GitLab Pages URL (e.g., `https://<YourUsername>.gitlab.io/<your-hugo-site>/`) and not a custom domain.

## Configure GitLab CI/CD

Define your [CI/CD](https://docs.gitlab.com/ee/ci/quick_start/) jobs by creating a `.gitlab-ci.yml` file in the root of your project.

{{< code file=".gitlab-ci.yml" >}}
image: registry.gitlab.com/pages/hugo/hugo_extended:latest

variables:
  GIT_SUBMODULE_STRATEGY: recursive

pages:
  script:
  - hugo
  artifacts:
    paths:
    - public
  rules:
  - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
{{< /code >}}

{{% note %}}
See [this list](https://gitlab.com/pages/hugo/container_registry) if you wish to use a particular Hugo version to build your site.
{{% /note %}}

## Push your Hugo website to GitLab

Next, create a new repository on GitLab. It is *not* necessary to make the repository public. In addition, you might want to add `/public` to your .gitignore file, as there is no need to push compiled assets to GitLab or keep your output website in version control.

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

## Wait for your page to build

That's it! You can now follow the CI agent building your page at `https://gitlab.com/<YourUsername>/<your-hugo-site>/pipelines`.

After the build has passed, your new website is available at `https://<YourUsername>.gitlab.io/<your-hugo-site>/`.

## Next steps

GitLab supports using custom CNAME's and TLS certificates. For more details on GitLab Pages, see the [GitLab Pages setup documentation](https://about.gitlab.com/2016/04/07/gitlab-pages-setup/).

[Quick Start]: /getting-started/quick-start/
