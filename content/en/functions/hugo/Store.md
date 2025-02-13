---
title: hugo.Store
description: Returns a globally scoped "scratch pad" to store and manipulate data.
categories: []
keywords: []
action:
  related:
    - methods/page/Store
    - methods/site/Store
    - methods/shortcode/Store
    - functions/collections/NewScratch
  returnType: maps.Scratch
  signatures: [hugo.Store]
toc: true
---

{{< new-in 0.139.0 />}}

Use the `hugo.Store` function to create a globally scoped [scratch pad](g) to store and manipulate data. To create a scratch pad with a different [scope](g), refer to the [scope](#scope) section below.

## Methods

###### Set

Sets the value of the given key.

```go-html-template
{{ hugo.Store.Set "greeting" "Hello" }}
```

###### Get

Gets the value of the given key.

```go-html-template
{{ hugo.Store.Set "greeting" "Hello" }}
{{ hugo.Store.Get "greeting" }} → Hello
```

###### Add

Adds the given value to the existing value(s) of the given key.

For single values, `Add` accepts values that support Go's `+` operator. If the first `Add` for a key is an array or slice, the following adds will be appended to that list.

```go-html-template
{{ hugo.Store.Set "greeting" "Hello" }}
{{ hugo.Store.Add "greeting" "Welcome" }}
{{ hugo.Store.Get "greeting" }} → HelloWelcome
```

```go-html-template
{{ hugo.Store.Set "total" 3 }}
{{ hugo.Store.Add "total" 7 }}
{{ hugo.Store.Get "total" }} → 10
```

```go-html-template
{{ hugo.Store.Set "greetings" (slice "Hello") }}
{{ hugo.Store.Add "greetings" (slice "Welcome" "Cheers") }}
{{ hugo.Store.Get "greetings" }} → [Hello Welcome Cheers]
```

###### SetInMap

Takes a `key`, `mapKey` and `value` and adds a map of `mapKey` and `value` to the given `key`.

```go-html-template
{{ hugo.Store.SetInMap "greetings" "english" "Hello" }}
{{ hugo.Store.SetInMap "greetings" "french" "Bonjour" }}
{{ hugo.Store.Get "greetings" }} → map[english:Hello french:Bonjour]
```

###### DeleteInMap

Takes a `key` and `mapKey` and removes the map of `mapKey` from the given `key`.

```go-html-template
{{ hugo.Store.SetInMap "greetings" "english" "Hello" }}
{{ hugo.Store.SetInMap "greetings" "french" "Bonjour" }}
{{ hugo.Store.DeleteInMap "greetings" "english" }}
{{ hugo.Store.Get "greetings" }} → map[french:Bonjour]
```

###### GetSortedMapValues

Returns an array of values from `key` sorted by `mapKey`.

```go-html-template
{{ hugo.Store.SetInMap "greetings" "english" "Hello" }}
{{ hugo.Store.SetInMap "greetings" "french" "Bonjour" }}
{{ hugo.Store.GetSortedMapValues "greetings" }} → [Hello Bonjour]
```

###### Delete

Removes the given key.

```go-html-template
{{ hugo.Store.Set "greeting" "Hello" }}
{{ hugo.Store.Delete "greeting" }}
```

{{% include "_common/scratch-pad-scope.md" %}}

## Determinate values

The `Store` method is often used to set scratch pad values within a shortcode, a partial template called by a shortcode, or by a Markdown render hook. In all three cases, the scratch pad values are indeterminate until Hugo renders the page content.

If you need to access a scratch pad value from a parent template, and the parent template has not yet rendered the page content, you can trigger content rendering by assigning the returned value to a [noop](g) variable:

```go-html-template
{{ $noop := .Content }}
{{ hugo.Store.Get "mykey" }}
```

You can also trigger content rendering with the `ContentWithoutSummary`, `FuzzyWordCount`, `Len`, `Plain`, `PlainWords`, `ReadingTime`, `Summary`, `Truncated`, and `WordCount` methods. For example:

```go-html-template
{{ $noop := .WordCount }}
{{ hugo.Store.Get "mykey" }}
```
