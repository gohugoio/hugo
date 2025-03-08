---
title: Quick start
description: Create a Hugo site in minutes.
categories: []
keywords: []
params:
  minVersion: v0.128.0
weight: 10
aliases: [/quickstart/,/overview/quickstart/]
---

In this tutorial you will:

1. Create a site
1. Add content
1. Configure the site
1. Publish the site

## Prerequisites

Before you begin this tutorial you must:

1. [Install Hugo] (extended or extended/deploy edition, {{% param "minVersion" %}} or later)
1. [Install Git]

You must also be comfortable working from the command line.

## Create a site

### Commands

> [!note]
> **If you are a Windows user:**
>
> - Do not use the Command Prompt
> - Do not use Windows PowerShell
> - Run these commands from [PowerShell] or a Linux terminal such as WSL or Git > Bash
>
> PowerShell and Windows PowerShell [are different applications].

Verify that you have installed Hugo {{% param "minVersion" %}} or later.

```text
hugo version
```

Run these commands to create a Hugo site with the [Ananke] theme. The next section provides an explanation of each command.

```text
hugo new site quickstart
cd quickstart
git init
git submodule add https://github.com/theNewDynamic/gohugo-theme-ananke.git themes/ananke
echo "theme = 'ananke'" >> hugo.toml
hugo server
```

View your site at the URL displayed in your terminal. Press `Ctrl + C` to stop Hugo's development server.

### Explanation of commands

Create the [directory structure] for your project in the `quickstart` directory.

```text
hugo new site quickstart
```

Change the current directory to the root of your project.

```text
cd quickstart
```

Initialize an empty Git repository in the current directory.

```text
git init
```

Clone the [Ananke] theme into the `themes` directory, adding it to your project as a [Git submodule].

```text
git submodule add https://github.com/theNewDynamic/gohugo-theme-ananke.git themes/ananke
```

Append a line to the site configuration file, indicating the current theme.

```text
echo "theme = 'ananke'" >> hugo.toml
```

Start Hugo's development server to view the site.

```text
hugo server
```

Press `Ctrl + C` to stop Hugo's development server.

## Add content

Add a new page to your site.

```text
hugo new content content/posts/my-first-post.md
```

Hugo created the file in the `content/posts` directory. Open the file with your editor.

```text
+++
title = 'My First Post'
date = 2024-01-14T07:07:07+01:00
draft = true
+++
```

Notice the `draft` value in the [front matter] is `true`. By default, Hugo does not publish draft content when you build the site. Learn more about [draft, future, and expired content].

Add some [Markdown] to the body of the post, but do not change the `draft` value.

```text
+++
title = 'My First Post'
date = 2024-01-14T07:07:07+01:00
draft = true
+++
## Introduction

This is **bold** text, and this is *emphasized* text.

Visit the [Hugo](https://gohugo.io) website!
```

Save the file, then start Hugo's development server to view the site. You can run either of the following commands to include draft content.

```text
hugo server --buildDrafts
hugo server -D
```

View your site at the URL displayed in your terminal. Keep the development server running as you continue to add and change content.

When satisfied with your new content, set the front matter `draft` parameter to `false`.

> [!note]
> Hugo's rendering engine conforms to the CommonMark [specification] for Markdown. The CommonMark organization provides a useful [live testing tool] powered by the reference implementation.

## Configure the site

With your editor, open the [site configuration] file (`hugo.toml`) in the root of your project.

```text
baseURL = 'https://example.org/'
languageCode = 'en-us'
title = 'My New Hugo Site'
theme = 'ananke'
```

Make the following changes:

1. Set the `baseURL` for your production site. This value must begin with the protocol and end with a slash, as shown above.
1. Set the `languageCode` to your language and region.
1. Set the `title` for your production site.

Start Hugo's development server to see your changes, remembering to include draft content.

```text
hugo server -D
```

> [!note]
> Most theme authors provide configuration guidelines and options. Make sure to visit your theme's repository or documentation site for details.
>
> [The New Dynamic], authors of the Ananke theme, provide [documentation] for configuration and usage. They also provide a [demonstration site].

## Publish the site

In this step you will _publish_ your site, but you will not _deploy_ it.

When you _publish_ your site, Hugo creates the entire static site in the `public` directory in the root of your project. This includes the HTML files, and assets such as images, CSS files, and JavaScript files.

When you publish your site, you typically do _not_ want to include [draft, future, or expired content]. The command is simple.

```text
hugo
```

To learn how to _deploy_ your site, see the [host and deploy] section.

## Ask for help

Hugo's [forum] is an active community of users and developers who answer questions, share knowledge, and provide examples. A quick search of over 20,000 topics will often answer your question. Please be sure to read about [requesting help] before asking your first question.

## Other resources

For other resources to help you learn Hugo, including books and video tutorials, see the [external learning resources](/getting-started/external-learning-resources/) page.

[Ananke]: https://github.com/theNewDynamic/gohugo-theme-ananke
[are different applications]: https://learn.microsoft.com/en-us/powershell/scripting/whats-new/differences-from-windows-powershell?view=powershell-7.3
[demonstration site]: https://gohugo-ananke-theme-demo.netlify.app/
[directory structure]: /getting-started/directory-structure/
[documentation]: https://github.com/theNewDynamic/gohugo-theme-ananke#readme
[draft, future, and expired content]: /getting-started/usage/#draft-future-and-expired-content
[draft, future, or expired content]: /getting-started/usage/#draft-future-and-expired-content
[forum]: https://discourse.gohugo.io/
[front matter]: /content-management/front-matter/
[Git submodule]: https://git-scm.com/book/en/v2/Git-Tools-Submodules
[host and deploy]: /host-and-deploy/
[Install Git]: https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
[Install Hugo]: /installation/
[live testing tool]: https://spec.commonmark.org/dingus/
[Markdown]: https://daringfireball.net/projects/markdown
[PowerShell]: https://learn.microsoft.com/en-us/powershell/scripting/install/installing-powershell-on-windows
[requesting help]: https://discourse.gohugo.io/t/requesting-help/9132
[site configuration]: /configuration/
[specification]: https://spec.commonmark.org/
[The New Dynamic]: https://www.thenewdynamic.com/
