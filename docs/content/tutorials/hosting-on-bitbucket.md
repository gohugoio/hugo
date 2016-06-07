---
authors:
- Jason Gowans
lastmod: 2016-03-11
date: 2016-03-11
linktitle: Hosting on Bitbucket
toc: true
menu:
  main:
    parent: tutorials
next: /tutorials/github-pages-blog
prev: /tutorials/creating-a-new-theme
title: Continuous deployment with Bitbucket & Aerobatic
weight: 10
---

# Continuous deployment with Bitbucket & Aerobatic

## Introduction

In this tutorial, we will use [Bitbucket](https://bitbucket.org/) and [Aerobatic](https://www.aerobatic.com) to build, deploy, and host a Hugo site. Aerobatic is a static hosting service that is installed as an add-on to Bitbucket and provides a free hosting tier with custom domain and wildcard SSL certificate.

It is assumed that you know how to use git for version control and have a Bitbucket account. It is also assumed that you have gone through the [quickstart guide]({{< relref "overview/quickstart.md" >}}) and already have a Hugo site on your local machine.

## Create package.json

```bash
cd your-hugo-site
```

In the root directory of your Hugo site, create a `package.json` file. The `package.json` informs Aerobatic to build a Hugo site. 

To do so, declare the following snippet in your `package.json` manifest. You can, of course, use any [Hugo theme](http://themes.gohugo.io/) of your choosing with the `themeRepo` option. Just tell Aerobatic where the theme’s git repo is.

```json
{
  "_aerobatic": {
    "build": {
      "engine": "hugo",
      "themeRepo": "https://github.com/alexurquhart/hugo-geo.git"
    }
  }
}
```

## Push Hugo site to Bitbucket
We will now create a git repository and then push our code to Bitbucket. Because Aerobatic both builds *and* hosts Hugo sites, there is no need to push the compiled assets in the `/public` folder.

```bash
# initialize new git repository
git init

# add /public directory to our .gitignore file
echo "/public" >> .gitignore

# commit and push code to master branch
git commit -a -m "Initial commit"
git remote add origin git@bitbucket.org:YourUsername/your-hugo-site.git
git push -u origin master
```

## Install Aerobatic
Clicking [this link](https://aerobatic.io/bb) will automatically install the Aerobatic add-on to your Bitbucket account. Alternatively, you can also install Aerobatic from the Bitbucket [add-on directory](https://bitbucket.org/account/addon-directory). Click **Grant Access** in the install dialog.

![][1]

[1]: /img/tutorials/hosting-on-bitbucket/bitbucket-grant-access.png

## Setup hosting

Select your repository from the dropdown menu and click **Setup hosting**

![][2]

[2]: /img/tutorials/hosting-on-bitbucket/bitbucket-setup-hosting.png

You will then be directed to the **Create Website** screen. This is a one-time step. With each subsequent `git push` to Bitbucket, Aerobatic will automatically build and deploy a new version of your site instantly.

  - Give your website a name. 
  - In this example, we won't setup a custom domain, but [you can](https://www.aerobatic.com/docs/custom-domains-ssl). 
  - Leave the deploy branch as master.

Click the **Create website** button:

![][3]

[3]: /img/tutorials/hosting-on-bitbucket/bitbucket-create-website.png


In less than 30 seconds, your Hugo site will be built and live on the Internet at http://your-hugo-site.aerobatic.io

![][4]

[4]: /img/tutorials/hosting-on-bitbucket/bitbucket-site-built.png


![][5]

[5]: /img/tutorials/hosting-on-bitbucket/bitbucket-site-live.png

## Suggested next steps

The code for this example can be found in this [Bitbucket repository](https://bitbucket.org/dundonian/wee-hugo/src).

Aerobatic provides a number of plugins such as [custom error pages](https://www.aerobatic.com/docs/custom-error-pages), [custom redirects](https://www.aerobatic.com/docs/redirects), [basic authentication](https://www.aerobatic.com/docs/http-basic-authentication), and many other [features](https://www.aerobatic.com/features/). In the case of authentication, [this blog post](https://www.aerobatic.com/blog/password-protect-a-hugo-site) describes how to password protect all, or part, of your Hugo site.
