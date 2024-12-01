---
title: Store
linktitle: site.Store
description: Returns a persistent "scratch pad" on the given site to store and manipulate data.
categories: []
keywords: []
action:
  related:
  - methods/page/store
  - functions/hugo/store
  - functions/collections/NewScratch
  returnType: maps.Scratch
  signatures: [site.Store]
toc: true
---

{{< new-in 0.139.0 >}}

The `Store` method on a `Site` object creates a persistent [scratch pad] to store and manipulate data. To create a locally scoped scratch pad that is not attached to a `Site` object, use the [`newScratch`] function.

[`Scratch`]: /methods/site/scratch/
[`newScratch`]: /functions/collections/newscratch/
[scratch pad]: /getting-started/glossary/#scratch-pad

## Methods

###### Set

Sets the value of a given key.

```go-html-template
{{ site.Store.Set "greeting" "Hello" }}
```

###### Get

Gets the value of a given key.

```go-html-template
{{ site.Store.Set "greeting" "Hello" }}
{{ site.Store.Get "greeting" }} → Hello
```

###### Add

Adds a given value to existing value(s) of the given key.

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

###### SetInMap

Takes a `key`, `mapKey` and `value` and adds a map of `mapKey` and `value` to the given `key`.

```go-html-template
{{ site.Store.SetInMap "greetings" "english" "Hello" }}
{{ site.Store.SetInMap "greetings" "french" "Bonjour" }}
{{ site.Store.Get "greetings" }} → map[english:Hello french:Bonjour]
```

###### DeleteInMap

Takes a `key` and `mapKey` and removes the map of `mapKey` from the given `key`.

```go-html-template
{{ site.Store.SetInMap "greetings" "english" "Hello" }}
{{ site.Store.SetInMap "greetings" "french" "Bonjour" }}
{{ site.Store.DeleteInMap "greetings" "english" }}
{{ site.Store.Get "greetings" }} → map[french:Bonjour]
```

###### GetSortedMapValues

Returns an array of values from `key` sorted by `mapKey`.

```go-html-template
{{ site.Store.SetInMap "greetings" "english" "Hello" }}
{{ site.Store.SetInMap "greetings" "french" "Bonjour" }}
{{ site.Store.GetSortedMapValues "greetings" }} → [Hello Bonjour]
```

###### Delete

Removes the given key.

```go-html-template
{{ site.Store.Set "greeting" "Hello" }}
{{ site.Store.Delete "greeting" }}
```

## Determinate values

The `Store` method is often used to set scratch pad values within a shortcode, a partial template called by a shortcode, or by a Markdown render hook. In all three cases, the scratch pad values are indeterminate until Hugo renders the page content.

If you need to access a scratch pad value from a parent template, and the parent template has not yet rendered the page content, you can trigger content rendering by assigning the returned value to a [noop] variable:

[noop]: /getting-started/glossary/#noop

```go-html-template
{{ $noop := .Content }}
{{ site.Store.Get "mykey" }}
```

You can also trigger content rendering with the `ContentWithoutSummary`, `FuzzyWordCount`, `Len`, `Plain`, `PlainWords`, `ReadingTime`, `Summary`, `Truncated`, and `WordCount` methods. For example:

```go-html-template
{{ $noop := .WordCount }}
{{ site.Store.Get "mykey" }}
```
