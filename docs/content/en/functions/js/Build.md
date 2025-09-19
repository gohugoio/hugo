---
title: js.Build
description: Bundle, transpile, tree shake, and minify JavaScript resources.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: resource.Resource
    signatures: ['js.Build [OPTIONS] RESOURCE']
---

The `js.Build` function uses the [evanw/esbuild] package to:

- Bundle
- Transpile (TypeScript and JSX)
- Tree shake
- Minify
- Create source maps

```go-html-template
{{ with resources.Get "js/main.js" }}
  {{$opts := dict
    "minify" (not hugo.IsDevelopment)
    "sourceMap" (cond hugo.IsDevelopment "external" "")
    "targetPath" "js/main.js"
  }}
  {{ with . | js.Build $opts }}
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

## Options

targetPath
: (`string`) If not set, the source path will be used as the base target path. Note that the target path's extension may change if the target MIME type is different, e.g. when the source is TypeScript.

format
: (`string`) The output format. One of: `iife`, `cjs`, `esm`. Default is `iife`, a self-executing function, suitable for inclusion as a `<script>` tag.

{{% include "/_common/functions/js/options.md" %}}

## Import JS code from the assets directory

`js.Build` has full support for the virtual union file system in [Hugo Modules](/hugo-modules/). You can see some simple examples in this [test project](https://github.com/gohugoio/hugoTestProjectJSModImports), but in short this means that you can do this:

```js
import { hello } from 'my/module';
```

And it will resolve to the top-most `index.{js,ts,tsx,jsx}` inside `assets/my/module` in the layered file system.

```js
import { hello3 } from 'my/module/hello3';
```

Will resolve to `hello3.{js,ts,tsx,jsx}` inside `assets/my/module`.

Any imports starting with `.` are resolved relative to the current file:

```js
import { hello4 } from './lib';
```

For other files (e.g. `JSON`, `CSS`) you need to use the relative path including any extension, e.g:

```js
import * as data from 'my/module/data.json';
```

Any imports in a file outside `assets` or that does not resolve to a component inside `assets` will be resolved by [ESBuild](https://esbuild.github.io/) with the **project directory** as the resolve directory (used as the starting point when looking for `node_modules` etc.). Also see [hugo mod npm pack](/commands/hugo_mod_npm_pack/). If you have any imported npm dependencies in your project, you need to make sure to run `npm install` before you run `hugo`.

Also note the new `params` option that can be passed from template to your JS files, e.g.:

```go-html-template
{{ $js := resources.Get "js/main.js" | js.Build (dict "params" (dict "api" "https://example.org/api")) }}
```

And then in your JS file:

```js
import * as params from '@params';
```

Hugo will, by default, generate a `assets/jsconfig.json` file that maps the imports. This is useful for navigation/intellisense help inside code editors, but if you don't need/want it, you can [turn it off](/configuration/build/).

## Node.js dependencies

Use the `js.Build` function to include Node.js dependencies.

Any imports in a file outside `assets` or that does not resolve to a component inside `assets` will be resolved by [esbuild](https://esbuild.github.io/) with the **project directory** as the resolve directory (used as the starting point when looking for `node_modules` etc.). Also see [hugo mod npm pack](/commands/hugo_mod_npm_pack/). If you have any imported npm dependencies in your project, you need to make sure to run `npm install` before you run `hugo`.

The start directory for resolving npm packages (aka. packages that live inside a `node_modules` directory) is always the main project directory.

> [!note]
> If you're developing a theme/component that is supposed to be imported and depends on dependencies inside `package.json`, we recommend reading about [hugo mod npm pack](/commands/hugo_mod_npm_pack/), a tool to consolidate all the npm dependencies in a project.

## Examples

```go-html-template
{{ $built := resources.Get "js/index.js" | js.Build "main.js" }}
```

Or with options:

```go-html-template
{{ $externals := slice "react" "react-dom" }}
{{ $defines := dict "process.env.NODE_ENV" `"development"` }}

{{ $opts := dict "targetPath" "main.js" "externals" $externals "defines" $defines }}
{{ $built := resources.Get "scripts/main.js" | js.Build $opts }}
<script src="{{ $built.RelPermalink }}" defer></script>
```

[evanw/esbuild]: https://github.com/evanw/esbuild
