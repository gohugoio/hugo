---
title: Content Types
description: Hugo is built around content organized in sections.
date: 2017-02-01
categories: [content management]
keywords: [lists,sections,content types,types,organization]
menu:
  docs:
    parent: "content-management"
    weight: 60
weight: 60	#rem
draft: false
aliases: [/content/types]
toc: true
---

A **content type** is a way to organize your content. Hugo resolves the content type from either the `type` in front matter or, if not set, the first directory in the file path. E.g. `content/blog/my-first-event.md` will be of type `blog` if no `type` set.

A content type is used to

* Determine how the content is rendered. See [Template Lookup Order](/templates/lookup-order/) and [Content Views](https://gohugo.io/templates/views) for more.
* Determine which [archetype](/content-management/archetypes/) template to use for new content.


