---
title: js.Babel
description: Transpile JavaScript resources using Babel.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [babel]
    returnType: resource.Resource
    signatures: ['js.Babel [OPTIONS] RESOURCE']
aliases: [/functions/resources/babel/]
---

The `js.Babel` function transforms JavaScript using [Babel][].

## Setup

Step 1
: Install [Node.js][].

Step 2
: Install the required Node packages in the root of your project. For example, to install Babel's core compiler, its command-line interface, and the preset for transpiling modern JavaScript based on your target environments:

  ```sh
  npm install --save-dev @babel/core @babel/cli @babel/preset-env
  ```

Step 3
: Create a Babel configuration file in the root of your project. For example, to use the environment preset to target Google Chrome version 79 or later:

  ```js {file="babel.config.mjs" copy=true}
  export default {
    presets: [
      [
        '@babel/preset-env',
        {
          targets: {
            chrome: "79"
          }
        }
      ]
    ]
  };
  ```

Step 4
: Place your JS file within the `assets/js` directory.

Step 5
: Add the Babel executable to Hugo's `security.exec.allow` list in your project configuration:

  {{< code-toggle file=hugo >}}
  [security.exec]
    allow = ['^(dart-)?sass(-embedded)?$', '^go$', '^git$', '^node$', '^postcss$', '^tailwindcss$', '^babel$']
  {{< /code-toggle >}}

Step 6
: Create a partial template to process the JavaScript:

  ```go-html-template {file="layouts/_partials/js.html" copy=true}
  {{ with resources.Get "js/main.js" }}
    {{ $opts := dict
      "minified" (cond hugo.IsDevelopment false true)
      "noComments" (cond hugo.IsDevelopment false true)
      "sourceMap" (cond hugo.IsDevelopment "inline" "none")
    }}
    {{ with . | js.Babel $opts }}
      {{ if hugo.IsDevelopment }}
        <script src="{{ .RelPermalink }}"></script>
      {{ else }}
        {{ with . | fingerprint }}
          <script src="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous"></script>
        {{ end }}
      {{ end }}
    {{ end }}
  {{ end }}
  ```

Step 7
: Call the partial template from your base template:

  ```go-html-template {file="layouts/baseof.html" copy=true}
  <head>
    {{ partial "js.html" . }}
  </head>
  ```

## Options

The `js.Babel` function accepts an options map.

`compact`
: (`bool`) Whether to remove optional newlines and whitespace. Enabled when `minified` is `true`. Default is `false`.

`config`
: (`string`) The path to the Babel configuration file. By default, Hugo searches the root of the project directory and any modules for `babel.config.js`, `babel.config.mjs`, and `babel.config.cjs` in that order. Use this option only to point to a configuration file with a custom name or one located in a custom subdirectory.

`minified`
: (`bool`) Whether to minify transpiled code. Enables the `compact` option. Default is `false`.

`noBabelrc`
: (`bool`) Whether to ignore `.babelrc` and `.babelignore` files. Default is `false`.

`noComments`
: (`bool`) Whether to remove comments. Default is `false`.

`sourceMap`
: (`string`) Whether to generate source maps, one of `external`, `inline`, or `none`. Default is `none`.

`verbose`
: (`bool`) Whether to enable verbose logging. Default is `false`.

<!--
In the above, technically "none" is not one of the enumerated sourceMap
values but it has the same effect and is easier to document than an empty string.
-->

[Babel]: https://babeljs.io/
[Node.js]: https://nodejs.org/en/download
