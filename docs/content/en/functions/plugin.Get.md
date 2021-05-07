---
title: plugin.Get
description: Get a symbol from a plugin.
godocref:
date: 2018-08-07
publishdate: 2018-08-07
lastmod: 2018-08-07
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [plugin,module]
signature: ["plugin.Get pluginName symbolName"]
workson: []
hugoversion:
relatedfuncs: ["plugin.Exist", "plugin.Has"]
deprecated: false
aliases: []
---

First please [see how to manage plugins]({{< ref "templates/plugins" >}}).

Once plugins are ready to use. All public symbols can be retrieved.

```
{{ plugin.Get "my-plugin" "BuildTimestamp" }} → Get the BuildTimestamp variable in ./plugins/my-plugin.so
{{ if plugin.Has "my-plugin" "Language" }}{{ plugin.Get "my-plugin" "Language" }}{{ end }} → Get the Language variable if exist
{{ $fn := plugin.Get "my-plugin" "TheFunction" | plugin.Call }} → Call TheFunction() from ./plugins/my-plugin.so
```
