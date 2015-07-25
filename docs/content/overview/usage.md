---
aliases:
- /doc/usage/
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

Make sure either `hugo` is in your `PATH` or provide a path to it.

<pre><code class="hljs nohighlight">$ hugo help
	
Hugo is a Fast and Flexible Static Site Generator built with love by spf13 and friends in Go.

Complete documentation is available at http://gohugo.io

Usage: 
  hugo [flags]
  hugo [command]

Available Commands: 
  server          Hugo runs its own webserver to render the files
  version         Print the version number of Hugo
  config          Print the site configuration
  check           Check content in the source directory
  benchmark       Benchmark hugo by building a site a number of times
  new             Create new content for your site
  undraft         Undraft changes the content's draft status from 'True' to 'False'
  genautocomplete Generate shell autocompletion script for Hugo
  gendoc          Generate Markdown documentation for the Hugo CLI.
  help            Help about any command

Flags:
  -b, --baseUrl="": hostname (and path) to the root eg. http://spf13.com/
  -D, --buildDrafts=false: include content marked as draft
  -F, --buildFuture=false: include content with publishdate in the future
      --cacheDir="": filesystem path to cache directory. Defaults: $TMPDIR/hugo_cache/
      --config="": config file (default is path/config.yaml|json|toml)
  -d, --destination="": filesystem path to write files to
      --disableRSS=false: Do not build RSS files
      --disableSitemap=false: Do not build Sitemap file
      --editor="": edit new content with this editor, if provided
  -h, --help=false: help for hugo
      --ignoreCache=false: Ignores the cache directory for reading but still writes to it
      --log=false: Enable Logging
      --logFile="": Log File path (if set, logging enabled automatically)
      --noTimes=false: Don't sync modification time of files
      --pluralizeListTitles=true: Pluralize titles in lists using inflect
  -s, --source="": filesystem path to read files relative from
      --stepAnalysis=false: display memory and timing of different steps of the program
  -t, --theme="": theme to use (located in /themes/THEMENAME/)
      --uglyUrls=false: if true, use /filename.html instead of /filename/
  -v, --verbose=false: verbose output
      --verboseLog=false: verbose logging
  -w, --watch=false: watch filesystem for changes and recreate as needed


Additional help topics:
 hugo convert         Convert will modify your content to different formats hugo list            Listing out various types of content

Use "hugo help [command]" for more information about a command.
</code></pre>

## Common Usage Example

The most common use is probably to run `hugo` with your current directory being the input directory:

    $ hugo
    0 draft content
    0 future content
    99 pages created
    0 paginator pages created
    16 tags created
    0 groups created
    in 120 ms

This generates your web site to the `public/` directory,
ready to be deployed to your web server.


## Instant feedback as you develop your web site

If you are working on things and want to see the changes immediately, tell Hugo to watch for changes.
Hugo will watch the filesystem for changes, and rebuild your site as soon as a file is saved:

    $ hugo -s ~/Code/hugo/docs --watch
    0 draft content
    0 future content
    99 pages created
    0 paginator pages created
    16 tags created
    0 groups created
    in 120 ms
    Watching for changes in /Users/spf13/Code/hugo/docs/content
    Press Ctrl+C to stop

Hugo can even run a server and create a site preview at the same time!
Hugo implements [LiveReload](/extras/livereload/) technology to automatically
reload any open pages in all JavaScript-enabled browsers, including mobile.
This is the easiest and most common way to develop a Hugo web site:

    $ hugo server -ws ~/Code/hugo/docs
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


## Deploying your web site

After running `hugo server` for local web development,
you need to do a final `hugo` run **without the `server` command**
and **without `--watch` or `-w`** to rebuild your site.
You may then **deploy your site** by copying the `public/` directory
(by FTP, SFTP, WebDAV, Rsync, git push, etc.) to your production web server.

Since Hugo generates a static website, your site can be hosted anywhere,
including [Heroku][], [GoDaddy][], [DreamHost][], [GitHub Pages][],
[Amazon S3][] and [CloudFront][], or any other cheap or even free
static web hosting services.

[Apache][], [nginx][], [IIS][]...  Any web server software would do!

[Apache]: http://httpd.apache.org/ "Apache HTTP Server"
[nginx]: http://nginx.org/
[IIS]: http://www.iis.net/
[Heroku]: https://www.heroku.com/
[GoDaddy]: https://www.godaddy.com/
[DreamHost]: http://www.dreamhost.com/
[GitHub Pages]: https://pages.github.com/
[Amazon S3]: http://aws.amazon.com/s3/
[CloudFront]: http://aws.amazon.com/cloudfront/ "Amazon CloudFront"


### Alternatively, serve your web site with Hugo!

Yes, that's right!  Because Hugo is so blazingly fast both in web site creation
*and* in web serving (thanks to its concurrent and multi-threaded design and
its Go heritage), some users actually prefer using Hugo itself to serve their
web site *on their production server*!

No other web server software (Apache, nginx, IIS...) is necessary.

Here is the command:

    hugo server --watch \
                --baseUrl=http://yoursite.org/ --port=80 \
                --appendPort=false
		--bind=87.245.198.50

Note the `bind` option, which is the interface to which the server will bind (defaults to `127.0.0.1`, which is fine for most development use cases). Some hosts, like Amazon WS, runs network address translation and it can sometimes be hard to figure out the actual IP address. Using `--bind=0.0.0.0` will bind to all interfaces.

This way, you may actually deploy just the source files,
and Hugo on your server will generate the resulting web site
on-the-fly and serve them at the same time.

You may optionally add `--disableLiveReload=true` if you do not want
the JavaScript code for LiveReload to be added to your web pages.

Interested? Here are some great tutorials contributed by Hugo users:

* [hugo, syncthing](http://fredix.ovh/2014/10/hugo-syncthing/) (French) by Frédéric Logier (@fredix)
* [服务器上 hugo 的安装和配置 <small>(Installing and configuring Hugo on the server)</small>](http://hucsmn.com/post/hugo-tutorial-make-it-work/) (Chinese) by hucsmn
