---
title: Host on Netlify
linktitle: Host on Netlify
description: Netlify can host your Hugo site with CDN, continuous deployment, 1-click HTTPS, an admin GUI, and its own CLI.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-03-11
categories: [hosting and deployment]
keywords: [netlify,hosting,deployment]
authors: [Ryan Watters, Seth MacLeod]
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

[Netlify][netlify] provides continuous deployment services, global CDN, ultra-fast DNS, atomic deploys, instant cache invalidation, one-click SSL, a browser-based interface, a CLI, and many other features for managing your Hugo website.

## Assumptions

* You have an account with GitHub, GitLab, or Bitbucket.
* You have completed the [Quick Start][] or have Hugo website you are ready to deploy and share with the world.
* You do not already have a Netlify account.

## Create a Netlify account

Go to [app.netlify.com][] and select your preferred signup method. This will likely be a hosted Git provider, although you also have the option to sign up with an email address.

The following examples use GitHub, but other git providers will follow a similar process.

![Screenshot of the homepage for app.netlify.com, containing links to the most popular hosted git solutions.](/images/hosting-and-deployment/hosting-on-netlify/netlify-signup.jpg)

Selecting GitHub will bring up a typical modal you've seen through other application that use GitHub for authentication. Select "Authorize application."

![Screenshot of the authorization popup for Netlify and GitHub.](/images/hosting-and-deployment/hosting-on-netlify/netlify-first-authorize.jpg)

## Create a New Site with Continuous Deployment

You're now already a Netlify member and should be brought to your new dashboard. Select "New site from git."

![Screenshot of the blank Netlify admin panel with no sites and highlighted 'add new site' button'](/images/hosting-and-deployment/hosting-on-netlify/netlify-add-new-site.jpg)

Netlify will then start walking you through the steps necessary for continuous deployment. First, you'll need to select your git provider again, but this time you are giving Netlify added permissions to your repositories.

![Screenshot of step 1 of create a new site for Netlify: selecting the git provider](/images/hosting-and-deployment/hosting-on-netlify/netlify-create-new-site-step-1.jpg)

And then again with the GitHub authorization modal:

![Screenshot of step 1 of create a new site for Netlify: selecting the git provider](/images/hosting-and-deployment/hosting-on-netlify/netlify-authorize-added-permissions.jpg)

Select the repo you want to use for continuous deployment. If you have a large number of repositories, you can filter through them in real time using repo search:

![Screenshot of step 1 of create a new site for Netlify: selecting the git provider](/images/hosting-and-deployment/hosting-on-netlify/netlify-create-new-site-step-2.jpg)

Once selected, you'll be brought to a screen for basic setup. Here you can select the branch you wanted published, your [build command][], and your publish (i.e. deploy) directory. The publish directory should mirror that of what you've set in your [site configuration][config], the default of which is `public`. The following steps assume you are publishing from the `master` branch.

### Build with a Specific Hugo Version

Setting the build command to `hugo` will build your site according to the current default Hugo version used by Netlify. You can see the full list of [available Hugo versions in Netlify's Docker file][hugoversions].

If you want to tell Netlify to build with a specific version (hugo <= 0.20), you can append an underscore followed by the version number to the build command:

```
hugo_0.19
```

Your simple configuration should now look similar to the following:

![Screenshot of 3-step, basic continuous deployment setup with a new Hugo site on Netlify](/images/hosting-and-deployment/hosting-on-netlify/netlify-create-new-site-step-3.jpg)

For version hugo > 0.20 you have to [specify version hugo for testing and production](https://www.netlify.com/blog/2017/04/11/netlify-plus-hugo-0.20-and-beyond/) in `netlify.toml` file or set `HUGO_VERSION` as a build environment variable in the Netlify console.

For production:

```
[context.production.environment]
  HUGO_VERSION = "0.26"
```

For testing:

```
[context.deploy-preview.environment]
  HUGO_VERSION = "0.26"
```  

Selecting "Deploy site" will immediately take you to a terminal for your build:.

![Animated gif of deploying a site to Netlify, including the terminal read out for the build.](/images/hosting-and-deployment/hosting-on-netlify/netlify-deploying-site.gif)

Once the build is finished---this should only take a few seconds--you should now see a "Hero Card" at the top of your screen letting you know the deployment is successful. The Hero Card is the first element that you see in most pages. It allows you to see a quick summary of the page and gives access to the most common/pertinent actions and information. You'll see that the URL is automatically generated by Netlify. You can update the URL in "Settings."

![Screenshot of successful deploy badge at the top of a deployments screen from within the Netlify admin.](/images/hosting-and-deployment/hosting-on-netlify/netlify-deploy-published.jpg)

![Screenshot of homepage to https://hugo-netlify-example.netlify.com, which is mostly dummy text](/images/hosting-and-deployment/hosting-on-netlify/netlify-live-site.jpg)

[Visit the live site][visit].

Now every time you push changes to your hosted git repository, Netlify will rebuild and redeploy your site.

## Use Hugo Themes with Netlify

The [`git clone` method for installing themes][installthemes] is not supported by Netlify. If you were to use `git clone`, it would require you to recursively remove the `.git` subdirectory from the theme folder and would therefore prevent compatibility with future versions of the theme.

A *better* approach is to install a theme as a proper git submodule. You can [read the GitHub documentation for submodules][ghsm] or those found on [Git's website][gitsm] for more information, but the command is similar to that of `git clone`:

```
cd themes
git submodule add https://github.com/<THEMECREATOR>/<THEMENAME>
```

## Next Steps

You now have a live website served over https, distributed through CDN, and configured for continuous deployment. Dig deeper into the Netlify documentation:

1. [Using a Custom Domain][]
2. [Setting up HTTPS on Custom Domains][httpscustom]
3. [Redirects and Rewrite Rules][]


[app.netlify.com]: https://app.netlify.com
[build command]: /getting-started/usage/#the-hugo-command
[config]: /getting-started/configuration/
[ghsm]: https://github.com/blog/2104-working-with-submodules
[gitsm]: https://git-scm.com/book/en/v2/Git-Tools-Submodules
[httpscustom]: https://www.netlify.com/docs/ssl/
[hugoversions]: https://github.com/netlify/build-image/blob/master/Dockerfile#L166
[installthemes]: /themes/installing/
[netlify]: https://www.netlify.com/
[netlifysignup]: https://app.netlify.com/signup
[Quick Start]: /getting-started/quick-start/
[Redirects and Rewrite Rules]: https://www.netlify.com/docs/redirects/
[Using a Custom Domain]: https://www.netlify.com/docs/custom-domains/
[visit]: https://hugo-netlify-example.netlify.com
