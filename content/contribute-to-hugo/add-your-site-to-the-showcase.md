---
title: Add Your Site to the Showcase
linktitle: Add Your Site to the Showcase
description: Proud of a site you built with Hugo? Add it to the Official Hugo Site Showcase.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [contribute to hugo]
tags: [dev,showcase]
weight: 30
draft: false
aliases: []
toc: true
notesforauthors:
---


If you've built a site with Hugo and would like to have it featured on the official Hugo site, you can add your website with a few steps to the [Site Showcase][].

## 1. Create Your Fork

First, make sure that you created a [fork](https://help.github.com/articles/fork-a-repo/) of Hugo on Github and cloned your fork on your local computer. Next, create a separate branch for your additions. Note that you can choose a different descriptive branch name if you like:

```git
git checkout -b showcase-addition
```

## 2. Add Your Showcase File via the `showcase` Archetype

Let's create a new document that contains some metadata of your homepage. Replace `example` in the following examples with something unique like the name of your website. Inside the terminal enter the following commands:

```bash
cd docs
hugo new showcase/my-hugo-site-name.md
```

You should find the new file at `content/showcase/your-site-name.md`. Open the file in your preferred text editor. The file should contain front matter with predefined variables like below:

```yaml
---
description: ""
lastmod: ""
license: ""
licenseLink: ""
sitelink: ""
sourcelink: ""
categories: [showcase]
tags: []
image: "yourimage.jpg"
toc: false
title: my hugo site name
notesforauthors: "For the image, only include the file name *without* a directory/path, which is taken care of in the templating. See the showcase contribution page at gohugo.io/contribute-to-hugo/add-your-site-to-the-showcase/ for more details. As always, feel free to delete `notesforauthors` or modify for anyone in the future who may edit the content in this file."
---
```

Add at least values for `sitelink`, `title`,  `description`, and a path for `thumbnail`.

{{% note "Notes for Authors" %}}
You may notice a `notesforauthors` key-value in your new content file for the showcase. These notes are not required, but rather have been added to make it easier to fill out the required metadata without needing to refer to the Hugo docs website. You can delete this metadata before submitting a pull request.
{{% /note %}}

## 3. Add an Image

We need to create the thumbnail of your website. Give your thumbnail a name like `my-hugo-site-name.png`. Save it under [`docs/static/images/showcase/`][].

{{% warning "Thumbnail Size" %}}
It's important that the thumbnail has the required dimensions of 600px by 400px or the site will not render appropriately. Be sure to optimize your image as a matter of best practice.
{{% /warning %}}

Check one last time that everything looks complete. Start Hugo's built-in server in order to inspect your local copy of the showcase in the browser.

```bash
hugo server
```

## 4. Commit and Submit a Pull Request

If everything looks fine, we are ready to commit your additions. For the sake of best practices, please make sure that your commit follows our [code contribution guideline][].

{{% input "commit-site.sh" %}}
```git
git commit -m "docs: Add example.com to the showcase"
```
{{% /input %}}

Last but not least, we're ready to create a [pull request].

### Contributor License Agreement

Don't forget to accept the contributor license agreement. Click on the yellow badge in the automatically added comment in the pull request to accept.

[code contribution guideline]: https://github.com/spf13/hugo#code-contribution-guideline
[pull request]: https://github.com/spf13/hugo/compare
[Site Showcase]: /showcase/
[`docs/static/images/showcase/`]: https://github.com/spf13/hugo/tree/master/docs/static/images/showcase/