---
title: transform
description: Functions for transforming data.
categories: [
  "templates",
  "data-and-files"
]
menu:
  docs:
    parent: "templates"
weight: 120
toc: true
aliases: []
---

## transform.OpenAPIDocFromJSON

Loads an OpenAPI (Swagger) document from a file path or URL.

This function supports version 2.0 and 3.0 of the specification.

```go-html-template
{{ $doc := transform.OpenAPIDocFromJSON "openapi.json" }}
