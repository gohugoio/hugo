---
title: plugin.Call
description: Call a function from a plugin.
godocref:
date: 2018-08-07
publishdate: 2018-08-07
lastmod: 2018-08-07
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [plugin,module]
signature: ["plugin.Call fn [arguments...]"]
workson: []
hugoversion:
relatedfuncs: ["plugin.Get", "plugin.Has"]
deprecated: false
aliases: []
---

First please [see how to manage plugins](../templates/plugins).

Once plugins are ready to use. All public methods can be called using a combination of [`plugin.Get`]({{< relref plugin.Get.md >}}) and [`plugin.Call`](.) functions.

```
{{ $fn := plugin.Get "my-plugin" "TheFunction" }} → Get TheFunction So you can use plugin.Call on it
{{ $fn | plugin.Call }} → Call TheFunction
{{ "lastArgument" | plugin.Call $fn "firstArgument" }} → True if plugins/my-plugin.so is valid
{{ list( "varArg1" "varArg2" ...) | plugin.Call $fn "firstArgument" }} → True if plugins/my-plugin.so is valid
```

## Limitations

Not all methods can be called: methods must have at most 2 results and in this case, the second one must an [error type](https://pkg.go.dev/errors).

Functions can have as many arguments as you want and can be variadic.
