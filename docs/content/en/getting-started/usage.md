---
title: Basic usage
description: Hugo's command line interface (CLI) is fully featured but simple to use, even for those with limited experience working from the command line.
categories: [getting started]
keywords: [usage,livereload,command,flags]
menu:
  docs:
    parent: getting-started
    weight: 30
weight: 30
toc: true
aliases: [/overview/usage/,/extras/livereload/,/doc/usage/,/usage/]
---

## Test your installation

After [installing] Hugo, test your installation by running:

```sh
hugo version
```

You should see something like:

```text
hugo v0.123.0-3c8a4713908e48e6523f058ca126710397aa4ed5+extended linux/amd64 BuildDate=2024-02-19T16:32:38Z VendorInfo=gohugoio
```

## Display available commands

To see a list of the available commands and flags:

```sh
hugo help
```

To get help with a subcommand, use the `--help` flag. For example:

```sh
hugo server --help
```

## Build your site

To build your site, `cd` into your project directory and run:

```sh
hugo
```

The [`hugo`] command builds your site, publishing the files to the `public` directory. To publish your site to a different directory, use the [`--destination`] flag or set [`publishDir`] in your site configuration.

{{% note %}}
Hugo does not clear the `public` directory before building your site. Existing files are overwritten, but not deleted. This behavior is intentional to prevent the inadvertent removal of files that you may have added to the `public` directory after the build.

Depending on your needs, you may wish to manually clear the contents of the public directory before every build.
{{% /note %}}

## Draft, future, and expired content

Hugo allows you to set `draft`, `date`, `publishDate`, and `expiryDate` in the [front matter] of your content. By default, Hugo will not publish content when:

- The `draft` value is `true`
- The `date` is in the future
- The `publishDate` is in the future
- The `expiryDate` is in the past

{{< new-in 0.123.0 >}}

{{% note %}}
Hugo publishes descendants of draft, future, and expired [node] pages. To prevent publication of these descendants, use the [`cascade`] front matter field to cascade [build options] to the descendant pages.

[build options]: /content-management/build-options/
[`cascade`]: /content-management/front-matter/#cascade-field
[node]: /getting-started/glossary/#node
{{% /note %}}

You can override the default behavior when running `hugo` or `hugo server` with command line flags:

```sh
hugo --buildDrafts    # or -D
hugo --buildExpired   # or -E
hugo --buildFuture    # or -F
```

Although you can also set these values in your site configuration, it can lead to unwanted results unless all content authors are aware of, and understand, the settings.

{{% note %}}
As noted above, Hugo does not clear the `public` directory before building your site. Depending on the _current_ evaluation of the four conditions above, after the build your `public` directory may contain extraneous files from a previous build.

A common practice is to manually clear the contents of the `public` directory before each build to remove draft, expired, and future content.
{{% /note %}}

## Develop and test your site

To view your site while developing layouts or creating content, `cd` into your project directory and run:

```sh
hugo server
```

The [`hugo server`] command builds your site into memory, and serves your pages using a minimal HTTP server. When you run `hugo server` it will display the URL of your local site:

```text
Web Server is available at http://localhost:1313/ 
```

While the server is running, it watches your project directory for changes to assets, configuration, content, data, layouts, translations, and static files. When it detects a change, the server rebuilds your site and refreshes your browser using [LiveReload].

Most Hugo builds are so fast that you may not notice the change unless you are looking directly at your browser.

### LiveReload

While the server is running, Hugo injects JavaScript into the generated HTML pages. The LiveReload script creates a connection from the browser to the server via web sockets. You do not need to install any software or browser plugins, nor is any configuration required.

### Automatic redirection

When editing content, if you want your browser to automatically redirect to the page you last modified, run:

```sh
hugo server --navigateToChanged
```

## Deploy your site

{{% note %}}
As noted above, Hugo does not clear the public directory before building your site. Manually clear the contents of the public directory before each build to remove draft, expired, and future content.
{{% /note %}}

When you are ready to deploy your site, run:

```sh
hugo
```

This builds your site, publishing the files to the public directory. The directory structure will look something like this:

```text
public/
├── categories/
│   ├── index.html
│   └── index.xml  <-- RSS feed for this section
├── posts/
│   ├── my-first-post/
│   │   └── index.html
│   ├── index.html
│   └── index.xml  <-- RSS feed for this section
├── tags/
│   ├── index.html
│   └── index.xml  <-- RSS feed for this section
├── index.html
├── index.xml      <-- RSS feed for the site
└── sitemap.xml
```

In a simple hosting environment, where you typically `ftp`, `rsync`, or `scp` your files to the root of a virtual host, the contents of the `public` directory are all that you need.

Most of our users deploy their sites using a CI/CD workflow, where a push[^1] to their GitHub or GitLab repository triggers a build and deployment. Popular providers include [AWS Amplify], [CloudCannon], [Cloudflare Pages], [GitHub Pages], [GitLab Pages], and [Netlify].

Learn more in the [hosting and deployment] section.

[^1]: The Git repository contains the entire project directory, typically excluding the public directory because the site is built _after_ the push.

[`--destination`]: /commands/hugo/#options
[`hugo server`]: /commands/hugo_server/
[`hugo`]: /commands/hugo/
[`publishDir`]: /getting-started/configuration/#publishdir
[AWS Amplify]: https://aws.amazon.com/amplify/
[CloudCannon]: https://cloudcannon.com/
[Cloudflare Pages]: https://pages.cloudflare.com/
[commands]: /commands/
[front matter]: /content-management/front-matter/
[GitHub Pages]: https://pages.github.com/
[GitLab Pages]: https://docs.gitlab.com/ee/user/project/pages/
[hosting and deployment]: /hosting-and-deployment/
[hosting]: /hosting-and-deployment/
[installing]: /installation/
[LiveReload]: https://github.com/livereload/livereload-js
[Netlify]: https://www.netlify.com/
