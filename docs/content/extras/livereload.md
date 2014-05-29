---
date: 2014-05-26
menu:
  main:
    parent: extras
next: /extras/menus
prev: /extras/comments
title: Live Reload
weight: 15
---

Hugo may not be the first static site generator to utilize live reload
technology, but it’s the first to do it right.

The combination of Hugo’s insane build speed and live reload make
crafting your content pure joy. Virtually instantly after you hit save
your rebuilt content will appear in your browser.

## Using livereload

Hugo comes with livereload built in. There are no additional packages to
install. A common way to use hugo while developing a site is to have
hugo run a server and watch for changes.

    hugo server --watch

This will run a full functioning web server while simultaneously
watching your file system for additions, deletions or changes within
your:

 * static files
 * content
 * layouts
 * current theme

Whenever anything changes Hugo will rebuild the site, continue to serve
the content and as soon as the build is finished it will tell the
browser and silently reload the page. Because most hugo builds are so
fast they are barely noticeable, you merely need to glance at your open
browser and you will see the change already there.

This means that keeping the site open on a second monitor (or another
half of your current monitor), allows you to see exactly what your
content looks like without even leaving your text editor.

## Disabling livereload

Live reload accomplishes this by injecting javascript into the pages it
creates that creates a web socket client to the hugo web socket server.

Awesome for development, but not something you would want to do in
production. Since many people use `hugo server --watch` in production to
instantly display any updated content, we’ve made it easy to disable the
live reload functionality.

    hugo server --watch --disableLiveReload





