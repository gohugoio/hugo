---
title: try
description: Returns a TryValue object after evaluating the given expression.
categories: []
keywords: []
action:
  aliases: []
  related: []
  returnType: TryValue
  signatures: ['try EXPRESSION']
toc: true
---

{{< new-in 0.141.0 />}}

The `try` statement is a non-standard extension to Go's [text/template] package. It introduces a mechanism for handling errors within templates, mimicking the `try-catch` constructs found in other programming languages.

[text/template]: https://pkg.go.dev/text/template

## Methods

The `TryValue` object encapsulates the result of evaluating the expression, and provides two methods:

Err
: (`string`) Returns a string representation of the error thrown by the expression, if an error occurred, or returns `nil` if the expression evaluated without errors.

Value
: (`any`) Returns the result of the expression if the evaluation was successful, or returns `nil` if an error occurred while evaluating the expression.

## Explanation

By way of example, let's divide a number by zero:

```go-html-template
{{ $x := 1 }}
{{ $y := 0 }}
{{ $result := div $x $y }}
{{ printf "%v divided by %v equals %v" $x $y .Value }}
```

As expected, the example above throws an error and fails the build:

```terminfo
Error: error calling div: can't divide the value by 0
```

Instead of failing the build, we can catch the error and emit a warning:

```go-html-template
{{ $x := 1 }}
{{ $y := 0 }}
{{ with try (div $x $y) }}
  {{ with .Err }}
    {{ warnf "%s" . }}
  {{ else }}
    {{ printf "%v divided by %v equals %v" $x $y .Value }}
  {{ end }}
{{ end }}
```

The error thrown by the expression is logged to the console as a warning:

```terminfo
WARN error calling div: can't divide the value by 0
```

Now let's change the arguments to avoid dividing by zero:

```go-html-template
{{ $x := 42 }}
{{ $y := 6 }}
{{ with try (div $x $y) }}
  {{ with .Err }}
    {{ warnf "%s" . }}
  {{ else }}
    {{ printf "%v divided by %v equals %v" $x $y .Value }}
  {{ end }}
{{ end }}
```

Hugo renders the above to:

```html
42 divided by 6 equals 7
```

## Example

Error handling is essential when using the [`resources.GetRemote`] function to capture remote resources such as data or images. When calling this function, if the HTTP request fails, Hugo will fail the build.

[`resources.GetRemote`]: /functions/resources/getremote/

Instead of failing the build, we can catch the error and emit a warning:

```go-html-template
{{ $url := "https://broken-example.org/images/a.jpg" }}
{{ with try (resources.GetRemote $url) }}
  {{ with .Err }}
    {{ warnf "%s" . }}
  {{ else with .Value }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ else }}
    {{ warnf "Unable to get remote resource %q" $url }}
  {{ end }}
{{ end }}
```
In the above, note that the [context](g) within the last conditional block is the `TryValue` object returned by the `try` statement. At this point neither the `Err` nor `Value` methods returned anything, so the current context is not useful. Use the `$` to access the [template context] if needed.

[template context]: /templates/introduction/#template-context

{{% note %}}
Hugo does not classify an HTTP response with status code 404 as an error. In this case `resources.GetRemote` returns nil.
{{% /note %}}
