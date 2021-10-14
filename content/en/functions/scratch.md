---
title: .Scratch
description: Acts as a "scratchpad" to store and manipulate data.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
keywords: [iteration]
categories: [functions]
menu:
  docs:
    parent: "functions"
toc:
signature: []
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: [/extras/scratch/,/doc/scratch/]
---

Scratch is a Hugo feature designed to conveniently manipulate data in a Go Template world. It is either a Page or Shortcode method for which the resulting data will be attached to the given context, or it can live as a unique instance stored in a variable.

{{% note %}}
Note that Scratch was initially created as a workaround for a [Go template scoping limitation](https://github.com/golang/go/issues/10608) that affected Hugo versions prior to 0.48. For a detailed analysis of `.Scratch` and contextual use cases, see [this blog post](https://regisphilibert.com/blog/2017/04/hugo-scratch-explained-variable/).
{{% /note %}}

### Contexted `.Scratch` vs. local `newScratch`

Since Hugo 0.43, there are two different ways of using Scratch:

#### The Page's `.Scratch`

`.Scratch` is available as a Page method or a Shortcode method and attaches the "scratched" data to the given page. Either a Page or a Shortcode context is required to use `.Scratch`.

```go-html-template
{{ .Scratch.Set "greeting" "bonjour" }}
{{ range .Pages }}
  {{ .Scratch.Set "greeting" (print "bonjour" .Title) }}
{{ end }}
```

#### The local `newScratch`

{{< new-in "0.43" >}} A Scratch instance can also be assigned to any variable using the `newScratch` function. In this case, no Page or Shortcode context is required and the scope of the scratch is only local. The methods detailed below are available from the variable the Scratch instance was assigned to.

```go-html-template
{{ $data := newScratch }}
{{ $data.Set "greeting" "hola" }}
```

### Methods

A Scratch has the following methods:

{{% note %}}
Note that the following examples assume a [local Scratch instance](#the-local-newscratch) has been stored in `$scratch`.
{{% /note %}}

#### .Set

Set the value of a given key.

```go-html-template
{{ $scratch.Set "greeting" "Hello" }}
```

#### .Get

Get the value of a given key.

```go-html-template
{{ $scratch.Set "greeting" "Hello" }}
----
{{ $scratch.Get "greeting" }} > Hello
```

#### .Add

Add a given value to existing value(s) of the given key. 

For single values, `Add` accepts values that support Go's `+` operator. If the first `Add` for a key is an array or slice, the following adds will be appended to that list.

```go-html-template
{{ $scratch.Add "greetings" "Hello" }}
{{ $scratch.Add "greetings" "Welcome" }}
----
{{ $scratch.Get "greetings" }} > HelloWelcome
```

```go-html-template
{{ $scratch.Add "total" 3 }}
{{ $scratch.Add "total" 7 }}
----
{{ $scratch.Get "total" }} > 10
```

```go-html-template
{{ $scratch.Add "greetings" (slice "Hello") }}
{{ $scratch.Add "greetings" (slice "Welcome" "Cheers") }}
----
{{ $scratch.Get "greetings" }} > []interface {}{"Hello", "Welcome", "Cheers"}
```

#### .SetInMap

Takes a `key`, `mapKey` and `value` and adds a map of `mapKey` and `value` to the given `key`.

```go-html-template
{{ $scratch.SetInMap "greetings" "english" "Hello" }}
{{ $scratch.SetInMap "greetings" "french" "Bonjour" }}
----
{{ $scratch.Get "greetings" }} > map[french:Bonjour english:Hello]
```

#### .DeleteInMap
Takes a `key` and `mapKey` and removes the map of `mapKey` from the given `key`.

```go-html-template
{{ .Scratch.SetInMap "greetings" "english" "Hello" }}
{{ .Scratch.SetInMap "greetings" "french" "Bonjour" }}
----
{{ .Scratch.DeleteInMap "greetings" "english" }}
----
{{ .Scratch.Get "greetings" }} > map[french:Bonjour]
```

#### .GetSortedMapValues

Return an array of values from `key` sorted by `mapKey`.

```go-html-template
{{ $scratch.SetInMap "greetings" "english" "Hello" }}
{{ $scratch.SetInMap "greetings" "french" "Bonjour" }}
----
{{ $scratch.GetSortedMapValues "greetings" }} > [Hello Bonjour]
```

#### .Delete

{{< new-in "0.38" >}} Remove the given key.

```go-html-template
{{ $scratch.Set "greeting" "Hello" }}
----
{{ $scratch.Delete "greeting" }}
```

#### .Values

Return the raw backing map. Note that you should only use this method on the locally scoped Scratch instances you obtain via [`newScratch`](#the-local-newscratch), not `.Page.Scratch` etc., as that will lead to concurrency issues.


[pagevars]: /variables/page/
