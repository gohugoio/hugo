---
author: "Michael Henderson"
lastmod: 2015-12-01
date: 2015-11-26
linktitle: Creating a New Theme
toc: true
menu:
  main:
    parent: tutorials
next: /tutorials/github-pages-blog
prev: /tutorials/automated-deployments
title: Creating a New Theme
weight: 10
---


## Introduction

This tutorial will show you how to create a simple theme in Hugo. I assume that you are familiar with HTML, the bash command line, and that you are comfortable using Markdown to format content. I'll explain how Hugo uses templates and how you can organize your templates to create a theme. I won't cover using CSS to style your theme.

We'll start with creating a new site with a very basic template. Then we'll add in a few articles and an about page. With small variations on that, you will be able to create many different types of web sites.

In this tutorial, commands that you enter will start with the `$` prompt. The output will follow. Lines that start with `#` are comments that I've added to explain a point. When I show updates to a file, the `:wq` on the last line means to save the file.

Here's an example:

```bash
# this is a comment
$ echo this is a command
this is a command

# edit the file
$ vi foo.md
+++
date = "2015-11-26"
title = "creating a new theme"
+++

bah and humbug
:wq

# show it
$ cat foo.md
+++
date = "2015-11-26"
title = "creating a new theme"
+++

bah and humbug
$
```


## Some Definitions

There are a few concepts that you need to understand before creating a theme.

### Skins

Skins are the files responsible for the look and feel of your site. It’s the CSS that controls colors and fonts, it’s the Javascript that determines actions and reactions. It’s also the rules that Hugo uses to transform your content into the HTML that the site will serve to visitors.

You have two ways to create a skin. The simplest way is to create it in the `layouts/` directory. If you do, then you don’t have to worry about configuring Hugo to recognize it. The first place that Hugo will look for rules and files is in the `layouts/` directory so it will always find the skin.

Your second choice is to create it in a sub-directory of the `themes/` directory. If you do, then you must always tell Hugo where to search for the skin. It’s extra work, though, so why bother with it?

The difference between creating a skin in `layouts/` and creating it in `themes/` is very subtle. A skin in `layouts/` can’t be customized without updating the templates and static files that it is built from. A skin created in `themes/`, on the other hand, can be and that makes it easier for other people to use it.

The rest of this tutorial will call a skin created in the `themes/` directory a theme.

Note that you can use this tutorial to create a skin in the `layouts/` directory if you wish to. The main difference will be that you won’t need to update the site’s configuration file to use a theme.

### The Home Page

The home page, or landing page, is the first page that many visitors to a site see. It is the `index.html` file in the root directory of the web site. Since Hugo writes files to the `public/` directory, our home page is `public/index.html`.

### Site Configuration File

When Hugo runs, it looks for a configuration file that contains settings that override default values for the entire site. The file can use TOML, YAML, or JSON. I prefer to use TOML for my configuration files. If you prefer to use JSON or YAML, you’ll need to translate my examples. You’ll also need to change the name of the file since Hugo uses the extension to determine how to process it.

Hugo translates Markdown files into HTML. By default, Hugo expects to find Markdown files in your `content/` directory and template files in your `themes/` directory. It will create HTML files in your `public/` directory. You can change this by specifying alternate locations in the configuration file.

### Content

Content is stored in text files that contain two sections. The first section is the "front matter," which is the meta-information on the content. The second section contains Markdown that will be converted to HTML.

#### Front Matter

The front matter is information about the content. Like the configuration file, it can be written in TOML, YAML, or JSON. Unlike the configuration file, Hugo doesn’t use the file’s extension to know the format. It looks for markers to signal the type. TOML is surrounded by "`+++`", YAML by "`---`", and JSON is enclosed in curly braces. I prefer to use TOML, so you’ll need to translate my examples if you prefer YAML or JSON.

The information in the front matter is passed into the template before the content is rendered into HTML.

#### Markdown

Content is written in Markdown which makes it easier to create the content. Hugo runs the content through a Markdown engine to create the HTML which will be written to the output file.

### Template Files

Hugo uses template files to render content into HTML. Template files are a bridge between the content and presentation. Rules in the template define what content is published, where it's published to, and how it will rendered to the HTML file. The template guides the presentation by specifying the style to use.

There are three types of templates: single, list, and partial. Each type takes a bit of content as input and transforms it based on the commands in the template.

Hugo uses its knowledge of the content to find the template file used to render the content. If it can’t find a template that is an exact match for the content, it will shift up a level and search from there. It will continue to do so until it finds a matching template or runs out of templates to try. If it can’t find a template, it will use the default template for the site.

Please note that you can use the front matter to influence Hugo’s choice of templates.

#### Single Template

A single template is used to render a single piece of content. For example, an article or post would be a single piece of content and use a single template.

#### List Template

A list template renders a group of related content. That could be a summary of recent postings or all articles in a category. List templates can contain multiple groups.

The homepage template is a special type of list template. Hugo assumes that the home page of your site will act as the portal for the rest of the content in the site.

#### Partial Template

A partial template is a template that can be included in other templates. Partial templates must be called using the "partial" template command. They are very handy for rolling up common behavior. For example, your site may have a banner that all pages use. Instead of copying the text of the banner into every single and list template, you could create a partial with the banner in it. That way if you decide to change the banner, you only have to change the partial template.

## Create a New Site

Let's use Hugo to create a new web site. The `hugo new site` command will create a skeleton of a site. It will give you the basic directory structure and a useable configuration file.

```bash
$ hugo new site hugo-0.16
$ ls -l hugo-0.16
total 8
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:27 archetypes
-rw-r--r--  1 mdhender  wheel  107 Nov 27 20:27 config.toml
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:27 content
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:27 data
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:27 layouts
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:27 static
$ 
```

Take a look in the `content/` directory to confirm that it is empty.

The other directories (`archetypes/`, `data/`, `layouts/`, and `static/`) are used when customizing a named theme. That's a topic for a different tutorial, so please ignore them for now.

### Generate the HTML For the New Site

Running the `hugo` command with no options will read all the available content and generate the HTML files. It will also copy all static files (that's everything that's not content). Since we have an empty site, it won't do much, but it will do it very quickly.

```bash
$ cd hugo-0.16
$ hugo --verbose
INFO: 2015/11/27 Using config file: /tmp/hugo-0.16/config.toml
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/static/ to /tmp/hugo-0.16/public/
WARN: 2015/11/27 Unable to locate layout for homepage: [index.html _default/list.html]
WARN: 2015/11/27 Unable to locate layout for 404 page: [404.html]
0 draft content
0 future content
0 pages created
0 paginator pages created
0 tags created
0 categories created
in 4 ms
$ 
```

The "`--verbose`" flag gives extra information that will be helpful when we build the template. Every line of the output that starts with "INFO:" or "WARN:" is present because we used that flag. The lines that start with "WARN:" are warning messages. We'll go over them later.

We can verify that the command worked by looking at the directory again.

```bash
$ ls -l
total 8
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:27 archetypes
-rw-r--r--  1 mdhender  wheel  107 Nov 27 20:27 config.toml
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:27 content
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:27 data
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:27 layouts
drwxr-xr-x  6 mdhender  wheel  204 Nov 27 20:29 public
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:27 static
$
```

See that new `public/` directory? Hugo placed all generated content there. When you're ready to publish your web site, that's the place to start. For now, though, let's just confirm that we have what we'd expect from a site with no content.

```bash
$ ls -l public/
total 16
-rw-r--r--  1 mdhender  wheel    0 Nov 27 20:29 404.html
-rw-r--r--  1 mdhender  wheel    0 Nov 27 20:29 index.html
-rw-r--r--  1 mdhender  wheel  511 Nov 27 20:29 index.xml
-rw-r--r--  1 mdhender  wheel  237 Nov 27 20:29 sitemap.xml
$ 
```

Hugo created two XML files, which is standard, and empty HTML files. The XML files are used for RSS feeds. Hugo has an opinion on what those feeds should contain, so it populates those files. Hugo has no opinion on what your web site looks like (or contains), so it leaves those files empty.

If you look back over the output from the `hugo server` command, you will notice that Hugo said:

```bash
0 pages created
```

That's because Hugo doesn't count the homepage, the 404 error page, or the RSS feed files as pages.

### Test the New Site

Verify that you can run the built-in web server. It will dramatically shorten your development cycle if you do. Start it by running the `hugo server` command. If it is successful, you will see output similar to the following:

```bash
$ hugo server --verbose
INFO: 2015/11/27 Using config file: /tmp/hugo-0.16/config.toml
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/static/ to /
WARN: 2015/11/27 Unable to locate layout for homepage: [index.html _default/list.html]
WARN: 2015/11/27 Unable to locate layout for 404 page: [404.html]
0 draft content
0 future content
0 pages created
0 paginator pages created
0 tags created
0 categories created
in 3 ms
Watching for changes in /tmp/hugo-0.16/{data,content,layouts,static}
Serving pages from memory
Web Server is available at http://localhost:1313/ (bind address 127.0.0.1)
Press Ctrl+C to stop
```

Connect to the listed URL (it's on the line that starts with `Web Server is available`.). If everything is working correctly, you should get a page that shows nothing.

Let’s go back and look at those warnings again.

```bash
WARN: 2015/11/27 Unable to locate layout for homepage: [index.html _default/list.html]
WARN: 2015/11/27 Unable to locate layout for 404 page: [404.html]
```

That second warning is easier to explain. We haven’t created a template to be used to generate "page not found errors." The 404 message is a topic for a separate tutorial.

Now for the first warning. It is for the home page. You can tell because the first layout that it looked for was `index.html`. That’s only used by the home page.

I like that the verbose flag causes Hugo to list the files that it's searching for. For the home page, they are `index.html` and `_default/list.html`. There are some rules that we'll cover later that explain the names and paths. For now, just remember that Hugo couldn't find a template for the home page and it told you so.

At this point, you've got a working installation and site that we can build upon. All that’s left is to add some content and a theme to display it.

## Create a New Theme

Hugo doesn't ship with a default theme. There are a few available (I counted a dozen when I first installed Hugo) and Hugo comes with a command to create new themes.

We're going to create a new theme called "zafta." Since the goal of this tutorial is to show you how to fill out the files to pull in your content, the theme will not contain any CSS. In other words, ugly but functional.

All themes have opinions on content and layout. For example, Zafta uses "article" over "blog" or "post." Strong opinions make for simpler templates but differing opinions make it tougher to use themes. When you build a theme, consider using the terms that other themes do.

### Create a Skeleton

Use the `hugo new theme` command to create the skeleton of a theme. This creates the directory structure and places empty files for you to fill out.

```bash
$ hugo new theme zafta

$ ls -l
total 8
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:27 archetypes
-rw-r--r--  1 mdhender  wheel  107 Nov 27 20:27 config.toml
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:27 content
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:27 data
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:27 layouts
drwxr-xr-x  6 mdhender  wheel  204 Nov 27 20:29 public
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:27 static
drwxr-xr-x  3 mdhender  wheel  102 Nov 27 20:35 themes


$ find themes -type f | xargs ls -l
-rw-r--r--  1 mdhender  wheel  1081 Nov 27 20:35 themes/zafta/LICENSE.md
-rw-r--r--  1 mdhender  wheel     8 Nov 27 20:35 themes/zafta/archetypes/default.md
-rw-r--r--  1 mdhender  wheel     0 Nov 27 20:35 themes/zafta/layouts/404.html
-rw-r--r--  1 mdhender  wheel     0 Nov 27 20:35 themes/zafta/layouts/_default/list.html
-rw-r--r--  1 mdhender  wheel     0 Nov 27 20:35 themes/zafta/layouts/_default/single.html
-rw-r--r--  1 mdhender  wheel     0 Nov 27 20:35 themes/zafta/layouts/index.html
-rw-r--r--  1 mdhender  wheel     0 Nov 27 20:35 themes/zafta/layouts/partials/footer.html
-rw-r--r--  1 mdhender  wheel     0 Nov 27 20:35 themes/zafta/layouts/partials/header.html
-rw-r--r--  1 mdhender  wheel   450 Nov 27 20:35 themes/zafta/theme.toml
$ 
```

The skeleton includes templates (the files ending in `.html`), license file, a description of your theme (the `theme.toml` file), and a default archetype file.

When you're creating a real theme, please remember to fill out the `theme.toml` and `LICENSE.md` files. They're optional, but if you're going to be distributing your theme, it tells the world who to praise (or blame). It's also nice to declare the license so that people will know how they can use the theme.

Note that the theme skeleton's template files are empty. Don't worry, we'll be changing that shortly.

```bash
$ find themes/zafta -name '*.html' | xargs ls -l
-rw-r--r--  1 mdhender  wheel  0 Nov 27 20:35 themes/zafta/layouts/404.html
-rw-r--r--  1 mdhender  wheel  0 Nov 27 20:35 themes/zafta/layouts/_default/list.html
-rw-r--r--  1 mdhender  wheel  0 Nov 27 20:35 themes/zafta/layouts/_default/single.html
-rw-r--r--  1 mdhender  wheel  0 Nov 27 20:35 themes/zafta/layouts/index.html
-rw-r--r--  1 mdhender  wheel  0 Nov 27 20:35 themes/zafta/layouts/partials/footer.html
-rw-r--r--  1 mdhender  wheel  0 Nov 27 20:35 themes/zafta/layouts/partials/header.html
$
```

### Update the Configuration File to Use the Theme

Now that we've got a theme to work with, it's a good idea to add the theme name to the configuration file. This is optional, because you can always add "-t zafta" on all your commands. I like to put it the configuration file because I like shorter command lines. If you don't put it in the configuration file or specify it on the command line, you won't use the template that you're expecting to.

Edit the file to add the named theme.

```bash
$ vi config.toml
theme = "zafta"
baseurl = "http://replace-this-with-your-hugo-site.com/"
languageCode = "en-us"
title = "My New Hugo Site"
:wq

$
```

### Generate the Site

Now that we have an empty theme, let's generate the site again.

```bash
$ hugo --verbose
INFO: 2015/11/27 Using config file: /tmp/hugo-0.16/config.toml
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/themes/zafta/static to /tmp/hugo-0.16/public/
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/static/ to /tmp/hugo-0.16/public/
0 draft content
0 future content
0 pages created
0 paginator pages created
0 tags created
0 categories created
in 4 ms
$
```

Did you notice that the output is different? The warning message for the home page has disappeared and we have an additional information line saying that Hugo is syncing from the theme's directory (`themes/zafta/`).

Let's check the `public/` directory to see what Hugo's created.

```bash
$ ls -l public
total 16
-rw-r--r--  1 mdhender  wheel    0 Nov 27 20:42 404.html
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:35 css
-rw-r--r--  1 mdhender  wheel    0 Nov 27 20:42 index.html
-rw-r--r--  1 mdhender  wheel  511 Nov 27 20:42 index.xml
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:35 js
-rw-r--r--  1 mdhender  wheel  237 Nov 27 20:42 sitemap.xml
$
```

It's similar to what we had without a theme. We'd expect that since our theme has only empty templates. But notice that Hugo created `css/` and `js/` directories. That's due to our template having those in its `static/` directory:

```bash
$ ls -l themes/zafta/static/
total 0
drwxr-xr-x  2 mdhender  wheel  68 Nov 27 20:35 css
drwxr-xr-x  2 mdhender  wheel  68 Nov 27 20:35 js
$ 
```

The rule with static files is simple: Hugo copies them over without any changes.

#### The Home Page

Hugo supports many different types of templates. The home page is special because it gets its own type of template and its own template file. The file `layouts/index.html` is used to generate the HTML for the home page. The Hugo documentation says that this is the only required template, but that depends. Hugo's warning message shows that it looks for two different templates:

```bash
WARN: 2015/11/27 Unable to locate layout for homepage: [index.html _default/list.html]
```

When Hugo created our theme, it created an empty home page template. Now, when we build the site, Hugo finds the template and uses it to generate the HTML for the home page. Since the template file is empty, the HTML file is empty, too. If the template had any rules in it, then Hugo would have used them to generate the home page.

```bash
$ find . -name index.html | xargs ls -l
-rw-r--r--  1 mdhender  wheel  0 Nov 27 20:42 ./public/index.html
-rw-r--r--  1 mdhender  wheel  0 Nov 27 20:35 ./themes/zafta/layouts/index.html
$ 
```

#### The Magic of Static

Hugo does two things when generating the site. It uses templates to transform content into HTML and it copies static files into the site. Unlike content, static files are not transformed. They are copied exactly as they are.

Hugo assumes that your site will use both CSS and JavaScript, so it creates directories in your theme to hold them. Remember opinions? Well, Hugo's opinion is that you'll store your CSS in a directory named `css/` and your JavaScript in a directory named `js/`. If you don't like that, you can change the directory names in your theme's `static/` directory or even delete them completely. Hugo's nice enough to offer its opinion, then behave nicely if you disagree.

```bash
$ find themes/zafta -type d | xargs ls -ld
drwxr-xr-x  7 mdhender  wheel  238 Nov 27 20:35 themes/zafta
drwxr-xr-x  3 mdhender  wheel  102 Nov 27 20:35 themes/zafta/archetypes
drwxr-xr-x  6 mdhender  wheel  204 Nov 27 20:35 themes/zafta/layouts
drwxr-xr-x  4 mdhender  wheel  136 Nov 27 20:35 themes/zafta/layouts/_default
drwxr-xr-x  4 mdhender  wheel  136 Nov 27 20:35 themes/zafta/layouts/partials
drwxr-xr-x  4 mdhender  wheel  136 Nov 27 20:35 themes/zafta/static
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:35 themes/zafta/static/css
drwxr-xr-x  2 mdhender  wheel   68 Nov 27 20:35 themes/zafta/static/js
$ 
```

## The Theme Development Cycle

When you're working on a theme, you will make changes in the theme's directory, rebuild the site, and check your changes in the browser. Hugo makes this very easy:

1. Purge the `public/` directory (optional, but useful if you want to start with a clean slate).
2. Run the built in web server.
3. Open your site in a browser.
4. Update the theme.
5. Glance at your browser window to see changes.
6. Return to step 4.

I’ll throw in one more opinion: never work on a theme on a live site. Always work on a copy of your site. Make changes to your theme, test them, then copy them up to your site. For added safety, use a tool like Git to keep a revision history of your content and your theme. Believe me when I say that it is too easy to lose both your mind and your changes.

Check the main Hugo site for information on using Git with Hugo.

### Purge the public/ Directory

When rendering the site, Hugo will create new files and update existing ones in the `public/` directory. It will not delete files that are no longer used. For example, files that were created in the wrong directory or with the wrong title will remain. If you leave them, you might get confused by them later. Cleaning out your public files prior to rendering can help.

As of version 0.15, Hugo doesn't write files when running in server mode. Instead, it keeps all the rendered files in memory. You can "clean" up files by stopping and restarting the server.

### Hugo's Watch Option

Hugo's "`--watch`" option will monitor the content/ and your theme directories for changes and rebuild the site automatically. Since version 0.15, this has been the default option for `hugo server`.

### Live Reload

Hugo's built in web server supports live reload. As pages are saved on the server, the browser is told to refresh the page. Usually, this happens faster than you can say, "Wow, that's totally amazing."

### Development Commands

Use the following commands as the basis for your workflow.

```bash
# purge old files. hugo will recreate the public directory.
#
$ rm -rf public
#
# run hugo in watch mode with live reload
#
$ hugo server --verbose
#
# hit Control+C to kill the server when you're done
#
```

Here's sample output showing Hugo detecting a change to the template for the home page. Once generated, the web browser automatically reloaded the page. I've said this before, it's amazing.


```bash
$ hugo server --verbose
INFO: 2015/11/27 Using config file: /tmp/hugo-0.16/config.toml
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/themes/zafta/static to /
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/static/ to /
0 draft content
0 future content
0 pages created
0 paginator pages created
0 tags created
0 categories created
in 4 ms
Watching for changes in /tmp/hugo-0.16/{data,content,layouts,static,themes}
Serving pages from memory
Web Server is available at http://localhost:1313/ (bind address 127.0.0.1)
Press Ctrl+C to stop

INFO: 2015/11/27 File System Event: ["/tmp/hugo-0.16/themes/zafta/layouts/index.html": CHMOD "/tmp/hugo-0.16/themes/zafta/layouts/index.html": WRITE]

Change detected, rebuilding site
2015-11-27 20:57 -0600
0 draft content
0 future content
0 pages created
0 paginator pages created
0 tags created
0 categories created
in 3 ms
```

## Update the Home Page Template

The home page is one of a few special pages that Hugo creates automatically. As mentioned earlier, it looks for one of two files in the theme's `layout/` directory:

1. `index.html`
2. `_default/list.html`

We could update the default templates, but a good design decision is to update the most specific template available. That's not a hard and fast rule (in fact, we'll break it a few times in this tutorial), but it is a good generalization.

### Make a Static Home Page

Right now, that page is empty because we don't have any content and we don't have any logic in the template. Let's change that by adding some text to the template.

```bash
$ vi themes/zafta/layouts/index.html
<!DOCTYPE html>
<html>
<body>
  <p>hugo says hello!</p>
</body>
</html>
:wq

$
```

Build the web site and then verify the results.

```bash
$ hugo --verbose
INFO: 2015/11/27 Using config file: /tmp/hugo-0.16/config.toml
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/themes/zafta/static to /tmp/hugo-0.16/public/
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/static/ to /tmp/hugo-0.16/public/
0 draft content
0 future content
0 pages created
0 paginator pages created
0 tags created
0 categories created
in 4 ms

$ ls -l public/index.html 
-rw-r--r--  1 mdhender  wheel  72 Nov 27 21:03 public/index.html
$ cat public/index.html 
<!DOCTYPE html>
<html>
<body>
  <p>hugo says hello!</p>
</body>
</html>

$ 
```

### Build a "Dynamic" Home Page

"Dynamic home page?" Hugo's a static web site generator, so this seems an odd thing to say. I mean let's have the home page automatically reflect the content in the site every time Hugo builds it. We'll use iteration in the template to do that.

#### Create New Articles

Now that we have the home page generating static content, let's add some content to the site. We'll display these articles as a list on the home page and on their own page, too.

Hugo has a command to generate a skeleton entry for new content, just like it does for sites and themes.

```bash
 hugo --verbose new article/first.md
INFO: 2015/11/27 Using config file: /tmp/hugo-0.16/config.toml
INFO: 2015/11/27 attempting to create  article/first.md of article
INFO: 2015/11/27 curpath: /tmp/hugo-0.16/themes/zafta/archetypes/default.md
INFO: 2015/11/27 creating /tmp/hugo-0.16/content/article/first.md
/tmp/hugo-0.16/content/article/first.md created

$ ls -l content/article/
total 8
-rw-r--r--  1 mdhender  wheel  61 Nov 27 21:06 first.md
$ 
```

Let's create a second article while we're here.

```bash
$ hugo --verbose new article/second.md
INFO: 2015/11/27 Using config file: /tmp/hugo-0.16/config.toml
INFO: 2015/11/27 attempting to create  article/second.md of article
INFO: 2015/11/27 curpath: /tmp/hugo-0.16/themes/zafta/archetypes/default.md
INFO: 2015/11/27 creating /tmp/hugo-0.16/content/article/second.md
/tmp/hugo-0.16/content/article/second.md created

$ ls -l content/article/
total 16
-rw-r--r--  1 mdhender  wheel  61 Nov 27 21:06 first.md
-rw-r--r--  1 mdhender  wheel  62 Nov 27 21:08 second.md
```

Edit both of those articles to put some text into the body.

```bash
$ cat content/article/first.md 
+++
date = "2015-11-27T21:06:38-06:00"
title = "first"
+++
In vel ligula tortor. Aliquam erat volutpat.
Pellentesque at felis eu quam tincidunt dignissim.
Nulla facilisi.

Pellentesque tempus nisi et interdum convallis.
In quam ante, vulputate at massa et, rutrum
gravida dui. Phasellus tristique libero at ex.

$ cat content/article/second.md 
+++
date = "2015-11-27T21:08:08-06:00"
title = "second"
+++
Fusce lacus magna, maximus nec sapien eu,
porta efficitur neque. Aliquam erat volutpat.
Vestibulum enim nibh, posuere eu diam nec,
varius sagittis turpis.

Praesent quis sapien egestas mauris accumsan
pulvinar. Ut mattis gravida venenatis. Vivamus
lobortis risus id nisi rutrum, at iaculis.
$ 
```

Build the web site and then verify the results.

```bash
$ rm -rf public/
$ hugo --verbose
INFO: 2015/11/27 Using config file: /tmp/hugo-0.16/config.toml
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/themes/zafta/static to /tmp/hugo-0.16/public/
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/static/ to /tmp/hugo-0.16/public/
INFO: 2015/11/27 found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
0 draft content
0 future content
2 pages created
0 paginator pages created
0 categories created
0 tags created
in 7 ms
$
```

The output says that it created 2 pages. Those are our new articles:

```bash
$ find public -type f -name '*.html' | xargs ls -l
-rw-r--r--  1 mdhender  wheel   0 Nov 27 21:12 public/404.html
-rw-r--r--  1 mdhender  wheel   0 Nov 27 21:12 public/article/first/index.html
-rw-r--r--  1 mdhender  wheel   0 Nov 27 21:12 public/article/index.html
-rw-r--r--  1 mdhender  wheel   0 Nov 27 21:12 public/article/second/index.html
-rw-r--r--  1 mdhender  wheel  72 Nov 27 21:12 public/index.html
$ 
```

The new files are empty because because the templates used to generate the content are empty. The homepage doesn't show the new content, either.

```bash
$ cat public/index.html 
<!DOCTYPE html>
<html>
<body>
  <p>hugo says hello!</p>
</body>
</html>
$ 
```

We have to update the templates to add the articles.

### List and Single Templates

In Hugo, we have three major kinds of templates. There's the home page template that we updated previously. It is used only by the home page. We also have "single" templates which are used to generate output for a single content file. We also have "list" templates that are used to group multiple pieces of content before generating output.

Generally speaking, list templates are named "list.html" and single templates are named "single.html."

There are three other types of templates: partials, content views, and terms. We will give an example of some partials, but really won't go into much detail on these.

### Add Content to the Homepage

The home page will contain a list of articles. Let's update its template to add the articles that we just created. The logic in the template will run every time we build the site.

```bash
$ vi themes/zafta/layouts/index.html
<!DOCTYPE html>
<html>
<body>
  {{ range first 10 .Data.Pages }}
    <h1><a href="{{ .Permalink }}">{{ .Title }}</a></h1>
  {{ end }}
</body>
</html>
:wq

$
```

Hugo uses the Go template engine. That engine scans the template files for commands which are enclosed between "{{" and "}}" (affectionately called "moustaches").

In our template, the commands are:

1. range
2. .Permalink
3. .Title
4. end

The `range` command is an iterator. We use it to go through the first ten pages. Every HTML file that Hugo creates is treated as a page, so looping through the list of pages will look at every file that will be created.

`.Permalink` prints the URL to link to an article and `.Title` prints the value of the "title" variable. Hugo pulls it from the front matter in the Markdown file.

The `end` command signals the end of the range iterator. The engine loops back to the top of the iteration when it finds `end.` Everything between `range` and `end` is evaluated each time the engine goes through the iteration. In this template, that would cause the title from the first ten pages to be output as heading level one tags along. Because of the permalink, the heading will link to the actual article.

It's helpful to remember that some variables, like `.Data`, are created before any output files. Hugo loads every content file into the variable and then gives the template a chance to process before creating the HTML files.

Build the web site and then verify the results.

```bash
$ rm -rf public/
$ hugo --verbose
INFO: 2015/11/27 Using config file: /tmp/hugo-0.16/config.toml
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/themes/zafta/static to /tmp/hugo-0.16/public/
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/static/ to /tmp/hugo-0.16/public/
INFO: 2015/11/27 found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
0 draft content
0 future content
2 pages created
0 paginator pages created
0 tags created
0 categories created
in 7 ms

$ find public -type f -name '*.html' | xargs ls -l
-rw-r--r--  1 mdhender  wheel   0 Nov 27 21:18 public/404.html
-rw-r--r--  1 mdhender  wheel   0 Nov 27 21:18 public/article/first/index.html
-rw-r--r--  1 mdhender  wheel   0 Nov 27 21:18 public/article/index.html
-rw-r--r--  1 mdhender  wheel   0 Nov 27 21:18 public/article/second/index.html
-rw-r--r--  1 mdhender  wheel  94 Nov 27 21:18 public/index.html

$ cat public/index.html 
<!DOCTYPE html>
<html>
<body>
  
    <h1><a href="http://replace-this-with-your-hugo-site.com/article/second/">second</a></h1>
  
    <h1><a href="http://replace-this-with-your-hugo-site.com/article/first/">first</a></h1>
  
</body>
</html>
$ 
```

Congratulations, the home page shows the title of the two articles and the links to them. The articles themselves are still empty, but let's take a moment to appreciate what we've done. Your template now generates output dynamically. Believe it or not, by inserting the range command inside of those curly braces, you've learned everything you need to know to build a theme. All that's really left is understanding which template will be used to generate each content file and becoming familiar with the commands for the template engine.

Well, if that were entirely true, this tutorial would be much shorter. There are a few things to know that will make creating a new template much easier. Don't worry, though, that's all to come.

There's also a few things to understand about developing and testing your theme. Notice that the links in the `public/index.html` file use the full `baseurl` from the `config.toml` file. That's because the rendered files are intended to be deployed to your web server. If you're testing your them, you'd run `hugo server` and connect to your browser. That command is smart enough to replace the `baseurl` with `http://localhost:1313` on the fly so that links automatically work for you. That's another reason we recommend testing with the built in server.

### Add Content to the Articles

We're working with articles, which are in the `content/article/` directory. That means that their section (as far as templates are concerned) is "article." If we don't do something weird, their type is also "article."

Hugo uses the section and type to find the template file for every piece of content that it renders. Hugo will first look for a template file that matches the section or type name. If it can't find one, then it will look in the `_default/` directory. There are some twists that we'll cover when we get to categories and tags, but for now we can assume that Hugo will try `article/single.html`, then `_default/single.html`.

Now that we know the search rule, let's see what we actually have available:

```bash
$ find themes/zafta -name single.html | xargs ls -l
-rw-r--r--  1 mdhender  wheel  0 Nov 27 20:35 themes/zafta/layouts/_default/single.html
```

If you look back at the articles that we've rendered, you can see that they're empty because Hugo looked for `article/single.html` but couldn't find it. The `_default/single.html` template is empty, so the rendered article file was empty, too.

We could create a new template, `article/single.html`, or change the default, `_default/single.html`. Since we don't know of any other content types, let's start with updating the default.

We mentioned earlier that you should always change the most specific template first to avoid accidentally change other content. We're breaking that rule intentionally just to explore how the default is used.

Remember, any content that we haven't created a template for will end up using this template. That can be good or bad. Bad because I know that we're going to be adding different types of content and we're going to end up undoing some of the changes we've made. It's good because we'll be able to see immediate results. It's also good to start here because we can start to build the basic layout for the site. As we add more content types, we'll refactor this file and move logic around. Hugo makes that fairly painless, so we'll accept the cost and proceed.

Please see the Hugo documentation on template rendering for all the details on determining which template to use. And, as the docs mention, if you're building a single page application (SPA) web site, you can delete all of the other templates and work with just the default single page. That's a refreshing amount of joy right there.

#### Update the Template File

```bash
$ vi themes/zafta/layouts/_default/single.html
<!DOCTYPE html>
<html>
<head>
  <title>{{ .Title }}</title>
</head>
<body>
  <h1>{{ .Title }}</h1>
  <h2>{{ .Date.Format "Mon, Jan 2, 2006" }}</h2>
  {{ .Content }}
  <p><a href="{{ .Site.BaseURL }}">Home</a></p>
</body>
</html>
:wq

$
```

Build the web site and verify the results.

```bash
$ rm -rf public/
$ hugo --verbose
INFO: 2015/11/27 Using config file: /tmp/hugo-0.16/config.toml
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/themes/zafta/static to /tmp/hugo-0.16/public/
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/static/ to /tmp/hugo-0.16/public/
INFO: 2015/11/27 found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
0 draft content
0 future content
2 pages created
0 paginator pages created
0 tags created
0 categories created
in 7 ms

$ find public -type f -name '*.html' | xargs ls -l
-rw-r--r--  1 mdhender  wheel    0 Nov 27 21:39 public/404.html
-rw-r--r--  1 mdhender  wheel  472 Nov 27 21:39 public/article/first/index.html
-rw-r--r--  1 mdhender  wheel    0 Nov 27 21:39 public/article/index.html
-rw-r--r--  1 mdhender  wheel  513 Nov 27 21:39 public/article/second/index.html
-rw-r--r--  1 mdhender  wheel  241 Nov 27 21:39 public/index.html
```

Note that the we have a "list" file for our articles, `public/article/index.html`. The file is empty because we don't have a template for it, but the other files contain our HTML.

```bash
$ cat public/article/first/index.html
<!DOCTYPE html>
<html>
<head>
  <title>first</title>
</head>
<body>
  <h1>first</h1>
  <h2>Fri, Nov 27, 2015</h2>
  <p>In vel ligula tortor. Aliquam erat volutpat.
Pellentesque at felis eu quam tincidunt dignissim.
Nulla facilisi.</p>

<p>Pellentesque tempus nisi et interdum convallis.
In quam ante, vulputate at massa et, rutrum
gravida dui. Phasellus tristique libero at ex.</p>

  <p><a href="http://replace-this-with-your-hugo-site.com/">Home</a></p>
</body>
</html>

$ cat public/article/second/index.html
<!DOCTYPE html>
<html>
<head>
  <title>second</title>
</head>
<body>
  <h1>second</h1>
  <h2>Fri, Nov 27, 2015</h2>
  <p>Fusce lacus magna, maximus nec sapien eu,
porta efficitur neque. Aliquam erat volutpat.
Vestibulum enim nibh, posuere eu diam nec,
varius sagittis turpis.</p>

<p>Praesent quis sapien egestas mauris accumsan
pulvinar. Ut mattis gravida venenatis. Vivamus
lobortis risus id nisi rutrum, at iaculis.</p>

  <p><a href="http://replace-this-with-your-hugo-site.com/">Home</a></p>
</body>
</html>
$ 
```

Notice that the articles now have content. You can run `hugo server` and use your browser to confirm. You should see a home page with the title of both articles. Each title should link you to the article. There should be a link at the bottom of the article to take you back to the home page.

### Create a List of Articles

We have the articles displaying on the home page and on their own page. We also have the empty file `public/article/index.html` file. Let's make it show a list of all articles (not just the first ten). Key to this is that individual pages use "single" templates. Pages that show collections (or lists) of other pages use "list" templates.

We need to decide which template to update. This will be a listing, so it should be a list template. Let's take a quick look and see which list templates are available.

```bash
$ find themes/zafta -name list.html | xargs ls -l
-rw-r--r--  1 mdhender  wheel  0 Nov 27 20:35 themes/zafta/layouts/_default/list.html
```

As with the single article, we have to decide to update `_default/list.html` or create `section/article.html`. We still don't have multiple content types, so let's stay consistent and update the default list template.

```bash
$ vi themes/zafta/layouts/_default/list.html
<!DOCTYPE html>
<html>
<body>
  {{ range first 10 .Data.Pages }}
    <h1><a href="{{ .Permalink }}">{{ .Title }}</a></h1>
  {{ end }}
  <p><a href="{{ .Site.BaseURL }}">Home</a></p>
</body>
</html>
:wq

$
```

Go ahead and render everything again.

```bash
$ rm -rf public
$ hugo --verbose
INFO: 2015/11/27 Using config file: /tmp/hugo-0.16/config.toml
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/themes/zafta/static to /tmp/hugo-0.16/public/
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/static/ to /tmp/hugo-0.16/public/
INFO: 2015/11/27 found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
0 draft content
0 future content
2 pages created
0 paginator pages created
0 tags created
0 categories created
in 7 ms

$ find public -type f -name '*.html' | xargs ls -l
-rw-r--r--  1 mdhender  wheel    0 Nov 27 21:56 public/404.html
-rw-r--r--  1 mdhender  wheel  472 Nov 27 21:56 public/article/first/index.html
-rw-r--r--  1 mdhender  wheel  241 Nov 27 21:56 public/article/index.html
-rw-r--r--  1 mdhender  wheel  513 Nov 27 21:56 public/article/second/index.html
-rw-r--r--  1 mdhender  wheel  241 Nov 27 21:56 public/index.html

$ 
```

We now have a list of articles. You can start `hugo server` and use your browser to confirm.


## Creating Top Level Pages

Let's add an "about" page and display it at the top level (as opposed to a sub-level like we did with articles).

The default in Hugo is to use the directory structure of the `content/` directory to guide the location of the generated html in the `public/` directory. Let's verify that by creating an "about" page at the top level:

```bash
$ hugo new about.md
/tmp/hugo-0.16/content/about.md created
$ ls -l content/
total 8
drwxr-xr-x   4 mdhender  wheel  136 Nov 27 22:01 .
drwxr-xr-x  10 mdhender  wheel  340 Nov 27 21:56 ..
-rw-r--r--   1 mdhender  wheel   61 Nov 27 22:01 about.md
drwxr-xr-x   4 mdhender  wheel  136 Nov 27 21:11 article

$ vi content/about.md
+++
date = "2015-11-27T22:01:00-06:00"
title = "about"
+++
Neque porro quisquam est qui dolorem
ipsum quia dolor sit amet consectetur
adipisci velit.
$ 
:wq
```

Render the web site and verify the results.

```bash
$ rm -rf public/
$ hugo --verbose
INFO: 2015/11/27 Using config file: /tmp/hugo-0.16/config.toml
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/themes/zafta/static to /tmp/hugo-0.16/public/
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/static/ to /tmp/hugo-0.16/public/
INFO: 2015/11/27 found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
0 draft content
0 future content
3 pages created
0 paginator pages created
0 tags created
0 categories created
in 9 ms

$ ls -l public/
total 24
-rw-r--r--  1 mdhender  wheel     0 Nov 27 22:04 404.html
drwxr-xr-x  3 mdhender  wheel   102 Nov 27 22:04 about
drwxr-xr-x  6 mdhender  wheel   204 Nov 27 22:04 article
drwxr-xr-x  2 mdhender  wheel    68 Nov 27 20:35 css
-rw-r--r--  1 mdhender  wheel   328 Nov 27 22:04 index.html
-rw-r--r--  1 mdhender  wheel  2221 Nov 27 22:04 index.xml
drwxr-xr-x  2 mdhender  wheel    68 Nov 27 20:35 js
-rw-r--r--  1 mdhender  wheel   708 Nov 27 22:04 sitemap.xml

$ ls -l public/about/
total 8
-rw-r--r--  1 mdhender  wheel  304 Nov 27 22:04 index.html

$ cat public/about/index.html 
<!DOCTYPE html>
<html>
<head>
  <title>about</title>
</head>
<body>
  <h1>about</h1>
  <h2>Fri, Nov 27, 2015</h2>
  <p>Neque porro quisquam est qui dolorem
ipsum quia dolor sit amet consectetur
adipisci velit.</p>

  <p><a href="http://replace-this-with-your-hugo-site.com/">Home</a></p>
</body>
</html>
$ 
```

Notice that the page wasn't created at the top level. It was created in a sub-directory named 'about/'. That name came from the name of our Markdown file, `about.md`.

One other thing. Take a look at the home page.

```bash
$ cat public/index.html 
<!DOCTYPE html>
<html>
<body>
  
    <h1><a href="http://replace-this-with-your-hugo-site.com/about/">about</a></h1>
  
    <h1><a href="http://replace-this-with-your-hugo-site.com/article/second/">second</a></h1>
  
    <h1><a href="http://replace-this-with-your-hugo-site.com/article/first/">first</a></h1>
  
</body>
</html>
$ 
```

Notice that the "about" link is listed with the articles? That's not desirable, so let's change that first.

```bash
$ vi themes/zafta/layouts/index.html
<!DOCTYPE html>
<html>
<body>
  <h1>Articles</h1>
  {{ range first 10 .Data.Pages }}
    {{ if eq .Type "article"}}
      <h2><a href="{{ .Permalink }}">{{ .Title }}</a></h2>
    {{ end }}
  {{ end }}

  <h1>Pages</h1>
  {{ range .Data.Pages }}
    {{ if eq .Type "page" }}
      <h2><a href="{{ .Permalink }}">{{ .Title }}</a></h2>
    {{ end }}
  {{ end }}

</body>
</html>
:wq
```

Render the web site and verify the results.

```bash
$ rm -rf public/
$ hugo --verbose
INFO: 2015/11/27 Using config file: /tmp/hugo-0.16/config.toml
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/themes/zafta/static to /tmp/hugo-0.16/public/
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/static/ to /tmp/hugo-0.16/public/
INFO: 2015/11/27 found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
0 draft content
0 future content
3 pages created
0 paginator pages created
0 tags created
0 categories created
in 9 ms

$ cat public/index.html 
<!DOCTYPE html>
<html>
<body>
  <h1>Articles</h1>
      <h2><a href="http://replace-this-with-your-hugo-site.com/article/second/">second</a></h2>
      <h2><a href="http://replace-this-with-your-hugo-site.com/article/first/">first</a></h2>
  <h1>Pages</h1>
      <h2><a href="http://replace-this-with-your-hugo-site.com/about/">about</a></h2>
</body>
</html>
```

The home page has two sections, Articles and Pages, and each section has the right set of headings and links in it.

## Sharing Templates

If you've been following along, you probably noticed that articles have titles in the browser and the home page doesn't. That's because we didn't put the title in the homepage's template (`layouts/index.html`). That's an easy thing to do, but let's look at a better option.

We can put the common bits into a shared template that's stored in the `themes/zafta/layouts/partials/` directory.

### Create the Header and Footer Partials

In Hugo, a partial is a template that's intended to be used within other templates. We're going to use partials to create a single header template that other templates will use. That gives us one place to maintain the header information, which makes maintenance much easier. So much easier, in fact, that we'll jump in and do the same for the footer, too.

```bash
$ vi themes/zafta/layouts/partials/header.html
<!DOCTYPE html>
<html>
<head>
  <title>{{ .Title }}</title>
</head>
<body>
:wq

$ vi themes/zafta/layouts/partials/footer.html
</body>
</html>
:wq
```

### Update the Home Page Template to Use the Partials

The most noticeable difference between a template call and a partials call is the lack of path:

```bash
{{ template "theme/partials/header.html" . }}
```

versus

```bash
{{ partial "header.html" . }}
```

Both pass in the context (that's the period just before the closing moustache).

Let's change the home page template to use these new partials we just created.

```bash
$ vi themes/zafta/layouts/index.html
{{ partial "header.html" . }}

  <h1>Articles</h1>
  {{ range first 10 .Data.Pages }}
    {{ if eq .Type "article"}}
      <h2><a href="{{ .Permalink }}">{{ .Title }}</a></h2>
    {{ end }}
  {{ end }}

  <h1>Pages</h1>
  {{ range .Data.Pages }}
    {{ if eq .Type "page" }}
      <h2><a href="{{ .Permalink }}">{{ .Title }}</a></h2>
    {{ end }}
  {{ end }}

{{ partial "footer.html" . }}
:wq
```

Render the web site and verify the results. The title on the home page is now "My New Hugo Site", which comes from the "title" variable in the `config.toml` file.

### Update the Default Templates to Use the Partials

```bash
$ vi themes/zafta/layouts/_default/single.html
{{ partial "header.html" . }}
  <h1>{{ .Title }}</h1>
  <h2>{{ .Date.Format "Mon, Jan 2, 2006" }}</h2>
  {{ .Content }}
  <p><a href="{{ .Site.BaseURL }}">Home</a></p>
{{ partial "footer.html" . }}
:wq

$ vi themes/zafta/layouts/_default/list.html
{{ partial "header.html" . }}
  {{ range first 10 .Data.Pages }}
    <h1><a href="{{ .Permalink }}">{{ .Title }}</a></h1>
  {{ end }}
  <p><a href="{{ .Site.BaseURL }}">Home</a></p>
{{ partial "footer.html" . }}
:wq
```

Generate the web site and verify the results. The title on the articles and the about page should both reflect the value in the markdown file.

## Addressing the "Date Published" on the About page

It's common to have articles display the date that they were written or published, so let's add that. The front matter of our articles has a variable named "date." It's usually the date the content was created, but let's pretend that's the value we want to display.

We display it by printing the formatted date in the template.

```bash
{{ .Date.Format "Mon, Jan 2, 2006" }}
```

Articles use the `_default/single.html` template, which includes this, so they show a nice looking date. Unfortunately, the "about" page uses the same default template, so it shows the date, too.

There are a couple of ways to make the date display only for articles. We could do an "if" statement to only display the date when the type equals "article." That would work, and is acceptable for sites that have just a couple of content types. It aligns with the principle of "code for today," too.

Let's assume, though, that we've made our site so complex that we feel we have to create a new template type. In Hugo-speak, we're going to create a section template for our articles.

Let's restore the default single template before we forget.

```bash
$ vi themes/zafta/layouts/_default/single.html
{{ partial "header.html" . }}

  <h1>{{ .Title }}</h1>
  {{ .Content }}

{{ partial "footer.html" . }}
:wq
```

Now we'll update the articles's version of the single template. If you remember Hugo's rules, the template engine will use this version over the default.

```bash
$ vi themes/zafta/layouts/_default/single.html
{{ partial "header.html" . }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
  <p><a href="{{ .Site.BaseURL }}">Home</a></p>
{{ partial "footer.html" . }}
:wq
```

Now let's create the section template. First step is to create the directory for the new section. Then we just create a "single" template in it.

```bash
$ mkdir themes/zafta/layouts/article
$ vi themes/zafta/layouts/article/single.html
{{ partial "header.html" . }}
  <h1>{{ .Title }}</h1>
  <h2>{{ .Date.Format "Mon, Jan 2, 2006" }}</h2>
  {{ .Content }}
  <p><a href="{{ .Site.BaseURL }}">Home</a></p>
{{ partial "footer.html" . }}
$ 
:wq

```

Note that we removed the date logic from the default template and put it in the "single" template for `layouts/article/`.

Render the site and verify the results. Articles have dates and the about page doesn't.

```bash
$ rm -rf public/
$ hugo --verbose
INFO: 2015/11/27 Using config file: /tmp/hugo-0.16/config.toml
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/themes/zafta/static to /tmp/hugo-0.16/public/
INFO: 2015/11/27 syncing from /tmp/hugo-0.16/static/ to /tmp/hugo-0.16/public/
INFO: 2015/11/27 found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
0 draft content
0 future content
3 pages created
0 paginator pages created
0 tags created
0 categories created
in 10 ms

$ cat public/about/index.html 
<!DOCTYPE html>
<html>
<head>
  <title>about</title>
</head>
<body>

  <h1>about</h1>
  <p>Neque porro quisquam est qui dolorem
ipsum quia dolor sit amet consectetur
adipisci velit.</p>

  <p><a href="http://replace-this-with-your-hugo-site.com/">Home</a></p>
</body>
</html>

$ cat public/article/first/index.html 
<!DOCTYPE html>
<html>
<head>
  <title>first</title>
</head>
<body>

  <h1>first</h1>
  <h2>Fri, Nov 27, 2015</h2>
  <p>In vel ligula tortor. Aliquam erat volutpat.
Pellentesque at felis eu quam tincidunt dignissim.
Nulla facilisi.</p>

<p>Pellentesque tempus nisi et interdum convallis.
In quam ante, vulputate at massa et, rutrum
gravida dui. Phasellus tristique libero at ex.</p>

  <p><a href="http://replace-this-with-your-hugo-site.com/">Home</a></p>
</body>
</html>

$ 
```

### Don't Repeat Yourself

DRY is a good design goal and Hugo does a great job supporting it. Part of the art of a good template is knowing when to add a new template and when to update an existing one. While you're figuring that out, accept that you'll be doing some refactoring. Hugo makes that easy and fast, so it's okay to delay splitting up a template.
