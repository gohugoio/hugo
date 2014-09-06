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
command: `hugo new [relative new content path]` you can start a content file
with the date and title automatically set. This is a welcome feature, but 
active writers need more. 

Hugo presents the concept of archetypes which are archetypal content files.

## Example archetype

In this example scenario I have a blog with a single content type (blog post).
I use ‘tags’ and ‘categories’ for my taxonomies.

### archetypes/default.md

    +++
    tags = ["x", "y"]
    categories = ["x", "y"]
    +++


## Using archetypes

If I wanted to create a new post in the `posts` section, I would run the following command:

`hugo new posts/my-new-post.md`

Hugo would create the file with the following contents:

### contents/posts/my-new-post.md

    +++
    title = "my new post"
    date = 2014-05-14T02:13:50Z
    tags = ["x", "y"]
    categories = ["x", "y"]
    +++


## Using a different front matter format

By default, the front matter will be created in the TOML format
regardless of what format the archetype is using.

You can specify a different default format in your config file using
the `MetaDataFormat` directive. Possible values are `toml`, `yaml` and `json`.


## Which archetype is being used

The following rules apply:

* If an archetype with a filename that matches the content type being created, it will be used.
* If no match is found, `archetypes/default.md` will be used.
* If neither are present and a theme is in use, then within the theme:
    * If an archetype with a filename that matches the content type being created, it will be used.
    * If no match is found, `archetypes/default.md` will be used.
* If no archetype files are present, then the one that ships with Hugo will be used.

Hugo provides a simple archetype which sets the title (based on the
file name) and the date based on `now()`.

Content type is automatically detected based on the path. You are welcome to declare which 
type to create using the `--kind` flag during creation.

