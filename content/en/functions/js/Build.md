---
title: js.Build
description: Bundles, transpiles, tree shakes, and minifies JavaScript resources.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/js/Babel
    - functions/resources/Fingerprint
    - functions/resources/Minify
  returnType: resource.Resource
  signatures: ['js.Build [OPTIONS] RESOURCE']
toc: true
---

The `js.Build` function uses the [evanw/esbuild] package to:

- Bundle
- Transpile (TypeScript and JSX)
- Tree shake
- Minify
- Create source maps

[evanw/esbuild]: https://github.com/evanw/esbuild

```go-html-template
{{ with resources.Get "js/main.js" }}
  {{ if hugo.IsDevelopment }}
    {{ with . | js.Build }}
      <script src="{{ .RelPermalink }}"></script>
    {{ end }}
  {{ else }}
    {{ $opts := dict "minify" true }}
    {{ with . | js.Build $opts | fingerprint }}
      <script src="{{ .RelPermalink }}" integrity="{{ .Data.Integrity }}" crossorigin="anonymous"></script>
    {{ end }}
  {{ end }}
{{ end }}
```

## Options

targetPath
: (`string`) If not set, the source path will be used as the base target path.
Note that the target path's extension may change if the target MIME type is different, e.g. when the source is TypeScript.

params
: (`map` or `slice`) Params that can be imported as JSON in your JS files, e.g.

```go-html-template
{{ $js := resources.Get "js/main.js" | js.Build (dict "params" (dict "api" "https://example.org/api")) }}
```
And then in your JS file:

```js
import * as params from '@params';
```

Note that this is meant for small data sets, e.g. configuration settings. For larger data, please put/mount the files into `/assets` and import them directly.

minify
: (`bool`)Let `js.Build` handle the minification.

inject
: (`slice`) This option allows you to automatically replace a global variable with an import from another file. The path names must be relative to `assets`. See https://esbuild.github.io/api/#inject

shims
: (`map`) This option allows swapping out a component with another. A common use case is to load dependencies like React from a CDN  (with _shims_) when in production, but running with the full bundled `node_modules` dependency during development:

```go-html-template
{{ $shims := dict "react" "js/shims/react.js"  "react-dom" "js/shims/react-dom.js" }}
{{ $js = $js | js.Build dict "shims" $shims }}
```

The _shim_ files may look like these:

```js
// js/shims/react.js
module.exports = window.React;
```

```js
// js/shims/react-dom.js
module.exports = window.ReactDOM;
```

With the above, these imports should work in both scenarios:

```js
import * as React from 'react';
import * as ReactDOM from 'react-dom/client';
```

target
: (`string`) The language target. One of: `es5`, `es2015`, `es2016`, `es2017`, `es2018`, `es2019`, `es2020` or `esnext`. Default is `esnext`.

externals
: (`slice`) External dependencies. Use this to trim dependencies you know will never be executed. See https://esbuild.github.io/api/#external

defines
: (`map`) Allow to define a set of string replacement to be performed when building. Should be a map where each key is to be replaced by its value.

```go-html-template
{{ $defines := dict "process.env.NODE_ENV" `"development"` }}
```

format
: (`string`) The output format. One of: `iife`, `cjs`, `esm`. Default is `iife`, a self-executing function, suitable for inclusion as a `<script>` tag.

sourceMap
: (`string`) Whether to generate `inline` or `external` source maps from esbuild. External source maps will be written to the target with the output file name + ".map". Input source maps can be read from js.Build and node modules and combined into the output source maps. By default, source maps are not created.

JSX {{< new-in 0.124.0 >}}
: (`string`) How to handle/transform JSX syntax. One of: `transform`, `preserve`, `automatic`. Default is `transform`. Notably, the `automatic` transform was introduced in React 17+ and will cause the necessary JSX helper functions to be imported automatically. See https://esbuild.github.io/api/#jsx

JSXImportSource {{< new-in 0.124.0 >}}
: (`string`) Which library to use to automatically import its JSX helper functions from. This only works if `JSX` is set to `automatic`. The specified library needs to be installed through npm and expose certain exports. See https://esbuild.github.io/api/#jsx-import-source

The combination of `JSX` and `JSXImportSource` is helpful if you want to use a non-React JSX library like Preact, e.g.:

```go-html-template
{{ $js := resources.Get "js/main.jsx" | js.Build (dict "JSX" "automatic" "JSXImportSource" "preact") }}
```

With the above, you can use Preact components and JSX without having to manually import `h` and `Fragment` every time:

```jsx
import { render } from 'preact';

const App = () => <>Hello world!</>;

const container = document.getElementById('app');
if (container) render(<App />, container);
```

### Import JS code from /assets

`js.Build` has full support for the virtual union file system in [Hugo Modules](/hugo-modules/). You can see some simple examples in this [test project](https://github.com/gohugoio/hugoTestProjectJSModImports), but in short this means that you can do this:

```js
import { hello } from 'my/module';
```

And it will resolve to the top-most `index.{js,ts,tsx,jsx}` inside `assets/my/module` in the layered file system.

```js
import { hello3 } from 'my/module/hello3';
```

Will resolve to `hello3.{js,ts,tsx,jsx}` inside `assets/my/module`.

Any imports starting with `.` is resolved relative to the current file:

```js
import { hello4 } from './lib';
```

For other files (e.g. `JSON`, `CSS`) you need to use the relative path including any extension, e.g:

```js
import * as data from 'my/module/data.json';
```

Any imports in a file outside `/assets` or that does not resolve to a component inside `/assets` will be resolved by [ESBuild](https://esbuild.github.io/) with the **project directory** as the resolve directory (used as the starting point when looking for `node_modules` etc.). Also see [hugo mod npm pack](/commands/hugo_mod_npm_pack/). If you have any imported npm dependencies in your project, you need to make sure to run `npm install` before you run `hugo`.

Also note the new `params` option that can be passed from template to your JS files, e.g.:

```go-html-template
{{ $js := resources.Get "js/main.js" | js.Build (dict "params" (dict "api" "https://example.org/api")) }}
```
And then in your JS file:

```js
import * as params from '@params';
```

Hugo will, by default, generate a `assets/jsconfig.json` file that maps the imports. This is useful for navigation/intellisense help inside code editors, but if you don't need/want it, you can [turn it off](/getting-started/configuration/#configure-build).

## Node.js dependencies

Use the `js.Build` function to include Node.js dependencies.

Any imports in a file outside `/assets` or that does not resolve to a component inside `/assets` will be resolved by [esbuild](https://esbuild.github.io/) with the **project directory** as the resolve directory (used as the starting point when looking for `node_modules` etc.). Also see [hugo mod npm pack](/commands/hugo_mod_npm_pack/). If you have any imported npm dependencies in your project, you need to make sure to run `npm install` before you run `hugo`.

The start directory for resolving npm packages (aka. packages that live inside a `node_modules` folder) is always the main project folder.

{{% note %}}
If you're developing a theme/component that is supposed to be imported and depends on dependencies inside `package.json`, we recommend reading about [hugo mod npm pack](/commands/hugo_mod_npm_pack/), a tool to consolidate all the npm dependencies in a project.
{{% /note %}}

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
