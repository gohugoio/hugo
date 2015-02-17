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

    $ hugo help
    A Fast and Flexible Static Site Generator
    built with love by spf13 and friends in Go.

    Complete documentation is available at http://gohugo.io

    Usage:
      hugo [flags]
      hugo [command]

    Available Commands:
      server                    Hugo runs its own webserver to render the files
      version                   Print the version number of Hugo
      check                     Check content in the source directory
      benchmark                 Benchmark hugo by building a site a number of times
      new [path]                Create new content for your site
      help [command]            Help about any command

     Available Flags:
      -b, --baseUrl="": hostname (and path) to the root eg. http://spf13.com/
      -D, --buildDrafts=false: build content marked as draft
      -F, --buildFuture=false: build content with PublishDate in the future
          --config="": config file (default is path/config.yaml|json|toml)
      -d, --destination="": filesystem path to write files to
          --disableRSS=false: Do not build RSS files
          --disableSitemap=false: Do not build Sitemap file
          --log=false: Enable Logging
          --logFile="": Log File path (if set, logging enabled automatically)
      -s, --source="": filesystem path to read files relative from
          --stepAnalysis=false: display memory and timing of different steps of the program
      -t, --theme="": theme to use (located in /themes/THEMENAME/)
          --uglyUrls=false: if true, use /filename.html instead of /filename/
      -v, --verbose=false: verbose output
          --verboseLog=false: verbose logging
      -w, --watch=false: watch filesystem for changes and recreate as needed

    Use "hugo help [command]" for more information about that command.

## Common Usage Example

The most common use is probably to run `hugo` with your current directory being the input directory.

    $ hugo
    > X pages created
      in 8 ms

If you are working on things and want to see the changes immediately, tell Hugo to watch for changes.

Hugo will watch the filesystem for changes, rebuild your site as soon as a file is saved.

    $ hugo -s ~/mysite --watch
       28 pages created
       in 18 ms
       Watching for changes in /Users/spf13/Code/hugo/docs/content
       Press Ctrl+C to stop

Hugo can even run a server and create a site preview at the same time! Hugo
implements [LiveReload](/extras/livereload/) technology to automatically reload any open pages in all browsers (including mobile). (Note that you'll need to run without -w before you deploy your site.)

    $ hugo server -ws ~/mysite
       Watching for changes in /Users/spf13/Code/hugo/docs/content
       Web Server is available at http://localhost:1313
       Press Ctrl+C to stop
       28 pages created
       0 tags created
       in 18 ms