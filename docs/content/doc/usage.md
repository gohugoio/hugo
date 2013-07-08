---
title: "Using Hugo"
Pubdate: "2013-07-01"
---

Make sure either hugo is in your path or provide a path to it.

    $ hugo --help
    usage: hugo [flags] []
      -b="": hostname (and path) to the root eg. http://spf13.com/
      -c="config.json": config file (default is path/config.json)
      -d=false: include content marked as draft
      -h=false: show this help
      -k=false: analyze content and provide feedback
      -p="": filesystem path to read files relative from
      -w=false: watch filesystem for changes and recreate as needed
      -s=false: a (very) simple webserver
      -port="1313": port for webserver to run on

## Common Usage Example:

The most common use is probably to run hugo with your current 
directory being the input directory.


    $ hugo
    > X pages created
    > Y indexes created


If you are working on things and want to see the changes 
immediately, tell Hugo to watch for changes. 
<br>
**It will 
recreate the site faster than you can tab over to 
your browser to view the changes.**

    $ hugo -p ~/mysite -w
       Watching for changes. Press ctrl+c to stop
       15 pages created
       0 tags created

Hugo can even run a server and create your site at the same time!

    $hugo -p ~/mysite -w -s
       Watching for changes. Press ctrl+c to stop
       15 pages created
       0 tags created
       Web Server is available at http://localhost:1313
       Press ctrl+c to stop


