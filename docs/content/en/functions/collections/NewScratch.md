---
title: collections.NewScratch
description: Returns a locally scoped persistent data structure for storing and manipulating keyed values.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [newScratch]
    returnType: maps.Scratch
    signatures: [collections.NewScratch]
---

Use the `collections.NewScratch` function to create a locally scoped persistent data structure for storing and manipulating keyed values. To create a data structure with a different [scope](g), refer to the [scope](#scope) section below.

## Methods

Use these methods on the data structure.

`Set`
: Sets the value of the given key.

  ```go-html-template
  {{ $s := newScratch }}
  {{ $s.Set "greeting" "Hello" }}
  ```

`Get`
: (`any`) Gets the value of the given key.

  ```go-html-template
  {{ $s := newScratch }}
  {{ $s.Set "greeting" "Hello" }}
  {{ $s.Get "greeting" }} → Hello
  ```

`Add`
: Adds the given value to existing value(s) of the given key.

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

`SetInMap`
: Takes a `key`, `mapKey` and `value` and adds a map of `mapKey` and `value` to the given `key`.

  ```go-html-template
  {{ $s := newScratch }}
  {{ $s.SetInMap "greetings" "english" "Hello" }}
  {{ $s.SetInMap "greetings" "french" "Bonjour" }}
  {{ $s.Get "greetings" }} → map[english:Hello french:Bonjour]
  ```

`DeleteInMap`
: Takes a `key` and `mapKey` and removes the map of `mapKey` from the given `key`.

  ```go-html-template
  {{ $s := newScratch }}
  {{ $s.SetInMap "greetings" "english" "Hello" }}
  {{ $s.SetInMap "greetings" "french" "Bonjour" }}
  {{ $s.DeleteInMap "greetings" "english" }}
  {{ $s.Get "greetings" }} → map[french:Bonjour]
  ```

`GetSortedMapValues`
: (`[]any`) Returns an array of values from `key` sorted by `mapKey`.

  ```go-html-template
  {{ $s := newScratch }}
  {{ $s.SetInMap "greetings" "english" "Hello" }}
  {{ $s.SetInMap "greetings" "french" "Bonjour" }}
  {{ $s.GetSortedMapValues "greetings" }} → [Hello Bonjour]
  ```

`Delete`
: Removes the given key.

  ```go-html-template
  {{ $s := newScratch }}
  {{ $s.Set "greeting" "Hello" }}
  {{ $s.Delete "greeting" }}
  ```

`Values`
: (`map`) Returns the raw backing map. Do not use with `Store` methods on a `Page` object due to concurrency issues.

  ```go-html-template
  {{ $s := newScratch }}
  {{ $s.SetInMap "greetings" "english" "Hello" }}
  {{ $s.SetInMap "greetings" "french" "Bonjour" }}

  {{ $map := $s.Values }}
  ```

{{% include "_common/store-scope.md" %}}
