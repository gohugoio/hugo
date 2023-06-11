---
title: fn
description: Calls a custom, external Typescript/Jacascript function.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [function, external, code]
signature: ["external.Function FILE.FUNCTION ARGUMENTS…", "fn FILE.FUNCTION ARGUMENTS…"]
relatedfuncs: []
---

This function allows the use of custom, portable Typescript/Javascript code stored within the hugo site files.

For example if the file `functions/hello.ts` exists in a site with the exported function `Name`:

```typescript
export function Name(name: string): string {
  return `Hello ${name || "World"}!`
}
```

Then a Hugo template can call the external function `hello.Name`:

```go-html-template
<ul>
  <li>No arguments: {{ external.Function "hello.Name" }}</li>
  <li>Using alias: {{ fn "hello.Name" }}</li>
  <li>With arguments: {{ fn "hello.Name" "you" }}</li>
  <li>Piped arguments: {{ printf "from elsewhere" | fn "hello.Name" }}</li>
</ul>
```

Which will call the custom code and placing the results in the page:

```html
<ul>
  <li>No arguments: Hello World!</li>
  <li>Using alias: Hello World!</li>
  <li>With arguments: Hello you!</li>
  <li>Piped arguments: Hello from elsewhere!</li>
</ul>
```

## Naming

Capitalisation matters in the file/function name; calling `{{ fn "hello.name" }}` in the example above would fail with:

```plain
error calling fn: the function named name does not exist in hello
```

Similarly calling `{{ fn "heLLo.Name" }}` would fail with:

```plain
error calling fn: the function file named heLLo has not been loaded
```

## Function signatures

The exported functions can accept any number of arguments, of any type, but must only return a single string. any exceptions `throw`n will be handled gracefully.

Arguments are converted using [goja's `ToValue` method](https://pkg.go.dev/github.com/dop251/goja#Runtime.ToValue)). Note in particular that [dates do not work as expected](https://pkg.go.dev/github.com/dop251/goja#hdr-Handling_of_time_Time) in order to provide access to timezone information. If timezone is unimportant you can use `.UnixNano()`:

```typescript
export function WorkWithDates(myDateObj: object): string {
  const myDate = new Date(myDateObj.UnixNano()/1e6);

  return myDate.toString();
}
```

Which could be called with:

```go-template
{{ fb "yours.WorkWithDates" .Date }} → Sun Jun 11 2023 12:00:00 GMT+0100 (British Summer Time)
```

## Imports

Though the [Almond AMD loader](https://github.com/requirejs/almond) is readily available (via [clarkmcc's go-typescript](https://github.com/clarkmcc/go-typescript)) the Typescript execution environment does not have access to the filesystem.

(Perhaps encourage to use webpack or similar to build single file?)
