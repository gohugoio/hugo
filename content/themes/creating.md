---
author: "Michael Henderson"
lastmod: 2016-09-01
date: 2015-11-26
linktitle: Creating a New Theme
menu:
  main:
    parent: tutorials
next: /tutorials/github-pages-blog
prev: /tutorials/automated-deployments
title: Creating a New Theme
weight: 10
---
## Introduction

This tutorial will show you how to create a simple theme in Hugo.

I'll introduce Hugo's use of templates,
and explain how to organize them into a theme.
The theme will grow, minimizing effort while meeting evolving needs.
To promote this focus, and to keep everything simple, I'll omit CSS styling.

We'll start by creating a tiny, blog-like web site.
We'll implement this blog with just one &mdash; quite basic &mdash; template.
Then we'll add an About page, and a few articles.
Overall, this web site (along with what you learn here)
will provide a good basis for you to continue working with Hugo in the future.
By making small variations,
you'll be able to create many different kinds of web sites.

I will assume you're comfortable with HTML, Markdown formatting,
and the Bash command line (possibly using [Git for
Windows](https://git-for-windows.github.io/)).

A few symbols might call for explanation: in this tutorial,
the commands you'll enter will be preceded by a `$` prompt &mdash;
and their output will follow.
`vi` means to open your editor; then `:wq` means to save the file.
Sometimes I'll add comments to explain a point &mdash; these start with `#`.
So, for example:
```bash
# this is a comment
$ echo this is a command
this is a command

# edit the file
$ vi foo.md
+++
date = "2040-01-18"
title = "creating a new theme"

+++
Bah! Humbug!
:wq

# show it
$ cat foo.md
+++
date = "2040-01-18"
title = "creating a new theme"

+++
Bah! Humbug!
```
## Definitions

Three concepts:

1. _Non-content_ files;
1. _Templates_ (as Hugo defines them); and
1. _Front-matter_

are essential for creating your first Hugo theme,
as well as your first Hugo website.
### Non-Content

The source files of a web site (destined to be rendered by Hugo)
are divided into two kinds:

1. The files containing its textual content (and nothing else &mdash;
except Hugo front-matter: see below, and Markdown styling); and
1. All other files. (These contain ***no*** textual content &mdash; ideally.)

Temporarily, let's affix the adjective _non-content_
to the latter kind of source files.

Non-content files are responsible for your web site's look and feel.
(Follow these article links from [Bop
Design](https://www.bopdesign.com/bop-blog/2013/11/what-is-the-look-and-feel-of-a-website-and-why-its-important/)
and
[Wikipedia](https://en.wikipedia.org/w/index.php?title=Look_and_feel&oldid=731052704)
if you wish for more information.)
They comprise its images, its CSS (for the sizes, colors and fonts),
its JavaScript (for the actions and reactions), and its Hugo templates
(which contain the rules Hugo uses to transform your content into HTML).

Given these files, Hugo will render a static web site &mdash;
informed by your content &mdash;
which contains the above images, HTML, CSS and JavaScript,
ready to be served to visitors.

Actually, a few of your invariant textual snippets
could reside in non-content files as well.
However, because someone might reuse your theme (eventually),
preferably you should keep those textual snippets in their own content files.
#### Where

Regarding where to create your non-content files, you have two options.
The simplest is the `./layouts/` and `./static/` filesystem trees.
If you choose this way,
then you needn't worry about configuring Hugo to find them.
Invariably, these are the first two places Hugo seeks for templates
(as well as images, CSS and JavaScript);
so in that case, it's guaranteed to find all your non-content files.

The second option is to create them in a filesystem tree
located somewhere under the `./themes/` directory.
If you choose that way,
then you must always tell Hugo where to search for them &mdash;
that's extra work, though. So, why bother?
#### Theme

Well &mdash; the difference between creating your non-content files under
`./layouts/` and `./static/` and creating them under `./themes/`
is admittedly very subtle.
Non-content files created under `./layouts/` and `./static/`
cannot be customized without editing them directly.
On the other hand, non-content files created under `./themes/`
can be customized, in another way. That way is both conventional
(for Hugo web sites) and non-destructive. Therefore,
creating your non-content files under `./themes/`
makes it easier for other people to use them.

The rest of this tutorial will call a set of non-content files a ***theme***
if they comprise a filesystem tree rooted anywhere under the
`./themes/` directory.

Note that you can use this tutorial to create your set of non-content files
under `./layouts/` and `./static/` if you wish. The only difference is that
you wouldn't need to edit your web site's configuration file
in order to select a theme.
### Home

The home page, or landing page,
is the first page that many visitors to a web site will see.
Often this is `/index.html`, located at the root URL of the web site.
Since Hugo writes files into the `./public/` tree,
your home page will reside in file `./public/index.html`.
### Configure

When Hugo runs, it first looks for an overall configuration file,
in order to read its settings, and applies them to the entire web site.
These settings override Hugo's default values.

The file can be in TOML, YAML, or JSON format.
I prefer TOML for my configuration files.
If you prefer JSON or YAML, you'll need to translate my examples.
You'll also need to change the basename, since Hugo uses its extension
to determine how to process it.

Hugo translates Markdown files into HTML.
By default, Hugo searches for Markdown files in the `./content/` tree
and template files under the `./themes/` directory.
It will render HTML files to the `./public/` tree.
You can override any of these defaults by specifying alternative locations
in the configuration file.
### Template

_Templates_ direct Hugo in rendering content into HTML;
they bridge content and presentation.

Rules in template files determine which content is published and where,
and precisely how it will be rendered into HTML files.
Templates also guide your web site's presentation
by specifying the CSS styling to use.

Hugo uses its knowledge of each piece of content
to seek a template file to use in rendering it.
If it can't find a template that matches the content, it will zoom out,
one conceptual level; it will then resume the search from there.
It will continue to do so, till it finds a matching template,
or runs out of templates to try.
Its last resort is your web site's default template,
which could conceivably be missing. If it finds no suitable template,
it simply forgoes rendering that piece of content.

It's important to note that _front-matter_ (see next)
can influence Hugo's template file selection process.
### Content

Content is stored in text files which contain two sections.
The first is called _front-matter_: this is information about the content.
The second contains Markdown-formatted text,
destined for conversion to HTML format.
#### Front-Matter

The _front-matter_ is meta-information describing the content.
Like the web site's configuration file, it can be written in the
TOML, YAML, or JSON formats.
Unlike the configuration file, Hugo doesn't use the file's extension
to determine the format.
Instead, it looks for markers in the file which signal this.
TOML is surrounded by "`+++`" and YAML by "`---`", but
JSON is enclosed in curly braces. I prefer to use TOML.
You'll need to translate my examples if you prefer YAML or JSON.

Hugo informs its chosen template files with the front-matter information
before rendering the content in HTML.
#### Markdown

Content is written in Markdown format, which makes it easy to create.
Hugo runs the content through a Markdown engine to transform it into HTML,
which it then renders to the output file.
### Template Kinds

Here I'll discuss three kinds of Hugo templates:
_Single_, _List_, and _Partial_.
All these kinds take one or more pieces of content as input,
and transform the pieces, based on commands in the template.
#### Single

A _Single_ template is used to render one piece of content.
For example, an article or a post is a single piece of content;
thus, it uses a Single template.
#### List

A _List_ template renders a group of related content items.
This could be a summary of recent postings,
or all of the articles in a category.
List templates can contain multiple groups (or categories).

The home page template is a special kind of List template.
This is because Hugo assumes that your home page will act as a portal
to all of the remaining content on your web site.
#### Partial

A _Partial_ template is a template that's incapable of producing a web page,
by itself. To include a Partial template in your web site,
another template must call it, using the `partial` command.

Partial templates are very handy for rolling up common behavior.
For example, you might want the same banner to appear on all
of your web site's pages &mdash; so, rather than copy your banner's text
into multiple content files,
as well as the other information relevant to your banner
into multiple template files (both Single and List),
you can instead create just one content file and one Partial template.
That way, whenever you decide to change the banner, you can do so
by editing one file only (or maybe two).
## Site

Let's let Hugo help you create your new web site.
The `hugo new site` command will generate a skeleton &mdash;
it will give you a basic directory structure, along with
a usable configuration file:
```bash
$ cd /tmp/

$ hugo new site mySite

$ cd mySite/

$ ls -l
total 8
drwxr-xr-x  2 {user} {group}   68 {date} archetypes
-rw-r--r--  1 {user} {group}  107 {date} config.toml
drwxr-xr-x  2 {user} {group}   68 {date} content
drwxr-xr-x  2 {user} {group}   68 {date} data
drwxr-xr-x  2 {user} {group}   68 {date} layouts
drwxr-xr-x  2 {user} {group}   68 {date} static
drwxr-xr-x  2 {user} {group}   68 {date} themes
```
Take a look in the `./content/` and `./themes/` directories to confirm
they are empty.

The other directories
(`./archetypes/`, `./data/`, `./layouts/` and `./static/`)
are used for customizing a named theme.
That's a topic for a different tutorial, so please ignore them for now.
### Render

Running the `hugo` command with no options will read
all of the available content and render the HTML files. Also, it will copy
all the static files (that's everything besides content).
Since we have an empty web site, Hugo won't be doing much.
However, generally speaking, Hugo does this very quickly:
```bash
$ hugo --verbose
INFO: {date} {source} Using config file: /tmp/mySite/config.toml
WARN: {date} {source} No theme set
INFO: {date} {source} /tmp/mySite/static/ is the only static directory available to sync from
INFO: {date} {source} syncing static files to /tmp/mySite/public/
Started building site
WARN: {date} {source} Unable to locate layout for homepage: [index.html _default/list.html]
WARN: {date} {source} "/" is rendered empty
=============================================================
Your rendered home page is blank: /index.html is zero-length
 * Did you specify a theme on the command-line or in your
   "config.toml" file?  (Current theme: "")
=============================================================
WARN: {date} {source} Unable to locate layout for 404 page: [404.html]
WARN: {date} {source} "404.html" is rendered empty
0 draft content
0 future content
0 expired content
0 pages created
0 non-page files copied
0 paginator pages created
0 tags created
0 categories created
in 4 ms
```
The "`--verbose`" flag gives extra information that will be helpful
whenever we are developing a template.
Every line of the output starting with "INFO:" or "WARN:" is present
because we used that flag. The lines that start with "WARN:"
are warning messages. We'll go over them later.

We can verify that the command worked by looking at the directory again:
```bash
$ ls -l
total 8
drwxr-xr-x  2 {user} {group}   68 {date} archetypes
-rw-r--r--  1 {user} {group}  107 {date} config.toml
drwxr-xr-x  2 {user} {group}   68 {date} content
drwxr-xr-x  2 {user} {group}   68 {date} data
drwxr-xr-x  2 {user} {group}   68 {date} layouts
drwxr-xr-x  6 {user} {group}  204 {date} public
drwxr-xr-x  2 {user} {group}   68 {date} static
drwxr-xr-x  2 {user} {group}   68 {date} themes
```
See that new `./public/` directory?
Hugo placed all its rendered content there.
When you're ready to publish your web site, that's the place to start.
For now, though, let's just confirm we have the files we expect
for a web site with no content:
```bash
$ ls -l public/
total 16
-rw-r--r--  1 {user} {group}    0 {date} 404.html
-rw-r--r--  1 {user} {group}    0 {date} index.html
-rw-r--r--  1 {user} {group}  511 {date} index.xml
-rw-r--r--  1 {user} {group}  210 {date} sitemap.xml
```
Hugo rendered two XML files and some empty HTML files.
The XML files are used for RSS feeds. Hugo has an opinion about what
those feeds should contain, so it populated those files.
Hugo has no opinion on the look or content of your web site,
so it left those files empty.

If you look back at the output from the `hugo server` command,
you'll notice that Hugo said:
```bash
0 pages created
```
That's because Hugo doesn't count the home page, the 404 error page,
or the RSS feed files as pages.
### Serve

Let's verify you can run the built-in web server &mdash;
that'll shorten your development cycle, dramatically.
Start it, by running the `hugo server` command.
If successful, you'll see output similar to the following:
```bash
$ hugo server --verbose
INFO: {date} {source} Using config file: /tmp/mySite/config.toml
WARN: {date} {source} No theme set
INFO: {date} {source} /tmp/mySite/static/ is the only static directory available to sync from
INFO: {date} {source} syncing static files to /
WARN: {date} {source} Unable to locate layout for homepage: [index.html _default/list.html]
WARN: {date} {source} "/" is rendered empty
=============================================================
Your rendered home page is blank: /index.html is zero-length
 * Did you specify a theme on the command-line or in your
   "config.toml" file?  (Current theme: "")
=============================================================
WARN: {date} {source} Unable to locate layout for 404 page: [404.html]
WARN: {date} {source} "404.html" is rendered empty
0 draft content
0 future content
0 expired content
0 pages created
0 non-page files copied
0 paginator pages created
0 tags created
0 categories created
in 3 ms
Watching for changes in /tmp/mySite/{data,content,layouts,static}
Serving pages from memory
Web Server is available at http://localhost:1313/ (bind address 127.0.0.1)
Press Ctrl+C to stop
```
Connect to the listed URL (it's on the line that begins with
`Web Server is available`). If everything's working correctly,
you should get a page that shows nothing.
### Warnings

Let's go back and look at some of those warnings again:
```bash
WARN: {date} {source} Unable to locate layout for 404 page: [404.html]
WARN: {date} {source} Unable to locate layout for homepage: [index.html _default/list.html]
```
The 404 warning is easy to explain &mdash; it's because we haven't created
the template file `layouts/404.html`. Hugo uses this to render an HTML file
which serves "page not found" errors. However,
the 404 page is a topic for a separate tutorial.

Regarding the home page warning: the first layout Hugo looked for was
`layouts/index.html`. Note that Hugo uses this file for the home page only.

It's good that Hugo lists the files it seeks, when
we give it the verbose flag. For the home page, these files are
`layouts/index.html` and `layouts/_default/list.html`.
Later, we'll cover some rules which explain these paths
(including their basenames). For now, just remember that
Hugo couldn't find a template to use for the home page, and it said so.

All right! So, now &mdash; after these few steps &mdash; you have a working
installation, and a web site foundation you can build upon.
All that's left is to add some content, as well as a theme to display it.
## Theme

Hugo doesn't ship with a default theme. However, a large number of themes
are easily available: for example, at
[hugoThemes](https://github.com/spf13/hugoThemes).
Also, Hugo comes with a command to generate them.

We're going to generate a new theme called Zafta.
The goal of this tutorial is simply to show you how to create
(in a theme) the minimal files Hugo needs in order to display your content.
Therefore, the theme will exclude CSS &mdash;
it'll be functional, not beautiful.

Every theme has its own opinions on content and layout. For example, this
Zafta theme prefers the Type "article" over the Types "blog" or "post."
Strong opinions make for simpler templates, but unconventional opinions
make themes tougher for other users. So when you develop a theme, you should
consider the value of adopting the terms used by themes similar to yours.
### Skeleton

Let's press Ctrl+C and use the `hugo new theme` command
to generate the skeleton of a theme. The result is a directory structure
containing empty files for you to fill out:
```bash
$ hugo new theme zafta

$ find themes -type f | xargs ls -l
-rw-r--r--  1 {user} {group}     8 {date} themes/zafta/archetypes/default.md
-rw-r--r--  1 {user} {group}     0 {date} themes/zafta/layouts/404.html
-rw-r--r--  1 {user} {group}     0 {date} themes/zafta/layouts/_default/list.html
-rw-r--r--  1 {user} {group}     0 {date} themes/zafta/layouts/_default/single.html
-rw-r--r--  1 {user} {group}     0 {date} themes/zafta/layouts/index.html
-rw-r--r--  1 {user} {group}     0 {date} themes/zafta/layouts/partials/footer.html
-rw-r--r--  1 {user} {group}     0 {date} themes/zafta/layouts/partials/header.html
-rw-r--r--  1 {user} {group}  1081 {date} themes/zafta/LICENSE.md
-rw-r--r--  1 {user} {group}   450 {date} themes/zafta/theme.toml
```
The skeleton includes templates (files ending in `.html`), a license file,
a description of your theme (`theme.toml`), and a default archetype file.

When you're developing a real theme, please remember to fill out files
`theme.toml` and `LICENSE.md`. They're optional, but if you're going to
distribute your theme, it tells the world who to praise (or blame).
It's also important to declare your choice of license, so people will know
whether (or where) they can use your theme.

Note that the skeleton theme's template files are empty. Don't worry;
we'll change that shortly:
```bash
$ find themes/zafta -name '*.html' | xargs ls -l
-rw-r--r--  1 {user} {group}  0 {date} themes/zafta/layouts/404.html
-rw-r--r--  1 {user} {group}  0 {date} themes/zafta/layouts/_default/list.html
-rw-r--r--  1 {user} {group}  0 {date} themes/zafta/layouts/_default/single.html
-rw-r--r--  1 {user} {group}  0 {date} themes/zafta/layouts/index.html
-rw-r--r--  1 {user} {group}  0 {date} themes/zafta/layouts/partials/footer.html
-rw-r--r--  1 {user} {group}  0 {date} themes/zafta/layouts/partials/header.html
```
### Select

Now that we've created a theme we can work with, it's a good idea
to add its name to the configuration file. This is optional, because
it's possible to add "-t zafta" to all your commands.
I like to put it in the configuration file because I like
shorter command lines. If you don't put it in the configuration file,
or specify it on the command line, sometimes you won't get the template
you're expecting.

So, let's edit your configuration file to add the theme name:
```toml
$ vi config.toml
theme = "zafta"
baseURL = "http://example.org/"
title = "My New Hugo Site"
languageCode = "en-us"
:wq
```
### Themed Render

Now that we have a theme (albeit empty), let's render the web site again:
```bash
$ hugo --verbose
INFO: {date} {source} Using config file: /tmp/mySite/config.toml
INFO: {date} {source} using a UnionFS for static directory comprised of:
INFO: {date} {source} Base: /tmp/mySite/themes/zafta/static
INFO: {date} {source} Overlay: /tmp/mySite/static/
INFO: {date} {source} syncing static files to /tmp/mySite/public/
Started building site
WARN: {date} {source} "/" is rendered empty
=============================================================
Your rendered home page is blank: /index.html is zero-length
 * Did you specify a theme on the command-line or in your
   "config.toml" file?  (Current theme: "zafta")
=============================================================
WARN: {date} {source} "404.html" is rendered empty
0 draft content
0 future content
0 expired content
0 pages created
0 non-page files copied
0 paginator pages created
0 tags created
0 categories created
in 4 ms
```
Did you notice the output is different?
Two previous warning messages have disappeared, which contained the words
"Unable to locate layout" for your home page and the 404 page.
And, a new informational message tells us Hugo is accessing your theme's tree
(`./themes/zafta/`).

Let's check the `./public/` directory to see what Hugo rendered:
```bash
$ ls -l public/
total 16
-rw-r--r--  1 {user} {group}    0 {date} 404.html
drwxr-xr-x  2 {user} {group}   68 {date} css
-rw-r--r--  1 {user} {group}    0 {date} index.html
-rw-r--r--  1 {user} {group}  511 {date} index.xml
drwxr-xr-x  2 {user} {group}   68 {date} js
-rw-r--r--  1 {user} {group}  210 {date} sitemap.xml
```
It's similar to what we had before, without a theme.
We'd expect so, since all your theme's templates are empty. But notice:
in `./public/`, Hugo created the `css/` and `js/` directories.
That's because Hugo found them in your theme's `static/` directory:
```bash
$ ls -l themes/zafta/static/
total 0
drwxr-xr-x  2 {user} {group}  68 {date} css
drwxr-xr-x  2 {user} {group}  68 {date} js
```
#### Home

In a Hugo web site, each kind of page is informed (primarily) by just one
of the many different kinds of templates available;
yet the home page is special, because it gets its own kind of template,
and its own template file.

Hugo uses template file `layouts/index.html` to render the home page's HTML.
Although Hugo's documentation may state that this file is the home page's
only required template, Hugo's earlier warning message showed it actually
looks for two different templates:
```bash
WARN: {date} {source} Unable to locate layout for homepage: [index.html _default/list.html]
```
#### Empty

When Hugo generated your theme, it included an empty home page template.
Whenever Hugo renders your web site, it seeks that same template and uses it
to render the HTML for the home page. Currently, the template file is empty,
so the output HTML file is empty, too. Whenever we add rules to that template,
Hugo will use them in rendering the home page:
```bash
$ find * -name index.html | xargs ls -l
-rw-r--r--  1 {user} {group}  0 {date} public/index.html
-rw-r--r--  1 {user} {group}  0 {date} themes/zafta/layouts/index.html
```
As we'll see later, Hugo follows this same pattern for all its templates.
## Static Files

Hugo does two things when it renders your web site.
Besides using templates to transform your content into HTML,
it also incorporates your static files. Hugo's rule is simple:
unlike with templates and content, static files aren't transformed.
Hugo copies them over, exactly as they are.

Hugo assumes that your web site will use both CSS and JavaScript,
so it generates some directories in your theme to hold them.
Remember opinions? Well, Hugo's opinion is that you'll store your CSS
in directory `static/css/`, and your JavaScript in directory `static/js/`.
If you don't like that, you can relocate these directories
or change their names (as long as they remain in your theme's `static/` tree),
or delete them completely.
Hugo is nice enough to offer its opinion; yet it still behaves nicely,
if you disagree:
```bash
$ find themes/zafta/* -type d | xargs ls -dl
drwxr-xr-x  3 {user} {group}  102 {date} themes/zafta/archetypes
drwxr-xr-x  6 {user} {group}  204 {date} themes/zafta/layouts
drwxr-xr-x  4 {user} {group}  136 {date} themes/zafta/layouts/_default
drwxr-xr-x  4 {user} {group}  136 {date} themes/zafta/layouts/partials
drwxr-xr-x  4 {user} {group}  136 {date} themes/zafta/static
drwxr-xr-x  2 {user} {group}   68 {date} themes/zafta/static/css
drwxr-xr-x  2 {user} {group}   68 {date} themes/zafta/static/js
```
## Theme Development

Generally (using any kind of software), working on a theme means
changing your files, serving your web site again, and then verifying
the resulting improvements in your browser.
With Hugo, this way of working is quite easy:

- First purge the `./public/` tree. (This is optional but useful,
if you want to start with a clean slate.)
- Run the built-in Hugo web server.
- Open your web site in a browser &mdash; and then:

1. Edit your theme;
1. Glance at your browser window to see your changes; and
1. Repeat.

I'll throw in one more opinion: ***never*** directly edit a theme on a live
web site. Instead, always develop ***using a copy***. First, make some changes
to your theme and test them. Afterwards, **when you've got them working,**
copy them to your web site. For added safety, use a tool like Git to keep
some revision history of your content, and of your theme. Believe me:
it's too easy to lose your changes, and your mind!

Check out the main Hugo web site for information about using Git with Hugo.
### Purge

When rendering your web site, Hugo will create new files in the `./public/`
tree and update existing ones. But it won't delete files that are
no longer used. For example, files previously rendered with
(what is now) the wrong basename, or in the wrong directory, will remain.
Later, if you leave them, they'll likely confuse you.
Cleaning out your `./public/` files prior to rendering can help.

When Hugo is running in web server mode (as of version 0.15),
it doesn't actually write the files. Instead,
it keeps all the rendered files in memory. So, you can "clean" up
your files simply by stopping and restarting the web server.
### Serve
#### Watch

Hugo's watch functionality monitors the relevant content, theme and
(overriding) site trees for filesystem changes,
and renders your web site again automatically, when changes are detected.

By default, watch is
enabled when in web server mode (`hugo server`),
but disabled for the web site renderer (`hugo`).

In some use cases,
Hugo's web site renderer should continue running and watch &mdash; simply
type `hugo --watch` on the command line.

Sometimes with Docker containers (and Heroku slugs),
the site sources may live on a read-only filesystem.
In that scenario, it makes no sense
for Hugo's web server to watch for file changes &mdash; so
use `hugo server --watch=false`.
#### Reload

Hugo's built in web server includes
[LiveReload](/extras/livereload/) functionality. When any page is updated
in the filesystem, the web browser is told to refresh its currently-open tabs
from your web site. Usually, this happens faster than you can say,
"Wow, that's totally amazing!"
### Workflow

Again,
I recommend you use the following commands as the basis for your workflow:
```bash
# purge old files. Hugo will recreate the public directory
$ rm -rf public/

# run Hugo in watch mode with LiveReload;
# when you're done, stop the web server
$ hugo server --verbose
Press Ctrl+C to stop
```
Below is some sample output showing Hugo detecting a change in the home page
template. (Actually, the change is the edit we're about to do.) Once it's
rendered again, the web browser automatically reloads the page.

(As I said above &mdash; it's amazing:)
```bash
$ rm -rf public/

$ hugo server --verbose
INFO: {date} {source} Using config file: /tmp/mySite/config.toml
INFO: {date} {source} using a UnionFS for static directory comprised of:
INFO: {date} {source} Base: /tmp/mySite/themes/zafta/static
INFO: {date} {source} Overlay: /tmp/mySite/static/
INFO: {date} {source} syncing static files to /
Started building site
WARN: {date} {source} "/" is rendered empty
=============================================================
Your rendered home page is blank: /index.html is zero-length
 * Did you specify a theme on the command-line or in your
   "config.toml" file?  (Current theme: "")
=============================================================
WARN: {date} {source} "404.html" is rendered empty
0 draft content
0 future content
0 expired content
0 pages created
0 non-page files copied
0 paginator pages created
0 tags created
0 categories created
in 4 ms
Watching for changes in /tmp/mySite/{data,content,layouts,static,themes}
Serving pages from memory
Web Server is available at http://localhost:1313/ (bind address 127.0.0.1)
Press Ctrl+C to stop
INFO: {date} {source} Received System Events: ["/tmp/mySite/themes/zafta/layouts/index.html": WRITE]

Change detected, rebuilding site
{date}
Template changed /tmp/mySite/themes/zafta/layouts/index.html
WARN: {date} {source} "404.html" is rendered empty
0 draft content
0 future content
0 expired content
0 pages created
0 non-page files copied
0 paginator pages created
0 tags created
0 categories created
in 3 ms
```
## Home Template

The home page is one of the few special pages Hugo renders automatically.
As mentioned earlier, it looks in your theme's `layouts/` tree for one
of two files:

1. `index.html`
1. `_default/list.html`

We could edit the default template, but a good design principle is to edit
the most specific template available. That's not a hard-and-fast rule
(in fact, in this tutorial, we'll break it a few times),
but it's a good generalization.
### Static

Right now, your home page is empty because you've added no content,
and because its template includes no logic. Let's change that by adding
some text to your home page template (`layouts/index.html`):
```html
$ vi themes/zafta/layouts/index.html
<!DOCTYPE html>
<html>
<body>
  <p>Hugo says hello!</p>
</body>
</html>
:wq
```
Let's press Ctrl+C and render the web site, and then verify the results:
```html
$ rm -rf public/

$ hugo --verbose
INFO: {date} {source} Using config file: /tmp/mySite/config.toml
INFO: {date} {source} using a UnionFS for static directory comprised of:
INFO: {date} {source} Base: /tmp/mySite/themes/zafta/static
INFO: {date} {source} Overlay: /tmp/mySite/static/
INFO: {date} {source} syncing static files to /tmp/mySite/public/
Started building site
WARN: {date} {source} "404.html" is rendered empty
0 draft content
0 future content
0 expired content
0 pages created
0 non-page files copied
0 paginator pages created
0 tags created
0 categories created
in 4 ms

$ ls -l public/index.html
-rw-r--r--  1 {user} {group}  72 {date} public/index.html

$ cat public/index.html
<!DOCTYPE html>
<html>
<body>
  <p>Hugo says hello!</p>
</body>
</html>
```
### Dynamic

A ***dynamic*** home page? Because Hugo is a _static web site_ generator,
the word _dynamic_ seems odd, doesn't it? But this means arranging for your
home page to reflect the content in your web site automatically,
each time Hugo renders it.

To accomplish that, later we'll add an iterator to your home page template.
## Article

Now that Hugo is successfully rendering your home page with static content,
let's add more pages to your web site. We'll display some new articles
as a list on your home page; and we'll display each article
on its own page, too.

Hugo has a command to generate an entry skeleton for new content,
just as it does for web sites and themes:
```bash
$ hugo --verbose new article/First.md
INFO: {date} {source} Using config file: /tmp/mySite/config.toml
INFO: {date} {source} attempting to create  article/First.md of article
INFO: {date} {source} curpath: /tmp/mySite/themes/zafta/archetypes/default.md
INFO: {date} {source} creating /tmp/mySite/content/article/First.md
/tmp/mySite/content/article/First.md created

$ ls -l content/article/
total 8
-rw-r--r--  1 {user} {group}  61 {date} First.md
```
Let's generate a second article, while we're here:
```bash
$ hugo --verbose new article/Second.md
INFO: {date} {source} Using config file: /tmp/mySite/config.toml
INFO: {date} {source} attempting to create  article/Second.md of article
INFO: {date} {source} curpath: /tmp/mySite/themes/zafta/archetypes/default.md
INFO: {date} {source} creating /tmp/mySite/content/article/Second.md
/tmp/mySite/content/article/Second.md created

$ ls -l content/article/
total 16
-rw-r--r--  1 {user} {group}  61 {date} First.md
-rw-r--r--  1 {user} {group}  62 {date} Second.md
```
Let's edit both those articles. Be careful to preserve their front-matter,
but append some text to their bodies, as follows:
```bash
$ vi content/article/First.md
In vel ligula tortor. Aliquam erat volutpat.
Pellentesque at felis eu quam tincidunt dignissim.
Nulla facilisi.

Pellentesque tempus nisi et interdum convallis.
In quam ante, vulputate at massa et, rutrum
gravida dui. Phasellus tristique libero at ex.
:wq

$ vi content/article/Second.md
Fusce lacus magna, maximus nec sapien eu,
porta efficitur neque. Aliquam erat volutpat.
Vestibulum enim nibh, posuere eu diam nec,
varius sagittis turpis.

Praesent quis sapien egestas mauris accumsan
pulvinar. Ut mattis gravida venenatis. Vivamus
lobortis risus id nisi rutrum, at iaculis.
:wq
```
So, for example, `./content/article/Second.md` becomes:
```toml
$ cat content/article/Second.md
+++
date = "2040-01-18T21:08:08-06:00"
title = "Second"

+++
Fusce lacus magna, maximus nec sapien eu,
porta efficitur neque. Aliquam erat volutpat.
Vestibulum enim nibh, posuere eu diam nec,
varius sagittis turpis.

Praesent quis sapien egestas mauris accumsan
pulvinar. Ut mattis gravida venenatis. Vivamus
lobortis risus id nisi rutrum, at iaculis.
```
Let's render the web site, and then verify the results:
```bash
$ rm -rf public/

$ hugo --verbose
INFO: {date} {source} Using config file: /tmp/mySite/config.toml
INFO: {date} {source} using a UnionFS for static directory comprised of:
INFO: {date} {source} Base: /tmp/mySite/themes/zafta/static
INFO: {date} {source} Overlay: /tmp/mySite/static/
INFO: {date} {source} syncing static files to /tmp/mySite/public/
Started building site
INFO: {date} {source} found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
WARN: {date} {source} "article" is rendered empty
WARN: {date} {source} "article/Second.html" is rendered empty
WARN: {date} {source} "article/First.html" is rendered empty
WARN: {date} {source} "404.html" is rendered empty
0 draft content
0 future content
0 expired content
2 pages created
0 non-page files copied
0 paginator pages created
0 tags created
0 categories created
in 7 ms
```
The output says Hugo rendered ("created") two pages.
Those pages are your new articles:
```bash
$ find public -type f -name '*.html' | xargs ls -l
-rw-r--r--  1 {user} {group}   0 {date} public/404.html
-rw-r--r--  1 {user} {group}   0 {date} public/article/First/index.html
-rw-r--r--  1 {user} {group}   0 {date} public/article/index.html
-rw-r--r--  1 {user} {group}   0 {date} public/article/Second/index.html
-rw-r--r--  1 {user} {group}  72 {date} public/index.html
```
The new pages are empty, because Hugo rendered their HTML from empty
template files. The home page doesn't show us the new content, either:
```html
$ cat public/index.html
<!DOCTYPE html>
<html>
<body>
  <p>Hugo says hello!</p>
</body>
</html>
```
So, we have to edit the templates, in order to pick up the articles.
### Single & List

Here again I'll discuss three kinds of Hugo templates. One kind is
the home page template we edited previously; it's applicable only to
the home page. Another kind is Single templates, which render output for
just one content file. The third kind are List templates, which group
multiple pieces of content before rendering output.

It's important to note that, generally, List templates
(except the home page template) are named `list.html`;
and Single templates are named `single.html`.

Hugo also has three other kinds of templates:
Partials, _Content Views_, and _Terms_.
We'll give examples of some Partial templates; but otherwise,
we won't go into much detail about these.
### Home

You'll want your home page to list the articles you just created.
So, let's alter its template file (`layouts/index.html`) to show them.
Hugo runs each template's logic whenever it renders that template's web page
(of course):
```html
$ vi themes/zafta/layouts/index.html
<!DOCTYPE html>
<html>
<body>
  {{- range first 10 .Data.Pages }}
    <h4><a href="{{ .Permalink }}">{{ .Title }}</a></h4>
  {{- end }}
</body>
</html>
:wq
```
#### Engine

Hugo uses the [Go language's template
engine](https://gohugo.io/templates/go-templates/).
That engine scans your template files for commands enclosed between
"{{" and "}}" (these are doubled, curly braces &mdash; affectionately
known as "mustaches").

BTW, a hyphen, if placed immediately after an opening mustache, or
immediately before a closing one, will prevent extraneous newlines.
(This can make Hugo's output look better, when viewed as text.)

So, the mustache commands in your newly-altered template are:

1. &nbsp;`range ...`
1. &nbsp;`.Permalink`
1. &nbsp;`.Title`
1. &nbsp;`end`

The `range` command is an iterator. We're using it to go through the latest
ten pages. (Hugo characterizes some of its HTML output files as "pages,"
but not all &mdash; see above.)

Looping through the list of data pages will consider each such HTML file
that Hugo renders (or rather &mdash; to speak more precisely &mdash; each
such HTML file that Hugo currently calculates it _will_ render).

It's helpful to remember that Hugo sets some variables, such as `.Data`, quite
early in its overall processing. Hugo loads information from every content
file into that variable, and gives all the templates a chance to process that
variable's contents, before actually rendering any HTML output files.

`.Permalink` supplies the URL which links to that article's page, and
`.Title` supplies the value of its "title" variable. Hugo obtains this
from the front-matter in the article's Markdown file.

Automatically, the pages are considered in descending order of the generation
times of their Markdown files (actually, based on the value of the "date"
variable in their front-matter) so that the latest is first (naturally).

The `end` command signals the end of the range iterator. The engine
loops back to the top of the iterator, whenever it finds `end.`
Everything between `range` and `end` is reevaluated,
each time the engine goes through the iterator.

For the present template, this means that the titles of your latest
ten pages (or however many exist, if that's less) become the
[textContent](https://developer.mozilla.org/en-US/docs/Web/API/Node/textContent)
of an equivalent number of copies Hugo makes, of your level-four
subheading tags (and anchor tags). `.Permalink` enables these to link
to the actual articles.

Let's render your web site, and then verify the results:
```html
$ rm -rf public/

$ hugo --verbose
INFO: {date} {source} Using config file: /tmp/mySite/config.toml
INFO: {date} {source} using a UnionFS for static directory comprised of:
INFO: {date} {source} Base: /tmp/mySite/themes/zafta/static
INFO: {date} {source} Overlay: /tmp/mySite/static/
INFO: {date} {source} syncing static files to /tmp/mySite/public/
Started building site
INFO: {date} {source} found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
WARN: {date} {source} "article" is rendered empty
WARN: {date} {source} "article/Second.html" is rendered empty
WARN: {date} {source} "article/First.html" is rendered empty
WARN: {date} {source} "404.html" is rendered empty
0 draft content
0 future content
0 expired content
2 pages created
0 non-page files copied
0 paginator pages created
0 tags created
0 categories created
in 7 ms

$ find public -type f -name '*.html' | xargs ls -l
-rw-r--r--  1 {user} {group}   0 {date} public/404.html
-rw-r--r--  1 {user} {group}   0 {date} public/article/First/index.html
-rw-r--r--  1 {user} {group}   0 {date} public/article/index.html
-rw-r--r--  1 {user} {group}   0 {date} public/article/Second/index.html
-rw-r--r--  1 {user} {group} 232 {date} public/index.html

$ cat public/index.html
<!DOCTYPE html>
<html>
<body>
    <h4><a href="http://example.org/article/Second/">Second</a></h4>
    <h4><a href="http://example.org/article/First/">First</a></h4>
</body>
</html>
```
### All Done

Congratulations! Your home page shows the titles of your two articles, along
with the links to them. The articles themselves are still empty. But,
let's take a moment to appreciate what we've done, so far!

Your home page template (`layouts/index.html`) now renders output dynamically.
Believe it or not, by inserting the range command inside those doubled
curly braces, you've learned everything you need to know &mdash;
essentially &mdash; about developing a theme.

All that's left is understanding which of your templates renders each content
file, and becoming more familiar with the commands for the template engine.
## More

Well &mdash; if things were so simple, this tutorial would be much shorter!

Some things are still useful to learn, because they'll make creating new
templates _much_ easier &mdash; so, I'll cover them, now.
### Base URL

While developing and testing your theme, did you notice that the links in the
rendered `./public/index.html` file use the full "baseURL" from your
`./config.toml` file? That's because those files are intended to be deployed
to your web server.

Whenever you test your theme, you start Hugo in web server mode
(with `hugo server`) and connect to it with your web browser.
That command is smart enough to replace the "baseURL" with
`http://localhost:1313` on the fly, so that the links automatically
work for you.

That's another reason why we recommend testing with the built-in web server.
### Content

The articles you've been working with are in your `./content/article/`
directory. That means their _Section_ (as far as templates are concerned)
is "article". Unless we do something unusual in their front-matter, their
_Type_ is also "article".
#### Search

Hugo uses the Section and Type to find a template file for every piece of
content it renders. Hugo first will seek a template file in subdirectories of
`layouts/` that match its Section or Type name (i.e., in `layouts/SECTION/`
or `layouts/TYPE/`). If it can't find a file there, then it will look in the
`layouts/_default/` directory. Other documentation covers some twists about
categories and tags, but we won't use those in this tutorial. Therefore,
we can assume that Hugo will try first `layouts/article/single.html`, then
`layouts/_default/single.html`.

Now that we know the search rule, let's see what's available:
```bash
$ find themes/zafta -name single.html | xargs ls -l
-rw-r--r--  1 {user} {group}  0 {date} themes/zafta/layouts/_default/single.html
```
If you look back at the articles Hugo has rendered, you can see that
they were empty. Now we can see that this is because Hugo sought
`layouts/article/single.html` but couldn't find it, and template
`layouts/_default/single.html` was empty. Therefore, the rendered article
file was empty, too.

So, we could either create a new template, `layouts/article/single.html`,
or edit the default one.
#### Default Single

Since we know of no other content Types, let's start by editing the default
template file, `layouts/_default/single.html`.

As we mentioned earlier, you always should edit (or create) the most
specific template first, in order to avoid accidentally changing how other
content is displayed. However, we're breaking that rule intentionally,
just so we can explore how the default is used.

Remember, any content &mdash; for which we don't create a specific template
&mdash; will end up using this default template. That can be good or bad.
Bad, because I know we'll be adding different Types of content, and we'll
eventually undo some of the changes we've made. Good, because then we'll be
able to see some results immediately. It's also good to create the default
template first, because with it, we can start to develop the basic layout
for the web site.

As we add more content Types, we'll refactor this file and move its logic
around. Hugo makes this fairly painless, so we'll accept the cost and proceed.

Please see Hugo's documentation on template rendering, for all the details on
determining which template to use. And, as the documentation mentions, if
your web site is a single-page application (SPA), you can delete all the
other templates and work with just the default Single one. By itself,
that fact provides a refreshing amount of joy.

Let's edit the default template file (`layouts/_default/single.html`):
```html
$ vi themes/zafta/layouts/_default/single.html
<!DOCTYPE html>
<html>
<head>
  <title>{{ .Title }}</title>
</head>
<body>
  <h1>{{ .Title }}</h1>
  <h6>{{ .Date.Format "Mon, Jan 2, 2006" }}</h6>
  {{ .Content }}
  <h4><a href="{{ .Site.BaseURL }}">Home</a></h4>
</body>
</html>
:wq
```
#### Verify

Let's render the web site, and verify the results:
```bash
$ rm -rf public/

$ hugo --verbose
INFO: {date} {source} Using config file: /tmp/mySite/config.toml
INFO: {date} {source} using a UnionFS for static directory comprised of:
INFO: {date} {source} Base: /tmp/mySite/themes/zafta/static
INFO: {date} {source} Overlay: /tmp/mySite/static/
INFO: {date} {source} syncing static files to /tmp/mySite/public/
Started building site
INFO: {date} {source} found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
WARN: {date} {source} "article" is rendered empty
WARN: {date} {source} "404.html" is rendered empty
0 draft content
0 future content
0 expired content
2 pages created
0 non-page files copied
0 paginator pages created
0 tags created
0 categories created
in 7 ms

$ find public -type f -name '*.html' | xargs ls -l
-rw-r--r--  1 {user} {group}    0 {date} public/404.html
-rw-r--r--  1 {user} {group}  473 {date} public/article/First/index.html
-rw-r--r--  1 {user} {group}    0 {date} public/article/index.html
-rw-r--r--  1 {user} {group}  514 {date} public/article/Second/index.html
-rw-r--r--  1 {user} {group}  232 {date} public/index.html
```
Note that although Hugo rendered a file, to list your articles:
`./public/article/index.html`, the file is empty, because we don't have
a template for it. (However: see next.) The other HTML files contain your
content, as we can see below:
```html
$ cat public/article/First/index.html
<!DOCTYPE html>
<html>
<head>
  <title>First</title>
</head>
<body>
  <h1>First</h1>
  <h6>Wed, Jan 18, 2040</h6>
  <p>In vel ligula tortor. Aliquam erat volutpat.
Pellentesque at felis eu quam tincidunt dignissim.
Nulla facilisi.</p>

<p>Pellentesque tempus nisi et interdum convallis.
In quam ante, vulputate at massa et, rutrum
gravida dui. Phasellus tristique libero at ex.</p>

  <h4><a href="http://example.org/">Home</a></h4>
</body>
</html>

$ cat public/article/Second/index.html
<!DOCTYPE html>
<html>
<head>
  <title>Second</title>
</head>
<body>
  <h1>Second</h1>
  <h6>Wed, Jan 18, 2040</h6>
  <p>Fusce lacus magna, maximus nec sapien eu,
porta efficitur neque. Aliquam erat volutpat.
Vestibulum enim nibh, posuere eu diam nec,
varius sagittis turpis.</p>

<p>Praesent quis sapien egestas mauris accumsan
pulvinar. Ut mattis gravida venenatis. Vivamus
lobortis risus id nisi rutrum, at iaculis.</p>

  <h4><a href="http://example.org/">Home</a></h4>
</body>
</html>
```
Again, notice that your rendered article files have content.
You can run `hugo server` and use your browser to confirm this.
You should see your home page, and it should contain the titles of both
articles. Each title should be a link to its respective article.

Each article should be displayed fully on its own page. And at the bottom of
each article, you should see a link which takes you back to your home page.
### Article List

Your home page still lists your most recent articles. However &mdash;
remember, from above, that I mentioned an empty file,
`./public/article/index.html`?
Let's make that show a list of ***all*** of your articles
(not just the latest ten).

We need to decide which template to edit. Key to this, is that
individual pages always come from Single templates. On the other hand,
only List templates are capable of rendering pages which display collections
(or lists) of other pages.

Because the new page will show a listing, we should select a List template.
Let's take a quick look to see which List templates are available already:
```bash
$ find themes/zafta -name list.html | xargs ls -l
-rw-r--r--  1 {user} {group}  0 {date} themes/zafta/layouts/_default/list.html
```
So, just as before with the single articles, so again now with the list of
articles, we must decide: whether to edit `layouts/_default/list.html`,
or to create `layouts/article/list.html`.
#### Default List

We still don't have multiple content Types &mdash; so, remaining consistent,
let's edit the default List template:
```html
$ vi themes/zafta/layouts/_default/list.html
<!DOCTYPE html>
<html>
<body>
  <h1>Articles</h1>
  {{- range first 10 .Data.Pages }}
    <h4><a href="{{ .Permalink }}">{{ .Title }}</a></h4>
  {{- end }}
  <h4><a href="{{ .Site.BaseURL }}">Home</a></h4>
</body>
</html>
:wq
```
Let's render everything again:
```bash
$ rm -rf public/

$ hugo --verbose
INFO: {date} {source} Using config file: /tmp/mySite/config.toml
INFO: {date} {source} using a UnionFS for static directory comprised of:
INFO: {date} {source} Base: /tmp/mySite/themes/zafta/static
INFO: {date} {source} Overlay: /tmp/mySite/static/
INFO: {date} {source} syncing static files to /tmp/mySite/public/
Started building site
INFO: {date} {source} found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
WARN: {date} {source} "404.html" is rendered empty
0 draft content
0 future content
0 expired content
2 pages created
0 non-page files copied
0 paginator pages created
0 categories created
0 tags created
in 7 ms

$ find public -type f -name '*.html' | xargs ls -l
-rw-r--r--  1 {user} {group}    0 {date} public/404.html
-rw-r--r--  1 {user} {group}  473 {date} public/article/First/index.html
-rw-r--r--  1 {user} {group}  327 {date} public/article/index.html
-rw-r--r--  1 {user} {group}  514 {date} public/article/Second/index.html
-rw-r--r--  1 {user} {group}  232 {date} public/index.html
```
Now (as you can see), we have a list of articles. To confirm it,
type `hugo server`; then, in your browser, navigate to `/article/`.
(Later, we'll link to it.)
## About

Let's add an About page, and try to display it at the top level
(as opposed to the next level down, where we placed your articles).
### Guide

Hugo's default goal is to let the directory structure of the `./content/`
tree guide the location of the HTML it renders to the `./public/` tree.
Let's check this, by generating an About page at the content's top level:
```toml
$ hugo new About.md
/tmp/mySite/content/About.md created

$ ls -l content/
total 8
-rw-r--r--   1 {user} {group}   61 {date} About.md
drwxr-xr-x   4 {user} {group}  136 {date} article

$ vi content/About.md
+++
date = "2040-01-18T22:01:00-06:00"
title = "About"

+++
Neque porro quisquam est qui dolorem
ipsum quia dolor sit amet consectetur
adipisci velit.
:wq
```
### Check

Let's render your web site, and check the results:
```html
$ rm -rf public/

$ hugo --verbose
INFO: {date} {source} Using config file: /tmp/mySite/config.toml
INFO: {date} {source} using a UnionFS for static directory comprised of:
INFO: {date} {source} Base: /tmp/mySite/themes/zafta/static
INFO: {date} {source} Overlay: /tmp/mySite/static/
INFO: {date} {source} syncing static files to /tmp/mySite/public/
Started building site
INFO: {date} {source} found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
WARN: {date} {source} "404.html" is rendered empty
0 draft content
0 future content
0 expired content
3 pages created
0 non-page files copied
0 paginator pages created
0 tags created
0 categories created
in 9 ms

$ ls -l public/
total 24
-rw-r--r--  1 {user} {group}     0 {date} 404.html
drwxr-xr-x  3 {user} {group}   102 {date} About
drwxr-xr-x  6 {user} {group}   204 {date} article
drwxr-xr-x  2 {user} {group}    68 {date} css
-rw-r--r--  1 {user} {group}   316 {date} index.html
-rw-r--r--  1 {user} {group}  2221 {date} index.xml
drwxr-xr-x  2 {user} {group}    68 {date} js
-rw-r--r--  1 {user} {group}   681 {date} sitemap.xml

$ ls -l public/About/
total 8
-rw-r--r--  1 {user} {group}  305 {date} index.html

$ cat public/About/index.html
<!DOCTYPE html>
<html>
<head>
  <title>About</title>
</head>
<body>
  <h1>About</h1>
  <h6>Wed, Jan 18, 2040</h6>
  <p>Neque porro quisquam est qui dolorem
ipsum quia dolor sit amet consectetur
adipisci velit.</p>

  <h4><a href="http://example.org/">Home</a></h4>
</body>
</html>
```
Oh, well. &mdash; Did you notice that your page wasn't rendered at the
top level? It was rendered to a subdirectory named `./public/About/`.
That name came from the basename of your Markdown file `./content/About.md`.
Interesting &mdash; but, we'll let that go, for now.
### Home

One other thing &mdash; let's take a look at your home page:
```html
$ cat public/index.html
<!DOCTYPE html>
<html>
<body>
    <h4><a href="http://example.org/About/">About</a></h4>
    <h4><a href="http://example.org/article/Second/">Second</a></h4>
    <h4><a href="http://example.org/article/First/">First</a></h4>
</body>
</html>
```
Did you notice that the About link is listed with your articles?
That's not exactly where we want it; so, let's edit your home page template
(`layouts/index.html`):
```html
$ vi themes/zafta/layouts/index.html
<!DOCTYPE html>
<html>
<body>
  <h2>Articles</h2>
  {{- range first 10 .Data.Pages -}}
    {{- if eq .Type "article"}}
      <h4><a href="{{ .Permalink }}">{{ .Title }}</a></h4>
    {{- end -}}
  {{- end }}
  <h2>Pages</h2>
  {{- range first 10 .Data.Pages -}}
    {{- if eq .Type "page" }}
      <h4><a href="{{ .Permalink }}">{{ .Title }}</a></h4>
    {{- end -}}
  {{- end }}
</body>
</html>
:wq
```
Let's render your web site, and verify the results:
```html
$ rm -rf public/

$ hugo --verbose
INFO: {date} {source} Using config file: /tmp/mySite/config.toml
INFO: {date} {source} using a UnionFS for static directory comprised of:
INFO: {date} {source} Base: /tmp/mySite/themes/zafta/static
INFO: {date} {source} Overlay: /tmp/mySite/static/
INFO: {date} {source} syncing static files to /tmp/mySite/public/
Started building site
INFO: {date} {source} found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
WARN: {date} {source} "404.html" is rendered empty
0 draft content
0 future content
0 expired content
3 pages created
0 non-page files copied
0 paginator pages created
0 tags created
0 categories created
in 9 ms

$ cat public/index.html
<!DOCTYPE html>
<html>
<body>
  <h2>Articles</h2>
      <h4><a href="http://example.org/article/Second/">Second</a></h4>
      <h4><a href="http://example.org/article/First/">First</a></h4>
  <h2>Pages</h2>
      <h4><a href="http://example.org/About/">About</a></h4>
</body>
</html>
```
Good! This time, your home page has two Sections: "article" and "page", and
each Section contains the correct set of headings and links.
## Template Sharing

If you've been following along on your computer, you might've noticed that
your home page doesn't show its title in your browser, although both of your
article pages do. That's because we didn't add your home page's title to its
template (`layouts/index.html`). That would be easy to do &mdash; but instead,
let's look at a better option.

We can put the common information into a shared template.
These reside in the `layouts/partials/` directory.
### Header & Footer

In Hugo (as elsewhere), a Partial is a template that's intended to be used
within other templates. We're going to create a Partial template that will
contain a header, for all of your page templates to use. That Partial will
enable us to maintain the header information in a single place, thus easing
our maintenance. Let's create both the header (`layouts/partials/header.html`)
and the footer (`layouts/partials/footer.html`):
```html
$ vi themes/zafta/layouts/partials/header.html
<!DOCTYPE html>
<html>
<head>
  <title>{{ .Title }}</title>
</head>
<body>
:wq

$ vi themes/zafta/layouts/partials/footer.html
  <h4><a href="{{ .Site.BaseURL }}">Home</a></h4>
</body>
</html>
:wq
```
### Calling

Any `partial` is called relative to its conventional location
`layouts/partials/`. So, you pass just the basename, followed by the context
(the period before the closing mustache). For example:
```bash
{{ partial "header.html" . }}
```
#### From Home

Let's change your home page template (`layouts/index.html`)
in order to use the new header Partial we just created:
```html
$ vi themes/zafta/layouts/index.html
{{ partial "header.html" . }}
  <h2>Articles</h2>
  {{- range first 10 .Data.Pages -}}
    {{- if eq .Type "article"}}
      <h4><a href="{{ .Permalink }}">{{ .Title }}</a></h4>
    {{- end -}}
  {{- end }}
  <h2>Pages</h2>
  {{- range first 10 .Data.Pages -}}
    {{- if eq .Type "page" }}
      <h4><a href="{{ .Permalink }}">{{ .Title }}</a></h4>
    {{- end -}}
  {{- end }}
</body>
</html>
:wq
```
Render your web site and verify the results. Now, the title on your home page
should be "My New Hugo Site". This comes from the "title" variable
in the `./config.toml` file.
#### From Default

Let's also edit the default templates (`layouts/_default/single.html` and
`layouts/_default/list.html`) to use your new Partials:
```html
$ vi themes/zafta/layouts/_default/single.html
{{ partial "header.html" . }}
  <h1>{{ .Title }}</h1>
  <h6>{{ .Date.Format "Mon, Jan 2, 2006" }}</h6>
  {{ .Content }}
{{ partial "footer.html" . -}}
:wq

$ vi themes/zafta/layouts/_default/list.html
{{ partial "header.html" . -}}
  <h1>Articles</h1>
  {{- range first 10 .Data.Pages }}
    <h4><a href="{{ .Permalink }}">{{ .Title }}</a></h4>
  {{- end }}
{{ partial "footer.html" . -}}
:wq
```
Render your web site and verify the results.
Now, the title of your About page should reflect the value of the "title"
variable in its corresponding Markdown file (`./content/About.md`).
The same should be true for each of your article pages as well (i.e.,
`./content/article/First.md` and `./content/article/Second.md`).
### DRY

Don't Repeat Yourself (also known as DRY) is a desirable goal,
in any kind of source code development &mdash;
and Hugo's partials do a fine job to help with that.

Part of the art of good templates is knowing when to add new ones, and when
to edit existing ones. While you're still figuring out the art of templates,
you should accept that you'll do some refactoring &mdash; Hugo makes this
easy and fast. And it's okay to delay splitting your templates into Partials.
## Section
### Date

Articles commonly display the date they were published
(or finalized) &mdash; so, here, let's do the same.

The front-matter of your articles contains a "date" variable
(as discussed above). Hugo sets this, when it creates each content file.
Now, sometimes an article requires many days to prepare, so its actual
publishing date might be later than the front-matter's "date". However, for
simplicity's sake, let's pretend this is the date we want to display, each time.

In Hugo, in order to format a variable date (or time),
we must do it by formatting the Go language [reference
time](https://golang.org/pkg/time/); for example:
```bash
{{ .Date.Format "Mon, Jan 2, 2006" }}
```
Now, your articles use the `layouts/_default/single.html` template (see above).
Because that template includes a date-formatting snippet, they show a
nice looking date. However, your About page uses the same default template.
Unfortunately, now it too shows its creation date (which makes no sense)!

There are a couple of ways to make the date display only for articles.
We could use an "if" statement, to display the date only when the Type equals
"article." That is workable, and acceptable for web sites with only a couple
of content Types. It aligns with the principle of "code for today," too.
### Template

Let's assume, though (for didactic purposes), that you've made your web site so
complex that you feel you must create a new template Type. In Hugo-speak, this
will be a new Section. It will contain your new, "article" Single template.

Let's restore your default Single template (`layouts/_default/single.html`)
to its earlier state (before we forget):
```html
$ vi themes/zafta/layouts/_default/single.html
{{ partial "header.html" . }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
{{ partial "footer.html" . -}}
:wq
```
Now, let's create your new template. If you remember Hugo's rules,
the template engine will prefer this version over the default. The first step
is to create (within your theme) its Section's directory: `layouts/article/`.
Then, create a Single template (`layouts/article/single.html`) within it:
```html
$ mkdir themes/zafta/layouts/article

$ vi themes/zafta/layouts/article/single.html
{{ partial "header.html" . }}
  <h1>{{ .Title }}</h1>
  <h6>{{ .Date.Format "Mon, Jan 2, 2006" }}</h6>
  {{ .Content }}
{{ partial "footer.html" . -}}
:wq
```
Basically, we moved the date logic &mdash; from the default template, to the
new "article" Section, Single template: `layouts/article/single.html`.

Let's render your web site and verify the results:
```html
$ rm -rf public/

$ hugo --verbose
INFO: {date} {source} Using config file: /tmp/mySite/config.toml
INFO: {date} {source} using a UnionFS for static directory comprised of:
INFO: {date} {source} Base: /tmp/mySite/themes/zafta/static
INFO: {date} {source} Overlay: /tmp/mySite/static/
INFO: {date} {source} syncing static files to /tmp/mySite/public/
Started building site
INFO: {date} {source} found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
WARN: {date} {source} "404.html" is rendered empty
0 draft content
0 future content
0 expired content
3 pages created
0 non-page files copied
0 paginator pages created
0 tags created
0 categories created
in 10 ms

$ cat public/article/First/index.html
<!DOCTYPE html>
<html>
<head>
  <title>First</title>
</head>
<body>

  <h1>First</h1>
  <h6>Wed, Jan 18, 2040</h6>
  <p>In vel ligula tortor. Aliquam erat volutpat.
Pellentesque at felis eu quam tincidunt dignissim.
Nulla facilisi.</p>

<p>Pellentesque tempus nisi et interdum convallis.
In quam ante, vulputate at massa et, rutrum
gravida dui. Phasellus tristique libero at ex.</p>

  <h4><a href="http://example.org/">Home</a></h4>
</body>
</html>

$ cat public/About/index.html
<!DOCTYPE html>
<html>
<head>
  <title>About</title>
</head>
<body>

  <h1>About</h1>
  <p>Neque porro quisquam est qui dolorem
ipsum quia dolor sit amet consectetur
adipisci velit.</p>

  <h4><a href="http://example.org/">Home</a></h4>
</body>
</html>
```
Now, as you can see, your articles show their dates,
and your About page (sensibly) doesn't.
