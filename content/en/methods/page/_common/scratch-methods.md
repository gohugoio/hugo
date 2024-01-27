---
# Do not remove front matter.
---

## Methods

###### Set

Sets the value of a given key.

```go-html-template
{{ .Scratch.Set "greeting" "Hello" }}
```

###### Get

Gets the value of a given key.

```go-html-template
{{ .Scratch.Set "greeting" "Hello" }}
{{ .Scratch.Get "greeting" }} → Hello
```

###### Add

Adds a given value to existing value(s) of the given key.

For single values, `Add` accepts values that support Go's `+` operator. If the first `Add` for a key is an array or slice, the following adds will be appended to that list.

```go-html-template
{{ .Scratch.Set "greeting" "Hello" }}
{{ .Scratch.Add "greeting" "Welcome" }}
{{ .Scratch.Get "greeting" }} → HelloWelcome
```

```go-html-template
{{ .Scratch.Set "total" 3 }}
{{ .Scratch.Add "total" 7 }}
{{ .Scratch.Get "total" }} → 10
```

```go-html-template
{{ .Scratch.Set "greetings" (slice "Hello") }}
{{ .Scratch.Add "greetings" (slice "Welcome" "Cheers") }}
{{ .Scratch.Get "greetings" }} → [Hello Welcome Cheers]
```

###### SetInMap

Takes a `key`, `mapKey` and `value` and adds a map of `mapKey` and `value` to the given `key`.

```go-html-template
{{ .Scratch.SetInMap "greetings" "english" "Hello" }}
{{ .Scratch.SetInMap "greetings" "french" "Bonjour" }}
{{ .Scratch.Get "greetings" }} → map[english:Hello french:Bonjour]
```

###### DeleteInMap

Takes a `key` and `mapKey` and removes the map of `mapKey` from the given `key`.

```go-html-template
{{ .Scratch.SetInMap "greetings" "english" "Hello" }}
{{ .Scratch.SetInMap "greetings" "french" "Bonjour" }}
{{ .Scratch.DeleteInMap "greetings" "english" }}
{{ .Scratch.Get "greetings" }} → map[french:Bonjour]
```

###### GetSortedMapValues

Returns an array of values from `key` sorted by `mapKey`.

```go-html-template
{{ .Scratch.SetInMap "greetings" "english" "Hello" }}
{{ .Scratch.SetInMap "greetings" "french" "Bonjour" }}
{{ .Scratch.GetSortedMapValues "greetings" }} → [Hello Bonjour]
```

###### Delete

Removes the given key.

```go-html-template
{{ .Scratch.Set "greeting" "Hello" }}
{{ .Scratch.Delete "greeting" }}
```
