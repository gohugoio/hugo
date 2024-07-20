---
title: collections.NewScratch
description: Returns a locally scoped "scratch pad" to store and manipulate data.
categories: []
keywords: []
action:
  aliases: [newScratch]
  related:
    - methods/page/scratch
    - methods/page/store
    - methods/shortcode/scratch
  returnType: maps.Scratch
  signatures: [collections.NewScratch ]
---

The `collections.NewScratch` function creates a locally scoped [scratch pad] to store and manipulate data. To create a scratch pad that is attached to a `Page` object, use the [`Scratch`] or [`Store`] method.

[`Scratch`]: /methods/page/scratch/
[`Store`]: /methods/page/store/
[scratch pad]: /getting-started/glossary/#scratch-pad

## Methods

###### Set

Sets the value of a given key.

```go-html-template
{{ $s := newScratch }}
{{ $s.Set "greeting" "Hello" }}
```

###### Get

Gets the value of a given key.

```go-html-template
{{ $s := newScratch }}
{{ $s.Set "greeting" "Hello" }}
{{ $s.Get "greeting" }} → Hello
```

###### Add

Adds a given value to existing value(s) of the given key.

For single values, `Add` accepts values that support Go's `+` operator. If the first `Add` for a key is an array or slice, the following adds will be appended to that list.

```go-html-template
{{ $s := newScratch }}
{{ $s.Set "greeting" "Hello" }}
{{ $s.Add "greeting" "Welcome" }}
{{ $s.Get "greeting" }} → HelloWelcome
```

```go-html-template
{{ $s := newScratch }}
{{ $s.Set "total" 3 }}
{{ $s.Add "total" 7 }}
{{ $s.Get "total" }} → 10
```

```go-html-template
{{ $s := newScratch }}
{{ $s.Set "greetings" (slice "Hello") }}
{{ $s.Add "greetings" (slice "Welcome" "Cheers") }}
{{ $s.Get "greetings" }} → [Hello Welcome Cheers]
```

###### SetInMap

Takes a `key`, `mapKey` and `value` and adds a map of `mapKey` and `value` to the given `key`.

```go-html-template
{{ $s := newScratch }}
{{ $s.SetInMap "greetings" "english" "Hello" }}
{{ $s.SetInMap "greetings" "french" "Bonjour" }}
{{ $s.Get "greetings" }} → map[english:Hello french:Bonjour]
```

###### DeleteInMap

Takes a `key` and `mapKey` and removes the map of `mapKey` from the given `key`.

```go-html-template
{{ $s := newScratch }}
{{ $s.SetInMap "greetings" "english" "Hello" }}
{{ $s.SetInMap "greetings" "french" "Bonjour" }}
{{ $s.DeleteInMap "greetings" "english" }}
{{ $s.Get "greetings" }} → map[french:Bonjour]
```

###### GetSortedMapValues

Returns an array of values from `key` sorted by `mapKey`.

```go-html-template
{{ $s := newScratch }}
{{ $s.SetInMap "greetings" "english" "Hello" }}
{{ $s.SetInMap "greetings" "french" "Bonjour" }}
{{ $s.GetSortedMapValues "greetings" }} → [Hello Bonjour]
```

###### Delete

Removes the given key.

```go-html-template
{{ $s := newScratch }}
{{ $s.Set "greeting" "Hello" }}
{{ $s.Delete "greeting" }}
```

###### Values

Returns the raw backing map. Do not use with `Scratch` or `Store` methods on a `Page` object due to concurrency issues.

```go-html-template
{{ $s := newScratch }}
{{ $s.SetInMap "greetings" "english" "Hello" }}
{{ $s.SetInMap "greetings" "french" "Bonjour" }}

{{ $map := $s.Values }}
```
