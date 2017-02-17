---
title: Add Your Site to the Showcase
linktitle: Add Your Site to the Showcase
description: Proud of a site you built with Hugo? Add it to the Official Hugo Site Showcase.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [contribute to hugo]
tags: [dev,showcase]
weight:
draft: false
slug:
aliases: []
toc: false
notes:
---

## Showcase Additions

You got your new website running and it's powered by Hugo? Great. You can add your website with a few steps to the [showcase](/showcase/).

First, make sure that you created a [fork](https://help.github.com/articles/fork-a-repo/) of Hugo on Github and cloned your fork on your local computer. Next, create a separate branch for your additions:

```
# You can choose a different descriptive branch name if you like
git checkout -b showcase-addition
```

Let's create a new document that contains some metadata of your homepage. Replace `example` in the following examples with something unique like the name of your website. Inside the terminal enter the following commands:

```
cd docs
hugo new showcase/example.md
```

You should find the new file at `content/showcase/example.md`. Open it in an editor. The file should contain a frontmatter with predefined variables like below:

```
---
date: 2016-02-12T21:01:18+01:00
description: ""
license: ""
licenseLink: ""
sitelink: http://spf13.com/
sourceLink: https://github.com/spf13/spf13.com
tags:
- personal
- blog
thumbnail: /images/spf13-tn.jpg
title: example
---
```

Add at least values for `sitelink`, `title`,  `description` and a path for `thumbnail`.

Furthermore, we need to create the thumbnail of your website. **It's important that the thumbnail has the required dimensions of 600px by 400px.** Give your thumbnail a name like `example-tn.png`. Save it under `docs/static/img/`.

Check a last time that everything works as expected. Start Hugo's built-in server in order to inspect your local copy of the showcase in the browser:

hugo server

If everything looks fine, we are ready to commit your additions. For the sake of best practices, please make sure that your commit follows our [code contribution guideline][].

{{% input "commit-site.sh" %}}
```git
git commit -m "docs: Add example.com to the showcase"
```
{{% /input %}}

Last but not least, we're ready to create a [pull request](https://github.com/spf13/hugo/compare).

Don't forget to accept the contributor license agreement. Click on the yellow badge in the automatically added comment in the pull request.

[code contribution guideline]: https://github.com/spf13/hugo#code-contribution-guideline