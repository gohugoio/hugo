---
title: templates.Inner
description: Executes the content block enclosed by a partial call.
categories: []
keywords: [decorator]
params:
  functions_and_methods:
    aliases: [inner]
    returnType: any
    signatures: ['templates.Inner [CONTEXT]']
---

{{< new-in 0.154.0 />}}

The `templates.Inner` function defines the injection point for code nested within a block style partial call. This is the core mechanism used to create a [partial decorator][].

## Overview

The `templates.Inner` function acts as a placeholder within a partial template. When a partial is called as a decorator, it captures a block of code from the calling template rather than rendering it immediately. The `templates.Inner` function tells Hugo exactly where to inject that captured content.

This signals a reversal of execution where the callee becomes the caller. The partial manages the outer structure while the calling template remains in control of the inner content.

## Usage

To use this function, the calling template must use the block style syntax with a [`with`][] statement. This allows decorators to be deeply nested.

```go-html-template {file="layouts/home.html"}
{{ with partial "components/card.html" . }}
  <p>This content is passed to the partial.</p>
{{ end }}
```

Inside the partial, call `templates.Inner` to render the captured block.

```go-html-template {file="layouts/_partials/components/card.html"}
<div class="card-frame">
  {{ templates.Inner . }}
</div>
```

## Arguments

The function accepts one optional argument: the [context](g). This argument determines the value of the dot (`.`) inside the captured block when it is rendered.

- If you provide an argument, such as `{{ templates.Inner .SomeData }}`, the dot inside the captured block is rebound to that specific data.
- If you do not provide an argument, the captured block uses the context of the caller where the partial was first invoked.

## Context and scope

When using decorators, the `with` statement creates a new [scope](g). Variables defined outside the with block in the calling template are not automatically available inside the captured block.

By passing a context to `templates.Inner`, you ensure that the injected content has access to the correct data even when nested inside multiple layers of wrappers. This is critical when the decorator is used inside a loop or a specific data overlay.

## Repeated execution

A decorator can execute the captured content zero or more times. This is useful when the wrapper needs to repeat the same decoration for a collection of items, such as a list or a grid.

```go-html-template {file="layouts/_partials/list-decorator.html"}
<ul class="styled-list">
  {{ range .items }}
    <li>
      {{ templates.Inner . }}
    </li>
  {{ end }}
</ul>
```

In this example, the code provided by the caller is rendered once for every item in the .items collection, with the dot . updated to the current item in each iteration.

[`with`]: /functions/go-template/with/
[partial decorator]: /templates/partial-decorators/
