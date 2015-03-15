---
date: 2014-05-14T02:13:50Z
menu:
  main:
    parent: content
next: /content/ordering
prev: /content/types
title: Archetypes
weight: 50
---

Hugo v0.11 introduced the concept of a content builder. Using the
command: <code>hugo new <em>[relative new content path]</em></code>,
you can start a content file with the date and title automatically set.
While this is a welcome feature, active writers need more.

Hugo presents the concept of archetypes, which are archetypal content files
with pre-configured [front matter](content/front-matter) which will
populate each new content file whenever you run the `hugo new` command.


## Example

### Step 1. Creating an archetype

In this example scenario, we have a blog with a single content type (blog post).
We will use ‘tags’ and ‘categories’ for our taxonomies, so let's create an archetype file with ‘tags’ and ‘categories’ pre-defined, as follows:

#### archetypes/default.md

    +++
    tags = ["x", "y"]
    categories = ["x", "y"]
    +++

> __CAVEAT:__  Some editors (e.g. Sublime, Emacs) do not insert an EOL (end-of-line) character at the end of the file (i.e. EOF).  If you get a [strange EOF error](/troubleshooting/strange-eof-error/) when using `hugo new`, please open each archetype file (i.e.&nbsp;`archetypes/*.md`) and press <kbd>Enter</kbd> to type a carriage return after the closing `+++` or `---` as necessary.


### Step 2. Using the archetype

Now, with `archetypes/default.md` in place, let's create a new post in the `post` section with the `hugo new` command:

    $ hugo new post/my-new-post.md

Hugo would create the file with the following contents:

#### content/post/my-new-post.md

    +++
    title = "my new post"
    date = "2015-01-12T19:20:04-07:00"
    tags = ["x", "y"]
    categories = ["x", "y"]
    +++

We see that the `title` and `date` variables have been added, in addition to the `tags` and `categories` variables which were carried over from `archetype/default.md`.

Congratulations!  We have successfully created an archetype and used it for our new contents.  That's all there is to it!


## Using a different front matter format

By default, the front matter will be created in the TOML format
regardless of what format the archetype is using.

You can specify a different default format in your site-wide config file
(e.g. `config.toml`) using the `MetaDataFormat` directive.
Possible values are `"toml"`, `"yaml"` and `"json"`.


## Which archetype is being used

The following rules apply:

* If an archetype with a filename that matches the content type being created, it will be used.
* If no match is found, `archetypes/default.md` will be used.
* If neither is present and a theme is in use, then within the theme:
    * If an archetype with a filename that matches the content type being created, it will be used.
    * If no match is found, `archetypes/default.md` will be used.
* If no archetype files are present, then the one that ships with Hugo will be used.

Hugo provides a simple archetype which sets the `title` (based on the
file name) and the `date` in RFC&nbsp;3339 format based on
[`now()`](http://golang.org/pkg/time/#Now), which returns the current time.

> *Note: `hugo new` does not automatically add `draft = true` when the user
> provides an archetype.  This is by design, rationale being that
> the archetype should set its own value for all fields.
> `title` and `date`, which are dynamic and unique for each piece of content,
> are the sole exceptions.*

Content type is automatically detected based on the path. You are welcome to declare which type to create using the `--kind` flag during creation.
