---
title: plugin.Exist
description: Checks that a plugin exist.
godocref:
date: 2018-08-07
publishdate: 2018-08-07
lastmod: 2018-08-07
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [plugin,module]
signature: ["plugin.Exist pluginName"]
workson: []
hugoversion:
relatedfuncs: ["plugin.Get", "plugin.Has"]
deprecated: false
aliases: []
---

Check wether a plugin is [installed and ready to use]({{< ref "templates/plugins" >}}).

```
{{ plugin.Exist "my-plugin" }} → True if plugins/my-plugin.so is valid
{{ plugin.Exist "unknown-plugin" }} → False
{{ plugin.Exist "plugins/my-plugin.so" }} → False (expected "my-plugin", not the full path)
```
