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

Typically, each piece of content you create within a Hugo project will have [front matter](/content/front-matter/) that follows a consistent structure. If you write blog posts, for instance, you might use the following front matter for the vast majority of those posts:

```toml
+++
title = ""
date = ""
slug = ""
tags = [
  ""
]
categories = [
  ""
]
draft = true
+++
```

You can always add non-typical front matter to any piece of content, but since it takes extra work to develop a theme that handles unique metadata, consistency is simpler.

With this in mind, Hugo has a convenient feature known as *archetypes* that allows users to define default front matter for new pieces of content.

By using archetypes, we can:

1. **Save time**. Stop writing the same front matter over and over again.
2. **Avoid errors**. Reduce the odds of typos, improperly formatted syntax, and other simple mistakes.
3. **Focus on more important things**. Avoid having to remember all of the fields that need to be associated with each piece of content. (This is particularly important for larger projects with complex front matter and a variety of content types.)

Let's explore how they work.

## Built-in Archetypes

If you've been using Hugo for a while, there's a decent chance you've come across archetypes without even realizing it. This is because Hugo includes a basic, built-in archetype that is used by default whenever it generates a content file.

To see this in action, open the command line, navigate into your project's directory, and run the following command:

```bash
hugo new hello-world.md
```

This `hugo new` command creates a new content file inside the project's `content` directory — in this case, a file named `hello-world.md` — and if you open this file, you'll notice it contains the following front matter:

```toml
+++
date = "2017-05-31T15:18:11+10:00"
draft = true
title = "hello world"
+++
```

Here, we can see that three fields have been added to the document: a `title` field that is based on the file name we defined, a `draft` field that ensures this content won't be published by default, and a `date` field that is auto-populated with the current date and time in the [RFC 3339](https://stackoverflow.com/questions/522251/whats-the-difference-between-iso-8601-and-rfc-3339-date-formats) format.

This, in its most basic form, is an example of an archetype. To understand how useful they can be though, it's best if we create our own.

## Creating Archetypes

In this section, we're going to create an archetype that will override the built-in archetype, allowing us to define custom front matter that will be included in any content files that we generate with the `hugo new` command.

To achieve this, create a file named `default.md` inside the `archetypes` folder of a Hugo project. (If the folder doesn't exist, create it.)

Then, inside this file, define the following front matter:

```toml
+++
slug = ""
tags = []
categories = []
draft = true
+++
```

You'll notice that we haven't defined a `title` or `date` field. This is because Hugo will automatically add these fields to the beginning of the front matter. We do, however, need to define the `draft` field if we want it to exist in our front matter.

You'll also notice that we're writing the front matter in the TOML format. It's possible to define archetype front matter in other formats, but a setting needs to be changed in the configuration file for this to be possible. See the "[Archetype Formats](#archetype-formats)" section of this article for more details.

Next, run the following command:

```bash
hugo new my-archetype-example.md
```

This command will generate a file named `my-archetype-example.md` inside the `content` directory, and this file will contain the following output:

```toml
+++
categories = []
date = "2017-05-31T15:21:13+10:00"
draft = true
slug = ""
tags = []
title = "my archetype example"
+++
```

As we can see, the file contains the `title` and `date` property that Hugo created for us, along with the front matter that we defined in the `archetypes/default.md` file.

You'll also notice that the fields have been sorted into alphabetical order. This is an unintentional side-effect that stems from the underlying code libraries that Hugo relies upon. It is, however, [a known issue that is actively being discussed](https://github.com/gohugoio/hugo/issues/452).

## Section Archetypes

By creating the `archetypes/default.md` file, we've created a default archetype that is more useful than the built-in archetype, but since Hugo encourages us to [organize our content into sections](/content/sections/), each of which will likely have different front matter requirements, a "one-size-fits-all" archetype isn't necessarily the best approach.

To accommodate for this, Hugo allows us to create archetypes for each section of our project. This means, whenever we generate content for a certain section, the appropriate front matter for that section will be automatically included in the generated file.

To see this in action, create a "photo" section by creating a directory named "photo" inside the `content` directory.

Then create a file named `photo.md` inside the `archetypes` directory and include the following front matter inside this file:

```toml
+++
image_url = ""
camera = ""
lens = ""
aperture = ""
iso = ""
draft = true
+++
```

Here, the critical detail is that the `photo.md` file in the `archetypes` directory is named after the `photo` section that we just created. By sharing a name, Hugo can understand that there's a relationship between them.

Next, run the following command:

```bash
hugo new photo/my-pretty-cat.md
```

This command will generate a file named `my-pretty-cat.md` inside the `content/photo` directory, and this file will contain the following output:

```toml
+++
aperture = ""
camera = ""
date = "2017-05-31T15:25:18+10:00"
draft = true
image_url = ""
iso = ""
lens = ""
title = "my pretty cat"
+++
```

As we can see, the `title` and `date` fields are still included by Hugo, but the rest of the front matter is being generated from the `photo.md` archetype instead of the `default.md` archetype.

### Tip: Default Values

To make archetypes more useful, define default values for any fields that will always be set to a range of limited options. In the case of the `photo.md` archetype, for instance, you could include lists of the various cameras and lenses that you own:

```toml
+++
image_url = ""
camera = [
  "Sony RX100 Mark IV",
  "Canon 5D Mark III",
  "iPhone 6S"
]
lens = [
  "Canon EF 50mm f/1.8",
  "Rokinon 14mm f/2.8"
]
aperture = ""
iso = ""
draft = true
+++
```

Then, after generating a content file, simply remove the values that aren't relevant. This saves you from typing out the same options over and over again while ensuring consistency in how they're written.

## Scaffolding Content

Archetypes aren't limited to defining default front matter. They can also be used to define a default structure for the body of Markdown documents.

For example, imagine creating a `review.md` archetype for the purpose of writing camera reviews. This is what the front matter for such an archetype might look like:

```toml
+++
manufacturer = ""
model = ""
price = ""
releaseDate = ""
rating = ""
+++
```

But reviews tend to follow strict formats and need to answer specific questions, and it's with these expectations of precise structure that archetypes can prove to be even more useful.

For the sake of writing reviews, for instance, we could define the structure of a review beneath the front matter of the `review.md` file:

```markdown
+++
manufacturer = ""
model = ""
price = ""
releaseDate = ""
rating = ""
+++

## Introduction

## Sample Photos

## Conclusion
```

Then, whenever we use the `hugo new` command to create a new review, not only will the default front matter be copied into the newly created Markdown document, but the body of the `review.md` archetype will also be copied.

To take this further though — and to ensure authors on multi-author websites are on the same page about how content should be written — we could include notes and reminders within the archetype:

```markdown
+++
manufacturer = ""
model = ""
price = ""
releaseDate = ""
rating = ""
+++

## Introduction

<!--
  What is the selling point of the camera?
  What has changed since last year's model?
  Include a bullet-point list of key features.
-->

## Sample Photos

<!-- TODO: Take at least 12 photos in a variety of situations. -->

## Conclusion

<!--
  Is this camera worth the money?
  Does it accomplish what it set out to achieve?
  Are there any specific groups of people who should/shouldn't buy it?
  Would you recommend it to a friend?
  Are there alternatives on the horizon?
-->

```

That way, each time we generate a new content file, we have a series of handy notes to push us closer to a piece of writing that's suitable for publishing.

(If you're wondering why the notes are wrapped in the HTML comment syntax, it's to ensure they won't appear inside the preview window of whatever Markdown editor the author happens to be using. They're not strictly necessary though.)

This is still a fairly simple example, but if your content usually contains a variety of components — headings, bullet-points, images, [short-codes](/extras/shortcodes/), etc — it's not hard to see the time-saving benefits of placing these components in the body of an archetype file.

## Theme Archetypes

Whenever you generate a content file with the `hugo new` command, Hugo will start by searching for archetypes in the `archetypes` directory, initially looking for an archetype that matches the content's section and falling-back on the `default.md` archetype (if one is present). If no archetypes are found in this directory, Hugo will continue its search in the `archetypes` directory of the currently active theme. In other words, it's possible for themes to come packaged with their own archetypes, ensuring that users of that theme format their content files with correctly structured front matter.

To allow Hugo to use archetypes from a theme, [that theme must be activated via the project's configuration file](/themes/usage/):

```toml
theme = "ThemeNameGoesHere"
```

If an archetype doesn't exist in the `archetypes` directory at the top-level of a project or inside the `archetypes` directory of an active theme, the built-in archetype will be used.

{{< figure src="/img/content/archetypes/archetype-hierarchy.png" alt="How Hugo Decides Which Archetype To Use" >}}

## Archetype Formats

By default, the `hugo new` command will generate front matter in the TOML format. This means, even if we define the front matter in our archetype files as YAML or JSON, it will be converted to the TOML format before it ends up in our content files.

Fortunately, this functionality can be overwritten.

Inside the project's configuration file, simply define a `metaDataFormat` property:

```toml
metaDataFormat = ""
```

Then set this property to any of the following values:

* toml
* yaml
* json

By defining this option, any front matter will be generated in your preferred format.

It's worth noting, however, that when generating front matter in the TOML format, you might encounter the following error:

```bash
Error: cannot convert type <nil> to TomlTree
```

This is because, to generate TOML, all of the fields in the front matter need to have a default value, even if that default value is just an empty string.

For example, this YAML would *not* successfully compile into the TOML format:

```yaml
---
slug:
tags:
categories:
draft:
---
```

But this YAML *would* successfully compile:

```yaml
---
slug: ""
tags:
  -
categories:
  -
draft: true
---
```

It's a subtle yet important detail to remember.

## Notes

* Prior to Hugo v0.13, some users received [an "EOF" error when using archetypes](https://github.com/gohugoio/hugo/issues/776), related to what text editor they used to create the archetype. As of Hugo v0.13, this error has been [resolved](https://github.com/gohugoio/hugo/pull/785).
