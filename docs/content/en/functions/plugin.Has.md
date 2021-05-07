---
title: plugin.Has
description: Check if a symbol exist in a plugin.
godocref:
date: 2018-08-07
publishdate: 2018-08-07
lastmod: 2018-08-07
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [plugin,module]
signature: ["plugin.Has pluginName symbolName"]
workson: []
hugoversion:
relatedfuncs: ["plugin.Exist", "plugin.Get"]
deprecated: false
aliases: []
---

First please [see how to manage plugins]({{< ref "templates/plugins" >}}).

Once plugins are ready to use. You can check wether a symbol exist in a plugin.

```
{{ plugin.Has "my-plugin" "BuildTimestamp" }} → True if the BuildTimestamp variable exist in my-plugin
{{ plugin.Has "my-plugin" "UnknownSymbol" }} -> False
{{ plugin.Has "my-plugin" "AnExistingMethod" }} -> True (symbols are variables and functions)
{{ if plugin.Exist "my-plugin" "Language" }}{{ plugin.Has "my-plugin" "Language" }}{{ end }} → True (Check wether my-plugin exist before checking for the symbol Language)
```
