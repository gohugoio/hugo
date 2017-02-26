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
draft: true
aliases: []
toc: true
notesforauthors:
---


If you've built a site with Hugo and would like to have it featured on the official Hugo site, you can add your website with a few steps to the [Site Showcase][].



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
notesforauthors: "For the image, only include the file name *without* a directory/path, which is taken care of in the templating. See the showcase contribution page at gohugo.io/contribute/add-your-site-to-the-showcase/ for more details. As always, feel free to delete `notesforauthors` or modify for anyone in the future who may edit the content in this file."
---
```


