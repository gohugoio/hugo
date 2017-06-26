---
authors:
- Arjen Schwarz
- Samuel Debruyn
lastmod: 2017-02-26
date: 2015-01-12
linktitle: Automated deployments
toc: true
menu:
  main:
    parent: tutorials
next: /tutorials/creating-a-new-theme
prev: /community/contributing
title: Automated deployments with Wercker
weight: 10
---

# Automated deployments with Wercker

In this tutorial we will set up a basic Hugo project and then configure a free tool called Wercker to automatically deploy the generated site any time we add an article. We will deploy it to GitHub pages as that is easiest to set up, but you will see that we can use anything. This tutorial takes you through every step of the process, complete with screenshots and is fairly long.

The  assumptions made for this tutorial are that you know how to use git for version control, and have a GitHub account. In case you are unfamiliar with these, in their [help section](https://help.github.com/articles/set-up-git/) GitHub has an explanation of how to install and use git and you can easily sign up for a free GitHub account as well.

## Creating a basic Hugo site

There are already [pages](http://gohugo.io/overview/quickstart/) dedicated to describing how to set up a Hugo site so we will only go through the most basic steps required to get a site up and running before we dive into the Wercker configuration. All the work for setting up the project is done using the command line, and kept as simple as possible.

Create the new site using the `hugo new site` command, and we move into it.

```bash
hugo new site hugo-wercker-example
cd hugo-wercker-example
```

Add the herring-cove theme by cloning it into the theme directory using the following commands.

```bash
mkdir themes
cd themes
git clone https://github.com/spf13/herring-cove.git
```

Cloning the project like this will conflict with our own version control, so we remove the external git configuration.

```bash
rm -rf herring-cove/.git
```

Let's add a quick **about** page.

```bash
hugo new about.md
```

Now we'll edit contents/about.md to ensure it's no longer a draft and add some text to it.

```bash
hugo undraft content/about.md
```

Once completed it's a good idea to do a quick check if everything is working by running

```bash
hugo server --theme=herring-cove
```

If everything is fine, you should be able to see something similar to the image below when you go to localhost:1313 in your browser.

![][1]

[1]: /img/tutorials/automated-deployments/creating-a-basic-hugo-site.png

## Setting up version control

Adding git to our project is done by running the `git init` command from the root directory of the project.

```bash
git init
```

Running `git status` at this point will show you p entries: the **config.toml** file, the **themes** directory, the **contents** directory, and the **public** directory. We don't want the **public** directory version controlled however, as we will use wercker to generate that later on. Therefore, we'll add a gitignore file that will exclude this using the following command.

```bash
echo "/public" >> .gitignore
```

As we currently have no static files outside of the theme directory, Wercker might complain when we try to build the site later on. To prevent this, we simply have to add any file to the static folder. To keep it simple for now we'll add a robots.txt file that will give all search engines full access to the site when it's up.

```bash
echo "User-agent: *\nDisallow:" > static/robots.txt
```

After this we can add everything to the repository.

```bash
git commit -a -m "Initial commit"
```

## Adding the project to GitHub

First we'll create a new repository. You can do this by clicking on the **+** sign at the top right, or by going to https://github.com/new

We then choose a name for the project (**hugo-wercker-example**). When clicking on create repository GitHub displays the commands for adding an existing project to the site. The commands shown below are the ones used for this site, if you're following along you will need to use the ones shown by GitHub. Once we've run those commands the project is in GitHub and we can move on to setting up the Wercker configuration.

```bash
git remote add origin git@github.com:YourUsername/hugo-wercker-example.git
git push -u origin master
```

![][2]

[2]: /img/tutorials/automated-deployments/adding-the-project-to-github.png

## Welcome to wercker

Let's start by setting up an account for Wercker. To do so we'll go to http://wercker.com and click on the **Sign up** button.

![][3]

[3]: /img/tutorials/automated-deployments/wercker-sign-up.png

## Register

To make life easier for ourselves, we will then register using GitHub. If you don't have a GitHub account, or don't want to use it for your account, you can of course register with a username and password as well.

![][4]

[4]: /img/tutorials/automated-deployments/wercker-sign-up-page.png

## Connect GitHub/Bitbucket

After you are registered, you will need to link your GitHub and/or Bitbucket account to Wercker. You do this by going to your profile settings, and then "Git connections" If you registered using GitHub it will most likely look like the image below. To connect a missing service, simply click on the connect button which will then send you to either GitHub or Bitbucket where you might need to log in and approve their access to your account.

![][5]

[5]: /img/tutorials/automated-deployments/wercker-git-connections.png

## Add your project

Now that we've got all the preliminaries out of the way, it's time to set up our application. For this we click on the **+ Create** button next to Applications. Create a new application, and choose to use GitHub.

![][6]

[6]: /img/tutorials/automated-deployments/wercker-add-app.png

## Select a repository

Clicking this will make Wercker show you all the repositories you have on GitHub, but you can easily filter them as well. So we search for our repository, select it, and then click on "Use selected repo".

![][7]

[7]: /img/tutorials/automated-deployments/wercker-select-repository.png

## Configure access

As Wercker doesn't access to check out your private projects by default, it will ask you what you want to do. When your project is public, as needs to be the case if you wish to use GitHub Pages, the top choice is recommended. When you use this it will simply check out the code in the same way anybody visiting the project on GitHub can do.

![][8]

[8]: /img/tutorials/automated-deployments/wercker-access.png

## Public or not

This is a personal choice; you can make an app public so that everyone can see more details about it. This doesn't give you any real benefits either way in general, although as part of the tutorial I have of course made this app public so you can see it in action [yourself](https://app.wercker.com/#applications/5586dcbdaf7de9c51b02b0d5).

![][9]

[9]: /img/tutorials/automated-deployments/public-or-not.png

## Wercker.yml

Choose Default for your programming language. Wercker will now attempt to create an initial *wercker.yml* file for you. Or rather, it will create the code you can copy into it yourself. Because there is nothing special about our project according to Wercker, we will simply get the `debian` box. So what we do now is create a *wercker.yml* file in the local root of our project that contains the provided configuration, and after we finish setting up the app we will expand this file to make it actually do something.

![][10]

[10]: /img/tutorials/automated-deployments/werckeryml.png

## And we've got an app

The application is added now, and Wercker will be offering you the chance to trigger a build. As we haven't pushed up the **wercker.yml** file however, we will politely decline this option. Wercker has automatically added a build pipeline to your application.

## Adding build step

And now we're going to add the build step to the build pipeline. First, we go to the "Registry" action in the top menu and then search for "hugo build". Find the **Hugo-Build** task by Arjen and select it.

![][11]

[11]: /img/tutorials/automated-deployments/wercker-search.png

## Using Hugo-Build

Inside the details of this step you will see how to use it. At the top is a summary for very basic usage, but when scrolling down you go through the README of the step which will usually contain more details about the advanced options available and a full example of using the step.

We're not going to use any of the advanced features in this tutorial, so we'll return to our project and add the details we need to our wercker.yml file so that it looks like below: 

```yaml
box: debian
build:
  steps:
    - install-packages:
        packages: git
    - script:
        name: download theme
        code: |
            $(git clone https://github.com/spf13/herring-cove ./themes/herring-cove)
    - arjen/hugo-build:
        version: "0.14"
        theme: herring-cove
        flags: --buildDrafts=true
```

As you can see, we have two steps in the build pipeline. The first step downloads the theme, and the second step runs arjen's hugo-build step. To use a different theme, you can replace the link to the herring-cove source with another theme's repository - just make sure the name of the folder you download the theme to (./themes/your-theme-name) matches the theme name you tell arjen/hugo-build to use (theme: your-theme-name). Now we'll test that it all works as it should by pushing up our wercker.yml file to GitHub and seeing the magic at work.

```bash
git commit -a -m "Add wercker.yml"
git push origin master
```

Once completed a nice tick should have appeared in front of your first build, and if you want you can look at the details by clicking on it. However, we're not done yet as we still need to deploy it to GitHub Pages.

![][12]

[12]: /img/tutorials/automated-deployments/using-hugo-build.png

## Adding a deploy pipeline

In order to deploy to GitHub Pages we need to add a deploy pipeline. 

1. First, go to your Wercker application's page. Go to the "Workflows" tab and click on "Add new pipeline." Name it whatever you want; "deploy-production" or "deploy" works fine. For your YML Pipeline name, type in "deploy" without quotes. Leave the hook type as "Default" and hit the Create button. 

2. Now you need to link the deploy pipeline to your build pipeline. In the workflow editor, click on the + next to your build pipeline and add the deploy pipeline you've just made. Now the deploy pipeline will be run automatically whenever the build pipeline is completed successfully.

![][13]

[13]: /img/tutorials/automated-deployments/adding-a-deploy-pipeline.png

## Adding a deploy step

Next, we need to add a step to our deploy pipeline that will deploy the Hugo-built website to your GitHub pages repository. Once again searching through the Steps registry, we find that the most popular step is the **lukevevier/gh-pages** step so we add the configuration for that to our wercker.yml file. Additionally, we need to ensure that the box we run on has git and ssh installed. We can do this using the **install-packages** command, which then turns the wercker.yml file into this:

```yaml
box: debian
build:
  steps:
    - arjen/hugo-build:
        version: "0.14"
        theme: herring-cove
        flags: --buildDrafts=true
deploy:
  steps:
    - install-packages:
        packages: git ssh-client
    - lukevivier/gh-pages@0.2.1:
        token: $GIT_TOKEN
        domain: hugo-wercker.ig.nore.me
        basedir: public
```

How does the GitHub Pages configuration work? We've selected a couple of things, first the domain we want to use for the site. Configuring this here will ensure that GitHub Pages is aware of the domain you want to use.

Secondly we've configured the basedir to **public**, this is the directory that will be used as the website on GitHub Pages.

And lastly, you can see here that this has a **$GIT_TOKEN** variable. This is used for pushing our changes up to GitHub and we will need to configure this before we can do that. To do this, go to your application page and click on the "Environment" tab. Under Application Environment Variables, put **$GIT_TOKEN** for the Key. Now you'll need to create an access token in GitHub. How to do that is described on a [GitHub help page](https://help.github.com/articles/creating-an-access-token-for-command-line-use/). Copy and paste the access token you generated from GitHub into the Value box. **Make sure you check Protected** and then hit Add. 

![][14]

[14]: /img/tutorials/automated-deployments/adding-a-deploy-step.png

With the deploy step configured in Wercker, we can push the updated wercker.yml file to GitHub and it will create the GitHub pages site for us. The example site we used here is accessible under hugo-wercker.ig.nore.me

## Conclusion

From now on, any time you want to put a new post on your blog all you need to do is push your new page to GitHub and the rest will happen automatically. The source code for the example site used here is available on [GitHub](https://github.com/ArjenSchwarz/hugo-wercker-example), as is the [Hugo Build step](https://github.com/ArjenSchwarz/wercker-step-hugo-build) itself.

If you want to see an example of how you can deploy to S3 instead of GitHub pages, take a look at [Wercker's documentation](http://devcenter.wercker.com/docs/deploy/s3.html) about how to set that up.
