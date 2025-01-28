---
title: js.Babel
description: Compile the given JavaScript resource with Babel.
categories: []
keywords: []
action:
  aliases: [babel]
  related:
    - functions/js/Batch
    - functions/js/Build
    - functions/resources/Fingerprint
    - functions/resources/Minify
  returnType: resource.Resource
  signatures: ['js.Babel [OPTIONS] RESOURCE']
weight: 30
toc: true
---

```go-html-template
{{ with resources.Get "js/main.js" }}
  {{ $opts := dict
    "minified" hugo.IsProduction
    "noComments" hugo.IsProduction
    "sourceMap" (cond hugo.IsProduction "none" "external")
  }}
  {{ with . | js.Babel $opts }}
    {{ if hugo.IsProduction }}
      {{ with . | fingerprint }}
        <script src="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous"></script>
      {{ end }}
    {{ else }}
      <script src="{{ .RelPermalink }}"></script>
    {{ end }}
  {{ end }}
{{ end }}
```

## Setup

Step 1
: Install [Node.js](https://nodejs.org/en/download)

Step 2
: Install the required Node.js packages in the root of your project.

```sh
npm install --save-dev @babel/core @babel/cli
```

Step 3
: Add the babel executable to Hugo's `security.exec.allow` list in your site configuration:

{{< code-toggle file=hugo >}}
[security.exec]
  allow = ['^(dart-)?sass(-embedded)?$', '^go$', '^npx$', '^postcss$', '^babel$']
{{< /code-toggle >}}

## Configuration

We add the main project's `node_modules` to `NODE_PATH` when running Babel and similar tools. There are some known [issues](https://github.com/babel/babel/issues/5618) with Babel in this area, so if you have a `babel.config.js` living in a Hugo Module (and not in the project itself), we recommend using `require` to load the presets/plugins, e.g.:

```js
module.exports = {
  presets: [
    [
      require("@babel/preset-env"),
      {
        useBuiltIns: "entry",
        corejs: 3,
      },
    ],
  ],
};
```

## Options

###### compact

(`bool`) Whether to remove optional newlines and whitespace. Enabled when `minified` is `true`. Default is `false`

###### config

(`string`) Path to the Babel configuration file. Hugo will, by default, look for a `babel.config.js` file in the root of your project. See [details](https://babeljs.io/docs/en/configuration).

###### minified

(`bool`) Whether to minify the compiled code. Enables the `compact` option. Default is `false`.

###### noBabelrc

(`string`) Whether to ignore `.babelrc` and `.babelignore` files. Default is `false`.

###### noComments

(`bool`) Whether to remove comments. Default is `false`.

###### sourceMap

(`string`) Whether to generate source maps, one of `external`, `inline`, or `none`. Default is `none`.

<!-- In the above, technically "none" is not one of the enumerated values, but it has the same effect and is easier to document than an empty string. -->

###### verbose

(`bool`) Whether to enable verbose logging. Default is `fasle`
