---
title: Basic Usage
linktitle: Basic Usage
description: Hugo's CLI is fully featured but simple to use, even for those who have very limited experience working from the command line.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [getting started]
tags: [usage,livereload,command line,flags]
menu:
  main:
    parent: "Getting Started"
    weight: 40
weight: 40
sections_weight: 40
draft: false
aliases: [/overview/usage/,/extras/livereload/,/doc/usage/,/usage/]
toc: true
---

The following is a description of the most command commands you will use while developing your Hugo project. See the [Command Line Reference][commands] for a comprehensive view of Hugo's CLI.

## Testing Installation

Once you have [installed Hugo][install], make sure it is in your `PATH`. You can test that Hugo has been installed correctly via the `help` command:

```bash
hugo help
```

The output you see in your console should be similar to the following:

```bash
hugo is the main command, used to build your Hugo site.

Hugo is a Fast and Flexible Static Site Generator
built with love by spf13 and friends in Go.

Complete documentation is available at http://gohugo.io/.

Usage:
  hugo [flags]
  hugo [command]

Available Commands:
  benchmark   Benchmark Hugo by building a site a number of times.
  check       Contains some verification checks
  config      Print the site configuration
  convert     Convert your content to different formats
  env         Print Hugo version and environment info
  gen         A collection of several useful generators.
  import      Import your site from others.
  list        Listing out various types of content
  new         Create new content for your site
  server      A high performance webserver
  undraft     Undraft changes the content's draft status from 'True' to 'False'
  version     Print the version number of Hugo

Flags:
  -b, --baseURL string             hostname (and path) to the root, e.g. http://spf13.com/
  -D, --buildDrafts                include content marked as draft
  -E, --buildExpired               include expired content
  -F, --buildFuture                include content with publishdate in the future
      --cacheDir string            filesystem path to cache directory. Defaults: $TMPDIR/hugo_cache/
      --canonifyURLs               if true, all relative URLs will be canonicalized using baseURL
      --cleanDestinationDir        Remove files from destination not found in static directories
      --config string              config file (default is path/config.yaml|json|toml)
  -c, --contentDir string          filesystem path to content directory
  -d, --destination string         filesystem path to write files to
      --disable404                 Do not render 404 page
      --disableKinds stringSlice   Disable different kind of pages (home, RSS etc.)
      --disableRSS                 Do not build RSS files
      --disableSitemap             Do not build Sitemap file
      --enableGitInfo              Add Git revision, date and author info to the pages
      --forceSyncStatic            Copy all files when static is changed.
      --i18n-warnings              Print missing translations
      --ignoreCache                Ignores the cache directory
  -l, --layoutDir string           filesystem path to layout directory
      --log                        Enable Logging
      --logFile string             Log File path (if set, logging enabled automatically)
      --noChmod                    Don't sync permission mode of files
      --noTimes                    Don't sync modification time of files
      --pluralizeListTitles        Pluralize titles in lists using inflect (default true)
      --preserveTaxonomyNames      Preserve taxonomy names as written ("Gérard Depardieu" vs "gerard-depardieu")
      --quiet                      build in quiet mode
      --renderToMemory             render to memory (only useful for benchmark testing)
  -s, --source string              filesystem path to read files relative from
      --stepAnalysis               display memory and timing of different steps of the program
  -t, --theme string               theme to use (located in /themes/THEMENAME/)
      --themesDir string           filesystem path to themes directory
      --uglyURLs                   if true, use /filename.html instead of /filename/
  -v, --verbose                    verbose output
      --verboseLog                 verbose logging
  -w, --watch                      watch filesystem for changes and recreate as needed
```

## The `hugo` Command

The most common usage is probably to run `hugo` with your current directory being the input directory.

This generates your website to the `public/` directory by default, although you can customize the output directory in your [site configuration][config] by changing the `publishDir` field.

The site Hugo renders into `public/` is ready to be deployed to your web server:

```bash
hugo
0 draft content
0 future content
99 pages created
0 paginator pages created
16 tags created
0 groups created
in 90 ms
```

## Draft, Future, and Expired Content

Hugo allows you to set `draft`, `publishdate`, and even `expirydate` in your content's [front matter][]. By default, Hugo will not publish:

1. Content with a future `publishdate` value
2. Content with `draft: true` status
3. Content with a past `expirydate` value

All three of these can be overridden during both local development *and* deployment by adding the following flags to `hugo` and `hugo server`, respectively, or by changing the boolean values assigned to the fields of the same name (without `--`) in your [configuration][config]:

1. `--buildFuture`
2. `--buildDrafts`
3. `--buildExpired`

## LiveReload

Hugo comes with [LiveReload](https://github.com/livereload/livereload-js) built in. There are no additional packages to install. A common way to use Hugo while developing a site is to have Hugo run a server with the `hugo server` command and watch for changes:

```bash
hugo server
0 draft content
0 future content
99 pages created
0 paginator pages created
16 tags created
0 groups created
in 120 ms
Watching for changes in /Users/yourname/sites/yourhugosite/{data,content,layouts,static}
Serving pages from /Users/yourname/sites/yourhugosite/public
Web Server is available at http://localhost:1313/
Press Ctrl+C to stop
```

This will run a fully functioning web server while simultaneously watching your file system for additions, deletions, or changes within the following areas of your [project organization][dirs]:

* `/static/*`
* `/content/*`
* `/data/*`
* `/layouts/*`
* `/themes/<CURRENT-THEME>/*`
* `config`

Whenever you make changes, Hugo will simultaneously rebuild the site and continue to serve content. As soon as the build is finished, LiveReload tells the browser to silently reload the page.

Most Hugo builds are so fast that you may not notice the change unless looking directly at the site in your browser. This means that keeping the site open on a second monitor (or another half of your current monitor) allows you to see the most up-to-date version of your website without the need to leave your text editor.

{{% note "Closing `</body>` Tag"%}}
Hugo injects the LiveReload `<script>` before the closing `</body>` in your templates and will therefore not work if this tag is not present..
{{% /note %}}

### Disabling LiveReload

LiveReload works by injecting JavaScript into the pages Hugo generates. The script creates a connection from the browser's web socket client to the Hugo web socket server.

LiveReload is awesome for development. However, some Hugo users may use `hugo server` in production to instantly display updated content. The following methods make it easy to disable LiveReload:

```bash
hugo server --watch=false
```

Or...

```bash
hugo server --disableLiveReload
```

The latter flag can be omitted by adding the following key-value to  your `config.toml` or `config.yml` file, respectively:

```toml
disableLiveReload = true
```

```yaml
disableLiveReload: true
```

## Deploying Your Website

After running `hugo server` for local web development, you need to do a final `hugo` run *without the `server` part of the command* to rebuild your site. You may then deploy your site by copying the `public/` directory to your production web server.

Since Hugo generates a static website, your site can be hosted *anywhere* using any web server. See [Hosting and Deployments][] for methods for hosting and automating deployments contributed by the Hugo community.

{{% warning "Generated Files are **NOT** Removed on Site Build" %}}
Running `hugo` *does not* remove generated files before building. This means that you should delete your `public/` directory (or the publish directory you specified via flag or configuration file) before running the `hugo` command. If you do not remove these files, you run the risk of the wrong files (e.g., drafts or future posts) being left in the generated site.
{{% /warning %}}

### Dev vs Deploy Destinations

Hugo does not remove generated files before building. An easy workaround is to use different directories for development and production.

To start a server that builds draft content (helpful for editing), you can specify a different destination; e.g., a `dev/` directory:

```bash
hugo server -wDs ~/Code/hugo/docs -d dev
```

When the content is ready for publishing, use the default `public/` dir:

```bash
hugo -s ~/Code/hugo/docs
```

This prevents draft content from accidentally becoming available.

### Using Hugo's Server in Production

Because Hugo is so blazingly fast both in website creation *and* in web serving (thanks to its concurrent, multi-threaded design and Golang heritage), some users prefer to use Hugo itself to serve their website *on their production server*.

No other web server software (e.g., Apache, nginx, IIS) is necessary.

Here is the command:

```bash
hugo server --baseURL=http://yoursite.org/ \
--port=80 \
--appendPort=false \
--bind=87.245.198.50
```

Note the `bind` option, which is the interface to which the server will bind (defaults to `127.0.0.1`: fine for most development use cases). Some hosts, such as Amazon Web Services, run NAT (network address translation); sometimes it can be hard to figure out the actual IP address. Using `--bind=0.0.0.0` will bind to all interfaces.

By using Hugo's server in production, you are able to deploy just the source files. Hugo, running on your server, will generate the resulting website on the fly and serve them at the same time.


[commands]: /commands/
[config]: /getting-started/configuration/
[dirs]: /getting-started/directory-structure/
[front matter]: /content-management/front-matter/
[hosting]: /hosting-and-deployments/
[install]: /getting-started/installing/