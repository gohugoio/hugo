---
title: js.Batch
description: Build JavaScript bundle groups with global code splitting and flexible hooks/runners setup.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/js/Build
    - functions/js/Babel
    - functions/resources/Fingerprint
    - functions/resources/Minify
  returnType: js.Batcher
  signatures: ['js.Batch [ID]']
weight: 20
toc: true
---

{{% note %}}
For a runnable example of this feature, see [this test and demo repo](https://github.com/bep/hugojsbatchdemo/).
{{% /note %}}

The Batch `ID` is used to create the base directory for this batch. Forward slashes are allowed. `js.Batch` returns an object with an API with this structure:

* [Group]
  * [Script]
    * [SetOptions]
  * [Instance]
    * [SetOptions]
  * [Runner]
    * [SetOptions]
  * [Config]
    * [SetOptions]

## Group

The `Group` method take an `ID` (`string`) as argument. No slashes. It returns an object with these methods:

#### Script

The `Script` method takes an `ID` (`string`) as argument. No slashes. It returns an [OptionsSetter] that can be used to set [script options] for this script.

```go-html-template
{{ with js.Batch "js/mybatch" }}
  {{ with .Group "mygroup" }}
      {{ with .Script "myscript" }}
          {{ .SetOptions (dict "resource" (resources.Get "myscript.js")) }}
      {{ end }}
  {{ end }}
{{ end }}
```

`SetOptions` takes a [script options] map. Note that if you want the  script to be handled by a [runner], you need to set the `export` option to match what you want to pass on to the runner (default is `*`).

#### Instance

The `Instance` method takes two `string` arguments `SCRIPT_ID` and `INSTANCE_ID`. No slashes. It returns an [OptionsSetter] that can be used to set [params options] for this instance.

```go-html-template
{{ with js.Batch "js/mybatch" }}
  {{ with .Group "mygroup" }}
      {{ with .Instance "myscript" "myinstance" }}
          {{ .SetOptions (dict "params" (dict "param1" "value1")) }}
      {{ end }}
  {{ end }}
{{ end }}
```

`SetOptions` takes a [params options] map. The instance options will be passed to any [runner] script in the same group, as JSON.

#### Runner

The `Runner` method takes an `ID` (`string`) as argument. No slashes. It returns an [OptionsSetter] that can be used to set [script options] for this runner.

```go-html-template
{{ with js.Batch "js/mybatch" }}
  {{ with .Group "mygroup" }}
      {{ with .Runner "myrunner" }}
          {{ .SetOptions (dict "resource" (resources.Get "myrunner.js")) }}
      {{ end }}
  {{ end }}
{{ end }}
```

`SetOptions` takes a [script options] map.

The runner will receive a data structure with all instances for that group with a live binding of the [JavaScript import] of the defined `export`.

The runner script's export must be a function that takes one argument, the group data structure. An example of a group data structure as JSON is:

```json
{
    "id": "leaflet",
    "scripts": [
        {
            "id": "mapjsx",
            "binding": JAVASCRIPT_BINDING,
            "instances": [
                {
                    "id": "0",
                    "params": {
                        "c": "h-64",
                        "lat": 48.8533173846729,
                        "lon": 2.3497416090232535,
                        "r": "map.jsx",
                        "title": "Cathédrale Notre-Dame de Paris",
                        "zoom": 23
                    }
                },
                {
                    "id": "1",
                    "params": {
                        "c": "h-64",
                        "lat": 59.96300872062237,
                        "lon": 10.663529183196863,
                        "r": "map.jsx",
                        "title": "Holmenkollen",
                        "zoom": 3
                    }
                }
            ]
        }
    ]
}
```

Below is an example of a runner script that uses React to render elements. Note that the export (`default`) must match the `export` option in the [script options] (`default` is the default value for runner scripts) (runnable versions of examples on this page can be found at [js.Batch Demo Repo]):

```js
import * as ReactDOM from 'react-dom/client';
import * as React from 'react';

export default function Run(group) {
	console.log('Running react-create-elements.js', group);
	const scripts = group.scripts;
	for (const script of scripts) {
		for (const instance of script.instances) {
			/* This is a convention in this project. */
			let elId = `${script.id}-${instance.id}`;
			let el = document.getElementById(elId);
			if (!el) {
				console.warn(`Element with id ${elId} not found`);
				continue;
			}
			const root = ReactDOM.createRoot(el);
			const reactEl = React.createElement(script.binding, instance.params);
			root.render(reactEl);
		}
	}
}
```

#### Config

Returns an [OptionsSetter] that can be used to set [build options] for the batch.

These are mostly the same as for `js.Build`, but note that:

* `targetPath` is set automatically (there may be multiple outputs).
* `format` must be `esm`, currently the only format supporting [code splitting].
* `params` will be available in the `@params/config` namespace in the scripts. This way you can import both the [script] or [runner] params and the [config] params with:

```js
import * as params from "@params";
import * as config from "@params/config";
```

Setting the `Config` for a batch can be done from any template (including shortcode templates), but will only be set once (the first will win):

```go-html-template
{{ with js.Batch "js/mybatch" }}
  {{ with .Config }}
       {{ .SetOptions (dict
        "target" "es2023"
        "format" "esm"
        "jsx" "automatic"
        "loaders" (dict ".png" "dataurl")
        "minify" true
        "params" (dict "param1" "value1")
        )
      }}
  {{ end }}
{{ end }}
```

## Options

### Build Options

format
: (`string`) Currently only `esm` is supported in [ESBuild's code splitting].

{{% include "./_common/options.md" %}}

### Script Options

resource
: The resource to build. This can be a file resource or a virtual resource.

export
: The export to bind the runner to. Set it to `*` to export the [entire namespace](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/import#namespace_import). Default is `default` for [runner] scripts and `*` for other [scripts](#script).

importContext
: An additional context for resolving imports. Hugo will always check this one first before falling back to `assets` and `node_modules`. A common use of this is to resolve imports inside a page bundle. See [import context](#import-context).

params
: A map of parameters that will be passed to the script as JSON. These gets bound to the `@params` namespace:
```js
import * as params from '@params';
```

### Script Options

### Params Options

params
: A map of parameters that will be passed to the script as JSON. 

### Import Context

Hugo will, by default, first try to resolve any import in [assets](/hugo-pipes/introduction/#asset-directory) and, if not found, let [ESBuild] resolve it (e.g. from `node_modules`). The  `importContext` option can be used to set the first context for resolving imports. A common use of this is to resolve imports inside a [page bundle](/content-management/page-bundles/).

```go-html-template
{{ $common := resources.Match "/js/headlessui/*.*" }}
{{ $importContext := (slice $.Page ($common.Mount "/js/headlessui" ".")) }}
```

You can pass any object that implements [Resource.Get](/methods/page/resources/#get). Pass a slice to set multiple contexts.

The example above uses [`Resources.Mount`] to resolve a directory inside `assets` relative to the page bundle.

### OptionsSetter

An `OptionsSetter` is a special object that is returned once only. This means that you should wrap it with [with]:

```go-html-template
{{ with .Script "myscript" }}
    {{ .SetOptions (dict "resource" (resources.Get "myscript.js"))}}
{{ end }}
```

## Build

The `Build` method returns an object with the following structure:

* Groups (map)
  * [`Resources`]

Eeach [`Resource`] will be of media type `application/javascript` or `text/css`.

 In a template you would typically handle one group with a given `ID` (e.g. scripts for the current section). Because of the concurrent build, this needs to be done in a [`templates.Defer`] block:

{{% note %}}
The [`templates.Defer`] acts as a synchronisation point to handle scripts added concurrently by different templates. If you have a setup with where the batch is created in one go (in one template), you don't need it.

See [this discussion](https://discourse.gohugo.io/t/js-batch-with-simple-global-script/53002/5?u=bep) for more.

[`templates.Defer`]: /functions/templates/defer/
{{% /note %}}

```go-html-template
{{ $group := .group }}
{{ with (templates.Defer (dict "key" $group "data" $group )) }}
  {{ with (js.Batch "js/mybatch") }}
    {{ with .Build }}
      {{ with index .Groups $ }}
        {{ range . }}
          {{ $s := . }}
          {{ if eq $s.MediaType.SubType "css" }}
            <link href="{{ $s.RelPermalink }}" rel="stylesheet" />
          {{ else }}
            <script src="{{ $s.RelPermalink }}" type="module"></script>
          {{ end }}
        {{ end }}
      {{ end }}
  {{ end }}
{{ end }}
```

## Known Issues

In the official documentation for [ESBuild's code splitting], there's a warning note in the header. The two issues are:

* `esm` is currently the only implemented output format. This means that it will not work for very old browsers. See [caniuse](https://caniuse.com/?search=ESM).
* There's a known import ordering issue.

We have not seen the ordering issue as a problem during our [extensive testing](https://github.com/bep/hugojsbatchdemo) of this new feature with different libraries. There are two main cases:

1. Undefined execution order of imports, see [this comment](https://github.com/evanw/esbuild/issues/399#issuecomment-1458680887)
1. Only one execution order of imports, see [this comment](https://github.com/evanw/esbuild/issues/399#issuecomment-735355932)

Many would say that both of the above are [code smells](https://en.wikipedia.org/wiki/Code_smell). The first one has a simple workaround in Hugo. Define the import order in its own script and make sure it gets passed early to ESBuild, e.g. by putting it in a script group with a name that comes early in the alphabet.

```js
import './lib2.js';
import './lib1.js';

console.log('entrypoints-workaround.js');

```

[build options]: #build-options
[`Resource`]: /methods/resource/
[`Resources`]: /methods/page/resources/
[`Resources.Mount`]: /methods/page/resources/#mount
[`templates.Defer`]: /functions/templates/defer/
[code splitting]: https://esbuild.github.io/api/#splitting
[config]: #config
[ESBuild's code splitting]: https://esbuild.github.io/api/#splitting
[ESBuild]: https://github.com/evanw/esbuild
[group]: #group
[instance]: #instance
[JavaScript import]: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/import
[js.Batch Demo Repo]: https://github.com/bep/hugojsbatchdemo/
[map]: /functions/collections/dictionary/
[OptionsSetter]: #optionssetter
[page bundles]: /content-management/page-bundles/
[params options]: #params-options
[runner]: #runner
[script options]: #script-options
[script]: #script
[SetOptions]: #optionssetter
[with]: /functions/go-template/with/
