---
title: js.Babel
description: Compiles the given JavaScript resource with Babel.
categories: []
keywords: []
action:
  aliases: [babel]
  related:
    - functions/js/Build
    - functions/resources/Fingerprint
    - functions/resources/Minify
  returnType: resource.Resource
  signatures: ['js.Babel [OPTIONS] RESOURCE']
toc: true
---

{{< new-in 0.128.0 >}}

```go-html-template
{{ with resources.Get "js/main.js" }}
  {{ if hugo.IsDevelopment }}
    {{ with . | babel }}
      <script src="{{ .RelPermalink }}"></script>
    {{ end }}
  {{ else }}
    {{ $opts := dict "minified" true }}
    {{ with . | babel $opts | fingerprint }}
      <script src="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous"></script>
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

config
: (`string`) Path to the Babel configuration file. Hugo will, by default, look for a `babel.config.js` in your project. More information on these configuration files can be found here: [babel configuration](https://babeljs.io/docs/en/configuration).

minified
: (`bool`) Save as many bytes as possible when printing

noComments
: (`bool`) Write comments to generated output (true by default)

compact
: (`bool`) Do not include superfluous whitespace characters and line terminators. Defaults to `auto` if not set.

verbose
: (`bool`) Log everything

sourceMap
: (`string`) Output `inline` or `external` sourcemap from the babel compile. External sourcemaps will be written to the target with the output file name + ".map". Input sourcemaps can be read from js.Build and node modules and combined into the output sourcemaps.
