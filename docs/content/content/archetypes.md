---
lastmod: 2016-10-01
date: 2014-05-14T02:13:50Z
menu:
  main:
    parent: content
next: /content/ordering
prev: /content/types
title: Archetypes
weight: 50
toc: true
---

In Hugo v0.11, we introduced the concept of a content builder. Using the CLI 
command <code>hugo new <em>[path/to/my/content]</em></code>, an author could 
create an empty content file, with the date and title automatically defined in 
the front matter of the post. While this was a welcome feature, active writers 
need more flexibility.

When defining a custom content type, you can use an **archetype** as a way to 
define the default metadata for a new post of that type. 

**Archetypes** are quite literally archetypal content files with pre-configured 
[front matter](/content/front-matter). An archetype will populate each new 
content file of a given type with any default metadata you've defined whenever 
you run the `hugo new` command.

## Example

### Step 1. Creating an archetype

In the following example scenario, suppose we have a blog with a single content 
type (blog post). Our imaginary blog will use ‘tags’ and ‘categories’ for its 
taxonomies, so let's create an archetype file with ‘tags’ and ‘categories’ 
pre-defined:

#### archetypes/default.md

```toml
+++
tags = ["x", "y"]
categories = ["x", "y"]
+++
```

> __CAVEAT:__  Some editors (e.g. Sublime, Emacs) do not insert an EOL (end-of-line) character at the end of the file (i.e. EOF).  If you get a [strange EOF error](/troubleshooting/strange-eof-error/) when using `hugo new`, please open each archetype file (i.e.&nbsp;`archetypes/*.md`) and press <kbd>Enter</kbd> to type a carriage return after the closing `+++` or `---` as necessary.


### Step 2. Using the archetype

Now, with `archetypes/default.md` in place, let's create a new post in the `post` 
section with the `hugo new` command:

    $ hugo new post/my-new-post.md

Hugo will now create the file with the following contents:

#### content/post/my-new-post.md

```toml
+++
title = "my new post"
date = "2015-01-12T19:20:04-07:00"
tags = ["x", "y"]
categories = ["x", "y"]
+++
```

We see that the `title` and `date` variables have been added, in addition to the 
`tags` and `categories` variables which were carried over from `archetype/default.md`.

Congratulations! We have successfully created an archetype and used it to 
quickly scaffold out a new post. But wait, what if we want to create some content
that isn't exactly a blog post, like a profile for a musician? Let's see how
using **archetypes** can help us out.

### Creating custom archetypes

Previously, we had created a new content type by adding a new subfolder to the 
content directory. In this case, its name would be `content/musician`. To begin 
using a `musician` archetype for each new `musician` post, we simply need to 
create a file named after the content type called `musician.md`, and put it in 
the `archetypes` directory, similar to the one below.

#### archetypes/musician.md

```toml
+++
name = ""
bio = ""
genre = ""
+++
```

Now, let's create a new musician.

    $ hugo new musician/mozart.md

This time, Hugo recognizes our custom `musician` archetype and uses it instead of 
the default one. Take a look at the new `musician/mozart.md` post. You should see 
that the generated file's front matter now includes the variables `name`, `bio`, 
and `genre`.

#### content/musician/mozart.md

```toml
+++
title = "mozart"
date = "2015-08-24T13:04:37+02:00"
name = ""
bio = ""
genre = ""
+++
```

## Using a different front matter format

By default, the front matter will be created in the TOML format
regardless of what format the archetype is using.

You can specify a different default format in your site-wide config file
(e.g. `config.toml`) using the `MetaDataFormat` directive.
Possible values are `"toml"`, `"yaml"` and `"json"`.

## Which archetype is being used

The following rules apply when creating new content:

* If an archetype with a filename matching the new post's [content type](/content/types) exists, it will be used.
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

The content type is automatically detected based on the file path passed to the 
Hugo CLI command <code>hugo new <em>[my-content-type/post-name]</em></code>. To 
override the content type for a new post, include the `--kind` flag during creation.

> *Note: if you wish to use archetypes that ship with a theme, the theme MUST be specified in your `config.toml`.*
