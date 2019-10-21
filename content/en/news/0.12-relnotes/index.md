---
date: 2015-01-21T20:35:00Z
description: "Hugo 0.12 focused on polishing the theme engine and adding critical functionality to the templates."
title: "Hugo 0.12"
categories: ["Releases"]
---

A lot has happened since Hugo v0.11.0 was released. Most of the work has been
focused on polishing the theme engine and adding critical functionality to the
templates.

This release represents over 90 code commits from 28 different contributors.

- 10 [new themes](https://github.com/spf13/hugoThemes) created by the community
- fully themable [partials](/templates/partials/)
- [404 template](/templates/404/) support in themes
- [shortcode](/extras/shortcodes/) support in themes
- [views](/templates/views/) support in themes
- inner [shortcode](/extras/shortcodes/) content now treated as markdown
- support for header ids in markdown (# header {#myid})
- [where](/templates/lists/#where) template function to filter lists of content, taxonomies, etc
- [groupby](/templates/list) & [groupbydate](/templates/list) methods to group pages
- taxonomy [pages list](/taxonomies/methods/) now sortable, filterable, limitable & groupable
- general cleanup to taxonomies & documentation to make it more clear and consistent
- [showcase](/showcase/) returned and has been expanded
- pretty links now always have trailing slashes
- [baseurl](/overview/configuration/) can now include a subdirectory
- better feedback about draft & future post rendering
- a variety of improvements to [the website](/)
