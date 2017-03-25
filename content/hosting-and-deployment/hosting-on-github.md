---
title: Hosting on GitHub
linktitle: Hosting on GitHub
description: Deploy Hugo as a GitHub Pages project or personal/organizational site and automate the whole process with a simple shell script.
date: 2014-03-21
publishdate: 2014-03-21
lastmod: 2017-02-01
categories: [hosting and deployment]
tags: [github,git,deployment,hosting]
authors: [Spencer Lyon, Gunnar Morling]
weight: 30
draft: false
toc: true
aliases: [/tutorials/github-pages-blog/]
---

{{% note %}}
This hosting and deployment guide was originally contributed as a tutorial by [Spencer Lyon](http://spencerlyon.com/) and [Gunnar Morling](https://github.com/gunnarmorling/).
{{% /note %}}

## Assumptions

The following sections are based on the assumption that you are working with a "Project Pages Site". This means that you'll have your Hugo sources and the generated HTML output within a single repository (in contrast, with a "User/Organization Pages Site", you'd have one repo for the sources and another repo for the published HTML files; refer to the [GitHub Pages docs](https://help.github.com/articles/user-organization-and-project-pages/) to learn more).

## Deployment via `/docs` Folder on Master Branch

[As described](https://help.github.com/articles/configuring-a-publishing-source-for-github-pages/#publishing-your-github-pages-site-from-a-docs-folder-on-your-master-branch) in the GitHub Pages docs, you can deploy from a folder called _docs_ on your master branch. This requires to change the Hugo publish directory in your [site configuration[config]; e.g., in `config.toml` and `config.yaml`, respectively:

```yaml
publishDir: docs
```

```toml
publishDir = "docs"
```

After running `hugo`, push your master branch to the remote repo and choose the _docs_ folder as the website source of your repo (in your GitHub project, go to "Settings " -> "GitHub Pages" -> "Source" -> Select "master branch /docs folder"). If that option isn't enabled, you likely haven't pushed your _docs_ folder yet.

This is the simplest approach but requires the usage of a non-standard publish directory (GitHub Pages cannot be configured to use another directory than _docs_ currently). Also the presence of generated files on the master branch may not be to eveyone's taste.

## Deployment via `gh-pages` Branch

Alternatively, you can deploy site through a separate branch called "gh_pages". That approach is a bit more complex but has some advantages:

* It keeps sources and generated HTML in two different branches
* It uses the default _public_ folder
* It keeps the histories of source branch and gh-pages branch fully separated from each other

### Preparations

These steps only need to be done once (replace "upstream" with the name of your remote; e.g., `origin`):

#### Add the Public Folder

First, add the _public_ folder to _.gitignore_ so it's ignored on the master branch:

```bash
echo "public" >> .gitignore
```

#### Initialize Your `gh-pages` Branch

Then initialize the gh-pages branch as an empty [orphan branch](https://git-scm.com/docs/git-checkout/#git-checkout---orphanltnewbranchgt):

```bash
git checkout --orphan gh-pages
git reset --hard
git commit --allow-empty -m "Initializing gh-pages branch"
git push upstream gh-pages
git checkout master
```

### Building and Deployment

Now check out the gh-pages branch into your _public_ folder, using git's [worktree feature](https://git-scm.com/docs/git-worktree)
(essentially, it allows you to have multiple branches of the same local repo to be checked out in different directories):

    rm -rf public
    git worktree add -B gh-pages public upstream/gh-pages

Regenerate the site using Hugo and commit the generated files on the gh-pages branch:

```bash
hugo
cd public && git add --all && git commit -m "Publishing to gh-pages" & cd ..
```

If the changes in your local gh-pages branch look alright, push them to the remote repo:

```bash
git push upstream gh-pages
```

After a short while you'll see the updated contents on your GitHub Pages site.

### Putting it into a script

To automate these steps, you can create a script with the following contents:

{{% code file="publish_to_ghpages.sh" %}}
```sh
#!/bin/sh

DIR=$(dirname "$0")

cd $DIR/..

if [[ $(git status -s) ]]
then
    echo "The working directory is dirty. Please commit any pending changes."
    exit 1;
fi

echo "Deleting old publication"
rm -rf public
mkdir public
git worktree prune
rm -rf .git/worktrees/public/

echo "Checking out gh-pages branch into public"
git worktree add -B gh-pages public upstream/gh-pages

echo "Removing existing files"
rm -rf public/*

echo "Generating site"
hugo

echo "Updating gh-pages branch"
cd public && git add --all && git commit -m "Publishing to gh-pages (publish.sh)"
```
{{% /code %}}

This will abort if there are pending changes in the working directory and also makes sure that all previously existing output files are removed. Adjust the script to taste, e.g. to include the final push to the remote repository if you don't need to take a look at the gh-pages branch before pushing. Or adding `echo yourdomainname.com >> CNAME` if you set up for your gh-pages to use customize domain.

## Deployment with Git 2.4 and Earlier

The `worktree` command was only introduced in Git 2.5. If you are still on an earlier version and cannot update, you can simply clone your local repo into the _public_ directory, only keeping the gh-pages branch:

```sh
git clone .git --branch gh-pages public
```

Having re-generated the site, you'd push back the gh-pages branch to your primary local repo:

```sh
cd public && git add --all && git commit -m "Publishing to gh-pages" && git push origin gh-pages
```

The other steps are the same as with the worktree approach.

## Hosting Personal or Organization Pages

As mentioned [in this GitHub Help article](https://help.github.com/articles/user-organization-and-project-pages/), you can host a user/organization page in addition to project pages. Here are the key differences in GitHub Pages websites for Users and Organizations:

1. You must use the `<username>.github.io` naming scheme.
2. Content from the `master` branch will be used to build and publish your GitHub Pages site.

It becomes much simpler in this case: we'll create two separate repos, one for Hugo's content, and a git submodule with the `public` folder's content in it.

### Step-by-step Instructions

1. Create on GitHub `<your-project>-hugo` repository (it will host Hugo's content)
2. Create on GitHub `<username>.github.io` repository (it will host the `public` folder: the static website)
3. `git clone <<your-project>-hugo-url> && cd <your-project>-hugo`
4. Make your website work locally (`hugo server -t <yourtheme>`)
5. Once you are happy with the results, <kbd>Ctrl</kbd>+<kbd>C</kbd> (kill server) and `rm -rf public` (don't worry, it can always be regenerated with `hugo -t <yourtheme>`)
6. `git submodule add -b master git@github.com:<username>/<username>.github.io.git public`
7. Almost done: add a `deploy.sh` script to help you (and make it executable: `chmod +x deploy.sh`):

```sh
#!/bin/bash

echo -e "\033[0;32mDeploying updates to GitHub...\033[0m"

# Build the project.
hugo # if using a theme, replace by `hugo -t <yourtheme>`

# Go To Public folder
cd public
# Add changes to git.
git add -A

# Commit changes.
msg="rebuilding site `date`"
if [ $# -eq 1 ]
  then msg="$1"
fi
git commit -m "$msg"

# Push source and build repos.
git push origin master

# Come Back
cd ..
```
7. `./deploy.sh "Your optional commit message"` to send changes to `<username>.github.io` (careful, you may also want to commit changes on the `<your-project>-hugo` repo).

That's it! Your personal page is running (after up to 10 minutes delay).

## Using a Custom Domain

If you'd like to use a custom domain for your GitHub Pages site, create a file _static/CNAME_ with the domain name as its sole contents. This will put the CNAME file to the root of the published site as required by GitHub Pages.

Refer to the [official documentation](https://help.github.com/articles/using-a-custom-domain-with-github-pages/) for further information.

[config]: /getting-started/configuration/