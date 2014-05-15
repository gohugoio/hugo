---
title: "Using Hugo"
date: "2013-07-01"
aliases: ["/doc/usage/"]
weight: 30
notoc: true
menu:
  main:
    parent: 'getting started'
---

Make sure either hugo is in your path or provide a path to it.

    $ hugo help
    A Fast and Flexible Static Site Generator
    built with love by spf13 and friends in Go.

    Complete documentation is available at http://hugo.spf13.com

    Usage:
      hugo [flags]
      hugo [command]

    Available Commands:
      server          :: Hugo runs it's own a webserver to render the files
      version         :: Print the version number of Hugo
      check           :: Check content in the source directory
      benchmark       :: Benchmark hugo by building a site a number of times
      help [command]  :: Help about any command

     Available Flags:
      -b, --base-url="": hostname (and path) to the root eg. http://spf13.com/
      -D, --build-drafts=false: include content marked as draft
          --config="": config file (default is path/config.yaml|json|toml)
      -d, --destination="": filesystem path to write files to
          --log=false: Enable Logging
          --logfile="": Log File path (if set, logging enabled automatically)
      -s, --source="": filesystem path to read files relative from
          --uglyurls=false: if true, use /filename.html instead of /filename/
      -v, --verbose=false: verbose output
          --verboselog=false: verbose logging
      -w, --watch=false: watch filesystem for changes and recreate as needed

    Use "hugo help [command]" for more information about that command.

## Common Usage Example:

The most common use is probably to run hugo with your current
directory being the input directory.


    $ hugo
    > X pages created
    > Y indexes created
      in 8 ms


If you are working on things and want to see the changes
immediately, tell Hugo to watch for changes. **It will
recreate the site faster than you can tab over to
your browser to view the changes.**

    $ hugo -s ~/mysite --watch
       28 pages created
       0 tags index created
       in 18 ms
       Watching for changes in /Users/spf13/Code/hugo/docs/content
       Press ctrl+c to stop

Hugo can even run a server and create your site at the same time!

    $ hugo server -ws ~/mysite
       Watching for changes in /Users/spf13/Code/hugo/docs/content
       Web Server is available at http://localhost:1313
       Press ctrl+c to stop
       28 pages created
       0 tags created
       in 18 ms

