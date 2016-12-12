---
lastmod: 2016-08-09
date: 2014-05-26
menu:
  main:
    parent: extras
next: /extras/menus
prev: /extras/gitinfo
title: LiveReload
---

Hugo may not be the first static site generator to utilize LiveReload
technology, but it’s the first to do it right.

The combination of Hugo’s insane build speed and LiveReload make
crafting your content pure joy. Virtually instantly after you hit save
your rebuilt content will appear in your browser.

## Using LiveReload

Hugo comes with LiveReload built in. There are no additional packages to
install. A common way to use Hugo while developing a site is to have
Hugo run a server and watch for changes:

{{< nohighlight >}}$ hugo server
{{< /nohighlight >}}

This will run a full functioning web server while simultaneously
watching your file system for additions, deletions or changes within
your:

 * static files
 * content
 * data files
 * layouts
 * current theme
 * configuration files

Whenever anything changes, Hugo will rebuild the site while continuing to serve
the content. As soon as the build is finished, it will tell the
browser and silently reload the page. Because most Hugo builds are so
fast they are barely noticeable, you merely need to glance at your open
browser and you will see the change, already there.

This means that keeping the site open on a second monitor (or another
half of your current monitor) allows you to see exactly what your
content looks like, without even leaving your text editor.

## Disabling Watch

If for some reason you don't want the Hugo server's watch functionality,
just do:

{{< nohighlight >}}$ hugo server --watch=false
{{< /nohighlight >}}

## Disabling LiveReload

LiveReload works by injecting JavaScript into the pages Hugo generates,
which creates a connection from the browser web socket client to the
Hugo web socket server.

Awesome for development, but not something you would want to do in
production. Since many people use `hugo server` in production to
instantly display any updated content, we’ve made it easy to disable the
LiveReload functionality:

{{< nohighlight >}}$ hugo server --disableLiveReload
{{< /nohighlight >}}

## Notes

You must have a closing `</body>` tag for LiveReload to work.
Hugo injects the LiveReload `<script>` before this tag.
