---
title: Quick start
description: Create your first Hugo project.
categories: []
keywords: []
params:
  minVersion: v0.156.0
weight: 10
aliases: [/quickstart/,/overview/quickstart/]
---

In this tutorial you will:

1. Create a project
1. Add content
1. Configure the project
1. Publish the project

## Prerequisites

Before you begin this tutorial you must:

1. [Install Hugo] (extended or extended/deploy edition, {{% param "minVersion" %}} or later)
1. [Install Git]

You must also be comfortable working from the command line.

## Create a project

### Commands

> [!note]
> **If you are a Windows user:**
>
> - Do not use the Command Prompt
> - Do not use Windows PowerShell
> - Run these commands from [PowerShell][] or a Linux terminal such as WSL or Git > Bash
>
> PowerShell and Windows PowerShell [are different applications][].

Verify that you have installed Hugo {{% param "minVersion" %}} or later.

```text
hugo version
```

Run these commands to create a Hugo project with the [Ananke][] theme. The next section provides an explanation of each command.

```text
hugo new project quickstart
cd quickstart
git init
git submodule add https://github.com/theNewDynamic/gohugo-theme-ananke.git themes/ananke
echo "theme = 'ananke'" >> hugo.toml
hugo server
```

View your project at the URL displayed in your terminal. Press `Ctrl + C` to stop Hugo's development server.

### Explanation of commands

Create the [project skeleton][] for your project in the `quickstart` directory.

```text
hugo new project quickstart
```

Change the current directory to the root of your project.

```text
cd quickstart
```

Initialize an empty Git repository in the current directory.

```text
git init
```

Clone the [Ananke][] theme into the `themes` directory, adding it to your project as a [Git submodule][].

```text
git submodule add https://github.com/theNewDynamic/gohugo-theme-ananke.git themes/ananke
```

Append a line to your project configuration file, indicating the current theme.

```text
echo "theme = 'ananke'" >> hugo.toml
```

Start Hugo's development server.

```text
hugo server
```

Press `Ctrl + C` to stop Hugo's development server.

## Add content

Add a new page to your project.

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

Notice the `draft` value in the [front matter][] is `true`. By default, Hugo does not publish draft content when you build the project. Learn more about [draft, future, and expired content][].

Add some [Markdown][] to the body of the post, but do not change the `draft` value.

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

Save the file, then start Hugo's development server. You can run either of the following commands to include draft content.

```text
hugo server --buildDrafts
hugo server -D
```

View your project at the URL displayed in your terminal. Keep the development server running as you continue to add and change content.

When satisfied with your new content, set the front matter `draft` parameter to `false`.

> [!note]
> Hugo's rendering engine conforms to the CommonMark [specification][] for Markdown. The CommonMark organization provides a useful [live testing tool][] powered by the reference implementation.

## Configure the project

With your editor, open your [project configuration][] file (`hugo.toml`) in the root of your project.

```text
baseURL = 'https://example.org/'
languageCode = 'en-us'
title = 'My New Hugo Project'
theme = 'ananke'
```

Make the following changes:

1. Set the `baseURL` for your project. This value must begin with the protocol and end with a slash, as shown above.
1. Set the `languageCode` to your locale.
1. Set the `title` for your project.

Start Hugo's development server to see your changes, remembering to include draft content.

```text
hugo server -D
```

> [!note]
> Most theme authors provide configuration guidelines and options. Make sure to visit your theme's repository or documentation site for details.
>
> [The New Dynamic][], authors of the Ananke theme, provide [documentation][] for configuration and usage. They also provide a [demonstration site][].

## Publish the project

In this step you will _publish_ your project, but you will not _deploy_ it.

When you publish your project, Hugo renders all build artifacts to the `public` directory in the root of your project. This includes the HTML files for every site, along with assets such as images, CSS, and JavaScript. The command is simple.

```text
hugo
```

To learn how to _deploy_ your project, see the [host and deploy][] section.

## Ask for help

Hugo's [forum][] is an active community of users and developers who answer questions, share knowledge, and provide examples. A quick search of over 20,000 topics will often answer your question. Please be sure to read about [requesting help] before asking your first question.

## Other resources

For other resources to help you learn Hugo, including books and video tutorials, see the [external learning resources][] page.

[Ananke]: https://github.com/theNewDynamic/gohugo-theme-ananke
[are different applications]: https://learn.microsoft.com/en-us/powershell/scripting/whats-new/differences-from-windows-powershell?view=powershell-7.3
[demonstration site]: https://gohugo-ananke-theme-demo.netlify.app/
[documentation]: https://github.com/theNewDynamic/gohugo-theme-ananke#readme
[draft, future, and expired content]: /getting-started/usage/#draft-future-and-expired-content
[external learning resources]: /getting-started/external-learning-resources/
[forum]: https://discourse.gohugo.io/
[front matter]: /content-management/front-matter/
[Git submodule]: https://git-scm.com/book/en/v2/Git-Tools-Submodules
[host and deploy]: /host-and-deploy/
[Install Git]: https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
[Install Hugo]: /installation/
[live testing tool]: https://spec.commonmark.org/dingus/
[Markdown]: https://daringfireball.net/projects/markdown
[PowerShell]: https://learn.microsoft.com/en-us/powershell/scripting/install/installing-powershell-on-windows
[project configuration]: /configuration/
[project skeleton]: /getting-started/directory-structure/#project-skeleton
[requesting help]: https://discourse.gohugo.io/t/requesting-help/9132
[specification]: https://spec.commonmark.org/
[The New Dynamic]: https://www.thenewdynamic.com/
