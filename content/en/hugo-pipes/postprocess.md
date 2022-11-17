---
title: PostProcess
description: Allows delaying of resource transformations to after the build.
date: 2020-04-09
categories: [asset management]
keywords: []
menu:
  docs:
    parent: "pipes"
    weight: 39
weight: 39
sections_weight: 39
---

Marking a resource with `resources.PostProcess` delays any transformations to after the build, typically because one or more of the steps in the transformation chain depends on the result of the build (e.g. files in `public`).

A prime use case for this is [CSS purging with PostCSS](#css-purging-with-postcss).

There are currently two limitations to this:

1. This only works in `*.html` templates (i.e. templates that produces HTML files).
2. You cannot manipulate the values returned from the resource's methods. E.g. the `upper` in this example will not work as expected:

    ```go-html-template
    {{ $css := resources.Get "css/main.css" }}
    {{ $css = $css | resources.PostCSS | minify | fingerprint | resources.PostProcess }}
    {{ $css.RelPermalink | upper }}
    ```

## CSS purging with PostCSS

{{% note %}}
There are several ways to set up CSS purging with PostCSS in Hugo. If you have a simple project, you should consider going the simpler route and drop the use of `resources.PostProcess` and just extract keywords from the templates. See the [Tailwind documentation](https://tailwindcss.com/docs/controlling-file-size/#app) for some examples.
{{% /note %}}

The below configuration will write a `hugo_stats.json` file to the project root as part of the build. If you're only using this for the production build, you should consider placing it below [config/production](/getting-started/configuration/#configuration-directory).

{{< code-toggle file="config" >}}
[build]
  writeStats = true
{{< /code-toggle >}}

`postcss.config.js`

```js
const purgecss = require('@fullhuman/postcss-purgecss')({
    content: [ './hugo_stats.json' ],
    defaultExtractor: (content) => {
        let els = JSON.parse(content).htmlElements;
        return els.tags.concat(els.classes, els.ids);
    }
});

module.exports = {
     plugins: [
         ...(process.env.HUGO_ENVIRONMENT === 'production' ? [ purgecss ] : [])
     ]
 };
```

Note that in the example above, the "CSS purge step" will only be applied to the production build. This means that you need to do something like this in your head template to build and include your CSS:

```go-html-template
{{ $css := resources.Get "css/main.css" }}
{{ $css = $css | resources.PostCSS }}
{{ if hugo.IsProduction }}
{{ $css = $css | minify | fingerprint | resources.PostProcess }}
{{ end }}
<link href="{{ $css.RelPermalink }}" rel="stylesheet" />
```
