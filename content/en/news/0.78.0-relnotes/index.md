
---
date: 2020-11-03
title: "Hugo 0.78.0: Full Hugo Modules Support in js.Build"
description: "Resolve JavaScript imports top-down in the layered filesystem, pass parameters from template to JS, new JS intellisense helper, improved JS build errors."
categories: ["Releases"]
---

This release finally brings full [Hugo Modules](https://gohugo.io/hugo-modules/) support to [js.Build](https://gohugo.io/hugo-pipes/js/), curtsy of he new plugin API in the really, really fast [ESBuild](https://github.com/evanw/esbuild) by [@evanw](https://github.com/evanw).

Some notes on the improvements in this release:

* Now `js.Build` fully supports the virtual union filesystem in [Hugo Modules](https://gohugo.io/hugo-modules/). Any import inside your JavaScript components will resolve starting from the top component mount inside `/assets` with a fallback to the traditional "JS way" (`node_modules` etc.)
* You can now pass configuration data from the templates to your scripts via a new `params` option.
* Hugo now writes a `jsconfig.json` file inside `/assets` (you can turn it off) with import mappings to help editors such as VS Code with intellisense/navigation, which is especially useful when there is no common root and the source lives inside some temporary directory.
* We have also improved the build errors you get from `js.Build`. In server mode you will get a preview of the failing lines and in the console you will get a link to the location.

Read more about this in [the documentation](https://gohugo.io/hugo-pipes/js/), but a short usage example would look like:

In the template:

```go-html-template
{{ $js := resources.Get "js/main.js" | js.Build (dict "params" (dict "api" "https://example.org/api" ) }}
```

And then in a JavaScript component:

```js
import * as params from '@params';

// Will resolve to one of `hello.{js,ts,tsx,jsx}` inside `assets/my/module`.
import { hello } from 'my/module/hello';

var api = params.api;

hello();

```

## Changes

* Add avoidTDZ option [3b2fe3cd](https://github.com/gohugoio/hugo/commit/3b2fe3cd33b74166c3debec9826826f2b5a54fd9) [@bep](https://github.com/bep) [#7865](https://github.com/gohugoio/hugo/issues/7865)
* Make js.Build fully support modules [85e4dd73](https://github.com/gohugoio/hugo/commit/85e4dd7370eae97ae367e596aa6a10ba42fd4b7c) [@bep](https://github.com/bep) [#7816](https://github.com/gohugoio/hugo/issues/7816)[#7777](https://github.com/gohugoio/hugo/issues/7777)[#7916](https://github.com/gohugoio/hugo/issues/7916)
* Generate tsconfig files [3089fc0b](https://github.com/gohugoio/hugo/commit/3089fc0ba171be14670b19439bc2eab6b077b6c3) [@richtera](https://github.com/richtera) [#7777](https://github.com/gohugoio/hugo/issues/7777)






