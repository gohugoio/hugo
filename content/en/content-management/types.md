---
title: Content types
description: Hugo is built around content organized in sections.
categories: [content management]
keywords: [lists,sections,content types,types,organization]
menu:
  docs:
    parent: content-management
    weight: 130
weight: 130
toc: true
aliases: [/content/types]
---

A **content type** is a way to organize your content. Hugo resolves the content type from either the `type` in front matter or, if not set, the first directory in the file path. E.g. `content/blog/my-first-event.md` will be of type `blog` if no `type` is set.

A content type is used to

- Determine how the content is rendered. See [Template Lookup Order](/templates/lookup-order/) and [Content Views](/templates/content-view) for more.
- Determine which [archetype](/content-management/archetypes/) template to use for new content.
