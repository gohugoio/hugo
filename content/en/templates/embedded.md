---
title: Embedded partial templates
description: Hugo provides embedded partial templates for common use cases.
categories: []
keywords: []
weight: 180
aliases: [/templates/internal]
---

## Disqus

> [!NOTE]
> To override Hugo's embedded Disqus template, copy the [source code](<{{% eturl disqus %}}>) to a file with the same name in the `layouts/_partials` directory, then call it from your templates using the [`partial`][] function:
>
> `{{ partial "disqus.html" . }}`

Hugo includes an embedded template for [Disqus][], a commenting system for both static and dynamic websites. To use this template, you must first obtain a Disqus shortname by [signing up][] for the free service.

To include the embedded template:

```go-html-template
{{ partial "disqus.html" . }}
```

### Configuration {#configuration-disqus}

To use Hugo's Disqus template, first set up a single configuration value:

{{< code-toggle file=hugo >}}
[services.disqus]
shortname = 'your-disqus-shortname'
{{</ code-toggle >}}

You can also set the following in the [front matter][] for a given page:

`disqus_identifier`
: (`string`) A unique identifier for the page's discussion thread. Set this to preserve comment threads across URL changes.

`disqus_title`
: (`string`) The title of the discussion thread.

`disqus_url`
: (`string`) The canonical URL for the discussion thread. Use this to override the URL Disqus uses to identify the thread, for example when the same content is served at multiple URLs.

{{< code-toggle file=content/blog/my-post.md fm=true >}}
[params]
disqus_identifier = 'unique-identifier'
disqus_title = 'Post title'
disqus_url = 'https://example.org/blog/my-post/'
{{</ code-toggle >}}

> [!NOTE]
> When previewing your site locally, Hugo replaces the Disqus widget with the message "Disqus comments not available by default when the website is previewed locally."

### Privacy {#privacy-disqus}

Adjust the relevant privacy settings in your project configuration.

{{< code-toggle config=privacy.disqus />}}

`disable`
: (`bool`) Whether to disable the template. Default is `false`.

## Google Analytics

> [!NOTE]
> To override Hugo's embedded Google Analytics template, copy the [source code](<{{% eturl google_analytics %}}>) to a file with the same name in the `layouts/_partials` directory, then call it from your templates using the [`partial`][] function:
>
> `{{ partial "google_analytics.html" . }}`

Hugo includes an embedded template for [Google Analytics 4][].

To include the embedded template:

```go-html-template
{{ partial "google_analytics.html" . }}
```

### Configuration {#configuration-google-analytics}

Provide your tracking ID in your configuration file:

{{< code-toggle file=hugo >}}
[services.googleAnalytics]
id = 'G-MEASUREMENT_ID'
{{</ code-toggle >}}

> [!NOTE]
> If the configured ID begins with `ua-` (case-insensitive), Hugo logs a warning and renders nothing. Google Universal Analytics (UA) was replaced by Google Analytics 4 (GA4) effective 1 July 2023. Create a GA4 property and data stream, then update your project configuration with the new measurement ID.

### Privacy {#privacy-google-analytics}

Adjust the relevant privacy settings in your project configuration.

{{< code-toggle config=privacy.googleAnalytics />}}

`disable`
: (`bool`) Whether to disable the template. Default is `false`.

`respectDoNotTrack`
: (`bool`) Whether to respect the browser's "do not track" setting. Default is `true`.

## Open Graph

> [!NOTE]
> To override Hugo's embedded Open Graph template, copy the [source code](<{{% eturl opengraph %}}>) to a file with the same name in the `layouts/_partials` directory, then call it from your templates using the [`partial`][] function:
>
> `{{ partial "opengraph.html" . }}`

Hugo includes an embedded template for the [Open Graph protocol][]. This metadata transforms your pages into rich objects when shared across major social media and messaging platforms.

To include the embedded template:

```go-html-template
{{ partial "opengraph.html" . }}
```

### Configuration {#configuration-open-graph}

Hugo's Open Graph template is configured using a mix of configuration settings and [front matter][] values on individual pages.

{{< code-toggle file=hugo >}}
title = 'My cool site'
[params]
  description = 'Text about my cool site'
  images = ['site-feature-image.jpg']
  [params.social]
  facebook_app_id = '12345678'
[taxonomies]
  series = 'series'
{{</ code-toggle >}}

{{< code-toggle file=content/blog/my-post.md fm=true >}}
title = 'Post title'
description = 'Text about this post'
date = 2024-03-08T08:18:11-08:00
images = ["post-cover.png"]
audio = []
videos = []
series = []
tags = []
locale = 'en-US'
{{</ code-toggle >}}

### Metadata {#metadata-open-graph}

Hugo emits the following metadata:

`og:url`
: The page permalink.

`og:site_name`
: The site title, falling back to the site configuration's `params.title` value.

`og:title`
: The page title, falling back to the site title, then the site configuration's `params.title` value.

`og:description`
: The page description, falling back to the page summary, then the site configuration's `params.description` value.

`og:locale`
: The `locale` front matter value, falling back to the site language's `locale`; hyphens are replaced with underscores (e.g. `en-US` → `en_US`).

`og:type`
: The value is `article` for pages and `website` for list and home pages.

For article pages, Hugo also emits:

`article:section`
: The page's top-level section.

`article:published_time`
: The page's publish date.

`article:modified_time`
: The page's last modified date.

`article:tag`
: The first 6 tags.

For image metadata, Hugo emits up to 6 `og:image` tags.

{{% include "/_common/embedded-get-page-images.md" %}}

`audio` and `videos` are `[]string` front matter parameters. Hugo emits up to 6 `og:audio` and `og:video` tags, passing each value through `absURL`, which converts relative paths to absolute URLs. Unlike `images`, Hugo does not search page resources or global resources for these values.

The `series` taxonomy is used to populate `og:see_also` metadata. Hugo emits up to 7 `og:see_also` tags using the first 7 pages in the same series as the current page, excluding the current page itself.

For Facebook metadata, if the site configuration's `params.social.facebook_app_id` value is set, Hugo emits `fb:app_id`. Otherwise, if the site configuration's `params.social.facebook_admin` value is set, Hugo emits `fb:admins`.

## Pagination

> [!NOTE]
> To override Hugo's embedded pagination template, copy the [source code](<{{% eturl pagination %}}>) to a file with the same name in the `layouts/_partials` directory, then call it from your templates using the [`partial`][] function:
>
> `{{ partial "pagination.html" . }}`

Hugo includes an embedded template for rendering navigation links between pagers. To include the embedded template:

```go-html-template
{{ partial "pagination.html" . }}
```

The embedded pagination template has two formats: `default` and `terse`. The `terse` format has fewer controls and page slots, consuming less space when styled as a horizontal list. See [pagination][] for details.

## Schema

> [!NOTE]
> To override Hugo's embedded Schema template, copy the [source code](<{{% eturl schema %}}>) to a file with the same name in the `layouts/_partials` directory, then call it from your templates using the [`partial`][] function:
>
> `{{ partial "schema.html" . }}`

Hugo includes an embedded template to render [microdata][] `meta` elements within the `head` element of your templates.

To include the embedded template:

```go-html-template
{{ partial "schema.html" . }}
```

### Configuration {#configuration-schema}

Hugo's Schema template uses a mix of page data and [front matter][] values on individual pages.

{{< code-toggle file=hugo >}}
title = 'My cool site'
[params]
  description = 'Text about my cool site'
{{</ code-toggle >}}

{{< code-toggle file=content/blog/my-post.md fm=true >}}
title = 'Post title'
description = 'Text about this post'
date = 2024-03-08T08:18:11-08:00
lastmod = 2024-03-09T12:00:00-08:00
images = ['post-cover.png']
keywords = ['ssg', 'hugo']
{{</ code-toggle >}}

### Metadata {#metadata-schema}

Hugo emits the following microdata:

`name`
: The page title, falling back to the site title.

`description`
: The page description, falling back to the page summary, then the site configuration's `params.description` value.

`datePublished`
: The page's publish date.

`dateModified`
: The page's last modified date.

`wordCount`
: The page's word count.

For image metadata, Hugo emits up to 6 `image` tags.

{{% include "/_common/embedded-get-page-images.md" %}}

For keyword metadata, Hugo uses the following order of precedence:

1. Titles of `keywords` taxonomy terms, if `keywords` is defined as a taxonomy
1. The `keywords` front matter value, if `keywords` is not a taxonomy
1. Titles of `tags` taxonomy terms
1. Titles of all taxonomy terms

## X (Twitter) Cards

> [!NOTE]
> To override Hugo's embedded Twitter Cards template, copy the [source code](<{{% eturl twitter_cards %}}>) to a file with the same name in the `layouts/_partials` directory, then call it from your templates using the [`partial`][] function:
>
> `{{ partial "twitter_cards.html" . }}`

Hugo includes an embedded template for [X (Twitter) Cards][], metadata used to attach rich media to Tweets linking to your site.

To include the embedded template:

```go-html-template
{{ partial "twitter_cards.html" . }}
```

### Configuration {#configuration-x-cards}

Hugo's X (Twitter) Card template is configured using a mix of configuration settings and [front matter][] values on individual pages.

{{< code-toggle file=hugo >}}
[params]
  description = 'Text about my cool site'
  images = ["site-feature-image.jpg"]
  [params.social]
  twitter = 'GoHugoIO'
{{</ code-toggle >}}

{{< code-toggle file=content/blog/my-post.md fm=true >}}
title = 'Post title'
description = 'Text about this post'
images = ["post-cover.png"]
{{</ code-toggle >}}

### Metadata {#metadata-x-cards}

If an image is found, Hugo sets `twitter:card` to `summary_large_image` and emits a `twitter:image` tag using the first image found. If no image is found, Hugo sets `twitter:card` to `summary` and omits the image tag.

{{% include "/_common/embedded-get-page-images.md" %}}

Hugo also emits the following metadata:

`twitter:title`
: The page title, falling back to the site title, then the site configuration's `params.title` value.

`twitter:description`
: The page description, falling back to the page summary, then the site configuration's `params.description` value.

`twitter:site`
: The site configuration's `params.social.twitter` value. The `@` prefix is added automatically if not already present. For example, with `twitter = 'GoHugoIO'` in your configuration, Hugo renders:

  ```html
  <meta name="twitter:site" content="@GoHugoIO"/>
  ```

[Disqus]: https://disqus.com
[Google Analytics 4]: https://support.google.com/analytics/answer/10089681
[Open Graph protocol]: https://ogp.me/
[X (Twitter) Cards]: https://developer.x.com/en/docs/twitter-for-websites/cards/overview/abouts-cards
[`partial`]: /functions/partials/include/
[front matter]: /content-management/front-matter/
[microdata]: https://html.spec.whatwg.org/multipage/microdata.html#microdata
[pagination]: /templates/pagination/
[signing up]: https://disqus.com/profile/signup/
