---
title: Basic Usage
linktitle: Basic Usage
description: Hugo's CLI is fully featured but simple. You do not need a high level of expertise on the command line to get up and running.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [getting started]
tags: [usage,livereload,command line]
weight: 40
draft: false
aliases: [/overview/usage/,/extras/livereload/,/getting-started/using-hugo/,/doc/usage/]
toc: true
needsreview: true
notesforauthors:
---

## Testing Installation with Hugo Help

Once you have [installed Hugo][install], make sure it is in your `PATH`. You can test that Hugo has been installed correctly via the `help` command:

```bash
hugo help
```

The output you see in your console should be similar to the following:

```bash
hugo is the main command, used to build your Hugo site.

Hugo is a Fast and Flexible Static Site Generator built with love by spf13 and friends in Go.

Complete documentation is available at http://gohugo.io/.

Usage:
  hugo [flags]
  hugo [command]

Available Commands:
  benchmark   Benchmark Hugo by building a site a number of times.
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
  -b, --baseURL string          hostname (and path) to the root, e.g. http://spf13.com/
  -D, --buildDrafts             include content marked as draft
  -E, --buildExpired            include expired content
  -F, --buildFuture             include content with publishdate in the future
      --cacheDir string         filesystem path to cache directory. Defaults: $TMPDIR/hugo_cache/
      --canonifyURLs            if true, all relative URLs will be canonicalized using baseURL
      --cleanDestinationDir     Remove files from destination not found in static directories
      --config string           config file (default is path/config.yaml|json|toml)
  -c, --contentDir string       filesystem path to content directory
  -d, --destination string      filesystem path to write files to
      --disable404              Do not render 404 page
      --disableRSS              Do not build RSS files
      --disableSitemap          Do not build Sitemap file
      --enableGitInfo           Add Git revision, date and author info to the pages
      --forceSyncStatic         Copy all files when static is changed.
      --i18n-warnings           Print missing translations
      --ignoreCache             Ignores the cache directory
  -l, --layoutDir string        filesystem path to layout directory
      --log                     Enable Logging
      --logFile string          Log File path (if set, logging enabled automatically)
      --noChmod                 Don't sync permission mode of files
      --noTimes                 Don't sync modification time of files
      --pluralizeListTitles     Pluralize titles in lists using inflect (default true)
      --preserveTaxonomyNames   Preserve taxonomy names as written ("Gérard Depardieu" vs "gerard-depardieu")
      --quiet                   build in quiet mode
      --renderToMemory          render to memory (only useful for benchmark testing)
  -s, --source string           filesystem path to read files relative from
      --stepAnalysis            display memory and timing of different steps of the program
  -t, --theme string            theme to use (located in /themes/THEMENAME/)
      --uglyURLs                if true, use /filename.html instead of /filename/
  -v, --verbose                 verbose output
      --verboseLog              verbose logging
  -w, --watch                   watch filesystem for changes and recreate as needed

Additional help topics:
  hugo check     Contains some verification checks

Use "hugo [command] --help" for more information about a command.
```


The most common use is probably to run `hugo` with your current directory being the input directory:

```bash
hugo
0 draft content
0 future content
99 pages created
0 paginator pages created
16 tags created
0 groups created
in 120 ms
```

This generates your website to the `public/` directory, ready to be deployed to your web server.

## Instant feedback as you develop your web site

If you are working on things and want to see the changes immediately, by default Hugo will watch the filesystem for changes and rebuild your site as soon as a file is saved:

```bash
hugo -s ~/Code/hugo/docs
0 draft content
0 future content
99 pages created
0 paginator pages created
16 tags created
0 groups created
in 120 ms
Watching for changes in /Users/spf13/Code/hugo/docs/content
Press Ctrl+C to stop
```

Hugo can even run a server and create a site preview at the same time! Hugo implements [LiveReload](#LiveReload) technology to automatically reload any open pages in all JavaScript-enabled browsers, including mobile. This is the easiest and most common way to develop a Hugo web site:

```bash
hugo server -ws ~/Code/hugo/docs
0 draft content
0 future content
99 pages created
0 paginator pages created
16 tags created
0 groups created
in 120 ms
Watching for changes in /Users/spf13/Code/hugo/docs/content
Serving pages from /Users/spf13/Code/hugo/docs/public
Web Server is available at http://localhost:1313/
Press Ctrl+C to stop
```
Hugo may not be the first static site generator to utilize LiveReload
technology, but it’s the first to do it right.

The combination of Hugo’s insane build speed and LiveReload make
crafting your content pure joy. Your updated content appears virtually instantly in your browser as soon as you save your changes.

### LiveReload

Hugo comes with [LiveReload](https://github.com/livereload/livereload-js) built in. There are no additional packages to install. A common way to use Hugo while developing a site is to have Hugo run a server and watch for changes:

```bash
hugo server
```

This will run a fully functioning web server while simultaneously watching your file system for additions, deletions, or changes within the following the following areas of your [project organization][directorystructure]

* `/static/*`
* `/content/*`
* `/data/*`
* `/layouts/*`
* `/themes/<only your current theme>/*`
* `config`

Whenever you make changes, Hugo will simultaneously rebuild the site and continue to serve content. As soon as the build is finished, LiveReload tells the browser to silently reload the page.

Most Hugo builds are so fast that you may not notice the change unless looking directly at the site in your browser. This means that keeping the site open on a second monitor (or another half of your current monitor) allows you to see the most up-to-date version of your website without the need to leave your text editor.

{{% note "Closing `</body>` Tag"%}}
You must have a closing `</body>` tag for LiveReload to work.
Hugo injects the LiveReload `<script>` before this tag.
{{% /note %}}

### Disabling LiveReload's Watch Feature

LiveReload works by injecting JavaScript into the pages Hugo generates. The script creates a connection from the browser's web socket client to the
Hugo web socket server.

LiveReload is awesome for development. However, some Hugo users may use `hugo server` in production to instantly display updated content. As such, we’ve provided multiple methods  made it easy to disable the LiveReload functionality:

```bash
hugo server --watch=false
```
Or...

```bash
hugo server --disableLiveReload
```
The latter flag can be omitted by adding the following key-value to  your `config` file:

```toml
disableLiveReload = true
```

```yaml
disableLiveReload: true
```

## Deploying Your Website

After running `hugo server` for local web development, you need to do a final `hugo` run **without the `server` part of the command** to rebuild your site. You may then **deploy your site** by copying the `public/` directory (by FTP, SFTP, WebDAV, Rsync, `git push`, etc.) to your production web server.

Since Hugo generates a static website, your site can be hosted anywhere, including [Heroku][], [GoDaddy][], [DreamHost][], [GitHub Pages][], [Amazon S3][] with [CloudFront][], [Firebase Hosting][], or any other cheap (or even free) static web hosting service.

[Apache][], [nginx][], [IIS][]...any web server software will work.

{{% warning "Generated Files are **NOT** Removed on Site Build" %}}
Running `hugo` *does not* remove generated files before building. This means that you should delete your `public/` directory (or the directory you specified with `-d`/`--destination`) before running the `hugo` command, or you run the risk of the wrong files (e.g., drafts or future posts) being left in the generated site.
{{% /warning %}}

### Destination Directories for Dev vs Deploy

Hugo does not remove generated files before building. An easy workaround is to use different directories for development and production.

To start a server that builds draft content (helpful for editing), you can specify a different destination; e.g., a `dev/` directory:

```bash
hugo server -wDs ~/Code/hugo/docs -d dev
```

When the content is ready for publishing, use the default `public/` dir:

```bash
hugo -s ~/Code/hugo/docs
```

This prevents content you're not yet ready to share from accidentally becoming available.

### Using Hugo's Server in Production

Because Hugo is so blazingly fast both in website creation *and* in web serving (thanks to its concurrent, multi-threaded design and Golang heritage), some users prefer to use Hugo itself to serve their website *on their production server*.

No other web server software (e.g., Apache, nginx, IIS) is necessary.

Here is the command:

{{% input file="hugo-production-server.sh" %}}
```bash
hugo server --baseURL=http://yoursite.org/ \
--port=80 \
--appendPort=false \
--bind=87.245.198.50
```
{{% /input %}}

Note the `bind` option, which is the interface to which the server will bind (defaults to `127.0.0.1`: fine for most development use cases). Some hosts, such as Amazon Web Services, run NAT (network address translation); sometimes it can be hard to figure out the actual IP address. Using `--bind=0.0.0.0` will bind to all interfaces.

By using Hugo's server in production, you are able to deploy just the source files. Hugo, running on your server, will generate the resulting website on the fly and serve them at the same time.

Interested? Here are some great tutorials contributed by Hugo users:

* [hugo, syncthing](http://fredix.xyz/2014/10/hugo-syncthing/) (French) by Frédéric Logier (@fredix)

[Amazon S3]: http://aws.amazon.com/s3/
[Apache]: http://httpd.apache.org/ "Apache HTTP Server"
[CloudFront]: http://aws.amazon.com/cloudfront/ "Amazon CloudFront"
[directorystructure]: /getting-started/directory-structure/
[DreamHost]: http://www.dreamhost.com/
[Firebase Hosting]: https://firebase.google.com/docs/hosting/
[GitHub Pages]: https://pages.github.com/
[GitLab]: https://about.gitlab.com
[GoDaddy]: https://www.godaddy.com/
[Heroku]: https://www.heroku.com/
[IIS]: http://www.iis.net/
[install]: /getting-started/install-hugo/
[nginx]: http://nginx.org/