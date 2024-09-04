---
title: Features
description: Hugo's rich and powerful feature set provides the framework and tools to create static sites that build in seconds, often less.
categories: [about]
keywords: []
menu:
  docs:
    parent: about
    weight: 30
weight: 30
toc: true
---

## Framework

[Multiplatform]
: Install Hugo's single executable on Linux, macOS, Windows, and more.

[Multilingual]
: Localize your project for each language and region, including translations, images, dates, currencies, numbers, percentages, and collation sequence. Hugo's multilingual framework supports single-host and multihost configurations.

[Output formats]
: Render each page of your site to one or more output formats, with granular control by page kind, section, and path. While HTML is the default output format, you can add JSON, RSS, CSV, and more. For example, create a REST API to access content.

[Templates]
: Create templates using variables, functions, and methods to transform your content, resources, and data into a published page. While HTML templates are the most common, you can create templates for any output format.

[Themes]
: Reduce development time and cost by using one of the hundreds of themes contributed by the Hugo community. Themes are available for corporate sites, documentation projects, image portfolios, landing pages, personal and professional blogs, resumes, CVs, and more.

[Modules]
: Reduce development time and cost by creating or importing packaged combinations of archetypes, assets, content, data, templates, translation tables, static files, or configuration settings. A module may serve as the basis for a new site, or to augment an existing site.

[Privacy]
: Configure the behavior of Hugo's embedded templates and shortcodes to facilitate compliance with regional privacy regulations, including the [GDPR] and [CCPA].

[Security]
: Hugo's security model is based on the premise that template and configuration authors are trusted, but content authors are not. This model enables generation of HTML output safe against code injection. Other protections prevent "shelling out" to arbitrary applications, limit access to specific environment variables, prevent connections to arbitrary remote data sources, and more.

## Content authoring

[Content formats]
: Create your content using Markdown, HTML, AsciiDoc, Emacs Org Mode, Pandoc, or reStructuredText. Markdown is the default content format, conforming to the [CommonMark] and [GitHub Flavored Markdown] specifications.

[Markdown attributes]
: Apply HTML attributes such as `class` and `id` to Markdown images and block elements including blockquotes, fenced code blocks, headings, horizontal rules, lists, paragraphs, and tables.

[Markdown extensions]
: Leverage the embedded Markdown extensions to create tables, definition lists, footnotes, task lists, inserted text, mark text, subscripts, superscripts, and more.

[Markdown render hooks]
: Override the conversion of Markdown to HTML when rendering blockquotes, fenced code blocks, headings, images, links, and tables. For example, render every standalone image as an HTML `figure` element.

[Diagrams]
: Use fenced code blocks and Markdown render hooks to include diagrams in your content.

[Mathematics]
: Include mathematical equations and expressions in Markdown using LaTeX or TeX typesetting syntax.

[Syntax highlighting]
: Syntactically highlight code examples using Hugo's embedded syntax highlighter, enabled by default for fenced code blocks in Markdown. The syntax highlighter supports hundreds of code languages and dozens of styles.

[Shortcodes]
: Use Hugo's embedded shortcodes, or create your own, to insert complex content. For example, use shortcodes to include `audio` and `video` elements, render tables from local or remote data sources, insert snippets from other  pages, and more.

## Content management

[Content adapters]
: Create content adapters to dynamically add content when building your site. For example, use a content adapter to create pages from a remote data source such as JSON, TOML, YAML, or XML.

[Taxonomies]
: Classify content to establish simple or complex logical relationships between pages. For example, create an authors taxonomy, and assign one or more authors to each page. Among other uses, the taxonomy system provides an inverted, weighted index to render a list of related pages, ordered by relevance.

[Data]
: Augment your content using local or remote data sources including CSV, JSON, TOML, YAML, and XML. For example, create a shortcode to render an HTML table from a remote CSV file.

[Menus]
: Provide rapid access to content via Hugo's menu system, configured automatically, globally, or on a page-by-page basis. The menu system is a key component of Hugo's multilingual architecture.

[URL management]
: Serve any page from any path via global configuration or on a page-by-page basis.

## Asset pipelines

[Image processing]
: Convert, resize, crop, rotate,  adjust colors, apply filters, overlay text and images, and extract EXIF data.

[JavaScript bundling]
: Transpile TypeScript and JSX to JavaScript, bundle, tree shake, minify, create source maps, and perform SRI hashing.

[Sass processing]
: Transpile Sass to CSS, bundle, tree shake, minify, create source maps, perform SRI hashing, and integrate with PostCSS.

[Tailwind CSS processing]
: Compile Tailwind CSS utility classes into standard CSS, bundle, tree shake, optimize, minify, perform SRI hashing, and integrate with PostCSS.

## Performance

[Caching]
: Reduce build time and cost by rendering a partial template once then cache the result, either globally or within a given context. For example, cache the result of an asset pipeline to prevent reprocessing on every rendered page.

[Segmentation]
: Reduce build time and cost by partitioning your sites into segments. For example, render the home page and the "news section" every hour, and render the entire site once a week.

[Minification]
: Minify HTML, CSS, and JavaScript to reduce file size, bandwidth consumption, and loading times.

[CCPA]: https://en.wikipedia.org/wiki/California_Consumer_Privacy_Act
[Sass processing]: /functions/css/Sass/
[Caching]: /functions/partials/includecached/
[CommonMark]: https://spec.commonmark.org/current/
[Content adapters]: /content-management/content-adapters/
[Content formats]: /content-management/formats/
[Data]: /content-management/data-sources/
[Diagrams]: /content-management/diagrams/
[GDPR]: https://en.wikipedia.org/wiki/General_Data_Protection_Regulation
[GitHub Flavored Markdown]: https://github.github.com/gfm/
[Image processing]: /content-management/image-processing/
[JavaScript bundling]: /functions/js/build/
[Markdown attributes]: /content-management/markdown-attributes/
[Markdown extensions]: /getting-started/configuration-markup/#goldmark-extensions
[Markdown render hooks]: /render-hooks/introduction/
[Mathematics]: /content-management/mathematics/
[Menus]: /content-management/menus/
[Minification]: /getting-started/configuration/#configure-minify
[Modules]: https://gohugo.io/hugo-modules/
[Multilingual]: /content-management/multilingual/
[Multiplatform]: /installation/
[Output formats]: /templates/output-formats/
[Privacy]: /about/privacy/
[Security]: /about/security/
[Segmentation]: /getting-started/configuration/#configure-segments
[Shortcodes]: /content-management/shortcodes/
[Syntax highlighting]: /content-management/syntax-highlighting/
[Tailwind CSS processing]: /functions/css/tailwindcss/
[Taxonomies]: /content-management/taxonomies/
[Templates]: templates/introduction/
[Themes]: https://themes.gohugo.io/
[URL management]: /content-management/urls/
