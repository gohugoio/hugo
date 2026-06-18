---
title: Store
description: Returns a persistent data structure for storing and manipulating keyed values, scoped to the current site.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: maps.Scratch
    signatures: [site.Store]
---

{{< new-in 0.139.0 />}}

Use the `Store` method on a `Site` object to create a persistent data structure for storing and manipulating keyed values, scoped to the current site. To create a data structure with a different [scope](g), refer to the [scope](#scope) section below.

## Methods

Use these methods on the data structure.

`Set`
: Sets the value of a given key.

  ```go-html-template
  {{ site.Store.Set "greeting" "Hello" }}
  ```

`Get`
: (`any`) Gets the value of a given key.

  ```go-html-template
  {{ site.Store.Set "greeting" "Hello" }}
  {{ site.Store.Get "greeting" }} → Hello
  ```

`Add`
: Adds a given value to existing value(s) of the given key.

  For single values, `Add` accepts values that support Go's `+` operator. If the first `Add` for a key is an array or slice, the following adds will be appended to that list.

  ```go-html-template
  {{ site.Store.Set "greeting" "Hello" }}
  {{ site.Store.Add "greeting" "Welcome" }}
  {{ site.Store.Get "greeting" }} → HelloWelcome
  ```

  ```go-html-template
  {{ site.Store.Set "total" 3 }}
  {{ site.Store.Add "total" 7 }}
  {{ site.Store.Get "total" }} → 10
  ```

  ```go-html-template
  {{ site.Store.Set "greetings" (slice "Hello") }}
  {{ site.Store.Add "greetings" (slice "Welcome" "Cheers") }}
  {{ site.Store.Get "greetings" }} → [Hello Welcome Cheers]
  ```

`SetInMap`
: Takes a `key`, `mapKey` and `value` and adds a map of `mapKey` and `value` to the given `key`.

  ```go-html-template
  {{ site.Store.SetInMap "greetings" "english" "Hello" }}
  {{ site.Store.SetInMap "greetings" "french" "Bonjour" }}
  {{ site.Store.Get "greetings" }} → map[english:Hello french:Bonjour]
  ```

`DeleteInMap`
: Takes a `key` and `mapKey` and removes the map of `mapKey` from the given `key`.

  ```go-html-template
  {{ site.Store.SetInMap "greetings" "english" "Hello" }}
  {{ site.Store.SetInMap "greetings" "french" "Bonjour" }}
  {{ site.Store.DeleteInMap "greetings" "english" }}
  {{ site.Store.Get "greetings" }} → map[french:Bonjour]
  ```

`GetSortedMapValues`
: (`[]any`) Returns an array of values from `key` sorted by `mapKey`.

  ```go-html-template
  {{ site.Store.SetInMap "greetings" "english" "Hello" }}
  {{ site.Store.SetInMap "greetings" "french" "Bonjour" }}
  {{ site.Store.GetSortedMapValues "greetings" }} → [Hello Bonjour]
  ```

`Delete`
: Removes the given key.

  ```go-html-template
  {{ site.Store.Set "greeting" "Hello" }}
  {{ site.Store.Delete "greeting" }}
  ```

{{% include "_common/store-scope.md" %}}

## Determinate values

The `Store` method is often used to set values within a _shortcode_ template, a _partial_ template called by a _shortcode_ template, or by a _render hook_ template. In all three cases, the stored values are indeterminate until Hugo renders the page content.

If you need to access a stored value from a parent template, and the parent template has not yet rendered the page content, you can trigger content rendering by assigning the returned value to a [noop](g) variable:

```go-html-template
{{ $noop := .Content }}
{{ site.Store.Get "mykey" }}
```

You can also trigger content rendering with the `ContentWithoutSummary`, `FuzzyWordCount`, `Len`, `Plain`, `PlainWords`, `ReadingTime`, `Summary`, `Truncated`, and `WordCount` methods. For example:

```go-html-template
{{ $noop := .WordCount }}
{{ site.Store.Get "mykey" }}
```
