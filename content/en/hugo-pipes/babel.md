---
title: Babel
description: Hugo Pipes can process JS files with Babel.
categories: [asset management]
keywords: []
menu:
  docs:
    parent: hugo-pipes
    weight: 70
weight: 70
function:
  aliases: [babel]
  returnType: resource.Resource
  signatures: ['js.Babel [OPTIONS] RESOURCE']
---

## Usage

Any JavaScript resource file can be transpiled to another JavaScript version using `js.Babel` which takes for argument the resource object and an optional dict of options listed below. Babel uses the [babel cli](https://babeljs.io/docs/en/babel-cli).

{{% note %}}
Hugo Pipe's Babel requires the `@babel/cli` and `@babel/core` JavaScript packages to be installed in the project or globally (`npm install -g @babel/cli @babel/core`) along with any Babel plugin(s) or preset(s) used (e.g., `npm install @babel/preset-env --save-dev`).

If you are using the Hugo Snap package, Babel and plugin(s) need to be installed locally within your Hugo site directory, e.g., `npm install @babel/cli @babel/core --save-dev` without the `-g` flag.
{{% /note %}}

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

## Examples

```go-html-template
{{- $transpiled := resources.Get "scripts/main.js" | babel  -}}
```

Or with options:

```go-html-template
{{ $opts := dict "noComments" true }}
{{- $transpiled := resources.Get "scripts/main.js" | babel $opts -}}
```
