---
aliases:
- /doc/usage/
lastmod: 2016-08-19
date: 2013-07-01
menu:
  main:
    parent: getting started
next: /overview/configuration
notoc: true
prev: /overview/installing
title: Using Hugo
weight: 30
---

Make sure Hugo is in your `PATH` (or provide a path to it). Test this by:

{{< nohighlight >}}$ hugo help
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
  help        Help about any command
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
      --cleanDestinationDir        remove files from destination not found in static directories
      --config string              config file (default is path/config.yaml|json|toml)
  -c, --contentDir string          filesystem path to content directory
  -d, --destination string         filesystem path to write files to
      --disable404                 do not render 404 page
      --disableKinds stringSlice   disable different kind of pages (home, RSS etc.)
      --disableRSS                 do not build RSS files
      --disableSitemap             do not build Sitemap file
      --enableGitInfo              add Git revision, date and author info to the pages
      --forceSyncStatic            copy all files when static is changed.
  -h, --help                       help for hugo
      --i18n-warnings              print missing translations
      --ignoreCache                ignores the cache directory
  -l, --layoutDir string           filesystem path to layout directory
      --log                        enable Logging
      --logFile string             log File path (if set, logging enabled automatically)
      --noChmod                    don't sync permission mode of files
      --noTimes                    don't sync modification time of files
      --pluralizeListTitles        pluralize titles in lists using inflect (default true)
      --preserveTaxonomyNames      preserve taxonomy names as written ("Gérard Depardieu" vs "gerard-depardieu")
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

Use "hugo [command] --help" for more information about a command.
{{< /nohighlight >}}

## Common Usage Example

The most common use is probably to run `hugo` with your current directory being the input directory:

{{< nohighlight >}}$ hugo
0 draft content
0 future content
99 pages created
0 paginator pages created
16 tags created
0 groups created
in 120 ms
{{< /nohighlight >}}

This generates your web site to the `public/` directory,
ready to be deployed to your web server.


## Instant feedback as you develop your web site

If you are working on things and want to see the changes immediately, by default
Hugo will watch the filesystem for changes, and rebuild your site as soon as a file is saved:

{{< nohighlight >}}$ hugo -s ~/Code/hugo/docs
0 draft content
0 future content
99 pages created
0 paginator pages created
16 tags created
0 groups created
in 120 ms
Watching for changes in /Users/spf13/Code/hugo/docs/content
Press Ctrl+C to stop
{{< /nohighlight >}}

Hugo can even run a server and create a site preview at the same time!
Hugo implements [LiveReload](/extras/livereload/) technology to automatically
reload any open pages in all JavaScript-enabled browsers, including mobile.
This is the easiest and most common way to develop a Hugo web site:

{{< nohighlight >}}$ hugo server -ws ~/Code/hugo/docs
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
{{< /nohighlight >}}


## Deploying your web site

After running `hugo server` for local web development,
you need to do a final `hugo` run
**without the `server` part of the command**
to rebuild your site.
You may then **deploy your site** by copying the `public/` directory
(by FTP, SFTP, WebDAV, Rsync, `git push`, etc.)
to your production web server.

Since Hugo generates a static website, your site can be hosted anywhere,
including [Heroku][], [GoDaddy][], [DreamHost][], [GitHub Pages][],
[Amazon S3][] with [CloudFront][], [Firebase Hosting][],
or any other cheap (or even free) static web hosting service.

[Apache][], [nginx][], [IIS][]...  Any web server software would do!

[Apache]: http://httpd.apache.org/ "Apache HTTP Server"
[nginx]: http://nginx.org/
[IIS]: http://www.iis.net/
[Heroku]: https://www.heroku.com/
[GoDaddy]: https://www.godaddy.com/
[DreamHost]: http://www.dreamhost.com/
[GitHub Pages]: https://pages.github.com/
[GitLab]: https://about.gitlab.com
[Amazon S3]: http://aws.amazon.com/s3/
[CloudFront]: http://aws.amazon.com/cloudfront/ "Amazon CloudFront"
[Firebase Hosting]: https://firebase.google.com/docs/hosting/

### A note about deployment

Running `hugo` *does not* remove generated files before building. This means that you should delete your `public/` directory (or the directory you specified with `-d`/`--destination`) before running the `hugo` command, or you run the risk of the wrong files (e.g. drafts and/or future posts) being left in the generated site.

An easy way to work around this is to use different directories for development and production.

To start a server that builds draft content (helpful for editing), you can specify a different destination: the `dev/` dir.

{{< nohighlight >}}$ hugo server -wDs ~/Code/hugo/docs -d dev
{{< /nohighlight >}}

When the content is ready for publishing, use the default `public/` dir:

{{< nohighlight >}}$ hugo -s ~/Code/hugo/docs
{{< /nohighlight >}}

This prevents content you're not yet ready to share
from accidentally becoming available.

### Alternatively, serve your web site with Hugo!

Yes, that's right!  Because Hugo is so blazingly fast both in web site creation
*and* in web serving (thanks to its concurrent and multi-threaded design and
its Go heritage), some users actually prefer using Hugo itself to serve their
web site *on their production server*!

No other web server software (Apache, nginx, IIS...) is necessary.

Here is the command:

{{< nohighlight >}}$ hugo server --baseURL=http://yoursite.org/ \
              --port=80 \
              --appendPort=false \
              --bind=87.245.198.50
{{< /nohighlight >}}

Note the `bind` option,
which is the interface to which the server will bind
(defaults to `127.0.0.1`:
fine for most development use cases).
Some hosts, such as Amazon Web Services,
run NAT (network address translation);
sometimes it can be hard to figure out the actual IP address.
Using `--bind=0.0.0.0` will bind to all interfaces.

This way, you may actually deploy just the source files,
and Hugo on your server will generate the resulting web site
on-the-fly and serve them at the same time.

You may optionally add `--disableLiveReload=true` if you do not want
the JavaScript code for LiveReload to be added to your web pages.

Interested? Here are some great tutorials contributed by Hugo users:

* [hugo, syncthing](http://fredix.xyz/2014/10/hugo-syncthing/) (French) by Frédéric Logier (@fredix)
