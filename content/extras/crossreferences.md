---
date: 2014-11-25
menu:
  main:
    parent: extras
next: /extras/livereload
prev: /extras/comments
title: Cross-References
weight: 40
---

Hugo makes it easy to link documents together with the `ref` and `relref` shortcodes. These shortcodes are also used to safely provide links to headings inside of your content, whether across documents or within a document. The only difference between `ref` and `relref` is whether the resulting URL is absolute (`http://1.com/about/`) or relative (`/about/`).

## Using `ref` and `relref`

    {{</* ref "document" */>}}
    {{</* ref "#anchor" */>}}
    {{</* ref "document#anchor" */>}}
    {{</* relref "document" */>}}
    {{</* relref "#anchor" */>}}
    {{</* relref "document#anchor" */>}}

The single parameter to `ref` is a string with a content _document name_ (`about.md`), an in-document _anchor_ (`#who`), or both (`about.md#who`).

### Document Names

The _document name_ is the name of a document including the format extension; this may be just the filename, or the relative path from the `content/` directory. With a document `content/blog/post.md`, either format will produce the same result.

    {{</* relref "blog/post.md" */>}} ⇒ `/blog/post/`
    {{</* relref "post.md" */>}} ⇒ `/blog/post/`

If you have multiple sections with the same filename, you should only use the relative path format, because the behaviour is _undefined_. So, if I also have a document `link/post.md`, the output of `ref` is unknown for `post.md`.

    {{</* relref "blog/post.md" */>}} ⇒ `/blog/post/`
    {{</* relref "post.md" */>}} ⇒ `/blog/post/` (maybe)
    {{</* relref "post.md" */>}} ⇒ `/link/post/` (maybe)
    {{</* relref "link/post.md" */>}} ⇒ `/link/post/`

A relative document name must *not* begin with a slash (`/`).

    {{</* relref "/blog/post.md" */>}} ⇒ `""`

### Anchors

When an _anchor_ is provided by itself, the current page’s unique identifier will be appended; when an _anchor_ is provided with a document name, the found page's unique identifier will be appended.

    {{</* relref "#who" */>}} ⇒ `#who:9decaf7`
    {{</* relref "blog/post.md#who" */>}} ⇒ `/blog/post/#who:badcafe`

More information about document unique identifiers and headings can be found [below]({{< ref "#hugo-heading-anchors" >}}).

### Examples

* `{{</* ref "blog/post.md" */>}}` ⇒ `http://1.com/blog/post/`
* `{{</* ref "post.md#tldr" */>}}` ⇒ `http://1.com/blog/post/#tldr:caffebad`
* `{{</* relref "post.md" */>}}` ⇒ `/blog/post/`
* `{{</* relref "blog/post.md#tldr" */>}}` ⇒ `/blog/post/#tldr:caffebad`
* `{{</* ref "#tldr" */>}}` ⇒ `#tldr:badcaffe`
* `{{</* relref "#tldr" */>}}` ⇒ `#tldr:badcaffe`

## Hugo Heading Anchors

When using Markdown document types, Hugo generates heading anchors automatically. The generated anchor for this section is `hugo-heading-anchors`. Because the heading anchors are generated automatically, Hugo takes some effort to ensure that heading anchors are unique both inside a document and across the entire site.

Ensuring heading uniqueness across the site is accomplished with a unique identifier for each document based on its path. Unless a document is renamed or moved between sections *in the filesystem*, the unique identifier for the document will not change: `blog/post.md` will always have a unique identifier of `81df004c333b392d34a49fd3a91ba720`.

`ref` and `relref` were added so you can make these reference links without having to know the document’s unique identifier. (The links in document tables of contents are automatically up-to-date with this value.)

    {{</* relref "extras/crossreferences.md#hugo-heading-anchors" */>}}
    /extras/crossreferences/#hugo-heading-anchors:77cd9ea530577debf4ce0f28c8dca242

> What follows is a deeper discussion of *why* and *how* Hugo generates heading anchors. It is not necessary to know this to use `ref` and `relref`, but it may be useful in understanding how some anchors may not match your expectations.

### How to Generate a Heading Anchor

Convert the text of the heading to lowercase.

    Hugo: A Fast & Modern Static Web Engine
    hugo: a fast & modern static web engine

Replace anything that isn't an ASCII letter (`a-z`) or number (`0-9`) with a dash (`-`).

    hugo: a fast & modern static web engine
    hugo--a-fast---modern-static-web-engine

Get rid of extra dashes.

    hugo--a-fast---modern-static-web-engine
    hugo-a-fast-modern-static-web-engine

You have just converting the text of a heading to a suitable anchor. If your document has unique heading text, all of the anchors will be unique, too.

#### Specifying Heading Anchors

You can also tell Hugo to use a particular heading anchor.

    # Hugo: A Fast & Modern Static Web Engine {#hugo-main}

Hugo will use `hugo-main` as the heading anchor.

### What About Duplicate Heading Anchors?

The technique outlined above works well enough, but some documents have headings with identical text, like the [shortcodes](/extras/shortcodes/) page—there are three headings with the text "Example". You can specify heading anchors manually:

    ### Example {#example-1}
    ### Example {#example-2}
    ### Example {#example-3}

It’s easy to forget to do that all the time, and Hugo is smart enough to do it for you. It just adds `-x` to the end of each heading it has already seen.

* `### Example` ⇒ `example`
* `### Example` ⇒ `example-1`
* `### Example` ⇒ `example-2`

Sometimes it's a little harder, but Hugo can recover from those, too, by adding more suffixes:

* `# Heading` ⇒ `heading`
* `# Heading 1` ⇒ `heading-1`
* `# Heading` ⇒ `heading-1-1`
* `# Heading` ⇒ `heading-1-2`
* `# Heading 1` ⇒ `heading-2`

This can even affect specified heading anchors that come after a generated heading anchor.

* `# My Heading` ⇒ `my-heading`
* `# My Heading {#my-heading}` ⇒ `my-heading-1`

> This particular collision and override is unfortunate, but unavoidable because Hugo processes each heading for collision detection as it sees it during conversion.

This technique works well for documents rendered on individual pages, like blog posts. What about on Hugo list pages?

### Unique Heading Anchors in Lists

Hugo converts each document from Markdown independently. it doesn’t know that `blog/post.md` has an "Example" heading that will collide with the "Example" heading in `blog/post2.md`. Even if it did know this, the addition of `blog/post3.md` should not cause the anchors for the headings in the other blog posts to change.

Enter the document’s unique identifier. To prevent this sort of collision on
list pages, Hugo always appends the document's to a generated heading anchor.
So, the "Example" heading in `blog/post.md` actually turns into
`#example:81df004…`, and the "Example" heading in `blog/post2.md` actually
turns into `#example:8cf1599…`. All you have to know is the heading anchor that
was generated, not the document identifier; `ref` and `relref` take care of the
rest for you.

    <a href='{{</* relref "blog/post.md#example" */>}}'>Post Example</a>
    <a href='/blog/post.md#81df004…'>Post Example</a>

    [Post Two Example]({{</* relref "blog/post2.md#example" */>}})
    <a href='/blog/post2.md#8cf1599…'>Post Two Example</a>

Now you know.
