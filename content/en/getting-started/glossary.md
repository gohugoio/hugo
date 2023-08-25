---
title: Glossary of terms
description: Terms commonly used throughout the documentation.
keywords: [glossary]
menu:
  docs:
    parent: getting-started
    weight: 60
weight: 60
type: glossary
---

<!-- Use level 3 headings for each term in the glossary. -->

### action

See [template action](#template-action).

### archetype

A template for new content. See [details](/content-management/archetypes/).

### argument

A [scalar](#scalar), [array](#array), [slice](#slice), [map](#map), or [object](#object) passed to a [function](#function), [method](#method), or [shortcode](#shortcode).

### array

A numbered sequence of elements. Unlike Go's [slice](#slice) data type, an array has a fixed length. See the [Go&nbsp;documentation](https://go.dev/ref/spec#Array_types) for details.

### bool

See [boolean](#boolean).

### boolean

A data type with two possible values, either `true` or `false`.

### branch bundle

A [page bundle](#page-bundle) with an&nbsp;_index.md file and zero or more [resources](#resource). Analogous to a physical branch, a branch bundle may have descendants including regular pages, [leaf bundles](/getting-started/glossary/#leaf-bundle), and other branch bundles. See [details](/content-management/page-bundles/).

### build

To generate a static site that includes HTML files and assets such as images, CSS, and JavaScript. The build process includes rendering and resource transformations.

### bundle

See [page bundle](#page-bundle).

### cache

A software component that stores data so that future requests for the same data are faster.

### collection

Typically, a collection of pages, but may also refer to an [array](#array),  [slice](#slice), or [map](#map). For example, the pages within a site's "articles" section are a page collection.

### content format

A markup language for creating content. Typically markdown, but may also be HTML, AsciiDoc, Org, Pandoc, or reStructuredText. See [details](/content-management/formats/).

### content type

A classification of content inferred from the top-level directory name or the `type` set in [front matter](#front-matter). Pages in the root of the content directory, including the home page, are of type "page". Accessed via `.Page.Type` in [templates](#template). See&nbsp;[details](/content-management/types/).

### context

Represented by a period "." within a [template action](#template-action), context is the current location in a data structure. For example, while iterating over a [collection](#collection) of pages, the context within each iteration is the page's data structure. The context received by each template depends on template type and/or how it was called. See [details](/templates/introduction/#the-dot).

### flag

An option passed to a command-line program, beginning with one or two hyphens. See [details](/commands/hugo/).

### float

See [floating point](#floating-point).

### floating point

A numeric data type with a fractional component. For example, `3.14159`.

### function

Used within a [template action](#template-action), a function takes one or more [arguments](#argument) and returns a value. Unlike [methods](#method), functions are not associated with an [object](#object). See [details](/functions/).

### front matter

Metadata at the beginning of each content page, separated from the content by format-specific delimiters. See&nbsp;[details](/content-management/front-matter/).

### int

See [integer](#integer).

### integer

A numeric data type without a fractional component. For example, `42`.

### internationalization

Software design and development efforts that enable [localization](#localization). See the [W3C definition](https://www.w3.org/International/questions/qa-i18n). Abbreviated i18n.

### kind

See [page kind](#page-kind).

### layout

See [template](#template).

### leaf bundle

A [page bundle](#page-bundle) with an index.md file and zero or more [resources](#resource). Analogous to a physical leaf, a leaf bundle is at the end of a branch. Hugo ignores content (but not resources) beneath the leaf bundle. See [details](/content-management/page-bundles/).

### list page

Any [page kind](#page-kind) that receives a page [collection](#collection) in [context](#context). This includes the home page, [section pages](#section-page), [taxonomy pages](#taxonomy-page), and [term pages](#term-page).

### localization

Adaptation of a site to meet language and regional requirements. This includes translations, language-specific media, date and currency formats, etc. See [details](/content-management/multilingual/) and the [W3C definition](https://www.w3.org/International/questions/qa-i18n). Abbreviated l10n.

### map

An unordered group of elements, each indexed by a unique key. See the [Go&nbsp;documentation](https://go.dev/ref/spec#Map_types) for details.

### method

Used within a [template action](#template-action) and associated with an [object](#object), a method takes zero or more [arguments](#argument) and returns a value. For example, `.IsHome` is a method on the `.Page` object which returns `true` if the current page is the home page. See also [function](#function).

### module

Like a [theme](#theme), a module is a packaged combination of [archetypes](#archetype), assets, content, data, [templates](#template), translations, or configuration settings. A module may serve as the basis for a new site, or to augment an existing site. See [details](/hugo-modules/).

### object

A data structure with or without associated [methods](#method).

### page bundle

A directory that encapsulates both content and associated [resources](#resource). There are two types of page bundles: [leaf bundles](/getting-started/glossary/#leaf-bundle) and [branch bundles](/getting-started/glossary/#branch-bundle). See [details](/content-management/page-bundles/).

### page kind

A classification of rendered pages, one of "home", "page", "section", "taxonomy", or "term". Accessed via `.Page.Kind` in [templates](#template). See&nbsp;[details](/templates/section-templates/#page-kinds).

### pager

Created during [pagination](#pagination), a pager contains a subset of a section list, and navigation links to other pagers.

### paginate

To split a [section](#section) list into two or more [pagers](#pager) See&nbsp;[details](/templates/pagination/).

### pagination

The process of [paginating](#paginate) a [section](#section) list.

### parameter

Typically, a user-defined key/value pair at the site or page level, but may also refer to a configuration setting or an [argument](#argument).

### partial

A [template](#template) called from any other template including [shortcodes](#shortcode), [render hooks](#render-hook), and other partials. A partial either renders something or returns something. A partial can also call itself, for example, to [walk](#walk) a data structure.

### permalink

The absolute URL of a rendered page, including scheme and host.

### pipe

See [pipeline](#pipeline).

### pipeline

Within a [template action](#template-action), a pipeline is a possibly chained sequence of values, [function](#function) calls, or [method](#method) calls. Functions and methods in the pipeline may take multiple [arguments](#argument).

A pipeline may be *chained* by separating a sequence of commands with pipeline characters "|". In a chained pipeline, the result of each command is passed as the last argument to the following command. The output of the final command in the pipeline is the value of the pipeline. See the [Go&nbsp;documentation](https://pkg.go.dev/text/template#hdr-Pipelines) for details.

### publish

See [build](#build).

### regular page

Content with the "page" [page kind](#page-kind). See also [section page](#section-page).

### render hook

A [template](#template) that overrides standard markdown rendering. See [details](/templates/render-hooks/).

### resource

Any file consumed by the build process to augment or generate content, structure, behavior, or presentation. For example: images, videos, content snippets, CSS, Sass, Javascript, and data.

Hugo supports three types of resources: page resources (located in a [page bundle](/getting-started/glossary/#page-bundle)), global resources (located in the assets directory), and remote resources (typically accessed via https).

### scalar

A single value, one of [string](#string), [integer](#integer), [floating point](#floating-point), or [boolean](#boolean).

### section

A top-level content directory, or any content directory with an&nbsp;_index.md file. A content directory with an&nbsp;_index.md file is also known as a [branch bundle](/getting-started/glossary/#branch-bundle). Section templates receive one or more page [collections](#collection) in [context](#context). See [details](/content-management/sections/).

### section page

Content with the "section" [page kind](#page-kind). Typically a listing of [regular pages](#regular-page) and/or [section pages](#section-page) within the current [section](#section). See also [regular page](#regular-page).

### shortcode

A [template](#template) called from within markdown, taking zero or more [arguments](#argument). See [details](/content-management/shortcodes/).

### slice

A numbered sequence of elements. Unlike Go's [array](#array) data type, slices are dynamically sized. See the [Go&nbsp;documentation](https://go.dev/ref/spec#Slice_types) for details.

### string

A sequence of bytes. For example, `"What is 6 times 7?"`&nbsp;.

### taxonomy

A group of related [terms](#term) used to classify content. For example, a "colors" taxonomy might include the terms "red", "green", and "blue". See&nbsp;[details](/content-management/taxonomies/).

### taxonomy page

Content with the "taxonomy" [page kind](#page-kind). Typically a listing of [terms](#term) within a given [taxonomy](#taxonomy).

### template

An HTML file with [template actions](#template-action), located within the layouts directory of a project, theme, or module. See&nbsp;[details](/templates/).

### template action

A data evaluation or control structure within a [template](#template), delimited by "{{"&nbsp;and&nbsp;"}}". See the [Go&nbsp;documentation](https://pkg.go.dev/text/template#hdr-Actions) for details.

### term

A member of a [taxonomy](#taxonomy), used to classify content. See&nbsp;[details](/content-management/taxonomies/).

### term page

Content with the "term" [page kind](#page-kind). Typically a listing of [regular pages](#regular-page) and [section pages](#section-page) with a given [term](#term).

### theme

A packaged combination of [archetypes](#archetype), assets, content, data, [templates](#template), translations, or configuration settings. A theme may serve as the basis for a new site, or to augment an existing site. See also [module](#module).

### token

An identifier within a format string, beginning with a colon and replaced with a value when rendered. For example, use tokens in format strings for both [permalinks](/content-management/urls/#permalinks) and [dates](/functions/dateformat/#datetime-formatting-layouts).


### type

See [content type](#content-type).

### variable

A variable initialized within a [template action](#template-action).

### walk

To recursively traverse a nested data structure. For example, rendering a multilevel menu.
