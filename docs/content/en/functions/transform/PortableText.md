---
title: transform.PortableText
description: Converts Portable Text to Markdown.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [transform.PortableText MAP]
---

{{< new-in 0.145.0 />}}

[Portable Text][] is a JSON structure that represents rich text content in the [Sanity][] CMS. In Hugo, this function is typically used in a [content adapter][] that creates pages from Sanity data.

## Types supported

- `block` and `span`
- `image`. Note that the image handling is currently basic; we link to the `asset.url` using `asset.altText` as the image alt text and `asset.title` as the title. For more fine-grained control you may want to process the images in an [image render hook][].
- `code` (see the [code-input][] plugin). Code will be rendered as fenced code blocks with any file name provided passed as a Markdown attribute.

> [!NOTE]
> Since the Portable Text gets converted to Markdown before it gets passed to Hugo, rendering of links, headings, images and code blocks can be controlled with [render hooks][].

## Example

### Content Adapter

```go-html-template {file="content/_content.gotmpl" copy=true}
{{ $projectID := "mysanityprojectid" }}
{{ $useCached := true }}
{{ $api := "api" }}
{{ if $useCached }}
  {{/* See https://www.sanity.io/docs/api-cdn */}}
  {{ $api = "apicdn" }}
{{ end }}
{{ $url := printf "https://%s.%s.sanity.io/v2021-06-07/data/query/production" $projectID $api }}

{{ $q := `*[_type == 'post']{
  title, publishedAt, summary, slug, body[]{
    ...,
    _type == "image" => {
      ...,
      asset->{
        _id,
        path,
        url,
        altText,
        title,
        description,
        metadata {
          dimensions {
            aspectRatio,
            width,
            height
          }
        }
      }
    }
  },
  }`
}}
{{ $body := dict "query" $q | jsonify }}
{{ $opts := dict "method" "post" "body" $body }}
{{ $r := resources.GetRemote $url $opts }}
{{ $m := $r | transform.Unmarshal }}
{{ $result := $m.result }}
{{ range $result }}
  {{ if not .slug }}
    {{ continue }}
  {{ end }}
  {{ $markdown := transform.PortableText .body }}
  {{ $content := dict
    "mediaType" "text/markdown"
    "value" $markdown
  }}
  {{ $params := dict
    "portabletext" (.body | jsonify (dict "indent" " "))
  }}
  {{ $page := dict
    "content" $content
    "kind" "page"
    "path" .slug.current
    "title" .title
    "date" (.publishedAt | time )
    "summary" .summary
    "params" $params
  }}
  {{ $.AddPage $page }}
{{ end }}
```

### Sanity setup

The following outlines a suitable Sanity studio setup for the above example.

```ts {file="sanity.config.ts" copy=true}
import {defineConfig} from 'sanity'
import {structureTool} from 'sanity/structure'
import {visionTool} from '@sanity/vision'
import {schemaTypes} from './schemaTypes'
import {media} from 'sanity-plugin-media'
import {codeInput} from '@sanity/code-input'

export default defineConfig({
  name: 'default',
  title: 'my-sanity-project',

  projectId: 'mysanityprojectid',
  dataset: 'production',

  plugins: [structureTool(), visionTool(), media(),codeInput()],

  schema: {
    types: schemaTypes,
  },
})
```

Type/schema definition:

```ts {file="schemaTypes/postType.ts" copy=true}
import {defineField, defineType} from 'sanity'

export const postType = defineType({
  name: 'post',
  title: 'Post',
  type: 'document',
  fields: [
    defineField({
      name: 'title',
      type: 'string',
      validation: (rule) => rule.required(),
    }),
    defineField({
      name: 'summary',
      type: 'string',
      validation: (rule) => rule.required(),
    }),
    defineField({
      name: 'slug',
      type: 'slug',
      options: {source: 'title'},
      validation: (rule) => rule.required(),
    }),
    defineField({
      name: 'publishedAt',
      type: 'datetime',
      initialValue: () => new Date().toISOString(),
      validation: (rule) => rule.required(),
    }),
    defineField({
      name: 'body',
      type: 'array',
      of: [
        {
          type: 'block',
        },
        {
          type: 'image'
        },
        {
          type: 'code',
          options: {
            language: 'css',
            languageAlternatives: [
              {title: 'HTML', value: 'html'},
              {title: 'CSS', value: 'css'},
            ],
            withFilename: true,
          },
        },
      ],
    }),
  ],
})
```

Note that the above requires some additional plugins to be installed:

```sh
npm i sanity-plugin-media @sanity/code-input
```

```ts {file="schemaTypes/index.ts" copy=true}
import {postType} from './postType'

export const schemaTypes = [postType]
```

## Server setup

Unfortunately, Sanity's API does not support [RFC 7234][] and their output changes even if the data has not. A recommended setup is therefore to use their cached `apicdn` endpoint (see above) and then set up a reasonable polling and file cache strategy in your Hugo configuration, e.g:

<!-- markdownlint-disable MD049 -->
{{< code-toggle file=hugo >}}
[HTTPCache]
  [[HTTPCache.polls]]
    disable = false
    low = '30s'
    high = '3m'
    [HTTPCache.polls.for]
      includes = ['https://*.*.sanity.io/**']

[caches.getresource]
    dir    = ':cacheDir/:project'
    maxAge = "5m"
{{< /code-toggle >}}
<!-- markdownlint-enable MD049 -->

The polling above will be used when running the server/watch mode and rebuilds when you push new content to Sanity.

See [Caching in resources.GetRemote][] for more fine-grained control.

[Caching in resources.GetRemote]: /functions/resources/getremote/#caching
[Portable Text]: https://www.portabletext.org/
[RFC 7234]: https://tools.ietf.org/html/rfc7234
[Sanity]: https://www.sanity.io/
[code-input]: https://www.sanity.io/plugins/code-input
[content adapter]: /content-management/content-adapters/
[image render hook]: /render-hooks/images/
[render hooks]: /render-hooks/
