---
title: plugin.Open
description: Open a go plugin to Lookup symbols.
godocref:
date: 2018-08-07
publishdate: 2018-08-07
lastmod: 2018-08-07
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [plugin,module]
signature: ["plugin.Open pluginName"]
workson: []
hugoversion:
relatedfuncs: ["plugin.Exist", "plugin.Get", "plugin.Has"]
deprecated: false
aliases: []
---

First please [see how to manage plugins](../templates/plugins).

Once plugins are ready to use. You will be able to load it using [`plugin.Open`](.). This method returns the [Plugin go structure](https://golang.org/pkg/plugin/#Plugin).

It is recommended to use [`plugin.Get`]({{< relref "plugin.Get.md" >}}) instead because, at this moment, the only method available on the result is [`Lookup()`](https://golang.org/pkg/plugin/#Plugin.Lookup) which is equivalent to the [`plugin.Get` function]({{< relref "plugin.Get.md" >}}).

```
{{ plugin.Open "my-plugin" }} → Get the plugin from plugins/my-plugin.so
{{ (plugin.Open "my-plugin").Lookup "MyVariable" }} → Value of MyVariable from my-plugin
```
