---
_comment: Do not remove front matter.
---

###### params

(`map` or `slice`) Params that can be imported as JSON in your JS files, e.g.

```go-html-template
{{ $js := resources.Get "js/main.js" | js.Build (dict "params" (dict "api" "https://example.org/api")) }}
```
And then in your JS file:

```js
import * as params from '@params';
```

Note that this is meant for small data sets, e.g. configuration settings. For larger data, please put/mount the files into `assets` and import them directly.

###### minify

(`bool`) Let `js.Build` handle the minification.

###### loaders

{{< new-in 0.140.0 />}}

(`map`) Configuring a loader for a given file type lets you load that file type with an import statement or a require call. For example configuring the .png file extension to use the data URL loader means importing a .png file gives you a data URLcontaining the contents of that image. Loaders available are `none`, `base64`, `binary`, `copy`,  `css`,  `dataurl`, `default`, `empty`, `file`, `global-css`, `js`, `json`, `jsx`, `local-css`,  `text`, `ts`, `tsx`. See https://esbuild.github.io/api/#loader.

###### inject

(`slice`) This option allows you to automatically replace a global variable with an import from another file. The path names must be relative to `assets`. See https://esbuild.github.io/api/#inject.

###### shims

(`map`) This option allows swapping out a component with another. A common use case is to load dependencies like React from a CDN  (with _shims_) when in production, but running with the full bundled `node_modules` dependency during development:

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

###### target

(`string`) The language target. One of: `es5`, `es2015`, `es2016`, `es2017`, `es2018`, `es2019`, `es2020` or `esnext`. Default is `esnext`.

###### platform

{{< new-in 0.140.0 />}}

(`string`) One of `browser`, `node`, `neutral`.  Default is `browser`. See https://esbuild.github.io/api/#platform.

###### externals

(`slice`) External dependencies. Use this to trim dependencies you know will never be executed. See https://esbuild.github.io/api/#external.

###### defines

(`map`) Allow to define a set of string replacement to be performed when building. Should be a map where each key is to be replaced by its value.

```go-html-template
{{ $defines := dict "process.env.NODE_ENV" `"development"` }}
```

##### drop

Edit your source code before building to drop certain constructs: One of `debugger` or `console`.

{{< new-in 0.144.0 />}}

See https://esbuild.github.io/api/#drop

###### sourceMap

(`string`) Whether to generate `inline`, `linked` or `external` source maps from esbuild. Linked and external source maps will be written to the target with the output file name + ".map". When `linked` a `sourceMappingURL` will also be written to the output file. By default, source maps are not created. Note that the `linked` option was added in Hugo 0.140.0.

###### sourcesContent

{{< new-in 0.140.0 />}}

(`bool`) Whether to include the content of the source files in the source map. By default, this is `true`.

###### JSX

{{< new-in 0.124.0 />}}

(`string`) How to handle/transform JSX syntax. One of: `transform`, `preserve`, `automatic`. Default is `transform`. Notably, the `automatic` transform was introduced in React 17+ and will cause the necessary JSX helper functions to be imported automatically. See https://esbuild.github.io/api/#jsx.

###### JSXImportSource

{{< new-in 0.124.0 />}}

(`string`) Which library to use to automatically import its JSX helper functions from. This only works if `JSX` is set to `automatic`. The specified library needs to be installed through npm and expose certain exports. See https://esbuild.github.io/api/#jsx-import-source.

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
