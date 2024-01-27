---
title: Weekday
description:  Returns the day of the week of the given time.Time value.
categories: []
keywords: []
action:
  related:
    - functions/time/AsTime
  returnType: time.Weekday
  signatures: [TIME.Weekday]
---

To convert the `time.Weekday` value to a string:

```go-html-template
{{ $t := time.AsTime "2023-01-27T23:44:58-08:00" }}
{{ $t.Weekday.String }} → Friday
```

To convert the `time.Weekday` value to an integer.

```go-html-template
{{ $t := time.AsTime "2023-01-27T23:44:58-08:00" }}
{{ $t.Weekday | int }} → 5
