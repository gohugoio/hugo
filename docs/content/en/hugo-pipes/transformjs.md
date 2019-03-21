---
title: TransformJS
description: Hugo Pipes can process JS files with Babel.
date: 2019-03-21
publishdate: 2019-03-21
lastmod: 2019-03-21
categories: [asset management]
keywords: []
menu:
  docs:
    parent: "pipes"
    weight: 75
weight: 75
sections_weight: 75
draft: false
---

Any JavaScript resource file can be transpiled to another JavaScript version using `resources.TransformJS` which takes for argument the resource object and a slice of options listed below. TransformJS uses the [babel cli](https://babeljs.io/docs/en/babel-cli).


{{% note %}}
Hugo Pipe's TranspileJS requires the `@babel/cli` and `@babel/core` JavaScript packages to be installed in the environment (`npm install -g @babel/cli @babel/core`) along with any Babel plugin(s) or preset(s) used (e.g., `npm install -g @babel/preset-env`).

If you are using the Hugo Snap package, Babel and plugin(s) need to be installed locally within your Hugo site directory, e.g., `npm install @babel/cli @babel/core` without the `-g` flag.
{{% /note %}}
### Options

plugins [string]
: Comma seperated string of Babel plugins to use

presets [string]
: Comma seperated string of Babel presets to use

minified [bool]
: Save as much bytes as possible when printing

noComments [bool]
: Write comments to generated output (true by default)

compact [string]
: Do not include superfluous whitespace characters and line terminators (true/false/auto)

verbose [bool]
: Log everything


```go-html-template
{{- $transpileOpts := (dict "plugins" "" "presets" "@babel/preset-env" "minified" true "noComments" true "compact" "true" "verbose" true ) -}}
{{- $transpiled := resources.Get "scripts/main.js" | transpileJS $transpileOpts -}}
```