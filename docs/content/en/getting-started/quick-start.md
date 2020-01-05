---
title: Quick Start
linktitle: Quick Start
description: Create a Hugo site using the beautiful Ananke theme.
date: 2013-07-01
publishdate: 2013-07-01
categories: [getting started]
keywords: [quick start,usage]
authors: [Shekhar Gulati, Ryan Watters]
menu:
  docs:
    parent: "getting-started"
    weight: 10
weight: 10
sections_weight: 10
draft: false
aliases: [/quickstart/,/overview/quickstart/]
toc: true
---

{{% note %}}
This quick start uses `macOS` in the examples. For instructions about how to install Hugo on other operating systems, see [install](/getting-started/installing).

It is recommended to have [Git installed](https://git-scm.com/downloads) to run this tutorial.

For other approaches learning Hugo like book or a video tutorial refer to the [external learning resources](/getting-started/external-learning-resources/) page.
{{% /note %}}

## Step 1: Install Hugo

{{% note %}}
`Homebrew`, a package manager for `macOS`,  can be installed from [brew.sh](https://brew.sh/). See [install](/getting-started/installing) if you are running Windows etc.
{{% /note %}}

```bash
brew install hugo
```

To verify your new install:

```bash
hugo version
```

{{< asciicast ItACREbFgvJ0HjnSNeTknxWy9 >}}

## Step 2: Create a New Site

```bash
hugo new site quickstart
```

The above will create a new Hugo site in a folder named `quickstart`.

{{< asciicast 3mf1JGaN0AX0Z7j5kLGl3hSh8 >}}

## Step 3: Add a Theme

See [themes.gohugo.io](https://themes.gohugo.io/) for a list of themes to consider. This quickstart uses the beautiful [Ananke theme](https://themes.gohugo.io/gohugo-theme-ananke/).

First, download the theme from Github and add it to your site's `theme` directory:

```bash
cd quickstart
git init
git submodule add https://github.com/budparr/gohugo-theme-ananke.git themes/ananke
```

*Note for non-git users:*
   - If you do not have git installed, you can download the archive of the latest
     version of this theme from:
       https://github.com/budparr/gohugo-theme-ananke/archive/master.zip
   - Extract that .zip file to get a "gohugo-theme-ananke-master" directory.
   - Rename that directory to "ananke", and move it into the "themes/" directory.

Then, add the theme to the site configuration:

```bash
echo 'theme = "ananke"' >> config.toml
```

{{< asciicast 7naKerRYUGVPj8kiDmdh5k5h9 >}}

## Step 4: Add Some Content

You can manually create content files (for example as `content/<CATEGORY>/<FILE>.<FORMAT>`) and provide metadata in them, however you can use the `new` command to do few things for you (like add title and date):

```
hugo new posts/my-first-post.md
```

{{< asciicast eUojYCfRTZvkEiqc52fUsJRBR >}}

Edit the newly created content file if you want, it will start with something like this:

```markdown
---
title: "My First Post"
date: 2019-03-26T08:47:11+01:00
draft: true
---

```

## Step 5: Start the Hugo server

Now, start the Hugo server with [drafts](/getting-started/usage/#draft-future-and-expired-content) enabled:

{{< asciicast BvJBsF6egk9c163bMsObhuNXj >}}

```
▶ hugo server -D

                   | EN
+------------------+----+
  Pages            | 10
  Paginator pages  |  0
  Non-page files   |  0
  Static files     |  3
  Processed images |  0
  Aliases          |  1
  Sitemaps         |  1
  Cleaned          |  0

Total in 11 ms
Watching for changes in /Users/bep/quickstart/{content,data,layouts,static,themes}
Watching for config changes in /Users/bep/quickstart/config.toml
Environment: "development"
Serving pages from memory
Running in Fast Render Mode. For full rebuilds on change: hugo server --disableFastRender
Web Server is available at http://localhost:1313/ (bind address 127.0.0.1)
Press Ctrl+C to stop
```

**Navigate to your new site at [http://localhost:1313/](http://localhost:1313/).**

Feel free to edit or add new content and simply refresh in browser to see changes quickly (You might need to force refresh in webbrowser, something like Ctrl-R usually works).

## Step 6: Customize the Theme

Your new site already looks great, but you will want to tweak it a little before you release it to the public.

### Site Configuration

Open up `config.toml` in a text editor:

```
baseURL = "https://example.org/"
languageCode = "en-us"
title = "My New Hugo Site"
theme = "ananke"
```

Replace the `title` above with something more personal. Also, if you already have a domain ready, set the `baseURL`. Note that this value is not needed when running the local development server.

{{% note %}}
**Tip:** Make the changes to the site configuration or any other file in your site while the Hugo server is running, and you will see the changes in the browser right away, though you may need to [clear your cache](https://kb.iu.edu/d/ahic).
{{% /note %}}

For theme specific configuration options, see the [theme site](https://github.com/budparr/gohugo-theme-ananke).

**For further theme customization, see [Customize a Theme](/themes/customizing/).**

### Step 7: Build static pages

It is simple. Just call:

```
hugo -D
```

Output will be in `./public/` directory by default (`-d`/`--destination` flag to change it, or set `publishdir` in the config file).

{{% note %}}
Drafts do not get deployed; once you finish a post, update the header of the post to say `draft: false`. More info [here](/getting-started/usage/#draft-future-and-expired-content).
{{% /note %}}
