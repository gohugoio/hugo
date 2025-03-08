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

{{< new-in "0.145.0" />}}

[Portable Text](https://www.portabletext.org/) is a JSON structure that represent rich text content in the [Sanity](https://www.sanity.io/) CMS. In Hugo, this function is typically used in a [Content Adapter](https://gohugo.io/content-management/content-adapters/) that creates pages from Sanity data.

## Types supported

- `block` and `span`
- `image`. Note that the image handling is currently very simple; we link to the `asset.url` using `asset.altText` as the image alt text and `asset.title` as the title. For more fine grained control you may want to process the images in a [image render hook](/render-hooks/images/).
- `code` (see the [code-input](https://www.sanity.io/plugins/code-input) plugin). Code will be rendered as [fenced code blocks](/contribute/documentation/#fenced-code-blocks) with any file name provided passed on as a markdown attribute.

> [!note]
> Since the Portable Text gets converted to Markdown before it gets passed to Hugo, rendering of links, headings, images and code blocks can be controlled with [Render Hooks](https://gohugo.io/render-hooks/).

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
{{ $url := printf "https://%s.%s.sanity.io/v2021-06-07/data/query/production"  $projectID $api }}

{{/* prettier-ignore-start */ -}}
{{ $q :=  `*[_type == 'post']{
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
{{/* prettier-ignore-end */ -}}
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

Below outlines a suitable Sanity studio setup for the above example.

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

```bash
npm i sanity-plugin-media @sanity/code-input
```

```ts {file="schemaTypes/index.ts" copy=true}
import {postType} from './postType'

export const schemaTypes = [postType]
```

## Server setup

Unfortunately, Sanity's API does not support [RFC 7234](https://tools.ietf.org/html/rfc7234) and their output changes even if the data has not. A recommended setup is therefore to use their cached `apicdn` endpoint (see above) and then set up a reasonable polling and file cache strategy in your Hugo configuration, e.g:

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

The polling above will be used when running the server/watch mode and rebuild when you push new content in Sanity.

See [Caching in resources.GetRemote](/functions/resources/getremote/#caching) for more fine grained control.
