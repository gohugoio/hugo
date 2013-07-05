# Hugo

A really fast static site generator written in GoLang.

## Overview

Hugo is a static site generator written in GoLang. It is optimized for 
speed, easy use and configurability. Hugo takes a directory with content and
templates and renders them into a full html website.

Hugo makes use of markdown files with front matter for meta data.  

A typical website of moderate size can be 
rendered in a fraction of a second. It is written to work well with any
kind of website including blogs, tumbles and docs. 


# Getting Started

## Installing Hugo

Hugo is written in GoLang with support for Windows, Linux, FreeBSD and OSX.

The latest release can be found at [hugo releases](https://github.com/spf13/hugo/releases)
We currently build for Windows, Linux, FreeBSD and OS X for x64
and 386 architectures. 

Installation is very easy. Simply download the appropriate version for your
platform. Once downloaded it can be run from anywhere. You don't need to install
it into a global location. This works well for shared hosts and other systems
where you don't have a privileged account.

Ideally you should install it somewhere in your path for easy use. `/usr/local/bin` 
is the most probable location.

*Hugo has no external dependencies.*

## Installing from source

Make sure you have a recent version of go installed. Hugo requires go 1.1+.

    git clone https://github.com/spf13/hugo
    cd hugo
    go build -o hugo main.go


## Source Directory Organization

Hugo takes a single directory and uses it as the input for creating a complete website.

Hugo has a very small amount of configuration, while remaining highly customizable. 
It accomplishes by assuming that you will only provide templates with the intent of
using them.

An example directory may look like:

    .
    ├── config.json
    ├── content
    |   ├── post
    |   |   ├── firstpost.md
    |   |   └── secondpost.md
    |   └── quote
    |   |   ├── first.md
    |   |   └── second.md
    ├── layouts
    |   ├── chrome
    |   |   ├── header.html
    |   |   └── footer.html
    |   ├── indexes
    |   |   ├── category.html
    |   |   ├── post.html
    |   |   ├── quote.html
    |   |   └── tag.html
    |   ├── post
    |   |   ├── li.html
    |   |   ├── single.html
    |   |   └── summary.html
    |   ├── quote
    |   |   ├── li.html
    |   |   ├── single.html
    |   |   └── summary.html
    |   ├── shortcodes
    |   |   ├── img.html
    |   |   ├── vimeo.html
    |   |   └── youtube.html
    |   ├── index.html
    |   └── rss.xml
    └── public

This directory structure tells us a lot about this site:

1. the website intends to have two different types of content, posts and quotes.
2. It will also apply two different indexes to that content, categories and tags.
3. It will be displaying content in 3 different views, a list, a summary and a full page view.

Included with the repository is an example site ready to be rendered.

## Configuration

The directory structure and templates provide the majority of the
configuration for a site. In fact a config file isn't even needed for many websites
since the defaults used follow commonly used patterns.

The following is an example of a config file with the default values

    {
        "SourceDir" : "content",
        "LayoutDir" : "layouts",
        "PublishDir" : "public",
        "BuildDrafts" : false,
        "Tags" : { "category" : "categories", "tag" : "tags" },
        "BaseUrl"    : "http://yourSite.com/"
    }

## Usage 
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
      -p="1313": port for webserver to run on

The most common use is probably to run hugo with your current 
directory being the input directory.


    $ hugo
    > X pages created
    > Y indexes created


If you are working on things and want to see the changes 
immediately, tell Hugo to watch for changes. **It will 
recreate the site faster than you can tab over to 
your browser to view the changes.**

    $ hugo -p ~/mysite -w


# Layout

Hugo is very flexible about how you organize and structure your content.

## Templates

Hugo uses the excellent golang html/template library for it's template engine. It is an extremely
lightweight engine that provides a very small amount of logic. In our 
experience that it is just the right amount of logic to be able to create a good static website

This document will not cover how to use golang templates, but the [golang docs](http://golang.org/pkg/html/template/)
provide a good introduction.

### Template roles

There are 5 different kinds of templates that Hugo works with.

#### index.html
This file must exist in the layouts directory. It is the template used to render the 
homepage of your site.

#### rss.xml
This file must exist in the layouts directory. It will be used to render all rss documents.
The one provided in the example application will generate an ATOM format. 

*Important: Hugo will automatically add the following header line to this file.*

    <?xml version="1.0" encoding="utf-8" standalone="yes" ?>

#### Indexes
An index is a page that list multiple pieces of content. If you think of a typical blog, the tag 
pages are good examples of indexes.


#### Content Type(s)
Hugo supports multiple types of content. Another way of looking at this is that Hugo has the ability
to render content in a variety of ways as determined by the type.

#### Chrome
Chrome is simply the decoration of your site. It's not a requirement to have this, but in practice
it's very convenient. Hugo doesn't know anything about Chrome, it's simply a convention that you may
likely find beneficial. As you create the rest of your templates you will include templates from the 
/layout/chrome directory. I've found it helpful to include a header and footer template 
in Chrome so I can include those in the other full page layouts (index.html, indexes/ type/single.html).

### Adding a new content type

Adding a type is easy.

**Step 1:**
Create a directory with the name of the type in layouts.Type is always singular.  *Eg /layouts/post*.

**Step 2:**
Create a file called single.html inside your directory. *Eg /layouts/post/single.html*.

**Step 3:**
Create a file with the same name as your directory in /layouts/indexes/. *Eg /layouts/index/post.html*.

**Step 4:**
Many sites support rendering content in a few different ways, for instance a single page view and a 
summary view to be used when displaying a list of contents on a single page. Hugo makes no assumptions
here about how you want to display your content, and will support as many different views of a content
type as your site requires. All that is required for these additional views is that a template
exists in each layout/type directory with the same name.

For these, reviewing the example site will be very helpful in order to understand how these types work.

## Variables

Hugo makes a set of values available to the templates. Go templates are context based. The following
are available in the context for the templates.

**.Title**  The title for the content. <br>
**.Description** The description for the content.<br>
**.Keywords** The meta keywords for this content.<br>
**.Date** The date the content is published on.<br>
**.Indexes** These will use the field name of the plural form of the index (see tags and categories above)<br>
**.Permalink** The Permanent link for this page.<br>
**.FuzzyWordCount** The approximate number of words in the content.<br>
**.RSSLink** Link to the indexes' rss link <br>

Any value defined in the front matter, including indexes will be made available under `.Params`. 
Take for example I'm using tags and categories as my indexes. The following would be how I would access them:

**.Params.Tags** <br> 
**.Params.Categories** <br> 

Also available is `.Site` which has the following:

**.Site.BaseUrl** The base URL for the site as defined in the config.json file.<br>
**.Site.Indexes** The names of the indexes of the site.<br>
**.Site.LastChange** The date of the last change of the most recent content.<br>
**.Site.Recent** Array of all content ordered by Date, newest first<br>

# Content
Hugo uses markdown files with headers commonly called the front matter. Hugo respects the organization
that you provide for your content to minimize any extra configuration, though this can be overridden
by additional configuration in the front matter.

## Organization
In Hugo the content should be arranged in the same way they are intended for the rendered website.
Without any additional configuration the following will just work.

    .
    └── content
        ├── post
        |   ├── firstpost.md   // <- http://site.com/post/firstpost.html
        |   └── secondpost.md  // <- http://site.com/post/secondpost.html
        └── quote
            ├── first.md       // <- http://site.com/quote/first.html
            └── second.md      // <- http://site.com/quote/second.html


## Front Matter

The front matter is one of the features that gives Hugo it's strength. It enables
you to include the meta data of the content right with it. Hugo supports a few 
different formats. The main format supported is JSON. Here is an example:

    {
        "Title": "spf13-vim 3.0 release and new website",
        "Description": "spf13-vim is a cross platform distribution of vim plugins and resources for Vim.",
        "Tags": [ ".vimrc", "plugins", "spf13-vim", "vim" ],
        "Pubdate": "2012-04-06",
        "Categories": [ "Development", "VIM" ],
        "Slug": "spf13-vim-3-0-release-and-new-website"
    }

### Variables
There are a few predefined variables that Hugo is aware of and utilizes. The user can also create
any variable they want to. These will be placed into the `.Params` variable available to the templates.

#### Required

**Title**  The title for the content. <br>
**Description** The description for the content.<br>
**Pubdate** The date the content will be sorted by.<br>
**Indexes** These will use the field name of the plural form of the index (see tags and categories above)

#### Optional

**Draft** If true the content will not be rendered unless `hugo` is called with -d<br>
**Type** The type of the content (will be derived from the directory automatically if unset).<br>
**Slug** The token to appear in the tail of the url.<br>
  *or*<br>
**Url** The full path to the content from the web root.<br>
*If neither is present the filename will be used.*

## Example
Somethings are better shown than explained. The following is a very basic example of a content file:

**mysite/project/nitro.md  <- http://mysite.com/project/nitro.html**

    {
        "Title": "Nitro : A quick and simple profiler for golang",
        "Description": "",
        "Keywords": [ "Development", "golang", "profiling" ],
        "Tags": [ "Development", "golang", "profiling" ],
        "Pubdate": "2013-06-19",
        "Topics": [ "Development", "GoLang" ],
        "Slug": "nitro",
        "project_url": "http://github.com/spf13/nitro"
    }

    # Nitro

    Quick and easy performance analyzer library for golang.

    ## Overview

    Nitro is a quick and easy performance analyzer library for golang.
    It is useful for comparing A/B against different drafts of functions
    or different functions.

    ## Implementing Nitro

    Using Nitro is simple. First use go get to install the latest version
    of the library.

        $ go get github.com/spf13/nitro

    Next include nitro in your application.



# Extras

## Shortcodes
Because Hugo uses markdown for it's content format, it was clear that there's a lot of things that 
markdown doesn't support well. This is good, the simple nature of markdown is exactly why we chose it.

However we cannot accept being constrained by our simple format. Also unacceptable is writing raw
html in our markdown every time we want to include unsupported content such as a video. To do 
so is in complete opposition to the intent of using a bare bones format for our content and 
utilizing templates to apply styling for display.

To avoid both of these limitations Hugo has full support for shortcodes.

### What is a shortcode?
A shortcode is a simple snippet inside a markdown file that Hugo will render using a template.

Short codes are designated by the opening and closing characters of '{{%' and '%}}' respectively.
Short codes are space delimited. The first word is always the name of the shortcode.  Following the 
name are the parameters. The author of the shortcode can choose if the short code
will use positional parameters or named parameters (but not both). A good rule of thumb is that if a
short code has a single required value in the case of the youtube example below then positional
works very well. For more complex layouts with optional parameters named parameters work best.

The format for named parameters models that of html with the format name="value"

### Example: youtube

    {{% youtube 09jf3ow9jfw %}}

This would be rendered as 

    <div class="embed video-player">
    <iframe class="youtube-player" type="text/html"
        width="640" height="385" 
        src="http://www.youtube.com/embed/09jf3ow9jfw"
        allowfullscreen frameborder="0">
    </iframe>
    </div>

### Example: image with caption

    {{% img src="/media/spf13.jpg" title="Steve Francia" %}}

Would be rendered as:

    <figure >
        <img src="/media/spf13.jpg"  />
        <figcaption>
            <h4>Steve Francia</h4>
        </figcaption>
    </figure>


### Creating a shortcode

All that you need to do to create a shortcode is place a template in the layouts/shortcodes directory.

The template name will be the name of the shortcode.

**Inside the template**

To access a parameter by either position or name the index method can be used.

    {{ index .Params 0 }}
    or
    {{ index .Params "class" }}

To check if a parameter has been provided use the isset method provided by Hugo.

    {{ if isset .Params "class"}} class="{{ index .Params "class"}}" {{ end }}


# Meta

## Release Notes

* **0.7.0** July 4, 2013
  * Hugo now includes a simple server
  * First public release
* **0.6.0** July 2, 2013
  * Hugo includes an [example documentation site](http://hugo.spf13.com) which it builds
* **0.5.0** June 25, 2013
  * Hugo is quite usable and able to build [spf13.com](http://spf13.com)

## Roadmap
In no particular order, here is what I'm working on:

 * Pagination
 * Support for top level pages (other than homepage)
 * Series support
 * Syntax highlighting
 * Previous & Next
 * Related Posts
 * Support for TOML front matter
 * Proper YAML support for front matter
 * Support for other formats

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

## Contributors

* [spf13](https://github.com/spf13)


## License

Hugo is released under the Simple Public License. See [LICENSE.md](https://github.com/spf13/hugo/blob/master/LICENSE.md).
