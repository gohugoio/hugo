---
title: "Using Hugo"
Pubdate: "2013-07-01"
---

Make sure either hugo is in your path or provide a path to it.

    $ hugo --help
    usage: hugo [flags] []
      -b, --base-url="": hostname (and path) to the root eg. http://spf13.com/
      -D, --build-drafts=false: include content marked as draft
          --config="": config file (default is path/config.yaml|json|toml)
      -d, --destination="": filesystem path to write files to
      -h, --help=false: show this help
          --port="1313": port to run web server on, default :1313
      -S, --server=false: run a (very) simple web server
      -s, --source="": filesystem path to read files relative from
          --uglyurls=false: use /filename.html instead of /filename/
      -v, --verbose=false: verbose output
          --version=false: which version of hugo
      -w, --watch=false: watch filesystem for changes and recreate as needed

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
       Watching for changes. Press ctrl+c to stop
       15 pages created
       0 tags created
       in 8 ms

Hugo can even run a server and create your site at the same time!

    $hugo --server -ws ~/mysite
       Watching for changes. Press ctrl+c to stop
       15 pages created
       0 tags created
       in 8 ms
       Web Server is available at http://localhost:1313
       Press ctrl+c to stop

