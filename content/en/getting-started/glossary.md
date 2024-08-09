---
title: Glossary of terms
description: Terms commonly used throughout the documentation.
categories: [getting started]
keywords: [glossary]
menu:
  docs:
    parent: getting-started
    weight: 60
weight: 60
# Use level 6 headings for each term in the glossary.
---

[A](#action)&nbsp;
[B](#bool)&nbsp;
[C](#cache)&nbsp;
[D](#default-sort-order)&nbsp;
[E](#environment)&nbsp;
[F](#field)&nbsp;
[G](#global-resource)&nbsp;
[H](#headless-bundle)&nbsp;
[I](#identifier)&nbsp;
[K](#kind)&nbsp;
[L](#layout)&nbsp;
[M](#map)&nbsp;
[N](#node)&nbsp;
[O](#object)&nbsp;
[P](#page-bundle)&nbsp;
[R](#regular-page)&nbsp;
[S](#scalar)&nbsp;
[T](#taxonomic-weight)&nbsp;
[U](#unmarshal)&nbsp;
[V](#variable)&nbsp;
[W](#walk)&nbsp;
[Z](#zero-time)&nbsp;

###### action

See [template action](#template-action).

###### archetype

A template for new content. See&nbsp;[details](/content-management/archetypes/).

###### argument

A [scalar](#scalar), [array](#array), [slice](#slice), [map](#map), or [object](#object) passed to a [function](#function), [method](#method), or [shortcode](#shortcode).

###### array

A numbered sequence of elements. Unlike Go's [slice](#slice) data type, an array has a fixed length. [Elements](#element) within an array can be [scalars](#scalar), slices, [maps](#map), pages, or other arrays. See the [Go&nbsp;documentation](https://go.dev/ref/spec#Array_types) for details.

###### bool

See [boolean](#boolean).

###### boolean

A data type with two possible values, either `true` or `false`.

###### branch bundle

A directory that contains an _index.md file and zero or more [resources](#resource). Analogous to a physical branch, a branch bundle may have descendants including leaf bundles and other branch bundles. Top level directories with or without _index.md files are also branch bundles. This includes the home page. See&nbsp;[details](/content-management/page-bundles/).

###### build

To generate a static site that includes HTML files and assets such as images, CSS, and JavaScript. The build process includes rendering and resource transformations.

###### bundle

See [page bundle](#page-bundle).

###### cache

A software component that stores data so that future requests for the same data are faster.

###### chain

Within a template, to connect one or more [identifiers](#identifier) with a dot. An identifier can represent a method, object, or field. For example, `.Site.Params.author.name` or `.Date.UTC.Hour`.

###### CJK

A collective term for the Chinese, Japanese, and Korean languages. See [details](https://en.wikipedia.org/wiki/CJK_characters).

###### CLI

Command line interface.

###### collection

An [array](#array), [slice](#slice), or [map](#map).

###### content adapter

A template that dynamically creates pages when building a site. For example, use a content adapter to create pages from a remote data source such as JSON, TOML, YAML, or XML. See&nbsp;[details](/content-management/content-adapters/).

###### content format

A markup language for creating content. Typically Markdown, but may also be HTML, AsciiDoc, Org, Pandoc, or reStructuredText. See&nbsp;[details](/content-management/formats/).

###### content type

A classification of content inferred from the top-level directory name or the `type` set in [front matter](#front-matter). Pages in the root of the content directory, including the home page, are of type "page". Accessed via `.Page.Type` in [templates](#template). See&nbsp;[details](/content-management/types/).

###### content view

A template called with the `.Page.Render` method. See&nbsp;[details](/templates/content-view/).

###### context

Represented by a dot "." within a [template action](#template-action), context is the current location in a data structure. For example, while iterating over a [collection](#collection) of pages, the context within each iteration is the page's data structure. The context received by each template depends on template type and/or how it was called. See&nbsp;[details](/templates/introduction/#context).

###### default sort order

The default sort order for page collections. Hugo sorts by [weight](#weight), then by date (descending), then by link title, and then by file path.

###### element

A member of a slice or array.

###### environment

Typically one of `development`, `staging`, or `production`, each environment may exhibit different behavior depending on configuration and template logic. For example, in a production environment you might minify and fingerprint CSS, but that probably doesn't make sense in a development environment.

When running the built-in development server with the `hugo server` command, the environment is set to `development`. When building your site with the `hugo` command, the environment is set to `production`. To override the environment value, use the `--environment` command line flag or the `HUGO_ENVIRONMENT` environment variable.

To determine the current environment within a template, use the [`hugo.Environment`] function.

[`hugo.Environment`]: /functions/hugo/environment/

###### field

A predefined key-value pair in front matter such as `date` or `title`. See&nbsp;also&nbsp;[parameter](#parameter).

###### flag

An option passed to a command-line program, beginning with one or two hyphens. See&nbsp;[details](/commands/hugo/).

###### float

See [floating point](#floating-point).

###### floating point

A numeric data type with a fractional component. For example, `3.14159`.

###### fragment

The final segment of a URL, beginning with a hash (`#`) mark, that references an `id` attribute of an HTML element on the page.

###### front matter

Metadata at the beginning of each content page, separated from the content by format-specific delimiters. See&nbsp;[details](/content-management/front-matter/).

###### function

Used within a [template action](#template-action), a function takes one or more [arguments](#argument) and returns a value. Unlike [methods](#method), functions are not associated with an [object](#object). See&nbsp;[details](/functions/).

###### global resource

A file within the assets directory, or within any directory [mounted](/hugo-modules/configuration/#module-configuration-mounts) to the assets directory. Capture one or more global resources using the [`resources.Get`], [`resources.GetMatch`], [`resources.Match`], or [`resources.ByType`] functions.

[`resources.Get`]: /functions/resources/get/
[`resources.GetMatch`]: /functions/resources/getmatch/
[`resources.Match`]: /functions/resources/match/
[`resources.ByType`]: /functions/resources/byType/

###### headless bundle

An unpublished leaf or branch bundle whose content and resources you can include in other pages. See [build options](/content-management/build-options/).

###### identifier

A string that represents a variable, method, object, or field. It must conform to Go's [language specification](https://go.dev/ref/spec#Identifiers), beginning with a letter or underscore, followed by zero or more letters, digits, or underscores.

###### int

See [integer](#integer).

###### integer

A numeric data type without a fractional component. For example, `42`.

###### internationalization

Software design and development efforts that enable [localization](#localization). See the [W3C definition](https://www.w3.org/International/questions/qa-i18n). Abbreviated i18n.

###### interval

An [interval](https://en.wikipedia.org/wiki/Interval_(mathematics)) is a range of numbers between two endpoints: closed, open, or half-open.

- A _closed_ interval, denoted by brackets, includes its endpoints. For example, [0,&nbsp;1]&nbsp;is the interval where `0 <= x <= 1`.

- An _open_ interval, denoted by parentheses, excludes its endpoints. For example, (0,&nbsp;1)&nbsp;is the interval where `0 < x < 1`.

- A _half-open_ interval includes only one of its endpoints. For example, (0,&nbsp;1]&nbsp;is the _left-open_ interval where `0 < x <= 1`, while [0,&nbsp;1)&nbsp;is the _right-open_ interval where `0 <= x < 1`.

###### kind

See [page kind](#page-kind).

###### layout

See [template](#template).

###### leaf bundle

A directory that contains an index.md file and zero or more [resources](#resource). Analogous to a physical leaf, a leaf bundle is at the end of a branch. It has no descendants. See&nbsp;[details](/content-management/page-bundles/).

###### list page

Any [page kind](#page-kind) that receives a page [collection](#collection) in [context](#context). This includes the home page, [section pages](#section-page), [taxonomy pages](#taxonomy-page), and [term pages](#term-page).

###### list template

Any template that renders a [list page](#list-page). This includes [home](/templates/types/#home), [section](/templates/types/#section), [taxonomy](/templates/types/#taxonomy), and [term](/templates/types/#term) templates.

###### localization

Adaptation of a site to meet language and regional requirements. This includes translations, language-specific media, date and currency formats, etc. See&nbsp;[details](/content-management/multilingual/) and the [W3C definition](https://www.w3.org/International/questions/qa-i18n). Abbreviated l10n.

###### logical path

{{< new-in 0.123.0 >}}

A page or page resource identifier derived from the file path, excluding its extension and language identifier. This value is neither a file path nor a URL. Starting with a file path relative to the content directory, Hugo determines the logical path by stripping the file extension and language identifier, converting to lower case, then replacing spaces with hyphens. <!-- You may also set this value using the `path` front matter field. --> See [examples](/methods/page/path/#examples).

###### map

An unordered group of elements, each indexed by a unique key. See the [Go&nbsp;documentation](https://go.dev/ref/spec#Map_types) for details.

###### Markdown attribute

A list of attributes, containing one or more key-value pairs, separated by spaces or commas, and wrapped by braces. Apply Markdown attributes to images and block-level elements including blockquotes, fenced code blocks, headings, horizontal rules, lists, paragraphs, and tables. See&nbsp;[details](/getting-started/configuration-markup/#goldmark).

###### marshal

To transform a data structure into a serialized object. For example, transforming a [map](#map) into a JSON string. See [unmarshal](#unmarshal).

###### method

Used within a [template action](#template-action) and associated with an [object](#object), a method takes zero or more [arguments](#argument) and either returns a value or performs an action. For example, `.IsHome` is a method on the `.Page` object which returns `true` if the current page is the home page. See also [function](#function).

###### module

Like a [theme](#theme), a module is a packaged combination of [archetypes](#archetype), assets, content, data, [templates](#template), translation tables, static files, or configuration settings. A module may serve as the basis for a new site, or to augment an existing site. See&nbsp;[details](/hugo-modules/).

###### node

A class of [page kinds](#page-kind) including `home`, `section`, `taxonomy`, and `term`.

###### noop

An abbreviated form of "no operation", a _noop_ is a statement that does nothing.

###### object

A data structure with or without associated [methods](#method).

###### ordered taxonomy

Created by invoking the [`Alphabetical`] or [`ByCount`] method on a [`Taxonomy`](#taxonomy-object) object, which is a [map](#map), an ordered taxonomy is a [slice](#slice), where each element is an object that contains the [term](#term) and a slice of its [weighted pages](#weighted-page).

[`Alphabetical`]: /methods/taxonomy/alphabetical/
[`ByCount`]: /methods/taxonomy/bycount/

###### output format

{{% include "methods/page/_common/output-format-definition.md" %}}

###### page bundle

A directory that encapsulates both content and associated [resources](#resource). There are two types of page bundles: [leaf bundles](#leaf-bundle) and [branch bundles](#branch-bundle). See&nbsp;[details](/content-management/page-bundles/).

###### page collection

A slice of `Page` objects.

###### page kind

A classification of pages, one of `home`, `page`, `section`, `taxonomy`, or `term`. See&nbsp;[details](/methods/page/kind/).

Note that there are also `RSS`, `sitemap`, `robotsTXT`, and `404` page kinds, but these are only available during the rendering of each of these respective page's kind and therefore *not* available in any of the `Pages` collections.

###### page resource

A file within a [page bundle](#page-bundle). Capture one or more page resources using any of the [`Resources`] methods on a `Page` object.

[`Resources`]: /methods/page/resources/#methods

###### pager

Created during [pagination](#pagination), a pager contains a subset of a list page and navigation links to other pagers.

###### paginate

To split a list page into two or more subsets.

###### pagination

The process of [paginating](#paginate) a list page. See&nbsp;[details](/templates/pagination/).

###### paginator

A collection of [pagers](#pager).

###### parameter

Typically, a user-defined key-value pair at the site or page level, but may also refer to a configuration setting or an [argument](#argument). See&nbsp;also&nbsp;[field](#field).

###### partial

A [template](#template) called from any other template including [shortcodes](#shortcode), [render hooks](#render-hook), and other partials. A partial either renders something or returns something. A partial can also call itself, for example, to [walk](#walk) a data structure.

###### permalink

The absolute URL of a published resource or a rendered page, including scheme and host.

###### pipe

See [pipeline](#pipeline).

###### pipeline

Within a [template action](#template-action), a pipeline is a possibly chained sequence of values, [function](#function) calls, or [method](#method) calls. Functions and methods in the pipeline may take multiple [arguments](#argument).

A pipeline may be *chained* by separating a sequence of commands with pipeline characters "|". In a chained pipeline, the result of each command is passed as the last argument to the following command. The output of the final command in the pipeline is the value of the pipeline. See the [Go&nbsp;documentation](https://pkg.go.dev/text/template#hdr-Pipelines) for details.

###### publish

See [build](#build).

###### regular page

Content with the "page" [page kind](#page-kind). See also [section page](#section-page).

###### relative permalink

The host-relative URL of a published resource or a rendered page.

###### render hook

A [template](#template) that overrides standard Markdown rendering. See&nbsp;[details](/render-hooks).

###### remote resource

A file on a remote server, accessible via HTTP or HTTPS with the [`resources.GetRemote`](/functions/resources/getremote) function.

###### resource

Any file consumed by the build process to augment or generate content, structure, behavior, or presentation. For example: images, videos, content snippets, CSS, Sass, JavaScript, and data.

Hugo supports three types of resources: [global](#global-resource), [page](#page-resource), and [remote](#remote-resource)

###### resource type

The main type of a resource's [media type]. Content files such as Markdown, HTML, AsciiDoc, Pandoc, reStructuredText, and Emacs Org Mode have resource type `page`. Other resource types include `image`, `video`, etc. Retrieve the resource type using the [`ResourceType`] method on a `Resource` object.

[media type]: /methods/resource/mediatype/
[`ResourceType`]: /methods/resource/resourcetype/

###### scalar

A single value, one of [string](#string), [integer](#integer), [floating point](#floating-point), or [boolean](#boolean).

###### scratch pad

Conceptually, a [map](#map) with [methods](#method) to set, get, update, and delete values. Attach the data structure to a `Page` object using the [`Scratch`] or [`Store`] methods, or create a locally scoped scratch pad using the [`newScratch`] function.

[`Scratch`]: /methods/page/scratch/
[`Store`]: /methods/page/store/
[`newScratch`]: /functions/collections/newscratch/

###### section

A top-level content directory, or any content directory with an&nbsp;_index.md file. A content directory with an&nbsp;_index.md file is also known as a [branch bundle](/getting-started/glossary/#branch-bundle). Section templates receive one or more page [collections](#collection) in [context](#context). See&nbsp;[details](/content-management/sections/).

###### section page

Content with the "section" [page kind](#page-kind). Typically a listing of [regular pages](#regular-page) and/or [section pages](#section-page) within the current [section](#section). See also [regular page](#regular-page).

###### shortcode

A [template](#template) called from within Markdown, taking zero or more [arguments](#argument). See&nbsp;[details](/content-management/shortcodes/).

###### slice

A numbered sequence of elements. Unlike Go's [array](#array) data type, slices are dynamically sized. [Elements](#element) within a slice can be [scalars](#scalar), [arrays](#array), [maps](#map), pages, or other slices. See the [Go&nbsp;documentation](https://go.dev/ref/spec#Slice_types) for details.

###### string

A sequence of bytes. For example, `"What is 6 times 7?"`&nbsp;.

###### string literal (interpreted)

Interpreted string literals are character sequences between double quotes, as in "foo". Within the quotes, any character may appear except a newline and an unescaped double quote. The text between the quotes forms the value of the literal, with backslash escapes interpreted. See [details](https://go.dev/ref/spec#String_literals).

###### string literal (raw)

Raw string literals are character sequences between backticks, as in \`bar\`. Within the backticks, any character may appear except a backtick. Backslashes have no special meaning and the string may contain newlines. Carriage return characters ('\r') inside raw string literals are discarded from the raw string value. See [details](https://go.dev/ref/spec#String_literals).

###### taxonomic weight

Defined in front matter and unique to each taxonomy, this [weight](#weight) determines the sort order of page collections contained within a [`Taxonomy`](#taxonomy-object) object. See&nbsp;[details](/content-management/taxonomies/#order-taxonomies).

###### taxonomy

A group of related [terms](#term) used to classify content. For example, a "colors" taxonomy might include the terms "red", "green", and "blue". See&nbsp;[details](/content-management/taxonomies/).

###### taxonomy object

A [map](#map) of [terms](#term) and the [weighted pages](#weighted-page) associated with each term.

###### taxonomy page

Content with the "taxonomy" [page kind](#page-kind). Typically a listing of [terms](#term) within a given [taxonomy](#taxonomy).

###### template

A file with [template actions](#template-action), located within the layouts directory of a project, theme, or module. See&nbsp;[details](/templates/).

###### template action

A data evaluation or control structure within a [template](#template), delimited by "{{"&nbsp;and&nbsp;"}}". See the [Go&nbsp;documentation](https://pkg.go.dev/text/template#hdr-Actions) for details.

###### term

A member of a [taxonomy](#taxonomy), used to classify content. See&nbsp;[details](/content-management/taxonomies/).

###### term page

Content with the "term" [page kind](#page-kind). Typically a listing of [regular pages](#regular-page) and [section pages](#section-page) with a given [term](#term).

###### theme

A packaged combination of [archetypes](#archetype), assets, content, data, [templates](#template), translation tables, static files, or configuration settings. A theme may serve as the basis for a new site, or to augment an existing site. See also [module](#module).

###### token

An identifier within a format string, beginning with a colon and replaced with a value when rendered. For example, use tokens in format strings for both [permalinks](/content-management/urls/#permalinks) and [dates](/functions/time/format/#localization).

###### type

See [content type](#content-type).

###### unmarshal

To transform a serialized object into a data structure. For example, transforming a JSON file into a [map](#map) that you can access within a template. See [marshal](#marshal).

###### variable

A user-defined [identifier](#identifier) prepended with a `$` symbol, representing a value of any data type, initialized or assigned within a [template action](#template-action). For example, `$foo`&nbsp;and&nbsp;`$bar` are variables.

###### walk

To recursively traverse a nested data structure. For example, rendering a multilevel menu.

###### weight

Used to position an element within a collection sorted by weight. Assign weights using non-zero integers. Lighter items float to the top, while heavier items sink to the bottom. Unweighted or zero-weighted elements are placed at the end of the collection. Weights are typically assigned to pages, menu entries, languages, and output formats.

###### weighted page

Contained within a [`Taxonomy`](#taxonomy-object) object, a weighted page is a [map](#map) with two elements: a `Page` object, and its [taxonomic weight](#taxonomic-weight) as defined in front matter. Access the elements using the `Page` and `Weight` keys.

###### zero time

The _zero time_ is January 1, 0001, 00:00:00 UTC. Formatted per [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) the _zero time_ is 0001-01-01T00:00:00-00:00.
