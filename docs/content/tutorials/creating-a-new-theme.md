---
author: "Michael Henderson"
date: 2014-09-28
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

This tutorial will show you how to create a simple theme in Hugo. I assume that you are familiar with HTML, the bash command line, and that you are comfortable using Markdown to format content. I'll explain how Hugo uses templates and how you can organize your templates to create a theme. I won't cover using CSS to style your theme.

We'll start with creating a new site with a very basic template. Then we'll add in a few pages and posts. With small variations on that, you will be able to create many different types of web sites.

In this tutorial, commands that you enter will start with the "$" prompt. The output will follow. Lines that start with "#" are comments that I've added to explain a point. When I show updates to a file, the ":wq" on the last line means to save the file.

Here's an example:

```
## this is a comment
$ echo this is a command
this is a command

## edit the file
$vi foo.md
+++
date = "2014-09-28"
title = "creating a new theme"
+++

bah and humbug
:wq

## show it
$ cat foo.md
+++
date = "2014-09-28"
title = "creating a new theme"
+++

bah and humbug
$
```


## Some Definitions

There are a few concepts that you need to understand before creating a theme.

### Skins

Skins are the files responsible for the look and feel of your site. It’s the CSS that controls colors and fonts, it’s the Javascript that determines actions and reactions. It’s also the rules that Hugo uses to transform your content into the HTML that the site will serve to visitors.

You have two ways to create a skin. The simplest way is to create it in the ```layouts/``` directory. If you do, then you don’t have to worry about configuring Hugo to recognize it. The first place that Hugo will look for rules and files is in the ```layouts/``` directory so it will always find the skin.

Your second choice is to create it in a sub-directory of the ```themes/``` directory. If you do, then you must always tell Hugo where to search for the skin. It’s extra work, though, so why bother with it?

The difference between creating a skin in ```layouts/``` and creating it in ```themes/``` is very subtle. A skin in ```layouts/``` can’t be customized without updating the templates and static files that it is built from. A skin created in ```themes/```, on the other hand, can be and that makes it easier for other people to use it.

The rest of this tutorial will call a skin created in the ```themes/``` directory a theme.

Note that you can use this tutorial to create a skin in the ```layouts/``` directory if you wish to. The main difference will be that you won’t need to update the site’s configuration file to use a theme.

### The Home Page

The home page, or landing page, is the first page that many visitors to a site see. It is the index.html file in the root directory of the web site. Since Hugo writes files to the public/ directory, our home page is public/index.html.

### Site Configuration File

When Hugo runs, it looks for a configuration file that contains settings that override default values for the entire site. The file can use TOML, YAML, or JSON. I prefer to use TOML for my configuration files. If you prefer to use JSON or YAML, you’ll need to translate my examples. You’ll also need to change the name of the file since Hugo uses the extension to determine how to process it.

Hugo translates Markdown files into HTML. By default, Hugo expects to find Markdown files in your ```content/``` directory and template files in your ```themes/``` directory. It will create HTML files in your ```public/``` directory. You can change this by specifying alternate locations in the configuration file.

### Content

Content is stored in text files that contain two sections. The first section is the “front matter,” which is the meta-information on the content. The second section contains Markdown that will be converted to HTML.

#### Front Matter

The front matter is information about the content. Like the configuration file, it can be written in TOML, YAML, or JSON. Unlike the configuration file, Hugo doesn’t use the file’s extension to know the format. It looks for markers to signal the type. TOML is surrounded by “`+++`”, YAML by “`---`”, and JSON is enclosed in curly braces. I prefer to use TOML, so you’ll need to translate my examples if you prefer YAML or JSON.

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

A partial template is a template that can be included in other templates. Partial templates must be called using the “partial” template command. They are very handy for rolling up common behavior. For example, your site may have a banner that all pages use. Instead of copying the text of the banner into every single and list template, you could create a partial with the banner in it. That way if you decide to change the banner, you only have to change the partial template.

## Create a New Site

Let's use Hugo to create a new web site. I'm a Mac user, so I'll create mine in my home directory, in the Sites folder. If you're using Linux, you might have to create the folder first.

The "new site" command will create a skeleton of a site. It will give you the basic directory structure and a useable configuration file.

```
$ hugo new site ~/Sites/zafta
$ cd ~/Sites/zafta
$ ls -l
total 8
drwxr-xr-x  7 quoha  staff  238 Sep 29 16:49 .
drwxr-xr-x  3 quoha  staff  102 Sep 29 16:49 ..
drwxr-xr-x  2 quoha  staff   68 Sep 29 16:49 archetypes
-rw-r--r--  1 quoha  staff   82 Sep 29 16:49 config.toml
drwxr-xr-x  2 quoha  staff   68 Sep 29 16:49 content
drwxr-xr-x  2 quoha  staff   68 Sep 29 16:49 layouts
drwxr-xr-x  2 quoha  staff   68 Sep 29 16:49 static
$
```

Take a look in the content/ directory to confirm that it is empty.

The other directories (archetypes/, layouts/, and static/) are used when customizing a theme. That's a topic for a different tutorial, so please ignore them for now.

### Generate the HTML For the New Site

Running the `hugo` command with no options will read all the available content and generate the HTML files. It will also copy all static files (that's everything that's not content). Since we have an empty site, it won't do much, but it will do it very quickly.

```
$ hugo --verbose
INFO: 2014/09/29 Using config file: config.toml
INFO: 2014/09/29 syncing from /Users/quoha/Sites/zafta/static/ to /Users/quoha/Sites/zafta/public/
WARN: 2014/09/29 Unable to locate layout: [index.html _default/list.html _default/single.html]
WARN: 2014/09/29 Unable to locate layout: [404.html]
0 draft content 
0 future content 
0 pages created 
0 tags created
0 categories created
in 2 ms
$ 
```

The "`--verbose`" flag gives extra information that will be helpful when we build the template. Every line of the output that starts with "INFO:" or "WARN:" is present because we used that flag. The lines that start with "WARN:" are warning messages. We'll go over them later.

We can verify that the command worked by looking at the directory again.

```
$ ls -l
total 8
drwxr-xr-x  2 quoha  staff   68 Sep 29 16:49 archetypes
-rw-r--r--  1 quoha  staff   82 Sep 29 16:49 config.toml
drwxr-xr-x  2 quoha  staff   68 Sep 29 16:49 content
drwxr-xr-x  2 quoha  staff   68 Sep 29 16:49 layouts
drwxr-xr-x  4 quoha  staff  136 Sep 29 17:02 public
drwxr-xr-x  2 quoha  staff   68 Sep 29 16:49 static
$
```

See that new public/ directory? Hugo placed all generated content there. When you're ready to publish your web site, that's the place to start. For now, though, let's just confirm that we have what we'd expect from a site with no content.

```
$ ls -l public
total 16
-rw-r--r--  1 quoha  staff  416 Sep 29 17:02 index.xml
-rw-r--r--  1 quoha  staff  262 Sep 29 17:02 sitemap.xml
$ 
```

Hugo created two XML files, which is standard, but there are no HTML files.



### Test the New Site

Verify that you can run the built-in web server. It will dramatically shorten your development cycle if you do. Start it by running the "server" command. If it is successful, you will see output similar to the following:

```
$ hugo server --verbose
INFO: 2014/09/29 Using config file: /Users/quoha/Sites/zafta/config.toml
INFO: 2014/09/29 syncing from /Users/quoha/Sites/zafta/static/ to /Users/quoha/Sites/zafta/public/
WARN: 2014/09/29 Unable to locate layout: [index.html _default/list.html _default/single.html]
WARN: 2014/09/29 Unable to locate layout: [404.html]
0 draft content 
0 future content 
0 pages created 
0 tags created
0 categories created
in 2 ms
Serving pages from /Users/quoha/Sites/zafta/public
Web Server is available at http://localhost:1313
Press Ctrl+C to stop
```

Connect to the listed URL (it's on the line that starts with "Web Server"). If everything is working correctly, you should get a page that shows the following:

```
index.xml
sitemap.xml
```

That's a listing of your public/ directory. Hugo didn't create a home page because our site has no content. When there's no index.html file in a directory, the server lists the files in the directory, which is what you should see in your browser.

Let’s go back and look at those warnings again.

```
WARN: 2014/09/29 Unable to locate layout: [index.html _default/list.html _default/single.html]
WARN: 2014/09/29 Unable to locate layout: [404.html]
```

That second warning is easier to explain. We haven’t created a template to be used to generate “page not found errors.” The 404 message is a topic for a separate tutorial.

Now for the first warning. It is for the home page. You can tell because the first layout that it looked for was “index.html.” That’s only used by the home page.

I like that the verbose flag causes Hugo to list the files that it's searching for. For the home page, they are index.html, _default/list.html, and _default/single.html. There are some rules that we'll cover later that explain the names and paths. For now, just remember that Hugo couldn't find a template for the home page and it told you so.

At this point, you've got a working installation and site that we can build upon. All that’s left is to add some content and a theme to display it.

## Create a New Theme

Hugo doesn't ship with a default theme. There are a few available (I counted a dozen when I first installed Hugo) and Hugo comes with a command to create new themes.

We're going to create a new theme called "zafta." Since the goal of this tutorial is to show you how to fill out the files to pull in your content, the theme will not contain any CSS. In other words, ugly but functional.

All themes have opinions on content and layout. For example, Zafta uses "post" over "blog". Strong opinions make for simpler templates but differing opinions make it tougher to use themes. When you build a theme, consider using the terms that other themes do.


### Create a Skeleton

Use the hugo "new" command to create the skeleton of a theme. This creates the directory structure and places empty files for you to fill out.

```
$ hugo new theme zafta

$ ls -l
total 8
drwxr-xr-x  2 quoha  staff   68 Sep 29 16:49 archetypes
-rw-r--r--  1 quoha  staff   82 Sep 29 16:49 config.toml
drwxr-xr-x  2 quoha  staff   68 Sep 29 16:49 content
drwxr-xr-x  2 quoha  staff   68 Sep 29 16:49 layouts
drwxr-xr-x  4 quoha  staff  136 Sep 29 17:02 public
drwxr-xr-x  2 quoha  staff   68 Sep 29 16:49 static
drwxr-xr-x  3 quoha  staff  102 Sep 29 17:31 themes

$ find themes -type f | xargs ls -l
-rw-r--r--  1 quoha  staff  1081 Sep 29 17:31 themes/zafta/LICENSE.md
-rw-r--r--  1 quoha  staff     0 Sep 29 17:31 themes/zafta/archetypes/default.md
-rw-r--r--  1 quoha  staff     0 Sep 29 17:31 themes/zafta/layouts/_default/list.html
-rw-r--r--  1 quoha  staff     0 Sep 29 17:31 themes/zafta/layouts/_default/single.html
-rw-r--r--  1 quoha  staff     0 Sep 29 17:31 themes/zafta/layouts/index.html
-rw-r--r--  1 quoha  staff     0 Sep 29 17:31 themes/zafta/layouts/partials/footer.html
-rw-r--r--  1 quoha  staff     0 Sep 29 17:31 themes/zafta/layouts/partials/header.html
-rw-r--r--  1 quoha  staff    93 Sep 29 17:31 themes/zafta/theme.toml
$ 
```

The skeleton includes templates (the files ending in .html), license file, a description of your theme (the theme.toml file), and an empty archetype.

Please take a minute to fill out the theme.toml and LICENSE.md files. They're optional, but if you're going to be distributing your theme, it tells the world who to praise (or blame). It's also nice to declare the license so that people will know how they can use the theme.

```
$ vi themes/zafta/theme.toml
author = "michael d henderson"
description = "a minimal working template"
license = "MIT"
name = "zafta"
source_repo = ""
tags = ["tags", "categories"]
:wq

## also edit themes/zafta/LICENSE.md and change
## the bit that says "YOUR_NAME_HERE"
```

Note that the the skeleton's template files are empty. Don't worry, we'll be changing that shortly.

```
$ find themes/zafta -name '*.html' | xargs ls -l
-rw-r--r--  1 quoha  staff  0 Sep 29 17:31 themes/zafta/layouts/_default/list.html
-rw-r--r--  1 quoha  staff  0 Sep 29 17:31 themes/zafta/layouts/_default/single.html
-rw-r--r--  1 quoha  staff  0 Sep 29 17:31 themes/zafta/layouts/index.html
-rw-r--r--  1 quoha  staff  0 Sep 29 17:31 themes/zafta/layouts/partials/footer.html
-rw-r--r--  1 quoha  staff  0 Sep 29 17:31 themes/zafta/layouts/partials/header.html
$
```



### Update the Configuration File to Use the Theme

Now that we've got a theme to work with, it's a good idea to add the theme name to the configuration file. This is optional, because you can always add "-t zafta" on all your commands. I like to put it the configuration file because I like shorter command lines. If you don't put it in the configuration file or specify it on the command line, you won't use the template that you're expecting to.

Edit the file to add the theme, add a title for the site, and specify that all of our content will use the TOML format.

```
$ vi config.toml
theme = "zafta"
baseurl = ""
languageCode = "en-us"
title = "zafta - totally refreshing"
MetaDataFormat = "toml"
:wq

$
```

### Generate the Site

Now that we have an empty theme, let's generate the site again.

```
$ hugo --verbose
INFO: 2014/09/29 Using config file: /Users/quoha/Sites/zafta/config.toml
INFO: 2014/09/29 syncing from /Users/quoha/Sites/zafta/themes/zafta/static/ to /Users/quoha/Sites/zafta/public/
INFO: 2014/09/29 syncing from /Users/quoha/Sites/zafta/static/ to /Users/quoha/Sites/zafta/public/
WARN: 2014/09/29 Unable to locate layout: [404.html theme/404.html]
0 draft content 
0 future content 
0 pages created 
0 tags created
0 categories created
in 2 ms
$
```

Did you notice that the output is different? The warning message for the home page has disappeared and we have an additional information line saying that Hugo is syncing from the theme's directory.

Let's check the public/ directory to see what Hugo's created.

```
$ ls -l public
total 16
drwxr-xr-x  2 quoha  staff   68 Sep 29 17:56 css
-rw-r--r--  1 quoha  staff    0 Sep 29 17:56 index.html
-rw-r--r--  1 quoha  staff  407 Sep 29 17:56 index.xml
drwxr-xr-x  2 quoha  staff   68 Sep 29 17:56 js
-rw-r--r--  1 quoha  staff  243 Sep 29 17:56 sitemap.xml
$
```

Notice four things:

1. Hugo created a home page. This is the file public/index.html.
2. Hugo created a css/ directory.
3. Hugo created a js/ directory.
4. Hugo claimed that it created 0 pages. It created a file and copied over static files, but didn't create any pages. That's because it considers a "page" to be a file created directly from a content file. It doesn't count things like the index.html files that it creates automatically.

#### The Home Page

Hugo supports many different types of templates. The home page is special because it gets its own type of template and its own template file. The file, layouts/index.html, is used to generate the HTML for the home page. The Hugo documentation says that this is the only required template, but that depends. Hugo's warning message shows that it looks for three different templates:

```
WARN: 2014/09/29 Unable to locate layout: [index.html _default/list.html _default/single.html]
```

If it can't find any of these, it completely skips creating the home page. We noticed that when we built the site without having a theme installed.

When Hugo created our theme, it created an empty home page template. Now, when we build the site, Hugo finds the template and uses it to generate the HTML for the home page. Since the template file is empty, the HTML file is empty, too. If the template had any rules in it, then Hugo would have used them to generate the home page.

```
$ find . -name index.html | xargs ls -l
-rw-r--r--  1 quoha  staff  0 Sep 29 20:21 ./public/index.html
-rw-r--r--  1 quoha  staff  0 Sep 29 17:31 ./themes/zafta/layouts/index.html
$ 
```

#### The Magic of Static

Hugo does two things when generating the site. It uses templates to transform content into HTML and it copies static files into the site. Unlike content, static files are not transformed. They are copied exactly as they are.

Hugo assumes that your site will use both CSS and JavaScript, so it creates directories in your theme to hold them. Remember opinions? Well, Hugo's opinion is that you'll store your CSS in a directory named css/ and your JavaScript in a directory named js/. If you don't like that, you can change the directory names in your theme directory or even delete them completely. Hugo's nice enough to offer its opinion, then behave nicely if you disagree.

```
$ find themes/zafta -type d | xargs ls -ld
drwxr-xr-x  7 quoha  staff  238 Sep 29 17:38 themes/zafta
drwxr-xr-x  3 quoha  staff  102 Sep 29 17:31 themes/zafta/archetypes
drwxr-xr-x  5 quoha  staff  170 Sep 29 17:31 themes/zafta/layouts
drwxr-xr-x  4 quoha  staff  136 Sep 29 17:31 themes/zafta/layouts/_default
drwxr-xr-x  4 quoha  staff  136 Sep 29 17:31 themes/zafta/layouts/partials
drwxr-xr-x  4 quoha  staff  136 Sep 29 17:31 themes/zafta/static
drwxr-xr-x  2 quoha  staff   68 Sep 29 17:31 themes/zafta/static/css
drwxr-xr-x  2 quoha  staff   68 Sep 29 17:31 themes/zafta/static/js
$ 
```

## The Theme Development Cycle

When you're working on a theme, you will make changes in the theme's directory, rebuild the site, and check your changes in the browser. Hugo makes this very easy:

1. Purge the public/ directory.
2. Run the built in web server in watch mode.
3. Open your site in a browser.
4. Update the theme.
5. Glance at your browser window to see changes.
6. Return to step 4.

I’ll throw in one more opinion: never work on a theme on a live site. Always work on a copy of your site. Make changes to your theme, test them, then copy them up to your site. For added safety, use a tool like Git to keep a revision history of your content and your theme. Believe me when I say that it is too easy to lose both your mind and your changes.

Check the main Hugo site for information on using Git with Hugo.

### Purge the public/ Directory

When generating the site, Hugo will create new files and update existing ones in the ```public/``` directory. It will not delete files that are no longer used. For example, files that were created in the wrong directory or with the wrong title will remain. If you leave them, you might get confused by them later. I recommend cleaning out your site prior to generating it.

Note: If you're building on an SSD, you should ignore this. Churning on a SSD can be costly.

### Hugo's Watch Option

Hugo's "`--watch`" option will monitor the content/ and your theme directories for changes and rebuild the site automatically.

### Live Reload

Hugo's built in web server supports live reload. As pages are saved on the server, the browser is told to refresh the page. Usually, this happens faster than you can say, "Wow, that's totally amazing."

### Development Commands

Use the following commands as the basis for your workflow.

```
## purge old files. hugo will recreate the public directory.
##
$ rm -rf public
##
## run hugo in watch mode
##
$ hugo server --watch --verbose
```

Here's sample output showing Hugo detecting a change to the template for the home page. Once generated, the web browser automatically reloaded the page. I've said this before, it's amazing.


```
$ rm -rf public
$ hugo server --watch --verbose
INFO: 2014/09/29 Using config file: /Users/quoha/Sites/zafta/config.toml
INFO: 2014/09/29 syncing from /Users/quoha/Sites/zafta/themes/zafta/static/ to /Users/quoha/Sites/zafta/public/
INFO: 2014/09/29 syncing from /Users/quoha/Sites/zafta/static/ to /Users/quoha/Sites/zafta/public/
WARN: 2014/09/29 Unable to locate layout: [404.html theme/404.html]
0 draft content 
0 future content 
0 pages created 
0 tags created
0 categories created
in 2 ms
Watching for changes in /Users/quoha/Sites/zafta/content
Serving pages from /Users/quoha/Sites/zafta/public
Web Server is available at http://localhost:1313
Press Ctrl+C to stop
INFO: 2014/09/29 File System Event: ["/Users/quoha/Sites/zafta/themes/zafta/layouts/index.html": MODIFY|ATTRIB]
Change detected, rebuilding site

WARN: 2014/09/29 Unable to locate layout: [404.html theme/404.html]
0 draft content 
0 future content 
0 pages created 
0 tags created
0 categories created
in 1 ms
```

## Update the Home Page Template

The home page is one of a few special pages that Hugo creates automatically. As mentioned earlier, it looks for one of three files in the theme's layout/ directory:

1. index.html
2. _default/list.html
3. _default/single.html

We could update one of the default templates, but a good design decision is to update the most specific template available. That's not a hard and fast rule (in fact, we'll break it a few times in this tutorial), but it is a good generalization.

### Make a Static Home Page

Right now, that page is empty because we don't have any content and we don't have any logic in the template. Let's change that by adding some text to the template.

```
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

```
$ hugo --verbose
INFO: 2014/09/29 Using config file: /Users/quoha/Sites/zafta/config.toml
INFO: 2014/09/29 syncing from /Users/quoha/Sites/zafta/themes/zafta/static/ to /Users/quoha/Sites/zafta/public/
INFO: 2014/09/29 syncing from /Users/quoha/Sites/zafta/static/ to /Users/quoha/Sites/zafta/public/
WARN: 2014/09/29 Unable to locate layout: [404.html theme/404.html]
0 draft content 
0 future content 
0 pages created 
0 tags created
0 categories created
in 2 ms

$ find public -type f -name '*.html' | xargs ls -l
-rw-r--r--  1 quoha  staff  78 Sep 29 21:26 public/index.html

$ cat public/index.html 
<!DOCTYPE html> 
<html> 
<body> 
  <p>hugo says hello!</p> 
</html>
```

#### Live Reload

Note: If you're running the server with the `--watch` option, you'll see different content in the file:

```
$ cat public/index.html 
<!DOCTYPE html> 
<html> 
<body> 
  <p>hugo says hello!</p> 
<script>document.write('<script src="http://' 
        + (location.host || 'localhost').split(':')[0] 
    + ':1313/livereload.js?mindelay=10"></' 
        + 'script>')</script></body> 
</html>
```

When you use `--watch`, the Live Reload script is added by Hugo. Look for live reload in the documentation to see what it does and how to disable it.

### Build a "Dynamic" Home Page

"Dynamic home page?" Hugo's a static web site generator, so this seems an odd thing to say. I mean let's have the home page automatically reflect the content in the site every time Hugo builds it. We'll use iteration in the template to do that.

#### Create New Posts

Now that we have the home page generating static content, let's add some content to the site. We'll display these posts as a list on the home page and on their own page, too.

Hugo has a command to generate a skeleton post, just like it does for sites and themes.

```
$ hugo --verbose new post/first.md
INFO: 2014/09/29 Using config file: /Users/quoha/Sites/zafta/config.toml
INFO: 2014/09/29 attempting to create  post/first.md of post
INFO: 2014/09/29 curpath: /Users/quoha/Sites/zafta/themes/zafta/archetypes/default.md
ERROR: 2014/09/29 Unable to Cast <nil> to map[string]interface{}

$ 
```

That wasn't very nice, was it?

The "new" command uses an archetype to create the post file. Hugo created an empty default archetype file, but that causes an error when there's a theme. For me, the workaround was to create an archetypes file specifically for the post type.

```
$ vi themes/zafta/archetypes/post.md
+++
Description = ""
Tags = []
Categories = []
+++
:wq

$ find themes/zafta/archetypes -type f | xargs ls -l
-rw-r--r--  1 quoha  staff   0 Sep 29 21:53 themes/zafta/archetypes/default.md
-rw-r--r--  1 quoha  staff  51 Sep 29 21:54 themes/zafta/archetypes/post.md

$ hugo --verbose new post/first.md
INFO: 2014/09/29 Using config file: /Users/quoha/Sites/zafta/config.toml
INFO: 2014/09/29 attempting to create  post/first.md of post
INFO: 2014/09/29 curpath: /Users/quoha/Sites/zafta/themes/zafta/archetypes/post.md
INFO: 2014/09/29 creating /Users/quoha/Sites/zafta/content/post/first.md
/Users/quoha/Sites/zafta/content/post/first.md created

$ hugo --verbose new post/second.md
INFO: 2014/09/29 Using config file: /Users/quoha/Sites/zafta/config.toml
INFO: 2014/09/29 attempting to create  post/second.md of post
INFO: 2014/09/29 curpath: /Users/quoha/Sites/zafta/themes/zafta/archetypes/post.md
INFO: 2014/09/29 creating /Users/quoha/Sites/zafta/content/post/second.md
/Users/quoha/Sites/zafta/content/post/second.md created

$ ls -l content/post
total 16
-rw-r--r--  1 quoha  staff  104 Sep 29 21:54 first.md
-rw-r--r--  1 quoha  staff  105 Sep 29 21:57 second.md

$ cat content/post/first.md 
+++
Categories = []
Description = ""
Tags = []
date = "2014-09-29T21:54:53-05:00"
title = "first"

+++
my first post

$ cat content/post/second.md 
+++
Categories = []
Description = ""
Tags = []
date = "2014-09-29T21:57:09-05:00"
title = "second"

+++
my second post

$ 
```

Build the web site and then verify the results.

```
$ rm -rf public
$ hugo --verbose
INFO: 2014/09/29 Using config file: /Users/quoha/Sites/zafta/config.toml
INFO: 2014/09/29 syncing from /Users/quoha/Sites/zafta/themes/zafta/static/ to /Users/quoha/Sites/zafta/public/
INFO: 2014/09/29 syncing from /Users/quoha/Sites/zafta/static/ to /Users/quoha/Sites/zafta/public/
INFO: 2014/09/29 found taxonomies: map[string]string{"category":"categories", "tag":"tags"}
WARN: 2014/09/29 Unable to locate layout: [404.html theme/404.html]
0 draft content 
0 future content 
2 pages created 
0 tags created
0 categories created
in 4 ms
$
```

The output says that it created 2 pages. Those are our new posts:

```
$ find public -type f -name '*.html' | xargs ls -l
-rw-r--r--  1 quoha  staff  78 Sep 29 22:13 public/index.html
-rw-r--r--  1 quoha  staff   0 Sep 29 22:13 public/post/first/index.html
-rw-r--r--  1 quoha  staff   0 Sep 29 22:13 public/post/index.html
-rw-r--r--  1 quoha  staff   0 Sep 29 22:13 public/post/second/index.html
$
```

The new files are empty because because the templates used to generate the content are empty. The homepage doesn't show the new content, either. We have to update the templates to add the posts.

### List and Single Templates

In Hugo, we have three major kinds of templates. There's the home page template that we updated previously. It is used only by the home page. We also have "single" templates which are used to generate output for a single content file. We also have "list" templates that are used to group multiple pieces of content before generating output.

Generally speaking, list templates are named "list.html" and single templates are named "single.html."

There are three other types of templates: partials, content views, and terms. We will not go into much detail on these.

### Add Content to the Homepage

The home page will contain a list of posts. Let's update its template to add the posts that we just created. The logic in the template will run every time we build the site.

```
$ vi themes/zafta/layouts/index.html 
<!DOCTYPE html>
<html>
<body>
  {{ range first 10 .Data.Pages }}
    <h1>{{ .Title }}</h1>
  {{ end }}
</body>
</html>
:wq

$
```

Hugo uses the Go template engine. That engine scans the template files for commands which are enclosed between "{{" and "}}". In our template, the commands are:

1. range
2. .Title
3. end

The "range" command is an iterator. We're going to use it to go through the first ten pages. Every HTML file that Hugo creates is treated as a page, so looping through the list of pages will look at every file that will be created.

The ".Title" command prints the value of the "title" variable. Hugo pulls it from the front matter in the Markdown file.

The "end" command signals the end of the range iterator. The engine loops back to the top of the iteration when it finds "end." Everything between the "range" and "end" is evaluated every time the engine goes through the iteration. In this file, that would cause the title from the first ten pages to be output as heading level one.

It's helpful to remember that some variables, like .Data, are created before any output files. Hugo loads every content file into the variable and then gives the template a chance to process before creating the HTML files.

Build the web site and then verify the results.

```
$ rm -rf public
$ hugo --verbose
INFO: 2014/09/29 Using config file: /Users/quoha/Sites/zafta/config.toml
INFO: 2014/09/29 syncing from /Users/quoha/Sites/zafta/themes/zafta/static/ to /Users/quoha/Sites/zafta/public/
INFO: 2014/09/29 syncing from /Users/quoha/Sites/zafta/static/ to /Users/quoha/Sites/zafta/public/
INFO: 2014/09/29 found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
WARN: 2014/09/29 Unable to locate layout: [404.html theme/404.html]
0 draft content 
0 future content 
2 pages created 
0 tags created
0 categories created
in 4 ms
$ find public -type f -name '*.html' | xargs ls -l 
-rw-r--r--  1 quoha  staff  94 Sep 29 22:23 public/index.html
-rw-r--r--  1 quoha  staff   0 Sep 29 22:23 public/post/first/index.html
-rw-r--r--  1 quoha  staff   0 Sep 29 22:23 public/post/index.html
-rw-r--r--  1 quoha  staff   0 Sep 29 22:23 public/post/second/index.html
$ cat public/index.html 
<!DOCTYPE html>
<html>
<body>
  
    <h1>second</h1>
  
    <h1>first</h1>
  
</body>
</html>
$
```

Congratulations, the home page shows the title of the two posts. The posts themselves are still empty, but let's take a moment to appreciate what we've done. Your template now generates output dynamically. Believe it or not, by inserting the range command inside of those curly braces, you've learned everything you need to know to build a theme. All that's really left is understanding which template will be used to generate each content file and becoming familiar with the commands for the template engine.

And, if that were entirely true, this tutorial would be much shorter. There are a few things to know that will make creating a new template much easier. Don't worry, though, that's all to come.

### Add Content to the Posts

We're working with posts, which are in the content/post/ directory. That means that their section is "post" (and if we don't do something weird, their type is also "post").

Hugo uses the section and type to find the template file for every piece of content. Hugo will first look for a template file that matches the section or type name. If it can't find one, then it will look in the _default/ directory. There are some twists that we'll cover when we get to categories and tags, but for now we can assume that Hugo will try post/single.html, then _default/single.html.

Now that we know the search rule, let's see what we actually have available:

```
$ find themes/zafta -name single.html | xargs ls -l
-rw-r--r--  1 quoha  staff  132 Sep 29 17:31 themes/zafta/layouts/_default/single.html
```

We could create a new template, post/single.html, or change the default. Since we don't know of any other content types, let's start with updating the default.

Remember, any content that we haven't created a template for will end up using this template. That can be good or bad. Bad because I know that we're going to be adding different types of content and we're going to end up undoing some of the changes we've made. It's good because we'll be able to see immediate results. It's also good to start here because we can start to build the basic layout for the site. As we add more content types, we'll refactor this file and move logic around. Hugo makes that fairly painless, so we'll accept the cost and proceed.

Please see the Hugo documentation on template rendering for all the details on determining which template to use. And, as the docs mention, if you're building a single page application (SPA) web site, you can delete all of the other templates and work with just the default single page. That's a refreshing amount of joy right there.

#### Update the Template File

```
$ vi themes/zafta/layouts/_default/single.html 
<!DOCTYPE html>
<html>
<head>
  <title>{{ .Title }}</title>
</head>
<body>
  <h1>{{ .Title }}</h1>
  {{ .Content }}
</body>
</html>
:wq

$
```

Build the web site and verify the results.

```
$ rm -rf public
$ hugo --verbose
INFO: 2014/09/29 Using config file: /Users/quoha/Sites/zafta/config.toml
INFO: 2014/09/29 syncing from /Users/quoha/Sites/zafta/themes/zafta/static/ to /Users/quoha/Sites/zafta/public/
INFO: 2014/09/29 syncing from /Users/quoha/Sites/zafta/static/ to /Users/quoha/Sites/zafta/public/
INFO: 2014/09/29 found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
WARN: 2014/09/29 Unable to locate layout: [404.html theme/404.html]
0 draft content 
0 future content 
2 pages created 
0 tags created
0 categories created
in 4 ms

$ find public -type f -name '*.html' | xargs ls -l
-rw-r--r--  1 quoha  staff   94 Sep 29 22:40 public/index.html
-rw-r--r--  1 quoha  staff  125 Sep 29 22:40 public/post/first/index.html
-rw-r--r--  1 quoha  staff    0 Sep 29 22:40 public/post/index.html
-rw-r--r--  1 quoha  staff  128 Sep 29 22:40 public/post/second/index.html

$ cat public/post/first/index.html 
<!DOCTYPE html>
<html>
<head>
  <title>first</title>
</head>
<body>
  <h1>first</h1>
  <p>my first post</p>

</body>
</html>

$ cat public/post/second/index.html 
<!DOCTYPE html>
<html>
<head>
  <title>second</title>
</head>
<body>
  <h1>second</h1>
  <p>my second post</p>

</body>
</html>
$
```

Notice that the posts now have content. You can go to localhost:1313/post/first to verify.

### Linking to Content

The posts are on the home page. Let's add a link from there to the post. Since this is the home page, we'll update its template.

```
$ vi themes/zafta/layouts/index.html
<!DOCTYPE html>
<html>
<body>
  {{ range first 10 .Data.Pages }}
    <h1><a href="{{ .Permalink }}">{{ .Title }}</a></h1>
  {{ end }}
</body>
</html>
```

Build the web site and verify the results.

```
$ rm -rf public
$ hugo --verbose
INFO: 2014/09/29 Using config file: /Users/quoha/Sites/zafta/config.toml
INFO: 2014/09/29 syncing from /Users/quoha/Sites/zafta/themes/zafta/static/ to /Users/quoha/Sites/zafta/public/
INFO: 2014/09/29 syncing from /Users/quoha/Sites/zafta/static/ to /Users/quoha/Sites/zafta/public/
INFO: 2014/09/29 found taxonomies: map[string]string{"tag":"tags", "category":"categories"}
WARN: 2014/09/29 Unable to locate layout: [404.html theme/404.html]
0 draft content 
0 future content 
2 pages created 
0 tags created
0 categories created
in 4 ms

$ find public -type f -name '*.html' | xargs ls -l
-rw-r--r--  1 quoha  staff  149 Sep 29 22:44 public/index.html
-rw-r--r--  1 quoha  staff  125 Sep 29 22:44 public/post/first/index.html
-rw-r--r--  1 quoha  staff    0 Sep 29 22:44 public/post/index.html
-rw-r--r--  1 quoha  staff  128 Sep 29 22:44 public/post/second/index.html

$ cat public/index.html 
<!DOCTYPE html>
<html>
<body>
  
    <h1><a href="/post/second/">second</a></h1>
  
    <h1><a href="/post/first/">first</a></h1>
  
</body>
</html>

$
```

### Create a Post Listing

We have the posts displaying on the home page and on their own page. We also have a file public/post/index.html that is empty. Let's make it show a list of all posts (not just the first ten).

We need to decide which template to update. This will be a listing, so it should be a list template. Let's take a quick look and see which list templates are available.

```
$ find themes/zafta -name list.html | xargs ls -l
-rw-r--r--  1 quoha  staff  0 Sep 29 17:31 themes/zafta/layouts/_default/list.html
```

As with the single post, we have to decide to update _default/list.html or create post/list.html. We still don't have multiple content types, so let's stay consistent and update the default list template.

## Creating Top Level Pages

Let's add an "about" page and display it at the top level (as opposed to a sub-level like we did with posts).

The default in Hugo is to use the directory structure of the content/ directory to guide the location of the generated html in the public/ directory. Let's verify that by creating an "about" page at the top level:

```
$ vi content/about.md 
+++
title = "about"
description = "about this site"
date = "2014-09-27"
slug = "about time"
+++

## about us

i'm speechless
:wq
```

Generate the web site and verify the results.

```
$ find public -name '*.html' | xargs ls -l
-rw-rw-r--  1 mdhender  staff   334 Sep 27 15:08 public/about-time/index.html
-rw-rw-r--  1 mdhender  staff   527 Sep 27 15:08 public/index.html
-rw-rw-r--  1 mdhender  staff   358 Sep 27 15:08 public/post/first-post/index.html
-rw-rw-r--  1 mdhender  staff     0 Sep 27 15:08 public/post/index.html
-rw-rw-r--  1 mdhender  staff   342 Sep 27 15:08 public/post/second-post/index.html
```

Notice that the page wasn't created at the top level. It was created in a sub-directory named 'about-time/'. That name came from our slug. Hugo will use the slug to name the generated content. It's a reasonable default, by the way, but we can learn a few things by fighting it for this file.

One other thing. Take a look at the home page.

```
$ cat public/index.html
<!DOCTYPE html>
<html>
<body>
    <h1><a href="http://localhost:1313/post/theme/">creating a new theme</a></h1>
    <h1><a href="http://localhost:1313/about-time/">about</a></h1>
    <h1><a href="http://localhost:1313/post/second-post/">second</a></h1>
    <h1><a href="http://localhost:1313/post/first-post/">first</a></h1>
<script>document.write('<script src="http://'
        + (location.host || 'localhost').split(':')[0]
		+ ':1313/livereload.js?mindelay=10"></'
        + 'script>')</script></body>
</html>
```

Notice that the "about" link is listed with the posts? That's not desirable, so let's change that first.

```
$ vi themes/zafta/layouts/index.html
<!DOCTYPE html>
<html>
<body>
  <h1>posts</h1>
  {{ range first 10 .Data.Pages }}
    {{ if eq .Type "post"}}
      <h2><a href="{{ .Permalink }}">{{ .Title }}</a></h2>
    {{ end }}
  {{ end }}

  <h1>pages</h1>
  {{ range .Data.Pages }}
    {{ if eq .Type "page" }}
      <h2><a href="{{ .Permalink }}">{{ .Title }}</a></h2>
    {{ end }}
  {{ end }}
</body>
</html>
:wq
```

Generate the web site and verify the results. The home page has two sections, posts and pages, and each section has the right set of headings and links in it.

But, that about page still renders to about-time/index.html.

```
$ find public -name '*.html' | xargs ls -l
-rw-rw-r--  1 mdhender  staff    334 Sep 27 15:33 public/about-time/index.html
-rw-rw-r--  1 mdhender  staff    645 Sep 27 15:33 public/index.html
-rw-rw-r--  1 mdhender  staff    358 Sep 27 15:33 public/post/first-post/index.html
-rw-rw-r--  1 mdhender  staff      0 Sep 27 15:33 public/post/index.html
-rw-rw-r--  1 mdhender  staff    342 Sep 27 15:33 public/post/second-post/index.html
```

Knowing that hugo is using the slug to generate the file name, the simplest solution is to change the slug. Let's do it the hard way and change the permalink in the configuration file.

```
$ vi config.toml
[permalinks]
	page = "/:title/"
	about = "/:filename/"
```

Generate the web site and verify that this didn't work. Hugo lets "slug" or "URL" override the permalinks setting in the configuration file. Go ahead and comment out the slug in content/about.md, then generate the web site to get it to be created in the right place.

## Sharing Templates

If you've been following along, you probably noticed that posts have titles in the browser and the home page doesn't. That's because we didn't put the title in the home page's template (layouts/index.html). That's an easy thing to do, but let's look at a different option.

We can put the common bits into a shared template that's stored in the themes/zafta/layouts/partials/ directory.

### Create the Header and Footer Partials

In Hugo, a partial is a sugar-coated template. Normally a template reference has a path specified. Partials are different. Hugo searches for them along a TODO defined search path. This makes it easier for end-users to override the theme's presentation.

```
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

```
{{ template "theme/partials/header.html" . }}
```
versus
```
{{ partial "header.html" . }}
```
Both pass in the context.

Let's change the home page template to use these new partials.

```
$ vi themes/zafta/layouts/index.html
{{ partial "header.html" . }}

  <h1>posts</h1>
  {{ range first 10 .Data.Pages }}
    {{ if eq .Type "post"}}
      <h2><a href="{{ .Permalink }}">{{ .Title }}</a></h2>
    {{ end }}
  {{ end }}

  <h1>pages</h1>
  {{ range .Data.Pages }}
    {{ if or (eq .Type "page") (eq .Type "about") }}
      <h2><a href="{{ .Permalink }}">{{ .Type }} - {{ .Title }} - {{ .RelPermalink }}</a></h2>
    {{ end }}
  {{ end }}

{{ partial "footer.html" . }}
:wq
```

Generate the web site and verify the results. The title on the home page is now "your title here", which comes from the "title" variable in the config.toml file.

### Update the Default Single Template to Use the Partials

```
$ vi themes/zafta/layouts/_default/single.html
{{ partial "header.html" . }}

  <h1>{{ .Title }}</h1>
  {{ .Content }}

{{ partial "footer.html" . }}
:wq
```

Generate the web site and verify the results. The title on the posts and the about page should both reflect the value in the markdown file.

## Add “Date Published” to Posts

It's common to have posts display the date that they were written or published, so let's add that. The front matter of our posts has a variable named "date." It's usually the date the content was created, but let's pretend that's the value we want to display.

### Add “Date Published” to the Template

We'll start by updating the template used to render the posts. The template code will look like:

```
{{ .Date.Format "Mon, Jan 2, 2006" }}
```

Posts use the default single template, so we'll change that file.

```
$ vi themes/zafta/layouts/_default/single.html
{{ partial "header.html" . }}

  <h1>{{ .Title }}</h1>
  <h2>{{ .Date.Format "Mon, Jan 2, 2006" }}</h2>
  {{ .Content }}

{{ partial "footer.html" . }}
:wq
```

Generate the web site and verify the results. The posts now have the date displayed in them. There's a problem, though. The "about" page also has the date displayed.

As usual, there are a couple of ways to make the date display only on posts. We could do an "if" statement like we did on the home page. Another way would be to create a separate template for posts.

The "if" solution works for sites that have just a couple of content types. It aligns with the principle of "code for today," too.

Let's assume, though, that we've made our site so complex that we feel we have to create a new template type. In Hugo-speak, we're going to create a section template.

Let's restore the default single template before we forget.

```
$ mkdir themes/zafta/layouts/post
$ vi themes/zafta/layouts/_default/single.html
{{ partial "header.html" . }}

  <h1>{{ .Title }}</h1>
  {{ .Content }}

{{ partial "footer.html" . }}
:wq
```

Now we'll update the post's version of the single template. If you remember Hugo's rules, the template engine will use this version over the default.

```
$ vi themes/zafta/layouts/post/single.html
{{ partial "header.html" . }}

  <h1>{{ .Title }}</h1>
  <h2>{{ .Date.Format "Mon, Jan 2, 2006" }}</h2>
  {{ .Content }}

{{ partial "footer.html" . }}
:wq

```

Note that we removed the date logic from the default template and put it in the post template. Generate the web site and verify the results. Posts have dates and the about page doesn't.

### Don't Repeat Yourself

DRY is a good design goal and Hugo does a great job supporting it. Part of the art of a good template is knowing when to add a new template and when to update an existing one. While you're figuring that out, accept that you'll be doing some refactoring. Hugo makes that easy and fast, so it's okay to delay splitting up a template.
