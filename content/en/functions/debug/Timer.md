---
title: debug.Timer
description: Creates a named timer that reports elapsed time to the console.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: debug.Timer
  signatures: [debug.Timer NAME] 
relatedFunctions: []
---

{{< new-in "0.120.0" >}}

Use the `debug.Timer` function to determine execution time for a block of code, useful for finding performance bottlenecks in templates.

The timer starts when you instantiate it, and stops when you call its `Stop` method.

```go-html-template
{{ $t := debug.Timer "TestSqrt" }}
{{ range seq 2000 }}
  {{ $f := math.Sqrt . }}
{{ end }}
{{ $t.Stop }}
```

Use the `--logLevel info` command line flag when you build the site.

```sh
hugo --logLevel info
```

The results are displayed in the console at the end of the build. You can have as many timers as you want and if you don't stop them, they will be stopped at the end of build.

```text
INFO  timer:  name TestSqrt total 12.429355ms
```
