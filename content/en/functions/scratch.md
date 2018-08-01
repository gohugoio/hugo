---
title: .Scratch
description: Acts as a "scratchpad" to allow for writable page- or shortcode-scoped variables.
godocref:
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

In most cases you can do okay without `Scratch`, but due to scoping issues, there are many use cases that aren't solvable in Go Templates without `Scratch`'s help.

`.Scratch` is available as methods on `Page` and `Shortcode`. Since Hugo 0.43 you can also create a locally scoped `Scratch` using the template func `newScratch`.


{{% note %}}
See [this Go issue](https://github.com/golang/go/issues/10608) for the main motivation behind Scratch.
{{% /note %}}

{{% note %}}
For a detailed analysis of `.Scratch` and in context use cases, see this [post](https://regisphilibert.com/blog/2017/04/hugo-scratch-explained-variable/).
{{% /note %}}

## Get a Scratch

From Hugo `0.43` you can also create a locally scoped `Scratch` by calling `newScratch`:

```go-html-template
$scratch := newScratch
$scratch.Set "greeting" "Hello"
```

A `Scratch` is also added to both `Page` and `Shortcode`. `Scratch` has the following methods:

#### .Set

Set the given value to a given key

```go-html-template
{{ .Scratch.Set "greeting" "Hello" }}
```
#### .Get
Get the value of a given key

```go-html-template
{{ .Scratch.Set "greeting" "Hello" }}
----
{{ .Scratch.Get "greeting" }} > Hello
```

#### .Add
Will add a given value to existing value of the given key. 

For single values, `Add` accepts values that support Go's `+` operator. If the first `Add` for a key is an array or slice, the following adds will be appended to that list.

```go-html-template
{{ .Scratch.Add "greetings" "Hello" }}
{{ .Scratch.Add "greetings" "Welcome" }}
----
{{ .Scratch.Get "greetings" }} > HelloWelcome
```

```go-html-template
{{ .Scratch.Add "total" 3 }}
{{ .Scratch.Add "total" 7 }}
----
{{ .Scratch.Get "total" }} > 10
```


```go-html-template
{{ .Scratch.Add "greetings" (slice "Hello") }}
{{ .Scratch.Add "greetings" (slice "Welcome" "Cheers") }}
----
{{ .Scratch.Get "greetings" }} > []interface {}{"Hello", "Welcome", "Cheers"}
```

#### .SetInMap
Takes a `key`, `mapKey` and `value` and add a map of `mapKey` and `value` to the given `key`.

```go-html-template
{{ .Scratch.SetInMap "greetings" "english" "Hello" }}
{{ .Scratch.SetInMap "greetings" "french" "Bonjour" }}
----
{{ .Scratch.Get "greetings" }} > map[french:Bonjour english:Hello]
```

#### .GetSortedMapValues
Returns array of values from `key` sorted by `mapKey`

```go-html-template
{{ .Scratch.SetInMap "greetings" "english" "Hello" }}
{{ .Scratch.SetInMap "greetings" "french" "Bonjour" }}
----
{{ .Scratch.GetSortedMapValues "greetings" }} > [Hello Bonjour]
```
#### .Delete
Removes the given key

```go-html-template
{{ .Scratch.Delete "greetings" }}
```

## Scope
The scope of the backing data is global for the given `Page` or `Shortcode`, and spans partial and shortcode includes.

Note that `.Scratch` from a shortcode will return the shortcode's `Scratch`, which in most cases is what you want. If you want to store it in the page scoped Scratch, then use `.Page.Scratch`.




[pagevars]: /variables/page/
