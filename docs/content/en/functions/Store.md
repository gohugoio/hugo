---
title: .Store
description: Returns a Scratch that is not reset on server rebuilds.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType:
  signatures: []
relatedFunctions:
  - .Store
  - .Scratch
---

The `.Store` method on `.Page` returns a [Scratch] to store and manipulate data. In contrast to the `.Scratch` method, this Scratch is not reset on server rebuilds.

[Scratch]: /functions/scratch/

### Methods

#### .Set

Sets the value of a given key.

```go-html-template
{{ .Store.Set "greeting" "Hello" }}
```

#### .Get

Gets the value of a given key.

```go-html-template
{{ .Store.Set "greeting" "Hello" }}

{{ .Store.Get "greeting" }} → Hello
```

#### .Add

Adds a given value to existing value(s) of the given key.

For single values, `Add` accepts values that support Go's `+` operator. If the first `Add` for a key is an array or slice, the following adds will be appended to that list.

```go-html-template
{{ .Store.Add "greetings" "Hello" }}
{{ .Store.Add "greetings" "Welcome" }}

{{ .Store.Get "greetings" }} → HelloWelcome
```

```go-html-template
{{ .Store.Add "total" 3 }}
{{ .Store.Add "total" 7 }}

{{ .Store.Get "total" }} → 10
```

```go-html-template
{{ .Store.Add "greetings" (slice "Hello") }}
{{ .Store.Add "greetings" (slice "Welcome" "Cheers") }}

{{ .Store.Get "greetings" }} → []interface {}{"Hello", "Welcome", "Cheers"}
```

#### .SetInMap

Takes a `key`, `mapKey` and `value` and adds a map of `mapKey` and `value` to the given `key`.

```go-html-template
{{ .Store.SetInMap "greetings" "english" "Hello" }}
{{ .Store.SetInMap "greetings" "french" "Bonjour" }}

{{ .Store.Get "greetings" }} → map[french:Bonjour english:Hello]
```

#### .DeleteInMap

Takes a `key` and `mapKey` and removes the map of `mapKey` from the given `key`.

```go-html-template
{{ .Store.SetInMap "greetings" "english" "Hello" }}
{{ .Store.SetInMap "greetings" "french" "Bonjour" }}
{{ .Store.DeleteInMap "greetings" "english" }}

{{ .Store.Get "greetings" }} → map[french:Bonjour]
```

#### .GetSortedMapValues

Returns an array of values from `key` sorted by `mapKey`.

```go-html-template
{{ .Store.SetInMap "greetings" "english" "Hello" }}
{{ .Store.SetInMap "greetings" "french" "Bonjour" }}

{{ .Store.GetSortedMapValues "greetings" }} → [Hello Bonjour]
```

#### .Delete

Removes the given key.

```go-html-template
{{ .Store.Set "greeting" "Hello" }}

{{ .Store.Delete "greeting" }}
```
